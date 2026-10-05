# Legacy Sprite Importer Inputs

This directory is kept for the old sheet importer and regression checks only.
The v0.1.5 runtime preview uses the release-scoped animal sprite sheets under
`assets/sprites/<animal>_setNN.png`; new production art should follow the
AnimalsDesktop asset-production flow instead of these legacy prompts.

Preferred current source is one accepted PNG per runtime frame in the animal's
own `docs/art-source/<animal>/motion-source` workspace, then imported through
the release asset pipeline.

Importer:

```powershell
go run ./cmd/importsheet
```

The legacy importer writes compatibility sheets, `assets/tray.ico`,
`docs/assets/legacy-sheet-preview.png`, and
`assets/source/import-report.json`. The report warns when background removal
finds no content or when source content touches an edge and may be cropped.
