package main

import (
	"errors"
	"flag"
	"log"
	"mask_of_the_tomb/internal/core/profiling"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	game *scenes.Game
}

func (a *App) Update() error {
	err := a.game.Update()
	return err
}

func (a *App) Draw(screen *ebiten.Image) {
	a.game.Draw()
	rendering.ScreenLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GAME_WIDTH * rendering.PIXEL_SCALE, rendering.GAME_HEIGHT * rendering.PIXEL_SCALE
}

func main() {
	flag.BoolVar(&resources.DebugMode, "debug", false, "enable debug mode")
	flag.StringVar(&scenes.InitLevelName, "initlevel", "", "Level in which to spawn the player")
	flag.IntVar(&scenes.SaveProfile, "saveprofile", 1, "Profile to use for saving/loading (99 for dev save)")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file (.prof ending)")

	flag.Parse()

	if *cpuprofile != "" {
		stopProfiling := profiling.StartProfiling(cpuprofile)
		defer stopProfiling()
	}

	ebiten.SetWindowSize(rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)
	ebiten.SetWindowTitle("Mask of the tomb")

	a := &App{scenes.NewGame()}
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		} else if err == scenes.ErrTerminated {
			return
		}
		log.Fatal(err)
	}
}
