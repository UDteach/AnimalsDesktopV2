package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareFrameKeepsTransparentSource(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "transparent.png")
	img := image.NewRGBA(image.Rect(0, 0, 160, 100))
	for y := 30; y < 82; y++ {
		for x := 28; x < 140; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 140, G: 140, B: 136, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	outPath := filepath.Join(root, "frame.png")
	report, err := prepareFrame(srcPath, outPath, 18)
	if err != nil {
		t.Fatalf("prepareFrame() error = %v", err)
	}
	if report.BackgroundMode != "source-alpha" || report.BackgroundRemoved {
		t.Fatalf("background report = %s/%v", report.BackgroundMode, report.BackgroundRemoved)
	}
	out := openTestPNG(t, outPath)
	if got := out.Bounds(); got.Dx() != frameW || got.Dy() != frameH {
		t.Fatalf("output bounds = %v", got)
	}
	if alphaBounds(out, out.Bounds()).Empty() {
		t.Fatalf("output has no visible alpha")
	}
	if out.RGBAAt(0, 0).A != 0 {
		t.Fatalf("output corner alpha = %d, want transparent", out.RGBAAt(0, 0).A)
	}
}

func TestPrepareFrameRemovesUniformOpaqueBackground(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "opaque.png")
	img := image.NewRGBA(image.Rect(0, 0, 180, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 12, G: 240, B: 12, A: 255})
		}
	}
	for y := 35; y < 95; y++ {
		for x := 42; x < 150; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	outPath := filepath.Join(root, "frame.png")
	report, err := prepareFrame(srcPath, outPath, 18)
	if err != nil {
		t.Fatalf("prepareFrame() error = %v", err)
	}
	if report.BackgroundMode != "uniform-edge-rgb" || !report.BackgroundRemoved {
		t.Fatalf("background report = %s/%v", report.BackgroundMode, report.BackgroundRemoved)
	}
	out := openTestPNG(t, outPath)
	if out.RGBAAt(0, 0).A != 0 {
		t.Fatalf("output corner alpha = %d, want transparent", out.RGBAAt(0, 0).A)
	}
	if alphaBounds(out, out.Bounds()).Empty() {
		t.Fatalf("output has no visible alpha")
	}
}

func TestPrepareFrameRemovesChromaGreenBackground(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "chroma.png")
	img := image.NewRGBA(image.Rect(0, 0, 180, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			g := byte(210 + (x+y)%35)
			img.SetRGBA(x, y, color.RGBA{R: 8, G: g, B: 12, A: 255})
		}
	}
	for y := 35; y < 95; y++ {
		for x := 42; x < 150; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	outPath := filepath.Join(root, "frame.png")
	report, err := prepareFrameWithMode(srcPath, outPath, "chroma-green", 18)
	if err != nil {
		t.Fatalf("prepareFrameWithMode() error = %v", err)
	}
	if report.BackgroundMode != "chroma-green" || !report.BackgroundRemoved {
		t.Fatalf("background report = %s/%v", report.BackgroundMode, report.BackgroundRemoved)
	}
	out := openTestPNG(t, outPath)
	if out.RGBAAt(0, 0).A != 0 {
		t.Fatalf("output corner alpha = %d, want transparent", out.RGBAAt(0, 0).A)
	}
	if alphaBounds(out, out.Bounds()).Empty() {
		t.Fatalf("output has no visible alpha")
	}
}

