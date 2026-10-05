# Budgerigar green yellow set00 visual QA - 2026-06-28

## Verdict

Accepted as frame-addressable `set00` source art for `budgerigar_green_yellow`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Parent run: `docs/art-source/one-frame-method-fullrun-20260627/budgerigar-green-yellow-footlock-set00-oneframe-62/`
- Accepted frames: `docs/art-source/budgerigar/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/budgerigar/motion-source/sheets/budgerigar-green-yellow-source-set00.png`

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62 and exit status 0. The report lists
  greenish pixel counts because the bird plumage is naturally green/yellow; this
  is treated as non-blocking by parent visual QA.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=80.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=80.
- `cmd/assemblemotion`: wrote `budgerigar-green-yellow-source-set00.png` at 5952x64.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
head, cropped feet, or cropped tail after the foot-lock normalization repairs.

## Slot Review

- `00-03 idle`: pass. Yellow head, green body, black wing barring, blue cheek
  patch, long blue tail, and small feet read clearly.
- `04-11 walk`: pass. Walking steps keep a stable baseline; the thin contact
  strokes are feet, not a floor prop.
- `12-19 fast`: pass. Low stretched poses remain one complete budgerigar and do
  not introduce a separate ground line.
- `20-25 sniff/nibble`: pass after replacing the earlier cross-lane contaminated
  `frame-20`; low head poses stay attached and bird-like.
- `26-31 species action`: pass. Wing/turn emphasis keeps the same species,
  palette, and tail.
- `32-39 turn`: pass. Front-biased turn frames remain a single budgerigar and
  return to right-facing without wrong-species drift.
- `40-43 eat`: pass. Lower poses stay compact with no food prop or floor band.
- `44-47 ground check`: pass. Low inspection frames keep the tail and feet inside
  the frame.
- `48-51 alert/rest`: pass. Scale, palette, and baseline stay consistent.
- `52-55 groom`: pass. Larger posture shifts remain one complete bird with the
  same markings.
- `56-61 reaction/recover`: pass. Final frames recover to the same compact
  right-facing budgerigar.

## Visual Notes

- Species read: green/yellow budgerigar.
- Coat read: bright yellow head, green body, black wing scalloping, blue cheek
  patch, and blue tail remain visible across the set.
- Anatomy: compact parakeet body, small hooked beak, long tapering tail, and
  tiny attached feet remain readable.
- Reject checks: no text, border, scenery, cage, perch, food prop, floor,
  shadow, duplicate animal, cropped head, cropped feet, or cropped tail.

## Release Status

`budgerigar_green_yellow` is accepted as local source art only. It should remain
outside the runtime/current Pages list until its one-animal patch version is
planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/budgerigar-green-yellow-set00-full-contact.png`
- Overview contact: `contacts/budgerigar-green-yellow-overview-00-61-2x.png`
- Slot contacts: `contacts/budgerigar-green-yellow-*-4x.png`
- Candidate sheet: `contacts/budgerigar-green-yellow-source-set00-candidate-sheet.png`
- Animated preview: `contacts/budgerigar-green-yellow-set00-preview.gif`
- One-frame QA copy: `qa/budgerigar-green-yellow-oneframe-review.json`
- Run audit copy: `qa/budgerigar-green-yellow-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
