package slambox

import (
	"mask_of_the_tomb/internal/core/maths"
)

// A box that can be slammed through a SlamboxEnvironment. Can be connected with other Slamboxes or joined via links. Can also
// be connected to non-interactive entities.
type Slambox struct {
	rect        *maths.Rect
	attachments []*Attachment
	tracker     *Tracker
}

func (s *Slambox) Update() {
	s.tracker.Update()
	x, y := s.tracker.GetPos()
	s.rect.SetPos(x, y)
	// Update attachment positions
}

func (s *Slambox) Slam(x, y float64) {
	s.tracker.SetTarget(x, y)
}

func (s *Slambox) GetRect() *maths.Rect {
	return s.rect
}

func NewSlambox(rect *maths.Rect, moveSpeed float64) *Slambox {
	newSlambox := Slambox{}
	newSlambox.rect = rect
	newSlambox.tracker = NewTracker(moveSpeed, rect.Left(), rect.Top())
	return &newSlambox
}

type SlamboxCollisionInfo struct{}
