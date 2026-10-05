//go:build darwin

package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"animals-desktop/internal/catalog"
)

func TestDarwinHorizontalMotionFramesUseStableSequence(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}

	for _, delay := range []int{1, 2} {
		for tick := 0; tick < 64; tick++ {
			got := seqFrameFrom(walkFrameSeq, tick, delay)
			if !allowed[got] {
				t.Fatalf("seqFrameFrom(walkFrameSeq, %d, %d) = %d, want stable horizontal walk frame", tick, delay, got)
			}
		}
	}
}

func TestDarwinTickWalksInDirectionWithStableFrame(t *testing.T) {
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "legacy")
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	a := &darwinPetApp{
		sceneW: 500,
		mode:   darwinModeRandom,
		pets: []darwinPet{
			{x: 20, dir: -1, speed: 2, nextPause: 50},
			{x: 80, dir: 1, speed: 2, nextPause: 50},
		},
	}

	a.legacyElapsed = float64(timerInterval) / 1000
	a.tickPets()
	if got := a.pets[0].x; got != 18 {
		t.Fatalf("left-moving pet x = %d, want 18", got)
	}
	if got := a.pets[1].x; got != 82 {
		t.Fatalf("right-moving pet x = %d, want 82", got)
	}
	for i, pet := range a.pets {
		if !allowed[pet.frame] {
			t.Fatalf("pet %d frame = %d, want stable horizontal walk frame", i, pet.frame)
		}
	}
}

func TestDarwinTickKeepsFerretPauseActionAnimated(t *testing.T) {
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "legacy")
	tests := []struct {
		variantID string
		action    int
		start     int
		count     int
	}{
		{variantID: "ferret_sable_panda", action: 0, start: ferretSniffStart, count: ferretSniffFrames},
		{variantID: "ferret_sable", action: 1, start: ferretGroomStart, count: ferretGroomFrames},
		{variantID: "ferret_albino", action: 2, start: ferretRestStart, count: ferretRestFrames},
	}

	for _, test := range tests {
		t.Run(test.variantID, func(t *testing.T) {
			variantIndex := -1
			for i, variant := range darwinVariants {
				if variant.ID == test.variantID {
					variantIndex = i
					break
				}
			}
			if variantIndex < 0 {
				t.Fatalf("runtime variant %q not found", test.variantID)
			}

			const startX = 120
			a := &darwinPetApp{
				sceneW: 500,
				mode:   darwinModeRandom,
				pets: []darwinPet{{
					x:           startX,
					dir:         1,
					speed:       3,
					variant:     variantIndex,
					pause:       6,
					pauseAction: test.action,
				}},
			}
			for step := 0; step < 4; step++ {
				a.legacyElapsed = float64(timerInterval) / 1000
				a.tickPets()
				if got := a.pets[0].x; got != startX {
					t.Fatalf("paused pet x after step %d = %d, want %d", step, got, startX)
				}
				got := a.pets[0].frame
				if got < test.start || got >= test.start+test.count {
					t.Fatalf("paused pet frame after step %d = %d, want range %d-%d", step, got, test.start, test.start+test.count-1)
				}
			}
		})
	}
}

func TestDarwinTickKeepsNonFerretPauseIdle(t *testing.T) {
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "legacy")
	variantIndex := -1
	for i, variant := range darwinVariants {
		if variant.ID == "hamster_golden_syrian" {
			variantIndex = i
			break
		}
	}
	if variantIndex < 0 {
		t.Fatal("runtime hamster variant not found")
	}

	a := &darwinPetApp{
		sceneW: 500,
		mode:   darwinModeRandom,
		pets: []darwinPet{{
			x:           120,
			dir:         1,
			speed:       3,
			variant:     variantIndex,
			pause:       4,
			pauseAction: 1,
		}},
	}
	for step := 0; step < 3; step++ {
		a.legacyElapsed = float64(timerInterval) / 1000
		a.tickPets()
		want := seqFrameFrom(idleFrameSeq, a.tick, 5)
		if got := a.pets[0].frame; got != want {
			t.Fatalf("non-ferret paused frame after step %d = %d, want unchanged idle frame %d", step, got, want)
		}
	}
}

