# Adding an Animal

This is the canonical integration path for a new AnimalsDesktop species,
breed, morph, or coat. It keeps source-art acceptance, local runtime preview,
public Pages, and release publication as separate decisions.

Use [current-asset-production-flow.md](current-asset-production-flow.md) for
the one-frame production gates, non-promotable coat-sheet experiment,
hash-backed parent approvals, contact review, and shared transparent-matte
policy that must be completed before this integration path begins.

## Prerequisites

- Go 1.24 or newer
- Python 3.10 or newer
- `python -m pip install -r requirements-dev.txt`

## Boundaries

- A new species or materially different silhouette needs its own ImageGen
  source family. Do not create it by recoloring another species.
- A coat variant may share a species only after that species' silhouette and
  `62/62` motion set are parent-approved. Every accepted coat frame still uses
  one ImageGen call. A generated multi-cell sheet is direction evidence only.
- Source-art acceptance does not automatically make an animal selectable.
- Runtime selection does not automatically update Pages, tags, downloads, or a
  GitHub Release. Those remain explicit release work.
- `cmd/importsheet` is the legacy Degu compatibility importer. Use
  `cmd/importanimals` for current animal families.

## 1. Prepare accepted source art

Create exactly 62 transparent `96x64` frames:

```text
docs/art-source/<family>/motion-source/accepted-frames/set00/
  frame-00.png
  ...
  frame-61.png
```

Run the mechanical gate and assemble the source sheet:

```powershell
go run ./cmd/auditframes `
  -frames-dir docs\art-source\<family>\motion-source\accepted-frames\set00 `
  -strict -artifact-warnings -motion-warnings `
  -motion-boundaries 4,12,20,26,32,40,48,56

python scripts\audit_frame_matte.py `
  --frames-dir docs\art-source\<family>\motion-source\accepted-frames\set00 `
  --format json `
  --output docs\art-source\<family>\motion-source\qa\matte-audit-set00.json

go run ./cmd/assemblemotion `
  -frames-dir docs\art-source\<family>\motion-source\accepted-frames\set00 `
  -out docs\art-source\<family>\motion-source\sheets\<sprite-base>-source-set00.png `
  -report docs\art-source\<family>\motion-source\qa\assemble-set00.json
```

Mechanical success is not visual acceptance. Review light, dark, and checker
contacts for species read, anatomy, baseline, crop, matte, and continuity
before registering the source as accepted.

For the bounded four-cell direction experiment, `cmd/coatbatch` can run
`build`, `measure`, `tone`, and diagnostic `slice` using the exact manifest and
command sequence in `current-asset-production-flow.md`. It binds the manifest,
prompt, call key, swatch, raw source, canonical base frames, generated sheets,
and slices with SHA-256. Every output remains non-promotable. Do not substitute
the retained experiment Python helpers or copy experiment pixels into accepted
source art.

The ordinary raw silhouette and tone bounds remain `IoU >= 98.5%`, centroid
movement `<= 1.25px`, and gain `0.95-1.05`. The production-flow document also
defines two evidence-bound recovery policies for this ferret run: an
Albino-only high-contrast raw-IoU floor and an exact-input Sable tone recovery.
They are named, fail-closed exceptions, not reusable numeric knobs. Any retained
diagnostic candidate remains non-promotable even after a successful replay.

For a standalone coat-variant frame of the same approved pose,
`prepareframe` can fail closed on a one-pixel normalization difference and
match the canonical output bbox:

```powershell
go run ./cmd/prepareframe `
  -background chroma-green `
  -src <candidate-source.png> `
  -out <candidate-frame.png> `
  -match-alpha-bbox <canonical-96x64-frame.png> `
  -report <prepare-report.json>
```

