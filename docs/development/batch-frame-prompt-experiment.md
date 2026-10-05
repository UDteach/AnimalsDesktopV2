# Batch Frame Prompt Experiment

Date: 2026-06-22

## Purpose

Evaluate whether AnimalsDesktop asset production can safely move faster than one-frame-per-image generation without lowering release-quality sprite standards.

This is an experiment only. It does not override the current accepted-frame gate. Accepted frames must still end as standalone `96x64` transparent PNGs, pass `cmd/prepareframe`, `cmd/auditframes`, and parent light/dark visual review.

## Current Decision

- Adopted for ordinary testing: 2-frame batch only.
- Trial-only fallback: 4-frame batch, only after 2-frame testing is proven clean.
- Parent-requested exploration: reference-guided Flow/Gemini/ChatGPT grids at `16`, `31`, or `62` cells, only through the `animalsdesktop-flow-grid-batch` Skill and only under `review/`.
- Not adopted for production: direct 8+ frame sheets/grids, or any grid output that is not parsed into standalone reviewed PNGs.
- Not accepted as a guarantee: "96x64 per cell" grid wording. It can guide layout, but it does not prove clean cell boundaries, exact cell dimensions, or safe cutouts.

## Rationale

The main production risks are not raw generation speed; they are rejection rate, cell-boundary noise, labels, dividers, checker/background residue, source-family drift, face collapse, disconnected alpha, and manual cutout overhead.

Large grids make the parent QA problem harder because one bad cell can spoil a batch and the model may trade individual frame quality for sheet composition. If cutting or cleanup is needed, the time saved by batch generation can disappear.

The current high-count grid experiment is therefore deliberately separated from accepted-frame production. It uses a prebuilt `2048x2048` placement grid, accepted-frame style anchors, pure chroma-green background, fixed-cell parsing, and normal `prepareframe`/contact-sheet/`auditframes` review after extraction.

Batch escalation is conservative: start at `16`, move to `31` only after clean `16` placement and visual QA, and treat `62` at `2048x2048` as exploration unless smaller grids have already shown zero severe defects.

For rough outline artifacts, prefer prevention over local repainting:

- render uploaded style anchors and grid seeds with CatmullRom scaling, not nearest-neighbor enlargement
- prefer transparent `style-anchors.png` for style reference, and use green `grid-seed.png` only as a placement cue so the model does not learn green edge color
- if transparent anchors render on black, checker, or with noisy edge specks, use `style-anchors-neutral.png` instead of green-matted anchors
- run `auditframes -artifact-warnings` before anchor reuse; do not use frames with `transparent pinholes` warnings as style anchors
- archive generated packs if they were built from a later-quarantined anchor
- start with the minimal `grid-guide.png` cell-outline guide; keep `grid-guide-detailed.png` for placement failures because internal guide lines are more likely to be reproduced
- request smooth antialiased source-size outlines with no jagged cutout, white fringe, green halo, noisy fur-edge pixels, or rough matte
- force pure flat `#00ff00` empty areas so `prepareframe -background chroma-green` can remove the background and despill green
- reject outputs that still show rough matte, halo, edge noise, green dot residue, or transparent pinholes on enlarged light/dark review
- do not locally smooth an accepted sprite silhouette unless a separate parent review explicitly approves that altered output

Current quarantine: `chinchilla set00/frame-00..43`, `frame-45`, and `frame-51` reported `transparent pinholes` and must not be used as style anchors. Flow-grid packs generated from `frame-00` were moved under `docs/art-source/review/flow-grid-experiments/_archive/pinhole-anchor-frame00-20260623/`.

## Flow Nano Banana 2 / Pro Comparison

2026-06-23 check:

