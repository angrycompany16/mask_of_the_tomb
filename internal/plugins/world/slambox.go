package world

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/autotile"
	"mask_of_the_tomb/internal/libraries/slambox"
	"math"
	"time"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

const moveSpeed = 10.0
const tileSize = 8.0
const slamDelay = time.Millisecond * 500

const (
	isGroupFieldName = "IsGroup"
	inChainFieldName = "InChain"
)

type slamCache struct {
	id  int
	dir maths.Direction
}

type SlamboxEntity struct {
	backendSlambox *slambox.Slambox
	sprite         *ebiten.Image
	inChain        bool
	isGroup        bool
	backendID      int
	rect           *maths.Rect
	slamCache      slamCache
	state          slamboxState
	slamTimer      *time.Timer
}

type slamboxState int

const (
	idle = iota
	waiting
	slamming
)

func (s *SlamboxEntity) Update() {
	switch s.state {
	case idle:
	case waiting:
		if _, done := threads.Poll(s.slamTimer.C); done {
			s.backendSlambox.RequestSlam(s.slamCache.dir)
			s.state = slamming
		}
	case slamming:
	}
}

func (s *SlamboxEntity) Draw(ctx rendering.Ctx) {
	x, y := s.rect.TopLeft()
	ebitenrenderutil.DrawAt(s.sprite, ctx.Dst, x, y)
	// s.particleSys.Draw(ctx)
}

func (s *SlamboxEntity) StartSlam(id int, dir maths.Direction) {
	s.slamTimer = time.NewTimer(slamDelay)
	s.state = waiting
	s.slamCache = slamCache{
		id:  id,
		dir: dir,
	}
}

func (s *SlamboxEntity) CreateSprite(slamboxTilemap *ebiten.Image) {
	autotile.CreateSprite(
		slamboxTilemap,
		s.sprite,
		autotile.GetDefaultTileRectData(0, 0, tileSize),
		autotile.GetDefaultTileRuleset(),
		tileSize,
		s.rect,
		autotile.WALL,
	)

	gemSprite := errs.Must(assettypes.GetImageAsset("slamboxGemRed"))
	x, y := s.GetGemPos()
	ebitenrenderutil.DrawAt(gemSprite, s.sprite, x, y, 0.5, 0.5)
}

func (s *SlamboxEntity) GetGemPos() (float64, float64) {
	sX, sY := s.rect.HalfSize()
	return sX, math.Min(sY, 8)
}

func NewSlambox(
	entity *ebitenLDTK.Entity,
	slamboxEnvironment *slambox.SlamboxEnvironment,
	levelLDTK *ebitenLDTK.Level,
) *SlamboxEntity {
	newSlambox := SlamboxEntity{}
	newBackendSlambox := slambox.NewSlambox(
		maths.NewRect(entity.Px[0], entity.Px[1], entity.Width, entity.Height), moveSpeed,
	)
	newSlambox.backendSlambox = newBackendSlambox

	isGroup := errs.Must(entity.GetFieldByName(isGroupFieldName)).Bool
	inChain := errs.Must(entity.GetFieldByName(inChainFieldName)).Bool

	var backendID int
	if inChain {
		// nothing
	} else {
		if isGroup {
			// nothing
		} else {
			backendID = slamboxEnvironment.AddSlambox(newBackendSlambox)
		}
	}

	newSlambox.backendID = backendID
	newSlambox.rect = newBackendSlambox.GetRect()
	newSlambox.sprite = ebiten.NewImage(int(newSlambox.rect.Width()), int(newSlambox.rect.Height()))

	return &newSlambox
}
