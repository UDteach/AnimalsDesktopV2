# Guinea Pig Tricolor set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/guinea-pig/motion-source/accepted-frames/set00/`
- `docs/art-source/guinea-pig/motion-source/sheets/guinea-pig-tricolor-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same
AnimalsDesktop/DeguDesktop slot-based motion rule used for the existing small
animal `set00` sources.

## Mechanical QA

- Raw one-frame run has 62/62 normalized candidate frames.
- `oneframe-review.json`: 62 frames, 0 mechanical issues.
- `cmd/auditframes` on promoted accepted frames: valid=62, missing=0, invalid=0, warnings=87.
- `cmd/assemblemotion`: wrote `guinea-pig-tricolor-source-set00.png`.
- `cmd/validatemotion -variant guinea_pig_tricolor -require-accepted`: accepted_source=true, runtime_sets=1, release_ready=false because only set00 exists.

## Slot Review

- `00-03 idle`: pass. Stable compact cavy body with no visible tail.
- `04-11 walk`: pass. Small waddle-step changes read while maintaining species silhouette.
- `12-19 fast`: pass. Body stretch and foot variation are readable without turning into a mouse or hamster.
- `20-25 sniff/nibble`: pass. Nose-down and compact foraging poses remain prop-free.
- `26-31 species action`: pass. Posture shifts stay inside the cavy silhouette.
- `32-39 turn`: pass with note. The turn phase uses frontal/angled cues but returns to side view.
- `40-43 chew/eat`: pass. Subtle mouth/head motion without food props.
- `44-47 ground check`: pass. Low body and feet remain intact.
- `48-51 alert/rest`: pass. Alert lift is readable and returns to neutral.
- `52-55 groom`: pass. Grooming poses remain one complete animal.
- `56-61 reaction/recover`: pass. Motion is stronger, but the animal remains complete and recovers by frame 61.

## Tricolor Pattern Audit

Date: 2026-06-26

Additional parent-side visual audit was performed after the variant was prepared
for preview-runtime integration.

Artifacts:

- `contacts/guinea-pig-tricolor-set00-pattern-audit-2x.png`
- `contacts/guinea-pig-tricolor-pattern-detail-4x.png`
- `qa/guinea-pig-tricolor-pattern-metrics.json`

Visual verdict: pass.

- The black rear/back patch, orange middle band, white underside, and dark/orange
  head patch remain readable across all 62 frames.
- No frame changes into a different coat, loses the tricolor read, or becomes a
  plain white/brown/black animal.
- `00-31` are especially stable in side-view profile.
- `32-34` intentionally move through a front/turn read; the head markings rotate
  but stay consistent with the same tricolor animal.
- `52-61` have the largest pose changes. Frames `53-55` make the black rear look
  larger and the orange band more diagonal, but this reads as posture/grooming
  change rather than a coat-pattern break.
- No cropped ears, feet, nose, body, or tail-like artifact was found in the
  enlarged review sheets.

Color-cluster sanity check found no rough ratio outliers for black, orange, or
white/cream markings. This check is supporting evidence only; the acceptance
decision is based on visual review.

## Rejected Alternate

The `guinea-pig-tricolor-template-lock-set00-62` pass is rejected for accepted
promotion. Its late frames, especially 53-61, introduce blurred/cut-out body
damage and partial silhouette loss. The accepted source uses the original
`guinea-pig-tricolor-set00-oneframe-62` candidate frames instead.

## Review Artifacts

- Full contact: `contacts/guinea-pig-tricolor-set00-full-contact.png`
- Animated preview: `contacts/guinea-pig-tricolor-set00-preview.gif`
- One-frame QA copy: `qa/guinea-pig-tricolor-oneframe-review.json`
