package entities

import (
	"image/color"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Lantern struct {
	sprite           *ebiten.Image
	x, y             float64
	anchorX, anchorY float64
}

func (l *Lantern) Update() {}

func (l *Lantern) Draw(ctx rendering.Ctx) {
	vector.StrokeLine(
		ctx.Dst,
		// centered
		float32(l.x+float64(l.sprite.Bounds().Dx()/2)),
		float32(l.y+float64(l.sprite.Bounds().Dy()/2)),
		float32(l.anchorX+float64(l.sprite.Bounds().Dx()/2)),
		float32(l.anchorY+float64(l.sprite.Bounds().Dy()/2)),
		2,
		color.RGBA{128, 128, 128, 255},
		false,
	)
	ebitenrenderutil.DrawAt(l.sprite, ctx.Dst, l.x, l.y)
}

func NewLantern(
	entity *ebitenLDTK.Entity,
	tileSize float64,
) *Lantern {
	newLantern := Lantern{}
	newLantern.x, newLantern.y = entity.Px[0], entity.Px[1]
	newLantern.sprite = errs.Must(assettypes.GetImageAsset("lanternSprite"))

	anchorPointField := errs.Must(entity.GetFieldByName("Anchor"))

	newLantern.anchorX = anchorPointField.Point.X * tileSize
	newLantern.anchorY = anchorPointField.Point.Y * tileSize

	return &newLantern
}
