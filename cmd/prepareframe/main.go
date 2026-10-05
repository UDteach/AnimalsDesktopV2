package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	stddraw "image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	xdraw "golang.org/x/image/draw"
)

const (
	frameW                           = 96
	frameH                           = 64
	targetW                          = 88
	targetH                          = 52
	baselineY                        = 58
	minNormalizedAlphaIoU            = 0.980
	maxNormalizedCentroidShiftPixels = 0.30
)

type prepareReport struct {
	Source               string    `json:"source"`
	SourceSHA256         string    `json:"source_sha256"`
	Output               string    `json:"output"`
	OutputSHA256         string    `json:"output_sha256"`
	SourceWidth          int       `json:"source_width"`
	SourceHeight         int       `json:"source_height"`
	OutputWidth          int       `json:"output_width"`
	OutputHeight         int       `json:"output_height"`
	BackgroundMode       string    `json:"background_mode"`
	BackgroundRemoved    bool      `json:"background_removed"`
	Content              rectJSON  `json:"content"`
	OutputContent        rectJSON  `json:"output_content"`
	Warnings             []string  `json:"warnings,omitempty"`
	GeometryMode         string    `json:"geometry_mode,omitempty"`
	GeometryReference    string    `json:"geometry_reference,omitempty"`
	GeometryReferenceSHA string    `json:"geometry_reference_sha256,omitempty"`
	PreLockOutputContent *rectJSON `json:"pre_lock_output_content,omitempty"`
	ReferenceContent     *rectJSON `json:"reference_content,omitempty"`
	GeometryAdjusted     *bool     `json:"geometry_adjusted,omitempty"`
	ReferenceAreaPX      *int      `json:"reference_area_px,omitempty"`
	CandidateAreaPX      *int      `json:"candidate_area_px,omitempty"`
	IntersectionPX       *int      `json:"intersection_px,omitempty"`
	UnionPX              *int      `json:"union_px,omitempty"`
	IoU                  *float64  `json:"iou,omitempty"`
	CentroidDXPX         *float64  `json:"centroid_dx_px,omitempty"`
	CentroidDYPX         *float64  `json:"centroid_dy_px,omitempty"`
	CentroidShiftPX      *float64  `json:"centroid_shift_px,omitempty"`
}

type rectJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type prepareFrameOptions struct {
	Background     string
	Tolerance      float64
	MatchAlphaBBox string
}

func main() {
	srcPath := flag.String("src", "", "source candidate PNG")
	outPath := flag.String("out", "", "output 96x64 transparent PNG")
	reportPath := flag.String("report", "", "optional JSON report path")
	background := flag.String("background", "auto", "background handling: auto or chroma-green")
	tolerance := flag.Float64("tolerance", 18, "RGB distance tolerance for uniform edge background removal")
	matchAlphaBBox := flag.String("match-alpha-bbox", "", "optional canonical 96x64 PNG whose alpha bounding box must be matched within 1px")
	flag.Parse()

	if *srcPath == "" {
		fatalf("-src is required")
	}
	if *outPath == "" {
		fatalf("-out is required")
	}

	options := prepareFrameOptions{
		Background:     *background,
		Tolerance:      *tolerance,
		MatchAlphaBBox: *matchAlphaBBox,
	}
	if err := requireDistinctPreparePaths(*srcPath, *outPath, *reportPath, *matchAlphaBBox); err != nil {
		fatalf("%v", err)
	}
	report, outputData, err := prepareFrameData(*srcPath, *outPath, options)
	if err != nil {
		fatalf("%v", err)
	}
	payloads := []filePayload{{Path: *outPath, Data: outputData}}
	if *reportPath != "" {
		reportData, encodeErr := encodePrepareReport(report)
		if encodeErr != nil {
			fatalf("encode report: %v", encodeErr)
		}
		payloads = append(payloads, filePayload{Path: *reportPath, Data: reportData})
	}
	if err := writeFileTransaction(payloads); err != nil {
		fatalf("write prepared output: %v", err)
	}
	fmt.Printf("prepared %s -> %s\n", filepath.ToSlash(*srcPath), filepath.ToSlash(*outPath))
}

func prepareFrame(srcPath string, outPath string, tolerance float64) (prepareReport, error) {
	return prepareFrameWithMode(srcPath, outPath, "auto", tolerance)
}

