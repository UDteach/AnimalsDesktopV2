# Upcoming Asset Pack - 2026-06-27

Planning asset pack only. These are not accepted runtime motion sources.

This pack exports the current upcoming-animal ImageGen sheet as reusable
planning assets outside the GitHub Pages directory.

## Contents

- `source/`: copied ImageGen source sheet.
- `color/`: transparent color cutouts, one animal per PNG.
- `silhouette/`: normalized 256x160 black silhouettes, one animal per PNG.
- `contact-sheets/`: color and silhouette review sheets.
- `manifest.json`: stable IDs, JP/EN names, source cells, and output paths.

## Animals

| Order | ID | JP | EN |
| --- | --- | --- | --- |
| 1 | `chipmunk` | シマリス | Striped chipmunk |
| 2 | `leucistic_sugar_glider` | リューシスティックモモンガ | Leucistic sugar glider |
| 3 | `african_dormouse` | アフリカヤマネ | African pygmy dormouse |
| 4 | `netherland_dwarf_himalayan` | ネザーランドドワーフ（ヒマラヤン） | Netherland Dwarf rabbit, Himalayan |
| 5 | `american_flying_squirrel` | アメリカモモンガ | American flying squirrel |
| 6 | `longhair_hamster_black_white` | 白黒長毛ハムスター | Black-and-white long-haired Syrian hamster |
| 7 | `djungarian_hamster_yellow` | イエロージャンガリアン | Yellow Djungarian hamster |
| 8 | `djungarian_hamster_pearl_white` | パールホワイトジャンガリアン | Pearl white Djungarian hamster |
| 9 | `fancy_rat_blue_hooded` | ファンシーラット（ブルーフーディッド） | Fancy rat, blue hooded |
| 10 | `fancy_rat_chocolate_self` | ファンシーラット（チョコレートセルフ） | Fancy rat, chocolate self |
| 11 | `fancy_rat_cream_agouti` | ファンシーラット（クリームアグーチ） | Fancy rat, cream agouti self |
| 12 | `rabbit_gray` | グレーうさぎ | Gray rabbit |
| 13 | `whites_tree_frog` | 水色イエアメガエル | Blue White's tree frog |
| 14 | `leopard_gecko` | ヒョウモントカゲモドキ | Leopard gecko |
| 15 | `budgerigar` | セキセイインコ | Green-and-yellow budgerigar |
| 16 | `cockatiel` | オカメインコ | Normal gray cockatiel |
| 17 | `java_sparrow` | 文鳥 | Java sparrow |
| 18 | `african_fat_tailed_gecko` | ニシアフリカトカゲモドキ | African fat-tailed gecko |

## Rebuild

```sh
python3 scripts/export_upcoming_asset_pack.py
```
