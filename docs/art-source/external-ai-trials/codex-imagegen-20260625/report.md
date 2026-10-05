# Codex Built-in ImageGen Trial Report - 2026-06-25

## Scope

This thread used only Codex built-in `image_gen` for review-only AnimalsDesktop art trials.

Write boundary:

- Wrote only under `docs/art-source/external-ai-trials/codex-imagegen-20260625/`.
- Did not write to `accepted-frames`.
- Did not edit catalog/runtime/importer code.
- Did not publish, tag, push, or touch GitHub Pages/Releases.

## Summary Decision

| Variant | 16-cell trial | 62-frame expansion | Parent integration recommendation |
| --- | --- | --- | --- |
| `hamster_golden_syrian` | Partial pass | Not attempted | Use cells 00-13/15 only as visual references. Do not integrate the sheet directly. |
| `macaroni_mouse_tan` | Best 16-cell pass | Rejected | Best next candidate for a stricter second 16-cell pass or individual frame generation. Do not integrate the 62 sheet. |
| `sugar_glider_gray` | Strong 16-cell pass | Rejected | Good prompt direction, especially grounded membrane stretch. Do not integrate the 62 sheet. |
| `rabbit_chestnut_agouti` | Partial pass | Not attempted | Use idle/rest/groom cells as references only. Hop/reaction cells need stricter ear margin and no fragment prompt. |

Main finding: the 16-cell prompts can produce useful species reads and small motion ideas. The 62-frame sheet prompts are not reliable because generated cell boundaries drift; equal-cell slicing produced row-wrapped, cropped, or duplicate fragments even though mechanical alpha checks passed.

## Generated Outputs

### 16-cell trials

| Variant | Built-in ImageGen source | Workspace raw copy | Contact sheet | QA |
| --- | --- | --- | --- | --- |
| `hamster_golden_syrian` | `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/ig_05d95f45a036f1d8016a3c5ae6eeec8191854554eee855cb30.png` | `raw/hamster_golden_syrian-16sheet-raw.png` | `contacts/hamster_golden_syrian-trial01-16cells-contact.png` | `qa/hamster_golden_syrian-trial01-qa.json` |
| `macaroni_mouse_tan` | `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/ig_05d95f45a036f1d8016a3c5b21bbd0819196745edcc2c3e32a.png` | `raw/macaroni_mouse_tan-16sheet-raw.png` | `contacts/macaroni_mouse_tan-trial01-16cells-contact.png` | `qa/macaroni_mouse_tan-trial01-qa.json` |
| `sugar_glider_gray` | `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/ig_05d95f45a036f1d8016a3c5b54493481919817cc913b76c121.png` | `raw/sugar_glider_gray-16sheet-raw.png` | `contacts/sugar_glider_gray-trial01-16cells-contact.png` | `qa/sugar_glider_gray-trial01-qa.json` |
| `rabbit_chestnut_agouti` | `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/ig_05d95f45a036f1d8016a3c5b8ac9008191b65411f26637f6fb.png` | `raw/rabbit_chestnut_agouti-16sheet-raw.png` | `contacts/rabbit_chestnut_agouti-trial01-16cells-contact.png` | `qa/rabbit_chestnut_agouti-trial01-qa.json` |

All four 16-cell trials split into 96x64 review frames with:

- empty cells: 0
- edge-alpha cells: 0
- outputs under `processed/<variant>/trial01/`

### 62-frame expansion attempts

| Variant | Built-in ImageGen source | Workspace raw copy | Contacts | Decision |
| --- | --- | --- | --- | --- |
| `macaroni_mouse_tan` | `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/ig_0945f3038facb811016a3c5c8268548191b2d126c353dd13ce.png` | `raw/macaroni_mouse_tan-set00-62sheet-raw.png` | `contacts/macaroni_mouse_tan-set00-62-trial01-light.png`, `-dark.png`, `-checker.png` | Rejected: equal-cell slicing produced row-wrapped/cropped animal fragments. |
| `sugar_glider_gray` | `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/ig_0945f3038facb811016a3c5cf1b1fc8191a931a9ebc3e46a8f.png` | `raw/sugar_glider_gray-set00-62sheet-raw.png` | `contacts/sugar_glider_gray-set00-62-trial01-light.png`, `-dark.png`, `-checker.png` | Rejected: equal-cell slicing produced row-wrapped/cropped animal fragments. |

