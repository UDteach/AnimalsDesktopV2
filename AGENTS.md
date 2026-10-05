# AnimalsDesktop Working Rules

## Project Goal

AnimalsDesktop is a Go + Win32 Windows taskbar pet app. Keep the existing lightweight Win32 architecture while improving it into a maintainable multi-species 2D sprite pet application.

## Product Constraints

- Keep Windows 10/11 amd64 support.
- Keep the transparent, always-on-top, click-through taskbar overlay.
- Keep tray operation, keyboard reaction, random stroll, typing wheel, foraging, grooming, GitHub Pages, and GitHub Release workflows working.
- Do not replace the app with a large GUI framework.
- Do not publish, push tags, or create GitHub Releases unless the user explicitly asks in the current task.

## Architecture Rules

- Separate species, coat, pet instance, animation, behavior profile, and render profile concepts.
- Do not make chinchillas or macaroni mice by recoloring degu sprites.
- Do not treat a simple recolor, vertical bob, or duplicated walk frame as a different action.
- Prefer adding migration-compatible structures before deleting old asset formats.
- Keep changes narrow and avoid unrelated refactors.

## Asset Quality Rules

- Each source frame must contain one complete animal on a transparent background.
- Frames must keep consistent camera, scale, anatomy, baseline, and contact points.
- Ears, feet, whiskers, and tails must not be cropped.
- Generated art must not include text, borders, scenery, shadows, costumes, multiple animals, or human-like poses.
- Generate each species as its own ImageGen source family. Coat variants may be expanded only when the species silhouette and motion set are already stable.
- For ImageGen animal production, use separate Codex project threads as lanes by animal/variant; do not use SubAgents to generate accepted animal pixels. SubAgents are only for read-only research, QA, or evidence review.
- Generate production motion art one frame per ImageGen call. Generated multi-cell grids/sheets may be kept as direction references only and must not be split into accepted production frames. Build contact sheets locally after normalization for review.

## Standard Checks

Run the relevant subset after each iteration:

```powershell
gofmt -w internal\catalog\catalog.go cmd\animalsdesktop\main_windows.go cmd\animalsdesktop\motion_windows_test.go cmd\importsheet\main.go cmd\importsheet\main_test.go cmd\importanimals\main.go cmd\importanimals\main_test.go
go run ./cmd/importsheet
go run ./cmd/importanimals
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

When importer behavior changes, compare generated asset diffs and ensure repeated imports are deterministic.

## Iteration Routine

1. Read the active `.codex/tasks/*.md` ledger, and read `docs/development/current-state.md` / `docs/development/iteration-log.md` when those files exist.
2. Inspect the current Git diff.
3. Pick one highest-value problem for the iteration and state the cause.
4. Implement the smallest useful change.
5. Verify with commands and visual review when UI or sprites changed.
6. Repair problems found in the same iteration when they are in scope.
7. Record results in `docs/development/iteration-log.md`.

## Codex Config Alignment

These local rules mirror the safe parts of `UDteach/codex_config` for this repository.

- For non-trivial work, start with a compact task ledger and keep exactly one critical-path item in progress.
- Use repo intake before broad edits when command structure, assets, or ownership boundaries are unclear.
- Prefer local source and tests first; use external/current docs only when the claim is version-sensitive or materially uncertain.
- If the same error survives two attempted fixes, stop patching and do an evidence pass before the next edit.
- Before finalizing non-trivial work, run an adversarial review for stale assumptions, missing tests, resource leaks, state-transition bugs, and release workflow impact.
- Keep bulky evidence in files, not chat output. Summarize only key lines.
- Do not restore global Codex config from `codex_config` into this repo environment; that backup can contain machine-specific paths. Apply only project-safe rules unless the user explicitly asks for a global restore.
- Do not push, create tags, publish releases, or change Cloudflare/GitHub production settings unless the user explicitly asks in the current task.
