# Chinchilla 8-frame ImageGen reproducibility run 02

Date: 2026-06-23

Purpose: consult ChatGPT Pro for improved prompts, then test reproducibility with multiple ImageGen outputs.

## Prompt Strategies

ChatGPT Pro suggested three controlled variants:

- `A`: flatter 2D sprite style and simpler silhouette.
- `B`: stronger per-frame gait differences while keeping identity stable.
- `C`: matte/background cleanliness and extraction success first.

Local follow-up added:

- `D`: hybrid of `A` and `C`, trying to keep flatter style while preserving extraction cleanliness.

Each variant was generated twice with Codex ImageGen. All outputs are review-only and must not be promoted to `accepted-frames`.

## Raw Outputs

| Variant | Sample files |
| --- | --- |
| A | `variant-a-flat-2d-01.png`, `variant-a-flat-2d-02.png` |
| B | `variant-b-gait-delta-01.png`, `variant-b-gait-delta-02.png` |
| C | `variant-c-cleanliness-01.png`, `variant-c-cleanliness-02.png` |
| D | `variant-d-flat-clean-hybrid-01.png`, `variant-d-flat-clean-hybrid-02.png` |

Every sample produced:

- native `1774x887` output
- exactly 8 detected animal components after diagnostic green normalization
- no visible grid lines, labels, or large extra debris

## Parse Summary

Raw fixed-cell parse failed for every sample because generated green was not parser-pure at the boundary. Diagnostic green-normalized parse results:

| Variant | Sample 1 | Sample 2 | Average | Full `8/8` passes |
| --- | ---: | ---: | ---: | ---: |
| A: flat 2D style | 6 / 8 | 7 / 8 | 6.5 | 0 / 2 |
| B: gait delta | 8 / 8 | 7 / 8 | 7.5 | 1 / 2 |
| C: extraction cleanliness | 8 / 8 | 8 / 8 | 8.0 | 2 / 2 |
| D: flat-clean hybrid | 4 / 8 | 6 / 8 | 5.0 | 0 / 2 |

Cell-recentered diagnostic parse matched the green-normalized counts; local recentering did not rescue additional frames.

## Visual Review

- `A` is the most consistently sprite-like among the original three, but it still creates pinholes.
- `B` gives the strongest walk-cycle differences and one full `8/8` diagnostic pass.
- `C` is the only variant with reproducible `8/8` diagnostic parse across both samples, but remains slightly realistic.
- `D` looks closest to the desired flat desktop mascot style, but its cleaner-looking art creates more pinholes and worse parse stability.

## Decision

Adopt `C` as the next prompt base for mechanical reproducibility. Add only a small amount of `B` gait wording. Do not use the `D` hybrid yet, because it improves visual style at the cost of extraction reliability.

Next prompt direction:

- Keep C's extraction-cleanliness ordering and matte/background language.
- Add B's explicit but modest foot/body/tail gait differences.
- Avoid extra "smaller than previous" or heavy style-compaction wording until parse stability stays at `8/8`.
- Continue to treat green normalization and cell recentering as diagnostics only, not accepted-frame cleanup.
