# Dog Source-Truth Art Review

## Output

- Final transparent source image: `docs/art-intake/dog/dog-source-truth-transparent.png`
- 96x64 transparent preview: `docs/art-intake/dog/dog-source-truth-transparent-96x64-preview.png`
- 96x64 light/dark background check: `docs/art-intake/dog/dog-source-truth-96x64-light-dark-check.png`
- Raw generated image with baked checkerboard background retained for traceability: `docs/art-intake/dog/dog-source-truth-raw-checker.png`

Alpha verification:

- Final image: 1536x1024, `Format32bppArgb`, corner alpha `0,0,0,0`
- Final non-transparent bounding box: `269,91-1496,931`, size `1228x841`
- 96x64 preview non-transparent bounding box: `21,14-74,52`, size `54x39`

## Final Prompt

Create a single source-truth sprite image for a desktop pet.

Subject: one complete small companion dog only, breed-neutral first pass, strict left-facing side view, standing naturally on all four feet. The dog should use only real-world dog colors suitable for normal mode, such as cream, tan, light brown, warm gray-tan, or subtle darker natural ear and tail accents. Do not use fantasy colors. The dog should read clearly as a small dog at Windows taskbar size: clear muzzle, complete ears, four distinct feet, readable gently curved tail, compact body, natural dog proportions. Avoid making it a specific breed in this first pass.

Style: taskbar-readable clean pixel-art / clean illustrated sprite hybrid, crisp silhouette, simple readable forms, limited shading, no painterly blur, no antialias-heavy softness. Designed to later normalize into a 96x64 runtime canvas. Keep a consistent low foot baseline: all four paws touch the same horizontal baseline near the bottom of the animal, with enough transparent margin around ears, nose, tail, and feet for animation frame derivation.

Composition: transparent alpha background, centered animal, full body visible, no cropped body parts. Animal should occupy most of the frame horizontally while leaving safe transparent padding.

Strict exclusions: no text, no border, no scenery, no ground line, no cast shadow, no costume, no collar, no accessories, no human-like pose, no sitting pose, no multiple animals, no props, no speech bubbles, no logo, no watermark.

## Species Readability

The current source image reads as a small companion dog: the side-view muzzle, drop ear, four-foot stance, compact body, and raised tail are clear at source size and remain readable in the 96x64 preview. The cream/tan body with darker natural ears and tail stays within real dog color ranges and is suitable for normal mode.

This first pass is intentionally breed-neutral. It leans toward a generic puppy / small companion-dog silhouette rather than a Shiba, Chihuahua, Dachshund, or Pomeranian variant.

## Real Color And Type Coverage Plan

Use real dog colors only for normal mode:

- White: useful high-contrast silhouette, but needs outline validation on light taskbars.
- Black: common and readable on light backgrounds, but needs internal highlight validation on dark taskbars.
- Brown / chocolate: strong all-purpose natural variant.
- Cream: good default companion-dog color, close to the first-pass source.
- Tan / fawn: good default small-dog color with clear muzzle and ear accents.
- Tricolor: black, tan, and white markings are real and useful, but fine markings may collapse at 96x64.
- Shiba-like: red, sesame, black-and-tan, or cream; needs sharper ears and curled tail, later verification required.
- Chihuahua-like: tan, cream, black, white, chocolate, or mixed coats; breed signal may be weak at 96x64, later verification required.
- Dachshund-like: red, black-and-tan, chocolate-and-tan, cream, or dapple; needs longer low body and short legs, later verification required.
- Pomeranian-like: orange, cream, white, black, or sable; needs fluffy coat silhouette, later verification required.

## 96x64 And Animation Risks

- The current tail is readable but high and curved; leave enough right-side padding when deriving walk and turn frames.
- Foot baseline is usable, but the front feet are close together; walk-cycle frame separation may need hand cleanup.
- Darker ears and tail help readability, but black variants will need extra highlights to avoid losing details on dark backgrounds.
- White and cream variants need a persistent dark outline for light taskbar readability.
- Breed-specific silhouettes may need separate source images instead of recoloring this breed-neutral source.
