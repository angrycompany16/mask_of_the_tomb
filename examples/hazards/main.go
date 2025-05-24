package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/internal/libraries/camera"
	"mask_of_the_tomb/internal/transformers/game"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	camera.ScreenLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return camera.GameWidth * camera.PixelScale, camera.GameHeight * camera.PixelScale
}

func main() {
	ebiten.SetWindowSize(camera.GameWidth*camera.PixelScale, camera.GameHeight*camera.PixelScale)
	ebiten.SetWindowTitle("Hazard test")

	a := &App{game.NewGame()}
	// world.LDTKMapPath = filepath.Join("assets", "LDTK", "hazard.ldtk")
	a.game.Load()

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