func TestPrepareFrameRejectsChromaGreenPinholes(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "chroma-pinholes.png")
	img := image.NewRGBA(image.Rect(0, 0, 180, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	for y := 35; y < 95; y++ {
		for x := 42; x < 150; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	for y := 58; y < 60; y++ {
		for x := 86; x < 88; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrameWithMode(srcPath, filepath.Join(root, "frame.png"), "chroma-green", 18)
	if err == nil {
		t.Fatalf("prepareFrameWithMode() succeeded for chroma pinholes")
	}
	if !strings.Contains(err.Error(), "transparent pinholes") {
		t.Fatalf("prepareFrameWithMode() error = %v, want transparent pinholes failure", err)
	}
}

func TestPrepareFrameRejectsCheckerBackground(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "checker.png")
	img := image.NewRGBA(image.Rect(0, 0, 120, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 120; x++ {
			if (x/8+y/8)%2 == 0 {
				img.SetRGBA(x, y, color.RGBA{R: 238, G: 238, B: 238, A: 255})
			} else {
				img.SetRGBA(x, y, color.RGBA{R: 190, G: 190, B: 190, A: 255})
			}
		}
	}
	for y := 24; y < 60; y++ {
		for x := 34; x < 92; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrame(srcPath, filepath.Join(root, "frame.png"), 18)
	if err == nil {
		t.Fatalf("prepareFrame() succeeded for checker background")
	}
	if !strings.Contains(err.Error(), "checker") {
		t.Fatalf("prepareFrame() error = %v, want checker/noisy failure", err)
	}
}

func TestPrepareFrameRejectsUnknownBackgroundMode(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "transparent.png")
	img := image.NewRGBA(image.Rect(0, 0, 120, 80))
	img.SetRGBA(40, 40, color.RGBA{R: 130, G: 130, B: 128, A: 255})
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrameWithMode(srcPath, filepath.Join(root, "frame.png"), "magic", 18)
	if err == nil {
		t.Fatalf("prepareFrameWithMode() succeeded for unknown mode")
	}
	if !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("prepareFrameWithMode() error = %v, want unknown-mode failure", err)
	}
}

func TestDespillGreenNeutralizesRemainingGreenEdges(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.SetRGBA(0, 0, color.RGBA{R: 50, G: 120, B: 55, A: 255})
	img.SetRGBA(1, 0, color.RGBA{R: 80, G: 82, B: 78, A: 255})

	despillGreen(img)

	if got := img.RGBAAt(0, 0); got.G != 55 {
		t.Fatalf("green spill pixel = %#v, want green channel clamped to max red/blue", got)
	}
	if got := img.RGBAAt(1, 0); got.G != 82 {
		t.Fatalf("neutral pixel changed to %#v", got)
	}
}

func TestClearTransparentRGBRemovesHiddenChroma(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.SetRGBA(0, 0, color.RGBA{R: 0, G: 255, B: 0, A: 0})

	clearTransparentRGB(img)

	if got := img.RGBAAt(0, 0); got != (color.RGBA{}) {
		t.Fatalf("transparent pixel = %#v, want zero RGB", got)
	}
}

func TestPrepareFrameRejectsIncompleteBackgroundRemoval(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "green-gradient.png")
	img := image.NewRGBA(image.Rect(0, 0, 140, 90))
	for y := 0; y < 90; y++ {
		for x := 0; x < 140; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 0, G: 240, B: 0, A: 255})
		}
	}
	for y := 8; y < 82; y++ {
		for x := 0; x < 9; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 0, G: 160, B: 0, A: 255})
		}
	}
	for y := 30; y < 62; y++ {
		for x := 42; x < 98; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrame(srcPath, filepath.Join(root, "frame.png"), 18)
	if err == nil {
		t.Fatalf("prepareFrame() succeeded for incomplete background removal")
	}
	if !strings.Contains(err.Error(), "content touches source canvas edge") {
		t.Fatalf("prepareFrame() error = %v, want source-edge failure", err)
	}
}

func TestPrepareFrameMatchAlphaBBoxAdjustsOnePixel(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	baselinePath := filepath.Join(root, "baseline.png")
	if _, err := prepareFrame(srcPath, baselinePath, 18); err != nil {
		t.Fatalf("prepare baseline: %v", err)
	}
	baseline := openTestPNG(t, baselinePath)
	baselineContent := alphaBounds(baseline, baseline.Bounds())
	referenceContent := image.Rect(
		baselineContent.Min.X,
		baselineContent.Min.Y,
		baselineContent.Max.X+1,
		baselineContent.Max.Y,
	)
	referencePath := filepath.Join(root, "reference.png")
	writeBBoxReference(t, referencePath, referenceContent)

	outPath := filepath.Join(root, "adjusted.png")
	report, err := prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: referencePath,
	})
	if err != nil {
		t.Fatalf("prepareFrameWithOptions() error = %v", err)
	}

	out := openTestPNG(t, outPath)
	if got := alphaBounds(out, out.Bounds()); got != referenceContent {
		t.Fatalf("adjusted alpha bounds = %v, want %v", got, referenceContent)
	}
	if report.OutputContent != rectToJSON(referenceContent) {
		t.Fatalf("report output content = %+v, want %+v", report.OutputContent, rectToJSON(referenceContent))
	}
	if report.GeometryAdjusted == nil || !*report.GeometryAdjusted {
		t.Fatalf("geometry adjusted = %v, want true", report.GeometryAdjusted)
	}
	assertSuccessfulGeometryReport(t, report)
}

