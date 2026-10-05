# Animal Variant Version Requirements

Last updated: 2026-06-21

## Goal

Create a multi-animal version of Degu Desktop where the existing desktop-pet behavior stays intact and the user can choose different animal art sets.

The first candidate animals are:

- Chinchilla
- Macaroni mouse
- Rabbit
- Dog
- Cat
- Gecko
- Hamster

The design must allow more animals to be added later without hard-coding each animal into unrelated behavior, settings, or release logic.

## Source Baseline

Use `D:\開発\DeguDesktop` as the behavior baseline. Do not modify that repository from this management task unless a later implementation task explicitly targets it.

Baseline behavior to preserve:

- Transparent always-on-top Windows taskbar overlay.
- System tray menu.
- Japanese and English settings UI.
- Pet count selection: 1, 2, 3, 5, 10.
- Keyboard reaction mode.
- Random stroll mode.
- Typing wheel behavior where supported by the animal art.
- Foraging, eating, digging, carrying, grooming, social pause, and hover name labels.
- Update checking and release artifact flow.

## Product Scope

The first deliverable is an art-selectable build, not separate apps per animal.

Required concepts:

- Species: the animal family, such as `degu`, `chinchilla`, `cat`, or `gecko`.
- Coat or variant: a visual variation within a species, such as color, breed, morph, or pattern.
- Behavior profile: motion weights and optional actions that fit the species.
- Render profile: frame size, baseline, scale, and optional prop compatibility.
- Asset source family: ImageGen source images and derived runtime sprites for one species.

Species and coat must be separate. Chinchilla, macaroni mouse, rabbit, dog, cat, gecko, and hamster must not be implemented as recolored degu sprites.

Color, coat, breed, and morph variants must be based on real animals. Do not invent fantasy colors for normal mode. For each species, cover several familiar real-world appearances over time rather than a single default look. Before accepting a variant set, check that the label and approximate color direction match real-world references for that animal.

## Real-World Variant Coverage

Use the first image as a readable representative seed. Use later images to cover real color, breed, coat, or morph diversity.

Initial coverage targets:

| Species | Real-world variant directions to consider |
| --- | --- |
| Chinchilla | Standard gray, black velvet, beige, white, ebony, violet, sapphire |
| Macaroni mouse | Real observed color directions only; cross-check against fat-tailed gerbil references and do not invent fancy mouse colors unless the species/source supports them |
| Rabbit | White, black, gray, brown, broken pattern, Dutch-like markings, Himalayan-like markings, lop-type silhouette |
| Dog | Breed-neutral small dog first, then white, black, brown, cream, tan, tricolor, Shiba-like, Chihuahua-like, Dachshund-like, Pomeranian-like directions |
| Cat | Brown tabby, mackerel tabby, silver tabby, black, white, calico, tortoiseshell, tuxedo, point pattern, longhair direction |
| Gecko | Japanese gecko-like gray-brown, leopard gecko normal, high yellow, tangerine, mack snow, and other real morphs only after species/morph labels are checked |
| Hamster | Golden Syrian, cream, cinnamon, black, gray, white, sable, banded, dominant spot, roan, winter white dwarf, Campbell's dwarf, Roborovski, Chinese hamster |

When a variant is visually too subtle at taskbar size, keep it as `needs_review` instead of promoting it to a normal selectable option.

## Initial Art Thread Contract

Each animal-image thread should produce one source-truth image first. This is a concept/source frame, not the full animation set.

Required image properties:

- Transparent background.
- One complete animal only.
- Strict side view.
- Pixel-art sprite or clean illustrated sprite that remains readable at Windows taskbar size.
- Consistent low foot or belly baseline suitable for later animation frames.
- Complete ears, feet, whiskers where anatomically relevant, and complete tail where relevant.
- No text, border, scenery, cast shadow, costume, human-like pose, or multiple animals.
- No cropped body parts.
- Output should be suitable for later normalization into a 96x64 runtime canvas unless the species needs an explicit larger render profile.

