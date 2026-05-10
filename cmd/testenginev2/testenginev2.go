package main

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/cmd/testenginev2/actors/demo"
	"mask_of_the_tomb/cmd/testenginev2/actors/switcher"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/assetviewer"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/inspector"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
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
	return a.game.Update()
}

func (a *App) Draw(screen *ebiten.Image) {
	a.game.Draw(screen)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	inputHandler := input.NewInputHandler()
	cmd := commands.NewCommands(
		renderer.NewRenderer(gw, gh, ps, false, false),
		assetloader.NewAssetLoader(assets.FS),
		inputHandler,
	)

	game := engine.NewGame(cmd)
	sceneManager, _ := commands.Get[engine.SceneManager](cmd)

	sceneManager.RegisterScene("TestScene1", CreateTestScene1)
	sceneManager.RegisterScene("TestScene2", CreateTestScene2)

	// Kinda cursed but this works?
	sceneManager.SpawnScene("TestScene1", cmd)

	app := &App{
		game: game,
	}

	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}

func CreateTestScene1(cmd *commands.Commands) *engine.Scene {
	scene := engine.NewScene("TestScene1", nodeactor.NewNode(), cmd)
	scene.SpawnActor("Switcher", switcher.NewSwitch(
		nodeactor.NewNode(),
		scene.GetName(),
	), cmd)
	node1 := scene.SpawnActor("Node1", demo.NewDemo(
		sprite.NewSprite(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
					transform2D.WithPos(gw*ps/2, gh*ps/2),
				),
			),
			renderer.RenderTarget{
				Type: renderer.SCREEN,
				Name: "Playerspace",
			},
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
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			renderer.RenderTarget{
				Type: renderer.SCREEN,
				Name: "Playerspace",
			},
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	), cmd)

	scene.SpawnActor("Inspector", inspector.NewInspector(
		nodeactor.NewNode(),
		inspector.WithPos(0, 0),
		inspector.WithSize(300, 400),
	), cmd)

	gw, gh := cmd.Renderer.GetGameSize()
	scene.SpawnActor("Camera", camera.NewCamera(
		transform2D.NewTransform2D(
			nodeactor.NewNode(),
		),
		camera.WithSize(gw, gh),
		camera.WithMargins(0, 0),
		camera.WithOffset(0, 0),
	), cmd)

	scene.SpawnActor("AssetViewer", assetviewer.NewAssetViewer(
		nodeactor.NewNode(),
	), cmd)

	scene.SetParent(node2, node1)
	scene.SetParent(node3, node2)
	return scene
}

func CreateTestScene2(cmd *commands.Commands) *engine.Scene {
	scene := engine.NewScene("TestScene2", nodeactor.NewNode(), cmd)
	scene.SpawnActor("Switcher", switcher.NewSwitch(
		nodeactor.NewNode(),
		scene.GetName(),
	), cmd)
	scene.SpawnActor("Node1", demo.NewDemo(
		sprite.NewSprite(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
					transform2D.WithPos(gw*ps/2, gh*ps/3),
				),
			),
			renderer.RenderTarget{
				Type: renderer.SCREEN,
				Name: "Playerspace",
			},
			"sprites/player/player.png",
			sprite.WithScaling(2.0),
		),
		demo.WithOnlyRotate(false),
	), cmd)

	scene.SpawnActor("Inspector", inspector.NewInspector(
		nodeactor.NewNode(),
		inspector.WithPos(0, 0),
		inspector.WithSize(300, 400),
	), cmd)
	return scene
}