func TestPrepareFrameMatchAlphaBBoxExactMatchIsPixelIdentical(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	referencePath := filepath.Join(root, "reference.png")
	if _, err := prepareFrame(srcPath, referencePath, 18); err != nil {
		t.Fatalf("prepare reference: %v", err)
	}
	reference := openTestPNG(t, referencePath)

	outPath := filepath.Join(root, "locked.png")
	report, err := prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: referencePath,
	})
	if err != nil {
		t.Fatalf("prepareFrameWithOptions() error = %v", err)
	}
	out := openTestPNG(t, outPath)
	if out.Rect != reference.Rect || out.Stride != reference.Stride || !bytes.Equal(out.Pix, reference.Pix) {
		t.Fatalf("exact geometry match changed output pixels")
	}
	if report.GeometryAdjusted == nil || *report.GeometryAdjusted {
		t.Fatalf("geometry adjusted = %v, want false", report.GeometryAdjusted)
	}
	assertSuccessfulGeometryReport(t, report)
}

func TestNormalizedAlphaGeometryGateBoundaries(t *testing.T) {
	iouBounds := image.Rect(0, 0, 12, 12)
	iouReference := filledRectImage(iouBounds, image.Rect(1, 1, 11, 11))
	iouCandidateAtLimit := cloneTestRGBA(iouReference)
	iouCandidateAtLimit.SetRGBA(5, 5, color.RGBA{})
	iouCandidateAtLimit.SetRGBA(6, 6, color.RGBA{})
	iouReference.SetRGBA(0, 0, color.RGBA{R: 100, G: 90, B: 80, A: 8})
	iouCandidateAtLimit.SetRGBA(0, 1, color.RGBA{R: 100, G: 90, B: 80, A: 8})
	atIoULimit, err := compareAlphaGeometry(iouReference, iouCandidateAtLimit, iouBounds)
	if err != nil {
		t.Fatalf("compare IoU boundary: %v", err)
	}
	if atIoULimit.ReferenceArea != 100 || atIoULimit.CandidateArea != 98 ||
		atIoULimit.Intersection != 98 || atIoULimit.Union != 100 ||
		atIoULimit.IoU != minNormalizedAlphaIoU ||
		atIoULimit.CentroidShift != 0 {
		t.Fatalf("IoU boundary metrics = %+v", atIoULimit)
	}
	if err := validateNormalizedAlphaGeometry(atIoULimit); err != nil {
		t.Fatalf("exact IoU boundary rejected: %v", err)
	}

	iouCandidateBelowLimit := cloneTestRGBA(iouCandidateAtLimit)
	iouCandidateBelowLimit.SetRGBA(5, 6, color.RGBA{})
	belowIoULimit, err := compareAlphaGeometry(iouReference, iouCandidateBelowLimit, iouBounds)
	if err != nil {
		t.Fatalf("compare below-IoU boundary: %v", err)
	}
	if err := validateNormalizedAlphaGeometry(belowIoULimit); err == nil ||
		!strings.Contains(err.Error(), "IoU") {
		t.Fatalf("below-IoU boundary error = %v", err)
	}

	centroidBounds := image.Rect(0, 0, 50, 24)
	centroidReference := filledRectImage(centroidBounds, image.Rect(10, 10, 30, 15))
	centroidCandidateAtLimit := cloneTestRGBA(centroidReference)
	centroidCandidateAtLimit.SetRGBA(10, 12, color.RGBA{})
	centroidCandidateAtLimit.SetRGBA(40, 12, color.RGBA{R: 130, G: 130, B: 128, A: 255})
	atCentroidLimit, err := compareAlphaGeometry(centroidReference, centroidCandidateAtLimit, centroidBounds)
	if err != nil {
		t.Fatalf("compare centroid boundary: %v", err)
	}
	if atCentroidLimit.IoU < minNormalizedAlphaIoU ||
		atCentroidLimit.CentroidDX != maxNormalizedCentroidShiftPixels ||
		atCentroidLimit.CentroidDY != 0 ||
		atCentroidLimit.CentroidShift != maxNormalizedCentroidShiftPixels {
		t.Fatalf("centroid boundary metrics = %+v", atCentroidLimit)
	}
	if err := validateNormalizedAlphaGeometry(atCentroidLimit); err != nil {
		t.Fatalf("exact centroid boundary rejected: %v", err)
	}

	centroidCandidateOverLimit := cloneTestRGBA(centroidCandidateAtLimit)
	centroidCandidateOverLimit.SetRGBA(40, 12, color.RGBA{})
	centroidCandidateOverLimit.SetRGBA(41, 12, color.RGBA{R: 130, G: 130, B: 128, A: 255})
	overCentroidLimit, err := compareAlphaGeometry(centroidReference, centroidCandidateOverLimit, centroidBounds)
	if err != nil {
		t.Fatalf("compare over-centroid boundary: %v", err)
	}
	if err := validateNormalizedAlphaGeometry(overCentroidLimit); err == nil ||
		!strings.Contains(err.Error(), "centroid shift") {
		t.Fatalf("over-centroid boundary error = %v", err)
	}
}

