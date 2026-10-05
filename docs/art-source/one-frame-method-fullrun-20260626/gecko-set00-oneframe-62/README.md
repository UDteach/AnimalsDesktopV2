# Gecko set00 one-frame 62-frame run

This run uses the accepted AnimalsDesktop/DeguDesktop workflow: generate one frame at a time, keep raw outputs, remove chroma key locally, normalize to 96x64 transparent PNG, run QA, then visually review motion slots before any promotion.

Target:

- family: `gecko`
- variant: gray-brown gecko
- set: `set00`
- frame count: 62
- method: built-in ImageGen, one call per frame
- source reference: `../gecko/frames/frame-00-attempt-01-review.png`

Directories:

- `raw/`: built-in ImageGen outputs copied from `$CODEX_HOME/generated_images`
- `alpha/`: chroma-key removed intermediate PNGs
- `frames/`: 96x64 candidate frames
- `contact/`: review sheets
- `prompts/`: base prompt and pose table
- `qa/`: JSON/TSV/Markdown QA records
- `metadata/`: generation markers and image mapping

Acceptance is not automatic. A complete set requires machine checks and visual review by motion slot.
