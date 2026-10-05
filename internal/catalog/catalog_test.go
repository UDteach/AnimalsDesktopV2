package catalog

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogInvariants(t *testing.T) {
	if err := Validate(); err != nil {
		t.Fatalf("catalog validation failed: %v", err)
	}

	seedCount := 0
	deguCount := 0
	for _, variant := range Variants {
		if variant.SeedStage {
			seedCount++
		}
		if variant.SpeciesID == "degu" {
			deguCount++
		}
	}
	if got := len(SeedVariants()); got != seedCount {
		t.Fatalf("SeedVariants() = %d, want %d catalog seed variants", got, seedCount)
	}
	if got := len(DeguVariantIDs()); got != deguCount {
		t.Fatalf("DeguVariantIDs() = %d, want %d catalog degu variants", got, deguCount)
	}
}

func TestVariantGroupsClassifyRuntimeAnimals(t *testing.T) {
	tests := map[string]string{
		"shoebill":     "bird",
		"cat":          "cat",
		"rabbit":       "rabbit",
		"chinchilla":   "chinchilla",
		"ferret":       "small_mammal",
		"gecko":        "reptile_amphibian",
		"sugar_glider": "sugar_glider",
	}
	for speciesID, want := range tests {
		if got := VariantGroupIDForSpecies(speciesID); got != want {
			t.Fatalf("VariantGroupIDForSpecies(%q) = %q, want %q", speciesID, got, want)
		}
	}
	group := VariantGroupForSpecies("shoebill")
	if group.LabelJA != "鳥" || group.LabelEN != "Birds" {
		t.Fatalf("shoebill group labels = %+v", group)
	}
}

func TestRuntimeVariantsAreReleaseScoped(t *testing.T) {
	runtime := RuntimeVariants()
	if got := len(runtime); got != len(runtimeVariantIDs) {
		t.Fatalf("runtime variants = %d, want %d release-scoped ids", got, len(runtimeVariantIDs))
	}
	for i, variant := range runtime {
		if variant.ID != runtimeVariantIDs[i] {
			t.Fatalf("runtime variant[%d] = %q, want %q", i, variant.ID, runtimeVariantIDs[i])
		}
		if variant.SpeciesID == "degu" {
			t.Fatalf("runtime variants must not include degu: %+v", variant)
		}
		if variant.SourceStatus != SourceStatusMotionAccepted {
			t.Fatalf("runtime variant %q source status = %q, want accepted", variant.ID, variant.SourceStatus)
		}
	}
}

