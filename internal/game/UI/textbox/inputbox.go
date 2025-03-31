package textbox

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Inputbox struct {
	Textbox *Textbox `yaml:"Textbox"`
	Name    string   `yaml:"Name"`
}

func (i *Inputbox) Update() (string, bool) {
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
		return i.Textbox.Text, true
	}

	return "", false
}

func (i *Inputbox) Draw() {
	i.Textbox.Draw()
}

func (i *Inputbox) Reset() {
	i.Textbox.Text = ""
}

func (i *Inputbox) Read() string {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return i.Textbox.Text
	}
	return ""
}
