//go:build windows

package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"animals-desktop/internal/catalog"
	"github.com/lxn/win"
)

func TestHorizontalMotionFramesUseStableRightFacingWalkSequence(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	states := []behaviorState{
		stateWalk,
		stateScurry,
		stateWheel,
		stateForage,
		stateCarry,
	}

	for _, state := range states {
		for frame := 0; frame < 64; frame++ {
			got := currentFrame(state, frame)
			if !allowed[got] {
				t.Fatalf("currentFrame(%v, %d) = %d, want stable right-facing walk frame", state, frame, got)
			}
		}
	}
}

func TestScurryUsesStableFastWalkFrames(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	for frame := 0; frame < scurryFrames*2; frame++ {
		got := currentFrame(stateScurry, frame)
		if !allowed[got] {
			t.Fatalf("currentFrame(stateScurry, %d) = %d, want stable fast walk frame", frame, got)
		}
	}
}

func TestWheelUsesStableWalkFramesForAnimalRuntime(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	for frame := 0; frame < wheelRunFrames*2; frame++ {
		got := currentFrame(stateWheel, frame)
		if !allowed[got] {
			t.Fatalf("currentFrame(stateWheel, %d) = %d, want stable walk frame", frame, got)
		}
	}
}

func TestWeakNibbleVariantsUseStableActionFallback(t *testing.T) {
	for _, variant := range []coatVariant{
		{ID: "sugar_glider_gray"},
		{ID: "rabbit_chestnut_agouti"},
	} {
		for frame := 0; frame < 32; frame++ {
			got := currentFrameForVariant(stateNibble, frame, variant)
			if got < hopStart || got >= hopStart+4 {
				t.Fatalf("%s nibble frame %d = %d, want stable action fallback", variant.ID, frame, got)
			}
			got = currentFrameForVariant(stateGroom, frame, variant)
			if got < groomStart || got >= groomStart+groomFrames {
				t.Fatalf("%s groom frame %d = %d, want groom fallback", variant.ID, frame, got)
			}
		}
	}

	hamster := coatVariant{ID: "hamster_golden_syrian"}
	for frame := 0; frame < 12; frame++ {
		got := currentFrameForVariant(stateNibble, frame, hamster)
		if got < nibbleStart || got >= nibbleStart+3 {
			t.Fatalf("hamster nibble frame %d = %d, want original nibble frames", frame, got)
		}
	}
}

func TestFerretMotionProfileUsesAcceptedFrameRanges(t *testing.T) {
	ferret := coatVariant{SpeciesID: "ferret"}
	tests := []struct {
		name    string
		states  []behaviorState
		start   int
		count   int
		divisor int
		clamped bool
	}{
		{name: "idle", states: []behaviorState{stateIdle}, start: 0, count: 4, divisor: 5},
		{name: "slink", states: []behaviorState{stateWalk, stateForage, stateCarry}, start: 4, count: 8, divisor: 2},
		{name: "scurry", states: []behaviorState{stateScurry, stateWheel}, start: 12, count: 8, divisor: 1},
		{name: "sniff", states: []behaviorState{stateNibble, stateEat, stateDig}, start: 20, count: 6, divisor: 3},
		{name: "groom", states: []behaviorState{stateGroom, stateFaceGroom}, start: 26, count: 6, divisor: 3},
		{name: "turn", states: []behaviorState{stateTurn}, start: 32, count: 8, divisor: 2, clamped: true},
		{name: "creep", states: []behaviorState{stateHop}, start: 40, count: 8, divisor: 2},
		{name: "rest", states: []behaviorState{stateRest}, start: 48, count: 8, divisor: 4},
		{name: "alert", states: []behaviorState{stateStand}, start: 56, count: 6, divisor: 4},
	}

	for _, tt := range tests {
		for _, state := range tt.states {
			t.Run(fmt.Sprintf("%s/state-%d", tt.name, state), func(t *testing.T) {
				seen := make(map[int]bool)
				for frame := 0; frame < tt.count*tt.divisor; frame++ {
					got := currentFrameForVariant(state, frame, ferret)
					if got < tt.start || got >= tt.start+tt.count {
						t.Fatalf("frame %d = %d, want %d..%d", frame, got, tt.start, tt.start+tt.count-1)
					}
					seen[got] = true
				}
				if len(seen) != tt.count {
					t.Fatalf("visited %d distinct frames, want all %d in %d..%d", len(seen), tt.count, tt.start, tt.start+tt.count-1)
				}
				if tt.clamped {
					if got := currentFrameForVariant(state, 999, ferret); got != tt.start+tt.count-1 {
						t.Fatalf("late frame = %d, want clamped final frame %d", got, tt.start+tt.count-1)
					}
				}
			})
		}
	}
}

func TestFerretMovingStatesNeverUseStationaryActionFrames(t *testing.T) {
	ferret := coatVariant{MotionProfile: catalog.MotionProfileFerretSlink}
	for _, state := range []behaviorState{
		stateWalk,
		stateScurry,
		stateWheel,
		stateForage,
		stateCarry,
	} {
		for frame := 0; frame < 128; frame++ {
			got := currentFrameForVariant(state, frame, ferret)
			if got < walkStart || got >= nibbleStart {
				t.Fatalf("moving state %d frame %d = %d, want slink/scurry range %d..%d", state, frame, got, walkStart, nibbleStart-1)
			}
		}
	}
	for frame := 0; frame < 128; frame++ {
		got := currentFrameForVariant(stateHop, frame, ferret)
		if got < eatStart || got >= standStart {
			t.Fatalf("moving creep frame %d = %d, want creep range %d..%d", frame, got, eatStart, standStart-1)
		}
	}
}

func TestFerretRandomActionSelectionCanReachRest(t *testing.T) {
	ferretCounts := make(map[behaviorState]int)
	defaultCounts := make(map[behaviorState]int)
	for roll := 0; roll < 100; roll++ {
		ferretCounts[randomActionState(catalog.MotionProfileFerretSlink, roll)]++
		defaultCounts[randomActionState(catalog.MotionProfileDegu, roll)]++
	}

	if got := ferretCounts[stateRest]; got != 2 {
		t.Fatalf("ferret rest rolls = %d, want 2", got)
	}
	if got := defaultCounts[stateRest]; got != 0 {
		t.Fatalf("default profile rest rolls = %d, want 0", got)
	}
	wantDefault := map[behaviorState]int{
		stateIdle:      30,
		stateWalk:      40,
		stateScurry:    14,
		stateNibble:    6,
		stateStand:     4,
		stateFaceGroom: 4,
		stateHop:       2,
	}
	for state, want := range wantDefault {
		if got := defaultCounts[state]; got != want {
			t.Fatalf("default state %d rolls = %d, want %d", state, got, want)
		}
	}
}

func TestRestIsStationaryAndCompletesSafely(t *testing.T) {
	a := &petApp{
		mode:     modeKeyboard,
		sceneW:   800,
		petSizes: defaultPetSizes(),
	}
	p := desktopPet{
		x:          120,
		variant:    0,
		item:       noItem,
		carryKind:  noItem,
		state:      stateRest,
		stateTicks: 0,
		moveSpeed:  9,
		dir:        1,
	}

	a.tickPet(0, &p)

	if p.x != 120 {
		t.Fatalf("rest completion moved pet to x=%d, want 120", p.x)
	}
	if p.state != stateIdle || p.moveSpeed != 0 {
		t.Fatalf("rest completion = state %d speed %d, want idle and stationary", p.state, p.moveSpeed)
	}
}

func TestDefaultMotionProfileFrameMappingIsUnchanged(t *testing.T) {
	degu := coatVariant{SpeciesID: "degu"}
	for _, state := range []behaviorState{
		stateIdle,
		stateWalk,
		stateScurry,
		stateNibble,
		stateHop,
		stateWheel,
		stateGroom,
		stateForage,
		stateCarry,
		stateTurn,
		stateEat,
		stateDig,
		stateStand,
		stateFaceGroom,
	} {
		for frame := 0; frame < 64; frame++ {
			got := currentFrameForVariant(state, frame, degu)
			want := currentFrame(state, frame)
			if got != want {
				t.Fatalf("state %d frame %d = %d, want legacy frame %d", state, frame, got, want)
			}
		}
	}
}

func TestFrameFromSeqHandlesEmptyAndBadDivisor(t *testing.T) {
	if got := frameFromSeq(nil, 12, 2); got != idleStart {
		t.Fatalf("frameFromSeq(nil) = %d, want %d", got, idleStart)
	}
	seq := []int{7, 9}
	if got := frameFromSeq(seq, 3, 0); got != 9 {
		t.Fatalf("frameFromSeq with zero divisor = %d, want 9", got)
	}
}

func TestFrameFromSeqClampedHoldsFinalFrame(t *testing.T) {
	seq := []int{7, 9, 11}
	if got := frameFromSeqClamped(seq, 999, 2); got != 11 {
		t.Fatalf("frameFromSeqClamped past end = %d, want 11", got)
	}
	if got := frameFromSeqClamped(seq, 3, 0); got != 11 {
		t.Fatalf("frameFromSeqClamped with zero divisor = %d, want 11", got)
	}
}

func TestTypingStartsAndExtendsWheelOnlyInKeyboardMode(t *testing.T) {
	// The legacy family retains its wheel. Timed assets have separate coverage.
	t.Setenv("ANIMALSDESKTOP_MOTIONS", "legacy")
	a := &petApp{
		mode:         modeKeyboard,
		wheelEnabled: true,
		wheelX:       400,
		sceneW:       1200,
		speed:        3,
		pets: []desktopPet{
			{state: stateWalk, stateTicks: 12, item: noItem},
			{state: stateWalk, stateTicks: 12, item: noItem},
		},
	}

	a.onTyping()
	if got := a.pets[0].state; got != stateWheel {
		t.Fatalf("first pet state = %v, want stateWheel", got)
	}
	if got := a.pets[0].stateTicks; got != wheelKeyHold {
		t.Fatalf("wheel hold ticks = %d, want %d", got, wheelKeyHold)
	}
	if got := a.pets[0].moveSpeed; got != 0 {
		t.Fatalf("wheel pet moveSpeed = %d, want 0", got)
	}
	wantX := clamp(a.wheelX-wheelSize/2, 0, max(0, a.sceneW-spriteW))
	if got := a.pets[0].x; got != wantX {
		t.Fatalf("wheel pet x = %d, want %d", got, wantX)
	}
	if got := a.pets[1].state; got != stateScurry {
		t.Fatalf("second pet state = %v, want stateScurry", got)
	}

	a.pets[0].frame = 7
	a.pets[0].stateTicks = 3
	a.onTyping()
	if got := a.pets[0].frame; got != 7 {
		t.Fatalf("wheel frame reset while extending: got %d, want 7", got)
	}
	if got := a.pets[0].stateTicks; got != wheelKeyHold {
		t.Fatalf("extended wheel hold ticks = %d, want %d", got, wheelKeyHold)
	}
}

