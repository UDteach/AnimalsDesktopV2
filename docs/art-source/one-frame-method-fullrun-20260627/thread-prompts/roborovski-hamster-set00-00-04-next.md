# Next child thread prompt: Roborovski hamster set00 00-04

Do not start this lane until the parent frees one of the current three child
slots. The current thread cap is parent plus three children.

## Scope

- Variant: `roborovski_hamster`
- Japanese label: ロボロフスキーハムスター
- Route: high-priority next-wave new source family.
- Output dir:
  `docs/art-source/one-frame-method-fullrun-20260627/roborovski-hamster-set00-oneframe-62/`
- Frame range: `00-04` only.
- Stop gate: stop after `00-04` and report parent-review evidence.

## Must Read First

1. `AGENTS.md`
2. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/SKILL.md`
3. `/Users/kyota/.codex/skills/animal-motion-imagegen-flow/references/animalsdesktop.md`
4. Existing hamster frames for scale reference only:
   - `docs/art-source/djungarian-hamster/motion-source/accepted-frames/set00/frame-00.png`
   - `docs/art-source/campbell-hamster/motion-source/accepted-frames/set00/frame-00.png`
   - `docs/art-source/hamster/motion-source/accepted-frames/set00/frame-00.png`

## Generation Direction

Create a Roborovski dwarf hamster: very small, round sandy-tan hamster with
white belly and distinct white eyebrow/cheek patches, tiny rounded ears, short
legs, small feet, compact body, and no strong dorsal stripe. It must read as a
Roborovski, not a Djungarian, Campbell, Syrian, mouse, or gerbil.

Avoid long Syrian-hamster body proportions, black-and-white patches, strong
Djungarian dorsal stripe, rat tail, scenery, props, shadows, text, borders,
duplicate animals, costumes, and human-like poses. Keep the animal complete,
right-facing, side-view, transparent or chroma-key-ready, and sized for a
96x64 desktop pet frame.

## Workflow

Use one-frame ImageGen attempts. Keep raw output, prompt text, normalized 96x64
transparent candidate frames, QA JSON, manifest rows, and contact sheets inside
the output dir only.

Do not edit catalog, runtime sprites, `docs/index.html`, release files,
workflows, or any other animal lane.

Report frame count, raw count, audit summary, contact sheet path, and any
visual risk such as losing eyebrow patches, becoming Djungarian/Campbell-like,
cropped ears/feet, or shadow/debris artifacts.
