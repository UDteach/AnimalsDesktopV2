//go:build windows

package main

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
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	appassets "animals-desktop/assets"
	"animals-desktop/internal/catalog"
	"animals-desktop/internal/motion"
	"github.com/lxn/win"
)

const (
	appName                     = "Animals Desktop"
	windowClass                 = "AnimalsDesktopPetWindow"
	wmTray                      = win.WM_APP + 1
	wmTyping                    = win.WM_APP + 2
	wmMouseClick                = win.WM_APP + 3
	wmUpdateReady               = win.WM_APP + 4
	wmUpdateFailed              = win.WM_APP + 5
	wmUpdateInstallReady        = win.WM_APP + 6
	timerID                     = 42
	timerInterval               = 55
	frameW                      = 96
	frameH                      = 64
	frameCount                  = 62
	motionSets                  = 10
	scale                       = 1
	spriteW                     = frameW * scale
	spriteH                     = frameH * scale
	forageW                     = 32
	forageH                     = 24
	sceneH                      = 92
	wheelSize                   = 72
	maxPetCount                 = 10
	maxForage                   = 5
	foragePropsEnabled          = false
	defaultOverlayOffsetY       = 10
	minOverlayOffsetY           = -48
	maxOverlayOffsetY           = 96
	overlayOffsetStep           = 4
	defaultWalkRangeStart       = 0
	defaultWalkRangeEnd         = 100
	minWalkRangeSpan            = 25
	walkRangeStep               = 5
	defaultPetSizePercent       = 100
	minPetSizePercent           = 70
	maxPetSizePercent           = 120
	petSizeStepPercent          = 10
	wheelKeyHold                = 18
	turnTicks                   = 16
	reactionTicks               = 54
	settingsClientW       int32 = 760
	settingsClientH       int32 = 620
	renameDialogClientW   int32 = 360
	renameDialogClientH   int32 = 158
	settingsDirName             = "AnimalsDesktop"
	settingsFileName            = "settings.json"
	settingsVersion             = 4
	updaterApplyArg             = "--animalsdesktop-apply-update"
	updaterCleanupArg           = "--animalsdesktop-cleanup"
	updateTempPrefix            = "animals-desktop-update-"
	monitorPrimaryFlag          = 1
)

const (
	whKeyboardLL = 13
	whMouseLL    = 14
)

const (
	idleStart      = 0
	idleFrames     = 4
	walkStart      = 4
	walkFrames     = 8
	scurryStart    = 12
	scurryFrames   = 8
	nibbleStart    = 20
	nibbleFrames   = 6
	hopStart       = 26
	hopFrames      = 6
	turnStart      = 32
	turnFrames     = 8
	eatStart       = 40
	eatFrames      = 4
	digStart       = 44
	digFrames      = 4
	standStart     = 48
	standFrames    = 4
	groomStart     = 52
	groomFrames    = 4
	wheelRunStart  = 56
	wheelRunFrames = 6
)

var (
	idleFrameSeq     = []int{idleStart, idleStart + 1, idleStart + 3, idleStart + 1}
	walkFrameSeq     = []int{walkStart, walkStart + 1, walkStart + 3, walkStart + 1}
	scurryFrameSeq   = []int{scurryStart, scurryStart + 1, scurryStart + 2, scurryStart + 3, scurryStart + 4, scurryStart + 5, scurryStart + 6, scurryStart + 7}
	nibbleFrameSeq   = []int{nibbleStart, nibbleStart + 1, nibbleStart + 2, nibbleStart + 1}
	hopFrameSeq      = []int{hopStart, hopStart + 1, hopStart + 2, hopStart + 3}
	turnFrameSeq     = []int{turnStart, turnStart + 1, turnStart + 2, turnStart + 3, turnStart + 4, turnStart + 5, turnStart + 6, turnStart + 7}
	eatFrameSeq      = []int{eatStart, eatStart + 1, eatStart + 2, eatStart + 3}
	digFrameSeq      = []int{digStart, digStart + 1, digStart + 2, digStart + 3}
	standFrameSeq    = []int{standStart, standStart + 1, standStart + 2, standStart + 3}
	groomFrameSeq    = []int{groomStart, groomStart + 1, groomStart + 2, groomStart + 3}
	wheelRunFrameSeq = []int{wheelRunStart, wheelRunStart + 1, wheelRunStart + 2, wheelRunStart + 3, wheelRunStart + 4, wheelRunStart + 5}
)

const (
	menuExit            uint16 = 100
	menuModeKeyboard    uint16 = 101
	menuModeRandom      uint16 = 102
	menuSpeedSlow       uint16 = 110
	menuSpeedNormal     uint16 = 111
	menuSpeedFast       uint16 = 112
	menuCount1          uint16 = 120
	menuCount2          uint16 = 121
	menuCount3          uint16 = 122
	menuCount5          uint16 = 123
	menuCount10         uint16 = 124
	menuWheelToggle     uint16 = 130
	menuCoatFixed       uint16 = 131
	menuCoatSelected    uint16 = 132
	menuCoatRandom      uint16 = 133
	menuSettings        uint16 = 140
	menuCheckUpdate     uint16 = 150
	menuInstallUpdate   uint16 = 151
	menuLangJapanese    uint16 = 160
	menuLangEnglish     uint16 = 161
	menuHideToggle      uint16 = 170
	menuRandomAll       uint16 = 180
	menuVariantBase     uint16 = 200
	menuRandomGroupBase uint16 = 400
)

const (
	ctrlTabAnimals       int32 = 1000
	ctrlTabMotion        int32 = 1001
	ctrlTabDisplay       int32 = 1009
	ctrlVariantCombo     int32 = 1002
	ctrlPetMinus         int32 = 1003
	ctrlPetPlus          int32 = 1004
	ctrlLanguageCombo    int32 = 1005
	ctrlCoatFixed        int32 = 1006
	ctrlCoatSelected     int32 = 1007
	ctrlCoatRandom       int32 = 1008
	ctrlModeKeyboard     int32 = 1011
	ctrlModeRandom       int32 = 1012
	ctrlSpeedSlow        int32 = 1021
	ctrlSpeedNormal      int32 = 1022
	ctrlSpeedFast        int32 = 1023
	ctrlTypingWheel      int32 = 1031
	ctrlBidirectional    int32 = 1032
	ctrlPositionTaskbar  int32 = 1033
	ctrlPositionBottom   int32 = 1034
	ctrlOffsetUp         int32 = 1035
	ctrlOffsetDown       int32 = 1036
	ctrlReset            int32 = 1041
	ctrlClose            int32 = 1042
	ctrlTopClose         int32 = 1043
	ctrlNameLabels       int32 = 1044
	ctrlPetVariantBase   int32 = 1050
	ctrlPetNameBase      int32 = 1070
	ctrlPetSizeBase      int32 = 1080
	ctrlRenameEdit       int32 = 1100
	ctrlRenameOK         int32 = 1101
	ctrlRenameCancel     int32 = 1102
	ctrlDisplayPrev      int32 = 1110
	ctrlDisplayNext      int32 = 1111
	ctrlRangeFull        int32 = 1112
	ctrlRangeNarrow      int32 = 1113
	ctrlRangeWide        int32 = 1114
	ctrlRangeLeft        int32 = 1115
	ctrlRangeRight       int32 = 1116
	ctrlRangeStartScroll int32 = 1117
	ctrlRangeEndScroll   int32 = 1118
	ctrlDisplaySingle    int32 = 1134
	ctrlDisplaySpan      int32 = 1135
	ctrlDisplaySpanLess  int32 = 1136
	ctrlDisplaySpanMore  int32 = 1137
)

var appVersion = "v2.0.0"

type behaviorMode int

const (
	modeKeyboard behaviorMode = iota
	modeRandom
)

type behaviorState int

const (
	stateIdle behaviorState = iota
	stateWalk
	stateScurry
	stateNibble
	stateHop
	stateWheel
	stateGroom
	stateForage
	stateCarry
	stateTurn
	stateEat
	stateDig
	stateStand
	stateFaceGroom
	stateRest
)

const (
	reservedItem = -2
	noItem       = -1
)

type coatVariant = catalog.Variant

var variants = catalog.RuntimeVariants()

type language int

const (
	langJapanese language = iota
	langEnglish
)

type settingsTab int

const (
	tabAnimals settingsTab = iota
	tabMotion
	tabDisplay
)

type coatMode int

const (
	coatFixed coatMode = iota
	coatSelected
	coatRandom
)

type petSlotMode int

const (
	petSlotFixed petSlotMode = iota
	petSlotRandom
)

type overlayPositionMode int

const (
	positionTaskbarEdge overlayPositionMode = iota
	positionScreenBottom
)

type displayScope int

const (
	displayScopeSingle displayScope = iota
	displayScopeSpan
)

type desktopPet struct {
	timed      motion.Player
	motionSet  int
	variant    int
	frame      int
	x          int
	laneOffset int
	item       int
	carryKind  int
	state      behaviorState
	prevState  behaviorState
	stateTicks int
	moveSpeed  int
	dir        int
	nextDir    int
}

type forageItem struct {
	x      int
	kind   int
	owner  int
	active bool
}

type petReaction struct {
	pet   int
	kind  int
	ticks int
}

type githubRelease struct {
	TagName    string               `json:"tag_name"`
	HTMLURL    string               `json:"html_url"`
	Draft      bool                 `json:"draft"`
	Prerelease bool                 `json:"prerelease"`
	Assets     []githubReleaseAsset `json:"assets"`
}

type githubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
	Size               int64  `json:"size"`
}

type updateState struct {
	mu         sync.Mutex
	latest     *githubRelease
	lastError  string
	checking   atomic.Bool
	installing atomic.Bool
}

type mouseHookStruct struct {
	pt          win.POINT
	mouseData   uint32
	flags       uint32
	time        uint32
	dwExtraInfo uintptr
}

type petApp struct {
	motionClock            motion.Clock
	legacyElapsed          float64
	hwnd                   win.HWND
	hinst                  win.HINSTANCE
	trayIcon               win.HICON
	keyHook                uintptr
	mouseHook              uintptr
	keyHookFailed          bool
	mouseHookFailed        bool
	overlayHwnds           []win.HWND
	frames                 *spriteCache
	forageSprites          []*image.RGBA
	wheel                  *image.RGBA
	pets                   []desktopPet
	forage                 []forageItem
	reactions              []petReaction
	variant                int
	coatMode               coatMode
	selectedCoats          [maxPetCount]int
	selectedSlotModes      [maxPetCount]petSlotMode
	selectedRandomGroupIDs [maxPetCount]string
	randomGroupID          string
	petSizes               [maxPetCount]int
	petNames               [maxPetCount]string
	nameLabels             bool
	speed                  int
	mode                   behaviorMode
	petCount               int
	wheelEnabled           bool
	bidirectional          bool
	positionMode           overlayPositionMode
	overlayOffsetY         int
	displayIndex           int
	displayScope           displayScope
	displaySpanEnd         int
	walkRangeStart         int
	walkRangeEnd           int
	overlayHidden          bool
	settingsHwnd           win.HWND
	settingsTab            settingsTab
	lang                   language
	settingsFont           win.HFONT
	settingsTitleFont      win.HFONT
	settingsSmallFont      win.HFONT
	settingsBrush          win.HBRUSH
	settingsCard           win.HBRUSH
	settingsX              int32
	settingsY              int32
	settingsSaveFailed     bool
	update                 updateState
	nameHwnd               win.HWND
	nameText               string
	hoverPet               int
	renameHwnd             win.HWND
	renameEdit             win.HWND
	renameIndex            int
	wheelX                 int
	sceneW                 int
	tickCount              int
	closing                atomic.Bool
}

type appSettings struct {
	Version                int      `json:"version"`
	Variant                int      `json:"variant"`
	VariantID              string   `json:"variantID,omitempty"`
	CoatMode               int      `json:"coatMode"`
	SelectedCoats          []int    `json:"selectedCoats"`
	SelectedCoatIDs        []string `json:"selectedCoatIDs,omitempty"`
	SelectedSlotModes      []int    `json:"selectedSlotModes,omitempty"`
	SelectedRandomGroupIDs []string `json:"selectedRandomGroupIDs,omitempty"`
	RandomGroupID          string   `json:"randomGroupID,omitempty"`
	PetSizes               []int    `json:"petSizes,omitempty"`
	Speed                  int      `json:"speed"`
	Mode                   int      `json:"mode"`
	PetCount               int      `json:"petCount"`
	WheelEnabled           bool     `json:"wheelEnabled"`
	Bidirectional          bool     `json:"bidirectional"`
	PositionMode           *int     `json:"positionMode,omitempty"`
	VerticalOffset         *int     `json:"verticalOffset,omitempty"`
	DisplayIndex           *int     `json:"displayIndex,omitempty"`
	DisplayScope           *int     `json:"displayScope,omitempty"`
	DisplaySpanEnd         *int     `json:"displaySpanEnd,omitempty"`
	WalkRangeStart         *int     `json:"walkRangeStart,omitempty"`
	WalkRangeEnd           *int     `json:"walkRangeEnd,omitempty"`
	Language               int      `json:"language"`
	SettingsX              int32    `json:"settingsX"`
	SettingsY              int32    `json:"settingsY"`
	NameLabels             bool     `json:"nameLabels"`
	PetNames               []string `json:"petNames,omitempty"`
}

var app *petApp

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	shcore                  = syscall.NewLazyDLL("shcore.dll")
	ntdll                   = syscall.NewLazyDLL("ntdll.dll")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procGetDlgCtrlID        = user32.NewProc("GetDlgCtrlID")
	procSetWindowTextW      = user32.NewProc("SetWindowTextW")
	procSetWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook   = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookExProc  = user32.NewProc("CallNextHookEx")
	procSetProcessDPIAware  = user32.NewProc("SetProcessDPIAware")
	procSetProcessDpiCtx    = user32.NewProc("SetProcessDpiAwarenessContext")
	procUpdateLayeredWin    = user32.NewProc("UpdateLayeredWindow")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetDpiForMonitor    = shcore.NewProc("GetDpiForMonitor")
	procAdjustWindowRectEx  = user32.NewProc("AdjustWindowRectEx")
	procRtlMoveMemory       = ntdll.NewProc("RtlMoveMemory")
)

const (
	acSrcOver       = 0
	defaultDPI      = 96
	monitorDPIType  = 0
	ulwAlpha        = 0x00000002
	spiGetWorkArea  = 0x0030
	sbsHorz         = 0x0000
	sbmSetPos       = 0x00E0
	sbmSetRange     = 0x00E2
	sbLineLeft      = 0
	sbLineRight     = 1
	sbPageLeft      = 2
	sbPageRight     = 3
	sbThumbPosition = 4
	sbThumbTrack    = 5
	sbLeft          = 6
	sbRight         = 7
	sbEndScroll     = 8
)

func enablePerMonitorDPIAwareness() {
	// Keep local developer builds DPI-aware even before a release manifest is embedded.
	if procSetProcessDpiCtx.Find() == nil {
		const dpiAwarenessContextPerMonitorAwareV2 = ^uintptr(3)
		ret, _, _ := procSetProcessDpiCtx.Call(dpiAwarenessContextPerMonitorAwareV2)
		if ret != 0 {
			return
		}
	}
	if procSetProcessDPIAware.Find() == nil {
		procSetProcessDPIAware.Call()
	}
}

func main() {
	if runMotionQA(os.Args[1:]) {
		return
	}
	if networkUpdatesEnabled {
		if runUpdaterUtility(os.Args[1:]) {
			return
		}
	}
	cleanupDir := ""
	if networkUpdatesEnabled {
		cleanupDir = updateCleanupDir(os.Args[1:])
	}

	enablePerMonitorDPIAwareness()
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())

	hinst := win.GetModuleHandle(nil)
	app = &petApp{
		hinst:          hinst,
		frames:         newSpriteCache(),
		forageSprites:  loadForageSprites(),
		wheel:          loadWheelSprite(),
		variant:        0,
		coatMode:       coatSelected,
		selectedCoats:  defaultSelectedCoats(),
		petSizes:       defaultPetSizes(),
		speed:          3,
		mode:           modeRandom,
		petCount:       5,
		wheelEnabled:   true,
		bidirectional:  true,
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: defaultOverlayOffsetY,
		displayIndex:   0,
		displayScope:   displayScopeSingle,
		displaySpanEnd: 0,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
		settingsTab:    tabAnimals,
		lang:           langJapanese,
		settingsX:      120,
		settingsY:      120,
		hoverPet:       -1,
		renameIndex:    -1,
	}
	_ = app.loadSettings()
	if cleanupDir != "" {
		go cleanupUpdateTempDirLater(cleanupDir)
	}

	className := syscall.StringToUTF16Ptr(windowClass)
	wc := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     hinst,
		HIcon:         win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION)),
		HCursor:       win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW)),
		HbrBackground: 0,
		LpszClassName: className,
	}
	if win.RegisterClassEx(&wc) == 0 {
		panic(fmt.Sprintf("RegisterClassEx failed: %v", syscall.GetLastError()))
	}

	app.hwnd = win.CreateWindowEx(
		win.WS_EX_LAYERED|win.WS_EX_TOPMOST|win.WS_EX_TOOLWINDOW|win.WS_EX_TRANSPARENT,
		className,
		syscall.StringToUTF16Ptr(appName),
		win.WS_POPUP,
		0, 0, 1, 1,
		0, 0, hinst, nil,
	)
	if app.hwnd == 0 {
		panic(fmt.Sprintf("CreateWindowEx failed: %v", syscall.GetLastError()))
	}

	app.resetPosition()
	app.installTray()
	app.showStartupToast()
	if networkUpdatesEnabled {
		app.startUpdateCheck(false)
	}
	app.installKeyboardHook()
	app.installMouseHook()
	win.ShowWindow(app.hwnd, win.SW_SHOWNOACTIVATE)
	win.SetTimer(app.hwnd, timerID, 16, 0)
	app.render()

	var msg win.MSG
	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		if app.settingsHwnd != 0 && win.IsWindowVisible(app.settingsHwnd) && win.IsDialogMessage(app.settingsHwnd, &msg) {
			continue
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	app.cleanup()
}

func wndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if app != nil && hwnd == app.nameHwnd {
		return app.nameWndProc(hwnd, msg, wParam, lParam)
	}
	if app != nil && hwnd == app.renameHwnd {
		return app.renameWndProc(hwnd, msg, wParam, lParam)
	}
	if app != nil && hwnd == app.settingsHwnd {
		return app.settingsWndProc(hwnd, msg, wParam, lParam)
	}
	switch msg {
	case win.WM_TIMER:
		if wParam == timerID {
			if app.tick() {
				app.render()
			}
			return 0
		}
	case wmTray:
		if lParam == win.WM_RBUTTONUP || lParam == win.WM_LBUTTONUP {
			app.showTrayMenu()
			return 0
		}
	case wmTyping:
		app.onTyping()
		return 0
	case wmMouseClick:
		x := int(int32(uint32(wParam)))
		y := int(int32(uint32(lParam)))
		app.onMouseClick(x, y)
		return 0
	case wmUpdateReady:
		app.onUpdateReady(wParam != 0)
		return 0
	case wmUpdateFailed:
		app.onUpdateFailed(wParam != 0)
		return 0
	case wmUpdateInstallReady:
		app.onUpdateInstallReady()
		return 0
	case win.WM_COMMAND:
		id := uint16(wParam & 0xffff)
		notify := uint16((wParam >> 16) & 0xffff)
		if app != nil && app.handleSettingsCommand(int32(id), notify) {
			return 0
		}
		if app != nil && app.handleMenuCommand(id) {
			return 0
		}
	case win.WM_DESTROY:
		if hwnd == app.hwnd {
			win.PostQuitMessage(0)
		}
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (a *petApp) resetPosition() {
	overlay := a.overlayRect()
	a.syncScene(overlay)
	a.setPetCount(a.petCount)
	a.arrangePetsForOverlay(overlay)
}

func settingsPath() (string, error) {
	if dir := os.Getenv("ANIMALSDESKTOP_SETTINGS_DIR"); dir != "" {
		return filepath.Join(dir, settingsFileName), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsDirName, settingsFileName), nil
}

func (a *petApp) loadSettings() error {
	path, err := settingsPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var settings appSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}
	if settings.Version < 1 || settings.Version > settingsVersion {
		return nil
	}
	if settings.Version >= 2 {
		a.variant = variantIndexFromSettings(settings.Variant, settings.VariantID, settings.Version, strings.Join(settings.PetNames, "\n"))
		a.coatMode = normalizeCoatMode(settings.CoatMode)
		for i := 0; i < min(len(a.selectedCoats), max(len(settings.SelectedCoats), len(settings.SelectedCoatIDs))); i++ {
			variant := a.selectedCoats[i]
			if i < len(settings.SelectedCoats) {
				variant = settings.SelectedCoats[i]
			}
			a.selectedCoats[i] = variantIndexFromSettings(variant, stringAt(settings.SelectedCoatIDs, i), settings.Version, stringAt(settings.PetNames, i))
		}
		for i, mode := range settings.SelectedSlotModes {
			if i >= len(a.selectedSlotModes) {
				break
			}
			a.selectedSlotModes[i] = normalizePetSlotMode(mode)
		}
		for i, groupID := range settings.SelectedRandomGroupIDs {
			if i >= len(a.selectedRandomGroupIDs) {
				break
			}
			a.selectedRandomGroupIDs[i] = normalizeVariantGroupID(groupID)
		}
		a.randomGroupID = normalizeVariantGroupID(settings.RandomGroupID)
		for i, size := range settings.PetSizes {
			if i >= len(a.petSizes) {
				break
			}
			a.petSizes[i] = normalizePetSizePercent(size)
		}
	}
	for i := range a.petSizes {
		a.petSizes[i] = normalizePetSizePercent(a.petSizes[i])
	}
	for i, name := range settings.PetNames {
		if i >= len(a.petNames) {
			break
		}
		a.petNames[i] = sanitizePetName(name)
	}
	a.speed = normalizeSpeed(settings.Speed)
	a.mode = normalizeBehaviorMode(settings.Mode)
	a.petCount = clamp(settings.PetCount, 1, maxPetCount)
	a.wheelEnabled = settings.WheelEnabled
	a.bidirectional = settings.Bidirectional
	if settings.PositionMode != nil {
		a.positionMode = normalizeOverlayPositionMode(*settings.PositionMode)
	}
	if settings.VerticalOffset != nil {
		a.overlayOffsetY = normalizeOverlayOffset(*settings.VerticalOffset)
	}
	if settings.DisplayIndex != nil {
		a.displayIndex = normalizeDisplayIndex(*settings.DisplayIndex, len(monitorAreas()))
	}
	if settings.DisplayScope != nil {
		a.displayScope = normalizeDisplayScope(*settings.DisplayScope)
	}
	if settings.DisplaySpanEnd != nil {
		a.displaySpanEnd = *settings.DisplaySpanEnd
	}
	a.normalizeDisplaySelection(len(displayAreaForScope(a.displayScope)))
	if settings.WalkRangeStart != nil || settings.WalkRangeEnd != nil {
		start := a.walkRangeStart
		end := a.walkRangeEnd
		if settings.WalkRangeStart != nil {
			start = *settings.WalkRangeStart
		}
		if settings.WalkRangeEnd != nil {
			end = *settings.WalkRangeEnd
		}
		a.walkRangeStart, a.walkRangeEnd = normalizeWalkRange(start, end)
	}
	a.nameLabels = settings.NameLabels
	a.lang = normalizeLanguage(settings.Language)
	if settings.SettingsX != 0 || settings.SettingsY != 0 {
		a.settingsX = settings.SettingsX
		a.settingsY = settings.SettingsY
		a.clampSettingsWindowPosition()
	}
	return nil
}

