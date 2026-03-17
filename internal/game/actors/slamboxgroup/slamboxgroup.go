package slamboxgroup

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
)

// A group of slamboxes connected as one.
// Contains slamboxes as children?
type SlamboxGroup struct {
	*transform2D.Transform2D
	backendID   int
	centerIndex int
	inChain     bool
}

func (sg *SlamboxGroup) Init(cmd *engine.Commands) {
	childSlamboxes := make([]*maths.Rect, 0)
	children := sg.Node.GetNode().GetChildren()
	for i, child := range children {
		childSlambox, ok := engine.GetActor[*slamboxactor.Slambox](child.GetValue())

		if !ok {
			continue
		}

		if childSlambox.IsCenter() {
			sg.centerIndex = i
		}

		childSlamboxes = append(childSlamboxes, childSlambox.GetRect())
	}
	sg.backendID = cmd.SlamboxEnv().AddSlamboxGroup(slambox.NewSlamboxGroup(childSlamboxes, sg.centerIndex))
}

func (sg *SlamboxGroup) Update(cmd engine.Commands) {
	children := sg.Node.GetNode().GetChildren()
	var slamDir maths.Direction
	for _, child := range children {
		childSlambox, ok := engine.GetActor[*slamboxactor.Slambox](child.GetValue())
		if !ok {
			continue
		}

		slamDir = childSlambox.GetSlamRequest()
	}
	if slamDir == maths.DirNone {
		return
	}

	newRects := cmd.SlamboxEnv().SlamSlamboxGroup(sg.backendID, slamDir)

	for i, child := range children {
		childSlambox, ok := engine.GetActor[*slamboxactor.Slambox](child.GetValue())
		if !ok {
			continue
		}

		childSlambox.SetTarget(newRects[i].Cx(), newRects[i].Cy())
	}
}

func NewSlamboxGroup(transform2D *transform2D.Transform2D, inChain bool) *SlamboxGroup {
	return &SlamboxGroup{
		Transform2D: transform2D,
	}
}
