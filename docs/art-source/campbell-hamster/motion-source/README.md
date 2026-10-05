# Campbell Hamster motion source

`campbell_hamster` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/campbell-hamster-source-set00.png`
- Full contact: `contacts/campbell-hamster-set00-full-contact.png`
- Animated preview: `contacts/campbell-hamster-set00-preview.gif`
- Visual QA: `qa/campbell-hamster-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/campbell-hamster/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant campbell_hamster -require-accepted
```

Frames `60` and `61` had usable anatomy but warmer color than the surrounding
frames. The parent pass preserved the original frames under
`qa/color-correction-20260626/`, corrected only those two visible-pixel color
distributions toward frames `56-59`, regenerated contact sheets, and reran
`cmd/auditframes` before promotion.

This canonical `set00` is the complete accepted 62-frame source. The importer
expands it into all ten runtime slots; publication still requires the separate
runtime, build, page, packaging, and release-approval gates.
