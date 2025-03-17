package main

import (
	"errors"
	"flag"
	"log"
	"mask_of_the_tomb/internal/app"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/rendering"
	save "mask_of_the_tomb/internal/libraries/savesystem"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	debugMode bool
)

func main() {
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")

	flag.Parse()

	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb")

	a := app.MakeApp()
	a.Init()

	if debugMode {
		// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		// game.State = game.StatePlaying
		// a.game.EnterPlayMode()
	} else {
		ebiten.SetFullscreen(true)
	}

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errs.ErrTerminated) {
			save.GlobalSave.SaveGame()
			return
		}
		log.Fatal(err)
	}
}
