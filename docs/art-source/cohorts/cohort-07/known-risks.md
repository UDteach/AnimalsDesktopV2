# Cohort 07 Known Risks

## Global

- Actual ImageGen PNGs were not generated or saved by this thread; this cohort is prompt-ready only.
- Parent integration must reject any generated frame with background, shadow, text, border, costume, props, or multiple animals.
- Parent integration must compare all frames at 96x64 because large-breed dog features can collapse into generic silhouettes.

## Capybara

- Risk: generated capybaras may look like guinea pigs, beavers, or generic rodents.
- Mitigation: require large barrel body, blunt muzzle, tiny round ears, short legs, and almost no tail.
- Motion risk: too much hop or bounce would feel unnatural.
- Mitigation: use slow lumber with stable baseline and small vertical movement.

## Tortoise

- Risk: generated tortoises may become aquatic turtles with flippers or may crop shell/head.
- Mitigation: require land tortoise legs, domed shell, low plod, head visible, and shell fully inside frame.
- Motion risk: hopping or rearing would violate low-crawler behavior.
- Mitigation: use plod, head extension, and leg stepping only.

## French Bulldog

- Risk: bat ears and compact head may crop or merge at 96x64.
- Mitigation: keep ears fully inside frame, side-view head large enough to read, and body compact.

## Labrador Retriever

- Risk: yellow and black Labrador variants can become generic dogs without retriever proportions.
- Mitigation: require broad head, floppy ears, otter tail, and athletic retriever body.
- Dark-preview risk: black coat may lose readability.
- Mitigation: parent should inspect dark preview and regenerate with readable highlights if needed.

## Golden Retriever

- Risk: fur feathering may create noisy edges after downscale.
- Mitigation: use clean sprite-friendly feathering and inspect 96x64 previews.

## German Shepherd Dog

- Risk: black/tan markings may drift between frames or ears may crop.
- Mitigation: lock upright ears, saddle marking placement, and tail extent across the 62-frame set.

## Dachshund

- Risk: generated dog may become a generic small dog rather than a long low dachshund.
- Mitigation: require long body, very short legs, long muzzle, floppy ears, and low baseline.
