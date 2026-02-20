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
	slamRequest maths.Direction
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

func (s *Slambox) RequestSlam(dir maths.Direction) {
	s.slamRequest = dir
}

func (s *Slambox) GetRequestedSlamDirection() maths.Direction {
	return s.slamRequest
}

func (s *Slambox) GetRect() *maths.Rect {
	return s.rect
}

func (s *Slambox) GetTracker() *Tracker {
	return s.tracker
}

// Immediately sets the slambox position to (x, y)
func (s *Slambox) SetPosDirect(x, y float64) {
	s.rect.SetPos(x, y)
	s.tracker.SetPos(x, y)
}

func NewSlambox(rect *maths.Rect, moveSpeed float64) *Slambox {
	newSlambox := Slambox{}
	newSlambox.rect = rect
	newSlambox.tracker = NewTracker(moveSpeed, rect.Left(), rect.Top())
	newSlambox.slamRequest = maths.DirNone
	return &newSlambox
}

type SlamboxCollisionInfo struct{}
