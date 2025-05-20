package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/examples/slamboxes/game"
	"mask_of_the_tomb/internal/libraries/rendering"
	"mask_of_the_tomb/internal/plugins/world"
	"path/filepath"

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
	return rendering.GameWidth * rendering.PixelScale, rendering.GameHeight * rendering.PixelScale
}

func main() {
	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Slambox test")

	a := &App{game.NewGame()}
	world.LDTKMapPath = filepath.Join("assets", "LDTK", "slambox.ldtk")
	a.game.Load()

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