func TestDarwinMovePetReflectsAtSceneEdges(t *testing.T) {
	a := &darwinPetApp{sceneW: 220}

	left := darwinPet{x: 1, dir: -1}
	a.movePet(&left, 4)
	if left.x != 0 || left.dir != 1 {
		t.Fatalf("left edge reflection = x:%d dir:%d, want x 0 dir +1", left.x, left.dir)
	}

	right := darwinPet{x: a.sceneW - spriteW - 1, dir: 1}
	a.movePet(&right, 4)
	if right.x != a.sceneW-spriteW || right.dir != -1 {
		t.Fatalf("right edge reflection = x:%d dir:%d, want max dir -1", right.x, right.dir)
	}
}

func TestDarwinSeqFrameFromHandlesEmptyAndBadDivisor(t *testing.T) {
	if got := seqFrameFrom(nil, 12, 2); got != idleStart {
		t.Fatalf("seqFrameFrom(nil) = %d, want %d", got, idleStart)
	}
	seq := []int{7, 9}
	if got := seqFrameFrom(seq, 3, 0); got != 9 {
		t.Fatalf("seqFrameFrom with zero divisor = %d, want 9", got)
	}
}

func TestDarwinRandomPauseAvoidsWeakNibbleFrames(t *testing.T) {
	for _, variantID := range []string{"sugar_glider_gray", "rabbit_chestnut_agouti"} {
		for tick := 0; tick < 32; tick++ {
			got := darwinRandomPauseFrame(variantID, 0, tick)
			if got < hopStart || got >= hopStart+4 {
				t.Fatalf("%s pause nibble tick %d = %d, want stable action fallback", variantID, tick, got)
			}
		}
	}

	for tick := 0; tick < 24; tick++ {
		got := darwinRandomPauseFrame("hamster_golden_syrian", 0, tick)
		if got < nibbleStart || got >= nibbleStart+3 {
			t.Fatalf("hamster pause nibble tick %d = %d, want original nibble frames", tick, got)
		}
	}
}

func TestDarwinFerretRandomPauseUsesProfileActions(t *testing.T) {
	actions := []struct {
		name   string
		action int
		start  int
		count  int
	}{
		{name: "sniff", action: 0, start: ferretSniffStart, count: ferretSniffFrames},
		{name: "groom", action: 1, start: ferretGroomStart, count: ferretGroomFrames},
		{name: "rest", action: 2, start: ferretRestStart, count: ferretRestFrames},
	}

	for _, variantID := range []string{"ferret_sable_panda", "ferret_sable", "ferret_albino"} {
		t.Run(variantID, func(t *testing.T) {
			if got := darwinMotionProfileForVariantID(variantID); got != catalog.MotionProfileFerretSlink {
				t.Fatalf("motion profile = %q, want %q", got, catalog.MotionProfileFerretSlink)
			}
			for _, variant := range darwinVariants {
				if variant.ID == variantID && variant.WheelCapable {
					t.Fatal("ferret must not be wheel-capable")
				}
			}
			for _, action := range actions {
				t.Run(action.name, func(t *testing.T) {
					seen := make(map[int]bool, action.count)
					for tick := 0; tick < 96; tick++ {
						got := darwinRandomPauseFrame(variantID, action.action, tick)
						if got < action.start || got >= action.start+action.count {
							t.Fatalf("tick %d = frame %d, want %s range %d-%d", tick, got, action.name, action.start, action.start+action.count-1)
						}
						seen[got] = true
					}
					for frame := action.start; frame < action.start+action.count; frame++ {
						if !seen[frame] {
							t.Fatalf("%s frame %d was never selected", action.name, frame)
						}
					}
				})
			}
		})
	}
}

func TestDarwinNonFerretRandomPauseRemainsGeneric(t *testing.T) {
	const variantID = "hamster_golden_syrian"
	if got := darwinMotionProfileForVariantID(variantID); got == catalog.MotionProfileFerretSlink {
		t.Fatalf("motion profile = %q, must not use ferret pause actions", got)
	}

	actions := []struct {
		action int
		seq    []int
		delay  int
	}{
		{action: 0, seq: nibbleFrameSeq, delay: 3},
		{action: 1, seq: hopFrameSeq, delay: 2},
		{action: 2, seq: idleFrameSeq, delay: 5},
	}
	for _, action := range actions {
		for tick := 0; tick < 64; tick++ {
			got := darwinRandomPauseFrame(variantID, action.action, tick)
			want := seqFrameFrom(action.seq, tick, action.delay)
			if got != want {
				t.Fatalf("action %d tick %d = %d, want unchanged generic frame %d", action.action, tick, got, want)
			}
		}
	}
}

