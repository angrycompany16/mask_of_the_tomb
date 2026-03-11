package vectorgraphic

import (
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"

	"github.com/hajimehoshi/ebiten/v2"
)

type VectorGraphic struct {
	*transform2D.Transform2D
	drawFunc  func(*ebiten.Image)
	image     *ebiten.Image
	layer     string `debug:"auto"`
	drawOrder int    `debug:"auto"`
}

// Note: In some cases this can be optimized by rendering only in init
func (v *VectorGraphic) Update(servers *engine.Servers) {
	v.Transform2D.Update(servers)
	v.image.Clear()
	v.drawFunc(v.image)

	gPosX, gPosY := v.Transform2D.GetPos(false)
	gAngle := v.Transform2D.GetAngle(false)
	gScaleX, gScaleY := v.Transform2D.GetScale(false)

	// Change this so that stuff is centered tbh
	servers.Renderer().Request(opgen.PosScaleRot(
		v.image,
		gPosX, gPosY,
		gAngle,
		gScaleX, gScaleY,
		0.5, 0.5,
	), v.image, v.layer, v.drawOrder)
}

func NewVectorGraphic(
	transform2d *transform2D.Transform2D,
	drawFunc func(*ebiten.Image),
	layer string,
	drawOrder int,
	width, height int, // not so great, but we need it
) *VectorGraphic {
	return &VectorGraphic{
		Transform2D: transform2d,
		drawFunc:    drawFunc,
		image:       ebiten.NewImage(width, height),
		layer:       layer,
		drawOrder:   drawOrder,
	}
}
