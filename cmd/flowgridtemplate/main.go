package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	stddraw "image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	xdraw "golang.org/x/image/draw"
)

const (
	defaultCanvas = 2048
	defaultMargin = 48
	defaultGutter = 16
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

func main() {
	outDir := flag.String("out-dir", "", "directory for the Flow grid experiment pack")
	animal := flag.String("animal", "chinchilla standard gray", "animal and coat label")
	sequence := flag.String("sequence", "set00", "sequence id")
	source := flag.String("source", "flow-output.png", "expected Flow output PNG filename")
	frames := flag.Int("frames", 16, "number of frame cells to request")
	frameStart := flag.Int("frame-start", 0, "first frame number")
	frameList := flag.String("frame-list", "", "comma-separated explicit frame numbers; overrides -frames and -frame-start")
	columns := flag.Int("columns", 0, "grid columns; default is derived from frame count")
	rows := flag.Int("rows", 0, "grid rows; default is derived from frame count")
	canvas := flag.Int("canvas", defaultCanvas, "square guide canvas size")
	margin := flag.Int("margin", defaultMargin, "outer margin")
	gutter := flag.Int("gutter", defaultGutter, "gutter between crop cells")
	anchors := flag.String("anchors", "", "comma-separated accepted PNG anchors to copy into the pack")
	flag.Parse()

	if *outDir == "" {
		fatalf("-out-dir is required")
	}
	frameNumbers, err := selectedFrameNumbers(*frameList, *frames, *frameStart)
	if err != nil {
		fatalf("%v", err)
	}
	frameCount := len(frameNumbers)
	cols, gridRows := gridShape(frameCount, *columns, *rows)
	if cols <= 0 || gridRows <= 0 || cols*gridRows < frameCount {
		fatalf("invalid grid shape: cols=%d rows=%d frames=%d", cols, gridRows, frameCount)
	}
	if *canvas <= 0 || *margin < 0 || *gutter < 0 {
		fatalf("invalid canvas/margin/gutter")
	}
	cellW := (*canvas - *margin*2 - (cols-1)*(*gutter)) / cols
	cellH := (*canvas - *margin*2 - (gridRows-1)*(*gutter)) / gridRows
	if cellW <= 0 || cellH <= 0 {
		fatalf("invalid cell size: %dx%d", cellW, cellH)
	}
	pitchX := cellW + *gutter
	pitchY := cellH + *gutter

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fatalf("mkdir: %v", err)
	}
	pack := manifest{
		Animal:      *animal,
		Sequence:    *sequence,
		Layout:      "grid-fixed",
		CellWidth:   cellW,
		CellHeight:  cellH,
		Columns:     cols,
		Rows:        gridRows,
		OriginX:     *margin,
		OriginY:     *margin,
		PitchX:      pitchX,
		PitchY:      pitchY,
		GuardPixels: 4,
		Background:  "chroma-green",
		Frames:      make([]manifestFrame, 0, frameCount),
	}
	for i, frameNo := range frameNumbers {
		id := fmt.Sprintf("%s-frame-%02d", *sequence, frameNo)
		pack.Frames = append(pack.Frames, manifestFrame{
			ID:     id,
			Source: *source,
			Cell:   i,
			Output: fmt.Sprintf("frame-%02d.png", frameNo),
		})
	}

	anchorPaths := splitList(*anchors)
	if err := validateAnchorImages(anchorPaths); err != nil {
		fatalf("%v", err)
	}
	guide := newGuide(*canvas, *margin, cols, gridRows, cellW, cellH, pitchX, pitchY, *frames, false)
	if err := writePNG(filepath.Join(*outDir, "grid-guide.png"), guide); err != nil {
		fatalf("write guide: %v", err)
	}
	detailedGuide := newGuide(*canvas, *margin, cols, gridRows, cellW, cellH, pitchX, pitchY, *frames, true)
	if err := writePNG(filepath.Join(*outDir, "grid-guide-detailed.png"), detailedGuide); err != nil {
		fatalf("write detailed guide: %v", err)
	}
	if len(anchorPaths) > 0 {
		if err := writeStyleAnchors(filepath.Join(*outDir, "style-anchors.png"), anchorPaths); err != nil {
			fatalf("write style anchors: %v", err)
		}
		if err := writeStyleAnchorsNeutral(filepath.Join(*outDir, "style-anchors-neutral.png"), anchorPaths); err != nil {
			fatalf("write neutral style anchors: %v", err)
		}
		if err := writeStyleAnchorsMatted(filepath.Join(*outDir, "style-anchors-green.png"), anchorPaths); err != nil {
			fatalf("write matted style anchors: %v", err)
		}
		seed := cloneRGBA(guide)
		if err := drawAnchorInFirstCell(seed, anchorPaths[0], *margin, *margin, cellW, cellH); err != nil {
			fatalf("write seed grid: %v", err)
		}
		if err := writePNG(filepath.Join(*outDir, "grid-seed.png"), seed); err != nil {
			fatalf("write seed grid: %v", err)
		}
	}
	if err := writeJSON(filepath.Join(*outDir, "batch.json"), pack); err != nil {
		fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(*outDir, "prompt.txt"), []byte(promptText(pack, *source, len(anchorPaths), *canvas)), 0o644); err != nil {
		fatalf("write prompt: %v", err)
	}
	fmt.Printf("wrote Flow grid experiment pack: %s frames=%d grid=%dx%d cell=%dx%d\n", *outDir, frameCount, cols, gridRows, cellW, cellH)
}

