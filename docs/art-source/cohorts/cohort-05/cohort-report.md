# Cohort 05 Report

Status: imagegen_prompts_ready
Owner boundary: `docs/art-source/cohorts/cohort-05/` only

## Output Status

This delegated thread could not save ImageGen PNG outputs directly into the workspace. The cohort therefore stops at the contract fallback: exact prompts, filenames, motion source plan, acceptance checklist, and known risks are recorded here for parent-thread or ImageGen-capable follow-up execution.

Required target layout when generated:

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00_idle_00.png` through `61_alert_05.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

## Shared Art Direction

All generated images must be transparent PNGs with one complete animal, side-view, right-facing, compact desktop-pet sprite posture, natural anatomy, no scene, no props, no shadows, no border, no text, no costume, no human pose. Keep camera, body scale, baseline, and contact points consistent across every frame in a variant. Leave safe padding around ears, paws, whiskers, and tails.

Use a clean small-animal sprite style: detailed enough to read at 96x64, soft but not plush/toy-like, realistic coat markings, crisp alpha edge, no background lighting.

## Variant Prompts

### `hamster_cream`

Source-truth prompt:

```text
One complete cream Syrian hamster, natural side-view facing right, compact desktop pet sprite pose, transparent background, realistic cream coat, small rounded ears, black eyes, short paws visible, no accessories, no scenery, no shadow, no text, full body including whiskers and feet with padding, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `small-rodent-scurry`

Motion source frames:

- `00_idle_00` to `03_idle_03`: breathing and tiny whisker twitch.
- `04_walk_00` to `11_walk_07`: short quick steps, low body drift, paws alternating.
- `12_scurry_00` to `19_scurry_07`: faster rodent dash, body stretched slightly forward.
- `20_eat_00` to `25_eat_05`: sitting low, nibbling with forepaws.
- `26_groom_00` to `31_groom_05`: face wash, forepaw up, no upright human pose.
- `32_turn_00` to `39_turn_07`: right-facing to near-front and back, same footprint.
- `40_sniff_00` to `47_sniff_07`: nose down, whisker twitch, tiny step.
- `48_rest_00` to `55_rest_07`: crouched pause, slow breathing.
- `56_alert_00` to `61_alert_05`: head lift and ear perk.

### `hamster_black_banded`

Source-truth prompt:

```text
One complete black-banded Syrian hamster, side-view facing right, cream base coat with a clear black band around the midbody, compact natural sprite stance, transparent background, small rounded ears, black eyes, visible paws and whiskers, no scenery, no shadow, no text, full body safely padded, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `small-rodent-scurry`

Motion source frames: use the `hamster_cream` sequence, preserving the black band registration so it does not slide or warp between frames.

### `hamster_white`

Source-truth prompt:

```text
One complete white Syrian hamster, side-view facing right, compact desktop pet sprite pose, transparent background, clean white coat with subtle warm shading, black eyes, pink ears and paws, visible whiskers, no accessories, no scenery, no shadow, no text, full body with padding, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `small-rodent-scurry`

Motion source frames: use the `hamster_cream` sequence, with white coat contrast checked on light and dark preview backgrounds.

### `hamster_cinnamon`

Source-truth prompt:

```text
One complete cinnamon Syrian hamster, side-view facing right, warm cinnamon brown coat, compact natural desktop pet sprite stance, transparent background, small ears, black eyes, tiny paws and whiskers visible, no scenery, no props, no shadow, no text, full body safely padded, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `small-rodent-scurry`

Motion source frames: use the `hamster_cream` sequence, preserving cinnamon tone and paw contact points.

### `ferret_sable`

Source-truth prompt:

