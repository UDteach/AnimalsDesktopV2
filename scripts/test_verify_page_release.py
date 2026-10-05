from __future__ import annotations

import io
import tempfile
import unittest
from contextlib import redirect_stderr
from pathlib import Path

from verify_page_release import english_animal_labels, verify_current_icons


class CurrentIconInventoryTests(unittest.TestCase):
    def test_accepts_exact_runtime_icon_set(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            icon_dir = Path(temp_dir)
            (icon_dir / "current-mouse_gray.png").touch()
            (icon_dir / "current-rabbit_white.png").touch()

            verify_current_icons(["mouse_gray", "rabbit_white"], icon_dir)

    def test_rejects_missing_and_stale_icons(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            icon_dir = Path(temp_dir)
            (icon_dir / "current-mouse_gray.png").touch()
            (icon_dir / "current-stale_variant.png").touch()

            stderr = io.StringIO()
            with redirect_stderr(stderr), self.assertRaises(SystemExit):
                verify_current_icons(["mouse_gray", "rabbit_white"], icon_dir)
            self.assertIn("rabbit_white", stderr.getvalue())
            self.assertIn("stale_variant", stderr.getvalue())


class EnglishAnimalLabelTests(unittest.TestCase):
    def test_extracts_labels_in_page_order(self) -> None:
        html = """
        <script>
          const animalTextEn = [
            "Sable panda ferret", "Sable ferret", "Albino ferret",
          ];
        </script>
        """

        self.assertEqual(
            english_animal_labels(html),
            ["Sable panda ferret", "Sable ferret", "Albino ferret"],
        )


if __name__ == "__main__":
    unittest.main()