- Google Flow opened successfully in the in-app browser at `https://labs.google/fx/tools/flow`.
- The account showed `PRO`.
- A new Flow project was created.
- Image generation mode exposed both `Nano Banana 2` and `Nano Banana Pro`.
- Square `1:1` / `1x` settings were available.
- Flow displayed `0 credits will be used` for generation in the selected image mode.
- Browser-level paste/type did not activate Flow's internal prompt state; the create button stayed `aria-disabled=true`.
- OS-level virtual key input did activate the prompt field and allowed generation.
- Nano Banana Pro was tested with `style-anchors.png` and `grid-seed.png`, square output, `x4`, 2K downloads. Outputs are archived in `docs/art-source/review/flow-grid-experiments/chinchilla-set00-16-anchor-frame50-clean/`.
- Result: visual style and face quality were stable, but all four 16-cell sheets mostly duplicated the same pose. Original guide manifest parsed `0/16`. After using an actual `512x512` grid and snapping green-dominant pixels to pure chroma green, the best candidate parsed `11/16`; remaining rejects were disconnected components and chroma pinholes. The workflow is not accepted.
- A follow-up Nano Banana Pro prompt variant with explicit per-cell limb deltas improved pose variation, but still reproduced visible guide lines and introduced many matte pinholes. Two 1024px cache outputs were recovered for diagnostics: one parsed `3/16` only after green normalization, the other parsed `0/16`. This is classified as `GUIDE_INK` plus `MAT_PINHOLE`, not a usable workflow.

Current interpretation:

- Do not use browser copy/paste as the Flow automation path. Use focused OS virtual key input, or ask the parent to type one short activation key sequence before automation continues.
- `Nano Banana 2` remains untested visually and is still the first Flow candidate for fast motion-idea batches: `walk`, `scurry`, and `turn`, starting at `16` cells.
- `Nano Banana Pro` is not good enough for this visible-grid 16-cell motion batch prompt. Keep it as a hard single-frame comparison candidate: eat, ground check, groom, reaction, and clean golden anchors. If grid batching is retried, change the guide strategy first: no visible grid seed, an 8-frame invisible walk grid, a smaller 4-frame strip, or independent images.
- Any future Flow output must be copied into a review directory, parsed or prepared as a standalone PNG, and compared against ImageGen through the same pinhole, face, source-family, baseline, and parent visual gates.

## Candidate Prompt A: Separate Images

Use this when the generator can return multiple independent images instead of one sheet.

```text
Generate exactly 2 separate transparent PNG images, not a sheet, collage, grid, strip, or comparison panel.

Each image contains exactly one complete side-view [ANIMAL] performing [FRAME INTENT].
Keep identical camera distance, body scale, proportions, contact baseline, facing direction, coat colors, lighting, outline style, and source family across both images.

Transparent background only. No floor, shadow, scenery, props, text, labels, borders, dividers, checkerboard, background residue, duplicate animal, cropped anatomy, disconnected alpha, or stray pixels.

Face quality guard: clean coherent eyes, muzzle, mouth, and nose; no distorted eye, giant eye, smeared muzzle, collapsed mouth/nose, face panel, red/open mouth, black mask, or mouth bar.
```

## Candidate Prompt B: Two-Frame Strip

Use only for experiments under `review/`, never directly for accepted frames.

```text
Create exactly one transparent 2-frame horizontal sprite strip for review, with two logical cells of 96x64 pixels each, total intended layout 192x64.

Each logical cell contains exactly one complete side-view [ANIMAL] performing the next two smooth motion frames: [FRAME A INTENT] and [FRAME B INTENT].
Keep identical camera distance, body scale, proportions, contact baseline, facing direction, coat colors, lighting, outline style, and source family across both cells.

No visible grid lines, no dividers, no labels, no numbers, no text, no borders, no floor, no shadow, no scenery, no props, no checkerboard, no background residue, no cell overlap, no pixels crossing between cells, no cropped ears/feet/tail/toes, no transparent holes, no disconnected alpha, and no stray pixels.

Face quality guard: clean coherent eyes, muzzle, mouth, and nose; no distorted eye, giant eye, smeared muzzle, collapsed mouth/nose, face panel, red/open mouth, black mask, or mouth bar.
```

## Adoption Gate

2-frame batching can replace one-frame generation only if an experiment batch proves all of the following:

- At least 100 candidate frames reviewed across multiple animals or motion blocks.
- Zero severe defects: labels, dividers, cropped animals, cell overlap, checker/background residue, floor/shadow bands, face collapse, disconnected alpha, or source-family breaks.
- Accepted-frame pass rate is no more than 3 percentage points worse than one-frame generation.
- Total parent+child time per accepted frame falls by at least 25%.
- Every accepted output still passes as a standalone `96x64` PNG after preparation and visual review.

## Immediate Reject Gate