func prepareFrameWithMode(srcPath string, outPath string, background string, tolerance float64) (prepareReport, error) {
	return prepareFrameWithOptions(srcPath, outPath, prepareFrameOptions{
		Background: background,
		Tolerance:  tolerance,
	})
}

func prepareFrameWithOptions(srcPath string, outPath string, options prepareFrameOptions) (prepareReport, error) {
	report, outputData, err := prepareFrameData(srcPath, outPath, options)
	if err != nil {
		return prepareReport{}, err
	}
	if err := writeFileTransaction([]filePayload{{Path: outPath, Data: outputData}}); err != nil {
		return prepareReport{}, err
	}
	return report, nil
}

func prepareFrameData(srcPath string, outPath string, options prepareFrameOptions) (prepareReport, []byte, error) {
	if err := requireDistinctPreparePaths(srcPath, outPath, "", options.MatchAlphaBBox); err != nil {
		return prepareReport{}, nil, err
	}
	sourceHash, err := fileSHA256(srcPath)
	if err != nil {
		return prepareReport{}, nil, fmt.Errorf("hash source: %w", err)
	}
	src, err := openPNG(srcPath)
	if err != nil {
		return prepareReport{}, nil, err
	}
	bounds := src.Bounds()
	report := prepareReport{
		Source:       filepath.ToSlash(srcPath),
		SourceSHA256: sourceHash,
		Output:       filepath.ToSlash(outPath),
		SourceWidth:  bounds.Dx(),
		SourceHeight: bounds.Dy(),
		OutputWidth:  frameW,
		OutputHeight: frameH,
	}

	cleaned := cloneRGBA(src)
	switch options.Background {
	case "auto":
		if hasTransparentAlpha(cleaned) {
			report.BackgroundMode = "source-alpha"
		} else {
			if err := removeUniformEdgeBackground(cleaned, options.Tolerance); err != nil {
				return prepareReport{}, nil, err
			}
			report.BackgroundMode = "uniform-edge-rgb"
			report.BackgroundRemoved = true
		}
	case "chroma-green":
		if err := removeChromaGreenBackground(cleaned); err != nil {
			return prepareReport{}, nil, err
		}
		despillGreen(cleaned)
		report.BackgroundMode = "chroma-green"
		report.BackgroundRemoved = true
	default:
		return prepareReport{}, nil, fmt.Errorf("unknown -background %q", options.Background)
	}
	if hasTransparentAlpha(cleaned) && report.BackgroundMode == "" {
		report.BackgroundMode = "source-alpha"
	}
	clearTransparentRGB(cleaned)

	content := alphaBounds(cleaned, cleaned.Bounds())
	if content.Empty() {
		return prepareReport{}, nil, fmt.Errorf("%s has no visible alpha after background preparation", srcPath)
	}
	if content == cleaned.Bounds() {
		return prepareReport{}, nil, fmt.Errorf("%s still has no transparent background after preparation", srcPath)
	}
	report.Content = rectToJSON(content)
	report.Warnings = frameWarnings(content, cleaned.Bounds())
	if len(report.Warnings) > 0 {
		return prepareReport{}, nil, fmt.Errorf("%s content touches source canvas edge after preparation; background removal or source crop is not clean", srcPath)
	}
	if holeCount, holeArea, largestHole := transparentHoleStats(cleaned, content); holeCount > 0 {
		return prepareReport{}, nil, fmt.Errorf("%s has transparent pinholes after background preparation: holes=%d area=%d largest=%d", srcPath, holeCount, holeArea, largestHole)
	}

	out := fitContent(cleaned, content)
	outContent := alphaBounds(out, out.Bounds())
	if outContent.Empty() {
		return prepareReport{}, nil, fmt.Errorf("%s produced empty output", srcPath)
	}
	if outContent == out.Bounds() {
		return prepareReport{}, nil, fmt.Errorf("%s produced output with no transparent margin", srcPath)
	}

	if options.MatchAlphaBBox != "" {
		referenceImage, referenceContent, err := geometryReferenceData(options.MatchAlphaBBox)
		if err != nil {
			return prepareReport{}, nil, err
		}
		referenceHash, err := fileSHA256(options.MatchAlphaBBox)
		if err != nil {
			return prepareReport{}, nil, fmt.Errorf("hash geometry reference: %w", err)
		}
		if deltas := rectDeltas(outContent, referenceContent); deltas.exceeds(1) {
			return prepareReport{}, nil, fmt.Errorf(
				"%s ordinary output alpha bounds %v differ from geometry reference %s bounds %v by more than 1px (x=%d y=%d w=%d h=%d)",
				srcPath,
				outContent,
				options.MatchAlphaBBox,
				referenceContent,
				deltas.X,
				deltas.Y,
				deltas.W,
				deltas.H,
			)
		}

		preLock := rectToJSON(outContent)
		reference := rectToJSON(referenceContent)
		adjusted := outContent != referenceContent
		report.GeometryMode = "match-alpha-bbox"
		report.GeometryReference = filepath.ToSlash(options.MatchAlphaBBox)
		report.GeometryReferenceSHA = referenceHash
		report.PreLockOutputContent = &preLock
		report.ReferenceContent = &reference
		report.GeometryAdjusted = &adjusted

		if adjusted {
			out = fitContentToRect(cleaned, content, referenceContent)
			clearTransparentRGB(out)
			outContent = alphaBounds(out, out.Bounds())
			if outContent != referenceContent {
				return prepareReport{}, nil, fmt.Errorf(
					"%s geometry lock produced alpha bounds %v, want exact reference bounds %v",
					srcPath,
					outContent,
					referenceContent,
				)
			}
		}

		geometry, err := compareAlphaGeometry(referenceImage, out, out.Bounds())
		if err != nil {
			return prepareReport{}, nil, fmt.Errorf("%s normalized geometry comparison: %w", srcPath, err)
		}
		setAlphaGeometryReport(&report, geometry)
		if err := validateNormalizedAlphaGeometry(geometry); err != nil {
			return prepareReport{}, nil, fmt.Errorf("%s normalized geometry gate failed: %w", srcPath, err)
		}
	}

	report.OutputContent = rectToJSON(outContent)
	report.Warnings = append(report.Warnings, frameWarnings(outContent, out.Bounds())...)

	outputData, err := encodePNG(out)
	if err != nil {
		return prepareReport{}, nil, fmt.Errorf("encode output PNG: %w", err)
	}
	report.OutputSHA256 = sha256Hex(outputData)
	return report, outputData, nil
}

