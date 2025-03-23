package player

// New task - gameplay / health system

import (
	"mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/animation"
	"mask_of_the_tomb/internal/game/camera"
	"mask_of_the_tomb/internal/game/events"
	"mask_of_the_tomb/internal/game/movebox"
	"mask_of_the_tomb/internal/game/rendering"
	"mask_of_the_tomb/internal/maths"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO: There's a bug which sometimes appears where the player moves a tiny bit out
// from the wall if they press the opposite move key right before hitting the wall
// TODO: Rewrite as much as possible using events rather than importing packages

const (
	moveSpeed             = 5.0
	defaultPlayerHealth   = 5.0
	invincibilityDuration = time.Second
	inputBufferDuration   = 0.1
)

// TODO: Allow for sprites which aren't exactly 16x16
// TODO: Rewrite player to use maths.direction instead of float angle
// TODO: turn animations into asset file? (i.e. assets for anims)
type Player struct {
	State                     playerState
	movebox                   *movebox.Movebox
	Hitbox                    *maths.Rect
	sprite                    *ebiten.Image
	jumpOffset, jumpOffsetvel float64
	direction                 maths.Direction
	animator                  *animation.Animator
	Invincible                bool
	Disabled                  bool
	damageOverlay             deathEffect
	InputBuffer               inputBuffer
	finishedMoveEventListener *events.EventListener
	finishedClipEventListener *events.EventListener
}

func (p *Player) Init(posX, posY float64) {
	p.SetPos(posX, posY)
	p.Hitbox = maths.RectFromImage(posX, posY, p.sprite)
	p.animator.SwitchClip(idleAnim)
}

func (p *Player) Draw() {
	posX, posY := p.movebox.GetPos()
	camX, camY := camera.GetPos()
	jumpOffsetX, jumpOffsetY := p.getJumpOffset()
	ebitenrenderutil.DrawAtRotated(
		p.animator.GetSprite(),
		rendering.RenderLayers.Playerspace,
		posX-camX-jumpOffsetX,
		posY-camY-jumpOffsetY,
		maths.ToRadians(p.direction),
		0.5,
		0.5,
	)

	p.damageOverlay.Draw()
}

func (p *Player) getJumpOffset() (float64, float64) {
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
	p.movebox.SetPos(x, y)
}

func (p *Player) SetRot(direction maths.Direction) {
	p.direction = direction
}

func (p *Player) GetPosCentered() (float64, float64) {
	s := p.sprite.Bounds().Size()
	x, y := p.movebox.GetPos()
	return x + float64(s.X)/2, y + float64(s.Y)/2
}

func (p *Player) SetTarget(x, y float64) {
	p.movebox.SetTarget(x, y)
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

func (p *Player) TakeDamage(damage float64) {
	// p.Invincible = true
	// p.Disabled = true
	// go p.hitHazard()
}

// func (p *Player) hitHazard() {
// 	time.Sleep(time.Millisecond * 300)
// 	p.Disabled = false
// 	p.damageOverlay.alpha = 1.0
// go p.disableInvincibility()
// }

// func (p *Player) disableInvincibility() {
// p.PosX = p.prevPosX
// p.PosY = p.prevPosY
// 	time.Sleep(invincibilityDuration)
// 	p.Invincible = false
// }

func (p *Player) EnterDashAnim() {
	p.animator.SwitchClip(dashInitAnim)
}

func (p *Player) EnterSlamAnim() {
	p.animator.SwitchClip(slamAnim)
}

func NewPlayer() *Player {
	player := &Player{
		movebox:       movebox.NewMovebox(moveSpeed),
		sprite:        errs.MustNewImageFromFile(playerSpritePath),
		animator:      animation.NewAnimator(playerAnimationMap),
		damageOverlay: newDamageOverlay(),
		Invincible:    false,
		InputBuffer:   newInputBuffer(inputBufferDuration),
		State:         Idle,
	}

	player.finishedMoveEventListener = events.NewEventListener(player.movebox.FinishedMoveEvent)
	player.finishedClipEventListener = events.NewEventListener(player.animator.FinishedClipEvent)
	// TODO: Shorten down these names holy flippin moly
	return player
}
