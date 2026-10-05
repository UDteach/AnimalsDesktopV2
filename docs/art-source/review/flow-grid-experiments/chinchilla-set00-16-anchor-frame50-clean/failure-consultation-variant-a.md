# Failure Consultation Packet: Variant A

Date: 2026-06-23

Animal / coat: standard gray chinchilla

Batch size: 16 cells, 4x4 sheet

Service / model: Google Flow, Nano Banana Pro

Anchors used:

- `style-anchors.png`
- `grid-seed.png`

Prompt:

- `prompt-variant-a-short.txt`
- Full variant set: `prompt-variants.md`

Parser result:

- Raw 1024 cache output `141219`: rejected `17` errors because visible guide/grid lines touched cell boundary guards.
- Raw 1024 cache output `141429`: rejected `17` errors because visible guide/grid lines touched cell boundary guards.
- Diagnostic green-normalized output `141219`: `parsed=3 rejected=13`.
- Diagnostic green-normalized output `141429`: `parsed=0 rejected=16`.

Visual result:

- `flow-nanobanana-pro-variant-a-cache-contact.png`
- `flow-nanobanana-pro-variant-a-141219-qa-contact.png`
- `flow-nanobanana-pro-variant-a-141429-qa-contact.png`

Failure type:

- `GUIDE_INK`: visible grid lines reproduced in generated sheet.
- `MAT_PINHOLE`: green/chroma pinholes inside animal after diagnostic normalization.
- `MAT_ROUGH`: edge/matte reliability is not good enough for release extraction.
- `LAY_CROSS`: some cells, especially low scurry frames, approach or cross fixed boundaries.
- `POSE_DUP`: improved versus previous Pro run, but still not enough to accept the workflow.

Question for prompt review:

What should change before the next attempt: guide image, anchor set, prompt wording, batch size, or output format? The goal is still one extracted frame per PNG after local QA. Do not recommend manual repainting as the acceptance path.

Current local recommendation:

1. Do not retry `style-anchors.png + grid-seed.png` with only pose wording changes.
2. Try a style-anchor-only invisible 4x4 prompt, accepting the risk of weaker placement.
3. If guide ink or pinholes repeat, downgrade to a 4-frame strip or four independent images.
4. Do not attempt 31 or 62 cells until 16 cells have zero guide ink and zero severe matte defects.