func TestTypingDoesNotStartWheelInRandomMode(t *testing.T) {
	a := &petApp{
		mode:         modeRandom,
		wheelEnabled: true,
		wheelX:       400,
		sceneW:       1200,
		pets: []desktopPet{
			{state: stateWalk, stateTicks: 12, item: noItem},
		},
	}

	a.onTyping()
	if got := a.pets[0].state; got == stateWheel {
		t.Fatalf("typing in random mode started wheel state")
	}
}

func TestRuntimeCatalogMatchesCatalogSource(t *testing.T) {
	want := catalog.RuntimeVariants()
	if got := len(variants); got != len(want) {
		t.Fatalf("Windows runtime variants = %d, want %d catalog variants", got, len(want))
	}
	for i, variant := range variants {
		if variant != want[i] {
			t.Fatalf("Windows runtime variant[%d] = %+v, want catalog variant %+v", i, variant, want[i])
		}
		if variant.SpeciesID == "degu" {
			t.Fatalf("runtime variants include degu: %+v", variant)
		}
	}
}

func TestSettingsVariantSelectionCoversEveryRuntimeVariant(t *testing.T) {
	a := &petApp{
		pets:          make([]desktopPet, maxPetCount),
		selectedCoats: defaultSelectedCoats(),
	}

	for i, variant := range variants {
		a.lang = langJapanese
		if got := a.variantLabel(i); got == "" {
			t.Fatalf("Japanese label for %q is empty", variant.ID)
		}
		a.lang = langEnglish
		if got := a.variantLabel(i); got == "" {
			t.Fatalf("English label for %q is empty", variant.ID)
		}

		a.setCoatMode(coatFixed)
		a.setFixedVariant(i)
		if got := a.variant; got != i {
			t.Fatalf("fixed variant index = %d, want %d", got, i)
		}
		for petIndex, pet := range a.pets {
			if pet.variant != i {
				t.Fatalf("fixed variant %q pet %d = %d, want %d", variant.ID, petIndex, pet.variant, i)
			}
		}

		a.setCoatMode(coatSelected)
		petIndex := i % maxPetCount
		a.setSelectedVariant(petIndex, i)
		if a.selectedCoats[petIndex] != i || a.pets[petIndex].variant != i {
			t.Fatalf("selected variant %q pet %d = selected:%d pet:%d, want %d", variant.ID, petIndex, a.selectedCoats[petIndex], a.pets[petIndex].variant, i)
		}
	}
}

func TestSettingsRoundTripPersistsCoreOptions(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)

	a := &petApp{
		variant:                4,
		coatMode:               coatSelected,
		selectedCoats:          [maxPetCount]int{1, 3, 5, 7, 9, 0, 2, 4, 6, 8},
		selectedSlotModes:      [maxPetCount]petSlotMode{petSlotFixed, petSlotRandom, petSlotFixed, petSlotRandom},
		selectedRandomGroupIDs: [maxPetCount]string{"", "rabbit", "", "chinchilla"},
		randomGroupID:          "rabbit",
		petSizes:               [maxPetCount]int{80, 90, 100, 110, 120, 70, 80, 90, 100, 110},
		petNames:               [maxPetCount]string{"モカ", "Sora", "  Nagi  ", "", "", "", "", "", "", ""},
		nameLabels:             true,
		speed:                  5,
		mode:                   modeKeyboard,
		petCount:               10,
		wheelEnabled:           false,
		bidirectional:          false,
		positionMode:           positionScreenBottom,
		overlayOffsetY:         24,
		displayIndex:           0,
		displayScope:           displayScopeSingle,
		displaySpanEnd:         0,
		walkRangeStart:         15,
		walkRangeEnd:           85,
		lang:                   langEnglish,
		settingsX:              220,
		settingsY:              180,
	}
	if err := a.saveSettings(); err != nil {
		t.Fatalf("saveSettings() error = %v", err)
	}

	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("settings file was not written: %v", err)
	}
	var saved appSettings
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("settings json is invalid: %v", err)
	}
	if saved.Version != settingsVersion || saved.PetCount != 10 || saved.Mode != int(modeKeyboard) {
		t.Fatalf("saved settings = %+v, want version %d petCount 10 keyboard mode", saved, settingsVersion)
	}
	if saved.VariantID != variantIDAt(a.variant) {
		t.Fatalf("saved VariantID = %q, want %q", saved.VariantID, variantIDAt(a.variant))
	}
	if len(saved.SelectedCoatIDs) != maxPetCount || saved.SelectedCoatIDs[1] != variantIDAt(3) {
		t.Fatalf("saved SelectedCoatIDs = %#v", saved.SelectedCoatIDs)
	}
	if len(saved.SelectedSlotModes) != maxPetCount || saved.SelectedSlotModes[1] != int(petSlotRandom) || saved.SelectedSlotModes[2] != int(petSlotFixed) {
		t.Fatalf("saved SelectedSlotModes = %#v", saved.SelectedSlotModes)
	}
	if len(saved.SelectedRandomGroupIDs) != maxPetCount || saved.SelectedRandomGroupIDs[1] != "rabbit" || saved.SelectedRandomGroupIDs[3] != "chinchilla" {
		t.Fatalf("saved SelectedRandomGroupIDs = %#v", saved.SelectedRandomGroupIDs)
	}
	if saved.RandomGroupID != "rabbit" {
		t.Fatalf("saved RandomGroupID = %q, want rabbit", saved.RandomGroupID)
	}
	if saved.Language != int(langEnglish) {
		t.Fatalf("saved Language = %d, want English", saved.Language)
	}
	if !saved.NameLabels {
		t.Fatalf("saved NameLabels = false, want true")
	}
	if saved.PositionMode == nil || *saved.PositionMode != int(positionScreenBottom) {
		t.Fatalf("saved PositionMode = %#v, want screen bottom", saved.PositionMode)
	}
	if saved.VerticalOffset == nil || *saved.VerticalOffset != 24 {
		t.Fatalf("saved VerticalOffset = %#v, want 24", saved.VerticalOffset)
	}
	if saved.DisplayScope == nil || *saved.DisplayScope != int(displayScopeSingle) || saved.DisplayIndex == nil || *saved.DisplayIndex != 0 || saved.DisplaySpanEnd == nil || *saved.DisplaySpanEnd != 0 {
		t.Fatalf("saved display selection = scope:%#v index:%#v span:%#v", saved.DisplayScope, saved.DisplayIndex, saved.DisplaySpanEnd)
	}
	if saved.WalkRangeStart == nil || *saved.WalkRangeStart != 15 || saved.WalkRangeEnd == nil || *saved.WalkRangeEnd != 85 {
		t.Fatalf("saved walk range = start:%#v end:%#v", saved.WalkRangeStart, saved.WalkRangeEnd)
	}
	if got := saved.PetNames[0]; got != "モカ" {
		t.Fatalf("saved pet name 0 = %q, want モカ", got)
	}
	if got := saved.PetNames[2]; got != "Nagi" {
		t.Fatalf("saved pet name 2 = %q, want sanitized Nagi", got)
	}
	if len(saved.PetSizes) != maxPetCount || saved.PetSizes[0] != 80 || saved.PetSizes[4] != 120 || saved.PetSizes[5] != 70 {
		t.Fatalf("saved pet sizes = %#v", saved.PetSizes)
	}

	b := &petApp{
		variant:        0,
		coatMode:       coatSelected,
		selectedCoats:  defaultSelectedCoats(),
		petSizes:       defaultPetSizes(),
		speed:          3,
		mode:           modeRandom,
		petCount:       2,
		wheelEnabled:   true,
		bidirectional:  true,
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: defaultOverlayOffsetY,
		displayScope:   displayScopeSingle,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
		lang:           langJapanese,
		settingsX:      120,
		settingsY:      120,
	}
	if err := b.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if b.variant != 4 || b.coatMode != a.coatMode || b.speed != a.speed || b.mode != a.mode || b.petCount != a.petCount {
		t.Fatalf("loaded scalar settings = variant:%d coat:%d speed:%d mode:%d count:%d", b.variant, b.coatMode, b.speed, b.mode, b.petCount)
	}
	if b.wheelEnabled != a.wheelEnabled || b.bidirectional != a.bidirectional || b.lang != a.lang {
		t.Fatalf("loaded flags = wheel:%v bidirectional:%v lang:%d", b.wheelEnabled, b.bidirectional, b.lang)
	}
	if b.nameLabels != a.nameLabels {
		t.Fatalf("loaded nameLabels = %v, want %v", b.nameLabels, a.nameLabels)
	}
	if b.positionMode != a.positionMode || b.overlayOffsetY != a.overlayOffsetY || b.displayScope != a.displayScope || b.displayIndex != a.displayIndex || b.displaySpanEnd != a.displaySpanEnd {
		t.Fatalf("loaded display settings = mode:%d offset:%d scope:%d index:%d span:%d", b.positionMode, b.overlayOffsetY, b.displayScope, b.displayIndex, b.displaySpanEnd)
	}
	if b.walkRangeStart != a.walkRangeStart || b.walkRangeEnd != a.walkRangeEnd {
		t.Fatalf("loaded walk range = %d-%d, want %d-%d", b.walkRangeStart, b.walkRangeEnd, a.walkRangeStart, a.walkRangeEnd)
	}
	wantCoats := [maxPetCount]int{1, 3, 5, 7, 9, 0, 2, 4, 6, 8}
	for i := 0; i < maxPetCount; i++ {
		if b.selectedCoats[i] != wantCoats[i] {
			t.Fatalf("selectedCoats[%d] = %d, want %d", i, b.selectedCoats[i], wantCoats[i])
		}
	}
	if b.selectedSlotModes[1] != petSlotRandom || b.selectedSlotModes[2] != petSlotFixed {
		t.Fatalf("loaded selectedSlotModes = %#v", b.selectedSlotModes)
	}
	if b.selectedRandomGroupIDs[1] != "rabbit" || b.selectedRandomGroupIDs[3] != "chinchilla" {
		t.Fatalf("loaded selectedRandomGroupIDs = %#v", b.selectedRandomGroupIDs)
	}
	if b.randomGroupID != "rabbit" {
		t.Fatalf("loaded randomGroupID = %q, want rabbit", b.randomGroupID)
	}
	if b.petNames[0] != "モカ" || b.petNames[1] != "Sora" || b.petNames[2] != "Nagi" {
		t.Fatalf("loaded pet names = %#v", b.petNames[:3])
	}
	for i, want := range a.petSizes {
		if got := b.petSizes[i]; got != want {
			t.Fatalf("loaded petSizes[%d] = %d, want %d", i, got, want)
		}
	}
}

