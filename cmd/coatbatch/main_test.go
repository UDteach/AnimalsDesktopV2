package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAndSliceExactCells(t *testing.T) {
	root := t.TempDir()
	colors := []color.NRGBA{
		{R: 220, G: 40, B: 30, A: 255},
		{R: 40, G: 80, B: 220, A: 255},
		{R: 230, G: 180, B: 30, A: 255},
		{R: 150, G: 70, B: 190, A: 255},
	}
	m := testManifest()
	for i := range m.Cells {
		width, height := 7+i, 5+i
		if i == 0 {
			width, height = cellWidth, cellHeight
		}
		source := image.NewNRGBA(image.Rect(0, 0, width, height))
		draw.Draw(source, source.Bounds(), &image.Uniform{C: colors[i]}, image.Point{}, draw.Src)
		writeTestPNG(t, filepath.Join(root, m.Cells[i].Source), source)
	}
	manifestPath := writeTestManifest(t, root, m)
	builtPath := filepath.Join(root, "built.png")
	buildReportPath := filepath.Join(root, "build-report.json")
	if err := execute(cliOptions{
		Mode:        "build",
		Manifest:    manifestPath,
		Out:         builtPath,
		Report:      buildReportPath,
		TargetRatio: math.NaN(),
	}); err != nil {
		t.Fatalf("build: %v", err)
	}
	built := readTestPNG(t, builtPath)
	if got := built.Bounds().Size(); got != (image.Point{X: sheetWidth, Y: sheetHeight}) {
		t.Fatalf("built dimensions = %v", got)
	}
	for i, want := range colors {
		bounds := boundsForCell(i)
		got := nrgbaAt(built, bounds.Min.X+bounds.Dx()/2, bounds.Min.Y+bounds.Dy()/2)
		if got != want {
			t.Errorf("cell %d center = %#v, want %#v", i, got, want)
		}
	}
	var buildReport coatReport
	readTestJSON(t, buildReportPath, &buildReport)
	assertReportProvenanceAndHashes(t, buildReport, manifestPath)
	if buildReport.OutputSHA256 != fileSHA256(t, builtPath) {
		t.Fatalf("build output hash = %q, written bytes = %q", buildReport.OutputSHA256, fileSHA256(t, builtPath))
	}
	if buildReport.OutputSHA256 == "" || len(buildReport.Cells) != cellCount {
		t.Fatalf("incomplete build report: %#v", buildReport)
	}
	for i, cell := range buildReport.Cells {
		if cell.Cell != i || cell.SourceSHA256 == "" || cell.BaseFrameSHA256 == "" ||
			cell.Bounds != reportRect(boundsForCell(i)) {
			t.Errorf("build report cell %d = %#v", i, cell)
		}
		wantFilter := "golang.org/x/image/draw.CatmullRom"
		if i == 0 {
			wantFilter = "none"
		}
		if cell.BuildFilter != wantFilter {
			t.Errorf("build report cell %d filter = %q, want %q", i, cell.BuildFilter, wantFilter)
		}
	}

	sliceSource := makeSubjectSheet([]color.NRGBA{
		{R: 100, G: 80, B: 60, A: 255},
		{R: 110, G: 90, B: 70, A: 255},
		{R: 120, G: 100, B: 80, A: 255},
		{R: 130, G: 110, B: 90, A: 255},
	})
	sliceSourcePath := filepath.Join(root, "generated.png")
	writeTestPNG(t, sliceSourcePath, sliceSource)
	sliceDir := filepath.Join(root, "slices")
	sliceReportPath := filepath.Join(root, "slice-report.json")
	if err := execute(cliOptions{
		Mode:        "slice",
		Manifest:    manifestPath,
		Sheet:       sliceSourcePath,
		Out:         sliceDir,
		Report:      sliceReportPath,
		TargetRatio: math.NaN(),
	}); err != nil {
		t.Fatalf("slice: %v", err)
	}
	for i, cell := range m.Cells {
		slicedPath := filepath.Join(sliceDir, cell.Output)
		sliced := readTestPNG(t, slicedPath)
		if got := sliced.Bounds().Size(); got != (image.Point{X: cellWidth, Y: cellHeight}) {
			t.Fatalf("slice %d dimensions = %v", i, got)
		}
		got := nrgbaAt(sliced, 300, 220)
		want := nrgbaAt(sliceSource, boundsForCell(i).Min.X+300, boundsForCell(i).Min.Y+220)
		if got != want {
			t.Errorf("slice %d subject pixel = %#v, want %#v", i, got, want)
		}
	}
	var sliceReport coatReport
	readTestJSON(t, sliceReportPath, &sliceReport)
	assertReportProvenanceAndHashes(t, sliceReport, manifestPath)
	if len(sliceReport.Cells) != cellCount {
		t.Fatalf("slice report cell count = %d", len(sliceReport.Cells))
	}
	for i, cell := range sliceReport.Cells {
		if cell.Cell != i || cell.OutputSHA256 == "" || cell.Bounds != reportRect(boundsForCell(i)) {
			t.Errorf("slice report cell %d = %#v", i, cell)
		}
		if cell.OutputSHA256 != fileSHA256(t, filepath.FromSlash(cell.Output)) {
			t.Errorf("slice report cell %d hash does not match written bytes", i)
		}
	}
}

func TestWrongCanvasRejectedBeforeSliceWrites(t *testing.T) {
	root := t.TempDir()
	manifestPath := writeTestManifest(t, root, testManifest())
	wrong := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	wrongPath := filepath.Join(root, "wrong.png")
	writeTestPNG(t, wrongPath, wrong)
	outDir := filepath.Join(root, "out")
	reportPath := filepath.Join(root, "report.json")
	err := execute(cliOptions{
		Mode:        "slice",
		Manifest:    manifestPath,
		Sheet:       wrongPath,
		Out:         outDir,
		Report:      reportPath,
		TargetRatio: math.NaN(),
	})
	if err == nil || !strings.Contains(err.Error(), "want exactly 1536x1024") {
		t.Fatalf("slice wrong canvas error = %v", err)
	}
	assertNotExists(t, outDir)
	assertNotExists(t, reportPath)
}

