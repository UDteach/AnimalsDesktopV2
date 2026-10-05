# Cohort 03 Report

Status: prompt-ready, image assets not generated in this delegated thread.
Owner boundary: only `docs/art-source/cohorts/cohort-03/`.

## Output Contract

Create the following when ImageGen file export is available:

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00_idle_00.png` through `61_alert_05.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

All frames must show one complete animal on transparent background, with no text, no border, no scenery, no shadow, no costume, and no human-like pose. Keep camera, scale, baseline, anatomy, contact points, and facing direction consistent across all frames for a variant.

## Shared Motion Frame Map

Use the existing 62-frame contract:

- 00-03: idle breathing
- 04-11: walk
- 12-19: faster locomotion
- 20-25: eat / nibble / sniff
- 26-31: species action
- 32-39: groom
- 40-47: turn / direction transition
- 48-55: rest / settle
- 56-61: alert / look-around

For rabbit-hop variants, the species action is a compact hop with hind feet as the contact anchor. For dog-trot variants, the faster locomotion is a four-beat trot with paws contacting the same baseline. For cat-stalk variants, locomotion is low, quiet, and shoulder-led with minimal vertical bob.

## Variant Prompts

### `rabbit_white`

Source truth prompt:

```text
Single white domestic rabbit, compact side view, full body visible, transparent background, natural proportions, long ears fully inside frame, small round tail visible, clean paws on a consistent ground baseline, no shadow, no scenery, no text, no border, no costume, no human pose, sprite source for a 2D Windows desktop pet, soft realistic fur texture readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same white rabbit, side view, constant camera and scale, rabbit-hop motion profile. Frames 00-03 idle breathing, 04-11 cautious hop-walk, 12-19 quick hop, 20-25 nibbling low, 26-31 compact hop with hind feet contact anchor, 32-39 grooming face and ears, 40-47 turn transition, 48-55 resting crouch, 56-61 alert ear movement. Keep ears, feet, whiskers, and tail uncropped.
```

### `rabbit_blue_gray`

Source truth prompt:

```text
Single blue gray domestic rabbit, compact side view, full body visible, transparent background, natural proportions, ears fully inside frame, small tail visible, paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source, soft gray-blue coat readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same blue gray rabbit, side view, constant camera and scale, rabbit-hop motion profile. Use the frame map exactly, with compact hop arcs and stable hind foot contact. Preserve blue-gray coat, ear shape, body volume, and baseline across all frames.
```

### `rabbit_fawn`

Source truth prompt:

```text
Single fawn domestic rabbit, warm tan coat, compact side view, full body visible, transparent background, natural proportions, ears and small tail fully visible, paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same fawn rabbit, side view, constant camera and scale, rabbit-hop motion profile. Keep the warm tan color consistent. Hop frames must feel rabbit-specific, not vertical bobbing. Avoid degu wheel or upright running poses.
```

### `dog_cream_tan`

Source truth prompt:

```text
Single small dog, cream and tan coat, side view, full body visible, transparent background, natural dog anatomy, ears and tail fully inside frame, four paws on a consistent ground baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source, readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same cream and tan small dog, side view, constant camera and scale, dog-trot motion profile. Frames 04-11 are relaxed walk, 12-19 are four-beat trot, 20-25 sniff/eat low, 26-31 playful paw lift or small stretch, 32-39 grooming, 40-47 turn, 48-55 rest, 56-61 alert head and ear motion. Keep paw contacts grounded.
```

### `dog_black_tan`

Source truth prompt:

```text
Single small black and tan dog, side view, full body visible, transparent background, natural dog anatomy, tan muzzle/eyebrow/chest accents, ears and tail uncropped, four paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same black and tan small dog, side view, constant camera and scale, dog-trot motion profile. Preserve black body and tan markings consistently across walk, trot, sniff, paw-lift, groom, turn, rest, and alert frame groups.
```

### `dog_white`

Source truth prompt:

```text
Single small white dog, side view, full body visible, transparent background, natural proportions, subtle off-white fur detail visible on light and dark backgrounds, ears and tail fully inside frame, paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same white small dog, side view, constant camera and scale, dog-trot motion profile. Maintain readable white fur edges with transparent background and no shadow. Paw contact and head height must remain stable.
```

### `dog_sable`

Source truth prompt:

```text
Single small sable dog, side view, full body visible, transparent background, natural dog anatomy, sable brown coat with darker back and muzzle, ears and tail fully visible, four paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same sable small dog, side view, constant camera and scale, dog-trot motion profile. Preserve sable gradient and darker muzzle/back marks while animating walk, trot, sniff, paw-lift, groom, turn, rest, and alert.
```

### `dog_gray`

Source truth prompt:

```text
Single small gray dog, side view, full body visible, transparent background, natural dog anatomy, soft gray coat, ears and tail fully inside frame, four paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same gray small dog, side view, constant camera and scale, dog-trot motion profile. Keep the trot dog-specific with alternating paw contacts and no rabbit hop or degu scurry.
```

### `cat_brown_tabby`

Source truth prompt:

```text
Single brown tabby cat, side view, full body visible, transparent background, natural feline anatomy, tabby stripes readable but not noisy at 96x64, ears, whiskers, paws, and tail fully inside frame, paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same brown tabby cat, side view, constant camera and scale, cat-stalk motion profile. Frames 04-11 are quiet walk, 12-19 are low stalking movement with little vertical bob, 20-25 sniff/eat, 26-31 crouch and tail flick, 32-39 grooming paw/face, 40-47 turn, 48-55 curled rest, 56-61 alert ear and tail movement. Keep tabby markings aligned.
```

### `cat_black`

Source truth prompt:

```text
Single black cat, side view, full body visible, transparent background, natural feline anatomy, readable black fur silhouette with subtle highlights, ears, whiskers, paws, and tail fully inside frame, paws on one baseline, no shadow, no scenery, no text, no border, no costume, no human pose, 2D desktop pet sprite source readable at 96x64.
```

Motion prompt:

```text
Generate 62 transparent PNG frames for the same black cat, side view, constant camera and scale, cat-stalk motion profile. Preserve readable silhouette and subtle highlight edges on transparent background. Movement should be feline, low, quiet, and shoulder-led.
```

## Parent Integration Notes

- Current status is prompt-ready only; no PNG assets were created by this delegated thread.
- Parent should run actual ImageGen/export into the required paths, then validate alpha, bounds, frame count, and 96x64 readability.
- For shared dog source families, do not accept simple hue-shift-only output as final unless anatomy, markings, and coat identity remain credible for each variant.
- For rabbit variants, accept coat-color family reuse only if the motion source preserves rabbit-specific hind-leg hop mechanics.
