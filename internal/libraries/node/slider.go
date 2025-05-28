package node

import (
	"fmt"
	"mask_of_the_tomb/internal/core/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Slider struct {
	Button    `yaml:",inline"`
	Min       float64 `yaml:"min"`
	Max       float64 `yaml:"max"`
	Increment float64 `yaml:"increment"`
	Default   float64 `yaml:"default"`
	val       float64
}

func (s *Slider) Update(confirmations map[string]ConfirmInfo) {
	if !s.selected {
		return
	}

	changed := false
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		s.val -= s.Increment
		changed = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		s.val += s.Increment
		changed = true
	}

	// lovely
	s.val = maths.Clamp(s.val, s.Min, s.Max)
	s.Text = fmt.Sprintf("< %d%% >", int(s.val))

	if changed {
		confirmations[s.Name] = ConfirmInfo{
			IsConfirmed: true,
			SliderVal:   s.val,
		}
	}
}

func (s *Slider) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	s.Button.Draw(offsetX, offsetY, parentWidth, parentHeight)
}

func (s *Slider) Reset(overWriteInfo map[string]OverWriteInfo) {
	if overWrite, ok := overWriteInfo[s.Name]; ok {
		s.Text = fmt.Sprintf("< %d%% >", int(overWrite.SliderVal))
		s.val = overWrite.SliderVal
		s.Button.Reset(overWriteInfo)
	} else {
		s.Text = fmt.Sprintf("< %d%% >", int(s.val))
		s.Button.Reset(overWriteInfo)
	}
}

func (s *Slider) SetSelected() {
	s.Button.SetSelected()
}

func (s *Slider) SetDeselected() {
	s.Button.SetDeselected()
}
