package game

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/scenes"
)

func CreateGame(gw, gh, ps int) *engine.Game {
	game := engine.NewGame(engine.NewCommands(
		engine.WithRenderer(gw, gh, ps, true, true),
	))

	game.RegisterScene("loading", scenes.LoadingScene)

	game.SpawnScene("loading")
	return game
}
