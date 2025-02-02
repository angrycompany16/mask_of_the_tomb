package world

import (
	"image"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	. "mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Collectible struct {
	sprite     *ebiten.Image
	posX, posY float64
	collected  bool
	iid        string
}

func (c *Collectible) draw(surf *ebiten.Image, camX, camY float64) {
	if c.collected {
		return
	}
	DrawAt(c.sprite, surf, c.posX-camX, c.posY-camY)
}

func newCollectible(
	collected bool,
	entityInstance ebitenLDTK.EntityInstance,
	defs ebitenLDTK.Defs,
) Collectible {
	newCollectible := Collectible{}
	newCollectible.collected = collected
	newCollectible.iid = entityInstance.Iid
	newCollectible.posX = entityInstance.Px[0]
	newCollectible.posY = entityInstance.Px[1]

	tileset, err := defs.GetTilesetByUid(entityInstance.Tile.TilesetUid)
	HandleLazy(err)
	tileSize := tileset.TileGridSize
	newCollectible.sprite = tileset.Image.SubImage(
		image.Rect(
			int(entityInstance.Tile.X),
			int(entityInstance.Tile.Y),
			int(entityInstance.Tile.X+tileSize),
			int(entityInstance.Tile.Y+tileSize),
		),
	).(*ebiten.Image)
	return newCollectible
}
