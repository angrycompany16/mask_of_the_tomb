package slambox

import (
	"errors"
	"image/color"
	"log"
	"mask_of_the_tomb/internal/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type testApp struct {
}

func (t *testApp) Update() error {
	return nil
}

func (t *testApp) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
}

func (t *testApp) Layout(outsideHeight, outsideWidth int) (int, int) {
	return GAME_WIDTH * PIXEL_SCALE, GAME_HEIGHT * PIXEL_SCALE
}

// Sets up a simple game loop for testing this package.
func RunTestEnv() {
	a := &testApp{}

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		} else if err == scenes.ErrTerminated {
			return
		}
		log.Fatal(err)
	}
}
