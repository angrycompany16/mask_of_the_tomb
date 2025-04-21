package node

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputField struct {
	Button `yaml:",inline"`
}

func (i *InputField) SetSelected() {
	i.selected = true
	i.Color = i.SelectedColor
}

func (i *InputField) SetDeselected() {
	i.selected = false
	i.Color = i.NormalColor
}

func (i *InputField) Confirm() {
	// Raise confirm event
}

func (i *InputField) Update(confirmations map[string]bool) {
	if !i.selected {
		return
	}

	keys := make([]rune, 0)
	keys = ebiten.AppendInputChars(keys)
	var b strings.Builder
	b.WriteString(i.Textbox.Text)
	// TODO: Raise event
	for _, key := range keys {
		b.WriteString(string(key))
	}
	i.Textbox.Text = b.String()

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(i.Textbox.Text) > 0 {
		i.Textbox.Text = i.Textbox.Text[0 : len(i.Textbox.Text)-1]
	}
}

func (i *InputField) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(i.Width, i.Height, parentWidth, parentHeight)
	i.DrawChildren(offsetX+i.PosX, offsetY+i.PosY, w, h)
	i.Button.Draw(offsetX, offsetY, w, h)
}

func (i *InputField) Reset() {
	i.Textbox.Text = ""
	i.ResetChildren()
}
