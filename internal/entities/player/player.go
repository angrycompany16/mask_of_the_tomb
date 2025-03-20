package player

// New task - gameplay / health system

import (
	"mask_of_the_tomb/internal/engine/advertisers"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/engine/events"
	"mask_of_the_tomb/internal/entities/animation"
	"mask_of_the_tomb/internal/entities/camera/pubcamera"
	pubgame "mask_of_the_tomb/internal/entities/game/pub"
	"mask_of_the_tomb/internal/entities/player/moveBox"
	"mask_of_the_tomb/internal/libraries/assets/ebitenrenderutil"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/gameplay/inputbuffer"
	"mask_of_the_tomb/internal/libraries/maths"
	"mask_of_the_tomb/internal/libraries/rendering"
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
	moveboxName           = "playerMoveBox"
)

// TODO: Reduce the pointer-y ness
// Not everything needs to be a pointer!
// TODO: Rewrite player to use maths.direction instead of float angle
// TODO: turn animations into asset file? (i.e. assets for anims)
type player struct {
	*entities.Entity
	// PosX, PosY                float64
	// targetPosX, targetPosY    float64
	// prevPosX, prevPosY        float64

	// moveDirX, moveDirY        float64
	// moveProgress              float64
	jumpOffset, jumpOffsetvel float64
	Hitbox                    *maths.Rect
	sprite                    *ebiten.Image
	animator                  *animation.Animator
	Invincible                bool
	Disabled                  bool
	// damageOverlay             deathEffect
	angle       float64
	InputBuffer inputbuffer.InputBuffer
	// State                     playerState
	finishedClipEventListener *events.EventListener
}

type playerInitMsg struct {
	Width, Height float64
}

func New(posX, posY float64) (*player, playerInitMsg) {
	_player := player{
		sprite:   errs.MustNewImageFromFile(playerSpritePath),
		animator: animation.NewAnimator(playerAnimationMap),
		// damageOverlay: newDamageOverlay(),
		Invincible: false,
		// InputBuffer: newInputBuffer(inputBufferDuration),
		// State:         Idle,
	}
	_player.Entity = entities.RegisterEntity(&_player, "Player")

	_player.finishedClipEventListener = events.NewEventListener(_player.animator.FinishedClipEvent)

	// _player.SetPos(posX, posY)
	_player.Hitbox = maths.RectFromImage(posX, posY, _player.sprite)
	_player.animator.SwitchClip(idleAnim)

	width, height := _player.GetSize()

	initMsg := playerInitMsg{
		Width:  width,
		Height: height,
	}

	_movebox := moveBox.NewMovebox(posX, posY, moveSpeed, moveboxName)
	_movebox.SetParent(_player.Entity)
	return &_player, initMsg
}

// TODO: Semplify getting advertisers?
func (p *player) Update() {
	val := advertisers.Get(pubgame.GameEntityName).Read().(pubgame.GameAdvertiser)
	if val.State != pubgame.StatePlaying {
		return
	}

	// switch p.State {
	// case Slamming:
	// 	if _, raised := p.finishedClipEventListener.Poll(); raised {
	// 		p.State = Idle
	// 		p.jumpOffset = 0
	// 		p.jumpOffsetvel = 0
	// 	}

	// 	if p.jumpOffsetvel > 0 {
	// 		p.jumpOffsetvel -= 0.1
	// 	} else {
	// 		p.jumpOffsetvel -= 0.25
	// 	}

	// 	p.jumpOffset += p.jumpOffsetvel
	// 	p.jumpOffset = maths.Clamp(p.jumpOffset, 0, 1000000)
	// case Idle:
	// 	p.animator.SwitchClip(idleAnim)
	// case Moving:
	// 	moveDir := p.InputBuffer.ReadSingle()
	// 	if moveDir != maths.DirNone && p.CanMove() && !p.Disabled {
	// 		pubplayer.OnMove.Raise(events.EventInfo{Data: pubplayer.MoveEvent{
	// 			Direction: moveDir,
	// 			Hitbox:    *p.Hitbox,
	// 		}})

	// 		// Get slamboxes via advertiser
	// 		slambox := g.world.ActiveLevel.GetSlamboxHit(p.Hitbox, moveDir)
	// 		if slambox != nil {
	// 			g.player.StartSlamming(playerMove)
	// 			if !slamming {
	// 				slamming = true
	// 				go g.DoSlam(slambox, playerMove)
	// 			}
	// 		} else {
	// 			newRect, _ := g.world.ActiveLevel.TilemapCollider.ProjectRect(g.player.Hitbox, playerMove, g.world.ActiveLevel.GetSlamboxColliders())
	// 			if newRect != *g.player.Hitbox {
	// 				g.player.EnterDashAnim()
	// 				g.player.SetRot(playerMove)
	// 				g.player.SetTarget(newRect.Left(), newRect.Top())
	// 				g.player.State = player.Moving
	// 			}
	// 		}

	// 	}
	// }
	// p.PosX += moveSpeed * p.moveDirX
	// p.PosY += moveSpeed * p.moveDirY

	// if p.moveDirX < 0 {
	// 	p.PosX = maths.Clamp(p.PosX, p.targetPosX, p.PosX)
	// } else if p.moveDirX > 0 {
	// 	p.PosX = maths.Clamp(p.PosX, p.PosX, p.targetPosX)
	// }
	// if p.moveDirY < 0 {
	// 	p.PosY = maths.Clamp(p.PosY, p.targetPosY, p.PosY)
	// } else if p.moveDirY > 0 {
	// 	p.PosY = maths.Clamp(p.PosY, p.PosY, p.targetPosY)
	// }

	// if p.PosX == p.targetPosX {
	// 	p.moveDirX = 0
	// }
	// if p.PosY == p.targetPosY {
	// 	p.moveDirY = 0
	// }

	// 	if p.PosX == p.targetPosX && p.PosY == p.targetPosY {
	// 		p.angle = p.angle - math.Pi
	// 		p.State = Idle
	// 	}
	// }

	// direction := p.getMoveInput()
	// if direction != maths.DirNone {
	// 	p.InputBuffer.set(direction)
	// }

	// p.InputBuffer.update()

	// p.damageOverlay.Update()
	// p.Hitbox.SetPos(p.PosX, p.PosY)

	// Turn into entity. Ideally you should never need to call an 'update' function in
	// udpate
	p.animator.Update()
}

