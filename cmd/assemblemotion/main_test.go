package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssembleMotionWritesSheetFromStandaloneFrames(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "frames")
	for frame := 0; frame < totalFrames; frame++ {
		img := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
		img.SetRGBA(4, 4, color.RGBA{R: byte(frame), G: 10, B: 20, A: 255})
		for y := 28; y < 52; y++ {
			for x := 20; x < 76; x++ {
				img.SetRGBA(x, y, color.RGBA{R: 145, G: 145, B: 140, A: 255})
			}
		}
		writeTestPNG(t, filepath.Join(framesDir, twoDigitsName(frame)), img)
	}

	outPath := filepath.Join(root, "sheet.png")
	report, err := assembleMotion(framesDir, "frame-%02d.png", outPath)
	if err != nil {
		t.Fatalf("assembleMotion() error = %v", err)
	}
	if report.FrameCount != totalFrames || len(report.Frames) != totalFrames {
		t.Fatalf("report frame count = %d/%d, want %d", report.FrameCount, len(report.Frames), totalFrames)
	}
	sheet := openTestPNG(t, outPath)
	if got := sheet.Bounds(); got.Dx() != frameW*totalFrames || got.Dy() != frameH {
		t.Fatalf("sheet bounds = %v, want %dx%d", got, frameW*totalFrames, frameH)
	}
	got := sheet.RGBAAt(12*frameW+4, 4)
	if got.R != 12 || got.A != 255 {
		t.Fatalf("frame marker = %#v, want frame 12 marker", got)
	}
}

func TestAssembleMotionRejectsOpaqueBackground(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "frames")
	for frame := 0; frame < totalFrames; frame++ {
		img := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
		fill := color.RGBA{R: 140, G: 140, B: 140, A: 255}
		if frame > 0 {
			fill = color.RGBA{R: 0, G: 0, B: 0, A: 0}
		}
		for y := 0; y < frameH; y++ {
			for x := 0; x < frameW; x++ {
				img.SetRGBA(x, y, fill)
			}
		}
		if frame > 0 {
			img.SetRGBA(20, 20, color.RGBA{R: 140, G: 140, B: 140, A: 255})
		}
		writeTestPNG(t, filepath.Join(framesDir, twoDigitsName(frame)), img)
	}

	_, err := assembleMotion(framesDir, "frame-%02d.png", filepath.Join(root, "sheet.png"))
	if err == nil {
		t.Fatalf("assembleMotion() succeeded for opaque frame")
	}
	if !strings.Contains(err.Error(), "transparent background") {
		t.Fatalf("assembleMotion() error = %v, want transparent-background failure", err)
	}
}

func TestAssembleMotionRejectsWrongFrameSize(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "frames")
	for frame := 0; frame < totalFrames; frame++ {
		w, h := frameW, frameH
		if frame == 7 {
			w = frameW + 1
		}
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		img.SetRGBA(10, 10, color.RGBA{R: 120, G: 120, B: 120, A: 255})
		writeTestPNG(t, filepath.Join(framesDir, twoDigitsName(frame)), img)
	}

	_, err := assembleMotion(framesDir, "frame-%02d.png", filepath.Join(root, "sheet.png"))
	if err == nil {
		t.Fatalf("assembleMotion() succeeded for wrong-sized frame")
	}
	if !strings.Contains(err.Error(), "bounds") {
		t.Fatalf("assembleMotion() error = %v, want bounds failure", err)
	}
}

func twoDigitsName(v int) string {
	return "frame-" + string([]byte{'0' + byte(v/10), '0' + byte(v%10)}) + ".png"
}

func writeTestPNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func openTestPNG(t *testing.T, path string) *image.RGBA {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	rgba, ok := img.(*image.RGBA)
	if ok {
		return rgba
	}
	out := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			out.Set(x, y, img.At(x, y))
		}
	}
	return out
}
