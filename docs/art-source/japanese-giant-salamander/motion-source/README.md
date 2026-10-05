# Japanese Giant Salamander Motion Source

Source-art accepted set for `japanese_giant_salamander` / オオサンショウウオ.

- Accepted frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Source sheet: `sheets/japanese-giant-salamander-source-set00.png`
- Full contact: `contacts/japanese-giant-salamander-set00-full-contact.png`
- Preview GIF: `contacts/japanese-giant-salamander-set00-preview.gif`

Visual contract:

- Brown mottled Japanese giant salamander.
- Long, low, flattened body with broad flat head, tiny eye, short sturdy limbs,
  and thick tail.
- Motion is intentionally subtle and ground-hugging.
- No scenery, floor, shadows, props, duplicate animals, or cropped head, feet,
  or tail.

QA summary:

- `oneframe_run.py review`: `frame_count=62`, `issues=[]`
- Source-run `auditframes`: `valid=62 missing=0 invalid=0 warnings=44`
- Accepted-source `auditframes`: `valid=62 missing=0 invalid=0 warnings=44`
- `assemblemotion`: wrote `sheets/japanese-giant-salamander-source-set00.png`
- `validatemotion`: deferred until the runtime catalog registers
  `japanese_giant_salamander`.