func (a *petApp) saveSettings() error {
	if a.settingsHwnd != 0 {
		a.rememberSettingsWindowPosition()
	}
	path, err := settingsPath()
	if err != nil {
		return err
	}
	coats := make([]int, len(a.selectedCoats))
	copy(coats, a.selectedCoats[:])
	for i := range coats {
		coats[i] = clamp(coats[i], 0, len(variants)-1)
	}
	coatIDs := variantIDsForIndices(coats)
	slotModes := make([]int, len(a.selectedSlotModes))
	for i, mode := range a.selectedSlotModes {
		slotModes[i] = int(normalizePetSlotMode(int(mode)))
	}
	slotRandomGroupIDs := make([]string, len(a.selectedRandomGroupIDs))
	for i, groupID := range a.selectedRandomGroupIDs {
		slotRandomGroupIDs[i] = normalizeVariantGroupID(groupID)
	}
	sizes := make([]int, len(a.petSizes))
	for i := range a.petSizes {
		sizes[i] = a.petSizePercent(i)
	}
	names := make([]string, len(a.petNames))
	for i := range a.petNames {
		names[i] = sanitizePetName(a.petNames[i])
	}
	positionMode := int(normalizeOverlayPositionMode(int(a.positionMode)))
	verticalOffset := normalizeOverlayOffset(a.overlayOffsetY)
	displayScope, displayIndex, displaySpanEnd := a.normalizedDisplaySelection(len(displayAreaForScope(a.displayScope)))
	displayScopeValue := int(displayScope)
	walkStart, walkEnd := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	settings := appSettings{
		Version:                settingsVersion,
		Variant:                clamp(a.variant, 0, len(variants)-1),
		VariantID:              variantIDAt(a.variant),
		CoatMode:               int(a.coatMode),
		SelectedCoats:          coats,
		SelectedCoatIDs:        coatIDs,
		SelectedSlotModes:      slotModes,
		SelectedRandomGroupIDs: slotRandomGroupIDs,
		RandomGroupID:          normalizeVariantGroupID(a.randomGroupID),
		PetSizes:               sizes,
		Speed:                  normalizeSpeed(a.speed),
		Mode:                   int(normalizeBehaviorMode(int(a.mode))),
		PetCount:               clamp(a.petCount, 1, maxPetCount),
		WheelEnabled:           a.wheelEnabled,
		Bidirectional:          a.bidirectional,
		PositionMode:           &positionMode,
		VerticalOffset:         &verticalOffset,
		DisplayIndex:           &displayIndex,
		DisplayScope:           &displayScopeValue,
		DisplaySpanEnd:         &displaySpanEnd,
		WalkRangeStart:         &walkStart,
		WalkRangeEnd:           &walkEnd,
		Language:               int(normalizeLanguage(int(a.lang))),
		SettingsX:              a.settingsX,
		SettingsY:              a.settingsY,
		NameLabels:             a.nameLabels,
		PetNames:               names,
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeFileAtomically(path, data)
}

func writeFileAtomically(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, settingsFileName+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func variantIndexFromSettings(index int, id string, version int, nameHint string) int {
	if len(variants) == 0 {
		return 0
	}
	if i, ok := variantIndexByID(id); ok {
		return i
	}
	if version < settingsVersion {
		if legacyID, ok := legacyVariantIDForIndex(index, nameHint); ok {
			if i, found := variantIndexByID(legacyID); found {
				return i
			}
		}
	}
	return clamp(index, 0, len(variants)-1)
}

func legacyVariantIDForIndex(index int, nameHint string) (string, bool) {
	if strings.Contains(nameHint, "ハシビロコウ") || strings.Contains(strings.ToLower(nameHint), "shoebill") {
		return "shoebill_stork", true
	}
	return "", false
}

func variantIndexByID(id string) (int, bool) {
	if id == "" {
		return 0, false
	}
	for i, variant := range variants {
		if variant.ID == id {
			return i, true
		}
	}
	return 0, false
}

func variantIDAt(index int) string {
	if len(variants) == 0 {
		return ""
	}
	return variants[clamp(index, 0, len(variants)-1)].ID
}

func variantIDsForIndices(indices []int) []string {
	out := make([]string, len(indices))
	for i, index := range indices {
		out[i] = variantIDAt(index)
	}
	return out
}

func stringAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}

func (a *petApp) persistSettings() {
	a.settingsSaveFailed = a.saveSettings() != nil
}

func normalizeCoatMode(mode int) coatMode {
	switch coatMode(mode) {
	case coatFixed, coatSelected, coatRandom:
		return coatMode(mode)
	default:
		return coatRandom
	}
}

func normalizePetSlotMode(mode int) petSlotMode {
	switch petSlotMode(mode) {
	case petSlotFixed, petSlotRandom:
		return petSlotMode(mode)
	default:
		return petSlotFixed
	}
}

func normalizeVariantGroupID(groupID string) string {
	groupID = strings.TrimSpace(groupID)
	if groupID == "" {
		return ""
	}
	if _, ok := catalog.VariantGroupByID(groupID); !ok {
		return ""
	}
	if len(variantIndicesForGroupID(groupID)) == 0 {
		return ""
	}
	return groupID
}

func normalizeBehaviorMode(mode int) behaviorMode {
	switch behaviorMode(mode) {
	case modeKeyboard, modeRandom:
		return behaviorMode(mode)
	default:
		return modeRandom
	}
}

func normalizeLanguage(lang int) language {
	switch language(lang) {
	case langJapanese, langEnglish:
		return language(lang)
	default:
		return langJapanese
	}
}

func normalizeSpeed(speed int) int {
	switch speed {
	case 2, 3, 5:
		return speed
	default:
		return 3
	}
}

func normalizeOverlayPositionMode(mode int) overlayPositionMode {
	switch overlayPositionMode(mode) {
	case positionTaskbarEdge, positionScreenBottom:
		return overlayPositionMode(mode)
	default:
		return positionTaskbarEdge
	}
}

func normalizeOverlayOffset(offset int) int {
	return clamp(offset, minOverlayOffsetY, maxOverlayOffsetY)
}

func normalizeDisplayIndex(index int, count int) int {
	if count <= 0 {
		return 0
	}
	return clamp(index, 0, count-1)
}

func normalizeDisplayScope(mode int) displayScope {
	switch displayScope(mode) {
	case displayScopeSingle, displayScopeSpan:
		return displayScope(mode)
	default:
		return displayScopeSingle
	}
}

func normalizeDisplaySpan(start, end, count int) (int, int) {
	if count <= 0 {
		return 0, 0
	}
	start = normalizeDisplayIndex(start, count)
	end = normalizeDisplayIndex(end, count)
	if end < start {
		start, end = end, start
	}
	return start, end
}

func normalizeWalkRange(start, end int) (int, int) {
	start = clamp(start, 0, 100)
	end = clamp(end, 0, 100)
	if end < start {
		start, end = end, start
	}
	if end-start >= minWalkRangeSpan {
		return start, end
	}
	mid := (start + end) / 2
	start = mid - minWalkRangeSpan/2
	end = start + minWalkRangeSpan
	if start < 0 {
		start = 0
		end = minWalkRangeSpan
	}
	if end > 100 {
		end = 100
		start = 100 - minWalkRangeSpan
	}
	return start, end
}

func (a *petApp) tick() bool {
	if a.closing.Load() {
		return false
	}
	dt := a.motionClock.Step(time.Now())
	a.legacyElapsed += dt
	legacyTick := a.legacyElapsed >= float64(timerInterval)/1000
	if legacyTick {
		a.legacyElapsed = math.Mod(a.legacyElapsed, float64(timerInterval)/1000)
	}
	overlay := a.overlayRect()
	a.syncScene(overlay)
	a.updateHoverName()
	if foragePropsEnabled {
		a.ensureForageItems()
	} else {
		a.clearForageItems()
	}
	segments := a.overlaySegments()
	timedActive := false
	for i := range a.pets {
		p := &a.pets[i]
		if v := timedVariant(variants[a.petVariant(p)].ID); v != nil {
			timedActive = true
			dpi := 96
			for _, segment := range segments {
				if p.x >= segment.SceneLeft && p.x < segment.SceneRight {
					dpi = segment.DPI
					break
				}
			}
			a.tickTimedPet(i, p, v, dt, dpi)
		} else if legacyTick {
			p.timed = motion.Player{}
			a.tickPet(i, p)
		}
	}
	if !legacyTick {
		return timedActive
	}
	a.tickReactions()
	a.syncNearbyWalkers()
	a.maybeStartSocial()
	a.tickCount++
	return true
}

func (a *petApp) tickReactions() {
	out := a.reactions[:0]
	for _, reaction := range a.reactions {
		reaction.ticks--
		if reaction.ticks > 0 && reaction.pet >= 0 && reaction.pet < len(a.pets) {
			out = append(out, reaction)
		}
	}
	a.reactions = out
}

func (a *petApp) tickPet(index int, p *desktopPet) {
	if p.stateTicks <= 0 {
		switch p.state {
		case stateWheel:
			a.leaveWheel(index, p)
		case stateTurn:
			a.finishTurn(p)
		case stateEat:
			if p.item >= 0 && p.item < len(a.forage) {
				a.finishEating(index, p)
			} else {
				a.chooseRandomAction(p)
			}
		case stateDig, stateStand, stateFaceGroom, stateRest:
			if a.mode == modeRandom {
				a.chooseRandomAction(p)
			} else {
				p.state = stateIdle
				p.moveSpeed = 0
				p.stateTicks = 12
			}
		case stateNibble:
			if p.item >= 0 && p.item < len(a.forage) {
				a.finishGnawing(index, p)
			} else if a.mode == modeRandom {
				a.chooseRandomAction(p)
			} else {
				p.state = stateIdle
				p.moveSpeed = 0
				p.stateTicks = 12
			}
		case stateGroom, stateCarry:
			a.releaseForage(index, p)
			if a.mode == modeRandom {
				a.chooseRandomAction(p)
			} else {
				p.state = stateIdle
				p.moveSpeed = 0
				p.stateTicks = 12
			}
		default:
			if a.mode != modeRandom {
				p.state = stateIdle
				p.moveSpeed = 0
				p.stateTicks = 12
				break
			}
			a.chooseRandomAction(p)
		}
	}

	speed := 0
	switch p.state {
	case stateWalk, stateScurry, stateHop, stateForage, stateCarry:
		speed = p.moveSpeed
	case stateWheel:
		w, _ := a.petSpriteSize(index)
		p.x = clamp(a.wheelX-wheelSize/2, 0, max(0, a.sceneW-w))
	}

	if speed > 0 {
		p.x += speed * p.dir
		if p.state == stateForage {
			a.maybeStartGnawing(index, p)
		}
	}

	p.stateTicks--
	w, _ := a.petSpriteSize(index)
	if p.x > a.sceneW+8 {
		a.resetPetAtEdge(index, p, 1)
	} else if p.x < -w-8 {
		a.resetPetAtEdge(index, p, -1)
	}
	p.frame++
}

func (a *petApp) chooseRandomAction(p *desktopPet) {
	roll := rand.Intn(100)
	p.frame = 0
	p.motionSet = rand.Intn(motionSets)
	p.item = noItem
	p.carryKind = noItem
	if p.dir == 0 {
		p.dir = 1
	}
	if a.bidirectional && p.state != stateTurn && rand.Intn(100) < 16 {
		a.startTurn(p, -p.dir, stateWalk)
		return
	}
	if foragePropsEnabled && roll < 18 && a.maybeAssignForageTarget(p) {
		return
	}
	switch randomActionState(a.petMotionProfile(p), roll) {
	case stateIdle:
		p.state = stateIdle
		p.moveSpeed = 0
		p.stateTicks = 24 + rand.Intn(58)
		return
	case stateWalk:
		p.state = stateWalk
		p.moveSpeed = max(1, a.speed-1+rand.Intn(2))
		p.stateTicks = 34 + rand.Intn(92)
	case stateScurry:
		p.state = stateScurry
		p.moveSpeed = a.speed + 1 + rand.Intn(2)
		p.stateTicks = 10 + rand.Intn(18)
	case stateNibble:
		p.state = stateNibble
		p.moveSpeed = 0
		p.stateTicks = 26 + rand.Intn(32)
	case stateRest:
		p.state = stateRest
		p.moveSpeed = 0
		p.stateTicks = 48 + rand.Intn(72)
	case stateStand:
		p.state = stateStand
		p.moveSpeed = 0
		p.stateTicks = 24 + rand.Intn(28)
	case stateFaceGroom:
		p.state = stateFaceGroom
		p.moveSpeed = 0
		p.stateTicks = 28 + rand.Intn(30)
	default:
		p.state = stateHop
		p.moveSpeed = max(1, a.speed-1)
		p.stateTicks = 14 + rand.Intn(16)
	}
}

func randomActionState(motionProfile string, roll int) behaviorState {
	if motionProfile == catalog.MotionProfileFerretSlink {
		switch {
		case roll < 30:
			return stateIdle
		case roll < 70:
			return stateWalk
		case roll < 84:
			return stateScurry
		case roll < 90:
			return stateNibble
		case roll < 92:
			return stateRest
		case roll < 95:
			return stateStand
		case roll < 98:
			return stateFaceGroom
		default:
			return stateHop
		}
	}
	switch {
	case roll < 30:
		return stateIdle
	case roll < 70:
		return stateWalk
	case roll < 84:
		return stateScurry
	case roll < 90:
		return stateNibble
	case roll < 94:
		return stateStand
	case roll < 98:
		return stateFaceGroom
	default:
		return stateHop
	}
}

func (a *petApp) render() {
	active := map[string]bool{}
	for i := range a.pets {
		active[variants[a.petVariant(&a.pets[i])].ID] = true
	}
	timedStore.Retain(active)
	if a.frames != nil {
		a.frames.retain(active)
	}
	if a.overlayHidden {
		a.hideOverlayWindows(0)
		return
	}
	areas, display := a.selectedDisplayAreasAndCombined()
	overlay := a.overlayRectFor(display.Work, display.Screen)
	a.syncScene(overlay)
	segments := a.overlaySegmentsForAreas(areas, overlay)
	if len(segments) == 0 {
		segments = []overlaySegment{{
			Rect:       overlay,
			Screen:     display.Screen,
			SceneLeft:  0,
			SceneRight: a.sceneW,
			DPI:        display.DPI,
		}}
	}
	for i, segment := range segments {
		hwnd := a.overlayWindow(i)
		if hwnd == 0 {
			continue
		}
		canvas := a.renderOverlaySegment(segment)
		updateLayeredWindow(hwnd, canvas, int(segment.Rect.Left), int(segment.Rect.Top))
		win.ShowWindow(hwnd, win.SW_SHOWNOACTIVATE)
	}
	a.hideOverlayWindows(len(segments))
}

func (a *petApp) renderOverlaySegment(segment overlaySegment) *image.RGBA {
	width := max(1, int(segment.Rect.Right-segment.Rect.Left))
	height := max(1, int(segment.Rect.Bottom-segment.Rect.Top))
	logicalWidth := max(1, unscaleForDPI(width, segment.DPI))
	canvas := image.NewRGBA(image.Rect(0, 0, logicalWidth, sceneH))
	draw.Draw(canvas, canvas.Bounds(), image.Transparent, image.Point{}, draw.Src)

	wheelActive := a.wheelEnabled && a.hasWheelRunner()
	wheelX := unscaleForDPI(a.wheelX-wheelSize/2-segment.SceneLeft, segment.DPI)
	wheelY := sceneH - wheelSize - 2
	if wheelActive {
		drawWheelBack(canvas, wheelX, wheelY, a.wheel)
	}

	if foragePropsEnabled {
		a.drawForageItems(canvas, segment.SceneLeft, segment.DPI)
	}

	for i := range a.pets {
		p := &a.pets[i]
		if p.state == stateWheel {
			continue
		}
		variant := variants[a.petVariant(p)]
		if v := timedVariant(variant.ID); v != nil {
			if p.timed.ID != v.ID {
				p.timed.Reset(v.ID, float64(p.x), p.dir < 0)
			}
			w, h := a.petSpriteSize(i)
			v.Draw(canvas, p.timed.Action, p.timed.Phase(), unscaleForDPI(p.x-segment.SceneLeft, segment.DPI), sceneH-p.laneOffset-3, w, h, p.dir)
			continue
		}
		frame := currentFrameForVariant(p.state, p.frame, variant)
		src := a.frames.frame(variant, p.motionSet, frame)
		w, h := a.petSpriteSize(i)
		y := sceneH - h - p.laneOffset
		localX := unscaleForDPI(p.x-segment.SceneLeft, segment.DPI)
		drawPetSprite(canvas, src, p, variant, frame, localX, y, w, h)
		if foragePropsEnabled && p.state == stateCarry && p.carryKind != noItem {
			propX := localX + w - 18
			if p.dir < 0 {
				propX = localX + 18
			}
			a.drawForageProp(canvas, propX, y+35, p.carryKind)
		} else if foragePropsEnabled && (p.state == stateEat || p.state == stateDig) && p.carryKind != noItem {
			propX := localX + w - 20
			if p.dir < 0 {
				propX = localX + 20
			}
			a.drawForageProp(canvas, propX, y+44, p.carryKind)
		}
	}

	if wheelActive {
		for i := range a.pets {
			p := &a.pets[i]
			if p.state != stateWheel {
				continue
			}
			variant := variants[a.petVariant(p)]
			frame := currentFrameForVariant(p.state, p.frame, variant)
			src := a.frames.frame(variant, p.motionSet, frame)
			drawWheelRunner(canvas, wheelX, wheelY, src, p.frame, drawDirectionForVariant(1, variant))
		}
		drawWheelFront(canvas, wheelX, wheelY, a.tickCount)
	}
	a.drawReactions(canvas, segment.SceneLeft, segment.SceneRight, segment.DPI)
	if width == logicalWidth && height == sceneH {
		return canvas
	}
	return scaleImageTo(canvas, width, height)
}

func currentFrame(state behaviorState, frame int) int {
	return currentFrameForVariant(state, frame, coatVariant{})
}

func currentFrameForVariant(state behaviorState, frame int, variant coatVariant) int {
	if catalog.MotionProfileForVariant(variant) == catalog.MotionProfileFerretSlink {
		return currentFerretFrame(state, frame)
	}
	switch state {
	case stateIdle:
		return frameFromSeq(idleFrameSeq, frame, 5)
	case stateWalk, stateForage, stateCarry:
		return frameFromSeq(walkFrameSeq, frame, 2)
	case stateScurry, stateWheel:
		return frameFromSeq(walkFrameSeq, frame, 1)
	case stateNibble:
		if usesStableActionFallback(variant) {
			return frameFromSeq(hopFrameSeq, frame, 3)
		}
		return frameFromSeq(nibbleFrameSeq, frame, 3)
	case stateHop:
		return frameFromSeq(hopFrameSeq, frame, 2)
	case stateGroom:
		if usesStableActionFallback(variant) {
			return frameFromSeq(groomFrameSeq, frame, 3)
		}
		return frameFromSeq(nibbleFrameSeq, frame, 4)
	case stateTurn:
		return frameFromSeqClamped(turnFrameSeq, frame, 2)
	case stateEat:
		return frameFromSeq(eatFrameSeq, frame, 3)
	case stateDig:
		return frameFromSeq(digFrameSeq, frame, 2)
	case stateStand:
		return frameFromSeq(standFrameSeq, frame, 4)
	case stateFaceGroom:
		return frameFromSeq(groomFrameSeq, frame, 3)
	}
	return idleStart
}

func currentFerretFrame(state behaviorState, frame int) int {
	switch state {
	case stateIdle:
		return frameFromRange(idleStart, idleFrames, frame, 5)
	case stateWalk, stateForage, stateCarry:
		return frameFromRange(walkStart, walkFrames, frame, 2)
	case stateScurry, stateWheel:
		return frameFromRange(scurryStart, scurryFrames, frame, 1)
	case stateNibble, stateEat, stateDig:
		return frameFromRange(nibbleStart, nibbleFrames, frame, 3)
	case stateGroom, stateFaceGroom:
		return frameFromRange(hopStart, hopFrames, frame, 3)
	case stateTurn:
		return frameFromRangeClamped(turnStart, turnFrames, frame, 2)
	case stateHop:
		return frameFromRange(eatStart, eatFrames+digFrames, frame, 2)
	case stateRest:
		return frameFromRange(standStart, standFrames+groomFrames, frame, 4)
	case stateStand:
		return frameFromRange(wheelRunStart, wheelRunFrames, frame, 4)
	default:
		return idleStart
	}
}

func usesStableActionFallback(variant coatVariant) bool {
	switch variant.ID {
	case "sugar_glider_gray", "rabbit_chestnut_agouti":
		return true
	default:
		return false
	}
}

func frameFromSeq(seq []int, frame, divisor int) int {
	if len(seq) == 0 {
		return idleStart
	}
	if divisor < 1 {
		divisor = 1
	}
	return seq[(frame/divisor)%len(seq)]
}

func frameFromSeqClamped(seq []int, frame, divisor int) int {
	if len(seq) == 0 {
		return idleStart
	}
	if divisor < 1 {
		divisor = 1
	}
	index := frame / divisor
	if index >= len(seq) {
		index = len(seq) - 1
	}
	return seq[index]
}

func frameFromRange(start, count, frame, divisor int) int {
	if count < 1 {
		return idleStart
	}
	if divisor < 1 {
		divisor = 1
	}
	return start + (frame/divisor)%count
}

func frameFromRangeClamped(start, count, frame, divisor int) int {
	if count < 1 {
		return idleStart
	}
	if divisor < 1 {
		divisor = 1
	}
	index := frame / divisor
	if index >= count {
		index = count - 1
	}
	return start + index
}

func (a *petApp) onTyping() {
	if a.mode != modeKeyboard {
		return
	}
	wheelStarted := false
	for i := range a.pets {
		p := &a.pets[i]
		if v := timedVariant(variants[a.petVariant(p)].ID); v != nil {
			if p.timed.ID != v.ID {
				p.timed.Reset(v.ID, float64(p.x), p.dir < 0)
			}
			p.timed.ReactWith(v)
			continue
		}
		if a.wheelEnabled && !wheelStarted && p.item == noItem && a.petWheelCapable(p) {
			a.enterWheelFromTyping(i, p)
			wheelStarted = true
			continue
		}
		p.state = stateScurry
		p.frame = rand.Intn(scurryFrames)
		p.motionSet = rand.Intn(motionSets)
		p.stateTicks = 18 + rand.Intn(16)
		p.moveSpeed = a.speed + 2 + rand.Intn(2)
	}
}

func (a *petApp) onMouseClick(screenX, screenY int) {
	if a.overlayHidden {
		return
	}
	index := a.petAtScreenPoint(screenX, screenY)
	if index < 0 {
		return
	}
	a.showPetReaction(index)
	a.render()
}

func (a *petApp) showPetReaction(index int) {
	if index < 0 || index >= len(a.pets) {
		return
	}
	p := &a.pets[index]
	if v := timedVariant(variants[a.petVariant(p)].ID); v != nil {
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
			return
		}
	}
	a.reactions = append(a.reactions, petReaction{pet: index, kind: kind, ticks: reactionTicks})
}

func (a *petApp) petAtScreenPoint(screenX, screenY int) int {
	for _, segment := range a.overlaySegments() {
		if screenX < int(segment.Rect.Left) || screenX >= int(segment.Rect.Right) ||
			screenY < int(segment.Rect.Top) || screenY >= int(segment.Rect.Bottom) {
			continue
		}
		localX := screenX - int(segment.Rect.Left)
		localY := screenY - int(segment.Rect.Top)
		return a.petAtSegmentPoint(segment, localX, localY)
	}
	return -1
}

func (a *petApp) petAtSegmentPoint(segment overlaySegment, localX, localY int) int {
	if localX < 0 || localY < 0 || localX >= int(segment.Rect.Right-segment.Rect.Left) || localY >= int(segment.Rect.Bottom-segment.Rect.Top) {
		return -1
	}
	for i := len(a.pets) - 1; i >= 0; i-- {
		w, h := a.petSpriteSize(i)
		scaledW := scaleForDPI(w, segment.DPI)
		scaledH := scaleForDPI(h, segment.DPI)
		scaledLaneOffset := scaleForDPI(a.pets[i].laneOffset, segment.DPI)
		petX := a.pets[i].x - segment.SceneLeft
		petY := scaleForDPI(sceneH, segment.DPI) - scaledH - scaledLaneOffset
		if scenePointInPet(a.pets[i], localX, localY, scaledW, scaledH, petX, petY) {
			return i
		}
	}
	return -1
}

func (a *petApp) petAtScenePoint(sceneX, sceneY int) int {
	if sceneX < 0 || sceneX >= a.sceneW || sceneY < 0 || sceneY >= sceneH {
		return -1
	}
	for i := len(a.pets) - 1; i >= 0; i-- {
		w, h := a.petSpriteSize(i)
		y := sceneH - h - a.pets[i].laneOffset
		if scenePointInPet(a.pets[i], sceneX, sceneY, w, h, a.pets[i].x, y) {
			return i
		}
	}
	return -1
}

func scenePointInPet(p desktopPet, sceneX, sceneY int, w int, h int, petX int, petY int) bool {
	if p.state == stateWheel {
		return false
	}
	insetX := min(6, max(0, w/8))
	insetTop := min(8, max(0, h/8))
	insetBottom := min(4, max(0, h/12))
	return sceneX >= petX+insetX && sceneX <= petX+w-insetX && sceneY >= petY+insetTop && sceneY <= petY+h-insetBottom
}

func (a *petApp) updateHoverName() {
	if a.overlayHidden {
		a.hoverPet = -1
		a.hideNameWindow()
		return
	}
	if !a.nameLabels {
		a.hoverPet = -1
		a.hideNameWindow()
		return
	}
	var pt win.POINT
	if !win.GetCursorPos(&pt) {
		a.hideNameWindow()
		return
	}
	index := a.petAtScreenPoint(int(pt.X), int(pt.Y))
	if index < 0 || index >= len(a.pets) {
		a.hoverPet = -1
		a.hideNameWindow()
		return
	}
	a.hoverPet = index
	a.showNameWindow(index)
}

func (a *petApp) showNameWindow(index int) {
	if index < 0 || index >= len(a.pets) {
		a.hideNameWindow()
		return
	}
	name := a.petDisplayName(index)
	if name == "" {
		a.hideNameWindow()
		return
	}
	display := a.selectedDisplayArea()
	overlay := a.overlayRectFor(display.Work, display.Screen)
	screen := display.Screen
	p := a.pets[index]
	runes := []rune(name)
	w := clamp(34+len(runes)*12, 72, 220)
	h := 30
	petW, petH := a.petSpriteSize(index)
	segment, ok := a.overlaySegmentForSceneX(p.x + petW/2)
	if ok {
		overlay = segment.Rect
		screen = segment.Screen
	}
	scaledPetW := scaleForDPI(petW, segment.DPI)
	scaledPetH := scaleForDPI(petH, segment.DPI)
	scaledBaseY := scaleForDPI(sceneH, segment.DPI) - scaledPetH - scaleForDPI(p.laneOffset, segment.DPI)
	x := int(overlay.Left) + p.x - segment.SceneLeft + scaledPetW/2 - w/2
	y := int(overlay.Top) + scaledBaseY - h - scaleForDPI(8, segment.DPI)
	x = clamp(x, int(overlay.Left), int(overlay.Right)-w)
	y = clamp(y, int(screen.Top), int(screen.Bottom)-h)
	if a.nameHwnd == 0 {
		a.nameHwnd = win.CreateWindowEx(
			win.WS_EX_TOOLWINDOW|win.WS_EX_TOPMOST|win.WS_EX_NOACTIVATE|win.WS_EX_TRANSPARENT,
			syscall.StringToUTF16Ptr(windowClass),
			syscall.StringToUTF16Ptr(name),
			win.WS_POPUP,
			int32(x), int32(y), int32(w), int32(h),
			a.hwnd, 0, a.hinst, nil,
		)
	}
	if a.nameHwnd == 0 {
		return
	}
	if a.nameText != name {
		a.nameText = name
		setWindowText(a.nameHwnd, name)
		win.InvalidateRect(a.nameHwnd, nil, true)
	}
	win.SetWindowPos(a.nameHwnd, win.HWND_TOPMOST, int32(x), int32(y), int32(w), int32(h), win.SWP_NOACTIVATE|win.SWP_SHOWWINDOW)
}

func (a *petApp) hideNameWindow() {
	if a.nameHwnd != 0 {
		win.ShowWindow(a.nameHwnd, win.SW_HIDE)
	}
}

func (a *petApp) temporaryVisibilityLabel() string {
	if a.overlayHidden {
		return a.localText("表示する", "Show")
	}
	return a.localText("一時的に非表示", "Hide temporarily")
}

func (a *petApp) setOverlayHidden(hidden bool) {
	a.overlayHidden = hidden
	if hidden {
		a.hideNameWindow()
	}
	a.applyOverlayVisibility()
}

func (a *petApp) applyOverlayVisibility() {
	if a.hwnd == 0 {
		return
	}
	if a.overlayHidden {
		a.hideOverlayWindows(0)
		return
	}
	win.ShowWindow(a.hwnd, win.SW_SHOWNOACTIVATE)
	a.render()
}

func (a *petApp) overlayWindow(index int) win.HWND {
	if index <= 0 {
		return a.hwnd
	}
	for len(a.overlayHwnds) < index {
		hwnd := win.CreateWindowEx(
			win.WS_EX_LAYERED|win.WS_EX_TOPMOST|win.WS_EX_TOOLWINDOW|win.WS_EX_TRANSPARENT,
			syscall.StringToUTF16Ptr(windowClass),
			syscall.StringToUTF16Ptr(appName),
			win.WS_POPUP,
			0, 0, 1, 1,
			0, 0, a.hinst, nil,
		)
		if hwnd == 0 {
			return 0
		}
		a.overlayHwnds = append(a.overlayHwnds, hwnd)
	}
	return a.overlayHwnds[index-1]
}

func (a *petApp) hideOverlayWindows(used int) {
	if used <= 0 && a.hwnd != 0 {
		win.ShowWindow(a.hwnd, win.SW_HIDE)
	}
	for i, hwnd := range a.overlayHwnds {
		if hwnd != 0 && i+1 >= used {
			win.ShowWindow(hwnd, win.SW_HIDE)
		}
	}
}

func (a *petApp) nameWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_PAINT:
		a.paintNameWindow(hwnd)
		return 0
	case win.WM_ERASEBKGND:
		return 1
	case win.WM_NCHITTEST:
		return ^uintptr(0)
	case win.WM_DESTROY:
		if hwnd == a.nameHwnd {
			a.nameHwnd = 0
			a.nameText = ""
		}
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (a *petApp) paintNameWindow(hwnd win.HWND) {
	a.ensureSettingsFonts()
	var ps win.PAINTSTRUCT
	hdc := win.BeginPaint(hwnd, &ps)
	if hdc == 0 {
		return
	}
	defer win.EndPaint(hwnd, &ps)
	var rect win.RECT
	win.GetClientRect(hwnd, &rect)
	drawRectFill(hdc, rect, rgb(255, 255, 248))
	body := win.RECT{Left: 1, Top: 1, Right: rect.Right - 1, Bottom: rect.Bottom - 1}
	drawRoundFill(hdc, body, rgb(255, 255, 248), 14)
	drawRoundOutline(hdc, body, rgb(64, 91, 72), 14)
	textRect := win.RECT{Left: 10, Top: 1, Right: rect.Right - 10, Bottom: rect.Bottom - 1}
	drawTextLine(hdc, a.nameText, textRect, a.settingsSmallFont, rgb(27, 36, 32), win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)
}

