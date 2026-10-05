#!/usr/bin/env python3
"""Fit a transparent generated animal cutout into a 96x64 review frame.

This is intentionally for review candidates only. It does not replace
cmd/prepareframe and must not be used as an accepted-frame promotion gate.
"""

from __future__ import annotations

import argparse
import json
from pathlib import Path

from PIL import Image, ImageDraw


def alpha_bbox(img: Image.Image) -> tuple[int, int, int, int] | None:
    alpha = img.getchannel("A")
    return alpha.getbbox()


def fit_frame(src: Path, out: Path, report: Path, contact: Path) -> None:
    img = Image.open(src).convert("RGBA")
    bbox = alpha_bbox(img)
    if bbox is None:
        raise SystemExit(f"{src} has no visible alpha")

    crop = img.crop(bbox)
    max_w = 92
    max_h = 60
    scale = min(max_w / crop.width, max_h / crop.height)
    new_w = max(1, round(crop.width * scale))
    new_h = max(1, round(crop.height * scale))
    resized = crop.resize((new_w, new_h), Image.Resampling.LANCZOS)

    frame = Image.new("RGBA", (96, 64), (0, 0, 0, 0))
    x = (96 - new_w) // 2
    y = 64 - new_h - 2
    if y < 1:
        y = 1
    frame.alpha_composite(resized, (x, y))

    out.parent.mkdir(parents=True, exist_ok=True)
    frame.save(out)

    content = alpha_bbox(frame)
    if content is None:
        raise SystemExit(f"{out} has no visible alpha after fit")

    make_contact(frame, contact)
    report.parent.mkdir(parents=True, exist_ok=True)
    report.write_text(
        json.dumps(
            {
                "source": src.as_posix(),
                "output": out.as_posix(),
                "contact": contact.as_posix(),
                "source_size": list(img.size),
                "source_alpha_bbox": list(bbox),
                "frame_size": [96, 64],
                "frame_alpha_bbox": list(content),
                "scale": scale,
                "status": "review_candidate_only",
                "promotion_note": "Do not promote without cmd/prepareframe/auditframes and visual QA.",
            },
            indent=2,
        )
        + "\n",
        encoding="utf-8",
    )


def make_contact(frame: Image.Image, out: Path) -> None:
    panel_w, panel_h = 128, 96
    sheet = Image.new("RGBA", (panel_w * 3, panel_h), (255, 255, 255, 255))
    backgrounds = [
        ("light", (245, 245, 245, 255)),
        ("dark", (36, 36, 36, 255)),
        ("checker", None),
    ]
    for idx, (_, bg) in enumerate(backgrounds):
        x0 = idx * panel_w
        if bg is None:
            draw_checker(sheet, x0, 0, panel_w, panel_h)
        else:
            ImageDraw.Draw(sheet).rectangle([x0, 0, x0 + panel_w, panel_h], fill=bg)
        sheet.alpha_composite(frame.resize((96, 64), Image.Resampling.NEAREST), (x0 + 16, 16))
    out.parent.mkdir(parents=True, exist_ok=True)
    sheet.convert("RGB").save(out)


def draw_checker(img: Image.Image, x: int, y: int, w: int, h: int) -> None:
    draw = ImageDraw.Draw(img)
    colors = [(235, 235, 235, 255), (205, 205, 205, 255)]
    step = 8
    for yy in range(y, y + h, step):
        for xx in range(x, x + w, step):
            color = colors[((xx - x) // step + (yy - y) // step) % 2]
            draw.rectangle([xx, yy, xx + step - 1, yy + step - 1], fill=color)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--src", required=True)
    parser.add_argument("--out", required=True)
    parser.add_argument("--report", required=True)
    parser.add_argument("--contact", required=True)
    args = parser.parse_args()
    fit_frame(Path(args.src), Path(args.out), Path(args.report), Path(args.contact))


if __name__ == "__main__":
    main()
