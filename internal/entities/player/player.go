package player

// New task - gameplay / health system

import (
	"mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/animation"
	"mask_of_the_tomb/internal/game/camera"
	"mask_of_the_tomb/internal/game/entities"
	"mask_of_the_tomb/internal/game/events"
	"mask_of_the_tomb/internal/game/rendering"
	"mask_of_the_tomb/internal/maths"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO: Rewrite literally everything with some kind of gameobject interface
// Basically we can create a struct of this type, and then the update function
// Will simply be called automatically without us needing to reference these things
// All over creating a lot of chaos

// TODO: Rewrite as much as possible using events rather than importing packages

var (
	_player = Player{}
)

const (
	moveSpeed             = 5.0
	defaultPlayerHealth   = 5.0
	invincibilityDuration = time.Second
	inputBufferDuration   = 0.1
)

// TODO: Rewrite player to use maths.direction instead of float angle
// TODO: turn animations into asset file? (i.e. assets for anims)
type Player struct {
	PosX, PosY                float64
	targetPosX, targetPosY    float64
	prevPosX, prevPosY        float64
	moveDirX, moveDirY        float64
	moveProgress              float64
	jumpOffset, jumpOffsetvel float64
	Hitbox                    *maths.Rect
	sprite                    *ebiten.Image
	animator                  *animation.Animator
	Invincible                bool
	Disabled                  bool
	damageOverlay             deathEffect
	angle                     float64
	InputBuffer               inputBuffer
	State                     playerState
	finishedClipEventListener *events.EventListener
}

type playerInitMsg struct {
	Width, Height float64
}

func Init(posX, posY float64) playerInitMsg {
	entities.RegisterEntity(&_player, "Player")
	_player.SetPos(posX, posY)
	_player.Hitbox = maths.RectFromImage(posX, posY, _player.sprite)
	_player.animator.SwitchClip(idleAnim)

	width, height := _player.GetSize()
	return playerInitMsg{
		Width:  width,
		Height: height,
	}
}

func (p *Player) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	jumpOffsetX, jumpOffsetY := p.getJumpOffset()
	ebitenrenderutil.DrawAtRotated(
		p.animator.GetSprite(),
		rendering.RenderLayers.Playerspace,
		p.PosX-camX-jumpOffsetX,
		p.PosY-camY-jumpOffsetY,
		p.angle,
		0.5,
		0.5,
	)

	p.damageOverlay.Draw()
}

func (p *Player) getJumpOffset() (float64, float64) {
	if p.angle == 0 {
		return 0, p.jumpOffset
	} else if p.angle == math.Pi/2 {
		return -p.jumpOffset, 0
	} else if p.angle == math.Pi {
		return 0, -p.jumpOffset
	} else if p.angle == 3*math.Pi/2 {
		return p.jumpOffset, 0
	}
	return 0, 0
}

func (p *Player) GetLevelSwapInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (p *Player) getMoveInput() maths.Direction {
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

func (p *Player) SetPos(x, y float64) {
	p.PosX, p.PosY = x, y
	p.targetPosX, p.targetPosY = x, y
}

func (p *Player) SetRot(direction maths.Direction) {
	switch direction {
	case maths.DirUp:
		p.angle = 0
	case maths.DirDown:
		p.angle = math.Pi
	case maths.DirRight:
		p.angle = math.Pi / 2
	case maths.DirLeft:
		p.angle = 3 * math.Pi / 2
	}
}

func (p *Player) GetPosCentered() (float64, float64) {
	s := p.sprite.Bounds().Size()
	return p.PosX + float64(s.X)/2, p.PosY + float64(s.Y)/2
}

func (p *Player) SetTarget(x, y float64) {
	p.targetPosX = x
	p.targetPosY = y
	p.prevPosX = p.PosX
	p.prevPosY = p.PosY
	p.moveDirX = math.Copysign(1, p.targetPosX-p.PosX)
	p.moveDirY = math.Copysign(1, p.targetPosY-p.PosY)
}

func (p *Player) GetSize() (float64, float64) {
	s := p.sprite.Bounds().Size()
	return float64(s.X), float64(s.Y)
}

func (p *Player) GetMovementSize() (float64, float64) {
	return moveSpeed * p.moveDirX, moveSpeed * p.moveDirY
}

func (p *Player) CanMove() bool {
	return p.State == Idle
}

func (p *Player) IsMoving() bool {
	return p.moveDirX != 0 || p.moveDirY != 0
}

func (p *Player) TakeDamage(damage float64) {
	p.Invincible = true
	p.Disabled = true
	go p.hitHazard()
}

func (p *Player) hitHazard() {
	time.Sleep(time.Millisecond * 300)
	p.Disabled = false
	p.damageOverlay.alpha = 1.0
	go p.disableInvincibility()
}

func (p *Player) disableInvincibility() {
	p.PosX = p.prevPosX
	p.PosY = p.prevPosY
	time.Sleep(invincibilityDuration)
	p.Invincible = false
}

func (p *Player) EnterDashAnim() {
	p.animator.SwitchClip(dashInitAnim)
}

func (p *Player) EnterSlamAnim() {
	p.animator.SwitchClip(slamAnim)
}

func NewPlayer() *Player {
	player := &Player{
		moveProgress:  1,
		sprite:        errs.MustNewImageFromFile(playerSpritePath),
		animator:      animation.NewAnimator(playerAnimationMap),
		damageOverlay: newDamageOverlay(),
		Invincible:    false,
		InputBuffer:   newInputBuffer(inputBufferDuration),
		State:         Idle,
	}

	player.finishedClipEventListener = events.NewEventListener(player.animator.FinishedClipEvent)
	return player
}
