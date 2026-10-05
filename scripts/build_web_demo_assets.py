#!/usr/bin/env python3
"""Build compact page-only sprite strips for the animated Top demo.

The runtime sheets stay authoritative and unchanged. This script copies only
the idle, walk, stand, and groom frames needed by the public web demo.
"""

from pathlib import Path

from PIL import Image


ROOT = Path(__file__).resolve().parents[1]
FRAME_W = 96
FRAME_H = 64
SOURCE_FRAMES = (0, 1, 3, 4, 5, 7, 48, 49, 50, 51, 52, 53, 54, 55)
WEB_VARIANTS = (
    "chinchilla_standard_gray",
    "hamster_golden_syrian",
    "macaroni_mouse_tan",
    "rabbit_chestnut_agouti",
    "guinea_pig_tricolor",
    "sugar_glider_gray",
    "fancy_rat_hooded",
    "holland_lop_broken_orange",
    "chipmunk_striped",
    "gecko_leopard",
)


def build_strip(variant: str) -> None:
    source_path = ROOT / "assets" / "sprites" / f"{variant}_set00.png"
    if not source_path.exists():
        raise SystemExit(f"missing runtime sprite sheet: {source_path}")

    source = Image.open(source_path).convert("RGBA")
    expected_size = (FRAME_W * 62, FRAME_H)
    if source.size != expected_size:
        raise SystemExit(f"unexpected runtime sheet size for {variant}: {source.size}, want {expected_size}")

    output = Image.new("RGBA", (FRAME_W * len(SOURCE_FRAMES), FRAME_H), (0, 0, 0, 0))
    for column, source_frame in enumerate(SOURCE_FRAMES):
        left = source_frame * FRAME_W
        frame = source.crop((left, 0, left + FRAME_W, FRAME_H))
        output.alpha_composite(frame, (column * FRAME_W, 0))

    output_dir = ROOT / "docs" / "assets" / "web-demo"
    output_dir.mkdir(parents=True, exist_ok=True)
    output.save(output_dir / f"{variant}.png", optimize=True)


def main() -> None:
    for variant in WEB_VARIANTS:
        build_strip(variant)
    print(f"wrote {len(WEB_VARIANTS)} compact web demo sprite strips")


if __name__ == "__main__":
    main()
