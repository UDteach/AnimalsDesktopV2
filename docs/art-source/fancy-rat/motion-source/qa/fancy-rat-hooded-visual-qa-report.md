# Fancy Rat Hooded set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/fancy-rat/motion-source/accepted-frames/set00/`
- `docs/art-source/fancy-rat/motion-source/sheets/fancy-rat-hooded-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=48.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=48.
- `cmd/assemblemotion`: wrote `fancy-rat-hooded-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They mostly correspond
to thin tail and foot/contact details being detected as pinholes, detached alpha
components, or lower-edge runs. Parent visual review did not find a persistent
floor, shelf, scenery, prop, text, or second animal.

## Slot Review

- `00-03 idle`: pass. Hooded head and white body are readable; tail is attached.
- `04-11 walk`: pass. Low rat body, small feet, and long thin tail remain intact.
- `12-19 fast`: pass. Body stretch increases but does not become mouse/gerbil.
- `20-25 sniff/nibble`: pass. Nose-down and upright sniff variants are prop-free.
- `26-31 species action`: pass. Tail remains inside frame; no cropped anatomy.
- `32-39 turn`: pass with note. Upright/angled frames read as the same hooded rat.
- `40-43 chew/eat`: pass. Head and forepaw motion stay subtle, with no food prop.
- `44-47 ground check`: pass. The rat remains complete and baseline stays plausible.
- `48-51 alert/rest`: pass. Alert posture is larger but still rat-like.
- `52-55 groom`: pass with note. Grooming raises the torso but remains one animal.
- `56-61 reaction/recover`: pass. Tail curve changes most here, but stays attached.

## Visual Notes

- Species read: hooded fancy rat, not a degu, mouse, or generic small rodent.
- Coat read: white body with dark hood over head and shoulders remains stable.
- Anatomy: pointed muzzle, rounded ears, slim attached tail, and small feet remain
  visible across the set.
- Reject checks: no text, border, scenery, food prop, costume, second animal,
  cropped ears, cropped feet, or cropped tail.

## Release Status

`fancy_rat_hooded` is accepted for current runtime-preview parity only. It should
remain `release_ready=false` until accepted `set00` through `set09` source
families exist.

## Review Artifacts

- Full contact: `contacts/fancy-rat-hooded-set00-full-contact.png`
- Candidate sheet: `contacts/fancy-rat-hooded-source-set00-candidate-sheet.png`
- Animated preview: `contacts/fancy-rat-hooded-set00-preview.gif`
- One-frame QA copy: `qa/fancy-rat-hooded-oneframe-review.json`
- Parent audit copy: `qa/fancy-rat-hooded-auditframes-parent.json`
