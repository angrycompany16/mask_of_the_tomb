package maths

import (
	"math/rand/v2"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
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
	if r.Left() < other.Right() &&
		r.Right() > other.Left() &&
		r.Top() < other.Bottom() &&
		r.Bottom() > other.Top() {
		return true
	}
	return false
}

func (r *Rect) RandomPointInside() (float64, float64) {
	return Lerp(r.Left(), r.Right(), rand.Float64()), Lerp(r.Bottom(), r.Top(), rand.Float64())
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

func RectFromEntity(entityInstance *ebitenLDTK.Entity) *Rect {
	return &Rect{
		x:      entityInstance.Px[0],
		y:      entityInstance.Px[1],
		width:  entityInstance.Width,
		height: entityInstance.Height,
	}
}
