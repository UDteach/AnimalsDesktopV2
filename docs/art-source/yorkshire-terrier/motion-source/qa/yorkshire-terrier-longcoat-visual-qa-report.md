# Yorkshire Terrier Longcoat set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/yorkshire-terrier/motion-source/accepted-frames/set00/`
- `docs/art-source/yorkshire-terrier/motion-source/sheets/yorkshire-terrier-longcoat-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=92.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=92.
- `cmd/assemblemotion`: wrote `yorkshire-terrier-longcoat-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They correspond mostly
to long hair fringes, feet, and lower-edge alpha details. Parent visual review
did not find a persistent floor, shelf, scenery, prop, text, second animal, or
cropped anatomy.

## Slot Review

- `00-03 idle`: pass. Low long-coated toy Yorkie body is readable.
- `04-11 walk`: pass. Long coat and short-leg silhouette remain stable.
- `12-19 fast`: pass. Body stretch increases but does not become shepherd-like.
- `20-25 sniff/nibble`: pass. Nose-down poses are prop-free.
- `26-31 species action`: pass. Hair and legs remain complete and in-frame.
- `32-39 turn`: pass with note. Turn frames show more front angle but keep the same Yorkie.
- `40-43 chew/eat`: pass. Head and forepaw motion stay subtle, with no food prop.
- `44-47 ground check`: pass. Low baseline and long coat remain stable.
- `48-51 alert/rest`: pass. The alert posture stays compact and toy-dog-like.
- `52-55 face groom`: pass. Grooming poses stay one animal with no duplicate head.
- `56-61 reaction/recover`: pass. Coat shape changes but stays attached and in-frame.

## Visual Notes

- Species read: long-coated Yorkshire Terrier, not the earlier shepherd/spitz-like retries.
- Coat read: black body coat with tan face/legs remains stable.
- Anatomy: short legs, small body, long coat, tan face, and compact toy-dog read
  remain visible across the set.
- Reject checks: no text, border, scenery, food prop, costume, second animal,
  cropped ears, cropped feet, or cropped coat/tail.

## Release Status

`yorkshire_terrier_longcoat` is accepted for current runtime-preview parity
only. It should remain `release_ready=false` until accepted `set00` through
`set09` source families exist.

## Review Artifacts

- Full contact: `contacts/yorkshire-terrier-longcoat-set00-full-contact.png`
- Candidate sheet: `contacts/yorkshire-terrier-longcoat-source-set00-candidate-sheet.png`
- Animated preview: `contacts/yorkshire-terrier-longcoat-set00-preview.gif`
- One-frame QA copy: `qa/yorkshire-terrier-longcoat-oneframe-review.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
