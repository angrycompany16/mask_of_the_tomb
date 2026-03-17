package main

import (
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/assetviewer"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/inspector"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/particles"
	"mask_of_the_tomb/internal/engine/actors/transform2D"

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
	game.RegisterScene("Scene1", CreateScene1)

	game.GetCmd().InputHandler().RegisterAction("toggleInspector", input.KeyJustPressedAction(ebiten.KeyTab))

	game.SpawnScene("Scene1")

	app := &App{
		game: game,
	}

	if err := ebiten.RunGame(app); err != nil {
		panic(err)
	}
}

func CreateScene1(cmd *engine.Commands) *engine.Scene {
	// Scene names aren't very useful atm
	scene := engine.NewScene("testScene1", nodeactor.NewNode(), cmd)

	// scene.SpawnActor("particleSystem",
	// 	particles.NewParticleSystem(
	// 		transform2D.NewTransform2D(
	// 			*nodeactor.NewNode(),
	// 			transform2D.WithPos(0, 0),
	// 		),
	// 		[]*particles.Burst{}, false, 1, 0,
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{1.0, 1.0},
	// 		maths.RandomFloat64{1.0, 1.0},
	// 		maths.RandomFloat64{2.0, 2.0},
	// 		maths.RandomFloat64{0, 0},
	// 		maths.RandomFloat64{0, 0},
	// 		[4]uint8{255, 255, 255, 255},
	// 		[4]uint8{255, 255, 255, 255},
	// 		128, 128,
	// 		"sprites/icons/square-64x64.png",
	// 		"Playerspace",
	// 		1,
	// 	), cmd,
	// )

	scene.SpawnActor("particleSystem",
		particles.NewParticleSystem(
			transform2D.NewTransform2D(
				*nodeactor.NewNode(),
				transform2D.WithPos(0, 0),
			),
			[]*particles.Burst{
				&particles.Burst{50, 0},
				&particles.Burst{50, 5},
			}, true, 10, 0,
			maths.RandomFloat64{-1, 1},
			maths.RandomFloat64{-100, 100},
			maths.RandomFloat64{-30, 30},
			maths.RandomFloat64{-40, 10},
			maths.RandomFloat64{0, 0},
			maths.RandomFloat64{0, 0},
			maths.RandomFloat64{0, 0},
			maths.RandomFloat64{0.1, 0.25},
			maths.RandomFloat64{0.0, 0.0},
			maths.RandomFloat64{1.0, 3.0},
			maths.RandomFloat64{0, 0},
			maths.RandomFloat64{0, 0},
			[4]uint8{255, 255, 255, 255},
			[4]uint8{255, 255, 255, 0},
			128, 128,
			"sprites/icons/square-64x64.png",
			"Playerspace",
			1,
		), cmd,
	)

	scene.SpawnActor("Inspector", inspector.NewInspector(0, 0, 300, 400), cmd)

	gw, gh := cmd.Renderer().GetGameSize()
	scene.SpawnActor("Camera", camera.NewCamera(
		transform2D.NewTransform2D(
			*nodeactor.NewNode(),
		),
		gw, gh, 0, 0, 0, 0,
	), cmd)

	scene.SpawnActor("AssetViewer", assetviewer.NewAssetViewer(
		nodeactor.NewNode(),
	), cmd)

	return scene
}
