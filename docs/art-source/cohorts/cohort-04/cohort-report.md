# Cohort 04 Report

Status: prompts_ready_no_png_output
Owner scope: `docs/art-source/cohorts/cohort-04/`

## Execution Note

This child thread could not emit local ImageGen PNG files directly into the repository. No source-truth PNGs, 62-frame motion source directories, or 96x64 preview PNGs were generated here. The report below is ready for the parent thread or an ImageGen-capable worker to execute without touching files outside this cohort directory.

## Shared Art Direction

- Transparent background, single complete animal, no ground plane, no shadow, no border, no scenery, no text.
- Pixel-sprite friendly 2D illustration with clean silhouette, soft natural anatomy, and readable limbs at 96x64.
- Consistent side-facing camera, animal facing right, baseline contact point stable across frames.
- Keep ears, paws/toes, whiskers, tail, and gecko toes fully inside the canvas.
- Avoid human-like poses, costumes, props, exaggerated expressions, or multiple animals.
- For source-truth stills, use a 4:3 transparent canvas around the animal and leave enough margin for motion offsets.

## Variant Source-Truth Prompts

### `cat_white`

Prompt:
`One complete white domestic shorthair cat, right-facing side profile, natural standing pose, full body visible, realistic white coat with subtle cream shading, visible ears whiskers paws and tail, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/cat_white.png`
- `preview/cat_white-96x64-light.png`
- `preview/cat_white-96x64-dark.png`

### `cat_orange_tabby`

Prompt:
`One complete orange tabby domestic cat, right-facing side profile, natural standing pose, full body visible, warm orange coat with readable tabby stripes on body tail and face, visible ears whiskers paws and tail, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/cat_orange_tabby.png`
- `preview/cat_orange_tabby-96x64-light.png`
- `preview/cat_orange_tabby-96x64-dark.png`

### `cat_gray`

Prompt:
`One complete gray domestic shorthair cat, right-facing side profile, natural standing pose, full body visible, blue-gray coat with subtle shading, visible ears whiskers paws and tail, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/cat_gray.png`
- `preview/cat_gray-96x64-light.png`
- `preview/cat_gray-96x64-dark.png`

### `cat_cream`

Prompt:
`One complete cream domestic shorthair cat, right-facing side profile, natural standing pose, full body visible, pale cream coat with warm beige shading, visible ears whiskers paws and tail, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/cat_cream.png`
- `preview/cat_cream-96x64-light.png`
- `preview/cat_cream-96x64-dark.png`

### `gecko_gray_brown`

Prompt:
`One complete gray-brown gecko, right-facing low side profile, natural crawling pose with four feet visible and toes inside canvas, full tail visible, mottled gray brown reptile skin, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/gecko_gray_brown.png`
- `preview/gecko_gray_brown-96x64-light.png`
- `preview/gecko_gray_brown-96x64-dark.png`

### `gecko_leopard`

Prompt:
`One complete leopard gecko, right-facing low side profile, natural crawling pose with four feet visible and toes inside canvas, full tail visible, yellow tan body with dark leopard spots, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/gecko_leopard.png`
- `preview/gecko_leopard-96x64-light.png`
- `preview/gecko_leopard-96x64-dark.png`

### `gecko_tangerine`

Prompt:
`One complete tangerine gecko, right-facing low side profile, natural crawling pose with four feet visible and toes inside canvas, full tail visible, vivid orange tangerine morph with soft lighter belly, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/gecko_tangerine.png`
- `preview/gecko_tangerine-96x64-light.png`
- `preview/gecko_tangerine-96x64-dark.png`

### `gecko_blizzard`

Prompt:
`One complete blizzard gecko, right-facing low side profile, natural crawling pose with four feet visible and toes inside canvas, full tail visible, pale gray white blizzard morph with subtle scale shading, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/gecko_blizzard.png`
- `preview/gecko_blizzard-96x64-light.png`
- `preview/gecko_blizzard-96x64-dark.png`

### `gecko_albino`

Prompt:
`One complete albino gecko, right-facing low side profile, natural crawling pose with four feet visible and toes inside canvas, full tail visible, pale cream yellow albino morph with soft pink undertones, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/gecko_albino.png`
- `preview/gecko_albino-96x64-light.png`
- `preview/gecko_albino-96x64-dark.png`

### `hamster_golden_syrian`

Prompt:
`One complete golden Syrian hamster, right-facing side profile, compact natural stance, full body visible with small ears paws whiskers and short tail area inside canvas, golden brown coat with cream underside, clean 2D sprite source art, transparent background, no shadow, no text, no props, no scenery.`

Target files:
- `source-truth/hamster_golden_syrian.png`
- `preview/hamster_golden_syrian-96x64-light.png`
- `preview/hamster_golden_syrian-96x64-dark.png`

## Motion Source Contract

Each variant must provide `motion-source/<variant-id>/frame-00.png` through `frame-61.png`.

### `cat-stalk` frames for cat variants

- 00-03 idle breathing: tiny ribcage rise, tail steady.
- 04-11 walk/stalk: quiet paw alternation, low head focus, minimal vertical bob.
- 12-19 scurry: faster paw cadence, tail counterbalance, no hopping.
- 20-25 eat/sniff: head dips toward paws, body baseline fixed.
- 26-31 stretch/pounce-prep: shoulder lift and small crouch, no airborne leap.
- 32-39 groom: paw-to-face or chest groom, body remains complete.
- 40-47 turn-source: mirrored intermediate posture, no cropped tail.
- 48-55 rest/loaf prep: body lowers slightly without hiding legs entirely.
- 56-61 alert: ears and head perk, tail tip motion.

### `gecko-crawl` frames for gecko variants

- 00-03 idle: still low posture, tiny throat/body pulse.
- 04-11 crawl: alternating feet, belly low, toes never cropped.
- 12-19 quick crawl: longer horizontal reach, no vertical hop.
- 20-25 lick/feed: head/tongue direction implied without text or props.
- 26-31 wall-check/body wave: subtle torso bend, tail counter-curve.
- 32-39 groom/shed rub analog: tiny head/body rub motion, no upright pose.
- 40-47 turn-source: low curved body turn, all limbs visible.
- 48-55 rest: flattened still posture, tail relaxed.
- 56-61 alert: head raise only slightly, feet stay planted.

### `small-rodent-scurry` frames for hamster

- 00-03 idle breathing: compact body pulse, whisker area stable.
- 04-11 walk: short paw steps, small vertical bob.
- 12-19 scurry: faster low run, body remains grounded.
- 20-25 eat: paws near mouth, cheek area readable.
- 26-31 rear/sniff: small lift but not human-like standing.
- 32-39 groom: face and chest grooming.
- 40-47 turn-source: rounded body turn, ears visible.
- 48-55 rest: tucked compact posture.
- 56-61 alert: ears perk and head raise.

## Parent Integration Checklist

- Generate PNGs only under `docs/art-source/cohorts/cohort-04/`.
- Run alpha/bounds checks on every `source-truth` and `motion-source` frame.
- Downscale each accepted frame to 96x64 previews on both light and dark backgrounds.
- Promote only accepted files into the parent importer manifest.
- Keep current catalog `SourceStatus` as `prototype_only` until parent confirms source-truth and 62 frames.
