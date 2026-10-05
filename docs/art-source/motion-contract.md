# AnimalsDesktop Motion Contract

Last updated: 2026-07-26

This file is the acceptance standard for promoting one animal family from seed/prototype art to accepted AnimalsDesktop motion art. The goal is to keep DeguDesktop compatibility while letting each species move like itself.

## Runtime Contract

Every selectable variant must keep the DeguDesktop sprite-sheet shape:

- 10 runtime motion sets per variant.
- 62 frames per set.
- 96x64 transparent canvas per frame.
- One complete animal per frame.
- The final sheet size for one set is 5952x64.
- Camera, scale, body volume, baseline, contact points, and facing direction stay stable inside a set.
- Ears, feet, whiskers, tail, shell, wings, or toes must not be cropped.
- No text, borders, scenery, shadows, costumes, props, multiple animals, or human-like poses.
- Source motion art is generated one pose per PNG. Multi-pose sheets, grids, contact sheets, and baked checker-background images are review references only and cannot be promoted as accepted source frames.

The ten runtime slots are a compatibility/storage contract, not a requirement
for ten independently generated source families. The normal accepted source is
one canonical 62-frame `set00`; `cmd/importanimals` expands it into all ten
runtime sheets. An optional independent `set00` through `set09` family is
allowed only when all ten sources exist and are byte-unique.

The app may still use mirrored drawing for left/right movement. Source frames should be clean right-facing side-view frames unless a profile explicitly documents another direction.

## Frame Slots

The frame numbers stay DeguDesktop-compatible. The meaning can be interpreted per motion profile:

| Frames | Count | DeguDesktop slot | Species interpretation |
| --- | ---: | --- | --- |
| 00-03 | 4 | idle | breathing, ear/head/tail micro motion, perch settle |
| 04-11 | 8 | walk | walk, trot, stalk, crawl, hop-walk, plod, slither travel |
| 12-19 | 8 | scurry | faster travel: scurry, bound, trot burst, crawl burst, slither wave |
| 20-25 | 6 | nibble | forage, sniff, nibble, peck, tongue flick, feed |
| 26-31 | 6 | hop | species action: hop, paw lift, play bow, pounce prep, low crawl shift, flutter |
| 32-39 | 8 | turn | body turn or turn-prep that stays readable when direction changes |
| 40-43 | 4 | eat | eat, chew, peck, lick, tongue feed, head-down feed |
| 44-47 | 4 | dig | dig, paw scratch, ground check, substrate push, shell/leg adjustment |
| 48-51 | 4 | stand | stand, sit, crouch, stretch, perch lift, alert posture |
| 52-55 | 4 | groom face | face groom, preen, whisker clean, paw clean, body adjustment |
| 56-61 | 6 | wheel run | wheel run only for wheel-capable profiles; otherwise alert or species-safe reaction |

Wheel-capable profiles are currently `degu` and `small-rodent-scurry`. Other profiles must never use wheel-only body poses in accepted art.

## Per-Animal Promotion Rule

Promote one animal family at a time. A family is complete for a version bump only when all target variants in that family have:

- source-truth PNGs reviewed at 96x64 on light and dark backgrounds;
- 62 standalone transparent motion-frame PNGs for every accepted variant;
- one canonical accepted 62-frame source set, expanded by the importer into all 10 runtime sheets;
- non-empty alpha-bounded frames with stable baseline and no cropped anatomy;
- a contact sheet or preview proving the full motion set is visually reviewable;
- catalog entries with species, breed/morph, color, popularity tier, source status, motion profile, and sprite base;
- runtime behavior review confirming wheel, hop, crawl, slither, perch, and low-body constraints are respected;
- `go run ./cmd/importsheet`, `go run ./cmd/importanimals`, `go run ./cmd/validatemotion -runtime-only -require-accepted`, `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, Windows build, and `git diff --check` passing.

If only source images exist, the family remains `prototype_only` or seed-stage even if it is selectable.

## Release Gate

Do not tag a new release just for planning, queue, page, or seed-art changes. A
content release should ship at least one animal family that has reached the
current accepted-source and runtime gate:

- the family is selectable in settings and tray UI;
- every accepted variant has one reviewed, accepted 62-frame canonical source set;
- the importer has generated all 10 runtime sheets from that canonical source;
- frame slots include real pose changes, not only recolor, bob, or duplicated walk frames;
- the family has species-appropriate motion for idle, walk, scurry/fast travel, forage/feed, action, turn, rest/stand, groom/preen/adjust, and alert/wheel-safe reaction;
- QA passes locally and in GitHub Actions;
- `cmd/validatemotion` reports `release_ready: true` for every runtime variant;
- the GitHub Pages page and `kdevelopk.pages.dev` works page describe the newly completed animal accurately.

Older `v0.1.x` preview exceptions that described accepted `set00` animals as
waiting for nine more independent source sets are historical and superseded by
this contract. Preview labels may still describe release maturity, packaging,
or public-testing scope, but they must not imply that canonical one-set animals
are missing source assets.

Public coming-soon silhouettes are release-managed: when an animal graduates
into a release, remove it from the planned silhouette/list and move it into the
current-animal section.

The parent QA pass may use color correction only for isolated color continuity
defects when the frame's anatomy, pose, alpha, scale, and species read are
otherwise acceptable. Preserve the original frame, write a correction report,
regenerate the contact sheet, rerun `cmd/auditframes`, and record the decision in
the visual QA report. Regenerate instead of correcting when the defect is wrong
species, wrong anatomy, cropped body parts, props, text, duplicate animals,
opaque/scenery background, or an unusable pose.

For future queues, normal striped chipmunk / シマリス is requested as a distinct
source family from `albino_chipmunk`. Coming-soon silhouettes remain optional
until page-specific generated future art exists.

## One-Animal Version Routine

Use this cycle repeatedly:

1. Pick the highest-priority family from `docs/art-source/asset-queue.md`.
2. Generate 4-6 source candidates for each target variant in that family.
3. Accept only the best 1-3 variants for the first family pass.
4. Generate 62 motion frames for those accepted variants using the slot table above.
5. Prepare single ImageGen candidates with `go run ./cmd/prepareframe` only when raw true-alpha output is unavailable; use `-background chroma-green` only for deliberate green-screen candidates and reject checker/noisy backgrounds.
6. Audit the standalone frame folders with `go run ./cmd/auditframes`.
7. Assemble each accepted 62-frame standalone PNG folder with `go run ./cmd/assemblemotion`.
8. Import, validate, visually review, and update source status.
9. Reflect the animal in settings, tray labels, docs, and public pages.
10. Run local QA and confirm GitHub Actions.
11. Commit, then tag a release only when the animal family meets the release gate above.

Do not publish or name a future animal queue until the next family is deliberately chosen by the parent thread. Color, breed, and morph variants should follow only after each accepted base motion is readable in runtime.

## Profile-Specific Notes

- Chinchilla: heavy soft body, rounded ears, fluffy tail, gentle scurry, no degu recolor.
- Macaroni mouse / fat-tailed gerbil: compact small body, thick tail, tiny feet, low fast scurry.
- Hamster: round cheek-forward body, very short tail, cheek/groom actions, wheel capable.
- Gecko: low body, splayed feet, toe pads if readable, crawl and tongue/alert actions, no wheel.
- Momonga / sugar glider: membrane silhouette, low skitter, brief stretch/glide-like action inside 96x64, no long airborne scene.
- Small birds: perch/ground hop, peck, preen, head turn, short flutter; no cage, branch, scenery, or floating flight loop.
