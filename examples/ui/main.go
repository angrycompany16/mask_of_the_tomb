package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/internal/game"
	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/UI/display"
	"mask_of_the_tomb/internal/game/UI/node/root"
	"mask_of_the_tomb/internal/game/UI/node/textbox"
	"mask_of_the_tomb/internal/game/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	ui *ui.UI
}

func (a *App) Update() error {
	a.ui.Update()
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.ui.Draw()
	rendering.RenderLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GameWidth * rendering.PixelScale, rendering.GameHeight * rendering.PixelScale
}

func main() {
	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("UI test")

	app := &App{ui.NewUI()}

	rootNode := root.New()
	rootNode.AddChild(textbox.New())
	app.ui.AddDisplayManual(&display.Display{
		Name: "Test display",
		Root: rootNode,
	})

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