Both 62 attempts produced non-empty alpha and no edge-alpha after normalization, but this is not sufficient. Visual contact sheets show that the generated sheet layout does not preserve real cell boundaries.

## Per-variant Visual Review

### `hamster_golden_syrian`

Decision: partial pass, reference-only.

Useful traits:

- Golden Syrian color read is strong.
- Cheek-forward compact body is mostly stable.
- Idle, walk, forage, eat, and groom concepts are readable.

Problems:

- Many cells are too similar for a real motion set.
- Cell 14 has airborne/reaction drift and less stable feet.
- Background was not truly flat chroma; local removal worked, but this remains a prompt failure.

Next prompt direction:

```text
Create a 4x4 review-only sprite sheet for hamster_golden_syrian with exact equal cells and wide pure #00ff00 gutters. Use the same compact golden Syrian hamster in every cell: round cheek-forward body, cream belly patches, blunt muzzle, rounded ears, tiny paws, and visually absent tail. Keep all poses grounded and wheel-safe. Limit differences to small idle, walk, forage, eat, groom, and alert changes. Do not create airborne reaction poses. No shadows, no floor, no effect marks, no debris, no labels, no dividers.
```

### `macaroni_mouse_tan`

Decision: best 16-cell pass, but 62 expansion rejected.

Useful traits:

- Thick pale club-shaped tail is clear in most cells.
- Low scurry and forage poses are useful.
- It reads more like a fat-tailed gerbil than a hamster.

Problems:

- A few cells have tiny side fragments after chroma removal.
- The 62 sheet failed because the generated layout did not preserve equal cells.
- Some 62 rows contain cropped duplicate animal parts.

Next prompt direction:

```text
Create a 4x4 review-only sprite sheet for macaroni_mouse_tan with exact equal cells and wide pure #00ff00 gutters. One complete right-facing tan fat-tailed gerbil per cell, same individual, compact sandy body, cream underside, large dark eye, small rounded ears, tiny paws, and a thick pale club-shaped tail that remains fat and low behind the body. Keep every pose grounded and fully inside its cell. Avoid thin mouse tail, hamster body, row wrapping, repeated partial animals, shadows, floor, effect marks, debris, labels, and dividers.
```

Recommended parent action:

- Use `contacts/macaroni_mouse_tan-trial01-16cells-contact.png` as the best visual prompt reference.
- Prefer individual frame generation or 8/16-cell batches, not 62-cell sheets.

### `sugar_glider_gray`

Decision: strong 16-cell pass, but 62 expansion rejected.

Useful traits:

- Large-eye sugar glider read is strong.
- Grounded membrane stretch in the 16-cell trial is useful.
- Low skitter body shape is better than the older primitive source direction.

Problems:

- Some cells contain small background debris and contact-shadow residue.
- The 62 sheet failed because generated cell boundaries drifted.
- It can become squirrel-like when membrane/tail cues are weak.

Next prompt direction:

```text
Create a 4x4 review-only sprite sheet for sugar_glider_gray with exact equal cells and wide pure #00ff00 gutters. One complete right-facing gray sugar glider per cell, same individual, gray back, cream underside, bold dark eye, small rounded ears, pointed pink nose, long soft tail, and visible natural patagium membrane as a grounded side fold. Include low skitter, forage, groom, alert, and one grounded membrane-stretch pose. Do not make bat wings or full flying poses. No shadows, no floor, no effect marks, no debris, no labels, no dividers.
```

Recommended parent action:

- Use `contacts/sugar_glider_gray-trial01-16cells-contact.png` as a prompt reference.
- Keep cell 14's grounded membrane-stretch idea, but regenerate as an individual frame or smaller controlled batch.

