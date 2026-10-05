#!/usr/bin/env python3
from __future__ import annotations

import argparse
import csv
import json
import subprocess
import os
from pathlib import Path

from PIL import Image

FIELDS = [
    "frame",
    "slot",
    "orientation",
    "pose_notes",
    "attempt",
    "source_imagegen_path",
    "prompt_path",
    "raw_path",
    "alpha_path",
    "frame_path",
    "prepare_report_path",
    "edge_alpha_pixels",
    "greenish_pixels",
    "nonzero_alpha",
    "alpha_bbox",
    "status",
    "notes",
]


def latest_after(pool: Path, marker: Path) -> Path:
    marker_time = marker.stat().st_mtime
    candidates = [p for p in pool.glob("*.png") if p.stat().st_mtime >= marker_time]
    if not candidates:
        raise SystemExit(f"no generated image found after marker: {marker}")
    return max(candidates, key=lambda p: p.stat().st_mtime)


def inspect(img: Image.Image) -> dict[str, object]:
    alpha = img.getchannel("A")
    edge = 0
    for x in range(96):
        edge += alpha.getpixel((x, 0)) > 8
        edge += alpha.getpixel((x, 63)) > 8
    for y in range(64):
        edge += alpha.getpixel((0, y)) > 8
        edge += alpha.getpixel((95, y)) > 8
    nonzero = 0
    greenish = 0
    for r, g, b, a in img.getdata():
        if a > 0:
            nonzero += 1
        if a > 8 and g > 90 and g > r * 1.18 and g > b * 1.18:
            greenish += 1
    return {
        "size": [96, 64],
        "nonzero_alpha": nonzero,
        "edge_alpha_pixels": edge,
        "greenish_pixels": greenish,
        "alpha_bbox": list(alpha.getbbox() or ()),
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--run-dir", required=True)
    parser.add_argument("--frame", type=int, required=True)
    parser.add_argument("--attempt", default="01")
    parser.add_argument("--slot", required=True)
    parser.add_argument("--pose-notes", required=True)
    parser.add_argument("--marker", required=True)
    parser.add_argument("--source", default="")
    parser.add_argument("--prompt-path", required=True)
    parser.add_argument("--target-w", type=int, default=84)
    parser.add_argument("--target-h", type=int, default=54)
    parser.add_argument("--baseline-y", type=int, default=58)
    parser.add_argument("--notes", default="")
    args = parser.parse_args()

    run = Path(args.run_dir)
    frame_s = f"{args.frame:02d}"
    attempt = f"{int(args.attempt):02d}"
    stem = f"frame-{frame_s}-attempt-{attempt}"
    if args.source:
        source = Path(args.source)
        if not source.exists():
            raise SystemExit(f"explicit generated image does not exist: {source}")
    else:
        pool = Path("/Users/kyota/.codex/generated_images/019f0767-666c-7573-a2a5-036c42e7c274")
        source = latest_after(pool, Path(args.marker))

    raw_path = run / "raw" / f"{stem}-raw.png"
    alpha_path = run / "alpha" / f"{stem}-alpha.png"
    frame_path = run / "frames" / f"frame-{frame_s}.png"
    report_path = run / "qa" / f"{stem}-qa.json"
    if frame_path.exists():
        raise SystemExit(f"refusing to overwrite existing canonical frame: {frame_path}")

    raw_path.parent.mkdir(parents=True, exist_ok=True)
    alpha_path.parent.mkdir(parents=True, exist_ok=True)
    frame_path.parent.mkdir(parents=True, exist_ok=True)
    report_path.parent.mkdir(parents=True, exist_ok=True)

    Image.open(source).convert("RGBA").save(raw_path)
    chroma = Path(os.environ.get("CODEX_HOME", str(Path.home() / ".codex"))) / "skills/.system/imagegen/scripts/remove_chroma_key.py"
    subprocess.run(
        [
            "python3",
            str(chroma),
            "--input",
            str(raw_path),
            "--out",
            str(alpha_path),
            "--auto-key",
            "border",
            "--soft-matte",
            "--transparent-threshold",
            "12",
            "--opaque-threshold",
            "220",
            "--despill",
            "--force",
        ],
        check=True,
        stdout=(run / "qa" / f"{stem}-chroma-key.log").open("w", encoding="utf-8"),
    )

    img = Image.open(alpha_path).convert("RGBA")
    bbox = img.getchannel("A").getbbox()
    if bbox is None:
        raise SystemExit(f"empty alpha: {alpha_path}")
    crop = img.crop(bbox)
    scale = min(args.target_w / crop.width, args.target_h / crop.height, 1.0)
    new_w = max(1, round(crop.width * scale))
    new_h = max(1, round(crop.height * scale))
    resized = crop.resize((new_w, new_h), Image.Resampling.LANCZOS)
    out = Image.new("RGBA", (96, 64), (0, 0, 0, 0))
    x = (96 - new_w) // 2
    y = max(0, args.baseline_y - new_h)
    out.alpha_composite(resized, (x, y))
    out.save(frame_path)

    qa = inspect(out)
    report = {
        "variant": "budgerigar_green_yellow",
        "frame": frame_s,
        "attempt": attempt,
        "status": "candidate",
        "source_imagegen_path": source.as_posix(),
        "raw_path": raw_path.as_posix(),
        "alpha_path": alpha_path.as_posix(),
        "frame_path": frame_path.as_posix(),
        "prompt_path": Path(args.prompt_path).as_posix(),
        "method": "magenta-key remove_chroma_key source-alpha no-hole-fill fit-96x64",
        "source_content": [bbox[0], bbox[1], bbox[2] - bbox[0], bbox[3] - bbox[1]],
        "fit": {"scale": scale, "placed_at": [x, y], "output_size": [96, 64], "output_content": qa["alpha_bbox"]},
        "mechanical_qa": qa,
    }
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    row = {
        "frame": frame_s,
        "slot": args.slot,
        "orientation": "right-facing side",
        "pose_notes": args.pose_notes,
        "attempt": attempt,
        "source_imagegen_path": source.as_posix(),
        "prompt_path": Path(args.prompt_path).as_posix(),
        "raw_path": raw_path.as_posix(),
        "alpha_path": alpha_path.as_posix(),
        "frame_path": frame_path.as_posix(),
        "prepare_report_path": report_path.as_posix(),
        "edge_alpha_pixels": qa["edge_alpha_pixels"],
        "greenish_pixels": qa["greenish_pixels"],
        "nonzero_alpha": qa["nonzero_alpha"],
        "alpha_bbox": json.dumps(qa["alpha_bbox"]),
        "status": "candidate",
        "notes": args.notes,
    }
    with (run / "manifest.csv").open("a", newline="", encoding="utf-8") as f:
        csv.DictWriter(f, fieldnames=FIELDS).writerow(row)
    print(json.dumps(row, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
