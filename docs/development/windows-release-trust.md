# Windows Release Trust

This note tracks the Windows AV and SmartScreen hardening steps for AnimalsDesktop.

## Implemented In Repo

- Windows builds embed `RT_VERSION` metadata with product name, company name, file description, original filename, copyright, file version, and product version.
- Windows builds embed a Windows 10+ application manifest with `as invoker`, no `uiAccess`, no auto-elevation, common controls v6, long-path awareness, segment heap, and high-resolution scrolling awareness.
- Windows builds generate and embed an app icon from the existing preview artwork.
- The release workflow publishes `SHA256SUMS.txt` for release ZIPs and the packaged `AnimalsDesktop.exe` files.
- Each release ZIP includes `SECURITY.txt` with the expected EXE hash, expected behavior, source/release URLs, and Microsoft/McAfee false-positive submission guidance.
- The app updater verifies GitHub release asset size and `sha256:` digest when GitHub provides it.
- The runtime updater no longer writes or launches a PowerShell script and no longer uses `ExecutionPolicy Bypass`; it copies the current EXE as a temporary helper and applies the update through the app's own constrained updater mode.
- The updater helper now accepts only a source `AnimalsDesktop.exe` inside the app-owned update temp directory and a target `AnimalsDesktop.exe` outside that temp directory.
- Release builds no longer use `-s -w`, leaving richer PE/debug metadata in the Go binary.
- Optional Authenticode signing is wired into the release workflow through Microsoft Azure Artifact Signing or fallback `.pfx` GitHub Secrets.

## Release Signing

### Preferred: Azure Artifact Signing

Apply for identity validation in Microsoft Azure Artifact Signing / Trusted Signing, create a signing account and certificate profile, then configure these GitHub Secrets:

- `AZURE_CLIENT_ID`
- `AZURE_TENANT_ID`
- `AZURE_SUBSCRIPTION_ID`
- `AZURE_ARTIFACT_SIGNING_ENDPOINT`
- `AZURE_ARTIFACT_SIGNING_ACCOUNT_NAME`
- `AZURE_ARTIFACT_SIGNING_CERTIFICATE_PROFILE_NAME`

The release workflow logs in with GitHub OIDC, signs both Windows EXEs with `azure/artifact-signing-action`, and then packages the signed EXEs into ZIPs.

### Fallback: PFX Signing

Use this only for a certificate/key that is legitimately exportable, such as an internal enterprise certificate or an older/private signing setup:

- `WINDOWS_CERTIFICATE_BASE64`: base64-encoded `.pfx`
- `WINDOWS_CERTIFICATE_PASSWORD`: `.pfx` password

Public-trust OV/EV code-signing certificates issued under current industry rules usually require key protection in hardware, HSM, or a cloud signing service. Do not assume a new public certificate can be exported to a `.pfx` and stored in GitHub Secrets.

If signing secrets are absent, the workflow builds and publishes checksums, but the EXE remains unsigned. Unsigned new binaries may still trigger reputation warnings.

## Certificate Choice

- **Azure Artifact Signing / Trusted Signing:** best fit for this repository if the account is eligible. It is CI-friendly and avoids exporting private keys.
- **OV code-signing certificate from a public CA:** normal choice when Artifact Signing is unavailable. Prefer CA cloud signing or HSM-backed signing integration rather than USB-token-only flows.
- **EV code-signing certificate:** only choose it if procurement/user trust requirements justify the cost. Do not rely on EV alone as a guaranteed SmartScreen bypass.
- **Self-signed certificate:** useful only for internal testing; it does not build public SmartScreen reputation.

## External Follow-Up

- Microsoft SmartScreen reputation is based on publisher/signing reputation and file-hash reputation. Use a consistent trusted signing identity for releases and avoid changing publisher/certificate identity unnecessarily.
- If Microsoft Defender or SmartScreen incorrectly blocks a clean build, submit the exact ZIP/EXE through the Microsoft Security Intelligence submission portal as a software developer.
- If McAfee flags a clean build, submit the release package through McAfee Dispute Detection & Allowlisting.

See also: `docs/development/windows-non-cert-trust-research.md` for the non-certificate mitigations that remain useful while signing/reputation is pending.

## References

- Microsoft SmartScreen reputation: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation
- Microsoft code-signing options: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/code-signing-options
- Azure Artifact Signing integrations: https://learn.microsoft.com/en-us/azure/artifact-signing/how-to-signing-integrations
- Microsoft Security Intelligence file submission: https://www.microsoft.com/en-us/wdsi/filesubmission
- McAfee Dispute Detection & Allowlisting: https://www.mcafee.com/en-us/consumer-support/dispute-detection-allowlisting.html
- GitHub release asset digest field: https://docs.github.com/en/rest/releases/assets
