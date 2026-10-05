# Sable final 62-frame gate

Status: **PARENT APPROVED 62/62 — FROZEN**

The classic Sable coat is approved as a separate `ferret_sable` variant. It
keeps the parent-approved Sable Panda silhouette, right-facing motion,
camera, scale, baseline, and contact points while reading as a dark warm-brown
body with the familiar lighter face mask and undercoat.

## Canonical result

- Canonical range: frame00 through frame61, 62/62 files
- Accepted-frame aggregate SHA-256:
  `33c2da1526d4bd71c7ce76f1e08ca6213406a590b6059d51849514118a4a69f0`
- Aggregate method: SHA-256 over accepted PNG bytes concatenated in numeric
  frame order
- Source-sheet SHA-256:
  `22322c1fc8a7b98ba79d603192ccc68e36040cea714eac6436f7b5b711482700`
- Source sheet dimensions: `5952x64`
- Accepted PNGs matching their recorded standalone winners: 62/62
- Unique accepted PNG hashes: 62/62

## Final alert packet

Frames56–61 used 10 fresh built-in ImageGen calls: six accepted winners and
four preserved rejects. Frame60 attempts 01 and 02 and frame61 attempts 01
and 02 were rejected; frame60 attempt03 and frame61 attempt03 were accepted.
The final geometry gate reached minimum IoU `0.987133667` and maximum centroid
shift `0.133973 px`.

Parent review inspected the full set and the alert loop on original-source,
checker, and dark contacts. The loop reads as a controlled rise from frame56
to the frame59 peak, a fall through frames60–61, and a coherent frame61→56
closure. The complete tail, four-limb support, padding, baseline, and classic
Sable coat remain readable.

## Mechanical QA

- Official action-aware audit with boundaries
  `4,12,20,26,32,40,48,56`:
  `valid=62 missing=0 invalid=0 warnings=0`
- Shared rounded-matte audit: 62/62 PASS, 0 failed frames
- Exact accepted-frame duplicate groups: 0
- Rest loop adjacent changed RGBA pixels, including closure:
  48→49 1213, 49→50 1227, 50→51 1277, 51→52 1273,
  52→53 1383, 53→54 1360, 54→55 1221, 55→48 1237
- Alert loop adjacent changed RGBA pixels, including closure:
  56→57 1417, 57→58 1535, 58→59 1540, 59→60 1689,
  60→61 1579, 61→56 1233

## Provenance and scope

Every promoted PNG came from one standalone built-in ImageGen call. No
generated sheet cell was split or promoted. No local pixel/alpha repair,
painting, recolor, warp, interpolation, or tone correction was used.

The full prompts, raw candidates, rejected attempts, provenance records, and
QA reports are preserved in:

`E:/Development/AnimalDesktop-asset-archives/20260724-ferret-three-variants/ferret-sable-set00-oneframe-62/`

This verdict authorizes local integration only. It does not authorize commit,
push, tag, release, or publication.
