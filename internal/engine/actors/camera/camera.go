package camera

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
)

var (
	shakePaddingX, shakePaddingY = 20.0, 20
)

// Control renderer.camera during update. This means that multiple cameras will
// break the game, but that's alright for now.

type Camera struct {
	transform2D.Transform2D
	width, height                  float64
	offsetX, offsetY               float64
	screenPaddingY, screenPaddingX float64
	shakeOffsetX, shakeOffsetY     float64
	shaking                        bool
	shakeDuration                  float64
	shakeStrength                  float64
	shakeTime                      float64
	damping                        float64
}

func (c *Camera) Init(width, height, offsetX, offsetY float64) {
	c.SetBorders(width, height)
	c.offsetX, c.offsetY = offsetX, offsetY
}

func (c *Camera) Update() {
	if !c.shaking {
		return
	}

	c.shakeTime += 1.0 / 60.0
	if c.shakeTime > c.shakeDuration {
		c.shaking = false
		c.shakeTime = 0
		c.shakeOffsetX = 0
		c.shakeOffsetY = 0
		return
	}

	c.shakeStrength = maths.Lerp(c.shakeStrength, 0, 0.05*c.damping)
	c.shakeOffsetX = maths.RandomRange(-c.shakeStrength, c.shakeStrength)
	c.shakeOffsetY = maths.RandomRange(-c.shakeStrength, c.shakeStrength)
}

// returns the position of the camera
func (c *Camera) GetPos(includeShake bool) (float64, float64) {
	transformX, transformY := c.Transform2D.GetPos(false)
	// Special adjustment for LDtk levels... probably not that great
	if c.height == 272 {
		transformY += 1
	}

	if includeShake {
		return transformX + c.shakeOffsetX, transformY + c.shakeOffsetY
	}
	return transformX, transformY
}

func (c *Camera) GetShake() (float64, float64) {
	return c.shakeOffsetX, c.shakeOffsetY
}

func (c *Camera) SetPos(x, y float64) {
	if c.height == 272 {
		c.screenPaddingY = 2
	} else {
		c.screenPaddingY = 0
	}
	// Need to fix
	// c.posX = maths.Clamp(x-c.offsetX, c.screenPaddingX/2, c.width-rendering.GAME_WIDTH-c.screenPaddingX/2)
	// c.posY = maths.Clamp(y-c.offsetY, c.screenPaddingY/2, c.height-rendering.GAME_HEIGHT-c.screenPaddingY/2)
}

func (c *Camera) SetPadding(paddingX, paddingY float64) {
	c.screenPaddingX = paddingX
	c.screenPaddingY = paddingY
}

func (c *Camera) SetBorders(width, height float64) {
	c.width = width
	c.height = height
}

func (c *Camera) Shake(duration, strength, damping float64) {
	c.shakeTime = 0
	c.shakeDuration = duration
	c.shakeStrength = strength
	c.shaking = true
	c.damping = damping
}
