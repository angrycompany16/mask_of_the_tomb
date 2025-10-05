package world

import (
	"fmt"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/autotile"
	"mask_of_the_tomb/internal/libraries/entities"
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
	chainNodes            []*entities.ChainNode
}

type Slambox struct {
	Collider                  *maths.Rect
	ConnectedBoxes            []*Slambox
	ChainedSlambox            *Slambox
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
	// TODO: Hazard should not really be in entities, instead there should be a general Hazard struct
	// that takes care of damage calculations and such, and then a specific hazard type
	attachedHazards []*entities.Hazard
	chainNodeID     string
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

func ProjectInChain(rect *maths.Rect, slamCtx *SlamContext) float64 {
	x, y := rect.Center()
	shortestDist := 0.0
	for i, chainNode := range slamCtx.chainNodes {
		hitNode, _, _ := chainNode.Rect.RaycastDirectional(x, y, slamCtx.direction)
		overlapsNode := chainNode.Rect.IsWithin(x, y)
		if hitNode {
			dist := maths.Length(chainNode.Rect.Cx()-x, chainNode.Rect.Cy()-y)
			shortestDist = dist
		} else if overlapsNode {
			shortestDist = 0
			for j, otherChainNode := range slamCtx.chainNodes {
				if i == j {
					continue
				}
				hitOtherNode, _, _ := otherChainNode.Rect.RaycastDirectional(x, y, slamCtx.direction)
				if hitOtherNode {
					dist := maths.Length(otherChainNode.Rect.Cx()-x, otherChainNode.Rect.Cy()-y)
					shortestDist = dist
				}
			}
		}
	}
	return shortestDist
}

// Three outcomes: Found no new node, found slambox, found new node
// No new node: dirnone + any bool
// Foudn new node: dir + false
// Foudn slambox: dir + true
func ExploreNode(dir maths.Direction, currRect *maths.Rect, chainNodes []*entities.ChainNode, chainedSlamboxRect *maths.Rect) (bool, maths.Direction) {
	fmt.Println("Exploring from position", currRect.Cx(), currRect.Cy())
	fmt.Println("Exploring direction", dir)

	for _, chainNode := range chainNodes {
		// Check if there is a chain node in the current direction
		hit, _, _ := chainNode.Rect.RaycastDirectional(currRect.Cx(), currRect.Cy(), dir)
		if hit {
			fmt.Println("Found a node")
			slamboxHit, _, _ := chainedSlamboxRect.RaycastDirectional(currRect.Cx(), currRect.Cy(), dir)
			if slamboxHit {
				fmt.Println("Found slambox with raycast", dir, currRect.Cx(), currRect.Cy())
				return true, dir
			}

			// Check if the chained slambox is in the current direction. If so, done.

			// Construct new directions to search in (all but where we came from)
			directions := make([]maths.Direction, 0)
			switch dir {
			case maths.DirUp:
				directions = []maths.Direction{maths.DirUp, maths.DirLeft, maths.DirRight}
			case maths.DirDown:
				directions = []maths.Direction{maths.DirDown, maths.DirLeft, maths.DirRight}
			case maths.DirRight:
				directions = []maths.Direction{maths.DirUp, maths.DirDown, maths.DirRight}
			case maths.DirLeft:
				directions = []maths.Direction{maths.DirUp, maths.DirDown, maths.DirLeft}
			}

			fmt.Println("I will search in directions:", directions)

			for _, direction := range directions {
				// Explore new node and return if we find chain slambox
				foundSlambox, newDirection := ExploreNode(direction, chainNode.Rect, chainNodes, chainedSlamboxRect)
				if foundSlambox {
					return true, newDirection
				}
			}
		}
	}

	fmt.Println("Found nothing")
	// If the loop terminates we did not find the slambox, so we return false + dirnone
	return false, maths.DirNone
}

// Explore in both directions and return the slambox direction depending on the direction
func ComputeChainedSlamboxDirection(startRect *maths.Rect, endRect *maths.Rect, dir maths.Direction, chainNodes []*entities.ChainNode) maths.Direction {
	for _, searchDirection := range []maths.Direction{maths.DirUp, maths.DirDown, maths.DirRight, maths.DirLeft} {
		foundSlambox, direction := ExploreNode(searchDirection, startRect, chainNodes, endRect)
		if foundSlambox {
			if searchDirection == dir {
				fmt.Println("Regular case")
				return direction
			} else {
				fmt.Println("Opposite case")
				return maths.Opposite(direction)
			}
		}
	}
	return dir
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

	otherSlamCtx := SlamContext{}
	if s.chainNodeID != "" {
		otherSlamCtx.direction = ComputeChainedSlamboxDirection(s.Collider, s.ChainedSlambox.Collider, slamCtx.direction, slamCtx.chainNodes)
		otherSlamCtx.tilemapCollider = slamCtx.tilemapCollider
		otherSlamCtx.disconnectedColliders = slamCtx.disconnectedColliders
		otherSlamCtx.chainNodes = slamCtx.chainNodes

		shortestDist = math.Min(shortestDist, ProjectInChain(s.Collider, &slamCtx))

		projectedChainDist := ProjectInChain(s.ChainedSlambox.Collider, &otherSlamCtx)
		CWChainDist := 0.0
		CCWChainDist := 0.0
		if projectedChainDist == 0 {
			fmt.Println("Projected chain dist 0")
			dir := otherSlamCtx.direction
			otherSlamCtx.direction = maths.RotateCW(dir)
			CWChainDist = ProjectInChain(s.ChainedSlambox.Collider, &otherSlamCtx)
			otherSlamCtx.direction = maths.RotateCCW(dir)
			CCWChainDist = ProjectInChain(s.ChainedSlambox.Collider, &otherSlamCtx)
			fmt.Println(CWChainDist)
			fmt.Println(CCWChainDist)
		}
		projectedChainDist = maths.Max(projectedChainDist, CWChainDist, CCWChainDist)

		fmt.Println(projectedChainDist)
		fmt.Println(shortestDist)
		shortestDist = math.Min(shortestDist, projectedChainDist)
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

	if s.chainNodeID != "" {
		chainedProjRect, _dist := otherSlamCtx.tilemapCollider.ProjectRect(
			s.ChainedSlambox.GetCollider(),
			otherSlamCtx.direction,
			slamCtx.disconnectedColliders,
		)

		offset := _dist - shortestDist

		fmt.Println(otherSlamCtx.direction)
		switch otherSlamCtx.direction {
		case maths.DirUp:
			chainedProjRect.SetPos(s.ChainedSlambox.Collider.Left(), chainedProjRect.Top()+offset)
		case maths.DirDown:
			chainedProjRect.SetPos(s.ChainedSlambox.Collider.Left(), chainedProjRect.Top()-offset)
		case maths.DirRight:
			chainedProjRect.SetPos(chainedProjRect.Left()-offset, s.ChainedSlambox.Collider.Top())
		case maths.DirLeft:
			chainedProjRect.SetPos(chainedProjRect.Left()+offset, s.ChainedSlambox.Collider.Top())
		}

		s.ChainedSlambox.SetTarget(chainedProjRect.Left(), chainedProjRect.Top())
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

func (s *Slambox) StartSlam(direction maths.Direction, tilemapCollider *physics.TilemapCollider, disconnectedColliders []*maths.Rect, chainNodes []*entities.ChainNode) {
	s.slamTimer = time.NewTimer(slamDelay)
	s.state = waiting
	s.currentSlamCtx = SlamContext{
		direction:             direction,
		tilemapCollider:       tilemapCollider,
		disconnectedColliders: disconnectedColliders,
		chainNodes:            chainNodes,
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

	chainNodeField := errs.Must(entity.GetFieldByName("ChainNode"))
	newSlambox.chainNodeID = chainNodeField.EntityRef.EntityIid

	return &newSlambox
}
