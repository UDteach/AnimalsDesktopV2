# Roborovski Hamster set00 Visual QA

Date: 2026-06-28

## Source

- Run directory: `docs/art-source/one-frame-method-fullrun-20260628/roborovski-hamster-set00-oneframe-62/`
- Accepted frames: `docs/art-source/roborovski-hamster/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/roborovski-hamster/motion-source/sheets/roborovski-hamster-source-set00.png`
- Review contacts:
  - `docs/art-source/roborovski-hamster/motion-source/contacts/roborovski-hamster-set00-full-light-4x.png`
  - `docs/art-source/roborovski-hamster/motion-source/contacts/roborovski-hamster-set00-full-dark-4x.png`
  - `docs/art-source/roborovski-hamster/motion-source/contacts/roborovski-hamster-set00-full-checker-4x.png`

## Mechanical QA

- Run-dir audit after final color correction: `valid=62 missing=0 invalid=0 warnings=79`
- Accepted-frame audit: `valid=62 missing=0 invalid=0 warnings=79`
- `assemblemotion`: passed, 62 frames assembled into source sheet.
- `validatemotion -variant roborovski_hamster -require-accepted`: passed with `accepted_source=true`, `release_ready=false` because only one source set exists.

## Color Correction

- Existing run-local correction had already normalized frames `29-52` using warm-coat RGB mean/std transfer.
- This pass added final corrections for frames `53-61`, then a stronger follow-up for frames `60-61`.
- Reports:
  - `docs/art-source/roborovski-hamster/motion-source/qa/color-correction-53-61-report.json`
  - `docs/art-source/roborovski-hamster/motion-source/qa/color-correction-60-61-report.json`

## Visual Review

- `00-03 idle`: pass. Roborovski white brow/cheek/belly pattern and compact round body are readable.
- `04-11 walk`: pass. Tiny feet remain attached; no tail-like appendage after retries.
- `12-19 fast movement`: pass. Low scurry is contained and does not introduce motion marks.
- `20-25 sniff / nibble`: pass. Forage pose is prop-free with compact body shape.
- `26-31 species action`: pass. Small alert lift stays hamster-like; frame `29` is a higher-detail but acceptable transition.
- `32-39 turn`: pass. Turn slot includes front/left-facing poses as expected; no duplicate animal or scenery.
- `40-43 chew / eat`: pass. Paw/muzzle gesture is subtle and attached.
- `44-47 ground check`: pass. Low head poses retain baseline and rounded body.
- `48-51 alert rest`: pass. Recovery stays compact.
- `52-55 face groom`: pass. Paw-to-face gesture remains small; color is within corrected band.
- `56-61 reaction / recover`: pass with note. Frames `60-61` remain slightly brighter and broader than early idle, but the final color correction reduces the yellow drift and the change is acceptable for a reaction/recover slot.

## Verdict

Accepted as source art for `roborovski_hamster` set00. This is source-art promotion only; runtime menu, GitHub Pages current list, release artifacts, tags, and public release links are unchanged.
