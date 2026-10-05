# Child thread prompt: African dormouse set00 00-04

You are a child production lane for AnimalsDesktop. Work only on the variant
below and only inside the output directory listed here.

## Scope

- Variant: `african_dormouse`
- Japanese label: アフリカヤマネ
- Route: new species/source-family early gate, basic wild color.
- Output dir:
  `docs/art-source/one-frame-method-fullrun-20260627/african-dormouse-set00-oneframe-62/`
- Frame range: `00-04` only.
- Stop gate: stop after `00-04` and report parent-review evidence. Do not
  continue to frame `05` until the parent approves.

## Must Read First

1. `AGENTS.md`
2. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/SKILL.md`
3. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/references/animalsdesktop.md`
4. Existing small-rodent accepted frames for scale reference only:
   - `docs/art-source/macaroni-mouse/motion-source/accepted-frames/set00/frame-00.png`
   - `docs/art-source/djungarian-hamster/motion-source/accepted-frames/set00/frame-00.png`
   - `docs/art-source/richardsons-ground-squirrel/motion-source/accepted-frames/set00/frame-00.png`

## Generation Direction

Create an African pygmy dormouse / African dormouse in its basic wild color:
very small brown-tan dormouse with rounded ears, large dark eyes, pale
underside, delicate feet, and a long soft bushy tail. It should read as a
dormouse, not a hamster, rat, chipmunk, squirrel, sugar glider, or gerbil.

Avoid leucistic/albino/white/fancy color variants, chipmunk stripes,
sugar-glider membranes, rat-length naked tail, oversized squirrel tail,
prairie-dog posture, scenery, props, shadows, text, borders, duplicate animals,
costumes, and human-like poses. Keep the animal complete, right-facing,
side-view, transparent or chroma-key-ready, and sized for a 96x64 desktop pet
frame.

## Workflow

Use one-frame ImageGen attempts. Keep raw output, prompt text, normalized 96x64
transparent candidate frames, QA JSON, manifest rows, and contact sheets inside
the output dir only.

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
- any visual risk such as species drift, tail scale drift, cropped ears/feet,
  detached toes, or shadow/debris artifacts