type rectangleDeltas struct {
	X int
	Y int
	W int
	H int
}

func rectDeltas(a image.Rectangle, b image.Rectangle) rectangleDeltas {
	return rectangleDeltas{
		X: absInt(a.Min.X - b.Min.X),
		Y: absInt(a.Min.Y - b.Min.Y),
		W: absInt(a.Dx() - b.Dx()),
		H: absInt(a.Dy() - b.Dy()),
	}
}

func (d rectangleDeltas) exceeds(limit int) bool {
	return d.X > limit || d.Y > limit || d.W > limit || d.H > limit
}

func geometryReferenceData(path string) (*image.RGBA, image.Rectangle, error) {
	reference, err := openPNG(path)
	if err != nil {
		return nil, image.Rectangle{}, fmt.Errorf("open geometry reference: %w", err)
	}
	if reference.Bounds().Dx() != frameW || reference.Bounds().Dy() != frameH {
		return nil, image.Rectangle{}, fmt.Errorf(
			"geometry reference %s is %dx%d, want %dx%d",
			path,
			reference.Bounds().Dx(),
			reference.Bounds().Dy(),
			frameW,
			frameH,
		)
	}

	content := alphaBounds(reference, reference.Bounds())
	if content.Empty() {
		return nil, image.Rectangle{}, fmt.Errorf("geometry reference %s has no visible alpha", path)
	}
	if content == reference.Bounds() {
		return nil, image.Rectangle{}, fmt.Errorf("geometry reference %s has no transparent margin", path)
	}
	if warnings := frameWarnings(content, reference.Bounds()); len(warnings) > 0 {
		return nil, image.Rectangle{}, fmt.Errorf("geometry reference %s alpha touches canvas edge: %s", path, strings.Join(warnings, "; "))
	}
	return reference, content, nil
}

type alphaGeometry struct {
	ReferenceArea int
	CandidateArea int
	Intersection  int
	Union         int
	IoU           float64
	CentroidDX    float64
	CentroidDY    float64
	CentroidShift float64
}

