package slider

import (
	"mask_of_the_tomb/internal/engine/actors/UI/selectable"
	"mask_of_the_tomb/internal/engine/commands"
)

type Slider struct {
	*selectable.Selectable
}

func (s *Slider) Update(cmd *commands.Commands) {
	s.Selectable.Update(cmd)

	// Check for key press and stuff
}
