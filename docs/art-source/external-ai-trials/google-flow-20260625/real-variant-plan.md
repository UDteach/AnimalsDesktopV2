# Real Variant Color Plan - 2026-06-25

This plan broadens the current accepted five AnimalsDesktop `set00` families with real-world coat, color, or morph directions only. It is a Google Flow preparation lane, not an accepted asset lane.

Do not promote any output from this directory into `accepted-frames`, runtime sprites, or catalog runtime scope until the local motion and visual QA gates pass.

## First Wave Pilot

Use one color per existing family first. Each pilot uses the same eight representative frames:

```text
00,04,12,26,32,40,52,56
```

| Family | Base accepted variant | First target | Pilot packet | Reason |
| --- | --- | --- | --- | --- |
| Chinchilla | `chinchilla_standard_gray` | beige | `chinchilla-beige-color-pilot/` | Real chinchilla color direction, close enough to base silhouette to test recolor stability. |
| Hamster | `hamster_golden_syrian` | cream | `hamster-cream-color-pilot/` | Real Syrian hamster coat direction and already the safest first Flow packet. |
| Macaroni mouse / fat-tailed gerbil | `macaroni_mouse_tan` | pale sandy cream | `macaroni-mouse-cream-color-pilot/` | Conservative light coat direction; must preserve the thick fat tail and avoid hamster drift. |
| Rabbit | `rabbit_chestnut_agouti` | white | `rabbit-white-color-pilot/` | Common domestic rabbit color direction; tests whether light coats keep outline readability. |
| Sugar glider / momonga | `sugar_glider_gray` | white-faced blonde | `sugar-glider-white-faced-blonde-color-pilot/` | Real pet sugar glider color direction; less destructive than all-white/leucistic as a first recolor. |

## Broader Candidate Backlog

These are candidates for later pilots after the first wave is visually reviewed.

| Family | Candidate directions | Notes |
| --- | --- | --- |
| Chinchilla | beige, ebony, white mosaic, black velvet, violet, sapphire | Beige/ebony/white mosaic already exist as seed catalog directions. Black velvet, violet, and sapphire need label/source confirmation before catalog addition. |
| Hamster | cream, cinnamon, black banded, white, sable, dominant spot, roan, dwarf hamster directions | Cream/cinnamon/black banded/white already exist as seed catalog directions. Dwarf species should not be treated as mere color swaps if body shape changes. |
| Macaroni mouse / fat-tailed gerbil | tan, pale sandy cream, gray-backed/sandy gray if source-confirmed | Keep this conservative. Do not import fancy mouse color names unless they are backed for fat-tailed gerbil/duprasi. |
| Rabbit | white, black, blue gray, fawn, broken orange, Dutch black-white, black otter | White/black/blue gray/fawn and several breed/color directions already exist as seed catalog directions. Breed silhouettes need separate review if ears/body shape should change. |
| Sugar glider / momonga | classic gray, white-faced blonde, mosaic, leucistic, platinum | Use pet sugar glider morph labels. Leucistic/platinum/mosaic can drift face/body markings; run after white-faced blonde succeeds. |

## QA Rules

Accept a pilot only if every generated frame preserves:

- one complete animal per used cell
- same species identity and source family
- same pose intent, camera distance, body scale, and contact baseline
- same eye, muzzle, ear, paw, tail, and membrane geometry where relevant
- clean outline, no transparent pinholes, no green edge debris
- flat transparent or chroma-green background with no floor, shadow, text, grid, labels, props, or duplicate animals

Hard reject if:

- the target color invents a fantasy coat for normal mode
- the species changes into a degu, mouse, hamster, rabbit, or bat-like silhouette
- a light coat loses dark outline readability on a light taskbar
- the macaroni mouse loses the thick fat tail
- the sugar glider membrane turns into bat wings or a flying scene
- rabbit ears, feet, whiskers, or tail are cropped

## Known Anchor Notes

The macaroni mouse accepted `frame-00` has a 1px transparent pinhole and should not be reused as a Flow style anchor until repaired. The first macaroni pilot instead uses frames `04,12,32` as style anchors.