func TestSettingsRoundTripUsesVariantIDsWhenRosterOrderChanges(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)
	shoebillIndex, ok := variantIndexByID("shoebill_stork")
	if !ok {
		t.Fatal("shoebill_stork missing from runtime variants")
	}
	catIndex, ok := variantIndexByID("domestic_shorthair_tabby_white_stocky")
	if !ok {
		t.Fatal("domestic_shorthair_tabby_white_stocky missing from runtime variants")
	}
	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir settings dir: %v", err)
	}
	data := []byte(fmt.Sprintf(`{
  "version": %d,
  "variant": %d,
  "variantID": "shoebill_stork",
  "coatMode": %d,
  "selectedCoats": [%d, %d],
  "selectedCoatIDs": ["domestic_shorthair_tabby_white_stocky", "shoebill_stork"],
  "petCount": 2,
  "speed": 5,
  "mode": 1,
  "wheelEnabled": true,
  "bidirectional": true,
  "language": 0
}`, settingsVersion, catIndex, int(coatSelected), catIndex, catIndex))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &petApp{
		selectedCoats: defaultSelectedCoats(),
		petSizes:      defaultPetSizes(),
		speed:         3,
		mode:          modeRandom,
		petCount:      1,
		wheelEnabled:  true,
		bidirectional: true,
	}
	if err := a.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if a.variant != shoebillIndex {
		t.Fatalf("variant = %d (%s), want shoebill index %d", a.variant, variantIDAt(a.variant), shoebillIndex)
	}
	if a.selectedCoats[0] != catIndex || a.selectedCoats[1] != shoebillIndex {
		t.Fatalf("selectedCoats = %v, want cat then shoebill", a.selectedCoats[:2])
	}
}

func TestLegacySettingsRecoverShoebillIndexDrift(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)
	shoebillIndex, ok := variantIndexByID("shoebill_stork")
	if !ok {
		t.Fatal("shoebill_stork missing from runtime variants")
	}
	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir settings dir: %v", err)
	}
	data := []byte(`{
  "version": 2,
  "variant": 41,
  "coatMode": 1,
  "selectedCoats": [41, 0],
  "petNames": ["爆速ハシビロコウさん", ""],
  "petCount": 2,
  "speed": 5,
  "mode": 1,
  "wheelEnabled": true,
  "bidirectional": true,
  "language": 0
}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &petApp{
		selectedCoats: defaultSelectedCoats(),
		petSizes:      defaultPetSizes(),
		speed:         3,
		mode:          modeRandom,
		petCount:      1,
		wheelEnabled:  true,
		bidirectional: true,
	}
	if err := a.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if a.variant != shoebillIndex || a.selectedCoats[0] != shoebillIndex {
		t.Fatalf("legacy shoebill migration = variant:%d selected:%d ids:%s/%s, want %d", a.variant, a.selectedCoats[0], variantIDAt(a.variant), variantIDAt(a.selectedCoats[0]), shoebillIndex)
	}
	if got := a.petNames[0]; got != "爆速ハシビロコウさん" {
		t.Fatalf("pet name = %q", got)
	}
}

func TestLegacySettingsKeepAmbiguousIndexWithoutShoebillName(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)
	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir settings dir: %v", err)
	}
	data := []byte(`{
  "version": 2,
  "variant": 41,
  "coatMode": 1,
  "selectedCoats": [41],
  "petNames": [""],
  "petCount": 1,
  "speed": 5,
  "mode": 1,
  "wheelEnabled": true,
  "bidirectional": true,
  "language": 0
}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &petApp{
		selectedCoats: defaultSelectedCoats(),
		petSizes:      defaultPetSizes(),
		speed:         3,
		mode:          modeRandom,
		petCount:      1,
		wheelEnabled:  true,
		bidirectional: true,
	}
	if err := a.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if a.variant != 41 || a.selectedCoats[0] != 41 {
		t.Fatalf("ambiguous legacy index changed to variant:%d selected:%d, want raw index 41", a.variant, a.selectedCoats[0])
	}
}

func TestVariantDisplayLabelsIncludeAnimalGroup(t *testing.T) {
	a := &petApp{lang: langJapanese}
	shoebillIndex, ok := variantIndexByID("shoebill_stork")
	if !ok {
		t.Fatal("shoebill_stork missing from runtime variants")
	}
	if got := a.variantDisplayLabel(shoebillIndex); got != "鳥 / ハシビロコウ" {
		t.Fatalf("shoebill display label = %q", got)
	}
	a.lang = langEnglish
	if got := a.variantDisplayLabel(shoebillIndex); got != "Birds / Shoebill" {
		t.Fatalf("English shoebill display label = %q", got)
	}
}

func TestSettingsLanguageLabelsSwitchToEnglish(t *testing.T) {
	a := &petApp{lang: langEnglish}
	if got := a.txt("settingsTitle"); got != "Animals Desktop Settings" {
		t.Fatalf("English settings title = %q", got)
	}
	if got := a.txt("language"); got != "Language" {
		t.Fatalf("English language label = %q", got)
	}
	if got := a.settingsButtonLabel(ctrlLanguageCombo); got != "English" {
		t.Fatalf("English language button = %q", got)
	}

	a.lang = langJapanese
	if got := a.txt("language"); got != "Language" {
		t.Fatalf("Japanese language label = %q, want Language", got)
	}
	if got := a.settingsButtonLabel(ctrlLanguageCombo); got == "English" || got == "" {
		t.Fatalf("Japanese language button should be a non-English label, got %q", got)
	}
}

func TestWindowsDefaultAppVersionTracksCurrentRelease(t *testing.T) {
	if appVersion != "v2.0.0" {
		t.Fatalf("appVersion = %q, want v2.0.0", appVersion)
	}
}

func TestTrayMenuLanguageCommandPersistsSelection(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)

	a := &petApp{
		coatMode:       coatSelected,
		selectedCoats:  defaultSelectedCoats(),
		petSizes:       defaultPetSizes(),
		speed:          3,
		mode:           modeRandom,
		petCount:       1,
		bidirectional:  true,
		positionMode:   positionTaskbarEdge,
		displayScope:   displayScopeSingle,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
		lang:           langJapanese,
	}

	a.handleMenu(menuLangEnglish)
	if a.lang != langEnglish {
		t.Fatalf("language after English menu command = %d, want English", a.lang)
	}
	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("settings file after English menu command was not written: %v", err)
	}
	var saved appSettings
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("settings json after English menu command is invalid: %v", err)
	}
	if saved.Language != int(langEnglish) {
		t.Fatalf("saved Language after English menu command = %d, want English", saved.Language)
	}

	a.handleMenu(menuLangJapanese)
	if a.lang != langJapanese {
		t.Fatalf("language after Japanese menu command = %d, want Japanese", a.lang)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("settings file after Japanese menu command was not written: %v", err)
	}
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("settings json after Japanese menu command is invalid: %v", err)
	}
	if saved.Language != int(langJapanese) {
		t.Fatalf("saved Language after Japanese menu command = %d, want Japanese", saved.Language)
	}
}

func TestTrayMenuTemporaryHideCommandDoesNotPersist(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)

	a := &petApp{lang: langJapanese}
	if got := a.temporaryVisibilityLabel(); got != "一時的に非表示" {
		t.Fatalf("visible temporary label = %q, want hide label", got)
	}

	a.handleMenu(menuHideToggle)
	if !a.overlayHidden {
		t.Fatalf("temporary hide menu command should hide overlay")
	}
	if got := a.temporaryVisibilityLabel(); got != "表示する" {
		t.Fatalf("hidden temporary label = %q, want show label", got)
	}
	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("temporary hide command should not persist settings")
	} else if !os.IsNotExist(err) {
		t.Fatalf("checking settings file after temporary hide: %v", err)
	}

	a.lang = langEnglish
	if got := a.temporaryVisibilityLabel(); got != "Show" {
		t.Fatalf("hidden English temporary label = %q, want Show", got)
	}
	a.handleMenu(menuHideToggle)
	if a.overlayHidden {
		t.Fatalf("second temporary hide menu command should show overlay")
	}
	if got := a.temporaryVisibilityLabel(); got != "Hide temporarily" {
		t.Fatalf("visible English temporary label = %q, want Hide temporarily", got)
	}
}

