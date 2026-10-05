# Current Asset Production Flow

This is the compact source of truth for producing and approving one
AnimalsDesktop animal motion family. It separates ImageGen production,
mechanical checks, parent visual acceptance, local integration, and release
publication.

## Ownership and hard boundaries

- Use one Codex project task per animal or coat variant. The task owns only its
  run directory until parent approval.
- Generate every accepted production frame with one ImageGen call for that
  frame. Multi-cell grids and sheets are direction references only and must
  never be split into accepted production frames.
- Do not use SubAgents to generate accepted pixels. They may perform read-only
  research, evidence review, or adversarial QA.
- Do not paint, erase, patch, or script alpha pixels to rescue a failed frame.
  Regenerate that single frame with the same visual invariants instead.
- Go build caches such as `.gocache`, `.gomodcache`, `.gotmp`, `go-cache`,
  `go-mod-cache`, and `go-tmp` are disposable tool output, not provenance.
  Keep them outside the run when practical and never promote or commit them.
- Do not replace one species with a recolor of another. Coat variants may share
  a stable species silhouette and motion contract, but each accepted frame
  still needs its own source provenance.
- Local source-art acceptance does not authorize Pages publication, a version
  change, a tag, a push, downloadable packages, or a GitHub Release.

## Completion and runtime expansion

- The normal completion unit for one animal or coat variant is one canonical,
  accepted `set00`: 62 standalone `96x64` transparent source frames assembled
  into one `5952x64` source sheet.
- `cmd/importanimals` intentionally expands that canonical source sheet into
  runtime files `<sprite-base>_set00.png` through
  `<sprite-base>_set09.png`. The ten runtime slots preserve the existing
  DeguDesktop-compatible storage and animation contract; they are not a
  requirement for nine additional independent source families.
- A valid accepted canonical source therefore reports
  `source_sets=1`, `runtime_sets=10`, and `release_ready=true`. This expansion
  is normal production behavior and must not produce a preview or
  incompleteness warning.
- An optional independent `set00` through `set09` source family remains
  supported. If any optional source beyond `set00` is introduced, all ten must
  exist and be byte-unique. Partial or duplicated optional families fail
  validation; keep only canonical `set00` unless ten genuinely distinct motion
  sets are ready.
- Source-art completion does not by itself authorize publication. Runtime,
  behavior, build, page, packaging, and explicit release-approval gates still
  apply separately.

## Non-promotable same-species coat-sheet experiment

The user asked to test whether a stable base could support four coat-only poses
in one ImageGen sheet. The experiment proved useful for coat direction and
measurement, but the repository production rule remains one accepted frame per
ImageGen call. No sheet cell or derivative of a sheet cell is promotable.

The failed pilot paths remain part of the contract:

- supplying the previous animal sheet as a color reference copied pose
  information and dropped raw silhouette IoU to `72.984-78.957%`;
- a flat palette-only reference preserved raw geometry, but two ImageGen calls
  still differed in exposure;
- comparing only normalized output to a re-sheeted control was misleading
  because an unchanged control roundtrip itself averaged `97.946%` IoU and
  suffered one-pixel quantization in frames `06` and `10`.

The mechanically successful pilot used a pose-free palette reference plus one bounded,
uniform tone calibration per four-cell call. Its contiguous `04-11` evidence
had raw silhouette IoU `98.971%` mean / `98.899%` minimum, exact final bboxes,
final-vs-canonical IoU `98.845%` mean / `98.036%` minimum, final centroid
movement at most `0.278px`, base-normalized adjacent tone-ratio change at most
`3.02%`, a cross-call boundary change of `1.97%`, and shared matte `8/8` pass.

The direction-reference experiment may run only when all of these are true:

1. The same species' base motion family is already parent-approved at `62/62`.
   A new species, new silhouette, new pose, or replacement base frame stays on
   the one-frame-per-call path.
2. The coat variant has its own Codex project task and run directory. The base
   frames and their hashes are read-only references.
3. Each ImageGen call receives one fixed `1536x1024` 2x2 input containing at
   most four sequential target poses in exact `768x512` cells. There are no
   labels, dividers, dynamic crops, or inferred cell boundaries.
4. The secondary reference is a pose-free palette or interior-fur swatch from
   an independently reviewed coat identity. Never pass a prior animal sheet as
   the color reference.
