package rendering

import "github.com/hajimehoshi/ebiten/v2"

type Ctx struct {
	Dst              *ebiten.Image
	CamX, CamY       float64
	PlayerX, PlayerY float64
}

func WithLayer(drawCtx Ctx, layer *ebiten.Image) Ctx {
	drawCtx.Dst = layer
	return drawCtx
}
