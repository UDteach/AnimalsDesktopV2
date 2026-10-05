# Macaroni Mouse Source-Truth Image

Status: accepted source seed

Canonical assets:

- `macaroni-mouse-source-truth.png`
- `macaroni-mouse-96x64-preview.png`
- `macaroni-mouse-96x64-on-light.png`
- `macaroni-mouse-96x64-on-dark.png`

Compatibility copies also exist one directory up:

- `../macaroni-mouse-source.png`
- `../macaroni-mouse-96x64-preview.png`

## Art Contract Check

- Transparent background: yes, verified as 32-bit ARGB PNG.
- Animal count: one complete animal only.
- View: strict right-facing side view.
- Style: clean illustrated sprite with simplified pixel-readable silhouette.
- Baseline: low body/foot baseline suitable for later animation normalization.
- Anatomy: complete ear, feet, whiskers, and tail are present.
- Exclusions: no text, border, scenery, cast shadow, costume, human-like pose, multiple animals, or cropped body parts.
- Runtime target: checked with a 96x64 transparent preview.
- Light/dark readability: checked with 96x64 light-background and dark-background QA previews. The body, eye, and fat tail remain readable; whiskers and small feet are the first details to degrade.

## Species Readability Review

The image is intended to read as a macaroni mouse / fat-tailed gerbil, not a generic mouse. The most important silhouette cue is the short, thick, club-shaped tail. The body is short and rounded with a nearly neckless low posture, a pointed muzzle, a large dark eye, small rounded ear, small hands, and small feet.

The representative color is a real-world direction: tan to gray-tan dorsal fur, pale/white underside and feet, dark eye, pinkish sparse-haired tail and ear. This follows fat-tailed gerbil references describing gray-to-tan fur, sometimes black-tipped dorsal hairs, white underside/feet, low-positioned pink ears, and a short thick club-shaped tail.

Primary references checked:

- Animal Diversity Web, `Pachyuromys duprasi`: https://animaldiversity.org/accounts/Pachyuromys_duprasi/
- DNA Zoo, `Pachyuromys duprasi`: https://www.dnazoo.org/assemblies/pachyuromys_duprasi
- DNA Zoo post, "A camel of a gerbil": https://www.dnazoo.org/post/a-camel-of-a-gerbil

## Final ImageGen Prompt

Use this if a regenerated source frame is needed:

```text
Transparent background PNG source sprite. One complete macaroni mouse / fat-tailed gerbil (Pachyuromys duprasi) only, strict right-facing side profile, orthographic view, centered with generous transparent margin. Clean illustrated sprite / taskbar-readable pixel-art hybrid, crisp silhouette, no scenery. Make it the smallest candidate animal: compact low posture, short rounded body, almost no neck, pointed muzzle, large glossy dark oval eye, small low rounded pink ear, fine whiskers, tiny front hands, tiny hind feet, all feet resting on one consistent low horizontal baseline for later animation.

Use real fat-tailed gerbil coloration only: tan to gray-tan sandy dorsal fur with a few subtle darker black-tipped dorsal hair accents, pale cream to white underside and white feet, pinkish sparsely haired ears and tail, dark eye. No fantasy colors, no fancy-mouse coat colors, no invented markings.

The tail is the key species signal: short, thick, plump, rounded club-shaped tail like a macaroni/fat-storage tail, attached at the rear and fully visible, clearly not a normal thin mouse tail. Complete ears, feet, whiskers, and tail. No text, no border, no cast shadow, no ground line, no scenery, no costume, no props, no human-like pose, no multiple animals, no cropped body parts.
```

## Real-World Variant Coverage Proposal

Promote only after reference checks and taskbar-size readability checks:

- `tan_gray_tan_white_belly`: representative seed, accepted.
- `sand_yellow_white_belly`: likely real-world direction, needs_review for exact label.
- `gray_tan_black_tipped_dorsal`: supported as a real-world direction, needs_review for visual distinction at 96x64.
- `lighter_sandy_brown_white_belly`: likely real-world direction, needs_review for exact label.

Do not promote until verified:

- Fancy mouse colors such as blue, chocolate, piebald, satin, merle, or dalmatian-like labels.
- Pure black, pure white, or high-contrast pet-morph labels unless a reliable fat-tailed gerbil source supports them.
- Subspecies-specific labels, because the current art pass did not verify that the visual differences are reliable or readable at taskbar size.

## 96x64 And Animation Risks

- Tail cropping: medium risk. The fat tail is essential and extends the sprite width; runtime frame boxes should reserve rear margin.
- Foot readability: medium risk. The feet are intentionally small and may need thicker runtime sprite pixels.
- Whisker readability: medium risk. Whiskers are visible in the source but may need simplified 1-2 pixel lines at runtime.
- Silhouette confusion: low-to-medium risk. The thick club-shaped tail separates it from a generic mouse; without the tail, it can read hamster-like.
- Baseline: low risk. The current source and preview keep a consistent low foot/belly baseline.
