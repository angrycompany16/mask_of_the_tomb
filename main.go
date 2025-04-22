package main

import (
	"errors"
	"flag"
	"log"
	"mask_of_the_tomb/internal/game"
	"mask_of_the_tomb/internal/game/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	debugMode bool
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
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")
	flag.StringVar(&game.InitLevelName, "initlevel", "", "Level in which to spawn the player")
	flag.IntVar(&game.SaveProfile, "saveprofile", 1, "Profile to use for saving/loading (99 for dev save)")

	flag.Parse()

	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb")

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
