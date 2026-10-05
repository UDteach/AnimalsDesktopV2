# Shoebill stork set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `shoebill_stork`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Parent run: `docs/art-source/one-frame-method-fullrun-20260628/shoebill-stork-set00-oneframe-62/`
- Accepted frames: `docs/art-source/shoebill-stork/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/shoebill-stork/motion-source/sheets/shoebill-stork-source-set00.png`

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=38.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=38.
- `cmd/assemblemotion`: wrote `shoebill-stork-source-set00.png` at 5952x64.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
head, cropped feet, or cropped tail.

## Slot Review

- `00-03 idle`: pass. Tall shoebill stance, heavy bill, attached legs, and gray-blue body read clearly.
- `04-11 walk`: pass. Steps are intentionally subtle and keep the body nearly still.
- `12-19 fast`: pass for this special lane. Motion remains restrained rather than bouncy.
- `20-25 sniff/nibble`: pass. Bill and head shifts stay attached with no prop or floor mark.
- `26-31 species action`: pass. The bird remains statue-like with small posture changes.
- `32-39 turn`: pass. The turn stays narrow and does not introduce a second face or wrong species read.
- `40-43 eat`: pass. Bill/mouth micro-motion is visible but the body stays stable.
- `44-47 ground check`: pass. The pose returns to narrow/quiet stance with stable feet.
- `48-51 alert/rest`: pass. Nearly static frames match the low-motion direction.
- `52-55 groom`: pass. Small bill changes remain coherent and do not widen the body.
- `56-61 reaction/recover`: pass. Final frames return to a stable shoebill silhouette.

## Visual Notes

- Species read: shoebill stork / ハシビロコウ.
- Motion direction: special low-motion lane; the set intentionally reads as
  statue-like with subtle bill, posture, and leg micro-shifts.
- Anatomy: tall body, long heavy bill, attached legs, visible feet, and folded
  wings remain readable.
- Reject checks: no text, border, scenery, perch, food prop, floor, shadow,
  duplicate bird, cropped head, cropped feet, or cropped bill.

## Release Status

`shoebill_stork` is accepted as local source art only. It should remain outside
the runtime/current Pages list until its one-animal patch version is planned,
cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/shoebill-stork-set00-full-contact.png`
- Overview contact: `contacts/shoebill-stork-overview-00-61-2x.png`
- Slot contacts: `contacts/shoebill-stork-*-4x.png`
- Candidate sheet: `contacts/shoebill-stork-source-set00-candidate-sheet.png`
- Animated preview: `contacts/shoebill-stork-set00-preview.gif`
- One-frame QA copy: `qa/shoebill-stork-oneframe-review.json`
- Run audit copy: `qa/shoebill-stork-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