func TestVersionOneSettingsKeepOptionsButResetOldAnimalSelection(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)
	dir := filepath.Join(configRoot, settingsDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir settings dir: %v", err)
	}
	path := filepath.Join(dir, settingsFileName)
	data := []byte(`{
  "version": 1,
  "variant": 8,
  "coatMode": 2,
  "selectedCoats": [8, 7, 6, 5, 4, 3, 2, 1, 0, 9],
  "speed": 5,
  "mode": 0,
  "petCount": 4,
  "wheelEnabled": false,
  "bidirectional": false,
  "language": 1,
  "nameLabels": true,
  "petNames": ["モカ"]
}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &petApp{
		variant:       0,
		coatMode:      coatSelected,
		selectedCoats: defaultSelectedCoats(),
		speed:         3,
		mode:          modeRandom,
		petCount:      5,
		wheelEnabled:  true,
		bidirectional: true,
		lang:          langJapanese,
	}
	if err := a.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if a.speed != 5 || a.mode != modeKeyboard || a.petCount != 4 || a.wheelEnabled || a.bidirectional || a.lang != langEnglish {
		t.Fatalf("loaded preserved settings = speed:%d mode:%d count:%d wheel:%v bidi:%v lang:%d", a.speed, a.mode, a.petCount, a.wheelEnabled, a.bidirectional, a.lang)
	}
	if a.variant != 0 || a.coatMode != coatSelected || a.selectedCoats != defaultSelectedCoats() {
		t.Fatalf("old animal selection was not reset: variant:%d coat:%d selected:%v", a.variant, a.coatMode, a.selectedCoats)
	}
	if !a.nameLabels || a.petNames[0] != "モカ" {
		t.Fatalf("loaded name settings = labels:%v names:%v", a.nameLabels, a.petNames[:1])
	}
	for i, size := range a.petSizes {
		if size != defaultPetSizePercent {
			t.Fatalf("legacy petSizes[%d] = %d, want default %d", i, size, defaultPetSizePercent)
		}
	}
}

func TestNormalizeWalkRangeKeepsMinimumSpan(t *testing.T) {
	start, end := normalizeWalkRange(48, 52)
	if end-start != minWalkRangeSpan {
		t.Fatalf("normalizeWalkRange narrow span = %d-%d, want %d point span", start, end, minWalkRangeSpan)
	}
	start, end = normalizeWalkRange(95, 10)
	if start != 10 || end != 95 {
		t.Fatalf("normalizeWalkRange reversed = %d-%d, want 10-95", start, end)
	}
}

func TestOverlayRectForAppliesScreenBottomOffsetAndWalkRange(t *testing.T) {
	a := &petApp{
		positionMode:   positionScreenBottom,
		overlayOffsetY: 24,
		walkRangeStart: 25,
		walkRangeEnd:   75,
	}
	work := win.RECT{Left: 0, Top: 0, Right: 1000, Bottom: 760}
	screen := win.RECT{Left: 0, Top: 0, Right: 1000, Bottom: 800}
	got := a.overlayRectFor(work, screen)
	if got.Left != 250 || got.Right != 750 {
		t.Fatalf("overlay x range = %d-%d, want 250-750", got.Left, got.Right)
	}
	if got.Top != int32(800-sceneH) || got.Bottom != 800 {
		t.Fatalf("overlay y range = %d-%d, want clamped to screen bottom", got.Top, got.Bottom)
	}
}

func TestOverlayRectForTaskbarOffsetStaysInsideScreen(t *testing.T) {
	a := &petApp{
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: -20,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
	}
	work := win.RECT{Left: 100, Top: 50, Right: 900, Bottom: 700}
	screen := win.RECT{Left: 0, Top: 0, Right: 1000, Bottom: 800}
	got := a.overlayRectFor(work, screen)
	if got.Left != work.Left || got.Right != work.Right {
		t.Fatalf("overlay x = %d-%d, want work area %d-%d", got.Left, got.Right, work.Left, work.Right)
	}
	wantTop := int32(700 - sceneH - 20)
	if got.Top != wantTop || got.Bottom != wantTop+sceneH {
		t.Fatalf("overlay y = %d-%d, want %d-%d", got.Top, got.Bottom, wantTop, wantTop+sceneH)
	}
}

func TestOverlaySegmentsUseEachMonitorBottomForMixedHeightSpan(t *testing.T) {
	a := &petApp{
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: defaultOverlayOffsetY,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
	}
	areas := []displayArea{
		{
			Work:    win.RECT{Left: 0, Top: 0, Right: 3840, Bottom: 2080},
			Screen:  win.RECT{Left: 0, Top: 0, Right: 3840, Bottom: 2160},
			Primary: true,
		},
		{
			Work:   win.RECT{Left: 3840, Top: 0, Right: 5760, Bottom: 1040},
			Screen: win.RECT{Left: 3840, Top: 0, Right: 5760, Bottom: 1080},
		},
	}
	combined := combineDisplayAreas(areas)
	overlay := a.overlayRectFor(combined.Work, combined.Screen)
	segments := a.overlaySegmentsForAreas(areas, overlay)
	if len(segments) != 2 {
		t.Fatalf("segments = %d, want 2: %+v", len(segments), segments)
	}
	if segments[0].SceneLeft != 0 || segments[0].SceneRight != 3840 ||
		segments[1].SceneLeft != 3840 || segments[1].SceneRight != 5760 {
		t.Fatalf("scene segments = %+v, want 0-3840 and 3840-5760", segments)
	}
	wantMainTop := int32(2080 - sceneH + defaultOverlayOffsetY)
	wantSubTop := int32(1040 - sceneH + defaultOverlayOffsetY)
	if segments[0].Rect.Top != wantMainTop || segments[0].Rect.Bottom != wantMainTop+sceneH {
		t.Fatalf("main segment y = %d-%d, want %d-%d", segments[0].Rect.Top, segments[0].Rect.Bottom, wantMainTop, wantMainTop+sceneH)
	}
	if segments[1].Rect.Top != wantSubTop || segments[1].Rect.Bottom != wantSubTop+sceneH {
		t.Fatalf("secondary segment y = %d-%d, want %d-%d", segments[1].Rect.Top, segments[1].Rect.Bottom, wantSubTop, wantSubTop+sceneH)
	}
	if segments[0].Rect.Top == segments[1].Rect.Top {
		t.Fatalf("mixed-height displays should not share one overlay top: %+v", segments)
	}
}

func TestOverlaySegmentsClampScreenBottomPerMonitor(t *testing.T) {
	a := &petApp{
		positionMode:   positionScreenBottom,
		overlayOffsetY: 48,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
	}
	areas := []displayArea{
		{
			Work:    win.RECT{Left: 0, Top: 0, Right: 3840, Bottom: 2080},
			Screen:  win.RECT{Left: 0, Top: 0, Right: 3840, Bottom: 2160},
			Primary: true,
		},
		{
			Work:   win.RECT{Left: 3840, Top: 0, Right: 5760, Bottom: 1040},
			Screen: win.RECT{Left: 3840, Top: 0, Right: 5760, Bottom: 1080},
		},
	}
	combined := combineDisplayAreas(areas)
	overlay := a.overlayRectFor(combined.Work, combined.Screen)
	segments := a.overlaySegmentsForAreas(areas, overlay)
	if len(segments) != 2 {
		t.Fatalf("segments = %d, want 2: %+v", len(segments), segments)
	}
	if got, want := segments[0].Rect.Bottom, int32(2160); got != want {
		t.Fatalf("main bottom = %d, want %d", got, want)
	}
	if got, want := segments[1].Rect.Bottom, int32(1080); got != want {
		t.Fatalf("secondary bottom = %d, want %d", got, want)
	}
}

func TestOverlaySegmentsScaleHeightForHighDPIMonitor(t *testing.T) {
	a := &petApp{
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: defaultOverlayOffsetY,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
	}
	areas := []displayArea{
		{
			Work:    win.RECT{Left: -2880, Top: 0, Right: -960, Bottom: 1032},
			Screen:  win.RECT{Left: -2880, Top: 0, Right: -960, Bottom: 1080},
			DPI:     144,
			Primary: false,
		},
	}
	overlay := a.overlayRectFor(areas[0].Work, areas[0].Screen)
	segments := a.overlaySegmentsForAreas(areas, overlay)
	if len(segments) != 1 {
		t.Fatalf("segments = %d, want 1: %+v", len(segments), segments)
	}
	wantH := int32(scaleForDPI(sceneH, 144))
	if got := segments[0].Rect.Bottom - segments[0].Rect.Top; got != wantH {
		t.Fatalf("high-DPI segment height = %d, want %d", got, wantH)
	}
	wantTop := int32(1032 - int(wantH) + scaleForDPI(defaultOverlayOffsetY, 144))
	if segments[0].Rect.Top != wantTop {
		t.Fatalf("high-DPI segment top = %d, want %d", segments[0].Rect.Top, wantTop)
	}
}

func TestRenderOverlaySegmentScalesCanvasForHighDPI(t *testing.T) {
	a := &petApp{
		sceneW:   1920,
		petSizes: defaultPetSizes(),
	}
	segment := overlaySegment{
		Rect:       win.RECT{Left: 0, Top: 0, Right: 1920, Bottom: int32(scaleForDPI(sceneH, 144))},
		SceneLeft:  0,
		SceneRight: 1920,
		DPI:        144,
	}
	canvas := a.renderOverlaySegment(segment)
	if canvas.Bounds().Dx() != 1920 || canvas.Bounds().Dy() != scaleForDPI(sceneH, 144) {
		t.Fatalf("canvas bounds = %v, want 1920x%d", canvas.Bounds(), scaleForDPI(sceneH, 144))
	}
}

func TestDrawReactionsSkipsPetsOutsideOverlaySegment(t *testing.T) {
	a := &petApp{
		sceneW: 200,
		pets: []desktopPet{
			{x: 150, state: stateIdle, dir: 1},
		},
		reactions: []petReaction{{pet: 0, kind: 0, ticks: reactionTicks}},
		petSizes:  defaultPetSizes(),
	}
	leftCanvas := image.NewRGBA(image.Rect(0, 0, 100, sceneH))
	a.drawReactions(leftCanvas, 0, 100, defaultDPI)
	if got := alphaSum(leftCanvas); got != 0 {
		t.Fatalf("left segment alpha = %d, want no reaction outside segment", got)
	}

	rightCanvas := image.NewRGBA(image.Rect(0, 0, 100, sceneH))
	a.drawReactions(rightCanvas, 100, 200, defaultDPI)
	if got := alphaSum(rightCanvas); got == 0 {
		t.Fatalf("right segment alpha = 0, want reaction drawn inside segment")
	}
}

func alphaSum(img *image.RGBA) int64 {
	var sum int64
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			sum += int64(img.RGBAAt(x, y).A)
		}
	}
	return sum
}

func TestPetScenePositionsDistributeFivePetsAcrossTwoDisplays(t *testing.T) {
	positions := petScenePositions(3840, 5, []sceneSegment{
		{Left: 0, Right: 1920},
		{Left: 1920, Right: 3840},
	})
	if len(positions) != 5 {
		t.Fatalf("positions = %d, want 5", len(positions))
	}
	mainCount := 0
	subCount := 0
	for i, x := range positions {
		switch {
		case x >= 0 && x+spriteW <= 1920:
			mainCount++
		case x >= 1920 && x+spriteW <= 3840:
			subCount++
		default:
			t.Fatalf("position[%d] = %d escapes monitor segments", i, x)
		}
	}
	if mainCount != 3 || subCount != 2 {
		t.Fatalf("pet distribution = main:%d sub:%d, want 3 and 2", mainCount, subCount)
	}
}

func TestPetScenePositionsAvoidMonitorGaps(t *testing.T) {
	positions := petScenePositions(3800, 4, []sceneSegment{
		{Left: 0, Right: 1600},
		{Left: 2200, Right: 3800},
	})
	for i, x := range positions {
		onLeft := x >= 0 && x+spriteW <= 1600
		onRight := x >= 2200 && x+spriteW <= 3800
		if !onLeft && !onRight {
			t.Fatalf("position[%d] = %d falls in the monitor gap or offscreen", i, x)
		}
	}
}

func TestSetPetCountPlacesAllPetsInsideCurrentScene(t *testing.T) {
	a := &petApp{
		sceneW:        3840,
		speed:         3,
		coatMode:      coatFixed,
		bidirectional: true,
		petCount:      2,
	}

	a.setPetCount(5)
	if len(a.pets) != 5 {
		t.Fatalf("pets = %d, want 5", len(a.pets))
	}
	subCount := 0
	for i, pet := range a.pets {
		if pet.x < 0 || pet.x+spriteW > a.sceneW {
			t.Fatalf("pet %d x = %d escapes scene width %d", i, pet.x, a.sceneW)
		}
		if pet.x >= 1920 {
			subCount++
		}
	}
	if subCount < 2 {
		t.Fatalf("sub-display pets = %d, want at least 2 after choosing 5 pets", subCount)
	}
}

func TestResetPositionDistributesPetsAcrossDetectedMultiMonitorSpan(t *testing.T) {
	areas := monitorAreasByPosition()
	if len(areas) < 2 {
		t.Skip("multi-monitor placement check requires at least two detected displays")
	}
	a := &petApp{
		speed:          3,
		coatMode:       coatFixed,
		bidirectional:  true,
		petCount:       5,
		displayScope:   displayScopeSpan,
		displayIndex:   0,
		displaySpanEnd: len(areas) - 1,
		positionMode:   positionTaskbarEdge,
		walkRangeEnd:   100,
	}
	a.resetPosition()
	overlay := a.overlayRect()
	segments := a.sceneSegmentsForOverlay(overlay)
	if len(segments) < 2 {
		t.Fatalf("detected display span produced %d visible segments, want at least 2", len(segments))
	}
	seen := make([]int, len(segments))
	for _, pet := range a.pets {
		for i, segment := range segments {
			if pet.x >= segment.Left && pet.x+spriteW <= segment.Right {
				seen[i]++
				break
			}
		}
	}
	if len(a.pets) >= len(segments) {
		for i, count := range seen {
			if count == 0 {
				t.Fatalf("segment %d received no pets; distribution=%v segments=%+v", i, seen, segments)
			}
		}
	}
}

func TestWalkRangeSummaryDescribesMultiDisplaySegments(t *testing.T) {
	a := &petApp{lang: langJapanese}
	segments := []sceneSegment{
		{Left: 0, Right: 1920},
		{Left: 1920, Right: 3840},
	}

	cases := []struct {
		name       string
		start, end int
		want       string
	}{
		{name: "all selected displays", start: 0, end: 100, want: "選択した画面ぜんぶ"},
		{name: "first display only", start: 0, end: 50, want: "画面1だけ"},
		{name: "second display only", start: 50, end: 100, want: "画面2だけ"},
		{name: "partial displays", start: 25, end: 75, want: "画面1-2の一部"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := a.walkRangeSummaryForSegments(tt.start, tt.end, segments); got != tt.want {
				t.Fatalf("summary = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDisplaySpanDefaultsWalkingRangeToAllDisplays(t *testing.T) {
	if len(monitorAreas()) < 2 {
		t.Skip("multi-monitor scope reset requires at least two detected displays")
	}
	a := &petApp{
		displayScope:   displayScopeSingle,
		displayIndex:   0,
		displaySpanEnd: 0,
		positionMode:   positionTaskbarEdge,
		walkRangeStart: 25,
		walkRangeEnd:   75,
		petCount:       2,
		speed:          3,
		coatMode:       coatFixed,
	}
	a.resetPosition()
	a.walkRangeStart = 25
	a.walkRangeEnd = 75
	a.setDisplayScope(displayScopeSpan)
	if a.walkRangeStart != 0 || a.walkRangeEnd != 100 {
		t.Fatalf("span walk range = %d-%d, want reset to 0-100", a.walkRangeStart, a.walkRangeEnd)
	}

	a.setWalkRange(50, 100)
	a.adjustDisplaySpan(1)
	if a.walkRangeStart != 50 || a.walkRangeEnd != 100 {
		t.Fatalf("adjusting span changed fine range = %d-%d, want 50-100", a.walkRangeStart, a.walkRangeEnd)
	}
}

func TestPetVariantRectsFitTenPetsInSettingsWindow(t *testing.T) {
	seen := map[[4]int]bool{}
	for i := 0; i < maxPetCount; i++ {
		numberRect, buttonRect := settingsPetVariantRects(i)
		sizeRect := settingsPetSizeRect(i)
		if buttonRect.Right > 708 || buttonRect.Bottom > 562 || sizeRect.Right > 708 || sizeRect.Bottom > 562 {
			t.Fatalf("pet variant button %d rect %+v overflows selected-coats panel", i, buttonRect)
		}
		if numberRect.Left < 238 || buttonRect.Left <= numberRect.Right || sizeRect.Left <= buttonRect.Right {
			t.Fatalf("pet variant %d row rects overlap or escape: number=%+v button=%+v size=%+v", i, numberRect, buttonRect, sizeRect)
		}
		key := [4]int{int(buttonRect.Left), int(buttonRect.Top), int(buttonRect.Right), int(buttonRect.Bottom)}
		if seen[key] {
			t.Fatalf("pet variant button %d duplicates another rect: %+v", i, buttonRect)
		}
		seen[key] = true
	}
}

func TestPetNameRectsFitTenPetsWithCoatPicker(t *testing.T) {
	for i := 0; i < maxPetCount; i++ {
		numberRect, nameRect := settingsPetNameRects(i)
		_, coatRect := settingsPetVariantRects(i)
		sizeRect := settingsPetSizeRect(i)
		if nameRect.Right >= coatRect.Left {
			t.Fatalf("pet %d name rect overlaps coat rect: name=%+v coat=%+v", i, nameRect, coatRect)
		}
		if coatRect.Right >= sizeRect.Left {
			t.Fatalf("pet %d coat rect overlaps size rect: coat=%+v size=%+v", i, coatRect, sizeRect)
		}
		if numberRect.Left < 238 || nameRect.Left <= numberRect.Right || sizeRect.Right > 708 || nameRect.Bottom > 562 || sizeRect.Bottom > 562 {
			t.Fatalf("pet %d name/coat/size row escapes panel: number=%+v name=%+v coat=%+v size=%+v", i, numberRect, nameRect, coatRect, sizeRect)
		}
	}
}

func TestRenameDialogControlsFitClientArea(t *testing.T) {
	editRect, okRect, cancelRect := renameDialogLayoutRects()
	rects := map[string]win.RECT{
		"edit":   editRect,
		"save":   okRect,
		"cancel": cancelRect,
	}
	for name, rect := range rects {
		if rect.Left < 24 || rect.Top < 20 || rect.Right > renameDialogClientW-24 || rect.Bottom > renameDialogClientH-24 {
			t.Fatalf("%s rect escapes rename dialog client area: %+v", name, rect)
		}
		if rectWidth(rect) <= 0 || rectHeight(rect) <= 0 {
			t.Fatalf("%s rect has invalid size: %+v", name, rect)
		}
	}
	if okRect.Right >= cancelRect.Left {
		t.Fatalf("rename action buttons overlap: save=%+v cancel=%+v", okRect, cancelRect)
	}
}

func TestCaptionedRenameDialogExpandsOuterWindowForClientArea(t *testing.T) {
	w, h := windowSizeForClient(renameDialogClientW, renameDialogClientH, renameDialogStyle(), uint32(win.WS_EX_TOOLWINDOW))
	if w <= renameDialogClientW || h <= renameDialogClientH {
		t.Fatalf("captioned rename dialog outer size = %dx%d, want larger than client %dx%d", w, h, renameDialogClientW, renameDialogClientH)
	}
}

func TestSettingsFooterControlsFitClientArea(t *testing.T) {
	rects := map[string]win.RECT{
		"language": {Left: 322, Top: 574, Right: 502, Bottom: 608},
		"reset":    {Left: 534, Top: 576, Right: 634, Bottom: 608},
		"close":    {Left: 646, Top: 576, Right: 724, Bottom: 608},
	}
	for name, rect := range rects {
		if rect.Left < 204 || rect.Right > settingsClientW-24 || rect.Bottom > settingsClientH-12 {
			t.Fatalf("%s footer control escapes settings client area: %+v", name, rect)
		}
		if rectWidth(rect) <= 0 || rectHeight(rect) <= 0 {
			t.Fatalf("%s footer control has invalid size: %+v", name, rect)
		}
	}
	if rects["reset"].Right >= rects["close"].Left {
		t.Fatalf("settings footer buttons overlap: reset=%+v close=%+v", rects["reset"], rects["close"])
	}
}

func TestUpdateVersionComparison(t *testing.T) {
	tests := []struct {
		latest  string
		current string
		want    bool
	}{
		{"v1.2.0", "v1.1.9", true},
		{"v1.2.0", "1.2.0", false},
		{"v1.2.0", "v1.3.0", false},
		{"v2.0.0", "dev", true},
		{"v2.0.0", "pages-abc123", true},
		{"not-semver", "v1.0.0", false},
	}
	for _, tt := range tests {
		if got := isNewerVersion(tt.latest, tt.current); got != tt.want {
			t.Fatalf("isNewerVersion(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
		}
	}
}

func TestSelectUpdateAssetFindsWindowsZip(t *testing.T) {
	rel := &githubRelease{Assets: []githubReleaseAsset{
		{Name: "notes.txt", BrowserDownloadURL: "https://example.test/notes.txt"},
		{Name: "AnimalsDesktop-windows-amd64.zip", BrowserDownloadURL: "https://example.test/app.zip"},
		{Name: "AnimalsDesktop-windows-386.zip", BrowserDownloadURL: "https://example.test/app-x86.zip"},
	}}
	asset := selectUpdateAsset(rel, "amd64")
	if asset == nil || asset.BrowserDownloadURL != "https://example.test/app.zip" {
		t.Fatalf("selectUpdateAsset(amd64) = %+v", asset)
	}
	asset = selectUpdateAsset(rel, "386")
	if asset == nil || asset.BrowserDownloadURL != "https://example.test/app-x86.zip" {
		t.Fatalf("selectUpdateAsset(386) = %+v", asset)
	}
}

func skipIfNetworkUpdatesDisabled(t *testing.T) {
	t.Helper()
	if !networkUpdatesEnabled {
		t.Skip("update runtime is excluded from the no-network build")
	}
}

func TestVerifyDownloadedAssetChecksSizeAndSHA256Digest(t *testing.T) {
	skipIfNetworkUpdatesDisabled(t)

	path := filepath.Join(t.TempDir(), "update.zip")
	data := []byte("trusted update bytes")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write update: %v", err)
	}
	sum := sha256.Sum256(data)
	asset := githubReleaseAsset{
		Size:   int64(len(data)),
		Digest: fmt.Sprintf("sha256:%x", sum[:]),
	}
	if err := verifyDownloadedAsset(path, asset); err != nil {
		t.Fatalf("verifyDownloadedAsset() error = %v", err)
	}
	asset.Size++
	if err := verifyDownloadedAsset(path, asset); err == nil {
		t.Fatalf("verifyDownloadedAsset accepted a size mismatch")
	}
	asset.Size = int64(len(data))
	asset.Digest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	if err := verifyDownloadedAsset(path, asset); err == nil {
		t.Fatalf("verifyDownloadedAsset accepted a digest mismatch")
	}
}

func TestParseUpdateApplyArgsRequiresSafeCleanupDir(t *testing.T) {
	skipIfNetworkUpdatesDisabled(t)

	cleanupDir := filepath.Join(os.TempDir(), updateTempPrefix+"unit-test")
	opts, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(cleanupDir, "AnimalsDesktop.exe"),
		"--target", filepath.Join(t.TempDir(), "AnimalsDesktop.exe"),
		"--parent-pid", "1234",
		"--cleanup-dir", cleanupDir,
	})
	if err != nil {
		t.Fatalf("parseUpdateApplyArgs() error = %v", err)
	}
	if opts.ParentPID != 1234 || opts.CleanupDir != cleanupDir {
		t.Fatalf("parseUpdateApplyArgs() = %+v", opts)
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", "a.exe",
		"--target", "b.exe",
		"--cleanup-dir", t.TempDir(),
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a non-update cleanup dir")
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(t.TempDir(), "AnimalsDesktop.exe"),
		"--target", filepath.Join(t.TempDir(), "AnimalsDesktop.exe"),
		"--cleanup-dir", cleanupDir,
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a source outside cleanup dir")
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(cleanupDir, "payload", "AnimalsDesktop.exe"),
		"--target", filepath.Join(cleanupDir, "installed", "AnimalsDesktop.exe"),
		"--cleanup-dir", cleanupDir,
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a target inside cleanup dir")
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(cleanupDir, "payload", "AnimalsDesktop.exe"),
		"--target", filepath.Join(t.TempDir(), "renamed.exe"),
		"--cleanup-dir", cleanupDir,
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a renamed target")
	}
}

func TestExtractUpdateExeUsesFixedPayloadPath(t *testing.T) {
	skipIfNetworkUpdatesDisabled(t)

	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "update.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(zipFile)
	w, err := zw.Create("release/nested/AnimalsDesktop.exe")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write([]byte("exe bytes")); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := zipFile.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	exePath, err := extractUpdateExe(zipPath, tmpDir)
	if err != nil {
		t.Fatalf("extractUpdateExe() error = %v", err)
	}
	want := filepath.Join(tmpDir, "payload", "AnimalsDesktop.exe")
	if exePath != want {
		t.Fatalf("extractUpdateExe() = %q, want %q", exePath, want)
	}
	data, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatalf("read extracted exe: %v", err)
	}
	if string(data) != "exe bytes" {
		t.Fatalf("extracted data = %q", data)
	}
}

func TestUpdaterCommandsInvokeAppExeDirectly(t *testing.T) {
	skipIfNetworkUpdatesDisabled(t)

	cleanupDir := filepath.Join(os.TempDir(), updateTempPrefix+"command-test")
	sourceExe := filepath.Join(cleanupDir, "payload", "AnimalsDesktop.exe")
	targetExe := filepath.Join(t.TempDir(), "AnimalsDesktop.exe")
	helperExe := filepath.Join(cleanupDir, "helper", "AnimalsDesktop.exe")

	applyCmd := newUpdaterHelperCommand(helperExe, cleanupDir, sourceExe, targetExe, 1234)
	assertCommandAvoidsPowerShell(t, applyCmd)
	if !strings.EqualFold(applyCmd.Path, helperExe) {
		t.Fatalf("apply command path = %q, want %q", applyCmd.Path, helperExe)
	}
	assertArgsContainInOrder(t, applyCmd.Args,
		updaterApplyArg,
		"--source", sourceExe,
		"--target", targetExe,
		"--parent-pid", "1234",
		"--cleanup-dir", cleanupDir,
	)
	if applyCmd.SysProcAttr == nil || !applyCmd.SysProcAttr.HideWindow {
		t.Fatalf("apply command should hide its helper window")
	}

	cleanupCmd := newUpdaterCleanupCommand(targetExe, cleanupDir)
	assertCommandAvoidsPowerShell(t, cleanupCmd)
	if !strings.EqualFold(cleanupCmd.Path, targetExe) {
		t.Fatalf("cleanup command path = %q, want %q", cleanupCmd.Path, targetExe)
	}
	assertArgsContainInOrder(t, cleanupCmd.Args, updaterCleanupArg, cleanupDir)
	if cleanupCmd.SysProcAttr == nil || !cleanupCmd.SysProcAttr.HideWindow {
		t.Fatalf("cleanup command should hide its helper window")
	}
}

func TestReleaseWorkflowPowerShellBlocksParseAfterGitHubSubstitution(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	scriptByStep := map[string]string{
		"Generate Windows assets": extractWorkflowRunBlock(t, workflowPath, "Generate Windows assets"),
		"Build Windows":           extractWorkflowRunBlock(t, workflowPath, "Build Windows"),
		"Package Windows":         extractWorkflowRunBlock(t, workflowPath, "Package Windows"),
	}
	for stepName, script := range scriptByStep {
		t.Run(stepName, func(t *testing.T) {
			script = strings.ReplaceAll(script, "${{ github.ref_name }}", "v0.2.1")
			assertPowerShellParses(t, script)
		})
	}
}

func TestReleaseWorkflowPackageIncludesSecurityManifestAndHashes(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	script := extractWorkflowRunBlock(t, workflowPath, "Package Windows")
	for _, want := range []string{
		"SECURITY.txt",
		"SHA256SUMS.txt",
		"AnimalsDesktop-windows-amd64.zip/AnimalsDesktop.exe",
		"AnimalsDesktop-windows-amd64-no-network.zip/AnimalsDesktop.exe",
		"AnimalsDesktop-windows-386.zip/AnimalsDesktop.exe",
		"No-network security-check edition",
		"The Go net/http update fetch/download implementation is excluded from this build.",
		"Microsoft Security Intelligence",
		"McAfee Dispute Detection & Allowlisting",
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("Package script does not contain %q", want)
		}
	}
}

func TestReleaseWorkflowPublishesMainLineWindowsTrustAssets(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	workflow := string(data)
	for _, want := range []string{
		"go run ./cmd/validatemotion -runtime-only -require-accepted",
		"go build -buildvcs=false",
		"-tags animalsdesktop_nonetwork",
		"AnimalsDesktop-windows-amd64-no-network.zip",
		"./cmd/animalsdesktop",
		"body_path: dist/RELEASE_NOTES.md",
		"github.ref_name == 'v0.2.1'",
		"github.ref_name == 'v0.2.2'",
		"github.ref_name == 'v0.2.3'",
		"github.ref_name == 'v0.2.4'",
		"github.ref_name == 'v0.2.6'",
		"github.ref_name == 'v0.2.7'",
		"docs/releases/${version}.md",
		"checksum_source=\"$(find release-assets -type f -name SHA256SUMS.txt -print -quit || true)\"",
		"sed 's/\\r$//' \"$checksum_source\" | grep '/AnimalsDesktop.exe$' >> dist/SHA256SUMS.txt || true",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("release workflow does not contain %q", want)
		}
	}
}

func TestUpdateCleanupDirOnlyAcceptsUpdateTempDirs(t *testing.T) {
	skipIfNetworkUpdatesDisabled(t)

	cleanupDir := filepath.Join(os.TempDir(), updateTempPrefix+"cleanup-test")
	if got := updateCleanupDir([]string{updaterCleanupArg, cleanupDir}); got != cleanupDir {
		t.Fatalf("updateCleanupDir() = %q, want %q", got, cleanupDir)
	}
	if got := updateCleanupDir([]string{updaterCleanupArg, t.TempDir()}); got != "" {
		t.Fatalf("updateCleanupDir accepted non-update dir %q", got)
	}
}

func TestTurnStateUsesGeneratedTurnFrames(t *testing.T) {
	if got := currentFrame(stateTurn, 0); got != turnStart {
		t.Fatalf("turn frame 0 = %d, want %d", got, turnStart)
	}
	if got := currentFrame(stateTurn, turnTicks-1); got != turnStart+turnFrames-1 {
		t.Fatalf("turn final active frame = %d, want %d", got, turnStart+turnFrames-1)
	}
	if got := currentFrame(stateTurn, turnTicks+10); got != turnStart+turnFrames-1 {
		t.Fatalf("turn frame after duration = %d, want held final frame %d", got, turnStart+turnFrames-1)
	}
}

func TestTurnDrawDirectionMirrorsOnlyLeftToRightTurns(t *testing.T) {
	if got := turnDrawDirection(1, -1); got != 1 {
		t.Fatalf("right-to-left turn draw direction = %d, want 1", got)
	}
	if got := turnDrawDirection(-1, 1); got != -1 {
		t.Fatalf("left-to-right turn draw direction = %d, want -1", got)
	}
}

func TestSetBidirectionalOffNormalizesPets(t *testing.T) {
	a := &petApp{
		bidirectional: true,
		speed:         3,
		pets: []desktopPet{
			{state: stateTurn, dir: -1, nextDir: -1, item: noItem},
			{state: stateWalk, dir: -1, nextDir: -1, item: noItem},
		},
	}

	a.setBidirectional(false)
	if a.bidirectional {
		t.Fatalf("bidirectional stayed enabled")
	}
	for i, pet := range a.pets {
		if pet.dir != 1 || pet.nextDir != 1 {
			t.Fatalf("pet %d direction = (%d,%d), want (1,1)", i, pet.dir, pet.nextDir)
		}
		if pet.state == stateTurn {
			t.Fatalf("pet %d remained in stateTurn", i)
		}
	}
}

func TestResetPetAtEdgeReentersFromOppositeSideWithMatchingDirection(t *testing.T) {
	a := &petApp{
		sceneW:        500,
		speed:         3,
		coatMode:      coatSelected,
		selectedCoats: defaultSelectedCoats(),
		forage: []forageItem{
			{owner: 0, active: true},
			{owner: reservedItem, active: true},
		},
	}

	right := desktopPet{dir: 1, nextDir: 1, item: 0, carryKind: 2, state: stateCarry}
	a.resetPetAtEdge(0, &right, 1)
	if right.x > -spriteW || right.dir != 1 || right.nextDir != 1 {
		t.Fatalf("right-moving reset = x:%d dir:%d next:%d, want off-left and direction +1", right.x, right.dir, right.nextDir)
	}
	if right.item != noItem || right.carryKind != noItem || right.state != stateWalk {
		t.Fatalf("right-moving reset state = item:%d carry:%d state:%v, want cleared walk", right.item, right.carryKind, right.state)
	}
	if a.forage[0].owner != noItem {
		t.Fatalf("owned forage was not released: owner=%d", a.forage[0].owner)
	}

	left := desktopPet{dir: -1, nextDir: -1, item: noItem, state: stateWalk}
	a.resetPetAtEdge(1, &left, -1)
	if left.x < a.sceneW || left.dir != -1 || left.nextDir != -1 {
		t.Fatalf("left-moving reset = x:%d dir:%d next:%d, want off-right and direction -1", left.x, left.dir, left.nextDir)
	}
}

func TestForagePropsDisabledClearsPropsAndStopsAssignment(t *testing.T) {
	if foragePropsEnabled {
		t.Fatalf("foragePropsEnabled = true, want false for preview release")
	}
	a := &petApp{
		sceneW: 500,
		speed:  3,
		forage: []forageItem{
			{x: 100, kind: 2, owner: 0, active: true},
			{x: 160, kind: 1, owner: reservedItem, active: true},
		},
		pets: []desktopPet{
			{state: stateCarry, item: 0, carryKind: 2, dir: 1, nextDir: 1},
			{state: stateForage, item: 1, carryKind: noItem, dir: -1, nextDir: -1},
		},
	}

	a.clearForageItems()

	for i, item := range a.forage {
		if item.active || item.owner != noItem {
			t.Fatalf("forage %d = active:%v owner:%d, want cleared", i, item.active, item.owner)
		}
	}
	for i, pet := range a.pets {
		if pet.item != noItem || pet.carryKind != noItem || pet.state != stateWalk {
			t.Fatalf("pet %d = item:%d carry:%d state:%v, want cleared walk", i, pet.item, pet.carryKind, pet.state)
		}
	}

	a.forage = []forageItem{{x: 140, kind: 2, owner: noItem, active: true}}
	p := desktopPet{state: stateWalk, item: noItem, carryKind: noItem, x: 40, dir: 1}
	if a.maybeAssignForageTarget(&p) {
		t.Fatalf("maybeAssignForageTarget assigned forage while props are disabled")
	}
	if p.item != noItem || p.state != stateWalk {
		t.Fatalf("pet after disabled assignment = item:%d state:%v, want unchanged", p.item, p.state)
	}
}

func TestTickPetMovesByDirectionAndWrapsPastEdges(t *testing.T) {
	a := &petApp{
		sceneW:        240,
		speed:         3,
		coatMode:      coatSelected,
		selectedCoats: defaultSelectedCoats(),
	}

	right := desktopPet{x: 20, dir: 1, nextDir: 1, state: stateWalk, moveSpeed: 4, stateTicks: 10, item: noItem, carryKind: noItem}
	a.tickPet(0, &right)
	if right.x != 24 {
		t.Fatalf("right-moving tick x = %d, want 24", right.x)
	}

	left := desktopPet{x: 20, dir: -1, nextDir: -1, state: stateWalk, moveSpeed: 4, stateTicks: 10, item: noItem, carryKind: noItem}
	a.tickPet(0, &left)
	if left.x != 16 {
		t.Fatalf("left-moving tick x = %d, want 16", left.x)
	}

	right.x = a.sceneW + 9
	a.tickPet(0, &right)
	if right.x > -spriteW || right.dir != 1 {
		t.Fatalf("right edge wrap = x:%d dir:%d, want off-left dir +1", right.x, right.dir)
	}

	left.x = -spriteW - 9
	a.tickPet(0, &left)
	if left.x < a.sceneW || left.dir != -1 {
		t.Fatalf("left edge wrap = x:%d dir:%d, want off-right dir -1", left.x, left.dir)
	}
}

func TestDrawFacingImageMirrorsNegativeDirection(t *testing.T) {
	red := color.RGBA{R: 255, A: 255}
	blue := color.RGBA{B: 255, A: 255}
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	src.SetRGBA(0, 0, red)
	src.SetRGBA(1, 0, blue)

	dst := image.NewRGBA(src.Bounds())
	drawFacingImage(dst, src, dst.Bounds(), 1)
	if got := dst.RGBAAt(0, 0); got != red {
		t.Fatalf("positive left pixel = %#v, want %#v", got, red)
	}
	if got := dst.RGBAAt(1, 0); got != blue {
		t.Fatalf("positive right pixel = %#v, want %#v", got, blue)
	}

	dst = image.NewRGBA(src.Bounds())
	drawFacingImage(dst, src, dst.Bounds(), -1)
	if got := dst.RGBAAt(0, 0); got != blue {
		t.Fatalf("negative left pixel = %#v, want %#v", got, blue)
	}
	if got := dst.RGBAAt(1, 0); got != red {
		t.Fatalf("negative right pixel = %#v, want %#v", got, red)
	}
}

func TestDrawDirectionForVariantCompensatesLeftFacingSource(t *testing.T) {
	sugarGlider := coatVariant{ID: "sugar_glider_gray"}
	hamster := coatVariant{ID: "hamster_golden_syrian"}

	if got := drawDirectionForVariant(1, sugarGlider); got != -1 {
		t.Fatalf("sugar glider right-moving draw direction = %d, want -1", got)
	}
	if got := drawDirectionForVariant(-1, sugarGlider); got != 1 {
		t.Fatalf("sugar glider left-moving draw direction = %d, want 1", got)
	}
	if got := drawDirectionForVariant(1, hamster); got != 1 {
		t.Fatalf("hamster right-moving draw direction = %d, want 1", got)
	}
	if got := drawDirectionForVariant(-1, hamster); got != -1 {
		t.Fatalf("hamster left-moving draw direction = %d, want -1", got)
	}
}

func TestFixedCoatModeRefreshesAllPets(t *testing.T) {
	a := &petApp{
		variant: 99,
		pets: []desktopPet{
			{variant: 0},
			{variant: 0},
			{variant: 0},
		},
	}

	a.setCoatMode(coatFixed)

	want := len(variants) - 1
	for i, pet := range a.pets {
		if pet.variant != want {
			t.Fatalf("pet %d variant = %d, want fixed variant %d", i, pet.variant, want)
		}
	}
}

func TestSelectedCoatModeUsesPerPetChoices(t *testing.T) {
	a := &petApp{
		selectedCoats: [maxPetCount]int{0, 3, 5, 7, 9},
		pets: []desktopPet{
			{variant: 0},
			{variant: 0},
			{variant: 0},
		},
	}

	a.setCoatMode(coatSelected)

	want := []int{0, 3, 5}
	for i := range []int{0, 1, 2} {
		if got := a.pets[i].variant; got != want[i] {
			t.Fatalf("pet %d variant = %d, want %d", i, got, want[i])
		}
	}
	a.setSelectedVariant(1, 8)
	if got, want := a.pets[1].variant, 8; got != want {
		t.Fatalf("selected variant update = %d, want %d", got, want)
	}
}

func TestSelectedCoatModeMixesFixedAndRandomSlots(t *testing.T) {
	rabbitIndex, ok := variantIndexByID("rabbit_chestnut_agouti")
	if !ok {
		t.Fatal("rabbit_chestnut_agouti missing from runtime variants")
	}
	chinchillaIndex, ok := variantIndexByID("chinchilla_standard_gray")
	if !ok {
		t.Fatal("chinchilla_standard_gray missing from runtime variants")
	}
	a := &petApp{
		coatMode:               coatSelected,
		selectedCoats:          [maxPetCount]int{rabbitIndex, chinchillaIndex},
		selectedSlotModes:      [maxPetCount]petSlotMode{petSlotFixed, petSlotFixed, petSlotRandom, petSlotRandom, petSlotRandom},
		selectedRandomGroupIDs: [maxPetCount]string{"", "", "rabbit", "chinchilla", ""},
		pets: []desktopPet{
			{variant: 0},
			{variant: 0},
			{variant: 0},
			{variant: 0},
			{variant: 0},
		},
	}

	a.refreshPetVariants()

	if got := a.pets[0].variant; got != rabbitIndex {
		t.Fatalf("fixed slot 0 = %d, want rabbit index %d", got, rabbitIndex)
	}
	if got := a.pets[1].variant; got != chinchillaIndex {
		t.Fatalf("fixed slot 1 = %d, want chinchilla index %d", got, chinchillaIndex)
	}
	if got := catalog.VariantGroupIDForSpecies(variants[a.pets[2].variant].SpeciesID); got != "rabbit" {
		t.Fatalf("random slot 2 group = %q, want rabbit", got)
	}
	if got := catalog.VariantGroupIDForSpecies(variants[a.pets[3].variant].SpeciesID); got != "chinchilla" {
		t.Fatalf("random slot 3 group = %q, want chinchilla", got)
	}
	if got := a.pets[4].variant; got < 0 || got >= len(variants) {
		t.Fatalf("random-all slot 4 = %d, want valid variant", got)
	}
}

func TestRandomCoatModeAssignsValidVariants(t *testing.T) {
	a := &petApp{coatMode: coatRandom}
	for i := 0; i < 100; i++ {
		got := a.variantForIndex(i)
		if got < 0 || got >= len(variants) {
			t.Fatalf("random variant = %d, want 0..%d", got, len(variants)-1)
		}
	}
}

func TestRandomCoatModeCanFilterByVariantGroup(t *testing.T) {
	a := &petApp{coatMode: coatRandom, randomGroupID: "rabbit"}
	for i := 0; i < 100; i++ {
		got := a.variantForIndex(i)
		if group := catalog.VariantGroupIDForSpecies(variants[got].SpeciesID); group != "rabbit" {
			t.Fatalf("random rabbit variant group = %q (%s), want rabbit", group, variants[got].ID)
		}
	}

	a.setRandomGroupID("chinchilla")
	for i := 0; i < 100; i++ {
		got := a.variantForIndex(i)
		if group := catalog.VariantGroupIDForSpecies(variants[got].SpeciesID); group != "chinchilla" {
			t.Fatalf("random chinchilla variant group = %q (%s), want chinchilla", group, variants[got].ID)
		}
	}
}

func TestFerretVariantsAreEligibleForSmallAnimalRandomGroup(t *testing.T) {
	indices := variantIndicesForGroupID("small_mammal")
	eligible := make(map[string]bool, len(indices))
	for _, index := range indices {
		if index < 0 || index >= len(variants) {
			t.Fatalf("small_mammal index = %d, want 0..%d", index, len(variants)-1)
		}
		variant := variants[index]
		if got := catalog.VariantGroupIDForSpecies(variant.SpeciesID); got != "small_mammal" {
			t.Fatalf("small_mammal contains %q in group %q", variant.ID, got)
		}
		eligible[variant.ID] = true
	}

	a := &petApp{lang: langJapanese}
	for _, variantID := range []string{
		"ferret_sable_panda",
		"ferret_sable",
		"ferret_albino",
	} {
		if !eligible[variantID] {
			t.Errorf("%s missing from small_mammal random group", variantID)
		}
		index, ok := variantIndexByID(variantID)
		if !ok {
			t.Fatalf("%s missing from runtime variants", variantID)
		}
		if got := a.variantDisplayLabel(index); !strings.HasPrefix(got, "小動物 / フェレット") {
			t.Errorf("Japanese %s display label = %q", variantID, got)
		}
		a.lang = langEnglish
		if got := a.variantDisplayLabel(index); !strings.HasPrefix(got, "Small animals / Ferret") {
			t.Errorf("English %s display label = %q", variantID, got)
		}
		a.lang = langJapanese
	}

	for draw := 0; draw < 256; draw++ {
		index := randomVariantForGroup("small_mammal")
		if index < 0 || index >= len(variants) {
			t.Fatalf("random small_mammal index = %d, want 0..%d", index, len(variants)-1)
		}
		if got := catalog.VariantGroupIDForSpecies(variants[index].SpeciesID); got != "small_mammal" {
			t.Fatalf("random small_mammal draw selected %q in group %q", variants[index].ID, got)
		}
	}
}

func TestPetAtScenePointFindsTopmostPet(t *testing.T) {
	a := &petApp{
		sceneW:   800,
		petSizes: defaultPetSizes(),
		pets: []desktopPet{
			{x: 100, laneOffset: 0, state: stateWalk},
			{x: 110, laneOffset: 0, state: stateIdle},
		},
	}

	got := a.petAtScenePoint(132, sceneH-spriteH+24)
	if got != 1 {
		t.Fatalf("petAtScenePoint overlap = %d, want topmost pet 1", got)
	}
	if got := a.petAtScenePoint(4, 4); got != -1 {
		t.Fatalf("petAtScenePoint outside = %d, want -1", got)
	}
}

func TestPetSizeAffectsHitTestingAndBounds(t *testing.T) {
	a := &petApp{
		sceneW:   800,
		petSizes: defaultPetSizes(),
		pets: []desktopPet{
			{x: 100, laneOffset: 0, state: stateWalk},
		},
	}
	a.setPetSize(0, 120)
	w, h := a.petSpriteSize(0)
	if w != 115 || h != 76 {
		t.Fatalf("petSpriteSize(120%%) = %dx%d, want 115x76", w, h)
	}
	if got := a.petAtScenePoint(100+w-8, sceneH-h+8); got != 0 {
		t.Fatalf("petAtScenePoint on enlarged pet = %d, want 0", got)
	}
	a.pets[0].x = 760
	a.setPetSize(0, 120)
	if a.pets[0].x+w > a.sceneW {
		t.Fatalf("setPetSize did not clamp x: x=%d w=%d scene=%d", a.pets[0].x, w, a.sceneW)
	}
}

func TestNormalizePetSizePercent(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{0, defaultPetSizePercent},
		{64, minPetSizePercent},
		{86, 90},
		{119, maxPetSizePercent},
		{200, maxPetSizePercent},
	}
	for _, tt := range tests {
		if got := normalizePetSizePercent(tt.in); got != tt.want {
			t.Fatalf("normalizePetSizePercent(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestShowPetReactionRefreshesExistingBubble(t *testing.T) {
	a := &petApp{
		pets: []desktopPet{{state: stateWalk}},
		reactions: []petReaction{
			{pet: 0, kind: 1, ticks: 3},
		},
	}

	a.showPetReaction(0)
	if len(a.reactions) != 1 {
		t.Fatalf("reaction count = %d, want 1 refreshed reaction", len(a.reactions))
	}
	if a.reactions[0].ticks != reactionTicks {
		t.Fatalf("reaction ticks = %d, want %d", a.reactions[0].ticks, reactionTicks)
	}
}

func TestTickReactionsDropsExpiredAndInvalid(t *testing.T) {
	a := &petApp{
		pets: []desktopPet{{state: stateWalk}},
		reactions: []petReaction{
			{pet: 0, ticks: 1},
			{pet: 3, ticks: 5},
		},
	}

	a.tickReactions()
	if len(a.reactions) != 0 {
		t.Fatalf("remaining reactions = %d, want 0", len(a.reactions))
	}
}

func TestSpriteCacheLoadsVariantOnDemand(t *testing.T) {
	cache := newSpriteCache()
	if got := len(cache.loaded); got != 0 {
		t.Fatalf("new cache loaded variants = %d, want 0", got)
	}
	frame := cache.frame(variants[0], 0, 0)
	if frame == nil || frame.Bounds().Dx() != frameW || frame.Bounds().Dy() != frameH {
		t.Fatalf("loaded frame bounds = %v", frame.Bounds())
	}
	if got := len(cache.loaded); got != 1 {
		t.Fatalf("cache loaded variants = %d, want 1", got)
	}
	_ = cache.frame(variants[0], motionSets+99, frameCount+99)
	if got := len(cache.loaded); got != 1 {
		t.Fatalf("cache reloaded same variant; loaded = %d", got)
	}
}

func assertCommandAvoidsPowerShell(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
	joined := strings.ToLower(strings.Join(cmd.Args, " "))
	for _, forbidden := range []string{"powershell", "pwsh", ".ps1", "executionpolicy", "bypass"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("command unexpectedly contains %q: %q", forbidden, joined)
		}
	}
}

func assertArgsContainInOrder(t *testing.T, args []string, want ...string) {
	t.Helper()
	next := 0
	for _, arg := range args {
		if next < len(want) && arg == want[next] {
			next++
		}
	}
	if next != len(want) {
		t.Fatalf("args %q did not contain %q in order", args, want)
	}
}

func extractWorkflowRunBlock(t *testing.T, workflowPath, stepName string) string {
	t.Helper()
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	stepIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "- name: "+stepName {
			stepIndex = i
			break
		}
	}
	if stepIndex < 0 {
		t.Fatalf("workflow step %q was not found", stepName)
	}
	runIndex := -1
	for i := stepIndex + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "- name: ") {
			break
		}
		if trimmed == "run: |" {
			runIndex = i
			break
		}
	}
	if runIndex < 0 {
		t.Fatalf("workflow step %q has no run block", stepName)
	}
	contentIndent := -1
	for i := runIndex + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		contentIndent = leadingSpaces(lines[i])
		break
	}
	if contentIndent < 0 {
		t.Fatalf("workflow step %q has an empty run block", stepName)
	}
	var out []string
	for i := runIndex + 1; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		if strings.TrimSpace(line) != "" && leadingSpaces(line) < contentIndent {
			break
		}
		if len(line) >= contentIndent {
			line = line[contentIndent:]
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func leadingSpaces(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}

func assertPowerShellParses(t *testing.T, script string) {
	t.Helper()
	powershell, err := exec.LookPath("powershell.exe")
	if err != nil {
		powershell, err = exec.LookPath("powershell")
	}
	if err != nil {
		t.Skipf("PowerShell was not found: %v", err)
	}
	scriptPath := filepath.Join(t.TempDir(), "script-under-test.ps1")
	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		t.Fatalf("write script under test: %v", err)
	}
	parser := `
$script = Get-Content -LiteralPath $env:ANIMALSDESKTOP_PARSE_SCRIPT_PATH -Raw
$tokens = $null
$errors = $null
[System.Management.Automation.Language.Parser]::ParseInput($script, [ref]$tokens, [ref]$errors) | Out-Null
if ($errors.Count -gt 0) {
  $errors | ForEach-Object { $_.Message }
  exit 1
}
`
	cmd := exec.Command(powershell, "-NoProfile", "-NonInteractive", "-Command", parser)
	cmd.Env = append(os.Environ(), "ANIMALSDESKTOP_PARSE_SCRIPT_PATH="+scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("PowerShell parser rejected script: %v\n%s", err, out)
	}
}
