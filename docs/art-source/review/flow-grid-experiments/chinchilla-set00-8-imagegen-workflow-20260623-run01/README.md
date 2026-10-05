# Chinchilla 8-frame ImageGen workflow run 01

Date: 2026-06-23

Purpose: run the proposed small-batch workflow once for a standard gray chinchilla walk cycle:

1. Generate an 8-frame invisible `2x4` sheet with ImageGen.
2. Generate a stricter layout-correction ImageGen pass.
3. Diagnose raw fixed-cell parse, green-normalized parse, and cell-recentered parse.
4. Keep everything review-only; do not promote any output to `accepted-frames`.

## Inputs

- Generator: Codex ImageGen.
- Animal: right-facing standard gray chinchilla.
- Layout request: `2048x1024`, invisible `2x4`, `512x512` cells, pure `#00ff00` background.
- Native ImageGen limitation: this run did not upload the first sheet into the second pass. The second pass is an independent stricter prompt, not a true edit of stage 1.

Original generated files remain in:

- `C:\Users\user\.codex\generated_images\019eef28-97a5-7513-a530-7b4501e063f3\ig_035b1848a80b6b20016a3a3ef95e648191ba07a86c8f3f17fb.png`
- `C:\Users\user\.codex\generated_images\019eef28-97a5-7513-a530-7b4501e063f3\ig_018257b07f382183016a3a3f4561588191b2e64fe643650fae.png`

## Outputs

- `imagegen-stage1-eight-walk.png`
- `imagegen-stage2-layout-correction.png`
- `imagegen-stage1-eight-walk-greenfixed.png`
- `imagegen-stage2-layout-correction-greenfixed.png`
- `imagegen-stage1-eight-walk-cell-recentered-2048x1024.png`
- `imagegen-stage2-layout-correction-cell-recentered-2048x1024.png`
- `imagegen-chinchilla-8-workflow-contact.png`
- `workflow-diagnostics.json`

## Parse Results

| Source | Raw fixed-cell parse | Green-normalized parse | Cell-recentered parse |
| --- | ---: | ---: | ---: |
| Stage 1 eight walk | 0 / 8 | 6 / 8 | 6 / 8 |
| Stage 2 layout correction | 0 / 8 | 7 / 8 | 7 / 8 |

Failure details:

- Stage 1 green-normalized rejects: `frame-01` and `frame-03`, both `MAT_PINHOLE`.
- Stage 2 green-normalized reject: `frame-04`, one `1px` `MAT_PINHOLE`.
- Both generated sheets had 8 detected animal components with no large extra debris.
- Raw parse failed because the green background was not pure enough at parser boundaries. Diagnostic green normalization is not an accepted cleanup path.

## Visual Review

The stage 2 sheet is the best mechanical result so far. It keeps eight animals separated, centered, and extractable after green normalization except for a single 1px pinhole.

Remaining blockers:

- Art is still more small realistic illustration than accepted DeguDesktop-equivalent sprite style.
- Motion deltas are readable but weak; several frames are close variants.
- The output still relies on diagnostic green normalization and still has a pinhole reject.

## Decision

This workflow is promising but not production-ready. The best current path for the next loop is:

1. Keep the 8-frame `2x4` batch size.
2. Use the stage 2 prompt as the base.
3. Add stronger wording for flatter 2D sprite style and visible gait deltas.
4. If using ChatGPT Pro or Flow, upload the best generated sheet as an explicit reference for a true layout-correction pass, after parent approval for that upload.
5. Continue to parse under `review/` only. Do not promote batch cuts into `accepted-frames`.
