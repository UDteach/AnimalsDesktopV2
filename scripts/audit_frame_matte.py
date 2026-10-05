#!/usr/bin/env python3
"""Audit chroma-green matte residue in transparent 96x64 animal frames."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path
from typing import Any

from PIL import Image, UnidentifiedImageError


FRAME_SIZE = (96, 64)
MIN_VISIBLE_ALPHA = 8
MIN_GREEN = 90
MIN_GREEN_DOMINANCE = 25
MAX_GREENISH_PIXELS = 8
MAX_COMPOSITED_GREEN_EXCESS_LEVEL = 1
MAX_GREENISH_COMPONENT = 2


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description=(
            "Reject visible or coherent chroma-green matte residue in top-level "
            "96x64 PNG frames."
        )
    )
    parser.add_argument("--frames-dir", required=True, type=Path)
    parser.add_argument(
        "--format",
        choices=("text", "json"),
        default="text",
        dest="output_format",
    )
    parser.add_argument(
        "--output",
        type=Path,
        help="write the report here instead of stdout",
    )
    return parser.parse_args()


def is_greenish(pixel: tuple[int, int, int, int]) -> bool:
    """Return whether a visible pixel matches the chroma-green residue probe."""

    red, green, blue, alpha = pixel
    return (
        alpha > MIN_VISIBLE_ALPHA
        and green >= MIN_GREEN
        and green - red >= MIN_GREEN_DOMINANCE
        and green - blue >= MIN_GREEN_DOMINANCE
    )


def composited_green_excess_raw(pixel: tuple[int, int, int, int]) -> float:
    """Return unrounded alpha-weighted green dominance."""

    red, green, blue, alpha = pixel
    return max(0, green - max(red, blue)) * alpha / 255


def composited_green_excess_level(pixel: tuple[int, int, int, int]) -> int:
    """Return rounded 8-bit green dominance after alpha compositing."""

    red, green, blue, alpha = pixel
    excess = max(0, green - max(red, blue))
    return (excess * alpha + 127) // 255


def largest_component(points: set[tuple[int, int]]) -> int:
    """Return the largest 8-connected greenish-pixel component."""

    remaining = set(points)
    largest = 0
    while remaining:
        start = min(remaining, key=lambda point: (point[1], point[0]))
        remaining.remove(start)
        stack = [start]
        size = 0
        while stack:
            x, y = stack.pop()
            size += 1
            for delta_y in (-1, 0, 1):
                for delta_x in (-1, 0, 1):
                    if delta_x == 0 and delta_y == 0:
                        continue
                    neighbor = (x + delta_x, y + delta_y)
                    if neighbor in remaining:
                        remaining.remove(neighbor)
                        stack.append(neighbor)
        largest = max(largest, size)
    return largest


def audit_frame(path: Path) -> dict[str, Any]:
    failures: list[str] = []
    try:
        with Image.open(path) as loaded:
            loaded.load()
            mode = loaded.mode
            size = loaded.size
            frame = loaded.convert("RGBA")
    except (OSError, UnidentifiedImageError) as error:
        return {
            "file": path.name,
            "passed": False,
            "failures": [f"unreadable PNG: {error}"],
        }

    if size != FRAME_SIZE:
        failures.append(
            f"size {size[0]}x{size[1]} != {FRAME_SIZE[0]}x{FRAME_SIZE[1]}"
        )

    greenish: dict[tuple[int, int], tuple[int, int, int, int]] = {}
    for y in range(frame.height):
        for x in range(frame.width):
            pixel = frame.getpixel((x, y))
            if is_greenish(pixel):
                greenish[(x, y)] = pixel

    greenish_count = len(greenish)
    max_excess_raw = max(
        (composited_green_excess_raw(pixel) for pixel in greenish.values()),
        default=0.0,
    )
    max_excess_level = max(
        (composited_green_excess_level(pixel) for pixel in greenish.values()),
        default=0,
    )
    component_size = largest_component(set(greenish))

    if greenish_count > MAX_GREENISH_PIXELS:
        failures.append(
            f"greenish pixels {greenish_count} > {MAX_GREENISH_PIXELS}"
        )
    if max_excess_level > MAX_COMPOSITED_GREEN_EXCESS_LEVEL:
        failures.append(
            "rounded composited green excess "
            f"{max_excess_level} > {MAX_COMPOSITED_GREEN_EXCESS_LEVEL} "
            f"(raw max {max_excess_raw:.6f})"
        )
    if component_size > MAX_GREENISH_COMPONENT:
        failures.append(
            f"largest greenish component {component_size} > "
            f"{MAX_GREENISH_COMPONENT}"
        )

    return {
        "file": path.name,
        "mode": mode,
        "size": list(size),
        "greenish_pixels": greenish_count,
        "max_composited_green_excess_raw": round(max_excess_raw, 6),
        "max_composited_green_excess_level": max_excess_level,
        "largest_greenish_component": component_size,
        "passed": not failures,
        "failures": failures,
    }


def audit_frames_dir(frames_dir: Path) -> dict[str, Any]:
    errors: list[str] = []
    if not frames_dir.is_dir():
        paths: list[Path] = []
        errors.append(f"frames directory does not exist: {frames_dir}")
    else:
        paths = sorted(
            (path for path in frames_dir.glob("*.png") if path.is_file()),
            key=lambda path: path.name,
        )
        if not paths:
            errors.append(f"no top-level PNG frames found: {frames_dir}")

    frames = [audit_frame(path) for path in paths]
    failed_frames = sum(not frame["passed"] for frame in frames)
    passed = not errors and failed_frames == 0
    return {
        "schema_version": 1,
        "passed": passed,
        "frame_count": len(frames),
        "failed_frame_count": failed_frames,
        "errors": errors,
        "contract": {
            "frame_size": list(FRAME_SIZE),
            "greenish_probe": {
                "alpha_greater_than": MIN_VISIBLE_ALPHA,
                "green_at_least": MIN_GREEN,
                "green_minus_red_at_least": MIN_GREEN_DOMINANCE,
                "green_minus_blue_at_least": MIN_GREEN_DOMINANCE,
                "connectivity": 8,
            },
            "limits": {
                "greenish_pixels_per_frame": MAX_GREENISH_PIXELS,
                "composited_green_excess_level_per_pixel": (
                    MAX_COMPOSITED_GREEN_EXCESS_LEVEL
                ),
                "greenish_component_pixels": MAX_GREENISH_COMPONENT,
            },
            "composited_green_excess_rounding": (
                "(max(0, G - max(R, B)) * A + 127) // 255"
            ),
            "visual_contact_review_required": True,
        },
        "frames": frames,
    }


def format_text(report: dict[str, Any]) -> str:
    verdict = "PASS" if report["passed"] else "FAIL"
    lines = [
        "audit_frame_matte: "
        f"{verdict} frames={report['frame_count']} "
        f"failed={report['failed_frame_count']}"
    ]
    lines.extend(f"error: {error}" for error in report["errors"])
    for frame in report["frames"]:
        if frame["passed"]:
            continue
        lines.append(
            f"{frame['file']}: FAIL; " + "; ".join(frame["failures"])
        )
    if report["passed"]:
        lines.append(
            "mechanical matte gate passed; checker/light/dark parent review "
            "is still required"
        )
    return "\n".join(lines) + "\n"


def render_report(report: dict[str, Any], output_format: str) -> str:
    if output_format == "json":
        return json.dumps(
            report,
            ensure_ascii=False,
            indent=2,
            sort_keys=True,
        ) + "\n"
    return format_text(report)


def write_report(text: str, output: Path | None) -> None:
    if output is None:
        sys.stdout.write(text)
        return
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(text, encoding="utf-8", newline="\n")


def main() -> int:
    args = parse_args()
    report = audit_frames_dir(args.frames_dir)
    write_report(render_report(report, args.output_format), args.output)
    return 0 if report["passed"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
