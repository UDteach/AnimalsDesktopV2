# Run 06 4-Frame 2x2 Summary

Date: 2026-06-23

Scope: review-only ImageGen 4-frame chinchilla sheets. These are not accepted assets.

## External Brainstorm

ChatGPT Pro text-only review and Gemini Pro text-only review agreed on the main direction:

- Background purity is a model limitation. Raw fixed-cell parse should not be expected until the generator can produce parser-pure `#00ff00`.
- The large, detailed, illustration-like look is prompt-sensitive and can be reduced.
- Degree-based foot/limb angle wording should be avoided for this test because it encourages complex separated limb geometry and pinholes.
- Prefer silhouette-safe phase labels: neutral compact pose, tiny forward lean, tiny opposite lean, soft recovery pose.
- Four-frame sheets are more promising than eight-frame sheets for layout and rough pose review, but not for accepted-frame production.

## Local Results

| Sample | Prompt type | Raw parse | Greenfixed parse | Recentered parse | Main failure |
| --- | --- | ---: | ---: | ---: | --- |
| H1 | attached foot angle prompt | 0/4 | 3/4 | 3/4 | one 1px chroma pinhole |
| H2 | attached foot angle prompt | 0/4 | 3/4 | 3/4 | one 1px chroma pinhole |
| I1 | silhouette-safe small sprite prompt | 0/4 | 2/4 | 2/4 | two frames with 1px chroma pinholes |

All three samples produced exactly four primary chinchilla components. The failure is not animal count or large detached debris; it is matte purity plus lingering 1px pinholes. H1/H2 are visually too large and still illustration-like. I1 is smaller and has better margin, but the stricter small-sprite prompt did not improve parse stability.

## Decision

Continue treating 4-frame sheets as review-only. They are useful for:

- checking whether a service can place exactly four animals in a simple grid
- comparing rough body lean / tail balance / family consistency
- testing prompt wording before returning to one-frame-per-PNG production

Do not use these outputs as accepted-frame sources. Do not locally repair pinholes or green-normalized frames for promotion.

Adoption threshold for another 4-frame retry:

- generate four H-next-style samples
- diagnostic greenfixed parse must reach `4/4` on at least three samples
- no detached alpha components
- no belly/feet/tail pinholes
- visual style must be closer to compact 2D sprite than current H samples

If the next four-frame retry stays at `3/4` or worse, stop the 4-frame sheet lane for chinchilla and return to one-frame-per-PNG or explicit reference-upload single-frame candidates.
