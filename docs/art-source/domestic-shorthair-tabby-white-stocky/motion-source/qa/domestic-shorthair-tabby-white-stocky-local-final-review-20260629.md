# Domestic Shorthair Tabby White Stocky Local Final Review - 2026-06-29

Scope: source-art accepted set only. Runtime catalog, runtime sprite sheets,
GitHub Pages current cards, release artifacts, tags, and public versions were
not changed.

Artifacts reviewed:

- Run directory:
  `docs/art-source/one-frame-method-fullrun-20260629/domestic-shorthair-tabby-white-stocky-set00-oneframe-62/`
- Accepted frames:
  `docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/accepted-frames/set00/`
- Full contact:
  `docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/contacts/domestic-shorthair-tabby-white-stocky-set00-full-contact.png`
- Final dark contact:
  `docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/contacts/domestic-shorthair-tabby-white-stocky-set00-final-dark-3x.png`
- Source sheet:
  `docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/sheets/domestic-shorthair-tabby-white-stocky-source-set00.png`

Visual verdict: pass.

The full contact keeps the supplied photo-cat read across all 62 frames: a
stocky adult domestic shorthair body, warm dark brown-gray mackerel tabby saddle
over the back/head/tail, white lower face/chest/belly/legs/paws, and dark tail
tip. Low frames `40-47` read as chew and ground-check poses rather than scale
loss. Reaction frames `57` and `59` are taller, but they keep the same baseline,
width, color blocking, and contact points and read as head-lift reaction poses
instead of camera zoom.

Frame `39` had a usable structure but a washed-out tabby saddle. It was repaired
with a dark-tabby visible-pixel color transfer from frames `37-38`; the original
and report are preserved under the run QA folder, and the accepted QA report is
copied to `qa/color-correction-frame39-20260629-report.json`.

Mechanical checks:

- `oneframe_run.py review --variant domestic_shorthair_tabby_white_stocky`:
  `frame_count=62`, `issues=[]`
- Source-run `auditframes -strict -artifact-warnings`: `valid=62 missing=0
  invalid=0 warnings=25`
- Accepted-source `auditframes -strict -artifact-warnings`: `valid=62 missing=0
  invalid=0 warnings=25`
- `assemblemotion`: produced
  `sheets/domestic-shorthair-tabby-white-stocky-source-set00.png`
- `validatemotion -variant domestic_shorthair_tabby_white_stocky
  -require-accepted`: deferred because the runtime catalog does not yet know the
  variant.
