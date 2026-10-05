# Richardson's Ground Squirrel motion source

`richardsons_ground_squirrel` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/richardsons-ground-squirrel-source-set00.png`
- Full contact: `contacts/richardsons-ground-squirrel-set00-full-contact.png`
- Animated preview: `contacts/richardsons-ground-squirrel-set00-preview.gif`
- Visual QA: `qa/richardsons-ground-squirrel-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/richardsons-ground-squirrel/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant richardsons_ground_squirrel -require-accepted
```

This canonical `set00` is the complete accepted 62-frame source. The importer
expands it into all ten runtime slots; publication still requires the separate
runtime, build, page, packaging, and release-approval gates.
