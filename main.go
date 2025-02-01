package main

import (
	"flag"
	"log"

	"mask_of_the_tomb/commons"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/game"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	debugMode bool
)

type App struct {
	world *game.Game
}

func (a *App) Init() {
	a.world.Init()
}

func (a *App) Update() error {
	a.world.Update()
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	DrawAt(a.world.Draw(), screen, 0, 0)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return commons.GameWidth * commons.PixelScale, commons.GameHeight * commons.PixelScale
}

func main() {
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")

	flag.Parse()

	ebiten.SetWindowSize(commons.GameWidth*commons.PixelScale, commons.GameHeight*commons.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb")

	if debugMode {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	} else {
		ebiten.SetFullscreen(true)
	}

	a := &App{game.NewGame()}
	a.Init()

	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}
