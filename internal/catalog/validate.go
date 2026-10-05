package catalog

import (
	"encoding/hex"
	"errors"
	"fmt"
	"path"
	"strings"
)

// Validate checks the catalog contract used by importers, runtime selection,
// and release tooling. Keep this free of filesystem access so callers can run
// it before generating or replacing assets.
func Validate() error {
	return validateCatalog(SpeciesList, Variants, runtimeVariantIDs)
}

func validateCatalog(speciesList []Species, variants []Variant, runtimeIDs []string) error {
	var issues []error

	speciesByID := make(map[string]Species, len(speciesList))
	for index, species := range speciesList {
		prefix := fmt.Sprintf("species[%d]", index)
		if !validCatalogID(species.ID) {
			issues = append(issues, fmt.Errorf("%s has invalid id %q", prefix, species.ID))
		}
		if strings.TrimSpace(species.Label) == "" {
			issues = append(issues, fmt.Errorf("%s %q has no label", prefix, species.ID))
		}
		if strings.TrimSpace(species.Profile) == "" {
			issues = append(issues, fmt.Errorf("%s %q has no profile", prefix, species.ID))
		}
		if _, exists := speciesByID[species.ID]; exists {
			issues = append(issues, fmt.Errorf("duplicate species id %q", species.ID))
		}
		speciesByID[species.ID] = species
	}

	variantByID := make(map[string]Variant, len(variants))
	spriteBases := make(map[string]string, len(variants))
	for index, variant := range variants {
		prefix := fmt.Sprintf("variant[%d]", index)
		if !validCatalogID(variant.ID) {
			issues = append(issues, fmt.Errorf("%s has invalid id %q", prefix, variant.ID))
		}
		if !validCatalogID(variant.SpriteBase) {
			issues = append(issues, fmt.Errorf("%s %q has invalid sprite base %q", prefix, variant.ID, variant.SpriteBase))
		}
		if _, exists := variantByID[variant.ID]; exists {
			issues = append(issues, fmt.Errorf("duplicate variant id %q", variant.ID))
		}
		variantByID[variant.ID] = variant
		if previous, exists := spriteBases[variant.SpriteBase]; exists {
			issues = append(issues, fmt.Errorf("variants %q and %q share sprite base %q", previous, variant.ID, variant.SpriteBase))
		}
		spriteBases[variant.SpriteBase] = variant.ID

		if _, exists := speciesByID[variant.SpeciesID]; !exists {
			issues = append(issues, fmt.Errorf("variant %q references unknown species %q", variant.ID, variant.SpeciesID))
		}
		if strings.TrimSpace(variant.BreedOrMorph) == "" {
			issues = append(issues, fmt.Errorf("variant %q has no breed or morph", variant.ID))
		}
		if strings.TrimSpace(variant.Color) == "" {
			issues = append(issues, fmt.Errorf("variant %q has no color", variant.ID))
		}
		if variant.PopularityTier < 1 || variant.PopularityTier > 3 {
			issues = append(issues, fmt.Errorf("variant %q has popularity tier %d, want 1-3", variant.ID, variant.PopularityTier))
		}
		if !knownMotionProfile(variant.MotionProfile) {
			issues = append(issues, fmt.Errorf("variant %q has unknown motion profile %q", variant.ID, variant.MotionProfile))
		}
		if !knownSourceStatus(variant.SourceStatus) {
			issues = append(issues, fmt.Errorf("variant %q has unknown source status %q", variant.ID, variant.SourceStatus))
		}
		if strings.TrimSpace(variant.LabelEN) == "" || strings.TrimSpace(variant.LabelJA) == "" {
			issues = append(issues, fmt.Errorf("variant %q needs both English and Japanese labels", variant.ID))
		}
		if issue := validateCatalogPath(variant.ID, "source", variant.SourcePath); issue != nil {
			issues = append(issues, issue)
		}
		if issue := validateCatalogPath(variant.ID, "motion source", variant.MotionSourcePath); issue != nil {
			issues = append(issues, issue)
		}
		if !validOptionalHexColor(variant.TintHex) {
			issues = append(issues, fmt.Errorf("variant %q has invalid tint color %q; want six hexadecimal digits", variant.ID, variant.TintHex))
		}
		if !validOptionalHexColor(variant.AccentHex) {
			issues = append(issues, fmt.Errorf("variant %q has invalid accent color %q; want six hexadecimal digits", variant.ID, variant.AccentHex))
		}

		isDegu := variant.SpeciesID == "degu"
		if isDegu {
			if variant.SeedStage {
				issues = append(issues, fmt.Errorf("degu variant %q cannot be marked SeedStage", variant.ID))
			}
			if variant.SourceStatus != SourceStatusDeguMotion {
				issues = append(issues, fmt.Errorf("degu variant %q has source status %q, want %q", variant.ID, variant.SourceStatus, SourceStatusDeguMotion))
			}
			if variant.MotionProfile != MotionProfileDegu {
				issues = append(issues, fmt.Errorf("degu variant %q has motion profile %q, want %q", variant.ID, variant.MotionProfile, MotionProfileDegu))
			}
		} else if !variant.SeedStage {
			issues = append(issues, fmt.Errorf("non-degu variant %q must be marked SeedStage", variant.ID))
		}

		if variant.SeedStage {
			if variant.SourcePath == "" && variant.Shape == "" {
				issues = append(issues, fmt.Errorf("seed variant %q has neither source path nor procedural shape", variant.ID))
			}
			if variant.SourceStatus == SourceStatusDeguMotion {
				issues = append(issues, fmt.Errorf("seed variant %q cannot use degu motion source status", variant.ID))
			}
		}

		hasMotionSource := variant.MotionSourcePath != ""
		needsMotionSource := variant.SourceStatus == SourceStatusMotionDraft || variant.SourceStatus == SourceStatusMotionAccepted
		if needsMotionSource && !hasMotionSource {
			issues = append(issues, fmt.Errorf("motion variant %q has no motion source path", variant.ID))
		}
		if hasMotionSource && !needsMotionSource {
			issues = append(issues, fmt.Errorf("variant %q has a motion source path with source status %q", variant.ID, variant.SourceStatus))
		}
	}

	if len(runtimeIDs) == 0 {
		issues = append(issues, errors.New("runtime variant roster is empty"))
	}
	runtimeSeen := make(map[string]bool, len(runtimeIDs))
	for index, id := range runtimeIDs {
		if runtimeSeen[id] {
			issues = append(issues, fmt.Errorf("runtime variant id %q is duplicated", id))
			continue
		}
		runtimeSeen[id] = true
		variant, exists := variantByID[id]
		if !exists {
			issues = append(issues, fmt.Errorf("runtime variant[%d] references unknown id %q", index, id))
			continue
		}
		if !variant.SeedStage {
			issues = append(issues, fmt.Errorf("runtime variant %q is not importable", id))
		}
		if variant.SourceStatus != SourceStatusMotionAccepted {
			issues = append(issues, fmt.Errorf("runtime variant %q has source status %q, want %q", id, variant.SourceStatus, SourceStatusMotionAccepted))
		}
	}

	return errors.Join(issues...)
}

