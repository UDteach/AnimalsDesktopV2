# Toy Poodle apricot set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `toy_poodle_apricot`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Source run: `docs/art-source/one-frame-method-fullrun-20260628/toy-poodle-apricot-set00-oneframe-62/`
- Accepted frames: `docs/art-source/toy-poodle-apricot/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/toy-poodle-apricot/motion-source/sheets/toy-poodle-apricot-source-set00.png`

## Mechanical QA

- Source run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Source run `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=9.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=9.
- `cmd/assemblemotion`: wrote `toy-poodle-apricot-source-set00.png`.
- `cmd/validatemotion -variant toy_poodle_apricot -require-accepted`:
  deferred until catalog integration. The variant is not registered in the
  catalog yet, so this is not a frame-quality failure.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
ears, cropped feet, cropped tail, alpha-edge crop, broken contact, duplicate
face, or color drift to pale cream in the accepted frames.

## Slot Review

- `00-03 idle`: pass. Apricot Toy Poodle read is clear; `03` is a crouch pose, not a scale jump.
- `04-11 walk`: pass. Walking frames stay grounded with attached legs and stable coat color.
- `12-19 fast`: pass. Stretched movement remains dog-like and in-frame.
- `20-25 sniff/nibble`: pass. Nose-down frames stay complete and prop-free.
- `26-31 species action`: pass. Earlier retry frames were replaced; accepted frames keep one complete dog.
- `32-39 turn`: pass. Turn frames keep saturated apricot color and stable silhouette.
- `40-43 eat`: pass with note. `41` is slightly larger, but reads as a posture change rather than camera zoom.
- `44-47 ground check`: pass. Low poses remain grounded and uncropped.
- `48-51 alert/rest`: pass. Alert/rest sequence has stable contact and no duplicate face.
- `52-55 groom`: pass. Grooming motion keeps the poodle silhouette intact.
- `56-61 reaction/recover`: pass. Final recover frames return to the side stance without a visible scale jump.

## Visual Notes

- Species read: Toy Poodle / apricot coat, with curly fur, rounded ears, and pom tail.
- Scale: accepted. The visible larger frames were checked in contact sheets and read as pose-driven changes.
- Style continuity: accepted. Coat remains saturated apricot, not pale cream.
- Motion variety: accepted for source art. Runtime timing should still be checked when catalog integration maps these slots.
- Contact: reviewed frames stay consistently grounded; no visible floating.
- Reject checks: no text, border, scenery, prop, floor band, cast shadow,
  duplicate animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`toy_poodle_apricot` is accepted as local source art only. It should remain
outside the runtime/current Pages list until its release or preview patch is
planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/toy-poodle-apricot-set00-full-contact.png`
- Overview contact: `contacts/toy-poodle-apricot-overview-00-61-2x.png`
- Slot contacts: `contacts/toy-poodle-apricot-*-4x.png`
- Candidate sheet: `contacts/toy-poodle-apricot-source-set00-candidate-sheet.png`
- Animated preview: `contacts/toy-poodle-apricot-set00-preview.gif`
- One-frame QA copy: `qa/toy-poodle-apricot-oneframe-review.json`
- Run audit copy: `qa/toy-poodle-apricot-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
