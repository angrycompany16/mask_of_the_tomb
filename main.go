package main

import (
	"errors"
	"flag"
	"log"
	"mask_of_the_tomb/internal/libraries/rendering"
	"mask_of_the_tomb/internal/transformers/game"

	"github.com/hajimehoshi/ebiten/v2"
)

// Q: Do we want to add everything in main?
// A: No, probably not. However, it would be nice if we could create these
// kinds of "bundles" for grouping together components that are frequently used
// together. Adding the components separately would be very nice though

// The ideal main.go
// Add these entity bundles
// Connect interactions like this

// Another thing to consider:
// It might be worth to modularize in terms of different stages of the game,
// so that one can for example easily create a game without a menu stage,
// without an intro stage of without a loading stage.

// The solution to it all?
// Separate modules into their own self-contained units, and then write
// separate interaction modules between these

// OR: structure into plugins and libraries
// plugins utilize libraries
// Libraries provide core functionality through public methods

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