func (a *petApp) syncScene(bounds win.RECT) {
	a.sceneW = max(1, int(bounds.Right-bounds.Left))
	nextWheelX := a.sceneW * 2 / 3
	a.wheelX = clamp(nextWheelX, wheelSize/2+24, max(wheelSize/2+24, a.sceneW-wheelSize/2-24))
}

func (a *petApp) overlayRect() win.RECT {
	display := a.selectedDisplayArea()
	return a.overlayRectFor(display.Work, display.Screen)
}

func (a *petApp) normalizedDisplaySelection(count int) (displayScope, int, int) {
	if count <= 0 {
		return displayScopeSingle, 0, 0
	}
	scope := normalizeDisplayScope(int(a.displayScope))
	start := normalizeDisplayIndex(a.displayIndex, count)
	end := normalizeDisplayIndex(a.displaySpanEnd, count)
	if scope == displayScopeSpan && count > 1 {
		start, end = normalizeDisplaySpan(start, end, count)
		if start == end {
			if end < count-1 {
				end++
			} else {
				start--
			}
		}
		return displayScopeSpan, start, end
	}
	return displayScopeSingle, start, start
}

func (a *petApp) normalizeDisplaySelection(count int) {
	scope, start, end := a.normalizedDisplaySelection(count)
	a.displayScope = scope
	a.displayIndex = start
	a.displaySpanEnd = end
}

func (a *petApp) overlayRectFor(work, screen win.RECT) win.RECT {
	base := a.applyWalkRange(a.overlayBaseFor(work, screen))
	y := a.overlayYForBase(base, screen)
	return win.RECT{
		Left:   base.Left,
		Top:    int32(y),
		Right:  base.Right,
		Bottom: int32(y + sceneH),
	}
}

func (a *petApp) overlayBaseFor(work, screen win.RECT) win.RECT {
	mode := normalizeOverlayPositionMode(int(a.positionMode))
	base := work
	if mode == positionScreenBottom {
		base = screen
	}
	if base.Right <= base.Left {
		base.Left = screen.Left
		base.Right = screen.Right
	}
	if base.Bottom <= base.Top {
		base.Top = screen.Top
		base.Bottom = screen.Bottom
	}
	return base
}

func (a *petApp) overlayYForBase(base, screen win.RECT) int {
	return a.overlayYForBaseDPI(base, screen, defaultDPI)
}

func (a *petApp) overlayYForBaseDPI(base, screen win.RECT, dpi int) int {
	height := scaleForDPI(sceneH, dpi)
	offset := normalizeOverlayOffset(a.overlayOffsetY)
	y := int(base.Bottom) - height + scaleForDPI(offset, dpi)
	minY := int(screen.Top)
	maxY := max(minY, int(screen.Bottom)-height)
	return clamp(y, minY, maxY)
}

func (a *petApp) applyWalkRange(base win.RECT) win.RECT {
	width := int(base.Right - base.Left)
	if width <= 0 {
		return base
	}
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	left := int(base.Left) + width*start/100
	right := int(base.Left) + width*end/100
	if right-left < spriteW {
		mid := (left + right) / 2
		left = mid - spriteW/2
		right = left + spriteW
	}
	left = clamp(left, int(base.Left), max(int(base.Left), int(base.Right)-spriteW))
	right = clamp(right, left+spriteW, int(base.Right))
	base.Left = int32(left)
	base.Right = int32(right)
	return base
}

func (a *petApp) resetOverlayPlacement() {
	a.positionMode = positionTaskbarEdge
	a.overlayOffsetY = defaultOverlayOffsetY
	a.displayScope = displayScopeSingle
	a.displayIndex = 0
	a.displaySpanEnd = 0
	a.walkRangeStart = defaultWalkRangeStart
	a.walkRangeEnd = defaultWalkRangeEnd
}

func (a *petApp) adjustOverlayOffset(delta int) {
	a.overlayOffsetY = normalizeOverlayOffset(a.overlayOffsetY + delta)
}

func (a *petApp) setWalkRange(start, end int) {
	a.walkRangeStart, a.walkRangeEnd = normalizeWalkRange(start, end)
	overlay := a.overlayRect()
	a.syncScene(overlay)
	a.arrangePetsForOverlay(overlay)
}

func (a *petApp) adjustWalkRangeWidth(delta int) {
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	a.setWalkRange(start-delta, end+delta)
}

func (a *petApp) shiftWalkRange(delta int) {
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	span := end - start
	start = clamp(start+delta, 0, 100-span)
	a.setWalkRange(start, start+span)
}

func (a *petApp) clampPetsToScene() {
	for i := range a.pets {
		w, _ := a.petSpriteSize(i)
		a.pets[i].x = clamp(a.pets[i].x, 0, max(0, a.sceneW-w))
	}
	for i := range a.forage {
		if a.forage[i].active {
			a.forage[i].x = clamp(a.forage[i].x, 24, max(24, a.sceneW-24))
		}
	}
}

func (a *petApp) adjustDisplayIndex(delta int) {
	areas := displayAreaForScope(a.displayScope)
	count := len(areas)
	if count <= 0 {
		a.displayIndex = 0
		a.displaySpanEnd = 0
		return
	}
	scope, start, end := a.normalizedDisplaySelection(count)
	if scope == displayScopeSpan {
		width := end - start
		nextStart := clamp(start+delta, 0, max(0, count-1-width))
		a.displayIndex = nextStart
		a.displaySpanEnd = nextStart + width
	} else {
		a.displayIndex = (start + delta + count) % count
		a.displaySpanEnd = a.displayIndex
	}
	a.resetPosition()
}

func (a *petApp) setDisplayScope(scope displayScope) {
	targetScope := normalizeDisplayScope(int(scope))
	wasSpan := normalizeDisplayScope(int(a.displayScope)) == displayScopeSpan
	if targetScope == displayScopeSpan {
		singleAreas := monitorAreas()
		spanAreas := monitorAreasByPosition()
		count := len(spanAreas)
		if count > 1 {
			start := normalizeDisplayIndex(a.displayIndex, count)
			if normalizeDisplayScope(int(a.displayScope)) != displayScopeSpan && len(singleAreas) > 0 {
				start = findDisplayAreaIndex(spanAreas, singleAreas[normalizeDisplayIndex(a.displayIndex, len(singleAreas))])
			}
			end := a.displaySpanEnd
			if normalizeDisplayScope(int(a.displayScope)) != displayScopeSpan || end == start {
				end = min(start+1, count-1)
				if end == start {
					start = max(0, start-1)
				}
			}
			a.displayScope = displayScopeSpan
			a.displayIndex, a.displaySpanEnd = normalizeDisplaySpan(start, end, count)
			if !wasSpan {
				a.walkRangeStart = defaultWalkRangeStart
				a.walkRangeEnd = defaultWalkRangeEnd
			}
			a.resetPosition()
			return
		}
	}
	singleAreas := monitorAreas()
	spanAreas := monitorAreasByPosition()
	count := len(singleAreas)
	index := normalizeDisplayIndex(a.displayIndex, count)
	if normalizeDisplayScope(int(a.displayScope)) == displayScopeSpan && len(spanAreas) > 0 {
		_, start, _ := a.normalizedDisplaySelection(len(spanAreas))
		index = findDisplayAreaIndex(singleAreas, spanAreas[start])
	}
	a.displayScope = displayScopeSingle
	a.displayIndex = index
	a.displaySpanEnd = a.displayIndex
	a.resetPosition()
}

func (a *petApp) adjustDisplaySpan(delta int) {
	areas := monitorAreasByPosition()
	count := len(areas)
	if count <= 1 {
		a.setDisplayScope(displayScopeSingle)
		return
	}
	if normalizeDisplayScope(int(a.displayScope)) != displayScopeSpan {
		a.setDisplayScope(displayScopeSpan)
	}
	_, start, end := a.normalizedDisplaySelection(count)
	if delta > 0 {
		if end < count-1 {
			end++
		} else if start > 0 {
			start--
		}
	} else if delta < 0 && end-start > 1 {
		end--
	}
	a.displayScope = displayScopeSpan
	a.displayIndex, a.displaySpanEnd = normalizeDisplaySpan(start, end, count)
	a.resetPosition()
}

func (a *petApp) setPetCount(count int) {
	count = clamp(count, 1, maxPetCount)
	a.petCount = count
	for len(a.pets) < count {
		a.pets = append(a.pets, a.newPet(len(a.pets)))
	}
	if len(a.pets) > count {
		a.pets = a.pets[:count]
	}
	for i := range a.forage {
		if a.forage[i].owner >= count || a.forage[i].owner == reservedItem {
			a.forage[i].owner = noItem
			a.forage[i].active = false
		}
	}
	for i := range a.pets {
		a.pets[i].laneOffset = (i % 3) * 5
		if a.coatMode != coatRandom {
			a.pets[i].variant = a.variantForIndex(i)
		} else if a.pets[i].variant < 0 || a.pets[i].variant >= len(variants) {
			a.pets[i].variant = a.variantForIndex(i)
		}
		if a.pets[i].dir == 0 {
			a.pets[i].dir = 1
			a.pets[i].nextDir = 1
		}
	}
	a.arrangePetsAcrossScene()
}

type sceneSegment struct {
	Left  int
	Right int
}

func (a *petApp) arrangePetsAcrossScene() {
	a.arrangePetsInSegments(nil)
}

func (a *petApp) arrangePetsForOverlay(overlay win.RECT) {
	a.arrangePetsInSegments(a.sceneSegmentsForOverlay(overlay))
}

func (a *petApp) arrangePetsInSegments(segments []sceneSegment) {
	positions := petScenePositions(a.sceneW, len(a.pets), segments)
	for i, x := range positions {
		p := &a.pets[i]
		w, _ := a.petSpriteSize(i)
		p.x = clamp(x, 0, max(0, a.sceneW-w))
		dir := 1
		if a.bidirectional && i%2 == 1 {
			dir = -1
		}
		p.dir = dir
		p.nextDir = dir
		if p.state == stateTurn {
			p.state = stateWalk
			p.moveSpeed = max(1, a.speed-1)
			p.stateTicks = 24
		}
	}
}

func (a *petApp) sceneSegmentsForOverlay(overlay win.RECT) []sceneSegment {
	areas := displayAreaForScope(a.displayScope)
	if len(areas) == 0 || overlay.Right <= overlay.Left {
		return nil
	}
	scope, start, end := a.normalizedDisplaySelection(len(areas))
	if scope != displayScopeSpan {
		end = start
	}
	mode := normalizeOverlayPositionMode(int(a.positionMode))
	segments := make([]sceneSegment, 0, end-start+1)
	for _, area := range areas[start : end+1] {
		base := area.Work
		if mode == positionScreenBottom {
			base = area.Screen
		}
		left := max(int(base.Left), int(overlay.Left))
		right := min(int(base.Right), int(overlay.Right))
		if right-left >= spriteW {
			segments = append(segments, sceneSegment{
				Left:  left - int(overlay.Left),
				Right: right - int(overlay.Left),
			})
		}
	}
	return mergeSceneSegments(segments)
}

func mergeSceneSegments(segments []sceneSegment) []sceneSegment {
	if len(segments) <= 1 {
		return segments
	}
	sort.SliceStable(segments, func(i, j int) bool {
		if segments[i].Left != segments[j].Left {
			return segments[i].Left < segments[j].Left
		}
		return segments[i].Right < segments[j].Right
	})
	merged := segments[:0]
	for _, segment := range segments {
		if segment.Right <= segment.Left {
			continue
		}
		if len(merged) == 0 || segment.Left >= merged[len(merged)-1].Right {
			merged = append(merged, segment)
			continue
		}
		if segment.Right > merged[len(merged)-1].Right {
			merged[len(merged)-1].Right = segment.Right
		}
	}
	return merged
}

func petScenePositions(sceneW, count int, segments []sceneSegment) []int {
	if count <= 0 || sceneW <= 0 {
		return nil
	}
	segments = normalizeSceneSegments(sceneW, segments)
	if len(segments) == 0 {
		segments = []sceneSegment{{Left: 0, Right: sceneW}}
	}
	allocations := allocatePetsToSegments(count, segments)
	positions := make([]int, 0, count)
	for i, segment := range segments {
		n := allocations[i]
		if n <= 0 {
			continue
		}
		leftLimit := segment.Left
		rightLimit := segment.Right - spriteW
		if rightLimit < leftLimit {
			rightLimit = leftLimit
		}
		margin := min(24, max(0, (rightLimit-leftLimit)/4))
		left := leftLimit + margin
		right := rightLimit - margin
		if right < left {
			left = leftLimit
			right = rightLimit
		}
		for j := 0; j < n; j++ {
			x := (left + right) / 2
			if n > 1 {
				x = left + (right-left)*j/(n-1)
			}
			positions = append(positions, clamp(x, 0, max(0, sceneW-spriteW)))
		}
	}
	for len(positions) < count {
		positions = append(positions, clamp(sceneW/2-spriteW/2, 0, max(0, sceneW-spriteW)))
	}
	return positions[:count]
}

func normalizeSceneSegments(sceneW int, segments []sceneSegment) []sceneSegment {
	out := make([]sceneSegment, 0, len(segments))
	for _, segment := range segments {
		left := clamp(segment.Left, 0, sceneW)
		right := clamp(segment.Right, 0, sceneW)
		if right-left >= spriteW || (sceneW < spriteW && right > left) {
			out = append(out, sceneSegment{Left: left, Right: right})
		}
	}
	return mergeSceneSegments(out)
}

func allocatePetsToSegments(count int, segments []sceneSegment) []int {
	allocations := make([]int, len(segments))
	if count <= 0 || len(segments) == 0 {
		return allocations
	}
	type remainder struct {
		index int
		value float64
	}
	totalWidth := 0
	for _, segment := range segments {
		totalWidth += max(0, segment.Right-segment.Left)
	}
	if totalWidth <= 0 {
		allocations[0] = count
		return allocations
	}
	remaining := count
	if count >= len(segments) {
		for i := range allocations {
			allocations[i] = 1
		}
		remaining -= len(segments)
	}
	remainders := make([]remainder, 0, len(segments))
	assigned := 0
	for i, segment := range segments {
		share := float64(max(0, segment.Right-segment.Left)) * float64(remaining) / float64(totalWidth)
		whole := int(math.Floor(share))
		allocations[i] += whole
		assigned += whole
		remainders = append(remainders, remainder{index: i, value: share - float64(whole)})
	}
	sort.SliceStable(remainders, func(i, j int) bool {
		if remainders[i].value != remainders[j].value {
			return remainders[i].value > remainders[j].value
		}
		return remainders[i].index < remainders[j].index
	})
	for i := 0; i < remaining-assigned; i++ {
		allocations[remainders[i%len(remainders)].index]++
	}
	return allocations
}

func (a *petApp) setCoatMode(mode coatMode) {
	a.coatMode = mode
	a.refreshPetVariants()
}

func (a *petApp) setFixedVariant(variant int) {
	a.variant = clamp(variant, 0, len(variants)-1)
	if a.coatMode == coatFixed {
		a.refreshPetVariants()
	}
}

func (a *petApp) setSelectedVariant(index int, variant int) {
	if index < 0 || index >= maxPetCount {
		return
	}
	a.selectedCoats[index] = clamp(variant, 0, len(variants)-1)
	a.selectedSlotModes[index] = petSlotFixed
	if a.coatMode == coatSelected && index < len(a.pets) {
		a.pets[index].variant = a.selectedCoats[index]
	}
}

func (a *petApp) setSelectedRandomGroup(index int, groupID string) {
	if index < 0 || index >= maxPetCount {
		return
	}
	a.selectedSlotModes[index] = petSlotRandom
	a.selectedRandomGroupIDs[index] = normalizeVariantGroupID(groupID)
	if a.coatMode == coatSelected && index < len(a.pets) {
		a.pets[index].variant = a.variantForIndex(index)
	}
}

func (a *petApp) setRandomGroupID(groupID string) {
	a.randomGroupID = normalizeVariantGroupID(groupID)
	if a.coatMode == coatRandom {
		a.refreshPetVariants()
	}
}

func (a *petApp) refreshPetVariants() {
	for i := range a.pets {
		a.pets[i].variant = a.variantForIndex(i)
	}
}

func (a *petApp) variantForIndex(index int) int {
	if len(variants) == 0 {
		return 0
	}
	switch a.coatMode {
	case coatRandom:
		return randomVariantForGroup(a.randomGroupID)
	case coatSelected:
		if index >= 0 && index < len(a.selectedCoats) {
			if normalizePetSlotMode(int(a.selectedSlotModes[index])) == petSlotRandom {
				return randomVariantForGroup(a.selectedRandomGroupIDs[index])
			}
			return clamp(a.selectedCoats[index], 0, len(variants)-1)
		}
	}
	return clamp(a.variant, 0, len(variants)-1)
}

func randomVariantForGroup(groupID string) int {
	if len(variants) == 0 {
		return 0
	}
	indices := variantIndicesForGroupID(normalizeVariantGroupID(groupID))
	if len(indices) == 0 {
		return rand.Intn(len(variants))
	}
	return indices[rand.Intn(len(indices))]
}

func defaultSelectedCoats() [maxPetCount]int {
	return [maxPetCount]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4}
}

func defaultPetSizes() [maxPetCount]int {
	var sizes [maxPetCount]int
	for i := range sizes {
		sizes[i] = defaultPetSizePercent
	}
	return sizes
}

func normalizePetSizePercent(size int) int {
	if size <= 0 {
		return defaultPetSizePercent
	}
	size = clamp(size, minPetSizePercent, maxPetSizePercent)
	return ((size + petSizeStepPercent/2) / petSizeStepPercent) * petSizeStepPercent
}

func petSpriteSizeForPercent(percent int) (int, int) {
	percent = normalizePetSizePercent(percent)
	return max(1, frameW*percent/100), max(1, frameH*percent/100)
}

func (a *petApp) petSizePercent(index int) int {
	if index < 0 || index >= len(a.petSizes) {
		return defaultPetSizePercent
	}
	return normalizePetSizePercent(a.petSizes[index])
}

func (a *petApp) petSpriteSize(index int) (int, int) {
	return petSpriteSizeForPercent(a.petSizePercent(index))
}

func (a *petApp) petDrawY(index int, p *desktopPet) int {
	_, h := a.petSpriteSize(index)
	return sceneH - h - p.laneOffset
}

func (a *petApp) setPetSize(index int, size int) {
	if index < 0 || index >= len(a.petSizes) {
		return
	}
	a.petSizes[index] = normalizePetSizePercent(size)
	if index < len(a.pets) {
		w, _ := a.petSpriteSize(index)
		a.pets[index].x = clamp(a.pets[index].x, 0, max(0, a.sceneW-w))
	}
}

func (a *petApp) petVariant(p *desktopPet) int {
	if len(variants) == 0 {
		return 0
	}
	return clamp(p.variant, 0, len(variants)-1)
}

func (a *petApp) petMotionProfile(p *desktopPet) string {
	if len(variants) == 0 {
		return ""
	}
	return catalog.MotionProfileForVariant(variants[a.petVariant(p)])
}

func (a *petApp) petWheelCapable(p *desktopPet) bool {
	if len(variants) == 0 {
		return false
	}
	return catalog.WheelCapableVariant(variants[a.petVariant(p)])
}

func (a *petApp) newPet(index int) desktopPet {
	w, _ := a.petSpriteSize(index)
	spread := max(w+24, a.sceneW/max(1, a.petCount+1))
	dir := 1
	if a.bidirectional && index%2 == 1 {
		dir = -1
	}
	x := -w - index*spread - rand.Intn(80)
	if dir < 0 {
		x = a.sceneW + index*spread + rand.Intn(80)
	}
	p := desktopPet{
		x:          x,
		laneOffset: (index % 3) * 5,
		variant:    a.variantForIndex(index),
		item:       noItem,
		carryKind:  noItem,
		motionSet:  rand.Intn(motionSets),
		state:      stateWalk,
		moveSpeed:  max(1, a.speed-1+rand.Intn(2)),
		stateTicks: 30 + rand.Intn(80),
		dir:        dir,
		nextDir:    dir,
	}
	if index == 0 {
		p.x = rand.Intn(max(1, a.sceneW-w))
	}
	a.chooseRandomAction(&p)
	return p
}

func (a *petApp) resetPetAtLeft(index int, p *desktopPet) {
	a.resetPetAtEdge(index, p, 1)
}

func (a *petApp) resetPetAtEdge(index int, p *desktopPet, dir int) {
	a.releaseForage(index, p)
	w, _ := a.petSpriteSize(index)
	if dir < 0 {
		p.x = a.sceneW + rand.Intn(120)
	} else {
		p.x = -w - rand.Intn(120)
	}
	p.frame = 0
	p.motionSet = rand.Intn(motionSets)
	p.variant = a.variantForIndex(index)
	p.item = noItem
	p.carryKind = noItem
	p.state = stateWalk
	p.prevState = stateWalk
	p.moveSpeed = max(1, a.speed-1+rand.Intn(2))
	p.stateTicks = 40 + rand.Intn(90)
	p.dir = normalizeDir(dir)
	p.nextDir = p.dir
}

func (a *petApp) startTurn(p *desktopPet, nextDir int, after behaviorState) {
	nextDir = normalizeDir(nextDir)
	if p.dir == 0 {
		p.dir = 1
	}
	if p.dir == nextDir {
		return
	}
	p.prevState = after
	p.state = stateTurn
	p.nextDir = nextDir
	p.moveSpeed = 0
	p.stateTicks = turnTicks
	p.frame = 0
	p.item = noItem
	p.carryKind = noItem
}

func (a *petApp) finishTurn(p *desktopPet) {
	p.dir = normalizeDir(p.nextDir)
	p.nextDir = p.dir
	p.frame = 0
	p.motionSet = rand.Intn(motionSets)
	p.state = p.prevState
	if p.state == stateTurn || p.state == stateWheel || p.state == stateGroom || p.state == stateForage || p.state == stateCarry {
		p.state = stateWalk
	}
	p.moveSpeed = max(1, a.speed-1+rand.Intn(2))
	p.stateTicks = 32 + rand.Intn(70)
}

func (a *petApp) setBidirectional(enabled bool) {
	a.bidirectional = enabled
	if enabled {
		return
	}
	for i := range a.pets {
		p := &a.pets[i]
		p.dir = 1
		p.nextDir = 1
		p.timed.Left = false
		if p.timed.Velocity < 0 {
			p.timed.Velocity = 0
		}
		if p.state == stateTurn {
			p.state = stateWalk
			p.moveSpeed = max(1, a.speed-1)
			p.stateTicks = 24 + rand.Intn(40)
		}
	}
}

func (a *petApp) ensureForageItems() {
	for len(a.forage) < maxForage {
		a.forage = append(a.forage, forageItem{owner: noItem})
	}
	if a.tickCount%90 != 0 && a.tickCount != 0 {
		return
	}
	for i := range a.forage {
		if a.forage[i].active || a.forage[i].owner != noItem {
			continue
		}
		if rand.Intn(100) > 45 {
			continue
		}
		x := 28 + rand.Intn(max(1, a.sceneW-56))
		if abs(x-a.wheelX) < wheelSize {
			x = clamp(x+wheelSize+24, 24, max(24, a.sceneW-24))
		}
		a.forage[i] = forageItem{
			x:      x,
			kind:   rand.Intn(3),
			owner:  noItem,
			active: true,
		}
	}
}

func (a *petApp) clearForageItems() {
	for i := range a.forage {
		a.forage[i].owner = noItem
		a.forage[i].active = false
	}
	for i := range a.pets {
		p := &a.pets[i]
		switch p.state {
		case stateForage, stateNibble, stateEat, stateDig, stateCarry:
			p.item = noItem
			p.carryKind = noItem
			p.state = stateWalk
			p.moveSpeed = max(1, a.speed-1)
			p.stateTicks = 24
		default:
			p.item = noItem
			p.carryKind = noItem
		}
	}
}

func (a *petApp) maybeAssignForageTarget(p *desktopPet) bool {
	if !foragePropsEnabled {
		return false
	}
	if p.item != noItem || p.state == stateWheel {
		return false
	}
	best := noItem
	bestDistance := a.sceneW + spriteW
	dir := normalizeDir(p.dir)
	for i, item := range a.forage {
		if !item.active || item.owner != noItem {
			continue
		}
		mouthX := p.x + spriteW - 22
		distance := item.x - mouthX
		if dir < 0 {
			mouthX = p.x + 22
			distance = mouthX - item.x
		}
		if distance < 12 || distance > bestDistance {
			continue
		}
		best = i
		bestDistance = distance
	}
	if best == noItem {
		return false
	}
	a.forage[best].owner = reservedItem
	p.item = best
	p.carryKind = noItem
	p.state = stateForage
	p.dir = dir
	p.nextDir = dir
	p.moveSpeed = max(1, a.speed-1)
	p.stateTicks = max(45, bestDistance/max(1, p.moveSpeed)+36)
	return true
}

func (a *petApp) maybeStartGnawing(index int, p *desktopPet) {
	if p.item < 0 || p.item >= len(a.forage) {
		a.releaseForage(index, p)
		a.chooseRandomAction(p)
		return
	}
	item := &a.forage[p.item]
	item.owner = index
	if !item.active {
		a.releaseForage(index, p)
		a.chooseRandomAction(p)
		return
	}
	mouthX := p.x + spriteW - 22
	if p.dir < 0 {
		mouthX = p.x + 22
		if mouthX > item.x {
			return
		}
		p.x = clamp(item.x-22, 0, max(0, a.sceneW-spriteW))
	} else {
		if mouthX < item.x {
			return
		}
		p.x = clamp(item.x-spriteW+22, 0, max(0, a.sceneW-spriteW))
	}
	p.state = stateNibble
	p.frame = 0
	p.moveSpeed = 0
	p.stateTicks = 28 + rand.Intn(34)
}

func (a *petApp) finishGnawing(index int, p *desktopPet) {
	item := &a.forage[p.item]
	kind := item.kind
	item.active = false
	item.owner = index
	if kind == 0 || rand.Intn(100) < 54 {
		p.state = stateEat
		p.frame = 0
		p.carryKind = kind
		p.moveSpeed = 0
		p.stateTicks = 28 + rand.Intn(24)
		return
	}
	if kind == 2 && rand.Intn(100) < 58 {
		p.state = stateDig
		p.frame = 0
		p.carryKind = kind
		p.moveSpeed = 0
		p.stateTicks = 30 + rand.Intn(26)
		return
	}
	if rand.Intn(100) < 58 {
		p.state = stateCarry
		p.frame = 0
		p.carryKind = kind
		p.moveSpeed = max(1, a.speed-1+rand.Intn(2))
		p.stateTicks = 26 + rand.Intn(44)
		return
	}
	a.releaseForage(index, p)
	a.chooseRandomAction(p)
}

func (a *petApp) finishEating(index int, p *desktopPet) {
	a.releaseForage(index, p)
	p.carryKind = noItem
	if a.mode == modeRandom {
		a.chooseRandomAction(p)
		return
	}
	p.state = stateIdle
	p.moveSpeed = 0
	p.stateTicks = 12
}

func (a *petApp) releaseForage(index int, p *desktopPet) {
	if p.item >= 0 && p.item < len(a.forage) && (a.forage[p.item].owner == index || a.forage[p.item].owner == reservedItem) {
		a.forage[p.item].owner = noItem
		a.forage[p.item].active = false
	}
	p.item = noItem
	p.carryKind = noItem
}

