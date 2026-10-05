from __future__ import annotations

import io
import json
import tempfile
import unittest
from contextlib import redirect_stdout
from pathlib import Path
from unittest.mock import patch

from PIL import Image

from audit_frame_matte import audit_frame, audit_frames_dir, main, render_report


FRAME_SIZE = (96, 64)
LOW_ALPHA_GREEN = (142, 170, 142, 9)


class FrameMatteAuditTests(unittest.TestCase):
    def make_frame(
        self,
        temp_dir: str,
        pixels: dict[tuple[int, int], tuple[int, int, int, int]],
    ) -> Path:
        path = Path(temp_dir) / "frame-00.png"
        image = Image.new("RGBA", FRAME_SIZE, (0, 0, 0, 0))
        for point, color in pixels.items():
            image.putpixel(point, color)
        image.save(path)
        return path

    def test_accepts_small_low_alpha_antialias_residue(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            path = self.make_frame(
                temp_dir,
                {
                    (20, 20): LOW_ALPHA_GREEN,
                    (21, 20): LOW_ALPHA_GREEN,
                },
            )

            result = audit_frame(path)

        self.assertTrue(result["passed"])
        self.assertEqual(result["greenish_pixels"], 2)
        self.assertEqual(result["largest_greenish_component"], 2)
        self.assertLessEqual(result["max_composited_green_excess_level"], 1)

    def test_rejects_more_than_eight_greenish_pixels(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            path = self.make_frame(
                temp_dir,
                {
                    (4 + column * 3, 10 + row * 3): LOW_ALPHA_GREEN
                    for row in range(3)
                    for column in range(3)
                },
            )

            result = audit_frame(path)

        self.assertFalse(result["passed"])
        self.assertEqual(result["greenish_pixels"], 9)
        self.assertTrue(
            any("greenish pixels 9 > 8" in failure for failure in result["failures"])
        )

    def test_accepts_raw_excess_above_one_when_rounded_level_is_one(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            path = self.make_frame(
                temp_dir,
                {(20, 20): (76, 102, 76, 10)},
            )

            result = audit_frame(path)

        self.assertTrue(result["passed"])
        self.assertEqual(result["max_composited_green_excess_level"], 1)
        self.assertAlmostEqual(
            result["max_composited_green_excess_raw"],
            1.019608,
            places=6,
        )

    def test_rejects_composited_green_excess_rounding_to_two(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            path = self.make_frame(
                temp_dir,
                {(20, 20): (100, 140, 100, 10)},
            )

            result = audit_frame(path)

        self.assertFalse(result["passed"])
        self.assertEqual(result["max_composited_green_excess_level"], 2)
        self.assertTrue(
            any(
                "rounded composited green excess 2 > 1" in failure
                for failure in result["failures"]
            )
        )

    def test_rejects_three_pixel_connected_streak(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            path = self.make_frame(
                temp_dir,
                {
                    (20, 20): LOW_ALPHA_GREEN,
                    (21, 21): LOW_ALPHA_GREEN,
                    (22, 22): LOW_ALPHA_GREEN,
                },
            )

            result = audit_frame(path)

        self.assertFalse(result["passed"])
        self.assertEqual(result["largest_greenish_component"], 3)
        self.assertTrue(
            any(
                "largest greenish component 3 > 2" in failure
                for failure in result["failures"]
            )
        )

    def test_directory_json_is_stable_and_cli_returns_failure(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            directory = Path(temp_dir)
            first = self.make_frame(
                temp_dir,
                {(20, 20): (100, 200, 100, 10)},
            )
            first.rename(directory / "frame-01.png")
            self.make_frame(temp_dir, {(20, 20): LOW_ALPHA_GREEN})

            report = audit_frames_dir(directory)
            first_json = render_report(report, "json")
            second_json = render_report(audit_frames_dir(directory), "json")
            stdout = io.StringIO()
            with (
                patch(
                    "sys.argv",
                    [
                        "audit_frame_matte.py",
                        "--frames-dir",
                        str(directory),
                        "--format",
                        "json",
                    ],
                ),
                redirect_stdout(stdout),
            ):
                exit_code = main()

        self.assertEqual(first_json, second_json)
        self.assertEqual(
            [frame["file"] for frame in report["frames"]],
            ["frame-00.png", "frame-01.png"],
        )
        self.assertEqual(exit_code, 1)
        self.assertFalse(json.loads(stdout.getvalue())["passed"])


if __name__ == "__main__":
    unittest.main()
