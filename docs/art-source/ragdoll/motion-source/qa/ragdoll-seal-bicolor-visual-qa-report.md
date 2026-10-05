# Ragdoll seal bicolor set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `ragdoll_seal_bicolor`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Source run: `docs/art-source/one-frame-method-fullrun-20260628/ragdoll-seal-bicolor-set00-oneframe-62/`
- Accepted frames: `docs/art-source/ragdoll-seal-bicolor/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/ragdoll-seal-bicolor/motion-source/sheets/ragdoll-seal-bicolor-source-set00.png`

## Mechanical QA

- Source run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Source run `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=21.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=21.
- `cmd/assemblemotion`: wrote `ragdoll-seal-bicolor-source-set00.png`.
- `cmd/validatemotion -variant ragdoll_seal_bicolor -require-accepted`: deferred
  until catalog integration. The current catalog entry has no motion source path,
  so this is not a frame-quality failure.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
ears, cropped feet, cropped tail, alpha-edge crop, or broken contact.

## Slot Review

- `00-03 idle`: pass. Seal bicolor Ragdoll read with stable baseline.
- `04-11 walk`: pass. Low walking frames keep the body, tail, and feet attached.
- `12-19 fast`: pass. Long low poses stay in-frame and grounded.
- `20-25 sniff/nibble`: pass. Nose-down frames remain complete with no prop.
- `26-31 species action`: pass. Action frames preserve seal face/ear/tail points.
- `32-39 turn`: pass. Turn/front-biased frames narrow and recover without crop.
- `40-43 eat`: pass. Low eat/action frames remain grounded.
- `44-47 ground check`: pass. Low frames keep face, feet, and tail inside frame.
- `48-51 alert/rest`: pass with note. Frames become fluffier and lighter but remain a seal bicolor Ragdoll.
- `52-55 groom`: pass with note. Groom/recovery poses keep one complete cat.
- `56-61 reaction/recover`: pass. Final frames stay grounded and in-frame.

## Visual Notes

- Species read: Ragdoll cat / seal bicolor.
- Scale: accepted; turn and low frames change width/height as pose rather than camera jump.
- Contact: reviewed frame bottoms stay consistently grounded; no visible floating.
- Style continuity: frames `47+` have a fluffier, slightly lighter paint read than earlier frames. This is accepted as watch-only because the seal face/ear/tail points, body silhouette, contact, and scale remain coherent.
- Reject checks: no text, border, scenery, prop, floor band, cast shadow, duplicate animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`ragdoll_seal_bicolor` is accepted as local source art only. It should remain
outside the runtime/current Pages list until its one-animal patch version is
planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/ragdoll-seal-bicolor-set00-full-contact.png`
- Overview contact: `contacts/ragdoll-seal-bicolor-overview-00-61-2x.png`
- Slot contacts: `contacts/ragdoll-seal-bicolor-*-4x.png`
- Candidate sheet: `contacts/ragdoll-seal-bicolor-source-set00-candidate-sheet.png`
- Animated preview: `contacts/ragdoll-seal-bicolor-set00-preview.gif`
- One-frame QA copy: `qa/ragdoll-seal-bicolor-oneframe-review.json`
- Run audit copy: `qa/ragdoll-seal-bicolor-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
