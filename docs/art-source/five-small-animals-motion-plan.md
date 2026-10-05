# Five Small Animals Motion Plan

Last updated: 2026-06-25

## Purpose

This plan keeps the user-requested animal set independent from the current chinchilla-only release path while preserving the DeguDesktop-quality bar. The target families are chinchilla, hamster, macaroni mouse / fat-tailed gerbil, momonga / sugar glider, and rabbit.

The near-term target is `set00` for one representative variant per family. A full release-ready family still needs `set00` through `set09`. Current user direction is ImageGen-only: do not use ChatGPT, Gemini, Flow, Nano Banana, or other external AI review/generation paths in this task.

## Representative Variants

| Family | Variant ID | Style anchor decision |
| --- | --- | --- |
| Chinchilla | `chinchilla_standard_gray` | Use the accepted chinchilla source family as the quality reference, but regenerate missing frames individually. |
| Hamster | `hamster_golden_syrian` | Keep round cheeks, short/no visible tail, tiny paws, and wheel-safe small-rodent poses. |
| Macaroni mouse / fat-tailed gerbil | `macaroni_mouse_tan` | Keep a compact body and thick tail; avoid normal fancy-mouse tail silhouettes. |
| Momonga / sugar glider | `sugar_glider_gray` | Regenerate source truth before motion production; current generated source is not enough. |
| Rabbit | `rabbit_chestnut_agouti` | Keep long ears, strong hindquarters, tiny tail, and hop-specific timing. |

## Production Rules

- One accepted source frame means one standalone transparent 96x64 PNG.
- Every frame must contain one complete right-facing side-view animal.
- Keep stable camera distance, scale, body volume, contact baseline, and anatomy across all 62 frames in a family.
- No text, labels, numbers, borders, dividers, scenery, props, costumes, shadows, floor bands, checker backgrounds, or multiple animals.
- No degu recolors and no copied silhouette across species.
- Do not promote any sheet, grid, contact sheet, or checker-background output into `accepted-frames`.
- Use chroma-green only as a candidate input fallback; it must be removed by `cmd/prepareframe -background chroma-green` before acceptance.

## Set00 Frame Contract

| Frame | Action | Prompt intent |
| ---: | --- | --- |
| 00 | idle | Calm neutral side-view stance, balanced baseline, all anatomy visible. |
| 01 | idle | Tiny breathing change, same stance and scale. |
| 02 | idle | Subtle ear, whisker, or head adjustment. |
| 03 | idle | Calm weight shift without changing body volume. |
| 04 | walk | Walk cycle start, front foot begins forward motion. |
| 05 | walk | Walk cycle, forefoot contacts shared baseline. |
| 06 | walk | Walk cycle, rear foot begins forward motion. |
| 07 | walk | Walk cycle, body moves softly but stays same size. |
| 08 | walk | Walk cycle, opposite forefoot motion. |
| 09 | walk | Walk cycle, opposite forefoot contacts baseline. |
| 10 | walk | Walk cycle, rear foot follows. |
| 11 | walk | Walk cycle return toward neutral. |
| 12 | scurry | Faster low travel pose, compact body. |
| 13 | scurry | Fast travel with stronger foot or body offset. |
| 14 | scurry | Fast travel mid-cycle, no jump unless species-specific. |
| 15 | scurry | Fast travel, tail or ear follows motion. |
| 16 | scurry | Fast travel alternate foot phase. |
| 17 | scurry | Fast travel contact phase. |
| 18 | scurry | Fast travel recovery phase. |
| 19 | scurry | Fast travel return phase. |
| 20 | forage | Sniff or nibble prep, head lowers slightly. |
| 21 | forage | Paws or muzzle move toward food-search position. |
| 22 | forage | Compact nibble pose, no visible food prop required. |
| 23 | forage | Cheek or muzzle motion while baseline stays stable. |
| 24 | forage | Nose-down search pose, no floor or shadow. |
| 25 | forage | Return from forage toward neutral. |
| 26 | species action | Chinchilla soft hop, hamster paw lift, macaroni quick crouch, sugar-glider membrane stretch, rabbit hop prep. |
| 27 | species action | Continue the species-specific action with stable scale. |
| 28 | species action | Peak of the action, no airborne scene unless rabbit needs tiny hop lift. |
| 29 | species action | Recover from action. |
| 30 | species action | Secondary small action variant. |
| 31 | species action | Return to neutral posture. |
| 32 | turn | Side-view turn preparation, head/shoulder rotates slightly. |
| 33 | turn | Early three-quarter turn, still one complete animal. |
| 34 | turn | Mid-turn body shift, no mirrored duplicate. |
| 35 | turn | Turn continuation, stable feet and volume. |
| 36 | turn | Return from turn, opposite phase. |
| 37 | turn | Head and body align back to side view. |
| 38 | turn | Final turn recovery. |
| 39 | turn | Neutral-ready side-view posture. |
| 40 | eat | Chew or eat pose, paws close to mouth where relevant. |
| 41 | eat | Chew continuation with small head movement. |
| 42 | eat | Compact eating posture, no food prop required. |
| 43 | eat | Return from eating pose. |
| 44 | ground check | Light paw scratch or substrate-check pose. |
| 45 | ground check | Ground-check continuation, no floor or dirt drawn. |
| 46 | ground check | Alternate paw or body adjustment. |
| 47 | ground check | Return from ground-check pose. |
| 48 | stand/rest | Alert seated, crouched, or rest posture matching species. |
| 49 | stand/rest | Slight taller/lower rest phase. |
| 50 | stand/rest | Calm hold, stable silhouette. |
| 51 | stand/rest | Return from rest/stand. |
| 52 | groom | Face groom, whisker clean, paw clean, or preen action. |
| 53 | groom | Groom continuation with coherent face. |
| 54 | groom | Groom alternate paw/head phase. |
| 55 | groom | Return from grooming. |
| 56 | reaction | Wheel-safe small-rodent burst or species-safe alert reaction. |
| 57 | reaction | Reaction continuation; non-wheel species must avoid wheel-only body pose. |
| 58 | reaction | Brief alert/stretch phase. |
| 59 | reaction | Reaction recovery. |
| 60 | reaction | Secondary alert or motion burst. |
| 61 | reaction | Return to neutral. |

