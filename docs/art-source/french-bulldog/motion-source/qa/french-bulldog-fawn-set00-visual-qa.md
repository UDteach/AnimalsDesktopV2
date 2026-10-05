# French Bulldog Fawn Set00 Visual QA

Date: 2026-06-28

Source run:

`docs/art-source/one-frame-method-fullrun-20260628/french-bulldog-fawn-set00-oneframe-62/`

Accepted source path:

`docs/art-source/french-bulldog/motion-source/accepted-frames/set00/`

## Verdict

Accepted as a reusable `set00` source asset only. This does not add the variant
to runtime variants, Pages, release notes, tags, downloads, or deploy outputs.

## Mechanical QA

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/french-bulldog/motion-source/accepted-frames/set00 -strict -artifact-warnings -report docs/art-source/french-bulldog/motion-source/accepted-frames/set00-auditframes-report.json
```

Result:

`valid=62 missing=0 invalid=0 warnings=6`

```sh
go run ./cmd/assemblemotion -frames-dir docs/art-source/french-bulldog/motion-source/accepted-frames/set00 -out docs/art-source/french-bulldog/motion-source/sheets/french-bulldog-fawn-source-set00.png -report docs/art-source/french-bulldog/motion-source/accepted-frames/set00-assemblemotion-report.json
```

Result:

`docs/art-source/french-bulldog/motion-source/sheets/french-bulldog-fawn-source-set00.png`

## Parent Visual QA

Reviewed contacts:

- `contacts/french-bulldog-fawn-set00-52-61-checker-contact.png`
- `contacts/french-bulldog-fawn-set00-45-61-checker-contact.png`
- `contacts/french-bulldog-fawn-set00-00-61-checker-contact.png`
- `contacts/french-bulldog-fawn-set00-52-61-dark-contact.png`
- `contacts/french-bulldog-fawn-set00-00-61-dark-contact.png`
- `contacts/french-bulldog-fawn-set00-00-61-light-contact.png`

Frames `53-57` form a consistent face-groom/front-facing band. Frames `58-61`
return toward the right-facing reaction/recover band without a one-frame body
size spike. The fawn coat, dark mask, upright ears, compact bulldog body,
attached paws, and baseline contact remain readable on checker, dark, and light
backgrounds.
