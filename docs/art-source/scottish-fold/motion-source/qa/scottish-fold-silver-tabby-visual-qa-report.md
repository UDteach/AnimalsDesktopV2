# Scottish Fold silver tabby set00 visual QA - 2026-06-29

## Verdict

Accepted as frame-addressable `set00` source art for `scottish_fold_silver_tabby`.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this promotion.

## Source Run

- Source run: `docs/art-source/one-frame-method-fullrun-20260628/scottish-fold-silver-tabby-set00-oneframe-62/`
- Accepted frames: `docs/art-source/scottish-fold-silver-tabby/motion-source/accepted-frames/set00/`
- Source sheet: `docs/art-source/scottish-fold-silver-tabby/motion-source/sheets/scottish-fold-silver-tabby-source-set00.png`

## Mechanical QA

- Source run has 62/62 normalized candidate frames.
- `oneframe_run.py review`: frame_count=62, issues=[].
- Source run `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=37.
- Promoted accepted-frame `cmd/auditframes`: valid=62, missing=0, invalid=0, warnings=37.
- `cmd/assemblemotion`: wrote `scottish-fold-silver-tabby-source-set00.png`.

The warning count is accepted for this source-art pass. Parent visual review did
not find persistent floor, shadow, scenery, prop, text, second animal, cropped
ears, cropped feet, cropped tail, alpha-edge crop, or broken contact.

## Slot Review

- `00-03 idle`: pass. Right-facing Scottish Fold read with folded ears and stable baseline.
- `04-11 walk`: pass. Steps keep the silver tabby silhouette grounded.
- `12-19 fast`: pass. Motion remains in-frame without sudden scale jumps.
- `20-25 sniff/nibble`: pass. Lower head poses keep body contact and tabby markings coherent.
- `26-31 species action`: pass. Low stretch/action poses remain grounded and complete.
- `32-39 turn`: pass. Turn/front-turn narrowing is accepted as pose change rather than scale error.
- `40-43 eat`: pass. Body and head stay connected with no duplicate face.
- `44-47 ground check`: pass. Low posture remains grounded and returns upward cleanly.
- `48-51 alert/rest`: pass. Folded ears, stripe pattern, and side-facing read remain stable.
- `52-55 groom`: pass. Face-grooming and front-leg stretch remain coherent and in-frame.
- `56-61 reaction/recover`: pass. Recovery frames keep stable contact; `58` is a low crouch, and `59-61` return to normal stance.

## Visual Notes

- Species read: Scottish Fold cat / silver tabby.
- Scale: accepted; turn frames narrow naturally, and recovery frames return without abrupt enlargement.
- Anatomy: folded ears, short rounded face, attached legs, visible tail, and tabby body pattern remain readable.
- Contact: reviewed frame bottoms stay consistently grounded; no visible floating.
- Reject checks: no text, border, scenery, prop, floor band, cast shadow, duplicate animal, cropped ears, cropped feet, or cropped tail.

## Release Status

`scottish_fold_silver_tabby` is accepted as local source art only. It should
remain outside the runtime/current Pages list until its one-animal patch version
is planned, cataloged, built, and verified.

## Review Artifacts

- Full contact: `contacts/scottish-fold-silver-tabby-set00-full-contact.png`
- Overview contact: `contacts/scottish-fold-silver-tabby-overview-00-61-2x.png`
- Slot contacts: `contacts/scottish-fold-silver-tabby-*-4x.png`
- Candidate sheet: `contacts/scottish-fold-silver-tabby-source-set00-candidate-sheet.png`
- Animated preview: `contacts/scottish-fold-silver-tabby-set00-preview.gif`
- One-frame QA copy: `qa/scottish-fold-silver-tabby-oneframe-review.json`
- Run audit copy: `qa/scottish-fold-silver-tabby-run-auditframes-report.json`
- Parent audit copy: `accepted-frames/set00-auditframes-report.json`
