# Chinchilla Flow Nano Banana Pro 16-cell test

Date: 2026-06-23

## Setup

- Service: Google Flow
- Model: Nano Banana Pro
- Output: square, x4, 2K download
- Inputs added to prompt:
  - `style-anchors.png`
  - `grid-seed.png`
- Prompt entry note: browser-level paste/type did not update Flow's internal prompt state. OS-level virtual key input activated the prompt and enabled generation.

## Outputs

- Original downloads:
  - `flow-nanobanana-pro-2ref-output-01.jpeg`
  - `flow-nanobanana-pro-2ref-output-02.jpeg`
  - `flow-nanobanana-pro-2ref-output-03.jpeg`
  - `flow-nanobanana-pro-2ref-output-04.jpeg`
- PNG conversions:
  - `flow-nanobanana-pro-2ref-output-01.png` through `-04.png`
- Green-normalized PNGs:
  - `flow-nanobanana-pro-2ref-output-01-greenfixed.png` through `-04-greenfixed.png`

## Parse Results

- Original `batch.json` fixed guide manifest: all candidates parsed `0/16`.
- Actual 512px grid manifest without green cleanup: all candidates rejected due boundary chroma noise.
- Actual 512px grid with green cleanup:
  - candidate 01: `11/16`
  - candidate 02: `9/16`
  - candidate 03: `7/16`
  - candidate 04: `9/16`

Remaining failures are mainly disconnected components and chroma pinholes. The generated sheets are visually consistent, but most cells duplicate the same pose, so this is not accepted as a motion-production workflow.

## Variant A: Explicit Limb Delta Prompt

Follow-up date: 2026-06-23

- Prompt file: `prompt-variant-a-short.txt`
- Prompt notes and next variants: `prompt-variants.md`
- Consultation packet: `failure-consultation-variant-a.md`
- Cache-recovered outputs:
  - `flow-nanobanana-pro-variant-a-cache-141219.jpeg`
  - `flow-nanobanana-pro-variant-a-cache-141429.jpeg`

Result:

- Pose variation improved versus the first Nano Banana Pro prompt.
- The visible `grid-seed.png` guide was still reproduced as cell lines, so raw output is rejected as `GUIDE_INK`.
- Diagnostic green normalization was used only to inspect animal content; it is not an acceptance cleanup path.
- `141219` green-normalized parse: `3/16`, mostly rejected for `MAT_PINHOLE` or detached alpha.
- `141429` green-normalized parse: `0/16`, rejected for boundary/grid residue.

Current conclusion:

Do not retry the same visible grid seed with only pose wording changes. The next prompt experiment should change the guide strategy: style-anchor-only invisible grid, a smaller 4-frame strip, or independent images.

## Variant G/I: 8-frame and two-stage layout correction

Follow-up date: 2026-06-23

Flow / Nano Banana Pro:

- Prompt file: `prompt-variant-g-eight-frame-short.txt`
- Input: `style-anchors.png` only, no `grid-seed.png`
- Result: `flow-nanobanana-pro-variant-g-eight-cache-145924.png`
- Outcome: ignored the requested 8-frame invisible `2x4` layout and returned a visible-line 4x4 sheet with 16 animals.
- Diagnostic parse: `0/8`

ChatGPT ImageGen:

- Result: `imagegen-variant-g-eight-textonly.png`
- Outcome: produced 8 separated animals in a clean 2x4 arrangement without visible grid lines.
- Fixed-cell parse: `0/8`, because the generated positions did not align to exact equal cells.
- Component diagnostic: 8 animal components.

Two-stage layout correction:

- Prompt file: `prompt-variant-i-layout-correction.txt`
- Result: `imagegen-variant-i-eight-layout-correction.png`
- Diagnostic green-normalized parse: `5/8`
- Remaining failure: three `MAT_PINHOLE` cells.

Conclusion:

The user's two-stage correction idea is promising. ImageGen can draw a coherent 8-frame sheet, and a strict layout-correction pass improves extraction. The remaining blocker is matte/pinhole quality and source-family verification, not the basic ability to produce eight frames.

## Review Images

- `flow-nanobanana-pro-2ref-contact.png`: four Flow outputs.
- `parsed-pro-01-greenfixed-contact.png`: best parsed candidate preview.
- `flow-nanobanana-pro-variant-a-cache-contact.png`: Variant A recovered outputs.
- `flow-nanobanana-pro-variant-a-141219-qa-contact.png`: Variant A diagnostic parse view.
- `flow-nanobanana-pro-variant-a-141429-qa-contact.png`: Variant A diagnostic parse view.
- `variant-g-eight-flow-vs-imagegen-contact.png`: Flow Pro versus ChatGPT ImageGen 8-frame comparison.
- `imagegen-eight-two-stage-correction-contact.png`: ImageGen 8-frame two-stage correction comparison.
