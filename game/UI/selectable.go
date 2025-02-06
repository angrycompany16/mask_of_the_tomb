package ui

import "mask_of_the_tomb/utils"

type selectable struct {
	textbox     *textbox
	normalColor utils.ColorPair
	hoverColor  utils.ColorPair
	selected    bool
}

func (s *selectable) setSelected() {
	s.selected = true
	s.textbox.color = s.hoverColor
}

func (s *selectable) setDeselected() {
	s.selected = false
	s.textbox.color = s.normalColor
}

func (s *selectable) draw() {
	s.textbox.draw()
}

func newSelectable(textbox *textbox, normalColor, hoverColor utils.ColorPair) *selectable {
	return &selectable{
		textbox:     textbox,
		normalColor: normalColor,
		hoverColor:  hoverColor,
	}
}
