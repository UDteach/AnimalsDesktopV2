# Local Final Review

- Variant: `fancy_rat_cream_agouti`
- Japanese label: ファンシーラット（クリームアグーチ）
- Source run: `docs/art-source/one-frame-method-fullrun-20260702/fancy-rat-cream-agouti-set00-oneframe-62`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/fancy-rat-cream-agouti-oneframe-review.json`, issues: none
- Run strict audit: `qa/fancy-rat-cream-agouti-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=35`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=35`
- Assembled sheet: `sheets/fancy-rat-cream-agouti-source-set00.png`
- Full contact: `contacts/fancy-rat-cream-agouti-set00-full-contact.png`
- Candidate sheet: `contacts/fancy-rat-cream-agouti-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/fancy-rat-cream-agouti-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. Pass candidate after turn repair. Cream agouti rat reads as a pale rat with long tail; frame 35 was reviewed by zoom and accepted as low side-view turn rather than vertical mascot drift.

This is source-art only. Runtime catalog, generated runtime sprite sheets,
GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

Known validation note: `go run ./cmd/validatemotion -variant fancy_rat_cream_agouti
-require-accepted` is deferred until catalog integration because this
pass only creates accepted source art for the Pages candidate queue.
