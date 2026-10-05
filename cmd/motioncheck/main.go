// motioncheck measures every frame through the native app's shared painter.
// Alpha envelopes are screening evidence, not anatomical body/foot labels.
package main

import (
	appassets "animals-desktop/assets"
	"animals-desktop/internal/motion"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type measurement struct {
	ID               string   `json:"id"`
	Action           string   `json:"action"`
	Frame            int      `json:"frame"`
	Height           int      `json:"height"`
	Direction        int      `json:"direction"`
	TimeMS           float64  `json:"timeMs"`
	SolidBounds      [4]int   `json:"solidBoundsRelativeToOriginGround"`
	AlphaBounds      [4]int   `json:"alphaBoundsRelativeToOriginGround"`
	SolidArea        int      `json:"solidArea"`
	BottomComponents [][2]int `json:"bottomComponentsRelativeX"`
	CanvasClipped    bool     `json:"canvasClipped"`
}

func bounds(im *image.RGBA, threshold uint8) (image.Rectangle, int) {
	r := image.Rectangle{Min: im.Bounds().Max}
	area := 0
	for y := 0; y < im.Bounds().Max.Y; y++ {
		for x := 0; x < im.Bounds().Max.X; x++ {
			if im.RGBAAt(x, y).A >= threshold {
				r.Min.X = min(r.Min.X, x)
				r.Min.Y = min(r.Min.Y, y)
				r.Max.X = max(r.Max.X, x+1)
				r.Max.Y = max(r.Max.Y, y+1)
				area++
			}
		}
	}
	if area == 0 {
		return image.Rectangle{}, 0
	}
	return r, area
}

func relative(r image.Rectangle, x, ground int) [4]int {
	return [4]int{r.Min.X - x, r.Min.Y - ground, r.Max.X - x, r.Max.Y - ground}
}

func main() {
	out := flag.String("out", "docs/motions/ver2/integration/qa/geometry/runtime", "output")
	root := flag.String("asset-root", "", "candidate root; default embedded runtime assets")
	manifestPath := flag.String("manifest", "motions/manifest.json", "manifest within root")
	exportFrames := flag.Bool("export-frames", false, "save every measured frame for named-foot and body review")
	flag.Parse()
	var files fs.FS = appassets.FS
	if *root != "" {
		files = os.DirFS(*root)
	}
	raw, err := fs.ReadFile(files, *manifestPath)
	must(err)
	var manifest motion.Manifest
	must(json.Unmarshal(raw, &manifest))
	store, err := motion.NewStore(files, *manifestPath)
	must(err)
	must(os.MkdirAll(*out, 0755))
	rows := []measurement{}
	exportedHashes := map[string]string{}
	for _, entry := range manifest.Variants {
		if !entry.Enabled {
			continue
		}
		v := store.Get(entry.ID)
		if v == nil {
			panic(store.Error(entry.ID))
		}
		actions := make([]string, 0, len(v.Frames))
		for action := range v.Frames {
			actions = append(actions, action)
		}
		sort.Strings(actions)
		for _, height := range []int{64, 96, 128} {
			for _, dir := range []int{-1, 1} {
				for _, action := range actions {
					phase := 0.0
					for frame, duration := range v.Motions[action].Durations {
						im := image.NewRGBA(image.Rect(0, 0, height*4, height*3))
						x, ground := height, height*2
						v.Draw(im, action, phase+duration/2, x, ground, height*3/2, height, dir)
						solid, area := bounds(im, 180)
						alpha, _ := bounds(im, 1)
						components := [][2]int{}
						start := -1
						for px := solid.Min.X; px <= solid.Max.X; px++ {
							present := false
							if px < solid.Max.X {
								for y := solid.Max.Y - max(2, (height+31)/32); y < solid.Max.Y; y++ {
									if im.RGBAAt(px, y).A >= 180 {
										present = true
										break
									}
								}
							}
							if present && start < 0 {
								start = px
							}
							if !present && start >= 0 {
								if px-start >= 2 {
									components = append(components, [2]int{start - x, px - 1 - x})
								}
								start = -1
							}
						}
						rows = append(rows, measurement{entry.ID, action, frame, height, dir, phase,
							relative(solid, x, ground), relative(alpha, x, ground), area, components,
							alpha.Min.X <= 0 || alpha.Min.Y <= 0 || alpha.Max.X >= im.Bounds().Max.X || alpha.Max.Y >= im.Bounds().Max.Y})
						if *exportFrames || frame == 0 && dir == 1 && height == 96 {
							name := entry.ID + "-" + action + ".png"
							if *exportFrames {
								name = fmt.Sprintf("%s-%s-h%d-d%d-f%03d.png", entry.ID, action, height, dir, frame)
							}
							f, err := os.Create(filepath.Join(*out, name))
							must(err)
							digest := sha256.New()
							must(png.Encode(io.MultiWriter(f, digest), im))
							must(f.Close())
							exportedHashes[name] = fmt.Sprintf("%x", digest.Sum(nil))
						}
						phase += duration
					}
				}
			}
		}
		store.Retain(map[string]bool{})
	}
	data, err := json.MarshalIndent(map[string]any{"schema": 2, "render": "shared Go motion.Draw; alpha>=180 solid; bottom ceil(3*height/96) rendered rows, min 2",
		"manifestSha256": fmt.Sprintf("%x", sha256.Sum256(raw)), "exportedFrameSha256": exportedHashes, "allMeasuredFramesExported": *exportFrames,
		"scope": "envelope screening only; no automatic stance/body inference", "measurements": rows}, "", "  ")
	must(err)
	must(os.WriteFile(filepath.Join(*out, "measurements.json"), data, 0644))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
