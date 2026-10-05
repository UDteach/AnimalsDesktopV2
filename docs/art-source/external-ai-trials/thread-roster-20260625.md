# External AI Trial Threads - 2026-06-25

All threads use the same local project folder:

`/Volumes/ZX20 II/Development/AnimalDesktop`

Working branch:

`imagegen-five-small-animals-assets`

## Threads

| Thread | ID | Scope | Write root |
| --- | --- | --- | --- |
| Google Flow image trials | `019efbc1-1c2a-75b1-8844-a6f379173d9d` | Google Flow 8/16-cell motion trials, sanitized external input only | `docs/art-source/external-ai-trials/google-flow-20260625/` |
| Codex ImageGen trials | `019efbc1-b6f1-7b32-bd44-e720ce478780` | Codex built-in image_gen 8/16-cell and possible 62-cell trials | `docs/art-source/external-ai-trials/codex-imagegen-20260625/` |
| ChatGPT Gemini consult | `019efbc3-eee2-79a1-b029-abd34b69150d` | Sanitized prompt and visual QA consultation; no direct asset promotion | `docs/art-source/external-ai-trials/chatgpt-gemini-consult-20260625/` |
| ChatGPT image trials | `019efbc4-851c-7da1-81af-204b7a5569b0` | ChatGPT image generation trials, isolated review outputs | `docs/art-source/external-ai-trials/chatgpt-imagegen-20260625/` |
| Hamster one-frame fullrun | `019efbf9-19a4-7881-9682-89730ce41e9c` | Codex built-in ImageGen one-frame retry-until-pass 62-frame run | `docs/art-source/one-frame-method-fullrun-20260625/hamster/` |
| Macaroni mouse one-frame fullrun | `019efbfa-24d0-7aa0-8db4-d2447df97aa5` | Codex built-in ImageGen one-frame retry-until-pass 62-frame run | `docs/art-source/one-frame-method-fullrun-20260625/macaroni-mouse/` |
| Sugar glider one-frame fullrun | `019efbfb-45fb-70b1-94a2-4c1ec2ce71c7` | Codex built-in ImageGen one-frame retry-until-pass 62-frame run | `docs/art-source/one-frame-method-fullrun-20260625/sugar-glider/` |
| Rabbit one-frame fullrun | `019efbfc-5bbf-71e0-96e0-160fe40066e9` | Codex built-in ImageGen one-frame retry-until-pass 62-frame run | `docs/art-source/one-frame-method-fullrun-20260625/rabbit/` |
| Hamster one-frame rescue | `019efc36-f49e-7400-a241-c813be169de6` | Codex built-in ImageGen rescue generation for hamster frames 32-61 only | `docs/art-source/one-frame-method-fullrun-20260625-rescue/hamster-32-61/` |
| Rabbit one-frame rescue | `019efc37-515b-78a3-aacc-847ecb5c246c` | Codex built-in ImageGen rescue generation for rabbit frames 40-61 only | `docs/art-source/one-frame-method-fullrun-20260625-rescue/rabbit-40-61/` |
| Macaroni mouse one-frame rescue | `019efc43-f614-78a3-8f80-c0295aae9fbc` | Codex built-in ImageGen rescue generation for macaroni mouse frames 56-61 only | `docs/art-source/one-frame-method-fullrun-20260625-rescue/macaroni-mouse-56-61/` |
| Sugar glider one-frame rescue | `019efc46-642f-7542-a942-b8fa021610f1` | Codex built-in ImageGen rescue generation for sugar glider frames 54-61 only; replaces `019efc44-3a4c-7063-850f-b3c11ac545b4`, which returned an ID but never appeared in thread lookup or created output files | `docs/art-source/one-frame-method-fullrun-20260625-rescue/sugar-glider-54-61/` |

## Parent Thread Rules

- Parent thread owns `docs/art-source/four-animals-imagegen-only-62/` and the current 62-frame candidate manifest.
- Parent thread also owns final adoption decisions between `accepted-candidates`, one-frame trials, and one-frame fullrun outputs.
- Trial threads must not write `accepted-frames`, catalog/runtime code, GitHub Pages, release assets, or external production settings.
- Trial outputs are hypotheses. Parent thread integrates only after local file inspection and QA.
- Rescue outputs are not automatic replacements. Parent thread must compare them with the main fullrun frame-by-frame before any final candidate adoption.
- Do not send secrets, tokens, local absolute paths, private repo state, or raw logs to external services.
