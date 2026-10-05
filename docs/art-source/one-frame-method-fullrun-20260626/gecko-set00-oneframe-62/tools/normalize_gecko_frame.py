#!/usr/bin/env python3
"""Normalize one gecko ImageGen raw frame to a 96x64 transparent candidate.

This is a run-local matte repair helper for the gecko set00 candidate. It keeps
the raw ImageGen output unchanged, removes chroma-green background pixels,
keeps the largest connected animal component, fills internal alpha holes, and
fits the result into the same 96x64 frame contract used by cmd/prepareframe.
"""

from __future__ import annotations

import argparse
import json
from collections import deque
from pathlib import Path

from PIL import Image

FRAME_W = 96
FRAME_H = 64
TARGET_W = 88
TARGET_H = 52
BASELINE_Y = 58


def is_chroma_green(r: int, g: int, b: int) -> bool:
    return g >= 90 and (g - r) >= 25 and (g - b) >= 25


def alpha_bbox(mask: list[list[bool]]) -> tuple[int, int, int, int] | None:
    h = len(mask)
    w = len(mask[0]) if h else 0
    min_x, min_y = w, h
    max_x, max_y = -1, -1
    for y, row in enumerate(mask):
        for x, visible in enumerate(row):
            if not visible:
                continue
            min_x = min(min_x, x)
            min_y = min(min_y, y)
            max_x = max(max_x, x)
            max_y = max(max_y, y)
    if max_x < min_x or max_y < min_y:
        return None
    return min_x, min_y, max_x + 1, max_y + 1


def neighbors8(x: int, y: int, w: int, h: int):
    for dy in (-1, 0, 1):
        for dx in (-1, 0, 1):
            if dx == 0 and dy == 0:
                continue
            nx, ny = x + dx, y + dy
            if 0 <= nx < w and 0 <= ny < h:
                yield nx, ny


def keep_largest_component(mask: list[list[bool]]) -> tuple[list[list[bool]], dict[str, int]]:
    h = len(mask)
    w = len(mask[0]) if h else 0
    seen = [[False] * w for _ in range(h)]
    components: list[list[tuple[int, int]]] = []
    for y in range(h):
        for x in range(w):
            if seen[y][x] or not mask[y][x]:
                continue
            q = deque([(x, y)])
            seen[y][x] = True
            comp: list[tuple[int, int]] = []
            while q:
                px, py = q.popleft()
                comp.append((px, py))
                for nx, ny in neighbors8(px, py, w, h):
                    if seen[ny][nx] or not mask[ny][nx]:
                        continue
                    seen[ny][nx] = True
                    q.append((nx, ny))
            components.append(comp)
    if not components:
        return mask, {"component_count": 0, "kept_area": 0, "removed_area": 0}

    components.sort(key=len, reverse=True)
    kept = set(components[0])
    out = [[False] * w for _ in range(h)]
    for x, y in kept:
        out[y][x] = True
    removed = sum(len(c) for c in components[1:])
    return out, {
        "component_count": len(components),
        "kept_area": len(components[0]),
        "removed_area": removed,
    }


def fill_internal_holes(
    rgba: Image.Image, mask: list[list[bool]]
) -> tuple[list[list[bool]], dict[str, int]]:
    h = len(mask)
    w = len(mask[0]) if h else 0
    bbox = alpha_bbox(mask)
    if bbox is None:
        return mask, {"holes_filled": 0, "hole_area": 0, "largest_hole": 0}
    min_x, min_y, max_x, max_y = bbox
    seen = [[False] * w for _ in range(h)]
    filled = [row[:] for row in mask]
    holes = 0
    hole_area = 0
    largest = 0
    px = rgba.load()

    for y in range(min_y, max_y):
        for x in range(min_x, max_x):
            if seen[y][x] or mask[y][x]:
                continue
            q = deque([(x, y)])
            seen[y][x] = True
            comp: list[tuple[int, int]] = []
            touches = False
            while q:
                cx, cy = q.popleft()
                comp.append((cx, cy))
                if cx in (min_x, max_x - 1) or cy in (min_y, max_y - 1):
                    touches = True
                for nx, ny in neighbors8(cx, cy, w, h):
                    if nx < min_x or nx >= max_x or ny < min_y or ny >= max_y:
                        continue
                    if seen[ny][nx] or mask[ny][nx]:
                        continue
                    seen[ny][nx] = True
                    q.append((nx, ny))
            if touches:
                continue

            boundary_colors: list[tuple[int, int, int]] = []
            comp_set = set(comp)
            for cx, cy in comp:
                for nx, ny in neighbors8(cx, cy, w, h):
                    if (nx, ny) in comp_set or not mask[ny][nx]:
                        continue
                    r, g, b, a = px[nx, ny]
                    if a:
                        boundary_colors.append((r, g, b))
            if boundary_colors:
                r = round(sum(c[0] for c in boundary_colors) / len(boundary_colors))
                g = round(sum(c[1] for c in boundary_colors) / len(boundary_colors))
                b = round(sum(c[2] for c in boundary_colors) / len(boundary_colors))
            else:
                r, g, b = 154, 132, 111
            for cx, cy in comp:
                filled[cy][cx] = True
                px[cx, cy] = (r, g, b, 255)
            holes += 1
            hole_area += len(comp)
            largest = max(largest, len(comp))

    return filled, {"holes_filled": holes, "hole_area": hole_area, "largest_hole": largest}


