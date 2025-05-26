package resources

var Time float64
var GrassWindSeed int64
var State GameState

type GameState int

const (
	Loading GameState = iota
	MainMenu
	Intro
	Playing
	Paused
)
