# One-Frame Method Trials - 2026-06-25

This directory stores one-frame-at-a-time ImageGen trials separately from the current sheet-extracted `accepted-candidates` baseline.

Evaluation target: decide whether per-frame retry-until-pass generation should replace or selectively repair the current four-family 62-frame candidate set.

Success means a frame has:

- a saved raw ImageGen source,
- an alpha-converted source,
- a normalized 96x64 transparent PNG,
- a light/dark/checker contact preview,
- a QA JSON record,
- a manifest row with prompt intent, attempt number, and decision.

This method is considered more operationally reproducible only if failed frame numbers can be retried independently without causing style drift across the rest of the set.