func TestDarwinRuntimeVariantsMirrorCatalog(t *testing.T) {
	runtimeVariants := catalog.RuntimeVariants()
	if len(darwinVariants) != len(runtimeVariants) {
		t.Fatalf("darwinVariants = %d, want catalog runtime count %d", len(darwinVariants), len(runtimeVariants))
	}
	for i, want := range runtimeVariants {
		got := darwinVariants[i]
		if got.ID != want.ID || got.SpriteBase != want.SpriteBase || got.LabelJA != want.LabelJA || got.LabelEN != want.LabelEN {
			t.Fatalf("darwinVariants[%d] = %+v, want catalog variant %+v", i, got, want)
		}
		if got.MotionProfile != catalog.MotionProfileForVariant(want) {
			t.Fatalf("darwinVariants[%d].MotionProfile = %q, want %q", i, got.MotionProfile, catalog.MotionProfileForVariant(want))
		}
		if got.LabelJA == "" || got.LabelEN == "" {
			t.Fatalf("darwinVariants[%d] has empty labels: %+v", i, got)
		}
	}
}

func TestDarwinCanSelectEveryRuntimeVariant(t *testing.T) {
	a := &darwinPetApp{
		sceneW:        900,
		speed:         darwinSpeedNormal,
		petCount:      maxPetCount,
		coatMode:      darwinCoatFixed,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
	}
	a.resetPets()
	if len(darwinVariants) != len(catalog.RuntimeVariants()) {
		t.Fatalf("darwinVariants = %d, want %d release-scoped selectable animals", len(darwinVariants), len(catalog.RuntimeVariants()))
	}

	for i, variant := range darwinVariants {
		a.setCoatMode(int(darwinCoatFixed))
		a.setFixedVariant(i)
		if got := a.variant; got != i {
			t.Fatalf("fixed variant index = %d, want %d", got, i)
		}
		for petIndex, pet := range a.pets {
			if pet.variant != i {
				t.Fatalf("fixed variant %q pet %d = %d, want %d", variant.ID, petIndex, pet.variant, i)
			}
		}
		if got := a.variantID(a.variant); got != variant.ID {
			t.Fatalf("fixed variant ID = %q, want %q", got, variant.ID)
		}
		if got := darwinVariantLabel(i, darwinLangJapanese); got == "" {
			t.Fatalf("Japanese label for %q is empty", variant.ID)
		}
		if got := darwinVariantLabel(i, darwinLangEnglish); got == "" {
			t.Fatalf("English label for %q is empty", variant.ID)
		}

		a.setCoatMode(int(darwinCoatSelected))
		petIndex := i % maxPetCount
		a.setSelectedVariant(petIndex, i)
		if a.selectedCoats[petIndex] != i || a.pets[petIndex].variant != i {
			t.Fatalf("selected variant %q pet %d = selected:%d pet:%d, want %d", variant.ID, petIndex, a.selectedCoats[petIndex], a.pets[petIndex].variant, i)
		}
	}
}

