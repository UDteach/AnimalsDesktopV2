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
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	xdraw "golang.org/x/image/draw"
)

const (
	sheetWidth                          = 1536
	sheetHeight                         = 1024
	cellWidth                           = 768
	cellHeight                          = 512
	cellCount                           = 4
	minRawSubjectIoU                    = 0.985
	minAlbinoHighContrastRawSubjectIoU  = 0.980
	maxRawCentroidShiftPixels           = 1.25
	defaultRawGeometryPolicyReportID    = "strict_default"
	albinoHighContrastRawGeometryPolicy = "ferret_albino_high_contrast_v1"
	defaultTonePolicyReportID           = "strict_default"
	sableExactRecoveryTonePolicy        = "ferret_sable_exact_recovery_v1"
	sableExactRecoveryApprovedInputSHA  = "e16d004799e033405d1404ca318b18a0550d7a36db2d1bdefb99b7fd97c6fbf5"
	sableExactRecoveryTargetRatio       = 0.573025
	defaultToneGainMin                  = 0.95
	sableExactRecoveryToneGainMin       = 0.85
	defaultToneGainMax                  = 1.05
)

var safeIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// coatManifest is intentionally small. Sources are relative to the manifest,
// and outputs are relative to the slice output directory.
type coatManifest struct {
	Species           string                  `json:"species"`
	Call              string                  `json:"call"`
	Prompt            string                  `json:"prompt"`
	Swatch            string                  `json:"swatch"`
	RawGeometryPolicy string                  `json:"raw_geometry_policy,omitempty"`
	TonePolicy        *coatManifestTonePolicy `json:"tone_policy,omitempty"`
	Cells             []coatManifestCell      `json:"cells"`

	manifestSHA256 string
}

type coatManifestTonePolicy struct {
	ID                  string  `json:"id"`
	ApprovedInputSHA256 string  `json:"approved_input_sha256"`
	TargetRatio         float64 `json:"target_ratio"`
}

type coatManifestCell struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Source    string `json:"source"`
	BaseFrame string `json:"base_frame"`
	Cell      int    `json:"cell"`
	Output    string `json:"output,omitempty"`
}

type cliOptions struct {
	Mode        string
	Manifest    string
	Sheet       string
	BaseSheet   string
	Out         string
	Report      string
	TargetRatio float64
}

type rectReport struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type artifactReport struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type provenanceArtifactReport struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type rawGeometryPolicyReport struct {
	ID                          string  `json:"id"`
	EffectiveMinIoU             float64 `json:"effective_min_iou"`
	EffectiveMaxCentroidShiftPX float64 `json:"effective_max_centroid_shift_px"`
}

type tonePolicyReport struct {
	ID                  string   `json:"id"`
	EffectiveGainMin    float64  `json:"effective_gain_min"`
	EffectiveGainMax    float64  `json:"effective_gain_max"`
	ApprovedInputSHA256 string   `json:"approved_input_sha256,omitempty"`
	ActualInputSHA256   string   `json:"actual_input_sha256,omitempty"`
	TargetRatio         *float64 `json:"target_ratio,omitempty"`
}

type cellReport struct {
	ID                      string      `json:"id"`
	Role                    string      `json:"role"`
	Cell                    int         `json:"cell"`
	Bounds                  rectReport  `json:"bounds"`
	BaseFrame               string      `json:"base_frame"`
	BaseFrameSHA256         string      `json:"base_frame_sha256"`
	Source                  string      `json:"source,omitempty"`
	SourceSHA256            string      `json:"source_sha256,omitempty"`
	SourceWidth             int         `json:"source_width,omitempty"`
	SourceHeight            int         `json:"source_height,omitempty"`
	BuildFilter             string      `json:"build_filter,omitempty"`
	Output                  string      `json:"output,omitempty"`
	OutputSHA256            string      `json:"output_sha256,omitempty"`
	BaseSubjectBounds       *rectReport `json:"base_subject_bounds,omitempty"`
	TorsoROI                *rectReport `json:"torso_roi,omitempty"`
	ReferenceAreaPX         *int        `json:"reference_area_px,omitempty"`
	CandidateAreaPX         *int        `json:"candidate_area_px,omitempty"`
	IntersectionPX          *int        `json:"intersection_px,omitempty"`
	UnionPX                 *int        `json:"union_px,omitempty"`
	IoU                     *float64    `json:"iou,omitempty"`
	CentroidDXPX            *float64    `json:"centroid_dx_px,omitempty"`
	CentroidDYPX            *float64    `json:"centroid_dy_px,omitempty"`
	CentroidShiftPX         *float64    `json:"centroid_shift_px,omitempty"`
	BaseMedianLuma          *float64    `json:"base_median_luma,omitempty"`
	CandidateMedianLuma     *float64    `json:"candidate_median_luma,omitempty"`
	Ratio                   *float64    `json:"ratio,omitempty"`
	PostCandidateMedianLuma *float64    `json:"post_candidate_median_luma,omitempty"`
	PostRatio               *float64    `json:"post_ratio,omitempty"`
}

type coatReport struct {
	Version            int                      `json:"version"`
	Mode               string                   `json:"mode"`
	Manifest           string                   `json:"manifest"`
	ManifestSHA256     string                   `json:"manifest_sha256"`
	Species            string                   `json:"species"`
	Call               string                   `json:"call"`
	RawGeometryPolicy  rawGeometryPolicyReport  `json:"raw_geometry_policy"`
	TonePolicy         tonePolicyReport         `json:"tone_policy"`
	Prompt             provenanceArtifactReport `json:"prompt"`
	Swatch             artifactReport           `json:"swatch"`
	BaseSheet          *artifactReport          `json:"base_sheet,omitempty"`
	InputSheet         *artifactReport          `json:"input_sheet,omitempty"`
	OutputSheet        *artifactReport          `json:"output_sheet,omitempty"`
	OutputDirectory    string                   `json:"output_directory,omitempty"`
	OutputSHA256       string                   `json:"output_sha256,omitempty"`
	PreSHA256          string                   `json:"pre_sha256,omitempty"`
	PostSHA256         string                   `json:"post_sha256,omitempty"`
	TargetRatio        *float64                 `json:"target_ratio,omitempty"`
	Gain               *float64                 `json:"gain,omitempty"`
	GroupMeanRatio     *float64                 `json:"group_mean_ratio,omitempty"`
	PostGroupMeanRatio *float64                 `json:"post_group_mean_ratio,omitempty"`
	Cells              []cellReport             `json:"cells"`
}

type decodedPNG struct {
	Image    image.Image
	Artifact artifactReport
}

type encodedOutput struct {
	Path string
	Data []byte
}

type manifestProvenance struct {
	ManifestSHA256 string
	Prompt         provenanceArtifactReport
	Swatch         artifactReport
	BaseFrames     map[int]provenanceArtifactReport
}

type rawGeometryPolicy struct {
	ID                 string
	MinIoU             float64
	MaxCentroidShiftPX float64
}

type tonePolicy struct {
	ID                  string
	GainMin             float64
	GainMax             float64
	ApprovedInputSHA256 string
	TargetRatio         *float64
}

type transactionFS struct {
	CreateTemp func(string, string) (*os.File, error)
	WriteTemp  func(*os.File, []byte) error
	Rename     func(string, string) error
	Remove     func(string) error
}

