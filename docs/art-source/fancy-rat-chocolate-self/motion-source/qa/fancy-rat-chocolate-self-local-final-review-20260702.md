# Local Final Review

- Variant: `fancy_rat_chocolate_self`
- Japanese label: ファンシーラット（チョコレートセルフ）
- Source run: `docs/art-source/one-frame-method-fullrun-20260702/fancy-rat-chocolate-self-set00-oneframe-62`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/fancy-rat-chocolate-self-oneframe-review.json`, issues: none
- Run strict audit: `qa/fancy-rat-chocolate-self-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=56`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=56`
- Assembled sheet: `sheets/fancy-rat-chocolate-self-source-set00.png`
- Full contact: `contacts/fancy-rat-chocolate-self-set00-full-contact.png`
- Candidate sheet: `contacts/fancy-rat-chocolate-self-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/fancy-rat-chocolate-self-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. Pass candidate after turn repair. Chocolate self rat stays complete and low enough across the final contact; watch tall but slot-acceptable frames 27-28 and noted late poses 52/54/59.

This is source-art only. Runtime catalog, generated runtime sprite sheets,
GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

Known validation note: `go run ./cmd/validatemotion -variant fancy_rat_chocolate_self
-require-accepted` is deferred until catalog integration because this
pass only creates accepted source art for the Pages candidate queue.
