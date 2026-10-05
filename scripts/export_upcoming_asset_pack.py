#!/usr/bin/env python3
"""Export upcoming-animal planning assets from the ImageGen source sheet."""

from __future__ import annotations

import json
import shutil
import sys
from collections import deque
from pathlib import Path

from PIL import Image, ImageFilter

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "scripts"))

import build_page_assets as page_assets  # noqa: E402


PACK_ID = "upcoming-asset-pack-20260627"
PACK_DIR = ROOT / "assets" / "source" / "upcoming" / "20260627"
COLOR_DIR = PACK_DIR / "color"
SILHOUETTE_DIR = PACK_DIR / "silhouette"
SOURCE_DIR = PACK_DIR / "source"
CONTACT_DIR = PACK_DIR / "contact-sheets"
SOURCE_COPY = SOURCE_DIR / page_assets.UPCOMING_SOURCE.name
SOURCE_NOTE = "Planning asset pack only. These are not accepted runtime motion sources."
CARD_W = page_assets.UPCOMING_W
CARD_H = page_assets.UPCOMING_H

DISPLAY_NAMES = {
    "chipmunk": ("シマリス", "Striped chipmunk"),
    "leucistic_sugar_glider": ("リューシスティックモモンガ", "Leucistic sugar glider"),
    "african_dormouse": ("アフリカヤマネ", "African pygmy dormouse"),
    "netherland_dwarf_himalayan": (
        "ネザーランドドワーフ（ヒマラヤン）",
        "Netherland Dwarf rabbit, Himalayan",
    ),
    "american_flying_squirrel": ("アメリカモモンガ", "American flying squirrel"),
    "longhair_hamster_black_white": (
        "白黒長毛ハムスター",
        "Black-and-white long-haired Syrian hamster",
    ),
    "djungarian_hamster_yellow": ("イエロージャンガリアン", "Yellow Djungarian hamster"),
    "djungarian_hamster_pearl_white": (
        "パールホワイトジャンガリアン",
        "Pearl white Djungarian hamster",
    ),
    "fancy_rat_blue_hooded": ("ファンシーラット（ブルーフーディッド）", "Fancy rat, blue hooded"),
    "fancy_rat_chocolate_self": (
        "ファンシーラット（チョコレートセルフ）",
        "Fancy rat, chocolate self",
    ),
    "fancy_rat_cream_agouti": (
        "ファンシーラット（クリームアグーチ）",
        "Fancy rat, cream agouti self",
    ),
    "rabbit_gray": ("グレーうさぎ", "Gray rabbit"),
    "whites_tree_frog": ("水色イエアメガエル", "Blue White's tree frog"),
    "leopard_gecko": ("ヒョウモントカゲモドキ", "Leopard gecko"),
    "budgerigar": ("セキセイインコ", "Green-and-yellow budgerigar"),
    "cockatiel": ("オカメインコ", "Normal gray cockatiel"),
    "java_sparrow": ("文鳥", "Java sparrow"),
    "african_fat_tailed_gecko": (
        "ニシアフリカトカゲモドキ",
        "African fat-tailed gecko",
    ),
}


def reset_pngs(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True)
    for png in path.glob("*.png"):
        png.unlink()


def crop_rgba_with_mask(cell: Image.Image, mask: Image.Image, pad: int = 10) -> tuple[Image.Image, tuple[int, int, int, int]]:
    bbox = mask.getbbox()
    if bbox is None:
        raise SystemExit("empty upcoming source cell")
    left, top, right, bottom = bbox
    left = max(0, left - pad)
    top = max(0, top - pad)
    right = min(cell.width, right + pad)
    bottom = min(cell.height, bottom + pad)
    crop_box = (left, top, right, bottom)
    animal = cell.convert("RGBA").crop(crop_box)
    alpha = mask.crop(crop_box)
    animal.putalpha(alpha)
    return animal, crop_box


