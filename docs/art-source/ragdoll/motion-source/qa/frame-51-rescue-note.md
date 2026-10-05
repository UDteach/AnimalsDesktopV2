# Ragdoll Seal Bicolor Frame 51 Rescue Note

Date: 2026-06-28

Scope: `ragdoll_seal_bicolor` frame 51 only. Frames 00-50 were not edited, and frame 52 was not generated.

Canonical output:
- `frames/frame-51.png`
- Source candidate: `deferred/frame-51-attempt-05-paused-not-canonical.png`
- Final bbox: `(4,18)-(92,58)`, size `88x40`, top `18`, baseline `57`

Visual gate:
- Pass. Frame 51 is a low-mid shoulder-rise / half-stand bridge from frame 50.
- It stays close to accepted frame 37's height band without becoming the rejected tall upright pose.
- Seal bicolor markings, dark ears/face mask/tail, pale fluffy body, white paws, long tail, and right-facing side-view are preserved.

Local repair:
- A single internal transparent pixel at `(70,53)` was filled using the run-local `tools/fill_alpha_holes.py`.
- Report: `qa/frame-51-canonical-pinhole-repair.json`
- Bbox stayed unchanged after repair.

Contact:
- `contact/ragdoll-seal-bicolor-set00-frames-49-51-checker-contact.png`
- `contact/ragdoll-seal-bicolor-set00-frames-49-51-light-contact.png`
- `contact/ragdoll-seal-bicolor-set00-frames-49-51-dark-contact.png`
- Metrics: `qa/frame-49-51-contact-metrics.json`

Audit:
- Command: `go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260627/ragdoll-seal-bicolor-set00-oneframe-62/frames -strict -artifact-warnings`
- Result: `valid=52 missing=10 invalid=0 warnings=12`, `audit_exit=1`
- The nonzero strict exit is expected for this child lane because frames 52-61 are still intentionally missing.
- JSON report: `qa/frame-51-auditframes.json`
- Text transcript: `qa/frame-51-auditframes.txt`
