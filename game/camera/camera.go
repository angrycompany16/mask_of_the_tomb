package camera

import (
	"mask_of_the_tomb/rect"
	"mask_of_the_tomb/utils"
)

type Camera struct {
	posX, posY            float64
	borders               rect.Rect
	halfWidth, halfHeight float64
}

func (c *Camera) Init(borders rect.Rect, halfWidth, halfHeight float64) {
	c.borders = borders
	c.halfWidth = halfWidth
	c.halfHeight = halfHeight
}

func (c *Camera) Update() {

}

func (c *Camera) GetX() float64 {
	return c.posY
}

func (c *Camera) GetY() float64 {
	return c.posY
}

func (c *Camera) SetPosition(x, y float64) {
	c.posX = utils.Clamp(c.posX, c.borders.Left()+c.halfWidth, c.borders.Right()-c.halfWidth)
	c.posY = utils.Clamp(c.posY, c.borders.Top()+c.halfHeight, c.borders.Bottom()-c.halfHeight)
}

func NewCamera() *Camera {
	return &Camera{}
}
