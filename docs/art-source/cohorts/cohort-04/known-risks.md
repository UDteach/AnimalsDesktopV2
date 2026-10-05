# Cohort 04 Known Risks

## Blocking Status

- Local PNG output from ImageGen was not available in this child thread, so this cohort is not visually complete.
- No `source-truth/*.png`, `motion-source/*/frame-XX.png`, or `preview/*.png` files were emitted here.
- Parent integration must treat this cohort as prompt-ready only.

## Species Risks

- Cat variants may drift into long-haired or breed-specific silhouettes. Keep them domestic shorthair style unless the catalog later changes breed labels.
- White and cream cats can disappear on light backgrounds. Require dark outline or shaded edges before acceptance.
- Orange tabby stripes can become noisy at 96x64. Preserve broad, readable stripes instead of many fine marks.
- Gecko toes and tail tips are high crop risks. Add extra transparent margin around all gecko frames.
- Gecko morph colors must not turn into fantasy lizards. Keep leopard, tangerine, blizzard, and albino morphs recognizable and natural.
- Hamster frames can become too round to read as walking. Ensure paws and head direction remain visible in the 96x64 preview.

## Motion Risks

- `cat-stalk` must stay quiet and grounded; do not use rabbit-hop, degu wheel, or exaggerated vertical bounce.
- `gecko-crawl` must keep the belly low with no hopping, rearing, or upright turn frames.
- `small-rodent-scurry` can use a small bob, but avoid degu-like wheel/run poses for the source motion set.
- Turn frames should preserve the same scale and baseline as walk frames so runtime mirroring and direction changes remain seamless.

## QA Needed By Parent

- Verify every source and motion frame has alpha transparency and non-empty visible pixels.
- Verify each 96x64 preview is readable on light and dark backgrounds.
- Verify all 62 frames preserve contact points and do not crop ears, whiskers, paws, toes, or tails.
- Verify generated filenames exactly match the cohort assignment variant IDs.
