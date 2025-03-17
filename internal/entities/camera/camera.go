package camera

import (
	"mask_of_the_tomb/internal/engine/advertisers"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/entities/camera/pubcamera"
	"mask_of_the_tomb/internal/libraries/maths"
	"mask_of_the_tomb/internal/libraries/rendering"
)

type camera struct {
	posX, posY       float64
	width, height    float64
	offsetX, offsetY float64
	cameraAdvertiser pubcamera.CameraAdvertiser
}

func New(width, height, offsetX, offsetY float64) *camera {
	_camera := camera{
		width:   width,
		height:  height,
		offsetX: offsetX,
		offsetY: offsetY,
	}
	entities.RegisterEntity(&_camera, pubcamera.CameraEntityName)
	advertisers.RegisterAdvertiser(&_camera.cameraAdvertiser, pubcamera.CameraEntityName)

	return &_camera
}

func (c *camera) Update() {

}

func (c *camera) PostUpdate() {
	c.cameraAdvertiser.PosX = c.posX
	c.cameraAdvertiser.PosY = c.posY
}

// func (c *camera) GetPos() (float64, float64) {
// 	return c.posX, c.posY
// }

func (c *camera) SetPos(x, y float64) {
	c.posX = maths.Clamp(x-c.offsetX, 0, c.width-rendering.GameWidth)
	c.posY = maths.Clamp(y-c.offsetY, 0, c.height-rendering.GameHeight)
}

func (c *camera) SetBorders(width, height float64) {
	c.width = width
	c.height = height
}
