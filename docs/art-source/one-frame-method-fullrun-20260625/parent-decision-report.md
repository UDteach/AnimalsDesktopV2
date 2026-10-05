# One-Frame Fullrun Parent Decision - 2026-06-25

## Outcome

The one-frame retry-until-pass lane produced a complete parent-selected review set for the four requested families:

| Family | Parent-selected frames | Source mix | Contact sheet | Source sheet |
| --- | ---: | --- | --- | --- |
| hamster | 62 | main 44, rescue 18 | `parent-selected/hamster/hamster-parent-selected-contact.png` | `parent-selected/hamster/hamster-parent-selected-source.png` |
| macaroni-mouse | 62 | main 62 | `parent-selected/macaroni-mouse/macaroni-mouse-parent-selected-contact.png` | `parent-selected/macaroni-mouse/macaroni-mouse-parent-selected-source.png` |
| sugar-glider | 62 | main 62 | `parent-selected/sugar-glider/sugar-glider-parent-selected-contact.png` | `parent-selected/sugar-glider/sugar-glider-parent-selected-source.png` |
| rabbit | 62 | main 56, rescue 6 | `parent-selected/rabbit/rabbit-parent-selected-contact.png` | `parent-selected/rabbit/rabbit-parent-selected-source.png` |

The previous sheet-extracted set remains preserved as the baseline under each family `accepted-candidates/set00/`.

## Reproducibility Decision

Built-in ImageGen does not expose deterministic seeds in this workflow, so exact pixel reproduction is not available. The most reproducible practical method is:

- generate one frame per call,
- keep the stable species spec and stable 62-frame pose table,
- retain raw ImageGen output, alpha output, normalized frame, QA JSON, prompt/report, and manifest rows,
- retry only failed frame numbers,
- use parent-selected overrides when a mechanically valid frame has species drift.

This is stronger than regenerating a whole sheet because one bad frame can be replaced without invalidating the other 61.

## QA Evidence

- `parent-combined-main-rescue-audit.json`: 62/62 for hamster, macaroni-mouse, sugar-glider, and rabbit; `issues=0`.
- `go run ./cmd/auditframes -frames-dir ... -strict -artifact-warnings` on promoted accepted frames:
  - hamster: valid=62, missing=0, invalid=0, warnings=70
  - macaroni-mouse: valid=62, missing=0, invalid=0, warnings=58
  - sugar-glider: valid=62, missing=0, invalid=0, warnings=42
  - rabbit: valid=62, missing=0, invalid=0, warnings=47
- `go run ./cmd/assemblemotion` built 5952x64 source sheets for all four promoted accepted sets.
- `baseline-vs-parent-selected-overview.png` compares the existing sheet-extracted baseline against the one-frame parent-selected set.
- `motion-visual-qa/motion-visual-qa-report.md` records the accepted-frame visual motion pass.

`auditframes` warnings are not hard failures here; they remain visual review prompts.

## Visual Review Notes

- Hamster: one-frame output is more readable than the baseline, but main frames 35, 39, 48, and 55 had species drift. These are overridden to rescue candidates in `parent-selected/selection-overrides.json`.
- Macaroni mouse: one-frame output is clearer than the baseline and keeps the short thick tail/pointed nose identity. Frame 61 is slightly darker and should remain a visual review note.
- Sugar glider: one-frame output is materially stronger than the baseline; membrane frames are readable and species identity is stable.
- Rabbit: one-frame output is more detailed and readable than the baseline. Some reaction/rest frames are stylized, but they remain rabbit-like and mechanically valid.

## Adoption Position

Use `parent-selected/` as the preserved provenance and review source for the four families.

The four selected sets have been promoted into:

- `docs/art-source/hamster/motion-source/accepted-frames/set00/`
- `docs/art-source/macaroni-mouse/motion-source/accepted-frames/set00/`
- `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/`
- `docs/art-source/rabbit/motion-source/accepted-frames/set00/`

They are accepted as `set00` source-frame assets for slot-based motion. They are not full release-ready families yet: release readiness still requires accepted `set00` through `set09`, runtime import, catalog validation, and app-side review.

## Chinchilla Follow-Up

Chinchilla was completed after this four-family parent-selected pass because the existing 16-pose draft showed extra debris and doubled-face/face-like artifacts.

The accepted chinchilla set is now:

- `docs/art-source/chinchilla/motion-source/accepted-frames/set00/`

Chinchilla evidence:

- `docs/art-source/one-frame-method-fullrun-20260625/chinchilla/manifest.csv`: 62 accepted rows.
- `docs/art-source/one-frame-method-fullrun-20260625/chinchilla/report.md`: rejection reason, QA evidence, and motion verdict.
- `docs/art-source/one-frame-method-fullrun-20260625/motion-visual-qa/chinchilla-accepted-motion-groups.png`: motion visual QA sheet.
- `docs/art-source/chinchilla/motion-source/accepted-frames/set00-auditframes-report.json`: valid=62, missing=0, invalid=0, warnings=87.
- `docs/art-source/chinchilla/motion-source/sheets/chinchilla-standard-gray-source-set00.png`: accepted set00 source sheet.

With this follow-up, the five requested families have accepted `set00` source-frame coverage. They are still not full release-ready families because release readiness requires accepted `set00` through `set09`, runtime import, catalog validation, and app-side review.
