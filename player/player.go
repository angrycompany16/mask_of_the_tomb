package player

import (
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
	sprite                 *ebiten.Image
}

func (p *Player) Init(posX, posY float64) {
	p.SetPos(posX, posY)
}

func (p *Player) GetInput() MoveDirection {
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

func (p *Player) SetTarget(x, y float64) {
	p.targetPosX = x
	p.targetPosY = y
	p.moveDirX = math.Copysign(1, p.targetPosX-p.posX)
	p.moveDirY = math.Copysign(1, p.targetPosY-p.posY)

	// if x == p.posX {
	// 	p.targetPosY = y
	// 	p.moveDirY = math.Copysign(1, p.targetPosY-p.posY)
	// 	p.moveDirX = 0
	// }
	// if y == p.posY {
	// 	p.targetPosX = x
	// 	p.moveDirX = math.Copysign(1, p.targetPosX-p.posX)
	// 	p.moveDirY = 0
	// }
	// fmt.Println(x, y, p.moveDirX, p.moveDirY)
}

func (p *Player) IsMoving() bool {
	return p.moveDirX != 0 || p.moveDirY != 0
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

func (p *Player) Draw(surf *ebiten.Image) {
	DrawAt(p.sprite, surf, p.posX, p.posY)
}

func NewPlayer() *Player {
	return &Player{
		posX:         0,
		posY:         0,
		targetPosX:   0,
		targetPosY:   0,
		moveDirX:     0,
		moveDirY:     0,
		moveProgress: 1,
		sprite:       files.LazyImage(files.PlayerSpritePath),
	}
}
