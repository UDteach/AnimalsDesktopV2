# Parent Thread Prompt Improvements

Use these as short additions to the next ImageGen prompts. They are intentionally prompt-only and do not imply accepted-frame promotion.

## Global Additions

```text
Create one single right-facing side-view 96x64 desktop pet sprite frame of [SPECIES], one complete animal, transparent background.
Match the frame-00 family anchor: same camera distance, body volume, contact baseline, head-to-body ratio, eye size, muzzle shape, foot scale, lighting, and soft outline.
Only change this one pose detail: [SMALL FRAME DELTA].
Keep all anatomy fully inside the canvas: ears, feet, whiskers, tail, and species-specific features.
No scene, prop, food item, wheel, floor band, cast shadow, text, border, duplicated animal, cropped anatomy, or human-like pose.
```

## Frames 01-03 Consistency Gate

```text
Frame 01: tiny breathing expansion only; body remains level on the same baseline.
Frame 02: subtle ear, whisker, or head adjustment only; same silhouette and scale.
Frame 03: tiny weight shift through feet and shoulders; no new source-family look.
```

## Per-Species Guards

```text
hamster_golden_syrian:
Preserve the golden Syrian hamster cheek-forward compact silhouette, blunt muzzle, rounded ears, tiny paws, and very short or hidden tail. Do not make it narrow-bodied, long-tailed, skinny, or sharp-muzzled.

macaroni_mouse_tan:
Preserve the compact tan fat-tailed gerbil silhouette with a thick pale tail that stays visible and attached. Do not make a thin tail, round hamster body, or pointed long-nose face.

sugar_glider_gray:
Preserve the gray sugar glider silhouette with large dark eyes, long tail, small rounded ears, and a visible patagium line. Keep feet near the shared baseline. Do not create a flying scene, spread-eagle poster pose, over-broad facial mask, or oversized black face patch.

rabbit_chestnut_agouti:
Preserve the chestnut agouti rabbit silhouette with long ears fully inside the canvas, strong hindquarters, compact body, small tail, and visible whiskers. Do not crop ears, over-stretch the jump, or turn the hop into a large airborne leap.
```

## First High-Risk Motion Recipes

```text
hamster_golden_syrian frame 04-11 walk:
Use low, tiny-foot walking phases. The cheek mass and blunt muzzle stay stable while one paw changes phase. No wheel, no big stride, no body stretch.

macaroni_mouse_tan frame 12-19 scurry:
Use low skitter phases with the thick tail counterbalancing behind the body. The tail remains thick, pale, attached, and inside the canvas.

sugar_glider_gray frame 26-31 species action:
Use a grounded patagium stretch: membrane line becomes slightly more visible while all feet remain near the baseline. No full flight, no branch, no scenery.

rabbit_chestnut_agouti frame 26-31 species action:
Use compact hop prep, tiny lift, landing, and recovery. Keep ears fully inside the frame with top margin and avoid a large airborne arc.
```
