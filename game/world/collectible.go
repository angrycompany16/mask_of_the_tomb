package world

import "github.com/hajimehoshi/ebiten/v2"

type Collectible struct {
	sprite     ebiten.Image
	posX, posY float64
}

func (c *Collectible) draw(surf ebiten.Image, camX, camY float64) {

}