def full_foreground_mask(source: Image.Image) -> Image.Image:
    rgba = source.convert("RGBA")
    w, h = rgba.size
    pixels = rgba.load()
    visited = bytearray(w * h)
    queue: deque[tuple[int, int]] = deque()

    def offset(x: int, y: int) -> int:
        return y * w + x

    def enqueue(x: int, y: int) -> None:
        i = offset(x, y)
        if visited[i] or not page_assets._background_like(pixels[x, y]):
            return
        visited[i] = 1
        queue.append((x, y))

    for x in range(w):
        enqueue(x, 0)
        enqueue(x, h - 1)
    for y in range(h):
        enqueue(0, y)
        enqueue(w - 1, y)

    while queue:
        x, y = queue.popleft()
        for nx, ny in ((x - 1, y), (x + 1, y), (x, y - 1), (x, y + 1)):
            if 0 <= nx < w and 0 <= ny < h:
                enqueue(nx, ny)

    mask = Image.new("L", (w, h), 0)
    mask_pixels = mask.load()
    for y in range(h):
        for x in range(w):
            _r, _g, _b, a = pixels[x, y]
            if a > 8 and not visited[offset(x, y)]:
                mask_pixels[x, y] = 255
    return mask.filter(ImageFilter.MaxFilter(3)).filter(ImageFilter.MinFilter(3))


def ordered_component_masks(mask: Image.Image, expected: int) -> list[Image.Image]:
    mask = mask.convert("L")
    w, h = mask.size
    pixels = mask.load()
    visited = bytearray(w * h)
    components: list[tuple[int, tuple[int, int, int, int], list[tuple[int, int]]]] = []

    def offset(x: int, y: int) -> int:
        return y * w + x

    for y in range(h):
        for x in range(w):
            start = offset(x, y)
            if visited[start] or pixels[x, y] == 0:
                continue
            visited[start] = 1
            queue: deque[tuple[int, int]] = deque([(x, y)])
            points: list[tuple[int, int]] = []
            while queue:
                cx, cy = queue.popleft()
                points.append((cx, cy))
                for nx, ny in ((cx - 1, cy), (cx + 1, cy), (cx, cy - 1), (cx, cy + 1)):
                    if not (0 <= nx < w and 0 <= ny < h):
                        continue
                    i = offset(nx, ny)
                    if visited[i] or pixels[nx, ny] == 0:
                        continue
                    visited[i] = 1
                    queue.append((nx, ny))
            if len(points) <= 500:
                continue
            xs = [point[0] for point in points]
            ys = [point[1] for point in points]
            components.append((len(points), (min(xs), min(ys), max(xs) + 1, max(ys) + 1), points))

    if len(components) != expected:
        raise SystemExit(f"expected {expected} source components, found {len(components)}")

    # Sort into visual rows first, then left-to-right within each row. This keeps
    # animals intact even when the ImageGen sheet is not a strict grid.
    components.sort(key=lambda item: (item[1][1] + item[1][3]) / 2)
    rows = [components[i : i + page_assets.UPCOMING_SOURCE_COLS] for i in range(0, expected, page_assets.UPCOMING_SOURCE_COLS)]
    ordered: list[tuple[int, tuple[int, int, int, int], list[tuple[int, int]]]] = []
    for row in rows:
        ordered.extend(sorted(row, key=lambda item: item[1][0]))

    out: list[Image.Image] = []
    for _area, _bbox, points in ordered:
        component_mask = Image.new("L", mask.size, 0)
        out_pixels = component_mask.load()
        for x, y in points:
            out_pixels[x, y] = pixels[x, y]
        out.append(component_mask)
    return out


def silhouette_from_mask(mask: Image.Image, kind: str) -> Image.Image:
    bbox = mask.getbbox()
    if bbox is None:
        raise SystemExit(f"empty ImageGen source component for upcoming silhouette: {kind}")
    cropped_mask = mask.crop(bbox)
    silhouette = Image.new("RGBA", cropped_mask.size, page_assets.SILHOUETTE_INK)
    silhouette.putalpha(cropped_mask)
    return fit_on_canvas(silhouette)


