package entities

import (
	"image"
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/events"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	"mask_of_the_tomb/internal/game/core/timer"
	"mask_of_the_tomb/internal/game/physics"
	"mask_of_the_tomb/internal/game/physics/movebox"
	"mask_of_the_tomb/internal/maths"
	"math"
	"time"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
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
	sprite                    *ebiten.Image
	movebox                   *movebox.Movebox
	state                     slamboxState
	moveFinishedEventListener *events.EventListener
	slamTimer                 *timer.Timer
	slamTimerEventListener    *events.EventListener
	currentSlamCtx            SlamContext
}

func (s *Slambox) Update() {
	s.movebox.Update()
	x, y := s.movebox.GetPos()
	s.Collider.SetPos(x, y)
	s.slamTimer.Update()
	switch s.state {
	case idle:
	case waiting:
		_, raised := s.slamTimerEventListener.Poll()
		if raised {
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

func (s *Slambox) Draw() {
	x, y := s.movebox.GetPos()
	camX, camY := camera.GetPos()
	ebitenrenderutil.DrawAt(s.sprite, rendering.RenderLayers.Playerspace, x-camX, y-camY)
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

	s.SetTarget(projectedSlamboxRect.Left(), projectedSlamboxRect.Top())
}

func (s *Slambox) DoSlam(direction maths.Direction, tilemapCollider *physics.TilemapCollider, disconnectedColliders []*physics.RectCollider) {
	s.slamTimer.Reset()
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

func NewSlambox(
	entity *ebitenLDTK.Entity,
) *Slambox {
	newSlambox := Slambox{}
	newSlambox.Collider = physics.NewRectCollider(*maths.RectFromEntity(entity))
	newSlambox.ConnectedBoxes = make([]*Slambox, 0)
	newSlambox.LinkID = entity.Iid
	newSlambox.movebox = movebox.NewMovebox(moveSpeed)
	newSlambox.SetPos(entity.Px[0], entity.Px[1])
	newSlambox.moveFinishedEventListener = events.NewEventListener(newSlambox.movebox.OnMoveFinished)
	newSlambox.slamTimer = timer.NewTimer(time.Millisecond * 500)
	newSlambox.slamTimer.Pause()
	// newSlambox.SlamTimerFinished = events.NewEvent()
	newSlambox.slamTimerEventListener = events.NewEventListener(newSlambox.slamTimer.TimedoutEvent)

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
		newSlambox.OtherLinkIDs = append(newSlambox.OtherLinkIDs, entityRef.EntityIid)
	}

	return &newSlambox
}
