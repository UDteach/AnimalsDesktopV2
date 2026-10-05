# Four Animals ImageGen Set00 Candidate Summary

Generated/verified at 2026-06-24T22:56:37.438528+00:00.

## Counts
- hamster: 62/62, missing=[], edge_alpha=0, bad_size=0, empty_alpha=0, greenish_gt_12=0
- macaroni_mouse: 62/62, missing=[], edge_alpha=0, bad_size=0, empty_alpha=0, greenish_gt_12=0
- sugar_glider: 62/62, missing=[], edge_alpha=0, bad_size=0, empty_alpha=0, greenish_gt_12=0
- rabbit: 62/62, missing=[], edge_alpha=0, bad_size=0, empty_alpha=0, greenish_gt_12=0

## Artifacts
- current-manifest.csv: exact current 248 candidate PNG rows.
- four-animals-set00-qa-report.json: size/alpha/edge/green-residue QA.
- four-animals-set00-overview.png: review contact sheet for all four families.
- `*-auditframes-report.json`: Go `cmd/auditframes -strict -artifact-warnings` reports for the current sheet-extracted candidate sets.
- `*-assemblemotion-report.json`: Go `cmd/assemblemotion` reports for review-only candidate source sheets.
- `*-set00-candidate-source.png`: review-only 5952x64 source sheets assembled from the current candidate frames.

## Notes
- These are accepted-candidates, not runtime-imported accepted-frames.
- Source ImageGen PNGs remain in /Users/kyota/.codex/generated_images/019efb9a-b711-77f2-b098-5dce8ddc678d.
- Go was installed with Homebrew during the follow-up task: `go version go1.26.4 darwin/arm64`.
- `go test -buildvcs=false ./...` and `go vet -buildvcs=false ./...` pass.
- Go `cmd/auditframes -strict` reports all four current candidate sets as valid=62, missing=0, invalid=0. `-artifact-warnings` is non-zero for all four and should be treated as visual-review guidance, not a hard failure.
- Go `cmd/assemblemotion` successfully assembled all four current candidate sets into review-only source sheets.
- Runtime import / `validatemotion -runtime-only -require-accepted` is still not appropriate for these candidate folders because they have not been promoted to `accepted-frames`.
