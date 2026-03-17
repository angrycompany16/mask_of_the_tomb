package main

import (
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/assetviewer"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/inspector"
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

func main() {
	game := engine.NewGame(engine.NewCommands(
		engine.WithRenderer(gw, gh, ps, false, false),
	))
	game.RegisterScene("TestScene1", CreateTestScene1)
	game.RegisterScene("TestScene2", CreateTestScene2)

	game.GetCmd().InputHandler().RegisterAction("toggleInspector", input.KeyJustPressedAction(ebiten.KeyTab))

	// Kinda cursed but this works?
	game.SpawnScene("TestScene1")

	app := &App{
		game: game,
	}

	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}

func CreateTestScene1(cmd *engine.Commands) *engine.Scene {
	scene := engine.NewScene("testScene1", nodeactor.NewNode(), cmd)
	node1 := scene.SpawnActor("Node1", demo.NewDemo(
		sprite.NewSprite(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(gw*ps/2, gh*ps/2),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(true),
	), cmd)

	node2 := scene.SpawnActor("Node2", transform2D.NewTransform2D(
		nodeactor.NewNode(),
		transform2D.WithPos(0, 100),
	), cmd)

	// This may get a little bit impractical for deeply nested stuff...
	// I guess we just have to wait and see
	node3 := scene.SpawnActor("Node3", demo.NewDemo(
		sprite.NewSprite(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	), cmd)

	scene.SpawnActor("Inspector", inspector.NewInspector(0, 0, 300, 400), cmd)

	gw, gh := cmd.Renderer().GetGameSize()
	scene.SpawnActor("Camera", camera.NewCamera(
		transform2D.NewTransform2D(
			nodeactor.NewNode(),
		),
		gw, gh, 0, 0, 0, 0,
	), cmd)

	scene.SpawnActor("AssetViewer", assetviewer.NewAssetViewer(
		nodeactor.NewNode(),
	), cmd)

	scene.SpawnActor("AssetViewer", assetviewer.NewAssetViewer(
		nodeactor.NewNode(),
	), cmd)

	scene.SetParent(node2, node1)
	scene.SetParent(node3, node2)
	return scene
}

func CreateTestScene2(cmd *engine.Commands) *engine.Scene {
	scene := engine.NewScene("testScene2", nodeactor.NewNode(), cmd)
	scene.SpawnActor("Node1", demo.NewDemo(
		sprite.NewSprite(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
				transform2D.WithPos(gw*ps/2, gh*ps/3),
			),
			"Playerspace",
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	), cmd)

	scene.SpawnActor("Inspector", inspector.NewInspector(0, 0, 300, 400), cmd)
	return scene
}
