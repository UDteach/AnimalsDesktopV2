package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAuditReportsPartialFrameSet(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	writeFrame(t, filepath.Join(framesDir, "frame-00.png"), frameW, frameH, false)
	writeFrame(t, filepath.Join(framesDir, "frame-01.png"), frameW+1, frameH, false)

	report, err := audit("", framesDir, "frame-%02d.png", false, false, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	if report.Valid != 1 || report.Invalid != 1 || report.Missing != totalFrames-2 {
		t.Fatalf("report counts valid/invalid/missing = %d/%d/%d", report.Valid, report.Invalid, report.Missing)
	}
	if report.Sets[0].Completed {
		t.Fatalf("partial set reported complete")
	}
}

func TestAuditRejectsOpaqueFrame(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	writeFrame(t, filepath.Join(framesDir, "frame-00.png"), frameW, frameH, true)

	report, err := audit("", framesDir, "frame-%02d.png", false, false, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	if frame.Status != "invalid" || frame.Error != "no transparent background" {
		t.Fatalf("frame status/error = %q/%q", frame.Status, frame.Error)
	}
}

func TestAuditRootScansTenSets(t *testing.T) {
	root := t.TempDir()
	writeFrame(t, filepath.Join(root, "set03", "frame-12.png"), frameW, frameH, false)

	report, err := audit(root, "", "frame-%02d.png", false, false, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	if report.SetCount != motionSets || report.FrameCount != motionSets*totalFrames {
		t.Fatalf("report set/frame count = %d/%d", report.SetCount, report.FrameCount)
	}
	if report.Valid != 1 {
		t.Fatalf("valid frames = %d, want 1", report.Valid)
	}
}

func TestAuditArtifactWarningsDetectsLowHorizontalRun(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addHorizontalRun(t, path, 16, 55, 64)

	report, err := audit("", framesDir, "frame-%02d.png", false, true, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	if frame.Status != "valid" {
		t.Fatalf("frame status = %q", frame.Status)
	}
	if len(frame.Warnings) == 0 {
		t.Fatalf("expected artifact warning")
	}
}

func TestAuditArtifactWarningsDetectsDisconnectedComponents(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addDetachedBlock(t, path, 8, 57, 4, 4)

	report, err := audit("", framesDir, "frame-%02d.png", false, true, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "disconnected alpha components") {
		t.Fatalf("warnings = %q, want disconnected alpha components warning", warnings)
	}
}

func TestAuditArtifactWarningsDetectsTransparentPinholes(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addTransparentHole(t, path, 42, 36, 2, 2)

	report, err := audit("", framesDir, "frame-%02d.png", false, true, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "transparent pinholes") {
		t.Fatalf("warnings = %q, want transparent pinholes warning", warnings)
	}
}

func TestAuditArtifactWarningsDetectsLowerShelf(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addLowerShelf(t, path)

	report, err := audit("", framesDir, "frame-%02d.png", false, true, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "possible lower ledge/shelf artifact") {
		t.Fatalf("warnings = %q, want lower shelf artifact warning", warnings)
	}
}

func TestAuditMotionWarningsDetectsSizeJump(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	writeRectFrame(t, filepath.Join(framesDir, "frame-00.png"), 30, 24, 28, 24)
	writeRectFrame(t, filepath.Join(framesDir, "frame-01.png"), 10, 10, 70, 48)

	report, err := audit("", framesDir, "frame-%02d.png", false, false, true)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[1]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "motion consistency: bbox width jumps") ||
		!strings.Contains(warnings, "motion consistency: bbox height jumps") ||
		!strings.Contains(warnings, "motion consistency: body alpha area jumps") {
		t.Fatalf("warnings = %q, want size and body-area motion warnings", warnings)
	}
}

func TestParseMotionBoundaries(t *testing.T) {
	t.Run("empty preserves default", func(t *testing.T) {
		got, err := parseMotionBoundaries("")
		if err != nil {
			t.Fatalf("parseMotionBoundaries() error = %v", err)
		}
		if got != nil {
			t.Fatalf("parseMotionBoundaries() = %v, want nil", got)
		}
	})

	t.Run("strictly increasing values", func(t *testing.T) {
		got, err := parseMotionBoundaries("4, 12,20,26,32,40,48,56,61")
		if err != nil {
			t.Fatalf("parseMotionBoundaries() error = %v", err)
		}
		want := []int{4, 12, 20, 26, 32, 40, 48, 56, 61}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("parseMotionBoundaries() = %v, want %v", got, want)
		}
	})

	for _, value := range []string{
		" ",
		"4,,12",
		"four",
		"+4",
		"1.5",
		"-1",
		"0",
		"62",
		"4,4",
		"12,4",
	} {
		t.Run("reject "+value, func(t *testing.T) {
			if got, err := parseMotionBoundaries(value); err == nil {
				t.Fatalf("parseMotionBoundaries(%q) = %v, want error", value, got)
			}
		})
	}
}

func TestAuditMotionBoundariesDefaultCompatibility(t *testing.T) {
	framesDir := filepath.Join(t.TempDir(), "set00")
	writeRectFrame(t, filepath.Join(framesDir, "frame-00.png"), 30, 24, 28, 24)
	writeRectFrame(t, filepath.Join(framesDir, "frame-01.png"), 10, 10, 70, 48)

	legacy, err := audit("", framesDir, "frame-%02d.png", false, false, true)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	boundaryAware, err := auditWithMotionBoundaries("", framesDir, "frame-%02d.png", false, false, true, nil)
	if err != nil {
		t.Fatalf("auditWithMotionBoundaries() error = %v", err)
	}
	if !reflect.DeepEqual(boundaryAware, legacy) {
		t.Fatalf("empty motion boundaries changed the audit report")
	}
	encoded, err := json.Marshal(boundaryAware)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), `"motion_boundaries"`) {
		t.Fatalf("default report unexpectedly contains motion_boundaries: %s", encoded)
	}
}

