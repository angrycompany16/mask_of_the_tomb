package selectable

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/cursor"
	"mask_of_the_tomb/internal/engine/actors/UI/textbox"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"
)

type Selectable struct {
	*textbox.Textbox
	NormalColor     colors.ColorPair
	SelectedColor   colors.ColorPair
	Selected        bool
	Hovered         bool
	OnHoverStart    *events.Event
	OnHoverEnd      *events.Event
	OnHoverStartBus *events.EventBus
	Callback        func(*commands.Commands)
}

func (b *Selectable) Update(cmd *commands.Commands) {
	b.Textbox.Update(cmd)
	scene, _ := commands.Get[engine.Scene](cmd)

	if cursorNode, ok := engine.GetNodeByType[*cursor.Cursor](scene); ok {
		cursor, _ := engine.As[*cursor.Cursor](cursorNode.GetValue())
		absRect := b.GetAbsRect()
		if absRect.Contains(cursor.Rect.Cx(), cursor.Rect.Cy()) {
			if !b.Hovered {
				b.OnHoverStart.Raise()
			}
			b.Hovered = true
		} else {
			if b.Hovered {
				b.OnHoverEnd.Raise()
			}
			b.Hovered = false
		}
	}

	UIControls := cmd.InputHandler.InputSchemes["UIControls"]
	if UIControls.PollAction("UIConfirm") && b.Selected {
		b.Callback(cmd)
	}

	if UIControls.PollAction("UIClick") && b.Selected {
		b.Callback(cmd)
	}
}

func (b *Selectable) SetSelectState(selected, suppressSound bool) {
	b.Selected = selected
	if selected {
		b.Textbox.Color = b.SelectedColor
	} else {
		b.Textbox.Color = b.NormalColor
	}
}

func defaultSelectable(textbox *textbox.Textbox) *Selectable {
	onHoverStart := events.NewEvent()
	return &Selectable{
		Textbox: textbox,
		NormalColor: colors.ColorPair{
			BrightColor: color.RGBA{255, 255, 255, 255},
			DarkColor:   color.RGBA{100, 100, 100, 255},
		},
		SelectedColor: colors.ColorPair{
			BrightColor: color.RGBA{255, 0, 0, 255},
			DarkColor:   color.RGBA{150, 50, 50, 255},
		},
		OnHoverStart:    onHoverStart,
		OnHoverEnd:      events.NewEvent(),
		OnHoverStartBus: events.NewBusFrom(onHoverStart),
		Callback:        func(c *commands.Commands) {},
	}
}

func NewSelectable(textbox *textbox.Textbox, options ...utils.Option[Selectable]) *Selectable {
	selectable := defaultSelectable(textbox)

	for _, option := range options {
		option(selectable)
	}

	return selectable
}

func WithNormalColor(color colors.ColorPair) utils.Option[Selectable] {
	return func(s *Selectable) {
		s.NormalColor = color
	}
}

func WithSelectedColor(color colors.ColorPair) utils.Option[Selectable] {
	return func(s *Selectable) {
		s.SelectedColor = color
	}
}

func WithCallback(callback func(cmd *commands.Commands)) utils.Option[Selectable] {
	return func(s *Selectable) {
		s.Callback = callback
	}
}
