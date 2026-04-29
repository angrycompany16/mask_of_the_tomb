package vectorgraphic

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type VectorGraphic struct {
	*graphic.Graphic
	drawFunc       func(*ebiten.Image)
	image          *ebiten.Image
	pivotX, pivotY float64 `debug:"auto"`
	drawOrder      int     `debug:"auto"`
	target         renderer.RenderTarget
}

// Note: In some cases this can be optimized by rendering only in init
func (v *VectorGraphic) Update(cmd *commands.Commands) {
	v.Graphic.Update(cmd)
	v.image.Clear()
	v.drawFunc(v.image)

	gPosX, gPosY := v.Transform2D.GetPos(false)
	camX, camY := v.GetCamera().WorldToCam(gPosX, gPosY, true)
	gAngle := v.Transform2D.GetAngle(false)
	gScaleX, gScaleY := v.Transform2D.GetScale(false)

	cmd.Renderer.Request(opgen.PosRotScale(
		v.image,
		camX, camY,
		gAngle,
		gScaleX, gScaleY,
		v.pivotX, v.pivotY,
	), v.image, v.target, v.drawOrder)
}

func NewDefaultVectorGraphic(graphic *graphic.Graphic) *VectorGraphic {
	return &VectorGraphic{
		Graphic:  graphic,
		drawFunc: func(i *ebiten.Image) { vector64.FillRect(i, 0, 0, 16, 16, color.RGBA{255, 0, 0, 255}, false) },
		image:    ebiten.NewImage(16, 16),
		target: renderer.RenderTarget{
			renderer.SCREEN,
			"Playerspace",
		},
		drawOrder: 0,
	}
}

func NewVectorGraphic(
	graphic *graphic.Graphic,
	options ...utils.Option[VectorGraphic],
) *VectorGraphic {
	vectorGraphic := NewDefaultVectorGraphic(graphic)

	for _, option := range options {
		option(vectorGraphic)
	}

	return vectorGraphic
}

func WithDrawFunc(drawFunc func(*ebiten.Image)) utils.Option[VectorGraphic] {
	return func(vg *VectorGraphic) {
		vg.drawFunc = drawFunc
	}
}

func WithImage(width, height int) utils.Option[VectorGraphic] {
	return func(vg *VectorGraphic) {
		vg.image = ebiten.NewImage(width, height)
	}
}

func WithDrawOrder(drawOrder int) utils.Option[VectorGraphic] {
	return func(vg *VectorGraphic) {
		vg.drawOrder = drawOrder
	}
}

func WithTarget(target renderer.RenderTarget) utils.Option[VectorGraphic] {
	return func(vg *VectorGraphic) {
		vg.target = target
	}
}

func WithPivot(x, y float64) utils.Option[VectorGraphic] {
	return func(vg *VectorGraphic) {
		vg.pivotX = x
		vg.pivotY = y
	}
}