func (a *petApp) syncNearbyWalkers() {
	for i := 0; i < len(a.pets); i++ {
		for j := i + 1; j < len(a.pets); j++ {
			pi := &a.pets[i]
			pj := &a.pets[j]
			if pi.timed.ID != "" || pj.timed.ID != "" {
				continue
			}
			if pi.state != stateWalk || pj.state != stateWalk {
				continue
			}
			if abs(pi.x-pj.x) > 72 {
				continue
			}
			speed := max(1, min(pi.moveSpeed, pj.moveSpeed))
			pi.moveSpeed = speed
			pj.moveSpeed = speed
			if normalizeDir(pj.dir) != normalizeDir(pi.dir) {
				a.startTurn(pj, pi.dir, stateWalk)
			}
		}
	}
}

func (a *petApp) maybeStartSocial() {
	if len(a.pets) < 2 || a.tickCount%24 != 0 || rand.Intn(100) > 28 {
		return
	}
	for i := 0; i < len(a.pets); i++ {
		for j := i + 1; j < len(a.pets); j++ {
			pi := &a.pets[i]
			pj := &a.pets[j]
			if pi.timed.ID != "" || pj.timed.ID != "" {
				continue
			}
			if !canSocialize(pi) || !canSocialize(pj) {
				continue
			}
			if abs(pi.x-pj.x) > 84 {
				continue
			}
			if normalizeDir(pi.dir) != 1 || normalizeDir(pj.dir) != -1 {
				if normalizeDir(pi.dir) != 1 {
					a.startTurn(pi, 1, stateWalk)
				}
				if normalizeDir(pj.dir) != -1 {
					a.startTurn(pj, -1, stateWalk)
				}
				return
			}
			anchor := min(pi.x, pj.x)
			piW, _ := a.petSpriteSize(i)
			pjW, _ := a.petSpriteSize(j)
			pi.x = clamp(anchor, 0, max(0, a.sceneW-piW-36))
			pj.x = clamp(pi.x+34+rand.Intn(14), 0, max(0, a.sceneW-pjW))
			pj.laneOffset = pi.laneOffset
			ticks := 44 + rand.Intn(70)
			pi.state = stateGroom
			pj.state = stateGroom
			pi.moveSpeed = 0
			pj.moveSpeed = 0
			pi.frame = 0
			pj.frame = 3
			pi.stateTicks = ticks
			pj.stateTicks = ticks + rand.Intn(16)
			return
		}
	}
}

func canSocialize(p *desktopPet) bool {
	return p.item == noItem && (p.state == stateIdle || p.state == stateWalk || p.state == stateNibble)
}

func (a *petApp) hasWheelRunner() bool {
	for i := range a.pets {
		if a.pets[i].state == stateWheel {
			return true
		}
	}
	return false
}

func (a *petApp) drawForageItems(dst *image.RGBA, sceneLeft int, dpi int) {
	if !foragePropsEnabled {
		return
	}
	y := sceneH - 9
	for _, item := range a.forage {
		if !item.active {
			continue
		}
		a.drawForageProp(dst, unscaleForDPI(item.x-sceneLeft, dpi), y, item.kind)
	}
}

func (a *petApp) drawReactions(dst *image.RGBA, sceneLeft, sceneRight int, dpi int) {
	for _, reaction := range a.reactions {
		if reaction.pet < 0 || reaction.pet >= len(a.pets) {
			continue
		}
		p := a.pets[reaction.pet]
		if p.state == stateWheel {
			continue
		}
		petW, petH := a.petSpriteSize(reaction.pet)
		baseY := sceneH - petH - p.laneOffset
		scaledPetW := scaleForDPI(petW, dpi)
		globalX := p.x + scaledPetW/2 - scaleForDPI(18, dpi)
		if globalX+scaleForDPI(42, dpi) <= sceneLeft || globalX >= sceneRight {
			continue
		}
		x := clamp(unscaleForDPI(globalX-sceneLeft, dpi), 2, max(2, dst.Bounds().Dx()-42))
		y := clamp(baseY-26-(reactionTicks-reaction.ticks)/8, 0, sceneH-32)
		drawReactionBubble(dst, x, y, reaction.kind, reaction.ticks)
	}
}

func (a *petApp) showSettings() {
	a.clampSettingsWindowPosition()
	if a.settingsHwnd == 0 {
		a.createSettingsWindow()
	}
	a.syncSettingsWindow()
	win.ShowWindow(a.settingsHwnd, win.SW_SHOW)
	win.SetForegroundWindow(a.settingsHwnd)
}

func (a *petApp) createSettingsWindow() {
	title := a.txt("settingsTitle")
	a.ensureSettingsBrushes()
	style := uint32(win.WS_POPUP | win.WS_VISIBLE | win.WS_CLIPCHILDREN)
	hwnd := win.CreateWindowEx(
		win.WS_EX_TOOLWINDOW,
		syscall.StringToUTF16Ptr(windowClass),
		syscall.StringToUTF16Ptr(title),
		style,
		a.settingsX, a.settingsY, settingsClientW, settingsClientH,
		a.hwnd, 0, a.hinst, nil,
	)
	if hwnd == 0 {
		return
	}
	a.settingsHwnd = hwnd
	a.ensureSettingsFonts()

	a.createButton(hwnd, ctrlTopClose, "x", 716, 18, 28, 28, 0)
	a.createButton(hwnd, ctrlTabAnimals, a.txt("tabAnimals"), 24, 122, 154, 38, win.WS_GROUP)
	a.createButton(hwnd, ctrlTabMotion, a.txt("tabMotion"), 24, 170, 154, 38, 0)
	a.createButton(hwnd, ctrlTabDisplay, a.txt("tabDisplay"), 24, 218, 154, 38, 0)

	if a.settingsTab == tabAnimals {
		a.createButton(hwnd, ctrlPetMinus, "-", 408, 150, 42, 32, 0)
		a.createButton(hwnd, ctrlPetPlus, "+", 506, 150, 42, 32, 0)

		a.createButton(hwnd, ctrlCoatFixed, a.txt("coatFixed"), 250, 250, 110, 30, win.WS_GROUP)
		a.createButton(hwnd, ctrlCoatSelected, a.txt("coatSelected"), 374, 250, 130, 30, 0)
		a.createButton(hwnd, ctrlCoatRandom, a.txt("coatRandom"), 518, 250, 110, 30, 0)

		if a.coatMode != coatSelected {
			a.createButton(hwnd, ctrlVariantCombo, "", 360, 298, 318, 34, 0)
		}
		a.createButton(hwnd, ctrlNameLabels, "", 548, 344, 132, 28, 0)
		for i := 0; i < a.petCount; i++ {
			_, nameRect := settingsPetNameRects(i)
			a.createButton(hwnd, ctrlPetNameBase+int32(i), "", nameRect.Left, nameRect.Top, nameRect.Right-nameRect.Left, nameRect.Bottom-nameRect.Top, 0)
			if a.coatMode == coatSelected {
				_, buttonRect := settingsPetVariantRects(i)
				a.createButton(hwnd, ctrlPetVariantBase+int32(i), "", buttonRect.Left, buttonRect.Top, buttonRect.Right-buttonRect.Left, buttonRect.Bottom-buttonRect.Top, 0)
			}
			sizeRect := settingsPetSizeRect(i)
			a.createButton(hwnd, ctrlPetSizeBase+int32(i), "", sizeRect.Left, sizeRect.Top, sizeRect.Right-sizeRect.Left, sizeRect.Bottom-sizeRect.Top, 0)
		}
	} else if a.settingsTab == tabMotion {
		a.createButton(hwnd, ctrlModeKeyboard, a.txt("modeKeyboard"), 250, 164, 210, 32, win.WS_GROUP)
		a.createButton(hwnd, ctrlModeRandom, a.txt("modeRandom"), 478, 164, 210, 32, 0)

		a.createButton(hwnd, ctrlSpeedSlow, a.txt("speedSlow"), 250, 270, 118, 32, win.WS_GROUP)
		a.createButton(hwnd, ctrlSpeedNormal, a.txt("speedNormal"), 384, 270, 118, 32, 0)
		a.createButton(hwnd, ctrlSpeedFast, a.txt("speedFast"), 518, 270, 118, 32, 0)

		a.createButton(hwnd, ctrlTypingWheel, a.txt("typingWheel"), 250, 378, 210, 32, win.WS_GROUP)
		a.createButton(hwnd, ctrlBidirectional, a.txt("naturalTurns"), 478, 378, 210, 32, 0)
	} else {
		a.createButton(hwnd, ctrlDisplaySingle, a.settingsButtonLabel(ctrlDisplaySingle), 250, 154, 82, 30, win.WS_GROUP)
		a.createButton(hwnd, ctrlDisplaySpan, a.settingsButtonLabel(ctrlDisplaySpan), 340, 154, 96, 30, 0)
		a.createButton(hwnd, ctrlDisplaySpanLess, a.settingsButtonLabel(ctrlDisplaySpanLess), 444, 154, 74, 30, 0)
		a.createButton(hwnd, ctrlDisplaySpanMore, a.settingsButtonLabel(ctrlDisplaySpanMore), 526, 154, 74, 30, 0)
		a.createButton(hwnd, ctrlDisplayPrev, a.settingsButtonLabel(ctrlDisplayPrev), 608, 154, 38, 30, 0)
		a.createButton(hwnd, ctrlDisplayNext, a.settingsButtonLabel(ctrlDisplayNext), 650, 154, 38, 30, 0)

		a.createScrollBar(hwnd, ctrlRangeStartScroll, 322, 258, 288, 18)
		a.createScrollBar(hwnd, ctrlRangeEndScroll, 322, 292, 288, 18)
		a.createButton(hwnd, ctrlRangeFull, a.settingsButtonLabel(ctrlRangeFull), 250, 326, 74, 30, win.WS_GROUP)
		a.createButton(hwnd, ctrlRangeNarrow, a.settingsButtonLabel(ctrlRangeNarrow), 334, 326, 74, 30, 0)
		a.createButton(hwnd, ctrlRangeWide, a.settingsButtonLabel(ctrlRangeWide), 418, 326, 74, 30, 0)
		a.createButton(hwnd, ctrlRangeLeft, a.settingsButtonLabel(ctrlRangeLeft), 502, 326, 74, 30, 0)
		a.createButton(hwnd, ctrlRangeRight, a.settingsButtonLabel(ctrlRangeRight), 586, 326, 74, 30, 0)

		a.createButton(hwnd, ctrlPositionTaskbar, a.settingsButtonLabel(ctrlPositionTaskbar), 250, 426, 134, 30, win.WS_GROUP)
		a.createButton(hwnd, ctrlPositionBottom, a.settingsButtonLabel(ctrlPositionBottom), 392, 426, 160, 30, 0)
		a.createButton(hwnd, ctrlOffsetUp, a.settingsButtonLabel(ctrlOffsetUp), 560, 426, 62, 30, 0)
		a.createButton(hwnd, ctrlOffsetDown, a.settingsButtonLabel(ctrlOffsetDown), 628, 426, 62, 30, 0)
	}

	a.createButton(hwnd, ctrlLanguageCombo, "", 322, 574, 180, 34, 0)
	a.createButton(hwnd, ctrlReset, a.txt("reset"), 534, 576, 100, 32, 0)
	a.createButton(hwnd, ctrlClose, a.txt("close"), 646, 576, 78, 32, 0)
}

func (a *petApp) ensureSettingsBrushes() {
	if a.settingsBrush == 0 {
		a.settingsBrush = solidBrush(244, 242, 235)
	}
	if a.settingsCard == 0 {
		a.settingsCard = solidBrush(255, 255, 251)
	}
}

func (a *petApp) ensureSettingsFonts() {
	if a.settingsFont == 0 {
		a.settingsFont = makeSettingsFont("Yu Gothic UI", -13, win.FW_NORMAL)
	}
	if a.settingsTitleFont == 0 {
		a.settingsTitleFont = makeSettingsFont("Yu Gothic UI", -19, win.FW_SEMIBOLD)
	}
	if a.settingsSmallFont == 0 {
		a.settingsSmallFont = makeSettingsFont("Yu Gothic UI", -12, win.FW_NORMAL)
	}
}

func makeSettingsFont(face string, height int32, weight int32) win.HFONT {
	var lf win.LOGFONT
	lf.LfHeight = height
	lf.LfWeight = weight
	lf.LfQuality = 5
	name := syscall.StringToUTF16(face)
	copy(lf.LfFaceName[:], name)
	return win.CreateFontIndirect(&lf)
}

func solidBrush(r, g, b byte) win.HBRUSH {
	return win.CreateBrushIndirect(&win.LOGBRUSH{
		LbStyle: win.BS_SOLID,
		LbColor: win.RGB(r, g, b),
	})
}

