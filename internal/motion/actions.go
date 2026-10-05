package motion

import (
	"fmt"
	"math"
	"sort"
)

// ActionSequence is an explicitly adopted stationary action. Each step names
// dedicated artwork; start/return drawings are steps too, never a cross-fade.
// Moving actions (hop/scurry etc.) need their own displacement contract and are
// deliberately not interpreted as stationary from an action's name.
type ActionSequence struct {
	Stationary bool         `json:"stationary"`
	Automatic  bool         `json:"automatic,omitempty"`
	Steps      []ActionStep `json:"steps"`
}

type ActionStep struct {
	Motion string `json:"motion"`
	Cycles int    `json:"cycles"`
}

func sequenceClips(v Variant) ([]string, error) {
	if len(v.Actions) > 32 {
		return nil, fmt.Errorf("%s: too many action sequences", v.ID)
	}
	names := []string{"walk", "idle"}
	seen := map[string]bool{"walk": true, "idle": true}
	extras := []string{}
	for name, sequence := range v.Actions {
		if name == "" || name == "walk" || name == "idle" || !sequence.Stationary || len(sequence.Steps) == 0 || len(sequence.Steps) > 16 {
			return nil, fmt.Errorf("%s: invalid stationary sequence %q", v.ID, name)
		}
		total := 0.0
		for _, step := range sequence.Steps {
			clip, ok := v.Motions[step.Motion]
			if !ok || step.Motion == "walk" || step.Motion == "idle" || step.Cycles < 1 || step.Cycles > 16 {
				return nil, fmt.Errorf("%s/%s: invalid action step", v.ID, name)
			}
			total += Duration(clip.Durations) * float64(step.Cycles)
			if !seen[step.Motion] {
				extras = append(extras, step.Motion)
				seen[step.Motion] = true
			}
		}
		if !finite(total) || total <= 0 || total > 120000 {
			return nil, fmt.Errorf("%s/%s: invalid total duration", v.ID, name)
		}
	}
	sort.Strings(extras)
	return append(names, extras...), nil
}

func (v *Loaded) HasAction(name string) bool {
	if v == nil {
		return false
	}
	_, ok := v.Actions[name]
	return ok
}

func (v *Loaded) AutomaticActions() []string {
	return append([]string(nil), v.automaticActions...)
}

func (p *Player) Busy() bool { return p.Sequence != "" || p.PendingAction != "" }

// RequestAction preserves an active sequence and its return drawings. Requests
// cannot restart, stack or interrupt a half-completed action.
func (p *Player) RequestAction(v *Loaded, name string) bool {
	if p.Busy() || !v.HasAction(name) {
		return false
	}
	p.PendingAction = name
	p.Reaction = 0
	return true
}

// ReactWith uses dedicated art only when an adopted sequence is present. Older
// variants retain the existing short-translation behavior.
func (p *Player) ReactWith(v *Loaded) {
	if p.Cooldown > 0 {
		return
	}
	if v.HasAction("reaction") {
		if p.RequestAction(v, "reaction") {
			p.Cooldown = 1.2
		}
		return
	}
	// Do not turn a stationary action into a moving reaction mid-pose.
	if !p.Busy() {
		p.React()
	}
}

func (p *Player) startSequence(v *Loaded) {
	sequence, ok := v.Actions[p.PendingAction]
	if !ok || len(sequence.Steps) == 0 {
		p.PendingAction = ""
		return
	}
	p.Sequence = p.PendingAction
	p.PendingAction = ""
	p.SequenceStep = 0
	p.SequencePhase = 0
	p.Action = sequence.Steps[0].Motion
	p.Walking = false
	p.Reaction = 0
}

func (p *Player) advanceSequence(v *Loaded, dt float64) {
	sequence, ok := v.Actions[p.Sequence]
	if !ok || p.SequenceStep < 0 || p.SequenceStep >= len(sequence.Steps) {
		p.Sequence = ""
		p.SequencePhase = 0
		p.SequenceStep = 0
		p.Action = "idle"
		return
	}
	p.SequencePhase += dt * 1000
	for p.SequenceStep < len(sequence.Steps) {
		step := sequence.Steps[p.SequenceStep]
		duration := Duration(v.Motions[step.Motion].Durations) * float64(step.Cycles)
		if p.SequencePhase < duration {
			p.Action = step.Motion
			return
		}
		p.SequencePhase -= duration
		p.SequenceStep++
	}
	// Consume only the small remainder in idle. Never advance movement while
	// displaying the action's exit, even if the caller is requesting a walk.
	p.IdlePhase = math.Mod(p.SequencePhase, Duration(v.Motions["idle"].Durations))
	p.Sequence = ""
	p.SequenceStep = 0
	p.SequencePhase = 0
	p.Action = "idle"
	p.Remaining = math.Max(p.Remaining, .5)
}