func TestAuditMotionBoundariesSuppressCrossBoundaryAdjacentWarnings(t *testing.T) {
	framesDir := filepath.Join(t.TempDir(), "set00")
	writeRectFrame(t, filepath.Join(framesDir, "frame-00.png"), 30, 24, 28, 24)
	framePath := filepath.Join(framesDir, "frame-01.png")
	writeRectFrame(t, framePath, 10, 10, 70, 48)
	addDetachedBlock(t, framePath, 3, 58, 4, 4)

	report, err := auditWithMotionBoundaries("", framesDir, "frame-%02d.png", false, true, true, []int{1})
	if err != nil {
		t.Fatalf("auditWithMotionBoundaries() error = %v", err)
	}
	warnings := strings.Join(report.Sets[0].Frames[1].Warnings, "\n")
	if strings.Contains(warnings, "motion consistency:") {
		t.Fatalf("cross-boundary warnings = %q, want none", warnings)
	}
	if !strings.Contains(warnings, "disconnected alpha components") {
		t.Fatalf("artifact warning missing at motion boundary: %q", warnings)
	}
}

func TestAuditMotionBoundariesSuppressIsolatedWarningsAcrossEitherEdge(t *testing.T) {
	for _, tc := range []struct {
		name     string
		boundary int
	}{
		{name: "previous-current edge", boundary: 1},
		{name: "current-next edge", boundary: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			framesDir := filepath.Join(t.TempDir(), "set00")
			writeRectFrame(t, filepath.Join(framesDir, "frame-00.png"), 30, 24, 28, 24)
			writeRectFrame(t, filepath.Join(framesDir, "frame-01.png"), 10, 10, 70, 48)
			writeRectFrame(t, filepath.Join(framesDir, "frame-02.png"), 30, 24, 28, 24)

			report, err := auditWithMotionBoundaries("", framesDir, "frame-%02d.png", false, false, true, []int{tc.boundary})
			if err != nil {
				t.Fatalf("auditWithMotionBoundaries() error = %v", err)
			}
			warnings := strings.Join(report.Sets[0].Frames[1].Warnings, "\n")
			if strings.Contains(warnings, "isolated") {
				t.Fatalf("isolated warnings across boundary = %q, want none", warnings)
			}
			if tc.boundary == 2 && !strings.Contains(warnings, "from previous frame 00 to 01") {
				t.Fatalf("within-segment adjacent warning missing: %q", warnings)
			}
		})
	}
}

