# Richardson's Ground Squirrel set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/richardsons-ground-squirrel/motion-source/accepted-frames/set00/`
- `docs/art-source/richardsons-ground-squirrel/motion-source/sheets/richardsons-ground-squirrel-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=57.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=57.
- `cmd/assemblemotion`: wrote `richardsons-ground-squirrel-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They correspond mostly
to low tail and foot/contact alpha details. Parent visual review did not find a
persistent floor, shelf, scenery, prop, text, second animal, bold chipmunk
stripes, high curled tree-squirrel tail, or cropped anatomy.

## Slot Review

- `00-03 idle`: pass. Tan ground-squirrel body and low tail are readable.
- `04-11 walk`: pass. Low movement frames keep the same ground-dwelling profile.
- `12-19 fast`: pass. Body stretch increases but does not become chipmunk-like.
- `20-25 sniff/nibble`: pass. Nose-down and upright sniff poses are prop-free.
- `26-31 species action`: pass. Tail remains low, attached, and inside frame.
- `32-39 turn`: pass with note. Upright/angled frames keep the same tan animal.
- `40-43 chew/eat`: pass. Head and forepaw motion stay subtle, with no food prop.
- `44-47 ground check`: pass. Low baseline and body proportions remain stable.
- `48-51 alert/rest`: pass. Alert posture is taller but still ground-squirrel-like.
- `52-55 face groom`: pass. Grooming poses stay one animal with no duplicate head.
- `56-61 reaction/recover`: pass. Tail shape changes but remains attached.

## Visual Notes

- Species read: Richardson's ground squirrel, not chipmunk or generic tree squirrel.
- Coat read: tan/gray-brown coat with no bold chipmunk stripes remains stable.
- Anatomy: compact body, short ears, pointed muzzle, small feet, and low tail
  remain readable across the set.
- Reject checks: no text, border, scenery, food prop, costume, second animal,
  cropped ears, cropped feet, or cropped tail.

## Release Status

`richardsons_ground_squirrel` is accepted for current runtime-preview parity
only. It should remain `release_ready=false` until accepted `set00` through
`set09` source families exist.

## Review Artifacts

- Full contact: `contacts/richardsons-ground-squirrel-set00-full-contact.png`
- Candidate sheet: `contacts/richardsons-ground-squirrel-source-set00-candidate-sheet.png`
- Animated preview: `contacts/richardsons-ground-squirrel-set00-preview.gif`
- One-frame QA copy: `qa/richardsons-ground-squirrel-oneframe-review.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
