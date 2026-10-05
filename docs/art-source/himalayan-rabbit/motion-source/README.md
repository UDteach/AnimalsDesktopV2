# Himalayan Rabbit motion source

`himalayan_rabbit` has an accepted `set00` source generated with the one-frame
ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/himalayan-rabbit-source-set00.png`
- Full contact: `contacts/himalayan-rabbit-set00-full-contact.png`
- Animated preview: `contacts/himalayan-rabbit-set00-preview.gif`
- Visual QA: `qa/himalayan-rabbit-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/himalayan-rabbit/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant himalayan_rabbit -require-accepted
```

This canonical `set00` is the complete accepted 62-frame source. The importer
expands it into all ten runtime slots; publication still requires the separate
runtime, build, page, packaging, and release-approval gates.
