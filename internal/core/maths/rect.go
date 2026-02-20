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

// Pixel perfect version. Should replace the usual right?
func (r *Rect) PPBottom() float64 {
	return r.y + r.height - 1
}

func (r *Rect) Center() (float64, float64) {
	return r.x + r.width/2, r.y + r.height/2
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

func (r *Rect) Width() float64 {
	return r.width
}

func (r *Rect) Height() float64 {
	return r.height
}

func (r *Rect) Size() (float64, float64) {
	return r.width, r.height
}

func (r *Rect) HalfSize() (float64, float64) {
	return r.width / 2, r.height / 2
}

func (r *Rect) Extended(dir Direction, length float64) *Rect {
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
	return &newRect
}

func (r *Rect) SetPos(x, y float64) {
	r.x, r.y = x, y
}

func (r *Rect) Translate(x, y float64) {
	r.x += x
	r.y += y
}

func (r *Rect) Translated(x, y float64) Rect {
	newRect := *r
	newRect.x += x
	newRect.y += y
	return newRect
}

func (r *Rect) Extend(x, y float64) {
	r.width += x
	r.height += y
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
		x:      Lerp(r.x, other.x, t),
		y:      Lerp(r.y, other.y, t),
		width:  Lerp(r.width, other.width, t),
		height: Lerp(r.height, other.height, t),
	}
}

// Returns the bounding box for the list of rects
func BB(rects []*Rect) *Rect {
	BBrect := NewRect(0, 0, 0, 0)
	minX := math.Inf(1)
	minY := math.Inf(1)
	for _, rect := range rects {
		if rect.x < minX {
			minX = rect.x
			BBrect.x = rect.x
		}
		if rect.y < minY {
			minY = rect.y
			BBrect.y = rect.y
		}
	}
	maxWidth := math.Inf(-1)
	maxHeight := math.Inf(-1)
	for _, rect := range rects {
		if rect.Right()-minX > maxWidth {
			maxWidth = rect.Right() - minX
			BBrect.width = rect.Right() - minX
		}

		if rect.Bottom()-minY > maxHeight {
			maxHeight = rect.Bottom() - minY
			BBrect.height = rect.Bottom() - minY
		}
	}
	return BBrect
}

// Returns a rect which is half of the original one, according
// to the passed in direction
func (r *Rect) GetHalved(dir Direction) *Rect {
	switch dir {
	case DirUp:
		return NewRect(r.x, r.y, r.width, r.height/2)
	case DirDown:
		return NewRect(r.x, r.y+r.height/2, r.width, r.height/2)
	case DirLeft:
		return NewRect(r.x, r.y, r.width/2, r.height)
	case DirRight:
		return NewRect(r.x+r.width/2, r.y, r.width/2, r.height)
	}
	return NewRect(0, 0, 0, 0)
}

// Checks if a ray starting in posX, posY and travelling in the given direction will
// intersect the rect.
// TODO: I think this should return the intersection point
// rather than the top left corner of the rect (how in the world
// was i able to use this successfully?)
func (r *Rect) RaycastDirectional(posX, posY float64, direction Direction) (bool, float64, float64) {
	if r.Contains(posX, posY) {
		// Should be true?
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
			return Rect{r.x, y, r.width, r.height}
		case LEFT:
			return Rect{x, r.y, r.width, r.height}
		case RIGHT:
			return Rect{x - r.width, r.y, r.width, r.height}
		case BOTTOM:
			return Rect{r.x, y - r.height, r.width, r.height}
		}
	}

	switch r.GetExteriorSection(x, y) {
	case TOP_LEFT:
		return Rect{x, y, r.width, r.height}
	case TOP_MIDDLE:
		return Rect{r.x, y, r.width, r.height}
	case TOP_RIGHT:
		return Rect{x - r.width, y, r.width, r.height}
	case MIDDLE_LEFT:
		return Rect{x, r.y, r.width, r.height}
	case MIDDLE_RIGHT:
		return Rect{x - r.width, r.y, r.width, r.height}
	case BOTTOM_LEFT:
		return Rect{x, y - r.height, r.width, r.height}
	case BOTTOM_MIDDLE:
		return Rect{r.x, y - r.height, r.width, r.height}
	case BOTTOM_RIGHT:
		return Rect{x - r.width, y - r.height, r.width, r.height}
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

func RectFromImage(x, y float64, image *ebiten.Image) *Rect {
	size := image.Bounds().Size()
	return &Rect{
		x:      x,
		y:      y,
		width:  float64(size.X),
		height: float64(size.Y),
	}
}
