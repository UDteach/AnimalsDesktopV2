# White Wagtail Motion Source

`white_wagtail` has an accepted `set00` source generated with the one-frame
ImageGen workflow.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/white-wagtail-source-set00.png`
- Full contact: `contacts/white-wagtail-set00-full-contact.png`
- Candidate sheet: `contacts/white-wagtail-source-set00-candidate-sheet.png`
- Animated preview: `contacts/white-wagtail-set00-preview.gif`
- Local final review: `qa/white-wagtail-local-final-review-20260629.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/white-wagtail/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/white-wagtail/motion-source/accepted-frames/set00 -out docs/art-source/white-wagtail/motion-source/sheets/white-wagtail-source-set00.png
```

`cmd/validatemotion -variant white_wagtail -require-accepted` is deferred
until catalog integration because the variant is not registered yet.