func TestManifestRejectsUnsafeAndDuplicateFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*coatManifest)
		want   string
	}{
		{
			name: "source escape",
			mutate: func(m *coatManifest) {
				m.Cells[0].Source = filepath.Join("..", "escape.png")
			},
			want: "source path escapes",
		},
		{
			name: "output escape",
			mutate: func(m *coatManifest) {
				m.Cells[0].Output = filepath.Join("..", "escape.png")
			},
			want: "output path escapes",
		},
		{
			name: "prompt escape",
			mutate: func(m *coatManifest) {
				m.Prompt = filepath.Join("..", "prompt.txt")
			},
			want: "prompt path escapes",
		},
		{
			name: "base frame escape",
			mutate: func(m *coatManifest) {
				m.Cells[0].BaseFrame = filepath.Join("..", "base.png")
			},
			want: "base_frame path escapes",
		},
		{
			name: "unsafe call",
			mutate: func(m *coatManifest) {
				m.Call = "../call"
			},
			want: "unsafe id",
		},
		{
			name: "invalid role",
			mutate: func(m *coatManifest) {
				m.Cells[0].Role = "reference"
			},
			want: "role must be target or filler",
		},
		{
			name: "filler output present",
			mutate: func(m *coatManifest) {
				m.Cells[0].Role = "filler"
			},
			want: "filler output must be absent",
		},
		{
			name: "no target",
			mutate: func(m *coatManifest) {
				for i := range m.Cells {
					m.Cells[i].Role = "filler"
					m.Cells[i].Output = ""
				}
			},
			want: "at least one target",
		},
		{
			name: "duplicate id case insensitive",
			mutate: func(m *coatManifest) {
				m.Cells[1].ID = strings.ToUpper(m.Cells[0].ID)
			},
			want: "duplicate cell id",
		},
		{
			name: "duplicate cell",
			mutate: func(m *coatManifest) {
				m.Cells[1].Cell = m.Cells[0].Cell
			},
			want: "duplicate cell index",
		},
		{
			name: "duplicate output normalized",
			mutate: func(m *coatManifest) {
				m.Cells[1].Output = strings.ToUpper(m.Cells[0].Output)
			},
			want: "duplicate output",
		},
		{
			name: "unsafe id",
			mutate: func(m *coatManifest) {
				m.Cells[0].ID = "../coat"
			},
			want: "unsafe id",
		},
		{
			name: "unknown raw geometry policy",
			mutate: func(m *coatManifest) {
				m.RawGeometryPolicy = "unreviewed_geometry_override"
			},
			want: "unknown raw_geometry_policy",
		},
		{
			name: "albino raw geometry policy on wrong species",
			mutate: func(m *coatManifest) {
				m.RawGeometryPolicy = albinoHighContrastRawGeometryPolicy
			},
			want: "valid only for species",
		},
		{
			name: "unknown tone policy",
			mutate: func(m *coatManifest) {
				m.TonePolicy = &coatManifestTonePolicy{ID: "unreviewed_tone_override"}
			},
			want: "unknown tone_policy",
		},
		{
			name: "sable tone policy on wrong species",
			mutate: func(m *coatManifest) {
				m.TonePolicy = sableRecoveryManifestPolicy()
			},
			want: "valid only for species",
		},
		{
			name: "sable tone policy wrong approved input hash",
			mutate: func(m *coatManifest) {
				m.Species = "ferret_sable"
				m.TonePolicy = sableRecoveryManifestPolicy()
				m.TonePolicy.ApprovedInputSHA256 = strings.ToUpper(sableExactRecoveryApprovedInputSHA)
			},
			want: "approved_input_sha256 must be exactly",
		},
		{
			name: "sable tone policy wrong target",
			mutate: func(m *coatManifest) {
				m.Species = "ferret_sable"
				m.TonePolicy = sableRecoveryManifestPolicy()
				m.TonePolicy.TargetRatio += 0.000001
			},
			want: "target_ratio must be exactly",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			m := testManifest()
			tt.mutate(&m)
			path := writeTestManifest(t, root, m)
			_, err := loadManifest(path)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("loadManifest error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestMeasureUsesCanonicalTorsoROIAndGroupMean(t *testing.T) {
	root := t.TempDir()
	manifestPath := writeTestManifest(t, root, testManifest())
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	writeTestPNG(t, basePath, makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255})))
	candidate := makeSubjectSheet(repeatedColor(color.NRGBA{R: 110, G: 110, B: 110, A: 255}))
	candidate.SetNRGBA(200, 100, color.NRGBA{G: 255, A: 255})
	writeTestPNG(t, candidatePath, candidate)
	reportPath := filepath.Join(root, "measure.json")
	if err := execute(cliOptions{
		Mode:        "measure",
		Manifest:    manifestPath,
		BaseSheet:   basePath,
		Sheet:       candidatePath,
		Report:      reportPath,
		TargetRatio: math.NaN(),
	}); err != nil {
		t.Fatalf("measure: %v", err)
	}
	var report coatReport
	readTestJSON(t, reportPath, &report)
	assertReportProvenanceAndHashes(t, report, manifestPath)
	if report.GroupMeanRatio == nil || math.Abs(*report.GroupMeanRatio-1.1) > 1e-9 {
		t.Fatalf("group mean = %v, want 1.1", report.GroupMeanRatio)
	}
	for i, cell := range report.Cells {
		if cell.BaseMedianLuma == nil || cell.CandidateMedianLuma == nil || cell.Ratio == nil {
			t.Fatalf("cell %d missing luma fields: %#v", i, cell)
		}
		if math.Abs(*cell.BaseMedianLuma-100) > 1e-9 ||
			math.Abs(*cell.CandidateMedianLuma-110) > 1e-9 ||
			math.Abs(*cell.Ratio-1.1) > 1e-9 {
			t.Errorf("cell %d luma report = base %.6f candidate %.6f ratio %.6f", i, *cell.BaseMedianLuma, *cell.CandidateMedianLuma, *cell.Ratio)
		}
		if cell.BaseSubjectBounds == nil || cell.TorsoROI == nil {
			t.Errorf("cell %d missing bbox/ROI", i)
		}
		if cell.ReferenceAreaPX == nil || cell.CandidateAreaPX == nil ||
			cell.IntersectionPX == nil || cell.UnionPX == nil ||
			cell.IoU == nil || cell.CentroidDXPX == nil ||
			cell.CentroidDYPX == nil || cell.CentroidShiftPX == nil {
			t.Fatalf("cell %d missing raw geometry fields: %#v", i, cell)
		}
		wantReferenceArea := 320 * 310
		wantCandidateArea := wantReferenceArea
		wantIntersection := wantReferenceArea
		if i == 0 {
			wantCandidateArea--
			wantIntersection--
		}
		if *cell.ReferenceAreaPX != wantReferenceArea ||
			*cell.CandidateAreaPX != wantCandidateArea ||
			*cell.IntersectionPX != wantIntersection ||
			*cell.UnionPX != wantReferenceArea {
			t.Errorf(
				"cell %d raw areas = reference %d candidate %d intersection %d union %d",
				i,
				*cell.ReferenceAreaPX,
				*cell.CandidateAreaPX,
				*cell.IntersectionPX,
				*cell.UnionPX,
			)
		}
		wantIoU := float64(wantIntersection) / float64(wantReferenceArea)
		if *cell.IoU != wantIoU {
			t.Errorf("cell %d IoU = %.17g, want %.17g", i, *cell.IoU, wantIoU)
		}
	}
}

