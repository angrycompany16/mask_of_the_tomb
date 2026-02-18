package world

import (
	"image/color"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/shaders"
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/core/threads"
	"mask_of_the_tomb/internal/libraries/autotile"
	"mask_of_the_tomb/internal/libraries/camera"
	"mask_of_the_tomb/internal/libraries/particles"
	"mask_of_the_tomb/internal/libraries/slambox"
	"math"
	"time"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const moveSpeed = 10.0
const tileSize = 8.0
const slamDelay = time.Millisecond * 500

const (
	isGroupFieldName      = "IsGroup"
	inChainFieldName      = "InChain"
	subSlamboxEntityName  = "SubSlambox"
	subSlamboxesFieldName = "SubSlamboxes"
)

type slamCache struct {
	id  int
	dir maths.Direction
}

type SlamboxEntity struct {
	slambox                   *slambox.Slambox
	slamboxGroup              *slambox.SlamboxGroup
	backendID                 int
	sprite                    *ebiten.Image
	inChain                   bool
	rect                      *maths.Rect
	slamCache                 slamCache
	state                     slamboxState
	slamTimer                 *time.Timer
	subSlamboxRects           []*maths.Rect
	subSlamboxSprites         []*ebiten.Image
	particleSys               *particles.ParticleSystem
	Light                     *shaders.Light
	landSound                 *sound.EffectPlayer
	moveFinishedEventListener *events.EventListener
}

type slamboxState int

const (
	idle = iota
	waiting
	slamming
)

func (s *SlamboxEntity) Update() {
	s.particleSys.Update()
	gemX, gemY := s.GetGemPos()
	s.Light.X = s.rect.Left() + gemX
	s.Light.Y = s.rect.Top() + gemY

	switch s.state {
	case idle:
	case waiting:
		if _, done := threads.Poll(s.slamTimer.C); done {
			s.slamboxGroup.RequestSlam(s.slamCache.dir)
			s.state = slamming
		}
	case slamming:
		if eventInfo, finished := s.moveFinishedEventListener.Poll(); finished {
			moveDir := eventInfo.Data.(maths.Direction)
			s.state = idle
			camera.Shake(0.4, 7, 1)
			// Also rotate into correct position
			s.PlayContactParticles(moveDir)
			s.landSound.Play()
		}
	}
}

func (s *SlamboxEntity) Draw(ctx rendering.Ctx) {
	x, y := s.rect.TopLeft()
	ebitenrenderutil.DrawAt(s.sprite, ctx.Dst, x, y)
	for i := range s.subSlamboxSprites {
		subX, subY := s.subSlamboxRects[i].TopLeft()
		ebitenrenderutil.DrawAt(s.subSlamboxSprites[i], ctx.Dst, subX, subY)
	}
	s.particleSys.Draw(ctx)
}

func (s *SlamboxEntity) DebugDraw(ctx rendering.Ctx) {
	vector.StrokeRect(ctx.Dst, float32(s.rect.Left()), float32(s.rect.Top()), float32(s.rect.Width()), float32(s.rect.Height()), 1.0, color.RGBA{255, 0, 0, 255}, false)
	for _, rect := range s.subSlamboxRects {
		vector.StrokeRect(ctx.Dst, float32(rect.Left()), float32(rect.Top()), float32(rect.Width()), float32(rect.Height()), 1.0, color.RGBA{255, 0, 0, 255}, false)
	}
}

func (s *SlamboxEntity) StartSlam(id int, dir maths.Direction) {
	s.slamTimer = time.NewTimer(slamDelay)
	s.state = waiting
	s.slamCache = slamCache{
		id:  id,
		dir: dir,
	}
}

func (s *SlamboxEntity) createSprite(slamboxTilemap *ebiten.Image) {
	autotile.CreateSprite(
		slamboxTilemap,
		s.sprite,
		autotile.GetDefaultTileRectData(0, 0, tileSize),
		autotile.GetDefaultTileRuleset(),
		tileSize,
		s.rect,
		autotile.WALL,
		autotile.RectList{
			List: s.slamboxGroup.GetSlamboxRects(),
			Kind: autotile.WALL,
		},
	)

	for i := range s.subSlamboxRects {
		autotile.CreateSprite(
			slamboxTilemap,
			s.subSlamboxSprites[i],
			autotile.GetDefaultTileRectData(0, 0, tileSize),
			autotile.GetDefaultTileRuleset(),
			tileSize,
			s.subSlamboxRects[i],
			autotile.WALL,
			autotile.RectList{
				List: s.slamboxGroup.GetSlamboxRects(),
				Kind: autotile.WALL,
			},
		)
	}

	gemSprite := errs.Must(assettypes.GetImageAsset("slamboxGemRed"))
	x, y := s.GetGemPos()
	ebitenrenderutil.DrawAt(gemSprite, s.sprite, x, y, 0.5, 0.5)
}

func (s *SlamboxEntity) GetGemPos() (float64, float64) {
	sX, sY := s.rect.HalfSize()
	return sX, math.Min(sY, 8)
}

// It is conceivable that this should be refactored
// Also this needs to take SubSlamboxes into account
func (s *SlamboxEntity) PlayContactParticles(moveDir maths.Direction) {
	w2, h2 := s.rect.HalfSize()
	spread := 20.0
	minSpeed := 2.0
	maxSpeed := 50.0
	switch moveDir {
	case maths.DirUp:
		s.particleSys.PosX = s.rect.Cx()
		s.particleSys.PosY = s.rect.Top()
		s.particleSys.SpawnPosY.Min = 0
		s.particleSys.SpawnPosY.Max = 0
		s.particleSys.SpawnPosX.Min = -w2
		s.particleSys.SpawnPosX.Max = w2
		s.particleSys.SpawnVelX.Min = -spread
		s.particleSys.SpawnVelX.Max = spread
		s.particleSys.SpawnVelY.Min = minSpeed
		s.particleSys.SpawnVelY.Max = maxSpeed
	case maths.DirDown:
		s.particleSys.PosX = s.rect.Cx()
		s.particleSys.PosY = s.rect.Bottom()
		s.particleSys.SpawnPosY.Min = 0
		s.particleSys.SpawnPosY.Max = 0
		s.particleSys.SpawnPosX.Min = -w2
		s.particleSys.SpawnPosX.Max = w2
		s.particleSys.SpawnVelX.Min = -spread
		s.particleSys.SpawnVelX.Max = spread
		s.particleSys.SpawnVelY.Max = -minSpeed
		s.particleSys.SpawnVelY.Min = -maxSpeed
	case maths.DirLeft:
		s.particleSys.PosX = s.rect.Left()
		s.particleSys.PosY = s.rect.Cy()
		s.particleSys.SpawnPosX.Min = 0
		s.particleSys.SpawnPosX.Max = 0
		s.particleSys.SpawnPosY.Min = -h2
		s.particleSys.SpawnPosY.Max = h2
		s.particleSys.SpawnVelY.Min = -spread
		s.particleSys.SpawnVelY.Max = spread
		s.particleSys.SpawnVelX.Min = minSpeed
		s.particleSys.SpawnVelX.Max = maxSpeed
	case maths.DirRight:
		s.particleSys.PosX = s.rect.Right()
		s.particleSys.PosY = s.rect.Cy()
		s.particleSys.SpawnPosX.Min = 0
		s.particleSys.SpawnPosX.Max = 0
		s.particleSys.SpawnPosY.Min = -h2
		s.particleSys.SpawnPosY.Max = h2
		s.particleSys.SpawnVelY.Min = -spread
		s.particleSys.SpawnVelY.Max = spread
		s.particleSys.SpawnVelX.Min = -minSpeed
		s.particleSys.SpawnVelX.Max = -maxSpeed
	case maths.DirNone:
	}
	s.particleSys.Play()
}

func NewSlamboxEntity(
	entity *ebitenLDTK.Entity,
	slamboxEnvironment *slambox.SlamboxEnvironment,
	levelLDTK *ebitenLDTK.Level,
) *SlamboxEntity {
	newSlamboxEntity := SlamboxEntity{}
	newSlamboxRect := maths.NewRect(entity.Px[0], entity.Px[1], entity.Width, entity.Height)
	newSlamboxEntity.rect = newSlamboxRect

	newSlambox := slambox.NewSlambox(newSlamboxRect, moveSpeed)

	// inChain := errs.Must(entity.GetFieldByName(inChainFieldName)).Bool

	subSlamboxRects := make([]*maths.Rect, 0)
	subSlamboxes := make([]*slambox.Slambox, 0)
	subSlamboxSprites := make([]*ebiten.Image, 0)
	subSlamboxesField := errs.Must(entity.GetFieldByName(subSlamboxesFieldName))
	for _, subSlamboxEntityRef := range subSlamboxesField.EntityRefArray {
		subSlamboxEntity := errs.Must(levelLDTK.GetEntityByIid(subSlamboxEntityRef.EntityIid))
		subSlamboxRect := maths.NewRect(subSlamboxEntity.Px[0], subSlamboxEntity.Px[1], subSlamboxEntity.Width, subSlamboxEntity.Height)
		subSlamboxRects = append(subSlamboxRects, subSlamboxRect)

		subSlambox := slambox.NewSlambox(subSlamboxRect, moveSpeed)
		subSlamboxes = append(subSlamboxes, subSlambox)

		subSlamboxSprite := ebiten.NewImage(int(subSlamboxRect.Width()), int(subSlamboxRect.Height()))
		subSlamboxSprites = append(subSlamboxSprites, subSlamboxSprite)
	}

	newSlamboxGroup := slambox.NewSlamboxGroup(
		append(subSlamboxes, newSlambox), len(subSlamboxes),
	)
	backendID := slamboxEnvironment.AddSlamboxGroup(newSlamboxGroup)

	gemX, gemY := newSlamboxEntity.GetGemPos()
	light := &shaders.Light{
		X:           newSlamboxRect.Left() + gemX,
		Y:           newSlamboxRect.Top() + gemY,
		InnerRadius: 0,
		OuterRadius: 50,
		ZOffset:     0,
		Intensity:   0.6,
		R:           1.0,
		G:           1.0,
		B:           1.0,
	}

	sprite := ebiten.NewImage(int(newSlamboxRect.Width()), int(newSlamboxRect.Height()))

	slamboxLandAudioStream := errs.Must(assettypes.GetWavStream("slamboxLandSound"))
	landSoundEffectPlayer := &sound.EffectPlayer{errs.Must(sound.FromStream(slamboxLandAudioStream)), 0.7}

	// NOTHING is a problem for the GOD og programming
	particleSys := *errs.Must(assettypes.GetYamlAsset("slamboxParticles")).(*particles.ParticleSystem)

	newSlamboxEntity.slamboxGroup = newSlamboxGroup
	newSlamboxEntity.subSlamboxRects = subSlamboxRects
	newSlamboxEntity.subSlamboxSprites = subSlamboxSprites
	newSlamboxEntity.backendID = backendID
	newSlamboxEntity.Light = light
	newSlamboxEntity.sprite = sprite
	newSlamboxEntity.landSound = landSoundEffectPlayer
	newSlamboxEntity.particleSys = &particleSys

	newSlamboxEntity.particleSys.Init()
	newSlamboxEntity.createSprite(errs.Must(assettypes.GetImageAsset("slamboxTilemap")))

	newSlamboxEntity.moveFinishedEventListener = events.NewEventListener(newSlamboxGroup.GetCenterSlambox().GetTracker().OnMoveFinished)

	return &newSlamboxEntity
}
