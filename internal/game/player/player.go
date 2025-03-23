package player

// New task - gameplay / health system

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/animation"
	"mask_of_the_tomb/internal/game/events"
	"mask_of_the_tomb/internal/game/movebox"
	"mask_of_the_tomb/internal/game/player/deathanim"
	"mask_of_the_tomb/internal/maths"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO: There's a bug which sometimes appears where the player moves a tiny bit out
// from the wall if they press the opposite move key right before hitting the wall
// Can maybe be solved with PreUpdate()?

const (
	moveSpeed             = 5.0
	defaultPlayerHealth   = 5.0
	invincibilityDuration = time.Second
	inputBufferDuration   = 0.1
)

// TODO: Allow for sprites which aren't exactly 16x16
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
	InputBuffer               inputBuffer
	deathAnim                 *deathanim.DeathAnim
	// Events
	OnDeath *events.Event
	// Listeners
	moveFinishedListener *events.EventListener
	clipFinishedListener *events.EventListener
}

func (p *Player) Init(posX, posY float64) {
	p.SetPos(posX, posY)
	p.Hitbox = maths.RectFromImage(posX, posY, p.sprite)
	p.animator.SwitchClip(idleAnim)
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

func (p *Player) Die() {
	p.Disabled = true
	p.State = Dying
	p.deathAnim.Play()
	// TODO: Make it centered
	p.deathAnim.SetPos(p.movebox.GetPos())

	// May not be necessary
	p.OnDeath.Raise(events.EventInfo{})
}

func (p *Player) Respawn() {
	p.Disabled = false
	p.direction = maths.DirUp
	p.State = Idle
}

func (p *Player) EnterDashAnim() {
	p.animator.SwitchClip(dashInitAnim)
}

func (p *Player) EnterSlamAnim() {
	p.animator.SwitchClip(slamAnim)
}

func NewPlayer() *Player {
	player := &Player{
		movebox:  movebox.NewMovebox(moveSpeed),
		sprite:   errs.MustNewImageFromFile(playerSpritePath),
		animator: animation.NewAnimator(playerAnimationMap),
		// damageOverlay: newDamageOverlay(),
		Invincible:  false,
		InputBuffer: newInputBuffer(inputBufferDuration),
		State:       Idle,
		OnDeath:     events.NewEvent(),
		deathAnim:   deathanim.NewDeathAnim(),
	}

	player.moveFinishedListener = events.NewEventListener(player.movebox.OnMoveFinished)
	player.clipFinishedListener = events.NewEventListener(player.animator.OnClipFinished)
	// TODO: Shorten down these names holy flippin moly
	return player
}
