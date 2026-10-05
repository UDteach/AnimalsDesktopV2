from __future__ import annotations

import json
from pathlib import Path

from PIL import Image, ImageDraw


RUN = Path(__file__).resolve().parents[1]
RAW = RUN / "raw"
ALPHA = RUN / "alpha"
FRAMES = RUN / "frames"
CONTACT = RUN / "contact"
QA = RUN / "qa"

CANVAS_W = 96
CANVAS_H = 64
TARGET_MAX_W = 54
TARGET_MAX_H = 47
BASELINE_Y = 57
ATTEMPT = 4
FRAME_RANGE = range(5, 13)


def alpha_from_magenta_background(img: Image.Image) -> Image.Image:
    rgba = img.convert("RGBA")
    pixels = rgba.load()
    w, h = rgba.size
    mask = Image.new("L", rgba.size, 0)
    mask_px = mask.load()
    for y in range(h):
        for x in range(w):
            r, g, b, _ = pixels[x, y]
            broad_magenta_bg = r > 210 and b > 180 and g < 90 and (r + b - 2 * g) > 300
            magenta_dist = ((255 - r) ** 2 + g**2 + (255 - b) ** 2) ** 0.5
            if broad_magenta_bg or magenta_dist <= 70:
                a = 0
            elif magenta_dist >= 125:
                a = 255
            else:
                a = int((magenta_dist - 70) / 55 * 255)
            mask_px[x, y] = a
    rgba.putalpha(mask)
    return rgba


def keep_largest_component(img: Image.Image) -> tuple[Image.Image, int]:
    alpha = img.getchannel("A")
    w, h = img.size
    data = alpha.load()
    seen = bytearray(w * h)
    best: list[tuple[int, int]] = []
    for y in range(h):
        for x in range(w):
            idx = y * w + x
            if seen[idx] or data[x, y] <= 12:
                continue
            stack = [(x, y)]
            seen[idx] = 1
            comp: list[tuple[int, int]] = []
            while stack:
                cx, cy = stack.pop()
                comp.append((cx, cy))
                for nx, ny in ((cx + 1, cy), (cx - 1, cy), (cx, cy + 1), (cx, cy - 1)):
                    if nx < 0 or ny < 0 or nx >= w or ny >= h:
                        continue
                    nidx = ny * w + nx
                    if seen[nidx] or data[nx, ny] <= 12:
                        continue
                    seen[nidx] = 1
                    stack.append((nx, ny))
            if len(comp) > len(best):
                best = comp

    out = Image.new("RGBA", img.size, (0, 0, 0, 0))
    src = img.load()
    dst = out.load()
    for x, y in best:
        dst[x, y] = src[x, y]
    return out, len(best)


def visible_bbox(img: Image.Image) -> tuple[int, int, int, int]:
    bbox = img.getchannel("A").point(lambda v: 255 if v > 12 else 0).getbbox()
    if bbox is None:
        raise ValueError("empty alpha after background removal")
    return bbox


def trim_and_frame(img: Image.Image, bbox: tuple[int, int, int, int]) -> tuple[Image.Image, dict]:
    cropped = img.crop(bbox)
    cw, ch = cropped.size
    scale = min(TARGET_MAX_W / cw, TARGET_MAX_H / ch)
    out_w = max(1, round(cw * scale))
    out_h = max(1, round(ch * scale))
    resized = cropped.resize((out_w, out_h), Image.Resampling.LANCZOS)

    canvas = Image.new("RGBA", (CANVAS_W, CANVAS_H), (0, 0, 0, 0))
    x = (CANVAS_W - out_w) // 2
    y = BASELINE_Y - out_h
    canvas.alpha_composite(resized, (x, y))

    alpha_values = list(canvas.getchannel("A").getdata())
    edge_alpha_pixels = 0
    for px in range(CANVAS_W):
        edge_alpha_pixels += int(canvas.getpixel((px, 0))[3] > 0)
        edge_alpha_pixels += int(canvas.getpixel((px, CANVAS_H - 1))[3] > 0)
    for py in range(1, CANVAS_H - 1):
        edge_alpha_pixels += int(canvas.getpixel((0, py))[3] > 0)
        edge_alpha_pixels += int(canvas.getpixel((CANVAS_W - 1, py))[3] > 0)

    residual_magenta = 0
    residual_white_gray = 0
    for r, g, b, a in canvas.getdata():
        if a != 0:
            continue
        if r > 200 and b > 200 and g < 90:
            residual_magenta += 1
        if abs(r - g) <= 4 and abs(g - b) <= 4 and 180 <= r <= 255:
            residual_white_gray += 1

    return canvas, {
        "source_bbox": bbox,
        "source_size": img.size,
        "canvas_size": canvas.size,
        "visible_bbox": canvas.getchannel("A").point(lambda v: 255 if v > 0 else 0).getbbox(),
        "visible_pixels": sum(1 for v in alpha_values if v > 0),
        "opaque_pixels": sum(1 for v in alpha_values if v > 180),
        "edge_alpha_pixels": edge_alpha_pixels,
        "transparent_residual_magenta_pixels": residual_magenta,
        "transparent_white_gray_pixels": residual_white_gray,
        "scale": scale,
        "placed_xy": [x, y],
        "baseline_y": BASELINE_Y,
    }


