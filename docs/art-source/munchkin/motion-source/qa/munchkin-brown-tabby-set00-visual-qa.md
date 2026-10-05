# Munchkin Brown Tabby Set00 Visual QA

Date: 2026-06-28

Source run:
`/Users/kyota/.codex/worktrees/1fe1/AnimalDesktop/docs/art-source/one-frame-method-fullrun-20260628/munchkin-brown-tabby-set00-oneframe-62/`

Parent checks:

- `cmd/auditframes -strict -artifact-warnings -motion-warnings`: `valid=62 missing=0 invalid=0 warnings=37`.
- Most warnings are lower ledge/floor candidates caused by the low, long Munchkin body and broad ground-contact fur. Contact review found connected body/leg pixels rather than a floor, prop, shadow, or detached shelf.
- Motion warnings at frame `20` and frame `44` were reviewed on slot contacts. `20` starts a lower sniff/ground-check band, and `44` is an intentionally low ground-check pose; neither reads as a sudden scale error when viewed with neighboring frames.
- Full `00-61` and tail `45-61` contacts keep a consistent short-legged brown-tabby Munchkin read, attached paws, ringed tail, stable baseline, and coherent body ratio.

Accepted for source promotion as `munchkin_brown_tabby` set00.