func compareAlphaGeometry(reference image.Image, candidate image.Image, bounds image.Rectangle) (alphaGeometry, error) {
	if bounds.Empty() || !bounds.In(reference.Bounds()) || !bounds.In(candidate.Bounds()) {
		return alphaGeometry{}, fmt.Errorf("comparison bounds %v are outside reference %v or candidate %v", bounds, reference.Bounds(), candidate.Bounds())
	}

	var referenceSumX, referenceSumY int64
	var candidateSumX, candidateSumY int64
	metrics := alphaGeometry{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			referenceSubject := alphaAboveGeometryThreshold(reference, x, y)
			candidateSubject := alphaAboveGeometryThreshold(candidate, x, y)
			localX := int64(x - bounds.Min.X)
			localY := int64(y - bounds.Min.Y)
			if referenceSubject {
				metrics.ReferenceArea++
				referenceSumX += localX
				referenceSumY += localY
			}
			if candidateSubject {
				metrics.CandidateArea++
				candidateSumX += localX
				candidateSumY += localY
			}
			if referenceSubject && candidateSubject {
				metrics.Intersection++
			}
			if referenceSubject || candidateSubject {
				metrics.Union++
			}
		}
	}
	if metrics.ReferenceArea == 0 {
		return alphaGeometry{}, fmt.Errorf("reference alpha mask is empty at alpha > 8")
	}
	if metrics.CandidateArea == 0 {
		return alphaGeometry{}, fmt.Errorf("candidate alpha mask is empty at alpha > 8")
	}
	if metrics.Union == 0 {
		return alphaGeometry{}, errors.New("alpha mask union is empty")
	}

	centroidDenominator := int64(metrics.ReferenceArea) * int64(metrics.CandidateArea)
	metrics.IoU = float64(metrics.Intersection) / float64(metrics.Union)
	metrics.CentroidDX = float64(
		candidateSumX*int64(metrics.ReferenceArea)-referenceSumX*int64(metrics.CandidateArea),
	) / float64(centroidDenominator)
	metrics.CentroidDY = float64(
		candidateSumY*int64(metrics.ReferenceArea)-referenceSumY*int64(metrics.CandidateArea),
	) / float64(centroidDenominator)
	metrics.CentroidShift = math.Hypot(metrics.CentroidDX, metrics.CentroidDY)
	return metrics, nil
}

func alphaAboveGeometryThreshold(img image.Image, x int, y int) bool {
	return color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA).A > 8
}

func setAlphaGeometryReport(report *prepareReport, metrics alphaGeometry) {
	report.ReferenceAreaPX = intPointer(metrics.ReferenceArea)
	report.CandidateAreaPX = intPointer(metrics.CandidateArea)
	report.IntersectionPX = intPointer(metrics.Intersection)
	report.UnionPX = intPointer(metrics.Union)
	report.IoU = floatPointer(metrics.IoU)
	report.CentroidDXPX = floatPointer(metrics.CentroidDX)
	report.CentroidDYPX = floatPointer(metrics.CentroidDY)
	report.CentroidShiftPX = floatPointer(metrics.CentroidShift)
}

