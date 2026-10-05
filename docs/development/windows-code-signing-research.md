# Windows Code Signing Research

Date: 2026-06-27

## Short Answer

Code signing is an identity-validation process, not just a file-generation step. For public Windows releases, AnimalsDesktop needs a public-trust Authenticode signing path.

The best fit is:

1. Use Microsoft Azure Artifact Signing / Trusted Signing if the publisher account is eligible.
2. Otherwise buy an OV code-signing certificate from a public CA and use that CA's cloud/HSM signing workflow.
3. Use EV only if the extra validation/cost is justified by user trust or procurement requirements; do not treat EV as a guaranteed SmartScreen bypass.

## Why PFX Is Not The Main Path

Modern public code-signing certificates generally require stronger private-key protection than an exportable `.pfx` stored in CI secrets. Some older/private/internal certificates can still be PFX-based, so the release workflow keeps PFX support as a fallback, but new public releases should prefer HSM/cloud signing.

## Practical Application Flow

### Azure Artifact Signing

1. Create or use an Azure subscription.
2. Set up Microsoft Entra app registration / federated credential for GitHub Actions OIDC.
3. Create an Artifact Signing account.
4. Complete identity validation.
5. Create a certificate profile.
6. Add GitHub Secrets:
   - `AZURE_CLIENT_ID`
   - `AZURE_TENANT_ID`
   - `AZURE_SUBSCRIPTION_ID`
   - `AZURE_ARTIFACT_SIGNING_ENDPOINT`
   - `AZURE_ARTIFACT_SIGNING_ACCOUNT_NAME`
   - `AZURE_ARTIFACT_SIGNING_CERTIFICATE_PROFILE_NAME`
7. Tag a release only after the app release gate passes.

### Public CA OV Certificate

1. Choose a CA that supports CI-friendly cloud signing or HSM-backed signing.
2. Complete organization or individual identity validation.
3. Configure the CA signing tool or service in CI.
4. Keep one stable publisher identity across releases.
5. Timestamp every signature.

## Additional Trust Work Still Worth Doing

- Publish signed EXE hashes and ZIP hashes; the workflow now publishes ZIP hashes plus inner packaged EXE hashes.
- Keep the release asset names stable (`AnimalsDesktop-windows-amd64.zip`, `AnimalsDesktop-windows-386.zip`).
- Avoid temporary scripts and suspicious command lines at runtime; the PowerShell updater was removed.
- Avoid packers/compressors such as UPX.
- Submit clean blocked builds to Microsoft and McAfee when false positives occur.
- Consider a signed installer or MSIX later if the app needs a more standard install/update channel.
- Consider Microsoft Store distribution later only if the product direction fits Store packaging and policies.

## References

- Microsoft SmartScreen reputation: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation
- Microsoft code-signing options: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/code-signing-options
- Azure Artifact Signing integrations: https://learn.microsoft.com/en-us/azure/artifact-signing/how-to-signing-integrations
- Microsoft Security Intelligence file submission: https://www.microsoft.com/en-us/wdsi/filesubmission
- McAfee Dispute Detection & Allowlisting: https://www.mcafee.com/en-us/consumer-support/dispute-detection-allowlisting.html
