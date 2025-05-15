package inputbuffer

import "mask_of_the_tomb/internal/core/maths"

type InputBuffer struct {
	direction maths.Direction
	t         float64
	duration  float64
}

func (ib *InputBuffer) Update() {
	ib.t -= 0.01666666667 // Update tick time should be constant

	if ib.t <= 0 {
		ib.direction = maths.DirNone
	}
}

func (ib *InputBuffer) Read() maths.Direction {
	return ib.direction
}

func (ib *InputBuffer) Set(direction maths.Direction) {
	ib.direction = direction
	ib.t = ib.duration
}

func (ib *InputBuffer) Clear() {
	ib.direction = maths.DirNone
}

func NewInputBuffer(duration float64) InputBuffer {
	return InputBuffer{
		direction: maths.DirNone,
		t:         0,
		duration:  duration,
	}
}
