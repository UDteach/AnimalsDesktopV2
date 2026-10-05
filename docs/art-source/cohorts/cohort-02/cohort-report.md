# Cohort 02 Report

Status: prompt-ready, PNG generation not completed in this delegated pass
Owner boundary: `docs/art-source/cohorts/cohort-02/`

## Scope

This cohort owns these 10 variants:

1. `cream_pied` - Degu, cream pied, motion `degu`
2. `chinchilla_standard_gray` - Chinchilla, standard gray, motion `small-rodent-scurry`
3. `chinchilla_beige` - Chinchilla, beige, motion `small-rodent-scurry`
4. `chinchilla_ebony` - Chinchilla, ebony, motion `small-rodent-scurry`
5. `chinchilla_white_mosaic` - Chinchilla, white mosaic, motion `small-rodent-scurry`
6. `macaroni_mouse_tan` - Macaroni mouse / fat-tailed gerbil, tan, motion `small-rodent-scurry`
7. `macaroni_mouse_gray` - Macaroni mouse / fat-tailed gerbil, gray, motion `small-rodent-scurry`
8. `macaroni_mouse_cream` - Macaroni mouse / fat-tailed gerbil, cream, motion `small-rodent-scurry`
9. `rabbit_chestnut_agouti` - Rabbit, chestnut agouti, motion `rabbit-hop`
10. `rabbit_black` - Rabbit, black, motion `rabbit-hop`

## Output Contract

When ImageGen is run, write only these paths:

- `source-truth/<variant-id>.png`
- `motion-source/<variant-id>/00_idle_00.png` through `61_alert_05.png`
- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

Do not write importer outputs, generated runtime sprites, catalog files, or workflow files from this cohort.

## Shared ImageGen Style Contract

Use this style clause in every prompt:

```text
transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers where anatomically present, full tail where present, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

## Source-Truth Prompts

### cream_pied

```text
Create source-truth art for an octodon degu desktop-pet sprite, cream pied coat with soft cream body and irregular lighter pied patches, natural degu proportions, rounded ears, visible whiskers, long tail with tuft, alert but animal-like posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### chinchilla_standard_gray

```text
Create source-truth art for a standard gray chinchilla desktop-pet sprite, dense plush gray fur, pale belly, large rounded ears, compact body, long whiskers, fluffy tail visible behind body, natural standing quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### chinchilla_beige

```text
Create source-truth art for a beige chinchilla desktop-pet sprite, warm beige dense fur, pale cream belly, large rounded ears, compact body, long whiskers, fluffy tail visible behind body, natural standing quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### chinchilla_ebony

```text
Create source-truth art for an ebony chinchilla desktop-pet sprite, deep charcoal black dense fur with subtle warm highlights, large rounded ears, compact body, long whiskers, fluffy tail visible behind body, natural standing quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### chinchilla_white_mosaic

```text
Create source-truth art for a white mosaic chinchilla desktop-pet sprite, mostly white and pale silver patchwork fur, large rounded ears, compact body, long whiskers, fluffy tail visible behind body, natural standing quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### macaroni_mouse_tan

```text
Create source-truth art for a fat-tailed gerbil, also known as a macaroni mouse, desktop-pet sprite, tan sandy coat, compact rounded body, short muzzle, small rounded ears, black eyes, distinctive thick club-shaped tail, tiny feet, natural low quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### macaroni_mouse_gray

```text
Create source-truth art for a fat-tailed gerbil, also known as a macaroni mouse, desktop-pet sprite, soft gray coat, compact rounded body, short muzzle, small rounded ears, black eyes, distinctive thick club-shaped tail, tiny feet, natural low quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### macaroni_mouse_cream

```text
Create source-truth art for a fat-tailed gerbil, also known as a macaroni mouse, desktop-pet sprite, pale cream coat, compact rounded body, short muzzle, small rounded ears, black eyes, distinctive thick club-shaped tail, tiny feet, natural low quadruped posture. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### rabbit_chestnut_agouti

```text
Create source-truth art for a chestnut agouti rabbit desktop-pet sprite, warm brown agouti coat with lighter belly, upright natural ears, compact rabbit body, visible hind legs, small tail, natural grounded quadruped posture ready to hop. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

### rabbit_black

```text
Create source-truth art for a black rabbit desktop-pet sprite, deep black coat with subtle visible highlights, upright natural ears, compact rabbit body, visible hind legs, small tail, natural grounded quadruped posture ready to hop. transparent background, single complete animal only, side-view 2D sprite source, no text, no border, no scenery, no shadow, no props, no costume, no human-like pose, full ears, full feet, full whiskers, full tail, consistent camera, consistent scale, consistent baseline, clean alpha edges, natural anatomy, compact desktop-pet silhouette, source art for 96x64 runtime sprite downscaling
```

## Motion Source Prompt Template

For each variant, generate 62 frames with identical camera, scale, and baseline. Use the source-truth prompt for the variant, then append the frame-specific action from this table.

| Frames | Action |
| --- | --- |
| `00_idle_00` to `03_idle_03` | calm breathing idle, feet planted, only subtle body compression |
| `04_walk_00` to `11_walk_07` | profile-specific walk loop, feet remain anatomically plausible |
| `12_scurry_00` to `19_scurry_07` | quicker profile-specific movement, no blur, no duplicated body |
| `20_eat_00` to `25_eat_05` | natural nibbling/eating posture without props |
| `26_action_00` to `31_action_05` | degu small hop for `cream_pied`, chinchilla/macaroni mouse quick scurry recoil, rabbit hop crouch-to-lift |
| `32_groom_00` to `39_groom_07` | natural grooming or face-wipe behavior |
| `40_turn_00` to `47_turn_07` | smooth body turn source sequence, no flipped text or duplicate animals |
| `48_rest_00` to `55_rest_07` | low rest/settle posture, still complete body |
| `56_alert_00` to `61_alert_05` | brief alert look, ears and whiskers intact |

## Acceptance Notes

- Chinchilla and macaroni mouse variants may share a species motion family, but each coat must preserve the same silhouette and baseline.
- Rabbit frames must have larger hind-leg compression and hop lift than small rodents.
- Degu `cream_pied` can reuse the existing degu motion action taxonomy, but new source frames must still be the cream pied coat and not a recolored runtime sheet.
- This report intentionally does not promote any asset to accepted ImageGen status; parent integration should only promote after PNG review.
