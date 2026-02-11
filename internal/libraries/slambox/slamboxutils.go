package slambox

import (
	"image/color"
	"mask_of_the_tomb/internal/core/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Debug tool for drawing a rect with border
func DrawRect(dst *ebiten.Image, rect *maths.Rect, borderColor color.RGBA, fillColor color.RGBA) {
	vector.DrawFilledRect(dst, float32(rect.Left()), float32(rect.Top()), float32(rect.Width()), float32(rect.Height()), fillColor, false)
	vector.StrokeRect(dst, float32(rect.Left()+1), float32(rect.Top()+1), float32(rect.Width()-1), float32(rect.Height()-1), 1.0, borderColor, false)
}
