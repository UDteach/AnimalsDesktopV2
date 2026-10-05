# Local Final Review

- Variant: `african_dormouse`
- Japanese label: アフリカヤマネ
- Source run: `docs/art-source/one-frame-method-fullrun-20260702/african-dormouse-set00-oneframe-62`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/african-dormouse-oneframe-review.json`, issues: none
- Run strict audit: `qa/african-dormouse-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=84`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=84`
- Assembled sheet: `sheets/african-dormouse-source-set00.png`
- Full contact: `contacts/african-dormouse-set00-full-contact.png`
- Candidate sheet: `contacts/african-dormouse-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/african-dormouse-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. Pass candidate. Reads as a small African dormouse with furry tail and rounded ears; watch matte pinhole warnings as visual-only warnings.

This is source-art only. Runtime catalog, generated runtime sprite sheets,
GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

Known validation note: `go run ./cmd/validatemotion -variant african_dormouse
-require-accepted` is deferred until catalog integration because this
pass only creates accepted source art for the Pages candidate queue.
