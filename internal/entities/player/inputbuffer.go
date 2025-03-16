package player

import "mask_of_the_tomb/internal/maths"

type inputBuffer struct {
	direction maths.Direction
	t         float64
	duration  float64
}

func (ib *inputBuffer) update() {
	ib.t -= 0.01666666667 // Update tick time should be constant

	if ib.t <= 0 {
		ib.direction = maths.DirNone
	}
}

func (ib *inputBuffer) Read() maths.Direction {
	return ib.direction
}

func (ib *inputBuffer) set(direction maths.Direction) {
	ib.direction = direction
	ib.t = inputBufferDuration
}

func (ib *inputBuffer) Clear() {
	ib.direction = maths.DirNone
}

func newInputBuffer(duration float64) inputBuffer {
	return inputBuffer{
		direction: maths.DirNone,
		t:         0,
		duration:  duration,
	}
}
