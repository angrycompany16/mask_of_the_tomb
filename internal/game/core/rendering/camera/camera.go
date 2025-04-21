package camera

import (
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/maths"
)

// TODO: Solve the single-screen height bug better

var (
	_camera = &Camera{}
)

type Camera struct {
	posX, posY       float64
	width, height    float64
	offsetX, offsetY float64
}

func Init(width, height, offsetX, offsetY float64) {
	SetBorders(width, height)
	_camera.offsetX, _camera.offsetY = offsetX, offsetY
}

func GetPos() (float64, float64) {
	if _camera.height == 272 {
		return _camera.posX, _camera.posY + 1
	}
	return _camera.posX, _camera.posY
}

func SetPos(x, y float64) {
	if _camera.height == 272 {
		return
	}
	_camera.posX = maths.Clamp(x-_camera.offsetX, 0, _camera.width-rendering.GameWidth)
	_camera.posY = maths.Clamp(y-_camera.offsetY, 0, _camera.height-rendering.GameHeight)
}

func SetBorders(width, height float64) {
	_camera.width = width
	_camera.height = height
}
