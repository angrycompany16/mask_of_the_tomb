package textbox

import (
	"mask_of_the_tomb/internal/game/UI/colorpair"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Selectable struct {
	Textbox     *Textbox            `yaml:"Textbox"`
	NormalColor colorpair.ColorPair `yaml:"NormalColor"`
	HoverColor  colorpair.ColorPair `yaml:"HoverColor"`
	selected    bool
	Name        string
}

func (s *Selectable) SetSelected() {
	s.selected = true
	s.Textbox.Color = s.HoverColor
}

func (s *Selectable) GetConfirm() bool {
	key := (inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.IsKeyJustReleased(ebiten.KeyEnter))
	return s.selected && key
}

func (s *Selectable) SetDeselected() {
	s.selected = false
	s.Textbox.Color = s.NormalColor
}

func (s *Selectable) Draw() {
	s.Textbox.Draw()
}

// func newSelectable(name string, textbox *Textbox, normalColor, hoverColor ColorPair) *Selectable {
// 	return &Selectable{
// 		Textbox:     textbox,
// 		NormalColor: normalColor,
// 		HoverColor:  hoverColor,
// 		Name:        name,
// 	}
// }