def checker(size: tuple[int, int], cell: int = 8) -> Image.Image:
    w, h = size
    img = Image.new("RGBA", size, (238, 238, 238, 255))
    draw = ImageDraw.Draw(img)
    for y in range(0, h, cell):
        for x in range(0, w, cell):
            if (x // cell + y // cell) % 2:
                draw.rectangle([x, y, x + cell - 1, y + cell - 1], fill=(204, 204, 204, 255))
    return img


def composite(frame: Image.Image, bg: Image.Image | tuple[int, int, int, int]) -> Image.Image:
    base = bg.copy() if isinstance(bg, Image.Image) else Image.new("RGBA", frame.size, bg)
    base.alpha_composite(frame)
    return base


def write_contact(frame_numbers: list[int], out_prefix: str) -> None:
    cell_w, cell_h = 216, 150
    labels_h = 14
    backgrounds = [
        ("light", (246, 247, 242, 255)),
        ("dark", (38, 42, 48, 255)),
        ("checker", None),
    ]
    frames = [(n, Image.open(FRAMES / f"frame-{n:02d}.png").convert("RGBA")) for n in frame_numbers]
    for name, color in backgrounds:
        sheet = Image.new("RGBA", (cell_w * len(frames), cell_h), (255, 255, 255, 255))
        draw = ImageDraw.Draw(sheet)
        for i, (frame_no, frame) in enumerate(frames):
            x0 = i * cell_w
            rendered = composite(frame, checker(frame.size) if color is None else color)
            rendered = rendered.resize((frame.width * 2, frame.height * 2), Image.Resampling.NEAREST)
            px = x0 + (cell_w - rendered.width) // 2
            py = labels_h + (cell_h - labels_h - rendered.height) // 2
            sheet.alpha_composite(rendered, (px, py))
            fill = (20, 20, 20, 255) if name != "dark" else (240, 240, 240, 255)
            draw.text((x0 + 6, 2), f"frame-{frame_no:02d}", fill=fill)
        sheet.save(CONTACT / f"{out_prefix}-{name}.png")


def main() -> None:
    ALPHA.mkdir(exist_ok=True)
    FRAMES.mkdir(exist_ok=True)
    CONTACT.mkdir(exist_ok=True)
    QA.mkdir(exist_ok=True)

    for idx in range(5):
        if not (FRAMES / f"frame-{idx:02d}.png").exists():
            raise FileNotFoundError(f"approved frame-{idx:02d}.png is missing")

    report = {
        "attempt": ATTEMPT,
        "frame_range": [5, 12],
        "normalization": {
            "background_key": "broad magenta: r > 210 and b > 180 and g < 90 and (r + b - 2*g) > 300, plus distance falloff",
            "target_max_w": TARGET_MAX_W,
            "target_max_h": TARGET_MAX_H,
            "scale_policy": "per-frame fit",
            "baseline_y": BASELINE_Y,
        },
        "frames": [],
        "rejects": [],
    }

    for idx in FRAME_RANGE:
        raw_path = RAW / f"attempt-{ATTEMPT:02d}-frame-{idx:02d}.png"
        alpha_path = ALPHA / f"attempt-{ATTEMPT:02d}-frame-{idx:02d}.png"
        frame_path = FRAMES / f"frame-{idx:02d}.png"
        img = Image.open(raw_path)
        alpha = alpha_from_magenta_background(img)
        alpha, component_pixels = keep_largest_component(alpha)
        bbox = visible_bbox(alpha)
        alpha.save(alpha_path)
        frame, stats = trim_and_frame(alpha, bbox)
        if frame.size != (96, 64):
            raise ValueError(f"frame-{idx:02d} has wrong size {frame.size}")
        if stats["visible_pixels"] <= 0:
            raise ValueError(f"frame-{idx:02d} has empty alpha")
        if stats["edge_alpha_pixels"] != 0:
            raise ValueError(f"frame-{idx:02d} has non-zero edge alpha")
        if stats["transparent_residual_magenta_pixels"] or stats["transparent_white_gray_pixels"]:
            raise ValueError(f"frame-{idx:02d} has residual background pixels")
        frame.save(frame_path)
        report["frames"].append(
            {
                "frame": idx,
                "raw_path": str(raw_path.relative_to(RUN)),
                "alpha_path": str(alpha_path.relative_to(RUN)),
                "frame_path": str(frame_path.relative_to(RUN)),
                "dimensions_ok": True,
                "non_empty_alpha": True,
                "largest_component_pixels": component_pixels,
                "stats": stats,
            }
        )

    write_contact(list(range(0, 13)), "lovebird-peach-faced-00-12")
    write_contact(list(range(5, 13)), "lovebird-peach-faced-05-12")
    (QA / "prepare-report-05-12.json").write_text(json.dumps(report, indent=2), encoding="utf-8")


if __name__ == "__main__":
    main()