func (a *petApp) paintSettingsWindow(hwnd win.HWND) {
	a.ensureSettingsBrushes()
	a.ensureSettingsFonts()
	var ps win.PAINTSTRUCT
	hdc := win.BeginPaint(hwnd, &ps)
	if hdc == 0 {
		return
	}
	defer win.EndPaint(hwnd, &ps)

	drawRectFill(hdc, win.RECT{Left: 0, Top: 0, Right: settingsClientW, Bottom: settingsClientH}, rgb(246, 248, 244))
	drawRectFill(hdc, win.RECT{Left: 0, Top: 0, Right: 204, Bottom: settingsClientH}, rgb(22, 45, 38))
	drawRectFill(hdc, win.RECT{Left: 204, Top: 0, Right: settingsClientW, Bottom: settingsClientH}, rgb(247, 248, 244))

	drawRoundFill(hdc, win.RECT{Left: 226, Top: 96, Right: 736, Bottom: 566}, rgb(255, 255, 251), 18)
	drawRoundFill(hdc, win.RECT{Left: 226, Top: 574, Right: 736, Bottom: 612}, rgb(255, 255, 251), 14)

	if a.settingsTab == tabAnimals {
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 142, Right: 708, Bottom: 192}, rgb(238, 242, 237), 14)
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 238, Right: 708, Bottom: 292}, rgb(238, 242, 237), 14)
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 338, Right: 708, Bottom: 562}, rgb(238, 242, 237), 14)
	} else if a.settingsTab == tabMotion {
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 150, Right: 708, Bottom: 208}, rgb(238, 242, 237), 14)
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 256, Right: 708, Bottom: 314}, rgb(238, 242, 237), 14)
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 364, Right: 708, Bottom: 424}, rgb(238, 242, 237), 14)
	} else {
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 150, Right: 708, Bottom: 212}, rgb(238, 242, 237), 14)
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 246, Right: 708, Bottom: 366}, rgb(238, 242, 237), 14)
		drawRoundFill(hdc, win.RECT{Left: 238, Top: 416, Right: 708, Bottom: 466}, rgb(238, 242, 237), 14)
	}

	drawTextLine(hdc, a.txt("settingsHeader"), win.RECT{Left: 24, Top: 28, Right: 178, Bottom: 56}, a.settingsTitleFont, rgb(245, 250, 244), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)
	drawTextLine(hdc, a.txt("settingsLead"), win.RECT{Left: 24, Top: 58, Right: 178, Bottom: 96}, a.settingsSmallFont, rgb(183, 207, 195), win.DT_LEFT|win.DT_WORDBREAK|win.DT_NOPREFIX)
	drawRoundFill(hdc, win.RECT{Left: 24, Top: 534, Right: 178, Bottom: 580}, rgb(35, 62, 52), 14)
	drawTextLine(hdc, a.settingsSidebarStatus(), win.RECT{Left: 38, Top: 546, Right: 164, Bottom: 568}, a.settingsSmallFont, rgb(231, 241, 233), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)

	drawTextLine(hdc, a.settingsPageTitle(), win.RECT{Left: 226, Top: 28, Right: 620, Bottom: 56}, a.settingsTitleFont, rgb(27, 36, 32), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
	drawTextLine(hdc, a.settingsPageLead(), win.RECT{Left: 226, Top: 56, Right: 660, Bottom: 82}, a.settingsSmallFont, rgb(91, 104, 96), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)

	labelColor := rgb(50, 61, 55)
	if a.settingsTab == tabAnimals {
		drawTextLine(hdc, a.txt("animalSection"), win.RECT{Left: 246, Top: 116, Right: 470, Bottom: 140}, a.settingsFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.txt("petCount"), win.RECT{Left: 250, Top: 154, Right: 390, Bottom: 184}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, fmt.Sprintf("%d", a.petCount), win.RECT{Left: 464, Top: 154, Right: 498, Bottom: 184}, a.settingsFont, rgb(25, 49, 40), win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.txt("coatMode"), win.RECT{Left: 250, Top: 214, Right: 420, Bottom: 238}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		if a.coatMode != coatSelected {
			drawTextLine(hdc, a.animalChoiceLabel(), win.RECT{Left: 250, Top: 302, Right: 352, Bottom: 328}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		}
		drawTextLine(hdc, a.petNameSectionLabel(), win.RECT{Left: 250, Top: 344, Right: 520, Bottom: 370}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		if a.petCount > 0 {
			for i := 0; i < a.petCount; i++ {
				numberRect, _ := settingsPetNameRects(i)
				drawTextLine(hdc, fmt.Sprintf("%d", i+1), numberRect, a.settingsSmallFont, rgb(69, 78, 72), win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
			}
		} else {
			drawTextLine(hdc, a.localText("オンにすると、名前の編集とカーソルホバー表示を使えます。", "Turn this on to edit names and show them on hover."), win.RECT{Left: 254, Top: 386, Right: 682, Bottom: 430}, a.settingsSmallFont, rgb(69, 78, 72), win.DT_LEFT|win.DT_WORDBREAK|win.DT_NOPREFIX)
		}
	} else if a.settingsTab == tabMotion {
		drawTextLine(hdc, a.txt("mode"), win.RECT{Left: 246, Top: 126, Right: 400, Bottom: 150}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.txt("speed"), win.RECT{Left: 246, Top: 232, Right: 400, Bottom: 256}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.txt("motion"), win.RECT{Left: 246, Top: 340, Right: 400, Bottom: 364}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
	} else {
		drawTextLine(hdc, a.localText("表示範囲", "Display scope"), win.RECT{Left: 246, Top: 126, Right: 400, Bottom: 150}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.displaySummary(), win.RECT{Left: 250, Top: 186, Right: 688, Bottom: 208}, a.settingsFont, rgb(27, 36, 32), win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)

		drawTextLine(hdc, a.walkRangeSectionLabel(), win.RECT{Left: 246, Top: 222, Right: 390, Bottom: 244}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.walkRangeSummary(), win.RECT{Left: 394, Top: 222, Right: 688, Bottom: 244}, a.settingsSmallFont, rgb(91, 104, 96), win.DT_RIGHT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)
		a.drawWalkRangePreview(hdc)
		drawTextLine(hdc, a.localText("左端", "Left edge"), win.RECT{Left: 250, Top: 254, Right: 316, Bottom: 280}, a.settingsSmallFont, labelColor, win.DT_RIGHT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.localText("右端", "Right edge"), win.RECT{Left: 250, Top: 288, Right: 316, Bottom: 314}, a.settingsSmallFont, labelColor, win.DT_RIGHT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
		drawTextLine(hdc, fmt.Sprintf("%d%%", start), win.RECT{Left: 616, Top: 254, Right: 680, Bottom: 280}, a.settingsSmallFont, rgb(91, 104, 96), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, fmt.Sprintf("%d%%", end), win.RECT{Left: 616, Top: 288, Right: 680, Bottom: 314}, a.settingsSmallFont, rgb(91, 104, 96), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)

		drawTextLine(hdc, a.localText("表示位置", "Position"), win.RECT{Left: 246, Top: 398, Right: 380, Bottom: 420}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		drawTextLine(hdc, a.positionSummary(), win.RECT{Left: 382, Top: 398, Right: 690, Bottom: 420}, a.settingsSmallFont, rgb(91, 104, 96), win.DT_RIGHT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX|win.DT_END_ELLIPSIS)
	}
	drawTextLine(hdc, a.txt("language"), win.RECT{Left: 226, Top: 582, Right: 314, Bottom: 606}, a.settingsSmallFont, labelColor, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
}

func (a *petApp) drawWalkRangePreview(hdc win.HDC) {
	track := win.RECT{Left: 250, Top: 246, Right: 688, Bottom: 252}
	drawRoundFill(hdc, track, rgb(215, 221, 213), 6)
	a.drawDisplaySpanMarkers(hdc, track)
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	width := int(track.Right - track.Left)
	selected := win.RECT{
		Left:   track.Left + int32(width*start/100),
		Top:    track.Top - 2,
		Right:  track.Left + int32(width*end/100),
		Bottom: track.Bottom + 2,
	}
	drawRoundFill(hdc, selected, rgb(70, 104, 87), 8)
	drawRoundOutline(hdc, win.RECT{Left: 246, Top: 242, Right: 692, Bottom: 256}, rgb(230, 233, 224), 10)
}

func (a *petApp) drawDisplaySpanMarkers(hdc win.HDC, track win.RECT) {
	areas := displayAreaForScope(a.displayScope)
	if len(areas) <= 1 {
		return
	}
	scope, start, end := a.normalizedDisplaySelection(len(areas))
	if scope != displayScopeSpan || start == end {
		return
	}
	selected := areas[start : end+1]
	combined := combineDisplayAreas(selected)
	totalW := max(1, int(combined.Screen.Right-combined.Screen.Left))
	trackW := max(1, int(track.Right-track.Left))
	for _, area := range selected[1:] {
		x := track.Left + int32(trackW*int(area.Screen.Left-combined.Screen.Left)/totalW)
		drawRectFill(hdc, win.RECT{Left: x - 1, Top: track.Top - 5, Right: x + 1, Bottom: track.Bottom + 5}, rgb(171, 181, 169))
	}
}

func (a *petApp) createStatic(parent win.HWND, text string, x, y, w, h int32) win.HWND {
	return a.createStaticWithFont(parent, text, x, y, w, h, a.settingsFont)
}

func (a *petApp) createStaticWithFont(parent win.HWND, text string, x, y, w, h int32, font win.HFONT) win.HWND {
	hwnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr("STATIC"),
		syscall.StringToUTF16Ptr(text),
		win.WS_CHILD|win.WS_VISIBLE|win.SS_LEFT,
		x, y, w, h,
		parent, 0, a.hinst, nil,
	)
	a.setControlFont(hwnd, font)
	return hwnd
}

func (a *petApp) createButton(parent win.HWND, id int32, text string, x, y, w, h int32, style uint32) win.HWND {
	buttonStyle := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.BS_OWNERDRAW)
	if style&win.WS_GROUP != 0 {
		buttonStyle |= win.WS_GROUP
	}
	hwnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr("BUTTON"),
		syscall.StringToUTF16Ptr(text),
		buttonStyle,
		x, y, w, h,
		parent, win.HMENU(uintptr(id)), a.hinst, nil,
	)
	a.setControlFont(hwnd, a.settingsFont)
	return hwnd
}

func (a *petApp) createScrollBar(parent win.HWND, id int32, x, y, w, h int32) win.HWND {
	hwnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr("SCROLLBAR"),
		nil,
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP|sbsHorz,
		x, y, w, h,
		parent, win.HMENU(uintptr(id)), a.hinst, nil,
	)
	a.setControlFont(hwnd, a.settingsFont)
	return hwnd
}

func (a *petApp) createEdit(parent win.HWND, id int32, text string, x, y, w, h int32) win.HWND {
	hwnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr("EDIT"),
		syscall.StringToUTF16Ptr(text),
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP|win.WS_BORDER|win.ES_LEFT|win.ES_AUTOHSCROLL,
		x, y, w, h,
		parent, win.HMENU(uintptr(id)), a.hinst, nil,
	)
	a.setControlFont(hwnd, a.settingsFont)
	win.SendMessage(hwnd, win.EM_LIMITTEXT, 24, 0)
	return hwnd
}

func (a *petApp) createCombo(parent win.HWND, id int32, x, y, w, h int32) win.HWND {
	hwnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr("COMBOBOX"),
		nil,
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP|win.WS_VSCROLL|win.CBS_DROPDOWNLIST,
		x, y, w, h,
		parent, win.HMENU(uintptr(id)), a.hinst, nil,
	)
	a.setControlFont(hwnd, a.settingsFont)
	return hwnd
}

func (a *petApp) drawSettingsButton(dis *win.DRAWITEMSTRUCT) bool {
	if dis == nil {
		return false
	}
	id := int32(dis.CtlID)
	label := a.settingsButtonLabel(id)
	if label == "" {
		return false
	}

	r := dis.RcItem
	enabled := dis.ItemState&win.ODS_DISABLED == 0 && win.IsWindowEnabled(dis.HwndItem)
	pressed := dis.ItemState&win.ODS_SELECTED != 0
	selected := a.settingsButtonSelected(id)
	selectField := a.settingsSelectButton(id)
	coatSelect := a.settingsCoatSelectButton(id)

	fill := rgb(252, 250, 244)
	text := rgb(35, 42, 35)
	if a.settingsSidebarButton(id) {
		fill = rgb(35, 62, 52)
		text = rgb(224, 238, 229)
		if selected {
			fill = rgb(232, 242, 231)
			text = rgb(22, 45, 38)
		} else if pressed {
			fill = rgb(48, 82, 68)
		}
	} else if id == ctrlTopClose {
		fill = rgb(239, 232, 228)
		text = rgb(122, 47, 38)
		if pressed {
			fill = rgb(224, 210, 204)
		}
	} else if selectField {
		fill = rgb(255, 255, 251)
	} else if selected {
		fill = rgb(48, 97, 73)
		text = rgb(255, 255, 251)
	} else if pressed {
		fill = rgb(226, 232, 220)
	}
	if !enabled {
		fill = rgb(230, 227, 216)
		text = rgb(126, 124, 112)
	}

	drawRectFill(dis.HDC, r, a.settingsButtonBackplate(id))
	inset := win.RECT{Left: r.Left + 1, Top: r.Top + 1, Right: r.Right - 1, Bottom: r.Bottom - 1}
	drawRoundFill(dis.HDC, inset, fill, 10)
	if enabled && dis.ItemState&win.ODS_FOCUS != 0 {
		focusRect := win.RECT{Left: r.Left + 4, Top: r.Top + 4, Right: r.Right - 4, Bottom: r.Bottom - 4}
		win.DrawFocusRect(dis.HDC, &focusRect)
	}

	textRect := win.RECT{Left: r.Left + 14, Top: r.Top, Right: r.Right - 14, Bottom: r.Bottom}
	if id == ctrlTopClose {
		textRect = win.RECT{Left: r.Left, Top: r.Top, Right: r.Right, Bottom: r.Bottom}
	}
	if coatSelect {
		idx := a.settingsSelectVariant(id)
		drawVariantSwatch(dis.HDC, r.Left+13, r.Top+(r.Bottom-r.Top-16)/2, idx, enabled)
		textRect.Left += 26
	}
	if selectField {
		arrowRect := win.RECT{Left: r.Right - 24, Top: r.Top, Right: r.Right - 8, Bottom: r.Bottom}
		if id >= ctrlPetSizeBase && id < ctrlPetSizeBase+maxPetCount {
			textRect = win.RECT{Left: r.Left + 6, Top: r.Top, Right: r.Right - 16, Bottom: r.Bottom}
			arrowRect = win.RECT{Left: r.Right - 16, Top: r.Top, Right: r.Right - 4, Bottom: r.Bottom}
		}
		drawTextLine(dis.HDC, "v", arrowRect, a.settingsFont, text, win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
		textRect.Right -= 24
		if id >= ctrlPetSizeBase && id < ctrlPetSizeBase+maxPetCount {
			textRect.Right = r.Right - 16
		}
	}

	flags := uint32(win.DT_VCENTER | win.DT_SINGLELINE | win.DT_NOPREFIX | win.DT_END_ELLIPSIS)
	if !selectField {
		flags |= win.DT_CENTER
	}
	drawTextLine(dis.HDC, label, textRect, a.settingsFont, text, flags)
	return true
}

func (a *petApp) settingsButtonLabel(id int32) string {
	switch id {
	case ctrlTabAnimals:
		return a.txt("tabAnimals")
	case ctrlTabMotion:
		return a.txt("tabMotion")
	case ctrlTabDisplay:
		return a.txt("tabDisplay")
	case ctrlPetMinus:
		return "-"
	case ctrlPetPlus:
		return "+"
	case ctrlCoatFixed:
		return a.txt("coatFixed")
	case ctrlCoatSelected:
		return a.txt("coatSelected")
	case ctrlCoatRandom:
		return a.txt("coatRandom")
	case ctrlVariantCombo:
		if a.coatMode == coatRandom {
			return a.randomGroupDisplayLabel(a.randomGroupID)
		}
		return a.variantDisplayLabel(a.variant)
	case ctrlModeKeyboard:
		return a.txt("modeKeyboard")
	case ctrlModeRandom:
		return a.txt("modeRandom")
	case ctrlSpeedSlow:
		return a.txt("speedSlow")
	case ctrlSpeedNormal:
		return a.txt("speedNormal")
	case ctrlSpeedFast:
		return a.txt("speedFast")
	case ctrlTypingWheel:
		return a.txt("typingWheel")
	case ctrlBidirectional:
		return a.txt("naturalTurns")
	case ctrlPositionTaskbar:
		return a.localText("タスクバー基準", "Taskbar edge")
	case ctrlPositionBottom:
		return a.localText("画面下端", "Screen bottom")
	case ctrlOffsetUp:
		return a.localText("上へ", "Up")
	case ctrlOffsetDown:
		return a.localText("下へ", "Down")
	case ctrlDisplayPrev:
		return "<"
	case ctrlDisplayNext:
		return ">"
	case ctrlDisplaySingle:
		return a.localText("1画面", "One")
	case ctrlDisplaySpan:
		return a.localText("複数画面", "Multi")
	case ctrlDisplaySpanLess:
		return a.localText("短く", "Less")
	case ctrlDisplaySpanMore:
		return a.localText("広く", "More")
	case ctrlRangeFull:
		if normalizeDisplayScope(int(a.displayScope)) == displayScopeSpan {
			return a.localText("全画面", "All")
		}
		return a.localText("全幅", "Full")
	case ctrlRangeNarrow:
		return a.localText("狭く", "Narrow")
	case ctrlRangeWide:
		return a.localText("広く", "Wider")
	case ctrlRangeLeft:
		return a.localText("左へ", "Left")
	case ctrlRangeRight:
		return a.localText("右へ", "Right")
	case ctrlLanguageCombo:
		if a.lang == langEnglish {
			return "English"
		}
		return "日本語"
	case ctrlTopClose:
		return "X"
	case ctrlReset:
		return a.txt("reset")
	case ctrlClose:
		return a.txt("close")
	case ctrlNameLabels:
		return a.localText("名前を表示", "Show names")
	case ctrlRenameOK:
		return a.localText("保存", "Save")
	case ctrlRenameCancel:
		return a.localText("キャンセル", "Cancel")
	}
	if id >= ctrlPetVariantBase && id < ctrlPetVariantBase+maxPetCount {
		return a.petAssignmentDisplayLabel(int(id - ctrlPetVariantBase))
	}
	if id >= ctrlPetNameBase && id < ctrlPetNameBase+maxPetCount {
		return a.petDisplayName(int(id - ctrlPetNameBase))
	}
	if id >= ctrlPetSizeBase && id < ctrlPetSizeBase+maxPetCount {
		return fmt.Sprintf("%d%%", a.petSizePercent(int(id-ctrlPetSizeBase)))
	}
	return ""
}

func (a *petApp) settingsButtonSelected(id int32) bool {
	switch id {
	case ctrlTabAnimals:
		return a.settingsTab == tabAnimals
	case ctrlTabMotion:
		return a.settingsTab == tabMotion
	case ctrlTabDisplay:
		return a.settingsTab == tabDisplay
	case ctrlCoatFixed:
		return a.coatMode == coatFixed
	case ctrlCoatSelected:
		return a.coatMode == coatSelected
	case ctrlCoatRandom:
		return a.coatMode == coatRandom
	case ctrlModeKeyboard:
		return a.mode == modeKeyboard
	case ctrlModeRandom:
		return a.mode == modeRandom
	case ctrlSpeedSlow:
		return a.speed == 2
	case ctrlSpeedNormal:
		return a.speed == 3
	case ctrlSpeedFast:
		return a.speed == 5
	case ctrlTypingWheel:
		return a.wheelEnabled
	case ctrlBidirectional:
		return a.bidirectional
	case ctrlPositionTaskbar:
		return normalizeOverlayPositionMode(int(a.positionMode)) == positionTaskbarEdge
	case ctrlPositionBottom:
		return normalizeOverlayPositionMode(int(a.positionMode)) == positionScreenBottom
	case ctrlDisplaySingle:
		return normalizeDisplayScope(int(a.displayScope)) == displayScopeSingle
	case ctrlDisplaySpan:
		return normalizeDisplayScope(int(a.displayScope)) == displayScopeSpan
	case ctrlNameLabels:
		return a.nameLabels
	}
	return false
}

func (a *petApp) settingsSelectButton(id int32) bool {
	return id == ctrlVariantCombo ||
		id == ctrlLanguageCombo ||
		(id >= ctrlPetVariantBase && id < ctrlPetVariantBase+maxPetCount) ||
		(id >= ctrlPetSizeBase && id < ctrlPetSizeBase+maxPetCount)
}

func (a *petApp) settingsSidebarButton(id int32) bool {
	return id == ctrlTabAnimals || id == ctrlTabMotion || id == ctrlTabDisplay
}

func (a *petApp) settingsButtonBackplate(id int32) settingsRGB {
	switch id {
	case ctrlTabAnimals, ctrlTabMotion, ctrlTabDisplay:
		return rgb(22, 45, 38)
	case ctrlTopClose:
		return rgb(247, 248, 244)
	case ctrlCoatFixed, ctrlCoatSelected, ctrlCoatRandom,
		ctrlModeKeyboard, ctrlModeRandom,
		ctrlSpeedSlow, ctrlSpeedNormal, ctrlSpeedFast,
		ctrlTypingWheel, ctrlBidirectional,
		ctrlPositionTaskbar, ctrlPositionBottom, ctrlOffsetUp, ctrlOffsetDown,
		ctrlDisplaySingle, ctrlDisplaySpan, ctrlDisplaySpanLess, ctrlDisplaySpanMore,
		ctrlDisplayPrev, ctrlDisplayNext,
		ctrlRangeFull, ctrlRangeNarrow, ctrlRangeWide, ctrlRangeLeft, ctrlRangeRight,
		ctrlNameLabels, ctrlRenameOK, ctrlRenameCancel:
		return rgb(235, 232, 220)
	}
	if (id >= ctrlPetVariantBase && id < ctrlPetVariantBase+maxPetCount) ||
		(id >= ctrlPetNameBase && id < ctrlPetNameBase+maxPetCount) ||
		(id >= ctrlPetSizeBase && id < ctrlPetSizeBase+maxPetCount) {
		return rgb(235, 232, 220)
	}
	return rgb(255, 255, 251)
}

func (a *petApp) settingsSidebarStatus() string {
	if a.settingsSaveFailed {
		if a.lang == langEnglish {
			return "Settings save warning"
		}
		return "設定保存に失敗"
	}
	if a.keyHookFailed || a.mouseHookFailed {
		if a.lang == langEnglish {
			return "Input hook warning"
		}
		return "入力反応に問題"
	}
	if a.lang == langEnglish {
		if a.coatMode == coatRandom {
			return fmt.Sprintf("%d pets / random", a.petCount)
		}
		if a.coatMode == coatSelected {
			return fmt.Sprintf("%d pets / custom", a.petCount)
		}
		return fmt.Sprintf("%d pets / fixed", a.petCount)
	}
	if a.coatMode == coatRandom {
		return fmt.Sprintf("%d匹 / ランダム", a.petCount)
	}
	if a.coatMode == coatSelected {
		return fmt.Sprintf("%d匹 / 個別指定", a.petCount)
	}
	return fmt.Sprintf("%d匹 / 固定", a.petCount)
}

func (a *petApp) settingsPageTitle() string {
	switch a.settingsTab {
	case tabMotion:
		return a.txt("motionPageTitle")
	case tabDisplay:
		return a.txt("displayPageTitle")
	default:
		return a.txt("animalPageTitle")
	}
}

func (a *petApp) settingsPageLead() string {
	switch a.settingsTab {
	case tabMotion:
		return a.txt("motionPageLead")
	case tabDisplay:
		return a.txt("displayPageLead")
	default:
		return a.txt("animalPageLead")
	}
}

func (a *petApp) displaySummary() string {
	areas := displayAreaForScope(a.displayScope)
	if len(areas) == 0 {
		return a.localText("現在の画面", "Current display")
	}
	scope, start, end := a.normalizedDisplaySelection(len(areas))
	if scope == displayScopeSpan {
		selected := areas[start : end+1]
		area := combineDisplayAreas(selected)
		width := max(0, int(area.Screen.Right-area.Screen.Left))
		height := max(0, int(area.Screen.Bottom-area.Screen.Top))
		if len(selected) == len(areas) {
			return fmt.Sprintf("%s / %d%s  %dx%d", a.localText("全モニタ", "All monitors"), len(selected), a.localText("画面", " displays"), width, height)
		}
		return fmt.Sprintf("%s %d-%d/%d  %dx%d", a.localText("複数画面", "Displays"), start+1, end+1, len(areas), width, height)
	}
	area := areas[start]
	width := max(0, int(area.Screen.Right-area.Screen.Left))
	height := max(0, int(area.Screen.Bottom-area.Screen.Top))
	primary := ""
	if area.Primary {
		primary = a.localText(" / メイン", " / primary")
	}
	return fmt.Sprintf("%s %d/%d%s  %dx%d", a.localText("ディスプレイ", "Display"), start+1, len(areas), primary, width, height)
}

func (a *petApp) walkRangeSummary() string {
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	if summary := a.walkRangeSummaryForSegments(start, end, a.displaySegmentsForSummary()); summary != "" {
		return summary
	}
	if start == 0 && end == 100 {
		return a.localText("全幅 (0%-100%)", "Full width (0%-100%)")
	}
	return fmt.Sprintf("%d%% - %d%%", start, end)
}

func (a *petApp) walkRangeSectionLabel() string {
	if normalizeDisplayScope(int(a.displayScope)) == displayScopeSpan {
		return a.localText("歩く画面", "Walking displays")
	}
	return a.localText("歩行範囲", "Walking range")
}

func (a *petApp) displaySegmentsForSummary() []sceneSegment {
	areas := displayAreaForScope(a.displayScope)
	if len(areas) == 0 {
		return nil
	}
	scope, start, end := a.normalizedDisplaySelection(len(areas))
	if scope != displayScopeSpan {
		end = start
	}
	selected := areas[start : end+1]
	combined := combineDisplayAreas(selected)
	base := combined.Work
	if normalizeOverlayPositionMode(int(a.positionMode)) == positionScreenBottom {
		base = combined.Screen
	}
	if base.Right <= base.Left {
		return nil
	}
	segments := make([]sceneSegment, 0, len(selected))
	for _, area := range selected {
		segmentBase := area.Work
		if normalizeOverlayPositionMode(int(a.positionMode)) == positionScreenBottom {
			segmentBase = area.Screen
		}
		left := max(int(segmentBase.Left), int(base.Left))
		right := min(int(segmentBase.Right), int(base.Right))
		if right-left >= spriteW {
			segments = append(segments, sceneSegment{
				Left:  left - int(base.Left),
				Right: right - int(base.Left),
			})
		}
	}
	return mergeSceneSegments(segments)
}

func (a *petApp) walkRangeSummaryForSegments(start, end int, segments []sceneSegment) string {
	start, end = normalizeWalkRange(start, end)
	segments = normalizeSceneSegments(segmentSpanWidth(segments), segments)
	if len(segments) == 0 {
		return ""
	}
	if start == 0 && end == 100 {
		if len(segments) == 1 {
			return a.localText("全幅", "Full width")
		}
		return a.localText("選択した画面ぜんぶ", "All selected displays")
	}
	totalLeft := segments[0].Left
	totalRight := segments[0].Right
	for _, segment := range segments[1:] {
		totalLeft = min(totalLeft, segment.Left)
		totalRight = max(totalRight, segment.Right)
	}
	totalW := max(1, totalRight-totalLeft)
	walkLeft := totalLeft + totalW*start/100
	walkRight := totalLeft + totalW*end/100
	covered := make([]int, 0, len(segments))
	fullCovered := true
	for i, segment := range segments {
		left := max(walkLeft, segment.Left)
		right := min(walkRight, segment.Right)
		if right-left <= 0 {
			continue
		}
		covered = append(covered, i+1)
		if left > segment.Left+2 || right < segment.Right-2 {
			fullCovered = false
		}
	}
	if len(covered) == 0 {
		return fmt.Sprintf("%d%% - %d%%", start, end)
	}
	if len(segments) == 1 {
		return a.localText("画面の一部", "Part of display")
	}
	if len(covered) == len(segments) && fullCovered {
		return a.localText("選択した画面ぜんぶ", "All selected displays")
	}
	first := covered[0]
	last := covered[len(covered)-1]
	if fullCovered {
		if first == last {
			return fmt.Sprintf("%s%d%s", a.localText("画面", "Display "), first, a.localText("だけ", " only"))
		}
		return fmt.Sprintf("%s%d-%d", a.localText("画面", "Displays "), first, last)
	}
	if first == last {
		if a.lang == langEnglish {
			return fmt.Sprintf("Part of display %d", first)
		}
		return fmt.Sprintf("画面%dの一部", first)
	}
	if a.lang == langEnglish {
		return fmt.Sprintf("Part of displays %d-%d", first, last)
	}
	return fmt.Sprintf("画面%d-%dの一部", first, last)
}

func segmentSpanWidth(segments []sceneSegment) int {
	if len(segments) == 0 {
		return 0
	}
	right := segments[0].Right
	for _, segment := range segments[1:] {
		right = max(right, segment.Right)
	}
	return right
}

func (a *petApp) positionSummary() string {
	mode := a.localText("タスクバー基準", "Taskbar edge")
	if normalizeOverlayPositionMode(int(a.positionMode)) == positionScreenBottom {
		mode = a.localText("画面下端基準", "Screen bottom")
	}
	offset := normalizeOverlayOffset(a.overlayOffsetY)
	if offset >= 0 {
		return fmt.Sprintf("%s / +%d px", mode, offset)
	}
	return fmt.Sprintf("%s / %d px", mode, offset)
}

func (a *petApp) settingsCoatSelectButton(id int32) bool {
	if id == ctrlVariantCombo {
		return a.coatMode == coatFixed
	}
	if id >= ctrlPetVariantBase && id < ctrlPetVariantBase+maxPetCount {
		index := int(id - ctrlPetVariantBase)
		return index >= 0 &&
			index < len(a.selectedSlotModes) &&
			normalizePetSlotMode(int(a.selectedSlotModes[index])) == petSlotFixed
	}
	return false
}

func settingsPetRowRects(index int) (win.RECT, win.RECT, win.RECT, win.RECT) {
	row := index
	if row < 0 {
		row = 0
	}
	left := int32(250)
	top := int32(372 + row*19)
	numberRect := win.RECT{Left: left, Top: top, Right: left + 24, Bottom: top + 18}
	nameRect := win.RECT{Left: left + 32, Top: top, Right: left + 104, Bottom: top + 18}
	variantRect := win.RECT{Left: left + 112, Top: top, Right: left + 344, Bottom: top + 18}
	sizeRect := win.RECT{Left: left + 350, Top: top, Right: left + 454, Bottom: top + 18}
	return numberRect, nameRect, variantRect, sizeRect
}

func settingsPetVariantRects(index int) (win.RECT, win.RECT) {
	numberRect, _, variantRect, _ := settingsPetRowRects(index)
	return numberRect, variantRect
}

func settingsPetNameRects(index int) (win.RECT, win.RECT) {
	numberRect, nameRect, _, _ := settingsPetRowRects(index)
	return numberRect, nameRect
}

func settingsPetSizeRect(index int) win.RECT {
	_, _, _, sizeRect := settingsPetRowRects(index)
	return sizeRect
}

func (a *petApp) settingsSelectVariant(id int32) int {
	if id >= ctrlPetVariantBase && id < ctrlPetVariantBase+maxPetCount {
		index := int(id - ctrlPetVariantBase)
		if index >= 0 && index < len(a.selectedCoats) {
			return clamp(a.selectedCoats[index], 0, len(variants)-1)
		}
	}
	return clamp(a.variant, 0, len(variants)-1)
}

func (a *petApp) petDisplayName(index int) string {
	if index < 0 || index >= maxPetCount {
		return ""
	}
	if name := sanitizePetName(a.petNames[index]); name != "" {
		return name
	}
	return a.localText(fmt.Sprintf("ペット%d", index+1), fmt.Sprintf("Pet %d", index+1))
}

func (a *petApp) petNameSectionLabel() string {
	if a.coatMode == coatSelected {
		return a.localText("名前 / 表示 / サイズ", "Names / assignment / size")
	}
	return a.localText("名前 / サイズ", "Names / size")
}

func (a *petApp) animalChoiceLabel() string {
	if a.coatMode == coatRandom {
		return a.localText("ランダム対象", "Random pool")
	}
	if a.coatMode == coatSelected {
		return a.localText("各ペットの表示", "Per-pet assignment")
	}
	return a.txt("coatColor")
}

func (a *petApp) petAssignmentDisplayLabel(index int) string {
	if index < 0 || index >= maxPetCount {
		return ""
	}
	if normalizePetSlotMode(int(a.selectedSlotModes[index])) == petSlotRandom {
		return a.randomGroupDisplayLabel(a.selectedRandomGroupIDs[index])
	}
	return a.variantDisplayLabel(a.selectedCoats[index])
}

func (a *petApp) randomGroupDisplayLabel(groupID string) string {
	groupID = normalizeVariantGroupID(groupID)
	if groupID == "" {
		return a.localText("ランダム / すべて", "Random / all")
	}
	return a.localText("ランダム / "+a.variantGroupLabelByID(groupID), "Random / "+a.variantGroupLabelByID(groupID))
}

func (a *petApp) variantGroupLabelByID(groupID string) string {
	group, ok := catalog.VariantGroupByID(groupID)
	if !ok {
		return ""
	}
	if a.lang == langEnglish {
		return group.LabelEN
	}
	return group.LabelJA
}

func (a *petApp) localText(ja, en string) string {
	if a.lang == langEnglish {
		return en
	}
	return ja
}

func sanitizePetName(name string) string {
	name = strings.Join(strings.Fields(name), " ")
	runes := []rune(name)
	if len(runes) > 24 {
		runes = runes[:24]
	}
	return string(runes)
}

type settingsRGB struct {
	r byte
	g byte
	b byte
}

func rgb(r, g, b byte) settingsRGB {
	return settingsRGB{r: r, g: g, b: b}
}

func drawRoundFill(hdc win.HDC, r win.RECT, c settingsRGB, radius int32) {
	brush := solidBrush(c.r, c.g, c.b)
	defer win.DeleteObject(win.HGDIOBJ(brush))
	oldBrush := win.SelectObject(hdc, win.HGDIOBJ(brush))
	oldPen := win.SelectObject(hdc, win.GetStockObject(win.NULL_PEN))
	win.RoundRect(hdc, r.Left, r.Top, r.Right, r.Bottom, radius, radius)
	win.SelectObject(hdc, oldPen)
	win.SelectObject(hdc, oldBrush)
}

func drawRoundOutline(hdc win.HDC, r win.RECT, c settingsRGB, radius int32) {
	lb := win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: win.RGB(c.r, c.g, c.b)}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, &lb, 0, nil)
	if pen == 0 {
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(pen))
	oldBrush := win.SelectObject(hdc, win.GetStockObject(win.NULL_BRUSH))
	oldPen := win.SelectObject(hdc, win.HGDIOBJ(pen))
	win.RoundRect(hdc, r.Left, r.Top, r.Right, r.Bottom, radius, radius)
	win.SelectObject(hdc, oldPen)
	win.SelectObject(hdc, oldBrush)
}

func drawRectFill(hdc win.HDC, r win.RECT, c settingsRGB) {
	brush := solidBrush(c.r, c.g, c.b)
	defer win.DeleteObject(win.HGDIOBJ(brush))
	oldBrush := win.SelectObject(hdc, win.HGDIOBJ(brush))
	oldPen := win.SelectObject(hdc, win.GetStockObject(win.NULL_PEN))
	win.Rectangle_(hdc, r.Left, r.Top, r.Right, r.Bottom)
	win.SelectObject(hdc, oldPen)
	win.SelectObject(hdc, oldBrush)
}

func drawTextLine(hdc win.HDC, text string, r win.RECT, font win.HFONT, c settingsRGB, flags uint32) {
	win.SetBkMode(hdc, win.TRANSPARENT)
	win.SetTextColor(hdc, win.RGB(c.r, c.g, c.b))
	oldFont := win.SelectObject(hdc, win.HGDIOBJ(font))
	chars := syscall.StringToUTF16(text)
	if len(chars) > 0 {
		win.DrawTextEx(hdc, &chars[0], int32(len(chars)-1), &r, flags, nil)
	}
	win.SelectObject(hdc, oldFont)
}

func drawVariantSwatch(hdc win.HDC, x, y int32, variant int, enabled bool) {
	base, patch, pied := variantSwatch(variant)
	if !enabled {
		base = rgb(168, 164, 151)
		patch = rgb(224, 220, 207)
	}
	brush := solidBrush(base.r, base.g, base.b)
	oldBrush := win.SelectObject(hdc, win.HGDIOBJ(brush))
	oldPen := win.SelectObject(hdc, win.GetStockObject(win.NULL_PEN))
	win.Ellipse(hdc, x, y, x+16, y+16)
	win.SelectObject(hdc, oldBrush)
	win.DeleteObject(win.HGDIOBJ(brush))
	if pied {
		patchBrush := solidBrush(patch.r, patch.g, patch.b)
		oldBrush = win.SelectObject(hdc, win.HGDIOBJ(patchBrush))
		win.Ellipse(hdc, x+7, y+3, x+15, y+12)
		win.SelectObject(hdc, oldBrush)
		win.DeleteObject(win.HGDIOBJ(patchBrush))
	}
	win.SelectObject(hdc, oldPen)
}

func variantSwatch(index int) (settingsRGB, settingsRGB, bool) {
	if index < 0 || index >= len(variants) {
		return rgb(128, 120, 105), rgb(240, 235, 220), false
	}
	if base, ok := swatchHex(variants[index].TintHex); ok {
		if patch, patchOK := swatchHex(variants[index].AccentHex); patchOK {
			return base, patch, true
		}
		return base, rgb(240, 235, 220), false
	}
	switch variants[index].ID {
	case "black":
		return rgb(42, 37, 33), rgb(240, 235, 220), false
	case "blue":
		return rgb(111, 116, 111), rgb(240, 235, 220), false
	case "gray":
		return rgb(128, 128, 120), rgb(240, 235, 220), false
	case "white_cream":
		return rgb(228, 218, 190), rgb(240, 235, 220), false
	case "sand_champagne":
		return rgb(185, 158, 115), rgb(240, 235, 220), false
	case "chocolate":
		return rgb(102, 69, 48), rgb(240, 235, 220), false
	case "black_pied":
		return rgb(42, 37, 33), rgb(236, 228, 204), true
	case "agouti_pied":
		return rgb(118, 96, 67), rgb(236, 228, 204), true
	case "blue_pied":
		return rgb(111, 116, 111), rgb(236, 228, 204), true
	case "cream_pied":
		return rgb(202, 175, 126), rgb(244, 238, 218), true
	case "chinchilla_standard_gray":
		return rgb(135, 139, 136), rgb(232, 232, 224), false
	case "macaroni_mouse_tan":
		return rgb(184, 145, 96), rgb(239, 224, 198), false
	case "rabbit_chestnut_agouti":
		return rgb(126, 91, 62), rgb(226, 205, 178), false
	case "dog_cream_tan":
		return rgb(213, 186, 133), rgb(246, 233, 203), false
	case "cat_brown_tabby":
		return rgb(122, 91, 58), rgb(216, 196, 160), true
	case "gecko_gray_brown":
		return rgb(117, 110, 92), rgb(213, 210, 190), false
	case "hamster_golden_syrian":
		return rgb(207, 139, 69), rgb(244, 223, 178), false
	default:
		return rgb(118, 96, 67), rgb(240, 235, 220), false
	}
}

func swatchHex(hex string) (settingsRGB, bool) {
	if len(hex) != 6 {
		return settingsRGB{}, false
	}
	val := func(i int) (byte, bool) {
		var out byte
		for _, ch := range hex[i : i+2] {
			out <<= 4
			switch {
			case ch >= '0' && ch <= '9':
				out += byte(ch - '0')
			case ch >= 'a' && ch <= 'f':
				out += byte(ch-'a') + 10
			case ch >= 'A' && ch <= 'F':
				out += byte(ch-'A') + 10
			default:
				return 0, false
			}
		}
		return out, true
	}
	r, ok := val(0)
	if !ok {
		return settingsRGB{}, false
	}
	g, ok := val(2)
	if !ok {
		return settingsRGB{}, false
	}
	b, ok := val(4)
	if !ok {
		return settingsRGB{}, false
	}
	return rgb(r, g, b), true
}

func (a *petApp) setControlFont(hwnd win.HWND, font win.HFONT) {
	if hwnd == 0 || font == 0 {
		return
	}
	win.SendMessage(hwnd, win.WM_SETFONT, uintptr(font), 1)
}

func (a *petApp) settingsWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		id := int32(uint16(wParam & 0xffff))
		notify := uint16((wParam >> 16) & 0xffff)
		if a.handleSettingsCommand(id, notify) {
			return 0
		}
	case win.WM_HSCROLL:
		if a.handleRangeScroll(getDlgCtrlID(win.HWND(lParam)), wParam) {
			return 0
		}
	case win.WM_PAINT:
		a.paintSettingsWindow(hwnd)
		return 0
	case win.WM_DRAWITEM:
		dis := drawItemStruct(lParam)
		if a.drawSettingsButton(&dis) {
			return 1
		}
		return 0
	case win.WM_ERASEBKGND:
		return 1
	case win.WM_CTLCOLORSTATIC:
		win.SetBkMode(win.HDC(wParam), win.TRANSPARENT)
		win.SetTextColor(win.HDC(wParam), win.RGB(32, 37, 31))
		return uintptr(win.GetStockObject(win.NULL_BRUSH))
	case win.WM_NCHITTEST:
		if a.settingsDragHit(lParam) {
			return uintptr(win.HTCAPTION)
		}
		return uintptr(win.HTCLIENT)
	case win.WM_CLOSE:
		a.rememberSettingsWindowPosition()
		a.persistSettings()
		win.ShowWindow(hwnd, win.SW_HIDE)
		return 0
	case win.WM_DESTROY:
		if hwnd == a.settingsHwnd {
			a.settingsHwnd = 0
		}
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (a *petApp) settingsDragHit(lParam uintptr) bool {
	var rect win.RECT
	if a.settingsHwnd == 0 || !win.GetWindowRect(a.settingsHwnd, &rect) {
		return false
	}
	x := int32(int16(lParam & 0xffff))
	y := int32(int16((lParam >> 16) & 0xffff))
	if y < rect.Top || y > rect.Top+88 {
		return false
	}
	if x >= rect.Right-64 && y <= rect.Top+58 {
		return false
	}
	return true
}

func drawItemStruct(lParam uintptr) win.DRAWITEMSTRUCT {
	var dis win.DRAWITEMSTRUCT
	if lParam != 0 {
		procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&dis)), lParam, unsafe.Sizeof(dis))
	}
	return dis
}

