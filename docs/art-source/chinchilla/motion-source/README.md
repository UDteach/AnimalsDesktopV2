# Chinchilla Motion Source

Status: set00_accepted
Generated: 2026-06-21; set00 regenerated 2026-06-25

## Source Batch

- Tool: ImageGen
- Thread image folder: `C:\Users\user\.codex\generated_images\019ee7c8-bfbb-7832-abcb-b8def1a72f27`
- Draft extraction source: first 16-pose chinchilla sheet
- Rejected generation: second sheet contained a non-animal symbol and is not used

The draft extraction source was generated with a checker background rather than a true transparent alpha channel. The raw checker source was kept outside the tracked tree under `.codex/tmp/chinchilla-imagegen-raw/`.

## Tracked Outputs

- `frames/chinchilla-standard-gray-pose-00.png` through `pose-15.png`
- `chinchilla-standard-gray-pose-frames-check.png`
- `chinchilla-standard-gray-frame-map.csv`
- `sheets/chinchilla-standard-gray-source-set00-draft.png`
- `chinchilla-standard-gray-source-set00-draft-check.png`

Future accepted single-frame work should be staged under:

- `accepted-frames/set00/frame-00.png` through `frame-61.png`
- `accepted-frames/set01/frame-00.png` through `frame-61.png`
- continue through `accepted-frames/set09/`

Each extracted pose frame is:

- 96x64
- transparent PNG
- one chinchilla per frame
- generated from the same ImageGen pose batch

## Current Use

The 2026-06-21 16-pose draft frames remain reference-only. They are not accepted source art and are not release-ready motion art.

The `source-set00-draft` sheet maps the 16 extracted poses into the 62-frame slot contract for review only. Do not promote it into accepted frames.

The draft was rejected for accepted use after visual QA found extra debris and doubled-face/face-like artifacts. Mechanically shifted draft `set01` through `set09` sheets were also rejected because they looked wrong and are not a substitute for accepted ImageGen motion sources.

Accepted `set00` now lives under:

- `accepted-frames/set00/frame-00.png` through `frame-61.png`

The accepted `set00` source sheet is:

- `sheets/chinchilla-standard-gray-source-set00.png`

The one-frame regeneration provenance is:

- `../one-frame-method-fullrun-20260625/chinchilla/manifest.csv`
- `../one-frame-method-fullrun-20260625/chinchilla/report.md`
- `../one-frame-method-fullrun-20260625/motion-visual-qa/chinchilla-accepted-motion-groups.png`

Accepted replacement work must be generated as one pose per PNG. Do not promote a multi-pose sheet, grid, contact sheet, or baked checker-background image as source art, even if it can be cut into 96x64 frames. A batch may contain up to four separate single-pose PNGs for review, but each accepted source frame must be a standalone transparent PNG before it is assembled into the 62-frame runtime sheet.

Use `cmd/assemblemotion` only after all 62 standalone frames in a set are accepted:

```powershell
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-set00.png -report docs\art-source\chinchilla\motion-source\accepted-frames\set00-assemblemotion-report.json
```

The assembler rejects wrong-sized frames, empty alpha, and fully opaque/checker-style backgrounds before writing a sheet.

Before release, this family still needs:

- accepted ImageGen `set01` through `set09` motion source variation
- runtime import and catalog validation
- local Go tests, vet, Windows build, and GitHub Actions success