func TestRawSubjectGeometryGateBoundaries(t *testing.T) {
	iouBounds := image.Rect(0, 0, 32, 16)
	iouReference := image.NewNRGBA(iouBounds)
	iouCandidateAtLimit := image.NewNRGBA(iouBounds)
	subjectColor := color.NRGBA{R: 120, G: 100, B: 80, A: 255}
	subjectRect := image.Rect(2, 2, 22, 12)
	draw.Draw(iouReference, subjectRect, &image.Uniform{C: subjectColor}, image.Point{}, draw.Src)
	draw.Draw(iouCandidateAtLimit, subjectRect, &image.Uniform{C: subjectColor}, image.Point{}, draw.Src)
	for _, point := range []image.Point{{X: 11, Y: 6}, {X: 12, Y: 6}, {X: 11, Y: 7}} {
		iouCandidateAtLimit.SetNRGBA(point.X, point.Y, color.NRGBA{})
	}
	iouReference.SetNRGBA(0, 0, color.NRGBA{R: 120, G: 100, B: 80, A: 8})
	iouCandidateAtLimit.SetNRGBA(1, 0, color.NRGBA{G: 255, A: 255})
	atIoULimit, err := compareSubjectGeometry(iouReference, iouCandidateAtLimit, iouBounds)
	if err != nil {
		t.Fatalf("compare IoU boundary: %v", err)
	}
	if atIoULimit.ReferenceArea != 200 || atIoULimit.CandidateArea != 197 ||
		atIoULimit.Intersection != 197 || atIoULimit.Union != 200 ||
		atIoULimit.IoU != minRawSubjectIoU {
		t.Fatalf("IoU boundary metrics = %+v", atIoULimit)
	}
	if err := validateRawSubjectGeometry(rawGeometryReportsForTest(atIoULimit)); err != nil {
		t.Fatalf("exact IoU boundary rejected: %v", err)
	}

	iouCandidateBelowLimit := cloneTestNRGBA(iouCandidateAtLimit)
	iouCandidateBelowLimit.SetNRGBA(12, 7, color.NRGBA{})
	belowIoULimit, err := compareSubjectGeometry(iouReference, iouCandidateBelowLimit, iouBounds)
	if err != nil {
		t.Fatalf("compare below-IoU boundary: %v", err)
	}
	if err := validateRawSubjectGeometry(rawGeometryReportsForTest(belowIoULimit)); err == nil ||
		!strings.Contains(err.Error(), "IoU") {
		t.Fatalf("below-IoU boundary error = %v", err)
	}

	centroidBounds := image.Rect(0, 0, 600, 16)
	centroidReference := image.NewNRGBA(centroidBounds)
	centroidCandidateAtLimit := image.NewNRGBA(centroidBounds)
	draw.Draw(centroidReference, image.Rect(100, 4, 500, 8), &image.Uniform{C: subjectColor}, image.Point{}, draw.Src)
	draw.Draw(centroidCandidateAtLimit, image.Rect(101, 4, 501, 8), &image.Uniform{C: subjectColor}, image.Point{}, draw.Src)
	centroidCandidateAtLimit.SetNRGBA(101, 4, color.NRGBA{})
	centroidCandidateAtLimit.SetNRGBA(501, 4, subjectColor)
	atCentroidLimit, err := compareSubjectGeometry(centroidReference, centroidCandidateAtLimit, centroidBounds)
	if err != nil {
		t.Fatalf("compare centroid boundary: %v", err)
	}
	if atCentroidLimit.CentroidDX != maxRawCentroidShiftPixels ||
		atCentroidLimit.CentroidDY != 0 ||
		atCentroidLimit.CentroidShift != maxRawCentroidShiftPixels {
		t.Fatalf("centroid boundary metrics = %+v", atCentroidLimit)
	}
	if err := validateRawSubjectGeometry(rawGeometryReportsForTest(atCentroidLimit)); err != nil {
		t.Fatalf("exact centroid boundary rejected: %v", err)
	}

	centroidCandidateOverLimit := cloneTestNRGBA(centroidCandidateAtLimit)
	centroidCandidateOverLimit.SetNRGBA(501, 4, color.NRGBA{})
	centroidCandidateOverLimit.SetNRGBA(502, 4, subjectColor)
	overCentroidLimit, err := compareSubjectGeometry(centroidReference, centroidCandidateOverLimit, centroidBounds)
	if err != nil {
		t.Fatalf("compare over-centroid boundary: %v", err)
	}
	if err := validateRawSubjectGeometry(rawGeometryReportsForTest(overCentroidLimit)); err == nil ||
		!strings.Contains(err.Error(), "centroid shift") {
		t.Fatalf("over-centroid boundary error = %v", err)
	}
}

func TestAlbinoHighContrastGeometryPolicyMeasureAndToneBoundaries(t *testing.T) {
	root := t.TempDir()
	m := testManifest()
	m.Species = "ferret_albino"
	m.RawGeometryPolicy = albinoHighContrastRawGeometryPolicy
	manifestPath := writeTestManifest(t, root, m)

	base := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}))
	candidateAtLimit := cloneTestNRGBA(base)
	for cell := 0; cell < cellCount; cell++ {
		bounds := boundsForCell(cell)
		// The canonical subject is 320x310 = 99,200 pixels. Removing this
		// centered 64x31 patch leaves exactly 97,216 pixels, so IoU is 0.980.
		draw.Draw(
			candidateAtLimit,
			image.Rect(bounds.Min.X+328, bounds.Min.Y+240, bounds.Min.X+392, bounds.Min.Y+271),
			&image.Uniform{C: color.NRGBA{G: 255, A: 255}},
			image.Point{},
			draw.Src,
		)
	}
	basePath := filepath.Join(root, "base.png")
	atLimitPath := filepath.Join(root, "candidate-at-limit.png")
	writeTestPNG(t, basePath, base)
	writeTestPNG(t, atLimitPath, candidateAtLimit)

	measureReportPath := filepath.Join(root, "measure-at-limit.json")
	if err := execute(cliOptions{
		Mode: "measure", Manifest: manifestPath, BaseSheet: basePath, Sheet: atLimitPath,
		Report: measureReportPath, TargetRatio: math.NaN(),
	}); err != nil {
		t.Fatalf("albino exact 0.980 measure rejected: %v", err)
	}
	var measureReport coatReport
	readTestJSON(t, measureReportPath, &measureReport)
	assertAlbinoPolicyReport(t, measureReport)
	for i, cell := range measureReport.Cells {
		if cell.IoU == nil || *cell.IoU != minAlbinoHighContrastRawSubjectIoU {
			t.Fatalf("cell %d IoU = %v, want %.3f", i, cell.IoU, minAlbinoHighContrastRawSubjectIoU)
		}
	}

	toneOutPath := filepath.Join(root, "tone-at-limit.png")
	toneReportPath := filepath.Join(root, "tone-at-limit.json")
	if err := execute(cliOptions{
		Mode: "tone", Manifest: manifestPath, BaseSheet: basePath, Sheet: atLimitPath,
		Out: toneOutPath, Report: toneReportPath, TargetRatio: 1,
	}); err != nil {
		t.Fatalf("albino exact 0.980 tone rejected: %v", err)
	}
	var toneReport coatReport
	readTestJSON(t, toneReportPath, &toneReport)
	assertAlbinoPolicyReport(t, toneReport)

	candidateBelowLimit := cloneTestNRGBA(candidateAtLimit)
	candidateBelowLimit.SetNRGBA(200, 100, color.NRGBA{G: 255, A: 255})
	belowLimitPath := filepath.Join(root, "candidate-below-limit.png")
	writeTestPNG(t, belowLimitPath, candidateBelowLimit)
	diagnosticPath := filepath.Join(root, "measure-below-limit.json")
	err := execute(cliOptions{
		Mode: "measure", Manifest: manifestPath, BaseSheet: basePath, Sheet: belowLimitPath,
		Report: diagnosticPath, TargetRatio: math.NaN(),
	})
	if err == nil || !strings.Contains(err.Error(), "IoU") {
		t.Fatalf("albino below-0.980 measure error = %v", err)
	}
	var diagnostic coatReport
	readTestJSON(t, diagnosticPath, &diagnostic)
	assertAlbinoPolicyReport(t, diagnostic)

	strictPolicy := effectiveRawGeometryPolicy(testManifest())
	if err := validateRawSubjectGeometryWithPolicy(
		rawGeometryReportsForTest(subjectGeometry{IoU: minAlbinoHighContrastRawSubjectIoU}),
		strictPolicy,
	); err == nil {
		t.Fatal("strict default unexpectedly accepted IoU 0.980")
	}
}

