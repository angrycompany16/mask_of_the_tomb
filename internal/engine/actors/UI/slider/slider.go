package slider

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine/actors/UI/selectable"
	"mask_of_the_tomb/internal/engine/commands"
)

type Slider struct {
	*selectable.Selectable
	Min     float64
	Max     float64
	Step    float64
	Default float64
	val     float64
}

func (s *Slider) Init(cmd *commands.Commands) {
	s.Selectable.Init(cmd)
	s.Text = fmt.Sprintf("< %d%% >", int(s.val))
}

func (s *Slider) Update(cmd *commands.Commands) {
	s.Selectable.Update(cmd)

	if !s.Selected {
		return
	}

	UIControls := cmd.InputHandler.InputSchemes["UIControls"]
	if UIControls.PollAction("UILeft") {
		s.val -= s.Step
	}

	if UIControls.PollAction("UIRight") {
		s.val += s.Step
	}

	// lovely
	s.val = maths.Clamp(s.val, s.Min, s.Max)
	s.Text = fmt.Sprintf("< %d%% >", int(s.val))
}

func NewSlider(selectable *selectable.Selectable, min, max, step, _default float64) *Slider {
	return &Slider{
		Selectable: selectable,
		Min:        min,
		Max:        max,
		Step:       step,
		Default:    _default,
		val:        _default,
	}
}
