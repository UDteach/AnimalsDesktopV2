# Windows Non-Certificate Trust Research

Date: 2026-06-27

## Scope

This note covers actions that help Microsoft Defender, SmartScreen, and McAfee outcomes without relying on a code-signing certificate. These do not replace Authenticode signing or publisher reputation, but they reduce avoidable false-positive signals and make vendor review easier.

## Findings

- SmartScreen is reputation-based. Without an established trusted publisher identity, each new unsigned binary may still need file/download reputation. Certificate-free changes can reduce suspicious signals, but cannot guarantee that SmartScreen will stop warning.
- Vendor submissions need exact artifacts. Publish stable release URLs, stable asset names, ZIP hashes, and inner EXE hashes so the same bytes can be submitted to Microsoft and McAfee.
- Avoid suspicious runtime behavior. Temporary PowerShell scripts, `ExecutionPolicy Bypass`, packers, obfuscators, and hidden update flows all make benign apps harder to classify.
- Keep update behavior narrow and user-triggered. AnimalsDesktop update installation should continue to be started from the tray, download only GitHub release ZIPs, verify available release metadata, and replace only the expected `AnimalsDesktop.exe` target.
- Make the binary self-describing. Product name, company name, file description, original filename, app manifest, and icon are not reputation substitutes, but they improve the file's inspection surface and user-facing identity.
- Use false-positive channels quickly. If a clean release is flagged, submit the exact ZIP/EXE to Microsoft Security Intelligence and McAfee Dispute Detection & Allowlisting; include the release URL, SHA256 values, app purpose, and expected behavior from `SECURITY.txt`.

## Implemented In This Iteration

- Release packages now include `SECURITY.txt` beside `AnimalsDesktop.exe` and `README.md`.
- `SHA256SUMS.txt` now contains both release ZIP hashes and inner packaged EXE hashes.
- Updater argument parsing now rejects:
  - a source EXE outside the app-owned update temp directory
  - a target EXE inside the update temp directory
  - renamed targets that are not `AnimalsDesktop.exe`
- Regression tests cover the stricter updater path constraints and fixed ZIP extraction path.

## Still Not Solved Without Certificates

- Unsigned binaries can still receive SmartScreen reputation prompts.
- A new unsigned release hash may need fresh reputation or vendor review even if the previous release was clean.
- McAfee and Microsoft classification can change after publication, so vendor submissions remain an operational step rather than a one-time code change.

## Operational Checklist

1. Build with embedded Windows metadata, manifest, icon, and no `-s -w` stripping.
2. Publish stable ZIP names:
   - `AnimalsDesktop-windows-amd64.zip`
   - `AnimalsDesktop-windows-386.zip`
3. Include `README.md` and `SECURITY.txt` inside each ZIP.
4. Publish `SHA256SUMS.txt` with ZIP and inner EXE hashes.
5. Run Microsoft Defender custom scan on the local built EXE before release.
6. If blocked, submit the exact ZIP/EXE to Microsoft and McAfee with the release URL and hashes.

## References

- Microsoft SmartScreen reputation: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation
- Microsoft Security Intelligence file submission: https://www.microsoft.com/en-us/wdsi/filesubmission
- Microsoft unwanted software criteria: https://learn.microsoft.com/en-us/unified-secops/criteria
- McAfee Dispute Detection & Allowlisting: https://www.mcafee.com/en-us/consumer-support/dispute-detection-allowlisting.html
- GitHub release asset digest field: https://docs.github.com/en/rest/releases/assets
