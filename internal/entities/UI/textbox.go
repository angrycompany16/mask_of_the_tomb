package ui

import (
	"bytes"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/rendering"

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

type textbox struct {
	text             string
	posX, posY       float64
	color            ColorPair
	font             *text.GoTextFaceSource
	fontSize         float64
	lineSpacing      float64
	primaryAlign     text.Align
	secondaryAlign   text.Align
	screenAlign      ScreenAlign
	shadowX, shadowY float64
}

func (t *textbox) draw() {
	opText := &text.DrawOptions{}

	if t.screenAlign == Centered {
		s := rendering.RenderLayers.UI.Bounds().Size()
		opText.GeoM.Translate(float64(s.X)/2, float64(s.Y)/2)
	}

	opText.LayoutOptions.LineSpacing = 10.0
	opText.LayoutOptions.PrimaryAlign = t.primaryAlign
	opText.LayoutOptions.SecondaryAlign = t.secondaryAlign
	opText.GeoM.Translate(t.posX+t.shadowX, t.posY)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.color.DarkColor)

	text.Draw(rendering.RenderLayers.UI, t.text, &text.GoTextFace{
		Source: t.font,
		Size:   t.fontSize,
	}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.color.BrightColor)

	opText.GeoM.Translate(-t.shadowX, t.shadowY)
	text.Draw(rendering.RenderLayers.UI, t.text, &text.GoTextFace{
		Source: t.font,
		Size:   t.fontSize,
	}, opText)
}

func newTextBoxSimple(_text string, fontSize, posX, posY, lineSpacing float64, textAlign text.Align, screenAlign ScreenAlign) *textbox {
	return &textbox{
		text:           _text,
		posX:           posX,
		posY:           posY,
		color:          TextColorNormal,
		font:           errs.Must(text.NewGoTextFaceSource(bytes.NewReader(mainFont))),
		fontSize:       fontSize,
		lineSpacing:    lineSpacing,
		primaryAlign:   textAlign,
		secondaryAlign: textAlign,
		screenAlign:    screenAlign,
		shadowX:        DefaultShadowX,
		shadowY:        DefaultShadowY,
	}
}
