package gamestate

type State int

const (
	Loading State = iota
	MainMenu
	Intro
	Playing
	Paused
)
