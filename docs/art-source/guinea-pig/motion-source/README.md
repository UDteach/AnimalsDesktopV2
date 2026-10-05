# Guinea Pig Motion Source

`guinea_pig_tricolor` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Current coverage:

- `accepted-frames/set00/frame-00.png` through `frame-61.png`
- `sheets/guinea-pig-tricolor-source-set00.png`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/guinea-pig/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/guinea-pig/motion-source/accepted-frames/set00 -out docs/art-source/guinea-pig/motion-source/sheets/guinea-pig-tricolor-source-set00.png -report docs/art-source/guinea-pig/motion-source/accepted-frames/set00-assemblemotion-report.json
go run ./cmd/validatemotion -variant guinea_pig_tricolor -require-accepted
```

Notes:

- This is accepted as slot-based motion, not a single continuous 62-frame loop.
- Canonical `set00` is the complete accepted source and expands into all ten runtime slots.
- Generation provenance and visual QA are in `docs/art-source/one-frame-method-fullrun-20260626/guinea-pig-tricolor-set00-oneframe-62/`.
- Do not use `guinea-pig-tricolor-template-lock-set00-62` as accepted source; its late frames were visually rejected because the template-lock pass damaged the animal silhouette.