def fit_to_frame(rgba: Image.Image, mask: list[list[bool]]) -> tuple[Image.Image, dict[str, object]]:
    bbox = alpha_bbox(mask)
    if bbox is None:
        raise SystemExit("no visible alpha after chroma removal")
    min_x, min_y, max_x, max_y = bbox
    crop = rgba.crop((min_x, min_y, max_x, max_y))
    scale = min(TARGET_W / crop.width, TARGET_H / crop.height, 1.0)
    new_w = max(1, round(crop.width * scale))
    new_h = max(1, round(crop.height * scale))
    resized = crop.resize((new_w, new_h), Image.Resampling.LANCZOS)
    out = Image.new("RGBA", (FRAME_W, FRAME_H), (0, 0, 0, 0))
    x = (FRAME_W - new_w) // 2
    y = BASELINE_Y - new_h
    if y < 0:
        y = 0
    out.alpha_composite(resized, (x, y))
    out_alpha = out.getchannel("A")
    out_bbox = out_alpha.getbbox()
    return out, {
        "source_content": [min_x, min_y, max_x - min_x, max_y - min_y],
        "output_content": list(out_bbox) if out_bbox else None,
        "scale": scale,
        "placed_at": [x, y],
        "output_size": [FRAME_W, FRAME_H],
    }


def normalize(src: Path, alpha_out: Path, frame_out: Path, report_out: Path) -> None:
    rgba = Image.open(src).convert("RGBA")
    w, h = rgba.size
    px = rgba.load()
    mask = [[False] * w for _ in range(h)]

    for y in range(h):
        for x in range(w):
            r, g, b, a = px[x, y]
            if a == 0 or is_chroma_green(r, g, b):
                px[x, y] = (0, 0, 0, 0)
                continue
            if g - max(r, b) >= 8:
                g = max(r, b)
            px[x, y] = (r, g, b, 255)
            mask[y][x] = True

    mask, component_stats = keep_largest_component(mask)
    for y in range(h):
        for x in range(w):
            if not mask[y][x]:
                px[x, y] = (0, 0, 0, 0)

    mask, hole_stats = fill_internal_holes(rgba, mask)
    for y in range(h):
        for x in range(w):
            if not mask[y][x]:
                px[x, y] = (0, 0, 0, 0)

    alpha_out.parent.mkdir(parents=True, exist_ok=True)
    rgba.save(alpha_out)
    frame, fit_stats = fit_to_frame(rgba, mask)
    frame_out.parent.mkdir(parents=True, exist_ok=True)
    frame.save(frame_out)

    report = {
        "source": src.as_posix(),
        "alpha": alpha_out.as_posix(),
        "output": frame_out.as_posix(),
        "method": "chroma-green largest-component internal-hole-fill fit-96x64",
        "source_size": [w, h],
        "component_stats": component_stats,
        "hole_stats": hole_stats,
        "fit": fit_stats,
    }
    report_out.parent.mkdir(parents=True, exist_ok=True)
    report_out.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--src", required=True)
    parser.add_argument("--alpha-out", required=True)
    parser.add_argument("--frame-out", required=True)
    parser.add_argument("--report", required=True)
    args = parser.parse_args()
    normalize(Path(args.src), Path(args.alpha_out), Path(args.frame_out), Path(args.report))


if __name__ == "__main__":
    main()
