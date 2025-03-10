package world

import (
	"image"
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
	moveSpeed                  = 10.0
	tilesRegionX, tilesRegionY = 0, 0
	tileSize                   = 8
)

var (
	TopLeft = image.Rect(
		int(tilesRegionX),
		int(tilesRegionY),
		int(tilesRegionX+tileSize),
		int(tilesRegionY+tileSize),
	)
	Top = image.Rect(
		int(tilesRegionX+tileSize),
		int(tilesRegionY),
		int(tilesRegionX+2*tileSize),
		int(tilesRegionY+tileSize),
	)
	TopRight = image.Rect(
		int(tilesRegionX+2*tileSize),
		int(tilesRegionY),
		int(tilesRegionX+3*tileSize),
		int(tilesRegionY+tileSize),
	)
	Left = image.Rect(
		int(tilesRegionX),
		int(tilesRegionY+tileSize),
		int(tilesRegionX+tileSize),
		int(tilesRegionY+2*tileSize),
	)
	Center = image.Rect(
		int(tilesRegionX+tileSize),
		int(tilesRegionY+tileSize),
		int(tilesRegionX+2*tileSize),
		int(tilesRegionY+2*tileSize),
	)
	Right = image.Rect(
		int(tilesRegionX+2*tileSize),
		int(tilesRegionY+tileSize),
		int(tilesRegionX+3*tileSize),
		int(tilesRegionY+2*tileSize),
	)
	BottomLeft = image.Rect(
		int(tilesRegionX),
		int(tilesRegionY+2*tileSize),
		int(tilesRegionX+tileSize),
		int(tilesRegionY+3*tileSize),
	)
	Bottom = image.Rect(
		int(tilesRegionX+tileSize),
		int(tilesRegionY+2*tileSize),
		int(tilesRegionX+2*tileSize),
		int(tilesRegionY+3*tileSize),
	)
	BottomRight = image.Rect(
		int(tilesRegionX+2*tileSize),
		int(tilesRegionY+2*tileSize),
		int(tilesRegionX+3*tileSize),
		int(tilesRegionY+3*tileSize),
	)
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

	tilemap := errs.MustNewImageFromFile(SlamboxTilemapPath)

	newSlambox.sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))

	// Render from tilemap
	// Draw corners
	ebitenrenderutil.DrawAt(tilemap.SubImage(TopLeft).(*ebiten.Image), newSlambox.sprite, 0, 0)
	ebitenrenderutil.DrawAt(tilemap.SubImage(TopRight).(*ebiten.Image), newSlambox.sprite, entity.Width-tileSize, 0)
	ebitenrenderutil.DrawAt(tilemap.SubImage(BottomLeft).(*ebiten.Image), newSlambox.sprite, 0, entity.Height-tileSize)
	ebitenrenderutil.DrawAt(tilemap.SubImage(BottomRight).(*ebiten.Image), newSlambox.sprite, entity.Width-tileSize, entity.Height-tileSize)

	// Draw edges
	for i := 1; i < int(entity.Width/tileSize)-1; i++ {
		ebitenrenderutil.DrawAt(tilemap.SubImage(Top).(*ebiten.Image), newSlambox.sprite, float64(i*tileSize), 0)
		ebitenrenderutil.DrawAt(tilemap.SubImage(Bottom).(*ebiten.Image), newSlambox.sprite, float64(i*tileSize), entity.Height-tileSize)
	}

	for i := 1; i < int(entity.Height/tileSize)-1; i++ {
		ebitenrenderutil.DrawAt(tilemap.SubImage(Left).(*ebiten.Image), newSlambox.sprite, 0, float64(i*tileSize))
		ebitenrenderutil.DrawAt(tilemap.SubImage(Right).(*ebiten.Image), newSlambox.sprite, entity.Width-tileSize, float64(i*tileSize))
	}

	// Draw interior
	for i := 1; i < int(entity.Width/tileSize)-1; i++ {
		for j := 1; j < int(entity.Height/tileSize)-1; j++ {
			ebitenrenderutil.DrawAt(tilemap.SubImage(Center).(*ebiten.Image), newSlambox.sprite, float64(i*tileSize), float64(j*tileSize))
		}
	}

	connectionField := errs.Must(entity.GetFieldByName(SlamboxConnectionFieldName))
	for _, entityRef := range connectionField.EntityRefArray {
		newSlambox.otherLinkIDs = append(newSlambox.otherLinkIDs, entityRef.EntityIid)
	}

	return &newSlambox
}
