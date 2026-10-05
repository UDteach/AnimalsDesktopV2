# Rabbit Chestnut Agouti Rescue Frames 40-61 - 2026-06-25

Dedicated output root:

`docs/art-source/one-frame-method-fullrun-20260625-rescue/rabbit-40-61/`

## Scope

- Species: `rabbit_chestnut_agouti`
- Frames: `40` through `61` only
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
- 22-frame contact sheet: `contacts/rabbit-40-61-rescue-contact.png`
- QA JSON: `qa/frame-NN-attempt-MM-qa.json`
- Prompt text: `prompts/frame-NN-attempt-MM-prompt.txt`
- Manifest: `manifest.csv`

## Mechanical QA Summary

| Frame | Group | Attempt | Alpha px | Edge alpha | Green px | Decision | Notes |
| ---: | --- | ---: | ---: | ---: | ---: | --- | --- |
| 40 | chew | 2 | 1838 | 0 | 0 | mechanical-pass | retry: visible white tail and right-facing chew |
| 41 | chew | 1 | 1912 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 42 | chew | 1 | 2048 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 43 | chew | 2 | 1905 | 0 | 0 | mechanical-pass | retry: remove green residue |
| 44 | ground_check | 1 | 2854 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 45 | ground_check | 1 | 2965 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 46 | ground_check | 1 | 2738 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 47 | ground_check | 1 | 2854 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 48 | alert_rest_ear_twitch | 1 | 1592 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 49 | alert_rest_ear_twitch | 1 | 1568 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 50 | alert_rest_ear_twitch | 1 | 1395 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 51 | alert_rest_ear_twitch | 2 | 2684 | 0 | 0 | mechanical-pass | retry: remove green residue |
| 52 | face_groom | 1 | 1455 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 53 | face_groom | 2 | 1537 | 0 | 0 | mechanical-pass | retry: thread-local source and remove green residue |
| 54 | face_groom | 1 | 1549 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 55 | face_groom | 1 | 1524 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 56 | reaction | 1 | 2022 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 57 | reaction | 1 | 2081 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 58 | reaction | 1 | 878 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 59 | reaction | 2 | 2479 | 0 | 0 | mechanical-pass | retry: stronger visible swept-back ears |
| 60 | reaction | 1 | 2343 | 0 | 0 | mechanical-pass | initial rescue candidate |
| 61 | reaction | 1 | 1675 | 0 | 0 | mechanical-pass | initial rescue candidate |

## Completion Check

- Candidate canonical frames present and mechanically clean: `22/22`
- Completed rescue range: `yes`

## Remaining Risk

- Mechanical QA catches size, alpha, edge contact, empty alpha, and green residue only.
- Parent thread still needs final visual comparison against the main fullrun for pose readability, anatomy, species identity, and style consistency.
