package app

import (
	"fmt"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/engine/events"
	pubui "mask_of_the_tomb/internal/entities/UI/pub"
	"mask_of_the_tomb/internal/entities/game"
	pubgame "mask_of_the_tomb/internal/entities/game/pub"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

// Connects engine and game
type App struct {
	UISelectListener *events.EventListener
}

func (a *App) Init() {
	// .. Init base entities
	fmt.Println("App init")
	events.InitEventManager()

	_game := game.NewGame()
	entities.RegisterEntity(_game, pubgame.GameEntityName)
}

func (a *App) Update() error {
	entities.PreUpdate()

	if info, raised := a.UISelectListener.Poll(); raised {
		if selectInfo, ok := info.Data.(pubui.UISelect); ok {
			if selectInfo == pubui.SelectQuit {
				return errs.ErrTerminated
			}
		}
	}

	entities.Update()
	entities.PostUpdate()

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	rendering.RenderLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GameWidth * rendering.PixelScale, rendering.GameHeight * rendering.PixelScale
}

func MakeApp() *App {
	return &App{
		UISelectListener: events.NewEventListener(pubui.UISelected),
	}
}
