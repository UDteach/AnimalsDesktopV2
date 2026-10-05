# One-Frame Method Report - Codex Built-in ImageGen - 2026-06-25

## Scope

This pass tested the parent-thread hypothesis:

> One-frame generation may be more robust than sheet generation for pose control, missing parts, and retry-until-pass reproducibility.

Execution rules followed:

- Used only Codex built-in `image_gen`.
- Wrote only under `docs/art-source/external-ai-trials/codex-imagegen-20260625/`.
- Did not write to `accepted-frames`, `accepted-candidates`, runtime code, catalog code, GitHub Pages, or release artifacts.
- Preserved raw, alpha, normalized 96x64, contact, QA, prompt, manifest, and visual-review artifacts.
- Did not overwrite existing trial outputs.

## Inputs Reviewed

| Input | Path | Use |
| --- | --- | --- |
| Four-animal sheet overview | `docs/art-source/four-animals-imagegen-only-62/four-animals-set00-overview.png` | Baseline sheet-extracted 62-frame comparison. |
| Parent one-frame pilot directory | `docs/art-source/one-frame-method-trials-20260625/` | Existing one-frame pilot for hamster, macaroni mouse, sugar glider, and rabbit. |
| Parent pilot comparison sheet | `docs/art-source/one-frame-method-trials-20260625/baseline-vs-one-frame-pilot.png` | Visual comparison of sheet baseline vs one-frame pilots. |
| Previous Codex-only trial report | `docs/art-source/external-ai-trials/codex-imagegen-20260625/report.md` | Earlier finding that 62-cell sheets fail equal-cell slicing. |

## New Trial Set

Target: `sugar_glider_gray`, because the sheet baseline has good general consistency but weak high-risk frames around low scurry, turn, grooming, and membrane/reaction poses.

Generated frames:

| Frame | Pose | Visual decision |
| ---: | --- | --- |
| 16 | `fast-low-scurry-alternate` | `visual-pass` |
| 20 | `forage-sniff-prep` | `visual-pass` |
| 26 | `grounded-membrane-stretch-prep` | `visual-pass` |
| 35 | `turn-continuation-three-quarter` | `visual-pass-with-style-watch` |
| 48 | `alert-seated-rest` | `visual-pass-with-scale-watch` |
| 52 | `face-groom` | `visual-pass` |
| 56 | `small-alert-reaction-burst` | `visual-pass` |
| 60 | `grounded-membrane-spread-reaction` | `visual-pass` |

Primary contact sheet:

- `docs/art-source/external-ai-trials/codex-imagegen-20260625/one-frame-method/sugar_glider_gray/set00/contacts/sugar_glider_gray-one-frame-8pose-contact-grid.png`

Visual decision artifacts:

- `docs/art-source/external-ai-trials/codex-imagegen-20260625/one-frame-method/sugar_glider_gray/set00/visual-review.csv`
- `docs/art-source/external-ai-trials/codex-imagegen-20260625/one-frame-method/sugar_glider_gray/set00/visual-review.json`

Manifest:

- `docs/art-source/external-ai-trials/codex-imagegen-20260625/one-frame-method/sugar_glider_gray/set00/manifest.csv`

## Saved Artifact Layout

For each generated frame:

- raw ImageGen copy: `one-frame-method/sugar_glider_gray/set00/raw/frame-XX-attempt-01-raw.png`
- alpha output: `one-frame-method/sugar_glider_gray/set00/alpha/frame-XX-attempt-01-alpha.png`
- normalized 96x64 PNG: `one-frame-method/sugar_glider_gray/set00/frames/frame-XX-attempt-01.png`
- contact preview: `one-frame-method/sugar_glider_gray/set00/contacts/frame-XX-attempt-01-contact.png`
- QA JSON: `one-frame-method/sugar_glider_gray/set00/qa/frame-XX-attempt-01-qa.json`
- prompt text: `one-frame-method/sugar_glider_gray/set00/prompts/frame-XX-attempt-01-prompt.txt`

Built-in ImageGen source files remain in:

- `/Users/kyota/.codex/generated_images/019efbc1-b6f1-7b32-bd44-e720ce478780/`

## Mechanical QA

All eight one-frame attempts:

- empty alpha: `0 / 8`
- edge-alpha contact: `0 / 8`
- greenish pixels after chroma removal: `0..4`
- normalized frame size: `96x64`
- normalized source width: `88px` for every frame in this trial

Per-frame mechanical summary:

