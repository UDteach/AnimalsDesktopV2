# Macaroni Mouse Rescue Frames 56-61 - 2026-06-25

Target: `macaroni_mouse_tan` set00 reaction frames 56-61.

Write boundary: this rescue folder only. This lane did not modify the main fullrun directory, `accepted-frames`, catalog/runtime code, GitHub Pages, or release assets.

## Scope

This rescue lane generated only the final reaction range:

- 56: small surprise reaction
- 57: quick recoil low quadruped
- 58: tiny startle stretch forward
- 59: settle shoulders lower
- 60: return-to-idle transition
- 61: final neutral recovery matching scale and stance

Parent reported the combined macaroni-mouse audit now has 62/62 coverage, so generation stopped after this rescue packet was finalized.

## Method

- Built-in Codex `image_gen`, one output per frame.
- Source images copied from `/Users/kyota/.codex/generated_images/019efc43-f614-78a3-8f80-c0295aae9fbc`.
- Background prompt: perfectly flat `#00ff00`.
- Alpha conversion: `/Users/kyota/.codex/skills/.system/imagegen/scripts/remove_chroma_key.py`.
- Normalization: subject fit inside 86x54, placed on a 96x64 transparent canvas.
- Artifacts written under this folder: `raw/`, `alpha/`, `frames/`, `contacts/`, `qa/`, `prompts/`, `manifest.csv`, and `macaroni-mouse-56-61-contact.png`.

## Results

| Frame | Final attempt | Decision | Notes |
| --- | ---: | --- | --- |
| 56 | 03 | accepted | Attempt 01 had 1 greenish pixel and a more upright reaction posture; attempt 02 still had 1 greenish pixel; attempt 03 used edge-contract chroma cleanup and passed mechanical QA. |
| 57 | 01 | accepted | Quick recoil, low quadruped, short thick pale tail retained. |
| 58 | 01 | accepted | Forward startle stretch, low body, short thick pale tail retained. |
| 59 | 01 | accepted | Settling posture with lower shoulders, complete animal. |
| 60 | 01 | accepted | Return-to-idle transition, stable baseline. |
| 61 | 01 | accepted | Neutral recovery, stable scale and stance. |

Mechanical QA summary:

- Final frames: 6
- Attempt rows: 8
- Accepted current frames: 6
- Retry-needed current frames: 0
- Final frame size: 96x64 PNG with alpha
- Final edge-alpha failures: 0
- Final green-residue failures: 0
- Empty-alpha failures: 0

## Visual Review

Contact sheet:

- `macaroni-mouse-56-61-contact.png`

The final six frames preserve a tan small rodent silhouette with a pointed nose, pale belly, and short thick pale tail. They avoid the hard rejects for this lane: long thin mouse tail, missing tail, hamster cheek pouches, bulky hamster body, text, props, borders, scenery, multiple animals, and cropped body parts.

This folder is a rescue candidate packet only. Parent integration should compare these frames against the main fullrun outputs before any promotion or replacement.