func validOptionalHexColor(value string) bool {
	if value == "" {
		return true
	}
	if len(value) != 6 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func validCatalogID(value string) bool {
	if value == "" || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for index := 1; index < len(value); index++ {
		char := value[index]
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' {
			continue
		}
		return false
	}
	return true
}

func validateCatalogPath(variantID string, label string, value string) error {
	if value == "" {
		return nil
	}
	if strings.Contains(value, `\`) || strings.HasPrefix(value, "/") || strings.Contains(value, ":") {
		return fmt.Errorf("variant %q %s path must be repository-relative with forward slashes: %q", variantID, label, value)
	}
	cleaned := path.Clean(value)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || cleaned != value {
		return fmt.Errorf("variant %q %s path is not clean and repository-relative: %q", variantID, label, value)
	}
	return nil
}

func knownSourceStatus(status string) bool {
	switch status {
	case SourceStatusDeguMotion,
		SourceStatusPrototypeOnly,
		SourceStatusImageGenQueued,
		SourceStatusMotionDraft,
		SourceStatusMotionAccepted:
		return true
	default:
		return false
	}
}

func knownMotionProfile(profile string) bool {
	switch profile {
	case MotionProfileDegu,
		MotionProfileSmallRodentScurry,
		MotionProfileRabbitHop,
		MotionProfileDogTrot,
		MotionProfileCatStalk,
		MotionProfileGeckoCrawl,
		MotionProfileTortoisePlod,
		MotionProfileFerretSlink,
		MotionProfileGuineaPigWaddle,
		MotionProfileHedgehogShuffle,
		MotionProfileSquirrelBound,
		MotionProfileFoxTrot,
		MotionProfileRedPandaAmble,
		MotionProfileOtterSlide,
		MotionProfileSugarGliderSkitter,
		MotionProfileCapybaraLumber,
		MotionProfileSnakeSlither,
		MotionProfileDragonPlod,
		MotionProfileFrogHop,
		MotionProfileBirdHop:
		return true
	default:
		return false
	}
}
