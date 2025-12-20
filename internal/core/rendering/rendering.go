package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Consider merging rendering and ebitenrenderutil

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type LayerList struct {
	Background2 *ebiten.Image
	Background  *ebiten.Image
	Midground   *ebiten.Image
	Playerspace *ebiten.Image
	Foreground  *ebiten.Image
	Overlay     *ebiten.Image
	GameplayUI  *ebiten.Image
	ScreenUI    *ebiten.Image
	layers      []*ebiten.Image
}

var ScreenLayers = NewLayerList(GAME_WIDTH, GAME_HEIGHT)

func NewLayerList(w, h int) LayerList {
	layerList := LayerList{
		Background2: ebiten.NewImage(w, h),
		Background:  ebiten.NewImage(w, h),
		Midground:   ebiten.NewImage(w, h),
		Playerspace: ebiten.NewImage(w, h),
		Foreground:  ebiten.NewImage(w, h),
		Overlay:     ebiten.NewImage(w, h),
		// TODO: Maybe find a better way to do this
		GameplayUI: ebiten.NewImage(w*PIXEL_SCALE, h*PIXEL_SCALE),
		ScreenUI:   ebiten.NewImage(w*PIXEL_SCALE, h*PIXEL_SCALE),
	}

	layerList.layers = []*ebiten.Image{
		layerList.Background2,
		layerList.Background,
		layerList.Midground,
		layerList.Playerspace,
		layerList.Foreground,
		layerList.Overlay,
		layerList.GameplayUI,
		layerList.ScreenUI,
	}

	return layerList
}

func (r *LayerList) Draw(screen *ebiten.Image) {
	for _, layer := range r.layers {
		scaleFactor := GAME_WIDTH * PIXEL_SCALE / float64(layer.Bounds().Dx())
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleFactor, scaleFactor)
		screen.DrawImage(layer, &op)
		layer.Clear()
	}
}

func (r *LayerList) DrawOnto(other *LayerList, x, y float64) {
	for i, layer := range r.layers {
		op := ebiten.DrawImageOptions{}
		other.layers[i].DrawImage(layer, &op)
	}
}
