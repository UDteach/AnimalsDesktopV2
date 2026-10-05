# Parrotlet blue/green set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `parrotlet_blue_or_green`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Source run: `docs/art-source/one-frame-method-fullrun-20260628/parrotlet-blue-green-set00-oneframe-62/`
- Accepted frames: `docs/art-source/parrotlet/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/parrotlet/motion-source/sheets/parrotlet-blue-green-source-set00.png`

## Mechanical QA

- Source run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Source run `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=116.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=116.
- `cmd/assemblemotion`: wrote `parrotlet-blue-green-source-set00.png`.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
beak, cropped feet, cropped tail, alpha-edge crop, or broken contact.

## Slot Review

- `00-03 idle`: pass. Right-facing compact parrotlet read with stable baseline.
- `04-11 walk`: pass. Walking frames keep feet attached and the small bird silhouette grounded.
- `12-19 fast`: pass. Lower and faster poses stay in-frame with consistent blue/green plumage.
- `20-25 sniff/nibble`: pass. Low head poses remain complete; the return at `25` reads as posture change, not scale jump.
- `26-31 species action`: pass. Upright/action poses keep the same bird identity and palette.
- `32-39 turn`: pass. Front-facing turn frames narrow as expected, then return to side-view without duplicate face or wrong-species drift.
- `40-43 eat`: pass. Head-lowered frames remain grounded with no food prop or floor band.
- `44-47 ground check`: pass. Low inspection frames keep beak, feet, and tail inside the frame.
- `48-51 alert/rest`: pass. Rest poses recover into stable side-view with consistent scale.
- `52-55 groom`: pass. Groom-like head turn stays attached and in-frame.
- `56-61 reaction/recover`: pass. Final recovery frames return to a stable right-facing parrotlet; replacement `57` and final `61` match the surrounding scale/contact.

## Visual Notes

- Species read: blue/green parrotlet.
- Scale: accepted; width changes are pose-driven, with the final reaction/recover sequence staying within the surrounding visual size.
- Anatomy: small hooked beak, compact body, short tail, attached feet, and wing markings remain readable.
- Contact: reviewed frame bottoms stay consistently grounded; no visible floating.
- Color and style continuity: blue/cyan body and darker wing/tail markings remain consistent enough for source-art acceptance.
- Reject checks: no text, border, scenery, perch, cage, prop, floor band, cast shadow, duplicate animal, cropped beak, cropped feet, or cropped tail.

## Release Status

`parrotlet_blue_or_green` is accepted as local source art only. It should remain
outside the runtime/current Pages list until its one-animal patch version is
planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/parrotlet-blue-green-set00-full-contact.png`
- Overview contact: `contacts/parrotlet-blue-green-overview-00-61-2x.png`
- Slot contacts: `contacts/parrotlet-blue-green-*-4x.png`
- Candidate sheet: `contacts/parrotlet-blue-green-source-set00-candidate-sheet.png`
- Animated preview: `contacts/parrotlet-blue-green-set00-preview.gif`
- One-frame QA copy: `qa/parrotlet-blue-green-oneframe-review.json`
- Run audit copy: `qa/parrotlet-blue-green-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
