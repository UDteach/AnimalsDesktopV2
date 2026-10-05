# Macaroni Mouse One-Frame Method Fullrun - 2026-06-25

Target: `macaroni_mouse_tan` set00, 62 one-frame ImageGen candidates.

Write boundary: this folder only. Do not promote to `accepted-frames`, do not alter runtime/catalog code, and do not touch GitHub/Pages/release assets.

## Method

- Built-in Codex `image_gen`, one output per frame.
- Raw generated PNGs are copied from `/Users/kyota/.codex/generated_images/019efbfa-24d0-7aa0-8db4-d2447df97aa5`.
- Background request: perfectly flat `#00ff00`.
- Alpha conversion: `${CODEX_HOME}/skills/.system/imagegen/scripts/remove_chroma_key.py`.
- Normalization: trim alpha subject, fit inside 86x54, place on 96x64 transparent canvas with shared lower baseline.
- QA artifacts: `raw/`, `alpha/`, `frames/`, `contacts/`, `qa/`, `prompts/`, `manifest.csv`.

## Species Guardrails

- Must read as fat-tailed gerbil / macaroni mouse, not hamster.
- Keep a slim small body, slightly pointed nose, small rounded ears, black eye, pale belly, and short thick pale club-shaped tail.
- Reject long thin mouse tails, missing tails, round hamster cheek pouches, bulky hamster body, cropped ears/feet/whiskers/tail, text, shadow, props, scenery, or multiple animals.

## Pose Table

See `poses.csv` for the authoritative 62-row pose table:

- idle: 00-03
- walk: 04-11
- fast scurry: 12-19
- sniff/nibble: 20-25
- paw/action: 26-31
- turn: 32-39
- chew: 40-43
- ground check: 44-47
- alert rest: 48-51
- face groom: 52-55
- reaction: 56-61

## Attempt Log

