package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"animals-desktop/internal/catalog"
	"animals-desktop/internal/motionsource"
	xdraw "golang.org/x/image/draw"
)

const (
	totalFrames        = 62
	motionSets         = 10
	frameW             = 96
	frameH             = 64
	defaultReportPath  = "assets/source/animals/seed-import-report.json"
	defaultPreviewPath = "docs/assets/animalsdesktop-seed-preview.png"
)

var defaultGeneratedSourceDir = filepath.FromSlash("assets/source/animals/generated")

type importOptions struct {
	outDir             string
	reportPath         string
	previewPath        string
	generatedSourceDir string
	variantID          string
	check              bool
	reportExplicit     bool
	previewExplicit    bool
}

type seedReport struct {
	Variant         string     `json:"variant"`
	Species         string     `json:"species"`
	BreedOrMorph    string     `json:"breed_or_morph"`
	Color           string     `json:"color"`
	PopularityTier  int        `json:"popularity_tier"`
	MotionProfile   string     `json:"motion_profile"`
	SourceStatus    string     `json:"source_status"`
	Source          string     `json:"source"`
	GeneratedSource string     `json:"generated_source,omitempty"`
	MotionSource    string     `json:"motion_source,omitempty"`
	MotionFrames    int        `json:"motion_frames,omitempty"`
	SourceSets      int        `json:"source_sets,omitempty"`
	RuntimeSets     int        `json:"runtime_sets,omitempty"`
	SpriteBase      string     `json:"sprite_base"`
	Shape           string     `json:"shape,omitempty"`
	TintHex         string     `json:"tint_hex,omitempty"`
	AccentHex       string     `json:"accent_hex,omitempty"`
	SourceWidth     int        `json:"source_width"`
	SourceHeight    int        `json:"source_height"`
	Content         rectReport `json:"content"`
	Outputs         []string   `json:"outputs"`
	Warnings        []string   `json:"warnings,omitempty"`
}

type rectReport struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type renderProfile struct {
	targetW  int
	targetH  int
	baseline int
	low      bool
}

