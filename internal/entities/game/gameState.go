package game

type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
)
