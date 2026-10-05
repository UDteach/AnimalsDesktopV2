# ChatGPT Pro follow-up consultation: run03 E prompt

Date: 2026-06-23

Context sent:

- A flat 2D style: raw parse `0/8`, diagnostic green-normalized parse `6/8` and `7/8`; most sprite-like but pinholes appeared.
- B gait delta: raw parse `0/8`, diagnostic green-normalized parse `8/8` and `7/8`; best gait readability.
- C extraction cleanliness: raw parse `0/8`, diagnostic green-normalized parse `8/8` and `8/8`; best extraction stability, slightly too realistic.
- D flat-clean hybrid: raw parse `0/8`, diagnostic green-normalized parse `4/8` and `6/8`; closer sprite feel but worse pinholes/parse.
- Common issue: raw parse fails because generated green backgrounds are not parser-pure `#00ff00` at edges and cell boundaries.

Advice received:

- Use C as the base and add only modest B-style gait deltas.
- Do not use D yet, because pushing the art style too hard appears to reintroduce matte and pinhole defects.
- Treat prompt-only pure `#00ff00` raw parser success as unrealistic for repeatable operation.
- For review-only sheets, diagnostic green normalization is acceptable as an evaluation aid.
- For accepted production frames, return to one-frame-per-PNG generation and the normal QA gate.

Recommended run03 focus:

- Generate E four times with the same prompt to measure reproducibility.
- Compare diagnostic normalized parse success, pinholes, gait readability, and frame-to-frame face/body/tail drift.
- Consider the sheet a motion-design and review artifact, not a direct accepted-frame source.

