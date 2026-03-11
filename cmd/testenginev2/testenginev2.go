package main

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/game/actors/demo"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	gw, gh = 960, 540
	ps     = 1
)

type App struct {
	game   *engine.Game
	toggle bool
}

// We can quite honestly get started at this point
func (a *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		a.toggle = !a.toggle
		if a.toggle {
			a.game.SpawnScene("TestScene2")
		} else {
			a.game.SpawnScene("TestScene1")
		}
	}
	return a.game.Update()
}

func (a *App) Draw(screen *ebiten.Image) {
	a.game.Draw(screen)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// TODO: Rewrite the scene system so that we operate with definitions
// + instances instead of how it is now, since at the current moment
// scenes "pick up where they left off" whenever they are reinstantiated
func main() {
	fmt.Println("START")

	game := engine.NewGame(engine.NewServers(
		engine.ServerArgs{
			GameWidth:  gw,
			GameHeight: gh,
			PixelScale: ps,
		},
	))
	game.RegisterScene("TestScene1", CreateTestScene1)
	game.RegisterScene("TestScene2", CreateTestScene2)

	// Kinda cursed but this works?
	game.SpawnScene("TestScene1")

	app := &App{
		game: game,
	}

	ebiten.SetWindowSize(gw, gh)

	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}

// Create various scenes
// YES
// note: Value receivers instead of pointers!
// But wait: How do we create scenes from LDtk?
// Probably want to use the function style anyway actually
// programmatic becomes too restrictive
func CreateTestScene1(servers *engine.Servers) *engine.Scene {
	scene := engine.NewScene("testScene1", nodeactor.NewNode(), servers)
	node1 := scene.SpawnActor("Node1", demo.NewDemo(
		*sprite.NewSprite(
			*transform2D.NewTransform2D(
				*nodeactor.NewNode(),
				transform2D.WithPos(gw*ps/2, gh*ps/2),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(true),
	), servers)

	node2 := scene.SpawnActor("Node2", transform2D.NewTransform2D(
		*nodeactor.NewNode(),
		transform2D.WithPos(0, 100),
	), servers)

	// This may get a little bit impractical for deeply nested stuff...
	// I guess we just have to wait and see
	node3 := scene.SpawnActor("Node3", demo.NewDemo(
		*sprite.NewSprite(
			*transform2D.NewTransform2D(
				*nodeactor.NewNode(),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	), servers)

	scene.SetParent(node2, node1)
	scene.SetParent(node3, node2)
	return scene
}

// Seems that servers contains some of the functionality that Commands
// has in Bevy...
func CreateTestScene2(servers *engine.Servers) *engine.Scene {
	scene := engine.NewScene("testScene2", nodeactor.NewNode(), servers)
	scene.SpawnActor("Node1", demo.NewDemo(
		*sprite.NewSprite(
			*transform2D.NewTransform2D(
				*nodeactor.NewNode(),
				transform2D.WithPos(gw*ps/2, gh*ps/3),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	), servers)
	return scene
}
