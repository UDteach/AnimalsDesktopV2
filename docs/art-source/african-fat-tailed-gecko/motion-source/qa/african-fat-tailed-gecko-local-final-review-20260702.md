# Local Final Review

- Variant: `african_fat_tailed_gecko`
- Japanese label: ニシアフリカトカゲモドキ
- Source run: `docs/art-source/one-frame-method-fullrun-20260702/african-fat-tailed-gecko-set00-oneframe-62`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/african-fat-tailed-gecko-oneframe-review.json`, issues: none
- Run strict audit: `qa/african-fat-tailed-gecko-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=23`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=23`
- Assembled sheet: `sheets/african-fat-tailed-gecko-source-set00.png`
- Full contact: `contacts/african-fat-tailed-gecko-set00-full-contact.png`
- Candidate sheet: `contacts/african-fat-tailed-gecko-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/african-fat-tailed-gecko-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. Pass candidate. Low fat-tailed gecko silhouette, broad bands, and fat tail remain readable; watch subdued far-side legs in 48/51-55, slightly tapered tail tip in 57, and warmer/oranger frame 58.

This is source-art only. Runtime catalog, generated runtime sprite sheets,
GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

Known validation note: `go run ./cmd/validatemotion -variant african_fat_tailed_gecko
-require-accepted` is deferred until catalog integration because this
pass only creates accepted source art for the Pages candidate queue.
