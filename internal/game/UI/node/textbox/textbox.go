package textbox

import (
	"bytes"
	"image/color"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/UI/colorpair"
	"mask_of_the_tomb/internal/game/UI/node"
	"mask_of_the_tomb/internal/game/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Textbox struct {
	*node.NodeData
	Text           string              `yaml:"Text"`
	Color          colorpair.ColorPair `yaml:"Color"`
	FontStr        string              `yaml:"Font"`
	font           *text.GoTextFaceSource
	FontSize       float64    `yaml:"FontSize"`
	LineSpacing    float64    `yaml:"LineSpacing"`
	PrimaryAlign   text.Align `yaml:"PrimaryAlign"`
	SecondaryAlign text.Align `yaml:"SecondaryAlign"`
	ShadowX        float64    `yaml:"ShadowX"`
	ShadowY        float64    `yaml:"ShadowY"`
}

func (t *Textbox) Update() {}

func (t *Textbox) Draw(offsetX, offsetY float64) {
	t.DrawChildren(offsetX+t.PosX, offsetY+t.PosY)

	opText := &text.DrawOptions{}
	opText.LayoutOptions.LineSpacing = 10.0
	opText.LayoutOptions.PrimaryAlign = t.PrimaryAlign
	opText.LayoutOptions.SecondaryAlign = t.SecondaryAlign
	opText.GeoM.Translate(t.PosX+t.ShadowX, t.PosY)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.DarkColor)

	text.Draw(rendering.RenderLayers.UI,
		t.Text,
		&text.GoTextFace{
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

func New() *Textbox {
	return &Textbox{
		NodeData: &node.NodeData{
			PosX:     0,
			PosY:     0,
			Children: make([]node.Node, 0),
		},
		Text: "Heisann!",
		Color: colorpair.ColorPair{
			BrightColor: color.RGBA{255, 255, 255, 255},
			DarkColor:   color.RGBA{255, 0, 0, 255},
		},
		font:           errs.Must(text.NewGoTextFaceSource(bytes.NewReader(assets.JSE_AmigaAMOS_ttf))),
		FontSize:       24,
		LineSpacing:    0,
		PrimaryAlign:   text.AlignStart,
		SecondaryAlign: text.AlignStart,
		ShadowX:        -4,
		ShadowY:        4,
	}
}
