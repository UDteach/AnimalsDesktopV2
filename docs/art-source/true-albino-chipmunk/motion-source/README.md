# True Albino Chipmunk Motion Source

`true_albino_chipmunk` has an accepted `set00` source generated with the
one-frame ImageGen workflow. This 2026-07-03 Direction B rebuild replaces the
prior disabled source.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/true-albino-chipmunk-source-set00.png`
- Full contact: `contacts/true-albino-chipmunk-set00-full-contact.png`
- Candidate sheet: `contacts/true-albino-chipmunk-source-set00-candidate-sheet.png`
- Animated preview: `contacts/true-albino-chipmunk-set00-preview.gif`
- Local final review: `qa/true-albino-chipmunk-local-final-review-20260703.md`

Visual target: true albino chipmunk with red eyes, no stripes, no dorsal
markings, no cheek stripe, smooth cream-white coat, subtle pink ears/paws/nose,
and readable white-background contrast.

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/true-albino-chipmunk/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/true-albino-chipmunk/motion-source/accepted-frames/set00 -out docs/art-source/true-albino-chipmunk/motion-source/sheets/true-albino-chipmunk-source-set00.png -report docs/art-source/true-albino-chipmunk/motion-source/accepted-frames/set00-assemblemotion-report.json
```
