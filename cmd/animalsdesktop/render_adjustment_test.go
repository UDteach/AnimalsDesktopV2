//go:build windows || darwin

package main

import (
	"image"
	"testing"
)

func TestRenderAdjustmentKeepsDefaultRect(t *testing.T) {
	rect := image.Rect(10, 20, 106, 84)
	got := renderAdjustmentForFrame("chinchilla_standard_gray", 33).adjustRect(rect)
	if got != rect {
		t.Fatalf("default adjustment rect = %v, want %v", got, rect)
	}
}

func TestLeopardGeckoTurnFrameAdjustmentShrinksAroundBaseline(t *testing.T) {
	rect := image.Rect(10, 20, 106, 84)
	tests := []struct {
		frame int
		want  image.Rectangle
	}{
		{frame: 32, want: image.Rect(15, 28, 101, 85)},
		{frame: 33, want: image.Rect(23, 40, 92, 86)},
		{frame: 34, want: image.Rect(22, 38, 94, 86)},
		{frame: 35, want: image.Rect(22, 39, 93, 86)},
	}
	for _, tt := range tests {
		got := renderAdjustmentForFrame("gecko_leopard", tt.frame).adjustRect(rect)
		if got != tt.want {
			t.Fatalf("gecko frame %02d adjusted rect = %v, want %v", tt.frame, got, tt.want)
		}
	}
	if got := renderAdjustmentForFrame("gecko_leopard", 36).adjustRect(rect); got != rect {
		t.Fatalf("gecko frame 36 adjusted rect = %v, want no adjustment %v", got, rect)
	}
}
