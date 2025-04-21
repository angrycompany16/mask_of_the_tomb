package node

import (
	"mask_of_the_tomb/internal/game/UI/colorpair"
	"mask_of_the_tomb/internal/game/UI/fonts"
	"mask_of_the_tomb/internal/game/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ScreenAlign int

const (
	screenCentered ScreenAlign = iota
	screenTopLeft
)

type Textbox struct {
	NodeData       `yaml:",inline"`
	Text           string              `yaml:"Text"`
	Color          colorpair.ColorPair `yaml:"Color"`
	Font           fonts.FontYAML      `yaml:"Font"`
	FontSize       float64             `yaml:"FontSize"`
	LineSpacing    float64             `yaml:"LineSpacing"`
	PrimaryAlign   text.Align          `yaml:"PrimaryAlign"`
	SecondaryAlign text.Align          `yaml:"SecondaryAlign"`
	ScreenAlign    ScreenAlign         `yaml:"ScreenAlign"`
	ShadowX        float64             `yaml:"ShadowX"`
	ShadowY        float64             `yaml:"ShadowY"`
}

func (t *Textbox) Update(confirmations map[string]bool) {
	t.UpdateChildren(confirmations)
}

func (t *Textbox) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(t.Width, t.Height, parentWidth, parentHeight)
	t.DrawChildren(offsetX+t.PosX, offsetY+t.PosY, w, h)

	opText := &text.DrawOptions{}
	opText.LayoutOptions.LineSpacing = t.LineSpacing
	opText.LayoutOptions.PrimaryAlign = t.PrimaryAlign
	opText.LayoutOptions.SecondaryAlign = t.SecondaryAlign
	opText.GeoM.Translate(t.PosX+t.ShadowX+offsetX, t.PosY+offsetY)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.DarkColor)
	if t.ScreenAlign == screenCentered {
		opText.GeoM.Translate(parentWidth/2, parentHeight/2)
	}

	text.Draw(rendering.RenderLayers.UI,
		t.Text,
		&text.GoTextFace{
			Source: t.Font.GoTextFaceSource,
			Size:   t.FontSize,
		}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.BrightColor)

	opText.GeoM.Translate(-t.ShadowX, t.ShadowY)
	text.Draw(rendering.RenderLayers.UI, t.Text, &text.GoTextFace{
		Source: t.Font.GoTextFaceSource,
		Size:   t.FontSize,
	}, opText)
}

func (t *Textbox) Reset() {
	t.ResetChildren()
}
