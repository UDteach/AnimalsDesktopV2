# ChatGPT Pro follow-up consultation: run04 F prompt after run03 E failure

Date: 2026-06-23

Run03 E result sent:

- Four identical E prompt samples all produced 8 detected animal components at `1774x887`.
- Raw parse stayed `0/8` for every sample because exact `#00ff00` pixels were effectively absent.
- Diagnostic greenfixed parse was only `3/8`, `5/8`, `5/8`, and `3/8`.
- Cell recentering did not materially improve the result: `4/8`, `5/8`, `5/8`, and `3/8`.
- Rejects were mainly `transparent/chroma pinholes`, with detached alpha components in some samples.
- Visual review showed coherent right-facing chinchillas, but belly/feet/tail details introduced small holes and specks.

Advice received:

- Treat E as a failed prompt direction for parse-clean batch output.
- Remove per-frame gait instructions for the next test.
- The likely trigger is wording such as `front foot slightly forward`, `passing pose`, `body bob`, `tail counter-balance`, and `visible walk cycle`, which encourages small gaps around belly, feet, and tail.
- Run F as near-identical clean review poses to test whether C's clean parse stability can be reproduced.
- If F cannot produce stable diagnostic greenfixed `8/8`, stop treating 8-frame ImageGen sheets as production candidates.

Decision rule:

- Generate F four times.
- Stop the 8-frame ImageGen sheet production-candidate lane if fewer than two of four samples reach greenfixed `8/8`, average greenfixed parse is below `7/8`, pinhole rejects continue, detached alpha recurs, or recentering does not help.
- Even on success, use 8-frame sheets only as review/layout/style/motion-reference artifacts. Accepted production frames remain one-frame-per-PNG.

