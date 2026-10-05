# Gecko set00 visual QA

Date: 2026-06-26

## Verdict

Accepted as a slot-based `set00` motion source candidate and promoted to:

- `docs/art-source/gecko/motion-source/accepted-frames/set00/`
- `docs/art-source/gecko/motion-source/sheets/gecko-gray-brown-source-set00.png`

This is not judged as a single smooth 62-frame loop. It is judged with the same AnimalsDesktop/DeguDesktop slot-based motion rule used for the five existing small-animal `set00` sets.

## Mechanical QA

- Raw ImageGen outputs copied: 62/62.
- `normalize_gecko_frame.py`: 62/62 produced 96x64 RGBA frames.
- `cmd/auditframes` on fullrun frames: valid=62, missing=0, invalid=0, warnings=31.
- `cmd/auditframes` on promoted accepted frames: valid=62, missing=0, invalid=0, warnings=31.
- `cmd/assemblemotion`: wrote `gecko-gray-brown-source-set00.png`.
- `cmd/validatemotion -variant gecko_gray_brown -require-accepted`: accepted_source=true, runtime_sets=1, release_ready=false because only set00 exists.

## Matte Repair Notes

The first strict `cmd/prepareframe -background chroma-green` pass accepted 41/62 frames and rejected 21/62 for transparent pinholes. The rejected pinholes were primarily natural negative spaces from gecko toes, leg gaps, and body arches after chroma-key removal.

The run-local helper `tools/normalize_gecko_frame.py` was used instead of modifying `cmd/prepareframe`. It performs:

- chroma-green removal,
- green despill,
- largest connected animal component selection,
- internal alpha-hole fill,
- 96x64 fit with the same 88x52 target and baseline used by `cmd/prepareframe`.

Repair summary:

- internal alpha holes filled: 21 frames.
- detached component cleanup: frame 44 removed small motion-line-like marks.
- raw outputs are retained unchanged in `raw/`.
- per-frame normalize reports are under `qa/frame-NN-normalize-gecko.json`.

## Slot Review

- `00-03 idle`: pass. Stable low gecko pose with small breathing/head changes.
- `04-11 walk`: pass. Crawl gait reads as slow movement; frame 10 has a tail wave but remains plausible.
- `12-19 fast`: pass. Stronger body stretch/recoil and tail movement; slot reads as quick scamper.
- `20-25 sniff/nibble`: pass. Nose-down and small mouth/head changes read correctly without props.
- `26-31 species action`: pass with note. Reads as climbing/raising front body; frames 26-30 are larger pose deltas but appropriate for a gecko-specific action slot.
- `32-39 turn`: pass with note. Frame 34 is a clear top-view turn midpoint; surrounding frames return to right-facing side view.
- `40-43 chew/eat`: pass. Subtle mouth/head movement without adding food props.
- `44-47 ground check`: pass after detached component cleanup on frame 44.
- `48-51 alert/rest`: pass. Head/chest raise and settle are readable.
- `52-55 face groom`: pass. Grooming foot is readable; no detached paw after normalization.
- `56-61 reaction/recover`: pass with note. Frames 56-57 use high tail arcs; recovery returns to neutral by 61.

## Review Artifacts

- Full contact: `contact/gecko-set00-full-contact.png`
- Slot review: `contact/gecko-set00-slot-review.png`
- Animated preview: `contact/gecko-set00-preview.gif`
- Promoted contact copy: `docs/art-source/gecko/motion-source/contacts/gecko-gray-brown-set00-slot-review.png`