func TestGeometryGateRejectsMeasureAndToneIncludingFillerWithDiagnosticReport(t *testing.T) {
	root := t.TempDir()
	m := testManifest()
	m.Cells[3].Role = "filler"
	m.Cells[3].Output = ""
	manifestPath := writeTestManifest(t, root, m)
	base := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}))
	candidate := makeSubjectSheet(repeatedColor(color.NRGBA{R: 110, G: 110, B: 110, A: 255}))
	fillerBounds := boundsForCell(3)
	draw.Draw(
		candidate,
		image.Rect(fillerBounds.Min.X+200, fillerBounds.Min.Y+100, fillerBounds.Min.X+210, fillerBounds.Min.Y+410),
		&image.Uniform{C: color.NRGBA{G: 255, A: 255}},
		image.Point{},
		draw.Src,
	)
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	writeTestPNG(t, basePath, base)
	writeTestPNG(t, candidatePath, candidate)

	for _, mode := range []string{"measure", "tone"} {
		t.Run(mode, func(t *testing.T) {
			reportPath := filepath.Join(root, mode+"-rejected.json")
			outPath := ""
			if mode == "tone" {
				outPath = filepath.Join(root, "tone-must-not-exist.png")
			}
			err := execute(cliOptions{
				Mode:        mode,
				Manifest:    manifestPath,
				BaseSheet:   basePath,
				Sheet:       candidatePath,
				Out:         outPath,
				Report:      reportPath,
				TargetRatio: map[string]float64{"measure": math.NaN(), "tone": 1}[mode],
			})
			if err == nil || !strings.Contains(err.Error(), "raw subject geometry gate failed") ||
				!strings.Contains(err.Error(), `cell "coat-3"`) {
				t.Fatalf("%s rejection error = %v", mode, err)
			}
			if outPath != "" {
				assertNotExists(t, outPath)
			}

			var report coatReport
			readTestJSON(t, reportPath, &report)
			assertReportProvenanceAndHashes(t, report, manifestPath)
			if len(report.Cells) != cellCount {
				t.Fatalf("%s diagnostic cell count = %d", mode, len(report.Cells))
			}
			for i, cell := range report.Cells {
				if cell.ReferenceAreaPX == nil || cell.CandidateAreaPX == nil ||
					cell.IntersectionPX == nil || cell.UnionPX == nil ||
					cell.IoU == nil || cell.CentroidDXPX == nil ||
					cell.CentroidDYPX == nil || cell.CentroidShiftPX == nil {
					t.Fatalf("%s diagnostic cell %d missing geometry: %#v", mode, i, cell)
				}
			}
			filler := report.Cells[3]
			if filler.Role != "filler" ||
				*filler.ReferenceAreaPX != 320*310 ||
				*filler.CandidateAreaPX != 310*310 ||
				*filler.IntersectionPX != 310*310 ||
				*filler.UnionPX != 320*310 ||
				*filler.IoU != 310.0/320.0 ||
				*filler.CentroidDXPX != 5 ||
				*filler.CentroidDYPX != 0 ||
				*filler.CentroidShiftPX != 5 {
				t.Fatalf("%s filler diagnostic geometry = %#v", mode, filler)
			}
		})
	}
}

func TestToneAppliesOneGainAndPreservesGreenAndAlpha(t *testing.T) {
	root := t.TempDir()
	m := testManifest()
	manifestPath := writeTestManifest(t, root, m)
	base := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 200}))
	candidateColor := color.NRGBA{R: 90, G: 80, B: 70, A: 200}
	candidate := makeSubjectSheet(repeatedColor(candidateColor))
	candidate.SetNRGBA(10, 10, color.NRGBA{R: 50, G: 60, B: 70, A: 8})
	base.SetNRGBA(10, 10, color.NRGBA{R: 50, G: 60, B: 70, A: 8})
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	writeTestPNG(t, basePath, base)
	writeTestPNG(t, candidatePath, candidate)

	provenance, err := loadManifestProvenance(manifestPath, m)
	if err != nil {
		t.Fatalf("load provenance: %v", err)
	}
	preCells, preMean, err := measureSheetPair(base, candidate, m, provenance)
	if err != nil {
		t.Fatalf("pre-measure: %v", err)
	}
	if len(preCells) != cellCount {
		t.Fatalf("pre-measure cells = %d", len(preCells))
	}
	target := preMean * 1.02
	outPath := filepath.Join(root, "calibrated.png")
	reportPath := filepath.Join(root, "tone.json")
	if err := execute(cliOptions{
		Mode:        "tone",
		Manifest:    manifestPath,
		BaseSheet:   basePath,
		Sheet:       candidatePath,
		Out:         outPath,
		Report:      reportPath,
		TargetRatio: target,
	}); err != nil {
		t.Fatalf("tone: %v", err)
	}
	output := readTestPNG(t, outPath)
	for y := 0; y < sheetHeight; y++ {
		for x := 0; x < sheetWidth; x++ {
			before := nrgbaAt(candidate, x, y)
			after := nrgbaAt(output, x, y)
			if after.A != before.A {
				t.Fatalf("alpha changed at (%d,%d): %d -> %d", x, y, before.A, after.A)
			}
			if isChromaGreen(before) && after != before {
				t.Fatalf("chroma green changed at (%d,%d): %#v -> %#v", x, y, before, after)
			}
		}
	}
	gotSubject := nrgbaAt(output, 300, 220)
	wantSubject := color.NRGBA{
		R: gainedChannel(candidateColor.R, 1.02),
		G: gainedChannel(candidateColor.G, 1.02),
		B: gainedChannel(candidateColor.B, 1.02),
		A: candidateColor.A,
	}
	if gotSubject != wantSubject {
		t.Errorf("gained subject = %#v, want %#v", gotSubject, wantSubject)
	}
	gotGreen := nrgbaAt(output, 1, 1)
	wantGreen := nrgbaAt(candidate, 1, 1)
	if gotGreen != wantGreen {
		t.Errorf("green changed = %#v, want %#v", gotGreen, wantGreen)
	}
	gotLowAlpha := nrgbaAt(output, 10, 10)
	wantLowAlpha := nrgbaAt(candidate, 10, 10)
	if gotLowAlpha != wantLowAlpha {
		t.Errorf("low-alpha non-subject changed = %#v, want %#v", gotLowAlpha, wantLowAlpha)
	}
	for i := 0; i < cellCount; i++ {
		point := boundsForCell(i).Min.Add(image.Pt(300, 220))
		if got := nrgbaAt(output, point.X, point.Y); got != wantSubject {
			t.Errorf("cell %d did not receive uniform gain: %#v", i, got)
		}
	}
	var report coatReport
	readTestJSON(t, reportPath, &report)
	assertReportProvenanceAndHashes(t, report, manifestPath)
	if report.PostSHA256 != fileSHA256(t, outPath) || report.OutputSHA256 != fileSHA256(t, outPath) {
		t.Fatalf("tone report hashes do not match written output")
	}
	if report.Gain == nil || math.Abs(*report.Gain-1.02) > 1e-12 {
		t.Fatalf("reported gain = %v, want 1.02", report.Gain)
	}
	if report.PreSHA256 == "" || report.PostSHA256 == "" || report.PreSHA256 == report.PostSHA256 {
		t.Fatalf("invalid tone hashes: pre=%q post=%q", report.PreSHA256, report.PostSHA256)
	}
	for i, cell := range report.Cells {
		if cell.PostRatio == nil || cell.PostCandidateMedianLuma == nil {
			t.Errorf("cell %d missing post-calibration values: %#v", i, cell)
		}
		if cell.IoU == nil || *cell.IoU != 1 ||
			cell.CentroidShiftPX == nil || *cell.CentroidShiftPX != 0 {
			t.Errorf("cell %d missing successful raw geometry values: %#v", i, cell)
		}
	}
}

