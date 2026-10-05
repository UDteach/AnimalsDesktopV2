# Google Flow External AI Trial - 2026-06-25

## Scope

This was a review-only Google Flow trial lane for AnimalsDesktop set00 source families:

- `chinchilla_standard_gray`
- `hamster_golden_syrian`
- `macaroni_mouse_tan`
- `sugar_glider_gray`
- `rabbit_chestnut_agouti`

The lane did not touch catalog/runtime code and did not write to `accepted-frames`.

## External Transmission Boundary

Allowed content was limited to public sprite-generation requirements and sanitized prompts. No secrets, tokens, local absolute paths, private repository information, raw logs, or raw browser page-state dumps were intentionally sent or saved.

No reference image was uploaded. This avoided an upload confirmation path and kept the first trial text-only.

## Google Flow State Reached

- Official entry point used: `https://labs.google/fx/tools/flow`
- Flow project creation reached successfully.
- UI locale: Japanese.
- Account state visible in UI: already signed in with a `PRO` badge.
- Model/settings visible in UI: image tab, `Nano Banana 2`, `16:9`, `x2`.
- The model menu displayed that the current generation would use `0 credits`.

## Attempted Trial

Primary target: `hamster_golden_syrian`, 8-cell motion study.

Prepared packet:

- `source-packets/hamster_golden_syrian_8cell_prompt.md`

Prompt content was sanitized and public. A long prompt and then shorter text-only variants were attempted in the Flow prompt box. No local paths or private repository details were included.

## Blocker

Generation did not start.

The Flow prompt box is a Slate rich-text editor. Browser automation could make text appear visually in the editor, but Flow did not treat it as valid prompt content. The create arrow stayed disabled (`aria-disabled=true`) across these attempts:

- clipboard paste of the full sanitized prompt
- clipboard paste of a shorter one-paragraph prompt
- direct typed short prompt after focusing the editor
- DOM-node click/type after page reload
- coordinate click/type after page reload

The browser console showed Slate node-resolution errors after the initial automated editor attempts. I did not save raw logs because the Flow page state can expose account/session data.

## Generated Images

None.

Because the create button never became enabled, no prompt was submitted for generation, no download was possible, and no 96x64 crops or contact sheets were produced.

## Prepared Follow-up Packets

The following manual/next-run packets are ready:

- `source-packets/hamster_golden_syrian_8cell_prompt.md`
- `source-packets/macaroni_mouse_tan_8cell_prompt.md`
- `source-packets/sugar_glider_gray_8cell_prompt.md`
- `source-packets/rabbit_chestnut_agouti_8cell_prompt.md`

Use them one species at a time. Do not run a 62-frame generation as the first Flow test.

## Prepared Color Variant Pilots

The first real-world color expansion pilots are ready as 8-frame anchor-based Flow packets:

- `chinchilla-beige-color-pilot/`
- `hamster-cream-color-pilot/`
- `macaroni-mouse-cream-color-pilot/`
- `rabbit-white-color-pilot/`
- `sugar-glider-white-faced-blonde-color-pilot/`

See `real-variant-plan.md` and `flow-pilot-manifest.tsv` for target colors, source anchors, and QA status. These packets are for Flow input only; no generated output has been accepted.

The macaroni mouse pilot intentionally does not use `frame-00` as a style anchor because the stricter Flow anchor validation found a 1px transparent pinhole in that accepted frame. It uses `04,12,32` instead.

## Recommended Next Flow Test

Start with `hamster-cream-color-pilot/` because the existing hamster read is already strong and should make quality drift easier to judge. After that, run one pilot per family before expanding a second color for any species.

If entering the prompt manually in the visible browser:

1. Create a new Flow project or use the project created for this trial.
2. Keep image generation on `Nano Banana 2`, `16:9`, and `x2` or reduce to `x1` if output volume should be minimized.
3. Upload the pilot pack images if the UI supports image anchors: `style-anchors.png`, `style-anchors-neutral.png` if transparency renders poorly, `grid-seed.png`, and `grid-guide.png`.
4. Paste only the sanitized prompt text from the target pilot `prompt.txt`.
5. Confirm the UI still reports `0 credits` or otherwise stop before generation.
6. Generate only the 8-cell sheet.
7. Save the original output under `raw/`.
8. Crop review-only candidates under `candidates/<pilot-id>/`.
9. Build a contact sheet under `contact-sheets/`.

## Adoption Notes

Expected success signs:

- one animal per cell
- stable right-facing side view
- stable camera distance and body scale
- shared contact baseline
- no text, numbers, borders, scenery, floor, or shadows
- alpha or flat pure green background

Hard reject signs:

- labels or grid dividers
- multiple animals in a cell
- cropped ears, feet, whiskers, tail, or body
- inconsistent scale or baseline
- toy/plush look, photo-real look, or degu-like silhouette
- hamster losing its short/no-tail shape
- macaroni mouse tail becoming thin or body becoming hamster-like
- sugar glider membrane becoming bat wings or a flying scene
- rabbit ears cropped or hop height requiring a taller canvas

## Result

This run produced reusable, sanitized Flow packets and confirmed that the logged-in Flow UI is reachable, but browser automation could not safely submit the prompt because the editor state did not accept automated input. The lane should be continued with manual prompt entry or a different browser-control surface before any image/crop/contact-sheet review can happen.
