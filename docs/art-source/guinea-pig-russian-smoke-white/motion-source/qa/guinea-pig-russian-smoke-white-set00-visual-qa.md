# Guinea pig Russian smoke white set00 visual QA

Date: 2026-06-28

## Verdict

Accepted as a distinct white Russian-smoke guinea pig `set00` motion source
candidate and promoted to:

- `docs/art-source/guinea-pig-russian-smoke-white/motion-source/accepted-frames/set00/`
- `docs/art-source/guinea-pig-russian-smoke-white/motion-source/sheets/guinea-pig-russian-smoke-white-source-set00.png`

This is a separate source family from the existing `guinea_pig_tricolor`
asset.

## Mechanical QA

- Parent final `auditframes -strict -artifact-warnings -motion-warnings`:
  `valid=62 missing=0 invalid=0 warnings=119`.
- Promoted accepted-frame audit:
  `valid=62 missing=0 invalid=0 warnings=119`.
- `cmd/assemblemotion` wrote a `5952x64` source sheet.

The high warning count is expected for this preview pass because the broad
white low body repeatedly triggers lower-edge floor/ledge heuristics. Parent
review checked the flagged width, height, baseline, alpha-area, and fill-ratio
changes on `53-61`, `45-61`, and full `00-61` contact sheets.

## Visual Review

- Species read: guinea pig, not hamster, mouse, rat, or rabbit.
- Coat read: white Russian-smoke direction with gray ear/nose, black eye, and
  pale rounded body.
- Anatomy: broad tailless body, small attached feet, rounded head, and low
  guinea pig stance remain readable.
- Size and body-ratio review: late frames keep width `84`, baseline `58`, and
  fill ratio in a stable range. The low horizontal body is intentional guinea
  pig posture, not a detached floor or sudden body-size jump.
- Final `58-61` recovery frames form a low right-facing band and keep the same
  coat/species read.
- Reject checks: no text, border, scenery, food prop, costume, second animal,
  tail, cropped ears, cropped feet, detached floor, or shadow prop.

## Review Artifacts

- `contacts/guinea-pig-russian-smoke-white-set00-53-61-checker-contact.png`
- `contacts/guinea-pig-russian-smoke-white-set00-45-61-checker-contact.png`
- `contacts/guinea-pig-russian-smoke-white-set00-00-61-checker-contact.png`
- `accepted-frames/set00-auditframes-report.json`
- `accepted-frames/set00-assemblemotion-report.json`