func selectedFrameNumbers(frameList string, frames int, frameStart int) ([]int, error) {
	if strings.TrimSpace(frameList) != "" {
		return parseFrameList(frameList)
	}
	if frames <= 0 {
		return nil, fmt.Errorf("-frames must be positive")
	}
	if frameStart < 0 {
		return nil, fmt.Errorf("-frame-start must be non-negative")
	}
	out := make([]int, 0, frames)
	for i := 0; i < frames; i++ {
		out = append(out, frameStart+i)
	}
	return out, nil
}

func parseFrameList(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	out := make([]int, 0, len(parts))
	seen := map[int]bool{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("-frame-list contains an empty frame number")
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("-frame-list contains invalid frame number %q", part)
		}
		if n < 0 {
			return nil, fmt.Errorf("-frame-list contains negative frame number %d", n)
		}
		if seen[n] {
			return nil, fmt.Errorf("-frame-list contains duplicate frame number %d", n)
		}
		seen[n] = true
		out = append(out, n)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("-frame-list must include at least one frame number")
	}
	return out, nil
}

func gridShape(frames int, columns int, rows int) (int, int) {
	if columns > 0 && rows > 0 {
		return columns, rows
	}
	switch {
	case frames <= 16:
		return 4, 4
	case frames <= 31:
		return 6, 6
	default:
		return 8, 8
	}
}

func newGuide(canvas int, margin int, cols int, rows int, cellW int, cellH int, pitchX int, pitchY int, usedCells int, detailed bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, canvas, canvas))
	fill(img, img.Bounds(), color.RGBA{G: 255, A: 255})
	border := color.RGBA{R: 10, G: 74, B: 180, A: 255}
	safe := color.RGBA{R: 255, G: 0, B: 255, A: 255}
	baseline := color.RGBA{R: 0, G: 112, B: 255, A: 255}
	inactive := color.RGBA{R: 0, G: 150, B: 0, A: 255}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell := row*cols + col
			x := margin + col*pitchX
			y := margin + row*pitchY
			rect := image.Rect(x, y, x+cellW, y+cellH)
			c := border
			if cell >= usedCells {
				c = inactive
			}
			drawOutsideRect(img, rect, c)
			if detailed && cell < usedCells {
				safeRect := rect.Inset(maxInt(8, cellW/14))
				drawRect(img, safeRect, safe)
				baseY := y + int(float64(cellH)*0.72)
				drawLine(img, x+maxInt(8, cellW/12), baseY, x+cellW-maxInt(8, cellW/12), baseY, baseline)
				drawCellTicks(img, x, y, cellW, cellH, c)
			}
		}
	}
	return img
}

func fill(img *image.RGBA, rect image.Rectangle, c color.RGBA) {
	stddraw.Draw(img, rect, &image.Uniform{C: c}, image.Point{}, stddraw.Src)
}

func drawOutsideRect(img *image.RGBA, rect image.Rectangle, c color.RGBA) {
	for i := 1; i <= 2; i++ {
		drawLine(img, rect.Min.X-i, rect.Min.Y-i, rect.Max.X+i-1, rect.Min.Y-i, c)
		drawLine(img, rect.Min.X-i, rect.Max.Y+i-1, rect.Max.X+i-1, rect.Max.Y+i-1, c)
		drawLine(img, rect.Min.X-i, rect.Min.Y-i, rect.Min.X-i, rect.Max.Y+i-1, c)
		drawLine(img, rect.Max.X+i-1, rect.Min.Y-i, rect.Max.X+i-1, rect.Max.Y+i-1, c)
	}
}

func drawRect(img *image.RGBA, rect image.Rectangle, c color.RGBA) {
	drawLine(img, rect.Min.X, rect.Min.Y, rect.Max.X-1, rect.Min.Y, c)
	drawLine(img, rect.Min.X, rect.Max.Y-1, rect.Max.X-1, rect.Max.Y-1, c)
	drawLine(img, rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y-1, c)
	drawLine(img, rect.Max.X-1, rect.Min.Y, rect.Max.X-1, rect.Max.Y-1, c)
}

