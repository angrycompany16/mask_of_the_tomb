package rendering

import (
	"mask_of_the_tomb/internal/core/ebitenrenderutil"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth, GameHeight = 480.0, 270.0
	PixelScale            = 4.0
)

type LayerList struct {
	Background2 *ebiten.Image
	Background  *ebiten.Image
	Midground   *ebiten.Image
	Playerspace *ebiten.Image
	Foreground  *ebiten.Image
	UI          *ebiten.Image
	Overlay     *ebiten.Image
	layers      []*ebiten.Image
}

var ScreenLayers = NewLayerList(GameWidth, GameHeight)

func NewLayerList(w, h int) LayerList {
	layerList := LayerList{
		Background2: ebiten.NewImage(w, h),
		Background:  ebiten.NewImage(w, h),
		Midground:   ebiten.NewImage(w, h),
		Playerspace: ebiten.NewImage(w, h),
		Foreground:  ebiten.NewImage(w, h),
		Overlay:     ebiten.NewImage(w, h),
		// TODO: Maybe find a better way to do this
		UI: ebiten.NewImage(w*PixelScale, h*PixelScale),
	}

	layerList.layers = []*ebiten.Image{
		layerList.Background2,
		layerList.Background,
		layerList.Midground,
		layerList.Playerspace,
		layerList.Foreground,
		layerList.Overlay,
		layerList.UI,
	}

	return layerList
}

func (r *LayerList) Draw(screen *ebiten.Image) {
	for _, layer := range r.layers {
		scaleFactor := GameWidth * PixelScale / float64(layer.Bounds().Dx())
		ebitenrenderutil.DrawAtScaled(layer, screen, 0, 0, scaleFactor, scaleFactor)
		layer.Clear()
	}
}

func (r *LayerList) DrawOnto(other *LayerList, x, y float64) {
	for i, layer := range r.layers {
		ebitenrenderutil.DrawAt(layer, other.layers[i], x, y)
	}
}
