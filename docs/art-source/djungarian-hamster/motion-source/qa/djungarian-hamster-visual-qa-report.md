# Djungarian Hamster set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/djungarian-hamster/motion-source/accepted-frames/set00/`
- `docs/art-source/djungarian-hamster/motion-source/sheets/djungarian-hamster-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=85.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=85.
- `cmd/assemblemotion`: wrote `djungarian-hamster-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They mostly reflect the
low belly/contact baseline and fine alpha details that the artifact warning mode
flags. Parent visual review did not find a persistent floor, shelf, scenery,
prop, text, second animal, or cropped anatomy.

## Slot Review

- `00-03 idle`: pass. Compact gray-white dwarf hamster body and dorsal stripe are readable.
- `04-11 walk`: pass. Low scurry frames keep a small hamster silhouette.
- `12-19 fast`: pass. Body stretch increases but does not become a mouse or rat.
- `20-25 sniff/nibble`: pass. Nose-down and upright sniff poses are prop-free.
- `26-31 species action`: pass. Short tail read remains appropriate for hamster.
- `32-39 turn`: pass with note. Rotated poses keep the gray back and pale belly read.
- `40-43 chew/eat`: pass. Head and forepaw motion stay subtle, with no food prop.
- `44-47 ground check`: pass. Low baseline and dwarf-hamster proportions remain stable.
- `48-51 alert/rest`: pass. Upright posture keeps rounded ears and compact body.
- `52-55 face groom`: pass. Grooming poses stay one animal with no duplicate head.
- `56-61 reaction/recover`: pass. Recovery frames keep the species read and attached anatomy.

## Visual Notes

- Species read: Djungarian / winter white dwarf hamster, not a generic mouse or Syrian hamster recolor.
- Coat read: cool gray back, white underside, and dorsal stripe remain visible across the set.
- Anatomy: rounded ears, compact body, cheek-forward head, tiny feet, and very short tail remain readable.
- Reject checks: no text, border, scenery, food prop, costume, second animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`djungarian_hamster` is accepted for current runtime-preview parity only. It
should remain `release_ready=false` until accepted `set00` through `set09`
source families exist.

## Review Artifacts

- Full contact: `contacts/djungarian-hamster-set00-full-contact.png`
- Candidate sheet: `contacts/djungarian-hamster-source-set00-candidate-sheet.png`
- Animated preview: `contacts/djungarian-hamster-set00-preview.gif`
- One-frame QA copy: `qa/djungarian-hamster-oneframe-review.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
