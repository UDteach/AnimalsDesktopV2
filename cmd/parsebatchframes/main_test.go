package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseIndependentCopiesAndReports(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "source.png")
	writeTestPNG(t, src, 80, 60, []image.Rectangle{image.Rect(18, 16, 62, 48)})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "independent",
		Frames: []manifestFrame{{
			ID:     "set01-frame-10",
			Source: "source.png",
		}},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), filepath.Join(root, "review", "report.json"))
	if err != nil {
		t.Fatalf("parseBatch() error = %v", err)
	}
	if report.Parsed != 1 || report.Rejected != 0 {
		t.Fatalf("parsed/rejected = %d/%d", report.Parsed, report.Rejected)
	}
	frame := report.Frames[0]
	if frame.Status != "parsed" || frame.OutputSHA256 == "" || frame.SourceSHA256 == "" {
		t.Fatalf("frame report = %+v", frame)
	}
	assertPNGSize(t, frame.Output, 80, 60)
}

func TestParseStrip2x1SlicesExactCells(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "strip.png")
	writeTestPNG(t, src, 192, 64, []image.Rectangle{
		image.Rect(16, 20, 70, 50),
		image.Rect(118, 18, 174, 52),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "strip-2x1",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "strip.png", Cell: 0},
			{ID: "set01-frame-11", Source: "strip.png", Cell: 1},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err != nil {
		t.Fatalf("parseBatch() error = %v", err)
	}
	if report.Parsed != 2 || report.Rejected != 0 {
		t.Fatalf("parsed/rejected = %d/%d", report.Parsed, report.Rejected)
	}
	for _, frame := range report.Frames {
		assertPNGSize(t, frame.Output, frameW, frameH)
		if frame.ComponentCount != 1 {
			t.Fatalf("component count for %s = %d, want 1", frame.ID, frame.ComponentCount)
		}
	}
}

func TestParseStripRejectsWrongDimensions(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "strip.png")
	writeTestPNG(t, src, 191, 64, []image.Rectangle{image.Rect(16, 20, 70, 50)})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "strip-2x1",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "strip.png", Cell: 0},
			{ID: "set01-frame-11", Source: "strip.png", Cell: 1},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want wrong-dimension rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "want 192x64") {
		t.Fatalf("report errors = %#v", report)
	}
}

func TestParseStripRejectsBoundaryAlpha(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "strip.png")
	writeTestPNG(t, src, 192, 64, []image.Rectangle{
		image.Rect(16, 20, 70, 50),
		image.Rect(95, 28, 97, 34),
		image.Rect(118, 18, 174, 52),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "strip-2x1",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "strip.png", Cell: 0},
			{ID: "set01-frame-11", Source: "strip.png", Cell: 1},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want boundary-alpha rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "x-boundary guard") {
		t.Fatalf("report errors = %#v", report)
	}
}

func TestParseStripRejectsCellEdgeAlpha(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "strip.png")
	writeTestPNG(t, src, 192, 64, []image.Rectangle{
		image.Rect(0, 20, 40, 50),
		image.Rect(118, 18, 174, 52),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "strip-2x1",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "strip.png", Cell: 0},
			{ID: "set01-frame-11", Source: "strip.png", Cell: 1},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want cell-edge rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "alpha touches horizontal canvas edge") {
		t.Fatalf("report errors = %#v", report)
	}
}

func TestParseStripRejectsDisconnectedAlpha(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "strip.png")
	writeTestPNG(t, src, 192, 64, []image.Rectangle{
		image.Rect(16, 20, 70, 50),
		image.Rect(8, 8, 12, 12),
		image.Rect(118, 18, 174, 52),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "strip-2x1",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "strip.png", Cell: 0},
			{ID: "set01-frame-11", Source: "strip.png", Cell: 1},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want disconnected-alpha rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "disconnected alpha components") {
		t.Fatalf("report errors = %#v", report)
	}
}

