# Local Final Review

- Variant: `fancy_rat_blue_hooded`
- Japanese label: ファンシーラット（ブルーフーディッド）
- Source run: `docs/art-source/one-frame-method-fullrun-20260702/fancy-rat-blue-hooded-set00-oneframe-62`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/fancy-rat-blue-hooded-oneframe-review.json`, issues: none
- Run strict audit: `qa/fancy-rat-blue-hooded-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=44`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=44`
- Assembled sheet: `sheets/fancy-rat-blue-hooded-source-set00.png`
- Full contact: `contacts/fancy-rat-blue-hooded-set00-full-contact.png`
- Candidate sheet: `contacts/fancy-rat-blue-hooded-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/fancy-rat-blue-hooded-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. Pass candidate after turn repair. Hooded rat identity, long tail, and blue hooded coat remain readable; watch tall turn frames 28-29 and high-ish repaired frame 34.

This is source-art only. Runtime catalog, generated runtime sprite sheets,
GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

Known validation note: `go run ./cmd/validatemotion -variant fancy_rat_blue_hooded
-require-accepted` is deferred until catalog integration because this
pass only creates accepted source art for the Pages candidate queue.
