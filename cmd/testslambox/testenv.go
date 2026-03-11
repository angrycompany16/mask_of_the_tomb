package main

import (
	"errors"
	"image/color"
	"log"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/slambox"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var testSlamboxChains = make([]*slambox.SlamboxChain, 0)

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

var testSlamboxGroups = make([]*slambox.SlamboxGroup, 0)

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

var testSlamboxes = []*slambox.Slambox{
	slambox.NewSlambox(maths.NewRect(64, 16, 16, 16), 5),
	slambox.NewSlambox(maths.NewRect(24, 16, 16, 8), 5),
	slambox.NewSlambox(maths.NewRect(16, 24, 32, 16), 5),
}

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type testApp struct {
	SlamboxEnvironment *slambox.SlamboxEnvironment
	environmentRects   []*maths.Rect
	surf               *ebiten.Image
}

func (t *testApp) Update() error {
	t.SlamboxEnvironment.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		t.SlamboxEnvironment.SlamSlambox(0, maths.DirLeft)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		t.SlamboxEnvironment.SlamSlambox(0, maths.DirRight)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		t.SlamboxEnvironment.SlamSlambox(0, maths.DirUp)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		t.SlamboxEnvironment.SlamSlambox(0, maths.DirDown)
	}

	return nil
}

func (t *testApp) Draw(screen *ebiten.Image) {
	t.surf.Fill(color.RGBA{0, 0, 0, 255})
	for _, rect := range t.environmentRects {
		slambox.DrawRect(t.surf, rect, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 200, 0, 255})
	}
	slamboxes := t.SlamboxEnvironment.GetSlamboxes()
	for _, slambox_ := range slamboxes {
		slambox.DrawRect(t.surf, slambox_.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
	}

	slamboxGroups := t.SlamboxEnvironment.GetSlamboxGroups()
	for _, slamboxGroup := range slamboxGroups {
		for _, slambox_ := range slamboxGroup.GetSlamboxes() {
			slambox.DrawRect(t.surf, slambox_.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
		}
	}

	slamboxChains := t.SlamboxEnvironment.GetSlamboxChains()
	for _, slamboxChain := range slamboxChains {
		for _, chainNode := range slamboxChain.GetNodes() {
			slambox.DrawRect(t.surf, chainNode.GetRect(), color.RGBA{255, 0, 255, 255}, color.RGBA{100, 0, 255, 255})
		}
		for _, slambox_ := range slamboxChain.GetSlamboxes() {
			slambox.DrawRect(t.surf, slambox_.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
		}
		for _, slamboxGroup := range slamboxChain.GetSlamboxGroups() {
			for _, slambox_ := range slamboxGroup.GetSlamboxes() {
				slambox.DrawRect(t.surf, slambox_.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
			}
		}
	}
	op := opgen.PosScale(t.surf, 16, 16, 8, 8)
	screen.DrawImage(t.surf, op)
}

func (t *testApp) Layout(outsideHeight, outsideWidth int) (int, int) {
	return GAME_WIDTH * PIXEL_SCALE, GAME_HEIGHT * PIXEL_SCALE
}

// Sets up a simple game loop for testing this package.
func main() {
	gridTiles := make([][]bool, 0)
	gridTiles = append(gridTiles, []bool{true, true, true, true, true, true})
	gridTiles = append(gridTiles, []bool{true, false, false, false, false, true})
	gridTiles = append(gridTiles, []bool{true, false, false, false, false, true})
	gridTiles = append(gridTiles, []bool{true, false, false, false, false, true})
	gridTiles = append(gridTiles, []bool{true, true, true, false, false, true})
	gridTiles = append(gridTiles, []bool{true, true, true, true, true, true})

	a := &testApp{
		SlamboxEnvironment: slambox.NewSlamboxEnvironment(16, gridTiles, testSlamboxes, testSlamboxGroups, testSlamboxChains),
		surf:               ebiten.NewImage(480, 270),
	}
	ebiten.SetWindowSize(480*2, 270*2)

	a.environmentRects = a.SlamboxEnvironment.Rectify()

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		}
		log.Fatal(err)
	}
}
