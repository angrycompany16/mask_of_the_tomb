package vector64

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func FillRect(dst *ebiten.Image, x float64, y float64, width float64, height float64, clr color.Color, antialias bool) {
	vector.FillRect(dst, float32(x), float32(y), float32(width), float32(height), clr, antialias)
}

func StrokeRect(dst *ebiten.Image, x float64, y float64, width float64, height float64, strokeWidth float64, clr color.Color, antialias bool) {
	vector.StrokeRect(dst, float32(x), float32(y), float32(width), float32(height), float32(strokeWidth), clr, antialias)
}

func FillCircle(dst *ebiten.Image, cx float64, cy float64, r float64, clr color.Color, antialias bool) {
	vector.FillCircle(dst, float32(cx), float32(cy), float32(r), clr, antialias)
}

func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float64, strokeWidth float64, clr color.Color, antialias bool) {
	vector.StrokeLine(dst, float32(x0), float32(y0), float32(x1), float32(y1), float32(strokeWidth), clr, antialias)
}