The option accepts only an ordinary output whose `X`, `Y`, width, and height
are each already within one pixel of the reference. Larger differences are
rejected, and an exact match is a pixel-identical no-op. The report records
the source/reference/output SHA-256 values, pre-lock and reference bboxes, and
whether resampling occurred. It also rejects the final output unless its alpha
mask has at least `98.0%` IoU and at most `0.30px` centroid movement from the
canonical frame. Source, output, reference, and report paths must all be
distinct. This is a normalization guard, not an anatomy repair: it does not
copy the reference alpha mask and cannot waive per-frame provenance, contact,
matte, motion-continuity, or parent visual review. It never makes a generated
sheet cell eligible for production.

## 2. Register the catalog entry

Update `internal/catalog/catalog.go`:

1. Add `SpeciesList`, picker-group, and default-motion-profile entries only
   when the species is new.
2. Add the source frame and `set00` sheet path constants.
3. Add one `acceptedMotionVariant(...)` entry with stable snake-case ID and
   `SpriteBase`, English and Japanese labels, breed/morph, color, popularity
   tier, and motion profile.
4. Do not add the ID to `runtimeVariantIDs` until the animal is approved for
   the selectable runtime roster.

Every importer run calls `catalog.Validate()`. It rejects duplicate or malformed
IDs, duplicate sprite bases, unknown species/profiles/statuses, unsafe paths,
incomplete metadata, invalid motion-source combinations, and invalid runtime
roster entries before assets are written.

## 3. Check only the new animal

Use the targeted check before writing runtime assets:

```powershell
go run ./cmd/importanimals -variant <variant_id> -check
go run ./cmd/validatemotion -variant <variant_id> -require-accepted
go test -buildvcs=false ./internal/catalog ./cmd/importanimals
```

`-check` runs the existing import pipeline in a temporary directory. It checks
the catalog entry, source PNG, motion-sheet dimensions, alpha/content, and
output generation without changing repository files.

## 4. Import only the new animal

```powershell
go run ./cmd/importanimals -variant <variant_id>
```

Targeted mode writes only that animal's normalized generated source and ten
runtime sprite sheets. It intentionally does not replace the aggregate
`seed-import-report.json` or `animalsdesktop-seed-preview.png` unless explicit
`-report` or `-preview` paths are supplied.

Confirm the targeted diff contains only the intended generated source and
`<sprite-base>_set00.png` through `<sprite-base>_set09.png`.

A single accepted `set00` source sheet is the normal complete source contract.
The importer expands it into all ten runtime sprite slots for compatibility,
without a preview or incompleteness warning. If any optional independent
`set01`-`set09` source is added, the family must contain all ten byte-unique
source sheets; a partial or duplicated family is rejected. Keep only canonical
`set00` unless ten genuinely different motion sources are ready.

## 5. Opt into runtime and release gates

After parent visual approval, add the variant ID once to `runtimeVariantIDs` at
the intended picker position. The Go catalog, Windows runtime, and page-asset
builder all consume that curated order; do not copy a second roster into tests
or scripts.

Run the full importer only at the integration/release boundary:

```powershell
go run ./cmd/importanimals
go run ./cmd/validatemotion -runtime-only -require-accepted
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go build -buildvcs=false -ldflags="-H=windowsgui" `
  -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

The no-flag importer remains the authority for the aggregate report and
preview. Run it twice and confirm the second run produces no asset diff.

Pages, version text, release workflows, tags, downloadable ZIPs, and GitHub
Releases are separate explicit-publication work. Do not infer publication
permission from local animal integration.

## Addition checklist

- [ ] New species uses its own source family; coat-only reuse is justified.
- [ ] `62/62` standalone source files pass mechanical and parent visual QA;
      no accepted file originates from a generated multi-cell sheet.
- [ ] Catalog metadata and paths pass `catalog.Validate()`.
- [ ] Targeted `importanimals -check` passes without repository writes.
- [ ] Targeted import changes only the intended animal assets.
- [ ] Runtime roster inclusion is an explicit parent decision.
- [ ] Full import is deterministic.
- [ ] Runtime validation, tests, vet, Windows build, and whitespace checks pass.
- [ ] Pages/release state is unchanged unless separately authorized.
