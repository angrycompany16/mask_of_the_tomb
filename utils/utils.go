package utils

import "github.com/hajimehoshi/ebiten/v2"

func Op() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{}
}

func DrawAt(src, dst *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(src, op)
}

func OpMove(op *ebiten.DrawImageOptions, x, y float64) {
	op.GeoM.Translate(x, y)
}

func OpScale(op *ebiten.DrawImageOptions, x, y float64) {
	op.GeoM.Scale(x, y)
}

func OpScaleCentered(src *ebiten.Image, op *ebiten.DrawImageOptions, x, y float64) {
	s := src.Bounds().Size()
	op.GeoM.Translate(float64(s.X)/2, float64(s.Y)/2)
	op.GeoM.Scale(x, y)
}
