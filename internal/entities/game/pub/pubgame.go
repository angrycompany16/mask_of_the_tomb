package pubgame

import "mask_of_the_tomb/internal/engine/events"

const (
	GameEntityName = "Game"
)

type GameState int

const (
	StateMainMenu = iota
	StatePlaying
	StatePaused
)

var (
	StateChanged = events.NewEvent()
)

type GameAdvertiser struct {
	State GameState
}

func (gadv *GameAdvertiser) Read() any {
	return *gadv
}
