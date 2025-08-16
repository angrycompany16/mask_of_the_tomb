package entities

import (
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

type Catcher struct {
	Sprite *ebiten.Image
	Hitbox *maths.Rect
}

func (c *Catcher) Draw(ctx rendering.Ctx) {
	// ebitenrenderutil.DrawAt(c.Sprite, ctx.Dst, c.Hitbox.Left(), c.Hitbox.Top())
}

func NewCatcher(
	entity *ebitenLDTK.Entity,
) *Catcher {
	newCatcher := &Catcher{}
	newCatcher.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)
	newCatcher.Sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))

	return newCatcher
}
