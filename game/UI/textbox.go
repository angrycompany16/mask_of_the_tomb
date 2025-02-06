package ui

import (
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/rendering"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type textbox struct {
	text           string
	posX, posY     float64
	color          utils.ColorPair
	font           *text.GoTextFaceSource
	fontSize       float64
	lineSpacing    float64
	primaryAlign   text.Align
	secondaryAlign text.Align
}

func (t *textbox) draw() {
	opText := &text.DrawOptions{}
	opText.LayoutOptions.LineSpacing = 10.0
	// opText.LayoutOptions.PrimaryAlign = text.AlignCenter
	// opText.LayoutOptions.SecondaryAlign = text.AlignCenter
	opText.GeoM.Translate(t.posX, t.posY)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.color.BrightColor)

	text.Draw(rendering.RenderLayers.UI, t.text, &text.GoTextFace{
		Source: t.font,
		Size:   t.fontSize,
	}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.color.DarkColor)

	opText.GeoM.Translate(4, -4)
	text.Draw(rendering.RenderLayers.UI, t.text, &text.GoTextFace{
		Source: t.font,
		Size:   t.fontSize,
	}, opText)
}

func newTextBoxSimple(_text string, fontSize, posX, posY, lineSpacing float64, align text.Align) *textbox {
	return &textbox{
		text:           _text,
		posX:           posX,
		posY:           posY,
		color:          utils.TextColorNormal,
		font:           files.LazyFont(mainFont),
		fontSize:       fontSize,
		lineSpacing:    lineSpacing,
		primaryAlign:   align,
		secondaryAlign: align,
	}
}