func main() {
	if err := runImporter(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}

func runImporter(args []string, stdout io.Writer, stderr io.Writer) error {
	options, err := parseImportOptions(args, stderr)
	if err == flag.ErrHelp {
		return nil
	}
	if err != nil {
		return err
	}
	if err := catalog.Validate(); err != nil {
		return fmt.Errorf("validate catalog: %w", err)
	}
	variants, err := selectImportVariants(options.variantID)
	if err != nil {
		return err
	}
	options = resolveOutputPaths(options)
	return executeImport(options, variants, stdout)
}

func parseImportOptions(args []string, stderr io.Writer) (importOptions, error) {
	options := importOptions{}
	flags := flag.NewFlagSet("importanimals", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.StringVar(&options.outDir, "out", filepath.FromSlash("assets/sprites"), "output sprite directory")
	flags.StringVar(&options.reportPath, "report", filepath.FromSlash(defaultReportPath), "JSON report path")
	flags.StringVar(&options.previewPath, "preview", filepath.FromSlash(defaultPreviewPath), "seed preview PNG path")
	flags.StringVar(&options.generatedSourceDir, "generated-source-dir", defaultGeneratedSourceDir, "directory for normalized/generated source PNGs")
	flags.StringVar(&options.variantID, "variant", "", "import exactly one catalog SeedStage variant")
	flags.BoolVar(&options.check, "check", false, "validate imports using temporary output without repository writes")
	if err := flags.Parse(args); err != nil {
		return importOptions{}, err
	}
	if flags.NArg() != 0 {
		return importOptions{}, fmt.Errorf("unexpected positional arguments: %s", strings.Join(flags.Args(), " "))
	}
	flags.Visit(func(flag *flag.Flag) {
		switch flag.Name {
		case "report":
			options.reportExplicit = true
		case "preview":
			options.previewExplicit = true
		}
	})
	return options, nil
}

func selectImportVariants(variantID string) ([]catalog.Variant, error) {
	if variantID == "" {
		variants := catalog.SeedVariants()
		if len(variants) == 0 {
			return nil, fmt.Errorf("catalog has no variants marked SeedStage")
		}
		return variants, nil
	}
	variant, ok := catalog.VariantByID(variantID)
	if !ok {
		return nil, fmt.Errorf("unknown variant %q; choose a catalog variant ID marked SeedStage", variantID)
	}
	if !variant.SeedStage {
		return nil, fmt.Errorf("variant %q is not importable: catalog entry is not marked SeedStage", variantID)
	}
	return []catalog.Variant{variant}, nil
}

func resolveOutputPaths(options importOptions) importOptions {
	if options.variantID == "" {
		return options
	}
	if !options.reportExplicit {
		options.reportPath = ""
	}
	if !options.previewExplicit {
		options.previewPath = ""
	}
	return options
}

func executeImport(options importOptions, variants []catalog.Variant, stdout io.Writer) error {
	if options.check {
		tempRoot, err := os.MkdirTemp("", "animalsdesktop-import-check-*")
		if err != nil {
			return fmt.Errorf("create check output: %w", err)
		}
		defer os.RemoveAll(tempRoot)
		options.outDir = filepath.Join(tempRoot, "sprites")
		options.generatedSourceDir = filepath.Join(tempRoot, "generated")
		if options.reportPath != "" {
			options.reportPath = filepath.Join(tempRoot, "seed-import-report.json")
		}
		if options.previewPath != "" {
			options.previewPath = filepath.Join(tempRoot, "seed-preview.png")
		}
	}

	reports := make([]seedReport, 0, len(variants))
	for _, variant := range variants {
		report, err := importVariant(variant, options.outDir, options.generatedSourceDir)
		if err != nil {
			return fmt.Errorf("import %s: %w", variant.ID, err)
		}
		reports = append(reports, report)
	}
	if options.reportPath != "" {
		if err := writeJSON(options.reportPath, reports); err != nil {
			return fmt.Errorf("write report: %w", err)
		}
	}
	if options.previewPath != "" {
		if err := writeSeedPreview(variants, options.outDir, options.previewPath); err != nil {
			return fmt.Errorf("write preview: %w", err)
		}
	}
	if stdout == nil {
		stdout = io.Discard
	}
	if options.check {
		fmt.Fprintf(stdout, "validated %d seed animal variants through %d-frame import pipeline\n", len(reports), totalFrames)
	} else {
		fmt.Fprintf(stdout, "imported %d seed animal variants into %d-frame sheets\n", len(reports), totalFrames)
	}
	return nil
}

func importVariant(variant catalog.Variant, outDir string, generatedSourceDir string) (seedReport, error) {
	src, sourceLabel, generatedSource, warnings, err := prepareVariantSource(variant, generatedSourceDir)
	if err != nil {
		return seedReport{}, err
	}
	bounds := src.Bounds()
	content := alphaBounds(src)
	report := seedReport{
		Variant:         variant.ID,
		Species:         variant.SpeciesID,
		BreedOrMorph:    variant.BreedOrMorph,
		Color:           variant.Color,
		PopularityTier:  variant.PopularityTier,
		MotionProfile:   catalog.MotionProfileForVariant(variant),
		SourceStatus:    variant.SourceStatus,
		Source:          sourceLabel,
		GeneratedSource: generatedSource,
		SpriteBase:      variant.SpriteBase,
		Shape:           variant.Shape,
		TintHex:         variant.TintHex,
		AccentHex:       variant.AccentHex,
		SourceWidth:     bounds.Dx(),
		SourceHeight:    bounds.Dy(),
		Content: rectReport{
			X: content.Min.X,
			Y: content.Min.Y,
			W: content.Dx(),
			H: content.Dy(),
		},
		Outputs:  make([]string, 0, motionSets),
		Warnings: warnings,
	}
	if content.Empty() {
		return seedReport{}, fmt.Errorf("source has no visible alpha: %s", sourceLabel)
	}
	if content == bounds {
		report.Warnings = append(report.Warnings, "source is fully opaque; verify transparent background before final animation import")
	}

	motionProfile := catalog.MotionProfileForVariant(variant)
	if variant.MotionSourcePath != "" {
		outputs, motionFrames, sourceSets, err := importMotionSourceSheet(variant, outDir)
		if err != nil {
			return seedReport{}, err
		}
		report.MotionSource = filepath.ToSlash(variant.MotionSourcePath)
		report.MotionFrames = motionFrames
		report.SourceSets = sourceSets
		report.RuntimeSets = len(outputs)
		report.Outputs = outputs
		return report, nil
	}

	normalized := normalizeSource(src, content, profileFor(motionProfile))
	for set := 0; set < motionSets; set++ {
		sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
		for frame := 0; frame < totalFrames; frame++ {
			sprite := seedFrame(normalized, frame, set, motionProfile)
			dst := image.Rect(frame*frameW, 0, (frame+1)*frameW, frameH)
			draw.Draw(sheet, dst, sprite, image.Point{}, draw.Over)
		}
		outPath := filepath.Join(outDir, fmt.Sprintf("%s_set%02d.png", variant.SpriteBase, set))
		if err := writePNG(outPath, sheet); err != nil {
			return seedReport{}, err
		}
		report.Outputs = append(report.Outputs, filepath.ToSlash(outPath))
	}
	return report, nil
}

func importMotionSourceSheet(variant catalog.Variant, outDir string) ([]string, int, int, error) {
	sourcePaths, err := motionSourceSheetPaths(variant.MotionSourcePath)
	if err != nil {
		return nil, 0, 0, err
	}
	outputs := make([]string, 0, motionSets)
	for set := 0; set < motionSets; set++ {
		sourcePath := sourcePaths[0]
		if len(sourcePaths) == motionSets {
			sourcePath = sourcePaths[set]
		}
		if _, err := loadMotionSourceSheet(sourcePath); err != nil {
			return nil, 0, 0, err
		}
		outPath := filepath.Join(outDir, fmt.Sprintf("%s_set%02d.png", variant.SpriteBase, set))
		if err := copyFile(sourcePath, outPath); err != nil {
			return nil, 0, 0, err
		}
		outputs = append(outputs, filepath.ToSlash(outPath))
	}
	return outputs, totalFrames, len(sourcePaths), nil
}

func motionSourceSheetPaths(set00Path string) ([]string, error) {
	return motionsource.ResolveSetPaths(set00Path, motionSets)
}

func loadMotionSourceSheet(path string) (*image.RGBA, error) {
	sheet, err := openPNG(path)
	if err != nil {
		return nil, err
	}
	bounds := sheet.Bounds()
	wantW := frameW * totalFrames
	if bounds.Dx() != wantW || bounds.Dy() != frameH {
		return nil, fmt.Errorf("motion source sheet %s bounds = %dx%d, want %dx%d", path, bounds.Dx(), bounds.Dy(), wantW, frameH)
	}
	for frame := 0; frame < totalFrames; frame++ {
		frameRect := image.Rect(bounds.Min.X+frame*frameW, bounds.Min.Y, bounds.Min.X+(frame+1)*frameW, bounds.Min.Y+frameH)
		content := alphaBounds(sheet.SubImage(frameRect))
		if content.Empty() {
			return nil, fmt.Errorf("motion source sheet %s has empty frame %02d", path, frame)
		}
		if content == frameRect {
			return nil, fmt.Errorf("motion source sheet %s frame %02d has no transparent background", path, frame)
		}
	}
	return sheet, nil
}

func prepareVariantSource(variant catalog.Variant, generatedSourceDir string) (*image.RGBA, string, string, []string, error) {
	warnings := []string{}
	var src *image.RGBA
	var err error
	sourceLabel := filepath.ToSlash(variant.SourcePath)
	switch {
	case variant.SourcePath != "":
		src, err = openPNG(variant.SourcePath)
		if err != nil {
			return nil, "", "", nil, err
		}
	case variant.Shape != "":
		src = proceduralSource(variant)
		sourceLabel = "procedural:" + variant.Shape
		warnings = append(warnings, "procedural source seed; replace with ImageGen source-truth before final animation import")
	default:
		return nil, "", "", nil, fmt.Errorf("variant has neither SourcePath nor Shape")
	}

	if variant.TintHex != "" {
		tint, err := parseHexColor(variant.TintHex)
		if err != nil {
			return nil, "", "", nil, err
		}
		var accent *color.RGBA
		if variant.AccentHex != "" {
			accentColor, err := parseHexColor(variant.AccentHex)
			if err != nil {
				return nil, "", "", nil, err
			}
			accent = &accentColor
		}
		src = tintSource(src, tint, accent)
	}

	generatedPath := filepath.Join(generatedSourceDir, variant.SpriteBase+"-source.png")
	if err := writePNG(generatedPath, src); err != nil {
		return nil, "", "", nil, err
	}
	return src, sourceLabel, filepath.ToSlash(generatedPath), warnings, nil
}

func tintSource(src *image.RGBA, tint color.RGBA, accent *color.RGBA) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	b := src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := src.RGBAAt(x, y)
			if c.A == 0 {
				continue
			}
			target := tint
			if accent != nil && accentMask(x-b.Min.X, y-b.Min.Y, b.Dx(), b.Dy()) {
				target = *accent
			}
			lum := int(c.R)*30 + int(c.G)*59 + int(c.B)*11
			lum /= 100
			shade := 0.52 + float64(lum)/255.0*0.66
			out.SetRGBA(x, y, color.RGBA{
				R: clampByte(float64(target.R) * shade),
				G: clampByte(float64(target.G) * shade),
				B: clampByte(float64(target.B) * shade),
				A: c.A,
			})
		}
	}
	return out
}

