package main

import (
	"image"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"animals-desktop/internal/catalog"
)

func TestSelectImportVariantsTargetsOneSeedVariant(t *testing.T) {
	seeds := catalog.SeedVariants()
	if len(seeds) == 0 {
		t.Fatal("catalog has no seed variants")
	}

	got, err := selectImportVariants(seeds[0].ID)
	if err != nil {
		t.Fatalf("selectImportVariants() error = %v", err)
	}
	if len(got) != 1 || got[0].ID != seeds[0].ID {
		t.Fatalf("selected variants = %#v, want exactly %q", got, seeds[0].ID)
	}
}

func TestRunImporterHelpSucceedsWithoutImporting(t *testing.T) {
	var stderr strings.Builder
	if err := runImporter([]string{"-h"}, io.Discard, &stderr); err != nil {
		t.Fatalf("runImporter(-h) error = %v", err)
	}
	for _, want := range []string{"Usage of importanimals", "-variant", "-check"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("help output = %q, want %q", stderr.String(), want)
		}
	}
}

func TestSelectImportVariantsRejectsUnknownAndNonSeedIDs(t *testing.T) {
	t.Run("unknown", func(t *testing.T) {
		_, err := selectImportVariants("definitely_not_a_catalog_variant")
		if err == nil {
			t.Fatal("selectImportVariants() succeeded for unknown variant")
		}
		if !strings.Contains(err.Error(), "unknown variant") || !strings.Contains(err.Error(), "SeedStage") {
			t.Fatalf("selectImportVariants() error = %v, want actionable unknown-variant error", err)
		}
	})

	t.Run("non-seed", func(t *testing.T) {
		var nonSeed catalog.Variant
		for _, variant := range catalog.Variants {
			if !variant.SeedStage {
				nonSeed = variant
				break
			}
		}
		if nonSeed.ID == "" {
			t.Fatal("catalog has no non-seed variant for rejection test")
		}

		_, err := selectImportVariants(nonSeed.ID)
		if err == nil {
			t.Fatalf("selectImportVariants() succeeded for non-seed variant %q", nonSeed.ID)
		}
		if !strings.Contains(err.Error(), "not importable") || !strings.Contains(err.Error(), "not marked SeedStage") {
			t.Fatalf("selectImportVariants() error = %v, want actionable non-importable error", err)
		}
	})
}

func TestResolveOutputPathsKeepsAggregateDefaultsAndProtectsTargetedArtifacts(t *testing.T) {
	aggregate, err := parseImportOptions(nil, io.Discard)
	if err != nil {
		t.Fatalf("parse aggregate options: %v", err)
	}
	aggregate = resolveOutputPaths(aggregate)
	if aggregate.reportPath != filepath.FromSlash(defaultReportPath) {
		t.Fatalf("aggregate report path = %q, want %q", aggregate.reportPath, filepath.FromSlash(defaultReportPath))
	}
	if aggregate.previewPath != filepath.FromSlash(defaultPreviewPath) {
		t.Fatalf("aggregate preview path = %q, want %q", aggregate.previewPath, filepath.FromSlash(defaultPreviewPath))
	}

	targeted, err := parseImportOptions([]string{"-variant", "target_variant"}, io.Discard)
	if err != nil {
		t.Fatalf("parse targeted options: %v", err)
	}
	targeted = resolveOutputPaths(targeted)
	if targeted.reportPath != "" || targeted.previewPath != "" {
		t.Fatalf("targeted default artifacts = report:%q preview:%q, want both suppressed", targeted.reportPath, targeted.previewPath)
	}

	explicit, err := parseImportOptions([]string{
		"-variant", "target_variant",
		"-report", filepath.FromSlash("tmp/target-report.json"),
		"-preview", filepath.FromSlash("tmp/target-preview.png"),
	}, io.Discard)
	if err != nil {
		t.Fatalf("parse targeted explicit options: %v", err)
	}
	explicit = resolveOutputPaths(explicit)
	if explicit.reportPath != filepath.FromSlash("tmp/target-report.json") {
		t.Fatalf("explicit report path = %q", explicit.reportPath)
	}
	if explicit.previewPath != filepath.FromSlash("tmp/target-preview.png") {
		t.Fatalf("explicit preview path = %q", explicit.previewPath)
	}
}

