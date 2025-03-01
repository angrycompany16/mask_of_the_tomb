package world

import (
	"image/color"
	"mask_of_the_tomb/ebitenLDTK"
	ebitenrenderutil "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/physics"
	"mask_of_the_tomb/rendering"
	"mask_of_the_tomb/utils"
	"mask_of_the_tomb/utils/rect"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: make a generalized LevelEntity or something to make this more logical
// TODO: make the moving box into a generalized entity?

const (
	moveSpeed = 10.0
)

type slambox struct {
	collider               physics.RectCollider
	sprite                 *ebiten.Image
	targetRect             rect.Rect
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
}

func (s *slambox) Update() {
	s.collider = physics.NewRectCollider(s.targetRect)
	s.posX += moveSpeed * s.moveDirX
	s.posY += moveSpeed * s.moveDirY

	if s.moveDirX < 0 {
		s.posX = utils.Clamp(s.posX, s.targetPosX, s.posX)
	} else if s.moveDirX > 0 {
		s.posX = utils.Clamp(s.posX, s.posX, s.targetPosX)
	}
	if s.moveDirY < 0 {
		s.posY = utils.Clamp(s.posY, s.targetPosY, s.posY)
	} else if s.moveDirY > 0 {
		s.posY = utils.Clamp(s.posY, s.posY, s.targetPosY)
	}

	if s.posX == s.targetPosX {
		s.moveDirX = 0
	}
	if s.posY == s.targetPosY {
		s.moveDirY = 0
	}

	s.collider.SetPos(s.posX, s.posY)
}

func (s *slambox) SetTarget(x, y float64) {
	s.targetPosX = x
	s.targetPosY = y
	s.moveDirX = math.Copysign(1, s.targetPosX-s.posX)
	s.moveDirY = math.Copysign(1, s.targetPosY-s.posY)
}

func (s *slambox) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	ebitenrenderutil.DrawAt(s.sprite, rendering.RenderLayers.Playerspace, s.posX-camX, s.posY-camY)
}

func (s *slambox) GetCollider() *physics.RectCollider {
	return &s.collider
}

func (s *slambox) SetPos(x, y float64) {
	s.posX, s.posY = x, y
	s.targetPosX, s.targetPosY = x, y
}

func newSlambox(
	entity *ebitenLDTK.Entity,
) *slambox {
	newSlambox := slambox{}
	newSlambox.collider = physics.NewRectCollider(*rect.FromEntity(entity))
	newSlambox.SetPos(entity.Px[0], entity.Px[1])
	newSlambox.targetRect = newSlambox.collider.Rect
	newSlambox.sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))
	newSlambox.sprite.Fill(color.RGBA{123, 245, 167, 255})

	return &newSlambox
}
