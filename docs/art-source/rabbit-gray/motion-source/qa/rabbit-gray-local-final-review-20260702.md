# Local Final Review

- Variant: `rabbit_gray`
- Japanese label: グレーうさぎ
- Source run: `docs/art-source/one-frame-method-fullrun-20260702/rabbit-gray-set00-oneframe-62`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/rabbit-gray-oneframe-review.json`, issues: none
- Run strict audit: `qa/rabbit-gray-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=71`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=71`
- Assembled sheet: `sheets/rabbit-gray-source-set00.png`
- Full contact: `contacts/rabbit-gray-set00-full-contact.png`
- Candidate sheet: `contacts/rabbit-gray-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/rabbit-gray-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. Pass candidate. Gray rabbit identity stays clear; watch rounder ground-check frames 44-45, upright alert/groom frames 48-55, and longer crouch frame 57.

This is source-art only. Runtime catalog, generated runtime sprite sheets,
GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

Known validation note: `go run ./cmd/validatemotion -variant rabbit_gray
-require-accepted` is deferred until catalog integration because this
pass only creates accepted source art for the Pages candidate queue.
