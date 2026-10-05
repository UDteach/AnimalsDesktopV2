package main

import (
	"animals-desktop/internal/motion"
	"image"
	"image/color"
	"image/draw"
	"math"
	"testing"
)

func TestStationaryPreviewNeverStartsWalkingAndPreservesDefaultSchedule(t *testing.T) {
	v := &motion.Loaded{Variant: motion.Variant{
		Profile: motion.Profile{Stride: .5},
		Motions: map[string]motion.Clip{"walk": {Durations: []float64{200, 200}}},
	}}
	for _, sample := range []struct{ seconds, want float64 }{{0, 120}, {3.99, 120}, {4, 0}, {5.99, 0}, {6, 162}, {6.59, 162}, {6.6, 120}, {12, 120}} {
		if got := previewTarget(v, 96, sample.seconds, false); math.Abs(got-sample.want) > 1e-9 {
			t.Fatalf("default schedule at%g: got%g want%g", sample.seconds, got, sample.want)
		}
		if got := previewTarget(v, 96, sample.seconds, true); got != 0 {
			t.Fatalf("stationary review moves at%g: %g", sample.seconds, got)
		}
	}
}

func TestPreviewRowsPreserveFullDrawAtTailAndRaisedPoseExtremes(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 24, 16))
	draw.Draw(src, src.Bounds(), image.NewUniform(color.RGBA{90, 60, 30, 255}), image.Point{}, draw.Src)
	for _, height := range []int{64, 96, 128} {
		v := &motion.Loaded{Variant: motion.Variant{
			Profile: motion.Profile{Ground: .8492822966507177, Facing: 1},
			Motions: map[string]motion.Clip{
				"idle":    {Durations: []float64{100}, Presentation: motion.Presentation{Scale: 1}},
				"raised":  {Durations: []float64{100}, Presentation: motion.Presentation{Scale: 1.2, Y: -.7}},
				"lowered": {Durations: []float64{100}, Presentation: motion.Presentation{Scale: 1.2, Y: .4}},
			},
		}, Frames: map[string]map[int][]*image.RGBA{}}
		for name := range v.Motions {
			v.Frames[name] = map[int][]*image.RGBA{64: {src}, 96: {src}}
		}
		rowHeight, baseline := previewRows([]*motion.Loaded{v}, height)
		for name := range v.Frames {
			for _, direction := range []int{-1, 1} {
				row := image.NewRGBA(image.Rect(0, 0, 512, rowHeight))
				probe := image.NewRGBA(image.Rect(0, 0, 512, 1024))
				v.Draw(row, name, 0, 128, baseline, height*3/2, height, direction)
				v.Draw(probe, name, 0, 128, baseline+400, height*3/2, height, direction)
				for y := 0; y < probe.Bounds().Dy(); y++ {
					for x := 0; x < probe.Bounds().Dx(); x++ {
						if got, want := row.RGBAAt(x, y-400), probe.RGBAAt(x, y); got != want {
							t.Fatalf("height%d %s direction%d clipped/different at%d,%d", height, name, direction, x, y-400)
						}
					}
				}
			}
		}
	}
}
