package node

import (
	"fmt"
	"mask_of_the_tomb/internal/core/colors"
	"mask_of_the_tomb/internal/core/sound"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Not great, but cannot find a better solution
var SelectSound *sound.EffectPlayer

type Button struct {
	Textbox       `yaml:",inline"`
	NormalColor   colors.ColorPair `yaml:"NormalColor"`
	SelectedColor colors.ColorPair `yaml:"SelectedColor"`
	Name          string           `yaml:"Name"`
	selectSound   *sound.EffectPlayer
	selected      bool
}

func (b *Button) Update(confirmations map[string]ConfirmInfo) {
	confirmations[b.Name] = ConfirmInfo{IsConfirmed: inpututil.IsKeyJustPressed(ebiten.KeyEnter) && b.selected}
	b.UpdateChildren(confirmations)
}

func (b *Button) SetSelected() {
	// TODO: Annoyingly, this also runs when we call
	// switchActiveDisplay(), because that calls reset which sets the
	// selector pos to 0. Pls fix
	if !b.selected {
		if b.selectSound != nil {
			b.selectSound.Play()
		} else {
			fmt.Println("selectSound was nil, no sound will be played")
		}
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
