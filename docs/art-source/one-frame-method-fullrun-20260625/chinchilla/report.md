# Chinchilla One-Frame Fullrun - 2026-06-25

## Outcome

Chinchilla `set00` was regenerated with the one-frame method and promoted into:

- `docs/art-source/chinchilla/motion-source/accepted-frames/set00/`

The previous 16-pose draft remains reference-only and is not used for accepted promotion.

## Rejection Of Existing Draft

The earlier draft under `docs/art-source/chinchilla/motion-source/frames/` and `docs/art-source/chinchilla/motion-source/sheets/chinchilla-standard-gray-source-set00-draft.png` was rejected for accepted-frame use.

Reasons:

- visual debris / extra fragments were visible in the source poses,
- some poses read as doubled-face or face-like duplicate anatomy,
- the 16-pose repeated draft did not meet the current one-frame-per-PNG provenance bar.

## QA Evidence

- `manifest.csv`: 62 accepted rows with chosen attempt, raw path, prepared path, QA report, and accepted path.
- `auditframes-report.json`: candidate fullrun frames pass with `valid=62`, `missing=0`, `invalid=0`, `warnings=87`.
- `docs/art-source/chinchilla/motion-source/accepted-frames/set00-auditframes-report.json`: promoted accepted frames pass with `valid=62`, `missing=0`, `invalid=0`, `warnings=87`.
- `docs/art-source/chinchilla/motion-source/accepted-frames/set00-assemblemotion-report.json`: accepted standalone frames assemble into `docs/art-source/chinchilla/motion-source/sheets/chinchilla-standard-gray-source-set00.png`.
- `../motion-visual-qa/chinchilla-accepted-motion-groups.png`: accepted-frame motion-slot visual QA sheet.

`auditframes` artifact warnings are not hard failures for this lane; they are visual review prompts. The accepted motion sheet was reviewed visually after the user pointed out the prior draft's debris and doubled-face defects.

## Motion Visual Verdict

Pass with notes. The accepted `set00` works as DeguDesktop-compatible slot-based motion:

- `00-03` idle: calm breathing / tiny weight changes.
- `04-11` walk: stable compact walk.
- `12-19` fast scurry: low quick movement with a trailing tail plume.
- `20-25` forage / nibble: head lowers and paws move near mouth.
- `26-31` species action: small paw lift / hop-safe alert action.
- `32-39` turn: side to front and back to side without doubled-face artifacts.
- `40-43` eat: paws near mouth and recovery.
- `44-47` ground check: low sniff / paw check and recovery.
- `48-51` alert/rest: head lift and compact rest.
- `52-55` face groom: paw-to-face motion and return.
- `56-61` reaction/recover: small surprise, recoil, and return to idle.

This is not a single polished all-62-frame loop. It is accepted as slot-based source-frame motion.
