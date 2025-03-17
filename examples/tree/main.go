package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/examples/tree/app"
	"mask_of_the_tomb/internal/entities/world"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/rendering"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb")

	app := app.MakeApp()
	app.Init()

	// Modify the data of world entity to use a different level - Note that this should never
	// be done in practice!
	world.LDTKMapPath = filepath.Join("assets", "LDTK", "tree.ldtk")

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, errs.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
