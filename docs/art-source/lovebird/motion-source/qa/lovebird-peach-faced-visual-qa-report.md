# Lovebird peach-faced set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `lovebird_peach_faced`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Parent run: `docs/art-source/one-frame-method-fullrun-20260628/lovebird-peach-faced-set00-oneframe-62/`
- Accepted frames: `docs/art-source/lovebird/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/lovebird/motion-source/sheets/lovebird-peach-faced-source-set00.png`

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62 and exit status 0. The report lists
  greenish pixel counts because the bird plumage is naturally green; this is
  treated as non-blocking by parent visual QA.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=73.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=73.
- `cmd/assemblemotion`: wrote `lovebird-peach-faced-source-set00.png` at 5952x64.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
head, cropped feet, or cropped tail.

## Slot Review

- `00-03 idle`: pass. Peach/orange face, pale beak, green body, darker wing, and
  short tail read clearly.
- `04-11 walk`: pass. Short walking steps keep feet attached and no perch,
  floor, or prop appears.
- `12-19 fast`: pass. Low stretched poses remain one complete lovebird with the
  same markings.
- `20-25 sniff/nibble`: pass. Low head poses remain attached to the body and no
  food prop or detached ground mark appears.
- `26-31 species action`: pass. Upright and front-biased poses remain a single
  peach-faced lovebird.
- `32-39 turn`: pass. Turn frames move through front-biased views and return to
  side-view without wrong-species drift or duplicate faces.
- `40-43 eat`: pass. Head and body stay connected with no floor band.
- `44-47 ground check`: pass. Low inspection frames keep feet and tail inside
  the frame.
- `48-51 alert/rest`: pass. Upright rest poses keep stable scale and palette.
- `52-55 groom`: pass. Groom-like posture changes remain one complete bird with
  peach face and green body intact.
- `56-61 reaction/recover`: pass. Final frames return to a stable right-facing
  lovebird.

## Visual Notes

- Species read: peach-faced lovebird.
- Coat read: peach/orange face and throat, pale beak, green body, darker wing,
  and short tail remain visible across the set.
- Anatomy: compact small parrot body, small beak, attached feet, and short tail
  remain readable.
- Reject checks: no text, border, scenery, cage, perch, food prop, floor,
  shadow, duplicate animal, cropped head, cropped feet, or cropped tail.

## Release Status

`lovebird_peach_faced` is accepted as local source art only. It should remain
outside the runtime/current Pages list until its one-animal patch version is
planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/lovebird-peach-faced-set00-full-contact.png`
- Overview contact: `contacts/lovebird-peach-faced-overview-00-61-2x.png`
- Slot contacts: `contacts/lovebird-peach-faced-*-4x.png`
- Candidate sheet: `contacts/lovebird-peach-faced-source-set00-candidate-sheet.png`
- Animated preview: `contacts/lovebird-peach-faced-set00-preview.gif`
- One-frame QA copy: `qa/lovebird-peach-faced-oneframe-review.json`
- Run audit copy: `qa/lovebird-peach-faced-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
