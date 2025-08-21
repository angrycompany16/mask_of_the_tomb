package scenes

import (
	"errors"
	"mask_of_the_tomb/internal/core/scene"
	"time"
)

// IMPORTANT NOTE: If there is a missing entity reference, the game world will not be able to load !!!
// Instead it becomes a nil pointer...

// Some problems:
// - When exiting to main menu, music volume remains low
//   - This follows from the game not playing the proper song on pressing play
// - When opening menus for the first time, the select sound plays
// - Somehow it seems that the game itself gets darker if we don't draw the UI???

var (
	ErrTerminated = errors.New("Terminatednow")

	// These should all be removed
	InitLevelName string
	SaveProfile   int
	initTime      time.Time
)

type Game struct {
	sceneStack *scene.SceneStack
}

func (g *Game) Update() error {
	quit := g.sceneStack.Update()
	if quit {
		return ErrTerminated
	}
	return nil
}

func (g *Game) Draw() {
	g.sceneStack.Draw()
}

func NewGame() *Game {
	game := &Game{
		sceneStack: scene.NewSceneStack(),
	}

	game.sceneStack.Push(MakeLoadingScene())
	return game
}
