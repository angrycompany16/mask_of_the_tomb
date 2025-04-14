package selectlist

import "mask_of_the_tomb/internal/game/UI/node"

type SelectList struct {
	*node.NodeData
	SelectorPos int
}

func (s *SelectList) Update() {
	// Loop through children and set selected to true based on selector pos
}
