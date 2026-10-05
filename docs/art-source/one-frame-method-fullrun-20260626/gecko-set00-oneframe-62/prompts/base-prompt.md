# Base Prompt

Use case: stylized-concept
Asset type: 96x64 transparent desktop-pet sprite frame, chroma-key source
Primary request: create one small gray-brown gecko frame for AnimalsDesktop `set00`.
Reference image: use the visible gecko seed only as a species/silhouette/style anchor; make a cleaner flat sprite version, not a photoreal copy.

Style:
- clean flat 2D desktop-pet sprite, right-facing side-view unless the pose table says this is a turn phase
- simple readable body shapes, smooth pixel-art-adjacent illustration, crisp outline, no realistic tiny scale texture
- gray-brown gecko body, cream underside, dark bead eye, tiny toes suggested as connected rounded pads
- low reptile body, four splayed feet, full tail inside canvas

Canvas/background:
- single centered animal only, generous padding
- perfectly flat solid #00ff00 chroma-key background
- no floor, no shadow, no contact shadow, no gradient, no texture, no props, no text, no watermark
- do not use #00ff00 anywhere in the animal

Hard rejects:
- duplicate face, extra head, extra animal, cropped tail or feet, detached toe/limb debris, floating specks, holes inside the animal, scenery, borders, UI, labels, checker background
- hair-like or fur-like detail; noisy photoreal scales; transparent pinholes; separated tiny toes that make green holes

Per-frame line to append:

```text
Frame NN of 62. Motion slot: <slot>. Pose: <pose_notes>. Keep the same gecko scale, palette, outline thickness, body proportions, and baseline as the other frames.
```
