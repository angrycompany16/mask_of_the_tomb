package pubplayer

import (
	"mask_of_the_tomb/internal/engine/events"
	"mask_of_the_tomb/internal/libraries/maths"
)

var OnMove = events.NewEvent()

type MoveEvent struct {
	Direction maths.Direction
	Hitbox    maths.Rect
}
