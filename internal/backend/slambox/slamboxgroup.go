package slambox

import "mask_of_the_tomb/internal/backend/maths"

type SlamboxGroup struct {
	slamboxes   []*maths.Rect
	centerIndex int
}

func NewSlamboxGroup(slamboxes []*maths.Rect, centerIndex int) *SlamboxGroup {
	newSlamboxGroup := SlamboxGroup{}
	newSlamboxGroup.slamboxes = slamboxes
	newSlamboxGroup.centerIndex = centerIndex
	return &newSlamboxGroup
}
