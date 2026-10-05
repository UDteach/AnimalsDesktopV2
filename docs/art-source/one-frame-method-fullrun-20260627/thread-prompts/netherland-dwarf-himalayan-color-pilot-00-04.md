# Child thread prompt: Netherland Dwarf Himalayan color pilot 00-04

You are a child production lane for AnimalsDesktop. Work only on the variant
below and only inside the output directory listed here.

## Scope

- Variant: `netherland_dwarf_himalayan`
- Japanese label: ネザーランドドワーフ（ヒマラヤン）
- Route: color-variant pilot from accepted `netherland_dwarf_chestnut`
  morphology.
- Output dir:
  `docs/art-source/one-frame-method-fullrun-20260627/netherland-dwarf-himalayan-color-pilot-00-04/`
- Frame range: `00-04` only.
- Stop gate: stop after `00-04` and report parent-review evidence. Do not
  continue to frame `05` until the parent approves.

## Must Read First

1. `AGENTS.md`
2. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/SKILL.md`
3. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/references/animalsdesktop.md`
4. Existing accepted reference frames:
   - `docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/frame-00.png`
   - `docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/frame-01.png`
   - `docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/frame-02.png`
   - `docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/frame-03.png`
   - `docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/frame-04.png`
5. Himalayan point color reference, but not body shape reference:
   - `docs/art-source/himalayan-rabbit/motion-source/accepted-frames/set00/frame-00.png`

## Generation Direction

Create a compact Netherland Dwarf rabbit with Himalayan point coloration:
cream-white body with darker ears, nose, feet, and tail. Preserve the accepted
Netherland Dwarf morphology: very compact body, short upright ears, round head,
small feet, short tail, same baseline/scale/camera as
`netherland_dwarf_chestnut`.

This is not the existing long-eared Himalayan rabbit. Do not lengthen the ears,
stretch the body, add lop ears, add a second animal, add scenery, props,
shadows, text, borders, costumes, or human-like poses.

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
- any visual risk such as long-ear drift, losing point coloration, cropped
  ears/feet/tail, detached feet, or shadow/debris artifacts
