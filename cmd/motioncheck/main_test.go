package main

import (
	"image"
	"image/color"
	"testing"
)

func TestBoundsExcludesTransparentCanvasAndUsesThreshold(t *testing.T) {
	im := image.NewRGBA(image.Rect(0, 0, 144, 96))
	im.SetRGBA(15, 22, color.RGBA{50, 50, 50, 100})
	im.SetRGBA(50, 80, color.RGBA{200, 200, 200, 200})
	im.SetRGBA(51, 81, color.RGBA{255, 255, 255, 255})
	r, n := bounds(im, 180)
	if r != image.Rect(50, 80, 52, 82) || n != 2 {
		t.Fatalf("solid %v / %d", r, n)
	}
	r, n = bounds(im, 1)
	if r != image.Rect(15, 22, 52, 82) || n != 3 {
		t.Fatalf("soft %v / %d", r, n)
	}
	r, n = bounds(image.NewRGBA(im.Bounds()), 180)
	if !r.Empty() || n != 0 {
		t.Fatalf("empty %v / %d", r, n)
	}
}
