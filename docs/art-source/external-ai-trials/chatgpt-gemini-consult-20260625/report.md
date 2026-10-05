# ChatGPT / Gemini Consultation Report

- Date: 2026-06-25
- Scope: advisory prompt and motion-breakdown review for four 96x64 desktop-pet sprite families
- Data class prepared: PUBLIC summary only
- External transmission: not performed in this thread
- External model answer received: none

## Local Triage Of Review Hypotheses

Because no external answer was received, the decisions below are local hypotheses derived from the existing motion plan and visual contact-sheet review. Treat them as proposed inputs for the parent thread, not as external-model findings.

### Adopt

- Generate frames 01-03 before any walk/scurry work, and use them as a family consistency gate. Cause: if subtle idle frames drift in scale, baseline, muzzle, ears, tail, or membrane, larger motion groups will drift more.
- Add a fixed "same family anchor" clause to every frame prompt: same body volume, camera distance, contact baseline, head size, eye size, muzzle shape, foot scale, and tail/ear/membrane identity as frame 00.
- Use group-specific motion verbs instead of asking for broad "animation" changes. Cause: broad motion requests tend to create new source families, props, or scene poses.
- Protect rabbit ears explicitly with "ears fully inside the 96x64 canvas with top margin" on every rabbit prompt.
- Protect macaroni mouse identity with "thick pale fat tail remains visible and attached, not a thin mouse tail" on every prompt.
- Keep sugar glider grounded except in explicitly tiny stretch frames; prefer "patagium line visible while feet stay near the shared baseline" over "gliding" or "flying."

### Hold

- Keep the full 62-slot structure for now. It matches the runtime contract, but frame prompts can compress visual variety within groups to avoid duplicates.
- Keep the side-view turn group, but reword it as "head/shoulder/body angle hint" rather than a full spin. A full turn may be too hard to keep readable at 96x64.
- Keep forage/eat/ground-check separate in the contract, but make the prompt deltas distinct: forage lowers nose, eat brings paws/muzzle close, ground-check moves one paw while the head stays more neutral.

### Reject

- Reject prompts that ask for scenery, visible food, wheel, floor contact shadow, flight scene, jump arc, or large body translation.
- Reject simple vertical bobbing as an idle or hop substitute.
- Reject species prompts that depend only on coat color. The silhouette marker must be present: hamster cheek mass, fat tail, patagium/tail, rabbit ears/hindquarters.
- Reject motion groups that crop ears, tail, paws, whiskers, or membrane even if the pose is otherwise attractive.

## Proposed Improved Prompt Rules

- Start every prompt with a hard canvas/task sentence: "Create one single right-facing side-view 96x64 desktop pet sprite frame of [species], one complete animal, transparent background."
- Immediately anchor continuity: "Match the frame-00 family anchor: same camera distance, body volume, baseline, head-to-body ratio, eye size, muzzle shape, foot scale, lighting, and soft outline."
- Describe exactly one small pose delta for frames 01-03.
- For motion groups, request a cycle phase rather than a scene: "walk phase 1/8 with front paw beginning forward while the body stays level."
- Put species identity before pose detail, and negative prompts after pose detail.
- For high-risk frames, mention protected anatomy in the positive prompt and again in the negative guard.

## Remaining Risks

- Without actual external ChatGPT/Gemini output, there are no external findings to adopt, hold, or reject.
- The contact-sheet review is based on frame 00 only; prompts for later frames may reveal new failure modes.
- Visual results must still be verified locally through transparent PNG checks, light/dark/checker review, and frame audit commands before any accepted-frame promotion.
