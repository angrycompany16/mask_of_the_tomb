package world

import (
	"image/color"
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/camera"
	"mask_of_the_tomb/internal/game/physics"
	"mask_of_the_tomb/internal/game/rendering"
	"mask_of_the_tomb/internal/maths"
	"math"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: make the moving box into a generalized entity?

const (
	moveSpeed = 10.0
)

type Slambox struct {
	Collider               physics.RectCollider
	ConnectedBoxes         []*Slambox
	LinkID                 string   // ID to check for linked boxes
	otherLinkIDs           []string // ID to check for linked boxes
	sprite                 *ebiten.Image
	targetRect             maths.Rect
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
}

func (s *Slambox) Update() {
	s.Collider = physics.NewRectCollider(s.targetRect)
	s.posX += moveSpeed * s.moveDirX
	s.posY += moveSpeed * s.moveDirY

	if s.moveDirX < 0 {
		s.posX = maths.Clamp(s.posX, s.targetPosX, s.posX)
	} else if s.moveDirX > 0 {
		s.posX = maths.Clamp(s.posX, s.posX, s.targetPosX)
	}
	if s.moveDirY < 0 {
		s.posY = maths.Clamp(s.posY, s.targetPosY, s.posY)
	} else if s.moveDirY > 0 {
		s.posY = maths.Clamp(s.posY, s.posY, s.targetPosY)
	}

	if s.posX == s.targetPosX {
		s.moveDirX = 0
	}
	if s.posY == s.targetPosY {
		s.moveDirY = 0
	}

	s.Collider.SetPos(s.posX, s.posY)
}

func (s *Slambox) SetTarget(x, y float64) {
	s.targetPosX = x
	s.targetPosY = y
	s.moveDirX = math.Copysign(1, s.targetPosX-s.posX)
	s.moveDirY = math.Copysign(1, s.targetPosY-s.posY)
}

func (s *Slambox) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	ebitenrenderutil.DrawAt(s.sprite, rendering.RenderLayers.Playerspace, s.posX-camX, s.posY-camY)
}

func (s *Slambox) GetCollider() *physics.RectCollider {
	return &s.Collider
}

func (s *Slambox) SetPos(x, y float64) {
	s.posX, s.posY = x, y
	s.targetPosX, s.targetPosY = x, y
}

func newSlambox(
	entity *ebitenLDTK.Entity,
) *Slambox {
	newSlambox := Slambox{}
	newSlambox.Collider = physics.NewRectCollider(*maths.RectFromEntity(entity))
	newSlambox.ConnectedBoxes = make([]*Slambox, 0)
	newSlambox.LinkID = entity.Iid
	newSlambox.SetPos(entity.Px[0], entity.Px[1])
	newSlambox.targetRect = newSlambox.Collider.Rect
	newSlambox.sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))
	newSlambox.sprite.Fill(color.RGBA{123, 245, 167, 255})

	connectionField := errs.Must(entity.GetFieldByName(SlamboxConnectionFieldName))
	for _, entityRef := range connectionField.EntityRefArray {
		newSlambox.otherLinkIDs = append(newSlambox.otherLinkIDs, entityRef.EntityIid)
	}

	return &newSlambox
}
