# Albino final 62-frame gate

Status: **PARENT APPROVED 62/62 — FROZEN**

The Albino coat is approved as a separate `ferret_albino` variant. It keeps
the parent-approved Sable Panda silhouette, right-facing motion, camera,
scale, baseline, and contact points while reading as warm white/cream fur with
ruby-red eyes and pink nose, ears, and paws.

## Canonical result

- Canonical range: frame00 through frame61, 62/62 files
- Accepted-frame aggregate SHA-256:
  `a16586fda8e858be602da4d4e37648488d2b734f7c01b8bda44d3766a6160590`
- Aggregate method: SHA-256 over accepted PNG bytes concatenated in numeric
  frame order
- Source-sheet SHA-256:
  `316fcc52f8c683455330c7e176a7d8a8183ab38e22a50b662ea799c6d0685b2a`
- Source sheet dimensions: `5952x64`
- Accepted PNGs matching their recorded standalone winners: 62/62
- Unique accepted PNG hashes: 62/62

## Final alert packet

Frames56–61 used seven fresh built-in ImageGen calls: six accepted winners and
one preserved reject. Frame60 attempt01 was rejected and attempt02 was
accepted. The final geometry gate reached minimum IoU `0.988205560` and
maximum centroid shift `0.196277 px`.

Parent review inspected raw final frames and the full set, alert loop, and
closure on checker and dark contacts. The loop reads as a staged rise from
frame56 to the frame59 peak, a controlled fall through frames60–61, and a
coherent frame61→56 closure. Full anatomy, padding, baseline, dark-background
readability, and true Albino identity remain intact.

## Mechanical QA

- Official action-aware audit with boundaries
  `4,12,20,26,32,40,48,56`:
  `valid=62 missing=0 invalid=0 warnings=0`
- Shared rounded-matte audit: 62/62 PASS, 0 failed frames
- Exact accepted-frame duplicate groups: 0
- Rest loop adjacent changed RGBA pixels, including closure:
  48→49 1217, 49→50 1227, 50→51 1279, 51→52 1275,
  52→53 1385, 53→54 1350, 54→55 1215, 55→48 1232
- Alert loop adjacent changed RGBA pixels, including closure:
  56→57 1425, 57→58 1514, 58→59 1541, 59→60 1686,
  60→61 1578, 61→56 1235

## Provenance and scope

Every promoted PNG came from one standalone built-in ImageGen call. No
generated sheet cell was split or promoted. No local pixel/alpha repair,
painting, recolor, warp, interpolation, or tone correction was used.

The full prompts, raw candidates, rejected attempts, provenance records, and
QA reports are preserved in:

`E:/Development/AnimalDesktop-asset-archives/20260724-ferret-three-variants/ferret-albino-set00-oneframe-62/`

This verdict authorizes local integration only. It does not authorize commit,
push, tag, release, or publication.
