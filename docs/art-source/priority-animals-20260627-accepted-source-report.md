# Priority animals accepted source report

Date: 2026-06-27

This parent pass promotes thirteen completed one-frame ImageGen runs into
frame-addressable accepted `set00` source assets. These variants are accepted
as source art only; they are not added to `runtimeVariantIDs` or public
downloads in this pass.

## Verdict

Accepted as slot-based `set00` motion sources:

- `chipmunk_striped` under `docs/art-source/chipmunk/motion-source/`
- `gecko_leopard` under `docs/art-source/leopard-gecko/motion-source/`
- `whites_tree_frog_blue` under `docs/art-source/whites-tree-frog-blue/motion-source/`
- `budgerigar_green_yellow` under `docs/art-source/budgerigar/motion-source/`
- `cockatiel_normal_gray` under `docs/art-source/cockatiel/motion-source/`
- `java_sparrow_normal` under `docs/art-source/java-sparrow/motion-source/`
- `parrotlet_green` under `docs/art-source/parrotlet/motion-source/`
- `lovebird_peach_faced` under `docs/art-source/lovebird/motion-source/`
- `ragdoll_seal_bicolor` under `docs/art-source/ragdoll/motion-source/`
- `scottish_fold_silver_tabby` under `docs/art-source/scottish-fold/motion-source/`
- `french_bulldog_fawn` under `docs/art-source/french-bulldog/motion-source/`
- `maine_coon_brown_tabby` under `docs/art-source/maine-coon/motion-source/`
- `domestic_shorthair_calico` under `docs/art-source/domestic-shorthair/motion-source/`

`budgerigar_green_yellow` was promoted after parent-side fallback review for
frames `40`, `41`, and `60`. Those frames are documented under
`docs/art-source/budgerigar/motion-source/qa/` and accepted for source-art
continuity only.

## Mechanical QA

All accepted directories were revalidated after copying into this branch.

| Variant | Frames | Audit result | Warnings | Sheet |
| --- | ---: | --- | ---: | --- |
| `chipmunk_striped` | 62/62 | valid=62 missing=0 invalid=0 | 52 | `sheets/chipmunk-striped-source-set00.png` |
| `gecko_leopard` | 62/62 | valid=62 missing=0 invalid=0 | 11 | `sheets/leopard-gecko-source-set00.png` |
| `whites_tree_frog_blue` | 62/62 | valid=62 missing=0 invalid=0 | 72 | `sheets/whites-tree-frog-blue-source-set00.png` |
| `budgerigar_green_yellow` | 62/62 | valid=62 missing=0 invalid=0 | 50 | `sheets/budgerigar-green-yellow-source-set00.png` |
| `cockatiel_normal_gray` | 62/62 | valid=62 missing=0 invalid=0 | 47 | `sheets/cockatiel-normal-gray-source-set00.png` |
| `java_sparrow_normal` | 62/62 | valid=62 missing=0 invalid=0 | 90 | `sheets/java-sparrow-normal-source-set00.png` |
| `parrotlet_green` | 62/62 | valid=62 missing=0 invalid=0 | 42 | `sheets/parrotlet-green-source-set00.png` |
| `lovebird_peach_faced` | 62/62 | valid=62 missing=0 invalid=0 | 92 | `sheets/lovebird-peach-faced-source-set00.png` |
| `ragdoll_seal_bicolor` | 62/62 | valid=62 missing=0 invalid=0 | 20 | `sheets/ragdoll-seal-bicolor-source-set00.png` |
| `scottish_fold_silver_tabby` | 62/62 | valid=62 missing=0 invalid=0 | 19 | `sheets/scottish-fold-silver-tabby-source-set00.png` |
| `french_bulldog_fawn` | 62/62 | valid=62 missing=0 invalid=0 | 6 | `sheets/french-bulldog-fawn-source-set00.png` |
| `maine_coon_brown_tabby` | 62/62 | valid=62 missing=0 invalid=0 | 21 | `sheets/maine-coon-brown-tabby-source-set00.png` |
| `domestic_shorthair_calico` | 62/62 | valid=62 missing=0 invalid=0 | 70 | `sheets/domestic-shorthair-calico-source-set00.png` |

