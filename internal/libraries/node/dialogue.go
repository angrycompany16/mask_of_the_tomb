package node

import (
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/core/threads"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var DialogueSound *sound.EffectPlayer

type Dialogue struct {
	Textbox         `yaml:",inline"`
	Name            string   `yaml:"Name"`
	RevealTime      float64  `yaml:"RevealTime"`
	Lines           []string `yaml:"Lines"`
	activeLine      int
	revealTicker    *time.Ticker
	revealIndicator int
}

// This is cursed but we will receive the input inside the UI node
func (d *Dialogue) Update(confirmations map[string]ConfirmInfo) {
	if d.revealTicker == nil {
		d.revealTicker = time.NewTicker(time.Duration(d.RevealTime * 1e9))
	}

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
				confirmations[d.Name] = ConfirmInfo{IsConfirmed: true}
				return
			}
		} else {
			d.Text = strings.Join([]string{
				d.Text, string(d.Lines[d.activeLine][d.revealIndicator:len(d.Lines[d.activeLine])]),
			}, "")
			d.revealIndicator = len(d.Lines[d.activeLine])
		}
	}

	if _, ticked := threads.Poll(d.revealTicker.C); ticked {
		if d.revealIndicator == len(d.Lines[d.activeLine]) {
			return
		}

		DialogueSound.Play()
		d.Text = strings.Join([]string{
			d.Text, string(d.Lines[d.activeLine][d.revealIndicator]),
		}, "")
		d.revealIndicator += 1
	}
}

func (d *Dialogue) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(d.Width, d.Height, parentWidth, parentHeight)
	d.DrawChildren(offsetX+d.PosX, offsetY+d.PosY, w, h)
	d.Textbox.Draw(offsetX, offsetY, parentWidth, parentHeight)
}

func (d *Dialogue) Reset(overWriteInfo map[string]OverWriteInfo) {
	d.ResetChildren(overWriteInfo)
}
