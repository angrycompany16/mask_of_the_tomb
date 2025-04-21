package node

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/game/UI/colorpair"
	"mask_of_the_tomb/internal/game/sound"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var buttonSound = sound.NewEffectPlayer(assets.Select_ogg, sound.Ogg)

type Button struct {
	Textbox       `yaml:",inline"`
	NormalColor   colorpair.ColorPair `yaml:"NormalColor"`
	SelectedColor colorpair.ColorPair `yaml:"SelectedColor"`
	Name          string              `yaml:"Name"`
	selected      bool
}

func (b *Button) Update(confirmations map[string]bool) {
	confirmations[b.Name] = inpututil.IsKeyJustPressed(ebiten.KeyEnter) && b.selected
	b.UpdateChildren(confirmations)
}

func (b *Button) SetSelected() {
	if !b.selected {
		buttonSound.Play()
	}

	b.selected = true
	b.Color = b.SelectedColor
}

func (b *Button) SetDeselected() {
	b.selected = false
	b.Color = b.NormalColor
}

func (b *Button) Confirm() {
	// Raise confirm event
}

func (b *Button) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(b.Width, b.Height, parentWidth, parentHeight)
	b.DrawChildren(offsetX+b.PosX, offsetY+b.PosY, w, h)
	b.Textbox.Draw(offsetX, offsetY, w, h)
}
