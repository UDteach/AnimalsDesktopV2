# Toy Poodle Apricot Set00 Visual QA

Date: 2026-06-28

Source run:
`/Users/kyota/.codex/worktrees/6572/AnimalDesktop/docs/art-source/one-frame-method-fullrun-20260628/toy-poodle-apricot-set00-oneframe-62/`

Parent checks:

- `cmd/auditframes -strict -artifact-warnings -motion-warnings`: `valid=62 missing=0 invalid=0 warnings=1`.
- The single warning is `frame-27` lower ledge/shelf candidate. Parent contact review found the alpha run is connected curly body/leg fur, not a floor, prop, shadow, or detached shelf.
- Full `00-61` checker contact keeps a consistent apricot Toy Poodle read, right-facing side-view except turn frames, attached paws, curly coat texture, rounded ears/muzzle, and stable baseline.
- Tail-slot `45-61` contact keeps size and body ratio stable; no abrupt growth/shrink or unnatural body fill jump was accepted.

Accepted for source promotion as `toy_poodle_apricot` set00.
