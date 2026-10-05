# Roborovski hamster set00 visual QA

Date: 2026-06-28

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/roborovski-hamster/motion-source/accepted-frames/set00/`
- `docs/art-source/roborovski-hamster/motion-source/sheets/roborovski-hamster-source-set00.png`

## Mechanical QA

- Parent final `auditframes -strict -artifact-warnings -motion-warnings`:
  `valid=62 missing=0 invalid=0 warnings=81`.
- Promoted accepted-frame audit:
  `valid=62 missing=0 invalid=0 warnings=81`.
- `cmd/assemblemotion` wrote a `5952x64` source sheet.

The warnings are review candidates, not automatic failures. Parent visual
review checked the flagged width, height, baseline, alpha-area, and fill-ratio
changes on `53-61`, `45-61`, and full `00-61` contact sheets.

## Color Correction

The user called out that the Roborovski tone should be corrected. Frames
`29-52` had usable anatomy and scale but were darker/grayer than the early
sandy frames. The parent pass preserved original frames under:

- `qa/color-correction-20260628/`

Accepted correction:

- Method: warm-coat masked RGB mean/std transfer toward frames `00-28`.
- Alpha, silhouette, white belly/cheeks/eyebrow, and pink ear/foot/nose pixels
  were preserved.
- Correction report:
  `qa/color-correction-20260628-report.json`.

This repair is accepted because it only normalizes coat continuity. It does not
turn another species into a Roborovski hamster and does not alter pose,
contact, or silhouette.

## Visual Review

- Species read: Roborovski dwarf hamster, not Syrian hamster, mouse, or rat.
- Coat read: sandy back with white belly, cheek, and eyebrow patches.
- Anatomy: compact round body, rounded ears, black eye, tiny pink feet, and no
  long tail remain readable.
- Size and body-ratio review: `45-52` and `53-61` include turn/upright/reaction
  bands with width changes, but parent review found these are continuous pose
  changes rather than isolated body-size spikes.
- Final `59-61` recovery frames are brighter and wider, but they form a
  continuous reaction/recover band and keep the same species read.
- Reject checks: no text, border, scenery, food prop, costume, second animal,
  cropped ears, cropped feet, long tail, detached floor, or shadow prop.

## Review Artifacts

- `contacts/roborovski-hamster-set00-53-61-checker-contact.png`
- `contacts/roborovski-hamster-set00-45-61-checker-contact.png`
- `contacts/roborovski-hamster-set00-00-61-checker-contact.png`
- `accepted-frames/set00-auditframes-report.json`
- `accepted-frames/set00-assemblemotion-report.json`