func (a *petApp) syncSettingsWindow() {
	if a.settingsHwnd == 0 {
		return
	}
	setWindowText(a.settingsHwnd, a.txt("settingsTitle"))
	a.setButtonChecked(ctrlTabAnimals, a.settingsTab == tabAnimals)
	a.setButtonChecked(ctrlTabMotion, a.settingsTab == tabMotion)
	a.setButtonChecked(ctrlTabDisplay, a.settingsTab == tabDisplay)
	a.setButtonChecked(ctrlModeKeyboard, a.mode == modeKeyboard)
	a.setButtonChecked(ctrlModeRandom, a.mode == modeRandom)
	a.setButtonChecked(ctrlSpeedSlow, a.speed == 2)
	a.setButtonChecked(ctrlSpeedNormal, a.speed == 3)
	a.setButtonChecked(ctrlSpeedFast, a.speed == 5)
	a.setButtonChecked(ctrlTypingWheel, a.wheelEnabled)
	a.setButtonChecked(ctrlBidirectional, a.bidirectional)
	a.setButtonChecked(ctrlPositionTaskbar, a.positionMode == positionTaskbarEdge)
	a.setButtonChecked(ctrlPositionBottom, a.positionMode == positionScreenBottom)
	a.setButtonChecked(ctrlDisplaySingle, a.displayScope == displayScopeSingle)
	a.setButtonChecked(ctrlDisplaySpan, a.displayScope == displayScopeSpan)
	a.setButtonChecked(ctrlNameLabels, a.nameLabels)
	a.setButtonChecked(ctrlCoatFixed, a.coatMode == coatFixed)
	a.setButtonChecked(ctrlCoatSelected, a.coatMode == coatSelected)
	a.setButtonChecked(ctrlCoatRandom, a.coatMode == coatRandom)
	a.syncSelectButton(ctrlVariantCombo)
	if h := win.GetDlgItem(a.settingsHwnd, ctrlVariantCombo); h != 0 {
		win.EnableWindow(h, a.coatMode == coatFixed || a.coatMode == coatRandom)
	}
	for i := 0; i < a.petCount; i++ {
		a.syncSelectButton(ctrlPetNameBase + int32(i))
		a.syncSelectButton(ctrlPetVariantBase + int32(i))
		a.syncSelectButton(ctrlPetSizeBase + int32(i))
	}
	a.syncSelectButton(ctrlLanguageCombo)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlPetMinus), a.petCount > 1)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlPetPlus), a.petCount < maxPetCount)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlOffsetUp), a.overlayOffsetY > minOverlayOffsetY)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlOffsetDown), a.overlayOffsetY < maxOverlayOffsetY)
	a.syncRangeScrollBars()
	monitorCount := len(monitorAreas())
	scope, spanStart, spanEnd := a.normalizedDisplaySelection(monitorCount)
	canMovePrev := monitorCount > 1
	canMoveNext := monitorCount > 1
	if scope == displayScopeSpan {
		canMovePrev = spanStart > 0
		canMoveNext = spanEnd < monitorCount-1
	}
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlDisplaySingle), monitorCount > 0)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlDisplaySpan), monitorCount > 1)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlDisplaySpanLess), scope == displayScopeSpan && spanEnd-spanStart > 1)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlDisplaySpanMore), scope == displayScopeSpan && monitorCount > 1 && spanEnd-spanStart+1 < monitorCount)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlDisplayPrev), canMovePrev)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlDisplayNext), canMoveNext)
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlRangeNarrow), end-start > minWalkRangeSpan)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlRangeWide), start > 0 || end < 100)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlRangeLeft), start > 0)
	win.EnableWindow(win.GetDlgItem(a.settingsHwnd, ctrlRangeRight), end < 100)
	win.InvalidateRect(a.settingsHwnd, nil, true)
}

func (a *petApp) setButtonChecked(id int32, checked bool) {
	h := win.GetDlgItem(a.settingsHwnd, id)
	if h == 0 {
		return
	}
	value := uintptr(win.BST_UNCHECKED)
	if checked {
		value = uintptr(win.BST_CHECKED)
	}
	win.SendMessage(h, win.BM_SETCHECK, value, 0)
	win.InvalidateRect(h, nil, true)
}

func (a *petApp) syncSelectButton(id int32) {
	h := win.GetDlgItem(a.settingsHwnd, id)
	if h == 0 {
		return
	}
	setWindowText(h, a.settingsButtonLabel(id))
	win.InvalidateRect(h, nil, true)
}

func (a *petApp) syncRangeScrollBars() {
	a.setScrollRangeAndPos(ctrlRangeStartScroll, a.walkRangeScrollValue(ctrlRangeStartScroll))
	a.setScrollRangeAndPos(ctrlRangeEndScroll, a.walkRangeScrollValue(ctrlRangeEndScroll))
}

func (a *petApp) setScrollRangeAndPos(id int32, pos int) {
	if a.settingsHwnd == 0 {
		return
	}
	h := win.GetDlgItem(a.settingsHwnd, id)
	if h == 0 {
		return
	}
	pos = clamp(pos, 0, 100)
	win.SendMessage(h, sbmSetRange, 0, 100)
	win.SendMessage(h, sbmSetPos, uintptr(pos), 1)
	win.InvalidateRect(h, nil, true)
}

func (a *petApp) walkRangeScrollValue(id int32) int {
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	if id == ctrlRangeEndScroll {
		return end
	}
	return start
}

func (a *petApp) handleRangeScroll(id int32, wParam uintptr) bool {
	if id != ctrlRangeStartScroll && id != ctrlRangeEndScroll {
		return false
	}
	code := int(uint16(wParam & 0xffff))
	thumb := int(uint16((wParam >> 16) & 0xffff))
	start, end := normalizeWalkRange(a.walkRangeStart, a.walkRangeEnd)
	pos := start
	if id == ctrlRangeEndScroll {
		pos = end
	}
	switch code {
	case sbLineLeft:
		pos--
	case sbLineRight:
		pos++
	case sbPageLeft:
		pos -= walkRangeStep
	case sbPageRight:
		pos += walkRangeStep
	case sbThumbPosition, sbThumbTrack:
		pos = thumb
	case sbLeft:
		pos = 0
	case sbRight:
		pos = 100
	case sbEndScroll:
		return true
	default:
		return true
	}
	if id == ctrlRangeStartScroll {
		start = clamp(pos, 0, end-minWalkRangeSpan)
	} else {
		end = clamp(pos, start+minWalkRangeSpan, 100)
	}
	a.setWalkRange(start, end)
	a.syncSettingsWindow()
	a.persistSettings()
	a.render()
	return true
}

func (a *petApp) handleSettingsCommand(id int32, notify uint16) bool {
	if id == ctrlNameLabels {
		a.nameLabels = !a.nameLabels
		if !a.nameLabels {
			a.hideNameWindow()
		}
		a.recreateSettingsWindow()
		a.persistSettings()
		a.render()
		return true
	}
	if id >= ctrlPetNameBase && id < ctrlPetNameBase+maxPetCount {
		a.showRenameDialog(int(id - ctrlPetNameBase))
		return true
	}
	if id >= ctrlPetSizeBase && id < ctrlPetSizeBase+maxPetCount {
		index := int(id - ctrlPetSizeBase)
		sel, ok := a.pickPetSizeFromMenu(id, a.petSizePercent(index))
		if ok {
			a.setPetSize(index, sel)
		}
		a.syncSettingsWindow()
		a.persistSettings()
		a.render()
		return true
	}
	switch id {
	case ctrlTabAnimals:
		a.settingsTab = tabAnimals
		a.recreateSettingsWindow()
	case ctrlTabMotion:
		a.settingsTab = tabMotion
		a.recreateSettingsWindow()
	case ctrlTabDisplay:
		a.settingsTab = tabDisplay
		a.recreateSettingsWindow()
	case ctrlVariantCombo:
		if a.coatMode == coatRandom {
			groupID, ok := a.pickRandomGroupFromMenu(id, a.randomGroupID)
			if ok {
				a.setRandomGroupID(groupID)
			}
		} else {
			sel, ok := a.pickVariantFromMenu(id, a.variant)
			if ok {
				a.setFixedVariant(sel)
			}
		}
	case ctrlCoatFixed:
		a.setCoatMode(coatFixed)
		a.recreateSettingsWindow()
	case ctrlCoatSelected:
		a.setCoatMode(coatSelected)
		a.recreateSettingsWindow()
	case ctrlCoatRandom:
		a.setCoatMode(coatRandom)
		a.recreateSettingsWindow()
	case ctrlLanguageCombo:
		if lang, ok := a.pickLanguageFromMenu(id); ok {
			a.lang = lang
			a.recreateSettingsWindow()
		}
	case ctrlPetMinus:
		a.setPetCount(a.petCount - 1)
		a.resetPosition()
		a.recreateSettingsWindow()
	case ctrlPetPlus:
		a.setPetCount(a.petCount + 1)
		a.resetPosition()
		a.recreateSettingsWindow()
	case ctrlModeKeyboard:
		a.handleMenu(menuModeKeyboard)
	case ctrlModeRandom:
		a.handleMenu(menuModeRandom)
	case ctrlSpeedSlow:
		a.handleMenu(menuSpeedSlow)
	case ctrlSpeedNormal:
		a.handleMenu(menuSpeedNormal)
	case ctrlSpeedFast:
		a.handleMenu(menuSpeedFast)
	case ctrlTypingWheel:
		a.handleMenu(menuWheelToggle)
	case ctrlBidirectional:
		a.setBidirectional(!a.bidirectional)
	case ctrlPositionTaskbar:
		a.positionMode = positionTaskbarEdge
	case ctrlPositionBottom:
		a.positionMode = positionScreenBottom
	case ctrlOffsetUp:
		a.adjustOverlayOffset(-overlayOffsetStep)
	case ctrlOffsetDown:
		a.adjustOverlayOffset(overlayOffsetStep)
	case ctrlDisplaySingle:
		a.setDisplayScope(displayScopeSingle)
	case ctrlDisplaySpan:
		a.setDisplayScope(displayScopeSpan)
	case ctrlDisplaySpanLess:
		a.adjustDisplaySpan(-1)
	case ctrlDisplaySpanMore:
		a.adjustDisplaySpan(1)
	case ctrlDisplayPrev:
		a.adjustDisplayIndex(-1)
	case ctrlDisplayNext:
		a.adjustDisplayIndex(1)
	case ctrlRangeFull:
		a.setWalkRange(defaultWalkRangeStart, defaultWalkRangeEnd)
	case ctrlRangeNarrow:
		a.adjustWalkRangeWidth(-walkRangeStep)
	case ctrlRangeWide:
		a.adjustWalkRangeWidth(walkRangeStep)
	case ctrlRangeLeft:
		a.shiftWalkRange(-walkRangeStep)
	case ctrlRangeRight:
		a.shiftWalkRange(walkRangeStep)
	case ctrlReset:
		a.resetOverlayPlacement()
		a.resetPosition()
		a.render()
	case ctrlClose, ctrlTopClose, int32(win.IDCANCEL):
		if a.settingsHwnd != 0 {
			a.rememberSettingsWindowPosition()
			win.ShowWindow(a.settingsHwnd, win.SW_HIDE)
		}
	default:
		if id >= ctrlPetVariantBase && id < ctrlPetVariantBase+maxPetCount {
			assignment, ok := a.pickPetAssignmentFromMenu(id, int(id-ctrlPetVariantBase))
			if ok {
				if assignment.Mode == petSlotRandom {
					a.setSelectedRandomGroup(int(id-ctrlPetVariantBase), assignment.GroupID)
				} else {
					a.setSelectedVariant(int(id-ctrlPetVariantBase), assignment.Variant)
				}
			}
			break
		}
		return false
	}
	a.syncSettingsWindow()
	a.persistSettings()
	a.render()
	return true
}

func (a *petApp) pickVariantFromMenu(id int32, selected int) (int, bool) {
	menu := win.CreatePopupMenu()
	a.appendGroupedVariantMenu(menu, selected, func(index int) uintptr {
		return uintptr(index + 1)
	})
	cmd := a.trackControlMenu(id, menu)
	win.DestroyMenu(menu)
	if cmd == 0 {
		return 0, false
	}
	choice := int(cmd) - 1
	if choice < 0 || choice >= len(variants) {
		return 0, false
	}
	return choice, true
}

type petAssignmentChoice struct {
	Mode    petSlotMode
	Variant int
	GroupID string
}

func (a *petApp) pickPetAssignmentFromMenu(id int32, index int) (petAssignmentChoice, bool) {
	menu := win.CreatePopupMenu()
	const (
		cmdRandomAll       uint32 = 1
		cmdRandomGroupBase uint32 = 100
		cmdVariantBase     uint32 = 1000
	)
	currentMode := petSlotFixed
	currentVariant := 0
	currentGroupID := ""
	if index >= 0 && index < maxPetCount {
		currentMode = normalizePetSlotMode(int(a.selectedSlotModes[index]))
		currentVariant = clamp(a.selectedCoats[index], 0, len(variants)-1)
		currentGroupID = normalizeVariantGroupID(a.selectedRandomGroupIDs[index])
	}
	appendMenu(menu, win.MF_STRING|checkedFlag(currentMode == petSlotRandom && currentGroupID == ""), uintptr(cmdRandomAll), syscall.StringToUTF16Ptr(a.randomAllMenuLabel()))
	groupMenu := win.CreatePopupMenu()
	groupCommandIDs := make(map[uint32]string)
	for groupIndex, group := range a.variantMenuGroups() {
		cmd := cmdRandomGroupBase + uint32(groupIndex)
		groupCommandIDs[cmd] = group.ID
		flags := uint32(win.MF_STRING)
		if currentMode == petSlotRandom && currentGroupID == group.ID {
			flags |= win.MF_CHECKED
		}
		appendMenu(groupMenu, flags, uintptr(cmd), syscall.StringToUTF16Ptr(a.randomGroupMenuLabel(group.Label)))
	}
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(groupMenu), syscall.StringToUTF16Ptr(a.txt("randomByType")))
	appendMenu(menu, win.MF_SEPARATOR, 0, nil)
	selectedVariant := currentVariant
	if currentMode == petSlotRandom {
		selectedVariant = -1
	}
	a.appendGroupedVariantMenu(menu, selectedVariant, func(variantIndex int) uintptr {
		return uintptr(cmdVariantBase + uint32(variantIndex))
	})
	cmd := a.trackControlMenu(id, menu)
	win.DestroyMenu(menu)
	if cmd == 0 {
		return petAssignmentChoice{}, false
	}
	if cmd == cmdRandomAll {
		return petAssignmentChoice{Mode: petSlotRandom}, true
	}
	if groupID, ok := groupCommandIDs[cmd]; ok {
		return petAssignmentChoice{Mode: petSlotRandom, GroupID: groupID}, true
	}
	if cmd >= cmdVariantBase && int(cmd-cmdVariantBase) < len(variants) {
		return petAssignmentChoice{Mode: petSlotFixed, Variant: int(cmd - cmdVariantBase)}, true
	}
	return petAssignmentChoice{}, false
}

func (a *petApp) pickRandomGroupFromMenu(id int32, selected string) (string, bool) {
	menu := win.CreatePopupMenu()
	const cmdRandomAll uint32 = 1
	selected = normalizeVariantGroupID(selected)
	appendMenu(menu, win.MF_STRING|checkedFlag(selected == ""), uintptr(cmdRandomAll), syscall.StringToUTF16Ptr(a.randomAllMenuLabel()))
	groupCommandIDs := make(map[uint32]string)
	for groupIndex, group := range a.variantMenuGroups() {
		cmd := uint32(groupIndex + 2)
		groupCommandIDs[cmd] = group.ID
		flags := uint32(win.MF_STRING)
		if selected == group.ID {
			flags |= win.MF_CHECKED
		}
		appendMenu(menu, flags, uintptr(cmd), syscall.StringToUTF16Ptr(group.Label))
	}
	cmd := a.trackControlMenu(id, menu)
	win.DestroyMenu(menu)
	if cmd == 0 {
		return "", false
	}
	if cmd == cmdRandomAll {
		return "", true
	}
	if groupID, ok := groupCommandIDs[cmd]; ok {
		return groupID, true
	}
	return "", false
}

func (a *petApp) randomAllMenuLabel() string {
	return a.localText("ランダム（すべて）", "Random: all")
}

func (a *petApp) randomGroupMenuLabel(group string) string {
	if a.lang == langEnglish {
		return "Random: " + group
	}
	return "ランダム（" + group + "）"
}

func (a *petApp) pickPetSizeFromMenu(id int32, selected int) (int, bool) {
	menu := win.CreatePopupMenu()
	choiceID := uintptr(1)
	for size := minPetSizePercent; size <= maxPetSizePercent; size += petSizeStepPercent {
		flags := uint32(win.MF_STRING)
		if size == normalizePetSizePercent(selected) {
			flags |= win.MF_CHECKED
		}
		appendMenu(menu, flags, choiceID, syscall.StringToUTF16Ptr(fmt.Sprintf("%d%%", size)))
		choiceID++
	}
	cmd := a.trackControlMenu(id, menu)
	win.DestroyMenu(menu)
	if cmd == 0 {
		return 0, false
	}
	size := minPetSizePercent + (int(cmd)-1)*petSizeStepPercent
	if size < minPetSizePercent || size > maxPetSizePercent {
		return 0, false
	}
	return size, true
}

func (a *petApp) pickLanguageFromMenu(id int32) (language, bool) {
	menu := win.CreatePopupMenu()
	appendMenu(menu, win.MF_STRING|checkedFlag(a.lang == langJapanese), 1, syscall.StringToUTF16Ptr("日本語"))
	appendMenu(menu, win.MF_STRING|checkedFlag(a.lang == langEnglish), 2, syscall.StringToUTF16Ptr("English"))
	cmd := a.trackControlMenu(id, menu)
	win.DestroyMenu(menu)
	switch cmd {
	case 1:
		return langJapanese, true
	case 2:
		return langEnglish, true
	}
	return a.lang, false
}

func checkedFlag(checked bool) uint32 {
	if checked {
		return win.MF_CHECKED
	}
	return 0
}

func (a *petApp) trackControlMenu(id int32, menu win.HMENU) uint32 {
	h := win.GetDlgItem(a.settingsHwnd, id)
	if h == 0 {
		return 0
	}
	var rect win.RECT
	if !win.GetWindowRect(h, &rect) {
		return 0
	}
	win.SetForegroundWindow(a.settingsHwnd)
	return win.TrackPopupMenu(menu, win.TPM_RETURNCMD|win.TPM_LEFTALIGN|win.TPM_TOPALIGN|win.TPM_RIGHTBUTTON, rect.Left, rect.Bottom+4, 0, a.settingsHwnd, nil)
}

func (a *petApp) showRenameDialog(index int) {
	if index < 0 || index >= maxPetCount || a.settingsHwnd == 0 {
		return
	}
	if a.renameHwnd != 0 {
		win.DestroyWindow(a.renameHwnd)
		a.renameHwnd = 0
	}
	a.ensureSettingsBrushes()
	a.ensureSettingsFonts()
	var parentRect win.RECT
	win.GetWindowRect(a.settingsHwnd, &parentRect)
	style := renameDialogStyle()
	exStyle := uint32(win.WS_EX_TOOLWINDOW)
	w, h := windowSizeForClient(renameDialogClientW, renameDialogClientH, style, exStyle)
	x := parentRect.Left + (parentRect.Right-parentRect.Left-w)/2
	y := parentRect.Top + (parentRect.Bottom-parentRect.Top-h)/2
	hwnd := win.CreateWindowEx(
		exStyle,
		syscall.StringToUTF16Ptr(windowClass),
		syscall.StringToUTF16Ptr(a.localText("名前を変更", "Rename pet")),
		style,
		x, y, w, h,
		a.settingsHwnd, 0, a.hinst, nil,
	)
	if hwnd == 0 {
		return
	}
	a.renameHwnd = hwnd
	a.renameIndex = index
	editRect, okRect, cancelRect := renameDialogLayoutRects()
	a.renameEdit = a.createEdit(hwnd, ctrlRenameEdit, a.petDisplayName(index), editRect.Left, editRect.Top, rectWidth(editRect), rectHeight(editRect))
	a.createButton(hwnd, ctrlRenameOK, "", okRect.Left, okRect.Top, rectWidth(okRect), rectHeight(okRect), 0)
	a.createButton(hwnd, ctrlRenameCancel, "", cancelRect.Left, cancelRect.Top, rectWidth(cancelRect), rectHeight(cancelRect), 0)
	win.SetForegroundWindow(hwnd)
	win.SetFocus(a.renameEdit)
}

func renameDialogStyle() uint32 {
	return win.WS_POPUP | win.WS_CAPTION | win.WS_SYSMENU | win.WS_VISIBLE | win.WS_CLIPCHILDREN
}

func renameDialogLayoutRects() (editRect, okRect, cancelRect win.RECT) {
	return win.RECT{Left: 24, Top: 56, Right: 336, Bottom: 84},
		win.RECT{Left: 144, Top: 104, Right: 222, Bottom: 134},
		win.RECT{Left: 234, Top: 104, Right: 336, Bottom: 134}
}

func windowSizeForClient(clientW, clientH int32, style, exStyle uint32) (int32, int32) {
	rect := win.RECT{Left: 0, Top: 0, Right: clientW, Bottom: clientH}
	if adjustWindowRectEx(&rect, style, false, exStyle) {
		return rect.Right - rect.Left, rect.Bottom - rect.Top
	}
	return clientW, clientH
}

func adjustWindowRectEx(rect *win.RECT, style uint32, menu bool, exStyle uint32) bool {
	if rect == nil {
		return false
	}
	hasMenu := uintptr(0)
	if menu {
		hasMenu = 1
	}
	ret, _, _ := procAdjustWindowRectEx.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(style),
		hasMenu,
		uintptr(exStyle),
	)
	return ret != 0
}

func rectWidth(rect win.RECT) int32 {
	return rect.Right - rect.Left
}

func rectHeight(rect win.RECT) int32 {
	return rect.Bottom - rect.Top
}

func (a *petApp) renameWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		id := int32(uint16(wParam & 0xffff))
		switch id {
		case ctrlRenameOK:
			a.commitRenameDialog()
			return 0
		case ctrlRenameCancel, int32(win.IDCANCEL):
			win.DestroyWindow(hwnd)
			return 0
		}
	case win.WM_PAINT:
		a.paintRenameWindow(hwnd)
		return 0
	case win.WM_DRAWITEM:
		dis := drawItemStruct(lParam)
		if a.drawSettingsButton(&dis) {
			return 1
		}
		return 0
	case win.WM_CTLCOLORSTATIC:
		win.SetBkMode(win.HDC(wParam), win.TRANSPARENT)
		win.SetTextColor(win.HDC(wParam), win.RGB(32, 37, 31))
		return uintptr(win.GetStockObject(win.NULL_BRUSH))
	case win.WM_CLOSE:
		win.DestroyWindow(hwnd)
		return 0
	case win.WM_DESTROY:
		if hwnd == a.renameHwnd {
			a.renameHwnd = 0
			a.renameEdit = 0
			a.renameIndex = -1
		}
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func (a *petApp) paintRenameWindow(hwnd win.HWND) {
	a.ensureSettingsFonts()
	var ps win.PAINTSTRUCT
	hdc := win.BeginPaint(hwnd, &ps)
	if hdc == 0 {
		return
	}
	defer win.EndPaint(hwnd, &ps)
	var rect win.RECT
	win.GetClientRect(hwnd, &rect)
	drawRectFill(hdc, rect, rgb(247, 248, 244))
	label := a.localText(fmt.Sprintf("%d匹目の名前", a.renameIndex+1), fmt.Sprintf("Pet %d name", a.renameIndex+1))
	drawTextLine(hdc, label, win.RECT{Left: 24, Top: 20, Right: rect.Right - 24, Bottom: 46}, a.settingsFont, rgb(27, 36, 32), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE|win.DT_NOPREFIX)
}

func (a *petApp) commitRenameDialog() {
	if a.renameIndex >= 0 && a.renameIndex < maxPetCount && a.renameEdit != 0 {
		a.petNames[a.renameIndex] = sanitizePetName(getWindowText(a.renameEdit))
		a.syncSettingsWindow()
		a.persistSettings()
		a.updateHoverName()
		a.render()
	}
	if a.renameHwnd != 0 {
		win.DestroyWindow(a.renameHwnd)
	}
}

func (a *petApp) recreateSettingsWindow() {
	if a.settingsHwnd != 0 {
		a.rememberSettingsWindowPosition()
		win.DestroyWindow(a.settingsHwnd)
		a.settingsHwnd = 0
	}
	a.showSettings()
}

func (a *petApp) rememberSettingsWindowPosition() {
	var rect win.RECT
	if a.settingsHwnd == 0 || !win.GetWindowRect(a.settingsHwnd, &rect) {
		return
	}
	a.settingsX = rect.Left
	a.settingsY = rect.Top
	a.clampSettingsWindowPosition()
}

func (a *petApp) clampSettingsWindowPosition() {
	work := workArea()
	maxX := work.Right - settingsClientW
	maxY := work.Bottom - settingsClientH
	if maxX < work.Left {
		maxX = work.Left
	}
	if maxY < work.Top {
		maxY = work.Top
	}
	a.settingsX = int32(clamp(int(a.settingsX), int(work.Left), int(maxX)))
	a.settingsY = int32(clamp(int(a.settingsY), int(work.Top), int(maxY)))
}

