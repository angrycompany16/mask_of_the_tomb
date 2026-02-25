package node

import (
	"mask_of_the_tomb/internal/core/colors"
	"mask_of_the_tomb/internal/core/sound_v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	Textbox       `yaml:",inline"`
	NormalColor   colors.ColorPair `yaml:"NormalColor"`
	SelectedColor colors.ColorPair `yaml:"SelectedColor"`
	Name          string           `yaml:"Name"`
	selected      bool
}

func (b *Button) Update(confirmations map[string]ConfirmInfo) {
	confirmations[b.Name] = ConfirmInfo{IsConfirmed: inpututil.IsKeyJustPressed(ebiten.KeyEnter) && b.selected}
	b.UpdateChildren(confirmations)
}

func (b *Button) SetSelected(suppresSound bool) {
	if !b.selected && !suppresSound {
		sound_v2.PlaySound("selectUI", "sfxMaster", 0)
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

func (b *Button) Reset(overWriteInfo map[string]OverWriteInfo) {
	b.selected = false
	b.Color = b.NormalColor
	b.ResetChildren(overWriteInfo)
}