func drawCellTicks(img *image.RGBA, x int, y int, w int, h int, c color.RGBA) {
	size := maxInt(6, minInt(w, h)/18)
	drawLine(img, x, y, x+size, y, c)
	drawLine(img, x, y, x, y+size, c)
	drawLine(img, x+w-1-size, y, x+w-1, y, c)
	drawLine(img, x+w-1, y, x+w-1, y+size, c)
	drawLine(img, x, y+h-1-size, x, y+h-1, c)
	drawLine(img, x, y+h-1, x+size, y+h-1, c)
	drawLine(img, x+w-1-size, y+h-1, x+w-1, y+h-1, c)
	drawLine(img, x+w-1, y+h-1-size, x+w-1, y+h-1, c)
}

func drawLine(img *image.RGBA, x1 int, y1 int, x2 int, y2 int, c color.RGBA) {
	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for y := y1; y <= y2; y++ {
			set(img, x1, y, c)
		}
		return
	}
	if y1 == y2 {
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			set(img, x, y1, c)
		}
		return
	}
}

func set(img *image.RGBA, x int, y int, c color.RGBA) {
	if image.Pt(x, y).In(img.Bounds()) {
		img.SetRGBA(x, y, c)
	}
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

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func promptText(m manifest, source string, anchorCount int, canvas int) string {
	var b strings.Builder
	if anchorCount > 0 {
		fmt.Fprintf(&b, "Upload style-anchors.png as the accepted transparent visual anchors. If the service renders transparency on black, checker, or with noisy edges, use style-anchors-neutral.png for style instead. If available, upload grid-seed.png as the base layout image and grid-guide.png as the placement guide.\n")
		fmt.Fprintf(&b, "Use the style anchor as the exact animal identity and drawing style. Use grid-seed.png only for placement, not for learning green edge color. Newly successful outputs should be reused as future anchors.\n")
	} else {
		fmt.Fprintf(&b, "Upload grid-guide.png as the placement guide. Add accepted animal PNGs as separate style anchors if possible.\n")
	}
	fmt.Fprintf(&b, "Use guide images only as placement references. Do not reproduce grid lines, labels, ticks, magenta safe boxes, blue baselines, or borders in the final artwork. Use grid-guide-detailed.png only if the minimal guide fails to keep cell placement.\n\n")
	fmt.Fprintf(&b, "Create one single %dx%d PNG named %s.\n", canvas, canvas, source)
	fmt.Fprintf(&b, "Background must be pure flat chroma green #00ff00 across all empty areas.\n")
	fmt.Fprintf(&b, "Animal: %s. One complete right-facing animal per used cell. Same camera distance, same body size, same contact baseline, same coat color, same face, same eye, same muzzle, same outline and shading style as the accepted style anchor.\n", m.Animal)
	fmt.Fprintf(&b, "Keep the animal outline clean and smooth with an antialiased silhouette at source size: no jagged cutout edge, no white fringe, no green halo, no noisy pixels around fur, no rough matte around ears, whiskers, feet, or tail, and no missing pixels, transparent pinholes, or green dots inside the animal.\n")
	fmt.Fprintf(&b, "Place one animal inside each used crop cell only. Keep ears, feet, whiskers, tail, toes, and body fully inside the crop cell. Leave empty green padding around every animal.\n")
	fmt.Fprintf(&b, "Do not draw text, numbers, grid lines, cell borders, dividers, scenery, props, floor, cast shadow, contact shadow, shelves, base, checkerboard, duplicate animals, or body parts crossing between cells.\n\n")
	fmt.Fprintf(&b, "Frame list, reading order left-to-right then top-to-bottom:\n")
	for _, frame := range m.Frames {
		fmt.Fprintf(&b, "- %s: %s\n", frame.ID, motionIntent(frame.ID))
	}
	fmt.Fprintf(&b, "\nIf any cell cannot be completed cleanly, keep the cell empty green rather than inventing a different animal style.\n")
	return b.String()
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	out := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func writeStyleAnchors(path string, anchorPaths []string) error {
	return writeStyleAnchorsWithMatte(path, anchorPaths, nil)
}

func writeStyleAnchorsMatted(path string, anchorPaths []string) error {
	green := color.RGBA{G: 255, A: 255}
	return writeStyleAnchorsWithMatte(path, anchorPaths, &green)
}

func writeStyleAnchorsNeutral(path string, anchorPaths []string) error {
	neutral := color.RGBA{R: 188, G: 192, B: 196, A: 255}
	return writeStyleAnchorsWithMatte(path, anchorPaths, &neutral)
}

func writeStyleAnchorsWithMatte(path string, anchorPaths []string, matte *color.RGBA) error {
	const cellW = 384
	const cellH = 256
	const pad = 32
	w := pad + len(anchorPaths)*(cellW+pad)
	h := cellH + pad*2
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if matte != nil {
		fill(img, img.Bounds(), *matte)
	}
	for i, anchorPath := range anchorPaths {
		anchor, err := readPNG(anchorPath)
		if err != nil {
			return err
		}
		x := pad + i*(cellW+pad)
		rect := image.Rect(x, pad, x+cellW, pad+cellH)
		drawScaledCentered(img, anchor, rect.Inset(16), matte)
	}
	return writePNG(path, img)
}

func drawAnchorInFirstCell(img *image.RGBA, anchorPath string, x int, y int, cellW int, cellH int) error {
	anchor, err := readPNG(anchorPath)
	if err != nil {
		return err
	}
	green := color.RGBA{G: 255, A: 255}
	drawScaledCentered(img, anchor, image.Rect(x, y, x+cellW, y+cellH).Inset(maxInt(8, minInt(cellW, cellH)/8)), &green)
	return nil
}

func readPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}

func validateAnchorImages(anchorPaths []string) error {
	for _, anchorPath := range anchorPaths {
		anchor, err := readPNG(anchorPath)
		if err != nil {
			return err
		}
		content := alphaBounds(anchor, anchor.Bounds())
		if content.Empty() {
			return fmt.Errorf("anchor %s has no visible alpha", anchorPath)
		}
		holeCount, holeArea, largestHole := transparentHoleStats(anchor, content)
		if holeCount > 0 {
			return fmt.Errorf("anchor %s has transparent pinholes; archive/regenerate before anchor reuse: holes=%d area=%d largest=%d", anchorPath, holeCount, holeArea, largestHole)
		}
	}
	return nil
}

func cloneRGBA(src *image.RGBA) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	stddraw.Draw(out, out.Bounds(), src, src.Bounds().Min, stddraw.Src)
	return out
}

