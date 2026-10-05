# Iteration Log

## 2026-06-21

- Restored `kdevelopk.pages.dev` from the recent icon-gallery Cloudflare deployment and added AnimalsDesktop with a chinchilla-based icon.
- Made `UDteach/works-gallery` the public GitHub source for the icon gallery.
- Replaced the AnimalsDesktop GitHub Pages draft with a Coming Soon page. Pet lists will appear only after chinchilla is complete.
- Removed GitHub Pages download ZIP generation for the Coming Soon phase so unfinished builds are not published through the site.
- Narrowed the GitHub Pages artifact to the public Coming Soon `index.html` only.
- Scoped the Windows runtime selectable variants to `chinchilla_standard_gray` only, so unfinished backlog variants and degu variants are not selectable in AnimalsDesktop builds.
- Generated the first ImageGen chinchilla pose-source batch, rejected the non-transparent raw sheet as final art, and extracted 16 transparent 96x64 pose frames for the motion-source pass.
- Built a non-release `set00` 62-frame draft sheet from the 16 extracted poses to make the next motion-review pass concrete.
- Recorded the chinchilla-first release loop and the no-degu rule in `.codex/tasks/chinchilla-v010-release.md`.
- Simplified the public AnimalsDesktop GitHub Pages page to plain Coming Soon only and strengthened workflow checks against unfinished animal names, images, download links, and work-in-progress lists.
- Connected `chinchilla_standard_gray` to the 62-frame draft motion source sheet in `cmd/importanimals`; runtime sprite sheets now come from ImageGen-derived motion frames and remain explicitly marked non-release draft until accepted `set00` through `set09` variation exists.
- Added draft `set01` through `set09` chinchilla motion-source sheets and taught `cmd/importanimals` to import a complete `set00` through `set09` source family when present.
- Added `cmd/validatemotion` and wired the release workflow to fail when runtime motion sources are still draft instead of accepted.
- Rejected the mechanically shifted chinchilla draft `set01` through `set09` sheets after runtime review; accepted motion work must come from real ImageGen pose/frame generation and visual QA.
- Rejected worker-generated chinchilla sheets and opaque/checker-background PNGs as review-only, and tightened importer/validator checks so full-opaque motion frames cannot pass as accepted transparent sources.
- Added `cmd/assemblemotion` to assemble exactly 62 standalone 96x64 transparent PNGs into one 5952x64 motion source sheet, blocking wrong dimensions, empty frames, and opaque/checker backgrounds.
- Tried a parent-thread one-pose ImageGen idle prompt; the result was `1536x1024` with `AlphaMin=255`, so it was rejected and kept out of the repo. Added single-frame prompt and accepted-frame staging docs for the next clean generation pass.
- Added `cmd/auditframes` so partial one-pose PNG progress can be measured without promoting bad ImageGen output: valid, missing, invalid, and edge-warning counts are reported per set.
- Added `cmd/prepareframe` for one-pose ImageGen candidates: true-alpha input is fitted to 96x64, uniform edge backgrounds can be removed, and checker/noisy backgrounds are rejected before any visual-review promotion.
- Tested a pure green single-pose ImageGen fallback. The first prepared output visibly retained green background, so `cmd/prepareframe` was tightened to fail when cleaned content still touches the source canvas edge; the candidate remains rejected.
- Added explicit `chroma-green` preparation, transparent-RGB cleanup, and green despill. The first visually reviewed chinchilla idle frame was promoted to `accepted-frames/set00/frame-00.png`; `cmd/auditframes` now reports `valid=1 missing=619`.

## 2026-06-25

- Ran a Codex built-in ImageGen-only review pass for `hamster_golden_syrian`, `macaroni_mouse_tan`, `sugar_glider_gray`, and `rabbit_chestnut_agouti` under `docs/art-source/external-ai-trials/codex-imagegen-20260625/`, keeping all outputs out of `accepted-frames` and runtime/catalog paths.
- Produced and locally split 16-cell trial sheets for all four variants; all split outputs had non-empty alpha and no edge-alpha cells after 96x64 normalization.
- Attempted 62-frame sheets only for the strongest 16-cell directions, `macaroni_mouse_tan` and `sugar_glider_gray`, then rejected both because visual contact sheets showed row-wrapped/cropped fragments despite passing mechanical empty/edge checks.
- Recorded prompts, saved paths, visual decisions, and parent-thread integration recommendations in `docs/art-source/external-ai-trials/codex-imagegen-20260625/report.md`.
- Tested a Codex built-in ImageGen one-frame retry-until-pass method for eight high-risk `sugar_glider_gray` poses under the same isolated external trial directory. All eight attempt-01 outputs produced complete 96x64 review candidates with no empty alpha and no edge-alpha; the method appears stronger than sheet extraction for high-risk pose repair, with style/scale matching as the remaining integration risk. See `docs/art-source/external-ai-trials/codex-imagegen-20260625/one-frame-method-report.md`.

## 2026-06-26

- Promoted the original `guinea_pig_tricolor` one-frame `set00` run to accepted source art under `docs/art-source/guinea-pig/motion-source/`, assembled `guinea-pig-tricolor-source-set00.png`, and changed the catalog entry from shape-only to accepted motion source. The template-lock alternate was rejected because late frames were visibly damaged.
- Verified the guinea pig source with `auditframes` (`valid=62 missing=0 invalid=0` with artifact warnings), `assemblemotion`, `validatemotion -variant guinea_pig_tricolor -require-accepted`, `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, and `git diff --check`. `release_ready` remains false because only one runtime set exists.
- Rechecked interrupted fancy-rat generation threads and reconciled local state to frames `00-19`; current `auditframes -strict -artifact-warnings` reports the expected incomplete state `valid=20 missing=42 invalid=0 warnings=10`.
- Confirmed the GitHub Pages / asset queue planned animals as guinea pig, fancy rat, albino chipmunk, Richardson's ground squirrel, and Yorkshire terrier. Started four scoped local Codex generation threads for the remaining planned work, each limited to its own `docs/art-source/one-frame-method-fullrun-20260626/` run directory and no catalog/runtime/page/Git edits.
- Created `.codex/tasks/20260626-release-asset-parity.md` to define "implemented asset parity" as the current preview-runtime level: accepted 62-frame `set00`, assembled source sheet, catalog accepted source status, runtime validation with one accepted set, and prepared page/queue updates without publishing. Added heartbeat monitor `animaldesktop-asset-parity-monitor` for child generation recovery and parent-side QA.
- Parent-side monitoring found no mechanical QA issues in current partial canonical frames, but detected duplicate `attempt=01` manifest rows in the fancy-rat child run. Sent all active child threads a correction to continue only from missing frames and use higher attempt numbers for retries.
- Performed a stricter visual audit for `guinea_pig_tricolor`, especially the tricolor coat pattern. Added 2x and 4x pattern-review sheets plus rough color-ratio metrics; no pattern-break outliers or visual rejects were found. Frames 53-55 have the strongest posture-driven marking angle change but still read as the same tricolor guinea pig.
- Added `guinea_pig_tricolor` to the local preview-runtime list, regenerated runtime sprites/import report, refreshed the current animal icon and preview/page assets, and removed guinea pig from the coming-soon silhouette image. This remains a local release-prep state; publishing still requires a next-version decision and matching release ZIPs.
- Added parent-monitor progress contact sheets for the four parallel generation lanes. Fancy rat, albino chipmunk, and Richardson's ground squirrel are mechanically clean so far; albino chipmunk still needs final tail/species-read review. The initial Yorkshire terrier run was rejected at partial QA as too shepherd/spitz-like, and the child thread was instructed to start a silkier toy-dog retry run.
- Started a separate Yorkshire terrier silky retry thread `019f02d9-1a2a-7250-9562-69ae3047ee5a` because the original Yorkshire child continued advancing before receiving the queued correction. Latest monitor snapshot: fancy rat 53/62, albino chipmunk 27/62, Richardson's ground squirrel 27/62, Yorkshire silky retry 0/62.
- Re-ran parent monitor contacts after compaction/resume. Latest snapshot recorded in `parent-monitor/progress-summary.json`: fancy rat 58/62, albino chipmunk 33/62, Richardson's ground squirrel 32/62, Yorkshire original 26/62 rejected-reference, Yorkshire silky 8/62 rejected-reference.
- Rejected the Yorkshire silky retry direction at partial visual QA because it still reads too much like an upright-eared large black-and-tan dog. Added `yorkshire-terrier-longcoat-retry-set00-00-61.md` and started thread `019f02e5-096a-7c21-bd45-1592f7258087` for a long-coated low rectangular toy Yorkie retry with an early 00-04 stop gate.
- Passed the Yorkshire longcoat retry through the parent early visual gate for frames 00-04. It reads as a low, long-coated toy Yorkie rather than shepherd/spitz, so the child thread was told to continue frames 05-61. Fancy rat is now waiting on only frame 61 before full review.
- Updated the public page layout so current animal icons use fixed image slots and the hero preview uses right-facing, edge-spaced animals. The sugar glider page asset was mirrored for page consistency only; runtime sprite sheets were not changed.
- Added Himalayan rabbit, Djungarian hamster, Campbell hamster, grayish chestnut Netherland Dwarf, and Holland Lop to the coming-soon queue. Started five scoped generation threads for them, each limited to its own run directory and required to generate real color ImageGen source frames rather than silhouette-only or reused existing images.
- Rebuilt `docs/assets/animalsdesktop-coming-soon-silhouettes.png` from a fresh Page-specific ImageGen source stored at `docs/art-source/one-frame-method-fullrun-20260626/page-coming-soon/coming-soon-eight-animals-imagegen-source.png`, then converted only that generated source into black silhouettes. Existing runtime/source/prototype animal images are not used for the coming-soon silhouette.
- Recorded the release policy change: completed animals should move through small preview version bumps after parent QA and matching artifacts, while full 10-set DeguDesktop-level completion remains a separate release gate.
- Deleted heartbeat automation `animaldesktop-asset-parity-monitor` during parent-thread takeover, then revalidated fancy rat and accepted albino chipmunk. Albino chipmunk passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=70`; visual QA accepts it as a faint-striped albino chipmunk rather than the earlier rejected white-squirrel seed.
- Promoted `albino_chipmunk` to accepted `set00` source art under `docs/art-source/albino-chipmunk/motion-source/`, assembled `albino-chipmunk-source-set00.png`, added it to the local runtime preview list, regenerated runtime sprites/page preview assets, and moved it from the coming-soon queue to current-page/runtime preview prep. External publication still requires explicit release/version approval and matching ZIP artifacts.
- Accepted Richardson's ground squirrel after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=57`; visual QA accepts it as a low-tailed ground squirrel without chipmunk stripes or tree-squirrel tail drift.
- Added a separate `ground_squirrel` catalog species, promoted `richardsons_ground_squirrel` to accepted `set00` source art under `docs/art-source/richardsons-ground-squirrel/motion-source/`, regenerated runtime sprites/page preview assets for 10 local preview animals, and moved Richardson's ground squirrel out of the coming-soon queue.
- Accepted the Yorkshire Terrier longcoat retry after it reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=92`; visual QA accepts it as a compact long-coated Yorkie rather than the earlier shepherd/spitz-like rejected runs.
- Promoted `yorkshire_terrier_longcoat` to accepted `set00` source art under `docs/art-source/yorkshire-terrier/motion-source/`, regenerated runtime sprites/page preview assets for 11 local preview animals, and moved Yorkshire Terrier out of the coming-soon queue.
- Accepted Djungarian hamster after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=85`; visual QA accepts it as a gray-white winter-white dwarf hamster with a visible dorsal stripe. Promoted `djungarian_hamster` to accepted `set00` source art under `docs/art-source/djungarian-hamster/motion-source/`, regenerated runtime sprites/page preview assets for 12 local preview animals, and moved Djungarian hamster out of the coming-soon queue.
- Accepted Holland Lop after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=62`; visual QA accepts it as a broken orange lop rabbit with consistently dropped ears. Promoted `holland_lop_broken_orange` to accepted `set00` source art under `docs/art-source/holland-lop/motion-source/`, regenerated runtime sprites/page preview assets for 13 local preview animals, moved Holland Lop out of the coming-soon queue, and kept the page-only sugar glider icon right-facing per user request without changing runtime sprites.
- Accepted Netherland Dwarf after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=81`; visual QA accepts it as a compact short-eared grayish chestnut dwarf rabbit. Promoted `netherland_dwarf_chestnut` to accepted `set00` source art under `docs/art-source/netherland-dwarf/motion-source/`, regenerated runtime sprites/page preview assets for 14 local preview animals, and moved Netherland Dwarf out of the coming-soon queue.
- Accepted Himalayan rabbit after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=28`; visual QA accepts it as a cream-bodied Himalayan rabbit with dark ears, nose, feet, and tail. Promoted `himalayan_rabbit` to accepted `set00` source art under `docs/art-source/himalayan-rabbit/motion-source/`, regenerated runtime sprites/page preview assets for 15 local preview animals, and moved Himalayan rabbit out of the coming-soon queue.
- Accepted Campbell hamster after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=71`; visual QA accepts it as a warm gray-brown Campbell dwarf hamster with a visible dorsal stripe. Frames `60` and `61` had a localized warm color drift, so the parent preserved originals, applied a color-only correction toward frames `56-59`, regenerated contact sheets, and reran audit before promotion. Promoted `campbell_hamster` to accepted `set00` source art under `docs/art-source/campbell-hamster/motion-source/`, regenerated runtime sprites/page preview assets for 16 local preview animals, prepared v0.1.5 page/docs/workflow checks, and moved the future queue to text-only normal striped chipmunk (`シマリス`) without reusing albino chipmunk art.
- Added `scripts/build_page_assets.py` to rebuild page icons and the hero preview from runtime sheets. It keeps `sugar_glider_gray` right-facing for page assets only, leaving runtime sprites unchanged.

