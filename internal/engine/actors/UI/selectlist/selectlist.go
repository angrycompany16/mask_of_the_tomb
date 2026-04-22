package selectlist

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/align"
	"mask_of_the_tomb/internal/engine/actors/UI/selectable"
	"mask_of_the_tomb/internal/engine/commands"
)

type SelectList struct {
	*align.Align
	selectIndex int
}

func (b *SelectList) Init(cmd *commands.Commands) {
	b.Align.Init(cmd)

	children := b.GetNode().GetChildren()
	b.selectIndex = maths.Clamp(b.selectIndex, 0, len(children))

	for i, child := range children {
		buttonActor, ok := engine.As[*selectable.Selectable](child.GetValue())
		if !ok {
			continue
		}
		if i == b.selectIndex {

			buttonActor.SetSelected(true)
		} else {
			buttonActor.SetDeselected()
		}
	}
}

func (b *SelectList) Update(cmd *commands.Commands) {
	b.Align.Update(cmd)
	if b.IsRow {
		if cmd.InputHandler.PollAction("UIRight") {
			b.selectIndex += 1
		} else if cmd.InputHandler.PollAction("UILeft") {
			b.selectIndex -= 1
		}
	} else {
		if cmd.InputHandler.PollAction("UIDown") {
			b.selectIndex += 1
		} else if cmd.InputHandler.PollAction("UIUp") {
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
		buttonActor, ok := engine.As[*selectable.Selectable](child.GetValue())
		if !ok {
			continue
		}
		if i == b.selectIndex {
			if !buttonActor.Selected {
				buttonActor.SetSelected(false)
			}
		} else {
			if buttonActor.Selected {
				buttonActor.SetDeselected()
			}
		}
	}
}

func NewButtonAlign(align *align.Align) *SelectList {
	return &SelectList{
		Align:       align,
		selectIndex: 0,
	}
}
