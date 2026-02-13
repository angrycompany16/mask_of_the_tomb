package slambox

import (
	"errors"
	"image/color"
	"log"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/scenes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var testSlamboxChains = []*SlamboxChain{
	NewSlamboxChain(
		[]float64{16 + 4, 16 + 4, 48 + 4, 48 + 4},
		[]float64{48 + 4, 16 + 4, 16 + 4, 64 + 4},
		[]*Slambox{
			NewSlambox(maths.NewRect(30, 18, 12, 12), 5),
			NewSlambox(maths.NewRect(18, 34, 12, 12), 5),
		},
		[]*SlamboxGroup{},
	),
}

var testSlamboxGroups = []*SlamboxGroup{}
var testSlamboxes = []*Slambox{}

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type testApp struct {
	SlamboxEnvironment *SlamboxEnvironment
	environmentRects   []*maths.Rect
	surf               *ebiten.Image
}

func (t *testApp) Update() error {
	t.SlamboxEnvironment.Update()
	if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		t.SlamboxEnvironment.SlamSlamboxChain(0, 0, false, maths.DirLeft)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		t.SlamboxEnvironment.SlamSlamboxChain(0, 0, false, maths.DirRight)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		t.SlamboxEnvironment.SlamSlamboxChain(0, 0, false, maths.DirUp)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		t.SlamboxEnvironment.SlamSlamboxChain(0, 0, false, maths.DirDown)
	}

	return nil
}

func (t *testApp) Draw(screen *ebiten.Image) {
	t.surf.Fill(color.RGBA{0, 0, 0, 255})
	for _, rect := range t.environmentRects {
		DrawRect(t.surf, rect, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 200, 0, 255})
	}
	slamboxes := t.SlamboxEnvironment.GetSlamboxes()
	for _, slambox := range slamboxes {
		DrawRect(t.surf, slambox.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
	}

	slamboxGroups := t.SlamboxEnvironment.GetSlamboxGroups()
	for _, slamboxGroup := range slamboxGroups {
		for _, slambox := range slamboxGroup.GetSlamboxes() {
			DrawRect(t.surf, slambox.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
		}
	}

	slamboxChains := t.SlamboxEnvironment.GetSlamboxChains()
	for _, slamboxChain := range slamboxChains {
		for _, chainNode := range slamboxChain.GetNodes() {
			DrawRect(t.surf, chainNode.GetRect(), color.RGBA{255, 0, 255, 255}, color.RGBA{100, 0, 255, 255})
		}
		for _, slambox := range slamboxChain.GetSlamboxes() {
			DrawRect(t.surf, slambox.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
		}
		for _, slamboxGroup := range slamboxChain.GetSlamboxGroups() {
			for _, slambox := range slamboxGroup.GetSlamboxes() {
				DrawRect(t.surf, slambox.GetRect(), color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})
			}
		}
	}

	ebitenrenderutil.DrawAtScaled(t.surf, screen, 16, 16, 8, 8)
}

func (t *testApp) Layout(outsideHeight, outsideWidth int) (int, int) {
	return GAME_WIDTH * PIXEL_SCALE, GAME_HEIGHT * PIXEL_SCALE
}

// Sets up a simple game loop for testing this package.
func RunTestEnv() {
	gridTiles := make([][]bool, 0)
	gridTiles = append(gridTiles, []bool{true, true, true, true, true, true})
	gridTiles = append(gridTiles, []bool{true, false, false, false, false, true})
	gridTiles = append(gridTiles, []bool{true, false, false, false, false, true})
	gridTiles = append(gridTiles, []bool{true, false, false, false, false, true})
	gridTiles = append(gridTiles, []bool{true, true, true, false, false, true})
	gridTiles = append(gridTiles, []bool{true, true, true, true, true, true})

	a := &testApp{
		SlamboxEnvironment: NewSlamboxEnvironment(16, gridTiles, testSlamboxes, testSlamboxGroups, testSlamboxChains),
		surf:               ebiten.NewImage(480, 270),
	}
	ebiten.SetWindowSize(480*2, 270*2)

	a.environmentRects = a.SlamboxEnvironment.Rectify()

	if err := ebiten.RunGame(a); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			return
		} else if err == scenes.ErrTerminated {
			return
		}
		log.Fatal(err)
	}
}
