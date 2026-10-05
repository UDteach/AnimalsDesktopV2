# lovebird_peach_faced 00-04 QA summary

Status: ready for parent 00-04 visual gate.

Canonical frames:

- `frames/frame-00.png`
- `frames/frame-01.png`
- `frames/frame-02.png`
- `frames/frame-03.png`
- `frames/frame-04.png`

Source:

- Canonical frames were normalized from `raw/attempt-02-frame-00.png` through `raw/attempt-02-frame-04.png`.
- Attempt 01 was rejected because checkerboard was baked into the raw RGB background.
- Attempt 02 was salvaged with broad magenta background keying and per-frame fit.
- Attempt 03 raw files are retained as unused fallback evidence only.

Mechanical checks:

- All canonical frames are `96x64`.
- All canonical frames have non-empty alpha.
- Edge alpha is `0` for all five canonical frames.
- Detected residual magenta / white-gray background pixels are `0` for all five canonical frames.
- Light, dark, and checker contact sheets were regenerated with uncropped `216x150` cells.

Visual notes:

- Frames 00-04 read as peach-faced lovebird / コザクラインコ at 96x64.
- Green body/wings, peach-orange face, pale beak, short tail, and visible attached feet are present.
- No frame 05 was generated.

## 05-12 continuation

Status: ready for parent 05-12 visual gate.

Canonical frames:

- `frames/frame-05.png`
- `frames/frame-06.png`
- `frames/frame-07.png`
- `frames/frame-08.png`
- `frames/frame-09.png`
- `frames/frame-10.png`
- `frames/frame-11.png`
- `frames/frame-12.png`

Source:

- Canonical frames 05-12 were normalized from `raw/attempt-04-frame-05.png` through `raw/attempt-04-frame-12.png`.
- Approved canonical frames 00-04 were not regenerated or overwritten during this continuation.

Mechanical checks:

- All canonical frames 05-12 are `96x64`.
- All canonical frames 05-12 have non-empty alpha.
- Edge alpha is `0` for frames 05-12.
- Detected residual magenta / white-gray background pixels are `0` for frames 05-12.
- Light, dark, and checker contact sheets were generated for `00-12` and `05-12`.

Audit:

- `go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260627/lovebird-peach-faced-set00-oneframe-62/frames -strict -artifact-warnings`
- Result: `valid=13 missing=49 invalid=0 warnings=17`
- The missing count is expected because this lane stopped at frame 12.

Visual notes:

- Frames 05-12 read as peach-faced lovebird / コザクラインコ at 96x64.
- Frame 08 is slightly more upright than neighbors but keeps species identity, feet/contact, and baseline.
- No frame 13 was generated.

## 13-20 continuation

Status: ready for parent 13-20 visual gate.

Canonical frames:

- `frames/frame-13.png`
- `frames/frame-14.png`
- `frames/frame-15.png`
- `frames/frame-16.png`
- `frames/frame-17.png`
- `frames/frame-18.png`
- `frames/frame-19.png`
- `frames/frame-20.png`

Source:

- Canonical frames 13-20 were normalized from `raw/attempt-05-frame-13.png` through `raw/attempt-05-frame-20.png`.
- Approved canonical frames 00-12 were not regenerated or overwritten during this continuation.

Mechanical checks:

- All canonical frames 13-20 are `96x64`.
- All canonical frames 13-20 have non-empty alpha.
- Edge alpha is `0` for frames 13-20.
- Detected residual magenta / white-gray background pixels are `0` for frames 13-20.
- Light, dark, and checker contact sheets were generated for `00-20` and `13-20`.

Audit:

- `go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260627/lovebird-peach-faced-set00-oneframe-62/frames -strict -artifact-warnings`
- Result: `valid=21 missing=41 invalid=0 warnings=24`
- The missing count is expected because this lane stopped at frame 20.

Visual notes:

- Frames 13-19 read as peach-faced lovebird / コザクラインコ and continue the quick walk/hop motion.
- Frame 20 is intentionally lower for sniff/nibble start and should get parent visual confirmation for slot transition strength.
- No frame 21 was generated.
