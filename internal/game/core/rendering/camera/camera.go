package camera

import (
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/maths"
)

// TODO: Solve the single-screen height bug better

var (
	_camera                      = &Camera{}
	shakePaddingX, shakePaddingY = 20.0, 20
)

type Camera struct {
	posX, posY                 float64
	width, height              float64
	offsetX, offsetY           float64
	screenPadding              float64
	shakeOffsetX, shakeOffsetY float64
	shaking                    bool
	shakeDuration              float64
	shakeStrength              float64
	shakeTime                  float64
	damping                    float64
}

func Init(width, height, offsetX, offsetY float64) {
	SetBorders(width, height)
	_camera.offsetX, _camera.offsetY = offsetX, offsetY
}

func Update() {
	if !_camera.shaking {
		return
	}

	_camera.shakeTime += 1.0 / 60.0
	if _camera.shakeTime > _camera.shakeDuration {
		_camera.shaking = false
		_camera.shakeTime = 0
		_camera.shakeOffsetX = 0
		_camera.shakeOffsetY = 0
		return
	}

	_camera.shakeStrength = maths.Lerp(_camera.shakeStrength, 0, 0.05*_camera.damping)
	_camera.shakeOffsetX = maths.RandomRange(-_camera.shakeStrength, _camera.shakeStrength)
	_camera.shakeOffsetY = maths.RandomRange(-_camera.shakeStrength, _camera.shakeStrength)
}

func GetPos() (float64, float64) {
	if _camera.height == 272 {
		return _camera.posX, _camera.posY + 1
	}
	return _camera.posX, _camera.posY
}

func SetPos(x, y float64) {
	if _camera.height == 272 {
		_camera.screenPadding = 2
	} else {
		_camera.screenPadding = 0
	}
	_camera.posX = maths.Clamp(x-_camera.offsetX, 0, _camera.width-rendering.GameWidth)
	_camera.posY = maths.Clamp(y-_camera.offsetY, _camera.screenPadding/2, _camera.height-rendering.GameHeight-_camera.screenPadding/2)
	_camera.posX += _camera.shakeOffsetX
	_camera.posY += _camera.shakeOffsetY
}

func SetBorders(width, height float64) {
	_camera.width = width
	_camera.height = height
}

func Shake(duration, strength, damping float64) {
	_camera.shakeTime = 0
	_camera.shakeDuration = duration
	_camera.shakeStrength = strength
	_camera.shaking = true
	_camera.damping = damping
}