func drawScaledCentered(dst *image.RGBA, src image.Image, rect image.Rectangle, matte *color.RGBA) {
	srcBounds := src.Bounds()
	if srcBounds.Empty() || rect.Empty() {
		return
	}
	scale := minFloat(float64(rect.Dx())/float64(srcBounds.Dx()), float64(rect.Dy())/float64(srcBounds.Dy()))
	if scale <= 0 {
		return
	}
	w := maxInt(1, int(float64(srcBounds.Dx())*scale))
	h := maxInt(1, int(float64(srcBounds.Dy())*scale))
	x0 := rect.Min.X + (rect.Dx()-w)/2
	y0 := rect.Min.Y + (rect.Dy()-h)/2
	scaled := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, srcBounds, xdraw.Over, nil)
	op := stddraw.Over
	drawSrc := image.Image(scaled)
	if matte != nil {
		matted := image.NewRGBA(scaled.Bounds())
		fill(matted, matted.Bounds(), *matte)
		stddraw.Draw(matted, matted.Bounds(), scaled, image.Point{}, stddraw.Over)
		drawSrc = matted
		op = stddraw.Src
	}
	stddraw.Draw(dst, image.Rect(x0, y0, x0+w, y0+h), drawSrc, image.Point{}, op)
}

func alphaBounds(img image.Image, rect image.Rectangle) image.Rectangle {
	rect = rect.Intersect(img.Bounds())
	minX, minY := rect.Max.X, rect.Max.Y
	maxX, maxY := rect.Min.X, rect.Min.Y
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if !alphaVisible(img, x, y) {
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

func transparentHoleStats(img image.Image, content image.Rectangle) (holeCount int, holeArea int, largestHole int) {
	width := content.Dx()
	height := content.Dy()
	if width <= 0 || height <= 0 {
		return 0, 0, 0
	}
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

func motionIntent(id string) string {
	n := -1
	fmt.Sscanf(id, "set00-frame-%02d", &n)
	if n < 0 {
		fmt.Sscanf(id, "set01-frame-%02d", &n)
	}
	switch {
	case n >= 0 && n <= 3:
		return "idle breathing / tiny ear or weight shift"
	case n >= 4 && n <= 11:
		return "walk cycle, small grounded step"
	case n >= 12 && n <= 19:
		return "scurry burst, faster low movement"
	case n >= 20 && n <= 25:
		return "forage, sniff, or nibble"
	case n >= 26 && n <= 31:
		return "species-safe action, no wheel or human pose"
	case n >= 32 && n <= 39:
		return "turn motion"
	case n >= 40 && n <= 43:
		return "eat"
	case n >= 44 && n <= 47:
		return "ground check, paw or foot action"
	case n >= 48 && n <= 51:
		return "rest-alert hold, not seated on a platform"
	case n >= 52 && n <= 55:
		return "groom"
	case n >= 56 && n <= 61:
		return "reaction"
	default:
		return "small continuation motion"
	}
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func minFloat(a float64, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
