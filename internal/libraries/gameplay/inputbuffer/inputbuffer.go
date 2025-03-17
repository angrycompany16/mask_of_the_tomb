package inputbuffer

import "mask_of_the_tomb/internal/libraries/maths"

type InputBuffer struct {
	direction maths.Direction
	t         float64
	duration  float64
}

func (buf *InputBuffer) Update() {
	buf.t -= 0.01666666667 // Update tick time should be constant

	if buf.t <= 0 {
		buf.direction = maths.DirNone
	}
}

func (buf *InputBuffer) Read() maths.Direction {
	return buf.direction
}

func (buf *InputBuffer) ReadSingle() maths.Direction {
	direction := buf.direction
	buf.Clear()
	return direction
}

func (buf *InputBuffer) Set(direction maths.Direction) {
	buf.direction = direction
	buf.t = buf.duration
}

func (buf *InputBuffer) Clear() {
	buf.direction = maths.DirNone
}

func NewInputBuffer(duration float64) InputBuffer {
	return InputBuffer{
		direction: maths.DirNone,
		t:         0,
		duration:  duration,
	}
}
