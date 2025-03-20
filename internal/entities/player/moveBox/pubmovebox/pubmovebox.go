package pubmovebox

import (
	"mask_of_the_tomb/internal/engine/events"
)

type PositionAdvertiser struct {
	PosX, PosY float64
}

func (padv *PositionAdvertiser) Read() any {
	return *padv
}

var (
	FinishedMoveEvent = events.NewEvent()
)
