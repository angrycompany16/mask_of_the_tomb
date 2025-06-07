package world

import (
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/autotile"
	"mask_of_the_tomb/internal/libraries/entities/hazard"
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
	HazardConnectionFieldName  = "ConnectedComponents" // The name is a lie
)

type slamboxState int

const (
	idle = iota
	waiting
	slamming
)

const (
	moveSpeed = 10.0
	tileSize  = 8
)

type SlamContext struct {
	direction             maths.Direction
	tilemapCollider       *physics.TilemapCollider
	disconnectedColliders []*maths.Rect
}

type Slambox struct {
	Collider                  *maths.Rect
	ConnectedBoxes            []*Slambox
	LinkID                    string   // ID to check for linked boxes
	OtherLinkIDs              []string // ID to check for linked boxes
	sprite                    *ebiten.Image
	movebox                   *movebox.Movebox
	state                     slamboxState
	moveFinishedEventListener *events.EventListener
	slamTimer                 *time.Timer
	slamTimerEventListener    *events.EventListener
	currentSlamCtx            SlamContext
	attachedHazardIDs         []string
	attachedHazards           []*hazard.Hazard
}

func (s *Slambox) Update() {
	s.movebox.Update()
	x, y := s.movebox.GetPos()
	s.Collider.SetPos(x, y)
	for _, hazard := range s.attachedHazards {
		hazard.Hitbox.SetPos(x+hazard.PosOffsetX, y+hazard.PosOffsetY)
	}

	switch s.state {
	case idle:
	case waiting:
		if _, done := threads.Poll(s.slamTimer.C); done {
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
		s.Collider,
		slamCtx.direction,
		slamCtx.disconnectedColliders,
	)
	shortestDist := dist

	for _, otherSlambox := range s.ConnectedBoxes {
		_, otherDist := slamCtx.tilemapCollider.ProjectRect(
			otherSlambox.GetCollider(),
			slamCtx.direction,
			slamCtx.disconnectedColliders,
		)

		if math.Abs(otherDist) < math.Abs(dist) {
			shortestDist = otherDist
		}
	}

	for _, otherSlambox := range s.ConnectedBoxes {
		otherProjRect, _dist := slamCtx.tilemapCollider.ProjectRect(
			otherSlambox.GetCollider(),
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
	s.SetTarget(projectedSlamboxRect.Left(), projectedSlamboxRect.Top())
}

func (s *Slambox) StartSlam(direction maths.Direction, tilemapCollider *physics.TilemapCollider, disconnectedColliders []*maths.Rect) {
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

func (s *Slambox) GetCollider() *maths.Rect {
	return s.Collider
}

func (s *Slambox) SetPos(x, y float64) {
	s.movebox.SetPos(x, y)
}

func (s *Slambox) CreateSprite(slamboxTilemap *ebiten.Image) {
	connectedRects := make([]*maths.Rect, 0)
	for _, connectedBox := range s.ConnectedBoxes {
		connectedRects = append(connectedRects, connectedBox.Collider)
	}
	autotile.CreateSprite(
		slamboxTilemap,
		s.sprite,
		autotile.GetDefaultTileRectData(0, 0, tileSize),
		autotile.GetDefaultTileRuleset(),
		tileSize,
		s.Collider,
		autotile.WALL,
		autotile.RectList{
			List: connectedRects,
			Kind: autotile.WALL,
		},
	)

	for _, hazard := range s.attachedHazards {
		autotile.CreateSprite(
			slamboxTilemap,
			hazard.Sprite,
			autotile.GetDefaultTileRectData(0, 0, tileSize),
			autotile.GetDefaultSpikeRules(),
			tileSize,
			hazard.Hitbox,
			autotile.SPIKE,
			autotile.RectList{
				List: append(connectedRects, s.Collider),
				Kind: autotile.WALL,
			},
		)
	}
}

func NewSlambox(
	entity *ebitenLDTK.Entity,
) *Slambox {
	newSlambox := Slambox{}
	newSlambox.Collider = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)
	newSlambox.ConnectedBoxes = make([]*Slambox, 0)
	newSlambox.LinkID = entity.Iid
	newSlambox.movebox = movebox.NewMovebox(moveSpeed)
	newSlambox.SetPos(entity.Px[0], entity.Px[1])
	newSlambox.moveFinishedEventListener = events.NewEventListener(newSlambox.movebox.OnMoveFinished)
	// TODO: Size dynamically by hazards
	newSlambox.sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))

	connectionField := errs.Must(entity.GetFieldByName(SlamboxConnectionFieldName))
	for _, entityRef := range connectionField.EntityRefArray {
		newSlambox.OtherLinkIDs = append(newSlambox.OtherLinkIDs, entityRef.EntityIid)
	}

	hazardField := errs.Must(entity.GetFieldByName(HazardConnectionFieldName))
	for _, entityRef := range hazardField.EntityRefArray {
		newSlambox.attachedHazardIDs = append(newSlambox.attachedHazardIDs, entityRef.EntityIid)
	}

	return &newSlambox
}
