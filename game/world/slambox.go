package world

import (
	"image/color"
	"mask_of_the_tomb/ebitenLDTK"
	ebitenrenderutil "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/physics"
	"mask_of_the_tomb/rendering"
	"mask_of_the_tomb/utils/rect"

	"github.com/hajimehoshi/ebiten/v2"
)

type slambox struct {
	collider physics.RectCollider
	sprite   *ebiten.Image
}

func (s *slambox) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	ebitenrenderutil.DrawAt(s.sprite, rendering.RenderLayers.Playerspace, s.collider.Rect.Left()-camX, s.collider.Rect.Top()-camY)
}

func newSlambox(
	entity *ebitenLDTK.Entity,
) slambox {
	newSlambox := slambox{}
	newSlambox.collider = physics.NewRectCollider(*rect.FromEntity(entity))
	newSlambox.sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))
	newSlambox.sprite.Fill(color.RGBA{123, 245, 167, 255})

	return newSlambox
}
