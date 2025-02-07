package ui

import (
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type selectable struct {
	textbox     *textbox
	normalColor utils.ColorPair
	hoverColor  utils.ColorPair
	selected    bool
	name        string
}

func (s *selectable) setSelected() {
	s.selected = true
	s.textbox.color = s.hoverColor
}

func (s *selectable) getConfirm() bool {
	key := (inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.IsKeyJustReleased(ebiten.KeyEnter))
	return s.selected && key
}

func (s *selectable) setDeselected() {
	s.selected = false
	s.textbox.color = s.normalColor
}

func (s *selectable) draw() {
	s.textbox.draw()
}

func newSelectable(name string, textbox *textbox, normalColor, hoverColor utils.ColorPair) *selectable {
	return &selectable{
		textbox:     textbox,
		normalColor: normalColor,
		hoverColor:  hoverColor,
		name:        name,
	}
}
