# Hamster Golden Syrian One-Frame Fullrun - 2026-06-25

Dedicated output root:

`docs/art-source/one-frame-method-fullrun-20260625/hamster/`

## Scope

- Species: `hamster_golden_syrian`
- Method: Codex built-in ImageGen, one frame at a time, retry failed frame numbers only.
- Background: perfectly flat `#00ff00` chroma key, then local alpha removal.
- Runtime/accepted-frame mutation: none intended.
- Baseline retained: existing sheet-extracted `accepted-candidates` remain untouched.

## Prompt Template

```text
Use case: stylized-concept
Asset type: 2D sprite source frame for AnimalsDesktop taskbar pet
Primary request: Create exactly one standalone sprite-source frame of a golden Syrian hamster for frame <NN>: <POSE>.
Scene/backdrop: perfectly flat solid #00ff00 chroma-key background only.
Subject: one complete golden Syrian hamster, warm golden-tan coat, cream belly, small rounded ears, black bead eyes, tiny paws, short tail if visible, compact soft body.
Style/medium: clean polished 2D game sprite illustration, same camera and scale as a small desktop pet, readable at 96x64 after downscaling.
Composition/framing: full body visible with generous padding, side-view or turn angle exactly as requested, animal centered, feet close to a shared baseline.
Lighting/mood: simple even lighting, no shadow.
Constraints: one animal only; no cropped ears, feet, whiskers, or tail; no props, food pieces, scenery, borders, text, watermark, floor plane, gradients, texture, cast shadow, contact shadow, reflection, costume, human pose, or extra animals. Do not use #00ff00 anywhere on the hamster.
Avoid: degu, chinchilla, mouse, rabbit, squirrel, long tail, oversized ears, realistic photo, 3D render, labels, frame numbers.
```

## Pose Table

| Frame | Group | Pose |
| ---: | --- | --- |
| 00 | idle | neutral standing rest, side view facing right |
| 01 | idle | soft breathing rest, side view facing right |
| 02 | idle | tiny weight shift, side view facing right |
| 03 | idle | settled crouch rest, side view facing right |
| 04 | walk | walk cycle contact, front paw forward, facing right |
| 05 | walk | walk cycle passing step, facing right |
| 06 | walk | walk cycle rear push, facing right |
| 07 | walk | walk cycle lifted forepaw, facing right |
| 08 | walk | walk cycle low body step, facing right |
| 09 | walk | walk cycle opposite contact, facing right |
| 10 | walk | walk cycle small stride, facing right |
| 11 | walk | walk cycle return to contact, facing right |
| 12 | fast_scurry | fast scurry stretched low body, facing right |
| 13 | fast_scurry | fast scurry compact gathered paws, facing right |
| 14 | fast_scurry | fast scurry long step, facing right |
| 15 | fast_scurry | fast scurry tucked step, facing right |
| 16 | fast_scurry | fast scurry forward lean, facing right |
| 17 | fast_scurry | fast scurry hind push, facing right |
| 18 | fast_scurry | fast scurry quick recovery step, facing right |
| 19 | fast_scurry | fast scurry low return, facing right |
| 20 | sniff_nibble | nose down sniffing ground, facing right |
| 21 | sniff_nibble | nose forward sniff, whiskers visible, facing right |
| 22 | sniff_nibble | tiny nibble with mouth near paws, facing right |
| 23 | sniff_nibble | head lowered nibble, paws planted, facing right |
| 24 | sniff_nibble | sniff pause with raised nose, facing right |
| 25 | sniff_nibble | small chew pause, side view facing right |
| 26 | paw_action | one front paw lifted near chest, facing right |
| 27 | paw_action | both front paws close to mouth, facing right |
| 28 | paw_action | short reaching paw gesture, facing right |
| 29 | paw_action | paw tap on ground, facing right |
| 30 | paw_action | upright tiny paw action, facing right |
| 31 | paw_action | settling after paw action, facing right |
| 32 | turn | begin turning from right-facing side toward viewer |
| 33 | turn | three-quarter front turn, compact body |
| 34 | turn | front-facing pause, both ears visible |
| 35 | turn | three-quarter left turn, compact body |
| 36 | turn | left-facing side view, small tail visible |
| 37 | turn | turning back from left toward viewer |
| 38 | turn | three-quarter right return |
| 39 | turn | right-facing side view restored |
| 40 | chew | sitting low, paws near mouth, chewing |
| 41 | chew | round cheek chew, paws tucked |
| 42 | chew | small bite motion, head slightly down |
| 43 | chew | chew pause, cheeks rounded |
| 44 | ground_check | nose to floor checking ground, facing right |
| 45 | ground_check | front paws low, inspecting ground, facing right |
| 46 | ground_check | head sweeping along ground, facing right |
| 47 | ground_check | nose down pause, facing right |
| 48 | alert_rest | alert rest with head raised, facing right |
| 49 | alert_rest | alert ears forward, still body, facing right |
| 50 | alert_rest | upright attentive rest, facing right |
| 51 | alert_rest | alert crouch ready to move, facing right |
| 52 | face_groom | front paws washing face, compact sitting pose |
| 53 | face_groom | one paw over cheek, compact sitting pose |
| 54 | face_groom | both paws near nose grooming |
| 55 | face_groom | finishing face groom, paws lowering |
| 56 | reaction | startled flinch, compact body, facing right |
| 57 | reaction | quick hop reaction, all paws close to body |
| 58 | reaction | surprised upright tiny pose |
| 59 | reaction | ducking reaction, low body |
| 60 | reaction | recovering from reaction, facing right |
| 61 | reaction | calm return after reaction, facing right |

