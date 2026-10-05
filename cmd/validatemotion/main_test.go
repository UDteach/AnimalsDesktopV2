package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"animals-desktop/internal/catalog"
)

func TestValidateVariantReportsDraftNotReleaseReady(t *testing.T) {
	root := t.TempDir()
	set00Path := writeMotionSourceFamily(t, root)
	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	}

	report, err := validateVariant(variant)
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if report.SourceSets != runtimeMotionSets || report.RuntimeSets != runtimeMotionSets || report.FramesPerSet != totalFrames {
		t.Fatalf("report source/runtime sets/frames = %d/%d/%d, want %d/%d/%d", report.SourceSets, report.RuntimeSets, report.FramesPerSet, runtimeMotionSets, runtimeMotionSets, totalFrames)
	}
	if report.ReleaseReady {
		t.Fatalf("draft motion source reported release-ready")
	}
	if report.AcceptedSource {
		t.Fatalf("draft motion source reported accepted")
	}
	if len(report.Warnings) == 0 {
		t.Fatalf("draft motion source should report a warning")
	}
}

func TestValidateVariantReportsAcceptedReleaseReady(t *testing.T) {
	root := t.TempDir()
	set00Path := writeMotionSourceFamily(t, root)
	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionAccepted,
		MotionSourcePath: set00Path,
	}

	report, err := validateVariant(variant)
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if !report.ReleaseReady {
		t.Fatalf("accepted motion source reported not release-ready")
	}
	if !report.AcceptedSource {
		t.Fatalf("accepted motion source reported not accepted")
	}
	if report.UniqueSetHashes != runtimeMotionSets {
		t.Fatalf("unique set hashes = %d, want %d", report.UniqueSetHashes, runtimeMotionSets)
	}
}

func TestValidateVariantRejectsDuplicateOptionalSourceFamily(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	for set := 0; set < runtimeMotionSets; set++ {
		writePNG(t, filepath.Join(root, "animal-set"+twoDigits(set)+"-source.png"), testMotionSheet())
	}

	report, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionAccepted,
		MotionSourcePath: set00Path,
	})
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if report.ReleaseReady {
		t.Fatalf("duplicate optional source family reported release-ready")
	}
	if report.UniqueSetHashes != 1 {
		t.Fatalf("unique set hashes = %d, want 1", report.UniqueSetHashes)
	}
	if len(report.Warnings) != 1 || !strings.Contains(report.Warnings[0], "byte-identical") {
		t.Fatalf("warnings = %v, want duplicate-family warning", report.Warnings)
	}
}

func TestValidateVariantAllowsSingleDraftSourceSet(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for frame := 0; frame < totalFrames; frame++ {
		x0 := frame * frameW
		sheet.SetRGBA(x0+6, 6, color.RGBA{R: 80, G: byte(frame), B: 80, A: 255})
		for y := 28; y < 52; y++ {
			for x := 20; x < 76; x++ {
				sheet.SetRGBA(x0+x, y, color.RGBA{R: 150, G: 150, B: 145, A: 255})
			}
		}
	}
	writePNG(t, set00Path, sheet)

	report, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	})
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if report.SourceSets != 1 || report.RuntimeSets != runtimeMotionSets {
		t.Fatalf("source/runtime sets = %d/%d, want 1/%d", report.SourceSets, report.RuntimeSets, runtimeMotionSets)
	}
	if report.ReleaseReady {
		t.Fatalf("single draft source set reported release-ready")
	}
	if report.AcceptedSource {
		t.Fatalf("single draft source set reported accepted")
	}
	if len(report.Warnings) != 1 || !strings.Contains(report.Warnings[0], "draft") {
		t.Fatalf("warnings = %v, want only draft warning", report.Warnings)
	}
}

func TestValidateVariantReportsSingleAcceptedSourceSetReleaseReady(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for frame := 0; frame < totalFrames; frame++ {
		x0 := frame * frameW
		sheet.SetRGBA(x0+6, 6, color.RGBA{R: 80, G: byte(frame), B: 80, A: 255})
		for y := 28; y < 52; y++ {
			for x := 20; x < 76; x++ {
				sheet.SetRGBA(x0+x, y, color.RGBA{R: 150, G: 150, B: 145, A: 255})
			}
		}
	}
	writePNG(t, set00Path, sheet)

	report, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionAccepted,
		MotionSourcePath: set00Path,
	})
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if !report.AcceptedSource {
		t.Fatalf("single accepted source set reported not accepted")
	}
	if !report.ReleaseReady {
		t.Fatalf("single accepted source set reported not release-ready")
	}
	if report.SourceSets != 1 || report.RuntimeSets != runtimeMotionSets {
		t.Fatalf("source/runtime sets = %d/%d, want 1/%d", report.SourceSets, report.RuntimeSets, runtimeMotionSets)
	}
	if len(report.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none for canonical accepted set00", report.Warnings)
	}
}

func TestValidateVariantRejectsPartialSourceFamily(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	writePNG(t, set00Path, testMotionSheet())
	writePNG(t, filepath.Join(root, "animal-set01-source.png"), testMotionSheet())

	_, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	})
	if err == nil || !strings.Contains(err.Error(), "incomplete motion source family") {
		t.Fatalf("validateVariant() error = %v, want incomplete-family failure", err)
	}
}

func TestValidateVariantRejectsOpaqueMotionFrameBackground(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for y := 0; y < frameH; y++ {
		for x := 0; x < frameW*totalFrames; x++ {
			sheet.SetRGBA(x, y, color.RGBA{R: 230, G: 230, B: 230, A: 255})
		}
	}
	writePNG(t, set00Path, sheet)

	_, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	})
	if err == nil {
		t.Fatalf("validateVariant() succeeded for opaque motion source")
	}
	if !strings.Contains(err.Error(), "transparent background") {
		t.Fatalf("validateVariant() error = %v, want transparent-background failure", err)
	}
}

func testMotionSheet() *image.RGBA {
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for frame := 0; frame < totalFrames; frame++ {
		x0 := frame * frameW
		sheet.SetRGBA(x0+6, 6, color.RGBA{R: 80, G: byte(frame), B: 80, A: 255})
		for y := 28; y < 52; y++ {
			for x := 20; x < 76; x++ {
				sheet.SetRGBA(x0+x, y, color.RGBA{R: 150, G: 150, B: 145, A: 255})
			}
		}
	}
	return sheet
}

func writeMotionSourceFamily(t *testing.T, root string) string {
	t.Helper()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	for set := 0; set < runtimeMotionSets; set++ {
		sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
		for frame := 0; frame < totalFrames; frame++ {
			x0 := frame * frameW
			sheet.SetRGBA(x0+6, 6, color.RGBA{R: byte(set), G: byte(frame), B: 80, A: 255})
			for y := 28; y < 52; y++ {
				for x := 20; x < 76; x++ {
					sheet.SetRGBA(x0+x, y, color.RGBA{R: 150, G: 150, B: 145, A: 255})
				}
			}
		}
		path := filepath.Join(root, "animal-set"+twoDigits(set)+"-source.png")
		writePNG(t, path, sheet)
	}
	return set00Path
}

func writePNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create png: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
}

func twoDigits(v int) string {
	return string([]byte{'0' + byte(v/10), '0' + byte(v%10)})
}
