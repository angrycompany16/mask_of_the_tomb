package player

// New task - gameplay / health system

import (
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/game/animation"
	"mask_of_the_tomb/game/camera"
	"mask_of_the_tomb/game/health"
	"mask_of_the_tomb/rendering"
	. "mask_of_the_tomb/utils" // This is bad
	"mask_of_the_tomb/utils/rect"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayerState int

// Finish player state machine
const (
	Normal PlayerState = iota
	TakingDamage
)

type MoveDirection int

const (
	DirNone MoveDirection = iota - 1
	DirUp
	DirDown
	DirLeft
	DirRight
)

const (
	moveSpeed             = 5.0
	defaultPlayerHealth   = 5.0
	invincibilityDuration = time.Second
)

// Convert to animation state machine, turn into asset file?
type Player struct {
	posX, posY             float64
	targetPosX, targetPosY float64
	prevPosX, prevPosY     float64
	moveDirX, moveDirY     float64
	moveProgress           float64
	hitbox                 *rect.Rect
	score                  int
	sprite                 *ebiten.Image
	playerTestAnim         *animation.Animation
	health                 *health.HealthComponent
	invincible             bool
	disabled               bool
	damageOverlay          damageOverlay
	angle                  float64
}

func (p *Player) Init(posX, posY float64) {
	p.SetPos(posX, posY)
	p.hitbox = rect.FromImage(posX, posY, p.sprite)
}

func (p *Player) Update() {
	p.posX += moveSpeed * p.moveDirX * GlobalTimeScale
	p.posY += moveSpeed * p.moveDirY * GlobalTimeScale

	if p.moveDirX < 0 {
		p.posX = Clamp(p.posX, p.targetPosX, p.posX)
	} else if p.moveDirX > 0 {
		p.posX = Clamp(p.posX, p.posX, p.targetPosX)
	}
	if p.moveDirY < 0 {
		p.posY = Clamp(p.posY, p.targetPosY, p.posY)
	} else if p.moveDirY > 0 {
		p.posY = Clamp(p.posY, p.posY, p.targetPosY)
	}

	if p.posX == p.targetPosX {
		p.moveDirX = 0
	}
	if p.posY == p.targetPosY {
		p.moveDirY = 0
	}

	p.damageOverlay.Update()
	p.hitbox.SetPos(p.posX, p.posY)

	p.playerTestAnim.Update()
}

func (p *Player) Draw() {
	camX, camY := camera.GlobalCamera.GetPos()
	DrawAtRotated(
		p.playerTestAnim.GetSprite(),
		rendering.RenderLayers.Playerspace,
		p.posX-camX,
		p.posY-camY,
		p.angle,
		0.5,
		0.5,
	)

	// DrawAt(p.sprite, rendering.RenderLayers.Playerspace, p.posX-camX, p.posY-camY)
	p.damageOverlay.Draw()
}

func (p *Player) GetLevelSwapInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (p *Player) GetMoveInput() MoveDirection {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.angle = 0
		return DirUp
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		p.angle = math.Pi
		return DirDown
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		p.angle = math.Pi / 2
		return DirRight
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		p.angle = 3 * math.Pi / 2
		return DirLeft
	}
	return DirNone
}

func (p *Player) GetPos() (float64, float64) {
	return p.posX, p.posY
}

func (p *Player) SetPos(x, y float64) {
	p.posX, p.posY = x, y
	p.targetPosX, p.targetPosY = x, y
}

func (p *Player) GetPosCentered() (float64, float64) {
	s := p.sprite.Bounds().Size()
	return p.posX + F64(s.X)/2, p.posY + F64(s.Y)/2
}

func (p *Player) SetTarget(x, y float64) {
	p.targetPosX = x
	p.targetPosY = y
	p.prevPosX = p.posX
	p.prevPosY = p.posY
	p.moveDirX = math.Copysign(1, p.targetPosX-p.posX)
	p.moveDirY = math.Copysign(1, p.targetPosY-p.posY)
}

func (p *Player) GetSize() (float64, float64) {
	s := p.sprite.Bounds().Size()
	return float64(s.X), float64(s.Y)
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) SetScore(score int) {
	p.score = score
}

func (p *Player) GetMovementSize() (float64, float64) {
	return moveSpeed * p.moveDirX, moveSpeed * p.moveDirY
}

func (p *Player) GetHitbox() *rect.Rect {
	return p.hitbox
}

func (p *Player) IsMoving() bool {
	return p.moveDirX != 0 || p.moveDirY != 0
}

func (p *Player) TakeDamage(damage float64) {
	p.health.TakeDamage(damage)
	p.invincible = true
	p.disabled = true
	go p.hitHazard()
}

func (p *Player) hitHazard() {
	time.Sleep(time.Millisecond * 300)
	p.disabled = false
	p.damageOverlay.alpha = 1.0
	go p.disableInvincibility()
}

func (p *Player) disableInvincibility() {
	p.posX = p.prevPosX
	p.posY = p.prevPosY
	time.Sleep(invincibilityDuration)
	p.invincible = false
}

func (p *Player) IsInvincible() bool {
	return p.invincible
}

func (p *Player) IsDisabled() bool {
	return p.disabled
}

func NewPlayer() *Player {
	return &Player{
		moveProgress: 1,
		sprite:       files.LazyImage(PlayerSpritePath),
		playerTestAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(files.LazyImage(IdleSpritesheetPath)),
			0.1666667,
			animation.Strip,
			animation.Loop,
		),
		health:        health.NewHealthComponent(defaultPlayerHealth),
		damageOverlay: newDamageOverlay(),
		invincible:    false,
	}
}
