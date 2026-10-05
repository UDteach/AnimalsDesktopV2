# Cohort 08 Low-Load Report

RAM protection mode is active. No ImageGen, image generation, bulk file generation, importer runs, or tests were executed.

## Variants

- `poodle_white`
- `beagle_tricolor`
- `bulldog_white_brindle`
- `shiba_inu_red`
- `pomeranian_orange`
- `corgi_sable`
- `maine_coon_brown_tabby`
- `ragdoll_seal_bicolor`
- `persian_white`
- `british_shorthair_blue`

## ImageGen Prompt Template

Create a complete, single animal sprite source for `{variant_id}`. Transparent background. Side-view desktop pet sprite. No text, border, scenery, shadow, costume, props, multiple animals, or human-like pose. Preserve full ears, feet, whiskers, and tail. Keep camera, scale, anatomy, baseline, and contact points consistent across all frames. Output one source-truth PNG and 62 motion source frames for the assigned motion profile.

## Required Files

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00.png` through `61.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

## Acceptance Checklist

- One complete animal per frame.
- Transparent background and alpha-contained content.
- No cropped ears, feet, whiskers, or tails.
- Dog variants use `dog-trot`; cat variants use `cat-stalk`.
- Breed identity is visible without scenery or labels.
