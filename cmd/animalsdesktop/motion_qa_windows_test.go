//go:build windows && motionqa

package main

import (
	"image"
	"testing"
)

func TestQACannotPassWithoutAcceptedCompletedActionForEachPet(t *testing.T) {
	for _, tc := range []struct {
		name            string
		requests        map[int]bool
		starts, returns map[int]int
	}{
		{"all rejected", map[int]bool{0: false, 1: false}, map[int]int{}, map[int]int{}},
		{"never entered", map[int]bool{0: true, 1: true}, map[int]int{}, map[int]int{}},
		{"never returned", map[int]bool{0: true, 1: true}, map[int]int{0: 1, 1: 1}, map[int]int{}},
		{"second pet missing", map[int]bool{0: true, 1: true}, map[int]int{0: 1, 1: 1}, map[int]int{0: 1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if len(actionCoverageFailures(2, tc.requests, tc.starts, tc.returns)) == 0 {
				t.Fatal("unexercised action falsely passed")
			}
		})
	}
	if got := actionCoverageFailures(2, map[int]bool{0: true, 1: false, 11: true}, map[int]int{0: 1, 1: 1}, map[int]int{0: 1, 1: 1}); len(got) != 0 {
		t.Fatal(got)
	}
}

func TestQACapturePropagatesWriteFailure(t *testing.T) {
	if err := writeQAPNG(t.TempDir(), image.NewRGBA(image.Rect(0, 0, 1, 1))); err == nil {
		t.Fatal("capture into directory silently passed")
	}
}

func TestQAResizeRequiresBothPausedAndPlayingCoverage(t *testing.T) {
	for _, counts := range []map[int]int{{}, {0: 4, 1: 4}, {0: 8, 1: 7}, {0: 9, 1: 8}} {
		if len(resizeCoverageFailures(2, counts)) == 0 {
			t.Fatal("partial or repeated resize coverage falsely passed", counts)
		}
	}
	if failures := resizeCoverageFailures(2, map[int]int{0: 8, 1: 8}); len(failures) != 0 {
		t.Fatal(failures)
	}
	if resizeStageDue(0, 5.6, 100, false, .2) || resizeStageDue(4, 5.6, 100, true, .2) || resizeStageDue(8, 14, 200, false, 0) {
		t.Fatal("resize ran in wrong phase or after coverage completed")
	}
	if !resizeStageDue(0, 5.6, 100, true, .2) || !resizeStageDue(7, 13.5, 260, false, 0) {
		t.Fatal("expected resize scenario was skipped")
	}
}

func TestQAResizePauseWaitsForEachActualActionEntry(t *testing.T) {
	starts := map[int]float64{}
	counts := []int{0, 0}
	// The second pet enters after the old 5.4..6.4s pause window expired.
	entry := []float64{5.0, 6.45}
	for tick := 0; tick < 180; tick++ {
		now := float64(tick) / 20
		for pet := range counts {
			active := now >= entry[pet]
			paused, pauseElapsed := resizePauseState(starts, pet, counts[pet], now, active)
			if !active {
				if _, seen := starts[pet]; paused || seen {
					t.Fatal("pause clock started before action entry", pet, now)
				}
			}
			if counts[pet] >= len(qaResizePercents) && paused {
				t.Fatal("pause held after four resizes", pet)
			}
			if resizeStageDue(counts[pet], now, 0, paused, pauseElapsed) {
				counts[pet]++
			}
		}
	}
	for pet, count := range counts {
		if count != 4 || starts[pet] != entry[pet] {
			t.Fatal("delayed pet missed independent pause/resize coverage", pet, count, starts)
		}
	}
}

func TestQAResizeCannotPassWithDroppedOrUnacknowledgedCapture(t *testing.T) {
	events := []map[string]any{{"tick": 80}, {"tick": 100}}
	for _, saved := range []map[int]bool{{}, {80: true}, {80: true, 100: false}} {
		if len(resizeCaptureFailures(events, saved)) == 0 {
			t.Fatal("missing resize PNG falsely passed", saved)
		}
	}
	if failures := resizeCaptureFailures(events, map[int]bool{80: true, 100: true}); len(failures) != 0 {
		t.Fatal(failures)
	}
	if events[1]["capture"] != "render.png" || events[1]["captureSaved"] != true {
		t.Fatal("resize evidence did not identify its successfully written PNG")
	}
}
