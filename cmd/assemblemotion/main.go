package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

const (
	frameW      = 96
	frameH      = 64
	totalFrames = 62
)

type assembleReport struct {
	FramesDir   string        `json:"frames_dir"`
	Pattern     string        `json:"pattern"`
	Output      string        `json:"output"`
	FrameWidth  int           `json:"frame_width"`
	FrameHeight int           `json:"frame_height"`
	FrameCount  int           `json:"frame_count"`
	SheetWidth  int           `json:"sheet_width"`
	SheetHeight int           `json:"sheet_height"`
	Frames      []frameReport `json:"frames"`
	Warnings    []string      `json:"warnings,omitempty"`
}

type frameReport struct {
	Frame    int      `json:"frame"`
	Path     string   `json:"path"`
	Content  rectJSON `json:"content"`
	Warnings []string `json:"warnings,omitempty"`
}

type rectJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func main() {
	framesDir := flag.String("frames-dir", "", "directory containing frame PNGs")
	outPath := flag.String("out", "", "output 5952x64 motion source sheet PNG")
	pattern := flag.String("pattern", "frame-%02d.png", "frame filename pattern with one integer verb")
	reportPath := flag.String("report", "", "optional JSON report path")
	flag.Parse()

	if *framesDir == "" {
		fatalf("-frames-dir is required")
	}
	if *outPath == "" {
		fatalf("-out is required")
	}

	report, err := assembleMotion(*framesDir, *pattern, *outPath)
	if err != nil {
		fatalf("%v", err)
	}
	if *reportPath != "" {
		if err := writeReport(*reportPath, report); err != nil {
			fatalf("write report: %v", err)
		}
	}
	fmt.Printf("assembled %d standalone frames into %s\n", report.FrameCount, filepath.ToSlash(*outPath))
}

func assembleMotion(framesDir string, pattern string, outPath string) (assembleReport, error) {
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	report := assembleReport{
		FramesDir:   filepath.ToSlash(framesDir),
		Pattern:     pattern,
		Output:      filepath.ToSlash(outPath),
		FrameWidth:  frameW,
		FrameHeight: frameH,
		FrameCount:  totalFrames,
		SheetWidth:  frameW * totalFrames,
		SheetHeight: frameH,
		Frames:      make([]frameReport, 0, totalFrames),
	}

	for frame := 0; frame < totalFrames; frame++ {
		framePath := filepath.Join(framesDir, fmt.Sprintf(pattern, frame))
		img, err := openPNG(framePath)
		if err != nil {
			return assembleReport{}, fmt.Errorf("frame %02d: %w", frame, err)
		}
		bounds := img.Bounds()
		if bounds.Dx() != frameW || bounds.Dy() != frameH {
			return assembleReport{}, fmt.Errorf("frame %02d bounds = %dx%d, want %dx%d: %s", frame, bounds.Dx(), bounds.Dy(), frameW, frameH, framePath)
		}
		content := alphaBounds(img, bounds)
		if content.Empty() {
			return assembleReport{}, fmt.Errorf("frame %02d has no visible alpha: %s", frame, framePath)
		}
		if content == bounds {
			return assembleReport{}, fmt.Errorf("frame %02d has no transparent background: %s", frame, framePath)
		}

		frameWarnings := frameWarnings(content, bounds)
		if len(frameWarnings) > 0 {
			report.Warnings = append(report.Warnings, fmt.Sprintf("frame %02d needs visual crop review", frame))
		}
		report.Frames = append(report.Frames, frameReport{
			Frame:    frame,
			Path:     filepath.ToSlash(framePath),
			Content:  rectToJSON(content),
			Warnings: frameWarnings,
		})
		dst := image.Rect(frame*frameW, 0, (frame+1)*frameW, frameH)
		draw.Draw(sheet, dst, img, bounds.Min, draw.Over)
	}

	if err := writePNG(outPath, sheet); err != nil {
		return assembleReport{}, err
	}
	return report, nil
}

func openPNG(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode PNG %s: %w", path, err)
	}
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), src, bounds.Min, draw.Src)
	return dst, nil
}

func writePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func writeReport(path string, report assembleReport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func alphaBounds(img image.Image, rect image.Rectangle) image.Rectangle {
	rect = rect.Intersect(img.Bounds())
	minX, minY := rect.Max.X, rect.Max.Y
	maxX, maxY := rect.Min.X, rect.Min.Y
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a <= 0x0800 {
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
	if maxX <= minX || maxY <= minY {
		return image.Rectangle{}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func frameWarnings(content image.Rectangle, bounds image.Rectangle) []string {
	warnings := []string{}
	if content.Min.X <= bounds.Min.X || content.Max.X >= bounds.Max.X {
		warnings = append(warnings, "alpha touches horizontal canvas edge")
	}
	if content.Min.Y <= bounds.Min.Y || content.Max.Y >= bounds.Max.Y {
		warnings = append(warnings, "alpha touches vertical canvas edge")
	}
	return warnings
}

func rectToJSON(rect image.Rectangle) rectJSON {
	return rectJSON{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
