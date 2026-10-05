# Cohort 07 Report

Status: imagegen_prompt_ready
Owner directory: `docs/art-source/cohorts/cohort-07/`

This cohort covers capybaras, tortoises, and first popular dog additions. No files outside this directory were edited by this cohort worker.

## Output Status

Actual ImageGen PNG output could not be saved directly from this thread into the repository workspace. This report therefore provides the exact prompts, file plan, motion notes, and acceptance checklist for the parent thread or an ImageGen-capable worker to generate and place assets under this cohort directory.

Expected directory structure:

```text
docs/art-source/cohorts/cohort-07/
  source-truth/
    <variant-id>.png
  motion-source/
    <variant-id>/
      00_idle_00.png ... 61_alert_05.png
  preview/
    <variant-id>-96x64-light.png
    <variant-id>-96x64-dark.png
  cohort-report.md
  known-risks.md
```

## Shared Prompt Constraints

Use these constraints for every generated image:

```text
single complete animal, transparent background, no text, no border, no scenery, no cast shadow, no costume, no props, no human pose, full body visible, ears feet whiskers and tail fully inside canvas, consistent side-view camera, consistent scale, clean 2D sprite source art, anatomy faithful to the named species or breed, neutral readable lighting, no duplicate animals
```

Negative prompt:

```text
text, watermark, border, frame, scenery, grass, floor shadow, costume, collar, harness, toy, food, multiple animals, cropped ears, cropped feet, cropped tail, human expression, standing upright like a human, degu recolor, generic rodent body for dog or tortoise
```

## Variant Prompts

### `capybara_brown`

Source-truth path: `source-truth/capybara_brown.png`

Prompt:

```text
Full-body adult capybara, natural brown coat, calm side-view walking-pet sprite pose, large rounded barrel body, blunt muzzle, tiny rounded ears, short legs, very small tail, realistic capybara anatomy, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/capybara_brown/`

- Profile: `capybara-lumber`
- Keep a heavy low body, slow weight-shift, small leg steps, minimal vertical bounce.
- Do not use hop or wheel motion.

### `capybara_sand`

Source-truth path: `source-truth/capybara_sand.png`

Prompt:

```text
Full-body adult capybara, warm sand tan coat, calm side-view walking-pet sprite pose, large rounded barrel body, blunt muzzle, tiny rounded ears, short legs, very small tail, realistic capybara anatomy, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/capybara_sand/`

- Profile: `capybara-lumber`
- Reuse capybara silhouette and contact points from `capybara_brown`; only coat color changes.
- Do not shrink into a guinea pig or hamster silhouette.

### `tortoise_olive`

Source-truth path: `source-truth/tortoise_olive.png`

Prompt:

```text
Full-body tortoise, olive shell and muted green-brown skin, side-view low crawling sprite pose, domed shell, visible head, four sturdy legs, tiny tail fully inside canvas, realistic tortoise proportions, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/tortoise_olive/`

- Profile: `tortoise-plod`
- Use a low plodding gait with almost no body bounce.
- Head can extend/retract slightly; shell baseline must stay stable.
- Do not hop, rear up, or use wheel motion.

### `tortoise_dark_shell`

Source-truth path: `source-truth/tortoise_dark_shell.png`

Prompt:

```text
Full-body tortoise, dark brown shell with muted olive skin, side-view low crawling sprite pose, domed shell, visible head, four sturdy legs, tiny tail fully inside canvas, realistic tortoise proportions, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/tortoise_dark_shell/`

- Profile: `tortoise-plod`
- Same silhouette and contact points as `tortoise_olive`; coat/shell color changes only.
- Avoid turtle-like flippers; keep land tortoise legs.

### `french_bulldog_fawn`

Source-truth path: `source-truth/french_bulldog_fawn.png`

Prompt:

```text
Full-body French Bulldog, fawn coat with dark muzzle, compact muscular body, square head, large bat ears fully visible, short muzzle, short tail, side-view dog sprite pose for a desktop pet, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/french_bulldog_fawn/`

