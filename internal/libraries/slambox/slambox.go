package slambox

import (
	"mask_of_the_tomb/internal/core/maths"
)

// A box that can be slammed through a SlamboxEnvironment. Can be connected with other Slamboxes or joined via links. Can also
// be connected to non-interactive entities.
type Slambox struct {
	Rect     *maths.Rect
	entities []*Attachment
	movebox  *Movebox
}

func (s *Slambox) Slam() {

}

type SlamboxCollisionInfo struct{}