5. The prompt changes pigmentation only and freezes silhouette, anatomy, limb
   count and placement, head/tail direction, baseline, scale, cell order,
   chroma background, and empty space.
6. Diagnostic slices may use only the declared cell coordinates. Record the input sheet,
   output sheet, prompt/call, swatch, cell coordinates, base-frame hash, raw
   cell hash, and final-frame hash. No dynamic cell detection may rescue a
   malformed output. Mark every slice `non-promotable`.
7. Establish one parent-reviewed target tone ratio from the first two pilot
   calls using a fixed canonical torso ROI and
   `candidate torso luma / base torso luma`. One uniform RGB gain may be
   applied to every subject pixel in a full 2x2 call only when the required
   gain is within `0.95-1.05`. The chroma background and alpha stay unchanged.
   Record the gain and pre/post hashes. Per-cell tuning, manual paint, alpha
   repair, selective masks, or a larger correction rejects the whole call.
8. A diagnostic normalization may run with `prepareframe -background chroma-green
   -match-alpha-bbox <base-frame>`. An ordinary bbox difference greater than
   one pixel is a hard reject, not permission to rescale the anatomy.
9. Require raw `768x512` silhouette IoU at least `98.5%` against the matching
   base cell, raw centroid movement at most `1.25px`, exact final bbox, final
   silhouette IoU at least `98.0%`, final centroid movement at most `0.30px`,
   and base-normalized adjacent/cross-call tone-ratio change at most `5%`.
   `coatbatch measure` and `coatbatch tone` enforce and report the raw geometry
   limits for all four cells, including fillers. `prepareframe
   -match-alpha-bbox` enforces and reports the final geometry limits before it
   writes either the diagnostic normalized PNG or its report.
   These remain the strict defaults. A reviewed recovery policy may change one
   bound only when `cmd/coatbatch` names and hard-codes the species and evidence
   binding. The current exceptions are:
   - `ferret_albino_high_contrast_v1`: `species=ferret_albino` only, raw IoU
     at least `98.0%`, with the same `1.25px` centroid limit.
   - `ferret_sable_exact_recovery_v1`: `species=ferret_sable` only, the exact
     approved input-sheet SHA-256
     `e16d004799e033405d1404ca318b18a0550d7a36db2d1bdefb99b7fd97c6fbf5`
     and target ratio `0.573025` only, with gain `0.85-1.05`.
   Unknown IDs, the wrong species, input hash, or target fail closed. There is
   no arbitrary numeric threshold override. A diagnostic produced before the
   named policy existed is not promotable; rerun the exact retained input
   through the diagnostic command path and all final geometry, matte, contact,
   tone-continuity, and runtime-loop gates. Passing still does not make a sheet
   cell promotable.
10. Require the shared matte audit, strict frame audit, light/dark/checker
    contacts, and a runtime-size animated loop. Reject visible anatomy, facing,
    coat-mask, brightness, foot-contact, tail, or dark-background readability
    loss even when numeric gates pass.
11. A failed target cell rejects all target cells from that ImageGen call.
    Do not mix selected cells from a tone-inconsistent call.
12. Do not extend this experiment into a 62-frame production run. Once the coat
    direction is approved, restart from the base as standalone one-frame
    ImageGen calls.

### Diagnostic coat-batch command path

`cmd/coatbatch` builds, measures, calibrates, and slices the bounded
non-promotable experiment. The retained Python
helpers under `.codex/experiments/ferret-sheet-recolor-contiguous8-20260724`
remain earlier experiment evidence. Neither path is an accepted-pixel
production entrypoint.

Every batch has one JSON manifest. All input paths are relative to that
manifest and must stay under its directory:

```json
{
  "species": "ferret_sable",
  "call": "frames-00-03-attempt-01",
  "prompt": "prompts/frames-00-03.txt",
  "swatch": "references/sable-palette.png",
  "cells": [
    {
      "id": "frame-00",
      "role": "target",
      "source": "base-raw/frame-00.png",
      "base_frame": "base-96x64/frame-00.png",
      "cell": 0,
      "output": "frame-00.png"
    },
    {
      "id": "frame-01",
      "role": "target",
      "source": "base-raw/frame-01.png",
      "base_frame": "base-96x64/frame-01.png",
      "cell": 1,
      "output": "frame-01.png"
    },
    {
      "id": "frame-02",
      "role": "target",
      "source": "base-raw/frame-02.png",
      "base_frame": "base-96x64/frame-02.png",
      "cell": 2,
      "output": "frame-02.png"
    },
    {
      "id": "frame-03",
      "role": "target",
      "source": "base-raw/frame-03.png",
      "base_frame": "base-96x64/frame-03.png",
      "cell": 3,
      "output": "frame-03.png"
    }
  ]
}
```

