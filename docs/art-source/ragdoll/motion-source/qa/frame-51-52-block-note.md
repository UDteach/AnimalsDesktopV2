# Ragdoll seal bicolor frame 51-52 block note

Date: 2026-06-27

Scope: recovery attempt for canonical frames 50-52 only.

Result:
- `frame-50.png` was accepted from attempt 02 with bbox `(4,20)-(90,58)`, size `86x38`, baseline `57`.
- `frame-51` attempt 02 normalized to bbox `(4,9)-(92,58)`, size `88x49`, baseline `57`.
- The requested `frame-51` ladder was width `86-88`, height `40-44`, top `14-18`, baseline `57`.
- Because attempt 02 landed in the next height band and would create an abrupt rise, it was moved to `rejected/frame-51-attempt-02-ladder-miss.png`.
- `frame-51.png` and `frame-52.png` are intentionally left missing. `frame-52` was not generated after the frame-51 ladder miss.

Follow-up retry:
- `frame-51` attempt 03 normalized to bbox `(4,26)-(92,58)`, size `88x32`, baseline `57`.
- The requested follow-up target was width `86-88`, height `40-43` preferred, hard max `44`, top roughly `14-17`, baseline `57`.
- Attempt 03 was too low/short, effectively returning to the `frame-49` height band instead of bridging upward from `frame-50`.
- It was moved to `rejected/frame-51-attempt-03-too-low-short.png`.
- `frame-51.png` remains intentionally missing; `frame-52` was not generated.

Frame-37 anchored retry:
- Accepted anchor `frame-37.png` bbox is `(4,18)-(92,58)`, size `88x40`, baseline `57`.
- Accepted neighbor `frame-50.png` bbox is `(4,20)-(90,58)`, size `86x38`, baseline `57`.
- `frame-51` attempt 04 normalized to bbox `(4,23)-(92,58)`, size `88x35`, baseline `57`.
- The requested anchor target was width `86-88`, height `40-42`, top `16-18`, baseline `57`; hard minimum height was `39`.
- Attempt 04 was too low/short and was moved to `rejected/frame-51-attempt-04-too-low-short.png`.