Reject the batch strategy and continue one-frame generation if any of these happen:

- A generated strip/sheet contains visible cell borders, labels, numbers, text, or grid lines.
- Any animal crosses a cell boundary or is cropped by the intended cell.
- Cutting introduces alpha chips, checker residue, background residue, floor/shadow bands, or transparent holes.
- Face-collapse rate increases compared with one-frame generation.
- Manual cutting/cleanup makes the total time saving less than 20%.

After one clear batch/grid failure, do not immediately try a larger grid or the same prompt again. Ask the in-app ChatGPT Pro thread for a short failure review first, including the source prompt, layout, parse report, light/dark contact sheets, reject reason, and whether the next attempt should return to one-frame-per-PNG.

## Current Operating Rule

Until the adoption gate is met, production threads remain on one-frame-per-PNG output. Batch outputs may be created only as review experiments and must not be placed directly under `accepted-frames`.

## Parsing Strategy

The parsing strategy must be conservative. Do not infer cut positions from the artwork when the output is meant to become production source.

Ranked options:

1. Hybrid parse plus audit, driven by a manifest.
2. Multiple independent downloaded images with manifest order and hashes.
3. Strict fixed-size slicing for `192x64` two-frame strips only.
4. Metadata or filename conventions as a helper, never as sole proof.
5. Gutter or cell-boundary detection as a reject check only.
6. Alpha connected-component segmentation as an anomaly detector only, not as the cutter.

### Preferred Path: Independent Images

When the generator can return two independent PNGs, this is the preferred experiment path.

Required checks:

- exact expected image count
- stable manifest order, not UI download order alone
- source SHA-256 for every image
- valid PNG decode
- individual `cmd/prepareframe`
- individual `cmd/auditframes`
- parent light/dark neighbor contact sheet
- visual source-family, face, baseline, scale, and motion-delta review

### Conditional Path: Strict `2x1` Strip

Use only when the source image is already exactly `192x64` and intended as two `96x64` cells.

Allowed fixed cuts:

- frame A: `x=0..95`, `y=0..63`
- frame B: `x=96..191`, `y=0..63`

Reject the whole strip immediately if:

- canvas dimensions are not exactly `192x64`
- the image must be resized to become `192x64`
- any visible divider, label, number, grid line, or border exists
- any alpha exists on or crosses the `x=96` boundary guard
- either cut frame has cropped anatomy, edge-touching content, disconnected alpha, transparent holes, floor/shadow bands, or background residue
- either cut frame fails standalone `prepareframe`, `auditframes`, or parent light/dark review

### Trial-Only Path: Strict `2x2` Grid

Only after the `2x1` path proves safe. Use exact `192x128` only. Cut positions are fixed at `x=0..95`, `x=96..191`, `y=0..63`, and `y=64..127`.

Do not use dynamic cell detection to rescue malformed grids.

### Reject-Only Helpers

Alpha connected-component analysis is useful for finding detached islands, stray pixels, holes, and merged animals. It is not safe for deciding crop boxes, because tails, whiskers, toes, or adjacent cells can be misclassified.

Gutter detection is useful for proving expected boundaries are empty. It is not safe for discovering boundaries, because transparent body gaps and tail curves can look like gutters.

### Minimal Pipeline Proposal

Use the experiment-only parser, separate from accepted-frame import:

```powershell
go run ./cmd/parsebatchframes `
  -manifest docs/art-source/<animal>/motion-source/review/<batch>/batch.json `
  -out docs/art-source/<animal>/motion-source/review/<batch>/parsed `
  -report docs/art-source/<animal>/motion-source/review/<batch>/parse-report.json