Each thread should save or report:

- The generated image or exact image-generation prompt if direct image output is unavailable.
- A short species-readability review.
- Any likely animation risks, such as tail cropping, foot baseline, silhouette confusion, or scale mismatch.

## Species Notes

### Chinchilla

- Larger and rounder than a degu.
- Very dense fluffy body.
- Very large rounded ears.
- Short muzzle.
- Thick fluffy tail.
- Movement should feel heavier and softer than a degu.

### Macaroni Mouse

- Smallest candidate.
- Short rounded body.
- Large dark eyes.
- Small hands and feet.
- Distinctive thick tail, not a normal thin mouse tail.
- Should also be checked against fat-tailed gerbil references to avoid species confusion.

### Rabbit

- Compact body with strong hind legs.
- Long ears are the primary silhouette signal.
- Tail should be small and not dominate the sprite.
- Hop animation will likely need species-specific timing.

### Dog

- Start with a small companion-dog silhouette rather than a large breed.
- Clear muzzle, four-foot stance, and readable tail.
- Breed-specific coats should come later as variants.
- Movement profile can support walk, scurry, sit/stand, and reaction bubbles.

### Cat

- Flexible low body, visible ears, whiskers, and expressive tail.
- Idle, walk, turn, sit/stand, groom, and reaction frames are high-value.
- Avoid making the silhouette dog-like at taskbar size.

### Gecko

- Low reptile body with four splayed feet.
- Distinct head, tapering tail, and toe pads where readable.
- No whiskers or mammal-style grooming assumptions.
- Wheel and foraging behaviors may need opt-in compatibility rather than default use.

### Hamster

- Small, compact body with short legs.
- Rounded ears, blunt muzzle, visible cheek area, and tiny paws.
- Tail should be extremely short or visually absent depending on species.
- Should not be confused with macaroni mouse: hamster silhouette is rounder, shorter-tailed, and cheekier.
- Initial variant coverage should consider real pet hamster types and colors such as golden Syrian, cream, cinnamon, black, gray, white, sable, banded, dominant spot, roan, winter white dwarf, Campbell's dwarf, Roborovski, and Chinese hamster where silhouette or markings are visually useful.
- Keep breed/species labels conservative if visual distinction is weak at taskbar size.

## Implementation Direction

When implementation starts, refactor toward data-driven species registration:

- `species_id`, display labels, default variants, and icon.
- `variant_id`, labels, source sprite sheet paths, and visual metadata.
- Per-species action availability.
- Per-species frame mapping for shared behavior states.
- Settings persistence migration from degu-only `variant` to species plus variant fields.

Keep migration compatibility with existing Degu Desktop settings. Existing degu coat users should keep their chosen coat after the upgrade.

## Acceptance Criteria

Requirements phase is complete when:

- This requirements document exists.
- A management task ledger lists every initial animal thread and its status.
- One image-production thread exists for each initial candidate animal.
- Each image thread has the shared art contract plus its species-specific notes.

Implementation phase is not complete until:

- Species and coat are separate in code.
- Existing Degu Desktop behaviors still pass tests.
- New species can be selected from settings or tray UI.
- At least one non-degu species can render without degu-only naming leaking into user-facing UI.
- Importer output is deterministic.
- Windows build still succeeds.

## Verification Plan

For later code work, use the AnimalsDesktop standard checks:

```powershell
gofmt -w cmd\animalsdesktop\main_windows.go cmd\animalsdesktop\motion_windows_test.go cmd\importsheet\main.go cmd\importsheet\main_test.go
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go run ./cmd/importsheet
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

For art intake:

- Visual review against the shared image contract.
- Check silhouette readability on light and dark backgrounds.
- Check that the animal remains recognizable when downscaled toward 96x64.
- Record any normalization or animation risks before accepting the image as an asset seed.
