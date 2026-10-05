# Cohort 05 Known Risks

## Blocking Status

ImageGen PNG files were not produced by this delegated thread because no file-saving ImageGen path was available here. Parent integration must treat this cohort as `imagegen_prompts_ready`, not as accepted source art.

## Variant Risks

- `hamster_black_banded`: band may drift across motion frames if generated one frame at a time. Use a reference sheet or image-to-image consistency pass.
- `hamster_white`: white coat may disappear on light preview. Require outline and warm shading without adding a background shadow.
- `ferret_albino`: pale body and ruby eyes may be hard to read at 96x64. Verify on both preview backgrounds.
- `ferret_sable` and `ferret_champagne`: long tails are likely to be cropped. Require extra horizontal padding.
- `guinea_pig_tricolor`: patch layout may change between frames. Lock patch positions before accepting a 62-frame sequence.
- `guinea_pig_black`: black coat can lose feet and face details on dark preview. Require controlled highlights.

## Motion Risks

- Hamster frames must use quick short steps and low scurry, not generic vertical bobbing.
- Ferret frames must use a long low body wave and tail follow-through, not a degu walk cycle.
- Guinea pig frames must use a rounded side-to-side waddle, not rabbit hopping.
- Grooming frames must stay animal-like and avoid upright or human hand-like poses.

## Parent Integration Requirements

- Do not promote any cohort-05 variant from `prototype_only` until PNG source-truth, 62 motion frames, and both previews exist.
- After files exist, run visual alpha/crop checks before changing importer manifest or catalog source status.
- Then regenerate sprites from the accepted source frames and compare deterministic importer output.
