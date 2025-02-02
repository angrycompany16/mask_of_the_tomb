package rendering

import (
	. "mask_of_the_tomb/ebitenRenderUtil"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth, GameHeight = 480.0, 270.0
	PixelScale            = 4.0
)

type renderLayers struct {
	Background  *ebiten.Image
	Midground   *ebiten.Image
	Playerspace *ebiten.Image
	Foreground  *ebiten.Image
	UI          *ebiten.Image
	Overlay     *ebiten.Image
}

var RenderLayers = renderLayers{
	Background:  ebiten.NewImage(GameWidth, GameHeight),
	Midground:   ebiten.NewImage(GameWidth, GameHeight),
	Playerspace: ebiten.NewImage(GameWidth, GameHeight),
	Foreground:  ebiten.NewImage(GameWidth, GameHeight),
	UI:          ebiten.NewImage(GameWidth*PixelScale, GameHeight*PixelScale),
	Overlay:     ebiten.NewImage(GameWidth, GameHeight),
}

func (r *renderLayers) Draw(screen *ebiten.Image) {
	// Probably stupid and bad
	layers := []*ebiten.Image{
		r.Background,
		r.Midground,
		r.Playerspace,
		r.Foreground,
		r.UI,
		r.Overlay,
	}

	for _, layer := range layers {
		scaleFactor := GameWidth * PixelScale / float64(layer.Bounds().Dx())
		DrawAtScaled(layer, screen, 0, 0, scaleFactor, scaleFactor)
		layer.Clear()
	}
}