func TestCheckModeDoesNotWriteConfiguredRepositoryPaths(t *testing.T) {
	repositoryRoot := filepath.Join(t.TempDir(), "repository")
	if err := os.MkdirAll(repositoryRoot, 0o755); err != nil {
		t.Fatalf("create repository root: %v", err)
	}
	options := importOptions{
		outDir:             filepath.Join(repositoryRoot, "assets", "sprites"),
		reportPath:         filepath.Join(repositoryRoot, "assets", "source", "animals", "seed-import-report.json"),
		previewPath:        filepath.Join(repositoryRoot, "docs", "assets", "animalsdesktop-seed-preview.png"),
		generatedSourceDir: filepath.Join(repositoryRoot, "assets", "source", "animals", "generated"),
		variantID:          "test_seed_variant",
		check:              true,
		reportExplicit:     true,
		previewExplicit:    true,
	}
	variant := catalog.Variant{
		ID:             "test_seed_variant",
		SpeciesID:      "mouse",
		BreedOrMorph:   "Test mouse",
		Color:          "brown",
		PopularityTier: 1,
		MotionProfile:  catalog.MotionProfileSmallRodentScurry,
		SourceStatus:   catalog.SourceStatusPrototypeOnly,
		SpriteBase:     "test_seed_variant",
		SeedStage:      true,
		Shape:          "small_rodent",
		TintHex:        "8a6748",
		AccentHex:      "eadbc0",
	}

	if err := executeImport(options, []catalog.Variant{variant}, io.Discard); err != nil {
		t.Fatalf("executeImport() check error = %v", err)
	}
	entries, err := os.ReadDir(repositoryRoot)
	if err != nil {
		t.Fatalf("read repository root: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("check mode wrote repository entries: %v", entries)
	}
}

func TestNormalizeSourceUsesFixedCanvas(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 200, 120))
	for y := 35; y < 95; y++ {
		for x := 20; x < 180; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 150, G: 110, B: 70, A: 255})
		}
	}
	content := alphaBounds(src)
	if content.Empty() {
		t.Fatalf("alphaBounds returned empty content")
	}
	got := normalizeSource(src, content, profileFor(catalog.MotionProfileSmallRodentScurry))
	if got.Bounds().Dx() != frameW || got.Bounds().Dy() != frameH {
		t.Fatalf("normalized bounds = %v, want %dx%d", got.Bounds(), frameW, frameH)
	}
	if alphaBounds(got).Empty() {
		t.Fatalf("normalized source has no visible pixels")
	}
}

func TestSeedFrameKeepsSpriteInCanvas(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	for y := 30; y < 52; y++ {
		for x := 20; x < 76; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 110, G: 110, B: 110, A: 255})
		}
	}
	for frame := 0; frame < totalFrames; frame++ {
		got := seedFrame(src, frame, 3, catalog.MotionProfileGeckoCrawl)
		if got.Bounds().Dx() != frameW || got.Bounds().Dy() != frameH {
			t.Fatalf("frame %d bounds = %v", frame, got.Bounds())
		}
		if alphaBounds(got).Empty() {
			t.Fatalf("frame %d has no visible pixels", frame)
		}
	}
}

func TestImportVariantUsesMotionSourceSheet(t *testing.T) {
	root := t.TempDir()
	sourcePath := filepath.Join(root, "source.png")
	src := image.NewRGBA(image.Rect(0, 0, 120, 80))
	for y := 20; y < 60; y++ {
		for x := 20; x < 100; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 140, G: 140, B: 135, A: 255})
		}
	}
	if err := writePNG(sourcePath, src); err != nil {
		t.Fatalf("write source: %v", err)
	}

	motionPath := filepath.Join(root, "motion-source.png")
	motion := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for frame := 0; frame < totalFrames; frame++ {
		xOffset := frame * frameW
		motion.SetRGBA(xOffset+8, 8, color.RGBA{R: byte(frame), G: 40, B: 90, A: 255})
		for y := 34; y < 54; y++ {
			for x := 24; x < 72; x++ {
				motion.SetRGBA(xOffset+x, y, color.RGBA{R: 160, G: 160, B: 150, A: 255})
			}
		}
	}
	if err := writePNG(motionPath, motion); err != nil {
		t.Fatalf("write motion source: %v", err)
	}

	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		BreedOrMorph:     "Chinchilla",
		Color:            "standard gray",
		PopularityTier:   1,
		MotionProfile:    catalog.MotionProfileSmallRodentScurry,
		SourceStatus:     catalog.SourceStatusMotionDraft,
		SpriteBase:       "test_chinchilla",
		SeedStage:        true,
		SourcePath:       sourcePath,
		MotionSourcePath: motionPath,
	}
	outDir := filepath.Join(root, "sprites")
	report, err := importVariant(variant, outDir, filepath.Join(root, "generated"))
	if err != nil {
		t.Fatalf("importVariant() error = %v", err)
	}
	if report.MotionFrames != totalFrames || report.SourceSets != 1 || report.RuntimeSets != motionSets || report.MotionSource == "" {
		t.Fatalf("motion report = frames:%d source sets:%d runtime sets:%d source:%q", report.MotionFrames, report.SourceSets, report.RuntimeSets, report.MotionSource)
	}
	if len(report.Outputs) != motionSets {
		t.Fatalf("outputs = %d, want %d", len(report.Outputs), motionSets)
	}
	if len(report.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none for canonical set00 runtime expansion", report.Warnings)
	}
	sheet, err := openPNG(filepath.Join(outDir, "test_chinchilla_set00.png"))
	if err != nil {
		t.Fatalf("open sheet: %v", err)
	}
	got := sheet.RGBAAt(12*frameW+8, 8)
	if got.R != 12 || got.G != 40 || got.B != 90 || got.A != 255 {
		t.Fatalf("frame 12 marker = %#v, want imported motion source marker", got)
	}
}