func (a *petApp) txt(key string) string {
	switch key {
	case "typingWheel":
		return a.localText("チンチラ/ハムスター回し車", "Chinchilla/hamster wheel")
	case "displayPageTitle":
		return a.localText("画面と歩行範囲", "Screen and walking area")
	case "displayPageLead":
		return a.localText("表示する画面、歩く範囲、縦位置を調整します。", "Choose which display area the pets use, their walking range, and their vertical placement.")
	case "tabDisplay":
		return a.localText("画面", "Display")
	}
	if a.lang == langEnglish {
		switch key {
		case "settingsTitle":
			return "Animals Desktop Settings"
		case "settingsHeader":
			return "Animals Desktop"
		case "settingsLead":
			return "Taskbar companion controls"
		case "animalPageTitle":
			return "Animals and size"
		case "animalPageLead":
			return "Set visible pets, animal assignment, names, and per-pet size."
		case "motionPageTitle":
			return "Motion behavior"
		case "motionPageLead":
			return "Tune keyboard reactions, random strolls, the chinchilla/hamster wheel, and turn behavior."
		case "displayPageTitle":
			return "Screen and walking area"
		case "displayPageLead":
			return "Choose which display area the pets use, their walking range, and their vertical placement."
		case "tabAnimals":
			return "Animals"
		case "tabMotion":
			return "Motion"
		case "tabDisplay":
			return "Display"
		case "animalSection":
			return "Animals"
		case "petCount":
			return "Visible pets"
		case "coatColor":
			return "Fixed animal"
		case "coatMode":
			return "Animal selection"
		case "coatFixed":
			return "Fixed"
		case "coatSelected":
			return "Choose each"
		case "coatRandom":
			return "Random"
		case "randomByType":
			return "Random by type"
		case "randomTarget":
			return "Random pool"
		case "selectedCoats":
			return "Per-pet animals"
		case "coatNote":
			return "Random gives each pet its own animal or color. Seed-stage species use source-truth sheets until full motion is drawn."
		case "language":
			return "Language"
		case "mode":
			return "Mode"
		case "modeKeyboard":
			return "Keyboard reaction"
		case "modeRandom":
			return "Random stroll"
		case "speed":
			return "Speed"
		case "speedSlow":
			return "Slow"
		case "speedNormal":
			return "Normal"
		case "speedFast":
			return "Fast"
		case "motion":
			return "Motion"
		case "typingWheel":
			return "Chinchilla/hamster wheel"
		case "naturalTurns":
			return "Natural left/right turns"
		case "reset":
			return "Reset"
		case "close":
			return "Close"
		case "exit":
			return "Exit"
		}
	}
	switch key {
	case "settingsTitle":
		return "AnimalsDesktop設定"
	case "settingsHeader":
		return "AnimalsDesktop"
	case "settingsLead":
		return "タスクバーで遊ぶペットの設定"
	case "animalPageTitle":
		return "ペットの数とサイズ"
	case "animalPageLead":
		return "ペットの数、種類の選び方、名前、1匹ごとのサイズを調整します。"
	case "motionPageTitle":
		return "動きかた"
	case "motionPageLead":
		return "キーボードへの反応、ランダム散歩、左右ターンを調整します。"
	case "tabAnimals":
		return "動物"
	case "tabMotion":
		return "動き"
	case "animalSection":
		return "ペット"
	case "petCount":
		return "出現数"
	case "coatColor":
		return "決まったペット"
	case "coatMode":
		return "ペットの選び方"
	case "coatFixed":
		return "固定"
	case "coatSelected":
		return "1匹ずつ選ぶ"
	case "coatRandom":
		return "ランダム"
	case "randomByType":
		return "種類でランダム"
	case "randomTarget":
		return "ランダム対象"
	case "selectedCoats":
		return "それぞれのペット"
	case "coatNote":
		return "「ランダム」では、完成済みのペットから1匹ごとに選ばれます。正式リリースまでは完成したペットだけを表示します。"
	case "language":
		return "Language"
	case "mode":
		return "モード"
	case "modeKeyboard":
		return "キーボード反応"
	case "modeRandom":
		return "ランダム散歩"
	case "speed":
		return "速度"
	case "speedSlow":
		return "ゆっくり"
	case "speedNormal":
		return "ふつう"
	case "speedFast":
		return "はやい"
	case "motion":
		return "動作"
	case "typingWheel":
		return "入力中だけ回し車"
	case "naturalTurns":
		return "自然な左右ターン"
	case "reset":
		return "整列"
	case "close":
		return "閉じる"
	case "exit":
		return "終了"
	}
	return key
}

func (a *petApp) variantLabel(i int) string {
	if i < 0 || i >= len(variants) {
		return ""
	}
	if a.lang == langEnglish {
		return variants[i].LabelEN
	}
	return variants[i].LabelJA
}

func (a *petApp) variantDisplayLabel(i int) string {
	label := a.variantLabel(i)
	if label == "" {
		return ""
	}
	group := a.variantGroupLabel(i)
	if group == "" {
		return label
	}
	if sameDisplayLabel(group, label) {
		return label
	}
	return group + " / " + label
}

func sameDisplayLabel(a, b string) bool {
	normalize := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		s = strings.TrimSuffix(s, "s")
		return s
	}
	return normalize(a) == normalize(b)
}

func (a *petApp) variantGroupLabel(i int) string {
	if i < 0 || i >= len(variants) {
		return ""
	}
	group := catalog.VariantGroupForSpecies(variants[i].SpeciesID)
	if a.lang == langEnglish {
		return group.LabelEN
	}
	return group.LabelJA
}

type variantMenuGroup struct {
	ID      string
	Label   string
	Indices []int
}

func (a *petApp) variantMenuGroups() []variantMenuGroup {
	byGroup := make(map[string][]int)
	for i, variant := range variants {
		groupID := catalog.VariantGroupIDForSpecies(variant.SpeciesID)
		byGroup[groupID] = append(byGroup[groupID], i)
	}
	groups := make([]variantMenuGroup, 0, len(byGroup))
	for _, group := range catalog.VariantGroups {
		indices := byGroup[group.ID]
		if len(indices) == 0 {
			continue
		}
		label := group.LabelJA
		if a.lang == langEnglish {
			label = group.LabelEN
		}
		groups = append(groups, variantMenuGroup{ID: group.ID, Label: label, Indices: indices})
	}
	return groups
}

func variantIndicesForGroupID(groupID string) []int {
	if groupID == "" {
		indices := make([]int, len(variants))
		for i := range variants {
			indices[i] = i
		}
		return indices
	}
	indices := make([]int, 0)
	for i, variant := range variants {
		if catalog.VariantGroupIDForSpecies(variant.SpeciesID) == groupID {
			indices = append(indices, i)
		}
	}
	return indices
}

func (a *petApp) appendGroupedVariantMenu(menu win.HMENU, selected int, commandID func(index int) uintptr) {
	for _, group := range a.variantMenuGroups() {
		submenu := win.CreatePopupMenu()
		for _, index := range group.Indices {
			flags := uint32(win.MF_STRING)
			if index == selected {
				flags |= win.MF_CHECKED
			}
			appendMenu(submenu, flags, commandID(index), syscall.StringToUTF16Ptr(a.variantLabel(index)))
		}
		appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(submenu), syscall.StringToUTF16Ptr(group.Label))
	}
}

func (a *petApp) drawForageProp(dst *image.RGBA, x, y, kind int) {
	if !foragePropsEnabled {
		return
	}
	if kind >= 0 && kind < len(a.forageSprites) && a.forageSprites[kind] != nil {
		src := a.forageSprites[kind]
		drawFacingImage(dst, src, image.Rect(x-forageW/2, y-forageH, x+forageW/2, y), 1)
		return
	}
	fillCircle(dst, x, y-2, 3, rgba(170, 150, 94, 220))
}

func drawPetSprite(dst *image.RGBA, src *image.RGBA, p *desktopPet, variant coatVariant, frame int, x, y int, w int, h int) {
	dir := normalizeDir(p.dir)
	if p.state == stateTurn {
		dir = turnDrawDirection(p.dir, p.nextDir)
	}
	rect := renderAdjustmentForFrame(variant.ID, frame).adjustRect(image.Rect(x, y, x+w, y+h))
	drawFacingImage(dst, src, rect, drawDirectionForVariant(dir, variant))
}

func turnDrawDirection(dir, nextDir int) int {
	if normalizeDir(dir) < 0 && normalizeDir(nextDir) > 0 {
		return -1
	}
	return 1
}

func drawDirectionForVariant(desiredDir int, variant coatVariant) int {
	return normalizeDir(desiredDir) * sourceFacingDirection(variant.ID)
}

func sourceFacingDirection(variantID string) int {
	switch variantID {
	case "sugar_glider_gray":
		return -1
	default:
		return 1
	}
}

func drawFacingImage(dst *image.RGBA, src *image.RGBA, r image.Rectangle, dir int) {
	if r.Empty() {
		return
	}
	sb := src.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		if y < dst.Bounds().Min.Y || y >= dst.Bounds().Max.Y {
			continue
		}
		sy := sb.Min.Y + (y-r.Min.Y)*sb.Dy()/max(1, r.Dy())
		for x := r.Min.X; x < r.Max.X; x++ {
			if x < dst.Bounds().Min.X || x >= dst.Bounds().Max.X {
				continue
			}
			dx := x - r.Min.X
			sx := sb.Min.X + dx*sb.Dx()/max(1, r.Dx())
			if dir < 0 {
				sx = sb.Max.X - 1 - dx*sb.Dx()/max(1, r.Dx())
			}
			c := src.RGBAAt(sx, sy)
			if c.A == 0 {
				continue
			}
			dst.SetRGBA(x, y, overRGBA(dst.RGBAAt(x, y), c))
		}
	}
}

func (a *petApp) enterWheelFromTyping(index int, p *desktopPet) {
	alreadyRunning := p.state == stateWheel
	p.state = stateWheel
	if !alreadyRunning {
		p.frame = 0
		p.motionSet = rand.Intn(motionSets)
	}
	p.item = noItem
	p.carryKind = noItem
	p.moveSpeed = 0
	p.stateTicks = wheelKeyHold
	w, _ := a.petSpriteSize(index)
	p.x = clamp(a.wheelX-wheelSize/2, 0, max(0, a.sceneW-w))
}

func (a *petApp) leaveWheel(index int, p *desktopPet) {
	p.state = stateScurry
	p.frame = 0
	p.motionSet = rand.Intn(motionSets)
	p.dir = 1
	p.nextDir = 1
	p.moveSpeed = a.speed + 1 + rand.Intn(2)
	p.stateTicks = 16 + rand.Intn(20)
	w, _ := a.petSpriteSize(index)
	p.x = clamp(a.wheelX+wheelSize/2-20, 0, max(0, a.sceneW-w))
}

type spriteCache struct {
	mu     sync.Mutex
	loaded map[string][][]*image.RGBA
}

func newSpriteCache() *spriteCache {
	return &spriteCache{loaded: make(map[string][][]*image.RGBA)}
}

func (c *spriteCache) retain(ids map[string]bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id := range c.loaded {
		if !ids[id] {
			delete(c.loaded, id)
		}
	}
}

func (c *spriteCache) frame(variant coatVariant, set int, frame int) *image.RGBA {
	sets := c.variantSets(variant)
	set = clamp(set, 0, len(sets)-1)
	frame = clamp(frame, 0, len(sets[set])-1)
	return sets[set][frame]
}

func (c *spriteCache) variantSets(variant coatVariant) [][]*image.RGBA {
	c.mu.Lock()
	defer c.mu.Unlock()
	if sets := c.loaded[variant.ID]; sets != nil {
		return sets
	}
	sets := loadSpriteVariant(variant)
	c.loaded[variant.ID] = sets
	return sets
}

func loadSpriteVariant(variant coatVariant) [][]*image.RGBA {
	sets := make([][]*image.RGBA, 0, motionSets)
	for set := 0; set < motionSets; set++ {
		name := fmt.Sprintf("sprites/%s_set%02d.png", variant.SpriteBase, set)
		data, err := fs.ReadFile(appassets.FS, name)
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		if img.Bounds().Dx() != frameW*frameCount || img.Bounds().Dy() != frameH {
			panic(fmt.Sprintf("%s must be %dx%d; run cmd/importsheet and cmd/importanimals", name, frameW*frameCount, frameH))
		}
		frames := make([]*image.RGBA, 0, frameCount)
		for i := 0; i < frameCount; i++ {
			r := image.Rect(i*frameW, 0, (i+1)*frameW, frameH)
			frame := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
			draw.Draw(frame, frame.Bounds(), img, r.Min, draw.Src)
			frames = append(frames, frame)
		}
		sets = append(sets, frames)
	}
	return sets
}

func loadWheelSprite() *image.RGBA {
	data, err := fs.ReadFile(appassets.FS, "sprites/wheel.png")
	if err != nil {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	dst := image.NewRGBA(image.Rect(0, 0, wheelSize, wheelSize))
	if img.Bounds().Dx() == wheelSize && img.Bounds().Dy() == wheelSize {
		draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
		return dst
	}
	src := image.NewRGBA(img.Bounds())
	draw.Draw(src, src.Bounds(), img, img.Bounds().Min, draw.Src)
	scaled := scaleImageTo(src, wheelSize, wheelSize)
	draw.Draw(dst, dst.Bounds(), scaled, image.Point{}, draw.Src)
	return dst
}

func loadForageSprites() []*image.RGBA {
	names := []string{"forage_hay", "forage_twig", "forage_seed"}
	out := make([]*image.RGBA, len(names))
	for i, name := range names {
		data, err := fs.ReadFile(appassets.FS, "sprites/"+name+".png")
		if err != nil {
			continue
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			continue
		}
		dst := image.NewRGBA(image.Rect(0, 0, forageW, forageH))
		draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
		out[i] = dst
	}
	return out
}

func (a *petApp) showStartupToast() {
	a.showTrayBalloon(
		a.localText("Animals Desktop 起動中", "Animals Desktop is running"),
		a.localText("タスクトレイから設定と終了ができます。", "Use the tray icon for settings and exit."),
	)
}

func (a *petApp) showTrayBalloon(title, body string) {
	if a.hwnd == 0 {
		return
	}
	var nid win.NOTIFYICONDATA
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = a.hwnd
	nid.UID = 1
	nid.UFlags = win.NIF_INFO
	nid.DwInfoFlags = win.NIIF_INFO | win.NIIF_RESPECT_QUIET_TIME
	copyUTF16Limited(nid.SzInfoTitle[:], title)
	copyUTF16Limited(nid.SzInfo[:], body)
	win.Shell_NotifyIcon(win.NIM_MODIFY, &nid)
}

func copyUTF16Limited(dst []uint16, text string) {
	if len(dst) == 0 {
		return
	}
	src := syscall.StringToUTF16(text)
	n := min(len(dst), len(src))
	copy(dst[:n], src[:n])
	if n == len(dst) {
		dst[len(dst)-1] = 0
	}
}

func (a *petApp) updateChecking() bool {
	return a.update.checking.Load() || a.update.installing.Load()
}

func (a *petApp) hasInstallableUpdate() bool {
	if !networkUpdatesEnabled {
		return false
	}
	rel := a.currentRelease()
	if rel == nil || !isNewerVersion(rel.TagName, appVersion) {
		return false
	}
	return selectUpdateAsset(rel, runtime.GOARCH) != nil
}

func (a *petApp) updateCheckMenuLabel() string {
	if !networkUpdatesEnabled {
		return a.localText("ネットワーク無効版", "Network disabled edition")
	}
	if a.update.installing.Load() {
		return a.localText("アップデート適用中...", "Installing update...")
	}
	if a.update.checking.Load() {
		return a.localText("アップデート確認中...", "Checking for updates...")
	}
	return a.localText("アップデートを確認", "Check for updates")
}

func (a *petApp) updateInstallMenuLabel() string {
	rel := a.currentRelease()
	if rel == nil {
		return a.localText("アップデートをインストール", "Install update")
	}
	return a.localText(
		fmt.Sprintf("%s をインストール", rel.TagName),
		fmt.Sprintf("Install %s", rel.TagName),
	)
}

func (a *petApp) startUpdateCheck(manual bool) {
	if !networkUpdatesEnabled {
		if manual {
			a.showTrayBalloon(
				a.localText("ネットワーク無効版", "Network disabled edition"),
				a.localText("このビルドでは自動更新確認とダウンロードを無効化しています。", "This build disables update checks and downloads."),
			)
		}
		return
	}
	if !a.update.checking.CompareAndSwap(false, true) {
		if manual {
			a.showTrayBalloon(
				a.localText("確認中です", "Already checking"),
				a.localText("最新アップデートの確認が進行中です。", "An update check is already in progress."),
			)
		}
		return
	}
	go func() {
		rel, err := fetchLatestRelease()
		a.setUpdateResult(rel, err)
		var notify uintptr
		if manual {
			notify = 1
		}
		if err != nil {
			win.PostMessage(a.hwnd, wmUpdateFailed, notify, 0)
			return
		}
		win.PostMessage(a.hwnd, wmUpdateReady, notify, 0)
	}()
}

func (a *petApp) onUpdateReady(manual bool) {
	a.update.checking.Store(false)
	rel := a.currentRelease()
	if rel == nil {
		if manual {
			a.showTrayBalloon(
				a.localText("確認できませんでした", "Update check failed"),
				a.localText("最新リリース情報が空でした。", "The latest release response was empty."),
			)
		}
		return
	}
	asset := selectUpdateAsset(rel, runtime.GOARCH)
	if asset == nil {
		if manual {
			a.showTrayBalloon(
				a.localText("配布 zip が見つかりません", "Update package not found"),
				a.localText(updateAssetName(runtime.GOARCH)+" が Release にありません。", updateAssetName(runtime.GOARCH)+" is missing from the release."),
			)
		}
		return
	}
	if !isNewerVersion(rel.TagName, appVersion) {
		if manual {
			a.showTrayBalloon(
				a.localText("最新版です", "Up to date"),
				a.localText("インストール済みの Animals Desktop は最新です。", "The installed Animals Desktop is current."),
			)
		}
		return
	}
	a.showTrayBalloon(
		a.localText("アップデートがあります", "Update available"),
		a.localText(
			fmt.Sprintf("%s をトレイメニューからインストールできます。", rel.TagName),
			fmt.Sprintf("%s can be installed from the tray menu.", rel.TagName),
		),
	)
}

func (a *petApp) onUpdateFailed(notify bool) {
	a.update.checking.Store(false)
	a.update.installing.Store(false)
	if !notify {
		return
	}
	a.showTrayBalloon(
		a.localText("アップデートに失敗しました", "Update failed"),
		a.currentUpdateError(),
	)
}

func (a *petApp) installLatestUpdate() {
	if !networkUpdatesEnabled {
		a.showTrayBalloon(
			a.localText("ネットワーク無効版", "Network disabled edition"),
			a.localText("このビルドではアップデートのダウンロードと適用を無効化しています。", "This build disables update downloads and installation."),
		)
		return
	}
	rel := a.currentRelease()
	if rel == nil || !isNewerVersion(rel.TagName, appVersion) {
		a.showTrayBalloon(
			a.localText("アップデートなし", "No update available"),
			a.localText("先にアップデートを確認してください。", "Check for updates first."),
		)
		return
	}
	asset := selectUpdateAsset(rel, runtime.GOARCH)
	if asset == nil {
		a.showTrayBalloon(
			a.localText("配布 zip が見つかりません", "Update package not found"),
			a.localText(updateAssetName(runtime.GOARCH)+" が Release にありません。", updateAssetName(runtime.GOARCH)+" is missing from the release."),
		)
		return
	}
	if !a.update.installing.CompareAndSwap(false, true) {
		a.showTrayBalloon(
			a.localText("適用中です", "Install already running"),
			a.localText("アップデートのダウンロードと適用を進めています。", "The update is already being downloaded and installed."),
		)
		return
	}
	a.showTrayBalloon(
		a.localText("アップデートを開始", "Starting update"),
		a.localText("ダウンロード後に Animals Desktop を再起動します。", "Animals Desktop will restart after download."),
	)
	go func(asset githubReleaseAsset) {
		err := downloadAndStartUpdater(asset)
		a.setUpdateResult(rel, err)
		if err != nil {
			win.PostMessage(a.hwnd, wmUpdateFailed, 1, 0)
			return
		}
		win.PostMessage(a.hwnd, wmUpdateInstallReady, 0, 0)
	}(*asset)
}

func (a *petApp) onUpdateInstallReady() {
	a.showTrayBalloon(
		a.localText("アップデートを適用します", "Applying update"),
		a.localText("本体を終了して、新しいバージョンで再起動します。", "The app will close and restart with the new version."),
	)
	a.closing.Store(true)
	win.DestroyWindow(a.hwnd)
}

func (a *petApp) setUpdateResult(rel *githubRelease, err error) {
	a.update.mu.Lock()
	defer a.update.mu.Unlock()
	if rel != nil {
		cp := *rel
		cp.Assets = append([]githubReleaseAsset(nil), rel.Assets...)
		a.update.latest = &cp
	}
	if err != nil {
		a.update.lastError = err.Error()
	} else {
		a.update.lastError = ""
	}
}

func (a *petApp) currentRelease() *githubRelease {
	a.update.mu.Lock()
	defer a.update.mu.Unlock()
	if a.update.latest == nil {
		return nil
	}
	cp := *a.update.latest
	cp.Assets = append([]githubReleaseAsset(nil), a.update.latest.Assets...)
	return &cp
}

func (a *petApp) currentUpdateError() string {
	a.update.mu.Lock()
	defer a.update.mu.Unlock()
	if a.update.lastError == "" {
		return a.localText("詳細なエラーはありません。", "No detailed error is available.")
	}
	return a.update.lastError
}

func selectUpdateAsset(rel *githubRelease, goarch string) *githubReleaseAsset {
	if rel == nil {
		return nil
	}
	want := updateAssetName(goarch)
	for i := range rel.Assets {
		asset := &rel.Assets[i]
		if strings.EqualFold(asset.Name, want) && asset.BrowserDownloadURL != "" {
			return asset
		}
	}
	return nil
}

func updateAssetName(goarch string) string {
	switch goarch {
	case "386":
		return "AnimalsDesktop-windows-386.zip"
	default:
		return "AnimalsDesktop-windows-amd64.zip"
	}
}

func isNewerVersion(latest, current string) bool {
	latest = normalizeVersionTag(latest)
	current = normalizeVersionTag(current)
	if latest == "" || latest == current {
		return false
	}
	latestParts, latestOK := parseVersionParts(latest)
	currentParts, currentOK := parseVersionParts(current)
	if latestOK && currentOK {
		for i := 0; i < max(len(latestParts), len(currentParts)); i++ {
			lv, cv := 0, 0
			if i < len(latestParts) {
				lv = latestParts[i]
			}
			if i < len(currentParts) {
				cv = currentParts[i]
			}
			if lv != cv {
				return lv > cv
			}
		}
		return false
	}
	return current == "" || current == "dev" || strings.HasPrefix(current, "pages-")
}

func normalizeVersionTag(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	return version
}

func parseVersionParts(version string) ([]int, bool) {
	if version == "" {
		return nil, false
	}
	if cut := strings.IndexAny(version, "-+"); cut >= 0 {
		version = version[:cut]
	}
	parts := strings.Split(version, ".")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			return nil, false
		}
		value := 0
		for _, r := range part {
			if r < '0' || r > '9' {
				return nil, false
			}
			value = value*10 + int(r-'0')
		}
		out = append(out, value)
	}
	return out, true
}

func (a *petApp) installTray() {
	iconPath := filepath.Join(os.TempDir(), "animals-desktop-tray.ico")
	if data, err := fs.ReadFile(appassets.FS, "tray.ico"); err == nil {
		_ = os.WriteFile(iconPath, data, 0o644)
	}
	a.trayIcon = win.HICON(win.LoadImage(0, syscall.StringToUTF16Ptr(iconPath), win.IMAGE_ICON, 0, 0, win.LR_LOADFROMFILE|win.LR_DEFAULTSIZE))
	if a.trayIcon == 0 {
		a.trayIcon = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	}
	var nid win.NOTIFYICONDATA
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = a.hwnd
	nid.UID = 1
	nid.UFlags = win.NIF_MESSAGE | win.NIF_ICON | win.NIF_TIP
	nid.UCallbackMessage = wmTray
	nid.HIcon = a.trayIcon
	copy(nid.SzTip[:], syscall.StringToUTF16(appName))
	win.Shell_NotifyIcon(win.NIM_ADD, &nid)
}

func (a *petApp) showTrayMenu() {
	menu := win.CreatePopupMenu()
	coatMenu := win.CreatePopupMenu()
	a.appendGroupedVariantMenu(coatMenu, a.variant, func(index int) uintptr {
		return uintptr(menuVariantBase + uint16(index))
	})
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(coatMenu), syscall.StringToUTF16Ptr(a.txt("coatColor")))
	coatModeMenu := win.CreatePopupMenu()
	appendChecked(coatModeMenu, menuCoatFixed, a.txt("coatFixed"), a.coatMode == coatFixed)
	appendChecked(coatModeMenu, menuCoatSelected, a.txt("coatSelected"), a.coatMode == coatSelected)
	appendChecked(coatModeMenu, menuCoatRandom, a.txt("coatRandom"), a.coatMode == coatRandom)
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(coatModeMenu), syscall.StringToUTF16Ptr(a.txt("coatMode")))
	randomGroupMenu := win.CreatePopupMenu()
	appendChecked(randomGroupMenu, menuRandomAll, a.randomAllMenuLabel(), normalizeVariantGroupID(a.randomGroupID) == "")
	for i, group := range a.variantMenuGroups() {
		appendChecked(randomGroupMenu, menuRandomGroupBase+uint16(i), group.Label, normalizeVariantGroupID(a.randomGroupID) == group.ID)
	}
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(randomGroupMenu), syscall.StringToUTF16Ptr(a.txt("randomTarget")))
	appendMenu(menu, win.MF_SEPARATOR, 0, nil)

	speedMenu := win.CreatePopupMenu()
	appendChecked(speedMenu, menuSpeedSlow, a.txt("speedSlow"), a.speed == 2)
	appendChecked(speedMenu, menuSpeedNormal, a.txt("speedNormal"), a.speed == 3)
	appendChecked(speedMenu, menuSpeedFast, a.txt("speedFast"), a.speed == 5)
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(speedMenu), syscall.StringToUTF16Ptr(a.txt("speed")))

	modeMenu := win.CreatePopupMenu()
	appendChecked(modeMenu, menuModeKeyboard, a.txt("modeKeyboard"), a.mode == modeKeyboard)
	appendChecked(modeMenu, menuModeRandom, a.txt("modeRandom"), a.mode == modeRandom)
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(modeMenu), syscall.StringToUTF16Ptr(a.txt("mode")))

	countMenu := win.CreatePopupMenu()
	appendChecked(countMenu, menuCount1, "1", a.petCount == 1)
	appendChecked(countMenu, menuCount2, "2", a.petCount == 2)
	appendChecked(countMenu, menuCount3, "3", a.petCount == 3)
	appendChecked(countMenu, menuCount5, "5", a.petCount == 5)
	appendChecked(countMenu, menuCount10, "10", a.petCount == 10)
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(countMenu), syscall.StringToUTF16Ptr(a.txt("petCount")))

	appendChecked(menu, menuWheelToggle, a.txt("typingWheel"), a.wheelEnabled)
	appendMenu(menu, win.MF_STRING, uintptr(menuHideToggle), syscall.StringToUTF16Ptr(a.temporaryVisibilityLabel()))
	appendMenu(menu, win.MF_SEPARATOR, 0, nil)
	languageMenu := win.CreatePopupMenu()
	appendChecked(languageMenu, menuLangJapanese, "日本語", a.lang == langJapanese)
	appendChecked(languageMenu, menuLangEnglish, "English", a.lang == langEnglish)
	appendMenu(menu, win.MF_POPUP|win.MF_STRING, uintptr(languageMenu), syscall.StringToUTF16Ptr(a.txt("language")))
	appendMenu(menu, win.MF_SEPARATOR, 0, nil)
	updateFlags := uint32(win.MF_STRING)
	if a.updateChecking() || !networkUpdatesEnabled {
		updateFlags |= win.MF_GRAYED
	}
	appendMenu(menu, updateFlags, uintptr(menuCheckUpdate), syscall.StringToUTF16Ptr(a.updateCheckMenuLabel()))
	if networkUpdatesEnabled && a.hasInstallableUpdate() {
		appendMenu(menu, win.MF_STRING, uintptr(menuInstallUpdate), syscall.StringToUTF16Ptr(a.updateInstallMenuLabel()))
	}
	appendMenu(menu, win.MF_SEPARATOR, 0, nil)
	appendMenu(menu, win.MF_STRING, uintptr(menuSettings), syscall.StringToUTF16Ptr(a.txt("settingsTitle")))
	appendMenu(menu, win.MF_STRING, uintptr(menuExit), syscall.StringToUTF16Ptr(a.txt("exit")))

	var pt win.POINT
	win.GetCursorPos(&pt)
	win.SetForegroundWindow(a.hwnd)
	cmd := win.TrackPopupMenu(menu, win.TPM_RETURNCMD|win.TPM_RIGHTBUTTON, pt.X, pt.Y, 0, a.hwnd, nil)
	win.DestroyMenu(menu)
	if cmd == 0 {
		return
	}
	a.handleMenu(uint16(cmd))
}