var txFS = transactionFS{
	CreateTemp: os.CreateTemp,
	WriteTemp: func(file *os.File, data []byte) error {
		_, err := file.Write(data)
		return err
	},
	Rename: os.Rename,
	Remove: os.Remove,
}

func main() {
	if err := runCLI(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "coatbatch: %v\n", err)
		os.Exit(1)
	}
}

func runCLI(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("coatbatch", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	opts := cliOptions{TargetRatio: math.NaN()}
	fs.StringVar(&opts.Mode, "mode", "", "mode: build, measure, tone, or slice")
	fs.StringVar(&opts.Manifest, "manifest", "", "coat-batch manifest JSON")
	fs.StringVar(&opts.Sheet, "sheet", "", "generated or calibrated 1536x1024 PNG")
	fs.StringVar(&opts.BaseSheet, "base-sheet", "", "base 1536x1024 PNG")
	fs.StringVar(&opts.Out, "out", "", "output PNG for build/tone or directory for slice")
	fs.StringVar(&opts.Report, "report", "", "JSON report path")
	fs.Float64Var(&opts.TargetRatio, "target-ratio", math.NaN(), "required target group ratio for tone mode")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(fs.Args(), " "))
	}
	if err := execute(opts); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "coatbatch %s complete\n", strings.ToLower(strings.TrimSpace(opts.Mode)))
	return nil
}

func execute(opts cliOptions) error {
	opts.Mode = strings.ToLower(strings.TrimSpace(opts.Mode))
	if opts.Mode == "" {
		return errors.New("-mode is required")
	}
	if opts.Manifest == "" {
		return errors.New("-manifest is required")
	}
	if opts.Report == "" {
		return errors.New("-report is required")
	}

	m, err := loadManifest(opts.Manifest)
	if err != nil {
		return err
	}
	provenance, err := loadManifestProvenance(opts.Manifest, m)
	if err != nil {
		return err
	}
	for _, input := range provenanceInputs(opts.Manifest, provenance) {
		if err := requireDistinctPaths(input, namedPath{"report", opts.Report}); err != nil {
			return err
		}
		if opts.Out != "" && opts.Mode != "slice" {
			if err := requireDistinctPaths(input, namedPath{"output", opts.Out}); err != nil {
				return err
			}
		}
	}

	switch opts.Mode {
	case "build":
		if opts.Out == "" {
			return errors.New("-out is required for build")
		}
		if opts.Sheet != "" || opts.BaseSheet != "" || !math.IsNaN(opts.TargetRatio) {
			return errors.New("build accepts -manifest, -out, and -report only")
		}
		return runBuild(opts, m, provenance)
	case "measure":
		if opts.Sheet == "" || opts.BaseSheet == "" {
			return errors.New("-sheet and -base-sheet are required for measure")
		}
		if opts.Out != "" || !math.IsNaN(opts.TargetRatio) {
			return errors.New("measure does not accept -out or -target-ratio")
		}
		return runMeasure(opts, m, provenance)
	case "tone":
		if opts.Sheet == "" || opts.BaseSheet == "" || opts.Out == "" {
			return errors.New("-sheet, -base-sheet, and -out are required for tone")
		}
		if math.IsNaN(opts.TargetRatio) {
			return errors.New("-target-ratio is required for tone")
		}
		if math.IsInf(opts.TargetRatio, 0) || opts.TargetRatio <= 0 {
			return errors.New("-target-ratio must be a finite positive number")
		}
		return runTone(opts, m, provenance)
	case "slice":
		if opts.Sheet == "" || opts.Out == "" {
			return errors.New("-sheet and -out are required for slice")
		}
		if opts.BaseSheet != "" || !math.IsNaN(opts.TargetRatio) {
			return errors.New("slice does not accept -base-sheet or -target-ratio")
		}
		return runSlice(opts, m, provenance)
	default:
		return fmt.Errorf("unsupported -mode %q", opts.Mode)
	}
}

