# Four New Animal Frame-00 Seed Run - 2026-06-26

## Scope

This run used the successful one-frame ImageGen workflow to create first-pass `frame-00` seed candidates for four next-priority real animals:

- guinea pig / tricolor
- fancy rat / hooded
- gray-brown gecko
- albino chipmunk

The seed outputs in this top-level run are review candidates only. Nothing from the four `frame-00` seed candidates was promoted directly to `accepted-frames`, runtime sprites, catalog runtime variants, or release assets.

A later follow-up under `gecko-set00-oneframe-62/` used the gecko seed direction to generate a full 62-frame gray-brown gecko `set00`. That follow-up was separately normalized, QA-reviewed, and promoted to `docs/art-source/gecko/motion-source/accepted-frames/set00/` as accepted source art.

## Method

Each animal was generated as one standalone right-facing side-view animal on a flat `#00ff00` chroma-key background. Raw outputs were copied into this repo, chroma-keyed to alpha with the system ImageGen helper, then fit into 96x64 review frames using `tools/fit_review_frame.py`.

`cmd/prepareframe` was also run as the strict promotion gate. All candidates failed strict mode because the photoreal/fur-detail-heavy matte contains small transparent pinholes. The failures are recorded in `strict-prepareframe-results.tsv`.

## Output Overview

```text
docs/art-source/one-frame-method-fullrun-20260626/four-new-animal-frame00-contact-overview.png
```

## Visual Verdict

| Family | Verdict | Notes |
| --- | --- | --- |
| Guinea pig | seed candidate | Strong species read and complete body. Too photoreal/fur-detail-heavy for accepted promotion without a more sprite-like retry. |
| Fancy rat | seed candidate | Long tail and hooded marking read well. Needs stronger outline for light-background readability. |
| Gecko | best seed candidate | Strongest output in this run. Low reptile body, feet, and tail survive the 96x64 contact review. |
| Albino chipmunk attempt 01 | rejected | Large curled tail makes it read as a white squirrel. |
| Albino chipmunk attempt 02 | needs review | Better low tail and compact pose, but stripe/read is still weak; likely needs a more stylized retry. |

## Next Prompt Direction

For accepted-frame production, the next prompts should reduce pinhole-prone detail:

- ask for cleaner sprite-like forms with fewer individual fur strands
- require a slightly stronger dark outline, especially for white/light animals
- avoid separated hair detail around feet, belly, tail, and whiskers
- for chipmunk, emphasize compact ground-chipmunk body, cheek stripe hints, and low slim tail

## QA

Local checks completed:

- visual review of light/dark/checker contact sheets
- `cmd/prepareframe` strict gate attempted for every candidate
- no accepted-frame promotion performed
