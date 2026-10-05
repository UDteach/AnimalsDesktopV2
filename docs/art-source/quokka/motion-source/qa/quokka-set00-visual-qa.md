# Quokka Set00 Visual QA

Date: 2026-06-28

Source run:

`/Volumes/ZX20 II/Development/AnimalDesktop/docs/art-source/one-frame-method-fullrun-20260628/quokka-set00-oneframe-62/`

Accepted source path:

`docs/art-source/quokka/motion-source/accepted-frames/set00/`

## Verdict

Accepted as the 19th reusable `set00` source asset in the expanded release
lane.

## Rejected Attempts

The first `00-12` direction was rejected as too rodent/kangaroo-rat-like
because the tail was long and thin and the body read low like a rat. The final
run uses a warmer upright wallaby/quokka read with a short stout attached tail.

During the final block, `frame-57` attempts `01` and `02` were rejected and
replaced; the first unclaimed `frame-57` source had a too-long tail.

## Mechanical QA

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/quokka/motion-source/accepted-frames/set00 -strict -artifact-warnings -motion-warnings -report docs/art-source/quokka/motion-source/accepted-frames/set00-auditframes-report.json
go run ./cmd/assemblemotion -frames-dir docs/art-source/quokka/motion-source/accepted-frames/set00 -out docs/art-source/quokka/motion-source/sheets/quokka-source-set00.png -report docs/art-source/quokka/motion-source/accepted-frames/set00-assemblemotion-report.json
```

Result:

`valid=62 missing=0 invalid=0 warnings=90`

The warning count is accepted after contact review. Most warnings are lower
ledge/floor heuristics from the low feet/tail contact, plus expected width
changes between crouched, turning, and upright wallaby poses. The baseline
stays at `y=57` through the final block.

## Parent Visual QA

Reviewed contacts:

- `contacts/quokka-set00-53-61-checker-contact.png`
- `contacts/quokka-set00-45-61-checker-contact.png`
- `contacts/quokka-set00-00-61-checker-contact.png`
- `contacts/quokka-set00-53-61-light-contact.png`
- `contacts/quokka-set00-53-61-dark-contact.png`

Frames `00-24` keep the low moving body mass and short attached tail. Frames
`26-43` intentionally move through upright/turning wallaby poses. Frames
`45-61` recover from crouched to upright right-facing reaction poses. The final
block preserves compact rounded ears, warm brown coat, visible hind feet,
attached stout tail, and stable contact.
