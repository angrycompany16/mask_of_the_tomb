package main

import (
	"image"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/wfc"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/engine/enginebundles"
	"mask_of_the_tomb/internal/game/actors/overlappingmodel"
	"mask_of_the_tomb/internal/game/actors/simpletiledmodel"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gw, gh = 120, 68
	ps     = 16
)

type App struct {
	game *engine.Game
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

	sceneManager.RegisterScene("Scene", CreateSimpleTiledScene)
	sceneManager.RegisterScene("Scene2", CreateOverlappingMethodScene)
	sceneManager.SpawnScene("Scene2", cmd)

	app := &App{
		game: game,
	}

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}

func CreateOverlappingMethodScene(cmd *commands.Commands) *engine.Scene {
	scene := engine.NewScene("TestScene2", nodeactor.NewNode(), cmd)
	scene.SpawnBundle(
		cmd,
		enginebundles.MakeDefaultBundle(gw, gh, ps),
	)

	wfcSrc := assetloader.LoadImmediate[image.RGBA](
		cmd.AssetLoader,
		"wfcSrc",
		assettypes.NewRGBAImageAsset("sprites/concept-art/wfc-test.png"),
	)
	wfcSetup := wfc.NewOverlappingModel(wfcSrc, 3, 64, 64, true, true)

	scene.SpawnActor("Halla kompis",
		overlappingmodel.NewOverlappingModel(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			wfcSetup,
			renderer.ScreenTarget("Playerspace"),
			0,
		),
		cmd,
	)
	return scene
}

func CreateSimpleTiledScene(cmd *commands.Commands) *engine.Scene {
	scene := engine.NewScene("TestScene1", nodeactor.NewNode(), cmd)
	scene.SpawnBundle(
		cmd,
		enginebundles.MakeDefaultBundle(gw, gh, ps),
	)

	tileset := assetloader.LoadImmediate[ebiten.Image](
		cmd.AssetLoader,
		"TestTileset",
		assettypes.NewImageAsset("sprites/concept-art/simple_tiled_test.png"),
	)
	wfcSetup := wfc.NewSimpleTiled(16, 10, 10, tileset, 14)

	// 0 - SAND
	sand := wfc.NewModule(image.Rect(0, 0, 16, 16),
		wfc.NewDirectionalRule(maths.DirUp, 0, 1, 2),
		wfc.NewDirectionalRule(maths.DirDown, 0, 1, 2),
		wfc.NewDirectionalRule(maths.DirLeft, 0, 2),
		wfc.NewDirectionalRule(maths.DirRight, 0, 1, 2),
	)

	// 1 - GRASS
	grass := wfc.NewModule(image.Rect(16, 0, 32, 16),
		wfc.NewDirectionalRule(maths.DirUp, 1, 3),
		wfc.NewDirectionalRule(maths.DirDown, 0, 1, 3),
		wfc.NewDirectionalRule(maths.DirLeft, 1, 3),
		wfc.NewDirectionalRule(maths.DirRight, 0, 1, 3),
	)

	// 2 - WATER
	water := wfc.NewModule(image.Rect(0, 16, 16, 32),
		wfc.NewDirectionalRule(maths.DirUp, 0),
		wfc.NewDirectionalRule(maths.DirDown, 0),
		wfc.NewDirectionalRule(maths.DirLeft, 0, 2),
		wfc.NewDirectionalRule(maths.DirRight, 0, 2),
	)

	// 3 - STONE
	stone := wfc.NewModule(image.Rect(16, 16, 32, 32),
		wfc.NewDirectionalRule(maths.DirUp, 1, 3),
		wfc.NewDirectionalRule(maths.DirDown, 1, 3),
		wfc.NewDirectionalRule(maths.DirLeft, 1, 3),
		wfc.NewDirectionalRule(maths.DirRight, 1, 3),
	)

	wfcSetup.
		AddModule(sand).
		AddModule(grass).
		AddModule(water).
		AddModule(stone)

	wfcSetup.InitTiles()

	scene.SpawnActor("Halla kompis",
		simpletiledmodel.NewSimpleTiledWFC(
			graphic.NewGraphic(
				transform2D.NewTransform2D(
					nodeactor.NewNode(),
				),
			),
			wfcSetup,
			renderer.ScreenTarget("Playerspace"),
			0,
		),
		cmd)

	return scene
}
