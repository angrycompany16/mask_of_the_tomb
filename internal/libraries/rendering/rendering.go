package rendering

import (
	"mask_of_the_tomb/internal/libraries/assets/ebitenrenderutil"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth, GameHeight = 480.0, 270.0
	PixelScale            = 4.0
)

type renderLayers struct {
	Background2 *ebiten.Image
	Background  *ebiten.Image
	Midground   *ebiten.Image
	Playerspace *ebiten.Image
	Foreground  *ebiten.Image
	UI          *ebiten.Image
	Overlay     *ebiten.Image
	layers      []*ebiten.Image
}

var RenderLayers = newRenderLayers()

func newRenderLayers() (rl renderLayers) {
	rl = renderLayers{
		Background:  ebiten.NewImage(GameWidth, GameHeight),
		Midground:   ebiten.NewImage(GameWidth, GameHeight),
		Playerspace: ebiten.NewImage(GameWidth, GameHeight),
		Foreground:  ebiten.NewImage(GameWidth, GameHeight),
		UI:          ebiten.NewImage(GameWidth*PixelScale, GameHeight*PixelScale),
		Overlay:     ebiten.NewImage(GameWidth, GameHeight),
	}

	rl.layers = []*ebiten.Image{
		rl.Background,
		rl.Midground,
		rl.Playerspace,
		rl.Foreground,
		rl.UI,
		rl.Overlay,
	}

	return
}

func (r *renderLayers) Draw(screen *ebiten.Image) {
	for _, layer := range r.layers {
		scaleFactor := GameWidth * PixelScale / float64(layer.Bounds().Dx())
		ebitenrenderutil.DrawAtScaled(layer, screen, 0, 0, scaleFactor, scaleFactor)
		layer.Clear()
	}
}
