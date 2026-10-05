# 白黒長毛ハムスター2 Motion Source

`longhair_hamster_black_white_masked` has an accepted `set00` source generated
with the one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/longhair-hamster-black-white-masked-source-set00.png`
- Full contact: `contacts/longhair-hamster-black-white-masked-set00-full-contact.png`
- Candidate sheet: `contacts/longhair-hamster-black-white-masked-source-set00-candidate-sheet.png`
- Animated preview: `contacts/longhair-hamster-black-white-masked-set00-preview.gif`
- Local final review: `qa/longhair-hamster-black-white-masked-local-final-review-20260703.md`

Visual target: a second black/white longhair hamster type with black ears,
black crown/head cap, dark smoky black-to-gray long back coat, broad white face
blaze, white muzzle/chest/belly/lower sides, black eyes, and pink nose. This is
separate from the earlier spotted/panda-like `longhair_hamster_black_white`
source family.

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/longhair-hamster-black-white-masked/motion-source/accepted-frames/set00 -strict -artifact-warnings -motion-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/longhair-hamster-black-white-masked/motion-source/accepted-frames/set00 -out docs/art-source/longhair-hamster-black-white-masked/motion-source/sheets/longhair-hamster-black-white-masked-source-set00.png -report docs/art-source/longhair-hamster-black-white-masked/motion-source/accepted-frames/set00-assemblemotion-report.json
```
