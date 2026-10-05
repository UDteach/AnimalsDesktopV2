// motionpreview renders an offline, actual-size preview using the same player
// and painter as the native application. It is not a performance benchmark.
package main

import (
	appassets "animals-desktop/assets"
	"animals-desktop/internal/motion"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
)

func main() {
	out := flag.String("out", "docs/motions/ver2/integration/qa/final-preview", "output directory")
	assetRoot := flag.String("asset-root", "", "optional read-only candidate filesystem root; defaults to embedded runtime assets")
	manifestPath := flag.String("manifest", "motions/manifest.json", "manifest path within asset root")
	heightFlag := flag.Int("height", 64, "display canvas height for calibration")
	frameCount := flag.Int("frames", 300, "number of 40ms frames (1..3000); extend for long actions to complete twice")
	actionFlag := flag.String("action", "", "request an adopted action at 5s and 10s in this offline review")
	stationary := flag.Bool("stationary", false, "hold target speed at zero and record idle/action phases; does not certify gait")
	traceFlag := flag.Bool("trace", false, "record gait/idle samples without requesting an action")
	flag.Parse()
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binary, err := os.Open(executable)
	if err != nil {
		panic(err)
	}
	binaryDigest := sha256.New()
	_, copyErr := io.Copy(binaryDigest, binary)
	closeErr := binary.Close()
	if copyErr != nil {
		panic(copyErr)
	}
	if closeErr != nil {
		panic(closeErr)
	}
	if *heightFlag < 32 || *heightFlag > 256 {
		panic("height must be 32..256")
	}
	if *frameCount < 1 || *frameCount > 3000 {
		panic("frames must be 1..3000")
	}
	height := *heightFlag
	width := height * 3 / 2
	cell := height * 5 / 2
	if err := os.MkdirAll(*out, 0755); err != nil {
		panic(err)
	}
	var files fs.FS = appassets.FS
	if *assetRoot != "" {
		files = os.DirFS(*assetRoot)
	}
	s, err := motion.NewStore(files, *manifestPath)
	if err != nil {
		panic(err)
	}
	data, err := fs.ReadFile(files, *manifestPath)
	if err != nil {
		panic(err)
	}
	manifestSHA := fmt.Sprintf("%x", sha256.Sum256(data))
	var manifest motion.Manifest
	if err = json.Unmarshal(data, &manifest); err != nil {
		panic(err)
	}
	players := make([]motion.Player, len(manifest.Variants))
	variants := make([]*motion.Loaded, len(players))
	for i, v := range manifest.Variants {
		variants[i] = s.Get(v.ID)
		if variants[i] == nil {
			panic(s.Error(v.ID))
		}
		players[i].Reset(v.ID, float64(i*cell+12), i%2 == 1)
	}
	rowHeight, baseline := previewRows(variants, height)
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	var actionTrace []map[string]any
	traceEnabled := *actionFlag != "" || *stationary || *traceFlag
	frameHashes := make(map[string]string, *frameCount)
	for frame := 0; frame < *frameCount; frame++ {
		seconds := float64(frame) / 25
		canvas := image.NewRGBA(image.Rect(0, 0, len(players)*cell, rowHeight*2))
		draw.Draw(canvas, image.Rect(0, 0, canvas.Bounds().Dx(), rowHeight), image.NewUniform(color.RGBA{245, 245, 245, 255}), image.Point{}, draw.Src)
		draw.Draw(canvas, image.Rect(0, rowHeight, canvas.Bounds().Dx(), rowHeight*2), image.NewUniform(color.RGBA{24, 29, 37, 255}), image.Point{}, draw.Src)
		for _, ground := range []int{baseline, rowHeight + baseline} {
			for x := 0; x < canvas.Bounds().Dx(); x++ {
				canvas.SetRGBA(x, ground, color.RGBA{128, 128, 128, 255})
			}
		}
		for i, v := range variants {
			p := &players[i]
			requested := false
			if *actionFlag != "" && (frame == 125 || frame == 250) {
				requested = p.RequestAction(v, *actionFlag)
			}
			target := previewTarget(v, height, seconds, *stationary)
			p.Advance(v, .04, target, float64(height), float64(i*cell+4), float64((i+1)*cell-width-4), seconds >= 8 && seconds < 9)
			if traceEnabled {
				actionTrace = append(actionTrace, map[string]any{"seconds": seconds, "id": v.ID, "x": p.X, "velocity": p.Velocity,
					"clip": p.Action, "phase": p.Phase(), "sequence": p.Sequence, "pending": p.PendingAction, "requestAccepted": requested,
					"left": p.Left, "drawX": int(p.X), "frameIndex": motion.FrameAt(p.Phase(), v.Motions[p.Action].Durations),
					"walkPhase": p.WalkPhase, "targetSpeed": target, "paused": seconds >= 8 && seconds < 9})
			}
			dir := 1
			if p.Left {
				dir = -1
			}
			v.Draw(canvas, p.Action, p.Phase(), int(p.X), baseline, width, height, dir)
			v.Draw(canvas, p.Action, p.Phase(), int(p.X), rowHeight+baseline, width, height, dir)
		}
		f, err := os.Create(filepath.Join(*out, fmt.Sprintf("%04d.png", frame)))
		if err != nil {
			panic(err)
		}
		digest := sha256.New()
		if err = encoder.Encode(io.MultiWriter(f, digest), canvas); err != nil {
			panic(err)
		}
		if err = f.Close(); err != nil {
			panic(err)
		}
		frameHashes[fmt.Sprintf("%04d.png", frame)] = fmt.Sprintf("%x", digest.Sum(nil))
	}
	meta := map[string]any{"schema": 2, "manifestSha256": manifestSHA, "exportedFrameSha256": frameHashes, "frames": *frameCount, "frameDurationMs": 40, "render": "offline native shared Go painter/player", "height": height, "width": width, "sequence": "walk 0-4s, deceleration/idle 4-6s, short faster walk 6-6.6s, walk, pause 8-9s, resume", "variants": manifest.Variants}
	meta["stationaryScenario"] = *stationary
	meta["traceEnabled"] = traceEnabled
	meta["executableSha256"] = fmt.Sprintf("%x", binaryDigest.Sum(nil))
	meta["rowHeight"] = rowHeight
	meta["baselineWithinRow"] = baseline
	meta["verticalFraming"] = "all loaded clips' full transformed canvases plus padding; no animal rescale or ground shift between poses"
	if *stationary {
		meta["sequence"] = "zero target speed for idle/action review; pause8-9s then resume; gait untested"
	}
	if traceEnabled {
		data, err = json.MarshalIndent(map[string]any{"schema": 2, "manifestSha256": manifestSHA, "height": height, "frameDurationMs": 40, "requestedAction": *actionFlag, "trace": actionTrace}, "", "  ")
		if err != nil {
			panic(err)
		}
		if err = os.WriteFile(filepath.Join(*out, "action-trace.json"), data, 0644); err != nil {
			panic(err)
		}
		meta["actionTraceSha256"] = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	data, err = json.MarshalIndent(meta, "", "  ")
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(filepath.Join(*out, "preview.json"), data, 0644); err != nil {
		panic(err)
	}
}

// Keep both background rows independent even when a tail extends below the
// support baseline, or a whole-clip presentation rises above the usual canvas.
// Full transformed canvases are conservative bounds; using feet or only idle
// alpha bounds would clip other poses. Baseline stays identical for all clips.
func previewRows(variants []*motion.Loaded, height int) (rowHeight, baseline int) {
	minY, maxY := 0, 0
	for _, v := range variants {
		for action := range v.Frames {
			p := v.Motions[action].Presentation
			top := float64(height) * (p.Y - v.Profile.Ground)
			bottom := top + float64(height)*p.Scale
			minY = min(minY, int(math.Floor(top)))
			maxY = max(maxY, int(math.Ceil(bottom)))
		}
	}
	baseline = max(height+22, -minY+4)
	rowHeight = max(height+28, baseline+maxY+4)
	return rowHeight, baseline
}

func previewTarget(v *motion.Loaded, height int, seconds float64, stationary bool) float64 {
	if stationary || seconds >= 4 && seconds < 6 {
		return 0
	}
	target := v.Profile.Stride * float64(height) / (motion.Duration(v.Motions["walk"].Durations) / 1000)
	if seconds >= 6 && seconds < 6.6 {
		target *= 1.35
	}
	return target
}
