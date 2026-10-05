# ChatGPT Image Generation Trial - 2026-06-25

## Scope

This is a review-only ChatGPT image-generation lane for four AnimalsDesktop set00 motion candidates:

- `hamster_golden_syrian`
- `macaroni_mouse_tan`
- `sugar_glider_gray`
- `rabbit_chestnut_agouti`

No files were written to `accepted-frames`, catalog files, runtime code, GitHub Pages, or release workflows. Existing files were not overwritten.

Only public-safe prompt content was used: species descriptions, sprite constraints, motion intents, and avoid lists. No secrets, tokens, raw logs, private GitHub data, or local absolute paths were included in prompts.

## Prompt And Output Index

| Variant | Prompt | Generated sheet | 96x64 cutouts | Contact sheets |
| --- | --- | --- | --- | --- |
| `hamster_golden_syrian` | `prompts/hamster_golden_syrian_set00_8cell_prompt.txt` | `generated/hamster_golden_syrian_set00_8cell_sheet.png` | `cutouts/hamster_golden_syrian/frame-00-96x64.png` through `frame-07-96x64.png` | `contact-sheets/hamster_golden_syrian_checker_contact.png`, `contact-sheets/hamster_golden_syrian_transparent_contact.png` |
| `macaroni_mouse_tan` | `prompts/macaroni_mouse_tan_set00_8cell_prompt.txt` | `generated/macaroni_mouse_tan_set00_8cell_sheet.png` | `cutouts/macaroni_mouse_tan/frame-00-96x64.png` through `frame-07-96x64.png` | `contact-sheets/macaroni_mouse_tan_checker_contact.png`, `contact-sheets/macaroni_mouse_tan_transparent_contact.png` |
| `sugar_glider_gray` | `prompts/sugar_glider_gray_set00_8cell_prompt.txt` | `generated/sugar_glider_gray_set00_8cell_sheet.png` | `cutouts/sugar_glider_gray/frame-00-96x64.png` through `frame-07-96x64.png` | `contact-sheets/sugar_glider_gray_checker_contact.png`, `contact-sheets/sugar_glider_gray_transparent_contact.png` |
| `rabbit_chestnut_agouti` | `prompts/rabbit_chestnut_agouti_set00_8cell_prompt.txt` | `generated/rabbit_chestnut_agouti_set00_8cell_sheet.png` | `cutouts/rabbit_chestnut_agouti/frame-00-96x64.png` through `frame-07-96x64.png` | `contact-sheets/rabbit_chestnut_agouti_checker_contact.png`, `contact-sheets/rabbit_chestnut_agouti_transparent_contact.png` |

Combined local review sheet:

- `contact-sheets/all_species_checker_contact.png`

## Local Processing

Generated sheets were copied into this folder, then split as uniform 4 columns by 2 rows. Each cell was keyed against broad chroma-green background criteria, trimmed to the detected subject bbox, scaled into a 96x64 transparent review candidate, and placed on a shared lower baseline.

Mechanical results:

- 32 total 96x64 PNG review candidates were written under `cutouts/`.
- All 32 cutouts are non-empty RGBA PNGs.
- All 32 cutouts have zero alpha pixels touching the outer 96x64 canvas edge.
- These are review candidates only. The sheet generation and automatic keying are not sufficient for promotion to `accepted-frames`.

## Visual Notes

| Variant | Result | Notes |
| --- | --- | --- |
| `hamster_golden_syrian` | promising review candidate | Strong species read: compact golden body, round cheek, very short/no tail. The 8 cells are visually consistent, but the motion delta is mild and several poses read nearly identical after 96x64 scaling. |
| `macaroni_mouse_tan` | promising review candidate | Thick pale tail is preserved and avoids the hamster no-tail problem. The contact-sheet scale is usable, but later prompts must keep the tail thick without making the body hamster-like. |
| `sugar_glider_gray` | useful but size-risky | Membrane, mask, long tail, and grounded stance are readable. The animal is very horizontal; at 96x64 it spends most of the width and loses fine pose differences. Next prompt should shorten/curve the tail or lift it slightly while keeping a sugar-glider silhouette. |
| `rabbit_chestnut_agouti` | promising review candidate | Long ears stay inside the frame and the chestnut agouti body reads well. Motion differences are subtle; next pass should exaggerate front-paw and hind-leg phase slightly without creating an oversized leap. |

## What Worked

- A 4x2 sheet request is a good small trial size. It produced coherent source families without jumping to a fragile 62-cell sheet.
- The shared prompt structure kept camera, scale, facing direction, and style more consistent than asking for isolated frames.
- Species-specific negative prompts helped: short/no tail for hamster, thick pale tail for macaroni mouse, membrane plus grounded stance for sugar glider, and protected ear margin for rabbit.
- The chroma-green background was good enough for local review extraction, even though it was not perfectly uniform.

## Failure Causes / Risks

- The generated background was not perfectly flat #00ff00. Local keying worked for review, but this would need stricter cleanup before any acceptance path.
- Motion deltas are too small for a production set00 walk cycle. The trial is closer to source-family consistency testing than final animation.
- The model did not reliably understand invisible cell boundaries; uniform 4x2 cutting worked here, but future runs should keep generous padding and avoid more than 8 or 16 cells at once.
- The sugar-glider tail makes the runtime silhouette too wide. The next prompt should explicitly request a compact tail arc inside the frame.
- None of these frames should be promoted directly. They need prepareframe/auditframes style validation and visual QA against light/dark/checker backgrounds after any real adoption pass.

## Recommended Next Prompt Direction

Run another 8-cell pass per species, still review-only, with these changes:

- Ask for visibly larger but controlled pose deltas: front paw raised, forefoot contact, rear foot swing, and recovery.
- Keep the same family anchor words from the best current sheet.
- Ask for more empty green padding and "no off-green background variation" explicitly.
- For `sugar_glider_gray`, add: "tail curves gently upward and stays compact inside each cell; do not draw a fully stretched horizontal tail."
- For `rabbit_chestnut_agouti`, add: "keep ears fully inside the cell with extra top padding; make leg phases clearer but no airborne leap."

## Status

Completion condition for this lane is met:

- `report.md` exists.
- All four requested variants have at least one ChatGPT image generation trial result.
- Prompts, generated sheets, extracted 96x64 candidates, contact sheets, and adoption notes are saved under this dedicated folder.
- The parent thread can decide whether to request another trial pass or adapt the stronger prompt patterns into the main ImageGen-only lane.
