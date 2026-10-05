# Domestic Shorthair Calico Set00 Visual QA

Date: 2026-06-28

Source run:

`docs/art-source/one-frame-method-fullrun-20260628/domestic-shorthair-calico-set00-oneframe-62/`

Accepted source path:

`docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00/`

## Verdict

Accepted as a reusable `set00` source asset only. This does not add the variant
to runtime variants, Pages, release notes, tags, downloads, or deploy outputs.

## Mechanical QA

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00 -strict -artifact-warnings -report docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00-auditframes-report.json
```

Result:

`valid=62 missing=0 invalid=0 warnings=70`

```sh
go run ./cmd/assemblemotion -frames-dir docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00 -out docs/art-source/domestic-shorthair/motion-source/sheets/domestic-shorthair-calico-source-set00.png -report docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00-assemblemotion-report.json
```

Result:

`docs/art-source/domestic-shorthair/motion-source/sheets/domestic-shorthair-calico-source-set00.png`

## Parent Visual QA

Reviewed contacts:

- `contacts/domestic-shorthair-calico-set00-53-61-checker-contact.png`
- `contacts/domestic-shorthair-calico-set00-45-61-checker-contact.png`
- `contacts/domestic-shorthair-calico-set00-00-61-checker-contact.png`
- `contacts/domestic-shorthair-calico-set00-53-61-dark-contact.png`
- `contacts/domestic-shorthair-calico-set00-00-61-dark-contact.png`
- `contacts/domestic-shorthair-calico-set00-00-61-light-contact.png`

Frames `53-61` form a consistent reaction/recover band. Frames `54-55`,
`57-58`, and `60` are the taller alert/reaction poses; frames `56`, `59`, and
`61` return toward the smaller right-facing recovery poses. Across the full set,
the calico orange/black/white patches, attached paws, tail, ears, and baseline
contact remain readable on checker, dark, and light backgrounds. No isolated
body-size spike was found.

## 2026-06-28 Pre-Release Contact QA

User review flagged that the domestic shorthair might read as floating. The
parent rechecked threshold-visible contact and found that frames `07`, `27`,
and `53` had visible body contact above the baseline even though faint alpha
extended lower.

Original frames were preserved under:

`docs/art-source/domestic-shorthair/motion-source/qa/pre-release-contact-repair-20260628/originals/`

Frames `07`, `27`, and `53` were mechanically re-anchored to bottom contact
line `y=57` without changing the rest of the set. Repair report:

`docs/art-source/domestic-shorthair/motion-source/qa/pre-release-contact-repair-20260628/report.json`

Re-run commands:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00 -strict -artifact-warnings -motion-warnings -report docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00-auditframes-report.json
go run ./cmd/assemblemotion -frames-dir docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00 -out docs/art-source/domestic-shorthair/motion-source/sheets/domestic-shorthair-calico-source-set00.png -report docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00-assemblemotion-report.json
```

Result:

`valid=62 missing=0 invalid=0 warnings=182`

Reviewed contacts:

- `contacts/domestic-shorthair-calico-set00-00-12-pre-release-repaired-checker-baseline.png`
- `contacts/domestic-shorthair-calico-set00-24-28-pre-release-repaired-checker-baseline.png`
- `contacts/domestic-shorthair-calico-set00-52-54-pre-release-repaired-checker-baseline.png`

Verdict: accepted for the current release lane. The flagged frames now visibly
touch the baseline, while the remaining warning count is mostly fine alpha,
white-body pinholes, and pose-group size differences that were visually
reviewed on the contact sheets.
