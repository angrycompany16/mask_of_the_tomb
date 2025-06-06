package node

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/colors"
	"mask_of_the_tomb/internal/core/rendering"

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
	Text           string               `yaml:"Text"`
	Color          colors.ColorPair     `yaml:"Color"`
	Font           assetloader.FontYAML `yaml:"Font"`
	FontSize       float64              `yaml:"FontSize"`
	LineSpacing    float64              `yaml:"LineSpacing"`
	PrimaryAlign   text.Align           `yaml:"PrimaryAlign"`
	SecondaryAlign text.Align           `yaml:"SecondaryAlign"`
	ScreenAlign    ScreenAlign          `yaml:"ScreenAlign"`
	ShadowX        float64              `yaml:"ShadowX"`
	ShadowY        float64              `yaml:"ShadowY"`
}

func (t *Textbox) Update(confirmations map[string]ConfirmInfo) {
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

	text.Draw(rendering.ScreenLayers.ScreenUI,
		t.Text,
		&text.GoTextFace{
			Source: t.Font.GoTextFaceSource,
			Size:   t.FontSize,
		}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.BrightColor)

	opText.GeoM.Translate(-t.ShadowX, t.ShadowY)
	text.Draw(rendering.ScreenLayers.ScreenUI, t.Text, &text.GoTextFace{
		Source: t.Font.GoTextFaceSource,
		Size:   t.FontSize,
	}, opText)
}

func (t *Textbox) Reset(overWriteInfo map[string]OverWriteInfo) {
	t.ResetChildren(overWriteInfo)
}
