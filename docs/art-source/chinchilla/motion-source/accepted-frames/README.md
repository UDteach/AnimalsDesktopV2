# Chinchilla Accepted Frame Staging

This directory is reserved for reviewed one-pose-per-PNG motion source frames.

Do not place ImageGen sheets, grids, contact sheets, or opaque checker-background images here. Every accepted file must be:

- `96x64`
- a PNG with real transparent alpha around the animal
- one complete standard gray chinchilla
- no text, border, shadow, scenery, ground, props, costumes, or multiple animals
- stable camera, scale, baseline, facing direction, and contact points within the set

Single ImageGen candidates should first be staged outside this directory. If a candidate has a simple uniform background, prepare a review frame with:

```powershell
go run ./cmd/prepareframe -src path\to\candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00.png -report docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00-report.json
go run ./cmd/prepareframe -background chroma-green -src path\to\green-candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00.png -report docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00-report.json
```

Do not copy `prepareframe` output into this accepted-frame directory until it has passed visual review on light and dark backgrounds. The command rejects checker/noisy backgrounds instead of trying to clean them.

Required canonical set folder:

- `set00/frame-00.png` through `set00/frame-61.png`

Optional independent motion sources may continue through
`set09/frame-00.png` through `set09/frame-61.png`, but only as a complete,
byte-unique ten-set family. They are not required for canonical source
completion.

Current accepted progress:

- `set00/frame-00.png` through `set00/frame-61.png`: accepted after one-frame ImageGen regeneration, `prepareframe`, `auditframes`, source-sheet assembly, and motion-slot visual QA.
- The earlier single accepted `frame-00` was preserved under `stale/pre-one-frame-fullrun-20260625/` before the regenerated set replaced it for consistency.

After all 62 frames in a set pass visual review:

```powershell
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -strict
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-set00.png -report docs\art-source\chinchilla\motion-source\accepted-frames\set00-assemblemotion-report.json
```

The assembler is a gate, not an art fixer. If it rejects a frame, regenerate or manually curate the source PNG before assembling.

For optional full-family tracking across all 10 source sets:

```powershell
go run ./cmd/auditframes -root docs\art-source\chinchilla\motion-source\accepted-frames -report docs\art-source\chinchilla\motion-source\accepted-frames\audit-report.json
```
