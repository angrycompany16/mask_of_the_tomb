package simpletiledmodel

import (
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/wfc"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type SimpleTiledWFC struct {
	*graphic.Graphic
	SimpleTileWFC *wfc.SimpleTileWFC
	result        *ebiten.Image
	target        renderer.RenderTarget
	drawOrder     int
}

func (w *SimpleTiledWFC) Init(cmd *commands.Commands) {
	w.Graphic.Init(cmd)
	w.SimpleTileWFC.Collapse(
		rand.IntN(w.SimpleTileWFC.Height),
		rand.IntN(w.SimpleTileWFC.Width),
	)
	w.result = w.SimpleTileWFC.MakeImage()
}

func (w *SimpleTiledWFC) Update(cmd *commands.Commands) {
	w.Graphic.Update(cmd)

	gPosX, gPosY := w.Transform2D.GetPos(false)
	gAngle := w.Transform2D.GetAngle(false)
	gScaleX, gScaleY := w.Transform2D.GetScale(false)

	cmd.Renderer.Request(opgen.PosRotScale(
		w.result,
		gPosX, gPosY,
		gAngle,
		gScaleX, gScaleY,
		0.5, 0.5,
	), w.result, w.target, w.drawOrder)
}

func NewSimpleTiledWFC(
	graphic *graphic.Graphic,
	wfcSetup *wfc.SimpleTileWFC,
	target renderer.RenderTarget,
	drawOrder int,
) *SimpleTiledWFC {
	return &SimpleTiledWFC{
		Graphic:       graphic,
		SimpleTileWFC: wfcSetup,
		target:        target,
		drawOrder:     drawOrder,
	}
}
