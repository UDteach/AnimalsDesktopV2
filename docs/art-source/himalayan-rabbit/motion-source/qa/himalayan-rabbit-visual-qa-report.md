# Himalayan Rabbit set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/himalayan-rabbit/motion-source/accepted-frames/set00/`
- `docs/art-source/himalayan-rabbit/motion-source/sheets/himalayan-rabbit-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=28.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=28.
- `cmd/assemblemotion`: wrote `himalayan-rabbit-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They mostly reflect
fine alpha details around ears, feet, tail, and the low contact baseline.
Parent visual review did not find a persistent floor, shelf, scenery, prop,
text, second animal, or cropped anatomy.

## Slot Review

- `00-03 idle`: pass. Pale body and dark points are readable.
- `04-11 walk`: pass. Hop/walk frames keep a rabbit silhouette.
- `12-19 fast`: pass. Body stretch increases but remains rabbit-like.
- `20-25 sniff/nibble`: pass. Head-down poses stay prop-free.
- `26-31 species action`: pass. Dark ears, nose, feet, and tail remain visible.
- `32-39 turn`: pass. Turn poses keep the Himalayan point markings.
- `40-43 chew/eat`: pass. No food prop appears.
- `44-47 ground check`: pass. Baseline and feet remain stable.
- `48-51 alert/rest`: pass. Upright poses keep dark ears and nose readable.
- `52-55 face groom`: pass. Grooming poses remain one rabbit with no duplicate head.
- `56-61 reaction/recover`: pass. Recovery frames preserve the same pointed coat.

## Visual Notes

- Species read: Himalayan rabbit, not a plain white rabbit.
- Coat read: cream body with dark ears, nose, feet, and tail remains readable across the set.
- Anatomy: upright ears, compact body, small feet, and short tail remain readable.
- Reject checks: no text, border, scenery, food prop, costume, second animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`himalayan_rabbit` is accepted for current runtime-preview parity only. It
should remain `release_ready=false` until accepted `set00` through `set09`
source families exist.

## Review Artifacts

- Full contact: `contacts/himalayan-rabbit-set00-full-contact.png`
- Candidate sheet: `contacts/himalayan-rabbit-source-set00-candidate-sheet.png`
- Animated preview: `contacts/himalayan-rabbit-set00-preview.gif`
- One-frame QA copy: `qa/himalayan-rabbit-oneframe-review.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
