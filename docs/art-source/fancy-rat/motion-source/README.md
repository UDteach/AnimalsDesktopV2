# Fancy Rat Hooded motion source

`fancy_rat_hooded` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/fancy-rat-hooded-source-set00.png`
- Full contact: `contacts/fancy-rat-hooded-set00-full-contact.png`
- Animated preview: `contacts/fancy-rat-hooded-set00-preview.gif`
- Visual QA: `qa/fancy-rat-hooded-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/fancy-rat/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant fancy_rat_hooded -require-accepted
```

This canonical `set00` is the complete accepted 62-frame source. The importer
expands it into all ten runtime slots; publication still requires the separate
runtime, build, page, packaging, and release-approval gates.