func TestAcceptedFerretVariantsAreRuntimeScopedAtEnd(t *testing.T) {
	want := []struct {
		id               string
		labelEN          string
		labelJA          string
		sourcePath       string
		motionSourcePath string
		tintHex          string
		accentHex        string
		color            string
		popularityTier   int
	}{
		{
			id:               "ferret_sable_panda",
			labelEN:          "Ferret - sable panda",
			labelJA:          "フェレット（セーブルパンダ）",
			sourcePath:       "docs/art-source/ferret-sable-panda/motion-source/accepted-frames/set00/frame-00.png",
			motionSourcePath: "docs/art-source/ferret-sable-panda/motion-source/sheets/ferret-sable-panda-source-set00.png",
			tintHex:          "eee7d9",
			accentHex:        "6e5140",
			color:            "sable panda",
			popularityTier:   2,
		},
		{
			id:               "ferret_sable",
			labelEN:          "Ferret - sable",
			labelJA:          "フェレット（セーブル）",
			sourcePath:       "docs/art-source/ferret-sable/motion-source/accepted-frames/set00/frame-00.png",
			motionSourcePath: "docs/art-source/ferret-sable/motion-source/sheets/ferret-sable-source-set00.png",
			tintHex:          "8b6746",
			accentHex:        "ece0c8",
			color:            "sable",
			popularityTier:   1,
		},
		{
			id:               "ferret_albino",
			labelEN:          "Ferret - albino",
			labelJA:          "フェレット（アルビノ）",
			sourcePath:       "docs/art-source/ferret-albino/motion-source/accepted-frames/set00/frame-00.png",
			motionSourcePath: "docs/art-source/ferret-albino/motion-source/sheets/ferret-albino-source-set00.png",
			tintHex:          "eadcc7",
			accentHex:        "c58b78",
			color:            "albino",
			popularityTier:   2,
		},
	}

	runtime := RuntimeVariants()
	if len(runtime) < len(want) {
		t.Fatalf("runtime variants = %d, need at least %d ferret variants", len(runtime), len(want))
	}
	runtimeTail := runtime[len(runtime)-len(want):]
	for i, expected := range want {
		if runtimeTail[i].ID != expected.id {
			t.Fatalf("runtime tail[%d] = %q, want %q", i, runtimeTail[i].ID, expected.id)
		}

		variant, ok := VariantByID(expected.id)
		if !ok {
			t.Fatalf("missing accepted ferret variant %q", expected.id)
		}
		if variant.SpeciesID != "ferret" ||
			variant.BreedOrMorph != "Ferret" ||
			variant.LabelEN != expected.labelEN ||
			variant.LabelJA != expected.labelJA ||
			variant.SpriteBase != expected.id ||
			variant.SourcePath != expected.sourcePath ||
			variant.MotionSourcePath != expected.motionSourcePath ||
			variant.TintHex != expected.tintHex ||
			variant.AccentHex != expected.accentHex ||
			variant.Color != expected.color ||
			variant.PopularityTier != expected.popularityTier ||
			!variant.SeedStage ||
			MotionProfileForVariant(variant) != MotionProfileFerretSlink ||
			variant.SourceStatus != SourceStatusMotionAccepted {
			t.Fatalf("accepted ferret variant %q metadata = %+v", expected.id, variant)
		}
	}

	champagne, ok := VariantByID("ferret_champagne")
	if !ok {
		t.Fatal("missing prototype ferret_champagne")
	}
	if champagne.SourceStatus != SourceStatusPrototypeOnly {
		t.Fatalf("ferret_champagne source status = %q, want prototype", champagne.SourceStatus)
	}
}

func TestRuntimeSpritesMatchAcceptedMotionSources(t *testing.T) {
	const runtimeMotionSets = 10
	for _, variant := range RuntimeVariants() {
		sourcePaths := expectedRuntimeMotionSources(t, variant.MotionSourcePath, runtimeMotionSets)
		for set := 0; set < runtimeMotionSets; set++ {
			runtimePath := repoPath("assets", "sprites", fmt.Sprintf("%s_set%02d.png", variant.SpriteBase, set))
			runtimeImg, err := readPNG(runtimePath)
			if err != nil {
				t.Fatalf("read runtime sprite %s: %v", runtimePath, err)
			}
			sourceImg, err := readPNG(sourcePaths[set])
			if err != nil {
				t.Fatalf("read motion source %s: %v", sourcePaths[set], err)
			}
			if !imagesEqual(runtimeImg, sourceImg) {
				t.Fatalf("runtime sprite %s does not match accepted source %s", runtimePath, sourcePaths[set])
			}
		}
	}
}

func readPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func imagesEqual(a image.Image, b image.Image) bool {
	ab := a.Bounds()
	bb := b.Bounds()
	if ab != bb {
		return false
	}
	for y := ab.Min.Y; y < ab.Max.Y; y++ {
		for x := ab.Min.X; x < ab.Max.X; x++ {
			ar, ag, ablu, aa := a.At(x, y).RGBA()
			br, bg, bblu, ba := b.At(x, y).RGBA()
			if ar != br || ag != bg || ablu != bblu || aa != ba {
				return false
			}
		}
	}
	return true
}

