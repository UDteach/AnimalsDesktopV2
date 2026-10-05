package motion

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"testing"
	"testing/fstest"
)

// Deliberately synthetic colored pixels. These are test fixtures, not animal
// assets, generated art, or evidence of an accepted production motion.
func actionFixture(t *testing.T) (fstest.MapFS, Manifest) {
	t.Helper()
	files := fstest.MapFS{}
	m := Manifest{Schema: 1, Hashes: map[string]string{}}
	v := Variant{ID: "fixture-action", Enabled: true, Profile: Profile{Stride: .3, Ground: .85, Facing: 1},
		Motions: map[string]Clip{}, Actions: map[string]ActionSequence{
			"reaction": {Stationary: true, Steps: []ActionStep{{"react_entry", 1}, {"react_hold", 2}, {"react_exit", 1}}},
			"groom":    {Stationary: true, Automatic: true, Steps: []ActionStep{{"react_entry", 1}, {"react_exit", 1}}},
		}}
	names := []string{"walk", "idle", "react_entry", "react_hold", "react_exit"}
	times := [][]float64{{20, 30}, {20, 30}, {10, 20}, {20, 30}, {15, 25}}
	for n, name := range names {
		clip := Clip{Durations: times[n], Presentation: Presentation{Scale: 1}, Tiers: map[string][]string{"64": {}}}
		for frame := 0; frame < 2; frame++ {
			im := image.NewRGBA(image.Rect(0, 0, 96, 64))
			for y := 20; y < 54; y++ {
				for x := 30; x < 60; x++ {
					im.SetRGBA(x, y, color.RGBA{uint8(40 + n*35), uint8(80 + frame*80), 20, 255})
				}
			}
			var b bytes.Buffer
			if err := png.Encode(&b, im); err != nil {
				t.Fatal(err)
			}
			path := fmt.Sprintf("fixture/%s-%d.png", name, frame)
			files[path] = &fstest.MapFile{Data: b.Bytes()}
			hash := sha256.Sum256(b.Bytes())
			m.Hashes[path] = hex.EncodeToString(hash[:])
			clip.Tiers["64"] = append(clip.Tiers["64"], path)
		}
		v.Motions[name] = clip
	}
	m.Variants = []Variant{v}
	return files, m
}