func validateNormalizedAlphaGeometry(metrics alphaGeometry) error {
	failures := make([]string, 0, 2)
	if math.IsNaN(metrics.IoU) || math.IsInf(metrics.IoU, 0) ||
		math.IsNaN(metrics.CentroidShift) || math.IsInf(metrics.CentroidShift, 0) {
		return errors.New("geometry metrics are not finite")
	}
	if metrics.IoU < minNormalizedAlphaIoU {
		failures = append(failures, fmt.Sprintf("IoU %g is below %g", metrics.IoU, minNormalizedAlphaIoU))
	}
	if metrics.CentroidShift > maxNormalizedCentroidShiftPixels {
		failures = append(failures, fmt.Sprintf(
			"centroid shift %gpx exceeds %gpx",
			metrics.CentroidShift,
			maxNormalizedCentroidShiftPixels,
		))
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func sameResolvedPath(a string, b string) (bool, error) {
	resolvedA, err := resolvePathForComparison(a)
	if err != nil {
		return false, err
	}
	resolvedB, err := resolvePathForComparison(b)
	if err != nil {
		return false, err
	}
	if samePathString(resolvedA, resolvedB) {
		return true, nil
	}

	infoA, errA := os.Stat(a)
	infoB, errB := os.Stat(b)
	if errA == nil && errB == nil && os.SameFile(infoA, infoB) {
		return true, nil
	}
	return false, nil
}

type namedPreparePath struct {
	Name string
	Path string
}

func requireDistinctPreparePaths(srcPath string, outPath string, reportPath string, referencePath string) error {
	paths := []namedPreparePath{
		{Name: "source", Path: srcPath},
		{Name: "output", Path: outPath},
		{Name: "report", Path: reportPath},
		{Name: "geometry reference", Path: referencePath},
	}
	filtered := paths[:0]
	for _, item := range paths {
		if item.Path != "" {
			filtered = append(filtered, item)
		}
	}
	for i := range filtered {
		for j := i + 1; j < len(filtered); j++ {
			same, err := sameResolvedPath(filtered[i].Path, filtered[j].Path)
			if err != nil {
				return fmt.Errorf("compare %s and %s paths: %w", filtered[i].Name, filtered[j].Name, err)
			}
			if same {
				return fmt.Errorf(
					"%s and %s resolve to the same path: %s",
					filtered[i].Name,
					filtered[j].Name,
					filepath.ToSlash(filtered[i].Path),
				)
			}
		}
	}
	return nil
}

func resolvePathForComparison(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absolute = filepath.Clean(absolute)

	current := absolute
	missing := []string{}
	for {
		resolved, evalErr := filepath.EvalSymlinks(current)
		if evalErr == nil {
			for i := len(missing) - 1; i >= 0; i-- {
				resolved = filepath.Join(resolved, missing[i])
			}
			return filepath.Clean(resolved), nil
		}
		if !os.IsNotExist(evalErr) {
			return "", evalErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			return absolute, nil
		}
		missing = append(missing, filepath.Base(current))
		current = parent
	}
}

func samePathString(a string, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
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
	out := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	stddraw.Draw(out, out.Bounds(), src, bounds.Min, stddraw.Src)
	return out, nil
}

func cloneRGBA(src *image.RGBA) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	copy(out.Pix, src.Pix)
	return out
}

func hasTransparentAlpha(img *image.RGBA) bool {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if img.RGBAAt(x, y).A < 250 {
				return true
			}
		}
	}
	return false
}

func removeUniformEdgeBackground(img *image.RGBA, tolerance float64) error {
	bounds := img.Bounds()
	corners := []color.RGBA{
		img.RGBAAt(bounds.Min.X, bounds.Min.Y),
		img.RGBAAt(bounds.Max.X-1, bounds.Min.Y),
		img.RGBAAt(bounds.Min.X, bounds.Max.Y-1),
		img.RGBAAt(bounds.Max.X-1, bounds.Max.Y-1),
	}
	bg := averageColor(corners)
	for _, corner := range corners {
		if colorDistance(corner, bg) > tolerance {
			return fmt.Errorf("background corners differ; likely checker or noisy background")
		}
	}

	w, h := bounds.Dx(), bounds.Dy()
	seen := make([]bool, w*h)
	queue := make([]image.Point, 0, w*2+h*2)
	add := func(x, y int) {
		if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
			return
		}
		idx := (y-bounds.Min.Y)*w + (x - bounds.Min.X)
		if seen[idx] {
			return
		}
		seen[idx] = true
		queue = append(queue, image.Point{X: x, Y: y})
	}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		add(x, bounds.Min.Y)
		add(x, bounds.Max.Y-1)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		add(bounds.Min.X, y)
		add(bounds.Max.X-1, y)
	}

	removed := 0
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		c := img.RGBAAt(p.X, p.Y)
		if colorDistance(c, bg) > tolerance {
			continue
		}
		c = color.RGBA{}
		img.SetRGBA(p.X, p.Y, c)
		removed++
		add(p.X+1, p.Y)
		add(p.X-1, p.Y)
		add(p.X, p.Y+1)
		add(p.X, p.Y-1)
	}
	if removed == 0 {
		return fmt.Errorf("no edge background pixels matched")
	}
	return nil
}

func removeChromaGreenBackground(img *image.RGBA) error {
	bounds := img.Bounds()
	removed := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if isChromaGreen(c) {
				c = color.RGBA{}
				img.SetRGBA(x, y, c)
				removed++
			}
		}
	}
	if removed == 0 {
		return fmt.Errorf("no chroma-green background pixels matched")
	}
	return nil
}

