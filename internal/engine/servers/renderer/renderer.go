package renderer

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Make camera part of the renderer

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type DrawRequest struct {
	op        *ebiten.DrawImageOptions
	src       *ebiten.Image
	layer     string
	drawOrder int
}

type Renderer struct {
	drawRequests []*DrawRequest
	layers       map[string]*ebiten.Image
}

func (r *Renderer) Request(op *ebiten.DrawImageOptions, src *ebiten.Image, layer string, drawOrder int) {
	r.drawRequests = append(r.drawRequests, &DrawRequest{
		op:        op,
		src:       src,
		layer:     layer,
		drawOrder: drawOrder,
	})
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	// Sort the slice before rendering. Nodes with the same draw order will be
	// drawn randomly
	slices.SortFunc(r.drawRequests, func(a *DrawRequest, b *DrawRequest) int {
		return a.drawOrder - b.drawOrder
	})

	for _, drawRequest := range r.drawRequests {
		r.layers[drawRequest.layer].DrawImage(drawRequest.src, drawRequest.op)
	}

	for _, layer := range r.layers {
		scaleFactor := GAME_WIDTH * PIXEL_SCALE / float64(layer.Bounds().Dx())
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleFactor, scaleFactor)
		screen.DrawImage(layer, &op)
		layer.Clear()
	}
}

func NewRenderer(w, h int) Renderer {
	renderer := Renderer{
		drawRequests: make([]*DrawRequest, 0),
		layers:       make(map[string]*ebiten.Image),
	}
	renderer.layers["Background2"] = ebiten.NewImage(w, h)
	renderer.layers["Background"] = ebiten.NewImage(w, h)
	renderer.layers["Midground"] = ebiten.NewImage(w, h)
	renderer.layers["Playerspace"] = ebiten.NewImage(w, h)
	renderer.layers["Foreground"] = ebiten.NewImage(w, h)
	renderer.layers["WorldUI"] = ebiten.NewImage(w*PIXEL_SCALE, h*PIXEL_SCALE)
	renderer.layers["Overlay"] = ebiten.NewImage(w, h)

	renderer.layers["GameplayUI"] = ebiten.NewImage(w*PIXEL_SCALE, h*PIXEL_SCALE)
	renderer.layers["ScreenUI"] = ebiten.NewImage(w*PIXEL_SCALE, h*PIXEL_SCALE)

	return renderer
}
