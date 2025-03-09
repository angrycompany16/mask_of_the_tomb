package camera

import (
	"mask_of_the_tomb/internal/game/rendering"
	"mask_of_the_tomb/internal/maths"
)

var (
	GlobalCamera = NewCamera()
)

type Camera struct {
	posX, posY       float64
	width, height    float64
	offsetX, offsetY float64
}

func (c *Camera) Init(width, height, offsetX, offsetY float64) {
	c.width, c.height = width, height
	c.offsetX, c.offsetY = offsetX, offsetY
}

func (c *Camera) Update() {

}

func (c *Camera) GetPos() (float64, float64) {
	return c.posX, c.posY
}

func (c *Camera) SetPos(x, y float64) {
	c.posX = maths.Clamp(x-c.offsetX, 0, c.width-rendering.GameWidth)
	c.posY = maths.Clamp(y-c.offsetY, 0, c.height-rendering.GameHeight)
}

func (c *Camera) SetBorders(width, height float64) {
	c.width = width
	c.height = height
}

func NewCamera() *Camera {
	return &Camera{}
}
