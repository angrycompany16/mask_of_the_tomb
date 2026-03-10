package main

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/demo"
	"mask_of_the_tomb/internal/engine/actors/node"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/servers"

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
			a.game.SetActiveScene("testScene2")
		} else {
			a.game.SetActiveScene("testScene1")
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

	game := engine.NewGame(servers.NewServers(
		servers.ServerArgs{
			GameWidth:  gw,
			GameHeight: gh,
			PixelScale: ps,
		},
	), engine.NewEditor(gw*ps, gh*ps),
	)
	// Essentially scene defs
	testScene1 := game.MakeScene("testScene1", node.NewNode())
	testScene2 := game.MakeScene("testScene2", node.NewNode())
	SpawnTestScene1(testScene1)
	SpawnTestScene2(testScene2)

	// Would be nice to make this a bit more failsafe
	game.SetActiveScene(testScene1.GetName())

	app := &App{
		game: game,
	}

	ebiten.SetWindowSize(gw, gh)

	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}

// func TestScene1(servers *servers.Servers) *engine.Scene {
// 	scene := engine.NewScene("testScene1", node.NewNode(), servers)
// 	node1 := scene.Spawn("Node1", demo.NewDemo(
// 		*sprite.NewSprite(
// 			*transform2D.NewTransform2D(
// 				*node.NewNode(),
// 				transform2D.WithPos(gw*ps/2, gh*ps/2),
// 			),
// 			"Playerspace",
// 			"sprites/player/player.png",
// 			sprite.WithScaling(2.0),
// 		),
// 		demo.WithOnlyRotate(true),
// 	))

// 	node2 := scene.Spawn("Node2", transform2D.NewTransform2D(
// 		*node.NewNode(),
// 		transform2D.WithPos(0, 100),
// 	))

// 	// This may get a little bit impractical for deeply nested stuff...
// 	// I guess we just have to wait and see
// 	node3 := scene.Spawn("Node3", demo.NewDemo(
// 		*sprite.NewSprite(
// 			*transform2D.NewTransform2D(
// 				*node.NewNode(),
// 			),
// 			"Playerspace",
// 			"sprites/player/player.png",
// 			sprite.WithScaling(2.0),
// 		),
// 		demo.WithOnlyRotate(false),
// 	))

// 	scene.SetParent(node2, node1)
// 	scene.SetParent(node3, node2)
// 	return scene
// }

// func TestScene2(servers *servers.Servers) *engine.Scene {
// 	scene := engine.NewScene("testScene2", node.NewNode(), servers)
// 	testScene2.Spawn("Node1", demo.NewDemo(
// 		*sprite.NewSprite(
// 			*transform2D.NewTransform2D(
// 				*node.NewNode(),
// 				transform2D.WithPos(gw*ps/2, gh*ps/3),
// 			),
// 			"Playerspace",
// 			"sprites/player/player.png",
// 			sprite.WithScaling(2.0),
// 		),
// 		demo.WithOnlyRotate(false),
// 	))
// }

// In order for scenes to be reinstantiated at each
func SpawnTestScene1(testScene1 *engine.Scene) {
	node1 := testScene1.Spawn("Node1", demo.NewDemo(
		*sprite.NewSprite(
			*transform2D.NewTransform2D(
				*node.NewNode(),
				transform2D.WithPos(gw*ps/2, gh*ps/2),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(true),
	))

	node2 := testScene1.Spawn("Node2", transform2D.NewTransform2D(
		*node.NewNode(),
		transform2D.WithPos(0, 100),
	))

	// This may get a little bit impractical for deeply nested stuff...
	// I guess we just have to wait and see
	node3 := testScene1.Spawn("Node3", demo.NewDemo(
		*sprite.NewSprite(
			*transform2D.NewTransform2D(
				*node.NewNode(),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	))

	testScene1.SetParent(node2, node1)
	testScene1.SetParent(node3, node2)
}

func SpawnTestScene2(testScene2 *engine.Scene) {
	testScene2.Spawn("Node1", demo.NewDemo(
		*sprite.NewSprite(
			*transform2D.NewTransform2D(
				*node.NewNode(),
				transform2D.WithPos(gw*ps/2, gh*ps/3),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	))
}