func TestPrepareFrameMatchAlphaBBoxRejectsMaskDriftBeforeOutputWrite(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	baselinePath := filepath.Join(root, "baseline.png")
	if _, err := prepareFrame(srcPath, baselinePath, 18); err != nil {
		t.Fatalf("prepare baseline: %v", err)
	}
	reference := openTestPNG(t, baselinePath)
	referenceContent := alphaBounds(reference, reference.Bounds())
	for y := referenceContent.Min.Y; y < referenceContent.Max.Y; y++ {
		reference.SetRGBA(referenceContent.Min.X+20, y, color.RGBA{})
		reference.SetRGBA(referenceContent.Min.X+21, y, color.RGBA{})
	}
	referencePath := filepath.Join(root, "drifted-reference.png")
	writeTestPNG(t, referencePath, reference)

	outPath := filepath.Join(root, "must-not-exist.png")
	_, err := prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: referencePath,
	})
	if err == nil || !strings.Contains(err.Error(), "normalized geometry gate failed") ||
		!strings.Contains(err.Error(), "IoU") {
		t.Fatalf("mask-drift rejection error = %v", err)
	}
	if _, statErr := os.Stat(outPath); !os.IsNotExist(statErr) {
		t.Fatalf("mask-drift rejection wrote output: stat error = %v", statErr)
	}
}

func TestPrepareFrameMatchAlphaBBoxRejectsDeltaOverOnePixelWithoutOutput(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	baselinePath := filepath.Join(root, "baseline.png")
	if _, err := prepareFrame(srcPath, baselinePath, 18); err != nil {
		t.Fatalf("prepare baseline: %v", err)
	}
	baseline := openTestPNG(t, baselinePath)
	baselineContent := alphaBounds(baseline, baseline.Bounds())
	referenceContent := baselineContent.Add(image.Pt(2, 0))
	referencePath := filepath.Join(root, "reference.png")
	writeBBoxReference(t, referencePath, referenceContent)

	outPath := filepath.Join(root, "should-not-exist.png")
	_, err := prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: referencePath,
	})
	if err == nil {
		t.Fatalf("prepareFrameWithOptions() succeeded for 2px geometry delta")
	}
	if !strings.Contains(err.Error(), "more than 1px") {
		t.Fatalf("prepareFrameWithOptions() error = %v, want 1px limit failure", err)
	}
	if _, statErr := os.Stat(outPath); !os.IsNotExist(statErr) {
		t.Fatalf("rejected geometry wrote output: stat error = %v", statErr)
	}
}