## 2026-06-27

- Promoted five completed priority-animal one-frame runs to parent-owned
  accepted `set00` source assets without adding them to `runtimeVariantIDs`:
  `chipmunk_striped`, `gecko_leopard`, `whites_tree_frog_blue`,
  `cockatiel_normal_gray`, and `java_sparrow_normal`. Each has 62 accepted
  frames, an assembled source sheet, copied one-frame review JSON, contact
  sheet, and branch-local `auditframes` / `assemblemotion` reports.
- Added catalog metadata for those five accepted sources and introduced the
  `bird-hop` motion profile for cockatiel / Java sparrow source variants. The
  public v0.2.2 runtime list and Pages current-animal list remain at 16 animals;
  the next version bump/release lane should decide when to move accepted
  sources into runtime. Budgerigar remains in-progress and was not promoted.
- Reordered the GitHub Pages Coming Soon grid for the next animal wave while
  keeping the public downloads on `v0.2.2`. The display priority is now striped
  chipmunk, blue White's tree frog, leopard gecko, then the bird wave, followed
  by popular cat and dog candidates already covered by the page-specific
  silhouette source sheet.
- Split `scripts/build_page_assets.py` upcoming silhouette handling into the
  ImageGen source-sheet layout and the public display order, so future Pages
  ordering changes do not accidentally crop the wrong source cell. Strengthened
  `scripts/verify_page_release.py` to verify the Coming Soon priority order and
  all local `assets/...` references.
- Updated the public GitHub Pages Coming Soon section with page-only silhouette
  cards. The first cards are quick release candidates for
  `gecko_leopard`, `whites_tree_frog_blue`, and `chipmunk_striped`; the next
  generation wave prioritizes budgerigar, cockatiel, and Java sparrow, with
  lovebird and parrotlet as the 4-5 thread cap fillers.
- Kept the published v0.1.5 page scoped to sixteen current animals so the live
  download links do not claim unreleased animals. Leopard gecko, blue White's
  tree frog, and striped chipmunk remain local accepted `set00` candidates until
  a follow-up release/version step updates artifacts.
- Extended `scripts/build_page_assets.py` and both Pages workflows so upcoming
  silhouette PNGs are generated, verified, copied into the Pages artifact, and
  kept separate from runtime/accepted source frames.
- Corrected the GitHub Pages Coming Soon silhouette flow so it no longer uses
  deterministic hand-drawn silhouettes. `scripts/build_page_assets.py` now
  requires the page-specific ImageGen source sheet at
  `docs/art-source/one-frame-method-fullrun-20260627/page-coming-soon/coming-soon-fifteen-animals-imagegen-source.png`
  and extracts black silhouettes from that generated source only.
- Simplified the Coming Soon cards to name plus silhouette only. Added cat
  breed candidates from the 2026 iPet/Nyanpedia cat breed ranking top five:
  mixed cat, Scottish Fold, Munchkin, Ragdoll, and Minuet. The parent release
  policy is to bump preview versions in small increments as each animal
  graduates into the current runtime/page list.
- Updated the deployable public Pages links to the published `v0.2.1` Windows
  and Mac release assets. Leopard gecko, blue White's tree frog, and striped
  chipmunk are shown as Coming Soon silhouettes until their release artifacts
  and version links are aligned.
- Added JP/EN language switching for the static GitHub Pages site without a
  frontend framework. The page uses JP/EN buttons, saves the selection in
  `localStorage`, updates page metadata, navigation, download copy, animal
  names, upcoming names, feature copy, version notes,
  and keeps the release-verifier-sensitive Windows badge source markup intact.
- Added macOS JP/EN language persistence and menu/settings switching. The
  existing JSON settings now stores `language`, the Cocoa status menu can switch
  between Japanese and English, and the settings window is rebuilt after a
  language change so labels, tabs, option titles, animal names, and placeholders
  refresh consistently. Windows code and runtime assets were not changed.
- Expanded the macOS animal picker from the old five hardcoded labels to the
  runtime catalog, added fixed/selected/random animal menus in the status item,
  added per-pet size controls, and corrected the macOS bundle identifier to
  `com.udteach.animalsdesktop`.
- Added first-pass macOS multi-monitor support. The Mac overlay now persists a
  display ID, starts on the saved screen when it is available, falls back to the
  main screen when it is not, and exposes the display selector in both the status
  menu and the settings Motion tab. Local installed app
  `/Users/kyota/Applications/AnimalsDesktop.app` was checked and is still the
  published `v0.2.1` arm64 build until explicitly replaced.
- Strengthened macOS release QA for `v0.2.2`: Darwin tests now require the
  exact 16 release-scoped runtime animals, exercise every animal through fixed
  and per-pet selection, and verify every visible size step from 70% through
  120%. `scripts/verify_page_release.py` now also checks that the Pages current
  animal grid matches `catalog.RuntimeVariants()` exactly.
- Prepared the `v0.2.2` release docs and Pages copy for the Mac parity release.
  This release keeps leopard gecko, blue White's tree frog, and normal striped
  chipmunk in Coming Soon until their runtime assets are promoted in a later
  animal-addition lane. Verified `go run ./cmd/importanimals`,
  `python3 scripts/verify_page_release.py`, `go run ./cmd/validatemotion
  -runtime-only -require-accepted`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, `git diff --check`, and macOS arm64/amd64
  `VERSION=v0.2.2` ZIP builds. The motion validator remains `release_ready=false`
  for all runtime animals because this is a one-set preview, not the full
  10-set release gate.
- Reprioritized the GitHub Pages Coming Soon queue from the parent takeover
  thread. The first wave now favors シマリス, リューシスティックモモンガ,
  アフリカヤマネ, ネザーランドドワーフ（ヒマラヤン）, アメリカモモンガ,
  hamster / Djungarian color variants, fancy rat coat variants, and グレーうさぎ
  before the earlier frog / gecko / bird candidates. Reptile morph expansion is
  tracked as lower priority.
- Rebuilt Coming Soon silhouettes from a new page-specific 18-animal ImageGen
  source sheet at
  `docs/art-source/one-frame-method-fullrun-20260627/page-coming-soon/coming-soon-eighteen-animals-imagegen-source.png`.
  `scripts/build_page_assets.py`, `scripts/verify_page_release.py`, and both
  Pages workflows now verify the 18-card order.
- Prepared a bounded post-budgerigar asset planning lane in
  `docs/development/post-budgerigar-asset-queue-20260627.md`. The prep keeps
  Lane A's budgerigar run untouched, preserves the 16-animal runtime/public
  release boundary, queues cockatiel and Java sparrow accepted-source promotion
  before a new lovebird source lane, and records parallelization boundaries for
  later bird, cat, and dog work.
- Redirected the current phase from a Pages-specific goal to reusable asset
  output. Added `scripts/export_upcoming_asset_pack.py` and exported the current
  18 upcoming animals to `assets/source/upcoming/20260627/` as transparent color
  cutouts, normalized black silhouettes, copied source sheet, contact sheets,
  and `manifest.json`. The exporter uses connected-component extraction instead
  of equal grid cells so animals remain complete even when the ImageGen sheet
  is not a strict grid.
- Checked current ranking sources and recorded the implementation priority
  evidence in `docs/development/popularity-priority-sources-20260627.md`. Bird
  production keeps the user's explicit order of budgerigar, cockatiel, and Java
  sparrow before later lovebird/parrotlet work; cat and dog queues are recorded
  as breed-specific candidates for later lanes.
- Promoted `budgerigar_green_yellow` as a 62-frame accepted `set00` source
  asset. Parent-side visual QA replaced gray-face drift frames `40` and `41`
  with documented max-alpha blends from accepted neighbors and kept the earlier
  `60` neighbor fallback; originals were preserved under the fullrun
  `rejected/` directory. Final asset-only QA reports
  `auditframes -strict -artifact-warnings` as
  `valid=62 missing=0 invalid=0 warnings=50`, assembled
  `docs/art-source/budgerigar/motion-source/sheets/budgerigar-green-yellow-source-set00.png`
  sized `5952x64`, added catalog metadata as an accepted motion source, and
  regenerated the seed import outputs. This does not add the variant to
  runtime, Pages, release notes, tags, or downloadable artifacts.
- Reconfirmed the current production target as reusable source assets, not
  GitHub Pages deployment. Parent-gated `lovebird_peach_faced` frames `00-12`:
  `auditframes` reports `valid=13 missing=49 invalid=0 warnings=17`, every
  generated frame is `96x64`, and alpha bbox checks show a stable `47px` height
  and baseline `56` across the generated range. Checker/dark contacts showed no
  obvious background contamination, missing feet, clipping, or sudden scale
  jump, so the lovebird child lane was advanced to generate only `13-20` while
  preserving `00-12`. The parrotlet lane remains in progress after its `00-04`
  parent gate; no `05-12` canonical frames are visible yet.
- Started a third bounded asset-only lane for the top cat-breed queue:
  `scottish_fold_silver_tabby` frames `00-04` only under
  `docs/art-source/one-frame-method-fullrun-20260627/scottish-fold-silver-tabby-set00-oneframe-62/`.
  The lane is constrained to run-local raw/normalized frames, QA notes, and
  light/dark/checker contacts; parent review will gate folded ears, silver tabby
  markings, four attached paws, baseline stability, and non-bouncy `cat-stalk`
  motion before any `05-12` continuation.
- Parent-gated `parrotlet_green` frames `00-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=11`, every generated frame is
  `96x64`, and alpha bbox checks show stable `52px` height and baseline `58`
  across the generated range. Checker/dark contacts showed compact green
  parrotlet identity, visible attached feet, no background contamination, no
  clipping, and no sudden scale jump. The parrotlet child lane was advanced to
  generate only `13-20` while preserving `00-12`. The Scottish Fold lane has
  begun raw generation but has not reached the `00-04` parent gate yet.
- Early-rejected the first Scottish Fold canonical `frame-00` for scale, not
  mechanics: it was `96x64` and transparent, but the visible bbox was only
  `48x27` with baseline `57`, much smaller than the existing cat 96x64 preview
  evidence around `83x58`. A recovery instruction was sent to rebuild the
  anchor larger before generating or claiming `00-04`, targeting a cat-sized
  bbox around `72-86px` wide, `45-56px` high, and baseline `58-60`.
- Parent-gated `lovebird_peach_faced` frames `00-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=24`. Frames `13-19` keep the accepted
  lovebird scale, colors, baseline, and feet/contact. `frame-20` is much lower
  (`54x29` bbox) but was accepted as the start of sniff/nibble/forage rather
  than a scale failure; the lovebird child lane was advanced to generate only
  `21-28`, with `21-25` required to stay visually continuous with the low
  forage posture.
- Parent-gated corrected `scottish_fold_silver_tabby` frames `00-04`:
  `auditframes` reports `valid=5 missing=57 invalid=0 warnings=0`; corrected
  idle bboxes are `83-84px` wide and `45-47px` high with baseline `59`.
  `frame-04` is lower as cat-stalk start and was accepted. The child lane was
  advanced to generate only `05-12` while preserving `00-04`.
- Started a fourth bounded asset-only lane for `ragdoll_seal_bicolor` frames
  `00-04` under
  `docs/art-source/one-frame-method-fullrun-20260627/ragdoll-seal-bicolor-set00-oneframe-62/`.
  This lane uses the same cat-size gate as Scottish Fold and remains run-dir
  only: no catalog, runtime, Pages, release, tag, or Git operations.
- Parent-gated `parrotlet_green` frames `00-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=15`. Frames `13-17` and `19` preserve
  compact green parrotlet scale, baseline, and feet/contact; `frame-18` is a
  little longer and `frame-20` is wider/lower, but `20` was accepted as the
  start of sniff/nibble/forage rather than a scale failure. The parrotlet child
  lane was advanced to generate only `21-28`, with `21-25` required to stay
  visually continuous with the low forage posture.
- Reconfirmed the latest user correction as asset-only output, not GitHub
  Pages. Parent-gated `ragdoll_seal_bicolor` frames `00-04`:
  `auditframes` reports `valid=5 missing=57 invalid=0 warnings=0`, bbox checks
  show `00-03` around `88x48-50` with baseline `57`, and `04` is lower at
  `88x36` with the same baseline. Temporary parent checker/dark contacts show
  clear Ragdoll seal bicolor identity, attached paws, and no crop; the lane was
  advanced to generate only `05-12`. Parent also gated `lovebird_peach_faced`
  `21-28` mechanically (`valid=29 missing=33 invalid=0 warnings=40`) but
  rejected only `frame-26.png` for an abrupt wing-flare/tall-silhouette jump
  from the low forage frames. The lovebird lane was instructed to preserve
  `00-25` and `27-28`, replace only `26`, rebuild contacts, rerun
  `auditframes`, and stop before `29`. `parrotlet_green` still shows no visible
  `21-28` output after parent polls, so a light continuation reminder was sent;
  `scottish_fold_silver_tabby` has raw `05-10` visible but no canonical
  `05-12` gate yet.
