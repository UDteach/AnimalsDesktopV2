# Post-Budgerigar Asset Queue - 2026-06-27

## Scope

This is a background preparation lane for the next asset wave. It does not own
the active `budgerigar_green_yellow` frame run, runtime promotion, GitHub Pages
publication, release links, tags, or GitHub Releases.

Lane A owns budgerigar frame generation and promotion. This lane only records
the queue, prompt scaffolds, QA gates, and serialization rules for work after
budgerigar.

## Evidence Checked

- `docs/art-source/priority-animals-20260627-accepted-source-report.md`
- `docs/art-source/asset-queue.md`
- `docs/development/iteration-log.md`
- `internal/catalog/catalog.go`
- `scripts/build_page_assets.py`
- `scripts/verify_page_release.py`
- `docs/index.html`
- `docs/assets/upcoming-silhouettes/`

## Current State

- Public runtime and Pages current animals remain release-scoped at 16.
- `runtimeVariantIDs` excludes `chipmunk_striped`, `gecko_leopard`,
  `whites_tree_frog_blue`, `cockatiel_normal_gray`, and
  `java_sparrow_normal`.
- Accepted source already exists for `cockatiel_normal_gray` and
  `java_sparrow_normal`; these are not image-generation blockers.
- `budgerigar_green_yellow` is still in progress and should not be touched by
  this lane.
- Pages/verifier currently encode an 18-card Coming Soon order that starts with
  chipmunk and small-mammal variants before frog/gecko/birds. That public order
  should stay unchanged until the parent release/page lane explicitly updates
  it.

## Next-Asset Queue

This queue follows the current delegated priority: finish budgerigar first, then
birds, then popular cats/dogs. Earlier accepted-source requests for chipmunk,
blue White's tree frog, and leopard gecko remain release-ready candidates but
should not silently override the newer bird-first production priority.

| Order | Variant | Status | Next Lane Type | Notes |
| --- | --- | --- | --- | --- |
| 0 | `budgerigar_green_yellow` | Lane A active | Serialized frame production | Do not edit or reclaim Lane A run paths. |
| 1 | `cockatiel_normal_gray` | Accepted `set00` source | Runtime/page/release promotion prep | Source exists; validate and promote only after parent release decision. |
| 2 | `java_sparrow_normal` | Accepted `set00` source | Runtime/page/release promotion prep | Source exists; validate and promote only after parent release decision. |
| 3 | `lovebird_peach_faced` | New source candidate | Isolated source child lane | First new bird source lane after accepted bird promotion decisions. |
| 4 | `parrotlet_green` | New source candidate | Isolated source child lane | Good filler if the parent allows another bird lane. |
| 5 | `scottish_fold_silver_tabby` | Catalog source variant | New source child lane | Popular-cat wave candidate; needs cat-stalk motion identity. |
| 6 | `ragdoll_seal_bicolor` | Catalog source variant | New source child lane | Popular-cat wave candidate; preserve point markings and fluffy tail. |
| 7 | `shiba_inu_red` | Catalog source variant | New source child lane | Popular-dog wave candidate; avoid large-dog proportions. |

Release-ready exceptions:

- `chipmunk_striped` has accepted source and should remain distinct from
  `albino_chipmunk`.
- `whites_tree_frog_blue` has accepted source; high chroma-adjacent warning
  counts require visual contact review, not count-only rejection.
- `gecko_leopard` has accepted source; low-crawler anatomy, four attached
  underside legs, rounded toe pads, and thick tail remain required review gates.

These three can be promoted in a future runtime/version lane when the parent
chooses a release boundary, but they are not the next frame-generation lane
while the current instruction says birds first.

## Prompt Scaffolds

Use one exact animal per lane. Each child lane must write only under its
assigned run directory and stop after frames `00-04` for parent gate.

### `cockatiel_normal_gray`

Primary recommendation: do not regenerate. Use existing accepted source under
`docs/art-source/cockatiel/motion-source/` unless parent visual QA finds a real
defect.

