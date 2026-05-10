package overlappingmodel

import (
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/wfc"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
)

type OverlappingModel struct {
	*graphic.Graphic
	WFC       *wfc.OverlappingModelWFC
	result    *ebiten.Image
	target    renderer.RenderTarget
	drawOrder int
}

func (w *OverlappingModel) Init(cmd *commands.Commands) {
	w.Graphic.Init(cmd)
	w.WFC.Generate()
	w.result = w.WFC.DrawTileAtlas()
}

func (w *OverlappingModel) Update(cmd *commands.Commands) {
	w.Graphic.Update(cmd)

	gPosX, gPosY := w.Transform2D.GetPos(false)
	gAngle := w.Transform2D.GetAngle(false)
	gScaleX, gScaleY := w.Transform2D.GetScale(false)

	cmd.Renderer.Request(opgen.PosRotScale(
		w.result,
		gPosX, gPosY,
		gAngle,
		gScaleX, gScaleY,
		0.0, 0.0,
	), w.result, w.target, w.drawOrder)
}

func NewOverlappingModel(
	graphic *graphic.Graphic,
	wfcSetup *wfc.OverlappingModelWFC,
	target renderer.RenderTarget,
	drawOrder int,
) *OverlappingModel {
	return &OverlappingModel{
		Graphic:   graphic,
		WFC:       wfcSetup,
		target:    target,
		drawOrder: drawOrder,
	}
}
