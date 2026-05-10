package wfcentity

import (
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/wfc"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

// Problem: Adding modules must be done after the
type WFCEntity struct {
	*graphic.Graphic
	wfcSetup  *wfc.WaveFunctionSetup
	wfcResult *ebiten.Image
	target    renderer.RenderTarget
	drawOrder int
}

func (w *WFCEntity) Init(cmd *commands.Commands) {
	w.Graphic.Init(cmd)
	w.wfcSetup.Collapse(
		rand.IntN(w.wfcSetup.Height),
		rand.IntN(w.wfcSetup.Width),
	)
	w.wfcResult = w.wfcSetup.MakeImage()
}

func (w *WFCEntity) Update(cmd *commands.Commands) {
	w.Graphic.Update(cmd)

	gPosX, gPosY := w.Transform2D.GetPos(false)
	gAngle := w.Transform2D.GetAngle(false)
	gScaleX, gScaleY := w.Transform2D.GetScale(false)

	cmd.Renderer.Request(opgen.PosRotScale(
		w.wfcResult,
		gPosX, gPosY,
		gAngle,
		gScaleX, gScaleY,
		0.5, 0.5,
	), w.wfcResult, w.target, w.drawOrder)
}

func NewWFCEntity(
	graphic *graphic.Graphic,
	wfcSetup *wfc.WaveFunctionSetup,
	target renderer.RenderTarget,
	drawOrder int,
) *WFCEntity {
	return &WFCEntity{
		Graphic:   graphic,
		wfcSetup:  wfcSetup,
		target:    target,
		drawOrder: drawOrder,
	}
}
