# Google Flow Trial Packet - Sugar Glider 8-cell

Date: 2026-06-25

Transmission boundary:

- Public product/spec summary only.
- No repository path, token, private URL, raw log, or local absolute path.
- No reference image upload in the first attempt.
- Review artifact only; do not promote into accepted frames.

Prompt to send:

```text
Create a single horizontal 8-cell sprite motion study for a tiny desktop pet.

Subject: one gray sugar glider / Japanese flying squirrel inspired small animal, right-facing side view, large dark eye, small rounded ears, gray body, pale underside, long tail, and visible patagium side membrane when stretching. Keep it grounded and desktop-pet sized; do not create a flying scene.

Output format: one image containing 8 separate poses arranged left to right in a clean row. Each pose should be easy to crop into a 96x64 runtime frame. Keep the same camera distance, body scale, baseline, lighting, outline softness, and source family across all 8 poses.

Motion plan:
1 idle neutral low stance, calm and balanced
2 idle tiny breathing change, same silhouette
3 idle subtle ear, whisker, or head adjustment
4 idle slight weight shift without changing body volume
5 low skitter start, front paw forward
6 low skitter contact, paws touch the same baseline
7 small membrane stretch while grounded, no flying
8 return phase, long tail follows motion, body stays same size

Style: polished 2D desktop pet sprite, soft natural fur detail, clean readable anatomy, no realism-only photo texture, no toy/plush look, no degu silhouette.

Background: transparent background if available. If transparency is not available, use a perfectly flat pure #00ff00 chroma green background with no gradient, no texture, no shadow, and no floor.

Hard negatives: no text, labels, numbers, borders, dividers, scenery, props, costumes, multiple animals, cropped ears, cropped feet, cropped whiskers, cropped membrane, cropped tail, cropped body, disconnected stray pixels, checkerboard, floor band, cast shadow, white background, gray background, full flight scene, bat wings, human-like pose, giant eye, smeared face, black mask blob, mouth bar, broken paws.
```

Expected local review:

- Save original output under `raw/`.
- If a clean 8-cell row is produced, crop to 8 review PNGs under `candidates/sugar_glider_gray_8cell/`.
- Build a contact sheet under `contact-sheets/`.
- Reject if the face mask becomes a black blob, the membrane turns into bat wings, or any pose leaves the ground plane.