func accentMask(x, y, w, h int) bool {
	if w <= 0 || h <= 0 {
		return false
	}
	leftPatch := x > w/8 && x < w/3 && y > h/4 && y < h*3/4
	backPatch := x > w/2 && x < w*7/8 && y > h/6 && y < h/2
	return leftPatch || backPatch
}

func parseHexColor(hex string) (color.RGBA, error) {
	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("hex color %q must be 6 characters", hex)
	}
	val := func(i int) (byte, error) {
		var out byte
		for _, ch := range hex[i : i+2] {
			out <<= 4
			switch {
			case ch >= '0' && ch <= '9':
				out += byte(ch - '0')
			case ch >= 'a' && ch <= 'f':
				out += byte(ch-'a') + 10
			case ch >= 'A' && ch <= 'F':
				out += byte(ch-'A') + 10
			default:
				return 0, fmt.Errorf("invalid hex color %q", hex)
			}
		}
		return out, nil
	}
	r, err := val(0)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := val(2)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := val(4)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

func clampByte(v float64) byte {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return byte(v + 0.5)
}

func proceduralSource(variant catalog.Variant) *image.RGBA {
	base, _ := parseHexColor(defaultHex(variant.TintHex, "8a6748"))
	accent, _ := parseHexColor(defaultHex(variant.AccentHex, "eadbc0"))
	dark := darken(base, 0.45)
	img := image.NewRGBA(image.Rect(0, 0, 512, 384))
	switch variant.Shape {
	case "ferret":
		ellipse(img, 138, 198, 146, 42, base)
		ellipse(img, 276, 181, 50, 34, base)
		ellipse(img, 318, 174, 26, 24, accent)
		ellipse(img, 70, 205, 46, 18, dark)
		legs(img, []int{172, 248}, 230, dark)
		eye(img, 318, 166)
	case "guinea_pig":
		ellipse(img, 218, 210, 118, 64, base)
		ellipse(img, 310, 196, 48, 46, base)
		ellipse(img, 180, 190, 42, 30, accent)
		ellipse(img, 262, 225, 54, 34, accent)
		ellipse(img, 305, 158, 15, 18, base)
		legs(img, []int{174, 258, 316}, 258, dark)
		eye(img, 326, 187)
	case "hedgehog":
		ellipse(img, 226, 216, 112, 56, base)
		ellipse(img, 314, 202, 44, 34, accent)
		for i := 0; i < 18; i++ {
			x := 128 + i*10
			triangle(img, image.Pt(x, 162+(i%3)*4), image.Pt(x+16, 190), image.Pt(x-8, 190), dark)
		}
		legs(img, []int{182, 262}, 258, dark)
		eye(img, 327, 194)
	case "squirrel":
		ellipse(img, 220, 216, 92, 46, base)
		ellipse(img, 298, 194, 42, 36, base)
		ellipse(img, 108, 168, 48, 92, accent)
		ellipse(img, 124, 146, 34, 62, base)
		ellipse(img, 296, 156, 13, 18, base)
		legs(img, []int{188, 250}, 250, dark)
		eye(img, 314, 188)
	case "fox":
		ellipse(img, 210, 214, 110, 44, base)
		ellipse(img, 306, 190, 48, 34, base)
		triangle(img, image.Pt(335, 188), image.Pt(382, 199), image.Pt(337, 214), base)
		triangle(img, image.Pt(294, 154), image.Pt(306, 124), image.Pt(318, 158), base)
		ellipse(img, 92, 198, 70, 32, base)
		ellipse(img, 55, 198, 25, 18, accent)
		ellipse(img, 314, 205, 24, 16, accent)
		legs(img, []int{176, 248}, 250, dark)
		eye(img, 322, 183)
	case "red_panda":
		ellipse(img, 206, 216, 106, 48, base)
		ellipse(img, 302, 190, 50, 42, base)
		ellipse(img, 95, 198, 78, 28, base)
		for i := 0; i < 4; i++ {
			ellipse(img, 52+i*26, 198, 9, 24, accent)
		}
		ellipse(img, 302, 194, 24, 18, accent)
		ellipse(img, 330, 196, 14, 16, accent)
		legs(img, []int{174, 252}, 254, dark)
		eye(img, 318, 184)
	case "otter":
		ellipse(img, 208, 220, 132, 38, base)
		ellipse(img, 318, 203, 46, 32, base)
		ellipse(img, 89, 232, 72, 15, dark)
		ellipse(img, 332, 214, 28, 14, accent)
		legs(img, []int{188, 270}, 250, dark)
		eye(img, 330, 194)
	case "sugar_glider":
		ellipse(img, 220, 210, 72, 40, base)
		ellipse(img, 286, 190, 36, 34, base)
		triangle(img, image.Pt(164, 200), image.Pt(234, 220), image.Pt(178, 242), accent)
		triangle(img, image.Pt(250, 216), image.Pt(326, 196), image.Pt(314, 238), accent)
		ellipse(img, 112, 218, 76, 10, dark)
		ellipse(img, 282, 160, 12, 18, base)
		eye(img, 298, 184)
	case "capybara":
		ellipse(img, 214, 222, 132, 56, base)
		ellipse(img, 322, 202, 52, 38, base)
		ellipse(img, 328, 216, 32, 16, accent)
		legs(img, []int{168, 232, 300}, 262, dark)
		eye(img, 334, 194)
	case "tortoise":
		ellipse(img, 220, 218, 112, 52, dark)
		ellipse(img, 220, 206, 92, 42, base)
		ellipse(img, 326, 218, 32, 22, accent)
		ellipse(img, 118, 228, 24, 14, dark)
		legs(img, []int{160, 224, 286}, 260, dark)
		eye(img, 336, 214)
	case "small_rodent":
		ellipse(img, 220, 216, 82, 44, base)
		ellipse(img, 288, 196, 36, 32, base)
		ellipse(img, 112, 223, 80, 9, dark)
		ellipse(img, 285, 166, 13, 18, base)
		ellipse(img, 202, 206, 30, 22, accent)
		legs(img, []int{188, 248}, 252, dark)
		eye(img, 301, 189)
	case "prairie_dog":
		ellipse(img, 226, 224, 70, 50, base)
		ellipse(img, 244, 180, 42, 36, base)
		ellipse(img, 238, 158, 12, 16, base)
		ellipse(img, 262, 158, 12, 16, base)
		ellipse(img, 256, 198, 26, 16, accent)
		legs(img, []int{202, 250}, 262, dark)
		eye(img, 258, 174)
	case "chipmunk":
		ellipse(img, 218, 216, 86, 42, base)
		ellipse(img, 292, 194, 38, 32, base)
		ellipse(img, 112, 170, 42, 82, accent)
		ellipse(img, 126, 148, 30, 56, base)
		for i := 0; i < 3; i++ {
			ellipse(img, 204+i*17, 204, 4, 34, dark)
		}
		ellipse(img, 290, 164, 12, 17, base)
		legs(img, []int{188, 246}, 250, dark)
		eye(img, 308, 187)
	case "dragon":
		ellipse(img, 208, 222, 118, 34, base)
		ellipse(img, 318, 205, 46, 26, base)
		triangle(img, image.Pt(284, 188), image.Pt(300, 166), image.Pt(310, 196), accent)
		triangle(img, image.Pt(246, 190), image.Pt(258, 171), image.Pt(270, 199), accent)
		ellipse(img, 88, 232, 70, 11, dark)
		legs(img, []int{174, 262}, 248, dark)
		eye(img, 330, 198)
	case "snake":
		for i := 0; i < 9; i++ {
			x := 104 + i*28
			y := 224 + []int{0, -9, -13, -8, 2, 10, 12, 5, -2}[i]
			ellipse(img, x, y, 32, 16, base)
		}
		ellipse(img, 340, 214, 36, 20, base)
		triangle(img, image.Pt(370, 214), image.Pt(400, 205), image.Pt(395, 225), base)
		for i := 1; i < 8; i += 2 {
			ellipse(img, 104+i*28, 218, 4, 16, accent)
		}
		eye(img, 354, 207)
	case "frog":
		ellipse(img, 224, 224, 84, 42, base)
		ellipse(img, 292, 196, 50, 34, base)
		ellipse(img, 270, 166, 14, 15, accent)
		ellipse(img, 306, 166, 14, 15, accent)
		ellipse(img, 166, 246, 42, 12, dark)
		ellipse(img, 286, 246, 46, 13, dark)
		ellipse(img, 300, 210, 26, 12, accent)
		eye(img, 305, 164)
	default:
		ellipse(img, 220, 216, 110, 50, base)
		ellipse(img, 310, 198, 44, 36, base)
		legs(img, []int{180, 260}, 254, dark)
		eye(img, 326, 190)
	}
	return img
}

