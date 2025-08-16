package player

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/libraries/animation"
	"mask_of_the_tomb/internal/libraries/assettypes"
	"mask_of_the_tomb/internal/libraries/inputbuffer"
	"mask_of_the_tomb/internal/libraries/movebox"
	"mask_of_the_tomb/internal/libraries/particles"
	"path/filepath"

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

var (
	playerSpritePath       = filepath.Join(assets.PlayerFolder, "player.png")
	jumpParticlesBroadPath = filepath.Join("assets", "particlesystems", "player", "jump-broad.yaml")
	jumpParticlesTightPath = filepath.Join("assets", "particlesystems", "player", "jump-tight.yaml")
)

// TODO: We want the player plugin to instead become a bundle of separate
// interactions with libraries. Possible separation
// - Player movement - Links input and movement
// - Then all the other stuff - Rendering and sound and stuff

// - Here's a performance criterion: It shouldn't be infeasible to create a version
//   of the player that doesn't interact with slamboxes, instead it should be as easy
//   as just removing that interaction, and that interaction should be very clear

// Consider implementing the different components as interfaces?

// TODO: Allow for sprites which aren't exactly 16x16
type Player struct {
	// Movecomponent
	// Graphicscomponent
	//
	State                     playerState
	movebox                   *movebox.Movebox
	hitbox                    *maths.Rect
	sprite                    *ebiten.Image
	jumpOffset, jumpOffsetvel float64
	direction                 maths.Direction
	animator                  *animation.Animator
	Disabled                  bool
	canPlaySlamSound          bool // Ugly :(
	InputBuffer               inputbuffer.InputBuffer
	deathAnim                 *DeathAnim
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
		InputBuffer: inputbuffer.NewInputBuffer(inputBufferDuration),
		State:       Idle,
		OnDeath:     events.NewEvent(),
		OnMove:      events.NewEvent(),
		deathAnim:   NewDeathAnim(),
	}

	player.moveFinishedListener = events.NewEventListener(player.movebox.OnMoveFinished)
	player.clipFinishedListener = events.NewEventListener(player.animator.OnClipFinished)
	return player
}

// ------ INIT ------
func (p *Player) CreateAssets() {
	assetloader.Load("playerSprite", assettypes.MakeImageAsset(assets.Player_sprite))
	assetloader.Load("dashSound", assettypes.MakeSoundAsset(assets.Dash_wav, assettypes.Wav))
	assetloader.Load("slamSound", assettypes.MakeSoundAsset(assets.Slam_wav, assettypes.Wav))
	assetloader.Load("deathSound", assettypes.MakeSoundAsset(assets.Death_mp3, assettypes.Mp3))
	assetloader.Load("jumpParticlesBroad", particles.NewParticleSystemAsset(jumpParticlesBroadPath, rendering.ScreenLayers.Playerspace))
	assetloader.Load("jumpParticlesTight", particles.NewParticleSystemAsset(jumpParticlesTightPath, rendering.ScreenLayers.Playerspace))
}

func (p *Player) Init(posX, posY float64, direction maths.Direction) {
	p.sprite = errs.Must(assettypes.GetImageAsset("playerSprite"))
	p.dashSound = errs.Must(assettypes.GetEffectPlayerAsset("dashSound"))
	p.slamSound = errs.Must(assettypes.GetEffectPlayerAsset("slamSound"))
	p.deathSound = errs.Must(assettypes.GetEffectPlayerAsset("deathSound"))
	p.jumpParticlesBroad = errs.Must(particles.GetParticleSystemAsset("jumpParticlesBroad"))
	p.jumpParticlesTight = errs.Must(particles.GetParticleSystemAsset("jumpParticlesTight"))

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

func (p *Player) SetHitboxPos(x, y float64) {
	p.movebox.SetPos(x, y)
	p.hitbox.SetPos(x, y)
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