func (p *player) Draw() {
	moveboxAdv := advertisers.Get(moveboxName).Read().(pubcamera.CameraAdvertiser)
	posX, posY := moveboxAdv.PosX, moveboxAdv.PosY

	gameAdv := advertisers.Get(pubgame.GameEntityName).Read().(pubgame.GameAdvertiser)
	if gameAdv.State == pubgame.StateMainMenu {
		return
	}

	camAdv := advertisers.Get(pubcamera.CameraEntityName).Read().(pubcamera.CameraAdvertiser)
	camX, camY := camAdv.PosX, camAdv.PosY

	jumpOffsetX, jumpOffsetY := p.getJumpOffset()
	ebitenrenderutil.DrawAtRotated(
		p.animator.GetSprite(),
		rendering.RenderLayers.Playerspace,
		posX-camX-jumpOffsetX,
		posY-camY-jumpOffsetY,
		p.angle,
		0.5,
		0.5,
	)

	// p.damageOverlay.Draw()
}

func (p *player) getJumpOffset() (float64, float64) {
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

func (p *player) GetLevelSwapInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (p *player) getMoveInput() maths.Direction {
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

// func (p *player) SetPos(x, y float64) {
// 	p.PosX, p.PosY = x, y
// 	p.targetPosX, p.targetPosY = x, y
// }

func (p *player) SetRot(direction maths.Direction) {
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

// func (p *player) GetPosCentered() (float64, float64) {
// 	s := p.sprite.Bounds().Size()
// 	return p.PosX + float64(s.X)/2, p.PosY + float64(s.Y)/2
// }

func (p *player) GetSize() (float64, float64) {
	s := p.sprite.Bounds().Size()
	return float64(s.X), float64(s.Y)
}

// func (p *player) GetMovementSize() (float64, float64) {
// 	return moveSpeed * p.moveDirX, moveSpeed * p.moveDirY
// }

// func (p *player) CanMove() bool {
// 	return p.State == Idle
// }

// func (p *player) IsMoving() bool {
// 	return p.moveDirX != 0 || p.moveDirY != 0
// }

// func (p *player) TakeDamage(damage float64) {
// 	p.Invincible = true
// 	p.Disabled = true
// 	go p.hitHazard()
// }

// func (p *player) hitHazard() {
// 	time.Sleep(time.Millisecond * 300)
// 	p.Disabled = false
// 	p.damageOverlay.alpha = 1.0
// 	go p.disableInvincibility()
// }

// func (p *player) disableInvincibility() {
// 	p.PosX = p.prevPosX
// 	p.PosY = p.prevPosY
// 	time.Sleep(invincibilityDuration)
// 	p.Invincible = false
// }

func (p *player) EnterDashAnim() {
	p.animator.SwitchClip(dashInitAnim)
}

func (p *player) EnterSlamAnim() {
	p.animator.SwitchClip(slamAnim)
}

// func NewPlayer() *player {
// 	player := &player{
// 		moveProgress:  1,
// 		sprite:        errs.MustNewImageFromFile(playerSpritePath),
// 		animator:      animation.NewAnimator(playerAnimationMap),
// 		damageOverlay: newDamageOverlay(),
// 		Invincible:    false,
// 		InputBuffer:   newInputBuffer(inputBufferDuration),
// 		State:         Idle,
// 	}

// 	player.finishedClipEventListener = events.NewEventListener(player.animator.FinishedClipEvent)
// 	return player
// }