func TestDarwinSettingsUpdateRuntimeStateAndPersist(t *testing.T) {
	oldSettingsPath := darwinSettingsPath
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) {
		return settingsPath, nil
	}
	t.Cleanup(func() {
		darwinSettingsPath = oldSettingsPath
	})

	a := &darwinPetApp{
		sceneW:        640,
		speed:         darwinSpeedNormal,
		petCount:      5,
		mode:          darwinModeRandom,
		coatMode:      darwinCoatRandom,
		selectedCoats: defaultDarwinSelectedCoats(),
		wheelEnabled:  true,
	}
	a.resetPets()

	a.setSpeed(darwinSpeedFast)
	if got := a.speed; got != darwinSpeedFast {
		t.Fatalf("speed = %d, want %d", got, darwinSpeedFast)
	}
	if got := a.pets[0].speed; got != 4 {
		t.Fatalf("first pet speed = %d, want 4", got)
	}
	if got := a.pets[1].speed; got != 5 {
		t.Fatalf("second pet speed = %d, want 5", got)
	}

	a.setPetCount(3)
	if got := len(a.pets); got != 3 {
		t.Fatalf("pet count = %d, want 3", got)
	}
	a.setPetCount(9)
	if got := len(a.pets); got != 9 {
		t.Fatalf("pet count = %d, want 9", got)
	}
	a.setPetCount(3)
	a.setMode(int(darwinModeKeyboard))
	if got := a.mode; got != darwinModeKeyboard {
		t.Fatalf("mode = %d, want %d", got, darwinModeKeyboard)
	}
	a.setCoatMode(int(darwinCoatSelected))
	a.setSelectedVariant(1, 7)
	if a.coatMode != darwinCoatSelected || a.selectedCoats[1] != 7 || a.pets[1].variant != 7 {
		t.Fatalf("selected coat state = mode:%d selected:%d pet:%d, want selected mode and variant 7", a.coatMode, a.selectedCoats[1], a.pets[1].variant)
	}
	a.setPetSize(1, 117)
	a.setPetSize(2, 60)
	if a.petSizePercent(1) != 120 || a.petSizePercent(2) != 70 {
		t.Fatalf("pet sizes = %d/%d, want rounded 120 and clamped 70", a.petSizePercent(1), a.petSizePercent(2))
	}

	a.keyHold = wheelKeyHold
	a.setWheelEnabled(false)
	if a.wheelEnabled {
		t.Fatal("wheelEnabled = true, want false")
	}
	if got := a.keyHold; got != 0 {
		t.Fatalf("keyHold = %d, want 0 after disabling keyboard reaction", got)
	}
	a.nameLabels = true
	a.setPetName(0, "  モカ  ")
	a.setPetName(2, "abcdefghijklmnopqrstuvwxyz")
	a.lang = darwinLangEnglish
	a.setDisplayID(123456789)

	a.saveSettings()
	b := &darwinPetApp{
		speed:         darwinSpeedNormal,
		petCount:      5,
		mode:          darwinModeRandom,
		coatMode:      darwinCoatRandom,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
		lang:          darwinLangJapanese,
	}
	b.loadSettings()
	if b.speed != darwinSpeedFast || b.petCount != 3 || b.lang != darwinLangEnglish || b.displayID != 123456789 || b.mode != darwinModeKeyboard || b.coatMode != darwinCoatSelected || b.selectedCoats[1] != 7 || b.wheelEnabled {
		t.Fatalf("loaded settings = speed:%d count:%d lang:%d display:%d mode:%d coat:%d selected:%d wheel:%v, want speed:%d count:3 english display 123456789 keyboard selected variant 7 wheel:false", b.speed, b.petCount, b.lang, b.displayID, b.mode, b.coatMode, b.selectedCoats[1], b.wheelEnabled, darwinSpeedFast)
	}
	if b.petSizePercent(1) != 120 || b.petSizePercent(2) != 70 {
		t.Fatalf("loaded pet sizes = %d/%d, want 120/70", b.petSizePercent(1), b.petSizePercent(2))
	}
	if !b.nameLabels || b.petNames[0] != "モカ" || b.petNames[2] != "abcdefghijklmnopqrstuvwx" {
		t.Fatalf("loaded names = labels:%v names:%#v", b.nameLabels, b.petNames[:3])
	}
	if got := b.petDisplayName(1); got != "Animal 2" {
		t.Fatalf("default display name = %q, want Animal 2", got)
	}
}

func TestDarwinSettingsPersistAndLoadVariantIDs(t *testing.T) {
	oldSettingsPath := darwinSettingsPath
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) {
		return settingsPath, nil
	}
	t.Cleanup(func() {
		darwinSettingsPath = oldSettingsPath
	})

	shoebillIndex, ok := darwinVariantIndexByID("shoebill_stork")
	if !ok {
		t.Fatal("shoebill_stork missing from runtime variants")
	}
	catIndex, ok := darwinVariantIndexByID("domestic_shorthair_tabby_white_stocky")
	if !ok {
		t.Fatal("domestic_shorthair_tabby_white_stocky missing from runtime variants")
	}
	data := []byte(fmt.Sprintf(`{
  "version": %d,
  "variant": %d,
  "variantID": "shoebill_stork",
  "coatMode": 1,
  "selectedCoats": [%d, %d],
  "selectedCoatIDs": ["domestic_shorthair_tabby_white_stocky", "shoebill_stork"],
  "petCount": 2,
  "speed": 5,
  "mode": 1,
  "wheelEnabled": true
}`, darwinSettingsVersion, catIndex, catIndex, catIndex))
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &darwinPetApp{
		speed:         darwinSpeedNormal,
		petCount:      1,
		coatMode:      darwinCoatSelected,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
	}
	a.loadSettings()
	if a.variant != shoebillIndex {
		t.Fatalf("variant = %d (%s), want shoebill index %d", a.variant, darwinVariantIDAt(a.variant), shoebillIndex)
	}
	if a.selectedCoats[0] != catIndex || a.selectedCoats[1] != shoebillIndex {
		t.Fatalf("selectedCoats = %v, want cat then shoebill", a.selectedCoats[:2])
	}
}

