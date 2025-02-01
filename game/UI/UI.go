package ui

import (
	"fmt"
	"image/color"
	"mask_of_the_tomb/files"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type UI struct {
	mainFont           *text.GoTextFaceSource
	mainTextColor      color.Color
	secondaryTextColor color.Color
	text               string
}

func (ui *UI) Init() {

}

func (ui *UI) Update() {

}

func (ui *UI) Draw(surf *ebiten.Image) {
	opText := &text.DrawOptions{}
	opText.LayoutOptions.LineSpacing = 10.0
	// opText.LayoutOptions.PrimaryAlign = text.AlignCenter
	// opText.LayoutOptions.SecondaryAlign = text.AlignCenter
	opText.GeoM.Translate(24, 24)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(ui.secondaryTextColor)

	text.Draw(surf, ui.text, &text.GoTextFace{
		Source: ui.mainFont,
		Size:   48,
	}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(ui.mainTextColor)

	opText.GeoM.Translate(4, -4)
	text.Draw(surf, ui.text, &text.GoTextFace{
		Source: ui.mainFont,
		Size:   48,
	}, opText)
}

func (ui *UI) GenerateScoreMessage(score int) string {
	return fmt.Sprintf("YOUR SCORE IS: %d", score)
}

func (ui *UI) SetText(text string) {
	ui.text = text
}

func NewUi() *UI {
	return &UI{
		mainFont:           files.LazyFont(mainFont),
		mainTextColor:      color.RGBA{205, 247, 226, 255},
		secondaryTextColor: color.RGBA{199, 176, 139, 255},
		text:               "",
	}
}
