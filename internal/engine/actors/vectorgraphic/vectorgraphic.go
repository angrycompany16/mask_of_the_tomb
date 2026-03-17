package vectorgraphic

import (
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"

	"github.com/hajimehoshi/ebiten/v2"
)

type VectorGraphic struct {
	*graphic.Graphic
	drawFunc  func(*ebiten.Image)
	image     *ebiten.Image
	layer     string `debug:"auto"`
	drawOrder int    `debug:"auto"`
}

// Note: In some cases this can be optimized by rendering only in init
func (v *VectorGraphic) Update(cmd *engine.Commands) {
	v.Graphic.Update(cmd)
	v.image.Clear()
	v.drawFunc(v.image)

	gPosX, gPosY := v.Transform2D.GetPos(false)
	camX, camY := v.GetCamera().WorldToCam(gPosX, gPosY, true)
	gAngle := v.Transform2D.GetAngle(false)
	gScaleX, gScaleY := v.Transform2D.GetScale(false)

	// Change this so that stuff is centered tbh
	cmd.Renderer().Request(opgen.PosScaleRot(
		v.image,
		camX, camY,
		gAngle,
		gScaleX, gScaleY,
		0.5, 0.5,
	), v.image, v.layer, v.drawOrder)
}

func NewVectorGraphic(
	graphic *graphic.Graphic,
	drawFunc func(*ebiten.Image),
	layer string,
	drawOrder int,
	width, height int, // not so great, but we need it
) *VectorGraphic {
	return &VectorGraphic{
		Graphic:   graphic,
		drawFunc:  drawFunc,
		image:     ebiten.NewImage(width, height),
		layer:     layer,
		drawOrder: drawOrder,
	}
}
