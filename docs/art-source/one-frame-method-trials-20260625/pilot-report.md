# One-Frame Method Pilot Report - 2026-06-25

## Summary

The parent thread tested one high-risk pose per family using Codex built-in ImageGen one-frame generation.

All four pilot frames passed mechanical QA:

- normalized size: 96x64 PNG,
- non-empty alpha,
- no edge-alpha contact,
- no green residue above threshold,
- raw source and alpha source preserved.

Comparison image:

- `docs/art-source/one-frame-method-trials-20260625/baseline-vs-one-frame-pilot.png`

## Pilot Frames

| Family | Frame | Pose | Result | Notes |
| --- | ---: | --- | --- | --- |
| hamster | 58 | face-groom / whisker-clean | pass | Complete body and readable pose. Single-frame style is sharper and slightly less front-symmetric than the current sheet baseline. |
| macaroni_mouse | 20 | low fast scurry with thick tail | pass | Tail and feet are complete. Better pose clarity than the sheet baseline, but the result is more character-like and must be kept from drifting into hamster. |
| sugar_glider | 60 | front membrane-spread reaction | pass | Strong pose clarity and complete membrane. Slightly higher detail level than the sheet baseline. |
| rabbit | 15 | fast hop with ears contained | pass | Strong clarity and complete ears/tail/feet. 96x64 normalization remains usable despite tall ears. |

## Decision

One-frame generation is better for high-risk frame repair because it gives:

- independent retries for failed frame numbers,
- stronger pose targeting,
- fewer layout/crop artifacts than multi-sprite sheets,
- better provenance per frame.

It is not automatically better as a blind full-set replacement because:

- each frame is a separate generation and can drift in detail level, camera angle, or anatomy,
- no visible seed control is available in the built-in tool,
- style consistency must be enforced by prompt template, contact-sheet review, and rejection/retry.

Best current production method:

1. Keep the sheet-extracted 248-frame set as the baseline.
2. Run one-frame retry-until-pass fullruns per species in isolated directories.
3. Adopt one-frame outputs only when the entire family contact sheet stays style-consistent.
4. If a fullrun drifts, use it selectively to replace weak baseline frames rather than replacing the whole family blindly.

This interprets "reproducibility" as operational reproducibility: a stable prompt table, saved raw sources, per-frame retry logs, and frame-local replacement without invalidating already-good frames.
