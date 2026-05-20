#!/usr/bin/env python3
"""Rasterize favicon.svg to PNG sizes (Pillow)."""

from pathlib import Path

from PIL import Image, ImageDraw

ROOT = Path(__file__).resolve().parents[1]
OUT_DIR = ROOT / "static" / "brand"


def draw_mark(size: int) -> Image.Image:
    img = Image.new("RGBA", (size, size), (8, 17, 31, 255))
    d = ImageDraw.Draw(img)
    s = size / 32
    cyan = (110, 231, 249, 255)
    blue = (70, 163, 255, 255)
    w = max(1, int(round(2 * s)))
    d.line([(7 * s, 10 * s), (17 * s, 10 * s)], fill=cyan, width=w)
    d.line([(7 * s, 16 * s), (19 * s, 16 * s)], fill=(*cyan[:3], 180), width=w)
    d.line([(7 * s, 22 * s), (14 * s, 22 * s)], fill=(*cyan[:3], 120), width=w)
    d.line([(23 * s, 10 * s), (23 * s, 22 * s)], fill=blue, width=max(1, int(round(2.5 * s))))
    return img


def main() -> None:
    for size, name in ((32, "favicon-32.png"), (180, "apple-touch-icon.png")):
        path = OUT_DIR / name
        draw_mark(size).save(path)
        print(f"Wrote {path}")


if __name__ == "__main__":
    main()