- Parent-gated `parrotlet_green` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=25`. Contact review showed `21-24`
  as low forage, `25-26` as a natural rise, and `27-28` as upright recovery
  while preserving compact green parrotlet identity, baseline, and feet/contact.
  The parrotlet lane was advanced to generate only `29-36`. Parent-gated
  `scottish_fold_silver_tabby` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=5`; bbox checks and contacts show a
  stable low cat-stalk sequence with baseline `59-60`, folded ears, silver
  tabby identity, attached paws, and no abrupt scale jump. The Scottish Fold
  lane was advanced to generate only `13-20`, with stripe shimmer and vertical
  bounce called out as risks to avoid.
- Parent-gated `lovebird_peach_faced` frames `21-28` after retrying only
  `frame-26`: attempt 01 was rejected for an abrupt wing-flare/tall-silhouette
  jump, attempt 02 was rejected for staying too low (`54x27`) and preserving a
  `26 -> 27` size jump, and attempt 03 passed as a usable intermediate posture
  (`54x36`, baseline `56`). Final `auditframes` for the current partial run
  reports `valid=29 missing=33 invalid=0 warnings=40`; contact review shows a
  smoother low-forage-to-upright transition, so the lovebird lane was advanced
  to generate only `29-36` while preserving `00-28`.
- Parent-gated `parrotlet_green` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=29`; bbox checks show stable height
  `54` and baseline `59`, and contacts preserve compact green parrotlet
  identity with visible feet/contact and no abrupt scale jump. The parrotlet
  lane was advanced to generate only `37-44`. Parent-gated
  `ragdoll_seal_bicolor` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=1`; bbox checks show stable baseline
  `57`, and contacts preserve Ragdoll seal bicolor identity, white paws, seal
  mask/tail, long tail, and no crop. The Ragdoll lane was advanced to generate
  only `13-20`.
- Parent-gated `scottish_fold_silver_tabby` frames `13-20`: `auditframes`
  reports `valid=21 missing=41 invalid=0 warnings=8`; bbox checks show stable
  width `86`, height `37-41`, and baseline `59-60`. Contact review, including
  a parent-created temporary `00-20` full-range contact because the child did
  not save one, showed stable Scottish Fold silver tabby identity, folded ears,
  tail, attached paws, and no abrupt scale jump. The Scottish Fold lane was
  advanced to generate only `21-28`, with an explicit reminder to save
  full-range contacts next time.
- Parent-gated `lovebird_peach_faced` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=49`; bbox checks show stable
  baseline `56`, and contacts preserve peach-faced lovebird color/body
  identity, visible feet/contact, and no abrupt wing flare or vertical scale
  jump. The lovebird lane was advanced to generate only `37-44`.
- Parent checked `parrotlet_green` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=34`, but contact review rejected only
  `frame-43.png` because it jumps upright between low `42` and low `44`,
  recreating the abrupt-size-change issue. The parrotlet lane was instructed to
  preserve `00-42` and `44`, replace only `43` with a low eat/ground-check
  intermediate, rebuild contacts, and stop before `45`. Parent-gated
  `ragdoll_seal_bicolor` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=5`; contacts preserve low
  cat-stalk/sniff posture, baseline `57`, seal mask/tail, white paws, long
  tail, and no crop. The Ragdoll lane was advanced to generate only `21-28`.
- Parent-gated the retried `parrotlet_green` `frame-43`: replacement bbox is
  `82x40` with baseline `59`, `37-44` contact no longer has the one-frame
  upright jump, and `auditframes` now reports
  `valid=45 missing=17 invalid=0 warnings=35`. The parrotlet lane was advanced
  to generate only `45-52`.
- Parent-gated `scottish_fold_silver_tabby` frames `21-28`: `auditframes`
  reports `valid=29 missing=33 invalid=0 warnings=11`; bbox checks show stable
  width `86` and baseline `59-60`. Contact review, including a parent-created
  temporary `00-28` full-range contact because the child again did not save one,
  showed stable folded ears, silver tabby identity, tail, attached paws, and no
  abrupt scale jump. The Scottish Fold lane was advanced to generate only
  `29-36`, with another explicit reminder to save full-range contacts.
- Parent checked `lovebird_peach_faced` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=58`, but contact review rejected only
  `frame-43.png` because it jumps upright between low `42` and low `44`,
  matching the same abrupt-size-change issue found in parrotlet. The lovebird
  lane was instructed to preserve `00-42` and `44`, replace only `43` with a low
  eat/ground-check intermediate, rebuild contacts, and stop before `45`.
- Reconfirmed the current output target as asset/source-frame production, not
  GitHub Pages. Parent-gated `ragdoll_seal_bicolor` frames `21-28`:
  `auditframes` reports `valid=29 missing=33 invalid=0 warnings=6`; bbox checks
  show stable baseline `56-57`, and contacts preserve Ragdoll seal bicolor
  identity, white paws, seal mask/tail, long tail, and no abrupt one-frame scale
  jump. The Ragdoll lane was advanced to generate only `29-36`.
- Parent-gated `parrotlet_green` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=41`; contact review shows `48-52` as
  a stable upright/rest band after the ground-check frames, with compact green
  parrotlet identity, visible feet/contact, and no isolated one-frame size
  spike. The parrotlet lane was advanced to generate final range `53-61`, then
  stop for parent final review.
- Parent-gated the retried `lovebird_peach_faced` `frame-43`: replacement bbox
  is `54x26` with baseline `56`, `37-44` contact no longer has the one-frame
  upright jump, and `auditframes` still reports
  `valid=45 missing=17 invalid=0 warnings=58`. The lovebird lane was advanced
  to generate only `45-52`.
- Parent-gated `scottish_fold_silver_tabby` frames `29-36`: `auditframes`
  reports `valid=37 missing=25 invalid=0 warnings=14`; bbox checks show stable
  baseline `59-60`. Contact review shows `34-36` as a continuous upright/rise
  band, not an isolated body-size spike, while preserving folded ears, silver
  tabby identity, tail, and attached paws. The lane was advanced to generate
  only `37-44`, with full `00-44` contact saving required.
- Parent final-gated `parrotlet_green` as a 62/62 draft source:
  `auditframes` reports `valid=62 missing=0 invalid=0 warnings=42`; parent
  temporary contacts for `53-61`, `45-61`, and `00-61` show stable compact
  green parrotlet identity, visible feet/contact, no pale/gray drift, and no
  one-frame size spike. The lane was told to save missing final contacts only,
  with no frame changes or promotion.
- `parrotlet_green` final contact artifact completion passed. The lane saved
  `53-61`, `45-61`, and `00-61` light/dark/checker contacts, and parent reran
  `auditframes` as `valid=62 missing=0 invalid=0 warnings=42`. This is ready
  for parent accepted-source promotion review as an asset-only source.
- Promoted `parrotlet_green` into the parent branch as an asset-only accepted
  source. The 62 frames, contact evidence, QA notes, assembled source sheet,
  catalog metadata, generated seed source, and 10 sprite sheets are present.
  Validation passed: `validatemotion -variant parrotlet_green
  -require-accepted`, `importanimals`, targeted Go tests, and
  `git diff --check`. Runtime variants, Pages, release, tags, and deploy were
  not changed.
- Parent-gated `lovebird_peach_faced` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=71`; contact review shows `47-52` as
  a stable upright/rest band after the low ground-check frames, preserving
  peach/orange face, green body, visible feet/contact, and no isolated
  one-frame size spike. The lovebird lane was advanced to generate final range
  `53-61`, then stop for parent final review.
- Parent-gated `ragdoll_seal_bicolor` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=6`; bbox checks show stable baseline
  `57`, and contacts preserve Ragdoll seal bicolor identity, white paws, seal
  mask/tail, long tail, attached paws, and no abrupt one-frame scale jump. The
  Ragdoll lane was advanced to generate only `37-44`.
- Parent checked `ragdoll_seal_bicolor` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=9`, but contact review rejected only
  `frame-43.png` because it jumps upright between low `42` and low `44`.
  The Ragdoll lane was instructed to preserve `00-42` and `44`, replace only
  `43` with a low ground-check intermediate, save the missing `37-44` and
  `00-44` contacts, and stop before `45`.
- Parent checked `scottish_fold_silver_tabby` frames `37-44`: `auditframes`
  reports `valid=45 missing=17 invalid=0 warnings=17`, but labeled contact
  review rejected only `frame-43.png` because it jumps upright between low `42`
  and low `44`. The Scottish Fold lane was instructed to preserve `00-42` and
  `44`, replace only `43` with a low head-down stalk/ground-check intermediate,
  save the missing full `00-44` contacts, and stop before `45`.
- Parent-gated the retried `ragdoll_seal_bicolor` `frame-43`: replacement bbox
  is `87x27` with top `31` and baseline `57`, saved `37-44` and `00-44`
  contacts no longer have the one-frame upright jump, and `auditframes` still
  reports `valid=45 missing=17 invalid=0 warnings=9`. The Ragdoll lane was
  advanced to generate only `45-52`.
- Reconfirmed the current deliverable as reusable local assets/source frames,
  not GitHub Pages. Parent work continues on
  `/Users/kyota/.codex/worktrees/42d3/AnimalDesktop`
  (`codex/upcoming-asset-pack`) to avoid mixing in unrelated dirty Pages/release
  changes from the older `/Volumes/.../AnimalDesktop` `main` worktree. Parent
  poll found `lovebird_peach_faced` still at canonical `00-52` and
  `ragdoll_seal_bicolor` still at canonical `00-44`, so reminder prompts were
  sent for their already-approved next ranges only. Parent-gated the retried
  `scottish_fold_silver_tabby` `frame-43`: current bbox is `86x39` with top
  `22` and baseline `60`, matching the low `42` / `44` neighbor band. The
  `37-44` checker contact no longer has the one-frame upright jump, and
  `auditframes` reports `valid=45 missing=17 invalid=0 warnings=17`
  (missing `45-61` expected). The Scottish Fold lane was advanced to generate
  only `45-52`.
- Promoted `lovebird_peach_faced` into the parent branch as an asset-only
  accepted source. Parent verified the final `53-61`, `45-61`, and full
  `00-61` contacts; `auditframes` reports
  `valid=62 missing=0 invalid=0 warnings=92`, and
  `assemblemotion` produced
  `docs/art-source/lovebird/motion-source/sheets/lovebird-peach-faced-source-set00.png`.
  The 62 accepted frames, contact evidence, QA notes, catalog metadata,
  generated seed source, and 10 sprite sheets are present. Validation passed:
  `validatemotion -variant lovebird_peach_faced -require-accepted`,
  `importanimals`, and targeted Go tests. Runtime variants, Pages, release,
  tags, downloads, and deploy were not changed.
- Ragdoll and Scottish Fold lanes remain active after lovebird promotion.
  Ragdoll has advanced to canonical `45-49` with low stable bboxes through
  `48` and a moderate `49` rise, but no `45-52` contacts yet. Scottish Fold
  has raw `frame-45-attempt-01.png` and no new canonical `45` yet. To keep the
  child pool near the requested 4-5 lane cap, two new first-gate asset lanes
  were queued: `maine_coon_brown_tabby` (`00-04` only, pending worktree
  `local:19a9757d-ffaa-4346-836c-1923ef1f7c3c`) and
  `french_bulldog_fawn` (`00-04` only, pending worktree
  `local:6a7f3365-2213-4d6f-9a11-f53659e2c3b4`). Both are source-frame lanes
  only, with no catalog/runtime/Pages/release/Git edits allowed in child
  threads.
- Parent-gated `french_bulldog_fawn` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=3`, bboxes are stable around
  `56-60x52` with baseline `57`, and the parent temporary checker contact
  shows a coherent fawn French Bulldog with dark muzzle, upright ears, sturdy
  body, attached paws, and no scale jump into `04`. The lane was instructed to
  save missing `00-04` contacts/audit evidence, preserve `00-04`, generate
  only `05-12`, save `05-12` and `00-12` contacts, rerun `auditframes`, and
  stop before `13`.
- Parent recovery prompt sent for `ragdoll_seal_bicolor` `50-52` after the
  child rejected `frame-50-attempt-01-height-spike` and
  `frame-51-attempt-01-upright-raw`. The prompt specified a gradual bbox
  ladder from accepted `49` into alert/rest: `50` about `35-38px` high, `51`
  about `40-44px`, and `52` about `45-49px`, all keeping baseline `57`.
  Scottish Fold now has raw `45-48` but no canonical `45` yet. Maine Coon has
  raw `00-04` and was asked to normalize or reject, with no weak canonical
  frames accepted by the parent.
- Parent-gated `maine_coon_brown_tabby` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=0`. Canonical bboxes keep baseline
  `57`; `04` is low at `88x33` with `4px` side margins and reads as the
  cat-stalk/walk start rather than a crop. Parent checker review shows a
  coherent Maine Coon brown tabby with longhair body, dark tabby stripes,
  fluffy tail/chest, prominent ears, attached paws, and no scale jump. The lane
  was instructed to save proper canonical `00-04` contacts/audit evidence,
  preserve `00-04`, generate only `05-12`, save `05-12` and `00-12` contacts,
  rerun `auditframes`, and stop before `13`.
- Parent-gated `scottish_fold_silver_tabby` frames `45-52`: `auditframes`
  reports `valid=53 missing=9 invalid=0 warnings=19`. Contact review shows
  `49-52` as a coherent upright/rest band after low `45-48`, not an isolated
  one-frame spike, while preserving folded ears, silver tabby markings, tail,
  attached paws, and baseline. The lane was instructed to save missing full
  `00-52` contacts, preserve `00-52`, generate only final `53-61`, save
  `53-61` / `45-61` / `00-61` contacts, rerun `auditframes`, and stop for
  final review.
