package vector64

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func FillCircle(dst *ebiten.Image, cx float64, cy float64, r float64, clr color.Color, antialias bool) {
	vector.FillCircle(dst, float32(cx), float32(cy), float32(r), clr, antialias)
}

func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float64, strokeWidth float64, clr color.Color, antialias bool) {
	vector.StrokeLine(dst, float32(x0), float32(y0), float32(x1), float32(y1), float32(strokeWidth), clr, antialias)
}
