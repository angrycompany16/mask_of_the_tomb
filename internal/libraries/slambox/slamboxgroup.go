package slambox

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
func (sg *SlamboxGroup) Slam(x, y float64, i int) {
	anchor := sg.slamboxes[i]
	anchorX, anchorY := anchor.GetRect().TopLeft()
	for j, slambox := range sg.slamboxes {
		if j == i {
			continue
		}
		posX, posY := slambox.GetRect().TopLeft()
		offsetX, offsetY := posX-anchorX, posY-anchorY
		slambox.Slam(x+offsetX, y+offsetY)
	}
	anchor.Slam(x, y)
}

func (sg *SlamboxGroup) GetSlamboxes() []*Slambox {
	return sg.slamboxes
}

func NewSlamboxGroup(slamboxes []*Slambox) *SlamboxGroup {
	newSlamboxGroup := SlamboxGroup{}
	newSlamboxGroup.slamboxes = slamboxes
	return &newSlamboxGroup
}
