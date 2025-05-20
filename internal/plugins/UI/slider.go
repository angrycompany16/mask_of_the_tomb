package ui

import (
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/libraries/colors"
	"mask_of_the_tomb/internal/libraries/rendering"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RectAlign int

const (
	rectCentered RectAlign = iota
	rectTopLeft
)

type Slider struct {
	val                     float64
	NodeData                `yaml:",inline"`
	Min                     float64          `yaml:"min"`
	Max                     float64          `yaml:"max"`
	KnobRadius              float64          `yaml:"knobRadius"`
	LineThickness           float64          `yaml:"lineThickness"`
	BackgroundColorNormal   colors.YAMLColor `yaml:"backgroundColorNormal"`
	BackgroundColorSelected colors.YAMLColor `yaml:"backgroundColorSelected"`
	KnobColor               colors.YAMLColor `yaml:"knobColor"`
	LineColor               colors.YAMLColor `yaml:"lineColor"`
	RectAlign               RectAlign        `yaml:"rectAlign"`
	ScreenAlign             ScreenAlign      `yaml:"screenAlign"`
	selected                bool
}

func (s *Slider) Update(confirmations map[string]bool) {
	if !s.selected {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.val -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.val += 1
	}

	// lovely
	s.val = maths.Clamp(s.val, s.Min, s.Max)
}

func (s *Slider) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	rectX := s.PosX + offsetX
	rectY := s.PosY + offsetY
	if s.RectAlign == rectCentered {
		rectX -= s.Width / 2
		rectY -= s.Height / 2
	}
	if s.ScreenAlign == screenCentered {
		rectX += parentWidth / 2
		rectY += parentHeight / 2
	}
	bgColor := s.BackgroundColorNormal
	if s.selected {
		bgColor = s.BackgroundColorSelected
	}

	vector.DrawFilledRect(
		rendering.ScreenLayers.UI,
		float32(rectX),
		float32(rectY),
		float32(s.Width),
		float32(s.Height),
		bgColor,
		false,
	)

	padding := min(s.Width/2, s.Height/2)
	vector.StrokeLine(
		rendering.ScreenLayers.UI,
		float32(rectX+padding),
		float32(rectY+padding),
		float32(rectX+s.Width-padding),
		float32(rectY+padding),
		float32(s.LineThickness),
		s.LineColor,
		false,
	)

	t := s.val / (s.Max - s.Min)
	vector.DrawFilledCircle(
		rendering.ScreenLayers.UI,
		float32((1-t)*(rectX+padding)+t*(rectX+s.Width-padding)),
		float32(rectY+padding),
		float32(s.KnobRadius),
		s.KnobColor,
		false,
	)
}

func (s *Slider) Reset() {
	s.ResetChildren()
}

func (s *Slider) SetSelected() {
	s.selected = true
}

func (s *Slider) SetDeselected() {
	s.selected = false
}