## Species Prompt Addenda

### Chinchilla

Dense fluffy standard gray body, very large rounded ears, short muzzle, thick fluffy tail, tiny feet, visible whiskers. Movement should feel soft and heavy, not degu-like.

### Hamster

Golden Syrian hamster, round cheek-forward compact body, rounded ears, blunt muzzle, tiny paws, very short or visually absent tail. Keep wheel-safe poses readable without drawing a wheel.

### Macaroni Mouse / Fat-Tailed Gerbil

Compact tan fat-tailed gerbil silhouette, large dark eyes, short rounded body, small hands and feet, thick pale tail. Avoid thin fancy-mouse tails and avoid making it a hamster.

### Momonga / Sugar Glider

Gray sugar glider / Japanese flying squirrel inspired silhouette, large dark eyes, small rounded ears, visible patagium membrane when stretching, long tail, low skitter stance. Keep grounded side-view desktop-pet poses; no full flying scene.

### Rabbit

Chestnut agouti rabbit, long ears, compact body, strong hind legs, small tail, visible whiskers. Hop frames should show rabbit timing without oversized leaps.

## Candidate Prompt Template

```text
Create one single right-facing side-view [FAMILY] desktop pet sprite pose for a 96x64 runtime frame.
Pose intent: [FRAME INTENT].
One complete animal only. Keep the same camera distance, body scale, proportions, contact baseline, facing direction, lighting, outline softness, and source family across the whole 62-frame set.
[SPECIES ADDENDUM]
Transparent background with real alpha channel if supported. If transparency is not supported, use a perfectly flat pure #00ff00 chroma background with no gradient, no texture, no shadow, and no floor.
No text, label, number, border, divider, scenery, prop, costume, multiple animals, cropped anatomy, disconnected stray pixels, checkerboard, floor band, cast shadow, white background, gray background, or human-like pose.
Face quality guard: coherent eye, nose, mouth, muzzle, ears, whiskers where relevant; no smeared face, giant eye, black mask, mouth bar, or broken paws.
```

## Verification

When Go is available:

```powershell
go run ./cmd/prepareframe -src <candidate.png> -out <family>/motion-source/prepared-candidates/set00/frame-00.png
go run ./cmd/auditframes -frames-dir <family>/motion-source/accepted-frames/set00
go run ./cmd/assemblemotion -frames-dir <family>/motion-source/accepted-frames/set00 -out <family>/motion-source/sheets/<variant>-source-set00.png
go run ./cmd/importanimals
go run ./cmd/validatemotion -variant <variant>
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
git diff --check
```
