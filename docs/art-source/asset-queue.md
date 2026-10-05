# AnimalsDesktop Asset Queue

This queue tracks asset production rules without exposing unverified animal
families to the public page or runtime picker. Promote or name a new family only
when the parent thread explicitly accepts it as the next production target.

## Current Released State

- Latest test-preview release: `v0.1.4`
- Next local release prep: `v0.1.5`
- Preview scope: sixteen accepted initial runtime animals with Mac ZIPs; Windows ZIP prepared on a Windows machine
- Current local preview animals: chinchilla, hamster, Djungarian hamster, Campbell hamster, macaroni mouse, sugar glider, rabbit, Holland Lop, Netherland Dwarf, Himalayan rabbit, gecko, guinea pig, fancy rat, albino chipmunk, Richardson's ground squirrel, Yorkshire Terrier.
- Public page scope: current preview animals plus text-only next candidate when no page-specific generated silhouette exists.
- Next requested production candidate after this batch: normal striped chipmunk / シマリス, distinct from `albino_chipmunk`.
- Incremental preview release policy: after an animal reaches accepted `set00`
  parity, it may ship in the next small preview version after explicit release
  approval and matching download artifacts. Accepted `set00` is the complete
  canonical source; the importer expands it into all ten runtime slots.
- Wheel-capable runtime animals: chinchilla and hamster only

## Batch Strategy

One thread should produce more images, but for fewer animals:

- Source candidate thread: one named family only after parent approval.
- Motion thread: one species family, 1-3 accepted variants, 62 source frames per accepted variant.
- Integration thread: parent only; verifies alpha bounds, 96x64 readability, deterministic import, catalog coverage, page assets, release docs, and runtime behavior.
- Frame progress thread: one species family only; uses `cmd/auditframes` to track valid/missing/invalid standalone PNG frames without promoting unfinished or opaque ImageGen outputs.

This gives each thread enough output to compare candidates without pushing 600+
unreviewed frames through one context.

## Completion Contract

Use `docs/art-source/motion-contract.md` as the promotion and release gate. The
runtime target stays DeguDesktop-compatible: 10 runtime slots, 62 frames per
slot, 96x64 transparent frames, and the DeguDesktop frame slots preserved. One
reviewed and accepted canonical 62-frame `set00` is the normal complete source;
the importer expands it into all ten runtime sheets. A family is not complete
just because it is selectable; it is complete for a version bump only when its
source-truth art, canonical motion set, importer output, settings/tray
reflection, visual review, and QA pass.

Work one family at a time for content releases. A completed canonical `set00`
does not need a preview exception for source coverage; publication still needs
explicit release approval and matching page/download artifacts. Public
coming-soon silhouettes are release-managed: when an animal ships, remove it
from the queue image/list and move it into the current-animal section. Do not
repurpose runtime, prototype, or accepted-frame art for coming-soon silhouettes;
either generate page-specific future art or keep the future queue text-only.

## Profile-Specific Notes

- Chinchilla: heavy soft body, rounded ears, fluffy tail, gentle scurry, no degu recolor.
- Macaroni mouse / fat-tailed gerbil: compact small body, thick tail, tiny feet, low fast scurry.
- Hamster: round cheek-forward body, very short tail, cheek/groom actions, wheel capable.
- Campbell hamster: warm gray-brown dwarf hamster with visible dorsal stripe; color-only drift may be corrected when anatomy and alpha are otherwise acceptable.
- Djungarian hamster: gray-white dwarf hamster with visible dorsal stripe; do not collapse it into the same visual read as Campbell.
- Chipmunk: the next requested normal striped chipmunk must be a separate source family from `albino_chipmunk`, not a recolor.
- Rabbit: readable ears, hop contact points, stable baseline.
- Sugar glider / momonga-style runtime family: membrane silhouette, low skitter, no long airborne scene.