`call` is a unique operator-assigned call key, not a claim about an external
provider ID. `prompt` contains the exact ImageGen prompt. `swatch` is the
reviewed pose-free color reference. `source` is the matching approved base
source used to construct that cell; `base_frame` is the canonical transparent
`96x64` geometry reference. Reports bind the raw manifest, prompt, swatch,
source, and base-frame bytes with SHA-256.

The two reviewed recovery fields are optional. Omit both for the strict path.
An Albino manifest may name only its raw-geometry policy:

```json
{
  "raw_geometry_policy": "ferret_albino_high_contrast_v1"
}
```

The exact reviewed Sable manifest may instead bind its tone recovery:

```json
{
  "tone_policy": {
    "id": "ferret_sable_exact_recovery_v1",
    "approved_input_sha256": "e16d004799e033405d1404ca318b18a0550d7a36db2d1bdefb99b7fd97c6fbf5",
    "target_ratio": 0.573025
  }
}
```

Reports record the selected policy IDs, effective bounds, actual input hash,
approved hash, target, and effective gain.

Run the four modes in order:

```powershell
go run ./cmd/coatbatch -mode build `
  -manifest batch.json `
  -out base-sheet.png `
  -report build-report.json

go run ./cmd/coatbatch -mode measure `
  -manifest batch.json `
  -base-sheet base-sheet.png `
  -sheet imagegen-output.png `
  -report measure-report.json

go run ./cmd/coatbatch -mode tone `
  -manifest batch.json `
  -base-sheet base-sheet.png `
  -sheet imagegen-output.png `
  -target-ratio 0.573025 `
  -out calibrated-sheet.png `
  -report tone-report.json

go run ./cmd/coatbatch -mode slice `
  -manifest batch.json `
  -sheet calibrated-sheet.png `
  -out raw-cells `
  -report slice-report.json

go run ./cmd/prepareframe `
  -background chroma-green `
  -src raw-cells/frame-00.png `
  -match-alpha-bbox base-96x64/frame-00.png `
  -out normalized/frame-00.png `
  -report normalized/frame-00-report.json
```

The numeric ratio above is an example only. Use the one parent-reviewed ratio
established from that coat's first two diagnostic pilot calls. A source that is
already `768x512` is copied without filtering. Other source sizes are scaled
once with Go `CatmullRom`, and that fact is recorded per cell; this does not
claim byte equivalence with the retained Pillow/Lanczos experiment.

The hardened Go path was replayed against the retained Sable `04-07` diagnostic.
It measured group ratio `0.549837` against the earlier
rounded `0.549750`, selected gain `1.042172` against `1.042337`, and produced
only one-level RGB rounding differences from the retained calibrated result
(`7,010` subject pixels, maximum channel delta `1`; alpha/chroma unchanged).
The raw geometry replay passed at `98.902%` minimum IoU and `1.044px` maximum
centroid movement. All four final geometry locks matched their canonical
bboxes and passed at `98.195%` minimum IoU and `0.154px` maximum centroid
movement. This is a bounded requalification, not permission to skip per-batch
gates or promote a diagnostic slice.

`prepareframe` rejects source/output/reference/report path collisions, stages
the output and report together, and records SHA-256 for its source, geometry
reference, and normalized output. The geometry report records mask area,
intersection, union, IoU, centroid delta, and centroid distance without
rounding the values used for acceptance. The lock proves only alpha geometry;
anatomy, direction, contact, matte, tone continuity, and parent visual gates
remain mandatory.

When testing a two-target call, the manifest may still declare two `filler`
cells. A filler omits `output`, participates in build/measure/tone and content
checks, and is never written by `slice`. A malformed filler still rejects the
diagnostic call because it can corrupt the shared tone measurement.

