# Chinchilla ImageGen run05: angle blueprint experiment

Date: 2026-06-23

Purpose:

- Test whether per-frame foot/tail angle wording can preserve run02 C's extraction cleanliness while adding motion variation.
- Draft a generic species pose profile flow that can apply to chinchilla, momonga, hamster, rabbit, and gecko.

Inputs:

- `species-pose-profile-draft.md`
- `pose-blueprint-spec.json`
- `pose-blueprint-bone.png`
- `pose-blueprint-silhouette.png`
- `prompt-g-angle-blueprint.txt`

Result:

- G1 and G2 both produced 8 detected animal components.
- Exact `#00ff00` background pixels were still effectively absent before diagnostic normalization.
- Raw fixed-cell parse stayed `0/8`.
- Greenfixed fixed-cell parse stayed `0/8` because the generated animals did not honor the fixed cell boundary positions.
- Cell-recentered diagnostics reached:
  - G1: `7/8`, one pinhole reject.
  - G2: `8/8`.

Visual decision:

- Angle-table wording creates some controlled foot/tail variation.
- The output is more upright/sitting than a usable chinchilla walk loop.
- The approach is useful for motion design and species-specific prompting, not direct production extraction.

Flow decision:

- Ask a species pose profile before any multi-frame prompt.
- Use bone-only and silhouette-only guides as local planning artifacts.
- Do not upload visible line/label guides by default because they can become guide ink.
- Keep accepted production on one-frame-per-PNG generation and normal QA.