func defaultHex(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func darken(c color.RGBA, amount float64) color.RGBA {
	return color.RGBA{
		R: clampByte(float64(c.R) * amount),
		G: clampByte(float64(c.G) * amount),
		B: clampByte(float64(c.B) * amount),
		A: c.A,
	}
}

func legs(img *image.RGBA, xs []int, y int, c color.RGBA) {
	for _, x := range xs {
		ellipse(img, x, y, 9, 25, c)
		ellipse(img, x+20, y, 9, 25, c)
	}
}

func eye(img *image.RGBA, x, y int) {
	ellipse(img, x, y, 5, 5, color.RGBA{R: 18, G: 15, B: 12, A: 255})
	ellipse(img, x+2, y-2, 1, 1, color.RGBA{R: 240, G: 240, B: 230, A: 255})
}

func ellipse(img *image.RGBA, cx, cy, rx, ry int, c color.RGBA) {
	if rx <= 0 || ry <= 0 {
		return
	}
	for y := cy - ry; y <= cy+ry; y++ {
		for x := cx - rx; x <= cx+rx; x++ {
			if !image.Pt(x, y).In(img.Bounds()) {
				continue
			}
			dx := float64(x-cx) / float64(rx)
			dy := float64(y-cy) / float64(ry)
			if dx*dx+dy*dy <= 1 {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func triangle(img *image.RGBA, a, b, c image.Point, fill color.RGBA) {
	minX := minInt(a.X, minInt(b.X, c.X))
	maxX := maxInt(a.X, maxInt(b.X, c.X))
	minY := minInt(a.Y, minInt(b.Y, c.Y))
	maxY := maxInt(a.Y, maxInt(b.Y, c.Y))
	area := edge(a, b, c)
	if area == 0 {
		return
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := image.Pt(x, y)
			if !p.In(img.Bounds()) {
				continue
			}
			w0 := edge(b, c, p)
			w1 := edge(c, a, p)
			w2 := edge(a, b, p)
			if (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0) {
				img.SetRGBA(x, y, fill)
			}
		}
	}
}

func edge(a, b, c image.Point) int {
	return (c.X-a.X)*(b.Y-a.Y) - (c.Y-a.Y)*(b.X-a.X)
}

func profileFor(profileOrSpecies string) renderProfile {
	profile := profileOrSpecies
	if _, ok := catalog.SpeciesByID(profileOrSpecies); ok {
		profile = catalog.DefaultMotionProfileForSpecies(profileOrSpecies)
	}
	switch profile {
	case catalog.MotionProfileGeckoCrawl, catalog.MotionProfileOtterSlide, catalog.MotionProfileTortoisePlod, catalog.MotionProfileSnakeSlither, catalog.MotionProfileDragonPlod:
		return renderProfile{targetW: 90, targetH: 30, baseline: 59, low: true}
	case catalog.MotionProfileRabbitHop, catalog.MotionProfileFrogHop, catalog.MotionProfileBirdHop:
		return renderProfile{targetW: 82, targetH: 56, baseline: 60}
	case catalog.MotionProfileDogTrot, catalog.MotionProfileCatStalk, catalog.MotionProfileFoxTrot, catalog.MotionProfileRedPandaAmble:
		return renderProfile{targetW: 84, targetH: 54, baseline: 59}
	case catalog.MotionProfileFerretSlink:
		return renderProfile{targetW: 90, targetH: 42, baseline: 59}
	case catalog.MotionProfileSmallRodentScurry, catalog.MotionProfileGuineaPigWaddle, catalog.MotionProfileHedgehogShuffle, catalog.MotionProfileSugarGliderSkitter:
		return renderProfile{targetW: 76, targetH: 48, baseline: 59}
	case catalog.MotionProfileSquirrelBound:
		return renderProfile{targetW: 86, targetH: 58, baseline: 60}
	case catalog.MotionProfileCapybaraLumber:
		return renderProfile{targetW: 88, targetH: 46, baseline: 60}
	default:
		return renderProfile{targetW: 84, targetH: 52, baseline: 59}
	}
}

func normalizeSource(src image.Image, content image.Rectangle, profile renderProfile) *image.RGBA {
	content = content.Intersect(src.Bounds())
	if content.Empty() {
		content = src.Bounds()
	}
	scale := minFloat(
		float64(profile.targetW)/float64(maxInt(1, content.Dx())),
		float64(profile.targetH)/float64(maxInt(1, content.Dy())),
	)
	w := maxInt(1, int(float64(content.Dx())*scale+0.5))
	h := maxInt(1, int(float64(content.Dy())*scale+0.5))
	x := (frameW - w) / 2
	y := profile.baseline - h
	if profile.low {
		y = minInt(frameH-h-1, maxInt(0, y+4))
	} else {
		y = maxInt(0, minInt(frameH-h-1, y))
	}

	scaled := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, content, draw.Over, nil)
	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	draw.Draw(out, image.Rect(x, y, x+w, y+h), withOutline(scaled), image.Point{}, draw.Over)
	return out
}

func withOutline(src *image.RGBA) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	shadow := color.RGBA{R: 34, G: 30, B: 25, A: 78}
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			if src.RGBAAt(x, y).A < 32 {
				continue
			}
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					px, py := x+dx, y+dy
					if !image.Pt(px, py).In(out.Bounds()) || out.RGBAAt(px, py).A >= shadow.A {
						continue
					}
					out.SetRGBA(px, py, shadow)
				}
			}
		}
	}
	draw.Draw(out, out.Bounds(), src, src.Bounds().Min, draw.Over)
	return out
}

