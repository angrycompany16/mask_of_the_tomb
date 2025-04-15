package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/internal/game"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/gamestate"
	"mask_of_the_tomb/internal/game/world"
	"path/filepath"

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
	rendering.RenderLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GameWidth * rendering.PixelScale, rendering.GameHeight * rendering.PixelScale
}

func main() {
	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Slambox test")

	a := &App{game.NewGame()}
	world.LDTKMapPath = filepath.Join("assets", "LDTK", "slambox.ldtk")
	a.game.Init()

	ebiten.SetFullscreen(true)
	a.game.State = gamestate.Playing
	a.game.EnterPlayMode()

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
