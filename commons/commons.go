package commons

import "errors"

const (
	GameWidth, GameHeight = 480, 270
	PixelScale            = 4
)

var (
	Terminated = errors.New("Terminated")
)
