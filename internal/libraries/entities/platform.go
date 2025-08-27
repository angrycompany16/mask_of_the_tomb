package entities

import (
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

type Platform struct {
	Sprite *ebiten.Image
	Hitbox *maths.Rect
	Up     bool
}

func (p *Platform) Draw(ctx rendering.Ctx) {
	// ebitenrenderutil.DrawAt(c.Sprite, ctx.Dst, c.Hitbox.Left(), c.Hitbox.Top())
}

func NewPlatform(
	entity *ebitenLDTK.Entity,
	tileSize float64,
) *Platform {
	newPlatform := &Platform{}
	newPlatform.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)
	newPlatform.Sprite = ebiten.NewImage(int(entity.Width), int(entity.Height))

	directionField := errs.Must(entity.GetFieldByName("Direction"))
	newPlatform.Up = directionField.Point.Y*tileSize < entity.Px[1]

	return newPlatform
}