```text
One complete sable ferret, long low body, natural side-view facing right, desktop pet sprite pose, transparent background, dark sable mask, tan body, small rounded ears, long tail fully visible, paws touching baseline, no scenery, no shadow, no text, no costume, full body with padding, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `ferret-slink`

Motion source frames:

- `00_idle_00` to `03_idle_03`: low breathing, tail still.
- `04_slink_00` to `11_slink_07`: elongated low walk with body wave and alternating paws.
- `12_scurry_00` to `19_scurry_07`: quick ferret dash, body close to baseline.
- `20_sniff_00` to `25_sniff_05`: nose down, neck extended.
- `26_groom_00` to `31_groom_05`: curl slightly to groom side, not upright.
- `32_turn_00` to `39_turn_07`: low turn with tail following body arc.
- `40_tunnel_00` to `47_tunnel_07`: crouch and creep, shoulders dipping.
- `48_rest_00` to `55_rest_07`: long body settled, tiny breathing.
- `56_alert_00` to `61_alert_05`: head lift, tail remains visible.

### `ferret_albino`

Source-truth prompt:

```text
One complete albino ferret, long low body, natural side-view facing right, white to cream fur, pink nose, red or ruby eyes, transparent background, small ears, long tail fully visible, paws on baseline, no scenery, no shadow, no text, no costume, full body with padding, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `ferret-slink`

Motion source frames: use the `ferret_sable` sequence, preserving pale-fur visibility on both light and dark previews.

### `ferret_champagne`

Source-truth prompt:

```text
One complete champagne ferret, long low body, natural side-view facing right, warm champagne tan coat with lighter face, transparent background, small ears, dark eyes, long tail fully visible, paws touching baseline, no scenery, no shadow, no text, no costume, full body with padding, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `ferret-slink`

Motion source frames: use the `ferret_sable` sequence, preserving mask and tail position between frames.

### `guinea_pig_tricolor`

Source-truth prompt:

```text
One complete tricolor guinea pig, rounded compact body, natural side-view facing right, patches of black, white, and warm brown, transparent background, small ears, visible feet, short muzzle, no scenery, no shadow, no text, no costume, full body safely padded, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `guinea-pig-waddle`

Motion source frames:

- `00_idle_00` to `03_idle_03`: round-body breathing, tiny ear motion.
- `04_waddle_00` to `11_waddle_07`: slow side-to-side weight shift with short steps.
- `12_scurry_00` to `19_scurry_07`: faster low run, body still rounded.
- `20_eat_00` to `25_eat_05`: hay/food nibble posture without drawn food prop.
- `26_groom_00` to `31_groom_05`: head dip and cheek rub.
- `32_turn_00` to `39_turn_07`: rounded body turn, feet stay grounded.
- `40_sniff_00` to `47_sniff_07`: nose forward/down, small shuffle.
- `48_rest_00` to `55_rest_07`: loafed low rest, slow breathing.
- `56_alert_00` to `61_alert_05`: head lift, slight body rise.

### `guinea_pig_cream`

Source-truth prompt:

```text
One complete cream guinea pig, rounded compact body, natural side-view facing right, warm cream coat with subtle shading, transparent background, small ears, visible feet, short muzzle, no scenery, no shadow, no text, no costume, full body safely padded, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `guinea-pig-waddle`

Motion source frames: use the `guinea_pig_tricolor` sequence, checking cream coat readability on light preview.

### `guinea_pig_black`

Source-truth prompt:

```text
One complete black guinea pig, rounded compact body, natural side-view facing right, black coat with readable soft highlights, transparent background, small ears, visible feet, short muzzle, no scenery, no shadow, no text, no costume, full body safely padded, consistent 2D sprite camera, PNG alpha.
```

Motion profile: `guinea-pig-waddle`

Motion source frames: use the `guinea_pig_tricolor` sequence, checking black coat readability on dark preview.

## Acceptance Checklist

- [ ] Every `source-truth/*.png` has alpha background and exactly one animal.
- [ ] Every variant has 62 files in `motion-source/<variant-id>/`.
- [ ] All frame canvases keep the same camera, scale, baseline, and contact points.
- [ ] Hamsters keep short quick scurry motion and never use ferret body waves.
- [ ] Ferrets keep a long, low slink and never become upright rodents.
- [ ] Guinea pigs keep rounded waddle motion and do not hop like rabbits.
- [ ] Light/dark 96x64 previews preserve eyes, feet, ears, whiskers, and tail where applicable.
- [ ] No generated frame contains text, border, scenery, shadow, costume, or multiple animals.