func TestDarwinLegacySettingsRecoverShoebillIndexDrift(t *testing.T) {
	oldSettingsPath := darwinSettingsPath
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) {
		return settingsPath, nil
	}
	t.Cleanup(func() {
		darwinSettingsPath = oldSettingsPath
	})

	shoebillIndex, ok := darwinVariantIndexByID("shoebill_stork")
	if !ok {
		t.Fatal("shoebill_stork missing from runtime variants")
	}
	data := []byte(`{
  "version": 3,
  "variant": 41,
  "coatMode": 1,
  "selectedCoats": [41, 0],
  "petNames": ["爆速ハシビロコウさん", ""],
  "petCount": 2,
  "speed": 5,
  "mode": 1,
  "wheelEnabled": true
}`)
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &darwinPetApp{
		speed:         darwinSpeedNormal,
		petCount:      1,
		coatMode:      darwinCoatSelected,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
	}
	a.loadSettings()
	if a.variant != shoebillIndex || a.selectedCoats[0] != shoebillIndex {
		t.Fatalf("legacy shoebill migration = variant:%d selected:%d ids:%s/%s, want %d", a.variant, a.selectedCoats[0], darwinVariantIDAt(a.variant), darwinVariantIDAt(a.selectedCoats[0]), shoebillIndex)
	}
	if got := a.petNames[0]; got != "爆速ハシビロコウさん" {
		t.Fatalf("pet name = %q", got)
	}
}

func TestDarwinLegacySettingsKeepAmbiguousIndexWithoutShoebillName(t *testing.T) {
	oldSettingsPath := darwinSettingsPath
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) {
		return settingsPath, nil
	}
	t.Cleanup(func() {
		darwinSettingsPath = oldSettingsPath
	})

	data := []byte(`{
  "version": 3,
  "variant": 41,
  "coatMode": 1,
  "selectedCoats": [41],
  "petNames": [""],
  "petCount": 1,
  "speed": 5,
  "mode": 1,
  "wheelEnabled": true
}`)
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &darwinPetApp{
		speed:         darwinSpeedNormal,
		petCount:      1,
		coatMode:      darwinCoatSelected,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
	}
	a.loadSettings()
	if a.variant != 41 || a.selectedCoats[0] != 41 {
		t.Fatalf("ambiguous legacy index changed to variant:%d selected:%d, want raw index 41", a.variant, a.selectedCoats[0])
	}
}

func TestDarwinVariantGroupLabels(t *testing.T) {
	shoebillIndex, ok := darwinVariantIndexByID("shoebill_stork")
	if !ok {
		t.Fatal("shoebill_stork missing from runtime variants")
	}
	if got := darwinVariantGroupLabel(shoebillIndex, darwinLangJapanese); got != "鳥" {
		t.Fatalf("Japanese shoebill group = %q", got)
	}
	if got := darwinVariantGroupLabel(shoebillIndex, darwinLangEnglish); got != "Birds" {
		t.Fatalf("English shoebill group = %q", got)
	}

	for _, variantID := range []string{
		"ferret_sable_panda",
		"ferret_sable",
		"ferret_albino",
	} {
		index, ok := darwinVariantIndexByID(variantID)
		if !ok {
			t.Fatalf("%s missing from runtime variants", variantID)
		}
		if got := darwinVariantGroupLabel(index, darwinLangJapanese); got != "小動物" {
			t.Errorf("Japanese %s group = %q", variantID, got)
		}
		if got := darwinVariantGroupLabel(index, darwinLangEnglish); got != "Small animals" {
			t.Errorf("English %s group = %q", variantID, got)
		}
	}
}