| Frame | Empty | Edge alpha | Greenish px | Nonzero alpha | Frame bbox |
| ---: | --- | ---: | ---: | ---: | --- |
| 16 | no | 0 | 1 | 1507 | `[4, 32, 91, 59]` |
| 20 | no | 0 | 0 | 1745 | `[4, 27, 91, 59]` |
| 26 | no | 0 | 1 | 1983 | `[4, 14, 91, 59]` |
| 35 | no | 0 | 3 | 1996 | `[4, 19, 92, 59]` |
| 48 | no | 0 | 3 | 2241 | `[4, 11, 91, 59]` |
| 52 | no | 0 | 4 | 2056 | `[4, 13, 91, 59]` |
| 56 | no | 0 | 2 | 1759 | `[4, 22, 91, 59]` |
| 60 | no | 0 | 0 | 2014 | `[4, 14, 92, 59]` |

## Comparison With Sheet Extraction

### Missing parts and cropped poses

One-frame generation is materially stronger.

The previous 62-sheet attempts in this thread failed because equal-grid slicing produced row-wrapped or cropped fragments. The parent overview also shows that some sheet frames can become pose families rather than controlled single-frame repairs. In the new eight-frame one-frame trial, all eight had one complete animal, no cropped tail/feet/ears, and no row fragments.

### Pose control

One-frame generation is stronger for high-risk poses.

Frames 26 and 60 produced clear grounded membrane poses, and frame 20 produced a clear nose-down forage pose. These are exactly the cases where sheet-level prompts tend to mix poses, front-facing spreads, or fragments.

Weakness: frame 56's reaction delta is visually safe but subtle. One-frame prompts still need sharper pose verbs when the target is a small animation change rather than a silhouette change.

### Style consistency

One-frame generation is usable, but this is the main risk.

The eight sugar glider outputs are internally fairly consistent because the prompt fixed species, camera, and source-family language. However, compared with the sheet baseline and parent pilot, one-frame outputs tend to be slightly sharper, more detailed, and sometimes larger/taller after 96x64 normalization. Frames 35 and 48 should be treated as `style-watch` examples.

Practical implication: one-frame repair should use a fixed normalization profile, shared prompt template, and a visual comparison contact against accepted neighbors before adoption.

### Re-generation and reproducibility

One-frame generation is operationally more reproducible.

The important reproducibility gain is not deterministic pixels; it is isolated retry scope. If frame 35 drifts too front-facing, only frame 35 needs another attempt. A 62-sheet retry can perturb every frame and creates a new slicing/selection problem.

In this trial, attempt 01 was enough for all eight frames to become review candidates. That supports the parent hypothesis that retry-until-pass is a better operational method than trying to get all poses right in one sheet.

### Throughput

One-frame generation is slower per frame but cheaper per accepted result for high-risk repairs.

For routine idle or walk sequences, small 8/16-cell sheets may still be useful as exploration. For acceptance-bound frames, especially species-action, turn, reaction, and tall/long-anatomy poses, one-frame generation is more reliable.

## Recommended Parent-Thread Operation

Adopt a hybrid workflow:

1. Use sheet or 8/16-cell batches only for exploration and style discovery.
2. Use one-frame `retry-until-pass` for any frame that may enter `accepted-frames`.
3. Keep a per-frame attempt ledger with raw, alpha, normalized, contact, QA, prompt, and visual decision.
4. Do not promote by mechanical QA alone; require visual contact review against neighboring frames.
5. Require these pass gates before parent integration:
   - one complete animal only
   - no cropped ears/feet/tail/membrane/whiskers
   - no detached fragments, duplicate body parts, props, floor, text, or effects
   - `96x64` normalized PNG
   - non-empty alpha
   - no edge-alpha contact
   - minimal green residue after chroma removal
   - style comparison against the target family contact sheet
6. Retry a frame independently until it passes, then freeze that frame and move to the next failed frame.

Recommended first production use:

- Use one-frame method to repair `sugar_glider_gray` high-risk frames `26`, `35`, `48`, `52`, `56`, and `60`.
- Use the new `frame-60-attempt-01` as the strongest proof that one-frame generation can outperform sheet extraction for membrane/reaction poses.

## Conclusion

Hypothesis result: supported.

One-frame generation is not automatically final-art quality, but it is clearly better than sheet extraction for high-risk pose completeness and retry control. The parent thread should treat one-frame retry-until-pass as the default acceptance-bound method, with sheets reserved for rough exploration.
