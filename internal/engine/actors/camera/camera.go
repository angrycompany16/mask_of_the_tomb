package camera

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
)

type Camera struct {
	*transform2D.Transform2D
	width, height float64
	// Represents the offset for the camera tracking.
	// For instance, this is used to center the camera on the player
	offsetX, offsetY float64
	// This is used to apply margins to the camera position
	screenMarginY, screenMarginX float64

	// For camera shake
	shakeOffsetX, shakeOffsetY float64
	shaking                    bool
	shakeDuration              float64
	shakeStrength              float64
	shakeTime                  float64
	damping                    float64
}

func (c *Camera) Update(cmd *engine.Commands) {
	c.Transform2D.Update(cmd)
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
func (c *Camera) WorldToCam(x, y float64, includeShake bool) (float64, float64) {
	camX, camY := c.Transform2D.GetPos(false)
	if includeShake {
		return x - (camX + c.shakeOffsetX - c.offsetX), y - (camY + c.shakeOffsetY - c.offsetY)
	}
	return x - (camX - c.offsetX), y - (camY - c.offsetY)
}

func (c *Camera) WorldToCamCustomOffset(x, y float64, offsetX, offsetY float64, includeShake bool) (float64, float64) {
	camX, camY := c.Transform2D.GetPos(false)
	if includeShake {
		return x - (camX + c.shakeOffsetX - offsetX), y - (camY + c.shakeOffsetY - offsetY)
	}
	return x - (camX - offsetX), y - (camY - offsetY)
}

func (c *Camera) GetShake() (float64, float64) {
	return c.shakeOffsetX, c.shakeOffsetY
}

// func (c *Camera) SetPos(x, y, gameWidth, gameHeight float64) {
// 	if c.height == 272 {
// 		c.screenMarginY = 2
// 	} else {
// 		c.screenMarginY = 0
// 	}

// 	c.x = maths.Clamp(x-c.offsetX, c.screenMarginX/2, c.width-gameWidth-c.screenMarginX/2)
// 	c.y = maths.Clamp(y-c.offsetY, c.screenMarginY/2, c.height-gameHeight-c.screenMarginY/2)
// }

func (c *Camera) Shake(duration, strength, damping float64) {
	c.shakeTime = 0
	c.shakeDuration = duration
	c.shakeStrength = strength
	c.shaking = true
	c.damping = damping
}

type Option func(*Camera)

func newDefaultCamera() *Camera {
	return &Camera{
		width:         480,
		height:        270,
		offsetX:       240,
		offsetY:       135,
		screenMarginX: 0,
		screenMarginY: 0,
	}
}

// 20, 20 are good default value for shake padding
func NewCamera(transform2D *transform2D.Transform2D, options ...Option) *Camera {
	camera := &Camera{
		Transform2D: transform2D,
	}

	for _, option := range options {
		option(camera)
	}

	return camera
}

func WithSize(width, height float64) Option {
	return func(c *Camera) {
		c.width = width
		c.height = height
	}
}

func WithOffset(offsetX, offsetY float64) Option {
	return func(c *Camera) {
		c.offsetX = offsetX
		c.offsetY = offsetY
	}
}

func WithMargins(marginX, marginY float64) Option {
	return func(c *Camera) {
		c.screenMarginX = marginX
		c.screenMarginY = marginY
	}
}
