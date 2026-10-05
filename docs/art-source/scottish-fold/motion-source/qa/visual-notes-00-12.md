# scottish_fold_silver_tabby set00 frames 00-12 QA notes

Scope: child production lane only. Frames 00-04 were parent-approved; this pass added frames 05-12 and stopped before frame 13.

Generated canonical frames:
- `frames/frame-05.png`
- `frames/frame-06.png`
- `frames/frame-07.png`
- `frames/frame-08.png`
- `frames/frame-09.png`
- `frames/frame-10.png`
- `frames/frame-11.png`
- `frames/frame-12.png`

Contact sheets:
- `contact/scottish_fold_silver_tabby-set00-frames-05-12-light.png`
- `contact/scottish_fold_silver_tabby-set00-frames-05-12-dark.png`
- `contact/scottish_fold_silver_tabby-set00-frames-05-12-checker.png`
- `contact/scottish_fold_silver_tabby-set00-frames-00-12-light.png`
- `contact/scottish_fold_silver_tabby-set00-frames-00-12-dark.png`
- `contact/scottish_fold_silver_tabby-set00-frames-00-12-checker.png`

Audit:
- Command: `go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260627/scottish-fold-silver-tabby-set00-oneframe-62/frames -strict -artifact-warnings -report docs/art-source/one-frame-method-fullrun-20260627/scottish-fold-silver-tabby-set00-oneframe-62/qa/auditframes-00-12-report.json`
- Result: `valid=13 missing=49 invalid=0 warnings=5`
- Missing `13-61` is expected because this lane stops at frame 12.

Warning frames:
- `frame-06`: lower-ledge heuristic and one 2 px transparent pinhole warning.
- `frame-07`: lower-ledge heuristic.
- `frame-08`: one 1 px transparent pinhole warning.
- `frame-12`: two 1 px transparent pinhole warnings.

Visual review:
- 05-12 are right-facing, transparent, and remain the same silver tabby Scottish Fold family.
- Folded-ear head, silver tabby stripes, thick curved tail, and stable baseline remain readable on light, dark, and checker contacts.
- 05-12 are visibly low compared with 00-03, matching the parent-approved low `frame-04` cat-stalk start. No crop, duplicate animal, background leak, floor plane, shadow, or sudden scale jump is visible in contacts.
