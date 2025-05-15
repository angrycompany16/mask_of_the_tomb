package gamestate

type GameState struct {
	S        State
	GameTime float64
}

type State int

const (
	Loading State = iota
	MainMenu
	Intro
	Playing
	Paused
)

// Would be nice to define more state variables if possible
