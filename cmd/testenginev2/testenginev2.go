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
			a.game.SpawnScene((&TestScene2{}).Name())
		} else {
			a.game.SpawnScene((&TestScene1{}).Name())
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
	game.AddScene(&TestScene1{})
	game.AddScene(&TestScene2{})

	// Kinda cursed but this works?
	game.SpawnScene((&TestScene1{}).Name())

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
type TestScene1 struct{}

func (t *TestScene1) Create(servers *servers.Servers) *engine.Scene {
	scene := engine.NewScene("testScene1", node.NewNode(), servers)
	node1 := scene.SpawnActor("Node1", demo.NewDemo(
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

	node2 := scene.SpawnActor("Node2", transform2D.NewTransform2D(
		*node.NewNode(),
		transform2D.WithPos(0, 100),
	))

	// This may get a little bit impractical for deeply nested stuff...
	// I guess we just have to wait and see
	node3 := scene.SpawnActor("Node3", demo.NewDemo(
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

	scene.SetParent(node2, node1)
	scene.SetParent(node3, node2)
	return scene
}

func (t *TestScene1) Name() string {
	return "TestScene1"
}

type TestScene2 struct{}

func (t *TestScene2) Create(servers *servers.Servers) *engine.Scene {
	scene := engine.NewScene("testScene2", node.NewNode(), servers)
	scene.SpawnActor("Node1", demo.NewDemo(
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
	return scene
}

func (t *TestScene2) Name() string {
	return "TestScene2"
}

// In order for scenes to be reinstantiated at each
// func SpawnTestScene1(testScene1 *engine.Scene) {
// 	node1 := testScene1.SpawnActor("Node1", demo.NewDemo(
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

// 	node2 := testScene1.SpawnActor("Node2", transform2D.NewTransform2D(
// 		*node.NewNode(),
// 		transform2D.WithPos(0, 100),
// 	))

// 	// This may get a little bit impractical for deeply nested stuff...
// 	// I guess we just have to wait and see
// 	node3 := testScene1.SpawnActor("Node3", demo.NewDemo(
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

// 	testScene1.SetParent(node2, node1)
// 	testScene1.SetParent(node3, node2)
// }

// func SpawnTestScene2(testScene2 *engine.Scene) {
// 	testScene2.SpawnActor("Node1", demo.NewDemo(
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
