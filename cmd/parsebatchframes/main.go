package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const (
	frameW = 96
	frameH = 64
)

type manifest struct {
	Animal      string          `json:"animal,omitempty"`
	Sequence    string          `json:"sequence,omitempty"`
	Layout      string          `json:"layout"`
	CellWidth   int             `json:"cell_width,omitempty"`
	CellHeight  int             `json:"cell_height,omitempty"`
	Columns     int             `json:"columns,omitempty"`
	Rows        int             `json:"rows,omitempty"`
	OriginX     int             `json:"origin_x,omitempty"`
	OriginY     int             `json:"origin_y,omitempty"`
	PitchX      int             `json:"pitch_x,omitempty"`
	PitchY      int             `json:"pitch_y,omitempty"`
	GuardPixels int             `json:"guard_pixels,omitempty"`
	Background  string          `json:"background,omitempty"`
	Frames      []manifestFrame `json:"frames"`
}

type manifestFrame struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Cell   int    `json:"cell,omitempty"`
	Output string `json:"output,omitempty"`
}

type parseReport struct {
	Manifest    string         `json:"manifest"`
	Animal      string         `json:"animal,omitempty"`
	Sequence    string         `json:"sequence,omitempty"`
	Layout      string         `json:"layout"`
	OutputDir   string         `json:"output_dir"`
	CellWidth   int            `json:"cell_width"`
	CellHeight  int            `json:"cell_height"`
	Columns     int            `json:"columns,omitempty"`
	Rows        int            `json:"rows,omitempty"`
	OriginX     int            `json:"origin_x,omitempty"`
	OriginY     int            `json:"origin_y,omitempty"`
	PitchX      int            `json:"pitch_x,omitempty"`
	PitchY      int            `json:"pitch_y,omitempty"`
	GuardPixels int            `json:"guard_pixels"`
	Background  string         `json:"background"`
	SourceCount int            `json:"source_count"`
	FrameCount  int            `json:"frame_count"`
	Parsed      int            `json:"parsed"`
	Rejected    int            `json:"rejected"`
	Sources     []sourceReport `json:"sources"`
	Frames      []frameReport  `json:"frames"`
	Errors      []string       `json:"errors,omitempty"`
}