func TestImportVariantUsesTenMotionSourceSheets(t *testing.T) {
	root := t.TempDir()
	sourcePath := filepath.Join(root, "source.png")
	src := image.NewRGBA(image.Rect(0, 0, 120, 80))
	for y := 20; y < 60; y++ {
		for x := 20; x < 100; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 140, G: 140, B: 135, A: 255})
		}
	}
	if err := writePNG(sourcePath, src); err != nil {
		t.Fatalf("write source: %v", err)
	}

	motionSet00Path := filepath.Join(root, "chinchilla-source-set00-draft.png")
	for set := 0; set < motionSets; set++ {
		motion := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
		for frame := 0; frame < totalFrames; frame++ {
			xOffset := frame * frameW
			motion.SetRGBA(xOffset+8, 8, color.RGBA{R: byte(set), G: byte(frame), B: 90, A: 255})
			for y := 34; y < 54; y++ {
				for x := 24; x < 72; x++ {
					motion.SetRGBA(xOffset+x, y, color.RGBA{R: 160, G: 160, B: 150, A: 255})
				}
			}
		}
		path := filepath.Join(root, "chinchilla-source-set"+twoDigits(set)+"-draft.png")
		if err := writePNG(path, motion); err != nil {
			t.Fatalf("write motion source set %02d: %v", set, err)
		}
	}

	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		BreedOrMorph:     "Chinchilla",
		Color:            "standard gray",
		PopularityTier:   1,
		MotionProfile:    catalog.MotionProfileSmallRodentScurry,
		SourceStatus:     catalog.SourceStatusMotionDraft,
		SpriteBase:       "test_chinchilla",
		SeedStage:        true,
		SourcePath:       sourcePath,
		MotionSourcePath: motionSet00Path,
	}
	outDir := filepath.Join(root, "sprites")
	report, err := importVariant(variant, outDir, filepath.Join(root, "generated"))
	if err != nil {
		t.Fatalf("importVariant() error = %v", err)
	}
	if report.MotionFrames != totalFrames || report.SourceSets != motionSets || report.RuntimeSets != motionSets || report.MotionSource == "" {
		t.Fatalf("motion report = frames:%d source sets:%d runtime sets:%d source:%q", report.MotionFrames, report.SourceSets, report.RuntimeSets, report.MotionSource)
	}
	if len(report.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none for complete set family", report.Warnings)
	}
	sheet, err := openPNG(filepath.Join(outDir, "test_chinchilla_set09.png"))
	if err != nil {
		t.Fatalf("open sheet: %v", err)
	}
	got := sheet.RGBAAt(12*frameW+8, 8)
	if got.R != 9 || got.G != 12 || got.B != 90 || got.A != 255 {
		t.Fatalf("set09 frame 12 marker = %#v, want imported set-specific motion source marker", got)
	}
}