- Parent-gated `french_bulldog_fawn` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=3`. Bboxes keep height `52` and
  baseline `57`; parent temporary `05-12` and `00-12` checker contacts show a
  stable fawn French Bulldog with dark muzzle, upright ears, sturdy body,
  attached paws, and no body-size spike. The lane was instructed to save
  missing `05-12` / `00-12` contacts, preserve `00-12`, generate only `13-20`,
  save `13-20` / `00-20` contacts, rerun `auditframes`, and stop before `21`.
- Ragdoll recovery is partially improved: `frame-50` is now accepted at
  `86x38 top=20 baseline=57`, but `frame-51` attempts 02 and 03 were rejected
  for landing too tall and recreating the abrupt-rise problem. Canonical frames
  currently stop at `50`; continue only with a low-mid `51` that meets the
  ladder target.
- Ragdoll `frame-51` repeated failure reached the split-lane threshold for
  this asset queue. Attempt 04 was also rejected as too low/short, while
  attempt 02 was too tall. Parent paused the original Ragdoll thread at
  canonical `00-50` and opened focused rescue thread
  `019f09e8-f775-7ea2-b9f5-70ba261d4735` for `frame-51` only, using existing
  worktree `/Users/kyota/.codex/worktrees/359e/AnimalDesktop`. The rescue
  target anchors to accepted `frame-37.png` (`88x40 top=18 baseline=57`) and
  must produce width `86-88`, height `40-42`, top `16-18`, baseline `57`,
  then stop for parent gate.
- Maine Coon `05-12` remains on hold: parent still saw old `frame-11.png` at
  `82x52 top=6 baseline=57`, a one-frame upright spike between `10` and `12`.
  The lane was re-instructed to retry only `frame-11` with target height
  `42-45`, top `13-16`, baseline `57`, then rebuild `05-12` and `00-12`
  contacts before parent gate.
- Promoted `ragdoll_seal_bicolor` as a reusable accepted `set00` source asset,
  not as a Pages/runtime/release change. Parent-side final QA accepted the
  `53-61`, `45-61`, and full `00-61` contacts; the 8-column full grid showed
  stable seal-bicolor mask, pale cream body, white paws, fluffy tail, baseline,
  and no isolated one-frame size spike. Mechanical QA passed with
  `auditframes` `valid=62 missing=0 invalid=0 warnings=20`, and
  `assemblemotion` produced
  `docs/art-source/ragdoll/motion-source/sheets/ragdoll-seal-bicolor-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `ragdoll_seal_bicolor` variant as `motion_source_accepted`; the runtime
  variant list, Pages, release notes, tags, downloads, and deploy outputs were
  not changed.
- Revalidated the parent branch after Ragdoll promotion:
  `go run ./cmd/validatemotion -variant ragdoll_seal_bicolor -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, `go run ./cmd/validatemotion -runtime-only
  -require-accepted`, and `git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets still have one set, not the
  full 10-set release gate.
- Cleaned AppleDouble `._*` files from the Git common object directory after
  they caused `git diff --check` to fail with `non-monotonic index`. A follow-up
  `find "$(git rev-parse --git-common-dir)" -name '._*'` returned empty.
- Remaining active asset lanes after the Ragdoll promotion: Scottish Fold is
  accepted through `00-54`, Maine Coon through `00-20`, and French Bulldog
  through `00-21`. Continue monitoring by local file state only and avoid
  large thread transcript reads.
- Scottish Fold reached mechanical `62/62` (`valid=62 missing=0 invalid=0
  warnings=19`) but was not promoted. Parent visual QA rejected `55-61` because
  `55` became too pale/low-contrast and `56-61` drifted into a lighter
  bullseye/swirl pattern that did not match the darker accepted `00-54`
  silver-tabby coat. The child lane was instructed to preserve `00-54` and
  retry only `55-61` with the prior dark stripe style.
- Queued a new bounded first-gate lane for `domestic_shorthair_calico` under
  pending worktree `local:e7c9af71-ba0a-4403-96ab-eca3a7f5fc23`. This covers
  the popular mixed-cat / domestic shorthair family using an existing catalog
  ID, with first parent gate `00-04` only and no catalog/runtime/Pages/release
  writes in the child lane.
- Promoted `scottish_fold_silver_tabby` as a reusable accepted `set00` source
  asset, not as a Pages/runtime/release change. Parent-side final QA accepted
  the retried `55-61`, `45-61`, and full `00-61` contacts; the final band
  restored darker silver-tabby contrast, folded ears, attached paws, ringed
  tail, and stable baseline. Mechanical QA passed with `auditframes`
  `valid=62 missing=0 invalid=0 warnings=19`, and `assemblemotion` produced
  `docs/art-source/scottish-fold/motion-source/sheets/scottish-fold-silver-tabby-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `scottish_fold_silver_tabby` variant as `motion_source_accepted`; runtime
  variants, Pages, release notes, tags, downloads, and deploy outputs were not
  changed.
- Revalidated the parent branch after Scottish Fold promotion:
  `go run ./cmd/validatemotion -variant scottish_fold_silver_tabby
  -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`,
  `go run ./cmd/validatemotion -runtime-only -require-accepted`, and
  `COPYFILE_DISABLE=1 git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets currently have one set, not
  the full 10-set release gate.
- Parent-gated `maine_coon_brown_tabby` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=9`. Parent reviewed the `29-36` and
  `21-36` checker contacts; `frame-34` is narrow (`48x52`) but reads as a
  front-turn frame between `33` and `35`, not an isolated scale jump. The lane
  was instructed to preserve `00-36`, generate only `37-44`, and stop before
  `45`.
- Parent-gated `french_bulldog_fawn` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=4`. Parent reviewed the `29-36` and
  `21-36` checker contacts; `34-35` form the front-turn band and `36` returns
  toward right-facing side view with stable fawn coat, dark mask, upright ears,
  attached paws, and baseline. The lane was instructed to preserve `00-36`,
  generate only `37-44`, and stop before `45`.
- Parent-gated `domestic_shorthair_calico` frames `05-12`: `auditframes`
  reports `valid=13 missing=49 invalid=0 warnings=5`. Parent reviewed the
  `05-12` and `00-12` checker contacts; `07` and `12` are taller walking
  poses but not body-scale spikes, and the calico patches remain coherent. The
  lane was instructed to preserve `00-12`, generate only `13-20`, and stop
  before `21`.
- Parent-gated `maine_coon_brown_tabby` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=13`. Parent reviewed the `37-44`
  and `29-44` checker contacts; `40-41` are the low ground-check band and
  `42-43` return to standing without an isolated scale jump. The lane was
  instructed to preserve `00-44`, generate only `45-52`, and stop before `53`.
- Parent-gated `domestic_shorthair_calico` frames `13-20`: `auditframes`
  reports `valid=21 missing=41 invalid=0 warnings=11`. Parent reviewed the
  `13-20`, `05-20`, and `00-20` checker contacts; `13-17` are a low fast
  movement band, `18-19` return upward, and `20` starts sniffing with the head
  lowered. The lane was instructed to preserve `00-20`, generate only `21-28`,
  and stop before `29`.
- `french_bulldog_fawn` has generated through `42`, but the `37-44` parent gate
  is still incomplete. Latest local audit is `valid=43 missing=19 invalid=0
  warnings=4`; a focused reminder was sent for only `43-44`, preserving
  `00-42` and stopping before `45`.
- Parent-gated `french_bulldog_fawn` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=5`. Parent reviewed the `37-44` and
  `29-44` checker contacts; `40-44` read as a low ground-check band with stable
  fawn coat, dark mask, upright ears, attached paws, and baseline. The lane was
  instructed to preserve `00-44`, generate only `45-52`, and stop before `53`.
- Current child-lane polling state after the French gate: `maine_coon_brown_tabby`
  remains at `valid=45` waiting for `45-52`, `french_bulldog_fawn` is also
  `valid=45` waiting for `45-52`, and `domestic_shorthair_calico` has generated
  through `25` (`valid=26`) while waiting for `26-28` before the next parent
  gate.
- Parent-gated `maine_coon_brown_tabby` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=17`. Parent reviewed the `45-52`,
  `37-52`, and `00-52` checker contacts; `45-46` are a low ground-check band,
  `47-51` return through alert/rest scale, and `52` is acceptable as the
  face-groom start. The lane was instructed to preserve `00-52`, generate only
  `53-61`, then create `52-61`, `45-61`, and `00-61` contacts plus final
  `auditframes`.
- `french_bulldog_fawn` remains at `valid=45 missing=17 invalid=0 warnings=5`.
  A focused reminder was sent to preserve `00-44`, generate only `45-52`, and
  stop before `53`.
- `domestic_shorthair_calico` is at `valid=28 missing=34 invalid=0 warnings=20`.
  Parent partial contacts for `21-27` and `13-27` showed no immediate scale or
  anatomy blocker, so the lane was instructed to preserve `00-27`, generate
  only missing `28`, and stop before `29` for the `21-28` parent gate.
- Parent-gated `french_bulldog_fawn` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=6`. Parent reviewed the `45-52`,
  `37-52`, and `00-52` checker contacts; `45-47` are the low-to-rising
  transition, `48-51` stay stable in the alert/rest band, and narrower `52` is
  acceptable as the face-groom start. The lane was instructed to preserve
  `00-52`, generate only `53-61`, then create `52-61`, `45-61`, and `00-61`
  contacts plus final `auditframes`.
- Parent-gated `domestic_shorthair_calico` frames `21-28`: `auditframes`
  reports `valid=29 missing=33 invalid=0 warnings=20`. Parent reviewed the
  `21-28`, `13-28`, and `00-28` checker contacts; `21-22` and `28` are low
  sniff/action poses, with coherent calico patches, visible paws/contact, intact
  tail, baseline continuity, and no isolated body-size spike. The lane was
  instructed to preserve `00-28`, generate only `29-36`, and stop before `37`.
- Promoted `french_bulldog_fawn` as a reusable accepted `set00` source asset,
  not as a runtime, Pages, release, tag, download, or deploy change. Parent
  final QA accepted the `52-61`, `45-61`, and full `00-61` contacts on checker,
  dark, and light backgrounds; the fawn coat, dark mask, upright ears, compact
  bulldog body, attached paws, and baseline remain readable. Mechanical QA
  passed with `auditframes` `valid=62 missing=0 invalid=0 warnings=6`, and
  `assemblemotion` produced
  `docs/art-source/french-bulldog/motion-source/sheets/french-bulldog-fawn-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `french_bulldog_fawn` variant as `motion_source_accepted`; runtime variants,
  Pages, release notes, tags, downloads, and deploy outputs were not changed.
- Revalidated the parent branch after French Bulldog promotion:
  `go run ./cmd/validatemotion -variant french_bulldog_fawn -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, and
  `go run ./cmd/validatemotion -runtime-only -require-accepted` all passed.
  `release_ready=false` remains expected because accepted source assets
  currently have one set, not the full 10-set release gate.
- Promoted `maine_coon_brown_tabby` as a reusable accepted `set00` source
  asset, not as a runtime, Pages, release, tag, download, or deploy change.
  Parent final QA accepted the `52-61`, `45-61`, and full `00-61` contacts on
  checker, dark, and light backgrounds; brown tabby stripes, longhair body,
  fluffy ringed tail, ears, attached paws, and baseline remain readable.
  Mechanical QA passed with `auditframes`
  `valid=62 missing=0 invalid=0 warnings=21`, and `assemblemotion` produced
  `docs/art-source/maine-coon/motion-source/sheets/maine-coon-brown-tabby-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `maine_coon_brown_tabby` variant as `motion_source_accepted`; runtime
  variants, Pages, release notes, tags, downloads, and deploy outputs were not
  changed.
- Revalidated the parent branch after Maine Coon promotion:
  `go run ./cmd/validatemotion -variant maine_coon_brown_tabby -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, and
  `go run ./cmd/validatemotion -runtime-only -require-accepted` all passed.
  `release_ready=false` remains expected because accepted source assets
  currently have one set, not the full 10-set release gate.
- Parent-gated `domestic_shorthair_calico` frames `29-36`: `auditframes`
  reports `valid=37 missing=25 invalid=0 warnings=37`. Parent reviewed the
  `29-36`, `21-36`, and `00-36` checker contacts; `32-36` form a coherent
  turn/front-facing band rather than an isolated body-size spike, while calico
  patches, visible paws/contact, tail, and baseline remain stable. The lane was
  instructed to preserve `00-36`, generate only `37-44`, and stop before `45`.
- Queued three asset-only child lanes, all capped at the first `00-04` parent
  gate and explicitly scoped away from Pages, runtime lists, releases, tags,
  downloads, and deploy output: `munchkin_brown_tabby` pending worktree
  `local:989cd7ad-752a-4242-9a78-cfcd7e2b8799`,
  `toy_poodle_apricot` pending worktree
  `local:b9029752-44bf-4a78-b43d-165529ca92e8`, and
  `british_shorthair_blue` pending worktree
  `local:fd8a46c7-2a07-450d-be98-cab5af70a4a6`.