func loadActionFixture(t *testing.T, files fstest.MapFS, manifest Manifest) *Store {
	t.Helper()
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	files["manifest.json"] = &fstest.MapFile{Data: raw}
	s, err := NewStore(files, "manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestActionLoadingAndDrawingUsesDedicatedFrames(t *testing.T) {
	files, m := actionFixture(t)
	s := loadActionFixture(t, files, m)
	v := s.Get("fixture-action")
	if v == nil {
		t.Fatal(s.Error("fixture-action"))
	}
	if len(v.Frames) != 5 || len(v.AutomaticActions()) != 1 || v.AutomaticActions()[0] != "groom" {
		t.Fatal("action clips/automatic list missing")
	}
	var p Player
	p.Reset(v.ID, 10, false)
	p.ReactWith(v)
	p.Advance(v, .016, 0, 64, 0, 500, false)
	if p.Action != "react_entry" || p.Phase() != 0 {
		t.Fatal("entry not selected", p)
	}
	canvas := image.NewRGBA(image.Rect(0, 0, 150, 100))
	v.Draw(canvas, p.Action, p.Phase(), 10, 80, 96, 64, 1)
	// Whole-motion renderer must select entry's own pixel color, not idle art.
	if got := canvas.RGBAAt(50, 60); got.R != 110 || got.G != 80 || got.A != 255 {
		t.Fatal("wrong dedicated frame", got)
	}
	list := v.AutomaticActions()
	list[0] = "broken"
	if v.AutomaticActions()[0] != "groom" {
		t.Fatal("caller mutated automatic list")
	}
}

func TestActionSequenceTimingPauseRepeatAndReturn(t *testing.T) {
	files, m := actionFixture(t)
	v := loadActionFixture(t, files, m).Get("fixture-action")
	var p Player
	p.Reset(v.ID, 100, false)
	if !p.RequestAction(v, "reaction") {
		t.Fatal("request rejected")
	}
	p.Advance(v, .016, 80, 64, 0, 500, false)
	if p.Sequence != "reaction" || p.Action != "react_entry" {
		t.Fatal(p)
	}
	before := p
	if p.RequestAction(v, "reaction") || p.RequestAction(v, "groom") || p != before {
		t.Fatal("busy action restarted/interrupted")
	}
	p.Advance(v, .04, 80, 128, 0, 500, true)
	p.Advance(v, 1, 80, 64, 0, 500, false)
	if p != before {
		t.Fatal("pause or long interruption advanced action")
	}
	for _, step := range []struct {
		dt     float64
		action string
		phase  float64
	}{{.029, "react_entry", 29}, {.001, "react_hold", 0}, {.075, "react_hold", 75}, {.025, "react_exit", 0}, {.04, "idle", 0}} {
		p.Advance(v, step.dt, 80, 64, 0, 500, false)
		if p.Action != step.action || math.Abs(p.Phase()-step.phase) > 1e-8 || p.X != 100 {
			t.Fatalf("step %+v, player %+v", step, p)
		}
	}
	if p.Busy() {
		t.Fatal("finished action stayed busy")
	}
	p.Advance(v, .016, 80, 64, 0, 500, false)
	if p.X <= 100 || p.Action != "walk" {
		t.Fatal("walk did not resume")
	}
}

func TestActionWaitsForBrakingAndPreservesWalkPhase(t *testing.T) {
	files, m := actionFixture(t)
	v := loadActionFixture(t, files, m).Get("fixture-action")
	var p Player
	p.Reset(v.ID, 200, false)
	for i := 0; i < 30; i++ {
		p.Advance(v, .016, 30, 64, 0, 500, false)
	}
	p.ReactWith(v)
	if p.Reaction != 0 || p.PendingAction != "reaction" {
		t.Fatal("dedicated action became translation")
	}
	for i := 0; i < 200 && p.Sequence == ""; i++ {
		old := p.X
		phase := p.WalkPhase
		p.Advance(v, .016, 90, 64, 0, 500, false)
		if p.Sequence == "" && p.X > old && p.Action != "walk" {
			t.Fatal("pose started while sliding")
		}
		if p.Sequence != "" && (p.Velocity != 0 || p.X != old || p.WalkPhase != phase) {
			t.Fatal("sequence consumed a moving tick")
		}
	}
	if p.Sequence != "reaction" {
		t.Fatal("braking did not start action")
	}
	phase := p.WalkPhase
	p.Advance(v, .016, 90, 128, 0, 100, false)
	if p.X != 100 || p.WalkPhase != phase || p.SequencePhase != 16 {
		t.Fatal("viewport/size change broke sequence", p)
	}
	p.Reset("other", 90, true)
	if p.Busy() || p.SequencePhase != 0 {
		t.Fatal("selection retained old action")
	}
}

func TestActionValidationAndCorruptionRemainAtomicPerAnimal(t *testing.T) {
	for name, mutate := range map[string]func(*Manifest, fstest.MapFS){
		"missing clip": func(m *Manifest, f fstest.MapFS) { delete(m.Variants[0].Motions, "react_entry") },
		"stationary undeclared": func(m *Manifest, f fstest.MapFS) {
			a := m.Variants[0].Actions["reaction"]
			a.Stationary = false
			m.Variants[0].Actions["reaction"] = a
		},
		"base clip disguise":     func(m *Manifest, f fstest.MapFS) { m.Variants[0].Actions["reaction"].Steps[0].Motion = "walk" },
		"zero cycles":            func(m *Manifest, f fstest.MapFS) { m.Variants[0].Actions["reaction"].Steps[0].Cycles = 0 },
		"unbounded duration":     func(m *Manifest, f fstest.MapFS) { m.Variants[0].Motions["react_hold"].Durations[0] = 60000 },
		"invalid extra duration": func(m *Manifest, f fstest.MapFS) { m.Variants[0].Motions["react_hold"].Durations[0] = -1 },
		"corrupt extra png": func(m *Manifest, f fstest.MapFS) {
			f["fixture/react_hold-1.png"] = &fstest.MapFile{Data: []byte("corrupt")}
		},
	} {
		t.Run(name, func(t *testing.T) {
			files, m := actionFixture(t)
			healthy := m.Variants[0]
			healthy.ID = "healthy-basic"
			healthy.Actions = nil
			m.Variants = append(m.Variants, healthy)
			mutate(&m, files)
			s := loadActionFixture(t, files, m)
			if s.Get("fixture-action") != nil || s.Error("fixture-action") == nil {
				t.Fatal("invalid action escaped loader")
			}
			if s.Get("healthy-basic") == nil {
				t.Fatal("unrelated base variant failed", s.Error("healthy-basic"))
			}
		})
	}
}

func TestMissingDedicatedReactionKeepsExistingTranslation(t *testing.T) {
	files, m := actionFixture(t)
	m.Variants[0].Actions = nil
	v := loadActionFixture(t, files, m).Get("fixture-action")
	var p Player
	p.Reset(v.ID, 100, false)
	p.ReactWith(v)
	if p.Reaction != .55 || p.Busy() {
		t.Fatal("old reaction changed")
	}
	before := p
	p.ReactWith(v)
	if p != before {
		t.Fatal("held input stacked fallback")
	}
	if p.RequestAction(v, "missing") {
		t.Fatal("missing action accepted")
	}
}
