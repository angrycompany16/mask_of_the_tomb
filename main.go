package main

import (
	"flag"
	"log"

	"mask_of_the_tomb/game"
	"mask_of_the_tomb/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	debugMode bool
)

type Game struct {
	world *game.World
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	utils.DrawAt(g.world.Draw(), screen, 0, 0)
}

func (g *Game) Layout(outsideHeight, outsideWidth int) (int, int) {
	return game.GameWidth * game.PixelScale, game.GameHeight * game.PixelScale
}

func main() {
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")

	flag.Parse()

	ebiten.SetWindowSize(game.GameWidth*game.PixelScale, game.GameHeight*game.PixelScale)
	ebiten.SetWindowTitle("Mask of the tomb")

	if debugMode {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}

	if err := ebiten.RunGame(&Game{game.MakeWorld()}); err != nil {
		log.Fatal(err)
	}
}
