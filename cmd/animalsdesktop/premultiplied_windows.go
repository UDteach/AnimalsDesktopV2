//go:build windows

package main

import "image"

// image.RGBA already stores premultiplied channels, exactly as Win32 expects.
// Multiplying by alpha here a second time makes fur/whisker edges dark.
func copyPremultipliedBGRA(dst []byte, img *image.RGBA) {
	i := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.RGBAAt(x, y)
			dst[i] = c.B
			dst[i+1] = c.G
			dst[i+2] = c.R
			dst[i+3] = c.A
			i += 4
		}
	}
}
