# Netherland Dwarf set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/`
- `docs/art-source/netherland-dwarf/motion-source/sheets/netherland-dwarf-chestnut-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop slot-based motion rule used for the existing preview-runtime
`set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- Parent run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=81.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=81.
- `cmd/assemblemotion`: wrote `netherland-dwarf-chestnut-source-set00.png` at 5952x64.

The warnings are accepted for this preview parity pass. They mostly reflect the
low body/contact baseline, compact feet, ear edge details, and fine alpha
regions that artifact warning mode flags. Parent visual review did not find a
persistent floor, shelf, scenery, prop, text, second animal, or cropped anatomy.

## Slot Review

- `00-03 idle`: pass. Compact dwarf rabbit body and short ears are readable.
- `04-11 walk`: pass. Hop/walk frames keep a low, rounded small-rabbit silhouette.
- `12-19 fast`: pass. Body stretch increases but remains rabbit-like.
- `20-25 sniff/nibble`: pass. Head-down poses stay prop-free.
- `26-31 species action`: pass. Rounded dwarf body remains consistent.
- `32-39 turn`: pass. Upright/turn poses keep short ears and compact proportions.
- `40-43 chew/eat`: pass. No food prop appears.
- `44-47 ground check`: pass. Baseline and small feet remain stable.
- `48-51 alert/rest`: pass. Alert poses keep a Netherland Dwarf read rather than a long-eared rabbit.
- `52-55 face groom`: pass. Grooming poses remain one rabbit with no duplicate head.
- `56-61 reaction/recover`: pass. Recovery frames preserve the compact short-eared silhouette.

## Visual Notes

- Species read: Netherland Dwarf rabbit, not a generic long-eared rabbit.
- Coat read: grayish chestnut / warm agouti body remains readable across the set.
- Anatomy: rounded compact body, short ears, small feet, short muzzle, and tiny tail remain readable.
- Reject checks: no text, border, scenery, food prop, costume, second animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`netherland_dwarf_chestnut` is accepted for current runtime-preview parity only.
It should remain `release_ready=false` until accepted `set00` through `set09`
source families exist.

## Review Artifacts

- Full contact: `contacts/netherland-dwarf-gray-chestnut-set00-full-contact.png`
- Candidate sheet: `contacts/netherland-dwarf-gray-chestnut-source-set00-candidate-sheet.png`
- Animated preview: `contacts/netherland-dwarf-gray-chestnut-set00-preview.gif`
- One-frame QA copy: `qa/netherland-dwarf-oneframe-review.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
