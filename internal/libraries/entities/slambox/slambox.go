package slambox

import (
	"image"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/concurrency"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/assettypes"
	"mask_of_the_tomb/internal/libraries/movebox"
	"mask_of_the_tomb/internal/libraries/physics"
	"math"
	"time"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	slamDelay                  = time.Millisecond * 500
	SlamboxConnectionFieldName = "ConnectedBoxes"
)

type slamboxState int

const (
	idle = iota
	waiting
	slamming
)

const (
	moveSpeed                  = 10.0
	tilesRegionX, tilesRegionY = 0, 0
	tileSize                   = 8
)

type TilePresence int

const (
	FREE TilePresence = iota
	NONE
	EXIST
)

type pattern [8]TilePresence
type tileData [8]bool

var (
	TopLeftInner = image.Rect(
		int(tilesRegionX+tileSize*3),
		int(tilesRegionY),
		int(tilesRegionX+tileSize*4),
		int(tilesRegionY+tileSize),
	)
	TopRightInner = image.Rect(
		int(tilesRegionX+tileSize*4),
		int(tilesRegionY),
		int(tilesRegionX+tileSize*5),
		int(tilesRegionY+tileSize),
	)
	BottomLeftInner = image.Rect(
		int(tilesRegionX+tileSize*3),
		int(tilesRegionY+tileSize),
		int(tilesRegionX+tileSize*4),
		int(tilesRegionY+tileSize*2),
	)
	BottomRightInner = image.Rect(
		int(tilesRegionX+tileSize*4),
		int(tilesRegionY+tileSize),
		int(tilesRegionX+tileSize*5),
		int(tilesRegionY+tileSize*2),
	)
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

	TopLeftInnerPattern = pattern{
		NONE, EXIST, FREE,
		EXIST /**/, FREE,
		FREE, FREE, EXIST}
	TopRightInnerPattern = pattern{
		FREE, EXIST, NONE,
		FREE /**/, EXIST,
		EXIST, FREE, FREE}
	BottomLeftInnerPattern = pattern{
		FREE, FREE, EXIST,
		EXIST /**/, FREE,
		NONE, EXIST, FREE}
	BottomRightInnerPattern = pattern{
		EXIST, FREE, FREE,
		FREE /**/, EXIST,
		FREE, EXIST, NONE}
	TopLeftPattern = pattern{
		NONE, NONE, FREE,
		NONE /**/, FREE,
		FREE, FREE, FREE}
	TopRightPattern = pattern{
		FREE, NONE, NONE,
		FREE /**/, NONE,
		FREE, FREE, FREE}
	BottomLeftPattern = pattern{
		FREE, FREE, FREE,
		NONE /**/, FREE,
		NONE, NONE, FREE}
	BottomRightPattern = pattern{
		FREE, FREE, FREE,
		FREE /**/, NONE,
		FREE, NONE, NONE}
	LeftPattern = pattern{
		FREE, FREE, FREE,
		NONE /**/, EXIST,
		FREE, FREE, FREE}
	RightPattern = pattern{
		FREE, FREE, FREE,
		EXIST /**/, NONE,
		FREE, FREE, FREE}
	TopPattern = pattern{
		FREE, NONE, FREE,
		FREE /**/, FREE,
		FREE, EXIST, FREE}
	BottomPattern = pattern{
		FREE, EXIST, FREE,
		FREE /**/, FREE,
		FREE, NONE, FREE}
)

type SlamContext struct {
	direction             maths.Direction
	tilemapCollider       *physics.TilemapCollider
	disconnectedColliders []*physics.RectCollider
}

type Slambox struct {
	Collider                  physics.RectCollider
	ConnectedBoxes            []*Slambox
	LinkID                    string   // ID to check for linked boxes
	OtherLinkIDs              []string // ID to check for linked boxes
	tilemap                   *ebiten.Image
	sprite                    *ebiten.Image
	movebox                   *movebox.Movebox
	state                     slamboxState
	moveFinishedEventListener *events.EventListener
	slamTimer                 *time.Timer
	slamTimerEventListener    *events.EventListener
	currentSlamCtx            SlamContext
}

