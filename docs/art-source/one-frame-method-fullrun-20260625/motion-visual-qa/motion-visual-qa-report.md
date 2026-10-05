# Motion Visual QA - Accepted Set00

## Scope

Reviewed the five `set00` families after promotion into `docs/art-source/<family>/motion-source/accepted-frames/set00/`.

The verdict is for DeguDesktop-compatible motion slots, not for playing all 62 frames as one continuous loop.

The motion contract defines frame ranges as separate behavior slots:

- 00-03 idle
- 04-11 walk / hop-walk
- 12-19 fast movement
- 20-25 sniff / nibble
- 26-31 species action
- 32-39 turn
- 40-43 chew / eat
- 44-47 ground check
- 48-51 alert rest / stand
- 52-55 face groom
- 56-61 reaction / recover

Visual review sheets:

- `hamster-accepted-motion-groups.png`
- `macaroni-mouse-accepted-motion-groups.png`
- `sugar-glider-accepted-motion-groups.png`
- `rabbit-accepted-motion-groups.png`
- `chinchilla-accepted-motion-groups.png`

Earlier parent-selected review sheets are also preserved:

- `hamster-motion-groups.png`
- `macaroni-mouse-motion-groups.png`
- `sugar-glider-motion-groups.png`
- `rabbit-motion-groups.png`

## Result

| Family | Motion verdict | Notes |
| --- | --- | --- |
| chinchilla | pass with notes | The previous 16-pose draft was rejected for debris and doubled-face/face-like artifacts. The regenerated accepted set has no visible debris or doubled-face defects. Idle, walk, scurry, nibble, turn, eat, ground check, alert/rest, groom, and reaction/recover read as slot-based chinchilla motion. It is not intended as one all-62-frame continuous loop. |
| hamster | pass with notes | Walk, fast movement, sniff, ground check, groom, and reaction read as hamster motion. Turn and chew include front-facing poses, so this should be used as behavior slots rather than one continuous 62-frame loop. |
| macaroni-mouse | pass | The strongest motion continuity of the four. Walk, fast scurry, ground check, and turn read cleanly while preserving the short thick tail identity. Frame 56 is slightly upright but acceptable as reaction. |
| sugar-glider | pass with notes | Low skitter and membrane reaction read clearly. Frames 56-58/60 are wide action poses and should be used as reaction/action slots, not a continuous travel loop. Frame 35 is visually larger but acceptable inside turn. |
| rabbit | pass with notes | Hop and fast hop read well. Frame 58 is a small upright reaction pose and frame 59 is a long low recover pose; both are acceptable as reaction/recover but would jump if played as an all-frames loop. |

## Promotion Decision

The five selected/regenerated sets are promoted to `motion-source/accepted-frames/set00` as source-frame assets for slot-based motion:

- `docs/art-source/chinchilla/motion-source/accepted-frames/set00/`
- `docs/art-source/hamster/motion-source/accepted-frames/set00/`
- `docs/art-source/macaroni-mouse/motion-source/accepted-frames/set00/`
- `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/`
- `docs/art-source/rabbit/motion-source/accepted-frames/set00/`

They should not be described as polished full release-ready families yet because release readiness still requires `set00` through `set09`, runtime import, catalog validation, and final app-side review.

Chinchilla was completed after the four-family promotion because the earlier draft had visible debris and doubled-face defects. The accepted chinchilla set comes from `docs/art-source/one-frame-method-fullrun-20260625/chinchilla/`, with the prior accepted `frame-00` preserved under `docs/art-source/chinchilla/motion-source/accepted-frames/stale/pre-one-frame-fullrun-20260625/`.

## Reproducibility Requirement

For accepted-frame promotion, exact ImageGen pixel determinism is not required. Required reproducibility is operational:

- each final frame has a stable frame number and accepted PNG path,
- each one-frame-generated candidate has raw/alpha/QA/prompt/manifest provenance in the fullrun or rescue folders,
- parent selection overrides are recorded in `parent-selected/selection-overrides.json`,
- a failed frame can be regenerated or replaced by frame number without rebuilding the other 61 frames.

This is the practical reproducibility level needed for sprite asset maintenance.
