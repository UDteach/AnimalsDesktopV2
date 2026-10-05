# Cohort 10 Low-Load Report

RAM protection mode is active. No ImageGen, image generation, bulk file generation, importer runs, or tests were executed.

## Variants

- `dutch_black_white`
- `fancy_rat_hooded`
- `fancy_mouse_white`
- `mongolian_gerbil_agouti`
- `prairie_dog_tan`
- `chipmunk_striped`
- `bearded_dragon_citrus`
- `crested_gecko_harlequin`
- `corn_snake_amelanistic`
- `whites_tree_frog_green`

## ImageGen Prompt Template

Create a complete, single animal sprite source for `{variant_id}`. Transparent background. Side-view desktop pet sprite. No text, border, scenery, shadow, costume, props, multiple animals, or human-like pose. Preserve full ears, feet, whiskers, and tail. Keep camera, scale, anatomy, baseline, and contact points consistent across all frames. Output one source-truth PNG and 62 motion source frames for the assigned motion profile.

## Required Files

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00.png` through `61.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

## Acceptance Checklist

- Small mammals use natural scurry/bound motion, not degu wheel motion.
- Bearded dragon and crested gecko stay low-crawling.
- Corn snake uses `snake-slither` with no legs or vertical hopping.
- White's tree frog uses `frog-hop` without entering wheel actions.
