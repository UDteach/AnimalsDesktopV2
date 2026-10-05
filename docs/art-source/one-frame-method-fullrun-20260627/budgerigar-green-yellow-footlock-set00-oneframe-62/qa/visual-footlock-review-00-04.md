# Budgerigar Green-Yellow Foot-Lock Early Gate 00-04

Variant: `budgerigar_green_yellow`
Range: `00-04`
Review date: 2026-06-27

## Verdict

Early-gate candidate pass for parent review.

All five normalized `96x64` frames show at least one attached visible
foot/toe contact point. Frames `00`, `01`, `02`, and `04` show two attached
feet or foot/contact pairs clearly. Frame `03` shows smaller gray foot/contact
shapes, still attached and visible on checker and dark backgrounds.

## Evidence

- Checker contact: `contact/budgerigar-green-yellow-footlock-early-gate-00-04-checker-6x.png`
- Light contact: `contact/budgerigar-green-yellow-footlock-early-gate-00-04-light-6x.png`
- Dark contact: `contact/budgerigar-green-yellow-footlock-early-gate-00-04-dark-6x.png`
- Manifest: `manifest.csv`
- Mechanical QA summary: `qa/early-gate-00-04-summary.json`
- Per-frame QA: `qa/frame-00-attempt-01-qa.json` through `qa/frame-04-attempt-01-qa.json`

## Rejects And Retries

No ImageGen frame retries were needed for `00-04`.

The repo strict `prepareframe` path rejected frame `00` because the
magenta-key alpha contained an internal transparent gap. This run preserves
that negative space with a no-hole-fill fit so the feet/contact shape is not
filled into a false underside arc. The generated green body intentionally
produces high `greenish_pixels` counts; for this green bird run, visual review
on checker/light/dark contacts is the acceptance gate for contact points.

## Stop Condition

Generation stopped after frame `04`. No catalog, runtime, docs page, workflow,
release, Git, or other animal directory edits were made by this lane.
