package turret

import (
	"fmt"
	"image/color"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/assettypes"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	directionFieldName = "Direction"
)

type Turret struct {
	Sprite           *ebiten.Image
	Hitbox           maths.Rect
	aimDirX, aimDirY float64
	RayEndX, RayEndY float64
	dead             bool
}

func (t *Turret) ShouldFire(target *maths.Rect) bool {
	if t.dead {
		return false
	}

	isLeft := target.Left() < t.Hitbox.Right()
	isRight := target.Right() > t.Hitbox.Left()

	isAbove := target.Top() < t.Hitbox.Bottom()
	isBelow := target.Bottom() > t.Hitbox.Top()

	LOShrz := (isLeft && t.aimDirX < 0 || isRight && t.aimDirX > 0) && isAbove && isBelow
	LOSvrt := (isAbove && t.aimDirY < 0 || isBelow && t.aimDirY > 0) && isLeft && isRight

	return LOShrz || LOSvrt
}

func (t *Turret) Draw(ctx rendering.Ctx) {
	if t.dead || (t.RayEndX == 0 && t.RayEndY == 0) {
		return
	}
	vector.DrawFilledCircle(ctx.Dst, float32(t.RayEndX), float32(t.RayEndY), 5.0, color.White, false)
	cx, cy := t.Hitbox.Center()
	// Note: Will depend on direction
	vector.StrokeLine(ctx.Dst, float32(cx+float64(t.Sprite.Bounds().Dx()/2)), float32(cy), float32(t.RayEndX), float32(t.RayEndY), 4.0, color.White, false)
	ebitenrenderutil.DrawAt(t.Sprite, ctx.Dst, t.Hitbox.Left(), t.Hitbox.Top())
}

func (t *Turret) Die() {
	t.dead = true
}

func (t *Turret) GetAimDir() maths.Direction {
	if t.aimDirX < 0 {
		return maths.DirLeft
	} else if t.aimDirX > 0 {
		return maths.DirRight
	} else if t.aimDirY < 0 {
		return maths.DirUp
	} else if t.aimDirY > 0 {
		return maths.DirDown
	}
	return maths.DirNone
}

func NewTurret(
	entity *ebitenLDTK.Entity,
	tileSize float64,
) *Turret {
	newTurret := &Turret{}
	newTurret.Hitbox = *maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)
	newTurret.Sprite = errs.Must(assettypes.GetImageAsset("turretSprite"))
	// newTurret.Sprite.Fill(color.RGBA{255, 0, 0, 255})
	directionField := errs.Must(entity.GetFieldByName(directionFieldName))

	newTurret.aimDirX = directionField.Point.X*tileSize - entity.Px[0]
	newTurret.aimDirY = directionField.Point.Y*tileSize - entity.Px[1]

	fmt.Println(newTurret.aimDirX, newTurret.aimDirY)

	return newTurret
}
