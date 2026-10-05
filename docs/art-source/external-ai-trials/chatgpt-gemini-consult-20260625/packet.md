# Sanitized External AI Consultation Packet

- Date: 2026-06-25
- Destination candidates: ChatGPT or Gemini
- Data classification: PUBLIC project/spec summary only
- External send authorization: user allowed trying ChatGPT/Gemini in a separate review thread
- Image upload: not included in this packet
- Local paths, private URLs, credentials, raw logs: removed

## Review Objective

We are building a lightweight Windows desktop pet app that uses 96x64 transparent 2D sprite frames. The next asset task is to make four small-animal sprite families readable and motion-consistent across one 62-frame animation set:

- `hamster_golden_syrian`
- `macaroni_mouse_tan` (fat-tailed gerbil / macaroni mouse style)
- `sugar_glider_gray`
- `rabbit_chestnut_agouti`

Please advise on prompt wording and motion breakdown only. Do not assume the output is final art. We will verify locally with transparent PNG checks, contact-sheet review, and frame audits.

## Fixed Constraints

- Each final frame is one standalone 96x64 transparent PNG.
- One complete animal only, right-facing side view.
- Stable camera distance, body scale, contact baseline, outline softness, and lighting across all 62 frames in a family.
- No text, labels, numbers, borders, dividers, scenery, props, costumes, floor bands, cast shadows, checker backgrounds, multiple animals, cropped anatomy, or human-like poses.
- Do not recolor one animal into another species.
- Do not treat a simple recolor, vertical bob, or duplicated walk frame as a different action.
- If true alpha is not available, a flat pure green background may be used as a local cleanup input, but the accepted frame must be transparent.
- Runtime readability matters more than painterly detail. The animal must still read at 96x64.

## Current Frame-00 Baseline, Described In Words

No image is attached. A local contact-sheet preview shows frame-00 anchors on light, dark, and checker backgrounds.

- Hamster: golden body, rounded cheek-forward shape, tiny paws, short/no visible tail. Reads as hamster, but later frames must preserve the blunt muzzle and compact cheek silhouette.
- Macaroni mouse / fat-tailed gerbil: compact tan rodent with a thick pale tail and small feet. Reads as distinct from hamster mainly through the thick tail and leaner body; later frames must not drift into normal mouse or hamster shapes.
- Sugar glider: gray grounded side-view animal with large dark eyes, tail, and visible membrane line. The mask/face and membrane help readability, but later frames must stay grounded and avoid a flying-scene pose.
- Rabbit: chestnut agouti body, long ears, strong hindquarters, small tail. Reads clearly, but the ears need protected top margin and hop frames can easily over-jump or crop.

## Existing 62-Frame Set Structure

The current set uses these motion groups:

- 00-03: idle, breathing, ear/whisker/head adjustments, tiny weight shift
- 04-11: walk cycle with shared baseline contacts
- 12-19: faster scurry or low travel cycle
- 20-25: forage/sniff/nibble/search poses
- 26-31: species-specific action
- 32-39: side-view turn or partial body/head rotation
- 40-43: eat/chew poses
- 44-47: ground check or light paw scratch
- 48-51: stand/rest/alert hold
- 52-55: groom/preen/face clean
- 56-61: reaction or alert burst and recovery

Species-specific action candidates:

- Hamster: tiny paw lift, cheek/muzzle sniff, or wheel-safe quick crouch without drawing a wheel.
- Macaroni mouse / fat-tailed gerbil: quick crouch, thick-tail balance shift, or low skitter emphasis.
- Sugar glider: grounded patagium stretch or short low skitter, not full flight.
- Rabbit: hop prep, compact hop peak, and landing/recovery with protected ear margin.

## Questions For The Reviewer

1. How should the prompt template be changed so each of the four animals remains readable at 96x64 while preserving consistent scale, baseline, and side-view pose across 62 frames?
2. Which motion groups should be simplified, merged, split, or reworded to avoid duplicate frames, vertical bobbing, or species drift?
3. What per-species negative prompts or guardrails would best prevent the common failures below?
   - hamster becoming mouse-like or losing cheek mass
   - macaroni mouse becoming hamster-like or thin-tailed
   - sugar glider becoming a flying scene, over-masked, or too poster-like
   - rabbit ears being cropped or hop frames becoming too airborne
4. Give a concise improved prompt recipe for frames 01-03 first, because those subtle idle variants are the consistency gate before walk/scurry.
5. Give a concise improved prompt recipe for the first high-risk motion group for each species.

## Requested Output Format

Please answer as advisory hypotheses, not final truth:

```text
Recommended global prompt changes:
- ...

Motion breakdown changes:
- adopt:
- hold:
- reject:

Per-species prompt guards:
- hamster_golden_syrian:
- macaroni_mouse_tan:
- sugar_glider_gray:
- rabbit_chestnut_agouti:

Frame 01-03 idle prompt recipe:
- ...

First high-risk group prompt recipe:
- hamster_golden_syrian:
- macaroni_mouse_tan:
- sugar_glider_gray:
- rabbit_chestnut_agouti:

Risks to verify locally:
- ...
```
