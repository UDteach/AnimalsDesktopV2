# Hamster Golden Syrian Rescue Frames 32-61 - 2026-06-25

Dedicated output root:

`docs/art-source/one-frame-method-fullrun-20260625-rescue/hamster-32-61/`

## Scope

- Species: `hamster_golden_syrian`
- Frames: `32` through `61` only
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
- 30-frame contact sheet: `contacts/hamster-32-61-rescue-contact.png`
- QA JSON: `qa/frame-NN-attempt-MM-qa.json`
- Prompt text: `prompts/frame-NN-attempt-MM-prompt.txt`
- Manifest: `manifest.csv`

## Mechanical QA Summary

| Frame | Group | Attempt | Alpha px | Edge alpha | Green px | Decision | Notes |
| ---: | --- | ---: | ---: | ---: | ---: | --- | --- |
| 32 | turn | 1 | 1700 | 0 | 0 | mechanical-pass | prompt:prompts/frame-32-attempt-01-prompt.txt |
| 33 | turn | 2 | 1773 | 0 | 0 | mechanical-pass | prompt:prompts/frame-33-attempt-02-prompt.txt |
| 34 | turn | 3 | 1669 | 0 | 0 | mechanical-pass | prompt:prompts/frame-34-attempt-03-prompt.txt |
| 35 | turn | 1 | 1487 | 0 | 0 | mechanical-pass | prompt:prompts/frame-35-attempt-01-prompt.txt |
| 36 | turn | 1 | 2746 | 0 | 0 | mechanical-pass | prompt:prompts/frame-36-attempt-01-prompt.txt |
| 37 | turn | 1 | 1929 | 0 | 0 | mechanical-pass | prompt:prompts/frame-37-attempt-01-prompt.txt |
| 38 | turn | 1 | 1558 | 0 | 0 | mechanical-pass | prompt:prompts/frame-38-attempt-01-prompt.txt |
| 39 | turn | 1 | 2637 | 0 | 0 | mechanical-pass | prompt:prompts/frame-39-attempt-01-prompt.txt |
| 40 | chew | 1 | 2096 | 0 | 0 | mechanical-pass | prompt:prompts/frame-40-attempt-01-prompt.txt |
| 41 | chew | 1 | 2006 | 0 | 0 | mechanical-pass | prompt:prompts/frame-41-attempt-01-prompt.txt |
| 42 | chew | 1 | 2192 | 0 | 0 | mechanical-pass | prompt:prompts/frame-42-attempt-01-prompt.txt |
| 43 | chew | 1 | 2103 | 0 | 0 | mechanical-pass | prompt:prompts/frame-43-attempt-01-prompt.txt |
| 44 | ground_check | 2 | 2785 | 0 | 0 | mechanical-pass | prompt:prompts/frame-44-attempt-02-prompt.txt |
| 45 | ground_check | 1 | 3040 | 0 | 0 | mechanical-pass | prompt:prompts/frame-45-attempt-01-prompt.txt |
| 46 | ground_check | 1 | 2475 | 0 | 0 | mechanical-pass | prompt:prompts/frame-46-attempt-01-prompt.txt |
| 47 | ground_check | 1 | 2968 | 0 | 0 | mechanical-pass | prompt:prompts/frame-47-attempt-01-prompt.txt |
| 48 | alert_rest | 1 | 2687 | 0 | 0 | mechanical-pass | prompt:prompts/frame-48-attempt-01-prompt.txt |
| 49 | alert_rest | 1 | 2685 | 0 | 0 | mechanical-pass | prompt:prompts/frame-49-attempt-01-prompt.txt |
| 50 | alert_rest | 1 | 1274 | 0 | 0 | mechanical-pass | prompt:prompts/frame-50-attempt-01-prompt.txt |
| 51 | alert_rest | 1 | 2428 | 0 | 0 | mechanical-pass | prompt:prompts/frame-51-attempt-01-prompt.txt |
| 52 | face_groom | 1 | 1817 | 0 | 0 | mechanical-pass | prompt:prompts/frame-52-attempt-01-prompt.txt |
| 53 | face_groom | 1 | 1723 | 0 | 0 | mechanical-pass | prompt:prompts/frame-53-attempt-01-prompt.txt |
| 54 | face_groom | 1 | 1739 | 0 | 0 | mechanical-pass | prompt:prompts/frame-54-attempt-01-prompt.txt |
| 55 | face_groom | 1 | 1788 | 0 | 0 | mechanical-pass | prompt:prompts/frame-55-attempt-01-prompt.txt |
| 56 | reaction | 2 | 1911 | 0 | 0 | mechanical-pass | prompt:prompts/frame-56-attempt-02-prompt.txt |
| 57 | reaction | 1 | 1809 | 0 | 0 | mechanical-pass | prompt:prompts/frame-57-attempt-01-prompt.txt |
| 58 | reaction | 1 | 1269 | 0 | 0 | mechanical-pass | prompt:prompts/frame-58-attempt-01-prompt.txt |
| 59 | reaction | 1 | 3010 | 0 | 0 | mechanical-pass | prompt:prompts/frame-59-attempt-01-prompt.txt |
| 60 | reaction | 1 | 2896 | 0 | 0 | mechanical-pass | prompt:prompts/frame-60-attempt-01-prompt.txt |
| 61 | reaction | 1 | 2488 | 0 | 0 | mechanical-pass | prompt:prompts/frame-61-attempt-01-prompt.txt |

## Completion Check

- Candidate canonical frames present and mechanically clean: `30/30`
- Completed rescue range: `yes`

## Remaining Risk

- Mechanical QA catches size, alpha, edge contact, empty alpha, and green residue only.
- Parent thread still needs final visual comparison against the main fullrun for pose readability, anatomy, species identity, and style consistency.
