//go:build darwin

package main

/*
#cgo darwin CFLAGS: -fblocks
#cgo darwin LDFLAGS: -framework Cocoa
#include "darwin_cocoa.h"
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	appassets "animals-desktop/assets"
	"animals-desktop/internal/catalog"
	"animals-desktop/internal/motion"
)

const (
	frameW        = 96
	frameH        = 64
	frameCount    = 62
	spriteW       = frameW
	spriteH       = frameH
	sceneH        = 92
	wheelSize     = 72
	timerInterval = 55
	maxPetCount   = 10
	wheelKeyHold  = 18
	reactionTicks = 54
	maxPetNameLen = 24

	defaultPetSizePercent = 100
	minPetSizePercent     = 70
	maxPetSizePercent     = 120
	petSizeStepPercent    = 10
)

const (
	idleStart    = 0
	idleFrames   = 4
	walkStart    = 4
	walkFrames   = 8
	scurryStart  = 12
	scurryFrames = 8
	nibbleStart  = 20
	nibbleFrames = 6
	hopStart     = 26
	hopFrames    = 6

	ferretSniffStart  = nibbleStart
	ferretSniffFrames = nibbleFrames
	ferretGroomStart  = hopStart
	ferretGroomFrames = hopFrames
	ferretRestStart   = 48
	ferretRestFrames  = 8
)

const (
	settingsDirName  = "AnimalsDesktop"
	settingsFileName = "settings.json"

	darwinSettingsVersion = 4

	darwinSpeedSlow   = 2
	darwinSpeedNormal = 3
	darwinSpeedFast   = 5
)

type darwinBehaviorMode int

const (
	darwinModeKeyboard darwinBehaviorMode = iota
	darwinModeRandom
)

type darwinCoatMode int

const (
	darwinCoatFixed darwinCoatMode = iota
	darwinCoatSelected
	darwinCoatRandom
)

type darwinLanguage int

const (
	darwinLangJapanese darwinLanguage = iota
	darwinLangEnglish
)

var (
	idleFrameSeq   = []int{idleStart, idleStart + 1, idleStart + 3, idleStart + 1}
	walkFrameSeq   = []int{walkStart, walkStart + 1, walkStart + 3, walkStart + 1}
	nibbleFrameSeq = []int{nibbleStart, nibbleStart + 1, nibbleStart + 2, nibbleStart + 1}
	hopFrameSeq    = []int{hopStart, hopStart + 1, hopStart + 2, hopStart + 3}
	ferretSniffSeq = []int{ferretSniffStart, ferretSniffStart + 1, ferretSniffStart + 2, ferretSniffStart + 3, ferretSniffStart + 4, ferretSniffStart + 5}
	ferretGroomSeq = []int{ferretGroomStart, ferretGroomStart + 1, ferretGroomStart + 2, ferretGroomStart + 3, ferretGroomStart + 4, ferretGroomStart + 5}
	ferretRestSeq  = []int{ferretRestStart, ferretRestStart + 1, ferretRestStart + 2, ferretRestStart + 3, ferretRestStart + 4, ferretRestStart + 5, ferretRestStart + 6, ferretRestStart + 7}
)

type darwinCoatVariant struct {
	ID            string
	SpeciesID     string
	SpriteBase    string
	LabelJA       string
	LabelEN       string
	MotionProfile string
	WheelCapable  bool
}

var darwinVariants = darwinRuntimeVariants()

var appVersion = "v2.0.0"

var darwinApp *darwinPetApp

type darwinPetApp struct {
	motionClock   motion.Clock
	legacyElapsed float64
	mu            sync.Mutex
	sceneW        int
	tick          int
	keyHold       int
	speed         int
	petCount      int
	lang          darwinLanguage
	displayID     int64
	mode          darwinBehaviorMode
	coatMode      darwinCoatMode
	variant       int
	selectedCoats [maxPetCount]int
	petNames      [maxPetCount]string
	petSizes      [maxPetCount]int
	nameLabels    bool
	wheelEnabled  bool
	reactions     []darwinReaction
	frames        map[string][]*image.RGBA
	wheel         *image.RGBA
	pets          []darwinPet
}

type darwinPet struct {
	timed       motion.Player
	x           int
	lane        int
	dir         int
	speed       int
	frame       int
	variant     int
	nextPause   int
	pause       int
	pauseAction int
}

type darwinReaction struct {
	pet   int
	kind  int
	ticks int
}

type darwinSettings struct {
	Version         int      `json:"version"`
	Variant         *int     `json:"variant,omitempty"`
	VariantID       string   `json:"variantID,omitempty"`
	CoatMode        *int     `json:"coatMode,omitempty"`
	SelectedCoats   []int    `json:"selectedCoats,omitempty"`
	SelectedCoatIDs []string `json:"selectedCoatIDs,omitempty"`
	Speed           int      `json:"speed"`
	Language        int      `json:"language,omitempty"`
	DisplayID       int64    `json:"displayID,omitempty"`
	Mode            *int     `json:"mode,omitempty"`
	PetSizes        []int    `json:"petSizes,omitempty"`
	PetCount        int      `json:"petCount"`
	WheelEnabled    *bool    `json:"wheelEnabled,omitempty"`
	NameLabels      bool     `json:"nameLabels"`
	PetNames        []string `json:"petNames,omitempty"`
}

var darwinSettingsPath = defaultDarwinSettingsPath

func darwinRuntimeVariants() []darwinCoatVariant {
	runtimeVariants := catalog.RuntimeVariants()
	out := make([]darwinCoatVariant, 0, len(runtimeVariants))
	for _, variant := range runtimeVariants {
		spriteBase := variant.SpriteBase
		if spriteBase == "" {
			spriteBase = variant.ID
		}
		labelJA := variant.LabelJA
		if labelJA == "" {
			labelJA = variant.LabelEN
		}
		labelEN := variant.LabelEN
		if labelEN == "" {
			labelEN = labelJA
		}
		out = append(out, darwinCoatVariant{
			ID:            variant.ID,
			SpeciesID:     variant.SpeciesID,
			SpriteBase:    spriteBase,
			LabelJA:       labelJA,
			LabelEN:       labelEN,
			MotionProfile: catalog.MotionProfileForVariant(variant),
			WheelCapable:  catalog.WheelCapableVariant(variant),
		})
	}
	return out
}

func darwinVariantLabel(index int, lang darwinLanguage) string {
	if index < 0 || index >= len(darwinVariants) {
		return ""
	}
	if lang == darwinLangEnglish {
		return darwinVariants[index].LabelEN
	}
	return darwinVariants[index].LabelJA
}

func darwinVariantGroupLabel(index int, lang darwinLanguage) string {
	if index < 0 || index >= len(darwinVariants) {
		return ""
	}
	group := catalog.VariantGroupForSpecies(darwinVariants[index].SpeciesID)
	if lang == darwinLangEnglish {
		return group.LabelEN
	}
	return group.LabelJA
}

func main() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
	darwinApp = newDarwinPetApp()
	icon := darwinApp.statusIconPNG()
	if len(icon) > 0 {
		C.startAnimalsDesktopApp(C.int(sceneH), (*C.uchar)(unsafe.Pointer(&icon[0])), C.int(len(icon)))
		runtime.KeepAlive(icon)
		return
	}
	C.startAnimalsDesktopApp(C.int(sceneH), nil, 0)
}

func newDarwinPetApp() *darwinPetApp {
	app := &darwinPetApp{
		sceneW:        900,
		speed:         darwinSpeedNormal,
		petCount:      5,
		lang:          darwinLangJapanese,
		mode:          darwinModeRandom,
		coatMode:      darwinCoatSelected,
		selectedCoats: defaultDarwinSelectedCoats(),
		petSizes:      defaultDarwinPetSizes(),
		wheelEnabled:  true,
		frames:        loadDarwinSprites(),
		wheel:         loadDarwinWheel(),
	}
	app.loadSettings()
	return app
}

//export goAnimalsDesktopSetSceneWidth
func goAnimalsDesktopSetSceneWidth(width C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	darwinApp.setSceneWidth(int(width))
}

//export goAnimalsDesktopKeyDown
func goAnimalsDesktopKeyDown() {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.keyHold = wheelKeyHold
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetSpeed
func goAnimalsDesktopSetSpeed(speed C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setSpeed(int(speed))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetPetCount
func goAnimalsDesktopSetPetCount(count C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setPetCount(int(count))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetWheelEnabled
func goAnimalsDesktopSetWheelEnabled(enabled C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setWheelEnabled(enabled != 0)
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetMode
func goAnimalsDesktopSetMode(mode C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setMode(int(mode))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetLanguage
func goAnimalsDesktopSetLanguage(lang C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.lang = normalizeDarwinLanguage(int(lang))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetDisplayID
func goAnimalsDesktopSetDisplayID(displayID C.longlong) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setDisplayID(int64(displayID))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetCoatMode
func goAnimalsDesktopSetCoatMode(mode C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setCoatMode(int(mode))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetVariant
func goAnimalsDesktopSetVariant(variant C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setFixedVariant(int(variant))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetSelectedCoat
func goAnimalsDesktopSetSelectedCoat(index C.int, variant C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setSelectedVariant(int(index), int(variant))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetNameLabels
func goAnimalsDesktopSetNameLabels(enabled C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.nameLabels = enabled != 0
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetPetName
func goAnimalsDesktopSetPetName(index C.int, value *C.char) {
	if darwinApp == nil {
		return
	}
	name := ""
	if value != nil {
		name = C.GoString(value)
	}
	darwinApp.mu.Lock()
	darwinApp.setPetName(int(index), name)
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopSetPetSize
func goAnimalsDesktopSetPetSize(index C.int, percent C.int) {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	darwinApp.setPetSize(int(index), int(percent))
	darwinApp.saveSettings()
	darwinApp.mu.Unlock()
}

//export goAnimalsDesktopClick
func goAnimalsDesktopClick(x C.int, y C.int) C.int {
	if darwinApp == nil {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	if darwinApp.addClickReaction(int(x), int(y)) {
		return C.int(1)
	}
	return C.int(0)
}

//export goAnimalsDesktopPetAt
func goAnimalsDesktopPetAt(x C.int, y C.int) C.int {
	if darwinApp == nil {
		return C.int(-1)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.petAtScenePoint(int(x), int(y)))
}

//export goAnimalsDesktopGetSpeed
func goAnimalsDesktopGetSpeed() C.int {
	if darwinApp == nil {
		return C.int(darwinSpeedNormal)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.speed)
}

//export goAnimalsDesktopGetPetCount
func goAnimalsDesktopGetPetCount() C.int {
	if darwinApp == nil {
		return C.int(5)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.petCount)
}

//export goAnimalsDesktopGetWheelEnabled
func goAnimalsDesktopGetWheelEnabled() C.int {
	if darwinApp == nil {
		return C.int(1)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	if darwinApp.wheelEnabled {
		return C.int(1)
	}
	return C.int(0)
}

//export goAnimalsDesktopGetMode
func goAnimalsDesktopGetMode() C.int {
	if darwinApp == nil {
		return C.int(darwinModeRandom)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.mode)
}

//export goAnimalsDesktopGetLanguage
func goAnimalsDesktopGetLanguage() C.int {
	if darwinApp == nil {
		return C.int(darwinLangJapanese)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.lang)
}

//export goAnimalsDesktopGetDisplayID
func goAnimalsDesktopGetDisplayID() C.longlong {
	if darwinApp == nil {
		return C.longlong(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.longlong(darwinApp.displayID)
}

//export goAnimalsDesktopGetCoatMode
func goAnimalsDesktopGetCoatMode() C.int {
	if darwinApp == nil {
		return C.int(darwinCoatRandom)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.coatMode)
}

//export goAnimalsDesktopGetVariant
func goAnimalsDesktopGetVariant() C.int {
	if darwinApp == nil {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.variant)
}

//export goAnimalsDesktopGetSelectedCoat
func goAnimalsDesktopGetSelectedCoat(index C.int) C.int {
	if darwinApp == nil {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	i := int(index)
	if i < 0 || i >= maxPetCount {
		return C.int(0)
	}
	return C.int(darwinApp.selectedCoats[i])
}

//export goAnimalsDesktopGetVariantCount
func goAnimalsDesktopGetVariantCount() C.int {
	return C.int(len(darwinVariants))
}

//export goAnimalsDesktopCopyVariantLabel
func goAnimalsDesktopCopyVariantLabel(index C.int, buffer *C.char, length C.int) C.int {
	lang := darwinLangJapanese
	if darwinApp != nil {
		darwinApp.mu.Lock()
		lang = darwinApp.lang
		darwinApp.mu.Unlock()
	}
	return copyDarwinCString(darwinVariantLabel(int(index), lang), buffer, length)
}

//export goAnimalsDesktopCopyVariantGroupLabel
func goAnimalsDesktopCopyVariantGroupLabel(index C.int, buffer *C.char, length C.int) C.int {
	lang := darwinLangJapanese
	if darwinApp != nil {
		darwinApp.mu.Lock()
		lang = darwinApp.lang
		darwinApp.mu.Unlock()
	}
	return copyDarwinCString(darwinVariantGroupLabel(int(index), lang), buffer, length)
}

//export goAnimalsDesktopGetPetSize
func goAnimalsDesktopGetPetSize(index C.int) C.int {
	if darwinApp == nil {
		return C.int(defaultPetSizePercent)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	return C.int(darwinApp.petSizePercent(int(index)))
}

//export goAnimalsDesktopGetNameLabels
func goAnimalsDesktopGetNameLabels() C.int {
	if darwinApp == nil {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	if darwinApp.nameLabels {
		return C.int(1)
	}
	return C.int(0)
}

//export goAnimalsDesktopCopyPetName
func goAnimalsDesktopCopyPetName(index C.int, buffer *C.char, length C.int) C.int {
	if darwinApp == nil || buffer == nil || length <= 0 {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	i := int(index)
	if i < 0 || i >= maxPetCount {
		return C.int(0)
	}
	name := sanitizeDarwinPetName(darwinApp.petNames[i])
	return copyDarwinCString(name, buffer, length)
}

func copyDarwinCString(value string, buffer *C.char, length C.int) C.int {
	if buffer == nil || length <= 0 {
		return C.int(0)
	}
	data := []byte(value)
	maxLen := int(length) - 1
	if maxLen < 0 {
		return C.int(0)
	}
	if len(data) > maxLen {
		data = data[:maxLen]
	}
	dst := unsafe.Slice((*byte)(unsafe.Pointer(buffer)), int(length))
	copy(dst, data)
	dst[len(data)] = 0
	return C.int(len(data))
}

//export goAnimalsDesktopGetPetDrawX
func goAnimalsDesktopGetPetDrawX(index C.int) C.int {
	if darwinApp == nil {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	i := int(index)
	if i < 0 || i >= len(darwinApp.pets) {
		return C.int(0)
	}
	return C.int(darwinApp.pets[i].x)
}

//export goAnimalsDesktopGetPetDrawY
func goAnimalsDesktopGetPetDrawY(index C.int) C.int {
	if darwinApp == nil {
		return C.int(0)
	}
	darwinApp.mu.Lock()
	defer darwinApp.mu.Unlock()
	i := int(index)
	if i < 0 || i >= len(darwinApp.pets) {
		return C.int(0)
	}
	_, h := darwinApp.petSpriteSize(i)
	return C.int(sceneH - h - darwinApp.pets[i].lane)
}

//export goAnimalsDesktopTick
func goAnimalsDesktopTick() {
	if darwinApp == nil {
		return
	}
	darwinApp.mu.Lock()
	if !darwinApp.tickPets() {
		darwinApp.mu.Unlock()
		return
	}
	img := darwinApp.render()
	darwinApp.mu.Unlock()

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return
	}
	data := buf.Bytes()
	if len(data) == 0 {
		return
	}
	C.updateAnimalsDesktopImage((*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)), C.int(img.Bounds().Dx()), C.int(img.Bounds().Dy()))
	runtime.KeepAlive(data)
}

func defaultDarwinSettingsPath() (string, error) {
	if dir := os.Getenv("ANIMALSDESKTOP_SETTINGS_DIR"); dir != "" {
		return filepath.Join(dir, settingsFileName), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsDirName, settingsFileName), nil
}

func (a *darwinPetApp) loadSettings() {
	path, err := darwinSettingsPath()
	if err != nil {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var settings darwinSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return
	}
	if settings.Speed != 0 {
		a.speed = normalizeDarwinSpeed(settings.Speed)
	}
	if settings.PetCount != 0 {
		a.petCount = normalizeDarwinPetCount(settings.PetCount)
	}
	a.lang = normalizeDarwinLanguage(settings.Language)
	a.displayID = normalizeDarwinDisplayID(settings.DisplayID)
	if settings.Version >= 3 && settings.Version <= darwinSettingsVersion {
		if settings.Variant != nil || settings.VariantID != "" {
			legacyVariant := a.variant
			if settings.Variant != nil {
				legacyVariant = *settings.Variant
			}
			a.variant = darwinVariantIndexFromSettings(legacyVariant, settings.VariantID, settings.Version, strings.Join(settings.PetNames, "\n"))
		}
		if settings.CoatMode != nil {
			a.coatMode = normalizeDarwinCoatMode(*settings.CoatMode)
		}
		if len(settings.SelectedCoats) > 0 || len(settings.SelectedCoatIDs) > 0 {
			for i := 0; i < min(maxPetCount, max(len(settings.SelectedCoats), len(settings.SelectedCoatIDs))); i++ {
				variant := a.selectedCoats[i]
				if i < len(settings.SelectedCoats) {
					variant = settings.SelectedCoats[i]
				}
				a.selectedCoats[i] = darwinVariantIndexFromSettings(variant, stringAt(settings.SelectedCoatIDs, i), settings.Version, stringAt(settings.PetNames, i))
			}
		}
		if len(settings.PetSizes) > 0 {
			for i, size := range settings.PetSizes {
				if i >= maxPetCount {
					break
				}
				a.petSizes[i] = normalizeDarwinPetSizePercent(size)
			}
		}
	}
	if settings.Mode != nil {
		a.mode = normalizeDarwinMode(*settings.Mode)
	}
	if settings.WheelEnabled != nil {
		a.wheelEnabled = *settings.WheelEnabled
	}
	a.nameLabels = settings.NameLabels
	for i, name := range settings.PetNames {
		if i >= maxPetCount {
			break
		}
		a.petNames[i] = sanitizeDarwinPetName(name)
	}
}

func (a *darwinPetApp) saveSettings() {
	path, err := darwinSettingsPath()
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	wheelEnabled := a.wheelEnabled
	selectedCoats := make([]int, maxPetCount)
	copy(selectedCoats, a.selectedCoats[:])
	for i := range selectedCoats {
		selectedCoats[i] = normalizeDarwinVariant(selectedCoats[i])
	}
	selectedCoatIDs := darwinVariantIDsForIndices(selectedCoats)
	petSizes := make([]int, maxPetCount)
	for i := range a.petSizes {
		petSizes[i] = normalizeDarwinPetSizePercent(a.petSizes[i])
	}
	petNames := make([]string, maxPetCount)
	for i := range a.petNames {
		petNames[i] = sanitizeDarwinPetName(a.petNames[i])
	}
	variant := normalizeDarwinVariant(a.variant)
	coatMode := int(normalizeDarwinCoatMode(int(a.coatMode)))
	mode := int(normalizeDarwinMode(int(a.mode)))
	settings := darwinSettings{
		Version:         darwinSettingsVersion,
		Variant:         &variant,
		VariantID:       darwinVariantIDAt(variant),
		CoatMode:        &coatMode,
		SelectedCoats:   selectedCoats,
		SelectedCoatIDs: selectedCoatIDs,
		Speed:           normalizeDarwinSpeed(a.speed),
		Language:        int(normalizeDarwinLanguage(int(a.lang))),
		DisplayID:       normalizeDarwinDisplayID(a.displayID),
		Mode:            &mode,
		PetSizes:        petSizes,
		PetCount:        normalizeDarwinPetCount(a.petCount),
		WheelEnabled:    &wheelEnabled,
		NameLabels:      a.nameLabels,
		PetNames:        petNames,
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, append(data, '\n'), 0o644)
}

func darwinVariantIndexFromSettings(index int, id string, version int, nameHint string) int {
	if len(darwinVariants) == 0 {
		return 0
	}
	if i, ok := darwinVariantIndexByID(id); ok {
		return i
	}
	if version < darwinSettingsVersion {
		if legacyID, ok := darwinLegacyVariantIDForIndex(index, nameHint); ok {
			if i, found := darwinVariantIndexByID(legacyID); found {
				return i
			}
		}
	}
	return normalizeDarwinVariant(index)
}

func darwinLegacyVariantIDForIndex(index int, nameHint string) (string, bool) {
	if strings.Contains(nameHint, "ハシビロコウ") || strings.Contains(strings.ToLower(nameHint), "shoebill") {
		return "shoebill_stork", true
	}
	return "", false
}

func darwinVariantIndexByID(id string) (int, bool) {
	if id == "" {
		return 0, false
	}
	for i, variant := range darwinVariants {
		if variant.ID == id {
			return i, true
		}
	}
	return 0, false
}

func darwinVariantIDAt(index int) string {
	if len(darwinVariants) == 0 {
		return ""
	}
	return darwinVariants[normalizeDarwinVariant(index)].ID
}

func darwinVariantIDsForIndices(indices []int) []string {
	out := make([]string, len(indices))
	for i, index := range indices {
		out[i] = darwinVariantIDAt(index)
	}
	return out
}

func stringAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}

func (a *darwinPetApp) setSceneWidth(width int) {
	if width < 320 {
		width = 320
	}
	a.sceneW = width
	if len(a.pets) == 0 {
		a.resetPets()
		return
	}
	for i := range a.pets {
		w, _ := a.petSpriteSize(i)
		a.pets[i].x = clamp(a.pets[i].x, 0, max(0, a.sceneW-w))
	}
}

func (a *darwinPetApp) resetPets() {
	count := normalizeDarwinPetCount(a.petCount)
	spacing := max(spriteW+28, a.sceneW/(count+1))
	a.pets = make([]darwinPet, count)
	for i := range a.pets {
		dir := 1
		if i%2 == 1 {
			dir = -1
		}
		w, _ := a.petSpriteSize(i)
		a.pets[i] = darwinPet{
			x:         clamp(spacing*(i+1)-w/2, 0, max(0, a.sceneW-w)),
			lane:      (i % 3) * 7,
			dir:       dir,
			speed:     a.petSpeed(i),
			variant:   a.variantForIndex(i),
			nextPause: 90 + rand.Intn(90),
		}
	}
}

func (a *darwinPetApp) setSpeed(speed int) {
	a.speed = normalizeDarwinSpeed(speed)
	for i := range a.pets {
		a.pets[i].speed = a.petSpeed(i)
	}
}

func (a *darwinPetApp) setPetCount(count int) {
	a.petCount = normalizeDarwinPetCount(count)
	a.resetPets()
	a.tickReactions()
}

func (a *darwinPetApp) setMode(mode int) {
	a.mode = normalizeDarwinMode(mode)
	if a.mode == darwinModeKeyboard {
		for i := range a.pets {
			a.pets[i].pause = 24
			a.pets[i].nextPause = 80 + rand.Intn(90)
		}
	}
}

func (a *darwinPetApp) setCoatMode(mode int) {
	a.coatMode = normalizeDarwinCoatMode(mode)
	a.refreshPetVariants()
}

func (a *darwinPetApp) setFixedVariant(variant int) {
	a.variant = normalizeDarwinVariant(variant)
	if a.coatMode == darwinCoatFixed {
		a.refreshPetVariants()
	}
}

func (a *darwinPetApp) setSelectedVariant(index int, variant int) {
	if index < 0 || index >= maxPetCount {
		return
	}
	a.selectedCoats[index] = normalizeDarwinVariant(variant)
	if a.coatMode == darwinCoatSelected && index < len(a.pets) {
		a.pets[index].variant = a.selectedCoats[index]
	}
}

func (a *darwinPetApp) setWheelEnabled(enabled bool) {
	a.wheelEnabled = enabled
	if !enabled {
		a.keyHold = 0
	}
}

func (a *darwinPetApp) setPetName(index int, name string) {
	if index < 0 || index >= maxPetCount {
		return
	}
	a.petNames[index] = sanitizeDarwinPetName(name)
}

func (a *darwinPetApp) setPetSize(index int, percent int) {
	if index < 0 || index >= maxPetCount {
		return
	}
	a.petSizes[index] = normalizeDarwinPetSizePercent(percent)
	if index < len(a.pets) {
		w, _ := a.petSpriteSize(index)
		a.pets[index].x = clamp(a.pets[index].x, 0, max(0, a.sceneW-w))
	}
}

func (a *darwinPetApp) setDisplayID(displayID int64) {
	a.displayID = normalizeDarwinDisplayID(displayID)
}

func (a *darwinPetApp) petDisplayName(index int) string {
	if index < 0 || index >= maxPetCount {
		return ""
	}
	if name := sanitizeDarwinPetName(a.petNames[index]); name != "" {
		return name
	}
	if a.lang == darwinLangEnglish {
		return fmt.Sprintf("Animal %d", index+1)
	}
	return fmt.Sprintf("どうぶつ%d", index+1)
}

func (a *darwinPetApp) petSpeed(index int) int {
	return max(1, a.speed-1+index%2)
}

func (a *darwinPetApp) refreshPetVariants() {
	for i := range a.pets {
		a.pets[i].variant = a.variantForIndex(i)
	}
}

func (a *darwinPetApp) variantForIndex(index int) int {
	if len(darwinVariants) == 0 {
		return 0
	}
	switch a.coatMode {
	case darwinCoatRandom:
		return rand.Intn(len(darwinVariants))
	case darwinCoatSelected:
		if index >= 0 && index < maxPetCount {
			return normalizeDarwinVariant(a.selectedCoats[index])
		}
	}
	return normalizeDarwinVariant(a.variant)
}

func (a *darwinPetApp) variantID(index int) string {
	index = normalizeDarwinVariant(index)
	return darwinVariants[index].ID
}

func (a *darwinPetApp) petWheelCapable(index int) bool {
	if index < 0 || index >= len(a.pets) || len(darwinVariants) == 0 {
		return false
	}
	variant := normalizeDarwinVariant(a.pets[index].variant)
	if timedVariant(darwinVariants[variant].ID) != nil {
		return false
	}
	return darwinVariants[variant].WheelCapable
}

func defaultDarwinSelectedCoats() [maxPetCount]int {
	return [maxPetCount]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
}

func defaultDarwinPetSizes() [maxPetCount]int {
	var sizes [maxPetCount]int
	for i := range sizes {
		sizes[i] = defaultPetSizePercent
	}
	return sizes
}

func normalizeDarwinPetSizePercent(size int) int {
	if size == 0 {
		return defaultPetSizePercent
	}
	size = clamp(size, minPetSizePercent, maxPetSizePercent)
	return ((size + petSizeStepPercent/2) / petSizeStepPercent) * petSizeStepPercent
}

func darwinPetSpriteSizeForPercent(percent int) (int, int) {
	percent = normalizeDarwinPetSizePercent(percent)
	return max(1, frameW*percent/100), max(1, frameH*percent/100)
}

func (a *darwinPetApp) petSizePercent(index int) int {
	if index < 0 || index >= len(a.petSizes) {
		return defaultPetSizePercent
	}
	return normalizeDarwinPetSizePercent(a.petSizes[index])
}

func (a *darwinPetApp) petSpriteSize(index int) (int, int) {
	return darwinPetSpriteSizeForPercent(a.petSizePercent(index))
}

func normalizeDarwinSpeed(speed int) int {
	switch speed {
	case darwinSpeedSlow, darwinSpeedNormal, darwinSpeedFast:
		return speed
	default:
		return darwinSpeedNormal
	}
}

func normalizeDarwinPetCount(count int) int {
	return clamp(count, 1, maxPetCount)
}

func normalizeDarwinMode(mode int) darwinBehaviorMode {
	switch darwinBehaviorMode(mode) {
	case darwinModeKeyboard, darwinModeRandom:
		return darwinBehaviorMode(mode)
	default:
		return darwinModeRandom
	}
}

func normalizeDarwinCoatMode(mode int) darwinCoatMode {
	switch darwinCoatMode(mode) {
	case darwinCoatFixed, darwinCoatSelected, darwinCoatRandom:
		return darwinCoatMode(mode)
	default:
		return darwinCoatRandom
	}
}

func normalizeDarwinLanguage(lang int) darwinLanguage {
	switch darwinLanguage(lang) {
	case darwinLangJapanese, darwinLangEnglish:
		return darwinLanguage(lang)
	default:
		return darwinLangJapanese
	}
}

func normalizeDarwinDisplayID(displayID int64) int64 {
	if displayID < 0 {
		return 0
	}
	return displayID
}

func normalizeDarwinVariant(variant int) int {
	return clamp(variant, 0, max(0, len(darwinVariants)-1))
}

func sanitizeDarwinPetName(name string) string {
	name = strings.Join(strings.Fields(name), " ")
	runes := []rune(name)
	if len(runes) > maxPetNameLen {
		runes = runes[:maxPetNameLen]
	}
	return string(runes)
}

func (a *darwinPetApp) movePet(p *darwinPet, speed int) {
	index := -1
	for i := range a.pets {
		if &a.pets[i] == p {
			index = i
			break
		}
	}
	w := spriteW
	if index >= 0 {
		w, _ = a.petSpriteSize(index)
	}
	p.x += p.dir * max(1, speed)
	if p.x <= 0 {
		p.x = 0
		p.dir = 1
	}
	if p.x >= max(0, a.sceneW-w) {
		p.x = max(0, a.sceneW-w)
		p.dir = -1
	}
}

func (a *darwinPetApp) tickPets() bool {
	dt := a.motionClock.Step(time.Now())
	a.legacyElapsed += dt
	legacyTick := a.legacyElapsed >= float64(timerInterval)/1000
	if legacyTick {
		a.legacyElapsed = math.Mod(a.legacyElapsed, float64(timerInterval)/1000)
	}
	timedActive := false
	for i := range a.pets {
		p := &a.pets[i]
		if v := timedVariant(a.variantID(p.variant)); v != nil {
			timedActive = true
			w, h := a.petSpriteSize(i)
			if p.timed.ID != v.ID {
				p.timed.Reset(v.ID, float64(p.x), p.dir < 0)
			}
			if a.keyHold > 0 && i == 0 {
				p.timed.ReactWith(v)
			}
			tickTimedPlayer(&p.timed, v, dt, p.x, h, max(0, a.sceneW-w), p.dir, float64(a.speed)/3, a.mode == darwinModeRandom, false)
			p.x = int(math.Round(p.timed.X))
			p.dir = 1
			if p.timed.Left {
				p.dir = -1
			}
		}
	}
	if !legacyTick {
		return timedActive
	}
	a.tick++
	if a.keyHold > 0 {
		a.keyHold--
	}

	for i := range a.pets {
		p := &a.pets[i]
		if timedVariant(a.variantID(p.variant)) != nil {
			continue
		}
		if a.keyHold > 0 && i == 0 {
			p.frame = seqFrameFrom(walkFrameSeq, a.tick, 1)
			if !a.wheelEnabled || !a.petWheelCapable(i) {
				a.movePet(p, p.speed+1)
			}
			continue
		}
		if a.mode == darwinModeKeyboard {
			p.frame = seqFrameFrom(idleFrameSeq, a.tick, 5)
			continue
		}
		if p.pause > 0 {
			p.pause--
			variantID := a.variantID(p.variant)
			if darwinMotionProfileForVariantID(variantID) == catalog.MotionProfileFerretSlink {
				p.frame = darwinRandomPauseFrame(variantID, p.pauseAction, a.tick)
			} else {
				p.frame = seqFrameFrom(idleFrameSeq, a.tick, 5)
			}
			continue
		}
		p.nextPause--
		if p.nextPause <= 0 {
			p.pause = 30 + rand.Intn(70)
			p.nextPause = 120 + rand.Intn(180)
			p.pauseAction = rand.Intn(3)
			p.frame = darwinRandomPauseFrame(a.variantID(p.variant), p.pauseAction, a.tick)
			continue
		}
		a.movePet(p, p.speed)
		p.frame = seqFrameFrom(walkFrameSeq, a.tick, 2)
	}
	a.tickReactions()
	return true
}

func (a *darwinPetApp) tickReactions() {
	out := a.reactions[:0]
	for _, reaction := range a.reactions {
		reaction.ticks--
		if reaction.ticks > 0 && reaction.pet >= 0 && reaction.pet < len(a.pets) {
			out = append(out, reaction)
		}
	}
	a.reactions = out
}

func (a *darwinPetApp) addClickReaction(sceneX, sceneY int) bool {
	index := a.petAtScenePoint(sceneX, sceneY)
	if index < 0 {
		return false
	}
	p := &a.pets[index]
	if v := timedVariant(a.variantID(p.variant)); v != nil {
		if p.timed.ID != v.ID {
			p.timed.Reset(v.ID, float64(p.x), p.dir < 0)
		}
		p.timed.ReactWith(v)
	}
	kind := rand.Intn(3)
	for i := range a.reactions {
		if a.reactions[i].pet == index {
			a.reactions[i].kind = kind
			a.reactions[i].ticks = reactionTicks
			return true
		}
	}
	a.reactions = append(a.reactions, darwinReaction{pet: index, kind: kind, ticks: reactionTicks})
	return true
}

func (a *darwinPetApp) petAtScenePoint(sceneX, sceneY int) int {
	if sceneX < 0 || sceneX >= a.sceneW || sceneY < 0 || sceneY >= sceneH {
		return -1
	}
	wheelActive := a.wheelEnabled && a.keyHold > 0 && len(a.pets) > 0 && a.petWheelCapable(0)
	for i := len(a.pets) - 1; i >= 0; i-- {
		if wheelActive && i == 0 {
			continue
		}
		p := a.pets[i]
		w, h := a.petSpriteSize(i)
		y := sceneH - h - p.lane
		insetX := max(2, w/14)
		insetTop := max(2, h/8)
		insetBottom := max(2, h/16)
		if sceneX >= p.x+insetX && sceneX <= p.x+w-insetX && sceneY >= y+insetTop && sceneY <= y+h-insetBottom {
			return i
		}
	}
	return -1
}

func (a *darwinPetApp) render() *image.RGBA {
	active := map[string]bool{}
	for _, p := range a.pets {
		active[a.variantID(p.variant)] = true
	}
	timedStore.Retain(active)
	w := max(320, a.sceneW)
	canvas := image.NewRGBA(image.Rect(0, 0, w, sceneH))
	draw.Draw(canvas, canvas.Bounds(), image.Transparent, image.Point{}, draw.Src)

	wheelActive := a.wheelEnabled && a.keyHold > 0 && len(a.pets) > 0 && a.petWheelCapable(0)
	if wheelActive {
		wheelX := clamp(w-116, 8, max(8, w-wheelSize-8))
		wheelY := sceneH - wheelSize - 4
		if a.wheel != nil {
			draw.Draw(canvas, image.Rect(wheelX, wheelY, wheelX+wheelSize, wheelY+wheelSize), a.wheel, image.Point{}, draw.Over)
		}
		variantID := a.variantID(a.pets[0].variant)
		if frames := a.frames[variantID]; len(frames) > a.pets[0].frame {
			runner := scaleNearest(frames[a.pets[0].frame], 66, 44)
			drawFacingImage(canvas, runner, image.Rect(wheelX+3, wheelY+22, wheelX+69, wheelY+66), darwinDrawDirection(1, variantID))
		}
	}

	for i := range a.pets {
		if wheelActive && i == 0 {
			continue
		}
		p := &a.pets[i]
		variantID := a.variantID(p.variant)
		if v := timedVariant(variantID); v != nil {
			if p.timed.ID != v.ID {
				p.timed.Reset(v.ID, float64(p.x), p.dir < 0)
			}
			w, h := a.petSpriteSize(i)
			v.Draw(canvas, p.timed.Action, p.timed.Phase(), p.x, sceneH-p.lane-3, w, h, p.dir)
			continue
		}
		frames := a.frames[variantID]
		if len(frames) <= p.frame {
			continue
		}
		w, h := a.petSpriteSize(i)
		y := sceneH - h - p.lane
		rect := renderAdjustmentForFrame(variantID, p.frame).adjustRect(image.Rect(p.x, y, p.x+w, y+h))
		drawFacingImage(canvas, frames[p.frame], rect, darwinDrawDirection(p.dir, variantID))
	}
	a.drawReactions(canvas)
	return canvas
}

func (a *darwinPetApp) drawReactions(dst *image.RGBA) {
	for _, reaction := range a.reactions {
		if reaction.pet < 0 || reaction.pet >= len(a.pets) {
			continue
		}
		wheelActive := a.wheelEnabled && a.keyHold > 0 && reaction.pet == 0 && a.petWheelCapable(0)
		if wheelActive {
			continue
		}
		p := a.pets[reaction.pet]
		w, h := a.petSpriteSize(reaction.pet)
		baseY := sceneH - h - p.lane
		x := clamp(p.x+w/2-18, 2, max(2, a.sceneW-42))
		y := clamp(baseY-26-(reactionTicks-reaction.ticks)/8, 0, sceneH-32)
		drawReactionBubble(dst, x, y, reaction.kind, reaction.ticks)
	}
}

func (a *darwinPetApp) statusIconPNG() []byte {
	if len(darwinVariants) == 0 {
		return nil
	}
	frames := a.frames[darwinVariants[0].ID]
	if len(frames) == 0 {
		return nil
	}
	visible := cropVisible(frames[idleStart])
	if visible.Bounds().Empty() {
		return nil
	}

	const iconW = 22
	const iconH = 18
	vb := visible.Bounds()
	targetW := iconW
	targetH := max(1, vb.Dy()*targetW/vb.Dx())
	if targetH > iconH {
		targetH = iconH
		targetW = max(1, vb.Dx()*targetH/vb.Dy())
	}
	scaled := scaleNearest(visible, targetW, targetH)
	icon := image.NewRGBA(image.Rect(0, 0, iconW, iconH))
	draw.Draw(icon, icon.Bounds(), image.Transparent, image.Point{}, draw.Src)
	draw.Draw(icon, image.Rect((iconW-targetW)/2, (iconH-targetH)/2, (iconW+targetW)/2, (iconH+targetH)/2), scaled, image.Point{}, draw.Over)

	var buf bytes.Buffer
	if err := png.Encode(&buf, icon); err != nil {
		return nil
	}
	return buf.Bytes()
}

func seqFrameFrom(seq []int, tick, delay int) int {
	if len(seq) == 0 {
		return idleStart
	}
	if delay <= 0 {
		delay = 1
	}
	return seq[(tick/delay)%len(seq)]
}

func darwinRandomPauseFrame(variantID string, action int, tick int) int {
	if darwinMotionProfileForVariantID(variantID) == catalog.MotionProfileFerretSlink {
		seq, delay := darwinFerretPauseSequence(action)
		return seqFrameFrom(seq, tick, delay)
	}

	switch action {
	case 0:
		seq, delay := darwinNibbleLikeSequence(variantID)
		return seqFrameFrom(seq, tick, delay)
	case 1:
		return seqFrameFrom(hopFrameSeq, tick, 2)
	default:
		return seqFrameFrom(idleFrameSeq, tick, 5)
	}
}

func darwinMotionProfileForVariantID(variantID string) string {
	for _, variant := range darwinVariants {
		if variant.ID == variantID {
			return variant.MotionProfile
		}
	}
	return ""
}

func darwinFerretPauseSequence(action int) ([]int, int) {
	switch action {
	case 0:
		return ferretSniffSeq, 3
	case 1:
		return ferretGroomSeq, 3
	default:
		return ferretRestSeq, 5
	}
}

func darwinNibbleLikeSequence(variantID string) ([]int, int) {
	switch variantID {
	case "sugar_glider_gray", "rabbit_chestnut_agouti":
		return hopFrameSeq, 2
	default:
		return nibbleFrameSeq, 3
	}
}

func loadDarwinSprites() map[string][]*image.RGBA {
	out := make(map[string][]*image.RGBA, len(darwinVariants))
	for _, variant := range darwinVariants {
		name := fmt.Sprintf("sprites/%s_set00.png", variant.SpriteBase)
		data, err := fs.ReadFile(appassets.FS, name)
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		if img.Bounds().Dx() != frameW*frameCount || img.Bounds().Dy() != frameH {
			panic(fmt.Sprintf("%s must be %dx%d", name, frameW*frameCount, frameH))
		}
		frames := make([]*image.RGBA, frameCount)
		for i := 0; i < frameCount; i++ {
			frame := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
			srcRect := image.Rect(i*frameW, 0, (i+1)*frameW, frameH)
			draw.Draw(frame, frame.Bounds(), img, srcRect.Min, draw.Src)
			frames[i] = frame
		}
		out[variant.ID] = frames
	}
	return out
}

func loadDarwinWheel() *image.RGBA {
	data, err := fs.ReadFile(appassets.FS, "sprites/wheel.png")
	if err != nil {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	dst := image.NewRGBA(image.Rect(0, 0, wheelSize, wheelSize))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
	return dst
}

func drawFacingImage(dst *image.RGBA, src *image.RGBA, r image.Rectangle, dir int) {
	if src == nil {
		return
	}
	if dir >= 0 {
		draw.Draw(dst, r, src, image.Point{}, draw.Over)
		return
	}
	flipped := image.NewRGBA(src.Bounds())
	b := src.Bounds()
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			flipped.Set(x, y, src.At(b.Min.X+b.Dx()-1-x, b.Min.Y+y))
		}
	}
	draw.Draw(dst, r, flipped, image.Point{}, draw.Over)
}

func darwinDrawDirection(desiredDir int, variantID string) int {
	dir := 1
	if desiredDir < 0 {
		dir = -1
	}
	return dir * darwinSourceFacingDirection(variantID)
}

func darwinSourceFacingDirection(variantID string) int {
	switch variantID {
	case "sugar_glider_gray":
		return -1
	default:
		return 1
	}
}

func scaleNearest(src *image.RGBA, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	if src == nil || width <= 0 || height <= 0 {
		return dst
	}
	sb := src.Bounds()
	for y := 0; y < height; y++ {
		sy := sb.Min.Y + y*sb.Dy()/height
		for x := 0; x < width; x++ {
			sx := sb.Min.X + x*sb.Dx()/width
			dst.SetRGBA(x, y, src.RGBAAt(sx, sy))
		}
	}
	return dst
}

func cropVisible(src *image.RGBA) *image.RGBA {
	if src == nil {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	b := src.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if src.RGBAAt(x, y).A <= 8 {
				continue
			}
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x+1 > maxX {
				maxX = x + 1
			}
			if y+1 > maxY {
				maxY = y + 1
			}
		}
	}
	if minX >= maxX || minY >= maxY {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	content := image.Rect(max(b.Min.X, minX-1), max(b.Min.Y, minY-1), min(b.Max.X, maxX+1), min(b.Max.Y, maxY+1))
	dst := image.NewRGBA(image.Rect(0, 0, content.Dx(), content.Dy()))
	draw.Draw(dst, dst.Bounds(), src, content.Min, draw.Src)
	return dst
}

func drawReactionBubble(dst *image.RGBA, x, y, kind, ticks int) {
	alpha := uint8(235)
	if ticks < 12 {
		alpha = uint8(max(0, ticks) * 235 / 12)
	}
	bg := rgba(255, 255, 248, alpha)
	edge := rgba(54, 89, 70, uint8(min(180, int(alpha))))
	shadow := rgba(22, 28, 24, uint8(min(70, int(alpha)/2)))
	drawRoundedRect(dst, x+2, y+3, 38, 24, 7, shadow)
	drawRoundedRect(dst, x, y, 38, 24, 7, bg)
	drawRoundedRectOutline(dst, x, y, 38, 24, 7, edge)
	drawBubbleTail(dst, x+17, y+23, alpha)
	switch kind % 3 {
	case 0:
		drawHeartIcon(dst, x+13, y+6, rgba(219, 72, 92, alpha))
	case 1:
		drawSmileIcon(dst, x+12, y+5, rgba(239, 184, 68, alpha), rgba(67, 62, 45, alpha))
	default:
		drawSparkleIcon(dst, x+11, y+5, rgba(70, 128, 104, alpha), rgba(239, 184, 68, alpha))
	}
}

func drawBubbleTail(dst *image.RGBA, x, y int, alpha uint8) {
	c := rgba(255, 255, 248, alpha)
	for row := 0; row < 6; row++ {
		for col := -row; col <= row; col++ {
			pt := image.Pt(x+col, y+row)
			if pt.In(dst.Bounds()) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), c))
			}
		}
	}
}

func drawHeartIcon(dst *image.RGBA, x, y int, c color.RGBA) {
	fillCircle(dst, x+5, y+4, 4, c)
	fillCircle(dst, x+12, y+4, 4, c)
	for row := 4; row < 16; row++ {
		half := max(0, 9-row/2)
		for px := x + 8 - half; px <= x+8+half; px++ {
			pt := image.Pt(px, y+row)
			if pt.In(dst.Bounds()) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), c))
			}
		}
	}
}

func drawSmileIcon(dst *image.RGBA, x, y int, face, ink color.RGBA) {
	fillCircle(dst, x+8, y+8, 8, face)
	fillCircle(dst, x+5, y+6, 1, ink)
	fillCircle(dst, x+11, y+6, 1, ink)
	for px := x + 4; px <= x+12; px++ {
		dx := px - (x + 8)
		py := y + 9 + abs(dx)/3
		for sy := 0; sy < 2; sy++ {
			pt := image.Pt(px, py+sy)
			if pt.In(dst.Bounds()) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), ink))
			}
		}
	}
}

func drawSparkleIcon(dst *image.RGBA, x, y int, main, accent color.RGBA) {
	drawDiamond(dst, x+8, y+8, 8, main)
	drawPixelLine(dst, x+8, y, x+8, y+16, main)
	drawPixelLine(dst, x, y+8, x+16, y+8, main)
	drawDiamond(dst, x+22, y+5, 4, accent)
	drawDiamond(dst, x+20, y+15, 3, accent)
}

func drawDiamond(dst *image.RGBA, cx, cy, r int, c color.RGBA) {
	for py := cy - r; py <= cy+r; py++ {
		half := r - abs(py-cy)
		for px := cx - half; px <= cx+half; px++ {
			pt := image.Pt(px, py)
			if pt.In(dst.Bounds()) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), c))
			}
		}
	}
}

func drawRoundedRect(dst *image.RGBA, x, y, w, h, radius int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			pt := image.Pt(px, py)
			if roundedRectContains(px-x, py-y, w, h, radius) && pt.In(dst.Bounds()) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), c))
			}
		}
	}
}

func drawRoundedRectOutline(dst *image.RGBA, x, y, w, h, radius int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			pt := image.Pt(px, py)
			if !pt.In(dst.Bounds()) || !roundedRectContains(px-x, py-y, w, h, radius) {
				continue
			}
			if px == x || px == x+w-1 || py == y || py == y+h-1 ||
				!roundedRectContains(px-x-1, py-y, w, h, radius) ||
				!roundedRectContains(px-x+1, py-y, w, h, radius) ||
				!roundedRectContains(px-x, py-y-1, w, h, radius) ||
				!roundedRectContains(px-x, py-y+1, w, h, radius) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), c))
			}
		}
	}
}

func roundedRectContains(x, y, w, h, radius int) bool {
	if x < 0 || y < 0 || x >= w || y >= h {
		return false
	}
	if (x >= radius && x < w-radius) || (y >= radius && y < h-radius) {
		return true
	}
	cx := radius
	if x >= w-radius {
		cx = w - radius - 1
	}
	cy := radius
	if y >= h-radius {
		cy = h - radius - 1
	}
	dx := x - cx
	dy := y - cy
	return dx*dx+dy*dy <= radius*radius
}

func drawPixelLine(dst *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		drawBlock(dst, x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func fillCircle(dst *image.RGBA, cx, cy, r int, c color.RGBA) {
	for py := cy - r; py <= cy+r; py++ {
		for px := cx - r; px <= cx+r; px++ {
			dx := px - cx
			dy := py - cy
			if dx*dx+dy*dy <= r*r {
				drawBlock(dst, px, py, c)
			}
		}
	}
}

func drawBlock(dst *image.RGBA, x, y int, c color.RGBA) {
	for py := y - 1; py <= y+1; py++ {
		for px := x - 1; px <= x+1; px++ {
			pt := image.Pt(px, py)
			if pt.In(dst.Bounds()) {
				dst.SetRGBA(pt.X, pt.Y, overRGBA(dst.RGBAAt(pt.X, pt.Y), c))
			}
		}
	}
}

func rgba(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: uint8(uint16(r) * uint16(a) / 255), G: uint8(uint16(g) * uint16(a) / 255), B: uint8(uint16(b) * uint16(a) / 255), A: a}
}

func overRGBA(dst, src color.RGBA) color.RGBA {
	a := int(src.A)
	inv := 255 - a
	return color.RGBA{
		R: uint8(int(src.R) + int(dst.R)*inv/255),
		G: uint8(int(src.G) + int(dst.G)*inv/255),
		B: uint8(int(src.B) + int(dst.B)*inv/255),
		A: uint8(a + int(dst.A)*inv/255),
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
