package pubui

import "mask_of_the_tomb/internal/engine/events"

type UISelect int

const (
	SelectPlay UISelect = iota
	SelectOpts
	SelectQuit
	SelectMainMenu
)

var (
	UISelected = events.NewEvent()
)