func clearTransparentRGB(img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A == 0 && (c.R != 0 || c.G != 0 || c.B != 0) {
				img.SetRGBA(x, y, color.RGBA{})
			}
		}
	}
}

func isChromaGreen(c color.RGBA) bool {
	if c.G < 90 {
		return false
	}
	return int(c.G)-int(c.R) >= 25 && int(c.G)-int(c.B) >= 25
}

func despillGreen(img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A == 0 {
				continue
			}
			maxRB := maxInt(int(c.R), int(c.B))
			if int(c.G)-maxRB < 8 {
				continue
			}
			c.G = byte(maxRB)
			img.SetRGBA(x, y, c)
		}
	}
}

func averageColor(colors []color.RGBA) color.RGBA {
	var r, g, b int
	for _, c := range colors {
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
	}
	n := len(colors)
	return color.RGBA{R: byte(r / n), G: byte(g / n), B: byte(b / n), A: 255}
}

func colorDistance(a color.RGBA, b color.RGBA) float64 {
	dr := float64(int(a.R) - int(b.R))
	dg := float64(int(a.G) - int(b.G))
	db := float64(int(a.B) - int(b.B))
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func fitContent(src *image.RGBA, content image.Rectangle) *image.RGBA {
	cropped := image.NewRGBA(image.Rect(0, 0, content.Dx(), content.Dy()))
	stddraw.Draw(cropped, cropped.Bounds(), src, content.Min, stddraw.Src)

	scale := math.Min(float64(targetW)/float64(content.Dx()), float64(targetH)/float64(content.Dy()))
	if scale > 1 {
		scale = 1
	}
	w := maxInt(1, int(math.Round(float64(content.Dx())*scale)))
	h := maxInt(1, int(math.Round(float64(content.Dy())*scale)))
	resized := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), cropped, cropped.Bounds(), xdraw.Over, nil)

	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	x := (frameW - w) / 2
	y := baselineY - h
	if y < 0 {
		y = 0
	}
	stddraw.Draw(out, image.Rect(x, y, x+w, y+h), resized, image.Point{}, stddraw.Over)
	return out
}

func fitContentToRect(src *image.RGBA, content image.Rectangle, destination image.Rectangle) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	xdraw.CatmullRom.Scale(out, destination, src, content, xdraw.Over, nil)
	return out
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

func transparentHoleStats(img image.Image, content image.Rectangle) (holeCount int, holeArea int, largestHole int) {
	content = content.Intersect(img.Bounds())
	if content.Empty() {
		return 0, 0, 0
	}
	width := content.Dx()
	height := content.Dy()
	visited := make([]bool, width*height)
	for y := content.Min.Y; y < content.Max.Y; y++ {
		for x := content.Min.X; x < content.Max.X; x++ {
			idx := (y-content.Min.Y)*width + (x - content.Min.X)
			if visited[idx] || alphaVisible(img, x, y) {
				continue
			}
			area, touchesBoundary := floodTransparentComponent(img, content, x, y, visited)
			if touchesBoundary {
				continue
			}
			holeCount++
			holeArea += area
			if area > largestHole {
				largestHole = area
			}
		}
	}
	return holeCount, holeArea, largestHole
}

func floodTransparentComponent(img image.Image, content image.Rectangle, startX int, startY int, visited []bool) (area int, touchesBoundary bool) {
	width := content.Dx()
	stack := []image.Point{{X: startX, Y: startY}}
	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !point.In(content) {
			continue
		}
		idx := (point.Y-content.Min.Y)*width + (point.X - content.Min.X)
		if visited[idx] || alphaVisible(img, point.X, point.Y) {
			continue
		}
		visited[idx] = true
		area++
		if point.X == content.Min.X || point.X == content.Max.X-1 || point.Y == content.Min.Y || point.Y == content.Max.Y-1 {
			touchesBoundary = true
		}
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				stack = append(stack, image.Point{X: point.X + dx, Y: point.Y + dy})
			}
		}
	}
	return area, touchesBoundary
}

func alphaVisible(img image.Image, x int, y int) bool {
	_, _, _, a := img.At(x, y).RGBA()
	return a > 0x0800
}

func rectToJSON(rect image.Rectangle) rectJSON {
	return rectJSON{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()}
}

func intPointer(value int) *int {
	return &value
}

func floatPointer(value float64) *float64 {
	return &value
}