func (s *Slambox) Update() {
	s.movebox.Update()
	x, y := s.movebox.GetPos()
	s.Collider.SetPos(x, y)
	switch s.state {
	case idle:
	case waiting:
		if _, done := concurrency.Poll(s.slamTimer.C); done {
			s.Slam(s.currentSlamCtx)
			s.state = slamming
		}
	case slamming:
		_, finished := s.moveFinishedEventListener.Poll()
		if finished {
			s.state = idle
		}
	}
}

func (s *Slambox) Draw(drawCtx rendering.Ctx) {
	x, y := s.movebox.GetPos()
	ebitenrenderutil.DrawAt(s.sprite, drawCtx.Dst, x, y)
}

// Projects a slambox through the environment given by slamctx
func (s *Slambox) Slam(slamCtx SlamContext) {
	projectedSlamboxRect, dist := slamCtx.tilemapCollider.ProjectRect(
		&s.Collider.Rect,
		slamCtx.direction,
		slamCtx.disconnectedColliders,
	)
	shortestDist := dist

	for _, otherSlambox := range s.ConnectedBoxes {
		_, otherDist := slamCtx.tilemapCollider.ProjectRect(
			&otherSlambox.GetCollider().Rect,
			slamCtx.direction,
			slamCtx.disconnectedColliders,
		)

		if math.Abs(otherDist) < math.Abs(dist) {
			shortestDist = otherDist
		}
	}

	for _, otherSlambox := range s.ConnectedBoxes {
		otherProjRect, _dist := slamCtx.tilemapCollider.ProjectRect(
			&otherSlambox.GetCollider().Rect,
			slamCtx.direction,
			slamCtx.disconnectedColliders,
		)

		offset := _dist - shortestDist

		switch slamCtx.direction {
		case maths.DirUp:
			otherProjRect.SetPos(otherSlambox.Collider.Left(), otherProjRect.Top()+offset)
		case maths.DirDown:
			otherProjRect.SetPos(otherSlambox.Collider.Left(), otherProjRect.Top()-offset)
		case maths.DirRight:
			otherProjRect.SetPos(otherProjRect.Left()-offset, otherSlambox.Collider.Top())
		case maths.DirLeft:
			otherProjRect.SetPos(otherProjRect.Left()+offset, otherSlambox.Collider.Top())
		}
		otherSlambox.SetTarget(otherProjRect.Left(), otherProjRect.Top())
		// TODO: set position of any connected components
	}

	offset := math.Abs(dist - shortestDist)

	switch slamCtx.direction {
	case maths.DirUp:
		projectedSlamboxRect.SetPos(s.Collider.Left(), projectedSlamboxRect.Top()+offset)
	case maths.DirDown:
		projectedSlamboxRect.SetPos(s.Collider.Left(), projectedSlamboxRect.Top()-offset)
	case maths.DirRight:
		projectedSlamboxRect.SetPos(projectedSlamboxRect.Left()-offset, s.Collider.Top())
	case maths.DirLeft:
		projectedSlamboxRect.SetPos(projectedSlamboxRect.Left()+offset, s.Collider.Top())
	}
	// TODO: set position of any connected components

	s.SetTarget(projectedSlamboxRect.Left(), projectedSlamboxRect.Top())
}

func (s *Slambox) StartSlam(direction maths.Direction, tilemapCollider *physics.TilemapCollider, disconnectedColliders []*physics.RectCollider) {
	s.slamTimer = time.NewTimer(slamDelay)
	s.state = waiting
	s.currentSlamCtx = SlamContext{
		direction:             direction,
		tilemapCollider:       tilemapCollider,
		disconnectedColliders: disconnectedColliders,
	}
}

func (s *Slambox) SetTarget(x, y float64) {
	s.movebox.SetTarget(x, y)
}

func (s *Slambox) GetCollider() *physics.RectCollider {
	return &s.Collider
}

func (s *Slambox) SetPos(x, y float64) {
	s.movebox.SetPos(x, y)
}

func matchPattern(in tileData, comp pattern) bool {
	for i := 0; i < 8; i++ {
		if comp[i] == FREE {
			continue
		} else if comp[i] == NONE && in[i] {
			return false
		} else if comp[i] == EXIST && !in[i] {
			return false
		}
	}
	return true
}

