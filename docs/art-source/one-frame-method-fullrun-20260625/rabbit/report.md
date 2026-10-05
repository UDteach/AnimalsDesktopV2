# Rabbit Chestnut Agouti One-Frame Fullrun - 2026-06-25

Dedicated output root:

`docs/art-source/one-frame-method-fullrun-20260625/rabbit/`

## Scope

- Species: `rabbit_chestnut_agouti`
- Method: Codex built-in ImageGen, one frame at a time, retry failed frame numbers only.
- Background: perfectly flat `#00ff00` chroma key, then local alpha removal.
- Runtime/accepted-frame mutation: none intended.
- Baseline retained: existing sheet-extracted `accepted-candidates` remain untouched.

## Prompt Template

See `prompts/frame-NN-attempt-MM.txt` for exact per-frame prompts. The stable base prompt requires a single complete chestnut agouti rabbit with cream chest, long ears, strong hind legs, and small white tail on flat `#00ff00`.

## Pose Table

| Frame | Group | Pose |
| ---: | --- | --- |
| 00 | idle | neutral standing rest, side view facing right |
| 01 | idle | soft breathing rest, side view facing right |
| 02 | idle | tiny weight shift, side view facing right |
| 03 | idle | settled crouch rest, side view facing right |
| 04 | hop_walk | hop-walk contact, front paws down and hind legs gathered, facing right |
| 05 | hop_walk | hop-walk lift, front paws just lifting, facing right |
| 06 | hop_walk | hop-walk mid hop, compact body, facing right |
| 07 | hop_walk | hop-walk hind push, strong rear legs extended, facing right |
| 08 | hop_walk | hop-walk landing, front paws near ground, facing right |
| 09 | hop_walk | hop-walk settle, hind feet catching up, facing right |
| 10 | hop_walk | hop-walk small stride, body low and long, facing right |
| 11 | hop_walk | hop-walk return to contact, facing right |
| 12 | fast_hop | fast hop stretched low body, ears swept slightly back, facing right |
| 13 | fast_hop | fast hop compact gathered paws, facing right |
| 14 | fast_hop | fast hop long rear push, strong hind legs visible, facing right |
| 15 | fast_hop | fast hop tucked airborne pose, facing right |
| 16 | fast_hop | fast hop forward lean, ears trailing, facing right |
| 17 | fast_hop | fast hop landing with front paws down, facing right |
| 18 | fast_hop | fast hop quick recovery step, facing right |
| 19 | fast_hop | fast hop low return, facing right |
| 20 | sniff_nibble | nose down sniffing ground, facing right |
| 21 | sniff_nibble | nose forward sniff, whiskers visible, facing right |
| 22 | sniff_nibble | tiny nibble with mouth near front paws, facing right |
| 23 | sniff_nibble | head lowered nibble, paws planted, facing right |
| 24 | sniff_nibble | sniff pause with raised nose and long ears, facing right |
| 25 | sniff_nibble | small chew pause, side view facing right |
| 26 | paw_hop_prep | one front paw lifted near chest, hind legs crouched to hop, facing right |
| 27 | paw_hop_prep | both front paws close to mouth, hind legs tucked, facing right |
| 28 | paw_hop_prep | short reaching paw gesture, rabbit crouch, facing right |
| 29 | paw_hop_prep | paw tap on ground, hind feet anchored, facing right |
| 30 | paw_hop_prep | upright hop preparation, front paws lifted, long ears visible |
| 31 | paw_hop_prep | settling after hop preparation, facing right |
| 32 | turn | begin turning from right-facing side toward viewer |
| 33 | turn | three-quarter front turn, compact rabbit body |
| 34 | turn | front-facing pause, both long ears visible |
| 35 | turn | three-quarter left turn, compact rabbit body |
| 36 | turn | left-facing side view, small white tail visible |
| 37 | turn | turning back from left toward viewer |
| 38 | turn | three-quarter right return, long ears visible |
| 39 | turn | right-facing side view restored |
| 40 | chew | sitting low, front paws near mouth, chewing |
| 41 | chew | round cheek chew, front paws tucked |
| 42 | chew | small bite motion, head slightly down |
| 43 | chew | chew pause, cheeks rounded |
| 44 | ground_check | nose to floor checking ground, facing right |
| 45 | ground_check | front paws low, inspecting ground, facing right |
| 46 | ground_check | head sweeping along ground, facing right |
| 47 | ground_check | nose down pause, long ears held back, facing right |
| 48 | alert_rest_ear_twitch | alert rest with head raised and long ears upright, facing right |
| 49 | alert_rest_ear_twitch | alert ears forward, still body, facing right |
| 50 | alert_rest_ear_twitch | upright attentive rest, one ear subtly twitching |
| 51 | alert_rest_ear_twitch | alert crouch ready to hop, facing right |
| 52 | face_groom | front paws washing face, compact sitting rabbit pose |
| 53 | face_groom | one paw over cheek, compact sitting rabbit pose |
| 54 | face_groom | both paws near nose grooming |
| 55 | face_groom | finishing face groom, paws lowering |
| 56 | reaction | startled flinch, compact body, facing right |
| 57 | reaction | quick hop reaction, all paws close to body |
| 58 | reaction | surprised upright tiny pose, long ears vertical |
| 59 | reaction | ducking reaction, low body and ears back |
| 60 | reaction | recovering from reaction, facing right |
| 61 | reaction | calm return after reaction, facing right |