func encodePNG(img image.Image) ([]byte, error) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func encodePrepareReport(report prepareReport) ([]byte, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return sha256Hex(data), nil
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

type filePayload struct {
	Path string
	Data []byte
}

type stagedFilePayload struct {
	filePayload
	TempPath    string
	BackupPath  string
	HadOriginal bool
	Committed   bool
}

var (
	createTempFile = os.CreateTemp
	renameFile     = os.Rename
	removeFile     = os.Remove
)

// writeFileTransaction is failure-transactional for ordinary filesystem
// errors: every payload is staged before any destination changes, and existing
// files are restored if a later rename fails. No portable API can make a
// multi-path transaction power-loss atomic, so a process or machine crash may
// leave hidden stage/backup files for manual recovery.
func writeFileTransaction(payloads []filePayload) error {
	if len(payloads) == 0 {
		return nil
	}
	for i := range payloads {
		if payloads[i].Path == "" {
			return fmt.Errorf("payload %d has an empty path", i)
		}
		for j := i + 1; j < len(payloads); j++ {
			same, err := sameResolvedPath(payloads[i].Path, payloads[j].Path)
			if err != nil {
				return fmt.Errorf("compare output paths: %w", err)
			}
			if same {
				return fmt.Errorf("transaction outputs resolve to the same path: %s", filepath.ToSlash(payloads[i].Path))
			}
		}
	}

	staged := make([]stagedFilePayload, 0, len(payloads))
	cleanupTemps := func() {
		for i := range staged {
			if staged[i].TempPath != "" {
				_ = removeFile(staged[i].TempPath)
			}
		}
	}
	for _, payload := range payloads {
		dir := filepath.Dir(payload.Path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			cleanupTemps()
			return err
		}
		temp, err := createTempFile(dir, "."+filepath.Base(payload.Path)+".stage-*")
		if err != nil {
			cleanupTemps()
			return err
		}
		tempPath := temp.Name()
		if _, err := temp.Write(payload.Data); err != nil {
			_ = temp.Close()
			_ = removeFile(tempPath)
			cleanupTemps()
			return err
		}
		if err := temp.Sync(); err != nil {
			_ = temp.Close()
			_ = removeFile(tempPath)
			cleanupTemps()
			return err
		}
		if err := temp.Close(); err != nil {
			_ = removeFile(tempPath)
			cleanupTemps()
			return err
		}
		staged = append(staged, stagedFilePayload{filePayload: payload, TempPath: tempPath})
	}

	restoreBackups := func() {
		for i := len(staged) - 1; i >= 0; i-- {
			item := &staged[i]
			if item.Committed {
				_ = removeFile(item.Path)
			}
			if item.HadOriginal && item.BackupPath != "" {
				_ = renameFile(item.BackupPath, item.Path)
				item.BackupPath = ""
			}
		}
		cleanupTemps()
	}

	for i := range staged {
		item := &staged[i]
		info, err := os.Stat(item.Path)
		switch {
		case err == nil:
			if !info.Mode().IsRegular() {
				restoreBackups()
				return fmt.Errorf("destination is not a regular file: %s", filepath.ToSlash(item.Path))
			}
			backup, createErr := createTempFile(filepath.Dir(item.Path), "."+filepath.Base(item.Path)+".backup-*")
			if createErr != nil {
				restoreBackups()
				return createErr
			}
			item.BackupPath = backup.Name()
			if closeErr := backup.Close(); closeErr != nil {
				_ = removeFile(item.BackupPath)
				item.BackupPath = ""
				restoreBackups()
				return closeErr
			}
			if removeErr := removeFile(item.BackupPath); removeErr != nil {
				item.BackupPath = ""
				restoreBackups()
				return removeErr
			}
			if renameErr := renameFile(item.Path, item.BackupPath); renameErr != nil {
				item.BackupPath = ""
				restoreBackups()
				return renameErr
			}
			item.HadOriginal = true
		case os.IsNotExist(err):
			// No backup is needed.
		default:
			restoreBackups()
			return err
		}
	}

	for i := range staged {
		item := &staged[i]
		if err := renameFile(item.TempPath, item.Path); err != nil {
			restoreBackups()
			return err
		}
		item.TempPath = ""
		item.Committed = true
	}
	for i := range staged {
		if staged[i].BackupPath != "" {
			_ = removeFile(staged[i].BackupPath)
			staged[i].BackupPath = ""
		}
	}
	return nil
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
