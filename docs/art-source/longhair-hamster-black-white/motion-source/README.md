# 白黒長毛ハムスター Motion Source

`longhair_hamster_black_white` has an accepted `set00` source generated with the one-frame
ImageGen workflow.

This promotion is source-art only. Runtime catalog, generated runtime
sprite sheets, GitHub Pages current-animal cards, release artifacts,
tags, and public release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/longhair-hamster-black-white-source-set00.png`
- Full contact: `contacts/longhair-hamster-black-white-set00-full-contact.png`
- Candidate sheet: `contacts/longhair-hamster-black-white-source-set00-candidate-sheet.png`
- Animated preview: `contacts/longhair-hamster-black-white-set00-preview.gif`
- Local final review: `qa/longhair-hamster-black-white-local-final-review-20260702.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/longhair-hamster-black-white/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/longhair-hamster-black-white/motion-source/accepted-frames/set00 -out docs/art-source/longhair-hamster-black-white/motion-source/sheets/longhair-hamster-black-white-source-set00.png -report docs/art-source/longhair-hamster-black-white/motion-source/accepted-frames/set00-assemblemotion-report.json
```

`cmd/validatemotion -variant longhair_hamster_black_white -require-accepted` is deferred
until catalog integration because this source-art-only pass does not
register the variant in the runtime catalog.
