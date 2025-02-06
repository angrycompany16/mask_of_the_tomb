package utils

import (
	"errors"
	"image/color"
)

const ()

var (
	GlobalTimeScale = 1.0
	Terminated      = errors.New("Terminated")
	TextColorNormal = ColorPair{
		BrightColor: color.RGBA{205, 247, 226, 255},
		DarkColor:   color.RGBA{199, 176, 139, 255},
	}
	TextColorSelected = ColorPair{
		BrightColor: color.RGBA{255, 255, 255, 255},
		DarkColor:   color.RGBA{0, 0, 0, 255},
	}
)

// Global state
type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
)

var (
	GlobalState GameState
)