func TestToneRejectsGainOverFivePercentWithoutWriting(t *testing.T) {
	root := t.TempDir()
	manifestPath := writeTestManifest(t, root, testManifest())
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	writeTestPNG(t, basePath, makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255})))
	writeTestPNG(t, candidatePath, makeSubjectSheet(repeatedColor(color.NRGBA{R: 80, G: 80, B: 80, A: 255})))
	outPath := filepath.Join(root, "must-not-exist.png")
	reportPath := filepath.Join(root, "must-not-exist.json")
	err := execute(cliOptions{
		Mode:        "tone",
		Manifest:    manifestPath,
		BaseSheet:   basePath,
		Sheet:       candidatePath,
		Out:         outPath,
		Report:      reportPath,
		TargetRatio: 0.90,
	})
	if err == nil || !strings.Contains(err.Error(), "outside allowed range") {
		t.Fatalf("tone rejection error = %v", err)
	}
	assertNotExists(t, outPath)
	assertNotExists(t, reportPath)
}

func TestSliceRejectsEdgeTouchingSubjectBeforeWrites(t *testing.T) {
	root := t.TempDir()
	manifestPath := writeTestManifest(t, root, testManifest())
	sheet := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 90, B: 80, A: 255}))
	sheet.SetNRGBA(0, 200, color.NRGBA{R: 100, G: 90, B: 80, A: 255})
	sheetPath := filepath.Join(root, "edge.png")
	writeTestPNG(t, sheetPath, sheet)
	outDir := filepath.Join(root, "slices")
	reportPath := filepath.Join(root, "report.json")
	err := execute(cliOptions{
		Mode:        "slice",
		Manifest:    manifestPath,
		Sheet:       sheetPath,
		Out:         outDir,
		Report:      reportPath,
		TargetRatio: math.NaN(),
	})
	if err == nil || !strings.Contains(err.Error(), "touching a cell edge") {
		t.Fatalf("slice edge error = %v", err)
	}
	assertNotExists(t, outDir)
	assertNotExists(t, reportPath)
}

func TestTwoTargetsTwoFillersAreMeasuredTonedAndOnlyTargetsSliced(t *testing.T) {
	root := t.TempDir()
	m := testManifest()
	m.Cells[0].Output = ""
	m.Cells[1].Role, m.Cells[1].Output = "filler", ""
	m.Cells[2].Output = filepath.Join("nested", "selected.png")
	m.Cells[3].Role, m.Cells[3].Output = "filler", ""
	manifestPath := writeTestManifest(t, root, m)
	base := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}))
	candidate := makeSubjectSheet(repeatedColor(color.NRGBA{R: 90, G: 80, B: 70, A: 255}))
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	writeTestPNG(t, basePath, base)
	writeTestPNG(t, candidatePath, candidate)

	measurePath := filepath.Join(root, "measure.json")
	if err := execute(cliOptions{
		Mode: "measure", Manifest: manifestPath, BaseSheet: basePath, Sheet: candidatePath,
		Report: measurePath, TargetRatio: math.NaN(),
	}); err != nil {
		t.Fatalf("measure: %v", err)
	}
	var measureReport coatReport
	readTestJSON(t, measurePath, &measureReport)
	for _, index := range []int{1, 3} {
		if measureReport.Cells[index].Role != "filler" ||
			measureReport.Cells[index].Ratio == nil ||
			measureReport.Cells[index].IoU == nil ||
			*measureReport.Cells[index].IoU != 1 ||
			measureReport.Cells[index].CentroidShiftPX == nil ||
			*measureReport.Cells[index].CentroidShiftPX != 0 {
			t.Fatalf("filler %d was not measured: %#v", index, measureReport.Cells[index])
		}
	}

	tonePath := filepath.Join(root, "toned.png")
	toneReportPath := filepath.Join(root, "tone.json")
	_, preMean, err := measureSheetPairForTest(manifestPath, m, base, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if err := execute(cliOptions{
		Mode: "tone", Manifest: manifestPath, BaseSheet: basePath, Sheet: candidatePath,
		Out: tonePath, Report: toneReportPath, TargetRatio: preMean * 1.02,
	}); err != nil {
		t.Fatalf("tone: %v", err)
	}
	toned := readTestPNG(t, tonePath)
	for _, index := range []int{1, 3} {
		point := boundsForCell(index).Min.Add(image.Pt(300, 220))
		want := color.NRGBA{
			R: gainedChannel(90, 1.02), G: gainedChannel(80, 1.02), B: gainedChannel(70, 1.02), A: 255,
		}
		if got := nrgbaAt(toned, point.X, point.Y); got != want {
			t.Fatalf("filler %d was not tone calibrated: got %#v want %#v", index, got, want)
		}
	}

	sliceDir := filepath.Join(root, "slices")
	sliceReportPath := filepath.Join(root, "slice.json")
	if err := execute(cliOptions{
		Mode: "slice", Manifest: manifestPath, Sheet: tonePath, Out: sliceDir,
		Report: sliceReportPath, TargetRatio: math.NaN(),
	}); err != nil {
		t.Fatalf("slice: %v", err)
	}
	if _, err := os.Stat(filepath.Join(sliceDir, "coat-0.png")); err != nil {
		t.Fatalf("default target output missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(sliceDir, "nested", "selected.png")); err != nil {
		t.Fatalf("explicit target output missing: %v", err)
	}
	for _, filler := range []string{"coat-1.png", "coat-3.png"} {
		assertNotExists(t, filepath.Join(sliceDir, filler))
	}
	var sliceReport coatReport
	readTestJSON(t, sliceReportPath, &sliceReport)
	for _, index := range []int{1, 3} {
		if sliceReport.Cells[index].Role != "filler" ||
			sliceReport.Cells[index].Output != "" ||
			sliceReport.Cells[index].OutputSHA256 != "" {
			t.Fatalf("filler %d unexpectedly has slice output: %#v", index, sliceReport.Cells[index])
		}
	}
}

func TestTorsoROIMatchesPythonBankersRoundingCoordinates(t *testing.T) {
	subject := image.Rect(10, 20, 25, 30)
	if got, want := torsoROI(subject), image.Rect(14, 22, 20, 26); got != want {
		t.Fatalf("torsoROI(%v) = %v, want %v", subject, got, want)
	}
}

func TestToneGainBoundaryAcceptedAndRejectedBeforeWrite(t *testing.T) {
	root := t.TempDir()
	m := testManifest()
	manifestPath := writeTestManifest(t, root, m)
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	sheet := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}))
	writeTestPNG(t, basePath, sheet)
	writeTestPNG(t, candidatePath, sheet)
	tests := []struct {
		name   string
		gain   float64
		accept bool
	}{
		{name: "lower exact", gain: 0.95, accept: true},
		{name: "upper exact", gain: 1.05, accept: true},
		{name: "below lower", gain: 0.95 - 1e-6},
		{name: "above upper", gain: 1.05 + 1e-6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPath := filepath.Join(root, tt.name+".png")
			reportPath := filepath.Join(root, tt.name+".json")
			err := execute(cliOptions{
				Mode: "tone", Manifest: manifestPath, BaseSheet: basePath, Sheet: candidatePath,
				Out: outPath, Report: reportPath, TargetRatio: tt.gain,
			})
			if tt.accept {
				if err != nil {
					t.Fatalf("gain %.9f rejected: %v", tt.gain, err)
				}
				if _, err := os.Stat(outPath); err != nil {
					t.Fatalf("accepted output missing: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), "outside allowed range") {
				t.Fatalf("gain %.9f error = %v", tt.gain, err)
			}
			assertNotExists(t, outPath)
			assertNotExists(t, reportPath)
		})
	}
}

