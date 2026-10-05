// Package motion plays audited, independently timed PNG sequences alongside the
// legacy sprite sheets. It is shared by the Win32 and Cocoa front ends.
package motion

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	xdraw "golang.org/x/image/draw"
)

type Presentation struct {
	Scale float64 `json:"scale"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}
type Clip struct {
	Durations    []float64           `json:"durations"`
	Tiers        map[string][]string `json:"tiers"`
	Presentation Presentation        `json:"presentation"`
}
type Profile struct {
	// Stride is distance in canvas-height units travelled per full walk cycle.
	Stride float64 `json:"stride"`
	Ground float64 `json:"ground"`
	Facing int     `json:"facing"`
}
type Variant struct {
	ID      string                    `json:"id"`
	Enabled bool                      `json:"enabled"`
	Profile Profile                   `json:"profile"`
	Motions map[string]Clip           `json:"motions"`
	Actions map[string]ActionSequence `json:"actions,omitempty"`
}
type Manifest struct {
	Schema   int               `json:"schema"`
	Variants []Variant         `json:"variants"`
	Hashes   map[string]string `json:"hashes"`
}
type Loaded struct {
	Variant
	Frames           map[string]map[int][]*image.RGBA
	renderMu         sync.Mutex
	rendered         map[renderKey]*image.RGBA
	renderOrder      []renderKey
	sideMargin       float64
	automaticActions []string
}
type renderKey struct {
	action               string
	frame, width, height int
	flip                 bool
}
type Store struct {
	mu       sync.Mutex
	fs       fs.FS
	manifest Manifest
	loaded   map[string]*Loaded
	failed   map[string]error
}

func NewStore(files fs.FS, manifestPath string) (*Store, error) {
	data, err := fs.ReadFile(files, manifestPath)
	if err != nil {
		return nil, err
	}
	s := &Store{fs: files, loaded: map[string]*Loaded{}, failed: map[string]error{}}
	if err = json.Unmarshal(data, &s.manifest); err != nil {
		return nil, err
	}
	if s.manifest.Schema != 1 {
		return nil, fmt.Errorf("unsupported motion schema %d", s.manifest.Schema)
	}
	seen := map[string]bool{}
	for _, v := range s.manifest.Variants {
		if v.ID == "" || seen[v.ID] {
			return nil, fmt.Errorf("duplicate or empty motion id %q", v.ID)
		}
		seen[v.ID] = true
	}
	return s, nil
}

// Get loads one complete variant atomically. A bad clip disables only this
// variant; callers keep its own legacy art for the entire session.
func (s *Store) Get(id string) *Loaded {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if v := s.loaded[id]; v != nil {
		return v
	}
	if s.failed[id] != nil {
		return nil
	}
	for _, v := range s.manifest.Variants {
		if v.ID == id && v.Enabled {
			loaded, err := s.load(v)
			if err != nil {
				s.failed[id] = err
				return nil
			}
			s.loaded[id] = loaded
			return loaded
		}
	}
	return nil
}
func (s *Store) Error(id string) error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.failed[id]
}

// Retain releases decoded frames after selection changes; pets of the same ID
// share the same immutable images. Negative cache entries do not hold images.
func (s *Store) Retain(ids map[string]bool) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for id := range s.loaded {
		if !ids[id] {
			delete(s.loaded, id)
		}
	}
}
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func (s *Store) load(v Variant) (*Loaded, error) {
	if !finite(v.Profile.Stride) || v.Profile.Stride <= 0 || v.Profile.Stride > 5 || !finite(v.Profile.Ground) || v.Profile.Ground <= 0 || v.Profile.Ground > 1.2 || (v.Profile.Facing != 1 && v.Profile.Facing != -1) {
		return nil, fmt.Errorf("%s: invalid gait profile", v.ID)
	}
	l := &Loaded{Variant: v, Frames: map[string]map[int][]*image.RGBA{}}
	clips, err := sequenceClips(v)
	if err != nil {
		return nil, err
	}
	for _, action := range clips {
		c, ok := v.Motions[action]
		if !ok || len(c.Durations) == 0 || len(c.Durations) > 1024 || len(c.Tiers) == 0 {
			return nil, fmt.Errorf("%s: missing %s", v.ID, action)
		}
		for _, d := range c.Durations {
			if !finite(d) || d <= 0 || d > 60000 {
				return nil, fmt.Errorf("%s: invalid duration", v.ID)
			}
		}
		p := c.Presentation
		if !finite(p.Scale) || p.Scale <= 0 || p.Scale > 3 || !finite(p.X) || !finite(p.Y) || math.Abs(p.X) > 2 || math.Abs(p.Y) > 2 {
			return nil, fmt.Errorf("%s: invalid presentation", v.ID)
		}
		l.Frames[action] = map[int][]*image.RGBA{}
		for tier, paths := range c.Tiers {
			h, err := strconv.Atoi(tier)
			if err != nil || h < 16 || h > 512 || len(paths) != len(c.Durations) {
				return nil, fmt.Errorf("%s: invalid tier", v.ID)
			}
			for _, path := range paths {
				if !fs.ValidPath(path) {
					return nil, fmt.Errorf("invalid frame path %q", path)
				}
				data, err := fs.ReadFile(s.fs, path)
				if err != nil {
					return nil, err
				}
				hash := sha256.Sum256(data)
				if hex.EncodeToString(hash[:]) != s.manifest.Hashes[path] {
					return nil, fmt.Errorf("hash mismatch: %s", path)
				}
				cfg, err := png.DecodeConfig(bytes.NewReader(data))
				if err != nil || cfg.Height != h || cfg.Width != h*3/2 {
					return nil, fmt.Errorf("invalid frame canvas: %s", path)
				}
				img, err := png.Decode(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				rgba := image.NewRGBA(img.Bounds())
				draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
				l.Frames[action][h] = append(l.Frames[action][h], rgba)
				left, right := rgba.Bounds().Dx(), 0
				for y := 0; y < rgba.Bounds().Dy(); y++ {
					for x := 0; x < rgba.Bounds().Dx(); x++ {
						if rgba.Pix[y*rgba.Stride+x*4+3] > 0 {
							left = min(left, x)
							right = max(right, x+1)
						}
					}
				}
				if right > left {
					l.sideMargin = math.Max(l.sideMargin, math.Max(-(float64(left)/float64(h)*p.Scale+p.X), float64(right)/float64(h)*p.Scale+p.X-1.5))
				}
			}
		}
	}
	for name, sequence := range v.Actions {
		if sequence.Automatic {
			l.automaticActions = append(l.automaticActions, name)
		}
	}
	sort.Strings(l.automaticActions)
	return l, nil
}

// A whole-motion transform can extend a tail past the nominal canvas. Reserve
// the union of both facings at either edge, without modifying the source art.
func (l *Loaded) HorizontalPadding(height int) int {
	return int(math.Ceil(l.sideMargin*float64(height))) + 2
}
func Duration(d []float64) float64 {
	var total float64
	for _, v := range d {
		total += v
	}
	return total
}
func FrameAt(phase float64, d []float64) int {
	total := Duration(d)
	if total <= 0 || !finite(phase) {
		return 0
	}
	phase = math.Mod(phase, total)
	if phase < 0 {
		phase += total
	}
	for i, v := range d {
		if phase < v {
			return i
		}
		phase -= v
	}
	return 0
}
func (l *Loaded) Frame(action string, phase float64, height int) *image.RGBA {
	tiers := l.Frames[action]
	sizes := make([]int, 0, len(tiers))
	for h := range tiers {
		sizes = append(sizes, h)
	}
	sort.Ints(sizes)
	if len(sizes) == 0 {
		return nil
	}
	tier := sizes[len(sizes)-1]
	for _, h := range sizes {
		if h >= height {
			tier = h
			break
		}
	}
	return tiers[tier][FrameAt(phase, l.Motions[action].Durations)]
}

// Draw preserves the upstream whole-motion transform, including its mirrored
// translation. Ground is measured for the animal, not copied from the degu.
func (l *Loaded) Draw(dst *image.RGBA, action string, phase float64, x, ground, width, height, dir int) {
	if action == "" {
		action = "idle"
	}
	src := l.Frame(action, phase, height)
	if src == nil {
		return
	}
	p := l.Motions[action].Presentation
	w, h := float64(width)*p.Scale, float64(height)*p.Scale
	xf := float64(x) + p.X*float64(height)
	if dir*l.Profile.Facing < 0 {
		xf = float64(x+width) - p.X*float64(height) - w
	}
	yf := float64(ground) - float64(height)*l.Profile.Ground + p.Y*float64(height)
	r := image.Rect(int(math.Round(xf)), int(math.Round(yf)), int(math.Round(xf+w)), int(math.Round(yf+h)))
	key := renderKey{action, FrameAt(phase, l.Motions[action].Durations), r.Dx(), r.Dy(), dir*l.Profile.Facing < 0}
	l.renderMu.Lock()
	defer l.renderMu.Unlock()
	if l.rendered == nil {
		l.rendered = map[renderKey]*image.RGBA{}
	}
	scaled := l.rendered[key]
	if scaled == nil {
		scaled = image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
		xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, src.Bounds(), draw.Src, nil)
		if key.flip {
			for y := 0; y < r.Dy(); y++ {
				for x := 0; x < r.Dx()/2; x++ {
					a, b := scaled.RGBAAt(x, y), scaled.RGBAAt(r.Dx()-1-x, y)
					scaled.SetRGBA(x, y, b)
					scaled.SetRGBA(r.Dx()-1-x, y, a)
				}
			}
		}
		// Bound resized images even if a user repeatedly changes size/DPI.
		if len(l.renderOrder) >= 512 {
			delete(l.rendered, l.renderOrder[0])
			l.renderOrder = l.renderOrder[1:]
		}
		l.rendered[key] = scaled
		l.renderOrder = append(l.renderOrder, key)
	}
	draw.Draw(dst, r, scaled, image.Point{}, draw.Over)
}

type Clock struct{ Last time.Time }

func (c *Clock) Step(now time.Time) float64 {
	d := now.Sub(c.Last)
	first := c.Last.IsZero()
	c.Last = now
	if first || d < 0 || d > 250*time.Millisecond {
		return 0
	}
	return math.Min(d.Seconds(), 0.075)
}

type Player struct {
	ID                                string
	X, Velocity, WalkPhase, IdlePhase float64
	Action                            string
	Left                              bool
	Remaining, Reaction, Cooldown     float64
	Walking                           bool
	PendingAction, Sequence           string
	SequenceStep                      int
	SequencePhase                     float64
}

func (p *Player) Reset(id string, x float64, left bool) {
	*p = Player{ID: id, X: x, Left: left, Action: "idle", Remaining: 1.5}
}

// React never stacks impulses or indefinitely extends motion under held input.
func (p *Player) React() {
	if p.Cooldown <= 0 {
		p.Reaction = .55
		p.Cooldown = 1.2
	}
}
func (p *Player) Phase() float64 {
	if p.Sequence != "" {
		return p.SequencePhase
	}
	if p.Action == "walk" {
		return p.WalkPhase
	}
	return p.IdlePhase
}
func (p *Player) Advance(v *Loaded, dt, target, height, minX, maxX float64, paused bool) {
	if paused || dt <= 0 {
		return
	}
	if dt > .25 || !finite(dt) {
		return
	}
	dt = math.Min(dt, .075)
	maxX = math.Max(minX, maxX)
	p.X = math.Max(minX, math.Min(maxX, p.X))
	if p.Sequence != "" {
		p.advanceSequence(v, dt)
		return
	}
	if maxX == minX {
		target = 0
	}
	if p.PendingAction != "" {
		target = 0
	}
	if p.Left {
		target = -target
	}
	old := p.X
	// Brake before changing facing, keeping stance feet in the travel direction.
	if p.Velocity*target < 0 {
		target = 0
	}
	p.Velocity += (target - p.Velocity) * (1 - math.Exp(-dt/.22))
	if target == 0 && math.Abs(p.Velocity) < .1 {
		p.Velocity = 0
	}
	p.X = math.Max(minX, math.Min(maxX, p.X+p.Velocity*dt))
	travel := math.Abs(p.X - old)
	if (p.X <= minX && p.Velocity < 0) || (p.X >= maxX && p.Velocity > 0) {
		p.Velocity = 0
		p.Left = p.X >= maxX
		p.Walking = false
		p.Remaining = .6
		p.Reaction = 0
	}
	if travel > 0.0001 {
		p.Action = "walk"
		p.WalkPhase = math.Mod(p.WalkPhase+travel/(v.Profile.Stride*height)*Duration(v.Motions["walk"].Durations), Duration(v.Motions["walk"].Durations))
	} else {
		p.Action = "idle"
		p.IdlePhase = math.Mod(p.IdlePhase+dt*1000, Duration(v.Motions["idle"].Durations))
		if p.PendingAction != "" && p.Velocity == 0 {
			p.startSequence(v)
		}
	}
}
