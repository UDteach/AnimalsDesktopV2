# Holland Lop set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/holland-lop/motion-source/accepted-frames/set00/`
- `docs/art-source/holland-lop/motion-source/sheets/holland-lop-broken-orange-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=62.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=62.
- `cmd/assemblemotion`: wrote `holland-lop-broken-orange-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They mostly reflect the
low body/contact baseline, large ear silhouette, and fine alpha details that
artifact warning mode flags. Parent visual review did not find a persistent
floor, shelf, scenery, prop, text, second animal, or cropped anatomy.

## Slot Review

- `00-03 idle`: pass. Broken orange lop body and dropped ears are readable.
- `04-11 walk`: pass. Hop/walk poses keep a compact rabbit silhouette.
- `12-19 fast`: pass. Body stretch increases but remains rabbit-like.
- `20-25 sniff/nibble`: pass. Head-down poses stay prop-free.
- `26-31 species action`: pass. Lopped ears remain visible and attached.
- `32-39 turn`: pass. Upright/turn poses keep the same orange-and-white rabbit read.
- `40-43 chew/eat`: pass. No food prop appears.
- `44-47 ground check`: pass. Baseline and compact body remain stable.
- `48-51 alert/rest`: pass. Alert poses still show dropped ears, not upright ears.
- `52-55 face groom`: pass. Grooming poses remain one rabbit with no duplicate head.
- `56-61 reaction/recover`: pass. Recovery frames preserve the lop silhouette.

## Visual Notes

- Species read: Holland Lop, not a generic upright-eared rabbit.
- Coat read: broken orange-and-white patches remain readable across the set.
- Anatomy: rounded compact body, short muzzle, small feet, and dropped ears remain readable.
- Reject checks: no text, border, scenery, food prop, costume, second animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`holland_lop_broken_orange` is accepted for current runtime-preview parity only.
It should remain `release_ready=false` until accepted `set00` through `set09`
source families exist.

## Review Artifacts

- Full contact: `contacts/holland-lop-set00-full-contact.png`
- Candidate sheet: `contacts/holland-lop-source-set00-candidate-sheet.png`
- Animated preview: `contacts/holland-lop-set00-preview.gif`
- One-frame QA copy: `qa/holland-lop-oneframe-review.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
