# Frames 40-41 Parent Fallback Note

Date: 2026-06-27

`frame-40.png` and `frame-41.png` are parent-created local fallback frames.

The previous canonical frames came from `attempt-02` and had correct scale,
baseline, and body shape, but retained a gray/white natural-budgerigar face and
pale wing/foot drift. Direct targeted recolor was inspected and not adopted
because it did not remove the gray face/wing drift cleanly.

Fallback method:

- `frame-40.png`: blend of `frame-39.png` and `frame-42.png`, weighted 70/30.
- `frame-41.png`: blend of `frame-39.png` and `frame-42.png`, weighted 25/75.
- alpha from the max of the blended neighbors.

Rejected originals were preserved under:

- `rejected/frame-40-attempt-02-face-color-drift-parent-replaced.png`
- `rejected/frame-41-attempt-02-face-color-drift-parent-replaced.png`

Parent QA contacts:

- `contact/budgerigar-footlock-38-43-parent-fallback-final-6x.png`
- `contact/budgerigar-footlock-00-61-parent-fallback-full-grid-v2-3x.png`

Promotion decision:

Parent accepted these fallback frames for the source asset after the `38-43`
and full-grid v2 contact reviews. They remove the gray/white face drift while
keeping visible feet/contact, scale, baseline, and motion continuity. This is
source-art acceptance only; it does not add the variant to runtime, Pages,
release notes, tags, or downloadable artifacts.
