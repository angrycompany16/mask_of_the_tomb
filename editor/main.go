package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/editor/editor"
	"mask_of_the_tomb/internal/game"
	"mask_of_the_tomb/internal/game/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	editor *editor.Editor
}

func (a *App) Update() error {
	err := a.editor.Update()
	if err == editor.ErrTerminated {
		return err
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.editor.Draw()
	rendering.RenderLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GameWidth * rendering.PixelScale, rendering.GameHeight * rendering.PixelScale
}

func main() {
	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb editor")

	a := &App{editor.NewEditor()}
	a.editor.Init()

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			// Save the data
			return
		}
		log.Fatal(err)
	}
}