func TestLoadMotionSourceSheetRejectsOpaqueFrameBackground(t *testing.T) {
	root := t.TempDir()
	motionPath := filepath.Join(root, "motion-source.png")
	motion := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for y := 0; y < frameH; y++ {
		for x := 0; x < frameW*totalFrames; x++ {
			motion.SetRGBA(x, y, color.RGBA{R: 238, G: 238, B: 238, A: 255})
		}
	}
	if err := writePNG(motionPath, motion); err != nil {
		t.Fatalf("write motion source: %v", err)
	}

	_, err := loadMotionSourceSheet(motionPath)
	if err == nil {
		t.Fatalf("loadMotionSourceSheet() succeeded for opaque motion source")
	}
	if !strings.Contains(err.Error(), "transparent background") {
		t.Fatalf("loadMotionSourceSheet() error = %v, want transparent-background failure", err)
	}
}

func TestProceduralSourceHasVisibleContent(t *testing.T) {
	src := proceduralSource(catalog.Variant{
		ID:        "ferret_sable",
		SpeciesID: "ferret",
		Shape:     "ferret",
		TintHex:   "8b6746",
		AccentHex: "ece0c8",
	})
	content := alphaBounds(src)
	if content.Empty() {
		t.Fatalf("procedural source content is empty")
	}
	if content.Dx() < 180 || content.Dy() < 50 {
		t.Fatalf("procedural source content too small: %v", content)
	}
}

func TestTintSourcePreservesAlpha(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 8, 8))
	src.SetRGBA(2, 2, color.RGBA{R: 200, G: 180, B: 120, A: 190})
	tint, err := parseHexColor("336699")
	if err != nil {
		t.Fatalf("parseHexColor() error = %v", err)
	}
	got := tintSource(src, tint, nil)
	if got.RGBAAt(2, 2).A != 190 {
		t.Fatalf("tinted alpha = %d, want 190", got.RGBAAt(2, 2).A)
	}
	if got.RGBAAt(0, 0).A != 0 {
		t.Fatalf("transparent pixel alpha = %d, want 0", got.RGBAAt(0, 0).A)
	}
}

func TestSeedVariantGeneratedAssetsExist(t *testing.T) {
	for _, variant := range catalog.SeedVariants() {
		source := filepath.Join("..", "..", "assets", "source", "animals", "generated", variant.SpriteBase+"-source.png")
		if _, err := os.Stat(source); err != nil {
			t.Fatalf("missing generated source for %s: %v", variant.ID, err)
		}
		for set := 0; set < motionSets; set++ {
			path := filepath.Join("..", "..", "assets", "sprites", variant.SpriteBase+"_set"+twoDigits(set)+".png")
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("missing sprite sheet for %s set %02d: %v", variant.ID, set, err)
			}
			sheet, err := openPNG(path)
			if err != nil {
				t.Fatalf("open sprite sheet for %s set %02d: %v", variant.ID, set, err)
			}
			if got := sheet.Bounds(); got.Dx() != frameW*totalFrames || got.Dy() != frameH {
				t.Fatalf("sheet bounds for %s set %02d = %v", variant.ID, set, got)
			}
			for frame := 0; frame < totalFrames; frame++ {
				frameRect := image.Rect(frame*frameW, 0, (frame+1)*frameW, frameH)
				if alphaBounds(sheet.SubImage(frameRect)).Empty() {
					t.Fatalf("empty frame for %s set %02d frame %02d", variant.ID, set, frame)
				}
			}
		}
	}
}

func TestMotionProfilesHaveDistinctEcologyOffsets(t *testing.T) {
	rabbitDx, rabbitDy := motionOffset(27, 0, catalog.MotionProfileRabbitHop)
	_, snakeDy := motionOffset(27, 0, catalog.MotionProfileSnakeSlither)
	_, tortoiseDy := motionOffset(14, 0, catalog.MotionProfileTortoisePlod)
	dogDx, dogDy := motionOffset(13, 0, catalog.MotionProfileDogTrot)

	if rabbitDy >= 0 || rabbitDx != 0 {
		t.Fatalf("rabbit action offset = (%d,%d), want upward hop without horizontal drift", rabbitDx, rabbitDy)
	}
	if snakeDy != 0 {
		t.Fatalf("snake action dy = %d, want low slither", snakeDy)
	}
	if tortoiseDy != 0 {
		t.Fatalf("tortoise fast dy = %d, want no vertical bob", tortoiseDy)
	}
	if dogDy >= 0 || dogDx == 0 {
		t.Fatalf("dog trot offset = (%d,%d), want horizontal trot with lift", dogDx, dogDy)
	}
}

func twoDigits(v int) string {
	return string([]byte{'0' + byte(v/10), '0' + byte(v%10)})
}
