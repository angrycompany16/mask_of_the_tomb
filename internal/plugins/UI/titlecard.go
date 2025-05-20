package ui

import (
	"bytes"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/libraries/rendering"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var textColor = []float32{1, 1, 1}

type TitleCard struct {
	text  string
	font  *text.GoTextFaceSource
	image *ebiten.Image
}

func (tc *TitleCard) Draw(t float64) {
	opText := &text.DrawOptions{}
	opText.LayoutOptions.LineSpacing = 40
	opText.LayoutOptions.PrimaryAlign = text.AlignCenter
	opText.LayoutOptions.SecondaryAlign = text.AlignCenter
	opText.ColorScale = ebiten.ColorScale{}

	opText.ColorScale.SetR(textColor[0] * float32(t))
	opText.ColorScale.SetG(textColor[1] * float32(t))
	opText.ColorScale.SetB(textColor[2] * float32(t))
	opText.ColorScale.SetA(float32(t))
	opText.GeoM.Translate(rendering.GameWidth*rendering.PixelScale/2, rendering.GameHeight*rendering.PixelScale/2)

	text.Draw(rendering.ScreenLayers.UI,
		tc.text,
		&text.GoTextFace{
			Source: tc.font,
			Size:   64,
		}, opText)
}

func (tc *TitleCard) ChangeText(text string) {
	tc.text = text
}

func NewTitleCard() OverlayContent {
	return &TitleCard{
		font:  errs.Must(text.NewGoTextFaceSource(bytes.NewReader(assets.JSE_AmigaAMOS_ttf))),
		image: ebiten.NewImage(rendering.GameWidth, rendering.GameHeight),
	}
}
