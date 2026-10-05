# ImageGen Cohort Contract

AnimalsDesktop v0.1.0 started with 10 ImageGen cohorts with 10 selectable variants each. That broad split is useful for manifest coverage, but the next asset-production pass should use fewer variants per thread and many more images per variant. Each cohort owns only its own directory:

`docs/art-source/cohorts/cohort-XX/`

Parent-thread integration promotes accepted cohort output into the importer manifest and regenerates `assets/sprites`. Cohort workers must not edit `internal/catalog`, `assets/sprites`, release workflows, or files outside their assigned cohort directory.

## Required Cohort Output

For each variant:

- `source-truth/<variant-id>.png`: one complete animal on transparent background.
- `motion-source/<variant-id>/`: 62 transparent source frames preserving camera, scale, anatomy, baseline, and contact points.
- `preview/<variant-id>-96x64-light.png`: 96x64 readability preview on light background.
- `preview/<variant-id>-96x64-dark.png`: 96x64 readability preview on dark background.

Motion sources must be generated as one pose per PNG. Do not create accepted source art by cutting a multi-pose sheet, grid, contact sheet, or checkerboard-background image. If a generator can only return grouped results, keep those files as `review-only` and regenerate the selected poses as separate transparent PNGs before parent integration.

For each cohort:

- `cohort-report.md`: source prompts, acceptance notes, and known risks.
- `known-risks.md`: cropping, anatomy, motion, alpha, scale, or species-identity risks that parent integration must review.

## Next Batch Shape

Use this for new ImageGen work so one thread produces enough alternatives to judge quality:

- Source candidate pass: 4-6 variants per thread, 4 source-truth candidates per variant, plus light/dark 96x64 previews for the best 2 candidates.
- Motion pass: 1 species family per thread, 1-3 accepted variants, 62 transparent motion frames per variant, plus light/dark previews and a contact sheet.
- ImageGen call shape: prefer one pose per request; at most request four related poses only when the tool returns them as four separate PNG files.
- Do not ask one thread to finish 10 variants x 62 frames. That is 620 frames before review and tends to waste RAM and review attention.
- Parent integration promotes only reviewed winners into `assets/sprites` and keeps rejected candidates inside the cohort report.

## Acceptance Rules

- Each source frame contains exactly one animal, no text, no borders, no scenery, no shadows, no costume, and no human-like pose.
- Each accepted source frame has true transparent alpha around the animal. Baked checkerboard, white, gray, or noisy backgrounds are rejected even when the pose is useful as reference.
- Do not make chinchillas, macaroni mice, rabbits, cats, dogs, reptiles, or other species by recoloring degu sprites.
- Coat variants may share a species source family only after the species silhouette and motion set are stable.
- Low-crawler, snake, frog, and tortoise profiles must not use degu wheel or upright hopping actions.
- Parent thread decides when a cohort is promoted from `prototype_only` to an accepted ImageGen source status.

## Cohort Map

- `cohort-01`: degu core coats 1-10.
- `cohort-02`: degu cream pied, chinchillas, macaroni mice, first rabbits.
- `cohort-03`: remaining baseline rabbits, baseline dogs, first cats.
- `cohort-04`: remaining baseline cats, geckos, first hamster.
- `cohort-05`: remaining hamsters, ferrets, guinea pigs.
- `cohort-06`: hedgehogs, squirrels, foxes, red pandas, otter, sugar glider.
- `cohort-07`: capybaras, tortoises, first popular dog additions.
- `cohort-08`: remaining popular dog additions, first popular cat additions.
- `cohort-09`: remaining popular cat additions, first popular rabbit additions.
- `cohort-10`: remaining rabbit addition, small animals, reptile/amphibian additions.

## Near-Term Source Candidate Queues

- `priority-small-pets`: chinchilla, macaroni mouse / fat-tailed gerbil, hamster, gecko.
- `momonga-family`: Japanese dwarf flying squirrel / momonga, sugar glider color variants.
- `small-birds`: sakura buncho / Java sparrow, budgerigar, white buncho, cockatiel, lovebird, zebra finch.
- `popular-companions`: top dog and cat additions already in the 100-variant manifest.
