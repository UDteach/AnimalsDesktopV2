package catalog

import (
	"strings"
	"testing"
)

func TestValidateCatalog(t *testing.T) {
	if err := Validate(); err != nil {
		t.Fatalf("Validate() error: %v", err)
	}
}

func TestValidateCatalogRejectsBrokenAdditionContracts(t *testing.T) {
	validSpecies := []Species{
		{ID: "mouse", Label: "Mouse", Profile: "small-mammal"},
	}
	validVariant := Variant{
		ID:               "mouse_gray",
		SpeciesID:        "mouse",
		BreedOrMorph:     "House mouse",
		Color:            "gray",
		PopularityTier:   1,
		MotionProfile:    MotionProfileSmallRodentScurry,
		SourceStatus:     SourceStatusMotionAccepted,
		LabelEN:          "Mouse - gray",
		LabelJA:          "マウス（グレー）",
		SpriteBase:       "mouse_gray",
		SeedStage:        true,
		SourcePath:       "docs/art-source/mouse/motion-source/accepted-frames/set00/frame-00.png",
		MotionSourcePath: "docs/art-source/mouse/motion-source/sheets/mouse-gray-source-set00.png",
	}

	tests := []struct {
		name       string
		species    []Species
		variants   []Variant
		runtimeIDs []string
		want       string
	}{
		{
			name:       "invalid variant id",
			species:    validSpecies,
			variants:   []Variant{withVariantID(validVariant, "Mouse Gray")},
			runtimeIDs: []string{"Mouse Gray"},
			want:       `invalid id "Mouse Gray"`,
		},
		{
			name:       "unknown species",
			species:    validSpecies,
			variants:   []Variant{withVariantSpecies(validVariant, "rat")},
			runtimeIDs: []string{"mouse_gray"},
			want:       `references unknown species "rat"`,
		},
		{
			name:       "accepted source without motion sheet",
			species:    validSpecies,
			variants:   []Variant{withMotionSource(validVariant, "")},
			runtimeIDs: []string{"mouse_gray"},
			want:       `has no motion source path`,
		},
		{
			name:       "unclean source path",
			species:    validSpecies,
			variants:   []Variant{withSourcePath(validVariant, "../mouse.png")},
			runtimeIDs: []string{"mouse_gray"},
			want:       `path is not clean and repository-relative`,
		},
		{
			name:       "duplicate variant",
			species:    validSpecies,
			variants:   []Variant{validVariant, validVariant},
			runtimeIDs: []string{"mouse_gray"},
			want:       `duplicate variant id "mouse_gray"`,
		},
		{
			name:       "unknown runtime id",
			species:    validSpecies,
			variants:   []Variant{validVariant},
			runtimeIDs: []string{"mouse_white"},
			want:       `references unknown id "mouse_white"`,
		},
		{
			name:       "duplicate runtime id",
			species:    validSpecies,
			variants:   []Variant{validVariant},
			runtimeIDs: []string{"mouse_gray", "mouse_gray"},
			want:       `runtime variant id "mouse_gray" is duplicated`,
		},
		{
			name:       "empty runtime roster",
			species:    validSpecies,
			variants:   []Variant{validVariant},
			runtimeIDs: nil,
			want:       `runtime variant roster is empty`,
		},
		{
			name:       "non-degu cannot use legacy source",
			species:    validSpecies,
			variants:   []Variant{withSeedAndStatus(validVariant, false, SourceStatusDeguMotion)},
			runtimeIDs: []string{"mouse_gray"},
			want:       `non-degu variant "mouse_gray" must be marked SeedStage`,
		},
		{
			name: "degu cannot use animal addition source",
			species: []Species{
				{ID: "degu", Label: "Degu", Profile: "degu"},
			},
			variants: []Variant{
				withVariantSpecies(validVariant, "degu"),
			},
			runtimeIDs: []string{"mouse_gray"},
			want:       `degu variant "mouse_gray" cannot be marked SeedStage`,
		},
		{
			name:       "invalid tint color",
			species:    validSpecies,
			variants:   []Variant{withTintHex(validVariant, "#8899aa")},
			runtimeIDs: []string{"mouse_gray"},
			want:       `invalid tint color "#8899aa"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateCatalog(test.species, test.variants, test.runtimeIDs)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validateCatalog() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func withVariantID(variant Variant, id string) Variant {
	variant.ID = id
	variant.SpriteBase = id
	return variant
}

func withVariantSpecies(variant Variant, speciesID string) Variant {
	variant.SpeciesID = speciesID
	return variant
}

func withMotionSource(variant Variant, source string) Variant {
	variant.MotionSourcePath = source
	return variant
}

func withSourcePath(variant Variant, source string) Variant {
	variant.SourcePath = source
	return variant
}

func withSeedAndStatus(variant Variant, seedStage bool, sourceStatus string) Variant {
	variant.SeedStage = seedStage
	variant.SourceStatus = sourceStatus
	variant.MotionSourcePath = ""
	return variant
}

func withTintHex(variant Variant, tintHex string) Variant {
	variant.TintHex = tintHex
	return variant
}