The tool stages every output and report in its destination directory and
rolls back ordinary write or rename failures. Multi-file renames cannot be
power-loss atomic on every filesystem, so a process or machine crash can leave
hidden `.coatbatch-stage-*` or `.coatbatch-backup-*` recovery files. Treat
those as a stopped batch, preserve them for diagnosis, and do not resume or
promote until the intended destination hashes are reconciled against the
reports.

## Frame contract and motion schedule

One accepted family contains exactly 62 individually generated transparent
`96x64` PNGs named `frame-00.png` through `frame-61.png`.

| Frames | Motion intent |
| --- | --- |
| `00-03` | idle |
| `04-11` | slink |
| `12-19` | scurry |
| `20-25` | sniff |
| `26-31` | groom |
| `32-39` | turn |
| `40-47` | creep |
| `48-55` | rest |
| `56-61` | alert |

For right-facing source families, facing is part of the frame contract:
frames `00-31` face right, frames `32-39` visibly turn from right to left, and
frame `40` begins a separate right-facing action through frame `61`. Keep
`sourceFacingDirection` at `+1`; do not compensate for a bad source sequence
with a per-variant runtime flip.

The numbered blocks are independent runtime actions, not one contiguous
`00 -> 61` animation. On Windows, the turn plays `32 -> 39` once, holding each
frame for two ticks, then flips direction and resumes at walk frame `04`.
Consequently, `31 -> 32` and `39 -> 40` must not be used as mechanical
size-continuity gates. Review continuity inside `32 -> 39`, then compare
frame `39` with a horizontally mirrored frame `04` as the actual right-to-left
turn exit. Review frame `40` as the start of the independent creep loop.
Darwin currently uses idle `00-03`, walk frames `04`, `05`, and `07`, sniff
`20-25`, groom `26-31`, and rest `48-55` for ferrets. It does not currently
select frame `06` or ranges `08-19`, `32-47`, and `56-61`. These frames remain
required source coverage for Windows and for future Darwin motion parity; do
not interpret their Darwin reachability gap as an asset-acceptance waiver.

Production stops for parent review after these bounded gates:

```text
00-04, 05-12, 13-20, 21-28, 29-36, 37-44, 45-52, 53-61
```

At a transition gate, both neighboring actions must read clearly. A repeated
pose with only recolor, translation, or vertical bob is not a new action.

## Base/new-species per-frame and per-gate procedure

1. Before generating a range, freeze a manifest of the existing canonical
   frame names and SHA-256 hashes. Record its aggregate SHA-256 and canonical
   frame count.
2. Generate one candidate frame in one ImageGen call. Keep the prior approved
   camera, direction, scale, baseline, lighting, anatomy, and coat invariants.
3. Normalize through the repository's deterministic source conversion only.
   Reject and regenerate crop, anatomy, shadow/floor, pinhole, or matte
   failures; do not repair their pixels.
   Prompted pose-height or width ranges are generation guidance unless the
   parent gate explicitly derives them from a runtime or canonical geometry
   contract. Do not reject an otherwise clean candidate only for a one- or
   two-pixel miss when its within-action contacts and motion audit are smoother
   than the nominal target. Record the miss and let the assembled action decide.
4. Run `auditframes` with
   `-motion-boundaries 4,12,20,26,32,40,48,56` and run the shared matte
   audit. The boundaries suppress false comparisons between independent
   actions while retaining every within-action warning. The report records the
   effective boundary list. A mechanical pass is necessary but is not visual
   approval. Loop closures and the turn-to-walk exit still require runtime
   contacts because `auditframes` does not compare loop endpoints.
5. Build checker, light, and dark contacts from the canonical individual PNGs.
   A contact is never an accepted source. A generated multi-cell sheet is never
   an accepted source.
6. Save a gate status artifact, then stop. The parent independently checks the
   hashes, mechanical reports, all three contacts, species read, anatomy,
   baseline, crop, matte, and motion continuity before authorizing the next
   range.

Run-local contact or candidate helpers may stay with raw evidence while the
lane is active, but they are not a second acceptance authority. Before
promotion, retain only helpers that explain durable provenance or are still
referenced; otherwise leave them out of the accepted source family. The shared
repository auditors and their saved reports remain authoritative.

Use `qa/gate-<start>-<end>-status.json` (or an equivalent existing lane
artifact) to make the handoff auditable. It must identify:

