# Maine Coon Brown Tabby Set00 Visual QA

Date: 2026-06-28

Source run:

`docs/art-source/one-frame-method-fullrun-20260628/maine-coon-brown-tabby-set00-oneframe-62/`

Accepted source path:

`docs/art-source/maine-coon/motion-source/accepted-frames/set00/`

## Verdict

Accepted as a reusable `set00` source asset only. This does not add the variant
to runtime variants, Pages, release notes, tags, downloads, or deploy outputs.

## Mechanical QA

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/maine-coon/motion-source/accepted-frames/set00 -strict -artifact-warnings -report docs/art-source/maine-coon/motion-source/accepted-frames/set00-auditframes-report.json
```

Result:

`valid=62 missing=0 invalid=0 warnings=21`

```sh
go run ./cmd/assemblemotion -frames-dir docs/art-source/maine-coon/motion-source/accepted-frames/set00 -out docs/art-source/maine-coon/motion-source/sheets/maine-coon-brown-tabby-source-set00.png -report docs/art-source/maine-coon/motion-source/accepted-frames/set00-assemblemotion-report.json
```

Result:

`docs/art-source/maine-coon/motion-source/sheets/maine-coon-brown-tabby-source-set00.png`

## Parent Visual QA

Reviewed contacts:

- `contacts/maine-coon-brown-tabby-set00-52-61-checker-contact.png`
- `contacts/maine-coon-brown-tabby-set00-45-61-checker-contact.png`
- `contacts/maine-coon-brown-tabby-set00-00-61-checker-contact.png`
- `contacts/maine-coon-brown-tabby-set00-52-61-dark-contact.png`
- `contacts/maine-coon-brown-tabby-set00-00-61-dark-contact.png`
- `contacts/maine-coon-brown-tabby-set00-00-61-light-contact.png`

Frames `52-54` form a consistent face-groom band. Frames `55-61` recover toward
right-facing alert/reaction poses without an isolated body-size spike. The brown
tabby stripes, longhair body, fluffy ringed tail, ears, attached paws, and
baseline contact remain readable on checker, dark, and light backgrounds.
