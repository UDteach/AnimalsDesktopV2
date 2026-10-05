# Cohort 06 Known Risks

## Blocking Risks

- No bitmap ImageGen output was saved by this delegated worker. `source-truth/`, `motion-source/`, and `preview/` PNG targets are still missing.
- Parent integration must not promote cohort-06 variants from `prototype_only` until all required PNGs exist and pass visual review.

## Species Identity Risks

- `hedgehog_salt_pepper` and `hedgehog_cinnamon`: generators may create porcupine-like long spines or crop tiny feet. Reject those frames.
- `squirrel_red` and `squirrel_gray`: tail cropping is likely because the tail is large relative to 96x64. Require full tail visibility or a consistent curled tail design.
- `fox_red` and `fox_silver`: may drift toward dog/wolf proportions. Require fox muzzle, triangular ears, and full tail.
- `red_panda_classic` and `red_panda_dark`: may drift toward raccoon. Require red panda head shape, cream face markings, compact body, and ringed tail.
- `otter_brown`: may become seal-like or upright. Require short legs, tapered tail, and low sliding posture.
- `sugar_glider_gray`: may become bat-like. Require sugar glider body, large eyes, thin tail, and subtle patagium, not wings.

## Motion Risks

- Hedgehog motion must stay low and shuffle-based; no hopping or wheel posture.
- Squirrel motion can bound, but the baseline must remain stable enough for taskbar placement.
- Fox motion should be a trot with small vertical lift, not a degu scurry.
- Red panda motion should be slower and heavier than squirrel and fox.
- Otter motion must be a low slide/crawl and must not enter degu wheel-like actions.
- Sugar glider motion should skitter and membrane-stretch, not flap like a bird or bat.

## Alpha And Framing Risks

- Reject opaque backgrounds, painted shadows, ground patches, scenery, or border artifacts.
- Reject any frame with cropped ears, feet, whiskers, or tails.
- Require consistent 62-frame camera, scale, baseline, and contact points per variant before parent importer promotion.
