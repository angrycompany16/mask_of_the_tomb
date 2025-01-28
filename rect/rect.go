package rect

import "github.com/hajimehoshi/ebiten/v2"

type Rect struct {
	x      float64
	y      float64
	width  float64
	height float64
}

func (r *Rect) Left() float64 {
	return r.x
}

func (r *Rect) Right() float64 {
	return r.x + r.width
}

func (r *Rect) Top() float64 {
	return r.y
}

func (r *Rect) Bottom() float64 {
	return r.y + r.height
}

func (r *Rect) Center() (float64, float64) {
	return r.x + r.width, r.y + r.height
}

func (r *Rect) Draw(surf *ebiten.Image) {
	// Render debug rect
}

func NewRect(x, y, width, height float64) *Rect {
	return &Rect{x, y, width, height}
}
