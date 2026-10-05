# Cohort 03 Known Risks

Status: no generated PNG files are present from this delegated thread.

## Cross-Variant Risks

- ImageGen may drift camera, scale, baseline, or body volume across 62 frames. Parent integration must reject frame sets with visible size jumps.
- Transparent-background requests can still produce halos, shadows, or checkerboard artifacts. Check alpha bounds and light/dark previews.
- Small features can be cropped at 96x64, especially rabbit ears, dog tails, cat whiskers, and cat tails.
- Recolor-only dog variants are not enough for final quality. Coat markings must be intentional and breed/coat credible.
- Frame group 26-31 must be species-specific. Reject degu-like wheel/run poses for rabbits, dogs, and cats.

## Variant-Specific Risks

- `rabbit_white`: white fur may lose edge readability on light backgrounds. Require subtle internal shading but no drop shadow.
- `rabbit_blue_gray`: blue-gray coat may become plain gray. Keep a cool gray tone while preserving rabbit anatomy.
- `rabbit_fawn`: fawn color can drift toward orange dog/cat coloring. Keep rabbit silhouette and compact body.
- `dog_cream_tan`: tan markings may become arbitrary patches. Keep dog anatomy and four-paw baseline.
- `dog_black_tan`: black coat can hide legs and muzzle. Require readable tan points and paw separation.
- `dog_white`: white fur can disappear on light preview. Require soft off-white edge detail without shadow.
- `dog_sable`: sable gradient may flicker across frames. Require stable darker back/muzzle markings.
- `dog_gray`: gray dog may resemble a cat if muzzle/tail proportions drift. Preserve dog head and paw shape.
- `cat_brown_tabby`: tabby stripes can shimmer frame to frame. Require stable stripe placement and readable tail.
- `cat_black`: silhouette can collapse into a single blob. Require subtle highlights and separated legs/tail.

## Parent QA Checklist

- Confirm 10 variants x 62 frames exist under `motion-source/`.
- Confirm `source-truth/` has exactly one source PNG per assigned variant.
- Confirm previews exist on both light and dark backgrounds.
- Confirm each PNG has transparent background and non-empty alpha.
- Confirm frame bounds fit inside 96x64 after importer normalization.
- Confirm `rabbit-hop`, `dog-trot`, and `cat-stalk` remain visually distinct.