func (s *Slambox) CreateSprite() {
	slamboxTilemap := assetloader.GetAsset("slamboxTilemap").(*assettypes.ImageAsset).Image

	check := func(i, j int) {
		localX, localY := float64(i*tileSize), float64(j*tileSize)
		worldX := localX + s.Collider.Left() + tileSize/2
		worldY := localY + s.Collider.Top() + tileSize/2

		var ul, um, ur, ml, mr, bl, bm, br bool
		for _, otherBox := range append(s.ConnectedBoxes, s) {
			ul = ul || otherBox.Collider.Rect.IsWithin(worldX-tileSize, worldY-tileSize)
			um = um || otherBox.Collider.Rect.IsWithin(worldX, worldY-tileSize)
			ur = ur || otherBox.Collider.Rect.IsWithin(worldX+tileSize, worldY-tileSize)
			ml = ml || otherBox.Collider.Rect.IsWithin(worldX-tileSize, worldY)
			mr = mr || otherBox.Collider.Rect.IsWithin(worldX+tileSize, worldY)
			bl = bl || otherBox.Collider.Rect.IsWithin(worldX-tileSize, worldY+tileSize)
			bm = bm || otherBox.Collider.Rect.IsWithin(worldX, worldY+tileSize)
			br = br || otherBox.Collider.Rect.IsWithin(worldX+tileSize, worldY+tileSize)
		}

		tileData := tileData{
			ul, um, ur, ml, mr, bl, bm, br,
		}

		// fmt.Println(slamboxTilemap)
		if matchPattern(tileData, TopLeftInnerPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(TopLeftInner).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, TopRightInnerPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(TopRightInner).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, BottomLeftInnerPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(BottomLeftInner).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, BottomRightInnerPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(BottomRightInner).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, TopLeftPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(TopLeft).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, TopRightPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(TopRight).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, BottomLeftPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(BottomLeft).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, BottomRightPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(BottomRight).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, LeftPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(Left).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, RightPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(Right).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, TopPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(Top).(*ebiten.Image), s.sprite, localX, localY)
		} else if matchPattern(tileData, BottomPattern) {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(Bottom).(*ebiten.Image), s.sprite, localX, localY)
		} else {
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(Center).(*ebiten.Image), s.sprite, localX, localY)
		}
	}

	for i := 0; i < int(s.Collider.Width()/tileSize); i++ {
		check(i, 0)
		check(i, int(s.Collider.Height()/tileSize-1))
	}

	for i := 0; i < int(s.Collider.Height()/tileSize); i++ {
		check(0, i)
		check(int(s.Collider.Width()/tileSize-1), i)
	}

	for i := 1; i < int(s.Collider.Width()/tileSize-1); i++ {
		for j := 1; j < int(s.Collider.Height()/tileSize-1); j++ {
			localX, localY := float64(i*tileSize), float64(j*tileSize)
			ebitenrenderutil.DrawAt(slamboxTilemap.SubImage(Center).(*ebiten.Image), s.sprite, localX, localY)
		}
	}
}

func NewSlambox(
	entity *ebitenLDTK.Entity,
) *Slambox {
	newSlambox := Slambox{}
	newSlambox.Collider = physics.NewRectCollider(*maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	))
	newSlambox.ConnectedBoxes = make([]*Slambox, 0)
	newSlambox.LinkID = entity.Iid
	newSlambox.movebox = movebox.NewMovebox(moveSpeed)
	newSlambox.SetPos(entity.Px[0], entity.Px[1])
	newSlambox.moveFinishedEventListener = events.NewEventListener(newSlambox.movebox.OnMoveFinished)
	// newSlambox.tilemap = assettypes.NewImageAsset(slamboxTilemapPath)

	newSlambox.sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))
	connectionField := errs.Must(entity.GetFieldByName(SlamboxConnectionFieldName))
	for _, entityRef := range connectionField.EntityRefArray {
		newSlambox.OtherLinkIDs = append(newSlambox.OtherLinkIDs, entityRef.EntityIid)
	}

	return &newSlambox
}
