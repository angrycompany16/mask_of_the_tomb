package player

import (
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	posX, posY float64
	sprite     *ebiten.Image
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.posY += 1
	}
}

func (p *Player) Draw(surf *ebiten.Image) {
	utils.DrawAt(p.sprite, surf, p.posX, p.posY)
}

func NewPlayer() *Player {
	return &Player{
		posX:   32,
		posY:   32,
		sprite: files.LazyImage(files.SpritePathA),
	}
}
