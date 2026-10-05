//go:build windows

package main

import (
	"animals-desktop/internal/motion"
	"math"
)

func (a *petApp) tickTimedPet(index int, p *desktopPet, v *motion.Loaded, dt float64, dpi int) {
	w, h := a.petSpriteSize(index)
	// Scene positions use physical pixels on Windows; use this segment's DPI
	// for displacement as well as the sprite so gait survives monitor changes.
	ph := scaleForDPI(h, dpi)
	pw := scaleForDPI(w, dpi)
	if !a.bidirectional && p.timed.ID == v.ID {
		p.timed.Left = false
		if p.timed.X >= float64(max(0, a.sceneW-pw-v.HorizontalPadding(ph))) && a.sceneW > pw {
			p.timed.X = 0
			p.x = 0
			p.timed.Velocity = 0
		}
	}
	tickTimedPlayer(&p.timed, v, dt, p.x, ph, max(0, a.sceneW-pw), p.dir, float64(a.speed)/3, a.mode == modeRandom, a.overlayHidden)
	p.x = int(math.Round(p.timed.X))
	p.dir = 1
	if p.timed.Left {
		p.dir = -1
	}
	p.state = stateIdle
	if p.timed.Action == "walk" {
		p.state = stateWalk
	}
	p.item = noItem
	p.carryKind = noItem
	p.moveSpeed = 0
}