func seedFrame(src *image.RGBA, frame int, set int, motionProfile string) *image.RGBA {
	dx, dy := motionOffset(frame, set, motionProfile)
	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	dst := src.Bounds().Add(image.Pt(dx, dy)).Intersect(out.Bounds())
	if dst.Empty() {
		return out
	}
	draw.Draw(out, dst, src, image.Pt(dst.Min.X-dx, dst.Min.Y-dy), draw.Over)
	return out
}

func motionOffset(frame int, set int, motionProfile string) (int, int) {
	local := frame
	phase := set % 4
	switch {
	case frame < 4:
		return idleOffset(local, phase, motionProfile)
	case frame < 12:
		local = frame - 4
		return walkOffset(local, phase, motionProfile)
	case frame < 20:
		local = frame - 12
		return fastOffset(local, phase, motionProfile)
	case frame < 26:
		local = frame - 20
		return feedOffset(local, phase, motionProfile)
	case frame < 32:
		local = frame - 26
		return actionOffset(local, phase, motionProfile)
	case frame < 40:
		local = frame - 32
		return groomOffset(local, phase, motionProfile)
	case frame < 48:
		local = frame - 40
		return turnOffset(local, phase, motionProfile)
	case frame < 56:
		local = frame - 48
		return restOffset(local, phase, motionProfile)
	default:
		local = frame - 56
		return alertOffset(local, phase, motionProfile)
	}
}

func idleOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileSnakeSlither, catalog.MotionProfileGeckoCrawl, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod, catalog.MotionProfileOtterSlide:
		return 0, 0
	default:
		return 0, []int{0, -1, 0, 0}[(local+phase)%4]
	}
}

func walkOffset(local int, phase int, profile string) (int, int) {
	step := []int{-2, -1, 0, 1, 2, 1, 0, -1}[(local+phase)%8]
	switch profile {
	case catalog.MotionProfileCatStalk:
		return step, []int{0, 0, -1, 0}[(local+phase)%4]
	case catalog.MotionProfileDogTrot, catalog.MotionProfileFoxTrot:
		return step, []int{0, -2, 0, -1}[(local+phase)%4]
	case catalog.MotionProfileRabbitHop, catalog.MotionProfileFrogHop, catalog.MotionProfileSquirrelBound, catalog.MotionProfileBirdHop:
		return step, []int{0, -2, -3, -1}[(local+phase)%4]
	case catalog.MotionProfileSnakeSlither:
		return []int{-3, -1, 1, 3, 2, 0, -2, -3}[(local+phase)%8], 0
	case catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod:
		return []int{-1, 0, 1, 1, 0, -1, -1, 0}[(local+phase)%8], 0
	case catalog.MotionProfileGeckoCrawl, catalog.MotionProfileOtterSlide:
		return step, 0
	case catalog.MotionProfileCapybaraLumber:
		return step / 2, []int{0, -1, 0, 0}[(local+phase)%4]
	default:
		return step, []int{0, -1, 0, 0}[(local+phase)%4]
	}
}

