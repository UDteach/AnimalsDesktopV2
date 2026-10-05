# ImageGen-Only Candidate Review - 2026-06-25

## Scope

User direction changed to ImageGen-only during the run. ChatGPT/Gemini consultation was not transmitted and is not part of this task.

Generated candidates are review artifacts only. Nothing in this pass was promoted into `accepted-frames`.

## Outputs

- Pass1 combined sheet: `docs/art-source/five-small-animals-imagegen-only-frame00-contact.png`
- Pass2 combined sheet: `docs/art-source/five-small-animals-imagegen-only-frame00-contact-pass2.png`
- Degu comparison sheet: `docs/art-source/five-small-animals-vs-degu-frame00-contact.png`
- Per-family pass2 raw/alpha/frame/contact outputs:
  - `docs/art-source/chinchilla/motion-source/review/imagegen-only-20260625-pass2/`
  - `docs/art-source/hamster/motion-source/review/imagegen-only-20260625-pass2/`
  - `docs/art-source/macaroni-mouse/motion-source/review/imagegen-only-20260625-pass2/`
  - `docs/art-source/rabbit/motion-source/review/imagegen-only-20260625-pass2/`
  - `docs/art-source/sugar-glider/motion-source/review/imagegen-only-20260625-pass2/`

## Mechanical Checks

All pass2 `frame00-review-96x64.png` candidates:

- are 96x64 PNGs
- have non-empty alpha
- have no alpha touching the outer canvas edge
- were reviewed on light, dark, and checker backgrounds

Rabbit also has a safer margin placement:

- `docs/art-source/rabbit/motion-source/review/imagegen-only-20260625-pass2/rabbit-frame00-review-96x64-safe.png`
- bbox after safe placement: `(14, 5, 81, 59)`, edge alpha pixels: `0`

## Visual Decision

| Family | Pass2 status | Notes |
| --- | --- | --- |
| Chinchilla | usable review anchor | Better than pass1; still softer than DeguDesktop but readable. |
| Hamster | usable review anchor | Good species read; keep short/no tail and rounded cheek shape in later frames. |
| Macaroni mouse | usable review anchor | Thick tail reads clearly; avoid letting later frames become hamster-like. |
| Rabbit | usable with safe placement | Use the `safe` 96x64 placement; ears need protected top margin in all future prompts. |
| Momonga / sugar glider | usable review anchor | Stronger than existing generated source; keep membrane and grounded stance. |

## Next ImageGen Loop

Generate `set00/frame-01` through `frame-03` for each family first. These should be subtle idle/breathing variants, not new source families.

Prompt rule for the next loop:

- Reference the pass2 candidate as the visual baseline in words.
- Keep identical camera, body scale, baseline, source family, and right-facing side view.
- Ask for small pose deltas only: breathing, ear adjustment, whisker/head shift, and slight weight shift.
- Keep chroma-green background and route through local chroma removal before review.

Do not start walk/scurry frames until the first four idle frames are visually consistent inside each family.
