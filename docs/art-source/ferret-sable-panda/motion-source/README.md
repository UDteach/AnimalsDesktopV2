# Ferret Sable Panda motion source

Status: parent-approved `set00`, accepted locally on 2026-07-24.

- Variant: `ferret_sable_panda`
- Motion profile: `ferret-slink`
- Frames: 62 standalone `96x64` transparent PNGs
- Full production archive:
  `E:/Development/AnimalDesktop-asset-archives/20260724-ferret-three-variants/ferret-sable-panda-set00-oneframe-62/`
- Source sheet SHA-256:
  `8870af7cbf98c0111b0e8064f34aca775feaf9839339dea7fbed21e2c7c8eabd`
- Accepted-PNG byte aggregate SHA-256:
  `fa8fff8d45c6ee27878f2dba31e688bed13fdf5fc1bd8371bc06fd242e590f9a`
- Original parent-gate path/hash aggregate SHA-256:
  `7947c23e4fcbcf304233cbdb09231584a9e7c479c1ee404a6712816087977ca9`

Every accepted frame was produced by its own built-in ImageGen call and copied
byte-for-byte from the parent-approved canonical set. The final action-aware
audit reports `valid=62`, `missing=0`, `invalid=0`, `warnings=0`; the shared
matte audit passes `62/62`; and the exact visual-duplicate scan reports no
pairs.

The `contacts/` directory contains the final full-set, rest-loop, alert-loop,
and turn-exit visual evidence. `provenance/manifest.csv` is a 62-row
winner-only index that maps each accepted frame to its prompt, built-in
ImageGen call, raw file, alpha candidate, final canonical winner, and parent
gate. Intermediate candidate statuses are retained separately from the final
acceptance basis, including the turn-sequence requalification history.
`qa/canonical-verification-set00.json` directly records all 62 accepted paths,
file hashes, dimensions, alpha bounds, the source sheet, both aggregate
methods, and the archived baseline match. Relative source paths are rooted at
the archived family above. Local acceptance does not authorize publication, a
release, tag, or push.
