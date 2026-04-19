package maths

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type RectExteriorSection int

const (
	TOP_LEFT RectExteriorSection = iota
	TOP_MIDDLE
	TOP_RIGHT
	MIDDLE_LEFT
	MIDDLE_RIGHT
	BOTTOM_LEFT
	BOTTOM_MIDDLE
	BOTTOM_RIGHT
	EXTERIOR_SECTION_INVALID
)

type RectInteriorSection int

const (
	TOP RectInteriorSection = iota
	LEFT
	RIGHT
	BOTTOM
	INTERIOR_SECTION_INVALID
)

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (r *Rect) Left() float64 {
	return r.X
}

func (r *Rect) Right() float64 {
	return r.X + r.Width
}

func (r *Rect) Top() float64 {
	return r.Y
}

func (r *Rect) Bottom() float64 {
	return r.Y + r.Height
}

// Pixel perfect version. Should replace the usual right?
func (r *Rect) PPBottom() float64 {
	return r.Y + r.Height - 1
}

func (r *Rect) Center() (float64, float64) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *Rect) Cx() float64 {
	cx, _ := r.Center()
	return cx
}

func (r *Rect) Cy() float64 {
	_, cy := r.Center()
	return cy
}

func (r *Rect) TopLeft() (float64, float64) {
	return r.Left(), r.Top()
}

// func (r *Rect) Width() float64 {
// 	return r.width
// }

// func (r *Rect) Height() float64 {
// 	return r.height
// }

func (r *Rect) Size() (float64, float64) {
	return r.Width, r.Height
}

func (r *Rect) HalfSize() (float64, float64) {
	return r.Width / 2, r.Height / 2
}

func (r *Rect) Extended(dir Direction, length float64) *Rect {
	newRect := *r
	switch dir {
	case DirUp:
		newRect.Y -= length
		newRect.Height += length
	case DirDown:
		newRect.Height += length
	case DirRight:
		newRect.Width += length
	case DirLeft:
		newRect.X -= length
		newRect.Width += length
	}
	return &newRect
}

func (r *Rect) SetSize(w, h float64) {
	r.Width = w
	r.Height = h
}

func (r *Rect) SetPos(x, y float64) {
	r.X, r.Y = x, y
}

func (r *Rect) Translate(x, y float64) {
	r.X += x
	r.Y += y
}

func (r *Rect) Translated(x, y float64) Rect {
	newRect := *r
	newRect.X += x
	newRect.Y += y
	return newRect
}

func (r *Rect) Extend(x, y float64) {
	r.Width += x
	r.Height += y
}

func (r *Rect) Overlapping(other *Rect) bool {
	return r.Left() < other.Right() &&
		r.Right() > other.Left() &&
		r.Top() < other.Bottom() &&
		r.Bottom() > other.Top()
}

// Checks if a point is inside the rect
func (r *Rect) Contains(x, y float64) bool {
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
		X:      Lerp(r.X, other.X, t),
		Y:      Lerp(r.Y, other.Y, t),
		Width:  Lerp(r.Width, other.Width, t),
		Height: Lerp(r.Height, other.Height, t),
	}
}

// Returns the bounding box for the list of rects
func BB(rects []*Rect) *Rect {
	BBrect := NewRect(0, 0, 0, 0)
	minX := math.Inf(1)
	minY := math.Inf(1)
	for _, rect := range rects {
		if rect.X < minX {
			minX = rect.X
			BBrect.X = rect.X
		}
		if rect.Y < minY {
			minY = rect.Y
			BBrect.Y = rect.Y
		}
	}
	maxWidth := math.Inf(-1)
	maxHeight := math.Inf(-1)
	for _, rect := range rects {
		if rect.Right()-minX > maxWidth {
			maxWidth = rect.Right() - minX
			BBrect.Width = rect.Right() - minX
		}

		if rect.Bottom()-minY > maxHeight {
			maxHeight = rect.Bottom() - minY
			BBrect.Height = rect.Bottom() - minY
		}
	}
	return BBrect
}

// Returns a rect which is half of the original one, according
// to the passed in direction
func (r *Rect) GetHalved(dir Direction) *Rect {
	switch dir {
	case DirUp:
		return NewRect(r.X, r.Y, r.Width, r.Height/2)
	case DirDown:
		return NewRect(r.X, r.Y+r.Height/2, r.Width, r.Height/2)
	case DirLeft:
		return NewRect(r.X, r.Y, r.Width/2, r.Height)
	case DirRight:
		return NewRect(r.X+r.Width/2, r.Y, r.Width/2, r.Height)
	}
	return NewRect(0, 0, 0, 0)
}

// Checks if a ray starting in posX, posY and travelling in the given direction will
// intersect the rect.
func (r *Rect) RaycastDirectional(posX, posY float64, direction Direction) bool {
	if r.Contains(posX, posY) {
		// Should be true?
		return false
	}

	switch direction {
	case DirUp:
		return posX >= r.Left() && posX <= r.Right() && posY >= r.Top()
	case DirDown:
		return posX >= r.Left() && posX <= r.Right() && posY <= r.Bottom()
	case DirLeft:
		return posY >= r.Top() && posY <= r.Bottom() && posX >= r.Left()
	case DirRight:
		return posY >= r.Top() && posY <= r.Bottom() && posX <= r.Right()
	}
	return false
}

// Returns whether the point (x, y) lies directly above/below/left/right to
// the center of the rect.
func (r *Rect) IsInDirection(x, y float64, direction Direction) bool {
	switch direction {
	case DirUp:
		return x == r.Cx() && y >= r.Cy()
	case DirDown:
		return x == r.Cx() && y <= r.Cy()
	case DirLeft:
		return x <= r.Cx() && y == r.Cy()
	case DirRight:
		return x >= r.Cx() && y == r.Cy()
	}
	return false
}

