# Windows Security-Check Edition

Date: 2026-07-02

This note defines the Windows no-network security-check edition of
AnimalsDesktop. It is intended to reduce avoidable reputation signals and make
the binary easier to explain when Microsoft Smart App Control, Microsoft
Defender SmartScreen, or antivirus review flags a fresh unsigned build.

## What This Edition Changes

- Built with `-tags animalsdesktop_nonetwork`.
- Automatic startup update checks are disabled.
- Tray update checks are disabled.
- Update ZIP downloads and update installation are disabled.
- The Go `net/http` update fetch/download implementation is excluded from the
  no-network build.
- The release ZIP name is stable:
  `AnimalsDesktop-windows-amd64-no-network.zip`.
- The release `SECURITY.txt` marks the package as
  `windows-amd64-no-network` and records the inner executable hash.

## What This Edition Does Not Solve

This is not a Smart App Control bypass. Microsoft documents that Smart App
Control uses cloud-powered safety prediction first. When it cannot make a
confident prediction, it checks whether the app has a valid signature; unsigned
or invalidly signed apps can still be treated as untrusted and blocked.

GlobalSign's April 2026 Japanese case study is consistent with this diagnosis:
Smart App Control can block unrecognized unsigned apps without giving the user a
run-anyway prompt, while code signing lets the app enter a more favorable trust
evaluation path. Treat it as supporting context from a certificate vendor, not
as the primary specification.

SmartScreen is also reputation-based. Microsoft documents that EV certificates
no longer provide an automatic SmartScreen bypass. Stable trusted signing and
reputation building are still the stronger path for public releases.

## Release Guidance

Use this edition as an additional Windows 64-bit download when users report
security-product friction with the normal build. Keep the normal build
available because it preserves update checks for users who want them.

For the best long-term result, continue with public-trust Authenticode signing,
preferably Azure Artifact Signing / Trusted Signing as already wired in the
release workflow.

## References

- Smart App Control FAQ: https://support.microsoft.com/en-us/windows/security/threat-malware-protection/smart-app-control-frequently-asked-questions
- SmartScreen reputation: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation
- Code signing options: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/code-signing-options
- GlobalSign Smart App Control case study: https://college.globalsign.com/blog/smartappcontrol_20260423/