func fastOffset(local int, phase int, profile string) (int, int) {
	step := []int{-3, -1, 1, 3, 2, 0, -2, -3}[(local+phase)%8]
	switch profile {
	case catalog.MotionProfileSnakeSlither:
		return []int{-4, -2, 1, 4, 3, 0, -2, -4}[(local+phase)%8], 0
	case catalog.MotionProfileGeckoCrawl, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod, catalog.MotionProfileOtterSlide:
		return step, 0
	case catalog.MotionProfileRabbitHop, catalog.MotionProfileFrogHop, catalog.MotionProfileSquirrelBound, catalog.MotionProfileBirdHop:
		return step, []int{0, -3, -6, -3}[(local+phase)%4]
	case catalog.MotionProfileDogTrot, catalog.MotionProfileFoxTrot:
		return step, []int{0, -2, 0, -2}[(local+phase)%4]
	case catalog.MotionProfileCatStalk:
		return step / 2, []int{0, 0, -1, 0}[(local+phase)%4]
	case catalog.MotionProfileCapybaraLumber:
		return step / 2, []int{0, -1, 0, -1}[(local+phase)%4]
	default:
		return step, []int{0, -1, 0, -1}[(local+phase)%4]
	}
}

func feedOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileSnakeSlither, catalog.MotionProfileGeckoCrawl, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod:
		return []int{0, 1, 1, 0, -1, 0}[(local+phase)%6], 0
	case catalog.MotionProfileFrogHop:
		return 0, []int{0, 2, 1, 2, 0, 0}[(local+phase)%6]
	default:
		return 0, []int{0, 1, 0, 1, 0, 0}[(local+phase)%6]
	}
}

func actionOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileRabbitHop, catalog.MotionProfileFrogHop, catalog.MotionProfileSquirrelBound, catalog.MotionProfileBirdHop:
		return 0, []int{0, -2, -5, -4, -1, 0}[(local+phase)%6]
	case catalog.MotionProfileSnakeSlither:
		return []int{-2, -1, 1, 2, 1, -1}[(local+phase)%6], 0
	case catalog.MotionProfileGeckoCrawl, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod:
		return []int{-2, -1, 0, 1, 2, 1}[(local+phase)%6], 0
	case catalog.MotionProfileOtterSlide:
		return []int{-3, -2, 0, 2, 3, 1}[(local+phase)%6], []int{0, 1, 0, 1, 0, 0}[(local+phase)%6]
	case catalog.MotionProfileDogTrot, catalog.MotionProfileFoxTrot:
		return []int{-1, 0, 1, 1, 0, -1}[(local+phase)%6], []int{0, -1, -2, -1, 0, 0}[(local+phase)%6]
	case catalog.MotionProfileCatStalk:
		return []int{-1, 0, 1, 0, -1, 0}[(local+phase)%6], []int{0, 0, -1, 0, 0, 0}[(local+phase)%6]
	default:
		return 0, []int{0, -2, -4, -3, -1, 0}[(local+phase)%6]
	}
}

func groomOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileSnakeSlither:
		return []int{0, 1, 2, 1, 0, -1, -2, -1}[(local+phase)%8], 0
	case catalog.MotionProfileTortoisePlod:
		return []int{0, 0, 1, 0, 0, 0, -1, 0}[(local+phase)%8], 0
	case catalog.MotionProfileGeckoCrawl, catalog.MotionProfileDragonPlod:
		return []int{0, 1, 1, 0, 0, -1, -1, 0}[(local+phase)%8], 0
	default:
		return []int{0, 1, 2, 1, 0, -1, -2, -1}[(local+phase)%8], 0
	}
}

func turnOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileSnakeSlither, catalog.MotionProfileGeckoCrawl, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod:
		return []int{0, 1, 0, 1, 0, -1, 0, -1}[(local+phase)%8], 0
	default:
		return []int{0, 1, 0, 1, 0, -1, 0, -1}[(local+phase)%8], []int{0, 1, 0, 1, 0, 1, 0, 1}[(local+phase)%8]
	}
}

func restOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileSnakeSlither, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod:
		return 0, 0
	default:
		return 0, []int{0, -1, -1, 0, 0, 1, 1, 0}[(local+phase)%8]
	}
}

func alertOffset(local int, phase int, profile string) (int, int) {
	switch profile {
	case catalog.MotionProfileSnakeSlither, catalog.MotionProfileGeckoCrawl, catalog.MotionProfileTortoisePlod, catalog.MotionProfileDragonPlod:
		return []int{-2, 0, 2, 1, -1, -2}[(local+phase)%6], 0
	case catalog.MotionProfileRabbitHop, catalog.MotionProfileFrogHop, catalog.MotionProfileBirdHop:
		return []int{-2, 0, 2, 1, -1, -2}[(local+phase)%6], []int{0, -2, -1, 1, 0, -1}[(local+phase)%6]
	default:
		return []int{-2, 0, 2, 1, -1, -2}[(local+phase)%6], []int{0, -1, 0, 1, 0, -1}[(local+phase)%6]
	}
}

func openPNG(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
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

func copyFile(srcPath string, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, data, 0o644)
}

func writeJSON(path string, report []seedReport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create report dir: %w", err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}

func writeSeedPreview(variants []catalog.Variant, outDir string, path string) error {
	if len(variants) == 0 {
		return nil
	}
	const (
		cols = 4
		pad  = 8
	)
	rows := (len(variants) + cols - 1) / cols
	w := cols*frameW + (cols+1)*pad
	h := rows*frameH + (rows+1)*pad
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	drawChecker(dst)
	for i, variant := range variants {
		sheet, err := openPNG(filepath.Join(outDir, fmt.Sprintf("%s_set00.png", variant.SpriteBase)))
		if err != nil {
			return err
		}
		cellX := pad + (i%cols)*(frameW+pad)
		cellY := pad + (i/cols)*(frameH+pad)
		draw.Draw(dst, image.Rect(cellX, cellY, cellX+frameW, cellY+frameH), sheet, image.Point{}, draw.Over)
	}
	return writePNG(path, dst)
}

func drawChecker(dst *image.RGBA) {
	a := color.RGBA{R: 235, G: 237, B: 233, A: 255}
	b := color.RGBA{R: 213, G: 217, B: 211, A: 255}
	for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {
		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			if (x/8+y/8)%2 == 0 {
				dst.SetRGBA(x, y, a)
			} else {
				dst.SetRGBA(x, y, b)
			}
		}
	}
}

func alphaBounds(img image.Image) image.Rectangle {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
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

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
