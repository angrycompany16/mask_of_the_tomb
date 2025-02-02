package player

import (
	"image/color"
	ebitenrenderutil "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/rendering"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	overlayColor = []uint8{20, 16, 19}
)

type damageOverlay struct {
	image       *ebiten.Image
	alpha       float64
	targetAlpha float64
}

func (d *damageOverlay) Update() {
	d.alpha = utils.Lerp(d.alpha, d.targetAlpha, 0.01)
}

func (d *damageOverlay) Draw() {
	alpha := uint8(d.alpha * 255)
	d.image.Fill(color.RGBA{overlayColor[0], overlayColor[1], overlayColor[2], alpha})
	ebitenrenderutil.DrawAt(d.image, rendering.RenderLayers.Overlay, 0, 0)
}

func newDamageOverlay() damageOverlay {
	return damageOverlay{
		image:       ebiten.NewImage(rendering.GameWidth, rendering.GameHeight),
		alpha:       0.0,
		targetAlpha: 0.0,
	}
}
