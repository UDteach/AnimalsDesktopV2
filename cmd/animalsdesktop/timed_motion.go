package main

import (
	appassets "animals-desktop/assets"
	"animals-desktop/internal/motion"
	"math"
	"math/rand"
	"os"
	"sync"
)

var timedOnce sync.Once
var timedStore *motion.Store

func timedVariant(id string) *motion.Loaded {
	// Classic launcher and regression/QA switch. IDs and settings never change.
	if os.Getenv("ANIMALSDESKTOP_MOTIONS") == "legacy" {
		return nil
	}
	timedOnce.Do(func() { timedStore, _ = motion.NewStore(appassets.FS, "motions/manifest.json") })
	return timedStore.Get(id)
}

func tickTimedPlayer(p *motion.Player, v *motion.Loaded, dt float64, x, height, maxX int, direction int, speed float64, stroll, paused bool) {
	if p.ID != v.ID {
		p.Reset(v.ID, float64(x), direction < 0)
	}
	// A settings/display change can reposition a pet outside this controller.
	if math.Abs(float64(x)-p.X) > 1.1 {
		p.X = float64(x)
		p.Velocity = 0
	}
	if paused || dt <= 0 {
		return
	}
	p.Cooldown = math.Max(0, p.Cooldown-dt)
	p.Reaction = math.Max(0, p.Reaction-dt)
	if p.Busy() {
		// Keep the idle/stroll clock still until all return drawings have played.
	} else if stroll {
		p.Remaining -= dt
		if p.Remaining <= 0 {
			p.Walking = !p.Walking
			if p.Walking {
				p.Remaining = 3 + rand.Float64()*5
			} else {
				p.Remaining = 2 + rand.Float64()*4
				requestAutomaticAction(p, v)
			}
		}
	} else {
		p.Walking = false
		if len(v.AutomaticActions()) > 0 && p.Reaction == 0 {
			p.Remaining -= dt
			if p.Remaining <= 0 {
				p.Remaining = 3 + rand.Float64()*5
				requestAutomaticAction(p, v)
			}
		}
	}
	target := 0.0
	if p.Walking || p.Reaction > 0 {
		target = v.Profile.Stride * float64(height) / (motion.Duration(v.Motions["walk"].Durations) / 1000) * speed
	}
	if p.Reaction > 0 {
		target *= 1.35
	}
	padding := float64(v.HorizontalPadding(height))
	minX := math.Min(padding, float64(maxX)/2)
	p.Advance(v, dt, target, float64(height), minX, math.Max(minX, float64(maxX)-padding), false)
}

func requestAutomaticAction(p *motion.Player, v *motion.Loaded) {
	actions := v.AutomaticActions()
	if len(actions) == 0 || p.Reaction > 0 {
		return
	}
	// One choice is quiet idle. Only explicitly adopted automatic actions run.
	choice := rand.Intn(len(actions) + 1)
	if choice < len(actions) {
		p.RequestAction(v, actions[choice])
	}
}