func TestPrepareFrameMatchAlphaBBoxRejectsInvalidReferences(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	tests := []struct {
		name      string
		image     *image.RGBA
		wantError string
	}{
		{
			name:      "wrong dimensions",
			image:     filledRectImage(image.Rect(0, 0, frameW-1, frameH), image.Rect(10, 10, 40, 40)),
			wantError: "want 96x64",
		},
		{
			name:      "empty alpha",
			image:     image.NewRGBA(image.Rect(0, 0, frameW, frameH)),
			wantError: "no visible alpha",
		},
		{
			name:      "no transparent margin",
			image:     filledRectImage(image.Rect(0, 0, frameW, frameH), image.Rect(0, 0, frameW, frameH)),
			wantError: "no transparent margin",
		},
		{
			name:      "edge contact",
			image:     filledRectImage(image.Rect(0, 0, frameW, frameH), image.Rect(0, 10, 30, 40)),
			wantError: "touches canvas edge",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			referencePath := filepath.Join(root, strings.ReplaceAll(test.name, " ", "-")+".png")
			writeTestPNG(t, referencePath, test.image)
			outPath := filepath.Join(root, strings.ReplaceAll(test.name, " ", "-")+"-out.png")
			_, err := prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
				Background:     "auto",
				Tolerance:      18,
				MatchAlphaBBox: referencePath,
			})
			if err == nil {
				t.Fatalf("prepareFrameWithOptions() succeeded for invalid reference")
			}
			if !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("prepareFrameWithOptions() error = %v, want %q", err, test.wantError)
			}
			if _, statErr := os.Stat(outPath); !os.IsNotExist(statErr) {
				t.Fatalf("invalid reference wrote output: stat error = %v", statErr)
			}
		})
	}
}

func TestPrepareFrameMatchAlphaBBoxRejectsOutputCollision(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	referencePath := filepath.Join(root, "reference.png")
	if _, err := prepareFrame(srcPath, referencePath, 18); err != nil {
		t.Fatalf("prepare reference: %v", err)
	}
	before, err := os.ReadFile(referencePath)
	if err != nil {
		t.Fatalf("read reference before collision test: %v", err)
	}

	_, err = prepareFrameWithOptions(srcPath, filepath.Join(root, ".", "reference.png"), prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: referencePath,
	})
	if err == nil {
		t.Fatalf("prepareFrameWithOptions() succeeded with reference/output collision")
	}
	if !strings.Contains(err.Error(), "same path") {
		t.Fatalf("prepareFrameWithOptions() error = %v, want same-path failure", err)
	}
	after, readErr := os.ReadFile(referencePath)
	if readErr != nil {
		t.Fatalf("read reference after collision test: %v", readErr)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("reference/output collision changed reference file")
	}
}

