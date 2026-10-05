# Toy Poodle Apricot motion source

`toy_poodle_apricot` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/toy-poodle-apricot-source-set00.png`
- Full contact: `contacts/toy-poodle-apricot-set00-full-contact.png`
- Overview contact: `contacts/toy-poodle-apricot-overview-00-61-2x.png`
- Animated preview: `contacts/toy-poodle-apricot-set00-preview.gif`
- Visual QA: `qa/toy-poodle-apricot-visual-qa-report.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/toy-poodle-apricot/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/toy-poodle-apricot/motion-source/accepted-frames/set00 -out docs/art-source/toy-poodle-apricot/motion-source/sheets/toy-poodle-apricot-source-set00.png
```

`cmd/validatemotion -variant toy_poodle_apricot -require-accepted` is deferred
until catalog integration because the variant is not registered yet.
