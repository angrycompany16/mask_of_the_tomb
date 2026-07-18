package slambox

import (
	"mask_of_the_tomb/internal/backend/maths"
)

type SlamboxGroup struct {
	slamboxes   []*maths.Rect
	CenterIndex int
}

func (s *SlamboxGroup) SetPos(x, y float64) {
	offsets := make([]maths.Vec2, len(s.slamboxes))

	center := s.slamboxes[s.CenterIndex]


	for i, slambox := range s.slamboxes {
		offsets[i] = maths.NewVec2(slambox.X - center.X, slambox.Y - center.Y)
	}

	for i := range s.slamboxes {
		pos := offsets[i].Plus(maths.NewVec2(x, y))
		s.slamboxes[i].SetPos(pos.XY())
	}
}

func NewSlamboxGroup(slamboxes []*maths.Rect, centerIndex int) *SlamboxGroup {
	newSlamboxGroup := SlamboxGroup{}
	newSlamboxGroup.slamboxes = slamboxes
	newSlamboxGroup.CenterIndex = centerIndex
	return &newSlamboxGroup
}