func TestParseGrid2x2SlicesExactCells(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "grid.png")
	writeTestPNG(t, src, 192, 128, []image.Rectangle{
		image.Rect(16, 20, 70, 50),
		image.Rect(118, 18, 174, 52),
		image.Rect(18, 82, 72, 114),
		image.Rect(120, 84, 176, 116),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "grid-2x2",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "grid.png", Cell: 0},
			{ID: "set01-frame-11", Source: "grid.png", Cell: 1},
			{ID: "set01-frame-12", Source: "grid.png", Cell: 2},
			{ID: "set01-frame-13", Source: "grid.png", Cell: 3},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err != nil {
		t.Fatalf("parseBatch() error = %v", err)
	}
	if report.Parsed != 4 || report.Rejected != 0 {
		t.Fatalf("parsed/rejected = %d/%d", report.Parsed, report.Rejected)
	}
	for _, frame := range report.Frames {
		assertPNGSize(t, frame.Output, frameW, frameH)
		if frame.ComponentCount != 1 {
			t.Fatalf("component count for %s = %d, want 1", frame.ID, frame.ComponentCount)
		}
	}
}

func TestParseGridRejectsHorizontalBoundaryAlpha(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "grid.png")
	writeTestPNG(t, src, 192, 128, []image.Rectangle{
		image.Rect(16, 20, 70, 50),
		image.Rect(118, 18, 174, 52),
		image.Rect(40, 63, 48, 65),
		image.Rect(18, 82, 72, 114),
		image.Rect(120, 84, 176, 116),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "grid-2x2",
		Frames: []manifestFrame{
			{ID: "set01-frame-10", Source: "grid.png", Cell: 0},
			{ID: "set01-frame-11", Source: "grid.png", Cell: 1},
			{ID: "set01-frame-12", Source: "grid.png", Cell: 2},
			{ID: "set01-frame-13", Source: "grid.png", Cell: 3},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want horizontal-boundary rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "y-boundary guard") {
		t.Fatalf("report errors = %#v", report)
	}
}

func TestParseFixedGridSlicesDeclaredCells(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "fixed-grid.png")
	writeTestPNG(t, src, 560, 360, []image.Rectangle{
		image.Rect(28, 32, 120, 96),
		image.Rect(188, 32, 280, 96),
		image.Rect(348, 32, 440, 96),
		image.Rect(28, 172, 120, 236),
		image.Rect(188, 172, 280, 236),
		image.Rect(348, 172, 440, 236),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout:     "grid-fixed",
		CellWidth:  128,
		CellHeight: 96,
		Columns:    3,
		Rows:       2,
		OriginX:    16,
		OriginY:    20,
		PitchX:     160,
		PitchY:     140,
		Frames: []manifestFrame{
			{ID: "set00-frame-00", Source: "fixed-grid.png", Cell: 0},
			{ID: "set00-frame-01", Source: "fixed-grid.png", Cell: 1},
			{ID: "set00-frame-02", Source: "fixed-grid.png", Cell: 2},
			{ID: "set00-frame-03", Source: "fixed-grid.png", Cell: 3},
			{ID: "set00-frame-04", Source: "fixed-grid.png", Cell: 4},
			{ID: "set00-frame-05", Source: "fixed-grid.png", Cell: 5},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err != nil {
		t.Fatalf("parseBatch() error = %v", err)
	}
	if report.Parsed != 6 || report.Rejected != 0 {
		t.Fatalf("parsed/rejected = %d/%d", report.Parsed, report.Rejected)
	}
	for _, frame := range report.Frames {
		assertPNGSize(t, frame.Output, 128, 96)
		if frame.ComponentCount != 1 {
			t.Fatalf("component count for %s = %d, want 1", frame.ID, frame.ComponentCount)
		}
	}
}

func TestParseFixedGridChromaGreenBackground(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "fixed-grid-green.png")
	writeGreenTestPNG(t, src, 320, 240, []image.Rectangle{
		image.Rect(24, 28, 110, 90),
		image.Rect(184, 28, 270, 90),
		image.Rect(24, 148, 110, 210),
		image.Rect(184, 148, 270, 210),
	})
	manifestPath := writeManifest(t, root, manifest{
		Layout:     "grid-fixed",
		CellWidth:  128,
		CellHeight: 96,
		Columns:    2,
		Rows:       2,
		OriginX:    16,
		OriginY:    16,
		PitchX:     160,
		PitchY:     120,
		Background: "chroma-green",
		Frames: []manifestFrame{
			{ID: "set00-frame-00", Source: "fixed-grid-green.png", Cell: 0},
			{ID: "set00-frame-01", Source: "fixed-grid-green.png", Cell: 1},
			{ID: "set00-frame-02", Source: "fixed-grid-green.png", Cell: 2},
			{ID: "set00-frame-03", Source: "fixed-grid-green.png", Cell: 3},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err != nil {
		t.Fatalf("parseBatch() error = %v", err)
	}
	if report.Parsed != 4 || report.Rejected != 0 {
		t.Fatalf("parsed/rejected = %d/%d", report.Parsed, report.Rejected)
	}
	for _, frame := range report.Frames {
		if frame.Content.W <= 0 || frame.Content.H <= 0 {
			t.Fatalf("empty content for %s: %+v", frame.ID, frame.Content)
		}
		assertPNGSize(t, frame.Output, 128, 96)
	}
}

func TestParseFixedGridRejectsChromaPinholes(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "fixed-grid-green-pinholes.png")
	writeGreenTestPNGWithHole(t, src, 180, 140, image.Rect(24, 28, 130, 96), image.Rect(70, 58, 73, 61))
	manifestPath := writeManifest(t, root, manifest{
		Layout:     "grid-fixed",
		CellWidth:  160,
		CellHeight: 120,
		Columns:    1,
		Rows:       1,
		OriginX:    10,
		OriginY:    10,
		Background: "chroma-green",
		Frames: []manifestFrame{
			{ID: "set00-frame-00", Source: "fixed-grid-green-pinholes.png", Cell: 0},
		},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want chroma pinhole rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "transparent/chroma pinholes") {
		t.Fatalf("report errors = %#v", report)
	}
}

func TestParseRejectsAcceptedFramesOutput(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "source.png")
	writeTestPNG(t, src, 80, 60, []image.Rectangle{image.Rect(18, 16, 62, 48)})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "independent",
		Frames: []manifestFrame{{ID: "frame", Source: "source.png"}},
	})

	_, err := parseBatch(manifestPath, filepath.Join(root, "docs", "art-source", "gecko", "motion-source", "accepted-frames", "parsed"), "")
	if err == nil || !strings.Contains(err.Error(), "accepted-frames") {
		t.Fatalf("parseBatch() error = %v, want accepted-frames refusal", err)
	}
}

func TestParseRejectsOutputPathEscape(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "source.png")
	writeTestPNG(t, src, 80, 60, []image.Rectangle{image.Rect(18, 16, 62, 48)})
	manifestPath := writeManifest(t, root, manifest{
		Layout: "independent",
		Frames: []manifestFrame{{
			ID:     "frame",
			Source: "source.png",
			Output: "../escape.png",
		}},
	})

	report, err := parseBatch(manifestPath, filepath.Join(root, "review", "parsed"), "")
	if err == nil {
		t.Fatalf("parseBatch() succeeded, want output escape rejection")
	}
	if report == nil || !strings.Contains(strings.Join(report.Errors, "\n"), "output must be a filename") {
		t.Fatalf("report errors = %#v", report)
	}
}

func writeManifest(t *testing.T, root string, m manifest) string {
	t.Helper()
	path := filepath.Join(root, "batch.json")
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return path
}

func writeTestPNG(t *testing.T, path string, w int, h int, rects []image.Rectangle) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i, rect := range rects {
		c := color.RGBA{R: uint8(80 + i*40), G: uint8(70 + i*30), B: uint8(60 + i*20), A: 255}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				img.SetRGBA(x, y, c)
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

func writeGreenTestPNG(t *testing.T, path string, w int, h int, rects []image.Rectangle) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	for i, rect := range rects {
		c := color.RGBA{R: uint8(80 + i*40), G: uint8(70 + i*30), B: uint8(60 + i*20), A: 255}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				img.SetRGBA(x, y, c)
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

func writeGreenTestPNGWithHole(t *testing.T, path string, w int, h int, rect image.Rectangle, hole image.Rectangle) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 120, G: 118, B: 116, A: 255})
		}
	}
	for y := hole.Min.Y; y < hole.Max.Y; y++ {
		for x := hole.Min.X; x < hole.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
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

func assertPNGSize(t *testing.T, path string, wantW int, wantH int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != wantW || bounds.Dy() != wantH {
		t.Fatalf("%s size = %dx%d, want %dx%d", path, bounds.Dx(), bounds.Dy(), wantW, wantH)
	}
}
