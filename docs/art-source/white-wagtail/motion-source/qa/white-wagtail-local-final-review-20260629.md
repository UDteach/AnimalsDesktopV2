# Local Final Review

- Variant: `white_wagtail`
- Frame count: 62/62 canonical frames present (`00-61`)
- One-frame review: `qa/white-wagtail-oneframe-review.json`, issues: none
- Run strict audit: `qa/white-wagtail-run-auditframes-report.json`, `valid=62 missing=0 invalid=0 warnings=6`
- Accepted strict audit: `qa/auditframes-set00.json`, `valid=62 missing=0 invalid=0 warnings=6`
- Assembled sheet: `sheets/white-wagtail-source-set00.png`
- Full contact: `contacts/white-wagtail-set00-full-contact.png`
- Candidate sheet: `contacts/white-wagtail-source-set00-candidate-sheet.png`
- Preview GIF: `contacts/white-wagtail-set00-preview.gif`

Visual verdict: parent pass for source-art promotion. The completed set reads
as a white wagtail with black cap and bib, white cheek and belly, gray back,
black wings with white wing bars, slim grounded legs, and a long attached tail.
The parent review accepted the larger tail wag/bob and low peck poses as
slot-appropriate motion. No duplicate animal, crop, floating body, scenery,
text, or species drift is visible in the final contact.

Known validation note: `go run ./cmd/validatemotion -variant white_wagtail
-require-accepted` currently reports `unknown variant "white_wagtail"` because
runtime catalog integration has not been done in this source-art-only pass.
