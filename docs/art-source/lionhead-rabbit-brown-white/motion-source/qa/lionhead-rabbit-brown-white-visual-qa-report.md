# Lionhead rabbit brown white set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for
`lionhead_rabbit_brown_white`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Source run: `docs/art-source/one-frame-method-fullrun-20260628/lionhead-rabbit-brown-white-set00-oneframe-62/`
- Accepted frames: `docs/art-source/lionhead-rabbit-brown-white/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/lionhead-rabbit-brown-white/motion-source/sheets/lionhead-rabbit-brown-white-source-set00.png`

## Mechanical QA

- Source run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Source run `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=106.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=106.
- `cmd/assemblemotion`: wrote `lionhead-rabbit-brown-white-source-set00.png`.
- `cmd/validatemotion -variant lionhead_rabbit_brown_white -require-accepted`:
  deferred until catalog integration. The variant is not registered in the
  catalog yet, so this is not a frame-quality failure.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
ears, cropped feet, cropped tail, alpha-edge crop, or broken contact.

## Slot Review

- `00-03 idle`: pass. Mostly white Lionhead rabbit read with brown eye and ear markings.
- `04-11 walk`: pass. Walking poses keep the compact rabbit grounded.
- `12-19 fast`: pass. Faster poses preserve the fluffy mane/ruff and brown back marks.
- `20-25 sniff/nibble`: pass. Nose-forward poses remain complete and prop-free.
- `26-31 species action`: pass with note. Frame `26` is narrower as a front-biased action pose, then returns to side-view scale without baseline drift.
- `32-39 turn`: pass. Turn frames narrow and recover naturally with ears and feet in-frame.
- `40-43 eat`: pass. Chew/eat poses keep the requested white body and brown markings.
- `44-47 ground check`: pass. Low inspection poses remain grounded and complete.
- `48-51 alert/rest`: pass with note. Width increases at `51` as the pose returns from front-biased alert to side-view; this is pose-driven, not a camera scale jump.
- `52-55 groom`: pass with note. Grooming frames narrow at `53` as the forepaw/face pose changes; contact and scale remain coherent.
- `56-61 reaction/recover`: pass. Final frames stay stable with non-empty adjacent diffs.

## Visual Notes

- Species read: Lionhead rabbit / lion-rabbit with visible mane/ruff.
- Requested coat read: mostly white rabbit with brown eye patch, brown ear
  markings, and small brown back markings matching the supplied direction.
- Scale: accepted. Mechanical width/area watches at `26 -> 27`, `50 -> 51`, and
  `52 -> 53` were reviewed visually and accepted as pose changes with unchanged
  baseline.
- Contact: reviewed frame bottoms stay consistently grounded; no visible floating.
- Reject checks: no text, border, scenery, prop, floor band, cast shadow,
  duplicate animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`lionhead_rabbit_brown_white` is accepted as local source art only. It should
remain outside the runtime/current Pages list until its one-animal patch version
is planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/lionhead-rabbit-brown-white-set00-full-contact.png`
- Overview contact: `contacts/lionhead-rabbit-brown-white-overview-00-61-2x.png`
- Slot contacts: `contacts/lionhead-rabbit-brown-white-*-4x.png`
- Candidate sheet: `contacts/lionhead-rabbit-brown-white-source-set00-candidate-sheet.png`
- Animated preview: `contacts/lionhead-rabbit-brown-white-set00-preview.gif`
- One-frame QA copy: `qa/lionhead-rabbit-brown-white-oneframe-review.json`
- Run audit copy: `qa/lionhead-rabbit-brown-white-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
