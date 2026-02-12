package slambox

import "mask_of_the_tomb/internal/core/maths"

// A group of slamboxes connected as one.
type SlamboxGroup struct {
	slamboxes []*Slambox
}

func (sg *SlamboxGroup) Update() {
	for _, slambox := range sg.slamboxes {
		slambox.Update()
	}
}

// Moves the group by moving the slambox at index i to the point x, y, and then
// maintaining the offset for the other slamboxes.
func (sg *SlamboxGroup) Slam(targetRects []maths.Rect) {
	for i, slambox := range sg.slamboxes {
		slambox.Slam(targetRects[i].Left(), targetRects[i].Top())
	}
}

func (sg *SlamboxGroup) GetSlamboxes() []*Slambox {
	return sg.slamboxes
}

func (sg *SlamboxGroup) GetSlamboxRects() []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for _, slambox := range sg.slamboxes {
		rects = append(rects, slambox.GetRect())
	}
	return rects
}

func NewSlamboxGroup(slamboxes []*Slambox) *SlamboxGroup {
	newSlamboxGroup := SlamboxGroup{}
	newSlamboxGroup.slamboxes = slamboxes
	return &newSlamboxGroup
}
