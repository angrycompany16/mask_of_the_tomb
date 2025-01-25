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

type FloatConvertible interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func F64[F FloatConvertible](num F) float64 {
	return float64(num)
}

func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

// The humble lerp
func Lerp(a, b, t float64) float64 {
	return a*(1.0-t) + b*t
}
