# Run07: Chinchilla 4-frame ImageGen Stability Loop

Date: 2026-06-23

Purpose: verify whether a 2x2 ImageGen sheet can become stable enough for review-speed batching without accepting the sheet as production art.

## Attempts

| Variant | Prompt change | Raw parse | Greenfixed parse | Recentered parse | Visual note |
| --- | --- | ---: | ---: | ---: | --- |
| J1 | Medium sprite size, silhouette-safe pose labels | 0/4 | 0/4 | 0/4 | Four animals, but 1px chroma pinholes remained. |
| K1 | Solid closed chunky shapes, no whiskers, fused feet/tail | 0/4 | 4/4 | 4/4 | Clean diagnostic parse; style simplified. |
| K2 | Same as K1 | 0/4 | 4/4 | 4/4 | Stable. |
| K3 | Same as K1 | 0/4 | 4/4 | 4/4 | Stable. |
| K4 | Same as K1 | 0/4 | 4/4 | 4/4 | Stable. |

## Interpretation

- K is the first run in this lane to reach greenfixed `4/4` on four consecutive ImageGen samples.
- All K samples have exactly four primary components.
- Raw parse still fails because ImageGen emits `1254x1254` sheets with non-parser-pure green and boundary content.
- The technical fix comes from removing pinhole-prone details: whiskers, fur texture, separated feet, tail gaps, and small negative spaces.
- K is not final art. It lacks DeguDesktop-equivalent detail and has weak motion deltas.

## Next Variable

Restore one visual detail at a time while preserving closed shapes:

1. Add only short attached muzzle markings drawn fully on the muzzle shape.
2. Keep fused foot pads, attached bean-like tail, and no internal green gaps.
3. Re-run four samples and require greenfixed `4/4` in at least three of four before considering the variable stable.