```

The parser accepts this manifest shape:

```json
{
  "animal": "gecko",
  "sequence": "set01",
  "layout": "strip-2x1",
  "cell_width": 96,
  "cell_height": 64,
  "guard_pixels": 1,
  "frames": [
    {
      "id": "set01-frame-08",
      "source": "batch-strip.png",
      "cell": 0
    },
    {
      "id": "set01-frame-09",
      "source": "batch-strip.png",
      "cell": 1
    }
  ]
}
```

For `layout: "independent"`, omit `cell` and point each frame at a separate source PNG. For `layout: "grid-2x2"`, `cell` is `0..3` in reading order. Optional `output` may be used, but it must be a filename only, not a path.

For the Flow/Gemini reference-grid experiment, the parser also accepts `layout: "grid-fixed"` with `columns`, `rows`, `origin_x`, `origin_y`, `pitch_x`, `pitch_y`, and `background: "chroma-green"`. The generated `batch.json` from `cmd/flowgridtemplate` owns those coordinates; do not hand-edit them to rescue a malformed sheet.

The parser:

- rejects output or report paths under `accepted-frames`
- rejects unexpected canvas dimensions for fixed strips/grids
- slices only from manifest-declared layouts: `independent`, `strip-2x1`, or `grid-2x2`
- additionally slices `grid-fixed` review sheets from manifest-declared fixed cells when produced by `cmd/flowgridtemplate`
- rejects fixed strip/grid sources with alpha in the configured boundary guard
- rejects fixed strip/grid cells whose alpha touches a cell edge
- rejects no-alpha, no-transparent-background, and disconnected-alpha outputs
- writes parsed PNGs only to the requested review output directory
- computes source and parsed SHA-256
- writes a JSON report with dimensions, cell bounds, alpha-boundary checks, connected-component summaries, and reject reasons
- never writes directly to `accepted-frames`

Current parser QA coverage:

- valid independent image parsing and reporting
- valid exact `strip-2x1` slicing
- valid exact `grid-2x2` slicing
- wrong fixed-layout canvas dimensions are rejected
- alpha in the `2x1` vertical boundary guard is rejected
- alpha in the `2x2` horizontal boundary guard is rejected
- fixed-layout cell-edge alpha is rejected
- disconnected alpha inside a parsed fixed-layout cell is rejected
- output/report paths under `accepted-frames` are rejected
- manifest `output` path escape is rejected

## Parser Design Notes

An in-app ChatGPT Pro parser follow-up on 2026-06-23 produced a usable advisory, but it is treated as brainstorming rather than proof. Its recommendation aligned with the local fail-closed strategy: independent images plus manifest and hybrid audit are safest; exact `2x1` fixed slicing can remain review-only; `2x2` is trial-only; dynamic crop, gutter detection, and alpha-component segmentation should not rescue malformed sheets. The current parser decision still relies on the local parser smoke tests, real `2x1` gecko visual rejections, and the fail-closed strategy below.

Current recommendation:

- Production remains one-frame-per-PNG.
- The separate same-species coat-only experiment in
  `current-asset-production-flow.md` may use fixed four-cell ImageGen output
  only after a `62/62` base family is approved. It uses high-resolution
  declared cells, base-geometry matching, bounded sheet-wide tone calibration,
  and per-cell provenance for direction evidence. Its cells and derivatives
  remain non-promotable; every accepted production frame still comes from its
  own ImageGen call.
- Safest review experiment: `independent` multi-image outputs with manifest order, source hashes, and hybrid audit, because no cell boundary has to be trusted.
- Conditional review experiment: exact `strip-2x1` only when the source canvas is exactly `192x64`, the manifest declares both cells, and both parsed cells pass all reject checks.
- Trial-only experiment: exact `grid-2x2` at `192x128`, after `2x1` quality is proven. Do not test `8`, `16`, or `32` frame grids for production acceptance.

Parser behavior should stay fail-closed:

- Use fixed manifest-declared cell bounds for strip/grid parsing.
- Use alpha/component detection only as a rejection signal, not as a dynamic crop or rescue mechanism.
- Do not infer cell positions from animal artwork for production candidates.
- Do not auto-trim transparent padding into a new crop, because that can hide baseline/camera drift and make neighbor comparison harder.
- Template/gridline detection may be useful only to reject visible dividers, labels, or frame borders.

Parent-review report should stay small and decision-oriented:

- source path, source SHA-256, output path, output SHA-256
- declared layout, source dimensions, expected dimensions, and cell bounds
- per-cell alpha bbox, edge-touch flags, boundary-guard alpha flags, connected-component count and largest component ratio
- pass/fail reject reasons in plain text
- neighbor contact sheet paths for light and dark backgrounds
- explicit note that parser success is not visual acceptance

Immediate reject rules for parsed outputs:

- source-family drift, blurry/soft source-family degradation, or changed camera/scale/baseline
- distorted eye, giant eye, smeared muzzle, collapsed mouth/nose, face panel, red/open mouth, black mask, or mouth bar
- alpha touching fixed cell edges, any alpha in strip/grid boundary guards, detached alpha chips, foot-like fragments, or animal crossing a cell boundary
- checker/noisy background, floor/shadow, scenery, labels, dividers, borders, or text
- disconnected alpha components, stray pixels, cropped anatomy, lower shelves, edge bars, transparent holes, green dot residue, or missing-pixel pinholes

Next worthwhile experiment, if any:

1. Try an 8-frame walk-only experiment before returning to 16 cells: no visible grid seed, `style-anchors.png` only, invisible `2x4` layout text prompt.
2. If a sheet is visually coherent but not extractable, run a second layout-correction pass: exact `2048x1024`, invisible `2x4`, `512x512` per cell, at least `80px` padding, and no green holes. This is promising with ChatGPT ImageGen; the first 2026-06-23 correction pass improved diagnostic parse to `5/8`, and a fresh stricter chinchilla workflow run improved to `7/8` after green normalization.
3. If the 8-frame sheet still shows guide ink, boundary residue, pinholes, or layout drift, try eight independent reference-guided PNGs instead of a sheet.
4. Parse under `review/`, never directly under `accepted-frames`.
5. Compare against the current single-frame acceptance rate and parent review time.
6. Adopt only if quality stays equal and accepted-frame time falls materially; otherwise keep single-frame production.

### 2026-06-23: Chinchilla 8-frame ImageGen workflow run 01

Artifact directory:

```text
docs/art-source/review/flow-grid-experiments/chinchilla-set00-8-imagegen-workflow-20260623-run01/
```

Result:

- Stage 1 text-only ImageGen produced a clean 8-component `2x4` sheet.
- Stage 2 stricter layout-correction prompt also produced a clean 8-component `2x4` sheet with better centering.
- Raw parse stayed `0/8` for both sheets because the green background was not parser-pure.
- Diagnostic green normalization improved stage 1 to `6/8` and stage 2 to `7/8`.
- Local cell recentering preserved the same counts: stage 1 `6/8`, stage 2 `7/8`.
- The remaining stage 2 reject was `frame-04`, a single `1px` `MAT_PINHOLE`.

Interpretation:

The 8-frame ImageGen workflow is close mechanically: layout and component separation are now good enough to inspect. It is not production-ready because the art is still too realistic for the accepted sprite family, pose deltas are weak, and green normalization/pinhole repair remain diagnostic-only. The next useful loop should keep the same 8-cell scale, make the style flatter and more sprite-like, and use a true reference-upload layout-correction pass when the browser lane is available and parent approval covers the upload.

### 2026-06-23: ChatGPT Pro prompt consultation and reproducibility run 02

Artifact directory:

```text
docs/art-source/review/flow-grid-experiments/chinchilla-set00-8-imagegen-repro-20260623-run02/
```

Setup:

- Consulted ChatGPT Pro with the prior `7/8` result, pinhole failure, realistic-style drift, and weak gait delta notes.
- Tested four prompt strategies twice each with Codex ImageGen:
  - `A`: flatter 2D sprite style.
  - `B`: stronger gait deltas.
  - `C`: extraction/matte cleanliness first.
  - `D`: local hybrid of `A` and `C`.

Result:

- All eight samples produced 8 detected animal components.
- Raw fixed-cell parse stayed `0/8` for every sample because the green background was still not parser-pure.
- Diagnostic green-normalized parse:
  - `A`: `6/8`, `7/8`
  - `B`: `8/8`, `7/8`
  - `C`: `8/8`, `8/8`
  - `D`: `4/8`, `6/8`

Interpretation:

`C` is the first prompt with reproducible `8/8` diagnostic extraction. `B` is useful for walk-cycle readability, but less stable. `A` improves sprite style but not extraction. `D` visually looks closer to the target mascot style but currently hurts pinhole/parse stability. The next prompt should use `C` as the base and add only modest `B` gait wording; do not adopt the `D` hybrid until extraction stays stable.

### 2026-06-23: ImageGen 8-frame run 03-05 follow-up

Artifact directories:

```text
docs/art-source/review/flow-grid-experiments/chinchilla-set00-8-imagegen-repro-20260623-run03/
docs/art-source/review/flow-grid-experiments/chinchilla-set00-8-imagegen-repro-20260623-run04/
docs/art-source/review/flow-grid-experiments/chinchilla-set00-8-imagegen-repro-20260623-run05/
```

Setup:

- `E`: C extraction-cleanliness base plus modest B gait wording and frame-by-frame walk sequence.
- `F`: removed gait wording and returned to near-identical clean review poses.
- `G`: angle-table prompt with attached-foot rules, informed by a local pose blueprint.

Result:

- `E` generated coherent 8-animal sheets, but diagnostic green-normalized parse dropped to `3/8`, `5/8`, `5/8`, and `3/8`. Rejects were mainly `transparent/chroma pinholes`, plus occasional detached alpha fragments.
- `F` reduced animation ambition but became near-static and sometimes upright. Diagnostic green-normalized parse was `7/8`, `5/8`, `5/8`, and one fixed-cell layout failure.
- `G` made angle-table prompting testable. Fixed-cell green-normalized parse still failed because the model did not preserve cell boundaries, but local recentering reached `7/8` and `8/8`.

Interpretation:

Adding walk-specific wording directly to an 8-frame ImageGen sheet is currently counterproductive: it improves pose intent but creates belly, foot, and tail matte defects. Angle tables are useful as a motion-design/specification tool, especially when each species has its own allowed-moving-parts profile, but they do not make 8-frame sheets production-safe. Keep ImageGen sheets as review-only layout, pose, and style references. Production accepted frames stay one-frame-per-PNG.

New flow:

1. Ask a species pose profile before prompting: movable parts, locked identity parts, risky anatomy, allowed angle/phase changes, and banned words.
2. For motion blocks, create a bone-only or silhouette-only blueprint locally to define phase and angle values.
3. Convert the blueprint into text prompts for review-only sheet experiments.
4. Use successful sheet poses only as references for one-frame-per-PNG production, not as direct accepted source frames.

## Trial Results

### 2026-06-22: ChatGPT Pro `2x1` Gecko Trial

Artifact directory:

```text
docs/art-source/gecko/motion-source/review/batch-imagegen-trial-20260622-2155/
```

Result:

- ChatGPT Pro produced an exact `192x64` PNG from a text-only two-frame gecko prompt.
- `cmd/parsebatchframes` accepted the strip: `parsed=2 rejected=0`.
- Parsed frames had one connected alpha component each.
- `cmd/auditframes` on parsed outputs returned `valid=2 missing=60 invalid=0 warnings=0`.
- Parent light/dark comparison against accepted gecko `set01/frame-06..07` rejected the result for source-family drift.

Interpretation:

The strict `2x1` parser path is viable as a mechanical workflow, but text-only batch prompting is not enough to preserve the accepted gecko family. Do not adopt `2x1` batching for production yet.

### 2026-06-22: ChatGPT Pro `2x1` Gecko Trial 2

Artifact directory:

```text
docs/art-source/gecko/motion-source/review/batch-imagegen-trial2-20260622-2205/
```

Result:

- A stricter text-only prompt also produced an exact `192x64` PNG.
- `cmd/parsebatchframes` accepted the strip: `parsed=2 rejected=0`.
- Parsed frames had one connected alpha component each.
- `cmd/auditframes` on parsed outputs returned `valid=2 missing=60 invalid=0 warnings=0`.
- Parent light/dark comparison again rejected the result for source-family drift.

Interpretation:

The second prompt reduced the oversized-head problem but still failed on mottled texture, tail/body silhouette, and leg/toe style. Text-only `2x1` batch prompting is not reliable enough for gecko production. Further `2x1` work should require reference-image context or switch back to independent single-frame generation.

After parsing, the existing flow remains authoritative:

```powershell
go run ./cmd/prepareframe -src <parsed.png> -out <review-prepared.png> -report <prepared-report.json>
go run ./cmd/auditframes -frames-dir <review-prepared-dir> -artifact-warnings
.codex/scripts/make-frame-contact-sheet.ps1 -Frames <neighbors and prepared outputs> -OutPrefix <review-sheet>
```

Promotion to `accepted-frames` remains a parent decision after visual review.
