# Lionhead Rabbit Brown White motion source

`lionhead_rabbit_brown_white` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/lionhead-rabbit-brown-white-source-set00.png`
- Full contact: `contacts/lionhead-rabbit-brown-white-set00-full-contact.png`
- Overview contact: `contacts/lionhead-rabbit-brown-white-overview-00-61-2x.png`
- Animated preview: `contacts/lionhead-rabbit-brown-white-set00-preview.gif`
- Visual QA: `qa/lionhead-rabbit-brown-white-visual-qa-report.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/lionhead-rabbit-brown-white/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/lionhead-rabbit-brown-white/motion-source/accepted-frames/set00 -out docs/art-source/lionhead-rabbit-brown-white/motion-source/sheets/lionhead-rabbit-brown-white-source-set00.png
```

`cmd/validatemotion -variant lionhead_rabbit_brown_white -require-accepted` is
deferred until catalog integration because the variant is not registered yet.