Commands run:

```sh
go run ./cmd/auditframes -frames-dir <family>/motion-source/accepted-frames/set00 -strict -artifact-warnings -report <family>/motion-source/accepted-frames/set00-auditframes-report.json
go run ./cmd/assemblemotion -frames-dir <family>/motion-source/accepted-frames/set00 -out <family>/motion-source/sheets/<sheet>.png -report <family>/motion-source/accepted-frames/set00-assemblemotion-report.json
```

## Visual Review

- `chipmunk_striped`: pass. Striped coat, compact body, and tail stay readable.
  Turn frames are the main angle change; no persistent crop, prop, or duplicate
  animal was found.
- `gecko_leopard`: pass. Low crawler anatomy keeps four underside legs, spotted
  yellow-tan body, and a thick tail. No feet appear on the back/top of the body.
- `whites_tree_frog_blue`: pass with chroma-adjacent warning note. The high
  warning count is accepted after contact-sheet review because the blue body,
  toe pads, and legs remain intact with no visible green matte damage.
- `budgerigar_green_yellow`: pass with fallback note. The yellow head, green
  body, blue cheek, black wing markings, pink feet/contact, scale, and baseline
  remain readable after parent fallbacks for `40`, `41`, and `60`.
- `cockatiel_normal_gray`: pass. Crest, orange cheek patch, gray body, and long
  tail remain identifiable. Low fast frames stay bird-like rather than becoming
  a rodent silhouette.
- `java_sparrow_normal`: pass with angle note. Black head, white cheek, red
  beak, and gray body remain readable across the set. Higher artifact warnings
  are tied to small feet and fine alpha details, not visible body breakage.
- `parrotlet_green`: pass. Compact green parrotlet identity, short tail,
  visible feet/contact, scale, and baseline remain stable across the full set.
  The low forage band and upright/rest band read as continuous pose groups
  without a one-frame body-size spike.
- `lovebird_peach_faced`: pass. Peach/orange face, green body, and visible
  feet/contact remain readable across the full set. The final reaction/recover
  band keeps a stable baseline and scale, and earlier one-frame size jumps were
  resolved before source promotion.
- `ragdoll_seal_bicolor`: pass. Seal bicolor mask, pale cream body, white paws,
  and fluffy tail remain readable. The final upright reaction band gets taller
  as a group rather than as a single-frame size spike, and the full grid keeps
  baseline and paw contact stable.
- `scottish_fold_silver_tabby`: pass after retrying the final `55-61` block.
  The first completed tail block was rejected for pale coat drift and a
  different swirl pattern. The accepted replacement restores darker
  silver-tabby contrast, folded ears, ringed tail, attached paws, and a stable
  baseline across the final reaction/recover band.
- `french_bulldog_fawn`: pass. The fawn coat, dark mask, upright ears, compact
  bulldog body, attached paws, and baseline remain readable. Frames `53-57`
  form a consistent face-groom/front-facing band, and `58-61` return toward the
  right-facing reaction/recover band without an isolated body-size spike.
- `maine_coon_brown_tabby`: pass. Brown tabby stripes, longhair body, fluffy
  ringed tail, ears, attached paws, and baseline remain readable. Frames `52-54`
  form the face-groom band, and `55-61` recover toward right-facing
  alert/reaction poses without an isolated body-size spike.
- `domestic_shorthair_calico`: pass. Orange/black/white calico patches, visible
  paws/contact, tail, ears, and baseline remain readable across low, walk,
  front-turn, and reaction/recover pose groups. Frames `53-61` form a coherent
  reaction/recover band; no isolated body-size spike was found.

## Release Status

These sources are ready for the next runtime/version lane. They should remain
off the public current-animal list until the corresponding app version,
Pages copy, release notes, and downloadable ZIP artifacts are aligned.
