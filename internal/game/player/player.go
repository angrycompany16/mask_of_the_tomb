package player

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/game/animation"
	"mask_of_the_tomb/internal/game/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/game/core/events"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/physics/movebox"
	"mask_of_the_tomb/internal/game/physics/particles"
	"mask_of_the_tomb/internal/game/player/deathanim"
	"mask_of_the_tomb/internal/game/sound"
	"mask_of_the_tomb/internal/maths"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	moveSpeed             = 5.0
	defaultPlayerHealth   = 5.0
	invincibilityDuration = time.Second
	inputBufferDuration   = 0.1
)

// TODO: Allow for sprites which aren't exactly 16x16
type Player struct {
	State                     playerState
	movebox                   *movebox.Movebox
	hitbox                    *maths.Rect
	sprite                    *ebiten.Image
	jumpOffset, jumpOffsetvel float64
	direction                 maths.Direction
	animator                  *animation.Animator
	Disabled                  bool
	canPlaySlamSound          bool // Ugly :(
	InputBuffer               inputBuffer
	deathAnim                 *deathanim.DeathAnim
	dashSound                 *sound.EffectPlayer
	slamSound                 *sound.EffectPlayer
	deathSound                *sound.EffectPlayer
	jumpParticlesBroad        *particles.ParticleSystem
	jumpParticlesTight        *particles.ParticleSystem
	// Events
	OnDeath *events.Event
	OnMove  *events.Event
	// Listeners
	moveFinishedListener *events.EventListener
	clipFinishedListener *events.EventListener
}

// ------ CONSTRUCTOR ------
func NewPlayer() *Player {
	player := &Player{
		movebox:     movebox.NewMovebox(moveSpeed),
		animator:    animation.NewAnimator(playerAnimationMap),
		InputBuffer: newInputBuffer(inputBufferDuration),
		State:       Idle,
		OnDeath:     events.NewEvent(),
		OnMove:      events.NewEvent(),
		deathAnim:   deathanim.NewDeathAnim(),
		dashSound:   sound.NewEffectPlayer(assets.Dash_wav, sound.Wav),
		slamSound:   sound.NewEffectPlayer(assets.Slam_wav, sound.Wav),
		deathSound:  sound.NewEffectPlayer(assets.Death_mp3, sound.Mp3),
	}

	player.moveFinishedListener = events.NewEventListener(player.movebox.OnMoveFinished)
	player.clipFinishedListener = events.NewEventListener(player.animator.OnClipFinished)
	return player
}

// ------ INIT ------
func (p *Player) CreateAssets() {
	p.sprite = assettypes.NewImageAsset(playerSpritePath)
	p.jumpParticlesBroad = assettypes.NewParticleSystemAsset(jumpParticlesBroadPath, rendering.RenderLayers.Playerspace)
	p.jumpParticlesTight = assettypes.NewParticleSystemAsset(jumpParticlesTightPath, rendering.RenderLayers.Playerspace)
}

func (p *Player) Init(posX, posY float64, direction maths.Direction) {
	p.SetPos(posX, posY)
	p.direction = direction
	p.hitbox = maths.RectFromImage(posX, posY, p.sprite)
	p.animator.SwitchClip(idleAnim)
}

// ------ GETTERS ------
func (p *Player) GetHitbox() *maths.Rect {
	return p.hitbox
}

func (p *Player) GetLevelSwapInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (p *Player) GetPosCentered() (float64, float64) {
	s := p.sprite.Bounds().Size()
	x, y := p.movebox.GetPos()
	return x + float64(s.X)/2, y + float64(s.Y)/2
}

func (p *Player) GetSize() (float64, float64) {
	s := p.sprite.Bounds().Size()
	return float64(s.X), float64(s.Y)
}

func (p *Player) GetMovementSize() (float64, float64) {
	moveDirX, moveDirY := p.movebox.GetMovedir()
	return moveSpeed * moveDirX, moveSpeed * moveDirY
}

func (p *Player) CanMove() bool {
	return p.State == Idle
}

func (p *Player) IsMoving() bool {
	moveDirX, moveDirY := p.movebox.GetMovedir()
	return moveDirX != 0 || moveDirY != 0
}

// ------ SETTERS ------
func (p *Player) SetPos(x, y float64) {
	p.movebox.SetPos(x, y)
}

func (p *Player) Die() {
	p.Disabled = true
	p.State = Dying
	p.deathAnim.Play()
	p.deathAnim.SetPos(p.hitbox.Center())
	p.deathSound.Play()

	// May not be necessary
	p.OnDeath.Raise(events.EventInfo{})
}

func (p *Player) Respawn() {
	p.Disabled = false
	p.direction = maths.DirUp
	p.State = Idle
}

func (p *Player) Dash(direction maths.Direction, x, y float64) {
	p.direction = direction
	p.State = Moving

	p.dashSound.Play()
	p.animator.SwitchClip(dashInitAnim)
	p.movebox.SetTarget(x, y)
	p.playJumpParticles(direction)
}

func (p *Player) EnterSlamAnim() {
	p.animator.SwitchClip(slamAnim)
}

// ------ INTERNAL ------
func (p *Player) calculateJumpOffset() (float64, float64) {
	angle := maths.ToRadians(p.direction)
	if angle == 0 {
		return 0, p.jumpOffset
	} else if angle == math.Pi/2 {
		return -p.jumpOffset, 0
	} else if angle == math.Pi {
		return 0, -p.jumpOffset
	} else if angle == 3*math.Pi/2 {
		return p.jumpOffset, 0
	}
	return 0, 0
}

func (p *Player) readMoveInput() maths.Direction {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return maths.DirUp
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return maths.DirDown
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return maths.DirRight
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return maths.DirLeft
	}
	return maths.DirNone
}

func (p *Player) playJumpParticles(direction maths.Direction) {
	centerX, centerY := p.hitbox.Center()
	switch direction {
	case maths.DirUp:
		p.jumpParticlesBroad.PosX = centerX
		p.jumpParticlesBroad.PosY = p.hitbox.Bottom()
		p.jumpParticlesTight.PosX = centerX
		p.jumpParticlesTight.PosY = p.hitbox.Bottom()
		p.jumpParticlesBroad.Angle = 0
		p.jumpParticlesTight.Angle = 0
	case maths.DirDown:
		p.jumpParticlesBroad.PosX = centerX
		p.jumpParticlesBroad.PosY = p.hitbox.Top()
		p.jumpParticlesTight.PosX = centerX
		p.jumpParticlesTight.PosY = p.hitbox.Top()
		p.jumpParticlesBroad.Angle = math.Pi
		p.jumpParticlesTight.Angle = math.Pi
	case maths.DirRight:
		p.jumpParticlesBroad.PosX = p.hitbox.Left()
		p.jumpParticlesBroad.PosY = centerY
		p.jumpParticlesTight.PosX = p.hitbox.Left()
		p.jumpParticlesTight.PosY = centerY
		p.jumpParticlesBroad.Angle = math.Pi / 2
		p.jumpParticlesTight.Angle = math.Pi / 2
	case maths.DirLeft:
		p.jumpParticlesBroad.PosX = p.hitbox.Right()
		p.jumpParticlesBroad.PosY = centerY
		p.jumpParticlesTight.PosX = p.hitbox.Right()
		p.jumpParticlesTight.PosY = centerY
		p.jumpParticlesBroad.Angle = 3 * math.Pi / 2
		p.jumpParticlesTight.Angle = 3 * math.Pi / 2
	}

	p.jumpParticlesBroad.Play()
	// p.jumpParticlesTight.Play()
}
