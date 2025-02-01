package player

import (
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/files"
	. "mask_of_the_tomb/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MoveDirection int

const (
	DirNone = iota - 1
	DirUp
	DirDown
	DirLeft
	DirRight
)

const (
	moveSpeed = 5.0
)

type Player struct {
	posX, posY             float64
	targetPosX, targetPosY float64
	moveDirX, moveDirY     float64
	moveProgress           float64
	score                  int
	sprite                 *ebiten.Image
}

func (p *Player) Init(posX, posY float64) {
	p.SetPos(posX, posY)
}

func (p *Player) Update() {
	p.posX += moveSpeed * p.moveDirX
	p.posY += moveSpeed * p.moveDirY

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
}

func (p *Player) Draw(surf *ebiten.Image, camX, camY float64) {
	DrawAt(p.sprite, surf, p.posX-camX, p.posY-camY)
}

func (p *Player) GetLevelSwapInput() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (p *Player) GetMoveInput() MoveDirection {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return DirUp
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return DirDown
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return DirRight
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
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

func (p *Player) IsMoving() bool {
	return p.moveDirX != 0 || p.moveDirY != 0
}

func NewPlayer() *Player {
	return &Player{
		moveProgress: 1,
		sprite:       files.LazyImage(PlayerSpritePath),
	}
}
