package player

// New task - gameplay / health system

import (
	"mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/animation"
	"mask_of_the_tomb/internal/game/camera"
	"mask_of_the_tomb/internal/game/rendering"
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

// TODO: Convert to animation state machine, turn into asset file? (i.e. assets for anims)
type Player struct {
	PosX, PosY             float64
	targetPosX, targetPosY float64
	prevPosX, prevPosY     float64
	moveDirX, moveDirY     float64
	moveProgress           float64
	Hitbox                 *maths.Rect
	sprite                 *ebiten.Image
	playerIdleAnim         *animation.Animation
	Invincible             bool
	Disabled               bool
	damageOverlay          deathEffect
	angle                  float64
	InputBuffer            inputBuffer
	State                  playerState
}

func (p *Player) Init(posX, posY float64) {
	p.SetPos(posX, posY)
	p.Hitbox = maths.RectFromImage(posX, posY, p.sprite)
}

func (p *Player) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	ebitenrenderutil.DrawAtRotated(
		p.playerIdleAnim.GetSprite(),
		rendering.RenderLayers.Playerspace,
		p.PosX-camX,
		p.PosY-camY,
		p.angle,
		0.5,
		0.5,
	)

	p.damageOverlay.Draw()
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

func NewPlayer() *Player {
	return &Player{
		moveProgress: 1,
		sprite:       errs.MustNewImageFromFile(PlayerSpritePath),
		playerIdleAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(IdleSpritesheetPath)),
			0.1666667,
			animation.Strip,
			animation.Loop,
		),
		damageOverlay: newDamageOverlay(),
		Invincible:    false,
		InputBuffer:   newInputBuffer(inputBufferDuration),
		State:         StateIdle,
	}
}