Fallback prompt if a regeneration lane is explicitly opened:

```text
Generate cockatiel_normal_gray set00 frames 00-04 only.
Write only under docs/art-source/one-frame-method-fullrun-YYYYMMDD/cockatiel-normal-gray-set00/.
Do not edit catalog, runtime sprites, docs page, release files, or other animals.
Use one-frame ImageGen with a complete right-facing 2D sprite cockatiel on transparent or chroma green background.
Keep a normal gray cockatiel: gray body, pale face, yellow crest, orange cheek patch, long narrow tail, small hooked beak, two feet on one baseline.
Keep camera, scale, baseline, crest size, cheek patch, tail length, and contact points consistent.
Reject/regenerate frames with text, scenery, props, floor, shadows, duplicate animals, cropped crest/tail/feet, parrot-like body drift, rodent silhouette, or left-facing side-view unless the slot is a turn.
Stop and report after frames 00-04 for parent gate, then continue only after parent approval.
```

Acceptance criteria:

- Reads as a cockatiel at 96x64, not a generic parrot or sparrow.
- Crest and orange cheek patch remain visible across motion slots.
- Bird-hop profile stays bird-like; no rodent scurry, dog trot, or exaggerated
  vertical bob.
- Tail is not cropped and does not become a second body component.

QA checklist:

- `go run ./cmd/validatemotion -variant cockatiel_normal_gray -require-accepted`
- `go run ./cmd/auditframes -frames-dir docs/art-source/cockatiel/motion-source/accepted-frames/set00 -strict -artifact-warnings`
- Contact-sheet review by set00 slots, especially frames `12-19`, `26-31`,
  `32-39`, and `56-61`.

### `java_sparrow_normal`

Primary recommendation: do not regenerate. Use existing accepted source under
`docs/art-source/java-sparrow/motion-source/` unless parent visual QA finds a
real defect.

Fallback prompt if a regeneration lane is explicitly opened:

```text
Generate java_sparrow_normal set00 frames 00-04 only.
Write only under docs/art-source/one-frame-method-fullrun-YYYYMMDD/java-sparrow-normal-set00/.
Do not edit catalog, runtime sprites, docs page, release files, or other animals.
Use one-frame ImageGen with a complete right-facing 2D sprite Java sparrow on transparent or chroma green background.
Keep a normal Java sparrow: compact gray body, black head cap, white cheek patches, red beak, small red feet, short tail, clean side-view posture.
Keep camera, scale, baseline, head cap, cheek patches, beak color, and foot contact points consistent.
Reject/regenerate frames with text, scenery, props, floor, shadows, duplicate animals, cropped beak/tail/feet, finch identity drift, rodent silhouette, or left-facing side-view unless the slot is a turn.
Stop and report after frames 00-04 for parent gate, then continue only after parent approval.
```

Acceptance criteria:

- Black head, white cheeks, red beak, and gray body remain readable at 96x64.
- Small feet stay attached and do not become detached alpha debris.
- Bird-hop profile stays compact and grounded.
- Fine alpha warnings are accepted only after visual review confirms no body
  breakage.

QA checklist:

- `go run ./cmd/validatemotion -variant java_sparrow_normal -require-accepted`
- `go run ./cmd/auditframes -frames-dir docs/art-source/java-sparrow/motion-source/accepted-frames/set00 -strict -artifact-warnings`
- Contact-sheet review by set00 slots, especially tiny feet/beak frames and
  turn/reaction slots.

### `lovebird_peach_faced`

This is the safest independent new source lane after budgerigar/cockatiel/Java
sparrow decisions because it is another small bird but does not share accepted
source paths.

Proposed run directory:

```text
docs/art-source/one-frame-method-fullrun-YYYYMMDD/lovebird-peach-faced-set00/
```

Initial child prompt:

