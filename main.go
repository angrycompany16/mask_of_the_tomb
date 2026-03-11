package main

import (
	_ "embed"
	"errors"
	"flag"
	"log"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game"
	"mask_of_the_tomb/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GAME_WIDTH, GAME_HEIGHT = 480, 270
	PIXEL_SCALE             = 4
)

type App struct {
	game *engine.Game
}

func (a *App) Update() error {
	return a.game.Update()
}

func (a *App) Draw(screen *ebiten.Image) {
	a.game.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return GAME_WIDTH * PIXEL_SCALE, GAME_HEIGHT * PIXEL_SCALE
}

func main() {
	// flag.BoolVar(&resources.DebugMode, "debug", false, "enable debug mode")
	// flag.StringVar(&scenes.InitLevelName, "initlevel", "", "Level in which to spawn the player")
	// flag.IntVar(&scenes.SaveProfile, "saveprofile", 1, "Profile to use for saving/loading (99 for dev save)")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file (.prof ending)")

	flag.Parse()

	if *cpuprofile != "" {
		stopProfiling := utils.StartProfiling(cpuprofile)
		defer stopProfiling()
	}

	ebiten.SetWindowSize(GAME_WIDTH*PIXEL_SCALE, GAME_HEIGHT*PIXEL_SCALE)
	ebiten.SetWindowTitle("Mask of the tomb")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	a := &App{
		game.CreateGame(GAME_WIDTH, GAME_HEIGHT, PIXEL_SCALE),
	}

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		} else if err == engine.ErrTerminated {
			return
		}
		log.Fatal(err)
	}
}
