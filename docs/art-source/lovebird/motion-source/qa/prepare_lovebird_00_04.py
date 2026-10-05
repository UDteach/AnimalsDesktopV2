from __future__ import annotations

import json
from pathlib import Path

from PIL import Image, ImageChops, ImageDraw


RUN = Path(__file__).resolve().parents[1]
RAW = RUN / "raw"
ALPHA = RUN / "alpha"
FRAMES = RUN / "frames"
CONTACT = RUN / "contact"
QA = RUN / "qa"

CANVAS_W = 96
CANVAS_H = 64
TARGET_MAX_W = 58
TARGET_MAX_H = 50
BASELINE_Y = 57
ATTEMPT = 2
COMMON_TARGET_MAX_W = 54
COMMON_TARGET_MAX_H = 47


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
    alpha = img.getchannel("A")
    bbox = alpha.point(lambda v: 255 if v > 12 else 0).getbbox()
    if bbox is None:
        raise ValueError("empty alpha after background removal")
    return bbox


def trim_and_frame(img: Image.Image, bbox: tuple[int, int, int, int]) -> tuple[Image.Image, dict]:
    cropped = img.crop(bbox)
    crop_alpha = cropped.getchannel("A")
    cbbox = crop_alpha.point(lambda v: 255 if v > 12 else 0).getbbox()
    if cbbox is None:
        raise ValueError("empty cropped alpha")

    cw, ch = cropped.size
    scale = min(COMMON_TARGET_MAX_W / cw, COMMON_TARGET_MAX_H / ch)
    out_w = max(1, round(cw * scale))
    out_h = max(1, round(ch * scale))
    resized = cropped.resize((out_w, out_h), Image.Resampling.LANCZOS)

    canvas = Image.new("RGBA", (CANVAS_W, CANVAS_H), (0, 0, 0, 0))
    x = (CANVAS_W - out_w) // 2
    y = BASELINE_Y - out_h
    canvas.alpha_composite(resized, (x, y))

    visible = canvas.getchannel("A").point(lambda v: 255 if v > 0 else 0)
    out_bbox = visible.getbbox()
    alpha_values = list(canvas.getchannel("A").getdata())
    visible_pixels = sum(1 for v in alpha_values if v > 0)
    opaque_pixels = sum(1 for v in alpha_values if v > 180)
    edge_alpha_pixels = 0
    for px in range(CANVAS_W):
        if canvas.getpixel((px, 0))[3] > 0:
            edge_alpha_pixels += 1
        if canvas.getpixel((px, CANVAS_H - 1))[3] > 0:
            edge_alpha_pixels += 1
    for py in range(1, CANVAS_H - 1):
        if canvas.getpixel((0, py))[3] > 0:
            edge_alpha_pixels += 1
        if canvas.getpixel((CANVAS_W - 1, py))[3] > 0:
            edge_alpha_pixels += 1
    transparent_residual_magenta = 0
    transparent_white_gray = 0
    for r, g, b, a in canvas.getdata():
        if a != 0:
            continue
        if r > 200 and b > 200 and g < 90:
            transparent_residual_magenta += 1
        if abs(r - g) <= 4 and abs(g - b) <= 4 and 180 <= r <= 255:
            transparent_white_gray += 1
    return canvas, {
        "source_bbox": bbox,
        "source_size": img.size,
        "canvas_size": canvas.size,
        "visible_bbox": out_bbox,
        "visible_pixels": visible_pixels,
        "opaque_pixels": opaque_pixels,
        "edge_alpha_pixels": edge_alpha_pixels,
        "transparent_residual_magenta_pixels": transparent_residual_magenta,
        "transparent_white_gray_pixels": transparent_white_gray,
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
    if isinstance(bg, tuple):
        base = Image.new("RGBA", frame.size, bg)
    else:
        base = bg.copy()
    base.alpha_composite(frame)
    return base


def write_contact(frames: list[Image.Image]) -> None:
    cell_w, cell_h = 216, 150
    labels_h = 14
    backgrounds = [
        ("light", (246, 247, 242, 255)),
        ("dark", (38, 42, 48, 255)),
        ("checker", None),
    ]
    for name, color in backgrounds:
        sheet = Image.new("RGBA", (cell_w * len(frames), cell_h), (255, 255, 255, 255))
        draw = ImageDraw.Draw(sheet)
        for i, frame in enumerate(frames):
            x0 = i * cell_w
            bg = checker(frame.size) if color is None else color
            rendered = composite(frame, bg)
            scale = 2
            rendered = rendered.resize((frame.width * scale, frame.height * scale), Image.Resampling.NEAREST)
            px = x0 + (cell_w - rendered.width) // 2
            py = labels_h + (cell_h - labels_h - rendered.height) // 2
            sheet.alpha_composite(rendered, (px, py))
            draw.text((x0 + 6, 2), f"frame-{i:02d}", fill=(20, 20, 20, 255) if name != "dark" else (240, 240, 240, 255))
        sheet.save(CONTACT / f"lovebird-peach-faced-00-04-{name}.png")


def main() -> None:
    ALPHA.mkdir(exist_ok=True)
    FRAMES.mkdir(exist_ok=True)
    CONTACT.mkdir(exist_ok=True)
    QA.mkdir(exist_ok=True)

    report = {
        "attempt": ATTEMPT,
        "normalization": {
            "background_key": "broad magenta: r > 210 and b > 180 and g < 90 and (r + b - 2*g) > 300, plus distance falloff",
            "target_max_w": COMMON_TARGET_MAX_W,
            "target_max_h": COMMON_TARGET_MAX_H,
            "scale_policy": "per-frame fit: min(54/source_width, 47/source_height)",
        },
        "frames": [],
        "rejects": [],
    }
    alpha_inputs = []
    bboxes = []
    for idx in range(5):
        raw_path = RAW / f"attempt-{ATTEMPT:02d}-frame-{idx:02d}.png"
        img = Image.open(raw_path)
        alpha = alpha_from_magenta_background(img)
        alpha, component_pixels = keep_largest_component(alpha)
        bbox = visible_bbox(alpha)
        alpha_inputs.append((idx, raw_path, alpha, bbox, component_pixels))
        bboxes.append(bbox)

    max_w = max(b[2] - b[0] for b in bboxes)
    max_h = max(b[3] - b[1] for b in bboxes)
    report["normalization"]["source_max_w"] = max_w
    report["normalization"]["source_max_h"] = max_h

    prepared = []
    for idx, raw_path, alpha, bbox, component_pixels in alpha_inputs:
        alpha_path = ALPHA / f"attempt-{ATTEMPT:02d}-frame-{idx:02d}.png"
        frame_path = FRAMES / f"frame-{idx:02d}.png"
        alpha.save(alpha_path)
        frame, stats = trim_and_frame(alpha, bbox)
        if stats["edge_alpha_pixels"] != 0:
            raise ValueError(f"frame {idx:02d} has non-zero edge alpha: {stats['edge_alpha_pixels']}")
        frame.save(frame_path)
        prepared.append(frame)
        report["frames"].append(
            {
                "frame": idx,
                "raw_path": str(raw_path.relative_to(RUN)),
                "alpha_path": str(alpha_path.relative_to(RUN)),
                "frame_path": str(frame_path.relative_to(RUN)),
                "dimensions_ok": frame.size == (96, 64),
                "non_empty_alpha": stats["visible_pixels"] > 0,
                "stats": stats,
                "largest_component_pixels": component_pixels,
                "visual_note": "candidate reads as compact peach-faced lovebird; parent visual gate required",
            }
        )

    write_contact(prepared)
    (QA / "prepare-report.json").write_text(json.dumps(report, indent=2), encoding="utf-8")


if __name__ == "__main__":
    main()