| Frame | Group | Pose summary | Accepted attempt | Decision notes |
| --- | --- | --- | ---: | --- |
| 00 | idle | neutral low side-view stance facing right; all four feet visible; calm breathing start | 1 | accepted; first generated frame; mechanical pass |
| 01 | idle | tiny inhale shift with chest slightly lifted; feet stay planted on same baseline | 1 | accepted; mechanical pass; batch 01-10 |
| 02 | idle | tiny exhale shift with body subtly lower; tail still short thick and visible | 1 | accepted; mechanical pass; batch 01-10 |
| 03 | idle | small whisker and ear attention shift; head still facing right | 1 | accepted; mechanical pass; batch 01-10 |
| 04 | walk | slow walk phase; front right paw reaches forward; rear foot planted | 1 | accepted; mechanical pass; batch 01-10 |
| 05 | walk | slow walk phase; front paw contacts ground; body moves slightly forward | 1 | accepted; mechanical pass; batch 01-10 |
| 06 | walk | slow walk phase; rear paw swings forward under body | 1 | accepted; mechanical pass; batch 01-10 |
| 07 | walk | slow walk phase; rear paw contacts ground; tail balances low | 1 | accepted; mechanical pass; batch 01-10 |
| 08 | walk | slow walk phase alternate; front paw lifts small amount | 1 | accepted; mechanical pass; batch 01-10 |
| 09 | walk | slow walk phase alternate; opposite forefoot reaches forward | 1 | accepted; mechanical pass; batch 01-10 |
| 10 | walk | slow walk phase alternate; rear foot pushes gently | 1 | accepted; mechanical pass; batch 01-10 |
| 11 | walk | slow walk return to neutral stride; all feet readable | 1 | accepted; mechanical pass; batch 11-20 |
| 12 | fast_scurry | fast low scurry phase; body stretched slightly longer; front paw extended | 1 | accepted; mechanical pass; batch 11-20 |
| 13 | fast_scurry | fast low scurry phase; compact body gathered under shoulders | 1 | accepted; mechanical pass; batch 11-20 |
| 14 | fast_scurry | fast low scurry phase; rear feet push back; nose pointed forward | 1 | accepted; mechanical pass; batch 11-20 |
| 15 | fast_scurry | fast low scurry phase; front feet touch down; tail thick and low | 1 | accepted; mechanical pass; batch 11-20 |
| 16 | fast_scurry | fast low scurry alternate; body low and streamlined | 1 | accepted; mechanical pass; batch 11-20 |
| 17 | fast_scurry | fast low scurry alternate; rear feet tucked under body | 1 | accepted; mechanical pass; batch 11-20 |
| 18 | fast_scurry | fast low scurry alternate; tiny forward lean with paws visible | 1 | accepted; mechanical pass; batch 11-20 |
| 19 | fast_scurry | fast low scurry recovery; returns to stable baseline | 1 | accepted; mechanical pass; batch 11-20 |
| 20 | sniff_nibble | nose down sniffing ground; forefeet planted; pointed nose visible | 1 | accepted; mechanical pass; batch 11-20 |
| 21 | sniff_nibble | nose slightly lower and forward; whiskers angled down | 1 | accepted; mechanical pass; batch 21-31 |
| 22 | sniff_nibble | small nibble posture with mouth near forepaws; no food prop | 1 | accepted; mechanical pass; batch 21-31 |
| 23 | sniff_nibble | compact nibble with head tucked; pale belly visible | 1 | accepted; mechanical pass; batch 21-31 |
| 24 | sniff_nibble | head lifts slightly from nibble; front paws close together | 1 | accepted; mechanical pass; batch 21-31 |
| 25 | sniff_nibble | nibble recovery with head halfway up; tail remains visible | 1 | accepted; mechanical pass; batch 21-31 |
| 26 | paw_action | one front paw raised near chest; other feet grounded | 1 | accepted; mechanical pass; batch 21-31 |
| 27 | paw_action | front paw reaches forward lightly as if tapping ground | 1 | accepted; mechanical pass; batch 21-31 |
| 28 | paw_action | front paw retracts; shoulders slightly lifted | 1 | accepted; mechanical pass; batch 21-31 |
| 29 | paw_action | small body weight shift over forefeet; no prop | 1 | accepted; mechanical pass; batch 21-31 |
| 30 | paw_action | front paw brushes muzzle briefly; not full groom | 1 | accepted; mechanical pass; batch 21-31 |
| 31 | paw_action | paw returns to ground; stable side stance | 1 | accepted; mechanical pass; batch 21-31 |
| 32 | turn | turn anticipation facing right; head begins rotating toward viewer | 1 | accepted; mechanical pass; batch 32-41 |
| 33 | turn | three-quarter view; body still compact; tail thick visible behind | 1 | accepted; mechanical pass; batch 32-41 |
| 34 | turn | near front-facing pivot; both ears visible; no standing upright | 1 | accepted; mechanical pass; batch 32-41 |
| 35 | turn | three-quarter view facing left direction briefly; feet grounded | 1 | accepted; mechanical pass; batch 32-41 |
| 36 | turn | left-facing side silhouette during turn; short thick tail still visible | 1 | accepted; mechanical pass; batch 32-41 |
| 37 | turn | three-quarter return; head rotating back to right | 1 | accepted; mechanical pass; batch 32-41 |
| 38 | turn | almost right-facing again; body aligns with baseline | 1 | accepted; mechanical pass; batch 32-41 |
| 39 | turn | right-facing turn completion; neutral scale and baseline | 1 | accepted; mechanical pass; batch 32-41 |
| 40 | chew | small chew motion; head slightly bobbed; mouth closed around invisible bite | 1 | accepted; mechanical pass; batch 32-41 |
| 41 | chew | chew phase with cheek not inflated; avoid hamster pouch look | 1 | accepted; mechanical pass; batch 32-41 |
| 42 | chew | chew phase with muzzle slightly forward; forepaws steady | 1 | accepted; mechanical pass; batch 42-51 |
| 43 | chew | chew recovery; head returns to neutral | 1 | accepted; mechanical pass; batch 42-51 |
| 44 | ground_check | low ground check; nose nearly touches baseline; back remains slim | 1 | accepted; mechanical pass; batch 42-51 |
| 45 | ground_check | ground check sweep forward; pointed nose and whiskers readable | 1 | accepted; mechanical pass; batch 42-51 |
| 46 | ground_check | ground check sweep back; body low but complete | 1 | accepted; mechanical pass; batch 42-51 |
| 47 | ground_check | ground check recovery; head rises to normal | 1 | accepted; mechanical pass; batch 42-51 |
| 48 | alert_rest | alert rest crouch; ears perked; body still low | 1 | accepted; mechanical pass; batch 42-51 |
| 49 | alert_rest | alert rest head lifted; black eye prominent | 1 | accepted; mechanical pass; batch 42-51 |
| 50 | alert_rest | alert rest listening; one ear angled slightly back | 1 | accepted; mechanical pass; batch 42-51 |
| 51 | alert_rest | alert rest relax; returns toward neutral | 1 | accepted; mechanical pass; batch 42-51 |
| 52 | face_groom | face groom start; front paw lifts to cheek; feet grounded | 1 | accepted; mechanical pass; immediate single-frame processing |
| 53 | face_groom | face groom paw touches muzzle; no human-like pose | 1 | accepted; mechanical pass; immediate single-frame processing |
| 54 | face_groom | face groom paw sweeps down; eye still visible | 1 | accepted; mechanical pass; immediate single-frame processing |
| 55 | face_groom | face groom finish; paw returns to ground | 1 | accepted; mechanical pass; immediate single-frame processing |
| 56 | reaction | small surprise reaction; head lifts and ears perk; all feet grounded | 1 | accepted; mechanical pass; immediate single-frame processing |
| 57 | reaction | quick recoil backward but low quadruped; no jump | 1 | accepted; mechanical pass; immediate single-frame processing |
| 58 | reaction | tiny startle stretch forward; tail thick and low | 1 | accepted; mechanical pass; immediate single-frame processing |
| 59 | reaction | reaction settle; shoulders lower | 1 | accepted; mechanical pass; immediate single-frame processing |
| 60 | reaction | return-to-idle transition; feet align to baseline | 1 | accepted; mechanical pass; immediate single-frame processing |
| 61 | reaction | final neutral recovery matching frame 00 scale and stance | 1 | accepted; mechanical pass; immediate single-frame processing |

## Final QA Summary

- Normalized 96x64 frames: 62/62
- Raw copied sources: 62/62
- Alpha sources: 62/62
- Per-frame contacts: 62/62
- QA JSON files: 62/62
- Prompt files: 62/62
- Manifest rows: 62/62
- Missing frames: []
- Bad size frames: []
- Empty alpha frames: []
- Edge alpha frames: []
- Green residue > 12 frames: []
- 62-frame contact sheet: `docs/art-source/one-frame-method-fullrun-20260625/macaroni-mouse/macaroni-mouse-set00-62-contact.png`

## Boundary Check

- Generation stopped after frame 61 per parent checkpoint.
- No promotion to accepted-frames was performed.
- Runtime/catalog/other species directories were not intentionally edited by this thread.
