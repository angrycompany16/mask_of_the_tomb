package selectlist

import (
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/align"
	"mask_of_the_tomb/internal/engine/actors/UI/selectable"
	"mask_of_the_tomb/internal/engine/commands"
)

// Hmmm. How do we make it possible to select this object?
// i.e. If we have an option like quality, with three different options:
// low
// mid
// high
// How do we do that? There is no way of making a conditional SelectList
// as of right now
type SelectList struct {
	*align.Align
	Active          bool
	selectIndex     int
	hoverEventBuses []*events.EventBus
}

func (b *SelectList) Init(cmd *commands.Commands) {
	b.Align.Init(cmd)

	children := b.GetNode().GetChildren()
	b.selectIndex = maths.Clamp(b.selectIndex, 0, len(children))

	for i, child := range children {
		selectable, ok := engine.As[*selectable.Selectable](child.GetValue())
		if !ok {
			continue
		}

		if !b.Active {
			continue
		}

		if i == b.selectIndex {
			selectable.SetSelectState(true, true)
		} else {
			selectable.SetSelectState(false, true)
		}
	}
}

func (b *SelectList) Update(cmd *commands.Commands) {
	b.Align.Update(cmd)
	if !b.Active {
		return
	}

	UIControls := cmd.InputHandler.InputSchemes["UIControls"]
	if b.IsRow {
		if UIControls.PollAction("UIRight") {
			b.selectIndex += 1
		} else if UIControls.PollAction("UILeft") {
			b.selectIndex -= 1
		}
	} else {
		if UIControls.PollAction("UIDown") {
			b.selectIndex += 1
		} else if UIControls.PollAction("UIUp") {
			b.selectIndex -= 1
		}
	}

	children := b.GetNode().GetChildren()
	if b.selectIndex < 0 {
		b.selectIndex = len(children) - 1
	} else if b.selectIndex >= len(children) {
		b.selectIndex = 0
	}

	for i, child := range children {
		selectableNodes := child.GetChildrenFunc(func(n *node.Node[engine.Actor]) bool {
			_, ok := engine.As[*selectable.Selectable](n.GetValue())
			return ok
		})

		for _, selectableNode := range selectableNodes {
			selectable, _ := engine.As[*selectable.Selectable](selectableNode.GetValue())
			if _, raised := selectable.OnHoverStartBus.Poll(); raised {
				b.selectIndex = i
			}
		}

		if selectable, ok := engine.As[*selectable.Selectable](child.GetValue()); ok {
			if _, raised := selectable.OnHoverStartBus.Poll(); raised {
				b.selectIndex = i
			}
		}
	}

	for i, child := range children {
		var selectValue bool
		if i == b.selectIndex {
			selectValue = true
		} else {
			selectValue = false
		}

		selectableNodes := child.GetChildrenFunc(func(n *node.Node[engine.Actor]) bool {
			_, ok := engine.As[*selectable.Selectable](n.GetValue())
			return ok
		})

		for _, selectableNode := range selectableNodes {
			selectable, _ := engine.As[*selectable.Selectable](selectableNode.GetValue())

			if selectValue == selectable.Selected {
				continue
			}
			selectable.SetSelectState(selectValue, false)
		}

		if selectable, ok := engine.As[*selectable.Selectable](child.GetValue()); ok {
			if selectValue != selectable.Selected {
				selectable.SetSelectState(selectValue, false)
			}
		}

		selectListNodes := child.GetChildrenFunc(func(n *node.Node[engine.Actor]) bool {
			_, ok := engine.As[*SelectList](n.GetValue())
			return ok
		})

		for _, selectListNode := range selectListNodes {
			selectList, _ := engine.As[*SelectList](selectListNode.GetValue())
			selectList.Active = selectValue
		}

		if selectList, ok := engine.As[*SelectList](child.GetValue()); ok {
			selectList.Active = selectValue
		}
	}
}

func NewSelectList(align *align.Align) *SelectList {
	return &SelectList{
		Align:       align,
		selectIndex: 0,
		Active:      true,
	}
}