func TestDarwinLanguageDefaultsAndNormalizes(t *testing.T) {
	if got := normalizeDarwinLanguage(-1); got != darwinLangJapanese {
		t.Fatalf("normalizeDarwinLanguage(-1) = %d, want Japanese", got)
	}
	if got := normalizeDarwinLanguage(1); got != darwinLangEnglish {
		t.Fatalf("normalizeDarwinLanguage(1) = %d, want English", got)
	}
	if got := normalizeDarwinLanguage(99); got != darwinLangJapanese {
		t.Fatalf("normalizeDarwinLanguage(99) = %d, want Japanese", got)
	}

	a := &darwinPetApp{lang: darwinLangJapanese}
	if got := a.petDisplayName(2); got != "どうぶつ3" {
		t.Fatalf("Japanese default display name = %q, want どうぶつ3", got)
	}
	a.lang = darwinLangEnglish
	if got := a.petDisplayName(2); got != "Animal 3" {
		t.Fatalf("English default display name = %q, want Animal 3", got)
	}
}

func TestDarwinDisplayIDDefaultsAndNormalizes(t *testing.T) {
	if got := normalizeDarwinDisplayID(-1); got != 0 {
		t.Fatalf("normalizeDarwinDisplayID(-1) = %d, want 0", got)
	}
	if got := normalizeDarwinDisplayID(0); got != 0 {
		t.Fatalf("normalizeDarwinDisplayID(0) = %d, want 0", got)
	}
	if got := normalizeDarwinDisplayID(42); got != 42 {
		t.Fatalf("normalizeDarwinDisplayID(42) = %d, want 42", got)
	}

	a := &darwinPetApp{}
	a.setDisplayID(-10)
	if a.displayID != 0 {
		t.Fatalf("negative display id = %d, want 0", a.displayID)
	}
	a.setDisplayID(98765)
	if a.displayID != 98765 {
		t.Fatalf("display id = %d, want 98765", a.displayID)
	}
}

func TestDarwinLoadsVersionTwoSettingsWithoutOldAnimalSelection(t *testing.T) {
	oldSettingsPath := darwinSettingsPath
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) {
		return settingsPath, nil
	}
	t.Cleanup(func() {
		darwinSettingsPath = oldSettingsPath
	})

	data := []byte(`{
  "version": 2,
  "variant": 8,
  "coatMode": 2,
  "selectedCoats": [8, 7, 6, 5, 4, 3, 2, 1, 0, 9],
  "speed": 5,
  "language": 1,
  "mode": 0,
  "petSizes": [120, 120, 120, 120, 120, 120, 120, 120, 120, 120],
  "petCount": 4,
  "wheelEnabled": false,
  "nameLabels": true,
  "petNames": ["モカ"]
}`)
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &darwinPetApp{
		speed:         darwinSpeedNormal,
		petCount:      5,
		mode:          darwinModeRandom,
		coatMode:      darwinCoatSelected,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
	}
	a.loadSettings()
	if a.speed != darwinSpeedFast || a.petCount != 4 || a.lang != darwinLangEnglish || a.mode != darwinModeKeyboard || a.wheelEnabled {
		t.Fatalf("loaded preserved settings = speed:%d count:%d lang:%d mode:%d wheel:%v", a.speed, a.petCount, a.lang, a.mode, a.wheelEnabled)
	}
	if a.coatMode != darwinCoatSelected || a.variant != 0 || a.selectedCoats != defaultDarwinSelectedCoats() {
		t.Fatalf("old animal selection was not reset: variant:%d coat:%d selected:%v", a.variant, a.coatMode, a.selectedCoats)
	}
	if got := a.petSizePercent(0); got != defaultPetSizePercent {
		t.Fatalf("old pet size = %d, want default %d", got, defaultPetSizePercent)
	}
	if !a.nameLabels || a.petNames[0] != "モカ" {
		t.Fatalf("loaded name settings = labels:%v names:%v", a.nameLabels, a.petNames[:1])
	}
}

func TestDarwinPetCountSupportsEveryVisibleCount(t *testing.T) {
	a := &darwinPetApp{
		sceneW:        900,
		speed:         darwinSpeedNormal,
		petCount:      5,
		mode:          darwinModeRandom,
		coatMode:      darwinCoatRandom,
		selectedCoats: defaultDarwinSelectedCoats(),
		wheelEnabled:  true,
	}
	for count := 1; count <= maxPetCount; count++ {
		a.setPetCount(count)
		if a.petCount != count || len(a.pets) != count {
			t.Fatalf("setPetCount(%d) = state:%d pets:%d", count, a.petCount, len(a.pets))
		}
	}
}