type sourceReport struct {
	Path     string   `json:"path"`
	SHA256   string   `json:"sha256,omitempty"`
	Width    int      `json:"width,omitempty"`
	Height   int      `json:"height,omitempty"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type frameReport struct {
	ID                 string   `json:"id"`
	Source             string   `json:"source"`
	Cell               int      `json:"cell,omitempty"`
	CellBounds         rectJSON `json:"cell_bounds"`
	Output             string   `json:"output,omitempty"`
	SourceSHA256       string   `json:"source_sha256,omitempty"`
	OutputSHA256       string   `json:"output_sha256,omitempty"`
	Content            rectJSON `json:"content,omitempty"`
	ComponentCount     int      `json:"component_count,omitempty"`
	DetachedComponents int      `json:"detached_components,omitempty"`
	DetachedArea       int      `json:"detached_area,omitempty"`
	LargestDetached    int      `json:"largest_detached,omitempty"`
	BackgroundHoles    int      `json:"background_holes,omitempty"`
	BackgroundHoleArea int      `json:"background_hole_area,omitempty"`
	LargestBgHole      int      `json:"largest_bg_hole,omitempty"`
	Status             string   `json:"status"`
	Warnings           []string `json:"warnings,omitempty"`
	Errors             []string `json:"errors,omitempty"`
}

type rectJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type sourceImage struct {
	report sourceReport
	img    image.Image
}

func main() {
	manifestPath := flag.String("manifest", "", "batch manifest JSON")
	outDir := flag.String("out", "", "review output directory for parsed PNGs")
	reportPath := flag.String("report", "", "optional parse report JSON path")
	flag.Parse()

	if *manifestPath == "" {
		fatalf("-manifest is required")
	}
	if *outDir == "" {
		fatalf("-out is required")
	}

	report, err := parseBatch(*manifestPath, *outDir, *reportPath)
	if *reportPath != "" && report != nil {
		if writeErr := writeReport(*reportPath, *report); writeErr != nil {
			fatalf("write report: %v", writeErr)
		}
	}
	if err != nil {
		fatalf("%v", err)
	}
	fmt.Printf("parsed batch frames: parsed=%d rejected=%d layout=%s\n", report.Parsed, report.Rejected, report.Layout)
}

func parseBatch(manifestPath string, outDir string, reportPath string) (*parseReport, error) {
	if pathHasSegment(outDir, "accepted-frames") {
		return nil, fmt.Errorf("-out must not be under accepted-frames: %s", outDir)
	}
	if reportPath != "" && pathHasSegment(reportPath, "accepted-frames") {
		return nil, fmt.Errorf("-report must not be under accepted-frames: %s", reportPath)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}
	if err := normalizeManifest(&m); err != nil {
		return nil, err
	}

	report := parseReport{
		Manifest:    filepath.ToSlash(manifestPath),
		Animal:      m.Animal,
		Sequence:    m.Sequence,
		Layout:      m.Layout,
		OutputDir:   filepath.ToSlash(outDir),
		CellWidth:   m.CellWidth,
		CellHeight:  m.CellHeight,
		Columns:     m.Columns,
		Rows:        m.Rows,
		OriginX:     m.OriginX,
		OriginY:     m.OriginY,
		PitchX:      m.PitchX,
		PitchY:      m.PitchY,
		GuardPixels: m.GuardPixels,
		Background:  m.Background,
		FrameCount:  len(m.Frames),
		Sources:     []sourceReport{},
		Frames:      []frameReport{},
	}

	manifestDir := filepath.Dir(manifestPath)
	sources := map[string]sourceImage{}
	sourceKeys := []string{}
	for _, frame := range m.Frames {
		sourcePath, err := resolveSource(manifestDir, frame.Source)
		if err != nil {
			report.Errors = append(report.Errors, err.Error())
			continue
		}
		if _, ok := sources[sourcePath]; ok {
			continue
		}
		src := readSource(sourcePath)
		src.report.Errors = append(src.report.Errors, validateSourceDimensions(m, src.img)...)
		src.report.Errors = append(src.report.Errors, validateSourceBoundaries(m, src.img)...)
		sources[sourcePath] = src
		sourceKeys = append(sourceKeys, sourcePath)
	}
	for _, key := range sourceKeys {
		report.Sources = append(report.Sources, sources[key].report)
	}
	report.SourceCount = len(report.Sources)

	seenCells := map[string]bool{}
	for index, frame := range m.Frames {
		frameReport := parseFrame(m, manifestDir, outDir, frame, index, sources, seenCells)
		if frameReport.Status == "parsed" {
			report.Parsed++
		} else {
			report.Rejected++
		}
		report.Frames = append(report.Frames, frameReport)
	}

	if len(report.Errors) > 0 {
		return &report, fmt.Errorf("parse batch rejected: %d manifest/source error(s)", len(report.Errors))
	}
	for _, src := range report.Sources {
		if len(src.Errors) > 0 {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", src.Path, strings.Join(src.Errors, "; ")))
		}
	}
	for _, frame := range report.Frames {
		if len(frame.Errors) > 0 {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", frame.ID, strings.Join(frame.Errors, "; ")))
		}
	}
	if len(report.Errors) > 0 {
		return &report, fmt.Errorf("parse batch rejected: %d error(s)", len(report.Errors))
	}
	return &report, nil
}

func normalizeManifest(m *manifest) error {
	m.Layout = strings.TrimSpace(m.Layout)
	if m.Layout == "" {
		return fmt.Errorf("manifest layout is required")
	}
	if m.CellWidth == 0 {
		m.CellWidth = frameW
	}
	if m.CellHeight == 0 {
		m.CellHeight = frameH
	}
	if m.GuardPixels == 0 {
		m.GuardPixels = 1
	}
	if strings.TrimSpace(m.Background) == "" {
		m.Background = "alpha"
	}
	switch m.Background {
	case "alpha", "chroma-green":
	default:
		return fmt.Errorf("unsupported background %q", m.Background)
	}
	if m.Layout != "grid-fixed" && (m.CellWidth != frameW || m.CellHeight != frameH) {
		return fmt.Errorf("cell size must be %dx%d, got %dx%d", frameW, frameH, m.CellWidth, m.CellHeight)
	}
	if m.CellWidth <= 0 || m.CellHeight <= 0 {
		return fmt.Errorf("cell size must be positive, got %dx%d", m.CellWidth, m.CellHeight)
	}
	if m.GuardPixels < 0 || m.GuardPixels > minInt(m.CellWidth/2, m.CellHeight/2) {
		return fmt.Errorf("guard_pixels out of range: %d", m.GuardPixels)
	}
	switch m.Layout {
	case "independent", "strip-2x1", "grid-2x2":
	case "grid-fixed":
		if m.Columns <= 0 || m.Rows <= 0 {
			return fmt.Errorf("grid-fixed requires positive columns and rows")
		}
		if m.OriginX < 0 || m.OriginY < 0 {
			return fmt.Errorf("grid-fixed origin must be non-negative")
		}
		if m.PitchX == 0 {
			m.PitchX = m.CellWidth
		}
		if m.PitchY == 0 {
			m.PitchY = m.CellHeight
		}
		if m.PitchX < m.CellWidth || m.PitchY < m.CellHeight {
			return fmt.Errorf("grid-fixed pitch must be at least cell size")
		}
	default:
		return fmt.Errorf("unsupported layout %q", m.Layout)
	}
	if len(m.Frames) == 0 {
		return fmt.Errorf("manifest has no frames")
	}
	return nil
}

func readSource(path string) sourceImage {
	report := sourceReport{Path: filepath.ToSlash(path)}
	data, err := os.ReadFile(path)
	if err != nil {
		report.Errors = append(report.Errors, err.Error())
		return sourceImage{report: report}
	}
	sum := sha256.Sum256(data)
	report.SHA256 = hex.EncodeToString(sum[:])
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("decode PNG: %v", err))
		return sourceImage{report: report}
	}
	bounds := img.Bounds()
	report.Width = bounds.Dx()
	report.Height = bounds.Dy()
	return sourceImage{report: report, img: img}
}

func validateSourceDimensions(m manifest, img image.Image) []string {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	switch m.Layout {
	case "strip-2x1":
		if bounds.Dx() != m.CellWidth*2 || bounds.Dy() != m.CellHeight {
			return []string{fmt.Sprintf("bounds = %dx%d, want %dx%d for strip-2x1", bounds.Dx(), bounds.Dy(), m.CellWidth*2, m.CellHeight)}
		}
	case "grid-2x2":
		if bounds.Dx() != m.CellWidth*2 || bounds.Dy() != m.CellHeight*2 {
			return []string{fmt.Sprintf("bounds = %dx%d, want %dx%d for grid-2x2", bounds.Dx(), bounds.Dy(), m.CellWidth*2, m.CellHeight*2)}
		}
	case "grid-fixed":
		wantW := m.OriginX + (m.Columns-1)*m.PitchX + m.CellWidth
		wantH := m.OriginY + (m.Rows-1)*m.PitchY + m.CellHeight
		if bounds.Dx() < wantW || bounds.Dy() < wantH {
			return []string{fmt.Sprintf("bounds = %dx%d, want at least %dx%d for grid-fixed", bounds.Dx(), bounds.Dy(), wantW, wantH)}
		}
	}
	return nil
}

func validateSourceBoundaries(m manifest, img image.Image) []string {
	if img == nil || m.Layout == "independent" {
		return nil
	}
	bounds := img.Bounds()
	errors := []string{}
	for _, x := range boundaryCentersX(m, bounds) {
		for _, guardX := range boundaryRange(x, m.GuardPixels) {
			if columnHasContent(m, img, bounds, guardX) {
				errors = append(errors, fmt.Sprintf("content exists in x-boundary guard column %d", guardX-bounds.Min.X))
				return errors
			}
		}
	}
	for _, y := range boundaryCentersY(m, bounds) {
		for _, guardY := range boundaryRange(y, m.GuardPixels) {
			if rowHasContent(m, img, bounds, guardY) {
				errors = append(errors, fmt.Sprintf("content exists in y-boundary guard row %d", guardY-bounds.Min.Y))
				return errors
			}
		}
	}
	return errors
}

func boundaryCentersX(m manifest, bounds image.Rectangle) []int {
	switch m.Layout {
	case "strip-2x1", "grid-2x2":
		return []int{bounds.Min.X + m.CellWidth}
	case "grid-fixed":
		values := []int{}
		if m.PitchX == m.CellWidth {
			for col := 1; col < m.Columns; col++ {
				values = append(values, bounds.Min.X+m.OriginX+col*m.PitchX)
			}
		}
		return values
	default:
		return nil
	}
}

func boundaryCentersY(m manifest, bounds image.Rectangle) []int {
	switch m.Layout {
	case "grid-2x2":
		return []int{bounds.Min.Y + m.CellHeight}
	case "grid-fixed":
		values := []int{}
		if m.PitchY == m.CellHeight {
			for row := 1; row < m.Rows; row++ {
				values = append(values, bounds.Min.Y+m.OriginY+row*m.PitchY)
			}
		}
		return values
	default:
		return nil
	}
}

func boundaryRange(center int, guard int) []int {
	if guard <= 0 {
		return nil
	}
	values := make([]int, 0, guard*2)
	for x := center - guard; x < center+guard; x++ {
		values = append(values, x)
	}
	return values
}

func parseFrame(m manifest, manifestDir string, outDir string, frame manifestFrame, index int, sources map[string]sourceImage, seenCells map[string]bool) frameReport {
	id := strings.TrimSpace(frame.ID)
	if id == "" {
		id = fmt.Sprintf("frame-%02d", index)
	}
	outName, outputErr := outputName(id, frame.Output)
	report := frameReport{
		ID:     id,
		Source: filepath.ToSlash(frame.Source),
		Cell:   frame.Cell,
		Status: "rejected",
	}
	if outputErr != nil {
		report.Errors = append(report.Errors, outputErr.Error())
		return report
	}
	sourcePath, err := resolveSource(manifestDir, frame.Source)
	if err != nil {
		report.Errors = append(report.Errors, err.Error())
		return report
	}
	src, ok := sources[sourcePath]
	if !ok {
		report.Errors = append(report.Errors, "source was not loaded")
		return report
	}
	report.Source = filepath.ToSlash(sourcePath)
	report.SourceSHA256 = src.report.SHA256
	if len(src.report.Errors) > 0 {
		report.Errors = append(report.Errors, src.report.Errors...)
		return report
	}

	cellBounds, err := cellBounds(m, frame.Cell, src.img.Bounds())
	if err != nil {
		report.Errors = append(report.Errors, err.Error())
		return report
	}
	report.CellBounds = rectToJSON(cellBounds)
	if m.Layout != "independent" {
		key := sourcePath + "#" + fmt.Sprint(frame.Cell)
		if seenCells[key] {
			report.Errors = append(report.Errors, fmt.Sprintf("duplicate source/cell entry %s", key))
			return report
		}
		seenCells[key] = true
	}

	cell := cropImage(src.img, cellBounds)
	content := contentBounds(m, cell, cell.Bounds())
	if content.Empty() {
		report.Errors = append(report.Errors, "no visible content in parsed frame")
		return report
	}
	report.Content = rectToJSON(content)
	if content == cell.Bounds() {
		report.Errors = append(report.Errors, "no transparent or chroma-key background in parsed frame")
	}
	if warnings := edgeWarnings(content, cell.Bounds()); len(warnings) > 0 {
		if m.Layout == "independent" {
			report.Warnings = append(report.Warnings, warnings...)
		} else {
			report.Errors = append(report.Errors, warnings...)
		}
	}
	componentCount, detachedCount, detachedArea, largestDetached := contentComponentStats(m, cell, content)
	report.ComponentCount = componentCount
	report.DetachedComponents = detachedCount
	report.DetachedArea = detachedArea
	report.LargestDetached = largestDetached
	if detachedCount > 0 {
		report.Errors = append(report.Errors, fmt.Sprintf("disconnected alpha components: components=%d detached=%d detached_area=%d largest_detached=%d", componentCount, detachedCount, detachedArea, largestDetached))
	}
	holeCount, holeArea, largestHole := backgroundHoleStats(m, cell, content)
	report.BackgroundHoles = holeCount
	report.BackgroundHoleArea = holeArea
	report.LargestBgHole = largestHole
	if holeCount > 0 {
		report.Errors = append(report.Errors, fmt.Sprintf("transparent/chroma pinholes: holes=%d area=%d largest=%d", holeCount, holeArea, largestHole))
	}
	if len(report.Errors) > 0 {
		return report
	}

	outPath := filepath.Join(outDir, outName)
	if pathHasSegment(outPath, "accepted-frames") {
		report.Errors = append(report.Errors, "resolved output path is under accepted-frames")
		return report
	}
	if err := writePNG(outPath, cell); err != nil {
		report.Errors = append(report.Errors, err.Error())
		return report
	}
	outputSHA, err := fileSHA256(outPath)
	if err != nil {
		report.Errors = append(report.Errors, err.Error())
		return report
	}
	report.Output = filepath.ToSlash(outPath)
	report.OutputSHA256 = outputSHA
	report.Status = "parsed"
	return report
}

func resolveSource(manifestDir string, source string) (string, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return "", fmt.Errorf("frame source is required")
	}
	if filepath.IsAbs(source) {
		return filepath.Clean(source), nil
	}
	return filepath.Clean(filepath.Join(manifestDir, source)), nil
}

func outputName(id string, output string) (string, error) {
	name := strings.TrimSpace(output)
	if name == "" {
		name = sanitizeFilename(id) + ".png"
	}
	if name == "." || name == ".." || strings.Contains(name, "/") || strings.Contains(name, "\\") || filepath.IsAbs(name) {
		return "", fmt.Errorf("output must be a filename, got %q", output)
	}
	if strings.Contains(name, "..") {
		return "", fmt.Errorf("output must not contain '..', got %q", output)
	}
	if !strings.HasSuffix(strings.ToLower(name), ".png") {
		name += ".png"
	}
	return name, nil
}

func sanitizeFilename(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "frame"
	}
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	name := strings.Trim(b.String(), "._-")
	if name == "" {
		return "frame"
	}
	return name
}

func cellBounds(m manifest, cell int, bounds image.Rectangle) (image.Rectangle, error) {
	switch m.Layout {
	case "independent":
		return bounds, nil
	case "strip-2x1":
		if cell < 0 || cell > 1 {
			return image.Rectangle{}, fmt.Errorf("cell %d out of range for strip-2x1", cell)
		}
		x := bounds.Min.X + cell*m.CellWidth
		return image.Rect(x, bounds.Min.Y, x+m.CellWidth, bounds.Min.Y+m.CellHeight), nil
	case "grid-2x2":
		if cell < 0 || cell > 3 {
			return image.Rectangle{}, fmt.Errorf("cell %d out of range for grid-2x2", cell)
		}
		col := cell % 2
		row := cell / 2
		x := bounds.Min.X + col*m.CellWidth
		y := bounds.Min.Y + row*m.CellHeight
		return image.Rect(x, y, x+m.CellWidth, y+m.CellHeight), nil
	case "grid-fixed":
		if cell < 0 || cell >= m.Columns*m.Rows {
			return image.Rectangle{}, fmt.Errorf("cell %d out of range for grid-fixed %dx%d", cell, m.Columns, m.Rows)
		}
		col := cell % m.Columns
		row := cell / m.Columns
		x := bounds.Min.X + m.OriginX + col*m.PitchX
		y := bounds.Min.Y + m.OriginY + row*m.PitchY
		return image.Rect(x, y, x+m.CellWidth, y+m.CellHeight), nil
	default:
		return image.Rectangle{}, fmt.Errorf("unsupported layout %q", m.Layout)
	}
}

func cropImage(src image.Image, rect image.Rectangle) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(out, out.Bounds(), src, rect.Min, draw.Src)
	return out
}

func contentBounds(m manifest, img image.Image, rect image.Rectangle) image.Rectangle {
	rect = rect.Intersect(img.Bounds())
	minX, minY := rect.Max.X, rect.Max.Y
	maxX, maxY := rect.Min.X, rect.Min.Y
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if !contentVisible(m, img, x, y) {
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

func edgeWarnings(content image.Rectangle, bounds image.Rectangle) []string {
	warnings := []string{}
	if content.Min.X <= bounds.Min.X || content.Max.X >= bounds.Max.X {
		warnings = append(warnings, "alpha touches horizontal canvas edge")
	}
	if content.Min.Y <= bounds.Min.Y || content.Max.Y >= bounds.Max.Y {
		warnings = append(warnings, "alpha touches vertical canvas edge")
	}
	return warnings
}

func contentComponentStats(m manifest, img image.Image, content image.Rectangle) (componentCount int, detachedCount int, detachedArea int, largestDetached int) {
	width := content.Dx()
	height := content.Dy()
	if width <= 0 || height <= 0 {
		return 0, 0, 0, 0
	}
	visited := make([]bool, width*height)
	componentAreas := []int{}
	for y := content.Min.Y; y < content.Max.Y; y++ {
		for x := content.Min.X; x < content.Max.X; x++ {
			idx := (y-content.Min.Y)*width + (x - content.Min.X)
			if visited[idx] || !contentVisible(m, img, x, y) {
				continue
			}
			area := floodContentComponent(m, img, content, x, y, visited)
			componentAreas = append(componentAreas, area)
		}
	}
	if len(componentAreas) <= 1 {
		return len(componentAreas), 0, 0, 0
	}
	mainAreaIndex := -1
	mainArea := 0
	for i, area := range componentAreas {
		if area > mainArea {
			mainAreaIndex = i
			mainArea = area
		}
	}
	for i, area := range componentAreas {
		if i == mainAreaIndex {
			continue
		}
		if area < 4 {
			continue
		}
		detachedCount++
		detachedArea += area
		if area > largestDetached {
			largestDetached = area
		}
	}
	return len(componentAreas), detachedCount, detachedArea, largestDetached
}

func floodContentComponent(m manifest, img image.Image, content image.Rectangle, startX int, startY int, visited []bool) int {
	width := content.Dx()
	stack := []image.Point{{X: startX, Y: startY}}
	area := 0
	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !point.In(content) {
			continue
		}
		idx := (point.Y-content.Min.Y)*width + (point.X - content.Min.X)
		if visited[idx] || !contentVisible(m, img, point.X, point.Y) {
			continue
		}
		visited[idx] = true
		area++
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				stack = append(stack, image.Point{X: point.X + dx, Y: point.Y + dy})
			}
		}
	}
	return area
}

func backgroundHoleStats(m manifest, img image.Image, content image.Rectangle) (holeCount int, holeArea int, largestHole int) {
	width := content.Dx()
	height := content.Dy()
	if width <= 0 || height <= 0 {
		return 0, 0, 0
	}
	visited := make([]bool, width*height)
	for y := content.Min.Y; y < content.Max.Y; y++ {
		for x := content.Min.X; x < content.Max.X; x++ {
			idx := (y-content.Min.Y)*width + (x - content.Min.X)
			if visited[idx] || contentVisible(m, img, x, y) {
				continue
			}
			area, touchesBoundary := floodBackgroundComponent(m, img, content, x, y, visited)
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

func floodBackgroundComponent(m manifest, img image.Image, content image.Rectangle, startX int, startY int, visited []bool) (area int, touchesBoundary bool) {
	width := content.Dx()
	stack := []image.Point{{X: startX, Y: startY}}
	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !point.In(content) {
			continue
		}
		idx := (point.Y-content.Min.Y)*width + (point.X - content.Min.X)
		if visited[idx] || contentVisible(m, img, point.X, point.Y) {
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

func contentVisible(m manifest, img image.Image, x int, y int) bool {
	if m.Background == "chroma-green" {
		return !isChromaGreen(img.At(x, y))
	}
	return alphaVisible(img, x, y)
}

func alphaVisible(img image.Image, x int, y int) bool {
	_, _, _, a := img.At(x, y).RGBA()
	return a > 0x0800
}

func isChromaGreen(c color.Color) bool {
	r16, g16, b16, a16 := c.RGBA()
	if a16 <= 0x0800 {
		return true
	}
	r := int(r16 >> 8)
	g := int(g16 >> 8)
	b := int(b16 >> 8)
	return absInt(r-0) <= 8 && absInt(g-255) <= 8 && absInt(b-0) <= 8
}

func columnHasContent(m manifest, img image.Image, bounds image.Rectangle, x int) bool {
	if x < bounds.Min.X || x >= bounds.Max.X {
		return false
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if contentVisible(m, img, x, y) {
			return true
		}
	}
	return false
}

func rowHasContent(m manifest, img image.Image, bounds image.Rectangle, y int) bool {
	if y < bounds.Min.Y || y >= bounds.Max.Y {
		return false
	}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		if contentVisible(m, img, x, y) {
			return true
		}
	}
	return false
}

func rectToJSON(rect image.Rectangle) rectJSON {
	return rectJSON{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()}
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

func writeReport(path string, report parseReport) error {
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

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func pathHasSegment(path string, segment string) bool {
	segment = strings.ToLower(segment)
	parts := strings.FieldsFunc(filepath.Clean(path), func(r rune) bool {
		return r == '/' || r == '\\'
	})
	for _, part := range parts {
		if strings.ToLower(part) == segment {
			return true
		}
	}
	return false
}

func minInt(a int, b int) int {
	if a < b {
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
