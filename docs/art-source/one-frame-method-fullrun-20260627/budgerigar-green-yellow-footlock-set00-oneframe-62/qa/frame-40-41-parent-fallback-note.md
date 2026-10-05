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

Treat the current set as a 62/62 draft motion asset. Before accepted-source
promotion, either approve these fallback frames explicitly or replace them with
true generated/hand-approved frames.
