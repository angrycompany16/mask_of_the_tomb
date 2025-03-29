package textbox

import (
	"mask_of_the_tomb/internal/game/UI/colorpair"
	"mask_of_the_tomb/internal/game/UI/fonts"
	"mask_of_the_tomb/internal/game/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TODO: Maybe make some effort to reduce the amount of variables in this, or rewrite into asset
// type object
type ScreenAlign int

const (
	Centered ScreenAlign = iota
	TopLeft
)

type Textbox struct {
	Text           string              `yaml:"Text"`
	PosX           float64             `yaml:"PosX"`
	PosY           float64             `yaml:"PosY"`
	Color          colorpair.ColorPair `yaml:"Color"`
	FontStr        string              `yaml:"Color"`
	font           *text.GoTextFaceSource
	FontSize       float64     `yaml:"FontSize"`
	LineSpacing    float64     `yaml:"LineSpacing"`
	PrimaryAlign   text.Align  `yaml:"PrimaryAlign"`
	SecondaryAlign text.Align  `yaml:"SecondaryAlign"`
	ScreenAlign    ScreenAlign `yaml:"ScreenAlign"`
	ShadowX        float64     `yaml:"ShadowX"`
	ShadowY        float64     `yaml:"ShadowX"`
}

func (t *Textbox) Draw() {
	opText := &text.DrawOptions{}

	if t.ScreenAlign == Centered {
		s := rendering.RenderLayers.UI.Bounds().Size()
		opText.GeoM.Translate(float64(s.X)/2, float64(s.Y)/2)
	}

	opText.LayoutOptions.LineSpacing = 10.0
	opText.LayoutOptions.PrimaryAlign = t.PrimaryAlign
	opText.LayoutOptions.SecondaryAlign = t.SecondaryAlign
	opText.GeoM.Translate(t.PosX+t.ShadowX, t.PosY)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.DarkColor)

	text.Draw(rendering.RenderLayers.UI, t.Text, &text.GoTextFace{
		Source: t.font,
		Size:   t.FontSize,
	}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.BrightColor)

	opText.GeoM.Translate(-t.ShadowX, t.ShadowY)
	text.Draw(rendering.RenderLayers.UI, t.Text, &text.GoTextFace{
		Source: t.font,
		Size:   t.FontSize,
	}, opText)
}

func (t *Textbox) GetFont() {
	t.font = fonts.FontRegistry.M[t.FontStr]
}