func TestSableExactRecoveryPolicyBindingGainAndPixelPreservation(t *testing.T) {
	m := testManifest()
	m.Species = "ferret_sable"
	m.TonePolicy = sableRecoveryManifestPolicy()
	if err := validateManifestPolicies(m); err != nil {
		t.Fatalf("valid recovery manifest rejected: %v", err)
	}
	policy := effectiveTonePolicy(m)
	if err := validateTonePolicyBinding(
		m,
		sableExactRecoveryApprovedInputSHA,
		sableExactRecoveryTargetRatio,
	); err != nil {
		t.Fatalf("exact recovery binding rejected: %v", err)
	}
	if err := validateTonePolicyBinding(m, strings.Repeat("0", 64), sableExactRecoveryTargetRatio); err == nil ||
		!strings.Contains(err.Error(), "does not match approved_input_sha256") {
		t.Fatalf("wrong sheet SHA error = %v", err)
	}
	if err := validateTonePolicyBinding(
		m,
		sableExactRecoveryApprovedInputSHA,
		sableExactRecoveryTargetRatio+0.000001,
	); err == nil || !strings.Contains(err.Error(), "requires -target-ratio exactly") {
		t.Fatalf("wrong target error = %v", err)
	}

	const reviewedGain = 0.850086610341972
	if err := validateToneGain(reviewedGain, policy); err != nil {
		t.Fatalf("reviewed gain rejected: %v", err)
	}
	if err := validateToneGain(sableExactRecoveryToneGainMin-1e-9, policy); err == nil {
		t.Fatal("gain below recovery lower bound accepted")
	}
	if err := validateToneGain(reviewedGain, effectiveTonePolicy(testManifest())); err == nil {
		t.Fatal("strict default accepted recovery-only gain")
	}

	report := newReport("tone", "manifest.json", m, manifestProvenance{})
	setReportActualInputSHA256(&report, m, sableExactRecoveryApprovedInputSHA)
	if report.TonePolicy.ID != sableExactRecoveryTonePolicy ||
		report.TonePolicy.EffectiveGainMin != sableExactRecoveryToneGainMin ||
		report.TonePolicy.EffectiveGainMax != defaultToneGainMax ||
		report.TonePolicy.ApprovedInputSHA256 != sableExactRecoveryApprovedInputSHA ||
		report.TonePolicy.ActualInputSHA256 != sableExactRecoveryApprovedInputSHA ||
		report.TonePolicy.TargetRatio == nil ||
		*report.TonePolicy.TargetRatio != sableExactRecoveryTargetRatio {
		t.Fatalf("recovery policy report = %#v", report.TonePolicy)
	}

	source := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	source.SetNRGBA(0, 0, color.NRGBA{R: 120, G: 100, B: 80, A: 201})
	source.SetNRGBA(1, 0, color.NRGBA{R: 5, G: 255, B: 3, A: 255})
	source.SetNRGBA(2, 0, color.NRGBA{R: 80, G: 70, B: 60, A: 8})
	got := applyUniformSubjectGain(source, reviewedGain)
	wantSubject := color.NRGBA{
		R: gainedChannel(120, reviewedGain),
		G: gainedChannel(100, reviewedGain),
		B: gainedChannel(80, reviewedGain),
		A: 201,
	}
	if pixel := got.NRGBAAt(0, 0); pixel != wantSubject {
		t.Fatalf("recovery subject pixel = %#v, want %#v", pixel, wantSubject)
	}
	for x := 0; x < 3; x++ {
		if before, after := source.NRGBAAt(x, 0), got.NRGBAAt(x, 0); before.A != after.A {
			t.Fatalf("alpha changed at x=%d: %#v -> %#v", x, before, after)
		}
	}
	if got.NRGBAAt(1, 0) != source.NRGBAAt(1, 0) {
		t.Fatalf("recovery changed chroma green: %#v -> %#v", source.NRGBAAt(1, 0), got.NRGBAAt(1, 0))
	}
	if got.NRGBAAt(2, 0) != source.NRGBAAt(2, 0) {
		t.Fatalf("recovery changed low-alpha pixel: %#v -> %#v", source.NRGBAAt(2, 0), got.NRGBAAt(2, 0))
	}
}

