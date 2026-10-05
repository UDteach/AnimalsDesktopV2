# ChatGPT Pro prompt consultation

Date: 2026-06-23

Purpose: improve the 8-frame chinchilla ImageGen batch prompt after the first workflow run reached `7/8` diagnostic parse but stayed too realistic and still had a `1px` matte pinhole.

Consultation packet sent:

- Review-only 8-frame sheets for a Windows desktop pet app.
- Final assets remain one-frame PNGs; this is only for fixed-cell extraction experiments.
- Prior run: `2048x1024`, invisible `2x4`, intended `512x512` cells, pure `#00ff00`.
- Stage 1: 8 separated components, raw parse `0/8`, green-normalized parse `6/8`.
- Stage 2: 8 separated components, raw parse `0/8`, green-normalized parse `7/8`, only reject was a `1px` pinhole.
- Visual failures: too realistic, weak gait deltas, background not parser-pure, tiny pinholes.

ChatGPT Pro recommendation:

1. Test `C` first to estimate extraction success ceiling.
2. Test `A` next for flatter sprite style and simpler silhouette.
3. Test `B` last for stronger gait deltas while keeping identity stable.

Local decision:

- Use A/B/C as prompt strategies.
- Run reproducibility samples through the same local parser diagnostics.
- Keep all outputs review-only and outside `accepted-frames`.