func TestDarwinClickReactionHitTestsPet(t *testing.T) {
	a := &darwinPetApp{
		sceneW:       400,
		wheelEnabled: true,
		pets: []darwinPet{
			{x: 50, lane: 0},
			{x: 70, lane: 7},
		},
	}

	if index := a.petAtScenePoint(90, 40); index != 1 {
		t.Fatalf("petAtScenePoint overlapping pets = %d, want topmost pet 1", index)
	}
	if !a.addClickReaction(90, 40) {
		t.Fatal("addClickReaction inside pet = false, want true")
	}
	if len(a.reactions) != 1 || a.reactions[0].pet != 1 || a.reactions[0].ticks != reactionTicks {
		t.Fatalf("reaction = %#v", a.reactions)
	}
	a.reactions[0].ticks = 10
	if !a.addClickReaction(90, 40) {
		t.Fatal("second addClickReaction inside pet = false, want true")
	}
	if len(a.reactions) != 1 || a.reactions[0].ticks != reactionTicks {
		t.Fatalf("updated reaction = %#v", a.reactions)
	}
	if a.addClickReaction(3, 3) {
		t.Fatal("addClickReaction outside pet = true, want false")
	}

	a.keyHold = 1
	a.pets = []darwinPet{{x: 50, lane: 0}}
	if index := a.petAtScenePoint(90, 40); index != -1 {
		t.Fatalf("wheel runner hit = %d, want ignored", index)
	}
}

func TestDarwinPetSizeAffectsBoundsAndHitTesting(t *testing.T) {
	a := &darwinPetApp{
		sceneW:   220,
		petSizes: defaultDarwinPetSizes(),
		pets: []darwinPet{
			{x: 200, lane: 0, dir: 1},
		},
	}

	a.setPetSize(0, 117)
	w, h := a.petSpriteSize(0)
	if w != 115 || h != 76 {
		t.Fatalf("120%% sprite size = %dx%d, want 115x76", w, h)
	}
	if wantX := a.sceneW - w; a.pets[0].x != wantX {
		t.Fatalf("large pet x = %d, want clamped %d", a.pets[0].x, wantX)
	}
	centerX := a.pets[0].x + w/2
	centerY := sceneH - h + h/2
	if got := a.petAtScenePoint(centerX, centerY); got != 0 {
		t.Fatalf("large pet hit = %d, want 0", got)
	}
	if got := a.petAtScenePoint(a.pets[0].x+w+2, centerY); got != -1 {
		t.Fatalf("outside large pet hit = %d, want -1", got)
	}

	a.pets[0].x = 200
	a.setPetSize(0, 60)
	w, h = a.petSpriteSize(0)
	if w != 67 || h != 44 {
		t.Fatalf("70%% sprite size = %dx%d, want 67x44", w, h)
	}
	if wantX := a.sceneW - w; a.pets[0].x != wantX {
		t.Fatalf("small pet x = %d, want clamped %d", a.pets[0].x, wantX)
	}
}

func TestDarwinPetSizeSupportsEveryVisibleStep(t *testing.T) {
	a := &darwinPetApp{
		sceneW:   300,
		petSizes: defaultDarwinPetSizes(),
		pets: []darwinPet{
			{x: 250, lane: 0, dir: 1},
		},
	}

	lastW := 0
	lastH := 0
	for _, size := range []int{70, 80, 90, 100, 110, 120} {
		a.setPetSize(0, size)
		if got := a.petSizePercent(0); got != size {
			t.Fatalf("pet size percent = %d, want %d", got, size)
		}
		w, h := a.petSpriteSize(0)
		if wantW, wantH := frameW*size/100, frameH*size/100; w != wantW || h != wantH {
			t.Fatalf("%d%% sprite size = %dx%d, want %dx%d", size, w, h, wantW, wantH)
		}
		if w <= lastW || h <= lastH {
			t.Fatalf("%d%% sprite size = %dx%d, want larger than previous %dx%d", size, w, h, lastW, lastH)
		}
		if a.pets[0].x > a.sceneW-w {
			t.Fatalf("%d%% pet x = %d, want clamped to <= %d", size, a.pets[0].x, a.sceneW-w)
		}
		lastW, lastH = w, h
	}
}