- lane and gate;
- baseline frame count and baseline manifest SHA-256;
- accepted frame count and accepted manifest SHA-256;
- SHA-256 for mechanical reports and checker/light/dark contacts;
- lane verdict (`ready_for_parent` or `rejected`);
- parent verdict, approved manifest SHA-256, and approved contact SHA-256
  values.

Lane readiness never implies parent approval. Any change to a canonical frame
or a contact after approval changes its hash and invalidates that approval.

## Interrupted-task recovery

An interrupted Codex task may stop after ImageGen returned or after a canonical
frame was written but before the manifest and gate summary were updated. Do not
assume the last chat message describes the filesystem, and do not immediately
regenerate the frame.

Recover in this order:

1. Treat the assigned run directory as the current-state authority. Inventory
   canonical frames, raw/alpha candidates, exact prompts, prepare/audit/matte
   reports, contacts, rejects, and manifest rows.
2. Recheck the last parent-approved SHA-256 baseline before accepting any
   post-interruption file. If an approved-prefix file changed, stop instead of
   continuing production.
3. Match each post-baseline canonical frame to one generated candidate and one
   exact prompt. Inspect its original-size anatomy and three-background
   contact, then rerun the repository audit and shared matte audit.
4. Reconcile missing manifest rows from existing evidence. Never invent a
   generated-source path, call count, or verdict; record
   `unknown_after_interruption` for an irrecoverable field and expose it at the
   parent gate.
5. Keep rejected candidates and their reasons. Do not silently promote them or
   overwrite them with a retry.
6. Resume from the first genuinely absent frame only after the canonical
   prefix, manifest, and QA evidence agree. Generate no replacement for a
   passing canonical frame merely because its producing turn was interrupted.

The next gate summary must name the interruption and recovery checks. Parent
approval is required for the recovered range exactly as for an uninterrupted
range.

## Shared transparent-matte gate

Run the same Pillow-based audit for every chroma-green animal lane:

```powershell
python scripts/audit_frame_matte.py `
  --frames-dir <canonical-frames-dir> `
  --format json `
  --output <qa-dir>\matte-audit.json
```

The command scans top-level PNGs in stable filename order, requires `96x64`,
and exits nonzero on failure. A pixel enters the greenish-residue probe when
alpha is greater than `8`, green is at least `90`, and green exceeds both red
and blue by at least `25`.

Tiny antialias residue passes only when all of these are true for each frame:

- no more than `8` greenish pixels;
- every greenish pixel has at most one rounded 8-bit level of composited green
  excess:
  `(max(0, G - max(R, B)) * A + 127) // 255 <= 1`;
- the largest 8-connected greenish component is at most `2` pixels;
- parent inspection of checker, light, and dark contacts finds no visible
  green fringe, dot, or coherent matte streak.

The JSON report preserves both the unrounded maximum and the rounded 8-bit
level so values near the quantization boundary remain reviewable.

The script deliberately records
`visual_contact_review_required: true`: passing numbers cannot waive contact
review. A failed threshold or visible contact defect means regenerate the
individual frame, not paint its alpha.

## Promotion and local integration

Only after all eight gates and the final `62/62` parent review pass:

1. Copy the canonical individual frames and durable provenance/QA evidence to
   `docs/art-source/<family>/motion-source/accepted-frames/set00/`.
2. Run the full strict frame audit and shared matte audit again.
3. Assemble the `5952x64` source sheet locally from the 62 accepted PNGs.
4. Register the accepted source in the catalog and use targeted
   `importanimals -variant <id> -check` before any repository-writing import.
5. Run the targeted import, inspect its diff, then run the full importer twice
   at the integration boundary and require the second run to be deterministic.
6. Validate runtime motion, tests, vet, the Windows GUI build, Pages asset
   consistency, and `git diff --check`.
7. Keep the tracked family compact: accepted frames, source sheet, final
   mechanical/matte reports, parent verdict, three-background contacts, and a
   provenance index that maps each accepted frame to its prompt/call evidence.
   Do not commit raw retries, rejected candidates, local helper caches, or Go
   caches.
8. After the compact package is verified, remove only the exact disposable
   cache directories from the completed run. Move the remaining full run
   directory to a named external archive when its raw/rejected evidence should
   be retained, then confirm the repository has no untracked run output. Never
   use a broad recursive cleanup against the repository or drive root.

Promotion makes the animal available for local integration only. Do not change
release version text or publish Pages, tags, packages, or Releases without
explicit authorization in the current task.
