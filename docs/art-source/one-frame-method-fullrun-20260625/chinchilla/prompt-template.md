# Chinchilla One-Frame Prompt Template

## Base Species Prompt

Create one small standard gray chinchilla desktop pet sprite source frame for `{frame_label}`: `{pose_description}`.

Use a clean small 2D illustrated chinchilla desktop pet sprite style. The animal must be exactly one complete standard gray chinchilla with one head, one face, rounded ears, compact gray fluffy body, tiny feet, and a readable chinchilla tail.

Use a perfectly flat `#00ff00` chroma-key background only. Center the full body with padding and keep the feet on a stable shared baseline.

## Defect Controls

The original chinchilla draft was rejected because visual QA found extra debris and double-face or face-like duplicate anatomy. Every retry in this lane used these hard controls:

- no duplicate head
- no second face
- no face-like body markings
- no extra eyes or extra nose outside the single face
- no detached paws or ears
- no floating fur, debris, loose fragments, dust, dots, or stray marks
- no props, ground, floor object, shadow, text, watermark, checkerboard, gradient, or texture
- no internal green holes in the animal silhouette

When `prepareframe` rejected a candidate for transparent pinholes, the retry prompt made the tail and feet more connected:

- tail short, thick, fluffy, and pressed directly against the body
- no curl loop, no tail/body gap, and no leg/belly gap
- feet tucked close to the body when the pose does not require a stride

## Reproducibility Note

Built-in ImageGen does not provide deterministic seeds here. Reproducibility is operational: frame numbers, pose descriptions, raw generated images, prepared `96x64` PNGs, QA reports, and accepted paths are all recorded in `manifest.csv`.
