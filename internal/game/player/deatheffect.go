package player

import (
	"image/color"
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/game/rendering"
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	overlayColor = []uint8{20, 16, 19}
)

type deathEffect struct {
	image       *ebiten.Image
	alpha       float64
	targetAlpha float64
}

func (d *deathEffect) Update() {
	d.alpha = maths.Lerp(d.alpha, d.targetAlpha, 0.01)
}

func (d *deathEffect) Draw() {
	alpha := uint8(d.alpha * 255)
	d.image.Fill(color.RGBA{overlayColor[0], overlayColor[1], overlayColor[2], alpha})
	ebitenrenderutil.DrawAt(d.image, rendering.RenderLayers.Overlay, 0, 0)
}

func newDamageOverlay() deathEffect {
	return deathEffect{
		image:       ebiten.NewImage(rendering.GameWidth, rendering.GameHeight),
		alpha:       0.0,
		targetAlpha: 0.0,
	}
}
