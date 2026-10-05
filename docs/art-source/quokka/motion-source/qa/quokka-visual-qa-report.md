# Quokka set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `quokka`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Parent combined run: `docs/art-source/one-frame-method-fullrun-20260628/quokka-parent-combined-set00-oneframe-62/`
- Original source frames used: `00-13` from `docs/art-source/one-frame-method-fullrun-20260628/quokka-set00-oneframe-62/`
- Repair source frames used: `14-61` from `docs/art-source/one-frame-method-fullrun-20260628/quokka-stylelock-repair-set00-oneframe-62/`
- Accepted frames: `docs/art-source/quokka/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/quokka/motion-source/sheets/quokka-source-set00.png`

## Mechanical QA

- Parent combined run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Parent combined run-dir `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=21.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=21.
- `cmd/assemblemotion`: wrote `quokka-source-set00.png`.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
ears, cropped feet, cropped tail, or alpha-edge crop.

## Slot Review

- `00-03 idle`: pass. Right-facing compact quokka read with stable baseline.
- `04-11 walk`: pass. Small stepping poses keep scale and contact consistent.
- `12-19 fast`: pass. The transition from original `13` to repaired `14` has a mild palette shift, but body size, contact, and quokka anatomy stay coherent.
- `20-25 sniff/nibble`: pass. Low head/body poses remain complete and grounded.
- `26-31 species action`: pass. Retried repair frames fixed the earlier tall/size jump.
- `32-39 turn`: pass. Turn poses remain in-frame and do not introduce a second face or wrong species read.
- `40-43 eat`: pass. Body stays compact with attached feet and visible tail.
- `44-47 ground check`: pass. Low posture is intentional and remains grounded.
- `48-51 alert/rest`: pass. Fur texture is a little stronger than early frames, but scale, anatomy, and contact stay stable.
- `52-55 groom`: pass. Grooming poses remain coherent and in-frame.
- `56-61 reaction/recover`: pass. Final frames keep stable bbox and contact, with no alpha-edge crop.

## Visual Notes

- Species read: quokka / クアッカ.
- Scale: accepted after parent repair review for the earlier size-jump problem.
- Anatomy: rounded body, small rounded ears, visible muzzle, attached feet, and tail remain readable.
- Color continuity: frames `14-61` are slightly more muted/warm-brown than `00-13`; this is accepted because the motion reads as one quokka set and the repair fixes the larger scale/contact issue.
- Reject checks: no text, border, scenery, prop, floor band, cast shadow, duplicate animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`quokka` is accepted as local source art only. It should remain outside the
runtime/current Pages list until its one-animal patch version is planned,
cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/quokka-set00-full-contact.png`
- Overview contact: `contacts/quokka-overview-00-61-2x.png`
- Slot contacts: `contacts/quokka-*-4x.png`
- Candidate sheet: `contacts/quokka-source-set00-candidate-sheet.png`
- Animated preview: `contacts/quokka-set00-preview.gif`
- One-frame QA copy: `qa/quokka-oneframe-review.json`
- Run audit copy: `qa/quokka-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
