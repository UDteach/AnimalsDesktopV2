# Google Flow Trial Packet - Macaroni Mouse / Fat-tailed Gerbil 8-cell

Date: 2026-06-25

Transmission boundary:

- Public product/spec summary only.
- No repository path, token, private URL, raw log, or local absolute path.
- No reference image upload in the first attempt.
- Review artifact only; do not promote into accepted frames.

Prompt to send:

```text
Create a single horizontal 8-cell sprite motion study for a tiny desktop pet.

Subject: one tan fat-tailed gerbil / macaroni mouse, right-facing side view, compact low body, small rounded ears, large dark eye, tiny paws, and a thick pale tail. The tail must be thick and soft, not a thin fancy-mouse tail. Do not make it look like a hamster.

Output format: one image containing 8 separate poses arranged left to right in a clean row. Each pose should be easy to crop into a 96x64 runtime frame. Keep the same camera distance, body scale, baseline, lighting, outline softness, and source family across all 8 poses.

Motion plan:
1 idle neutral stance, calm and balanced
2 idle tiny breathing change, same silhouette
3 idle subtle head or whisker adjustment
4 idle slight weight shift without changing body volume
5 quick crouch prep, compact body
6 low scurry contact, forefoot touches the same baseline
7 scurry alternate phase, thick tail follows motion
8 scurry return phase, body stays same size

Style: polished 2D desktop pet sprite, soft natural fur detail, clean readable anatomy, no realism-only photo texture, no toy/plush look, no degu silhouette.

Background: transparent background if available. If transparency is not available, use a perfectly flat pure #00ff00 chroma green background with no gradient, no texture, no shadow, and no floor.

Hard negatives: no text, labels, numbers, borders, dividers, scenery, props, costumes, multiple animals, cropped ears, cropped feet, cropped whiskers, cropped tail, cropped body, disconnected stray pixels, checkerboard, floor band, cast shadow, white background, gray background, human-like pose, giant eye, smeared face, black mask, mouth bar, broken paws, thin mouse tail, hamster silhouette.
```

Expected local review:

- Save original output under `raw/`.
- If a clean 8-cell row is produced, crop to 8 review PNGs under `candidates/macaroni_mouse_tan_8cell/`.
- Build a contact sheet under `contact-sheets/`.
- Reject if the tail becomes thin, the body becomes hamster-like, or the baseline/scale changes between cells.
