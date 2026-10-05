# Popularity Priority Sources - 2026-06-27

This note records current ranking sources for future AnimalsDesktop asset
priority decisions. The parent thread still owns final release order. User
priority overrides ranking ties when needed.

## Current Sources Checked

- Anicom `犬種ランキング2026`, published 2026-01-29:
  https://www.anicom-sompo.co.jp/news-release/2025/20260128/
- Anicom `猫種ランキング2026`, published 2026-02-19:
  https://www.anicom-sompo.co.jp/news-release/2025/20260219/
- iPet / Nyanpedia `人気飼育猫種ランキング2026年版`, published
  2026-01-21:
  https://nyanpedia.com/cat-ranking-2026/
- iPet `人気飼育犬種・猫種・小動物ランキング2026年版`, published from
  2025 policy starts:
  https://www.ipet-ins.com/info/40566/
- SBI Prism `2025年度人気ランキング` for small animals, birds, reptiles:
  https://www.sbiprism.co.jp/ranking/2025/ranking.html

## Bird Priority

SBI Prism ranks birds as:

1. セキセイインコ
2. オカメインコ
3. コザクラインコ
4. ブンチョウ
5. ウロコインコ / マメルリハ

Current implementation order keeps the user's explicit active priority:

1. `budgerigar_green_yellow`
2. `cockatiel_normal_gray`
3. `java_sparrow_normal`
4. `lovebird_peach_faced`
5. `parrotlet_green` or a green-cheek conure candidate

Reason: the ranking puts lovebird before Java sparrow, but the active user
instruction names budgerigar, cockatiel, and Java sparrow first.

## Cat Priority

Anicom ranks 2026 cat breeds as:

1. Scottish Fold
2. mixed cat
3. Munchkin
4. Ragdoll
5. Minuet
6. Siberian
7. British Shorthair
8. American Shorthair
9. Norwegian Forest Cat
10. Ragamuffin

iPet / Nyanpedia uses a broader `ミックス` category and reports the same stable
top-five family for product planning: mixed cat, Scottish Fold, Munchkin,
Ragdoll, and Minuet. Future cat lanes should present cats by breed, not as a
single generic `cat` family.

Candidate variant queue:

1. `scottish_fold_silver_tabby`
2. `mixed_cat`
3. `munchkin`
4. `ragdoll_seal_bicolor`
5. `minuet`
6. `siberian`
7. `british_shorthair_blue`
8. `american_shorthair_silver_tabby`

## Dog Priority

Anicom ranks 2026 dog breeds as:

1. MIX dog under 10kg
2. Toy Poodle
3. Chihuahua
4. Miniature Dachshund
5. Pomeranian
6. Shiba, including Mameshiba
7. Miniature Schnauzer
8. French Bulldog
9. Maltese
10. Kaninchen Dachshund

Candidate variant queue:

1. `mixed_small_dog`
2. `toy_poodle`
3. `chihuahua`
4. `miniature_dachshund_red`
5. `pomeranian_orange`
6. `shiba_inu_red`
7. `miniature_schnauzer_salt_pepper`
8. `french_bulldog_fawn`
9. `maltese_white`
10. `kaninchen_dachshund`

## Small Animal Priority

iPet's 2026 small-animal summary ranks Djungarian hamster and Netherland Dwarf
at the top and notes Golden hamster, Kinkuma hamster, chinchilla, degu, birds,
and ferret among ranked exotic animals. SBI Prism's category ranking puts
rabbit, hamster, ferret, chinchilla, guinea pig, hedgehog, degu, momonga,
squirrel, and prairie dog in the top ten.

The current runtime already covers many of these families. Color/morph variants
that still fit the ranking-driven queue include:

- `djungarian_hamster_yellow`
- `djungarian_hamster_pearl_white`
- `longhair_hamster_black_white`
- `netherland_dwarf_himalayan`
- `rabbit_gray`
- `leucistic_sugar_glider`
- `american_flying_squirrel`
- `african_dormouse`

## Current Asset-Pack Link

The 2026-06-27 reusable upcoming asset pack is:

- `assets/source/upcoming/20260627/manifest.json`
- `assets/source/upcoming/20260627/color/`
- `assets/source/upcoming/20260627/silhouette/`

It covers the current small-animal/bird/reptile upcoming set, but not the full
cat and dog breed queues above. Cat/dog source lanes should use this note plus
existing `assets/source/animals/generated/` inventory before opening new
ImageGen threads.
