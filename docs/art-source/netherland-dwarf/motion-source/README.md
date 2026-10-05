# Netherland Dwarf motion source

`netherland_dwarf_chestnut` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/netherland-dwarf-chestnut-source-set00.png`
- Full contact: `contacts/netherland-dwarf-gray-chestnut-set00-full-contact.png`
- Animated preview: `contacts/netherland-dwarf-gray-chestnut-set00-preview.gif`
- Visual QA: `qa/netherland-dwarf-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant netherland_dwarf_chestnut -require-accepted
```

This canonical `set00` is the complete accepted 62-frame source. The importer
expands it into all ten runtime slots; publication still requires the separate
runtime, build, page, packaging, and release-approval gates.
