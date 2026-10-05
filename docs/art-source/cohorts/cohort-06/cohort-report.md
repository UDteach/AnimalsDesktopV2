# Cohort 06 Report

Status: imagegen_prompt_pack_ready
Owner boundary: `docs/art-source/cohorts/cohort-06/`

## Scope

This cohort owns 10 variants:

| Variant | Species / coat | Motion profile | Source-truth target |
| --- | --- | --- | --- |
| `hedgehog_salt_pepper` | Hedgehog, salt and pepper | `hedgehog-shuffle` | `source-truth/hedgehog_salt_pepper.png` |
| `hedgehog_cinnamon` | Hedgehog, cinnamon | `hedgehog-shuffle` | `source-truth/hedgehog_cinnamon.png` |
| `squirrel_red` | Squirrel, red | `squirrel-bound` | `source-truth/squirrel_red.png` |
| `squirrel_gray` | Squirrel, gray | `squirrel-bound` | `source-truth/squirrel_gray.png` |
| `fox_red` | Fox, red | `fox-trot` | `source-truth/fox_red.png` |
| `fox_silver` | Fox, silver | `fox-trot` | `source-truth/fox_silver.png` |
| `red_panda_classic` | Red panda, classic | `red-panda-amble` | `source-truth/red_panda_classic.png` |
| `red_panda_dark` | Red panda, dark | `red-panda-amble` | `source-truth/red_panda_dark.png` |
| `otter_brown` | Otter, brown | `otter-slide` | `source-truth/otter_brown.png` |
| `sugar_glider_gray` | Sugar glider, gray | `sugar-glider-skitter` | `source-truth/sugar_glider_gray.png` |

## ImageGen Availability

This delegated worker could not save ImageGen bitmap output into the repository. Per the parent contract, no PNG files were fabricated. The cohort is left as a ready-to-run prompt pack with deterministic file targets, frame sequencing, acceptance criteria, and known risks.

## Shared Source-Truth Prompt Prefix

Use this exact prefix for each `source-truth/<variant-id>.png` request:

```text
Pixel-art desktop pet sprite source, one complete animal only, transparent background, no text, no border, no scenery, no cast shadow, no props, no costume, no human-like pose. Side-view three-quarter right-facing animal, calm natural posture, full body visible including ears, feet, whiskers, and tail when present. Consistent camera, consistent scale, clean readable silhouette, suitable for reducing to a 96x64 Windows taskbar pet sprite, natural anatomy, soft hand-painted pixel-art edges, transparent alpha.
```

## Variant Source Prompts

### `hedgehog_salt_pepper`

Append:

```text
Species: African pygmy hedgehog. Coat: salt and pepper quills with pale face and belly. Body low and compact, rounded back, tiny feet visible, small dark eyes, short snout. Avoid porcupine-like long spines; quills should read as short dense texture.
```

### `hedgehog_cinnamon`

Append:

```text
Species: African pygmy hedgehog. Coat: warm cinnamon quills with cream face and belly. Body low and compact, rounded back, tiny feet visible, small dark eyes, short snout. Avoid orange fox coloration; keep a hedgehog silhouette.
```

### `squirrel_red`

Append:

```text
Species: red squirrel. Coat: warm red-brown with cream belly. Long curved bushy tail fully visible behind body, pointed ears, small paws, alert but natural quadruped stance. Avoid chipmunk stripes and avoid standing upright like a mascot.
```

### `squirrel_gray`

Append:

```text
Species: gray squirrel. Coat: cool gray with pale belly. Long curved bushy tail fully visible behind body, pointed ears, small paws, alert but natural quadruped stance. Avoid red tones and avoid standing upright like a mascot.
```

### `fox_red`

Append:

```text
Species: red fox. Coat: orange-red body, white chest and tail tip, black lower legs. Long muzzle, triangular ears, long bushy tail fully visible, natural small canine trot-ready stance. Avoid dog puppy proportions.
```

### `fox_silver`

Append:

```text
Species: silver fox. Coat: dark charcoal and silver guard hairs, pale tail tip, black lower legs. Long muzzle, triangular ears, long bushy tail fully visible, natural small canine trot-ready stance. Avoid wolf proportions and avoid black cat silhouette.
```

### `red_panda_classic`

Append:

```text
Species: red panda. Coat: rusty red body, cream face markings, dark legs, ringed tail fully visible. Rounded ears, compact arboreal body, gentle ambling quadruped stance. Avoid raccoon-only face; keep red panda head shape and tail rings.
```

### `red_panda_dark`

Append:

```text
Species: red panda. Coat: darker rust-brown body, cream face markings, dark legs, ringed tail fully visible. Rounded ears, compact arboreal body, gentle ambling quadruped stance. Keep red panda identity clear even with darker coat.
```

### `otter_brown`

Append:

```text
Species: small-clawed otter. Coat: smooth brown with lighter muzzle and throat. Long low body, short legs, tapered tail trailing behind, rounded head, whiskers visible. Natural sliding low-crawler stance, not upright.
```

### `sugar_glider_gray`

Append:

```text
Species: sugar glider. Coat: soft gray with dark dorsal stripe, pale belly, large dark eyes, small rounded ears, long thin tail. Patagium membrane subtly visible between front and hind legs in a natural skitter-ready stance. Avoid bat wings.
```

## Motion Source Frame Plan

Create 62 transparent frames under `motion-source/<variant-id>/frame-00.png` through `frame-61.png`. Keep each variant's frame set at the same camera, scale, baseline, and contact points.

| Frame range | Motion use | Hedgehog | Squirrel | Fox | Red panda | Otter | Sugar glider |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 00-03 | idle | subtle breathing, quills stable | tail micro sway | ear/tail micro sway | slow breath | head/whisker micro sway | large-eye blink, tail micro sway |
| 04-11 | walk | shuffle steps, low bob | bound-ready steps | trot cycle | slow amble | low slide-crawl | quick skitter |
| 12-19 | fast | brisk shuffle | two-beat bound | quicker trot | cautious quick amble | sliding dash, no hop | rapid skitter |
| 20-25 | eat | snout dip | paws nibble | head dip | paw nibble | muzzle dip | paw nibble |
| 26-31 | species action | curl start, then uncurl | short bound | playful pounce-ready crouch | tail balance step | belly slide | membrane stretch |
| 32-39 | groom | face paw rub | face paw rub | cheek/ear rub | face rub | whisker rub | face grooming |
| 40-47 | turn | compact turn | tail-led turn | tail-led turn | tail-led turn | body ripple turn | tail-led turn |
| 48-55 | rest | settle low | crouch with tail | lie-low pause | slow settle | flattened rest | tucked rest |
| 56-61 | alert | nose twitch | ear/tail perk | ear perk | head lift | head lift | blink and ear perk |

## Preview Outputs

For every accepted source-truth image, export:

- `preview/<variant-id>-96x64-light.png`
- `preview/<variant-id>-96x64-dark.png`

Both previews must preserve readable full-body silhouette at 96x64.

## Parent Integration Notes

- Keep all accepted files in this cohort directory until the parent imports them.
- Parent should reject any frame where tails, ears, feet, or whiskers are cropped.
- Parent should keep `SourceStatus` as `prototype_only` until all required PNGs exist and pass alpha/preview review.
