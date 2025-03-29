package ebitenrenderutil

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// Contains functions for drawing images rotated and scaled around arbitrary pivots

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

func DrawAt(src, dst *ebiten.Image, x, y float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func DrawAtScaled(src, dst *ebiten.Image, x, y, scaleX, scaleY float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(pivotX*(1-scaleX), pivotY*(1-scaleY))

	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func DrawAtRotated(src, dst *ebiten.Image, x, y, angle float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(pivotX, pivotY)

	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func DrawAtRotatedScaled(src, dst *ebiten.Image, x, y, angle, scaleX, scaleY float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(pivotX*(1-scaleX), pivotY*(1-scaleY))

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(pivotX, pivotY)

	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func RotatedScaledOp(src *ebiten.Image, x, y, angle, scaleX, scaleY float64, pivot ...float64) *ebiten.DrawImageOptions {
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

func getUICoords(dst *ebiten.Image, x, y float64) (float64, float64) {
	s := dst.Bounds().Size()
	return x + float64(s.X)/2, -y + float64(s.Y)/2
}

func UIDrawAt(src, dst *ebiten.Image, x, y float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	x, y = getUICoords(dst, x, y)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func UIDrawAtScaled(src, dst *ebiten.Image, x, y, scaleX, scaleY float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	x, y = getUICoords(dst, x, y)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(-pivotX*scaleX, -pivotY*scaleY)

	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func UIDrawAtRotated(src, dst *ebiten.Image, x, y, angle float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	x, y = getUICoords(dst, x, y)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(angle)

	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func UIDrawAtRotatedScaled(src, dst *ebiten.Image, x, y, angle, scaleX, scaleY float64, pivot ...float64) {
	pivotX, pivotY := unpackPivot(pivot...)
	pivotX, pivotY = getPivotPx(src, pivotX, pivotY)
	x, y = getUICoords(dst, x, y)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(pivotX*(1-scaleX), pivotY*(1-scaleY))

	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(angle)

	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func OpConvert(op *ebiten.DrawImageOptions) *colorm.DrawImageOptions {
	return &colorm.DrawImageOptions{GeoM: op.GeoM, Blend: op.Blend, Filter: op.Filter}
}
