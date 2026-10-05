//go:build windows || darwin

package main

import "image"

type spriteRenderAdjustment struct {
	scalePermille  int
	baselineOffset int
}

func renderAdjustmentForFrame(variantID string, frame int) spriteRenderAdjustment {
	switch variantID {
	case "gecko_leopard":
		switch frame {
		case 32:
			return spriteRenderAdjustment{scalePermille: 900, baselineOffset: 1}
		case 33:
			return spriteRenderAdjustment{scalePermille: 720, baselineOffset: 2}
		case 34:
			return spriteRenderAdjustment{scalePermille: 760, baselineOffset: 2}
		case 35:
			return spriteRenderAdjustment{scalePermille: 740, baselineOffset: 2}
		}
	}
	return spriteRenderAdjustment{scalePermille: 1000}
}

func (adjustment spriteRenderAdjustment) adjustRect(rect image.Rectangle) image.Rectangle {
	scale := adjustment.scalePermille
	if scale <= 0 {
		scale = 1000
	}
	if scale == 1000 && adjustment.baselineOffset == 0 {
		return rect
	}
	w := max(1, rect.Dx()*scale/1000)
	h := max(1, rect.Dy()*scale/1000)
	x := rect.Min.X + (rect.Dx()-w)/2
	y := rect.Max.Y - h + adjustment.baselineOffset
	return image.Rect(x, y, x+w, y+h)
}