func TestAuditMotionBoundariesRetainWithinSegmentAndArtifactWarnings(t *testing.T) {
	framesDir := filepath.Join(t.TempDir(), "set00")
	writeRectFrame(t, filepath.Join(framesDir, "frame-03.png"), 30, 24, 28, 24)
	framePath := filepath.Join(framesDir, "frame-04.png")
	writeRectFrame(t, framePath, 10, 10, 70, 48)
	addDetachedBlock(t, framePath, 3, 58, 4, 4)

	report, err := auditWithMotionBoundaries("", framesDir, "frame-%02d.png", false, true, true, []int{1})
	if err != nil {
		t.Fatalf("auditWithMotionBoundaries() error = %v", err)
	}
	warnings := strings.Join(report.Sets[0].Frames[4].Warnings, "\n")
	if !strings.Contains(warnings, "motion consistency: bbox width jumps from previous frame 03 to 04") {
		t.Fatalf("within-segment motion warning missing: %q", warnings)
	}
	if !strings.Contains(warnings, "disconnected alpha components") {
		t.Fatalf("artifact warning missing: %q", warnings)
	}
	if !reflect.DeepEqual(report.MotionBoundaries, []int{1}) {
		t.Fatalf("report motion boundaries = %v, want [1]", report.MotionBoundaries)
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if !strings.Contains(string(encoded), `"motion_boundaries":[1]`) {
		t.Fatalf("report JSON lacks motion boundary provenance: %s", encoded)
	}
}

func writeFrame(t *testing.T, path string, w int, h int, opaque bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if opaque {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetRGBA(x, y, color.RGBA{R: 210, G: 210, B: 210, A: 255})
			}
		}
	} else {
		for y := 24; y < minInt(54, h); y++ {
			for x := 20; x < minInt(76, w); x++ {
				img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 130, A: 255})
			}
		}
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

func writeRectFrame(t *testing.T, path string, x int, y int, w int, h int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			img.SetRGBA(px, py, color.RGBA{R: 130, G: 130, B: 130, A: 255})
		}
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

func addDetachedBlock(t *testing.T, path string, x int, y int, w int, h int) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for py := y; py < y+h; py++ {
			for px := x; px < x+w; px++ {
				rgba.SetRGBA(px, py, color.RGBA{R: 70, G: 60, B: 50, A: 230})
			}
		}
	})
}

func addTransparentHole(t *testing.T, path string, x int, y int, w int, h int) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for py := y; py < y+h; py++ {
			for px := x; px < x+w; px++ {
				rgba.SetRGBA(px, py, color.RGBA{})
			}
		}
	})
}

func addLowerShelf(t *testing.T, path string) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for px := 8; px < 88; px++ {
			rgba.SetRGBA(px, 50, color.RGBA{R: 90, G: 84, B: 76, A: 190})
		}
		for py := 54; py < 58; py++ {
			for px := 30; px < 34; px++ {
				rgba.SetRGBA(px, py, color.RGBA{R: 80, G: 76, B: 70, A: 220})
			}
		}
	})
}

func addHorizontalRun(t *testing.T, path string, x int, y int, length int) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for px := x; px < x+length; px++ {
			rgba.SetRGBA(px, y, color.RGBA{R: 80, G: 70, B: 60, A: 180})
		}
	})
}

func updateFrame(t *testing.T, path string, draw func(*image.RGBA)) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	img, err := png.Decode(f)
	if closeErr := f.Close(); closeErr != nil {
		t.Fatalf("close: %v", closeErr)
	}
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	rgba := image.NewRGBA(img.Bounds())
	for py := img.Bounds().Min.Y; py < img.Bounds().Max.Y; py++ {
		for px := img.Bounds().Min.X; px < img.Bounds().Max.X; px++ {
			rgba.Set(px, py, img.At(px, py))
		}
	}
	draw(rgba)
	out, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer out.Close()
	if err := png.Encode(out, rgba); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
