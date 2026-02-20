package slambox

import "mask_of_the_tomb/internal/core/maths"

// A group of slamboxes connected as one.
type SlamboxGroup struct {
	slamboxes   []*Slambox
	centerIndex int
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

func (sg *SlamboxGroup) GetRequestedSlamDirection() maths.Direction {
	return sg.slamboxes[sg.centerIndex].GetRequestedSlamDirection()
}

func (sg *SlamboxGroup) RequestSlam(dir maths.Direction) {
	sg.slamboxes[sg.centerIndex].RequestSlam(dir)
}

func (sg *SlamboxGroup) GetSlamboxes() []*Slambox {
	return sg.slamboxes
}

func (sg *SlamboxGroup) GetCenterSlambox() *Slambox {
	return sg.slamboxes[sg.centerIndex]
}

func (sg *SlamboxGroup) GetSlamboxRects() []*maths.Rect {
	rects := make([]*maths.Rect, 0)
	for _, slambox := range sg.slamboxes {
		rects = append(rects, slambox.GetRect())
	}
	return rects
}

func (sg *SlamboxGroup) SetPosDirect(x, y float64) {
	for _, slambox := range sg.slamboxes {
		slambox.SetPosDirect(x, y)
	}
}

func NewSlamboxGroup(slamboxes []*Slambox, centerIndex int) *SlamboxGroup {
	newSlamboxGroup := SlamboxGroup{}
	newSlamboxGroup.slamboxes = slamboxes
	newSlamboxGroup.centerIndex = centerIndex
	return &newSlamboxGroup
}