func TestPrepareFrameMatchAlphaBBoxReportsProvenance(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	baselinePath := filepath.Join(root, "baseline.png")
	if _, err := prepareFrame(srcPath, baselinePath, 18); err != nil {
		t.Fatalf("prepare baseline: %v", err)
	}
	baseline := openTestPNG(t, baselinePath)
	preLockContent := alphaBounds(baseline, baseline.Bounds())
	referenceContent := image.Rect(
		preLockContent.Min.X,
		preLockContent.Min.Y,
		preLockContent.Max.X+1,
		preLockContent.Max.Y,
	)
	referencePath := filepath.Join(root, "reference.png")
	writeBBoxReference(t, referencePath, referenceContent)

	report, err := prepareFrameWithOptions(srcPath, filepath.Join(root, "locked.png"), prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: referencePath,
	})
	if err != nil {
		t.Fatalf("prepareFrameWithOptions() error = %v", err)
	}
	if report.GeometryMode != "match-alpha-bbox" {
		t.Fatalf("geometry mode = %q", report.GeometryMode)
	}
	if report.GeometryReference != filepath.ToSlash(referencePath) {
		t.Fatalf("geometry reference = %q, want %q", report.GeometryReference, filepath.ToSlash(referencePath))
	}
	if report.PreLockOutputContent == nil || *report.PreLockOutputContent != rectToJSON(preLockContent) {
		t.Fatalf("pre-lock content = %+v, want %+v", report.PreLockOutputContent, rectToJSON(preLockContent))
	}
	if report.ReferenceContent == nil || *report.ReferenceContent != rectToJSON(referenceContent) {
		t.Fatalf("reference content = %+v, want %+v", report.ReferenceContent, rectToJSON(referenceContent))
	}
	if report.GeometryAdjusted == nil || !*report.GeometryAdjusted {
		t.Fatalf("geometry adjusted = %v, want true", report.GeometryAdjusted)
	}
	assertSuccessfulGeometryReport(t, report)

	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal report: %v", err)
	}
	for _, field := range []string{
		`"geometry_mode":"match-alpha-bbox"`,
		`"geometry_reference":`,
		`"pre_lock_output_content":`,
		`"reference_content":`,
		`"geometry_adjusted":true`,
		`"reference_area_px":`,
		`"candidate_area_px":`,
		`"intersection_px":`,
		`"union_px":`,
		`"iou":1`,
		`"centroid_dx_px":0`,
		`"centroid_dy_px":0`,
		`"centroid_shift_px":0`,
	} {
		if !bytes.Contains(encoded, []byte(field)) {
			t.Fatalf("report JSON %s does not contain %s", encoded, field)
		}
	}

	plainReport, err := prepareFrame(srcPath, filepath.Join(root, "plain.png"), 18)
	if err != nil {
		t.Fatalf("prepareFrame() error = %v", err)
	}
	plainEncoded, err := json.Marshal(plainReport)
	if err != nil {
		t.Fatalf("marshal plain report: %v", err)
	}
	for _, field := range []string{
		"geometry_mode",
		"geometry_reference",
		"pre_lock_output_content",
		"reference_content",
		"geometry_adjusted",
		"reference_area_px",
		"candidate_area_px",
		"intersection_px",
		"union_px",
		"iou",
		"centroid_dx_px",
		"centroid_dy_px",
		"centroid_shift_px",
	} {
		if bytes.Contains(plainEncoded, []byte(field)) {
			t.Fatalf("plain report JSON %s unexpectedly contains %q", plainEncoded, field)
		}
	}
}

func TestPrepareFrameRejectsSourceOutputCollisionWithoutChangingSource(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	before, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("read source before collision test: %v", err)
	}

	_, err = prepareFrame(srcPath, filepath.Join(root, ".", filepath.Base(srcPath)), 18)
	if err == nil {
		t.Fatalf("prepareFrame() succeeded with source/output collision")
	}
	if !strings.Contains(err.Error(), "same path") {
		t.Fatalf("prepareFrame() error = %v, want same-path failure", err)
	}
	after, readErr := os.ReadFile(srcPath)
	if readErr != nil {
		t.Fatalf("read source after collision test: %v", readErr)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("source/output collision changed source file")
	}
}

func TestRequireDistinctPreparePathsRejectsReportCollisionMatrix(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	outPath := filepath.Join(root, "output.png")
	referencePath := filepath.Join(root, "reference.png")
	writeBBoxReference(t, referencePath, image.Rect(8, 8, 88, 56))

	tests := []struct {
		name   string
		report string
	}{
		{name: "source", report: srcPath},
		{name: "output", report: outPath},
		{name: "reference", report: referencePath},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := requireDistinctPreparePaths(srcPath, outPath, test.report, referencePath)
			if err == nil {
				t.Fatalf("requireDistinctPreparePaths() accepted report/%s collision", test.name)
			}
			if !strings.Contains(err.Error(), "same path") {
				t.Fatalf("requireDistinctPreparePaths() error = %v, want same-path failure", err)
			}
		})
	}
}

func TestPrepareFrameReportsSourceReferenceAndOutputHashes(t *testing.T) {
	root := t.TempDir()
	srcPath := writeGeometryTestSource(t, root)
	plainPath := filepath.Join(root, "plain.png")
	if _, err := prepareFrame(srcPath, plainPath, 18); err != nil {
		t.Fatalf("prepare plain reference: %v", err)
	}
	outPath := filepath.Join(root, "locked.png")
	report, err := prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
		Background:     "auto",
		Tolerance:      18,
		MatchAlphaBBox: plainPath,
	})
	if err != nil {
		t.Fatalf("prepareFrameWithOptions() error = %v", err)
	}

	for name, pair := range map[string][2]string{
		"source":    {report.SourceSHA256, srcPath},
		"reference": {report.GeometryReferenceSHA, plainPath},
		"output":    {report.OutputSHA256, outPath},
	} {
		want, hashErr := fileSHA256(pair[1])
		if hashErr != nil {
			t.Fatalf("hash %s: %v", name, hashErr)
		}
		if pair[0] != want {
			t.Fatalf("%s hash = %q, want %q", name, pair[0], want)
		}
	}
}

