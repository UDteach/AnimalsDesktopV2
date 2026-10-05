# Child thread prompt: leucistic sugar glider color pilot 00-04

You are a child production lane for AnimalsDesktop. Work only on the variant
below and only inside the output directory listed here.

## Scope

- Variant: `leucistic_sugar_glider`
- Japanese label: リューシスティックフクロモモンガ
- Route: color-variant pilot from accepted `sugar_glider_gray` morphology.
- Output dir:
  `docs/art-source/one-frame-method-fullrun-20260627/leucistic-sugar-glider-color-pilot-00-04/`
- Frame range: `00-04` only.
- Stop gate: stop after `00-04` and report parent-review evidence. Do not
  continue to frame `05` until the parent approves.

## Must Read First

1. `AGENTS.md`
2. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/SKILL.md`
3. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/references/animalsdesktop.md`
4. Existing accepted reference frames:
   - `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/frame-00.png`
   - `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/frame-01.png`
   - `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/frame-02.png`
   - `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/frame-03.png`
   - `docs/art-source/sugar-glider/motion-source/accepted-frames/set00/frame-04.png`

## Generation Direction

Create the leucistic color morph of the sugar glider / フクロモモンガ: white
to warm-cream coat, black eyes, pale pink nose/feet/ears, and very subtle
shading. Preserve the accepted sugar glider body plan: compact small marsupial
glider, large round eyes, rounded ears, side-view body, long fluffy tail,
visible gliding membrane when the pose shows it, small attached feet, and the
same baseline/scale/camera as `sugar_glider_gray`.

This is a white/leucistic フクロモモンガ pattern, not a separate rodent
flying squirrel. Do not add chipmunk stripes, American flying squirrel
proportions, scenery, props, shadows, text, borders, floor marks, duplicate
animals, costumes, or human-like poses.

## Workflow

Use one-frame ImageGen attempts or an equivalent color-variant workflow. Keep
raw output, prompt text, normalized 96x64 transparent candidate frames, QA JSON,
manifest rows, and contact sheets inside the output dir only.

Do not edit catalog, runtime sprites, `docs/index.html`, release files,
workflows, or any other animal lane.

For each accepted candidate, produce:

- `frames/frame-NN.png`
- raw source under `raw/`
- prompt under `prompts/`
- QA evidence under `qa/`
- a contact sheet for frames `00-04`

Run the local frame audit if available. Report:

- frame count
- raw count
- audit summary
- contact sheet path
- any visual risk such as color drift, membrane loss, cropped tail/ears/feet,
  detached feet, or species drift
