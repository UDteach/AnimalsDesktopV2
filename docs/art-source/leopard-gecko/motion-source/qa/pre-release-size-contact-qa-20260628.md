# Leopard Gecko Pre-Release Size/Contact QA

Date: 2026-06-28

## Trigger

User review flagged that the middle turn frames looked too large.

## Repair

Original frames were preserved under:

`docs/art-source/leopard-gecko/motion-source/qa/pre-release-size-contact-repair-20260628/originals/`

Frames `33`, `34`, and `35` were mechanically cropped, scaled down, and
re-anchored to the same bottom contact line. Other frames were left unchanged.

Repair report:

`docs/art-source/leopard-gecko/motion-source/qa/pre-release-size-contact-repair-20260628/report.json`

## Verification

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/leopard-gecko/motion-source/accepted-frames/set00 -strict -artifact-warnings -motion-warnings -report docs/art-source/leopard-gecko/motion-source/accepted-frames/set00-auditframes-report.json
go run ./cmd/assemblemotion -frames-dir docs/art-source/leopard-gecko/motion-source/accepted-frames/set00 -out docs/art-source/leopard-gecko/motion-source/sheets/leopard-gecko-source-set00.png -report docs/art-source/leopard-gecko/motion-source/accepted-frames/set00-assemblemotion-report.json
```

Result:

`valid=62 missing=0 invalid=0 warnings=15`

Reviewed contact:

`docs/art-source/leopard-gecko/motion-source/contacts/leopard-gecko-set00-28-40-pre-release-repaired-checker-baseline.png`

## Verdict

Accepted for the current release lane. The center turn now reads as a controlled
arched/turning pose instead of a sudden oversized animal, and the bottom contact
line stays stable through frames `28-40`.
