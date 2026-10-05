# Miniature Schnauzer Salt Pepper Motion Source

`miniature_schnauzer_salt_pepper` has an accepted `set00` source generated
with the one-frame ImageGen workflow.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/miniature-schnauzer-salt-pepper-source-set00.png`
- Full contact: `contacts/miniature-schnauzer-salt-pepper-set00-full-contact.png`
- Candidate sheet: `contacts/miniature-schnauzer-salt-pepper-source-set00-candidate-sheet.png`
- Animated preview: `contacts/miniature-schnauzer-salt-pepper-set00-preview.gif`
- Local final review: `qa/miniature-schnauzer-salt-pepper-local-final-review-20260629.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/miniature-schnauzer/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/miniature-schnauzer/motion-source/accepted-frames/set00 -out docs/art-source/miniature-schnauzer/motion-source/sheets/miniature-schnauzer-salt-pepper-source-set00.png
```

`cmd/validatemotion -variant miniature_schnauzer_salt_pepper -require-accepted`
is deferred until catalog integration because the variant is not registered yet.
