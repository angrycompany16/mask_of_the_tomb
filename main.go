// package main

// import (
// 	"log"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// )

// type Game struct{}

// func (g *Game) Update() error {
// 	return nil
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// 	ebitenutil.DebugPrint(screen, "Hello, World!")
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
// 	return 320, 240
// }

// func main() {
// 	ebiten.SetWindowSize(640, 480)
// 	ebiten.SetWindowTitle("Hello, World!")
// 	if err := ebiten.RunGame(&Game{}); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"errors"
	"flag"
	"log"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/transformers/game"

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
	rendering.ScreenLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GAME_WIDTH * rendering.PIXEL_SCALE, rendering.GAME_HEIGHT * rendering.PIXEL_SCALE
}

func main() {
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")
	flag.StringVar(&game.InitLevelName, "initlevel", "", "Level in which to spawn the player")
	flag.IntVar(&game.SaveProfile, "saveprofile", 1, "Profile to use for saving/loading (99 for dev save)")

	flag.Parse()

	ebiten.SetWindowSize(rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)
	ebiten.SetWindowTitle("Mask of the tomb")

	a := &App{game.NewGame()}
	a.game.Load()

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		}
		log.Fatal(err)
	}
}
