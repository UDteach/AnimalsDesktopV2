# Cohort 09 Low-Load Report

RAM protection mode is active. No ImageGen, image generation, bulk file generation, importer runs, or tests were executed.

## Variants

- `siamese_seal_point`
- `sphynx_pink`
- `scottish_fold_silver_tabby`
- `bengal_rosetted`
- `domestic_shorthair_calico`
- `domestic_shorthair_tuxedo`
- `holland_lop_broken_orange`
- `netherland_dwarf_chestnut`
- `mini_rex_black_otter`
- `lionhead_tort`

## ImageGen Prompt Template

Create a complete, single animal sprite source for `{variant_id}`. Transparent background. Side-view desktop pet sprite. No text, border, scenery, shadow, costume, props, multiple animals, or human-like pose. Preserve full ears, feet, whiskers, and tail. Keep camera, scale, anatomy, baseline, and contact points consistent across all frames. Output one source-truth PNG and 62 motion source frames for the assigned motion profile.

## Required Files

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00.png` through `61.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

## Acceptance Checklist

- Cat variants use `cat-stalk`; rabbit variants use `rabbit-hop`.
- Rabbits keep lop/dwarf/rex/lionhead silhouettes distinct.
- Coat patterns remain readable at 96x64.
- No degu recolor or duplicated degu motion.
