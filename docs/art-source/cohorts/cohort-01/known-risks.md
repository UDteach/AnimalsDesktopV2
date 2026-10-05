# Cohort 01 Known Risks

Status: assets_not_generated

## Blocking Risk

- Real ImageGen PNG output was not available in this thread. No `source-truth`, `motion-source`, or `preview` PNG files were created.

## Variant Risks

- `black`: dark body may lose anatomy on dark preview unless subtle charcoal highlights are added.
- `white_cream`: pale body may lose contour on light preview unless internal shading is strong enough.
- `black_pied`, `agouti_pied`, `blue_pied`: patch placement can drift between frames; parent review must compare all 62 frames side by side.
- `blue` and `gray`: colors may become too similar; parent review should compare swatches and small 96x64 previews.
- All variants: tail, whiskers, and feet are high-crop-risk details at 96x64.

## Motion Risks

- Wheel-capable run frames must remain degu-like and four-footed, not upright or cartoon treadmill posing.
- Forage/eat frames can accidentally introduce props or food; reject frames with visible props if they violate the no-props rule.
- Groom frames can become human-like if forepaws are posed too high; keep grooming compact and natural.
- Turn frames can shift baseline or scale; parent should verify contact points and silhouette continuity.

## Review Before Promotion

- Confirm all required filenames exist.
- Confirm every frame has non-empty alpha and transparent background.
- Confirm no frame contains scenery, shadows, text, borders, props, costume, or extra animals.
- Confirm frame-to-frame scale and baseline do not jump.
- Confirm source art remains anatomically degu-like, especially ears, tail, whiskers, and hind feet.
