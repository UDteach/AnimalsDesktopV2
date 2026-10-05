# parrotlet_green set00 frames 00-04 QA notes

## Result

Frames `00-04` are candidate parent-gate frames. No frame `05+` was generated.

## Rejected attempts

None. Attempts `01` for frames `00-04` were kept as parent-gate candidates.

## Raw / normalization note

All raw ImageGen files arrived as RGB PNGs with a baked checker-style transparency preview. The raw files are preserved under `raw/`. A run-local normalizer removed the preview background, kept the largest foreground bird component, and produced standalone transparent 96x64 PNGs under `frames/`.

## Visual review

- Species read: all five frames read as compact green Pacific parrotlet / マメルリハ candidates at 96x64.
- Exclusions: no yellow budgie head, black budgie barring, blue cheek spots, cockatiel crest, orange cheek patch, Java sparrow red beak/white cheek pattern, peach/orange lovebird face mask, bulky lovebird head, or macaw proportions observed.
- Composition: one complete right-facing bird per frame, short tail, pale beak, one dark eye, visible attached feet, no text, border, scenery, props, floor plane, cast shadow, duplicate bird, or cropped body parts.
- Continuity: scale, baseline, body volume, head size, eye size, beak size, and lighting are stable across `00-04`.
- Parent-gate caveat: frame `04` is a subtle grounded foot shift rather than a strong walk/perch step. It is mechanically valid and visually coherent, but the parent may request a stronger step before continuing.

## Mechanical QA

- Dimension/alpha check: `qa/frame-alpha-dimensions.json`; all five frames are `96x64` with non-empty alpha.
- Normalization report: `qa/normalize-report.json`.
- `auditframes` command:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260627/parrotlet-green-set00-oneframe-62/frames -strict -artifact-warnings -report docs/art-source/one-frame-method-fullrun-20260627/parrotlet-green-set00-oneframe-62/qa/auditframes-report.json
```

Result: `valid=5 missing=57 invalid=0 warnings=5`; exit status `1` is expected because frames `05-61` are intentionally missing for the early parent gate. Warnings are lower-content alpha runs/pinholes around tiny feet/baseline details; visual contact review did not show detached feet, floor plane, or body holes.

## Contact sheets

- `contact/parrotlet-green-set00-00-04-light-contact.png`
- `contact/parrotlet-green-set00-00-04-dark-contact.png`
- `contact/parrotlet-green-set00-00-04-checker-contact.png`

## 05-12 continuation

Parent approved frames `00-04`, then frames `05-12` were generated as the next bounded lane. Existing canonical frames `00-04` were not overwritten.

- Canonical range now present: `frames/frame-00.png` through `frames/frame-12.png`.
- Raw and alpha-prepared attempts for `05-12` are preserved under `raw/` and `alpha/`.
- Visual read: `05-12` remain compact green parrotlet candidates with short tail, pale beak, one dark eye, attached feet, and stable baseline/contact. No lovebird, budgerigar, cockatiel, Java sparrow, macaw, duplicate, crop, text, scenery, prop, floor plane, or cast-shadow drift was seen in contact review.
- Parent-gate caveat: frame `07` is slightly taller in head/body silhouette, and frame `12` is wider from the brisk-step lean. Both still read as the same compact parrotlet but should be parent-confirmed before continuing to `13+`.

Additional contact sheets:

- `contact/parrotlet-green-set00-05-12-light-contact.png`
- `contact/parrotlet-green-set00-05-12-dark-contact.png`
- `contact/parrotlet-green-set00-05-12-checker-contact.png`
- `contact/parrotlet-green-set00-00-12-light-contact.png`
- `contact/parrotlet-green-set00-00-12-dark-contact.png`
- `contact/parrotlet-green-set00-00-12-checker-contact.png`

Latest `auditframes` result:

```text
valid=13 missing=49 invalid=0 warnings=11
```

The command exited nonzero because frames `13-61` are intentionally missing at this stop gate. Warnings are lower-content alpha runs/pinholes around tiny feet/baseline details; visual contact review did not show detached feet, floor plane, or body holes.

## 13-20 continuation

Parent approved frames `00-12`, then frames `13-20` were generated as the next bounded lane. Existing canonical frames `00-12` were not overwritten.

- Canonical range present after this pass: `frames/frame-00.png` through `frames/frame-20.png`.
- Visual read: `13-17` and `19` maintain the compact green parrotlet scale, color, foot/contact, and short-tail identity.
- Parent-gate caveat recorded before approval: frame `18` is wider from the brisk movement pose, and frame `20` is a low sniff/nibble transition posture.
- Latest `auditframes` result for this pass: `valid=21 missing=41 invalid=0 warnings=15`.

## 21-28 continuation

Parent approved frames `00-20`, then frames `21-28` were generated as the next bounded lane. Existing canonical frames `00-20` were not overwritten.

- Canonical range now present: `frames/frame-00.png` through `frames/frame-28.png`.
- Raw and alpha-prepared attempts for `21-28` are preserved under `raw/` and `alpha/`.
- Visual read: `21-25` continue the low sniff/nibble/forage posture from frame `20`; `26-28` transition back toward species-action head/beak gestures while preserving compact green parrotlet color, short tail, attached feet, and baseline/contact.
- Parent-gate caveat: `21-25` are intentionally low and elongated; `27-28` are more upright action transition poses and should be parent-confirmed before continuing to `29+`.

Additional contact sheets:

- `contact/parrotlet-green-set00-21-28-light-contact.png`
- `contact/parrotlet-green-set00-21-28-dark-contact.png`
- `contact/parrotlet-green-set00-21-28-checker-contact.png`
- `contact/parrotlet-green-set00-00-28-light-contact.png`
- `contact/parrotlet-green-set00-00-28-dark-contact.png`
- `contact/parrotlet-green-set00-00-28-checker-contact.png`

Latest `auditframes` result:

```text
valid=29 missing=33 invalid=0 warnings=25
```

The command exited nonzero because frames `29-61` are intentionally missing at this stop gate. Warnings are lower-content alpha runs/pinholes around tiny feet/baseline details; visual contact review did not show detached feet, floor plane, body holes, crop, or wrong-species drift.

## Full 00-61 source-frame QA

Parent final QA passed the full 62-frame draft source for `parrotlet_green`.

- Canonical range present: `frames/frame-00.png` through `frames/frame-61.png`.
- Final `cmd/auditframes` result: `valid=62 missing=0 invalid=0 warnings=42`.
- Final audit report: `qa/auditframes-report-00-61.json`.
- Final contact sheets:
  - `contact/parrotlet-green-set00-53-61-light-contact.png`
  - `contact/parrotlet-green-set00-53-61-dark-contact.png`
  - `contact/parrotlet-green-set00-53-61-checker-contact.png`
  - `contact/parrotlet-green-set00-45-61-light-contact.png`
  - `contact/parrotlet-green-set00-45-61-dark-contact.png`
  - `contact/parrotlet-green-set00-45-61-checker-contact.png`
  - `contact/parrotlet-green-set00-00-61-light-contact.png`
  - `contact/parrotlet-green-set00-00-61-dark-contact.png`
  - `contact/parrotlet-green-set00-00-61-checker-contact.png`

Parent contact review found stable compact green parrotlet identity, visible feet/contact, no pale/gray drift, no cropped tail, and no one-frame body-size spike. This run is ready for parent promotion review only; no catalog, runtime, Pages, release, tag, download, commit, or push action was performed.