## Artifact Layout

- Raw built-in ImageGen copies: `raw/frame-NN-attempt-MM-raw.png`
- Alpha converted source: `alpha/frame-NN-attempt-MM-alpha.png`
- Normalized candidate frame: `frames/frame-NN.png`
- Attempt-specific normalized frame: `frames/frame-NN-attempt-MM.png`
- Per-frame contact preview: `contacts/frame-NN-attempt-MM-contact.png`
- 62-frame contact sheet: `contacts/hamster-set00-fullrun-contact.png`
- QA JSON: `qa/frame-NN-attempt-MM-qa.json`
- Manifest: `manifest.csv`

## Progress Log

- 2026-06-25: Created fullrun scaffold and processing helper.
- 2026-06-25: Generated/claimed/processed active candidate frames `00` through `61` locally.
- 2026-06-25: Recorded retry provenance for frame `17` attempt 01 rejected, frame `19` attempt 01 rejected, species-drift retries for frames `12`, `13`, `18`, and `28`, plus later local completion attempts for previously missing frames.
- 2026-06-25: Parent clarified that unclaimed external `generated_images` should be ignored unless explicitly generated and claimed by this fullrun thread.
- 2026-06-25: Parent later reported combined audit coverage at 62/62 for all four families with issues=0; local scan now also confirms 62/62 canonical hamster frames.

## Finalization Checkpoint

Final local fullrun counts, excluding macOS `._*` AppleDouble files:

- Raw attempt PNGs: 68
- Manifest rows: 68
- Attempt normalized PNGs: 68
- QA JSON files: 68
- Per-attempt contact PNGs: 68
- Active canonical normalized PNGs visible locally: 62
- Local canonical frames visible: `00-61`
- Local canonical frames missing: none
- Mechanical QA over local QA JSON: no failures
- Contact sheet rebuilt: `contacts/hamster-set00-fullrun-contact.png`

## Remaining Risk

- Mechanical QA catches size, alpha, edge contact, empty alpha, and green residue only.
- Parent thread still needs final visual adoption review for pose readability, anatomy, species identity, and style consistency.
- Many late-frame candidates are intentionally retained as candidates despite high-detail/realistic style risk; this is recorded in per-frame manifest/QA notes and should be considered during final adoption.
