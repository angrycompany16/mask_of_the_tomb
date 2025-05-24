package ui

import (
	"image/color"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	OverlayColor = []uint8{20, 16, 19}
)

type ScreenFade struct {
	image *ebiten.Image
}

func (d *ScreenFade) Draw(t float64) {
	alpha := uint8(t * 255)
	d.image.Fill(color.RGBA{OverlayColor[0], OverlayColor[1], OverlayColor[2], alpha})
	ebitenrenderutil.DrawAt(d.image, rendering.ScreenLayers.Overlay, 0, 0)
}

func NewScreenFade() OverlayContent {
	return &ScreenFade{
		image: ebiten.NewImage(rendering.GAME_WIDTH, rendering.GAME_HEIGHT),
	}
}
