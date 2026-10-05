//go:build windows

package main

import (
	"animals-desktop/internal/motion"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestLegacyOnlySceneKeepsLegacyRenderCadence(t *testing.T) {
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "")
	rabbit, _ := variantIndexByID("rabbit_gray")
	a := &petApp{petCount: 1, variant: rabbit, coatMode: coatFixed, mode: modeKeyboard, speed: 3, petSizes: defaultPetSizes(), selectedCoats: defaultSelectedCoats()}
	a.pets = []desktopPet{a.newPet(0)}
	if a.tick() {
		t.Fatal("legacy-only scene rendered before its 55ms update")
	}
	a.legacyElapsed = float64(timerInterval) / 1000
	if !a.tick() {
		t.Fatal("legacy scene missed its update")
	}
	a.variant, _ = variantIndexByID("chinchilla_standard_gray")
	a.legacyElapsed = 0
	a.motionClock = motion.Clock{}
	for i := range a.pets {
		a.pets[i].variant = a.variant
	}
	if !a.tick() {
		t.Fatal("timed scene lost its independent redraw cadence")
	}
}

func TestIDOnlySettingsRestoreTimedAndLegacySelection(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANIMALSDESKTOP_SETTINGS_DIR", dir)
	data := []byte(`{"version":4,"coatMode":1,"petCount":3,"selectedCoatIDs":["ferret_sable","chinchilla_standard_gray","rabbit_gray"],"petSizes":[70,100,120]}`)
	if err := os.WriteFile(filepath.Join(dir, "settings.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	a := &petApp{selectedCoats: defaultSelectedCoats()}
	if err := a.loadSettings(); err != nil {
		t.Fatal(err)
	}
	for i, id := range []string{"ferret_sable", "chinchilla_standard_gray", "rabbit_gray"} {
		if variantIDAt(a.selectedCoats[i]) != id {
			t.Fatalf("slot %d lost stable ID", i)
		}
	}
	if err := a.saveSettings(); err != nil {
		t.Fatal(err)
	}
	b := &petApp{selectedCoats: defaultSelectedCoats()}
	if err := b.loadSettings(); err != nil {
		t.Fatal(err)
	}
	if a.selectedCoats != b.selectedCoats || a.petSizes != b.petSizes {
		t.Fatal("selection/size changed on restart")
	}
}

func TestTimedKeyboardReactionKeepsStyleAndNoWheel(t *testing.T) {
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "")
	v := timedVariant("chinchilla_standard_gray")
	if v == nil {
		t.Fatal("timed variant missing")
	}
	a := &petApp{mode: modeKeyboard, wheelEnabled: true, sceneW: 600, speed: 3, petSizes: defaultPetSizes()}
	idx := 0
	for i, c := range variants {
		if c.ID == v.ID {
			idx = i
		}
	}
	a.pets = []desktopPet{{variant: idx, x: 100, dir: 1, item: noItem}}
	a.onTyping()
	p := &a.pets[0]
	if p.timed.Reaction <= 0 || p.state == stateWheel {
		t.Fatal("timed reaction not connected")
	}
	for i := 0; i < 200; i++ {
		a.onTyping()
		tickTimedPlayer(&p.timed, v, .016, p.x, 64, 500, p.dir, 1, false, false)
	}
	if p.timed.Reaction > .55 || p.timed.Velocity > 32 {
		t.Fatal("repeated input stacked")
	}
}

func TestAdoptedHamsterTypingReactionCompletesReturnUnderRepeatedInput(t *testing.T) {
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "")
	v := timedVariant("hamster_golden_syrian")
	if v == nil || !v.HasAction("reaction") {
		t.Fatal("adopted hamster reaction is missing from the compiled assets")
	}
	index, ok := variantIndexByID(v.ID)
	if !ok {
		t.Fatal("hamster public ID missing")
	}
	a := &petApp{mode: modeKeyboard, wheelEnabled: true, sceneW: 600, speed: 3, petSizes: defaultPetSizes()}
	a.pets = []desktopPet{{variant: index, x: 100, dir: 1, item: noItem}}
	a.onTyping()
	p := &a.pets[0]
	if p.timed.PendingAction != "reaction" || p.timed.Reaction != 0 || p.state == stateWheel {
		t.Fatal("typing did not choose the adopted stationary reaction")
	}
	seen := map[int]bool{}
	returns := 0
	for tick := 0; tick < 40; tick++ {
		before := p.timed
		a.onTyping()
		tickTimedPlayer(&p.timed, v, .02, p.x, 64, 500, p.dir, 1, false, false)
		if p.timed.Sequence == "reaction" {
			seen[motion.FrameAt(p.timed.Phase(), v.Motions[p.timed.Action].Durations)] = true
		}
		if before.Sequence == "reaction" && p.timed.Sequence == "" && p.timed.Action == "idle" {
			returns++
		}
		if tick == 8 {
			paused := p.timed
			tickTimedPlayer(&p.timed, v, .2, p.x, 64, 500, p.dir, 1, false, true)
			if p.timed != paused {
				t.Fatal("pause changed the adopted reaction")
			}
		}
		if p.timed.X != 100 || p.timed.Reaction != 0 {
			t.Fatal("stationary typing reaction moved or fell back to an impulse")
		}
	}
	if len(seen) != 3 || returns != 1 || p.timed.Busy() || p.timed.Action != "idle" {
		t.Fatalf("typing skipped/restarted the return: frames=%v returns=%d player=%+v", seen, returns, p.timed)
	}
}

func TestTimedPauseAndSelectionReset(t *testing.T) {
	v := timedVariant("chinchilla_standard_gray")
	var p motion.Player
	tickTimedPlayer(&p, v, .016, 100, 64, 500, 1, 1, true, false)
	before := p
	tickTimedPlayer(&p, v, .016, 100, 64, 500, 1, 1, true, true)
	if p != before {
		t.Fatal("hidden pet advanced")
	}
	p.ID = "old"
	p.WalkPhase = 600
	p.Reaction = 1
	tickTimedPlayer(&p, v, .016, 100, 64, 500, 1, 1, false, false)
	if p.ID != v.ID || p.WalkPhase != 0 || p.Reaction != 0 {
		t.Fatal("coat changed without reset")
	}
}
func TestPremultipliedBGRAKeepsSoftEdge(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.SetRGBA(0, 0, color.RGBA{100, 60, 20, 128})
	b := make([]byte, 4)
	copyPremultipliedBGRA(b, img)
	if b[0] != 20 || b[1] != 60 || b[2] != 100 || b[3] != 128 {
		t.Fatal(b)
	}
	got := overRGBA(color.RGBA{}, img.RGBAAt(0, 0))
	if got != img.RGBAAt(0, 0) {
		t.Fatal("transparent destination darkened edge", got)
	}
}
