package maths

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	return r.x + r.width/2, r.y + r.height/2
}

func (r *Rect) TopLeft() (float64, float64) {
	return r.Left(), r.Top()
}

func (r *Rect) Width() float64 {
	return r.width
}

func (r *Rect) Height() float64 {
	return r.height
}

func (r *Rect) Extended(dir Direction, length float64) Rect {
	newRect := *r
	switch dir {
	case DirUp:
		newRect.y -= length
		newRect.height += length
	case DirDown:
		newRect.height += length
	case DirRight:
		newRect.width += length
	case DirLeft:
		newRect.x -= length
		newRect.width += length
	}
	return newRect
}

func (r *Rect) SetPos(x, y float64) {
	r.x, r.y = x, y
}

// TODO: implement
func (r *Rect) Draw(surf *ebiten.Image) {
	// Render debug rect
}

func (r *Rect) Overlapping(other *Rect) bool {
	return r.Left() < other.Right() &&
		r.Right() > other.Left() &&
		r.Top() < other.Bottom() &&
		r.Bottom() > other.Top()
}

func (r *Rect) IsWithin(x, y float64) bool {
	return x <= r.Right() &&
		r.Left() <= x &&
		r.Top() <= y &&
		r.Bottom() >= y
}

func (r *Rect) RandomPointInside() (float64, float64) {
	return Lerp(r.Left(), r.Right(), rand.Float64()), Lerp(r.Bottom(), r.Top(), rand.Float64())
}

func (r *Rect) Lerp(other *Rect, t float64) Rect {
	return Rect{
		x:      Lerp(r.x, other.x, t),
		y:      Lerp(r.y, other.y, t),
		width:  Lerp(r.width, other.width, t),
		height: Lerp(r.height, other.height, t),
	}
}

func (r *Rect) RaycastDirectional(posX, posY float64, direction Direction) (bool, float64, float64) {
	if r.IsWithin(posX, posY) {
		return false, 0, 0
	}

	switch direction {
	case DirUp:
		return posX >= r.Left() && posX <= r.Right() && posY >= r.Top(), r.Left(), r.Top()
	case DirDown:
		return posX >= r.Left() && posX <= r.Right() && posY <= r.Bottom(), r.Left(), r.Top()
	case DirLeft:
		return posY >= r.Top() && posY <= r.Bottom() && posX >= r.Left(), r.Left(), r.Top()
	case DirRight:
		return posY >= r.Top() && posY <= r.Bottom() && posX <= r.Right(), r.Left(), r.Top()
	}
	return false, 0, 0
}

func NewRect(x, y, width, height float64) *Rect {
	return &Rect{x, y, width, height}
}

func RectFromImage(x, y float64, image *ebiten.Image) *Rect {
	size := image.Bounds().Size()
	return &Rect{
		x:      x,
		y:      y,
		width:  float64(size.X),
		height: float64(size.Y),
	}
}
