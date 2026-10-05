"""Read the curated runtime roster from the Go catalog without accepting comments.

The catalog intentionally keeps one quoted ID per line. This small parser fails
closed when that contract changes, which keeps page generation independent from
Go test assertions while avoiding a second hand-maintained roster.
"""

from __future__ import annotations

import re
from pathlib import Path


_RUNTIME_BLOCK = re.compile(
    r"var\s+runtimeVariantIDs\s*=\s*\[\]string\s*\{(.*?)\n\}",
    flags=re.S,
)
_RUNTIME_ENTRY = re.compile(r'^\s*"([a-z][a-z0-9_]*)",\s*(?://.*)?$')


def runtime_variant_ids(catalog_path: Path) -> list[str]:
    catalog = catalog_path.read_text(encoding="utf-8")
    blocks = _RUNTIME_BLOCK.findall(catalog)
    if len(blocks) != 1:
        raise ValueError(
            f"expected one runtime variant list in {catalog_path}, found {len(blocks)}"
        )

    ids: list[str] = []
    for line_number, raw_line in enumerate(blocks[0].splitlines(), start=1):
        stripped = raw_line.strip()
        if not stripped or stripped.startswith("//"):
            continue
        match = _RUNTIME_ENTRY.fullmatch(raw_line)
        if match is None:
            raise ValueError(
                f"unsupported runtime variant entry on block line {line_number}: "
                f"{stripped!r}"
            )
        ids.append(match.group(1))

    if not ids:
        raise ValueError(f"runtime variant list is empty in {catalog_path}")
    if len(ids) != len(set(ids)):
        raise ValueError(f"runtime variant list contains duplicate IDs in {catalog_path}")
    return ids
