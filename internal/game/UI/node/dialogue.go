package node

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/game/sound"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	dialogueSound = sound.NewEffectPlayer(assets.Text_scroll_ogg, sound.Ogg)
)

type Dialogue struct {
	Textbox         `yaml:",inline"`
	Name            string   `yaml:"Name"`
	RevealTime      float64  `yaml:"RevealTime"`
	Lines           []string `yaml:"Lines"`
	activeLine      int
	t               float64
	revealIndicator int
}

// This is cursed but we will receive the input inside the UI node
func (d *Dialogue) Update(confirmations map[string]bool) {
	d.UpdateChildren(confirmations)
	if d.activeLine == len(d.Lines) {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if d.revealIndicator == len(d.Lines[d.activeLine]) {
			d.Text = strings.Join([]string{
				d.Text, "\n",
			}, "")
			d.activeLine += 1
			d.revealIndicator = 0
			if d.activeLine == len(d.Lines) {
				confirmations[d.Name] = true
				return
			}
		} else {
			d.Text = strings.Join([]string{
				d.Text, string(d.Lines[d.activeLine][d.revealIndicator:len(d.Lines[d.activeLine])]),
			}, "")
			d.revealIndicator = len(d.Lines[d.activeLine])
		}
	}

	d.t += 1.0 / 60.0
	if d.t > d.RevealTime {
		if d.revealIndicator == len(d.Lines[d.activeLine]) {
			return
		}

		dialogueSound.Play()
		d.Text = strings.Join([]string{
			d.Text, string(d.Lines[d.activeLine][d.revealIndicator]),
		}, "")
		d.t = 0
		d.revealIndicator += 1
	}
}

func (d *Dialogue) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(d.Width, d.Height, parentWidth, parentHeight)
	d.DrawChildren(offsetX+d.PosX, offsetY+d.PosY, w, h)
	d.Textbox.Draw(offsetX, offsetY, parentWidth, parentHeight)
}

func (d *Dialogue) Reset() {
	d.ResetChildren()
}