func expectedRuntimeMotionSources(t *testing.T, set00Path string, sets int) []string {
	t.Helper()
	if set00Path == "" {
		t.Fatalf("runtime variant has no motion source path")
	}
	set00 := repoPath(filepath.FromSlash(set00Path))
	out := make([]string, sets)
	for set := 0; set < sets; set++ {
		candidate := set00
		if strings.Contains(set00, "set00") {
			candidate = strings.Replace(set00, "set00", fmt.Sprintf("set%02d", set), 1)
		}
		if _, err := os.Stat(candidate); err != nil {
			if set > 0 && os.IsNotExist(err) {
				candidate = set00
			} else {
				t.Fatalf("stat motion source %s: %v", candidate, err)
			}
		}
		out[set] = candidate
	}
	return out
}

func repoPath(parts ...string) string {
	all := append([]string{"..", ".."}, parts...)
	return filepath.Join(all...)
}

func TestWheelCapabilityIsLimitedToChinchillaAndHamster(t *testing.T) {
	wantBySpecies := map[string]bool{
		"chinchilla":     true,
		"hamster":        true,
		"degu":           false,
		"macaroni_mouse": false,
		"mouse":          false,
		"sugar_glider":   false,
		"rabbit":         false,
		"gecko":          false,
	}
	for species, want := range wantBySpecies {
		if got := WheelCapableSpecies(species); got != want {
			t.Fatalf("WheelCapableSpecies(%q) = %v, want %v", species, got, want)
		}
	}
	if WheelCapableMotionProfile(MotionProfileSmallRodentScurry) {
		t.Fatalf("small-rodent motion profile must not imply wheel capability")
	}
	wantByVariant := map[string]bool{
		"chinchilla_standard_gray": true,
		"chinchilla_beige":         true,
		"chinchilla_ebony":         true,
		"hamster_golden_syrian":    true,
		"macaroni_mouse_tan":       false,
		"sugar_glider_gray":        false,
		"rabbit_chestnut_agouti":   false,
		"wild_agouti":              false,
	}
	for id, want := range wantByVariant {
		variant, ok := VariantByID(id)
		if !ok {
			t.Fatalf("missing test variant %q", id)
		}
		if got := WheelCapableVariant(variant); got != want {
			t.Fatalf("WheelCapableVariant(%q) = %v, want %v", id, got, want)
		}
	}
}

func TestRequestedPopularVariantsArePresent(t *testing.T) {
	required := []string{
		"french_bulldog_fawn",
		"labrador_yellow",
		"golden_retriever_golden",
		"maine_coon_brown_tabby",
		"ragdoll_seal_bicolor",
		"holland_lop_broken_orange",
		"fancy_rat_hooded",
		"bearded_dragon_citrus",
		"corn_snake_amelanistic",
		"whites_tree_frog_green",
	}
	seen := map[string]bool{}
	for _, variant := range Variants {
		seen[variant.ID] = true
	}
	for _, id := range required {
		if !seen[id] {
			t.Fatalf("missing requested variant %q", id)
		}
	}
}

func TestMotionProfilesCoverCatalog(t *testing.T) {
	known := map[string]bool{
		MotionProfileDegu:               true,
		MotionProfileSmallRodentScurry:  true,
		MotionProfileRabbitHop:          true,
		MotionProfileDogTrot:            true,
		MotionProfileCatStalk:           true,
		MotionProfileGeckoCrawl:         true,
		MotionProfileTortoisePlod:       true,
		MotionProfileFerretSlink:        true,
		MotionProfileGuineaPigWaddle:    true,
		MotionProfileHedgehogShuffle:    true,
		MotionProfileSquirrelBound:      true,
		MotionProfileFoxTrot:            true,
		MotionProfileRedPandaAmble:      true,
		MotionProfileOtterSlide:         true,
		MotionProfileSugarGliderSkitter: true,
		MotionProfileCapybaraLumber:     true,
		MotionProfileSnakeSlither:       true,
		MotionProfileDragonPlod:         true,
		MotionProfileFrogHop:            true,
		MotionProfileBirdHop:            true,
	}
	for _, variant := range Variants {
		if !known[MotionProfileForVariant(variant)] {
			t.Fatalf("variant %q has unknown motion profile %q", variant.ID, MotionProfileForVariant(variant))
		}
	}
}
