# Sugar Glider Gray Rescue Frames 54-61 - 2026-06-25

Dedicated output root:

`docs/art-source/one-frame-method-fullrun-20260625-rescue/sugar-glider-54-61/`

## Scope

- Species: `sugar_glider_gray` / momonga
- Frames: `54` through `61` only
- Method: Codex built-in ImageGen, one frame at a time, retry failed frame numbers only
- Background: perfectly flat `#00ff00` chroma key, then local alpha removal
- Runtime/accepted-frame mutation: none
- Main fullrun mutation: none

## Artifact Layout

- Raw built-in ImageGen copies: `raw/frame-NN-attempt-MM-raw.png`
- Alpha converted source: `alpha/frame-NN-attempt-MM-alpha.png`
- Normalized candidate frame: `frames/frame-NN.png`
- Attempt-specific normalized frame: `frames/frame-NN-attempt-MM.png`
- Per-frame contact preview: `contacts/frame-NN-attempt-MM-contact.png`
- 8-frame contact sheet: `contacts/sugar-glider-54-61-rescue-contact.png`
- QA JSON: `qa/frame-NN-attempt-MM-qa.json`
- Prompt text: `prompts/frame-NN-attempt-MM-prompt.txt`
- Manifest: `manifest.csv`

## Mechanical QA Summary

| Frame | Group | Attempt | Alpha px | Edge alpha | Green px | Decision | Notes |
| ---: | --- | ---: | ---: | ---: | ---: | --- | --- |
| 54 | face_groom | 3 | 1462 | 0 | 0 | mechanical-pass | prompt:prompts/frame-54-attempt-03-prompt.txt; visual-precheck: paw touches ear, species identity acceptable |
| 55 | face_groom | 2 | 1399 | 0 | 0 | mechanical-pass | prompt:prompts/frame-55-attempt-02-prompt.txt; visual-precheck pending contact review |
| 56 | membrane_reaction | 3 | 1948 | 0 | 0 | mechanical-pass | prompt:prompts/frame-56-attempt-03-prompt.txt; visual-precheck pending contact review |
| 57 | membrane_reaction | 3 | 2007 | 0 | 0 | mechanical-pass | prompt:prompts/frame-57-attempt-03-prompt.txt; visual-precheck pending contact review |
| 58 | membrane_reaction | 2 | 2068 | 0 | 0 | mechanical-pass | prompt:prompts/frame-58-attempt-02-prompt.txt; visual-precheck pending contact review |
| 59 | membrane_reaction | 1 | 1682 | 0 | 0 | mechanical-pass | prompt:prompts/frame-59-attempt-01-prompt.txt |
| 60 | membrane_reaction | 2 | 1707 | 0 | 0 | mechanical-pass | prompt:prompts/frame-60-attempt-02-prompt.txt; visual-precheck: half-spread membrane side view readable |
| 61 | membrane_reaction | 2 | 2026 | 0 | 0 | mechanical-pass | prompt:prompts/frame-61-attempt-02-prompt.txt; visual-precheck: settle after reaction readable |

## Completion Check

- Candidate canonical frames present and mechanically clean: `8/8`
- Completed rescue range: `yes`

## Remaining Risk

- Mechanical QA catches size, alpha, edge contact, empty alpha, and green residue only.
- Parent thread still needs final visual comparison against the main fullrun for pose readability, anatomy, species identity, membrane clarity, and style consistency.
