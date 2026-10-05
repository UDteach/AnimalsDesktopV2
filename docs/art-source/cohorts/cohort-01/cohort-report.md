# Cohort 01 Report

Status: prompts_ready_assets_not_generated
Owner scope: `docs/art-source/cohorts/cohort-01/`

## Summary

This cohort covers the first 10 Degu variants using the `degu` motion profile. Real ImageGen file output was not available in this thread, so this report records generation prompts, expected filenames, acceptance checks, and unresolved risks for parent-thread integration.

## Shared Art Direction

Generate realistic 2D sprite source art for a degu desktop pet. Each image must contain exactly one complete degu on a transparent background. Use a consistent side-facing camera, compact desktop mascot scale, natural degu anatomy, visible ears, feet, whiskers, and tail. Do not include text, border, scenery, shadow, costume, props, multiple animals, or human-like pose.

Motion source frames must preserve the same camera, scale, baseline, contact points, and body proportions across all 62 frames. The action set follows the existing `degu` motion profile: idle, walk, scurry, eat/forage, hop, groom, turn, rest, alert, and wheel-capable run source when parent integration maps frames.

## Required Output Paths

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00_idle_00.png` through `motion-source/<variant-id>/61_alert_05.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

## Variant Prompts

### `wild_agouti`

Source-truth prompt:

```text
One complete realistic degu, wild agouti coat, warm brown gray guard hairs, cream belly, natural dark eyes, side-facing desktop pet sprite, transparent background, consistent neutral pose, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same wild agouti degu. Preserve camera, scale, anatomy, baseline, and contact points. Include natural degu idle breathing, quick four-foot walk, fast scurry, low forage nibble, compact hop, face grooming, turn anticipation, rest crouch, alert head lift, and wheel-capable run motion. One animal only, no background or shadows.
```

### `black`

Source-truth prompt:

```text
One complete realistic degu, solid black coat with subtle charcoal highlights so the silhouette remains readable, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same black degu. Keep highlights stable and readable on dark backgrounds. Preserve camera, scale, anatomy, baseline, and contact points across idle, walk, scurry, forage nibble, hop, groom, turn, rest, alert, and wheel-capable run actions.
```

### `blue`

Source-truth prompt:

```text
One complete realistic degu, blue slate gray coat, soft cool gray body with lighter belly, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same blue slate gray degu. Preserve camera, scale, anatomy, baseline, and contact points. Actions: idle breathing, walk, scurry, forage nibble, compact hop, groom, turn, rest, alert head lift, and wheel-capable run.
```

### `gray`

Source-truth prompt:

```text
One complete realistic degu, neutral gray coat, slightly lighter underside, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same gray degu. Preserve camera, scale, anatomy, baseline, and contact points. Keep the degu low and natural, with idle, walk, scurry, forage nibble, hop, groom, turn, rest, alert, and wheel-capable run actions.
```

### `white_cream`

Source-truth prompt:

```text
One complete realistic degu, white cream coat with gentle warm shading for readability, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same white cream degu. Preserve camera, scale, anatomy, baseline, and contact points. Keep the pale coat readable without outlines or shadows. Include idle, walk, scurry, forage nibble, hop, groom, turn, rest, alert, and wheel-capable run actions.
```

### `sand_champagne`

Source-truth prompt:

```text
One complete realistic degu, sand champagne coat, pale beige tan body with cream underside, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same sand champagne degu. Preserve camera, scale, anatomy, baseline, and contact points. Actions: idle, walk, scurry, forage nibble, compact hop, groom, turn, rest, alert head lift, and wheel-capable run.
```

### `chocolate`

Source-truth prompt:

```text
One complete realistic degu, chocolate brown coat, warm dark brown fur with lighter belly, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same chocolate degu. Preserve camera, scale, anatomy, baseline, and contact points. Keep dark fur readable on light and dark previews. Include idle, walk, scurry, forage nibble, hop, groom, turn, rest, alert, and wheel-capable run actions.
```

### `black_pied`

Source-truth prompt:

```text
One complete realistic degu, black pied coat with black body patches and clear cream white piebald areas, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same black pied degu. Preserve the same piebald patch placement across all frames. Preserve camera, scale, anatomy, baseline, and contact points. Include idle, walk, scurry, forage nibble, hop, groom, turn, rest, alert, and wheel-capable run actions.
```

### `agouti_pied`

Source-truth prompt:

```text
One complete realistic degu, agouti pied coat with wild agouti brown gray areas and cream white piebald patches, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same agouti pied degu. Preserve patch placement and readable agouti texture across frames. Preserve camera, scale, anatomy, baseline, and contact points. Include idle, walk, scurry, forage nibble, hop, groom, turn, rest, alert, and wheel-capable run actions.
```

### `blue_pied`

Source-truth prompt:

```text
One complete realistic degu, blue slate gray pied coat with cool gray patches and cream white piebald areas, natural dark eyes, side-facing desktop pet sprite, transparent background, visible ears, feet, whiskers, and tail, no text, no border, no scenery, no shadow, no costume, no props.
```

Motion prompt:

```text
Create a 62-frame transparent source set for the same blue pied degu. Preserve patch placement, camera, scale, anatomy, baseline, and contact points. Include idle breathing, walk, scurry, forage nibble, compact hop, groom, turn, rest, alert, and wheel-capable run actions.
```

## Acceptance Checklist

- [ ] Exactly one complete degu in every source image.
- [ ] Transparent background, no text, no border, no scenery, no shadow, no costume, no props.
- [ ] Ears, feet, whiskers, and tail are not cropped.
- [ ] Consistent camera, scale, baseline, and contact points across 62 motion frames.
- [ ] Pied variants keep stable patch placement across frames.
- [ ] Light and dark 96x64 previews remain readable.
- [ ] Parent integration can map the 62 source frames to the existing `degu` runtime frame contract.

## Parent Integration Notes

Do not promote this cohort until real PNG files exist at the required paths. Current status is prompt-ready only. The parent thread should keep catalog source status as `prototype_only` or existing degu source status until the actual ImageGen assets pass visual review.
