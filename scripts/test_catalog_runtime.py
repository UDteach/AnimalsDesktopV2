from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

from catalog_runtime import runtime_variant_ids


class RuntimeVariantParserTests(unittest.TestCase):
    def test_ignores_quoted_ids_in_comments(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            catalog = Path(temp_dir) / "catalog.go"
            catalog.write_text(
                """package catalog

var runtimeVariantIDs = []string{
\t// Temporarily held: "not_a_runtime_variant"
\t"mouse_gray",
\t"rabbit_white", // Approved.
}
""",
                encoding="utf-8",
                newline="\n",
            )

            self.assertEqual(
                runtime_variant_ids(catalog),
                ["mouse_gray", "rabbit_white"],
            )

    def test_rejects_noncanonical_entries(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            catalog = Path(temp_dir) / "catalog.go"
            catalog.write_text(
                """package catalog

var runtimeVariantIDs = []string{
\t"mouse_gray", "rabbit_white",
}
""",
                encoding="utf-8",
                newline="\n",
            )

            with self.assertRaisesRegex(ValueError, "unsupported runtime variant entry"):
                runtime_variant_ids(catalog)


if __name__ == "__main__":
    unittest.main()
