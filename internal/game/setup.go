package game

import (
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

func CreateGame(gw, gh, ps int) *engine.Game {
	game := engine.NewGame(engine.NewCommands(
		engine.WithRenderer(gw, gh, ps, true, true),
	))

	inputServer := game.GetCmd().InputHandler()
	inputServer.RegisterAction("toggleInspector", input.KeyJustPressedAction(ebiten.KeyTab))
	game.RegisterScene("loading", scenes.LoadingScene)

	game.SpawnScene("loading")
	return game
}