def fit_on_canvas(img: Image.Image, max_w: int = CARD_W - 28, max_h: int = CARD_H - 24) -> Image.Image:
    scale = min(max_w / img.width, max_h / img.height)
    resized = img.resize(
        (max(1, round(img.width * scale)), max(1, round(img.height * scale))),
        Image.Resampling.LANCZOS,
    )
    canvas = Image.new("RGBA", (CARD_W, CARD_H), (0, 0, 0, 0))
    canvas.alpha_composite(resized, ((CARD_W - resized.width) // 2, (CARD_H - resized.height) // 2))
    return canvas


def write_contact_sheet(images: list[Image.Image], output: Path, *, columns: int = 6, background: tuple[int, int, int] = (248, 248, 243)) -> None:
    rows = (len(images) + columns - 1) // columns
    sheet = Image.new("RGBA", (CARD_W * columns, CARD_H * rows), (*background, 255))
    for i, img in enumerate(images):
        sheet.alpha_composite(img, ((i % columns) * CARD_W, (i // columns) * CARD_H))
    sheet.convert("RGB").save(output)


def export_pack() -> None:
    if not page_assets.UPCOMING_SOURCE.exists():
        raise SystemExit(f"missing ImageGen source sheet: {page_assets.UPCOMING_SOURCE}")

    for path in (COLOR_DIR, SILHOUETTE_DIR, CONTACT_DIR):
        reset_pngs(path)
    SOURCE_DIR.mkdir(parents=True, exist_ok=True)
    shutil.copy2(page_assets.UPCOMING_SOURCE, SOURCE_COPY)

    source = Image.open(page_assets.UPCOMING_SOURCE).convert("RGBA")
    color_cards: list[Image.Image] = []
    silhouette_cards: list[Image.Image] = []
    items: list[dict[str, object]] = []

    component_masks = ordered_component_masks(full_foreground_mask(source), len(page_assets.UPCOMING_SOURCE_LAYOUT))

    for index, kind in enumerate(page_assets.UPCOMING_SOURCE_LAYOUT):
        if kind not in DISPLAY_NAMES:
            raise SystemExit(f"missing display name for upcoming asset: {kind}")
        mask = component_masks[index]
        color, bbox = crop_rgba_with_mask(source, mask)
        color_path = COLOR_DIR / f"{kind}.png"
        color.save(color_path)

        silhouette = silhouette_from_mask(mask, kind)
        silhouette_path = SILHOUETTE_DIR / f"{kind}.png"
        silhouette.save(silhouette_path)

        color_cards.append(fit_on_canvas(color))
        silhouette_cards.append(silhouette)
        name_ja, name_en = DISPLAY_NAMES[kind]
        items.append(
            {
                "id": kind,
                "order": index + 1,
                "name_ja": name_ja,
                "name_en": name_en,
                "source_cell": index,
                "source_component": index,
                "source_extraction": "connected-components",
                "source_bbox": list(bbox),
                "color": str(color_path.relative_to(PACK_DIR)),
                "color_size": list(color.size),
                "silhouette": str(silhouette_path.relative_to(PACK_DIR)),
                "silhouette_size": [CARD_W, CARD_H],
            }
        )

    write_contact_sheet(color_cards, CONTACT_DIR / "color-contact-sheet.png")
    write_contact_sheet(silhouette_cards, CONTACT_DIR / "silhouette-contact-sheet.png")

    manifest = {
        "pack_id": PACK_ID,
        "note": SOURCE_NOTE,
        "source_sheet": str(SOURCE_COPY.relative_to(PACK_DIR)),
        "source_sheet_original": str(page_assets.UPCOMING_SOURCE.relative_to(ROOT)),
        "source_grid": {
            "columns": page_assets.UPCOMING_SOURCE_COLS,
            "rows": page_assets.UPCOMING_SOURCE_ROWS,
        },
        "output_dirs": {
            "color": "color",
            "silhouette": "silhouette",
            "contact_sheets": "contact-sheets",
        },
        "items": items,
    }
    (PACK_DIR / "manifest.json").write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n")
    (PACK_DIR / "README.md").write_text(readme_text(items), encoding="utf-8")
    print(f"exported {len(items)} upcoming planning assets to {PACK_DIR.relative_to(ROOT)}")


def readme_text(items: list[dict[str, object]]) -> str:
    lines = [
        "# Upcoming Asset Pack - 2026-06-27",
        "",
        SOURCE_NOTE,
        "",
        "This pack exports the current upcoming-animal ImageGen sheet as reusable",
        "planning assets outside the GitHub Pages directory.",
        "",
        "## Contents",
        "",
        "- `source/`: copied ImageGen source sheet.",
        "- `color/`: transparent color cutouts, one animal per PNG.",
        "- `silhouette/`: normalized 256x160 black silhouettes, one animal per PNG.",
        "- `contact-sheets/`: color and silhouette review sheets.",
        "- `manifest.json`: stable IDs, JP/EN names, source cells, and output paths.",
        "",
        "## Animals",
        "",
        "| Order | ID | JP | EN |",
        "| --- | --- | --- | --- |",
    ]
    for item in items:
        lines.append(f"| {item['order']} | `{item['id']}` | {item['name_ja']} | {item['name_en']} |")
    lines.extend(
        [
            "",
            "## Rebuild",
            "",
            "```sh",
            "python3 scripts/export_upcoming_asset_pack.py",
            "```",
            "",
        ]
    )
    return "\n".join(lines)


def main() -> None:
    export_pack()


if __name__ == "__main__":
    main()
