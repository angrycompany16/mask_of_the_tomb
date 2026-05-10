package renderer

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	om "github.com/wk8/go-ordered-map/v2"
)

type requestType int

const (
	DRAW_IMAGE requestType = iota
	COLORM
)

type RenderType int

const (
	SCREEN RenderType = iota
	TEXTURE
)

type RenderTarget struct {
	Type RenderType
	Name string
}

type DrawRequest struct {
	requestType  requestType
	renderTarget RenderTarget
	op           *ebiten.DrawImageOptions
	colorm       colorm.ColorM
	colormOp     *colorm.DrawImageOptions
	src          *ebiten.Image
	drawOrder    int
}

type Renderer struct {
	gameWidth, gameHeight float64
	pixelScale            float64
	drawRequests          []*DrawRequest
	layers                *om.OrderedMap[string, *ebiten.Image]
	Textures              map[string]*ebiten.Image // Images that aren't rendered to the screen
}

func (r *Renderer) Request(op *ebiten.DrawImageOptions, src *ebiten.Image, renderTarget RenderTarget, drawOrder int) {
	r.drawRequests = append(r.drawRequests, &DrawRequest{
		requestType:  DRAW_IMAGE,
		op:           op,
		src:          src,
		renderTarget: renderTarget,
		drawOrder:    drawOrder,
	})
}

func (r *Renderer) RequestColorM(colorm colorm.ColorM, op *colorm.DrawImageOptions, src *ebiten.Image, renderTarget RenderTarget, drawOrder int) {
	r.drawRequests = append(r.drawRequests, &DrawRequest{
		requestType:  COLORM,
		colorm:       colorm,
		colormOp:     op,
		src:          src,
		renderTarget: renderTarget,
		drawOrder:    drawOrder,
	})
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	// Sort the slice before rendering. Nodes with the same draw order will be
	// drawn randomly
	// Stable sort would lowkey be nice as it stops Z-fighting. But at the same time idk
	// BETTER: It's most likely a lot more performant to sort the
	// list as it is being made, rather than every time we call draw!
	// NOTE: We actually need to do texture requests first.
	// Cause some screen stuff might need to use textures during the same
	// frame
	slices.SortFunc(r.drawRequests, func(a *DrawRequest, b *DrawRequest) int {
		return a.drawOrder - b.drawOrder
	})

	for _, drawRequest := range r.drawRequests {
		var target *ebiten.Image
		var ok bool

		switch drawRequest.renderTarget.Type {
		case SCREEN:
			target, ok = r.layers.Get(drawRequest.renderTarget.Name)
		case TEXTURE:
			target, ok = r.Textures[drawRequest.renderTarget.Name]
		}

		if !ok {
			fmt.Println("Draw request failed - layer does not exist")
			continue
		}

		switch drawRequest.requestType {
		case DRAW_IMAGE:
			target.DrawImage(drawRequest.src, drawRequest.op)
		case COLORM:
			colorm.DrawImage(target, drawRequest.src, drawRequest.colorm, drawRequest.colormOp)
		}
	}

	r.drawRequests = make([]*DrawRequest, 0)

	for pair := r.layers.Oldest(); pair != nil; pair = pair.Next() {
		layer := pair.Value
		scaleFactor := r.gameWidth * r.pixelScale / float64(layer.Bounds().Dx())
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleFactor, scaleFactor)
		screen.DrawImage(layer, &op)
		layer.Clear()
	}
}

func (r *Renderer) ClearTextures() {
	for _, texture := range r.Textures {
		texture.Clear()
	}
}

func (r *Renderer) GetGameSize() (float64, float64) {
	return r.gameWidth, r.gameHeight
}

func (r *Renderer) GetPixelScale() float64 {
	return r.pixelScale
}

func TextureTarget(name string) RenderTarget {
	return RenderTarget{
		Type: TEXTURE,
		Name: name,
	}
}

func ScreenTarget(name string) RenderTarget {
	return RenderTarget{
		Type: SCREEN,
		Name: name,
	}
}

func NewRenderer(gameWidth, gameHeight, pixelScale int, fullScreen, hideCursor bool) *Renderer {
	renderer := Renderer{
		gameWidth:    float64(gameWidth),
		gameHeight:   float64(gameHeight),
		pixelScale:   float64(pixelScale),
		drawRequests: make([]*DrawRequest, 0),
		layers:       om.New[string, *ebiten.Image](),
		Textures:     make(map[string]*ebiten.Image, 0),
	}

	renderer.layers.Set("Background2", ebiten.NewImage(gameWidth, gameHeight))
	renderer.layers.Set("Background", ebiten.NewImage(gameWidth, gameHeight))
	renderer.layers.Set("Midground", ebiten.NewImage(gameWidth, gameHeight))
	renderer.layers.Set("Playerspace", ebiten.NewImage(gameWidth, gameHeight))
	renderer.layers.Set("Foreground", ebiten.NewImage(gameWidth, gameHeight))
	renderer.layers.Set("WorldUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))
	renderer.layers.Set("Overlay", ebiten.NewImage(gameWidth, gameHeight))

	renderer.layers.Set("GameplayUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))
	renderer.layers.Set("ScreenUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))
	renderer.layers.Set("EditorUI", ebiten.NewImage(gameWidth*pixelScale, gameHeight*pixelScale))

	ebiten.SetWindowSize(gameWidth*pixelScale, gameHeight*pixelScale)

	ebiten.SetFullscreen(fullScreen)

	if hideCursor {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}

	return &renderer
}
