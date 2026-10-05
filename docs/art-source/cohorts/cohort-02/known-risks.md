# Cohort 02 Known Risks

Status: prompt-ready only; no accepted PNG files were generated in this delegated pass.

## Cohort-Level Risks

- Source-truth PNGs and 62-frame motion source sets are missing until ImageGen is run and reviewed.
- Coat consistency can drift across 62 generated frames, especially for cream pied degu, white mosaic chinchilla, and chestnut agouti rabbit.
- Chinchilla frames must preserve dense round body, large ears, long whiskers, and fluffy tail; ImageGen may crop whiskers or tail at 96x64 downscale.
- Macaroni mouse / fat-tailed gerbil identity is fragile; generated output may become a generic mouse unless the thick club-shaped tail is explicit.
- Rabbit frames must avoid rodent scurry anatomy and must preserve hind-leg hop mechanics.
- Transparent background must be verified after generation; prompts alone do not prove alpha containment.

## Variant-Specific Risks

| Variant | Risk | Parent Review |
| --- | --- | --- |
| `cream_pied` | Pied patch pattern may flicker between frames. | Check coat-map consistency across all 62 frames. |
| `chinchilla_standard_gray` | Tail or whiskers may be cropped. | Inspect full alpha bounds and 96x64 previews. |
| `chinchilla_beige` | Beige may collapse into hamster-like silhouette. | Compare large ears, plush round body, and tail against chinchilla reference. |
| `chinchilla_ebony` | Dark coat can lose readable body outline. | Check dark preview and edge contrast. |
| `chinchilla_white_mosaic` | Mosaic patches may look like costume markings. | Require natural fur patch distribution and no hard graphic shapes. |
| `macaroni_mouse_tan` | Tail may be drawn long and thin instead of fat-tailed. | Require thick club-shaped tail in every frame. |
| `macaroni_mouse_gray` | Gray coat may read as generic house mouse. | Check short muzzle and fat-tailed gerbil silhouette. |
| `macaroni_mouse_cream` | Pale coat may disappear on light preview. | Check light preview contrast and outline. |
| `rabbit_chestnut_agouti` | Hop frames may crop ears at top. | Check frame bounds before import. |
| `rabbit_black` | Dark coat can become a featureless blob. | Check dark and light previews for eye/body separation. |

## Required Follow-Up

1. Run ImageGen for `source-truth/<variant-id>.png` for all 10 variants.
2. Generate or manually curate 62 motion source frames for each variant using `cohort-report.md`.
3. Create light/dark 96x64 previews.
4. Parent thread should reject any frame with multiple animals, shadows, scenery, text, cropped anatomy, inconsistent baseline, or non-transparent background.
