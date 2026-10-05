# Domestic Shorthair Tabby White Stocky Motion Source

`domestic_shorthair_tabby_white_stocky` has an accepted `set00` source
generated with the one-frame ImageGen workflow from the supplied photo
reference.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this source promotion.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/domestic-shorthair-tabby-white-stocky-source-set00.png`
- Full contact: `contacts/domestic-shorthair-tabby-white-stocky-set00-full-contact.png`
- Final light contact: `contacts/domestic-shorthair-tabby-white-stocky-set00-final-light-3x.png`
- Final dark contact: `contacts/domestic-shorthair-tabby-white-stocky-set00-final-dark-3x.png`
- Final checker contact: `contacts/domestic-shorthair-tabby-white-stocky-set00-final-checker-3x.png`
- Candidate sheet: `contacts/domestic-shorthair-tabby-white-stocky-source-set00-candidate-sheet.png`
- Animated preview: `contacts/domestic-shorthair-tabby-white-stocky-set00-preview.gif`
- Local final review: `qa/domestic-shorthair-tabby-white-stocky-local-final-review-20260629.md`
- Photo reference: `../reference/source-photo-20260629.jpg`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/accepted-frames/set00 -out docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/sheets/domestic-shorthair-tabby-white-stocky-source-set00.png
```

`cmd/validatemotion -variant domestic_shorthair_tabby_white_stocky
-require-accepted` is deferred until catalog integration because the variant is
not registered yet.