## Artifact Layout

- Raw built-in ImageGen copies: `raw/frame-NN-attempt-MM-raw.png`
- Alpha converted source: `alpha/frame-NN-attempt-MM-alpha.png`
- Normalized candidate frame: `frames/frame-NN.png`
- Attempt-specific normalized frame: `frames/frame-NN-attempt-MM.png`
- Per-frame contact preview: `contacts/frame-NN-attempt-MM-contact.png`
- 62-frame contact sheet: `contacts/rabbit-set00-fullrun-contact.png`
- QA JSON: `qa/frame-NN-attempt-MM-qa.json`
- Prompt text: `prompts/frame-NN-attempt-MM.txt`
- Manifest: `manifest.csv`

## Progress Log

- Generated candidate rows: 62
- Current accepted candidate frames: 62/62
- Frame 00: candidate attempt 01; initial one-frame candidate
- Frame 01: candidate attempt 01; initial one-frame candidate
- Frame 02: candidate attempt 01; initial one-frame candidate
- Frame 03: candidate attempt 01; initial one-frame candidate
- Frame 04: candidate attempt 01; initial one-frame candidate
- Frame 05: candidate attempt 01; initial one-frame candidate
- Frame 06: candidate attempt 01; initial one-frame candidate
- Frame 07: candidate attempt 01; initial one-frame candidate
- Frame 08: candidate attempt 01; initial one-frame candidate
- Frame 09: candidate attempt 01; initial one-frame candidate
- Frame 10: candidate attempt 01; initial one-frame candidate
- Frame 11: candidate attempt 01; initial one-frame candidate
- Frame 12: candidate attempt 01; initial one-frame candidate
- Frame 13: candidate attempt 01; initial one-frame candidate
- Frame 14: candidate attempt 01; initial one-frame candidate
- Frame 15: candidate attempt 01; initial one-frame candidate
- Frame 16: candidate attempt 01; initial one-frame candidate
- Frame 17: candidate attempt 01; initial one-frame candidate
- Frame 18: candidate attempt 01; initial one-frame candidate
- Frame 19: candidate attempt 01; initial one-frame candidate
- Frame 20: candidate attempt 01; initial one-frame candidate
- Frame 21: candidate attempt 01; initial one-frame candidate
- Frame 22: candidate attempt 01; initial one-frame candidate
- Frame 23: candidate attempt 01; initial one-frame candidate
- Frame 24: candidate attempt 01; initial one-frame candidate
- Frame 25: candidate attempt 01; initial one-frame candidate
- Frame 26: candidate attempt 01; initial one-frame candidate
- Frame 27: candidate attempt 01; initial one-frame candidate
- Frame 28: candidate attempt 01; initial one-frame candidate
- Frame 29: candidate attempt 01; initial one-frame candidate
- Frame 30: candidate attempt 01; initial one-frame candidate
- Frame 31: candidate attempt 01; initial one-frame candidate
- Frame 32: candidate attempt 01; initial one-frame candidate
- Frame 33: candidate attempt 01; initial one-frame candidate
- Frame 34: candidate attempt 01; initial one-frame candidate
- Frame 35: candidate attempt 01; initial one-frame candidate
- Frame 36: candidate attempt 01; initial one-frame candidate
- Frame 37: candidate attempt 01; initial one-frame candidate
- Frame 38: candidate attempt 01; initial one-frame candidate
- Frame 39: candidate attempt 01; initial one-frame candidate
- Frame 40: candidate attempt 01; initial one-frame candidate
- Frame 41: candidate attempt 01; initial one-frame candidate; slightly realistic but complete rabbit
- Frame 42: candidate attempt 01; initial one-frame candidate; slightly realistic but complete rabbit
- Frame 43: candidate attempt 01; initial one-frame candidate; realistic but complete rabbit
- Frame 44: candidate attempt 01; initial one-frame candidate; realistic but complete rabbit
- Frame 45: candidate attempt 01; initial one-frame candidate; realistic but complete rabbit
- Frame 46: candidate attempt 01; initial one-frame candidate; realistic but complete rabbit
- Frame 47: candidate attempt 01; initial one-frame candidate; realistic but complete rabbit
- Frame 48: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 49: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 50: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 51: candidate attempt 01; completion fill; complete rabbit
- Frame 52: candidate attempt 01; completion fill; complete rabbit
- Frame 53: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 54: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 55: candidate attempt 01; completion fill; complete rabbit
- Frame 56: candidate attempt 01; completion fill; complete rabbit
- Frame 57: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 58: candidate attempt 01; completion fill; complete rabbit
- Frame 59: candidate attempt 01; completion fill; complete rabbit
- Frame 60: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit
- Frame 61: candidate attempt 01; initial one-frame candidate; realistic style, complete rabbit

## Remaining Risk

- Mechanical QA catches size, alpha, edge contact, empty alpha, and green residue only.
- Final parent-thread visual adoption review still needs to judge pose readability, anatomy, species identity, and style consistency.