func loadManifest(path string) (coatManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return coatManifest{}, fmt.Errorf("read manifest: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var m coatManifest
	if err := decoder.Decode(&m); err != nil {
		return coatManifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return coatManifest{}, errors.New("decode manifest: trailing JSON value")
		}
		return coatManifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	m.Species = strings.TrimSpace(m.Species)
	if m.Species == "" {
		return coatManifest{}, errors.New("manifest species is required")
	}
	if err := validateManifestPolicies(m); err != nil {
		return coatManifest{}, err
	}
	if !safeIDPattern.MatchString(m.Call) || m.Call == "." || m.Call == ".." || unsafeWindowsPathComponent(m.Call) {
		return coatManifest{}, fmt.Errorf("manifest call has unsafe id %q", m.Call)
	}
	prompt, err := validateRelativeFile(m.Prompt, "prompt", ".txt")
	if err != nil {
		return coatManifest{}, err
	}
	m.Prompt = prompt
	swatch, err := validateRelativeFile(m.Swatch, "swatch", ".png")
	if err != nil {
		return coatManifest{}, err
	}
	m.Swatch = swatch
	if len(m.Cells) != cellCount {
		return coatManifest{}, fmt.Errorf("manifest must declare exactly %d cells, got %d", cellCount, len(m.Cells))
	}

	seenIDs := make(map[string]bool, cellCount)
	seenCells := make(map[int]bool, cellCount)
	seenOutputs := make(map[string]bool, cellCount)
	targetCount := 0
	for i := range m.Cells {
		cell := &m.Cells[i]
		if !safeIDPattern.MatchString(cell.ID) || cell.ID == "." || cell.ID == ".." || unsafeWindowsPathComponent(cell.ID) {
			return coatManifest{}, fmt.Errorf("cell %d has unsafe id %q", i, cell.ID)
		}
		idKey := strings.ToLower(cell.ID)
		if seenIDs[idKey] {
			return coatManifest{}, fmt.Errorf("duplicate cell id %q", cell.ID)
		}
		seenIDs[idKey] = true
		if cell.Cell < 0 || cell.Cell >= cellCount {
			return coatManifest{}, fmt.Errorf("cell %q index must be 0..3, got %d", cell.ID, cell.Cell)
		}
		if seenCells[cell.Cell] {
			return coatManifest{}, fmt.Errorf("duplicate cell index %d", cell.Cell)
		}
		seenCells[cell.Cell] = true
		switch cell.Role {
		case "target":
			targetCount++
		case "filler":
			if cell.Output != "" {
				return coatManifest{}, fmt.Errorf("cell %q filler output must be absent", cell.ID)
			}
		default:
			return coatManifest{}, fmt.Errorf("cell %q role must be target or filler, got %q", cell.ID, cell.Role)
		}

		source, err := validateRelativePNG(cell.Source, "source")
		if err != nil {
			return coatManifest{}, fmt.Errorf("cell %q: %w", cell.ID, err)
		}
		cell.Source = source
		baseFrame, err := validateRelativePNG(cell.BaseFrame, "base_frame")
		if err != nil {
			return coatManifest{}, fmt.Errorf("cell %q: %w", cell.ID, err)
		}
		cell.BaseFrame = baseFrame
		if cell.Role == "target" {
			if cell.Output == "" {
				cell.Output = cell.ID + ".png"
			}
			output, err := validateRelativePNG(cell.Output, "output")
			if err != nil {
				return coatManifest{}, fmt.Errorf("cell %q: %w", cell.ID, err)
			}
			cell.Output = output
			outputKey := strings.ToLower(filepath.Clean(output))
			if seenOutputs[outputKey] {
				return coatManifest{}, fmt.Errorf("duplicate output %q", output)
			}
			seenOutputs[outputKey] = true
		}
	}
	if targetCount == 0 {
		return coatManifest{}, errors.New("manifest must declare at least one target cell")
	}
	for cell := 0; cell < cellCount; cell++ {
		if !seenCells[cell] {
			return coatManifest{}, fmt.Errorf("manifest is missing cell index %d", cell)
		}
	}
	sort.Slice(m.Cells, func(i, j int) bool { return m.Cells[i].Cell < m.Cells[j].Cell })
	m.manifestSHA256 = sha256Hex(data)
	return m, nil
}

func validateManifestPolicies(m coatManifest) error {
	switch m.RawGeometryPolicy {
	case "":
	case albinoHighContrastRawGeometryPolicy:
		if m.Species != "ferret_albino" {
			return fmt.Errorf(
				"raw_geometry_policy %q is valid only for species %q, got %q",
				m.RawGeometryPolicy,
				"ferret_albino",
				m.Species,
			)
		}
	default:
		return fmt.Errorf("unknown raw_geometry_policy %q", m.RawGeometryPolicy)
	}

	if m.TonePolicy == nil {
		return nil
	}
	switch m.TonePolicy.ID {
	case sableExactRecoveryTonePolicy:
		if m.Species != "ferret_sable" {
			return fmt.Errorf(
				"tone_policy %q is valid only for species %q, got %q",
				m.TonePolicy.ID,
				"ferret_sable",
				m.Species,
			)
		}
		if m.TonePolicy.ApprovedInputSHA256 != sableExactRecoveryApprovedInputSHA {
			return fmt.Errorf(
				"tone_policy %q approved_input_sha256 must be exactly %q",
				m.TonePolicy.ID,
				sableExactRecoveryApprovedInputSHA,
			)
		}
		if m.TonePolicy.TargetRatio != sableExactRecoveryTargetRatio {
			return fmt.Errorf(
				"tone_policy %q target_ratio must be exactly %.6f",
				m.TonePolicy.ID,
				sableExactRecoveryTargetRatio,
			)
		}
	default:
		return fmt.Errorf("unknown tone_policy id %q", m.TonePolicy.ID)
	}
	return nil
}

func validateRelativePNG(raw string, label string) (string, error) {
	return validateRelativeFile(raw, label, ".png")
}

func validateRelativeFile(raw string, label string, extension string) (string, error) {
	if raw == "" || strings.TrimSpace(raw) != raw {
		return "", fmt.Errorf("%s path is empty or has surrounding whitespace", label)
	}
	if filepath.IsAbs(raw) || filepath.VolumeName(raw) != "" {
		return "", fmt.Errorf("%s path must be relative: %q", label, raw)
	}
	clean := filepath.Clean(raw)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%s path escapes its root: %q", label, raw)
	}
	if strings.ContainsRune(clean, '\x00') {
		return "", fmt.Errorf("%s path contains NUL", label)
	}
	for _, component := range strings.FieldsFunc(strings.ReplaceAll(clean, `\`, `/`), func(r rune) bool { return r == '/' }) {
		if unsafeWindowsPathComponent(component) {
			return "", fmt.Errorf("%s path has an unsafe component %q", label, component)
		}
	}
	if !strings.EqualFold(filepath.Ext(clean), extension) {
		return "", fmt.Errorf("%s path must end in %s: %q", label, extension, raw)
	}
	return clean, nil
}

func loadManifestProvenance(manifestPath string, m coatManifest) (manifestProvenance, error) {
	manifestSHA256 := m.manifestSHA256
	if manifestSHA256 == "" {
		manifestData, err := os.ReadFile(manifestPath)
		if err != nil {
			return manifestProvenance{}, fmt.Errorf("read manifest provenance: %w", err)
		}
		manifestSHA256 = sha256Hex(manifestData)
	}
	manifestDir := filepath.Dir(manifestPath)
	promptPath, err := resolveConfinedExisting(manifestDir, m.Prompt)
	if err != nil {
		return manifestProvenance{}, fmt.Errorf("prompt: %w", err)
	}
	promptData, err := os.ReadFile(promptPath)
	if err != nil {
		return manifestProvenance{}, fmt.Errorf("read prompt: %w", err)
	}
	if len(promptData) == 0 || !utf8.Valid(promptData) {
		return manifestProvenance{}, errors.New("prompt must be a non-empty UTF-8 text file")
	}
	swatchPath, err := resolveConfinedExisting(manifestDir, m.Swatch)
	if err != nil {
		return manifestProvenance{}, fmt.Errorf("swatch: %w", err)
	}
	swatch, err := readPNG(swatchPath)
	if err != nil {
		return manifestProvenance{}, fmt.Errorf("swatch: %w", err)
	}
	provenance := manifestProvenance{
		ManifestSHA256: manifestSHA256,
		Prompt: provenanceArtifactReport{
			Path:   filepath.ToSlash(promptPath),
			SHA256: sha256Hex(promptData),
		},
		Swatch:     swatch.Artifact,
		BaseFrames: make(map[int]provenanceArtifactReport, cellCount),
	}
	for _, cell := range m.Cells {
		baseFramePath, err := resolveConfinedExisting(manifestDir, cell.BaseFrame)
		if err != nil {
			return manifestProvenance{}, fmt.Errorf("cell %q base_frame: %w", cell.ID, err)
		}
		baseFrame, err := readPNG(baseFramePath)
		if err != nil {
			return manifestProvenance{}, fmt.Errorf("cell %q base_frame: %w", cell.ID, err)
		}
		provenance.BaseFrames[cell.Cell] = provenanceArtifactReport{
			Path:   filepath.ToSlash(baseFramePath),
			SHA256: baseFrame.Artifact.SHA256,
		}
	}
	return provenance, nil
}

func effectiveRawGeometryPolicy(m coatManifest) rawGeometryPolicy {
	if m.RawGeometryPolicy == albinoHighContrastRawGeometryPolicy {
		return rawGeometryPolicy{
			ID:                 albinoHighContrastRawGeometryPolicy,
			MinIoU:             minAlbinoHighContrastRawSubjectIoU,
			MaxCentroidShiftPX: maxRawCentroidShiftPixels,
		}
	}
	return rawGeometryPolicy{
		ID:                 defaultRawGeometryPolicyReportID,
		MinIoU:             minRawSubjectIoU,
		MaxCentroidShiftPX: maxRawCentroidShiftPixels,
	}
}

func effectiveTonePolicy(m coatManifest) tonePolicy {
	if m.TonePolicy != nil && m.TonePolicy.ID == sableExactRecoveryTonePolicy {
		targetRatio := sableExactRecoveryTargetRatio
		return tonePolicy{
			ID:                  sableExactRecoveryTonePolicy,
			GainMin:             sableExactRecoveryToneGainMin,
			GainMax:             defaultToneGainMax,
			ApprovedInputSHA256: sableExactRecoveryApprovedInputSHA,
			TargetRatio:         &targetRatio,
		}
	}
	return tonePolicy{
		ID:      defaultTonePolicyReportID,
		GainMin: defaultToneGainMin,
		GainMax: defaultToneGainMax,
	}
}

func newReport(mode string, manifestPath string, m coatManifest, provenance manifestProvenance) coatReport {
	rawPolicy := effectiveRawGeometryPolicy(m)
	selectedTonePolicy := effectiveTonePolicy(m)
	return coatReport{
		Version:        2,
		Mode:           mode,
		Manifest:       filepath.ToSlash(manifestPath),
		ManifestSHA256: provenance.ManifestSHA256,
		Species:        m.Species,
		Call:           m.Call,
		RawGeometryPolicy: rawGeometryPolicyReport{
			ID:                          rawPolicy.ID,
			EffectiveMinIoU:             rawPolicy.MinIoU,
			EffectiveMaxCentroidShiftPX: rawPolicy.MaxCentroidShiftPX,
		},
		TonePolicy: tonePolicyReport{
			ID:                  selectedTonePolicy.ID,
			EffectiveGainMin:    selectedTonePolicy.GainMin,
			EffectiveGainMax:    selectedTonePolicy.GainMax,
			ApprovedInputSHA256: selectedTonePolicy.ApprovedInputSHA256,
			TargetRatio:         selectedTonePolicy.TargetRatio,
		},
		Prompt: provenance.Prompt,
		Swatch: provenance.Swatch,
	}
}

func setReportActualInputSHA256(report *coatReport, m coatManifest, sha256 string) {
	if m.TonePolicy != nil {
		report.TonePolicy.ActualInputSHA256 = sha256
	}
}

func validateTonePolicyBinding(m coatManifest, actualInputSHA256 string, targetRatio float64) error {
	if m.TonePolicy == nil {
		return nil
	}
	policy := effectiveTonePolicy(m)
	if policy.TargetRatio == nil {
		return fmt.Errorf("tone_policy %q has no approved target ratio", policy.ID)
	}
	if actualInputSHA256 != policy.ApprovedInputSHA256 {
		return fmt.Errorf(
			"tone_policy %q input sheet SHA-256 %q does not match approved_input_sha256 %q",
			policy.ID,
			actualInputSHA256,
			policy.ApprovedInputSHA256,
		)
	}
	if targetRatio != *policy.TargetRatio {
		return fmt.Errorf(
			"tone_policy %q requires -target-ratio exactly %.6f, got %.17g",
			policy.ID,
			*policy.TargetRatio,
			targetRatio,
		)
	}
	return nil
}

func validateToneGain(gain float64, policy tonePolicy) error {
	if gain < policy.GainMin || gain > policy.GainMax || math.IsNaN(gain) || math.IsInf(gain, 0) {
		return fmt.Errorf(
			"required gain %.6f is outside allowed range %.2f..%.2f for tone policy %q",
			gain,
			policy.GainMin,
			policy.GainMax,
			policy.ID,
		)
	}
	return nil
}

func newCellReport(cell coatManifestCell, provenance manifestProvenance) cellReport {
	baseFrame := provenance.BaseFrames[cell.Cell]
	return cellReport{
		ID:              cell.ID,
		Role:            cell.Role,
		Cell:            cell.Cell,
		Bounds:          reportRect(boundsForCell(cell.Cell)),
		BaseFrame:       baseFrame.Path,
		BaseFrameSHA256: baseFrame.SHA256,
	}
}

func provenanceInputs(manifestPath string, provenance manifestProvenance) []namedPath {
	inputs := []namedPath{
		{Name: "manifest", Path: manifestPath},
		{Name: "prompt", Path: provenance.Prompt.Path},
		{Name: "swatch", Path: provenance.Swatch.Path},
	}
	for cell := 0; cell < cellCount; cell++ {
		inputs = append(inputs, namedPath{
			Name: fmt.Sprintf("base frame %d", cell),
			Path: provenance.BaseFrames[cell].Path,
		})
	}
	return inputs
}

func runBuild(opts cliOptions, m coatManifest, provenance manifestProvenance) error {
	if err := requirePNGOutput(opts.Out); err != nil {
		return err
	}
	if err := requireDistinctPaths(
		namedPath{"manifest", opts.Manifest},
		namedPath{"output", opts.Out},
		namedPath{"report", opts.Report},
	); err != nil {
		return err
	}

	sheet := image.NewNRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))
	report := newReport("build", opts.Manifest, m, provenance)
	report.Cells = make([]cellReport, 0, cellCount)
	manifestDir := filepath.Dir(opts.Manifest)
	for _, cell := range m.Cells {
		sourcePath, err := resolveConfinedExisting(manifestDir, cell.Source)
		if err != nil {
			return fmt.Errorf("cell %q source: %w", cell.ID, err)
		}
		if err := requireDistinctPaths(
			namedPath{"source", sourcePath},
			namedPath{"output", opts.Out},
			namedPath{"report", opts.Report},
		); err != nil {
			return fmt.Errorf("cell %q: %w", cell.ID, err)
		}
		source, err := readPNG(sourcePath)
		if err != nil {
			return fmt.Errorf("cell %q: %w", cell.ID, err)
		}
		if source.Artifact.Width <= 0 || source.Artifact.Height <= 0 {
			return fmt.Errorf("cell %q source has non-positive dimensions", cell.ID)
		}
		bounds := boundsForCell(cell.Cell)
		buildFilter := "none"
		if source.Artifact.Width == cellWidth && source.Artifact.Height == cellHeight {
			stddraw.Draw(sheet, bounds, source.Image, source.Image.Bounds().Min, stddraw.Src)
		} else {
			// CatmullRom is this Go implementation's resampling filter. The report
			// intentionally does not claim byte or pixel equivalence with Pillow.
			xdraw.CatmullRom.Scale(sheet, bounds, source.Image, source.Image.Bounds(), stddraw.Src, nil)
			buildFilter = "golang.org/x/image/draw.CatmullRom"
		}
		cellReport := newCellReport(cell, provenance)
		cellReport.Source = filepath.ToSlash(sourcePath)
		cellReport.SourceSHA256 = source.Artifact.SHA256
		cellReport.SourceWidth = source.Artifact.Width
		cellReport.SourceHeight = source.Artifact.Height
		cellReport.BuildFilter = buildFilter
		report.Cells = append(report.Cells, cellReport)
	}
	pngData, err := encodePNG(sheet)
	if err != nil {
		return fmt.Errorf("encode output: %w", err)
	}
	outputArtifact := artifactFromData(opts.Out, pngData, sheetWidth, sheetHeight)
	report.OutputSheet = &outputArtifact
	report.OutputSHA256 = outputArtifact.SHA256
	if err := writeArtifactAndReport(opts.Out, pngData, opts.Report, report); err != nil {
		return err
	}
	return nil
}

func runMeasure(opts cliOptions, m coatManifest, provenance manifestProvenance) error {
	for _, input := range []namedPath{
		{"manifest", opts.Manifest},
		{"base sheet", opts.BaseSheet},
		{"sheet", opts.Sheet},
	} {
		if err := requireDistinctPaths(input, namedPath{"report", opts.Report}); err != nil {
			return err
		}
	}
	base, err := readExactSheet(opts.BaseSheet)
	if err != nil {
		return fmt.Errorf("base sheet: %w", err)
	}
	candidate, err := readExactSheet(opts.Sheet)
	if err != nil {
		return fmt.Errorf("sheet: %w", err)
	}
	cells, mean, err := measureSheetPair(base.Image, candidate.Image, m, provenance)
	if err != nil {
		return err
	}
	report := newReport("measure", opts.Manifest, m, provenance)
	report.BaseSheet = &base.Artifact
	report.InputSheet = &candidate.Artifact
	setReportActualInputSHA256(&report, m, candidate.Artifact.SHA256)
	report.GroupMeanRatio = floatPointer(mean)
	report.Cells = cells
	geometryErr := validateRawSubjectGeometryWithPolicy(cells, effectiveRawGeometryPolicy(m))
	if err := writeReport(opts.Report, report); err != nil {
		if geometryErr != nil {
			return fmt.Errorf("%v; write geometry diagnostic report: %w", geometryErr, err)
		}
		return err
	}
	return geometryErr
}

func runTone(opts cliOptions, m coatManifest, provenance manifestProvenance) error {
	if err := requirePNGOutput(opts.Out); err != nil {
		return err
	}
	inputs := []namedPath{
		{"manifest", opts.Manifest},
		{"base sheet", opts.BaseSheet},
		{"sheet", opts.Sheet},
	}
	if err := requireDistinctPaths(namedPath{"output", opts.Out}, namedPath{"report", opts.Report}); err != nil {
		return err
	}
	for _, input := range inputs {
		if err := requireDistinctPaths(input, namedPath{"output", opts.Out}); err != nil {
			return err
		}
		if err := requireDistinctPaths(input, namedPath{"report", opts.Report}); err != nil {
			return err
		}
	}
	base, err := readExactSheet(opts.BaseSheet)
	if err != nil {
		return fmt.Errorf("base sheet: %w", err)
	}
	candidate, err := readExactSheet(opts.Sheet)
	if err != nil {
		return fmt.Errorf("sheet: %w", err)
	}
	preCells, preMean, err := measureSheetPair(base.Image, candidate.Image, m, provenance)
	if err != nil {
		return err
	}
	if geometryErr := validateRawSubjectGeometryWithPolicy(preCells, effectiveRawGeometryPolicy(m)); geometryErr != nil {
		report := newReport("tone", opts.Manifest, m, provenance)
		report.BaseSheet = &base.Artifact
		report.InputSheet = &candidate.Artifact
		setReportActualInputSHA256(&report, m, candidate.Artifact.SHA256)
		report.PreSHA256 = candidate.Artifact.SHA256
		report.TargetRatio = floatPointer(opts.TargetRatio)
		report.GroupMeanRatio = floatPointer(preMean)
		report.Cells = preCells
		if preMean > 0 && !math.IsNaN(preMean) && !math.IsInf(preMean, 0) {
			report.Gain = floatPointer(opts.TargetRatio / preMean)
		}
		if reportErr := writeReport(opts.Report, report); reportErr != nil {
			return fmt.Errorf("%v; write geometry diagnostic report: %w", geometryErr, reportErr)
		}
		return geometryErr
	}
	if preMean <= 0 || math.IsNaN(preMean) || math.IsInf(preMean, 0) {
		return fmt.Errorf("invalid pre-calibration group mean ratio %.6f", preMean)
	}
	if err := validateTonePolicyBinding(m, candidate.Artifact.SHA256, opts.TargetRatio); err != nil {
		return err
	}
	gain := opts.TargetRatio / preMean
	// This rejection happens before output or report creation.
	if err := validateToneGain(gain, effectiveTonePolicy(m)); err != nil {
		return err
	}

	calibrated := applyUniformSubjectGain(candidate.Image, gain)
	postCells, postMean, err := measureSheetPair(base.Image, calibrated, m, provenance)
	if err != nil {
		return fmt.Errorf("measure calibrated sheet: %w", err)
	}
	for i := range preCells {
		preCells[i].PostCandidateMedianLuma = postCells[i].CandidateMedianLuma
		preCells[i].PostRatio = postCells[i].Ratio
	}
	pngData, err := encodePNG(calibrated)
	if err != nil {
		return fmt.Errorf("encode calibrated sheet: %w", err)
	}
	outputArtifact := artifactFromData(opts.Out, pngData, sheetWidth, sheetHeight)
	report := newReport("tone", opts.Manifest, m, provenance)
	report.BaseSheet = &base.Artifact
	report.InputSheet = &candidate.Artifact
	setReportActualInputSHA256(&report, m, candidate.Artifact.SHA256)
	report.OutputSheet = &outputArtifact
	report.OutputSHA256 = outputArtifact.SHA256
	report.PreSHA256 = candidate.Artifact.SHA256
	report.PostSHA256 = outputArtifact.SHA256
	report.TargetRatio = floatPointer(opts.TargetRatio)
	report.Gain = floatPointer(gain)
	report.GroupMeanRatio = floatPointer(preMean)
	report.PostGroupMeanRatio = floatPointer(postMean)
	report.Cells = preCells
	return writeArtifactAndReport(opts.Out, pngData, opts.Report, report)
}

func runSlice(opts cliOptions, m coatManifest, provenance manifestProvenance) error {
	if err := requireDistinctPaths(
		namedPath{"manifest", opts.Manifest},
		namedPath{"sheet", opts.Sheet},
		namedPath{"report", opts.Report},
	); err != nil {
		return err
	}
	sheet, err := readExactSheet(opts.Sheet)
	if err != nil {
		return fmt.Errorf("sheet: %w", err)
	}

	report := newReport("slice", opts.Manifest, m, provenance)
	report.InputSheet = &sheet.Artifact
	report.OutputDirectory = filepath.ToSlash(opts.Out)
	report.Cells = make([]cellReport, 0, cellCount)
	outputs := make([]encodedOutput, 0, cellCount)
	for _, cell := range m.Cells {
		bounds := boundsForCell(cell.Cell)
		if _, err := checkedSubjectBounds(sheet.Image, bounds, "sheet", cell.ID); err != nil {
			return err
		}
		cellReport := newCellReport(cell, provenance)
		if cell.Role == "filler" {
			report.Cells = append(report.Cells, cellReport)
			continue
		}
		target, err := resolveConfinedOutput(opts.Out, cell.Output)
		if err != nil {
			return fmt.Errorf("cell %q output: %w", cell.ID, err)
		}
		if err := requireDistinctPaths(
			namedPath{"sheet", opts.Sheet},
			namedPath{"cell output", target},
			namedPath{"report", opts.Report},
		); err != nil {
			return fmt.Errorf("cell %q: %w", cell.ID, err)
		}
		for _, input := range provenanceInputs(opts.Manifest, provenance) {
			if err := requireDistinctPaths(input, namedPath{"cell output", target}); err != nil {
				return fmt.Errorf("cell %q: %w", cell.ID, err)
			}
		}
		cropped := image.NewNRGBA(image.Rect(0, 0, cellWidth, cellHeight))
		stddraw.Draw(cropped, cropped.Bounds(), sheet.Image, bounds.Min, stddraw.Src)
		data, err := encodePNG(cropped)
		if err != nil {
			return fmt.Errorf("cell %q encode: %w", cell.ID, err)
		}
		hash := sha256Hex(data)
		outputs = append(outputs, encodedOutput{Path: target, Data: data})
		cellReport.Output = filepath.ToSlash(target)
		cellReport.OutputSHA256 = hash
		report.Cells = append(report.Cells, cellReport)
	}
	for i := range outputs {
		for j := i + 1; j < len(outputs); j++ {
			if samePath(outputs[i].Path, outputs[j].Path) {
				return fmt.Errorf("slice outputs resolve to the same path: %s", outputs[i].Path)
			}
		}
	}

	reportData, err := encodeReport(report)
	if err != nil {
		return err
	}
	outputs = append(outputs, encodedOutput{Path: opts.Report, Data: reportData})
	return writeTransaction(outputs)
}

func measureSheetPair(base image.Image, candidate image.Image, m coatManifest, provenance manifestProvenance) ([]cellReport, float64, error) {
	reports := make([]cellReport, 0, cellCount)
	sum := 0.0
	for _, cell := range m.Cells {
		cellBounds := boundsForCell(cell.Cell)
		baseBounds, err := checkedSubjectBounds(base, cellBounds, "base sheet", cell.ID)
		if err != nil {
			return nil, 0, err
		}
		if _, err := checkedSubjectBounds(candidate, cellBounds, "candidate sheet", cell.ID); err != nil {
			return nil, 0, err
		}
		geometry, err := compareSubjectGeometry(base, candidate, cellBounds)
		if err != nil {
			return nil, 0, fmt.Errorf("cell %q raw geometry: %w", cell.ID, err)
		}
		roi := torsoROI(baseBounds).Intersect(cellBounds)
		if roi.Empty() {
			return nil, 0, fmt.Errorf("cell %q canonical torso ROI is empty", cell.ID)
		}
		baseMedian, err := subjectMedianLuma(base, roi)
		if err != nil {
			return nil, 0, fmt.Errorf("cell %q base torso ROI: %w", cell.ID, err)
		}
		if baseMedian <= 0 {
			return nil, 0, fmt.Errorf("cell %q base torso median luma is zero", cell.ID)
		}
		candidateMedian, err := subjectMedianLuma(candidate, roi)
		if err != nil {
			return nil, 0, fmt.Errorf("cell %q candidate torso ROI: %w", cell.ID, err)
		}
		ratio := candidateMedian / baseMedian
		if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
			return nil, 0, fmt.Errorf("cell %q produced invalid luma ratio", cell.ID)
		}
		baseRect := reportRect(baseBounds)
		roiRect := reportRect(roi)
		cellReport := newCellReport(cell, provenance)
		cellReport.BaseSubjectBounds = &baseRect
		cellReport.TorsoROI = &roiRect
		setSubjectGeometryReport(&cellReport, geometry)
		cellReport.BaseMedianLuma = floatPointer(baseMedian)
		cellReport.CandidateMedianLuma = floatPointer(candidateMedian)
		cellReport.Ratio = floatPointer(ratio)
		reports = append(reports, cellReport)
		sum += ratio
	}
	return reports, sum / float64(cellCount), nil
}

type subjectGeometry struct {
	ReferenceArea int
	CandidateArea int
	Intersection  int
	Union         int
	IoU           float64
	CentroidDX    float64
	CentroidDY    float64
	CentroidShift float64
}

func compareSubjectGeometry(reference image.Image, candidate image.Image, bounds image.Rectangle) (subjectGeometry, error) {
	if bounds.Empty() || !bounds.In(reference.Bounds()) || !bounds.In(candidate.Bounds()) {
		return subjectGeometry{}, fmt.Errorf("comparison bounds %v are outside reference %v or candidate %v", bounds, reference.Bounds(), candidate.Bounds())
	}

	var referenceSumX, referenceSumY int64
	var candidateSumX, candidateSumY int64
	metrics := subjectGeometry{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			referenceSubject := isSubject(nrgbaAt(reference, x, y))
			candidateSubject := isSubject(nrgbaAt(candidate, x, y))
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
		return subjectGeometry{}, errors.New("reference subject mask is empty")
	}
	if metrics.CandidateArea == 0 {
		return subjectGeometry{}, errors.New("candidate subject mask is empty")
	}
	if metrics.Union == 0 {
		return subjectGeometry{}, errors.New("subject mask union is empty")
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

func setSubjectGeometryReport(report *cellReport, metrics subjectGeometry) {
	report.ReferenceAreaPX = intPointer(metrics.ReferenceArea)
	report.CandidateAreaPX = intPointer(metrics.CandidateArea)
	report.IntersectionPX = intPointer(metrics.Intersection)
	report.UnionPX = intPointer(metrics.Union)
	report.IoU = floatPointer(metrics.IoU)
	report.CentroidDXPX = floatPointer(metrics.CentroidDX)
	report.CentroidDYPX = floatPointer(metrics.CentroidDY)
	report.CentroidShiftPX = floatPointer(metrics.CentroidShift)
}

func validateRawSubjectGeometry(reports []cellReport) error {
	return validateRawSubjectGeometryWithPolicy(reports, rawGeometryPolicy{
		ID:                 defaultRawGeometryPolicyReportID,
		MinIoU:             minRawSubjectIoU,
		MaxCentroidShiftPX: maxRawCentroidShiftPixels,
	})
}

func validateRawSubjectGeometryWithPolicy(reports []cellReport, policy rawGeometryPolicy) error {
	if len(reports) != cellCount {
		return fmt.Errorf("raw geometry gate received %d cells, want %d", len(reports), cellCount)
	}
	failures := make([]string, 0)
	for _, report := range reports {
		if report.IoU == nil || report.CentroidShiftPX == nil {
			failures = append(failures, fmt.Sprintf("cell %q is missing geometry metrics", report.ID))
			continue
		}
		if math.IsNaN(*report.IoU) || math.IsInf(*report.IoU, 0) ||
			math.IsNaN(*report.CentroidShiftPX) || math.IsInf(*report.CentroidShiftPX, 0) {
			failures = append(failures, fmt.Sprintf("cell %q has non-finite geometry metrics", report.ID))
			continue
		}
		if *report.IoU < policy.MinIoU {
			failures = append(failures, fmt.Sprintf(
				"cell %q IoU %g is below %g",
				report.ID,
				*report.IoU,
				policy.MinIoU,
			))
		}
		if *report.CentroidShiftPX > policy.MaxCentroidShiftPX {
			failures = append(failures, fmt.Sprintf(
				"cell %q centroid shift %gpx exceeds %gpx",
				report.ID,
				*report.CentroidShiftPX,
				policy.MaxCentroidShiftPX,
			))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("raw subject geometry gate failed: %s", strings.Join(failures, "; "))
	}
	return nil
}

func checkedSubjectBounds(img image.Image, cell image.Rectangle, label string, id string) (image.Rectangle, error) {
	var bounds image.Rectangle
	found := false
	for y := cell.Min.Y; y < cell.Max.Y; y++ {
		for x := cell.Min.X; x < cell.Max.X; x++ {
			pixel := nrgbaAt(img, x, y)
			if !isSubject(pixel) {
				continue
			}
			if x == cell.Min.X || x == cell.Max.X-1 || y == cell.Min.Y || y == cell.Max.Y-1 {
				return image.Rectangle{}, fmt.Errorf("cell %q %s has non-green subject content touching a cell edge at (%d,%d)", id, label, x-cell.Min.X, y-cell.Min.Y)
			}
			pointRect := image.Rect(x, y, x+1, y+1)
			if !found {
				bounds = pointRect
				found = true
			} else {
				bounds = bounds.Union(pointRect)
			}
		}
	}
	if !found {
		return image.Rectangle{}, fmt.Errorf("cell %q %s has empty subject", id, label)
	}
	return bounds, nil
}

func torsoROI(subject image.Rectangle) image.Rectangle {
	width := subject.Dx()
	height := subject.Dy()
	minX := subject.Min.X + int(math.RoundToEven(float64(width)*0.30))
	maxX := subject.Min.X + int(math.RoundToEven(float64(width)*0.68))
	minY := subject.Min.Y + int(math.RoundToEven(float64(height)*0.15))
	maxY := subject.Min.Y + int(math.RoundToEven(float64(height)*0.65))
	if maxX <= minX {
		maxX = minX + 1
	}
	if maxY <= minY {
		maxY = minY + 1
	}
	if maxX > subject.Max.X {
		maxX = subject.Max.X
	}
	if maxY > subject.Max.Y {
		maxY = subject.Max.Y
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func subjectMedianLuma(img image.Image, roi image.Rectangle) (float64, error) {
	values := make([]float64, 0, roi.Dx()*roi.Dy())
	for y := roi.Min.Y; y < roi.Max.Y; y++ {
		for x := roi.Min.X; x < roi.Max.X; x++ {
			pixel := nrgbaAt(img, x, y)
			if isSubject(pixel) {
				values = append(values, luma(pixel))
			}
		}
	}
	if len(values) == 0 {
		return 0, errors.New("contains no subject pixels")
	}
	sort.Float64s(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle], nil
	}
	return (values[middle-1] + values[middle]) / 2, nil
}

func luma(pixel color.NRGBA) float64 {
	return 0.2126*float64(pixel.R) + 0.7152*float64(pixel.G) + 0.0722*float64(pixel.B)
}

func isSubject(pixel color.NRGBA) bool {
	return pixel.A > 8 && !isChromaGreen(pixel)
}

func isChromaGreen(pixel color.NRGBA) bool {
	return pixel.A > 8 &&
		int(pixel.G) >= 90 &&
		int(pixel.G)-int(pixel.R) >= 25 &&
		int(pixel.G)-int(pixel.B) >= 25
}

func applyUniformSubjectGain(src image.Image, gain float64) *image.NRGBA {
	bounds := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			pixel := nrgbaAt(src, bounds.Min.X+x, bounds.Min.Y+y)
			if isSubject(pixel) {
				pixel.R = gainedChannel(pixel.R, gain)
				pixel.G = gainedChannel(pixel.G, gain)
				pixel.B = gainedChannel(pixel.B, gain)
			}
			out.SetNRGBA(x, y, pixel)
		}
	}
	return out
}

func gainedChannel(value uint8, gain float64) uint8 {
	scaled := math.Round(float64(value) * gain)
	if scaled < 0 {
		return 0
	}
	if scaled > 255 {
		return 255
	}
	return uint8(scaled)
}

func readExactSheet(path string) (decodedPNG, error) {
	decoded, err := readPNG(path)
	if err != nil {
		return decodedPNG{}, err
	}
	if decoded.Artifact.Width != sheetWidth || decoded.Artifact.Height != sheetHeight {
		return decodedPNG{}, fmt.Errorf("dimensions = %dx%d, want exactly %dx%d", decoded.Artifact.Width, decoded.Artifact.Height, sheetWidth, sheetHeight)
	}
	return decoded, nil
}

func readPNG(path string) (decodedPNG, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return decodedPNG{}, fmt.Errorf("read PNG: %w", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return decodedPNG{}, fmt.Errorf("decode PNG: %w", err)
	}
	bounds := img.Bounds()
	return decodedPNG{
		Image: img,
		Artifact: artifactReport{
			Path:   filepath.ToSlash(path),
			SHA256: sha256Hex(data),
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
		},
	}, nil
}

func encodePNG(img image.Image) ([]byte, error) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func artifactFromData(path string, data []byte, width int, height int) artifactReport {
	return artifactReport{
		Path:   filepath.ToSlash(path),
		SHA256: sha256Hex(data),
		Width:  width,
		Height: height,
	}
}

func writeArtifactAndReport(artifactPath string, artifactData []byte, reportPath string, report coatReport) error {
	reportData, err := encodeReport(report)
	if err != nil {
		return err
	}
	return writeTransaction([]encodedOutput{
		{Path: artifactPath, Data: artifactData},
		{Path: reportPath, Data: reportData},
	})
}

func writeReport(path string, report coatReport) error {
	data, err := encodeReport(report)
	if err != nil {
		return err
	}
	return writeTransaction([]encodedOutput{{Path: path, Data: data}})
}

func encodeReport(report coatReport) ([]byte, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode report: %w", err)
	}
	return append(data, '\n'), nil
}

type stagedOutput struct {
	Destination string
	Temporary   string
	Backup      string
	Existed     bool
	Committed   bool
}

// writeTransaction stages every payload in its destination directory, moves
// all existing destinations aside, and then commits the staged files. Any
// ordinary write or rename error is rolled back. No portable filesystem API can
// make several independent renames crash-atomic: process/OS failure between
// commit renames may leave a mixed set plus .coatbatch-stage/.coatbatch-backup
// recovery files.
func writeTransaction(outputs []encodedOutput) error {
	if len(outputs) == 0 {
		return errors.New("empty output transaction")
	}
	for i := range outputs {
		if outputs[i].Path == "" {
			return errors.New("empty output path")
		}
		for j := i + 1; j < len(outputs); j++ {
			if samePath(outputs[i].Path, outputs[j].Path) {
				return fmt.Errorf("transaction destinations must be different paths: %s", outputs[i].Path)
			}
		}
	}

	staged := make([]stagedOutput, 0, len(outputs))
	cleanupTemps := func() {
		for _, item := range staged {
			if item.Temporary != "" {
				_ = txFS.Remove(item.Temporary)
			}
		}
	}
	for _, output := range outputs {
		directory := filepath.Dir(output.Path)
		if err := os.MkdirAll(directory, 0o755); err != nil {
			cleanupTemps()
			return fmt.Errorf("create output directory for %s: %w", output.Path, err)
		}
		file, err := txFS.CreateTemp(directory, ".coatbatch-stage-*")
		if err != nil {
			cleanupTemps()
			return fmt.Errorf("stage %s: %w", output.Path, err)
		}
		item := stagedOutput{Destination: output.Path, Temporary: file.Name()}
		staged = append(staged, item)
		if err := file.Chmod(0o644); err != nil {
			_ = file.Close()
			cleanupTemps()
			return fmt.Errorf("set staged permissions for %s: %w", output.Path, err)
		}
		if err := txFS.WriteTemp(file, output.Data); err != nil {
			_ = file.Close()
			cleanupTemps()
			return fmt.Errorf("stage %s: %w", output.Path, err)
		}
		if err := file.Sync(); err != nil {
			_ = file.Close()
			cleanupTemps()
			return fmt.Errorf("sync staged %s: %w", output.Path, err)
		}
		if err := file.Close(); err != nil {
			cleanupTemps()
			return fmt.Errorf("close staged %s: %w", output.Path, err)
		}
	}

	rollback := func(cause error) error {
		var rollbackErrors []string
		for i := len(staged) - 1; i >= 0; i-- {
			item := &staged[i]
			if item.Committed {
				if err := txFS.Remove(item.Destination); err != nil && !os.IsNotExist(err) {
					rollbackErrors = append(rollbackErrors, fmt.Sprintf("remove %s: %v", item.Destination, err))
				}
				item.Committed = false
			}
		}
		for i := len(staged) - 1; i >= 0; i-- {
			item := &staged[i]
			if item.Backup != "" {
				if err := txFS.Rename(item.Backup, item.Destination); err != nil {
					rollbackErrors = append(rollbackErrors, fmt.Sprintf("restore %s: %v", item.Destination, err))
				} else {
					item.Backup = ""
				}
			}
		}
		cleanupTemps()
		if len(rollbackErrors) != 0 {
			return fmt.Errorf("%w (rollback incomplete: %s)", cause, strings.Join(rollbackErrors, "; "))
		}
		return cause
	}

	for i := range staged {
		item := &staged[i]
		info, err := os.Stat(item.Destination)
		switch {
		case err == nil:
			if info.IsDir() {
				return rollback(fmt.Errorf("destination is a directory: %s", item.Destination))
			}
			backupFile, err := txFS.CreateTemp(filepath.Dir(item.Destination), ".coatbatch-backup-*")
			if err != nil {
				return rollback(fmt.Errorf("prepare backup for %s: %w", item.Destination, err))
			}
			backupPath := backupFile.Name()
			if err := backupFile.Close(); err != nil {
				_ = txFS.Remove(backupPath)
				return rollback(fmt.Errorf("close backup placeholder for %s: %w", item.Destination, err))
			}
			if err := txFS.Remove(backupPath); err != nil {
				return rollback(fmt.Errorf("prepare backup path for %s: %w", item.Destination, err))
			}
			if err := txFS.Rename(item.Destination, backupPath); err != nil {
				return rollback(fmt.Errorf("back up %s: %w", item.Destination, err))
			}
			item.Existed = true
			item.Backup = backupPath
		case os.IsNotExist(err):
		default:
			return rollback(fmt.Errorf("inspect destination %s: %w", item.Destination, err))
		}
	}

	for i := range staged {
		item := &staged[i]
		if err := txFS.Rename(item.Temporary, item.Destination); err != nil {
			return rollback(fmt.Errorf("commit %s: %w", item.Destination, err))
		}
		item.Temporary = ""
		item.Committed = true
	}
	for i := range staged {
		if staged[i].Backup != "" {
			if err := txFS.Remove(staged[i].Backup); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove transaction backup %s: %w", staged[i].Backup, err)
			}
			staged[i].Backup = ""
		}
	}
	return nil
}

func resolveConfinedExisting(root string, relative string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(filepath.Join(rootAbs, relative))
	if err != nil {
		return "", err
	}
	if !pathWithin(rootAbs, targetAbs) {
		return "", fmt.Errorf("path escapes manifest directory: %q", relative)
	}
	resolvedRoot, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", fmt.Errorf("resolve manifest directory: %w", err)
	}
	resolvedTarget, err := filepath.EvalSymlinks(targetAbs)
	if err != nil {
		return "", fmt.Errorf("resolve source: %w", err)
	}
	if !pathWithin(resolvedRoot, resolvedTarget) {
		return "", fmt.Errorf("path escapes manifest directory through a link: %q", relative)
	}
	return resolvedTarget, nil
}

func resolveConfinedOutput(root string, relative string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(filepath.Join(rootAbs, relative))
	if err != nil {
		return "", err
	}
	if !pathWithin(rootAbs, targetAbs) {
		return "", fmt.Errorf("path escapes output directory: %q", relative)
	}
	if _, err := os.Stat(rootAbs); err == nil {
		resolvedRoot, err := filepath.EvalSymlinks(rootAbs)
		if err != nil {
			return "", fmt.Errorf("resolve output directory: %w", err)
		}
		existing, suffix, err := nearestExistingAncestor(filepath.Dir(targetAbs))
		if err != nil {
			return "", err
		}
		resolvedExisting, err := filepath.EvalSymlinks(existing)
		if err != nil {
			return "", fmt.Errorf("resolve output parent: %w", err)
		}
		resolvedTarget := filepath.Join(append([]string{resolvedExisting}, suffix...)...)
		resolvedTarget = filepath.Join(resolvedTarget, filepath.Base(targetAbs))
		if !pathWithin(resolvedRoot, resolvedTarget) {
			return "", fmt.Errorf("path escapes output directory through a link: %q", relative)
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("inspect output directory: %w", err)
	}
	return targetAbs, nil
}

func nearestExistingAncestor(path string) (string, []string, error) {
	current := filepath.Clean(path)
	var suffix []string
	for {
		_, err := os.Lstat(current)
		if err == nil {
			return current, suffix, nil
		}
		if !os.IsNotExist(err) {
			return "", nil, fmt.Errorf("inspect output parent: %w", err)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", nil, fmt.Errorf("no existing ancestor for output path %q", path)
		}
		suffix = append([]string{filepath.Base(current)}, suffix...)
		current = parent
	}
}

func pathWithin(root string, target string) bool {
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)))
}

type namedPath struct {
	Name string
	Path string
}

func requireDistinctPaths(paths ...namedPath) error {
	for i := range paths {
		if paths[i].Path == "" {
			continue
		}
		for j := i + 1; j < len(paths); j++ {
			if paths[j].Path != "" && samePath(paths[i].Path, paths[j].Path) {
				return fmt.Errorf("%s and %s must be different paths", paths[i].Name, paths[j].Name)
			}
		}
	}
	return nil
}

func samePath(a string, b string) bool {
	aAbs, aErr := filepath.Abs(a)
	bAbs, bErr := filepath.Abs(b)
	if aErr != nil || bErr != nil {
		return false
	}
	aAbs = filepath.Clean(aAbs)
	bAbs = filepath.Clean(bAbs)
	if strings.EqualFold(aAbs, bAbs) {
		return true
	}
	aInfo, aStatErr := os.Stat(aAbs)
	bInfo, bStatErr := os.Stat(bAbs)
	if aStatErr == nil && bStatErr == nil && os.SameFile(aInfo, bInfo) {
		return true
	}
	aResolved, aResolveErr := resolveThroughExistingParent(aAbs)
	bResolved, bResolveErr := resolveThroughExistingParent(bAbs)
	return aResolveErr == nil && bResolveErr == nil && strings.EqualFold(aResolved, bResolved)
}

func resolveThroughExistingParent(path string) (string, error) {
	if _, err := os.Lstat(path); err == nil {
		return filepath.EvalSymlinks(path)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	existing, suffix, err := nearestExistingAncestor(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	resolvedExisting, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", err
	}
	parts := append([]string{resolvedExisting}, suffix...)
	parts = append(parts, filepath.Base(path))
	return filepath.Clean(filepath.Join(parts...)), nil
}

func requirePNGOutput(path string) error {
	if !strings.EqualFold(filepath.Ext(path), ".png") {
		return fmt.Errorf("-out must end in .png: %q", path)
	}
	return nil
}

func unsafeWindowsPathComponent(component string) bool {
	if component == "" || strings.HasSuffix(component, ".") || strings.HasSuffix(component, " ") ||
		strings.ContainsAny(component, `<>:"|?*`) {
		return true
	}
	stem := component
	if dot := strings.IndexByte(stem, '.'); dot >= 0 {
		stem = stem[:dot]
	}
	switch strings.ToUpper(stem) {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return true
	}
	return false
}

func boundsForCell(cell int) image.Rectangle {
	column := cell % 2
	row := cell / 2
	return image.Rect(column*cellWidth, row*cellHeight, (column+1)*cellWidth, (row+1)*cellHeight)
}

func reportRect(rect image.Rectangle) rectReport {
	return rectReport{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()}
}

func nrgbaAt(img image.Image, x int, y int) color.NRGBA {
	if nrgba, ok := img.(*image.NRGBA); ok {
		return nrgba.NRGBAAt(x, y)
	}
	return color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
}

func intPointer(value int) *int {
	return &value
}

func floatPointer(value float64) *float64 {
	return &value
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
