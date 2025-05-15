package main

import (
	"errors"
	"log"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/libraries/rendering"
	ui "mask_of_the_tomb/internal/plugins/UI"
	"mask_of_the_tomb/internal/transformers/game"
	"path/filepath"

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
	assetloader.LoadPreamble()
	// rootNode := root.New()
	// rootNode.AddChild(textbox.New())
	app.ui.AddDisplayManual(errs.Must(
		ui.FromFile(filepath.Join("assets", "menus", "example", "menu.yaml")),
	))
	// filepath.Join("assets", "menus", "example", "menu.yaml") |> display.FromFile |> errs.Must |> app.ui.AddDisplayManual

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, game.ErrTerminated) {
			return
		}
		log.Fatal(err)
	}
}