func TestDarwinPartialSettingsKeepRuntimeDefaults(t *testing.T) {
	oldSettingsPath := darwinSettingsPath
	settingsPath := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) {
		return settingsPath, nil
	}
	t.Cleanup(func() {
		darwinSettingsPath = oldSettingsPath
	})
	if err := os.WriteFile(settingsPath, []byte(`{"version":1}`), 0o644); err != nil {
		t.Fatal(err)
	}

	a := &darwinPetApp{
		speed:        darwinSpeedNormal,
		petCount:     5,
		wheelEnabled: true,
	}
	a.loadSettings()
	if a.speed != darwinSpeedNormal || a.petCount != 5 || !a.wheelEnabled {
		t.Fatalf("settings defaults = speed:%d count:%d wheel:%v, want speed:%d count:5 wheel:true", a.speed, a.petCount, a.wheelEnabled, darwinSpeedNormal)
	}
}

func TestDarwinDrawFacingImageMirrorsNegativeDirection(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	red := color.RGBA{R: 255, A: 255}
	blue := color.RGBA{B: 255, A: 255}
	src.SetRGBA(0, 0, red)
	src.SetRGBA(1, 0, blue)

	dst := image.NewRGBA(image.Rect(0, 0, 2, 1))
	drawFacingImage(dst, src, dst.Bounds(), 1)
	if got := dst.RGBAAt(0, 0); got != red {
		t.Fatalf("drawFacingImage positive left pixel = %#v, want %#v", got, red)
	}
	if got := dst.RGBAAt(1, 0); got != blue {
		t.Fatalf("drawFacingImage positive right pixel = %#v, want %#v", got, blue)
	}

	dst = image.NewRGBA(image.Rect(0, 0, 2, 1))
	drawFacingImage(dst, src, dst.Bounds(), -1)
	if got := dst.RGBAAt(0, 0); got != blue {
		t.Fatalf("drawFacingImage negative left pixel = %#v, want %#v", got, blue)
	}
	if got := dst.RGBAAt(1, 0); got != red {
		t.Fatalf("drawFacingImage negative right pixel = %#v, want %#v", got, red)
	}
}

func TestDarwinDrawDirectionCompensatesLeftFacingSource(t *testing.T) {
	if got := darwinDrawDirection(1, "sugar_glider_gray"); got != -1 {
		t.Fatalf("sugar glider right-moving draw direction = %d, want -1", got)
	}
	if got := darwinDrawDirection(-1, "sugar_glider_gray"); got != 1 {
		t.Fatalf("sugar glider left-moving draw direction = %d, want 1", got)
	}
	if got := darwinDrawDirection(1, "hamster_golden_syrian"); got != 1 {
		t.Fatalf("hamster right-moving draw direction = %d, want 1", got)
	}
	if got := darwinDrawDirection(-1, "hamster_golden_syrian"); got != -1 {
		t.Fatalf("hamster left-moving draw direction = %d, want -1", got)
	}
}

func TestDarwinIDOnlySettingsRestoreSelection(t *testing.T) {
	previous := darwinSettingsPath
	path := filepath.Join(t.TempDir(), "settings.json")
	darwinSettingsPath = func() (string, error) { return path, nil }
	t.Cleanup(func() { darwinSettingsPath = previous })
	data := []byte(`{"version":4,"variantID":"ferret_sable","coatMode":1,"petCount":3,"selectedCoatIDs":["ferret_sable","chinchilla_standard_gray","rabbit_gray"],"petSizes":[70,100,120]}`)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	a := &darwinPetApp{selectedCoats: defaultDarwinSelectedCoats()}
	a.loadSettings()
	if a.variantID(a.variant) != "ferret_sable" {
		t.Fatal("stable variant ID was ignored without its numeric index")
	}
	for i, id := range []string{"ferret_sable", "chinchilla_standard_gray", "rabbit_gray"} {
		if a.variantID(a.selectedCoats[i]) != id {
			t.Fatalf("slot %d lost stable ID", i)
		}
	}
	a.saveSettings()
	b := &darwinPetApp{selectedCoats: defaultDarwinSelectedCoats()}
	b.loadSettings()
	if a.variant != b.variant || a.selectedCoats != b.selectedCoats || a.petSizes != b.petSizes {
		t.Fatal("selection or size changed after saving")
	}
}
