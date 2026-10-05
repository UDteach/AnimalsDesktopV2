# Sable Panda final 62-frame gate

Status: **PARENT APPROVED 62/62 — FROZEN AS COAT-VARIANT BASE**

No frame62 or coat-variant generation was started. No catalog, runtime, import,
commit, push, tag, release, or publication action was performed.

## Canonical result

- Canonical range: frame00 through frame61, 62/62 files
- Canonical aggregate SHA-256:
  `7947C23E4FCBCF304233CBDB09231584A9E7C479C1EE404A6712816087977CA9`
- Final generation-anchor aggregate SHA-256:
  `C782AF07B88ADF51C7185AF73FAE9FD990CA2D0FD3CDFE85485804AEEBCC18A7`
- Parent-approved canonical 00–52 baseline aggregate:
  `C5B1CBD3580687E7467658CB562BC0D955D7EE0D304BFA21F1B044B252D72334`
- Protected 00–52 plus four baseline anchors: 57/57 match, 0 mismatches
- Final 00–61 plus five anchors: 67/67 match, 0 mismatches
- Protected original dirty checkout: 25,728 status entries. Historical
  PowerShell culture-sorted SHA-256
  `5682fcb75549f352404a544d34a30efe27bf21a62e3054b3bb1888b952128da4`;
  corrected ordinal-sorted SHA-256
  `0d7c471812644c51220002d5707f0def40902d66c389d873919a26ebf508ebc2`.
  Both are from the same unchanged line set; see
  `fingerprint-sort-order-correction.md`.

Frame53 is the parent-approved existing-call requalification of unused
`frame-52-attempt-01`. Its original frame52 prompt and provenance remain
unchanged. It was promoted once as frame53 without pixel modification, while
canonical frame52 remains attempt02.

## New ImageGen calls for frames54–61

Exactly 12 fresh built-in ImageGen calls were made for this bounded packet.
Eight candidates were accepted and four were rejected. Rejected pixels were
preserved and were not used as retry references.

| Frame | Attempt | Call source | Result |
| --- | --- | --- | --- |
| 54 | 01 | `call_LKZefnUjRSxDo7abTNgqKZP3.png` | reject: official lower alpha run |
| 54 | 02 | `call_Hc3BcyMtndOlkbCNYkaRpTms.png` | accepted |
| 55 | 01 | `call_uHqydMUgoq2nN4tD9IxKD2sc.png` | accepted |
| 56 | 01 | `call_nQ3sHBdza0AEoCgvq7abfh5X.png` | accepted |
| 57 | 01 | `call_E2NJgrI77hQ4J7QKiFQierUe.png` | reject: abrupt early alert rise |
| 57 | 02 | `call_nDtn2TEAUda5kaN3eZUozY2f.png` | accepted |
| 58 | 01 | `call_XBw8FNUVTFkSeWJHMe2NbT4e.png` | accepted |
| 59 | 01 | `call_GzF8QG6NOxy46vYDbcum3Skx.png` | accepted |
| 60 | 01 | `call_zcfTxq3zfzpnLZz7CBX1AHSI.png` | reject: pinhole and abrupt descent |
| 60 | 02 | `call_Nut4Ky7wNSPjXBNWsrxtGuXb.png` | accepted |
| 61 | 01 | `call_kuzFe7JSjH0vV0YKG97uSyTr.png` | reject: rest/creep collapse |
| 61 | 02 | `call_uhxWktZ8WvvNEkjCogme7eZt.png` | accepted |

Frame53 used 0 new calls. Every accepted fresh frame has one source call and
one exact prompt. No sheet, local painting, alpha fill, recolor, warp,
interpolation, or repaired candidate was promoted.

## Mechanical QA

- Official action-aware audit with boundaries
  `4,12,20,26,32,40,48,56`:
  `valid=62 missing=0 invalid=0 warnings=0`
- Shared rounded-matte audit: 62/62 PASS, 0 failed frames
- Exact canonical duplicate pairs: 0
- Transparent pinholes, crop, edge contact, lower shelf, and motion-warning
  findings in the final canonical set: 0
- Rest loop adjacent changed pixels, including closure:
  48→49 1212, 49→50 1225, 50→51 1274, 51→52 1276,
  52→53 1385, 53→54 1354, 54→55 1220, 55→48 1237
- Alert loop adjacent changed pixels, including closure:
  56→57 1422, 57→58 1523, 58→59 1536, 59→60 1689,
  60→61 1580, 61→56 1228

## Visual QA

The canonical full-set light, dark, and checker contacts were inspected at
nearest-neighbor scale. The result retains a white/cream-first Sable Panda
read, warm brown points, complete tail, four-limb support, stable baseline and
padding, and no visible green halo, crop, shelf, floor, or shadow.

The rest 48–55 contact reads as a low supported breathing loop: frame53 is the
approved quiet inhalation rebound, frame54 lowers midway, frame55 returns
toward frame48 without exact duplication. The alert 56–61 contact reads as a
low rise to frame59 peak and a controlled return through frames60–61. The
frame61→56 closure is same-height but visibly differs in head/ear angle, paw
support, and tail detail. Frame39→horizontally mirrored frame04 remains a
coherent left-facing runtime exit.

Parent review inspected the alert loop on checker and dark backgrounds plus
the saved `59 -> 60`, `60 -> 61`, and `61 -> 56` transition contacts. The
faster alert descent is intentional and remains readable; `61 -> 56` closes
without a visible pose jump. The complete checker and dark contacts also pass
for silhouette, anatomy, baseline, padding, and coat readability. This
canonical set is approved and frozen as the geometry base for the narrow
four-cell Sable and Albino coat-only exception.

## Evidence

- `action-aware-audit.json`
- `shared-matte.json`
- `visual-duplicate-loop-report.json`
- `canonical-00-61-baseline.json`
- `canonical-00-61-verification.json`
- `protected-00-52-verification.json`
- `protected-original-fingerprint.json`
- `fingerprint-sort-order-correction.md`
- `contacts/full-00-61-{light,dark,checker}.png`
- `contacts/rest-loop-48-55-{light,dark,checker}.png`
- `contacts/alert-loop-56-61-{light,dark,checker}.png`
- `runtime-exit-39-to-mirrored04/frame39-to-mirrored04-{light,dark,checker}.png`
- Per-candidate prompts, raw files, alpha files, reports, accept/reject records,
  and connection contacts under `qa/final-53-61/candidates/`

## Scope preservation

- Branch: `codex/animal-addition-foundation-20260723`
- All writes in this packet stayed inside the owned Sable Panda run directory
- Protected original `E:/Development/AnimalDesktop` was read only and its
  fingerprint is unchanged
- Parent approval was recorded on 2026-07-24 after the final contact review
