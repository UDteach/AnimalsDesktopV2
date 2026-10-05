# Gecko Motion Source

`gecko_gray_brown` has an accepted `set00` source generated with the one-frame ImageGen workflow.

Current coverage:

- `accepted-frames/set00/frame-00.png` through `frame-61.png`
- `sheets/gecko-gray-brown-source-set00.png`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/gecko/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/gecko/motion-source/accepted-frames/set00 -out docs/art-source/gecko/motion-source/sheets/gecko-gray-brown-source-set00.png -report docs/art-source/gecko/motion-source/accepted-frames/set00-assemblemotion-report.json
go run ./cmd/validatemotion -variant gecko_gray_brown -require-accepted
```

Notes:

- This is accepted as slot-based motion, not a single continuous 62-frame loop.
- Canonical `set00` is the complete accepted source and expands into all ten runtime slots.
- Generation provenance and visual QA are in `docs/art-source/one-frame-method-fullrun-20260626/gecko-set00-oneframe-62/`.
