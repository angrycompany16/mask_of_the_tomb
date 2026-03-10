package renderer

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	om "github.com/wk8/go-ordered-map/v2"
)

// TODO: Make camera part of the renderer

type DrawRequest struct {
	op        *ebiten.DrawImageOptions
	src       *ebiten.Image
	layer     string
	drawOrder int
}

type Renderer struct {
	gameWidth, gameHeight float64
	pixelScale            float64
	drawRequests          []*DrawRequest
	Olayers               *om.OrderedMap[string, *ebiten.Image]
	layers                map[string]*ebiten.Image
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
	// Stable sort would lowkey be nice as it stops Z-fighting. But at the same time idk
	slices.SortFunc(r.drawRequests, func(a *DrawRequest, b *DrawRequest) int {
		return a.drawOrder - b.drawOrder
	})

	for _, drawRequest := range r.drawRequests {
		layer, ok := r.Olayers.Get(drawRequest.layer)
		if !ok {
			fmt.Println("Draw request failed - layer does not exist")
			continue
		}
		layer.DrawImage(drawRequest.src, drawRequest.op)
	}

	r.drawRequests = make([]*DrawRequest, 0)

	for pair := r.Olayers.Oldest(); pair != nil; pair = pair.Next() {
		layer := pair.Value
		scaleFactor := r.gameWidth * r.pixelScale / float64(layer.Bounds().Dx())
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleFactor, scaleFactor)
		screen.DrawImage(layer, &op)
		layer.Clear()
	}
}

func (r *Renderer) GetGameSize() (float64, float64) {
	return r.gameWidth, r.gameHeight
}

func (r *Renderer) GetPixelScale() float64 {
	return r.pixelScale
}

func NewRenderer(gameWidth, gameHeight, pixelScale int) *Renderer {
	renderer := Renderer{
		gameWidth:    float64(gameWidth),
		gameHeight:   float64(gameHeight),
		pixelScale:   float64(pixelScale),
		drawRequests: make([]*DrawRequest, 0),
		layers:       make(map[string]*ebiten.Image),
		Olayers:      om.New[string, *ebiten.Image](),
	}

	renderer.Olayers.Set("Background2", ebiten.NewImage(gameWidth, gameHeight))
	renderer.Olayers.Set("Background", ebiten.NewImage(gameWidth, gameHeight))
	renderer.Olayers.Set("Midground", ebiten.NewImage(gameWidth, gameHeight))
	renderer.Olayers.Set("Playerspace", ebiten.NewImage(gameWidth, gameHeight))
	renderer.Olayers.Set("Foreground", ebiten.NewImage(gameWidth, gameHeight))
	renderer.Olayers.Set("WorldUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))
	renderer.Olayers.Set("Overlay", ebiten.NewImage(gameWidth, gameHeight))

	renderer.Olayers.Set("GameplayUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))
	renderer.Olayers.Set("ScreenUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))
	renderer.Olayers.Set("EditorUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))

	return &renderer
}
