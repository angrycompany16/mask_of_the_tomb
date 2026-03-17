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
	// Special adjustment for LDtk levels... probably not that great
	if c.height == 272 {
		y += 1
	}

	camX, camY := c.Transform2D.GetPos(false)
	if includeShake {
		return x - (camX + c.shakeOffsetX), y - (camY + c.shakeOffsetY)
	}
	return x - camX, y - camY
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

// 20, 20 are good default value for shake padding
func NewCamera(transform2D *transform2D.Transform2D, width, height, offsetX, offsetY, screenPaddingX, screenPaddingY float64) *Camera {
	return &Camera{
		Transform2D:   transform2D,
		width:         width,
		height:        height,
		offsetX:       offsetX,
		offsetY:       offsetY,
		screenMarginX: screenPaddingX,
		screenMarginY: screenPaddingY,
	}
}
