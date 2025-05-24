package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/examples/slamboxes/game"
	"mask_of_the_tomb/internal/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Add some kind of staging so that we can skip for instance the cutscene

type App struct {
	game *game.Game
}

func (a *App) Update() error {
	err := a.game.Update()
	if err == game.ErrTerminated {
		return err
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.game.Draw()
	rendering.ScreenLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GAME_WIDTH * rendering.PIXEL_SCALE, rendering.GAME_HEIGHT * rendering.PIXEL_SCALE
}

func main() {
	ebiten.SetWindowSize(rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)
	ebiten.SetWindowTitle("Slambox test")

	a := &App{game.NewGame()}
	a.game.Load()

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