- The queued lanes materialized as active child threads:
  `munchkin_brown_tabby` thread `019f0a7f-2439-77f3-974e-65025430c586` in
  `/Users/kyota/.codex/worktrees/1fe1/AnimalDesktop`,
  `toy_poodle_apricot` thread `019f0a7f-2449-7a92-8a59-c1c53ef5b278` in
  `/Users/kyota/.codex/worktrees/6572/AnimalDesktop`, and
  `british_shorthair_blue` thread `019f0a7f-2439-77f3-974e-64fd29322719` in
  `/Users/kyota/.codex/worktrees/f06b/AnimalDesktop`.
- Parent-gated `british_shorthair_blue` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=4`; bboxes keep height `52` and
  baseline `57`. The checker contact shows a coherent blue-gray British
  Shorthair with round head/body and visible paws/contact. The wider `02` and
  `04` frames are stride-width changes, not body-size spikes. The lane was
  instructed to preserve `00-04`, generate only `05-12`, and stop before `13`.
- Parent-gated `munchkin_brown_tabby` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=2`; bboxes are stable at `88x41` with
  baseline `57`. The checker contact shows a coherent short-legged brown tabby
  Munchkin with attached paws and ringed tail. Adjacent pixel diffs confirm the
  frames are not exact duplicates despite the stable silhouette. The lane was
  instructed to preserve `00-04`, generate only `05-12`, and stop before `13`.
- Parent-gated `domestic_shorthair_calico` frames `37-44`: `auditframes`
  reports `valid=45 missing=17 invalid=0 warnings=54`; all `37-44` frames keep
  height `52` and baseline `57`, with `39` as the wider stride frame. Parent
  reviewed the `37-44`, `29-44`, and `00-44` checker contacts; calico patches,
  visible paws/contact, tail, and baseline remain stable with no sudden
  body-size spike. The lane was instructed to preserve `00-44`, generate only
  `45-52`, and stop before `53`.
- Parent-gated `british_shorthair_blue` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=8`; `05-11` keep height `50-52` and
  baseline `57`, while `12` is an acceptable low sniff / ground-check pose at
  height `36` and the same baseline. Parent reviewed the `05-12` and `00-12`
  checker contacts and instructed the lane to preserve `00-12`, generate only
  `13-20`, and stop before `21`.
- Parent-gated `toy_poodle_apricot` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=0`; bboxes stay around `58-60x52`
  with baseline `57`. The checker contact shows a readable apricot Toy Poodle
  with curly coat, rounded muzzle/ears, attached paws, and no one-frame
  body-size spike. The lane was instructed to preserve `00-04`, generate only
  `05-12`, and stop before `13`.
- Parent-gated `munchkin_brown_tabby` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=9`; bboxes stay `88x37-41` with
  baseline `57`. Parent reviewed the `05-12` and `00-12` checker contacts;
  `06` and `10` read as low walking phases, with stable stripes, ringed tail,
  attached paws, and no body-size spike. The lane was instructed to preserve
  `00-12`, generate only `13-20`, and stop before `21`.
- Parent-gated `toy_poodle_apricot` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=0`; bboxes keep height `52` and
  baseline `57`, with `12` as a wider stride frame. Parent reviewed the
  `05-12` and `00-12` checker contacts; apricot coat, curly texture, rounded
  muzzle/ears, attached paws, and scale remain stable. The lane was instructed
  to preserve `00-12`, generate only `13-20`, and stop before `21`.
- Parent-gated `domestic_shorthair_calico` frames `45-52`: `auditframes`
  reports `valid=53 missing=9 invalid=0 warnings=64`; all `45-52` frames keep
  height `52` and baseline `57`. Parent reviewed the `45-52`, `37-52`, and
  `00-52` checker contacts. `49-52` read as a coherent reaction / upright band
  after the low `45-46` and walk `47-48` poses, not as an isolated scale jump.
  The lane was instructed to preserve `00-52`, generate only final `53-61`,
  and stop for final review.
- Parent-gated `british_shorthair_blue` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=14`; bboxes keep width `88`,
  baseline `57`, and height `35-44` across the low walking / ground-check band.
  Parent reviewed the `13-20` and `00-20` checker contacts and instructed the
  lane to preserve `00-20`, generate only `21-28`, and stop before `29`.
- Parent-gated `munchkin_brown_tabby` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=17`; bboxes keep width `88`,
  baseline `57`, and height `29-41`. Parent reviewed the `13-20` and `00-20`
  checker contacts; `20` reads as a low crouch / ground-check frame, not as a
  size jump. The lane was instructed to preserve `00-20`, generate only
  `21-28`, and stop before `29`.
- Parent-gated `toy_poodle_apricot` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=0`; bboxes keep height `52` and
  baseline `57`, with `20` as a wider low ground-check start frame. Parent
  created and reviewed the `13-20` and `00-20` checker contacts and instructed
  the lane to preserve `00-20`, generate only `21-28`, continue naturally from
  the low `20` pose, and stop before `29`.
- Promoted `domestic_shorthair_calico` as a reusable accepted `set00` source
  asset, not as a runtime, Pages, release, tag, download, or deploy change.
  Parent final QA accepted the `53-61`, `45-61`, and full `00-61` contacts on
  checker, dark, and light backgrounds; calico patches, visible paws/contact,
  tail, ears, and baseline remain readable. Mechanical QA passed with
  `auditframes` `valid=62 missing=0 invalid=0 warnings=70`, and
  `assemblemotion` produced
  `docs/art-source/domestic-shorthair/motion-source/sheets/domestic-shorthair-calico-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `domestic_shorthair_calico` variant as `motion_source_accepted`; runtime
  variants, Pages, release notes, tags, downloads, and deploy outputs were not
  changed.
- Revalidated the parent branch after Calico promotion:
  `go run ./cmd/validatemotion -variant domestic_shorthair_calico -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`,
  `go run ./cmd/validatemotion -runtime-only -require-accepted`, and
  `COPYFILE_DISABLE=1 git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets currently have one set, not
  the full 10-set release gate.
- Parent-gated `british_shorthair_blue` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=18`; `21-25` form the low
  ground-check band and `26-28` recover toward upright poses without a
  one-frame size spike. Parent reviewed the `21-28` and `00-28` checker
  contacts and instructed the lane to preserve `00-28`, generate only `29-36`,
  and stop before `37`.
- Parent-gated `toy_poodle_apricot` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=1`; `21-23` form a low ground-check
  band and `24-28` recover toward upright poses. Parent created and reviewed
  the `21-28` and `00-28` checker contacts and instructed the lane to preserve
  `00-28`, generate only `29-36`, and stop before `37`.
- Parent-gated `munchkin_brown_tabby` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=22`; bboxes keep width `88`,
  baseline `57`, and height `27-37` across the low ground-check band. Parent
  reviewed the `21-28` and `00-28` checker contacts and instructed the lane to
  preserve `00-28`, generate only `29-36`, and stop before `37`.
- Parent-gated `british_shorthair_blue` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=22`; `32-35` form a coherent
  front-facing band, and `36` returns toward the right-facing pose. Parent
  reviewed the `29-36` and `00-36` checker contacts and instructed the lane to
  preserve `00-36`, generate only `37-44`, and stop before `45`.
- Parent-gated `munchkin_brown_tabby` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=25`; `33-34` form a front-turn band
  and `35-36` return toward the right-facing pose. Parent reviewed the `29-36`
  and `00-36` checker contacts and instructed the lane to preserve `00-36`,
  generate only `37-44`, and stop before `45`.
- Parent-gated `toy_poodle_apricot` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=1`; `32-36` form a coherent
  front-facing band after the upright walk frames. Parent reviewed the `29-36`
  and `00-36` checker contacts and instructed the lane to preserve `00-36`,
  generate only `37-44`, and stop before `45`.
- Parent-gated `british_shorthair_blue` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=27`; `40-44` form a coherent low
  ground-check band. Parent reviewed the `37-44` and `00-44` checker contacts
  and instructed the lane to preserve `00-44`, generate only `45-52`, and stop
  before `53`.
- Parent-gated `munchkin_brown_tabby` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=26`; bboxes keep width `88`,
  baseline `57`, and height `26-41`. Parent reviewed the `37-44` and `00-44`
  checker contacts; `44` is a low ground-check pose, not a body-size spike. The
  lane was instructed to preserve `00-44`, generate only `45-52`, and stop
  before `53`.
- Parent-gated `toy_poodle_apricot` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=1`; all frames keep height `52` and
  baseline `57`, with width changes matching the turn and ground-check poses.
  Parent reviewed the `37-44` and `00-44` checker contacts and instructed the
  lane to preserve `00-44`, generate only `45-52`, and stop before `53`.
