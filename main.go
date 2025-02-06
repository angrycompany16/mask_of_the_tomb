package main

import (
	"errors"
	"flag"
	"log"

	"mask_of_the_tomb/game"
	"mask_of_the_tomb/game/save"
	"mask_of_the_tomb/rendering"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	debugMode bool
)

type App struct {
	game *game.Game
}

func (a *App) Init() {
	a.game.Init()
}

func (a *App) Update() error {
	err := a.game.Update()
	if err == utils.Terminated {
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

	flag.Parse()

	ebiten.SetWindowSize(rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb")

	if debugMode {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	} else {
		ebiten.SetFullscreen(true)
	}

	a := &App{game.NewGame()}
	a.Init()

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, utils.Terminated) {
			save.GlobalSave.SaveGame()
			return
		}
		log.Fatal(err)
	}
}