// Moves the rect so that it reaches (x, y) (or avoids the point)
func (r *Rect) Reach(x, y float64) Rect {
	if r.Contains(x, y) {
		switch r.GetInteriorSection(x, y) {
		case TOP:
			return Rect{r.X, y, r.Width, r.Height}
		case LEFT:
			return Rect{x, r.Y, r.Width, r.Height}
		case RIGHT:
			return Rect{x - r.Width, r.Y, r.Width, r.Height}
		case BOTTOM:
			return Rect{r.X, y - r.Height, r.Width, r.Height}
		}
	}

	switch r.GetExteriorSection(x, y) {
	case TOP_LEFT:
		return Rect{x, y, r.Width, r.Height}
	case TOP_MIDDLE:
		return Rect{r.X, y, r.Width, r.Height}
	case TOP_RIGHT:
		return Rect{x - r.Width, y, r.Width, r.Height}
	case MIDDLE_LEFT:
		return Rect{x, r.Y, r.Width, r.Height}
	case MIDDLE_RIGHT:
		return Rect{x - r.Width, r.Y, r.Width, r.Height}
	case BOTTOM_LEFT:
		return Rect{x, y - r.Height, r.Width, r.Height}
	case BOTTOM_MIDDLE:
		return Rect{r.X, y - r.Height, r.Width, r.Height}
	case BOTTOM_RIGHT:
		return Rect{x - r.Width, y - r.Height, r.Width, r.Height}
	}
	return *r
}

// Finds the point on the rect that is closest to (x, y)
func (r *Rect) FindClosestPoint(x, y float64) (float64, float64) {
	if r.Contains(x, y) {
		return x, y
	}
	switch r.GetExteriorSection(x, y) {
	case TOP_LEFT:
		return r.Left(), r.Top()
	case TOP_MIDDLE:
		return x, r.Top()
	case TOP_RIGHT:
		return r.Right(), r.Top()
	case MIDDLE_LEFT:
		return r.Left(), y
	case MIDDLE_RIGHT:
		return r.Right(), y
	case BOTTOM_LEFT:
		return r.Left(), r.Bottom()
	case BOTTOM_MIDDLE:
		return x, r.Bottom()
	case BOTTOM_RIGHT:
		return r.Right(), r.Bottom()
	}

	return x, y
}

func (r *Rect) FindClosestPointOnEdge(x, y float64) (float64, float64) {
	if !r.Contains(x, y) {
		return r.FindClosestPoint(x, y)
	}

	switch r.GetInteriorSection(x, y) {
	case TOP:
		return x, r.Top()
	case LEFT:
		return r.Left(), y
	case RIGHT:
		return r.Right(), y
	case BOTTOM:
		return x, r.Bottom()
	}

	return x, y
}

// Returns whether x, y is above the line from top left to bottom right
func (r *Rect) TL_BR_diag(x, y float64) bool {
	x_0, y_0 := r.Left(), r.Top()
	x_1, y_1 := r.Right(), r.Bottom()
	return (y-y_0)/(y_1-y_0) > (x-x_0)/(x_1-x_0)
}

// Returns whether x, y is above the line from bottom left to top right
func (r *Rect) BL_TR_diag(x, y float64) bool {
	x_0, y_0 := r.Left(), r.Bottom()
	x_1, y_1 := r.Right(), r.Top()
	return (y-y_0)/(y_1-y_0) > (x-x_0)/(x_1-x_0)
}

func (r *Rect) GetInteriorSection(x, y float64) RectInteriorSection {
	if r.TL_BR_diag(x, y) && r.BL_TR_diag(x, y) {
		return LEFT
	} else if !r.TL_BR_diag(x, y) && r.BL_TR_diag(x, y) {
		return TOP
	} else if r.TL_BR_diag(x, y) && !r.BL_TR_diag(x, y) {
		return BOTTOM
	} else if !r.TL_BR_diag(x, y) && !r.BL_TR_diag(x, y) {
		return RIGHT
	}
	return INTERIOR_SECTION_INVALID
}

func (r *Rect) GetExteriorSection(x, y float64) RectExteriorSection {
	if x < r.Left() && y < r.Top() {
		return TOP_LEFT
	} else if x > r.Right() && y < r.Top() {
		return TOP_RIGHT
	} else if x < r.Left() && y > r.Bottom() {
		return BOTTOM_LEFT
	} else if x > r.Right() && y > r.Bottom() {
		return BOTTOM_RIGHT
	}

	isWithinX := IsBetween(r.Left(), r.Right(), x)
	isWithinY := IsBetween(r.Top(), r.Bottom(), y)
	if isWithinX && !isWithinY {
		if y > r.Bottom() {
			return BOTTOM_MIDDLE
		} else if y < r.Top() {
			return TOP_MIDDLE
		}
	} else if isWithinY && !isWithinX {
		if x > r.Right() {
			return MIDDLE_RIGHT
		} else if x < r.Left() {
			return MIDDLE_LEFT
		}
	}
	return EXTERIOR_SECTION_INVALID
}

func NewRect(x, y, width, height float64) *Rect {
	return &Rect{x, y, width, height}
}

// Deprecated: Do this manually instead
func RectFromImage(x, y float64, image *ebiten.Image) *Rect {
	size := image.Bounds().Size()
	return &Rect{
		X:      x,
		Y:      y,
		Width:  float64(size.X),
		Height: float64(size.Y),
	}
}
