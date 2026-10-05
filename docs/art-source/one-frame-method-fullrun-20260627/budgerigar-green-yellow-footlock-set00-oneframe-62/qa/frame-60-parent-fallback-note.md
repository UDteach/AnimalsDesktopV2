# Frame 60 Parent Fallback Note

Date: 2026-06-27

`frame-60.png` is a parent-created local fallback, not a fresh accepted
ImageGen frame.

Attempts `01` through `05` repeatedly produced the correct 96x64 scale and
baseline, but kept drifting into a gray/white natural-budgerigar face, pale wing
bars, and gray feet. Parent QA rejected those because adjacent frames `59` and
`61` use the simplified yellow-head style.

The fallback frame was produced from accepted neighbors:

- 85% `frames/frame-61.png`
- 15% `frames/frame-59.png`
- alpha from the max of both neighbors

This avoids the gray/white face and gray feet while keeping a near-inbetween
shape. It should be treated as a draft completion fallback. Before promoting
`budgerigar_green_yellow` to accepted source, run a full visual QA pass and
replace this fallback with a true generated or hand-approved frame if needed.

Parent QA contacts:

- `contact/budgerigar-footlock-59-61-frame60-parent-fallback-10x.png`
- `contact/budgerigar-footlock-55-61-parent-fallback-final-6x.png`
- `contact/budgerigar-footlock-00-61-parent-fallback-full-grid-3x.png`
