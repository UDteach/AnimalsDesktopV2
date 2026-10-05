# AnimalsDesktop Ver2

AnimalsDesktop is a small desktop pet app for Windows and macOS. Ver2 introduces
individually timed animation frames and movement matched to the animals' steps.
The first Ver2 release uses the completed, integrated motion assets available now.

Creator: [kdevelopk on X](https://x.com/kdevelopk) · [More apps](https://kdevelopk.pages.dev/)

Public page: <https://udteach.github.io/AnimalsDesktopV2/>

Current app version: `v2.0.0`

Start `AnimalsDesktop.exe` on Windows or `AnimalsDesktop.app` on Mac for Ver2.
To use the established animations, including the typing wheel, quit the app and
start the included `Start-Classic.cmd` (Windows) or `Start-Classic.command` (Mac).
Classic uses the complete old motion families throughout the session. The new
seven families do not yet include wheel, foraging or grooming clips.

Ver1 settings are retained. Install Ver2 manually the first time; subsequent
update checks use the Ver2 repository. Quit Ver1 before starting Ver2.

Source and releases: <https://github.com/UDteach/AnimalsDesktopV2>

## Current Status

The picker retains 60 animals. Seven use the new walking and idle animation
system: gray, beige and ebony chinchillas, golden Syrian hamster, tan macaroni
mouse, striped chipmunk, and sable ferret. The golden hamster also has a new
three-frame input reaction. The other 53 animals retain their existing sprites.
This release does not claim that all 60 new motion families are finished.

The new animations preserve per-frame durations, scale and ground alignment.
Walking phase follows distance traveled through acceleration, deceleration and
stops. Existing animal IDs, saved selections and the version 4 settings format
remain compatible. Up to ten pets can be displayed together.

New-style pets use their completed walk/idle and input response instead of
switching into older-style wheel, grooming or foraging sprites. More dedicated
motions are still being produced and are not included until integrated and checked.

Selectable animals:

- chinchilla standard gray
- chinchilla beige
- chinchilla ebony
- golden Syrian hamster
- Djungarian hamster
- Campbell hamster
- tan macaroni mouse
- gray sugar glider
- chestnut agouti rabbit
- broken orange Holland Lop
- chestnut Netherland Dwarf
- Himalayan rabbit
- gray-brown gecko
- tricolor guinea pig
- hooded fancy rat
- Richardson's ground squirrel
- longcoat Yorkshire Terrier
- striped chipmunk
- true albino chipmunk
- leopard gecko
- blue White's tree frog
- green-yellow budgerigar
- normal gray cockatiel
- normal Java sparrow
- green parrotlet
- blue-green parrotlet
- peach-faced lovebird
- seal bicolor ragdoll
- silver tabby Scottish Fold
- fawn French Bulldog
- brown tabby Maine Coon
- calico domestic shorthair
- blue British Shorthair
- apricot Toy Poodle
- brown tabby Munchkin
- Roborovski hamster
- Russian smoke white guinea pig
- quokka
- Miniature Schnauzer
- Japanese giant salamander
- white wagtail
- tabby-white domestic shorthair
- brown-white Lionhead rabbit
- low-motion shoebill
- leucistic sugar glider
- African dormouse
- Himalayan Netherland Dwarf
- American flying squirrel
- black-and-white long-haired hamster
- black-and-white long-haired hamster 2
- yellow Djungarian hamster
- pearl white Djungarian hamster
- blue hooded fancy rat
- chocolate self fancy rat
- cream agouti fancy rat
- gray rabbit
- African fat-tailed gecko
- Sable Panda ferret
- Sable ferret
- Albino ferret

## Runtime Scope

The v2.0.0 picker exposes the same 60 stable animal IDs. The seven refreshed
animals use 1,330 checked-in PNG frames, including both display tiers and the
golden hamster reaction. Remaining animals keep the established runtime assets.
Production candidates and unfinished motion frames are not shipped.

Windows playback, selection, size changes, settings persistence and mixed old/new
pets have been exercised locally. macOS packages are built and tested in CI;
physical Mac behavior remains unverified. See [release notes](docs/releases/v2.0.0.md)
for this release's completed checks and limitations.

## Release Gate

Ver2 ships the seven integrated timed-motion families listed above. Each uses
hash-checked transparent frames, explicit frame durations, scale and ground
profiles, and distance-driven playback. The embedded manifest and every frame
are validated by `go test -buildvcs=false ./internal/motion`; application tests
cover playback, input reactions and legacy compatibility.

The remaining 53 families keep the legacy source contract: accepted canonical
62-frame source sets, generated compatibility sheets, and
`go run ./cmd/validatemotion -runtime-only -require-accepted`.

Release CI regenerates legacy assets without drift, runs Go tests and vet and
Python validation, checks the public download links, and builds Windows and
macOS packages. New motion candidates must pass source, geometry and playback
review before joining the runtime manifest.

## Development

Local asset work requires Go 1.24 or newer plus Python 3.10 or newer. Install
the pinned Python image dependency once:

```powershell
python -m pip install -r requirements-dev.txt
```

For the source-art, catalog, targeted-import, runtime, and release boundaries
for a new animal, see
[`docs/development/adding-an-animal.md`](docs/development/adding-an-animal.md).
The fast local loop validates and imports one cataloged animal without
rewriting the aggregate report or preview:

```powershell
go run ./cmd/importanimals -variant <variant_id> -check
go run ./cmd/importanimals -variant <variant_id>
go run ./cmd/validatemotion -variant <variant_id> -require-accepted
```

Useful checks during local development:

```powershell
go run ./cmd/importsheet
go run ./cmd/importanimals
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -strict -artifact-warnings -motion-warnings
go run ./cmd/validatemotion -runtime-only -require-accepted
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go run ./cmd/winresicon -src docs/assets/animalsdesktop-preview.png -out winres/icon.png
go run github.com/tc-hib/go-winres@v0.3.1 make --arch amd64 --out cmd/animalsdesktop/rsrc --file-version v2.0.0 --product-version v2.0.0
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

Windows release ZIPs are built on a Windows machine and can be uploaded to the
same release tag as the Mac ZIPs when they are ready:

```powershell
New-Item -ItemType Directory -Force dist | Out-Null
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -buildvcs=false -ldflags="-H=windowsgui -s -w -X main.appVersion=v2.0.0" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
Compress-Archive -Path dist\AnimalsDesktop.exe,README.md -DestinationPath dist\AnimalsDesktop-windows-amd64.zip -Force
```

macOS release ZIPs are built with:

```bash
VERSION=v2.0.0 GOARCH=arm64 scripts/build_macos.sh
VERSION=v2.0.0 GOARCH=amd64 scripts/build_macos.sh
```

Run `cmd/prepareframe` only on one-pose candidates, outside the standard QA loop.
It rejects checker/noisy backgrounds; prepared output still needs visual review
before it counts as accepted art:

```powershell
go run ./cmd/prepareframe -src path\to\candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\<set-id>\frame-00.png
go run ./cmd/prepareframe -background chroma-green -src path\to\green-candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\<set-id>\frame-00.png
```

Run `cmd/assemblemotion` only for a set after `cmd/auditframes -strict` passes:

```powershell
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\<set-id> -strict -artifact-warnings -motion-warnings
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\<set-id> -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-<set-id>-draft.png
```

In-progress sets are expected to fail strict mode until all 62 standalone
transparent frames exist.

## Windows Trust Notes

The Windows release workflow embeds file metadata, product metadata, a Windows 10+ manifest, and an app icon into `AnimalsDesktop.exe`. It also publishes `SHA256SUMS.txt` next to the release ZIPs and includes `SECURITY.txt` inside each ZIP with the expected EXE hash and false-positive submission notes.

The Windows release workflow also builds `AnimalsDesktop-windows-amd64-no-network.zip` as a security-check edition. That variant disables automatic update checks, tray update checks, update ZIP downloads, and update installation, and excludes the Go `net/http` update fetch/download implementation from the no-network build. It is a reduced-signal fallback for security-product review, not a Smart App Control bypass.

Microsoft's Smart App Control documentation remains the primary source for the
signing requirement. The GlobalSign April 2026 Japanese case study is useful
supporting context because it describes the same unsigned-app block pattern and
the code-signing mitigation path.

For the best Microsoft Defender SmartScreen and McAfee outcome, release builds should be Authenticode-signed with a timestamped public-trust code-signing certificate. The preferred CI path is Microsoft Azure Artifact Signing because the signing key stays in Microsoft-managed HSMs instead of being stored as a repository secret.

Configure these GitHub Secrets for Azure Artifact Signing:

- `AZURE_CLIENT_ID`
- `AZURE_TENANT_ID`
- `AZURE_SUBSCRIPTION_ID`
- `AZURE_ARTIFACT_SIGNING_ENDPOINT`
- `AZURE_ARTIFACT_SIGNING_ACCOUNT_NAME`
- `AZURE_ARTIFACT_SIGNING_CERTIFICATE_PROFILE_NAME`

If you have a legacy or private `.pfx` signing certificate, the workflow can use these fallback secrets:

- `WINDOWS_CERTIFICATE_BASE64`: base64-encoded `.pfx`
- `WINDOWS_CERTIFICATE_PASSWORD`: `.pfx` password

If signing secrets are missing, the workflow still builds and publishes checksums, but the EXE remains unsigned and may continue to receive reputation-based warnings until Microsoft/McAfee reputation or allowlisting catches up.

`v0.2.0` is retained as a mistaken Windows prerelease. `v0.2.1` is the main-line Windows trust-hardening release. `v0.2.2` keeps that trust-hardening work and adds Mac parity for animal selection, size controls, language, and display settings. `v0.2.3` is a Windows settings UI hotfix with the same sixteen-animal scope. `v0.2.4` expands the selectable animal roster to 35. `v0.2.5` fixes Windows mixed-DPI multi-monitor overlay size and placement for the 35-animal roster. `v0.2.6` expands it to 41 and adds true albino chipmunk, Miniature Schnauzer, Japanese giant salamander, white wagtail, tabby-white cat, and blue-green parrotlet. `v0.2.7` temporarily removes true albino chipmunk from the public runtime while the no-pattern albino repair lane improves white-background readability, leaving 40 selectable animals. `v0.2.8` adds the brown-white Lionhead rabbit and special low-motion shoebill, bringing the public roster to 42 animals while true albino chipmunk remains held for a new ImageGen repair lane. `v0.2.9` promotes the remaining 12 Pages candidate animals into runtime and Pages, bringing the public roster to 54 animals. `v0.2.10` adds the two chinchilla color variants, restores the true albino chipmunk runtime slot, and refreshes the blue hooded fancy rat tone, bringing the public roster to 56 animals. `v0.2.11` removes the true albino chipmunk public slot until a new ImageGen-only source passes review, bringing the public roster to 55 animals. `v0.2.12` adds a second black-and-white long-haired hamster and restores true albino chipmunk with a fresh ImageGen-only source, bringing the public roster to 57 animals. `v0.2.13` fixes Windows animal selection drift with stable animal IDs. `v0.2.14` extends stable ID selection to macOS, keeps named shoebill legacy recovery conservative, and groups animal pickers by broad type while keeping the public roster at 57 animals. `v0.2.15` adds mixed fixed/random Windows slots and animal-type filters for per-slot and global random selection while keeping the published roster at 57 animals. `v0.2.16` adds Sable Panda, Sable, and Albino ferrets, bringing the roster to 60 animals.

Do not create a stable/final release tag until the current animal target is complete.
