# Google Flow Trial Packet - Rabbit 8-cell

Date: 2026-06-25

Transmission boundary:

- Public product/spec summary only.
- No repository path, token, private URL, raw log, or local absolute path.
- No reference image upload in the first attempt.
- Review artifact only; do not promote into accepted frames.

Prompt to send:

```text
Create a single horizontal 8-cell sprite motion study for a tiny desktop pet.

Subject: one chestnut agouti rabbit, right-facing side view, long ears fully visible with safe top margin, compact body, strong hindquarters, small tail, visible whiskers, and small front paws. Keep the rabbit cute but anatomically coherent.

Output format: one image containing 8 separate poses arranged left to right in a clean row. Each pose should be easy to crop into a 96x64 runtime frame. Keep the same camera distance, body scale, baseline, lighting, outline softness, and source family across all 8 poses.

Motion plan:
1 idle neutral crouched stance, calm and balanced
2 idle tiny breathing change, same silhouette
3 idle subtle ear or whisker/head adjustment, ears fully visible
4 idle slight weight shift without changing body volume
5 hop prep, hind legs compress
6 hop contact, front paw touches the same baseline
7 tiny hop peak, very small lift only, ears still inside canvas
8 hop recovery, body stays same size

Style: polished 2D desktop pet sprite, soft natural fur detail, clean readable anatomy, no realism-only photo texture, no toy/plush look, no degu silhouette.

Background: transparent background if available. If transparency is not available, use a perfectly flat pure #00ff00 chroma green background with no gradient, no texture, no shadow, and no floor.

Hard negatives: no text, labels, numbers, borders, dividers, scenery, props, costumes, multiple animals, cropped ears, cropped feet, cropped whiskers, cropped tail, cropped body, disconnected stray pixels, checkerboard, floor band, cast shadow, white background, gray background, oversized leap, human-like pose, giant eye, smeared face, black mask, mouth bar, broken paws.
```

Expected local review:

- Save original output under `raw/`.
- If a clean 8-cell row is produced, crop to 8 review PNGs under `candidates/rabbit_chestnut_agouti_8cell/`.
- Build a contact sheet under `contact-sheets/`.
- Reject if ears are cropped, hop height changes canvas needs, or hindquarters shrink between cells.