### `rabbit_chestnut_agouti`

Decision: partial pass, reference-only.

Useful traits:

- Chestnut agouti color and rabbit silhouette read well.
- Idle, rest, groom, and forage cells are usable references.
- Ear shape is mostly coherent when not in hop/reaction cells.

Problems:

- Hop/reaction cells have clipped or floating fragments.
- Ears remain top-margin risk.
- Some poses over-extend or lose stable contact timing.

Next prompt direction:

```text
Create a 4x4 review-only sprite sheet for rabbit_chestnut_agouti with exact equal cells and wide pure #00ff00 gutters. One complete right-facing chestnut agouti rabbit per cell, same individual, long upright ears with large top margin, compact body, strong hindquarters, tiny white tail, visible whiskers, and small forepaws. Keep hop poses low and grounded; do not create airborne leaps. Avoid cropped ears, floating feet, detached fragments, shadows, floor, effect marks, debris, labels, and dividers.
```

## Exact Prompt Pattern Used

The prompts that worked best had this structure:

```text
Use case: stylized-concept
Asset type: AnimalsDesktop review-only sprite source sheet, 4 columns by 4 rows.
Primary request: Create a 16-cell pose test sheet for <variant>. The image must be one clean 4x4 sprite sheet: 16 equal rectangular cells, no visible grid lines, no labels, no numbers, no text. Every cell must have a perfectly flat pure #00ff00 chroma-key background, including gutters, with no gradient, no texture, no shadow, no floor, no checker pattern.
Subject: one complete right-facing side-view <species description> in each cell, same individual source family across all cells.
Style/medium: polished 2D desktop pet source art, soft natural fur detail, clean silhouette, downsample-friendly, like high-quality painted sprite source art rather than a photo.
Composition/framing: each cell contains exactly one animal, fully visible with generous padding; consistent camera distance, body scale, contact baseline, facing direction, lighting, outline softness, and anatomy across the sheet.
Pose cells in order left-to-right, top-to-bottom: neutral idle stance; tiny breathing shift; ear/whisker adjustment; slight weight shift; walk phase front paw forward; walk phase forefoot contact; walk phase rear foot forward; walk return; faster low scurry pose; scurry alternate foot phase; nose-down forage; compact eating pose; face groom with paw; alert rest pose; small species-safe reaction; return to neutral.
Avoid: text, labels, numbers, borders, visible dividers, scenery, props, costumes, multiple animals in one cell, cropped anatomy, disconnected stray pixels, cast shadow, floor band, white background, gray background, human-like pose, giant eye, smeared face, broken paws.
```

## Verification Run

Local processing performed:

- Copied all built-in ImageGen outputs from `/Users/kyota/.codex/generated_images/...` into the workspace `raw/` directory.
- Removed accidental AppleDouble `._*` files only inside this trial directory.
- Split 4x4 sheets into 16 cells, removed chroma background locally, normalized cells to 96x64, and wrote contacts/QA JSON.
- Split the two 8x8 62-frame attempts into 62 cells, removed chroma background locally, normalized cells to 96x64, and wrote contacts/QA JSON.

Mechanical QA result:

- 16-cell trials: all 4 variants have 0 empty cells and 0 edge-alpha cells.
- 62-frame attempts: both attempted variants have 0 empty frames and 0 edge-alpha frames, but visual QA rejects them.

## Final Recommendation

Do not promote any output from this folder directly. The best useful artifacts for the parent thread are:

1. `contacts/macaroni_mouse_tan-trial01-16cells-contact.png`
2. `contacts/sugar_glider_gray-trial01-16cells-contact.png`
3. `contacts/hamster_golden_syrian-trial01-16cells-contact.png`
4. `contacts/rabbit_chestnut_agouti-trial01-16cells-contact.png`

For the next parent integration loop, use 8-cell or 16-cell batches only, or generate individual frames. Avoid 62-cell ImageGen sheets unless there is a stronger external layout-control method, because equal-grid slicing is not reliable enough for accepted-frame production.
