package slambox

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/vector64"

	"github.com/hajimehoshi/ebiten/v2"
)

// Debug tool for drawing a rect with border
func DrawRect(dst *ebiten.Image, rect *maths.Rect, borderColor color.RGBA, fillColor color.RGBA) {
	vector64.FillRect(dst, rect.Left(), rect.Top(), rect.Width(), rect.Height(), fillColor, false)
	vector64.StrokeRect(dst, rect.Left(), rect.Top(), rect.Width(), rect.Height(), 1.0, borderColor, false)
}