func TestSableExactRecoveryToneRejectsUnapprovedSheetHash(t *testing.T) {
	root := t.TempDir()
	m := testManifest()
	m.Species = "ferret_sable"
	m.TonePolicy = sableRecoveryManifestPolicy()
	manifestPath := writeTestManifest(t, root, m)
	basePath := filepath.Join(root, "base.png")
	candidatePath := filepath.Join(root, "candidate.png")
	sheet := makeSubjectSheet(repeatedColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}))
	writeTestPNG(t, basePath, sheet)
	writeTestPNG(t, candidatePath, sheet)
	outPath := filepath.Join(root, "must-not-exist.png")
	reportPath := filepath.Join(root, "must-not-exist.json")
	err := execute(cliOptions{
		Mode: "tone", Manifest: manifestPath, BaseSheet: basePath, Sheet: candidatePath,
		Out: outPath, Report: reportPath, TargetRatio: sableExactRecoveryTargetRatio,
	})
	if err == nil || !strings.Contains(err.Error(), "does not match approved_input_sha256") {
		t.Fatalf("unapproved sheet error = %v", err)
	}
	assertNotExists(t, outPath)
	assertNotExists(t, reportPath)
}

func TestWriteTransactionRollsBackStageFailures(t *testing.T) {
	for _, preexisting := range []bool{false, true} {
		t.Run(map[bool]string{false: "new", true: "preexisting"}[preexisting], func(t *testing.T) {
			root := t.TempDir()
			paths := []string{filepath.Join(root, "first.bin"), filepath.Join(root, "second.bin")}
			if preexisting {
				for i, path := range paths {
					if err := os.WriteFile(path, []byte{byte('a' + i)}, 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			oldFS := txFS
			defer func() { txFS = oldFS }()
			writes := 0
			txFS.WriteTemp = func(file *os.File, data []byte) error {
				writes++
				if writes == 2 {
					return errors.New("injected stage failure")
				}
				return oldFS.WriteTemp(file, data)
			}
			err := writeTransaction([]encodedOutput{
				{Path: paths[0], Data: []byte("new-first")},
				{Path: paths[1], Data: []byte("new-second")},
			})
			if err == nil || !strings.Contains(err.Error(), "injected stage failure") {
				t.Fatalf("transaction error = %v", err)
			}
			assertTransactionDestinations(t, paths, preexisting)
			assertNoTransactionTemps(t, root)
		})
	}
}

func TestWriteTransactionRollsBackCommitFailures(t *testing.T) {
	for _, preexisting := range []bool{false, true} {
		t.Run(map[bool]string{false: "new", true: "preexisting"}[preexisting], func(t *testing.T) {
			root := t.TempDir()
			paths := []string{filepath.Join(root, "first.bin"), filepath.Join(root, "second.bin")}
			if preexisting {
				for i, path := range paths {
					if err := os.WriteFile(path, []byte{byte('a' + i)}, 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			oldFS := txFS
			defer func() { txFS = oldFS }()
			renameCalls := 0
			failAt := 2
			if preexisting {
				failAt = 4
			}
			txFS.Rename = func(oldPath string, newPath string) error {
				renameCalls++
				if renameCalls == failAt {
					return errors.New("injected commit failure")
				}
				return oldFS.Rename(oldPath, newPath)
			}
			err := writeTransaction([]encodedOutput{
				{Path: paths[0], Data: []byte("new-first")},
				{Path: paths[1], Data: []byte("new-second")},
			})
			if err == nil || !strings.Contains(err.Error(), "injected commit failure") {
				t.Fatalf("transaction error = %v", err)
			}
			assertTransactionDestinations(t, paths, preexisting)
			assertNoTransactionTemps(t, root)
		})
	}
}

func TestDistinctPathsDetectHardlinkAliases(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first.dat")
	second := filepath.Join(root, "second.dat")
	if err := os.WriteFile(first, []byte("same inode"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(first, second); err != nil {
		t.Skipf("hardlinks unsupported: %v", err)
	}
	if err := requireDistinctPaths(namedPath{"first", first}, namedPath{"second", second}); err == nil {
		t.Fatal("hardlink aliases were not rejected")
	}
}

func TestRelativePathRejectsADSReservedAndDriveRelativeNames(t *testing.T) {
	for _, path := range []string{
		"coat.png:stream",
		"CON.png",
		`C:coat.png`,
		"nested/AUX.png",
	} {
		t.Run(path, func(t *testing.T) {
			if _, err := validateRelativePNG(path, "test"); err == nil {
				t.Fatalf("unsafe path %q accepted", path)
			}
		})
	}
}

func testManifest() coatManifest {
	return coatManifest{
		Species: "test-species",
		Call:    "call-001",
		Prompt:  "prompt.txt",
		Swatch:  "swatch.png",
		Cells: []coatManifestCell{
			{ID: "coat-0", Role: "target", Source: "source-0.png", BaseFrame: "base-frame-0.png", Cell: 0, Output: "coat-0.png"},
			{ID: "coat-1", Role: "target", Source: "source-1.png", BaseFrame: "base-frame-1.png", Cell: 1, Output: "coat-1.png"},
			{ID: "coat-2", Role: "target", Source: "source-2.png", BaseFrame: "base-frame-2.png", Cell: 2, Output: "coat-2.png"},
			{ID: "coat-3", Role: "target", Source: "source-3.png", BaseFrame: "base-frame-3.png", Cell: 3, Output: "coat-3.png"},
		},
	}
}

func sableRecoveryManifestPolicy() *coatManifestTonePolicy {
	return &coatManifestTonePolicy{
		ID:                  sableExactRecoveryTonePolicy,
		ApprovedInputSHA256: sableExactRecoveryApprovedInputSHA,
		TargetRatio:         sableExactRecoveryTargetRatio,
	}
}

func assertAlbinoPolicyReport(t *testing.T, report coatReport) {
	t.Helper()
	if report.RawGeometryPolicy.ID != albinoHighContrastRawGeometryPolicy ||
		report.RawGeometryPolicy.EffectiveMinIoU != minAlbinoHighContrastRawSubjectIoU ||
		report.RawGeometryPolicy.EffectiveMaxCentroidShiftPX != maxRawCentroidShiftPixels {
		t.Fatalf("albino raw geometry policy report = %#v", report.RawGeometryPolicy)
	}
	if report.TonePolicy.ID != defaultTonePolicyReportID ||
		report.TonePolicy.EffectiveGainMin != defaultToneGainMin ||
		report.TonePolicy.EffectiveGainMax != defaultToneGainMax {
		t.Fatalf("albino tone policy report = %#v", report.TonePolicy)
	}
}

func repeatedColor(value color.NRGBA) []color.NRGBA {
	return []color.NRGBA{value, value, value, value}
}

func makeSubjectSheet(subjectColors []color.NRGBA) *image.NRGBA {
	sheet := image.NewNRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))
	draw.Draw(sheet, sheet.Bounds(), &image.Uniform{C: color.NRGBA{G: 255, A: 255}}, image.Point{}, draw.Src)
	for cell := 0; cell < cellCount; cell++ {
		bounds := boundsForCell(cell)
		subject := image.Rect(bounds.Min.X+200, bounds.Min.Y+100, bounds.Min.X+520, bounds.Min.Y+410)
		draw.Draw(sheet, subject, &image.Uniform{C: subjectColors[cell]}, image.Point{}, draw.Src)
	}
	return sheet
}

func cloneTestNRGBA(src *image.NRGBA) *image.NRGBA {
	clone := image.NewNRGBA(src.Bounds())
	draw.Draw(clone, clone.Bounds(), src, src.Bounds().Min, draw.Src)
	return clone
}

func rawGeometryReportsForTest(metrics subjectGeometry) []cellReport {
	reports := make([]cellReport, cellCount)
	for i := range reports {
		reports[i].ID = fmt.Sprintf("cell-%d", i)
		setSubjectGeometryReport(&reports[i], metrics)
	}
	return reports
}

func writeTestManifest(t *testing.T, root string, manifest coatManifest) string {
	t.Helper()
	if prompt, err := validateRelativeFile(manifest.Prompt, "prompt", ".txt"); err == nil {
		if err := os.WriteFile(filepath.Join(root, prompt), []byte("operator prompt\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if swatchPath, err := validateRelativePNG(manifest.Swatch, "swatch"); err == nil {
		swatch := image.NewNRGBA(image.Rect(0, 0, 2, 2))
		draw.Draw(swatch, swatch.Bounds(), &image.Uniform{C: color.NRGBA{R: 90, G: 80, B: 70, A: 255}}, image.Point{}, draw.Src)
		writeTestPNG(t, filepath.Join(root, swatchPath), swatch)
	}
	for _, cell := range manifest.Cells {
		if baseFramePath, err := validateRelativePNG(cell.BaseFrame, "base_frame"); err == nil {
			baseFrame := image.NewNRGBA(image.Rect(0, 0, 3, 3))
			draw.Draw(baseFrame, baseFrame.Bounds(), &image.Uniform{C: color.NRGBA{R: uint8(40 + cell.Cell), G: 30, B: 20, A: 255}}, image.Point{}, draw.Src)
			writeTestPNG(t, filepath.Join(root, baseFramePath), baseFrame)
		}
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "manifest.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeTestPNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func readTestPNG(t *testing.T, path string) image.Image {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

func readTestJSON(t *testing.T, path string, target any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatal(err)
	}
}

func measureSheetPairForTest(manifestPath string, manifest coatManifest, base image.Image, candidate image.Image) ([]cellReport, float64, error) {
	provenance, err := loadManifestProvenance(manifestPath, manifest)
	if err != nil {
		return nil, 0, err
	}
	return measureSheetPair(base, candidate, manifest, provenance)
}

func assertReportProvenanceAndHashes(t *testing.T, report coatReport, manifestPath string) {
	t.Helper()
	if report.Version != 2 || report.Call != "call-001" {
		t.Fatalf("report provenance header = version %d call %q", report.Version, report.Call)
	}
	if report.RawGeometryPolicy.ID != defaultRawGeometryPolicyReportID ||
		report.RawGeometryPolicy.EffectiveMinIoU != minRawSubjectIoU ||
		report.RawGeometryPolicy.EffectiveMaxCentroidShiftPX != maxRawCentroidShiftPixels {
		t.Fatalf("default raw geometry policy report = %#v", report.RawGeometryPolicy)
	}
	if report.TonePolicy.ID != defaultTonePolicyReportID ||
		report.TonePolicy.EffectiveGainMin != defaultToneGainMin ||
		report.TonePolicy.EffectiveGainMax != defaultToneGainMax ||
		report.TonePolicy.ApprovedInputSHA256 != "" ||
		report.TonePolicy.ActualInputSHA256 != "" ||
		report.TonePolicy.TargetRatio != nil {
		t.Fatalf("default tone policy report = %#v", report.TonePolicy)
	}
	if report.ManifestSHA256 != fileSHA256(t, manifestPath) {
		t.Fatalf("manifest hash = %q, written bytes = %q", report.ManifestSHA256, fileSHA256(t, manifestPath))
	}
	if report.Prompt.SHA256 == "" || report.Prompt.SHA256 != fileSHA256(t, filepath.FromSlash(report.Prompt.Path)) {
		t.Fatalf("prompt hash does not match written bytes: %#v", report.Prompt)
	}
	if report.Swatch.SHA256 == "" || report.Swatch.SHA256 != fileSHA256(t, filepath.FromSlash(report.Swatch.Path)) {
		t.Fatalf("swatch hash does not match written bytes: %#v", report.Swatch)
	}
	for label, artifact := range map[string]*artifactReport{
		"base sheet":   report.BaseSheet,
		"input sheet":  report.InputSheet,
		"output sheet": report.OutputSheet,
	} {
		if artifact != nil && artifact.SHA256 != fileSHA256(t, filepath.FromSlash(artifact.Path)) {
			t.Fatalf("%s hash does not match written bytes: %#v", label, artifact)
		}
	}
	if report.OutputSHA256 != "" && report.OutputSheet != nil && report.OutputSHA256 != report.OutputSheet.SHA256 {
		t.Fatalf("top-level output hash does not match output artifact")
	}
	if report.PreSHA256 != "" && report.InputSheet != nil && report.PreSHA256 != report.InputSheet.SHA256 {
		t.Fatalf("pre hash does not match input artifact")
	}
	if report.PostSHA256 != "" && report.OutputSheet != nil && report.PostSHA256 != report.OutputSheet.SHA256 {
		t.Fatalf("post hash does not match output artifact")
	}
	for i, cell := range report.Cells {
		if cell.BaseFrameSHA256 == "" ||
			cell.BaseFrameSHA256 != fileSHA256(t, filepath.FromSlash(cell.BaseFrame)) {
			t.Fatalf("cell %d base frame hash does not match written bytes: %#v", i, cell)
		}
		if cell.SourceSHA256 != "" && cell.SourceSHA256 != fileSHA256(t, filepath.FromSlash(cell.Source)) {
			t.Fatalf("cell %d source hash does not match written bytes: %#v", i, cell)
		}
		if cell.OutputSHA256 != "" && cell.OutputSHA256 != fileSHA256(t, filepath.FromSlash(cell.Output)) {
			t.Fatalf("cell %d output hash does not match written bytes: %#v", i, cell)
		}
	}
}

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum)
}

func assertTransactionDestinations(t *testing.T, paths []string, preexisting bool) {
	t.Helper()
	for i, path := range paths {
		if !preexisting {
			assertNotExists(t, path)
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read restored %s: %v", path, err)
		}
		if want := []byte{byte('a' + i)}; !bytes.Equal(data, want) {
			t.Fatalf("restored %s = %q, want %q", path, data, want)
		}
	}
}

func assertNoTransactionTemps(t *testing.T, root string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(root, ".coatbatch-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("transaction temp files remain: %v", matches)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		t.Fatalf("%s exists or stat failed unexpectedly: %v", path, err)
	}
}