- Parent-gated `british_shorthair_blue` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=31`; `45-48` stay low and `49-52`
  recover into upright alert/rest poses with stable blue-gray coat, round
  head/body, tail attachment, and visible paws/contact. Parent reviewed the
  `45-52` and `00-52` checker contacts and instructed the lane to preserve
  `00-52`, generate only final `53-61`, and stop for final review.
- Revised the active asset/release target after user confirmation: the target
  set is now 19 animals, adding `quokka`, `roborovski_hamster`, and
  `guinea_pig_russian_smoke_white` to the previous 16-count plan. The user also
  explicitly requested the lane to proceed through asset QA, settings UI/menu
  confirmation for every target animal, and Release. Child generation lanes
  remain asset-only; the parent lane now owns runtime/UI/page/release
  integration after all 19 assets pass parent QA.
- Promoted `british_shorthair_blue` as the 14th accepted source asset in the
  19-target plan. Parent accepted the full `00-61` contact evidence from the
  child lane and copied the 62 source frames to
  `docs/art-source/british-shorthair/motion-source/accepted-frames/set00/`.
  Mechanical QA passed with `auditframes` `valid=62 missing=0 invalid=0
  warnings=31`; `assemblemotion` produced
  `docs/art-source/british-shorthair/motion-source/sheets/british-shorthair-blue-source-set00.png`.
  Catalog metadata now marks `british_shorthair_blue` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Revalidated the parent branch after British Shorthair promotion:
  `go run ./cmd/validatemotion -variant british_shorthair_blue
  -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, and
  `COPYFILE_DISABLE=1 git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets currently have one set, not
  the full 10-set release gate.
- Sent final-frame child continuations for `munchkin_brown_tabby` and
  `toy_poodle_apricot`, each scoped to generate only `53-61`, create
  `53-61`/`45-61`/`00-61` contacts, run strict `auditframes`, and stop for
  parent final QA without catalog/runtime/UI/Pages/release edits.
- Added mechanical motion-consistency review to `cmd/auditframes` via
  `-motion-warnings`. The new warnings flag adjacent or isolated changes in
  bbox width/height, contact baseline, body alpha area, and bbox fill ratio so
  sudden size jumps and body-ratio outliers are review candidates before parent
  acceptance. This does not make warnings fatal; the parent must still inspect
  flagged frames on contact sheets and record whether they are intentional
  crouch/turn/groom poses or true rejects. `go test -buildvcs=false
  ./cmd/auditframes` passed, and British Shorthair's new motion report flags
  frames `12`, `25`, `26`, `40`, and `49` for contact-sheet review.
- Promoted `toy_poodle_apricot` as the 15th accepted source asset in the
  19-target plan. Parent `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=1`;
  the lone `frame-27` lower-ledge warning was reviewed on the `21-28` and full
  `00-61` contacts and accepted as connected curly body/leg fur, not a floor,
  prop, shadow, or detached shelf. The `45-61` contact keeps stable size, body
  ratio, and baseline. The accepted frames now live under
  `docs/art-source/toy-poodle/motion-source/accepted-frames/set00/`, and
  `assemblemotion` produced
  `docs/art-source/toy-poodle/motion-source/sheets/toy-poodle-apricot-source-set00.png`.
  Catalog metadata now marks `toy_poodle_apricot` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Revalidated after Toy Poodle promotion: `go run ./cmd/validatemotion -variant
  toy_poodle_apricot -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./cmd/auditframes ./cmd/importanimals
  ./internal/catalog`, and `COPYFILE_DISABLE=1 git diff --check` passed.
- Promoted `munchkin_brown_tabby` as the 16th accepted source asset in the
  19-target plan. Parent `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=37`;
  most warnings are lower ledge/floor candidates from the low, long Munchkin
  body. Motion warnings at `20` and `44` were reviewed on slot contacts and
  accepted as lower sniff/ground-check pose changes, not sudden scale errors.
  The accepted frames now live under
  `docs/art-source/munchkin/motion-source/accepted-frames/set00/`, and
  `assemblemotion` produced
  `docs/art-source/munchkin/motion-source/sheets/munchkin-brown-tabby-source-set00.png`.
  Catalog metadata now marks `munchkin_brown_tabby` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Revalidated after Munchkin promotion: `go run ./cmd/validatemotion -variant
  munchkin_brown_tabby -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./cmd/auditframes ./cmd/importanimals
  ./internal/catalog`, and `COPYFILE_DISABLE=1 git diff --check` passed.
- Parent-gated new 19-target lanes through `05-12` where usable:
  `roborovski_hamster` passed with parent `auditframes` `valid=13 missing=49
  invalid=0 warnings=3` and stable sandy Roborovski read;
  `guinea_pig_russian_smoke_white` passed with `valid=13 missing=49 invalid=0
  warnings=22`, where lower-body/floor heuristic warnings were accepted after
  checker/dark contact review confirmed a white guinea pig with gray ear/nose
  and no detached shadow. `quokka` current `00-12` was rejected at 4x review as
  too rodent/kangaroo-rat-like because of the long thin tail and low rat-like
  body; the worker was interrupted to preserve it as rejected-reference and
  restart `00-04` with stronger wallaby/quokka constraints.
- Parent-gated the remaining three 19-target lanes through later partial
  ranges. `roborovski_hamster` now passes through `45-52` with
  `auditframes -artifact-warnings -motion-warnings` reporting
  `valid=53 missing=9 invalid=0 warnings=63`; the `45-52` contact confirms the
  width changes are turn/upright pose changes, not isolated size jumps. The user
  noted the Roborovski tone should be corrected, so a parent candidate
  color-only correction for frames `29-52` was prepared toward the brighter
  `00-28` sandy coat while preserving white belly/cheek/eyebrow and pink
  ear/foot/nose pixels. It will be applied only after final `53-61` exists, with
  originals and a correction report preserved. `quokka` regenerated art now
  passes through `29-36` with `valid=37 missing=25 invalid=0 warnings=44`; the
  stout attached tail and upright wallaby/quokka read avoid the rejected
  long-tail rodent direction. `guinea_pig_russian_smoke_white` now passes
  through `45-52` with `valid=53 missing=9 invalid=0 warnings=101`; width stays
  `84`, fill ratio remains about `806-848` permille, and the warnings are broad
  white-body lower-edge heuristics rather than detached floor or body-size
  spikes. `numpy 2.5.0` was installed into the user-site Python 3.14 package
  directory after direct system pip install was blocked by the externally
  managed Python environment.
- Promoted `roborovski_hamster` as the 17th accepted source asset in the
  19-target plan. Parent final `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=81`.
  Frames `29-52` received a documented color-only correction toward the early
  `00-28` sandy coat after the user called out tone drift; originals and the
  correction report were preserved under
  `docs/art-source/roborovski-hamster/motion-source/qa/`. The accepted frames
  now live under
  `docs/art-source/roborovski-hamster/motion-source/accepted-frames/set00/`,
  and `assemblemotion` produced
  `docs/art-source/roborovski-hamster/motion-source/sheets/roborovski-hamster-source-set00.png`.
  Catalog metadata now marks `roborovski_hamster` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Promoted `guinea_pig_russian_smoke_white` as the 18th accepted source asset
  in the 19-target plan. Parent final `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=119`;
  the high warning count was accepted after contact review because the broad
  white low body triggers lower-edge heuristics while width, baseline, fill
  ratio, and Russian-smoke gray ear/nose read remain stable. The accepted frames
  now live under
  `docs/art-source/guinea-pig-russian-smoke-white/motion-source/accepted-frames/set00/`,
  and `assemblemotion` produced
  `docs/art-source/guinea-pig-russian-smoke-white/motion-source/sheets/guinea-pig-russian-smoke-white-source-set00.png`.
  Catalog metadata now marks `guinea_pig_russian_smoke_white` as a distinct
  `motion_source_accepted` variant rather than reusing existing tricolor guinea
  pig art. Validation passed: both new variants pass
  `validatemotion -require-accepted`, `go run ./cmd/importanimals` imported
  `105` seed variants, and `go test -buildvcs=false ./cmd/auditframes
  ./cmd/importanimals ./internal/catalog` passed.
- Added a Windows tray-menu temporary visibility toggle for `v0.2.2` follow-up
  testing. The right-click menu now offers `Hide temporarily` / `Show` (and the
  Japanese equivalents), hides only the transparent pet overlay, suppresses
  hover names and click reactions while hidden, and deliberately does not
  persist the hidden state to `settings.json`. Verified
  `go test -buildvcs=false ./cmd/animalsdesktop`, `go test -buildvcs=false
  ./...`, `git diff --check`, and a local Windows GUI build launched from
  `dist/AnimalsDesktop.exe` with `main.appVersion=v0.2.2`.
- Prepared `v0.2.3` as a Windows UI hotfix after the name-change dialog showed
  clipped Save / Cancel buttons. The dialog now creates a captioned window from
  a fixed client area instead of treating the outer window size as the client
  size. Added regression coverage for rename-dialog controls and settings
  footer controls, updated Pages/Release metadata to `v0.2.3`, and kept the
  runtime animal scope unchanged from `v0.2.2`.
- Completed the 19/19 accepted-source wave by promoting `quokka`. Parent final
  `auditframes -strict -artifact-warnings -motion-warnings` passed with
  `valid=62 missing=0 invalid=0 warnings=90`; frame `57` needed two rejected
  replacements before the final stout, attached-tail quokka read was accepted.
  Accepted frames, contact sheets, QA notes, assembled source sheet, catalog
  metadata, generated seed source, and 10 runtime sprite sheets now exist for
  `quokka`.
- Per the user pre-release visual flags, repaired the `gecko_leopard` middle
  turn scale jump by resizing/re-anchoring frames `33-35`, and repaired
  `domestic_shorthair_calico` contact drift by re-anchoring frames `07`, `27`,
  and `53`. Originals and JSON reports were preserved under each animal's
  `motion-source/qa/` repair directory. Post-repair `auditframes` passed for
  both accepted source sets, with contacts reviewed on checker/light/dark
  backgrounds.
- Integrated all 19 new accepted animals into the v0.2.4 runtime/page release
  scope, raising runtime selection to 35 variants. GitHub Pages now shows 35
  current animal cards, 12 remaining upcoming silhouettes, v0.2.4 download
  links, v0.2.4 version-history copy, and explicit future roadmap text for the
  remaining Pages candidates, the white lionhead-pattern rabbit, and the
  special low-motion shoebill. Release workflow prerelease handling and
  `docs/releases/v0.2.4.md` were updated.
- Release-prep validation passed locally: `go run ./cmd/importsheet`, `go run
  ./cmd/importanimals`, `python3 scripts/build_page_assets.py`,
  `python3 scripts/verify_page_release.py`, JS syntax check for
  `docs/index.html`, `go run ./cmd/validatemotion -runtime-only
  -require-accepted` (35/35 accepted source; 35 expected one-set preview
  warnings), `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`,
  Windows amd64/386 cross-builds, Windows amd64 compile-only test, macOS
  arm64/amd64 ZIP builds, and `COPYFILE_DISABLE=1 git diff --check`. Playwright
  browser QA confirmed the JP and EN Pages views, 35 animal cards, 12 upcoming
  cards, v0.2.4 version history, and future roadmap text render.
- Added `lionhead_rabbit` and `shoebill` as explicit public Upcoming cards after
  the v0.2.4 release. Both use page-only ImageGen source art under
  `docs/art-source/*/page-coming-soon/`, black silhouettes under
  `docs/assets/upcoming-silhouettes/`, and do not touch runtime sprite truth.
  The public Upcoming contract is now 14 silhouettes: the 12 remaining Pages
  candidates, the user-specified white lionhead-pattern lion rabbit, and the
  special low-motion shoebill. Local validation passed with
  `scripts/build_page_assets.py`, `scripts/verify_page_release.py`, JS syntax
  check, workflow silhouette existence check, `git diff --check`, and
  Playwright JP/EN card/image-load QA.
## 2026-06-29

- Fixed Windows mixed-height multi-monitor overlay placement for 4K + 1080p
  spans. The Windows renderer now keeps one horizontal scene coordinate system
  but draws separate thin layered overlay strips per selected monitor, so each
  display uses its own taskbar/screen bottom instead of forcing the 4K display
  to share the 1080p work-area bottom. Click, hover-name, forage, wheel, and
  reaction drawing were adjusted to convert through the per-strip scene offset;
  reaction bubbles are skipped outside their strip to avoid duplicate clamped
  edge bubbles. Windows DPI awareness was raised from `system` to `per monitor
  v2` in `winres/winres.json`, with a runtime `SetProcessDpiAwarenessContext`
  fallback for local builds before a manifest is embedded. Validation passed:
  `go test -buildvcs=false ./cmd/animalsdesktop`, `go test -buildvcs=false
  ./...`, `go vet -buildvcs=false ./...`, `go build -buildvcs=false
  -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop`,
  `go run ./cmd/winresicon -src docs/assets/animalsdesktop-preview.png -out
  winres/icon.png`, temp `go-winres` generation for the updated manifest,
  smoke launching `dist\AnimalsDesktop.exe`, and `git diff --check`.
  Follow-up live check after adding another low-resolution display detected
  three screens (`1920x1080` primary, `1920x1080` left secondary, `1024x768`
  right secondary); cache-bypassed multi-monitor tests passed and the rebuilt
  `dist\AnimalsDesktop.exe` launched successfully.
- Follow-up screenshot QA covered a left monitor set to 150% scaling. The first
  PowerShell screenshot path was misleading because the capture process itself
  was DPI-virtualized; after making the capture process per-monitor DPI aware,
  screenshots confirmed animals render on all three displays and the 150%
  display keeps the same apparent sprite size. The app process reported
  per-monitor DPI awareness, and the Windows renderer now scales each segment's
  logical 96dpi canvas to the target monitor DPI. Revalidation passed with
  `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, and the
  Windows GUI build; final `git diff --check` also passed. DPI-aware screenshot
  evidence is stored under
  `.codex/screenshots/mixed-dpi-20260629-195339/dpi-aware-capture-200839/`.
- Prepared the mixed-DPI multi-monitor fix as `v0.2.5`, a preview hotfix on top
  of the `v0.2.4` 35-animal roster. Updated the default Windows app version,
  release workflow prerelease condition, release notes, README, Pages download
  links, Pages version history, and Pages release verifier. Local release-prep
  validation passed with `go run ./cmd/importsheet`, `go run ./cmd/importanimals`,
  `py -3 scripts/build_page_assets.py`, `py -3 scripts/verify_page_release.py`,
  `go run ./cmd/validatemotion -runtime-only -require-accepted`, `go test
  -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, Windows resource
  generation via `cmd/winresicon` and `go-winres`, `go build -buildvcs=false
  -ldflags="-H=windowsgui -X main.appVersion=v0.2.5"`, and a launch smoke test
  of `dist/AnimalsDesktop.exe`.

## 2026-06-30

- Prepared the v0.2.6 release lane by combining the mixed-DPI multi-monitor
  hotfix with six additional selectable runtime variants:
  `parrotlet_blue_green`, `true_albino_chipmunk`,
  `miniature_schnauzer_salt_pepper`, `japanese_giant_salamander`,
  `white_wagtail`, and `domestic_shorthair_tabby_white_stocky`. The earlier
  `albino_chipmunk` label was corrected to black-eyed white chipmunk so the new
  true albino chipmunk remains a separate no-stripe red-eyed asset.
- Refreshed selected accepted assets for budgerigar, lovebird, quokka,
  Roborovski hamster, Ragdoll, Scottish Fold, and Toy Poodle after comparing
  local accepted sources against the v0.2.4 release state. The existing leopard
  gecko source sheet was kept, and runtime-only render adjustments now scale
  the middle turn frames instead of adopting the worse alternate sheet.
- Updated GitHub Pages and release copy to v0.2.6, 41 current animals, and
  retained the Upcoming contract as 14 cards: 12 remaining Pages candidates,
  the white lionhead-pattern rabbit, and the special low-motion shoebill.
- Local validation passed: `go run ./cmd/importsheet`, `go run
  ./cmd/importanimals`, `python3 scripts/build_page_assets.py`, `python3
  scripts/verify_page_release.py`, `go run ./cmd/validatemotion -runtime-only
  -require-accepted`, `go test -buildvcs=false ./...`, `go vet
  -buildvcs=false ./...`, `git diff --check`, macOS arm64/amd64 ZIP builds,
  Windows amd64/386 cross-builds, and Playwright JP/EN Pages smoke checks.
  Visual QA artifacts are under `.codex/qa/release-v025-visual/`.

## 2026-07-01

- Prepared the v0.2.8 release lane on top of published v0.2.7 by promoting the
  two complete unpublished ImageGen accepted sources into runtime and Pages:
  `lionhead_rabbit_brown_white` and `shoebill_stork`. Both source directories
  were copied from the asset-production worktree with accepted frames, source
  sheets, contact sheets, and QA reports; large raw/fullrun generation
  directories stayed out of the release commit.
- Catalog/runtime scope now exposes 42 animals. `true_albino_chipmunk` remains
  excluded from the public runtime under the ImageGen-only constraint until a
  new ImageGen repair or edit lane passes white-background readability review.
- Rebuilt runtime sprites and public Pages assets with 42 current animal icons,
  12 remaining upcoming silhouettes, v0.2.8 download links, and v0.2.8 version
  copy. `lionhead_rabbit_brown_white` and `shoebill_stork` moved from Upcoming
  to current animals.

## 2026-07-02

- Prepared the v0.2.9 release lane by promoting the 12 remaining Pages
  candidates into selectable runtime animals, bringing the public catalog to 54
  animals on Windows and macOS. The first pass still excluded the true albino
  chipmunk pending white-background readability review; the public slot was
  replaced later in this same iteration.
- Release validation covered runtime imports, accepted-motion validation,
  Pages asset rebuild and verifier, full Go tests/vet, macOS arm64/amd64 ZIP
  builds, local JP/EN Pages Playwright checks, public Pages checks, and
  GitHub Actions Release/Site Artifact/GitHub Pages success.
- During release verification, `SHA256SUMS.txt` was found to include only
  Windows package hashes. The v0.2.9 release asset was regenerated from the
  published ZIPs, and the release workflow now rebuilds checksums in the
  publish job so macOS and Windows ZIPs are covered together.
- After visual review of the public release page, the hero preview image was
  adjusted from three rows to four rows for the 54-animal roster so the animals
  have more horizontal breathing room while keeping the page image dimensions
  stable.
- Replaced the runtime and GitHub Pages black-eyed white chipmunk slot with
  `true_albino_chipmunk`, leaving the legacy black-eyed source only as cataloged
  source evidence. Rebuilt Pages icons and preview assets, and verified the
  runtime list still contains 54 animals with `albino_chipmunk` absent.
- Ran the first real-color correction pilot for `fancy_rat_blue_hooded` after
  checking rat color/marking references: 62 accepted frames had only the cool
  hood/spine pixels shifted toward slate gray-blue, preserving the clean white
  hooded body. QA artifacts are stored under
  `docs/art-source/fancy-rat-blue-hooded/motion-source/qa/slate-gray-correction-20260702/`.
- Prepared the v0.2.11 hotfix after the `true_albino_chipmunk` source was
  rejected for the current ImageGen-only release constraint. The public runtime,
  Pages current grid, page asset builder, release verifier, and workflow checks
  now expose 55 selectable animals and keep the albino chipmunk slot held for a
  fresh ImageGen-only candidate lane.
- Local hotfix validation passed: `go run ./cmd/importanimals`, `go run
  ./cmd/validatemotion -runtime-only -require-accepted`, `python3
  scripts/build_page_assets.py`, `python3 scripts/verify_page_release.py`,
  `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, macOS
  arm64/amd64 ZIP builds, and `git diff --check`.

## 2026-07-03

- Prepared the v0.2.12 release lane by promoting two highest-priority
  ImageGen-only accepted sources: `longhair_hamster_black_white_masked` as the
  second black-and-white longhair hamster type, and the rebuilt Direction B
  `true_albino_chipmunk` as a red-eyed, no-stripe albino chipmunk. The old
  rejected albino source was replaced in accepted frames and runtime sprites.
- Runtime and Pages now expose 57 selectable animals. The public copy, current
  animal grid, preview image, release verifier, release notes, catalog tests,
  and macOS variant mirror tests were updated to include
  `true_albino_chipmunk` and `longhair_hamster_black_white_masked`.
- The ImageGen production rule is now explicit in both the local workflow Skill
  and repo rules: production motion art is one ImageGen call per one frame;
  generated grids/sheets are reference-only and must not be split into accepted
  production frames.
- Local validation passed: `go run ./cmd/importsheet`, `go run
  ./cmd/importanimals`, `go run ./cmd/validatemotion -runtime-only
  -require-accepted`, `python3 scripts/build_page_assets.py`, `python3
  scripts/verify_page_release.py`, `go test -buildvcs=false ./...`, `go vet
  -buildvcs=false ./...`, `git diff --check`, macOS arm64/amd64 ZIP builds,
  extracted ZIP arch/version/signature checks, and an arm64 app launch smoke
  test using a temporary HOME.

- Prepared the v0.2.14 settings hotfix after a report that a previously
  selected low-motion shoebill could load as a cat. The cause was legacy
  settings storing only the runtime picker index; after v0.2.12 roster
  additions, the old shoebill index pointed at a cat slot. Windows and macOS now
  save stable runtime variant IDs alongside compatibility indices, recover
  named legacy shoebill settings, leave ambiguous legacy numeric indices
  unchanged, and group animal pickers by broad animal type. The public roster
  remains 57 animals with no new asset changes in this hotfix.

## 2026-07-04

- Prepared the v0.2.15 Windows settings release by adding per-pet fixed/random
  slot modes, per-slot random type filters, and a global random type filter
  reachable from settings and the tray menu. The public roster remains 57
  animals with no asset changes. Local validation covered importer
  determinism, full Go tests/vet, Windows amd64/no-network builds, and visible
  settings UI screenshots for mixed fixed/random, per-slot random menus, and
  global type-filtered random mode.

## 2026-07-23

- Preserved the existing dirty `main` checkout, fetched the latest
  `origin/main`, and created a separate clean implementation branch from the
  current remote head. No existing source art, generated evidence, or local
  work was deleted.
- Added the canonical `docs/development/adding-an-animal.md` integration path,
  catalog contract validation, and `cmd/importanimals -variant <id> -check`.
  Targeted mode validates or imports one animal without overwriting the
  aggregate seed report or preview.
- Made the Go catalog's curated runtime order the shared source for Windows
  and macOS tests plus Pages generation. A strict Python parser rejects
  formatting drift and quoted IDs in comments, while Pages and tagged-release
  workflows now run the independent runtime-versus-page and exact-icon-set
  verifier.
- Unified motion-source family resolution across import and release
  validation. `set00` alone remains an explicitly warned preview fallback, but
  any partially populated `set00`-`set09` family now fails instead of silently
  duplicating one sheet.
- Local development dependencies were completed with Go 1.26.5, Python 3.14.6,
  and pinned Pillow 12.3.0. A fresh 57-animal Pages asset rebuild was
  deterministic and produced no tracked image changes.
- Final local validation covered the uncached full Go suite, `go vet`, Windows
  GUI build, all 127 seed variants through no-write import checking, Python
  parser tests, Pages verification, and `git diff --check`. The current
  57-animal preview roster still uses the documented single-source-set release
  exception; all entries are accepted sources, but they are not the future
  full ten-unique-set content gate. No animal pixels, tags, releases, or
  production Pages state were published in this iteration.

## 2026-07-24

- Added and hardened the `cmd/coatbatch` four-cell coat-only workflow. Each
  call binds the manifest, prompt, palette swatch, source frames, canonical
  base frames, and outputs by SHA-256; supports explicit non-promotable filler
  cells; preserves alpha and chroma pixels during one uniform bounded tone
  adjustment; and stages output/report sets with ordinary-error rollback.
- Hardened `cmd/prepareframe` with an optional one-pixel canonical alpha-bbox
  geometry lock, source/reference/output hashes, collision rejection, and
  transactional output/report writes. Replaying the retained four-frame Sable
  pilot kept alpha/chroma unchanged, matched all canonical bboxes, and differed
  from the earlier calibrated result by at most one RGB level.
- Added fail-closed silhouette gates to the production coat tools.
  `cmd/coatbatch` now records and enforces raw per-cell IoU of at least `98.5%`
  and centroid movement of at most `1.25px`, including filler cells.
  `cmd/prepareframe -match-alpha-bbox` records and enforces final IoU of at
  least `98.0%` and centroid movement of at most `0.30px` before writing.
  The retained real Sable pilot replay passed with raw minimum IoU `98.902%`,
  raw maximum centroid movement `1.044px`, final minimum IoU `98.195%`, and
  final maximum centroid movement `0.154px`.
- Runtime review found that the ferret turn is a standalone Windows one-shot:
  frames `32-39` play twice each, then direction flips and walk frame `04`
  resumes. It is not a required `31 -> 32 -> ... -> 40` sequence, and Darwin
  currently does not select frames `32-47`.
- Added `cmd/auditframes -motion-boundaries` so production audits can suppress
  cross-action false positives while retaining within-action motion and
  artifact warnings. The ferret layout uses
  `4,12,20,26,32,40,48,56`; loop closures and turn-to-walk exits remain visual
  runtime-contact gates.
- Completed and parent-approved the corrected Sable Panda base at `62/62`.
  Final parent review included complete checker/dark contacts, the rest and
  alert loops, and the alert transitions `59 -> 60`, `60 -> 61`, and
  `61 -> 56`. The official action-aware audit reports `valid=62`,
  `missing=0`, `invalid=0`, `warnings=0`; shared matte passes `62/62`, and the
  exact-duplicate scan reports no pairs.
- Restarted the dedicated Sable and Albino tasks for the user-requested
  four-cell same-species coat-only experiment. Both start with bounded
  `00-03` and `04-07` calls. Sable uses the already approved `0.573025` tone
  target and runs through normalization and visual QA; Albino stops after raw
  measurement so its target tone can be approved from two independent calls.
  Existing one-frame variant evidence remains preserved. No commit,
  publication, tag, release, or production Pages change was made.
- Corrected a protected-checkout fingerprint provenance error. The unchanged
  25,728 porcelain lines hash to the historical
  `5682fcb75549f352404a544d34a30efe27bf21a62e3054b3bb1888b952128da4`
  with PowerShell culture sorting, but to
  `0d7c471812644c51220002d5707f0def40902d66c389d873919a26ebf508ebc2`
  with true ordinal sorting. Future gates record both instead of reporting
  the ordering-method difference as a filesystem drift.
- Added two named, fail-closed `coatbatch` recovery policies after independent
  review of the first fresh coat gates. Strict defaults remain raw IoU `98.5%`,
  centroid movement `1.25px`, and tone gain `0.95-1.05`.
  `ferret_albino_high_contrast_v1` permits raw IoU `98.0%` only for
  `ferret_albino`; `ferret_sable_exact_recovery_v1` permits gain down to `0.85`
  only for the reviewed Sable `04-07` input SHA-256 and target `0.573025`.
  Unknown IDs and mismatched species, hash, or target fail closed, and no
  arbitrary numeric override was added. Focused tests and vet passed. Formal
  replay can qualify the measurement method but cannot make a sheet derivative
  promotable.
- Completed the bounded four-cell direction test for the opening Sable and
  Albino poses. Sable `00-07` passed the diagnostic tone, geometry, matte,
  continuity, and runtime checks; Albino `00-07` passed the named raw-geometry
  measurements and established a reference tone ratio. The repository rule
  still requires one ImageGen call per accepted production frame, so no sheet
  cell or derived slice is eligible for promotion. Both coat lanes restarted
  `00-07` as standalone one-frame calls using the approved Panda poses and
  pose-free coat swatches.
- Parent-approved the standalone Sable and Albino `00-07` gates after direct
  checker/light/dark contact review. Sable used eight ImageGen calls for eight
  winners and passed exact-bbox, minimum IoU `0.992494`, maximum-centroid
  `0.116912px`, matte `8/8`, artifact, duplicate, identity, and continuity
  checks. Albino used nine calls for eight winners: frame `02` attempt 01 was
  rejected for a transparent pinhole and replaced by a fresh call, with no
  local repair. Its winners passed exact-bbox, minimum IoU `0.989983`,
  maximum-centroid `0.145234px`, matte `8/8`, artifact, duplicate, identity,
  and contrast checks. All accepted hashes differ from sheet-derived evidence.
- Promoted the parent-approved Sable Panda `62/62` base byte-for-byte into
  `docs/art-source/ferret-sable-panda/motion-source/`, including provenance,
  canonical/audit/matte/duplicate evidence, parent verdict, review contacts,
  and assembled source sheet. The promoted source passes action-aware audit
  (`valid=62`, `warnings=0`), matte `62/62`, assembly, targeted import check,
  and accepted-source validation. A targeted local import generated the source
  and ten runtime sheets; nothing was published, tagged, pushed, or released.
- Parent-approved standalone Sable and Albino frames `08-15` after reviewing
  the complete `00-15` light/checker/dark contacts and the slink-to-scurry
  boundary. Sable used eight calls for eight winners and passed minimum IoU
  `0.993735`, maximum centroid `0.053508px`, audit warnings `0`, matte `16/16`,
  duplicates `0`, and sheet-hash collisions `0`. Albino used nine calls for
  eight winners; frame `12` attempt 01 was rejected for a one-pixel transparent
  pinhole and replaced by a fresh call without repair. It passed minimum IoU
  `0.984615`, maximum centroid `0.253989px`, audit warnings `0`, matte `16/16`,
  duplicates `0`, and sheet-hash collisions `0`. Prefix and protected-checkout
  hashes remained unchanged.
- Parent-approved standalone Sable and Albino frames `16-23` after reviewing
  the complete `00-23` contacts, scurry closure `19 -> 12`, action exit
  `19 -> 20`, and sniff progression `20-23`. Sable used eight calls for eight
  winners with no retry and passed minimum IoU `0.994427`, maximum centroid
  `0.081928px`, audit warnings `0`, matte `24/24`, duplicates `0`, and
  sheet-hash collisions `0`. Albino used ten calls for eight winners; frame
  `17` attempts 01 and 02 were rejected for the same one-pixel transparent
  pinhole and attempt 03 was a fresh passing call without repair. It passed
  minimum IoU `0.986953`, maximum centroid `0.181948px`, audit warnings `0`,
  matte `24/24`, duplicates `0`, and sheet-hash collisions `0`; independent
  read-only QA found no blocking defect. Both lanes advanced only to the
  bounded `24-31` gate, with no promotion or publication.
- Parent-approved standalone Sable and Albino frames `24-31` after reviewing
  the complete `00-31` contacts, sniff closure `25 -> 20`, action exit
  `25 -> 26`, and six-pose groom progression `26-31`. Sable used ten calls
  for eight winners: frame `26` attempt 01 had five visible limbs and attempt
  02 failed geometry, so both remained rejected and fresh attempt 03 passed
  without repair. Sable passed minimum IoU `0.992165`, maximum centroid
  `0.138059px`, warnings `0`, and zero duplicate/sheet collisions. Albino
  used eight calls for eight winners with no retry and passed minimum IoU
  `0.983179`, maximum centroid `0.141627px`, warnings `0`, matte `32/32`,
  and zero duplicate/sheet collisions; independent read-only QA found no
  blocker. Both lanes advanced only to the one-shot turn gate `32-39`.
- Parent-approved standalone Sable and Albino frames `32-39` as one-shot
  turns after checking Panda intent, right-facing entry, frontal midpoint,
  left-facing exit, and the runtime transition to mirrored walk frames
  `04-05`. Sable used nine calls for eight winners; frame `37` attempt 01 was
  rejected for a transparent pinhole and fresh attempt 02 passed without
  repair. It passed minimum IoU `0.995413`, maximum centroid `0.056557px`,
  warnings `0`, and zero duplicate/sheet collisions. Albino used eight calls
  for eight winners and passed minimum IoU `0.992792`, maximum centroid
  `0.072064px`, warnings `0`, matte `40/40`, and zero duplicate/sheet
  collisions; independent read-only QA found no blocker. Both lanes advanced
  only to the independent right-facing creep gate `40-47`.
- Parent-approved standalone Albino frames `40-47` after checking the
  right-facing reset at frame `40`, eight distinct creep support phases,
  `47 -> 40` and `46 -> 47 -> 40 -> 41` closure, and light/dark readability.
  Eight calls produced eight winners with minimum IoU `0.982544`, maximum
  centroid `0.216529px`, warnings `0`, matte `48/48`, and zero
  duplicate/sheet collisions. Albino advanced only to the rest gate `48-55`.
  Sable remained in the `40-47` gate because fresh frame `41` candidates
  repeatedly closed a one-pixel Panda underbelly notch and correctly failed
  the lower-shelf warning; no local repair or threshold relaxation was used.
- Parent-approved standalone Sable frames `40-47` after the frame `41`
  recovery and complete creep-loop review. Eighteen calls produced eight
  winners: attempts 01 through 10 for frame `41` closed the source's one-pixel
  underbelly notch and remained rejected, while fresh attempt 11 preserved
  the exterior-green separation inside ImageGen and passed with no local
  repair. The complete set passed minimum IoU `0.991843`, maximum centroid
  `0.102941px`, warnings `0`, zero duplicate/sheet collisions, right-facing
  motion, and both loop closures. Sable advanced to rest `48-55`. Albino's
  `48-55` continuation remained paused because the exhausted task and two
  replacement starts system-errored before output, ImageGen, or file changes.
- Paused both coat lanes before `48-55` after the dedicated Codex task host
  began rejecting every continuation and replacement at startup. The failure
  reproduced with the original Sable task, a same-directory fork, new
  projectless tasks, a saved-project task, and a read-only one-line health
  check; all ended before assistant output, tool use, ImageGen, or workspace
  writes. Failed replacement tasks were archived. The approved `00-47`
  artifacts remain unchanged, and accepted pixels were not rerouted through
  the parent task or a SubAgent. While waiting for host recovery, the parent
  fixed the `56-61` alert reference contract and audited the exact winner
  mapping needed for later byte-for-byte promotion.
- Rechecked the paused boundary: neither coat has a `48-55` run directory, so
  the failed starts wrote zero generation artifacts. The protected original
  checkout remains exactly `25,728` entries with culture-sort SHA-256
  `5682fcb75549f352404a544d34a30efe27bf21a62e3054b3bb1888b952128da4`
  and ordinal-sort SHA-256
  `0d7c471812644c51220002d5707f0def40902d66c389d873919a26ebf508ebc2`.
  Focused Go tests passed for `animalsdesktop`, `coatbatch`, `prepareframe`, and
  `auditframes`; scoped vet and the six shared matte-audit Python tests passed.
  `internal/catalog` stopped only at the intentional partial-integration gate:
  the staged `ferret_sable` entry is accepted/runtime-scoped while its final
  source sheet does not yet exist. That test must pass after the remaining
  standalone frames are promoted; it is not being weakened or skipped.
  The Windows GUI binary nevertheless builds successfully from the staged
  runtime/catalog code (`dist/AnimalsDesktop.exe`, 37,149,184 bytes); final
  release verification still waits for the two accepted source sheets and
  regenerated runtime/page assets.
- Recovered the dedicated Codex ImageGen task host and completed the standalone
  Sable `48-55` rest gate. Ten built-in calls produced eight winners: frame 52
  attempt 01 was rejected for a 29px lower shelf/floor alpha run and attempt 02
  was rejected at IoU `0.961631`; fresh attempt 03 passed without local repair
  or threshold relaxation. Parent review of all original-size winners,
  light/checker/dark contacts, Panda intent comparison, `52 -> 53 -> 54`
  head lift, and `54 -> 55 -> 48 -> 49` closure found no blocker. The accepted
  range aggregate is
  `74d0dbad5bd406491508b47bbd73064e18b7fb584de00af3dd3f9756821cb404`;
  combined `00-55` is
  `22f927b8c42cedfe0c568356a80680386c1c72bf46838da7f52cc2cc8822ad5c`.
  The full prefix passes `valid=56`, `warnings=0`, matte `56/56`, duplicate
  `0`, and sheet collision `0`. Sable advanced only to the final `56-61`
  alert gate. Albino `48-55` also resumed; a mid-gate task system error was
  recovered from filesystem authority without repeating completed ImageGen
  calls.
- Parent-approved standalone Albino frames `48-55` after the resumed task
  completed all eight winners and recovered a second system interruption
  without repeating or adding ImageGen calls. Ten built-in calls produced
  eight winners: frame `52` attempts 01 and 02 were retained as rejects for a
  29px lower floor/shelf alpha run, while fresh attempt 03 passed without local
  repair. Parent review covered original-size winners, light/checker/dark
  contacts, the natural `52 -> 53 -> 54` head lift, and
  `54 -> 55 -> 48 -> 49` closure. The prefix passes `valid=56`,
  `warnings=0`, matte `56/56`, duplicate `0`, and sheet collision `0`, with
  minimum IoU `0.992281304` and maximum centroid shift `0.126642px`.
  Albino advanced only to the final `56-61` alert gate; no promotion,
  integration, or publication was performed.
- Parent-approved standalone Sable frames `56-61`, completing its independent
  one-frame source family at `62/62`. Ten built-in calls produced six winners
  and four fully retained rejects: frame `60` attempts 01 and 02 failed,
  respectively, a transparent pinhole and geometry; frame `61` attempts 01
  and 02 failed geometry. Fresh attempt 03 passed for both frames without
  local pixel repair or threshold relaxation. Parent review covered every
  original-size winner, the Panda pose-intent comparison, light/checker/dark
  contacts, `61 -> 56`, and `59 -> 60 -> 61 -> 56 -> 57`. The complete set
  passes `valid=62`, `warnings=0`, shared matte, exact-duplicate and historical
  sheet-disjoint gates. Minimum IoU is `0.987133667`, maximum centroid shift
  is `0.133973px`, and combined `00-61` aggregate is
  `d24265406e61dc1729305648180b05b4b7c5c43c744d03f56c638ecc647056dc`.
  Sable is frozen for parent-owned promotion after Albino also reaches 62;
  the lane performed no promotion, integration, or publication.
- Parent-approved standalone Albino frames `56-61`, completing all three
  ferret source families at `62/62`. Seven built-in calls produced six winners;
  frame `60` attempt 01 was retained as a geometry reject because its head and
  neck were three pixels lower than the Panda source, and fresh attempt 02
  passed without local repair. Parent review covered original-size winners,
  checker/dark contacts, the complete `00-61` set, and
  `59 -> 60 -> 61 -> 56 -> 57`. The complete set passes `valid=62`,
  `warnings=0`, shared matte `62/62`, duplicate `0`, and historical
  sheet-collision `0`. Minimum IoU is `0.988205560` and maximum centroid shift
  is `0.196277px`; approved prefix, Panda, lane, and protected-checkout hashes
  are unchanged. Albino is frozen for parent-owned promotion; the generation
  lane performed no promotion, integration, or publication.
- Promoted the parent-approved Sable and Albino sets byte-for-byte into compact
  `motion-source/` packages. Each package now contains 62 standalone accepted
  `96x64` PNGs, 62-row winner provenance, a `5952x64` assembled sheet,
  action-aware/matte/canonical/duplicate QA, parent verdict, and 12 final
  light/checker/dark contacts. All accepted PNGs match their recorded
  standalone winners; both sets have 62 unique hashes. Compact concatenated-PNG
  aggregates are
  `33c2da1526d4bd71c7ce76f1e08ca6213406a590b6059d51849514118a4a69f0`
  for Sable and
  `a16586fda8e858be602da4d4e37648488d2b734f7c01b8bda44d3766a6160590`
  for Albino.
- Replaced the procedural Sable/Albino placeholders and locally integrated
  Sable Panda, Sable, and Albino through the catalog, Windows/macOS runtime,
  generated source/runtime sprites, and development Page. Targeted imports are
  deterministic. The full 128-seed importer validated and imported, then
  preserved all 1,412 generated PNGs on its second run with aggregate digest
  `398ba1549e7ec61b3561fd71e851981854e296a3a3658863e2ba5b3bfbfeed91`.
  The local Page now lists 60 animals while public v0.2.15 download copy remains
  at 57.
- Verified the completed integration with targeted package tests, full
  `go test -buildvcs=false ./...`, full `go vet -buildvcs=false ./...`,
  Page release verification, and a Windows GUI build
  (`dist/AnimalsDesktop.exe`, 37,534,208 bytes,
  SHA-256
  `47d02b1eedaac626110e985fcb8cc136080542ee870d025b8163c2ae05c3a4db`).
  Isolated Playwright QA passed at desktop and `390x844`: all-60 expansion,
  small-animal 29 filter, three distinct Japanese ferret cards, English roster,
  all requested icon responses 200, and zero console errors.
- Moved the exact full three-family generation tree without deleting it from
  the working copy to
  `E:/Development/AnimalDesktop-asset-archives/20260724-ferret-three-variants/`.
  The compact package successfully revalidated all 124 Sable/Albino accepted
  winners from the archive. The protected original checkout still has exactly
  25,728 status entries and both fingerprints are unchanged. The only
  intentional release blocker is motion breadth: each ferret currently has one
  accepted source set, so local runtime repeats set00 into ten sheets and
  `release_ready` correctly remains false until ten independent accepted source
  sets exist.
- Closed the final adversarial-review blocker without changing accepted art.
  Replaced Sable Panda's 147-row attempt-history copy with a 62-row winner-only
  compact provenance index. An independent hash pass verified every accepted
  frame's prompt, built-in ImageGen source, raw, unique alpha, canonical
  winner, and accepted PNG. The self-contained canonical report now contains
  62 accepted paths/hashes, 62 unique frames, the `5952x64` sheet, accepted
  byte aggregate
  `fa8fff8d45c6ee27878f2dba31e688bed13fdf5fc1bd8371bc06fd242e590f9a`,
  and the archived parent-gate aggregate
  `7947c23e4fcbcf304233cbdb09231584a9e7c479c1ee404a6712816087977ca9`.
  The original turn rejection/requalification history remains explicit and
  separate from final parent acceptance.
- Added regression coverage for `ferret -> small_mammal`, all three ferrets in
  the Windows small-animal random pool, Japanese/English grouped display
  labels, and Darwin group labels. Targeted tests, full Go tests, full vet, and
  `git diff --check` pass on Windows. Independent final re-review cleared the
  provenance blocker. A macOS runner is still required to execute Darwin tests,
  and all three variants remain intentionally `release_ready=false` while only
  set00 is accepted.
- Corrected the current asset-production flow on 2026-07-26 so source coverage
  and runtime storage are reported separately. A canonical accepted 62-frame
  `set00` is the normal complete source for an animal; `cmd/importanimals`
  expands it into ten runtime sheets without an incompleteness warning.
  Optional `set00` through `set09` source families remain supported only when
  all ten sheets exist and are byte-unique. This dated correction supersedes
  the historical `release_ready=false` conclusions immediately above without
  rewriting the evidence of what the old validator reported.
- Verified the corrected contract against the live local catalog: all 60
  runtime variants report `source_sets=1`, `runtime_sets=10`,
  `release_ready=true`, and zero warnings. Regenerated the 128-entry aggregate
  import report; all 61 accepted motion-source entries have the new fields and
  zero motion warnings. A second full import preserved generated-output
  aggregate SHA-256
  `6dd1f55f7431b915e2444e5b88580703f6f4d4c289023cbb4f9f867bdf46fede`.
- Installed official Go 1.26.5 through WinGet, but Windows Smart App Control
  rejects its unsigned helper and newly built test executables. The policy was
  not weakened. Changed-package tests passed as Go WebAssembly under signed
  Node, every package compiled for Windows/amd64, the Windows GUI build passed,
  and `git diff --check` passed. Native Windows test execution and `go vet`
  remain a host-policy verification boundary.
- Aligned the v0.2.16 Pages download choices with Degu Desktop. The primary
  buttons now identify Windows x64 and Mac Apple Silicon / macOS 12+, the
  alternative links identify Windows x86 / 32-bit and Mac Intel, and
  Japanese/English disclosures explain which Mac build to choose. Extended
  `scripts/verify_page_release.py` to bind stable download selectors to the
  five platform ZIPs and check the new public copy. The release verifier,
  scoped public-copy scan, and `git diff --check` pass. Playwright verified
  JP/EN labels, both disclosure states, zero console warnings/errors, and no
  horizontal overflow at desktop, 390px, or 320px; mobile button spacing was
  tightened so the Windows label no longer leaves one Japanese character on a
  line by itself. Updated both Pages workflow copy guards to the new Apple
  Silicon label after the pre-deploy audit caught their stale exact-string
  assertion.

## 2026-10-01 — Creator links

- Added https://x.com/kdevelopk and https://kdevelopk.pages.dev/ to the current README and public page footer.
- Based this documentation-only change on latest remote main `426d6684a322`; preserved release versions, downloads, and existing page content.