func TestWriteFileTransactionRollsBackExistingFilesOnCommitFailure(t *testing.T) {
	root := t.TempDir()
	firstPath := filepath.Join(root, "first.bin")
	secondPath := filepath.Join(root, "second.bin")
	if err := os.WriteFile(firstPath, []byte("old-first"), 0o644); err != nil {
		t.Fatalf("write first fixture: %v", err)
	}
	if err := os.WriteFile(secondPath, []byte("old-second"), 0o644); err != nil {
		t.Fatalf("write second fixture: %v", err)
	}

	originalRename := renameFile
	renameCalls := 0
	renameFile = func(oldPath string, newPath string) error {
		renameCalls++
		if renameCalls == 4 {
			return errors.New("injected second commit failure")
		}
		return os.Rename(oldPath, newPath)
	}
	t.Cleanup(func() { renameFile = originalRename })

	err := writeFileTransaction([]filePayload{
		{Path: firstPath, Data: []byte("new-first")},
		{Path: secondPath, Data: []byte("new-second")},
	})
	if err == nil || !strings.Contains(err.Error(), "injected second commit failure") {
		t.Fatalf("writeFileTransaction() error = %v, want injected failure", err)
	}
	for path, want := range map[string]string{
		firstPath:  "old-first",
		secondPath: "old-second",
	} {
		got, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("read restored %s: %v", path, readErr)
		}
		if string(got) != want {
			t.Fatalf("restored %s = %q, want %q", path, got, want)
		}
	}
	leftovers, globErr := filepath.Glob(filepath.Join(root, ".*.*-*"))
	if globErr != nil {
		t.Fatalf("glob transaction leftovers: %v", globErr)
	}
	if len(leftovers) != 0 {
		t.Fatalf("transaction left temporary files: %v", leftovers)
	}
}

func assertSuccessfulGeometryReport(t *testing.T, report prepareReport) {
	t.Helper()
	if report.ReferenceAreaPX == nil || report.CandidateAreaPX == nil ||
		report.IntersectionPX == nil || report.UnionPX == nil ||
		report.IoU == nil || report.CentroidDXPX == nil ||
		report.CentroidDYPX == nil || report.CentroidShiftPX == nil {
		t.Fatalf("geometry report is incomplete: %#v", report)
	}
	if *report.ReferenceAreaPX <= 0 ||
		*report.ReferenceAreaPX != *report.CandidateAreaPX ||
		*report.ReferenceAreaPX != *report.IntersectionPX ||
		*report.ReferenceAreaPX != *report.UnionPX ||
		*report.IoU != 1 ||
		*report.CentroidDXPX != 0 ||
		*report.CentroidDYPX != 0 ||
		*report.CentroidShiftPX != 0 {
		t.Fatalf("successful geometry metrics = %#v", report)
	}
}

func cloneTestRGBA(src *image.RGBA) *image.RGBA {
	clone := image.NewRGBA(src.Bounds())
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			clone.SetRGBA(x, y, src.RGBAAt(x, y))
		}
	}
	return clone
}

func writeGeometryTestSource(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "geometry-source.png")
	img := filledRectImage(image.Rect(0, 0, 160, 100), image.Rect(30, 20, 130, 80))
	writeTestPNG(t, path, img)
	return path
}

func writeBBoxReference(t *testing.T, path string, content image.Rectangle) {
	t.Helper()
	writeTestPNG(t, path, filledRectImage(image.Rect(0, 0, frameW, frameH), content))
}

func filledRectImage(bounds image.Rectangle, content image.Rectangle) *image.RGBA {
	img := image.NewRGBA(bounds)
	for y := content.Min.Y; y < content.Max.Y; y++ {
		for x := content.Min.X; x < content.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	return img
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
