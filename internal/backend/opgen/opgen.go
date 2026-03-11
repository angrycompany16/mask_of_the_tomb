package opgen

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// Contains functions for drawing images rotated and scaled around arbitrary pivots.
func OpConvert(op *ebiten.DrawImageOptions) *colorm.DrawImageOptions {
	return &colorm.DrawImageOptions{GeoM: op.GeoM, Blend: op.Blend, Filter: op.Filter}
}

func getPivotPx(src *ebiten.Image, pivotX, pivotY float64) (float64, float64) {
	s := src.Bounds().Size()
	return float64(s.X) * pivotX, float64(s.Y) * pivotY
}

func unpackPivot(pivot ...float64) (float64, float64) {
	pivotX, pivotY := 0.0, 0.0
	if len(pivot) == 1 {
		pivotX, pivotY = pivot[0], pivot[0]
	} else if len(pivot) > 1 {
		pivotX, pivotY = pivot[0], pivot[1]
	}
	return pivotX, pivotY
}

func Pos(src *ebiten.Image, x, y float64, pivot ...float64) *ebiten.DrawImageOptions {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Translate(x, y)
	return op
}

func PosScale(src *ebiten.Image, x, y, scaleX, scaleY float64, pivot ...float64) *ebiten.DrawImageOptions {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(pivotX*(1-scaleX), pivotY*(1-scaleY))

	op.GeoM.Translate(x, y)
	return op
}

func PosRot(src *ebiten.Image, x, y, angle float64, pivot ...float64) *ebiten.DrawImageOptions {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(angle)

	op.GeoM.Translate(x, y)
	return op
}

func PosScaleRot(src *ebiten.Image, x, y, angle, scaleX, scaleY float64, pivot ...float64) *ebiten.DrawImageOptions {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(pivotX*(1-scaleX), pivotY*(1-scaleY))

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(pivotX, pivotY)

	op.GeoM.Translate(x, y)
	return op
}