func appendChecked(menu win.HMENU, id uint16, label string, checked bool) {
	flags := uint32(win.MF_STRING)
	if checked {
		flags |= win.MF_CHECKED
	}
	appendMenu(menu, flags, uintptr(id), syscall.StringToUTF16Ptr(label))
}

func (a *petApp) handleMenu(id uint16) {
	if !a.handleMenuCommand(id) {
		return
	}
	if id == menuHideToggle {
		return
	}
	a.syncSettingsWindow()
	a.persistSettings()
}

func (a *petApp) handleMenuCommand(id uint16) bool {
	switch {
	case id == menuExit:
		a.closing.Store(true)
		win.DestroyWindow(a.hwnd)
	case id == menuSettings:
		a.showSettings()
	case id == menuCheckUpdate:
		a.startUpdateCheck(true)
	case id == menuInstallUpdate:
		a.installLatestUpdate()
	case id == menuModeKeyboard:
		a.mode = modeKeyboard
		for i := range a.pets {
			a.pets[i].state = stateIdle
			a.pets[i].stateTicks = 1
			a.pets[i].moveSpeed = 0
		}
	case id == menuModeRandom:
		a.mode = modeRandom
		for i := range a.pets {
			a.chooseRandomAction(&a.pets[i])
		}
	case id == menuSpeedSlow:
		a.speed = 2
	case id == menuSpeedNormal:
		a.speed = 3
	case id == menuSpeedFast:
		a.speed = 5
	case id == menuCount1:
		a.setPetCount(1)
		a.resetPosition()
	case id == menuCount2:
		a.setPetCount(2)
		a.resetPosition()
	case id == menuCount3:
		a.setPetCount(3)
		a.resetPosition()
	case id == menuCount5:
		a.setPetCount(5)
		a.resetPosition()
	case id == menuCount10:
		a.setPetCount(10)
		a.resetPosition()
	case id == menuWheelToggle:
		a.wheelEnabled = !a.wheelEnabled
		for i := range a.pets {
			if a.pets[i].state == stateWheel {
				a.leaveWheel(i, &a.pets[i])
			}
		}
	case id == menuHideToggle:
		a.setOverlayHidden(!a.overlayHidden)
	case id == menuCoatFixed:
		a.setCoatMode(coatFixed)
	case id == menuCoatSelected:
		a.setCoatMode(coatSelected)
	case id == menuCoatRandom:
		a.setCoatMode(coatRandom)
	case id == menuRandomAll:
		a.setRandomGroupID("")
		a.setCoatMode(coatRandom)
		a.render()
	case id >= menuRandomGroupBase:
		groupIndex := int(id - menuRandomGroupBase)
		groups := a.variantMenuGroups()
		if groupIndex < 0 || groupIndex >= len(groups) {
			return false
		}
		a.setRandomGroupID(groups[groupIndex].ID)
		a.setCoatMode(coatRandom)
		a.render()
	case id == menuLangJapanese:
		a.lang = langJapanese
	case id == menuLangEnglish:
		a.lang = langEnglish
	case id >= menuVariantBase && int(id-menuVariantBase) < len(variants):
		a.setFixedVariant(int(id - menuVariantBase))
		a.setCoatMode(coatFixed)
		a.render()
	default:
		return false
	}
	return true
}

func (a *petApp) cleanup() {
	win.KillTimer(a.hwnd, timerID)
	a.persistSettings()
	if a.settingsHwnd != 0 {
		win.DestroyWindow(a.settingsHwnd)
		a.settingsHwnd = 0
	}
	if a.renameHwnd != 0 {
		win.DestroyWindow(a.renameHwnd)
		a.renameHwnd = 0
	}
	if a.nameHwnd != 0 {
		win.DestroyWindow(a.nameHwnd)
		a.nameHwnd = 0
	}
	if a.keyHook != 0 {
		unhookWindowsHookEx(a.keyHook)
	}
	if a.mouseHook != 0 {
		unhookWindowsHookEx(a.mouseHook)
	}
	var nid win.NOTIFYICONDATA
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = a.hwnd
	nid.UID = 1
	win.Shell_NotifyIcon(win.NIM_DELETE, &nid)
	if a.trayIcon != 0 {
		win.DestroyIcon(a.trayIcon)
	}
	if a.settingsBrush != 0 {
		win.DeleteObject(win.HGDIOBJ(a.settingsBrush))
		a.settingsBrush = 0
	}
	if a.settingsCard != 0 {
		win.DeleteObject(win.HGDIOBJ(a.settingsCard))
		a.settingsCard = 0
	}
	if a.settingsFont != 0 {
		win.DeleteObject(win.HGDIOBJ(a.settingsFont))
		a.settingsFont = 0
	}
	if a.settingsTitleFont != 0 {
		win.DeleteObject(win.HGDIOBJ(a.settingsTitleFont))
		a.settingsTitleFont = 0
	}
	if a.settingsSmallFont != 0 {
		win.DeleteObject(win.HGDIOBJ(a.settingsSmallFont))
		a.settingsSmallFont = 0
	}
}

func (a *petApp) installKeyboardHook() {
	cb := syscall.NewCallback(func(code int, wParam uintptr, lParam uintptr) uintptr {
		if code >= 0 && (wParam == win.WM_KEYDOWN || wParam == win.WM_SYSKEYDOWN) {
			a.postTypingFromHook()
		}
		return callNextHookEx(0, code, wParam, lParam)
	})
	a.keyHook = setWindowsHookEx(whKeyboardLL, cb, a.hinst, 0)
	a.keyHookFailed = a.keyHook == 0
}

func (a *petApp) installMouseHook() {
	cb := syscall.NewCallback(func(code int, wParam uintptr, lParam uintptr) uintptr {
		if code >= 0 && wParam == win.WM_LBUTTONDOWN {
			a.postMouseClickFromHook(lParam)
		}
		return callNextHookEx(0, code, wParam, lParam)
	})
	a.mouseHook = setWindowsHookEx(whMouseLL, cb, a.hinst, 0)
	a.mouseHookFailed = a.mouseHook == 0
}

func (a *petApp) postTypingFromHook() {
	defer recoverHookCallback()
	win.PostMessage(a.hwnd, wmTyping, 0, 0)
}

func (a *petApp) postMouseClickFromHook(lParam uintptr) {
	defer recoverHookCallback()
	pt := mouseHookPoint(lParam)
	win.PostMessage(a.hwnd, wmMouseClick, uintptr(uint32(pt.X)), uintptr(uint32(pt.Y)))
}

func recoverHookCallback() {
	_ = recover()
}

func mouseHookPoint(lParam uintptr) win.POINT {
	var hook mouseHookStruct
	if lParam != 0 {
		procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&hook)), lParam, unsafe.Sizeof(hook))
	}
	return hook.pt
}

func appendMenu(menu win.HMENU, flags uint32, item uintptr, text *uint16) bool {
	var textPtr uintptr
	if text != nil {
		textPtr = uintptr(unsafe.Pointer(text))
	}
	ret, _, _ := procAppendMenuW.Call(uintptr(menu), uintptr(flags), item, textPtr)
	return ret != 0
}

func getDlgCtrlID(hwnd win.HWND) int32 {
	ret, _, _ := procGetDlgCtrlID.Call(uintptr(hwnd))
	return int32(ret)
}

func setWindowText(hwnd win.HWND, text string) bool {
	ret, _, _ := procSetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))))
	return ret != 0
}

func getWindowText(hwnd win.HWND) string {
	if hwnd == 0 {
		return ""
	}
	length := int(win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0))
	buf := make([]uint16, max(1, length+1))
	win.SendMessage(hwnd, win.WM_GETTEXT, uintptr(len(buf)), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func setWindowsHookEx(idHook int, callback uintptr, module win.HINSTANCE, threadID uint32) uintptr {
	ret, _, _ := procSetWindowsHookExW.Call(uintptr(idHook), callback, uintptr(module), uintptr(threadID))
	return ret
}

func unhookWindowsHookEx(hook uintptr) bool {
	ret, _, _ := procUnhookWindowsHook.Call(hook)
	return ret != 0
}

func callNextHookEx(hook uintptr, code int, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := procCallNextHookExProc.Call(hook, uintptr(code), wParam, lParam)
	return ret
}

func updateLayeredWindowNative(hwnd win.HWND, dstDC win.HDC, dst *win.POINT, size *win.SIZE, srcDC win.HDC, src *win.POINT, key uint32, blend *win.BLENDFUNCTION, flags uint32) bool {
	ret, _, _ := procUpdateLayeredWin.Call(
		uintptr(hwnd),
		uintptr(dstDC),
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(size)),
		uintptr(srcDC),
		uintptr(unsafe.Pointer(src)),
		uintptr(key),
		uintptr(unsafe.Pointer(blend)),
		uintptr(flags),
	)
	return ret != 0
}

func updateLayeredWindow(hwnd win.HWND, img *image.RGBA, x, y int) bool {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	screenDC := win.GetDC(0)
	memDC := win.CreateCompatibleDC(screenDC)
	defer win.ReleaseDC(0, screenDC)
	defer win.DeleteDC(memDC)

	var bmi win.BITMAPINFOHEADER
	bmi.BiSize = uint32(unsafe.Sizeof(bmi))
	bmi.BiWidth = int32(w)
	bmi.BiHeight = -int32(h)
	bmi.BiPlanes = 1
	bmi.BiBitCount = 32
	bmi.BiCompression = win.BI_RGB

	var bits unsafe.Pointer
	bitmap := win.CreateDIBSection(memDC, &bmi, win.DIB_RGB_COLORS, &bits, 0, 0)
	if bitmap == 0 {
		return false
	}
	defer win.DeleteObject(win.HGDIOBJ(bitmap))
	dst := unsafe.Slice((*byte)(bits), w*h*4)
	copyPremultipliedBGRA(dst, img)

	old := win.SelectObject(memDC, win.HGDIOBJ(bitmap))
	defer win.SelectObject(memDC, old)

	ptDst := win.POINT{X: int32(x), Y: int32(y)}
	size := win.SIZE{CX: int32(w), CY: int32(h)}
	ptSrc := win.POINT{X: 0, Y: 0}
	blend := win.BLENDFUNCTION{BlendOp: acSrcOver, SourceConstantAlpha: 255, AlphaFormat: win.AC_SRC_ALPHA}
	return updateLayeredWindowNative(hwnd, screenDC, &ptDst, &size, memDC, &ptSrc, 0, &blend, ulwAlpha)
}

type displayArea struct {
	Work    win.RECT
	Screen  win.RECT
	Primary bool
	DPI     int
}

type overlaySegment struct {
	Rect       win.RECT
	Screen     win.RECT
	SceneLeft  int
	SceneRight int
	DPI        int
}

var (
	monitorEnumCallback = syscall.NewCallback(monitorEnumProc)
	monitorEnumMu       sync.Mutex
	monitorEnumAreas    *[]displayArea
)

func monitorEnumProc(hMonitor win.HMONITOR, _ win.HDC, _ *win.RECT, lParam uintptr) uintptr {
	if monitorEnumAreas == nil {
		return 0
	}
	info := win.MONITORINFO{CbSize: uint32(unsafe.Sizeof(win.MONITORINFO{}))}
	if win.GetMonitorInfo(hMonitor, &info) {
		*monitorEnumAreas = append(*monitorEnumAreas, displayArea{
			Work:    info.RcWork,
			Screen:  info.RcMonitor,
			Primary: info.DwFlags&monitorPrimaryFlag != 0,
			DPI:     monitorDPI(hMonitor),
		})
	}
	return 1
}

func monitorDPI(hMonitor win.HMONITOR) int {
	if hMonitor == 0 || procGetDpiForMonitor.Find() != nil {
		return defaultDPI
	}
	var dpiX uint32
	var dpiY uint32
	ret, _, _ := procGetDpiForMonitor.Call(uintptr(hMonitor), monitorDPIType, uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
	if ret != 0 || dpiX == 0 {
		return defaultDPI
	}
	return clamp(int(dpiX), 72, 384)
}

func (a *petApp) selectedDisplayArea() displayArea {
	_, combined := a.selectedDisplayAreasAndCombined()
	return combined
}

func (a *petApp) selectedDisplayAreasAndCombined() ([]displayArea, displayArea) {
	areas := displayAreaForScope(a.displayScope)
	if len(areas) == 0 {
		fallback := displayArea{Work: workArea(), Screen: screenArea(), Primary: true}
		return []displayArea{fallback}, fallback
	}
	scope, start, end := a.normalizedDisplaySelection(len(areas))
	a.displayScope = scope
	a.displayIndex = start
	a.displaySpanEnd = end
	selected := append([]displayArea(nil), areas[start:end+1]...)
	if scope == displayScopeSpan {
		return selected, combineDisplayAreas(selected)
	}
	return selected, selected[0]
}

func (a *petApp) overlaySegments() []overlaySegment {
	areas, display := a.selectedDisplayAreasAndCombined()
	overlay := a.overlayRectFor(display.Work, display.Screen)
	return a.overlaySegmentsForAreas(areas, overlay)
}

func (a *petApp) overlaySegmentsForAreas(areas []displayArea, overlay win.RECT) []overlaySegment {
	if overlay.Right <= overlay.Left {
		return nil
	}
	segments := make([]overlaySegment, 0, len(areas))
	for _, area := range areas {
		base := a.overlayBaseFor(area.Work, area.Screen)
		left := max(int(base.Left), int(overlay.Left))
		right := min(int(base.Right), int(overlay.Right))
		if right-left < spriteW {
			continue
		}
		dpi := normalizeDPI(area.DPI)
		height := scaleForDPI(sceneH, dpi)
		y := a.overlayYForBaseDPI(base, area.Screen, dpi)
		segments = append(segments, overlaySegment{
			Rect: win.RECT{
				Left:   int32(left),
				Top:    int32(y),
				Right:  int32(right),
				Bottom: int32(y + height),
			},
			Screen:     area.Screen,
			SceneLeft:  left - int(overlay.Left),
			SceneRight: right - int(overlay.Left),
			DPI:        dpi,
		})
	}
	return segments
}

func (a *petApp) overlaySegmentForSceneX(sceneX int) (overlaySegment, bool) {
	segments := a.overlaySegments()
	for _, segment := range segments {
		if sceneX >= segment.SceneLeft && sceneX < segment.SceneRight {
			return segment, true
		}
	}
	if len(segments) == 0 {
		return overlaySegment{}, false
	}
	nearest := segments[0]
	bestDistance := abs(sceneX - nearest.SceneLeft)
	for _, segment := range segments {
		distance := 0
		if sceneX < segment.SceneLeft {
			distance = segment.SceneLeft - sceneX
		} else if sceneX >= segment.SceneRight {
			distance = sceneX - segment.SceneRight + 1
		}
		if distance < bestDistance {
			nearest = segment
			bestDistance = distance
		}
	}
	return nearest, true
}

func monitorAreas() []displayArea {
	areas := make([]displayArea, 0, 4)
	monitorEnumMu.Lock()
	monitorEnumAreas = &areas
	ret, _, _ := procEnumDisplayMonitors.Call(0, 0, monitorEnumCallback, 0)
	monitorEnumAreas = nil
	monitorEnumMu.Unlock()
	if ret == 0 || len(areas) == 0 {
		return nil
	}
	sort.SliceStable(areas, func(i, j int) bool {
		if areas[i].Primary != areas[j].Primary {
			return areas[i].Primary
		}
		if areas[i].Screen.Left != areas[j].Screen.Left {
			return areas[i].Screen.Left < areas[j].Screen.Left
		}
		return areas[i].Screen.Top < areas[j].Screen.Top
	})
	return areas
}

func monitorAreasByPosition() []displayArea {
	areas := append([]displayArea(nil), monitorAreas()...)
	sort.SliceStable(areas, func(i, j int) bool {
		if areas[i].Screen.Left != areas[j].Screen.Left {
			return areas[i].Screen.Left < areas[j].Screen.Left
		}
		if areas[i].Screen.Top != areas[j].Screen.Top {
			return areas[i].Screen.Top < areas[j].Screen.Top
		}
		if areas[i].Primary != areas[j].Primary {
			return areas[i].Primary
		}
		return false
	})
	return areas
}

func displayAreaForScope(scope displayScope) []displayArea {
	if normalizeDisplayScope(int(scope)) == displayScopeSpan {
		return monitorAreasByPosition()
	}
	return monitorAreas()
}

func sameDisplayArea(a, b displayArea) bool {
	return a.Screen.Left == b.Screen.Left &&
		a.Screen.Top == b.Screen.Top &&
		a.Screen.Right == b.Screen.Right &&
		a.Screen.Bottom == b.Screen.Bottom
}

func findDisplayAreaIndex(areas []displayArea, target displayArea) int {
	for i, area := range areas {
		if sameDisplayArea(area, target) {
			return i
		}
	}
	return 0
}

func combineDisplayAreas(areas []displayArea) displayArea {
	if len(areas) == 0 {
		return displayArea{}
	}
	screen := areas[0].Screen
	workLeft := areas[0].Work.Left
	workRight := areas[0].Work.Right
	workTop := areas[0].Work.Top
	workBottom := areas[0].Work.Bottom
	primary := areas[0].Primary
	for _, area := range areas[1:] {
		screen.Left = min32(screen.Left, area.Screen.Left)
		screen.Top = min32(screen.Top, area.Screen.Top)
		screen.Right = max32(screen.Right, area.Screen.Right)
		screen.Bottom = max32(screen.Bottom, area.Screen.Bottom)
		workLeft = min32(workLeft, area.Work.Left)
		workRight = max32(workRight, area.Work.Right)
		workTop = max32(workTop, area.Work.Top)
		workBottom = min32(workBottom, area.Work.Bottom)
		primary = primary || area.Primary
	}
	work := win.RECT{Left: workLeft, Top: workTop, Right: workRight, Bottom: workBottom}
	if work.Right <= work.Left || work.Bottom <= work.Top {
		work = screen
	}
	return displayArea{Work: work, Screen: screen, Primary: primary}
}

func workArea() win.RECT {
	var rect win.RECT
	if !win.SystemParametersInfo(spiGetWorkArea, 0, unsafe.Pointer(&rect), 0) {
		rect = win.RECT{Left: 0, Top: 0, Right: 1280, Bottom: 720}
	}
	return rect
}

func screenArea() win.RECT {
	left := int32(win.GetSystemMetrics(win.SM_XVIRTUALSCREEN))
	top := int32(win.GetSystemMetrics(win.SM_YVIRTUALSCREEN))
	width := win.GetSystemMetrics(win.SM_CXVIRTUALSCREEN)
	height := win.GetSystemMetrics(win.SM_CYVIRTUALSCREEN)
	if width <= 0 || height <= 0 {
		left = 0
		top = 0
		width = win.GetSystemMetrics(win.SM_CXSCREEN)
		height = win.GetSystemMetrics(win.SM_CYSCREEN)
	}
	if width <= 0 || height <= 0 {
		return workArea()
	}
	return win.RECT{Left: left, Top: top, Right: left + width, Bottom: top + height}
}

func scaleImage(src *image.RGBA, factor int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx()*factor, src.Bounds().Dy()*factor))
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			c := src.RGBAAt(x, y)
			for sy := 0; sy < factor; sy++ {
				for sx := 0; sx < factor; sx++ {
					dst.SetRGBA(x*factor+sx, y*factor+sy, c)
				}
			}
		}
	}
	return dst
}

func scaleImageTo(src *image.RGBA, w, h int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	sw := max(1, src.Bounds().Dx())
	sh := max(1, src.Bounds().Dy())
	for y := 0; y < h; y++ {
		sy := src.Bounds().Min.Y + y*sh/h
		for x := 0; x < w; x++ {
			sx := src.Bounds().Min.X + x*sw/w
			dst.SetRGBA(x, y, src.RGBAAt(sx, sy))
		}
	}
	return dst
}

func fitVisibleImageTo(src *image.RGBA, w, h int) *image.RGBA {
	content := alphaBounds(src)
	if content.Empty() {
		return scaleImageTo(src, w, h)
	}
	cropped := image.NewRGBA(image.Rect(0, 0, content.Dx(), content.Dy()))
	draw.Draw(cropped, cropped.Bounds(), src, content.Min, draw.Src)
	scale := math.Min(float64(w)/float64(content.Dx()), float64(h)/float64(content.Dy()))
	outW := max(1, int(math.Round(float64(content.Dx())*scale)))
	outH := max(1, int(math.Round(float64(content.Dy())*scale)))
	scaled := scaleImageTo(cropped, outW, outH)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	offX := (w - outW) / 2
	offY := h - outH
	draw.Draw(dst, image.Rect(offX, offY, offX+outW, offY+outH), scaled, image.Point{}, draw.Over)
	return dst
}

func alphaBounds(img *image.RGBA) image.Rectangle {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y
	found := false
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.RGBAAt(x, y).A == 0 {
				continue
			}
			found = true
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
	if !found {
		return image.Rect(0, 0, 0, 0)
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func drawWheelBack(dst *image.RGBA, x, y int, wheel *image.RGBA) {
	cx := x + wheelSize/2
	cy := y + wheelSize/2
	outer := float64(wheelSize/2 - 2)
	inner := outer - 5
	rim := rgba(92, 86, 76, 210)
	shadow := rgba(44, 41, 38, 120)
	base := rgba(74, 67, 58, 210)

	if wheel != nil {
		draw.Draw(dst, image.Rect(x, y, x+wheelSize, y+wheelSize), wheel, image.Point{}, draw.Over)
		return
	}
	for py := y; py < y+wheelSize; py++ {
		for px := x; px < x+wheelSize; px++ {
			dx := float64(px - cx)
			dy := float64(py - cy)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= outer+1 && dist >= inner {
				if px < cx && py > cy {
					dst.SetRGBA(px, py, shadow)
				} else {
					dst.SetRGBA(px, py, rim)
				}
			}
		}
	}
	drawPixelLine(dst, cx-20, y+wheelSize-2, cx-30, sceneH-2, base)
	drawPixelLine(dst, cx+20, y+wheelSize-2, cx+30, sceneH-2, base)
	for px := cx - 38; px <= cx+38; px++ {
		for py := sceneH - 4; py <= sceneH-2; py++ {
			if image.Pt(px, py).In(dst.Bounds()) {
				dst.SetRGBA(px, py, base)
			}
		}
	}
}

func drawWheelRunner(dst *image.RGBA, x, y int, src *image.RGBA, frame int, dir int) {
	runnerW := 58
	runnerH := 40
	scaled := fitVisibleImageTo(src, runnerW, runnerH)
	bob := int(math.Sin(float64(frame)/2.0) * 2)
	dstX := x + (wheelSize-runnerW)/2
	dstY := y + wheelSize/2 - runnerH/2 + 3 + bob
	drawFacingImage(dst, scaled, image.Rect(dstX, dstY, dstX+runnerW, dstY+runnerH), dir)
}

func drawWheelFront(dst *image.RGBA, x, y, tick int) {
	cx := x + wheelSize/2
	cy := y + wheelSize/2
	inner := float64(wheelSize/2 - 7)
	spoke := rgba(132, 123, 106, 115)
	hub := rgba(86, 78, 68, 230)
	rim := rgba(92, 86, 76, 160)

	angle := float64(tick) * 0.34
	for i := 0; i < 8; i++ {
		a := angle + float64(i)*math.Pi/4
		x2 := cx + int(math.Cos(a)*(inner-2))
		y2 := cy + int(math.Sin(a)*(inner-2))
		drawThinLine(dst, cx, cy, x2, y2, spoke)
	}

	fillCircle(dst, cx, cy, 4, hub)
	for py := y; py < y+wheelSize; py++ {
		for px := x; px < x+wheelSize; px++ {
			dx := float64(px - cx)
			dy := float64(py - cy)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist >= float64(wheelSize/2-5) && dist <= float64(wheelSize/2-1) {
				dst.SetRGBA(px, py, rim)
			}
		}
	}
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
			if image.Pt(x+col, y+row).In(dst.Bounds()) {
				dst.SetRGBA(x+col, y+row, overRGBA(dst.RGBAAt(x+col, y+row), c))
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
			if image.Pt(px, y+row).In(dst.Bounds()) {
				dst.SetRGBA(px, y+row, overRGBA(dst.RGBAAt(px, y+row), c))
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
			if image.Pt(px, py+sy).In(dst.Bounds()) {
				dst.SetRGBA(px, py+sy, overRGBA(dst.RGBAAt(px, py+sy), ink))
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
			if image.Pt(px, py).In(dst.Bounds()) {
				dst.SetRGBA(px, py, overRGBA(dst.RGBAAt(px, py), c))
			}
		}
	}
}

func drawRoundedRect(dst *image.RGBA, x, y, w, h, radius int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if roundedRectContains(px-x, py-y, w, h, radius) && image.Pt(px, py).In(dst.Bounds()) {
				dst.SetRGBA(px, py, overRGBA(dst.RGBAAt(px, py), c))
			}
		}
	}
}

func drawRoundedRectOutline(dst *image.RGBA, x, y, w, h, radius int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if !image.Pt(px, py).In(dst.Bounds()) || !roundedRectContains(px-x, py-y, w, h, radius) {
				continue
			}
			if px == x || px == x+w-1 || py == y || py == y+h-1 ||
				!roundedRectContains(px-x-1, py-y, w, h, radius) ||
				!roundedRectContains(px-x+1, py-y, w, h, radius) ||
				!roundedRectContains(px-x, py-y-1, w, h, radius) ||
				!roundedRectContains(px-x, py-y+1, w, h, radius) {
				dst.SetRGBA(px, py, overRGBA(dst.RGBAAt(px, py), c))
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

func drawThinLine(dst *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
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
		if image.Pt(x0, y0).In(dst.Bounds()) {
			dst.SetRGBA(x0, y0, c)
		}
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

func rgba(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: uint8(uint16(r) * uint16(a) / 255), G: uint8(uint16(g) * uint16(a) / 255), B: uint8(uint16(b) * uint16(a) / 255), A: a}
}

func drawBlock(dst *image.RGBA, x, y int, c color.RGBA) {
	for py := y - 1; py <= y+1; py++ {
		for px := x - 1; px <= x+1; px++ {
			if image.Pt(px, py).In(dst.Bounds()) {
				dst.SetRGBA(px, py, c)
			}
		}
	}
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

func normalizeDir(dir int) int {
	if dir < 0 {
		return -1
	}
	return 1
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func normalizeDPI(dpi int) int {
	if dpi <= 0 {
		return defaultDPI
	}
	return clamp(dpi, 72, 384)
}

func scaleForDPI(v int, dpi int) int {
	dpi = normalizeDPI(dpi)
	if v == 0 || dpi == defaultDPI {
		return v
	}
	sign := 1
	if v < 0 {
		sign = -1
		v = -v
	}
	return sign * max(1, int(math.Round(float64(v)*float64(dpi)/float64(defaultDPI))))
}

func unscaleForDPI(v int, dpi int) int {
	dpi = normalizeDPI(dpi)
	if v == 0 || dpi == defaultDPI {
		return v
	}
	sign := 1
	if v < 0 {
		sign = -1
		v = -v
	}
	return sign * max(1, int(math.Round(float64(v)*float64(defaultDPI)/float64(dpi))))
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

func max32(a, b int32) int32 {
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