- Profile: `dog-trot`
- Compact trot with short steps and a slight shoulder bounce.
- Keep bat ears stable and uncropped.

### `labrador_yellow`

Source-truth path: `source-truth/labrador_yellow.png`

Prompt:

```text
Full-body Labrador Retriever, yellow coat, athletic medium-large dog body, broad head, floppy ears, otter tail fully visible, side-view dog sprite pose for a desktop pet, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/labrador_yellow/`

- Profile: `dog-trot`
- Smooth even trot, longer stride than French Bulldog.
- Tail follows body line without cropping.

### `labrador_black`

Source-truth path: `source-truth/labrador_black.png`

Prompt:

```text
Full-body Labrador Retriever, black coat with readable highlights, athletic medium-large dog body, broad head, floppy ears, otter tail fully visible, side-view dog sprite pose for a desktop pet, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/labrador_black/`

- Profile: `dog-trot`
- Match `labrador_yellow` silhouette and stride; keep coat readable on dark preview.
- Add subtle edge readability without shadows or outlines baked into the source.

### `golden_retriever_golden`

Source-truth path: `source-truth/golden_retriever_golden.png`

Prompt:

```text
Full-body Golden Retriever, rich golden coat with soft feathering, friendly side-view dog sprite pose, floppy ears, feathered tail fully visible, balanced medium-large retriever body, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/golden_retriever_golden/`

- Profile: `dog-trot`
- Smooth retriever trot with light tail sway.
- Preserve feathering without noisy fur clumps at 96x64.

### `german_shepherd_black_tan`

Source-truth path: `source-truth/german_shepherd_black_tan.png`

Prompt:

```text
Full-body German Shepherd Dog, black and tan coat, alert side-view dog sprite pose, upright ears, sloped athletic body, long bushy tail fully visible, realistic shepherd markings, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/german_shepherd_black_tan/`

- Profile: `dog-trot`
- Athletic trot with moderate stride and stable head carriage.
- Keep upright ears and tail inside frame.

### `dachshund_red`

Source-truth path: `source-truth/dachshund_red.png`

Prompt:

```text
Full-body Dachshund, red coat, long low body, very short legs, long muzzle, floppy ears, tail fully visible, side-view dog sprite pose for a desktop pet, transparent background. Apply shared prompt constraints.
```

Motion source plan: `motion-source/dachshund_red/`

- Profile: `dog-trot`
- Low short-legged trot; body stays close to baseline.
- Do not turn into generic small dog proportions.

## Frame Naming Plan

Use the same 62-frame semantic contract for every variant:

```text
00_idle_00.png
01_idle_01.png
02_idle_02.png
03_idle_03.png
04_walk_00.png ... 11_walk_07.png
12_scurry_00.png ... 19_scurry_07.png
20_eat_00.png ... 25_eat_05.png
26_action_00.png ... 31_action_05.png
32_groom_00.png ... 39_groom_07.png
40_turn_00.png ... 47_turn_07.png
48_rest_00.png ... 55_rest_07.png
56_alert_00.png ... 61_alert_05.png
```

For capybara/tortoise, `scurry` should be interpreted as faster lumber/plod, not running. For tortoise, `groom` should be a head/leg adjustment, not degu grooming. For dogs, `action` can be a natural sniff/play-bow-like weight shift, not a wheel action.

## Acceptance Checklist

- Source-truth PNG exists for all 10 variants.
- Each variant has exactly 62 motion-source PNGs.
- Every PNG has transparent background and one complete animal.
- All frames fit 96x64 reduction without cropped ears, feet, whiskers, or tails.
- Capybara body mass and tortoise low shell profile remain stable across frames.
- Dog breeds are visually distinguishable by silhouette, not just recolor.
- Light and dark 96x64 previews exist for all 10 variants.
- Parent thread reviews before promotion from `prototype_only`.
