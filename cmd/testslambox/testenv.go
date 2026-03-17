package main

import (
	"errors"
	"fmt"
	"log"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/inspector"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/actors/slamboxtilemap"
	"mask_of_the_tomb/internal/game/actors/tracker"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var testSlamboxChains = make([]*slambox.SlamboxChain, 0)

// var testSlamboxChains = []*SlamboxChain{
// 	NewSlamboxChain(
// 		[]float64{16 + 7, 16 + 7, 48 + 7, 48 + 7},
// 		[]float64{48 + 7, 16 + 7, 16 + 7, 64 + 7},

// 		[]*Slambox{NewSlambox(maths.NewRect(50, 50, 12, 12), 5)},
// 		[]*SlamboxGroup{
// 			NewSlamboxGroup(
// 				[]*Slambox{
// 					NewSlambox(maths.NewRect(18, 30, 12, 12), 5),
// 					NewSlambox(maths.NewRect(19, 42, 10, 10), 5),
// 				},
// 				0,
// 			),
// 		},
// 	),
// }

// var testSlamboxGroups = make([]*slambox.SlamboxGroup, 0)

// var testSlamboxGroups = []*SlamboxGroup{
// 	NewSlamboxGroup(
// 		[]*Slambox{
// 			NewSlambox(maths.NewRect(64, 16, 16, 16), 5),
// 		}, 0,
// 	),
// 	NewSlamboxGroup(
// 		[]*Slambox{
// 			NewSlambox(maths.NewRect(24, 16, 16, 8), 5),
// 			NewSlambox(maths.NewRect(16, 24, 32, 16), 5),
// 		}, 0,
// 	),
// }

type SlamboxController struct {
	*slamboxactor.Slambox
}

func (sc *SlamboxController) Update(cmd *engine.Commands) {
	sc.Slambox.Update(cmd)
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		sc.RequestSlam(maths.DirLeft)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		sc.RequestSlam(maths.DirRight)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		sc.RequestSlam(maths.DirUp)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		sc.RequestSlam(maths.DirDown)
	}
}

var testSlamboxes = []*maths.Rect{
	maths.NewRect(64, 16, 16, 16),
	maths.NewRect(24, 16, 16, 8),
	maths.NewRect(16, 24, 32, 16),
}

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type App struct {
	game *engine.Game
	surf *ebiten.Image
}

func (a *App) Update() error {
	return a.game.Update()
}

func (a *App) Draw(screen *ebiten.Image) {
	a.game.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return GAME_WIDTH * PIXEL_SCALE, GAME_HEIGHT * PIXEL_SCALE
}

// Sets up a simple game loop for testing this package.
func main() {
	game := engine.NewGame(engine.NewCommands(
		engine.WithRenderer(GAME_WIDTH, GAME_HEIGHT, PIXEL_SCALE, true, true),
		engine.WithSlamboxEnvironment(16),
	))

	game.RegisterScene("testScene", CreateTestScene)

	game.SpawnScene("testScene")

	game.GetCmd().InputHandler().RegisterAction("toggleInspector", input.KeyJustPressedAction(ebiten.KeyTab))

	a := &App{
		game: game,
	}

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		}
		log.Fatal(err)
	}
}

func CreateTestScene(cmd *engine.Commands) *engine.Scene {
	gridTiles := make([][]int, 0)
	gridTiles = append(gridTiles, []int{1, 1, 1, 1, 1, 1})
	gridTiles = append(gridTiles, []int{1, 0, 0, 0, 0, 1})
	gridTiles = append(gridTiles, []int{1, 0, 0, 0, 0, 1})
	gridTiles = append(gridTiles, []int{1, 0, 0, 0, 0, 1})
	gridTiles = append(gridTiles, []int{1, 1, 1, 0, 0, 1})
	gridTiles = append(gridTiles, []int{1, 1, 1, 1, 1, 1})

	scene := engine.NewScene("testScene", nodeactor.NewNode(), cmd)
	for i, slambox := range testSlamboxes {
		scene.SpawnActor(fmt.Sprintf("slambox %d", i),
			&SlamboxController{
				Slambox: slamboxactor.NewSlambox(
					tracker.NewTracker(
						transform2D.NewTransform2D(
							*nodeactor.NewNode(),
						),
						5.0,
						slambox.Left(), slambox.Top(),
					),
					slambox,
				),
			},
			cmd)
	}

	scene.SpawnActor("SlamboxTilemap", slamboxtilemap.NewSlamboxTilemap(
		transform2D.NewTransform2D(
			*nodeactor.NewNode(),
		),
		gridTiles,
		16,
	), cmd)

	// gw, gh := cmd.Renderer().GetGameSize()
	// scene.SpawnActor("MainCamera", camera.NewCamera(
	// 	transform2D.NewTransform2D(
	// 		*nodeactor.NewNode(),
	// 	),
	// 	0, 0, gw, gh,
	// ),
	// 	gridTiles,
	// 	16,
	// ), cmd)

	scene.SpawnActor("Inspector", inspector.NewInspector(0, 0, 300, 400), cmd)

	return scene
}
