//go:build windows && motionqa

package main

// A separate local-QA build exercises the real software renderer and Win32
// layered window without reading/writing settings, hooks, network or tray state.
import (
	appassets "animals-desktop/assets"
	"animals-desktop/internal/motion"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lxn/win"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var qaApp *petApp
var qaDir string
var qaStart time.Time
var qaClock motion.Clock
var qaTicks int
var qaLimit = 20.0
var qaFailures []string
var qaNativeDraws int
var qaTimes []float64
var qaFinishSeconds float64
var qaTickMilliseconds []float64
var qaCount = 10
var qaAction string
var qaCandidateManifest string
var qaManifestSHA string
var qaEmbeddedRuntime bool
var qaActionRequests = map[int]bool{}
var qaActionTransitions []map[string]any
var qaActionStarts = map[int]int{}
var qaActionReturns = map[int]int{}
var qaResizeDuringAction bool
var qaResizeCounts = map[int]int{}
var qaResizePauseStarts = map[int]float64{}
var qaResizeEvents []map[string]any
var qaResizePercents = []int{70, 100, 110, 120}

type qaSample struct {
	canvas *image.RGBA
	tick   int
}

var qaSamples chan qaSample
var qaWriter sync.WaitGroup

func writeQAPNG(path string, canvas image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	encodeErr := png.Encode(f, canvas)
	return errors.Join(encodeErr, f.Close())
}

func actionCoverageFailures(count int, requests map[int]bool, starts, returns map[int]int) []string {
	var failures []string
	for pet := 0; pet < count; pet++ {
		accepted := false
		for key, ok := range requests {
			if key%10 == pet && ok {
				accepted = true
			}
		}
		if !accepted || starts[pet] == 0 || returns[pet] == 0 {
			failures = append(failures, fmt.Sprintf("pet %d did not complete an accepted action through entry and return to idle", pet))
		}
	}
	return failures
}

func resizeCoverageFailures(count int, counts map[int]int) []string {
	var failures []string
	for pet := 0; pet < count; pet++ {
		if counts[pet] != len(qaResizePercents)*2 {
			failures = append(failures, fmt.Sprintf("pet %d missed paused/playing action resizes: %d/%d", pet, counts[pet], len(qaResizePercents)*2))
		}
	}
	return failures
}

func qaCaptureFilename(tick int) string {
	if tick == 100 {
		return "render.png"
	}
	return filepath.Join("frames", fmt.Sprintf("%06d.png", tick))
}

func resizeCaptureFailures(events []map[string]any, saved map[int]bool) []string {
	var failures []string
	for _, event := range events {
		tick := event["tick"].(int)
		event["capture"] = qaCaptureFilename(tick)
		event["captureSaved"] = saved[tick]
		if !saved[tick] {
			failures = append(failures, fmt.Sprintf("resize at tick %d has no successful PNG capture", tick))
		}
	}
	return failures
}

// Start each pet's pause after it actually enters its action. A fixed wall-clock
// window can expire while a faster/larger pet is still braking into the action.
func resizePauseState(starts map[int]float64, pet, stage int, elapsed float64, active bool) (bool, float64) {
	if !active || stage < 0 || stage >= len(qaResizePercents) {
		return false, 0
	}
	start, seen := starts[pet]
	if !seen {
		start = elapsed
		starts[pet] = start
	}
	return true, elapsed - start
}

func resizeStageDue(stage int, phase, sequencePhase float64, paused bool, pauseElapsed float64) bool {
	if stage < len(qaResizePercents) {
		return stage >= 0 && paused && pauseElapsed >= .1+float64(stage)*.2
	}
	return stage < len(qaResizePercents)*2 && !paused && phase >= 13 && sequencePhase >= float64(stage-len(qaResizePercents))*80
}

func qaSolidEnvelope(v *motion.Loaded, p *desktopPet, width, height int) [4]int {
	canvas := image.NewRGBA(image.Rect(0, 0, qaApp.sceneW, sceneH))
	ground := sceneH - p.laneOffset - 3
	v.Draw(canvas, p.timed.Action, p.timed.Phase(), p.x, ground, width, height, p.dir)
	box := image.Rectangle{Min: canvas.Bounds().Max}
	found := false
	for y := 0; y < canvas.Bounds().Dy(); y++ {
		for x := 0; x < canvas.Bounds().Dx(); x++ {
			if canvas.RGBAAt(x, y).A >= 180 {
				found = true
				box.Min.X, box.Min.Y = min(box.Min.X, x), min(box.Min.Y, y)
				box.Max.X, box.Max.Y = max(box.Max.X, x+1), max(box.Max.Y, y+1)
			}
		}
	}
	if !found {
		qaFailures = append(qaFailures, "resize produced empty solid envelope")
		return [4]int{}
	}
	// Screen X and Y relative to the same floor; this is not a named-foot label.
	return [4]int{box.Min.X, box.Min.Y - ground, box.Max.X, box.Max.Y - ground}
}

func runMotionQA(args []string) bool {
	if len(args) < 2 || args[0] != "--motion-qa" {
		return false
	}
	if args[len(args)-1] == "--resize-during-action" {
		qaResizeDuringAction = true
		args = args[:len(args)-1]
	}
	qaDir = args[1]
	if len(args) > 2 {
		qaLimit, _ = strconv.ParseFloat(args[2], 64)
	}
	if len(args) > 3 {
		qaCount, _ = strconv.Atoi(args[3])
		if qaCount < 1 || qaCount > 10 {
			qaCount = 10
		}
	}
	if qaLimit < 5 {
		qaLimit = 20
	}
	if err := os.MkdirAll(filepath.Join(qaDir, "frames"), 0755); err != nil {
		panic(err)
	}
	enablePerMonitorDPIAwareness()
	runtime.LockOSThread()
	qaApp = &petApp{frames: newSpriteCache(), sceneW: 1500, petSizes: defaultPetSizes(), speed: 3, mode: modeKeyboard, wheelEnabled: true}
	ids := []string{"chinchilla_standard_gray", "chinchilla_beige", "chinchilla_ebony", "hamster_golden_syrian", "ferret_sable", "chipmunk_striped", "macaroni_mouse_tan", "rabbit_gray", "quokka", "white_wagtail"}[:qaCount]
	// Optional candidate files are available only in this separate QA build:
	// --motion-qa OUT SECONDS COUNT ASSET_ROOT MANIFEST [ACTION] [--resize-during-action]
	// ASSET_ROOT=@embedded verifies the actual compiled runtime bundle and
	// selects only variants that have the requested adopted action.
	if len(args) > 4 {
		if len(args) < 6 || len(args) > 7 {
			panic("candidate QA needs ASSET_ROOT MANIFEST and optional ACTION")
		}
		if os.Getenv("ANIMALSDESKTOP_MOTIONS") == "legacy" {
			panic("candidate QA cannot run with the legacy-only override")
		}
		var files fs.FS = os.DirFS(args[4])
		qaEmbeddedRuntime = args[4] == "@embedded"
		if qaEmbeddedRuntime {
			files = appassets.FS
		}
		qaCandidateManifest = args[5]
		if len(args) == 7 {
			qaAction = args[6]
		}
		store, err := motion.NewStore(files, qaCandidateManifest)
		if err != nil {
			panic(err)
		}
		raw, err := fs.ReadFile(files, qaCandidateManifest)
		if err != nil {
			panic(err)
		}
		qaManifestSHA = fmt.Sprintf("%x", sha256.Sum256(raw))
		var manifest motion.Manifest
		if err = json.Unmarshal(raw, &manifest); err != nil {
			panic(err)
		}
		var available []string
		for _, entry := range manifest.Variants {
			if !entry.Enabled {
				continue
			}
			if _, ok := variantIndexByID(entry.ID); !ok {
				panic("QA candidate is not a public animal ID")
			}
			v := store.Get(entry.ID)
			if v == nil {
				panic("QA candidate action is missing or failed integrity")
			}
			if qaAction != "" && !v.HasAction(qaAction) {
				if qaEmbeddedRuntime {
					continue
				}
				panic("QA candidate action is missing or failed integrity")
			}
			available = append(available, entry.ID)
		}
		if len(available) == 0 {
			panic("no enabled QA candidates")
		}
		timedStore = store
		timedOnce.Do(func() {})
		ids = make([]string, qaCount)
		for i := range ids {
			ids[i] = available[i%len(available)]
		}
	}
	if qaResizeDuringAction && qaAction == "" {
		panic("resize-during-action QA requires an explicit candidate ACTION")
	}
	for i, id := range ids {
		index := 0
		for n, v := range variants {
			if v.ID == id {
				index = n
				break
			}
		}
		qaApp.petSizes[i] = []int{70, 100, 120}[i%3]
		if qaResizeDuringAction {
			qaApp.petSizes[i] = qaResizePercents[i%len(qaResizePercents)]
		}
		p := desktopPet{variant: index, x: 20 + i*150, dir: 1, item: noItem, carryKind: noItem, state: stateWalk, stateTicks: 100000, moveSpeed: 1}
		if i%2 == 1 {
			p.dir = -1
		}
		if v := timedVariant(id); v != nil {
			p.timed.Reset(id, float64(p.x), p.dir < 0)
		}
		qaApp.pets = append(qaApp.pets, p)
		if timedVariant(id) == nil {
			qaApp.frames.variantSets(variants[index])
		}
	}
	hinst := win.GetModuleHandle(nil)
	class := syscall.StringToUTF16Ptr("AnimalsDesktopMotionQA")
	wc := win.WNDCLASSEX{CbSize: uint32(unsafe.Sizeof(win.WNDCLASSEX{})), LpfnWndProc: syscall.NewCallback(qaWndProc), HInstance: hinst, LpszClassName: class}
	if win.RegisterClassEx(&wc) == 0 {
		panic(syscall.GetLastError())
	}
	hwnd := win.CreateWindowEx(win.WS_EX_LAYERED|win.WS_EX_TOPMOST|win.WS_EX_TOOLWINDOW|win.WS_EX_TRANSPARENT, class, syscall.StringToUTF16Ptr("AnimalsDesktop motion QA"), win.WS_POPUP, 60, 180, 1500, sceneH*2+16, 0, 0, hinst, nil)
	if hwnd == 0 {
		panic(syscall.GetLastError())
	}
	defer win.DestroyWindow(hwnd)
	qaSamples = make(chan qaSample, 4)
	var captureErrors []string // Only the writer accesses these until Wait returns.
	capturedFrames := 0
	capturedTicks := map[int]bool{} // Writer acknowledgements, read after Wait.
	qaWriter.Add(1)
	go func() {
		defer qaWriter.Done()
		for sample := range qaSamples {
			dest := filepath.Join(qaDir, qaCaptureFilename(sample.tick))
			if err := writeQAPNG(dest, sample.canvas); err != nil {
				captureErrors = append(captureErrors, fmt.Sprintf("capture %s: %v", qaCaptureFilename(sample.tick), err))
			} else {
				capturedFrames++
				capturedTicks[sample.tick] = true
			}
		}
	}()
	qaStart = time.Now()
	win.SetTimer(hwnd, 42, 16, 0)
	win.ShowWindow(hwnd, win.SW_SHOWNOACTIVATE)
	var msg win.MSG
	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

	close(qaSamples)
	qaWriter.Wait()
	qaFailures = append(qaFailures, captureErrors...)
	if capturedFrames == 0 {
		qaFailures = append(qaFailures, "no PNG evidence was saved")
	}
	if qaAction != "" {
		qaFailures = append(qaFailures, actionCoverageFailures(qaCount, qaActionRequests, qaActionStarts, qaActionReturns)...)
	}
	if qaResizeDuringAction {
		qaFailures = append(qaFailures, resizeCoverageFailures(qaCount, qaResizeCounts)...)
		qaFailures = append(qaFailures, resizeCaptureFailures(qaResizeEvents, capturedTicks)...)
	}
	report := map[string]any{"elapsedSeconds": qaFinishSeconds, "ticks": qaTicks, "nativeDraws": qaNativeDraws, "failures": qaFailures, "pets": qaCount, "sampleTimes": qaTimes, "tickMilliseconds": qaTickMilliseconds, "sizes": []int{70, 100, 120}, "notes": "Actual Win32 layered window; real renderOverlaySegment; PNG samples with measured timestamps. QA mode does not access user settings. Mac unverified."}
	report["capturedFrames"] = capturedFrames
	if qaResizeDuringAction {
		report["resizeEvents"] = qaResizeEvents
		report["sizes"] = qaResizePercents
		report["resizeScope"] = "Four resizes while action paused and four while playing, via real setPetSize and shared tick reconciliation; UI percentages 70/100/110/120 include the 64-to-96 tier seam. Center placement; left-X origin convention, not fixed world paw X or edge-clamp acceptance. Solid envelopes are screening only. Not physical mixed-DPI hardware verification."
	}
	if qaCandidateManifest != "" {
		report["sourceManifestSha256"] = qaManifestSHA
		if qaEmbeddedRuntime {
			report["embeddedManifest"] = qaCandidateManifest
			report["embeddedRuntimeAssets"] = true
			report["runtimeAdopted"] = true
		} else {
			report["candidateManifest"] = qaCandidateManifest
			report["runtimeAdopted"] = false
		}
	}
	if qaAction != "" {
		report["requestedAction"] = qaAction
		report["actionRequests"] = qaActionRequests
		report["actionTransitions"] = qaActionTransitions
		report["actionStarts"] = qaActionStarts
		report["actionReturnsToIdle"] = qaActionReturns
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(qaDir, "result.json"), data, 0644); err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	if len(qaFailures) != 0 {
		panic("motion QA failed; see result.json")
	}
	return true
}

func qaWndProc(hwnd win.HWND, msg uint32, w, l uintptr) uintptr {
	if msg == win.WM_TIMER {
		now := time.Now()
		elapsed := now.Sub(qaStart).Seconds()
		dt := qaClock.Step(now)
		if elapsed > qaLimit {
			qaFinishSeconds = elapsed
			win.KillTimer(hwnd, 42)
			win.PostQuitMessage(0)
			return 0
		}
		// Walk -> decelerate -> stationary -> short input reaction -> pause -> resume.
		phase := math.Mod(elapsed, 16)
		resizedThisTick := false
		for i := range qaApp.pets {
			p := &qaApp.pets[i]
			v := timedVariant(variants[p.variant].ID)
			var resizeEvent map[string]any
			if v != nil {
				paused := phase >= 10 && phase < 12 || !qaResizeDuringAction && qaAction != "" && phase >= 5.4 && phase < 6.4
				pauseElapsed := 0.0
				if qaResizeDuringAction {
					var resizePaused bool
					resizePaused, pauseElapsed = resizePauseState(qaResizePauseStarts, i, qaResizeCounts[i], elapsed, p.timed.Sequence == qaAction)
					paused = paused || resizePaused
				}
				if qaResizeDuringAction && p.timed.Sequence == qaAction && resizeStageDue(qaResizeCounts[i], phase, p.timed.SequencePhase, paused, pauseElapsed) {
					beforeResize := p.timed
					from := qaApp.petSizePercent(i)
					stage := qaResizeCounts[i]
					to := qaResizePercents[(i+stage+1)%len(qaResizePercents)]
					oldWidth, oldHeight := qaApp.petSpriteSize(i)
					oldEnvelope := qaSolidEnvelope(v, p, oldWidth, oldHeight)
					qaApp.setPetSize(i, to)
					width, height := qaApp.petSpriteSize(i)
					// Use the same platform reconciliation while holding time still.
					tickTimedPlayer(&p.timed, v, 0, p.x, height, max(0, qaApp.sceneW-width), p.dir, 1, false, paused)
					if p.timed != beforeResize || from == qaApp.petSizePercent(i) {
						qaFailures = append(qaFailures, "resize reset action or did not change size")
					}
					resizeEvent = map[string]any{"seconds": elapsed, "pet": i, "id": v.ID,
						"fromPercent": from, "toPercent": qaApp.petSizePercent(i), "width": width, "height": height,
						"paused": paused, "sequence": p.timed.Sequence, "clip": p.timed.Action, "phase": p.timed.Phase(),
						"resizePauseElapsed": pauseElapsed,
						"frame":              motion.FrameAt(p.timed.Phase(), v.Motions[p.timed.Action].Durations), "x": p.timed.X,
						"oldSolidEnvelope": oldEnvelope, "newSolidEnvelope": qaSolidEnvelope(v, p, width, height), "tick": qaTicks}
					qaResizeEvents = append(qaResizeEvents, resizeEvent)
					qaResizeCounts[i]++
					resizedThisTick = true
				}
				_, h := qaApp.petSpriteSize(i)
				target := v.Profile.Stride * float64(h) / (motion.Duration(v.Motions["walk"].Durations) / 1000) * []float64{.65, 1, 1.65}[i%3]
				if phase > 4 && phase < 8 {
					target = 0
				}
				if phase >= 8 && phase < 8.55 {
					target *= 1.35
				}
				if qaAction != "" && phase >= 5 {
					slot := int(elapsed/16) * 2
					if phase >= 13 {
						slot++
					}
					key := slot*10 + i
					if _, seen := qaActionRequests[key]; !seen {
						qaActionRequests[key] = p.timed.RequestAction(v, qaAction)
					}
				}
				before := p.timed
				p.timed.Advance(v, dt, target, float64(h), float64(i*150+5), float64(i*150+40), paused)
				if paused && p.timed != before {
					qaFailures = append(qaFailures, "pause moved")
				}
				if qaAction != "" {
					if before.Sequence == "" && p.timed.Sequence == qaAction {
						qaActionStarts[i]++
					}
					if before.Sequence == qaAction && p.timed.Sequence == "" && p.timed.Action == "idle" {
						qaActionReturns[i]++
					}
					if before.Sequence != "" && before.X != p.timed.X {
						qaFailures = append(qaFailures, "stationary action moved")
					}
					if before.Sequence != p.timed.Sequence || before.Action != p.timed.Action {
						qaActionTransitions = append(qaActionTransitions, map[string]any{"seconds": elapsed, "pet": i, "id": v.ID,
							"x": p.timed.X, "velocity": p.timed.Velocity, "clip": p.timed.Action, "sequence": p.timed.Sequence, "phase": p.timed.Phase()})
					}
				}
				p.x = int(math.Round(p.timed.X))
				p.dir = 1
				if p.timed.Left {
					p.dir = -1
				}
				p.state = stateIdle
				if p.timed.Action == "walk" {
					p.state = stateWalk
				}
				if resizeEvent != nil {
					// The saved canvas is drawn after this tick advances. Keep that
					// phase separate from the zero-time resize comparison above.
					resizeEvent["capturePhase"] = p.timed.Phase()
					resizeEvent["captureClip"] = p.timed.Action
					resizeEvent["captureSequence"] = p.timed.Sequence
					resizeEvent["captureFrame"] = motion.FrameAt(p.timed.Phase(), v.Motions[p.timed.Action].Durations)
					resizeEvent["captureX"] = p.x
				}
			} else if qaTicks%4 == 0 {
				qaApp.tickPet(i, p)
			}
		}
		segment := overlaySegment{Rect: win.RECT{Right: 1500, Bottom: sceneH}, SceneRight: 1500, DPI: 96}
		layer := qaApp.renderOverlaySegment(segment)
		canvas := image.NewRGBA(image.Rect(0, 0, 1500, sceneH*2+16))
		draw.Draw(canvas, image.Rect(0, 0, 1500, sceneH+8), image.NewUniform(color.RGBA{245, 245, 245, 255}), image.Point{}, draw.Src)
		draw.Draw(canvas, image.Rect(0, sceneH+8, 1500, sceneH*2+16), image.NewUniform(color.RGBA{24, 29, 37, 255}), image.Point{}, draw.Src)
		draw.Draw(canvas, layer.Bounds(), layer, image.Point{}, draw.Over)
		draw.Draw(canvas, layer.Bounds().Add(image.Pt(0, sceneH+8)), layer, image.Point{}, draw.Over)
		if updateLayeredWindow(hwnd, canvas, 60, 180) {
			qaNativeDraws++
		} else {
			qaFailures = append(qaFailures, "UpdateLayeredWindow failed")
		}
		if qaTicks == 100 {
			select {
			case qaSamples <- qaSample{canvas, qaTicks}:
			default:
			}
		}
		if (qaTicks%4 == 0 || resizedThisTick) && elapsed < 16 {
			select {
			case qaSamples <- qaSample{canvas, qaTicks}:
				qaTimes = append(qaTimes, elapsed)
			default:
			}
		}

		qaTickMilliseconds = append(qaTickMilliseconds, time.Since(now).Seconds()*1000)
		qaTicks++
		return 0
	}
	return win.DefWindowProc(hwnd, msg, w, l)
}
