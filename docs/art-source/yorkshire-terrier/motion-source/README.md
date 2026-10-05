# Yorkshire Terrier Longcoat motion source

`yorkshire_terrier_longcoat` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/yorkshire-terrier-longcoat-source-set00.png`
- Full contact: `contacts/yorkshire-terrier-longcoat-set00-full-contact.png`
- Animated preview: `contacts/yorkshire-terrier-longcoat-set00-preview.gif`
- Visual QA: `qa/yorkshire-terrier-longcoat-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/yorkshire-terrier/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant yorkshire_terrier_longcoat -require-accepted
```

This canonical `set00` is the complete accepted 62-frame source. The importer
expands it into all ten runtime slots; publication still requires the separate
runtime, build, page, packaging, and release-approval gates.
