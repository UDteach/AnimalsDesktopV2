package motion

import (
	appassets "animals-desktop/assets"
	"encoding/json"
	"image"
	"math"
	"testing"
	"testing/fstest"
	"time"
)

func BenchmarkDrawMotion(b *testing.B) {
	s, err := NewStore(appassets.FS, "motions/manifest.json")
	if err != nil {
		b.Fatal(err)
	}
	v := s.Get("chinchilla_standard_gray")
	dst := image.NewRGBA(image.Rect(0, 0, 180, 92))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Draw(dst, "walk", float64(i%30)*42, 10, 89, 114, 76, 1)
	}
}

func fixture(t *testing.T) (*Store, *Loaded) {
	t.Helper()
	s, err := NewStore(appassets.FS, "motions/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	v := s.Get("chinchilla_standard_gray")
	if v == nil {
		t.Fatal(s.Error("chinchilla_standard_gray"))
	}
	return s, v
}
func TestOriginalTimingAndAtomicFallback(t *testing.T) {
	s, v := fixture(t)
	if len(v.Motions["walk"].Durations) != 30 || len(v.Motions["idle"].Durations) != 96 {
		t.Fatal("original frame count lost")
	}
	if Duration(v.Motions["walk"].Durations) != 1250 || Duration(v.Motions["idle"].Durations) != 4000 {
		t.Fatal("original cycle lost")
	}
	files := fstest.MapFS{}
	data, _ := appassets.FS.ReadFile("motions/manifest.json")
	files["manifest.json"] = &fstest.MapFile{Data: data}
	for _, c := range v.Motions {
		for _, tier := range c.Tiers {
			for _, p := range tier {
				b, _ := appassets.FS.ReadFile(p)
				files[p] = &fstest.MapFile{Data: b}
			}
		}
	}
	other := s.Get("chinchilla_beige")
	for _, c := range other.Motions {
		for _, tier := range c.Tiers {
			for _, path := range tier {
				data, _ := appassets.FS.ReadFile(path)
				files[path] = &fstest.MapFile{Data: data}
			}
		}
	}
	broken := v.Motions["idle"].Tiers["64"][95]
	files[broken] = &fstest.MapFile{Data: []byte("broken")}
	bad, err := NewStore(files, "manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if bad.Get(v.ID) != nil || bad.Error(v.ID) == nil {
		t.Fatal("partial variant must never escape loader")
	}
	if bad.Get(other.ID) == nil {
		t.Fatal("one corrupt animal disabled a healthy animal", bad.Error(other.ID))
	}
	delete(files, broken)
	missing, err := NewStore(files, "manifest.json")
	if err != nil || missing.Get(v.ID) != nil || missing.Get(other.ID) == nil {
		t.Fatal("missing frame was not isolated to its animal", err)
	}
	// Repeated reads share images, then selection pruning releases ownership.
	if s.Get(v.ID) != v {
		t.Fatal("images not shared")
	}
	s.Retain(map[string]bool{})
	if len(s.loaded) != 0 {
		t.Fatal("unselected cache retained")
	}
}
func TestValidationRejectsInvalidManifest(t *testing.T) {
	_, v := fixture(t)
	for _, modify := range []func(*Manifest){
		func(m *Manifest) { m.Variants[0].Motions["walk"].Durations[0] = 0 },
		func(m *Manifest) {
			c := m.Variants[0].Motions["walk"]
			c.Tiers["64"] = c.Tiers["64"][:1]
			m.Variants[0].Motions["walk"] = c
		},
		func(m *Manifest) { m.Variants[0].Profile.Stride = 0 },
	} {
		data, _ := appassets.FS.ReadFile("motions/manifest.json")
		var m Manifest
		json.Unmarshal(data, &m)
		modify(&m)
		data, _ = json.Marshal(m)
		files := fstest.MapFS{"manifest.json": &fstest.MapFile{Data: data}}
		s, err := NewStore(files, "manifest.json")
		if err != nil {
			t.Fatal(err)
		}
		if s.Get(v.ID) != nil || s.Error(v.ID) == nil {
			t.Fatal("invalid variant enabled")
		}
	}
}
func TestClockDoesNotIntegrateResume(t *testing.T) {
	var c Clock
	t0 := time.Unix(1000, 0)
	if c.Step(t0) != 0 || c.Step(t0.Add(16*time.Millisecond)) != .016 || c.Step(t0.Add(time.Hour)) != 0 {
		t.Fatal("resume clock jumped")
	}
	if c.Step(t0.Add(time.Hour+16*time.Millisecond)) != .016 {
		t.Fatal("clock failed to resume")
	}
}
func TestFrameBoundaries(t *testing.T) {
	d := []float64{42, 41, 42}
	for _, c := range []struct {
		phase float64
		frame int
	}{{0, 0}, {41.99, 0}, {42, 1}, {83, 2}, {125, 0}, {-1, 2}} {
		if FrameAt(c.phase, d) != c.frame {
			t.Fatalf("phase %v", c.phase)
		}
	}
}
func TestDistanceGaitAndPause(t *testing.T) {
	_, v := fixture(t)
	var p Player
	p.Reset(v.ID, 200, false)
	for i := 0; i < 60; i++ {
		p.Advance(v, .016, 30, 64, 0, 1000, false)
	}
	want := math.Mod((p.X-200)/(64*v.Profile.Stride)*1250, 1250)
	if math.Abs(p.WalkPhase-want) > 1e-8 {
		t.Fatal("feet not locked to distance")
	}
	before := p
	p.Advance(v, .016, 30, 64, 0, 1000, true)
	if p != before {
		t.Fatal("pause advanced")
	}
	p.Advance(v, 60, 30, 64, 0, 1000, false)
	if p != before {
		t.Fatal("long interruption advanced")
	}
	for i := 0; i < 180; i++ {
		p.Advance(v, .016, 0, 64, 0, 1000, false)
	}
	phase := p.WalkPhase
	if p.Action != "idle" {
		t.Fatal("did not settle")
	}
	p.Advance(v, .016, 30, 64, 0, 1000, false)
	if p.WalkPhase < phase && phase < 1200 {
		t.Fatal("restart reset gait phase")
	}
	p.X = 999.9
	p.Advance(v, .016, 100, 64, 0, 1000, false)
	if p.X > 1000 || !p.Left || p.Velocity != 0 {
		t.Fatal("edge failed")
	}
}
func TestRepeatedReactionDoesNotStack(t *testing.T) {
	var p Player
	p.React()
	first := p
	p.React()
	if p != first {
		t.Fatal("reaction extended")
	}
}

func TestGaitDistanceScalesWithSizeAndRetainsPhaseOnResize(t *testing.T) {
	v := &Loaded{Variant: Variant{Profile: Profile{Stride: .3}, Motions: map[string]Clip{"walk": {Durations: []float64{100, 200, 300, 400}}}}}
	for _, left := range []bool{false, true} {
		var p Player
		p.Reset("fixture", 5000, left)
		want := 0.0
		for _, height := range []float64{64, 128, 96, 64} {
			for tick := 0; tick < 100; tick++ {
				before := p.X
				p.Advance(v, .02, height*.3, height, 0, 10000, false)
				want = math.Mod(want+math.Abs(p.X-before)/(height*.3)*1000, 1000)
				if math.Abs(want-p.WalkPhase) > 1e-8 {
					t.Fatalf("height %v left %v: phase %v want %v", height, left, p.WalkPhase, want)
				}
			}
			before := p
			p.Advance(v, .02, 30, height*2, 0, 10000, true)
			if p != before {
				t.Fatal("resize while paused moved or reset gait")
			}
		}
	}
}

func TestAdoptedFramesKeepCanvasAndContactEnvelope(t *testing.T) {
	s, err := NewStore(appassets.FS, "motions/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range s.manifest.Variants {
		v := s.Get(entry.ID)
		if v == nil {
			t.Fatal(entry.ID, s.Error(entry.ID))
		}
		for _, height := range []int{44, 64, 76} {
			for _, dir := range []int{-1, 1} {
				bottomMin, bottomMax := 1000, 0
				for action, clip := range v.Motions {
					phase := 0.0
					for _, duration := range clip.Durations {
						img := image.NewRGBA(image.Rect(0, 0, 200, 100))
						v.Draw(img, action, phase+duration/2, v.HorizontalPadding(height), 89, height*3/2, height, dir)
						bottom := 0
						for y := 0; y < 100; y++ {
							for x := 0; x < 200; x++ {
								alpha := img.RGBAAt(x, y).A
								if alpha >= 128 && y > bottom {
									bottom = y
								}
								if alpha > 0 && (x == 0 || x == 199 || y == 0 || y == 99) {
									t.Fatalf("%s %s clipped", v.ID, action)
								}
							}
						}
						if bottom < 85 || bottom > 91 {
							t.Fatalf("%s/%s height%d baseline=%d", v.ID, action, height, bottom)
						}
						bottomMin = min(bottomMin, bottom)
						bottomMax = max(bottomMax, bottom)
						phase += duration
					}
				}
				// Beige's original short loop includes an airborne phase. Do not
				// flatten that motion by aligning each frame's lowest alpha pixel.
				t.Logf("%s height=%d dir=%d contact envelope=%d..%d", v.ID, height, dir, bottomMin, bottomMax)
			}
		}
		s.Retain(map[string]bool{})
	}
}
