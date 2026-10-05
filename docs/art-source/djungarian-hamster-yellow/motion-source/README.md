# イエロージャンガリアン Motion Source

`djungarian_hamster_yellow` has an accepted `set00` source generated with the one-frame
ImageGen workflow.

This promotion is source-art only. Runtime catalog, generated runtime
sprite sheets, GitHub Pages current-animal cards, release artifacts,
tags, and public release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/djungarian-hamster-yellow-source-set00.png`
- Full contact: `contacts/djungarian-hamster-yellow-set00-full-contact.png`
- Candidate sheet: `contacts/djungarian-hamster-yellow-source-set00-candidate-sheet.png`
- Animated preview: `contacts/djungarian-hamster-yellow-set00-preview.gif`
- Local final review: `qa/djungarian-hamster-yellow-local-final-review-20260702.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/djungarian-hamster-yellow/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/djungarian-hamster-yellow/motion-source/accepted-frames/set00 -out docs/art-source/djungarian-hamster-yellow/motion-source/sheets/djungarian-hamster-yellow-source-set00.png -report docs/art-source/djungarian-hamster-yellow/motion-source/accepted-frames/set00-assemblemotion-report.json
```

`cmd/validatemotion -variant djungarian_hamster_yellow -require-accepted` is deferred
until catalog integration because this source-art-only pass does not
register the variant in the runtime catalog.