```text
Generate lovebird_peach_faced set00 frames 00-04 only.
Write only under docs/art-source/one-frame-method-fullrun-YYYYMMDD/lovebird-peach-faced-set00/.
Do not edit catalog, runtime sprites, docs page, release files, or other animals.
Use one-frame ImageGen with a complete right-facing 2D sprite peach-faced lovebird on transparent or chroma green background.
Keep a peach-faced lovebird: compact small parrot body, green wings/body, peach-orange face and forehead, pale beak, short rounded tail, two small feet on one baseline.
Keep camera, scale, baseline, head size, beak, face patch, wing color, tail length, and contact points consistent.
Reject/regenerate frames with text, scenery, props, floor, shadows, duplicate animals, cropped beak/tail/feet, budgerigar stripes, cockatiel crest, Java sparrow cheek pattern, oversized macaw/parrot proportions, or left-facing side-view unless the slot is a turn.
Stop and report after frames 00-04 for parent gate, then continue only after parent approval.
```

Acceptance criteria:

- Reads as a small peach-faced lovebird, not budgerigar, cockatiel, Java
  sparrow, macaw, or generic green parrot.
- Peach face patch and green body stay stable across frames.
- Short rounded tail stays inside the 96x64 frame.
- Bird-hop motion remains compact with grounded feet.

QA checklist:

- Run-local `oneframe_run.py review` or equivalent manifest/contacts review.
- `go run ./cmd/auditframes -frames-dir <accepted-frame-dir> -strict -artifact-warnings`
- `go run ./cmd/assemblemotion -frames-dir <accepted-frame-dir> -out <sheet> -report <report>`
- `go run ./cmd/validatemotion -variant lovebird_peach_faced -require-accepted` only after catalog metadata is added by the parent integration lane.
- Slot review: `00-03` idle, `04-11` hop-walk, `12-19` faster hop, `26-31`
  species action, `32-39` turn, `56-61` reaction/recover.

## Parallelization Boundaries

Safe to parallelize later:

- Read-only validation of accepted cockatiel and Java sparrow source contacts.
- A new `lovebird_peach_faced` source child lane after parent approves the exact
  variant id and run directory.
- Separate read-only Pages queue review comparing public order to product queue.
- Cat/dog prompt research for `scottish_fold_silver_tabby`,
  `ragdoll_seal_bicolor`, and `shiba_inu_red`, without touching runtime/page.

Must stay serialized:

- Budgerigar frame ownership and promotion until Lane A reports completion.
- Runtime `runtimeVariantIDs` changes.
- Public Pages current/upcoming order changes.
- Release docs, release links, tags, GitHub Releases, and publish/deploy steps.
- Promotion from run directories into `docs/art-source/<family>/motion-source/`.

## Exact Future Lane Prompt

Use this if the parent wants a safe independent new source lane after budgerigar
and the accepted cockatiel/Java sparrow promotion decision is queued separately:

```text
AnimalsDesktop bird source lane: generate lovebird_peach_faced set00 frames 00-04 only.

Use the repo AGENTS.md rules and animal-motion-imagegen-flow. Write only under:
docs/art-source/one-frame-method-fullrun-YYYYMMDD/lovebird-peach-faced-set00/

Do not edit budgerigar, cockatiel, java-sparrow, catalog, runtime sprites, docs page, release files, workflows, tags, GitHub Releases, or other animals.

Generate one complete right-facing 2D sprite peach-faced lovebird per frame on transparent or chroma green background. Preserve a compact small-parrot body, green wings/body, peach-orange face/forehead, pale beak, short rounded tail, and two small feet on one baseline. Reject text, scenery, props, floor, shadow, duplicate animals, cropped beak/tail/feet, budgerigar stripes, cockatiel crest, Java sparrow cheek pattern, oversized parrot proportions, or left-facing side-view unless the slot is a turn.

Stop after frames 00-04 and report the run directory, manifest rows, normalized frame paths, QA JSON paths, and contact sheet for parent visual gate.
```

