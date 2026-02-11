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

var slamboxes = make([]*maths.Rect, 0)

const (
	GAME_WIDTH, GAME_HEIGHT = 480.0, 270.0
	PIXEL_SCALE             = 4.0
)

type testApp struct {
	SlamboxEnvironment *SlamboxEnvironment
	environmentRects   []*maths.Rect
	surf               *ebiten.Image
	movable            *maths.Rect
}

func (t *testApp) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		projRect, _ := t.SlamboxEnvironment.ProjectRect(*t.movable, maths.DirLeft)
		t.movable = &projRect
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		projRect, _ := t.SlamboxEnvironment.ProjectRect(*t.movable, maths.DirRight)
		t.movable = &projRect
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		projRect, _ := t.SlamboxEnvironment.ProjectRect(*t.movable, maths.DirDown)
		t.movable = &projRect
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		projRect, _ := t.SlamboxEnvironment.ProjectRect(*t.movable, maths.DirUp)
		t.movable = &projRect
	}
	return nil

}

func (t *testApp) Draw(screen *ebiten.Image) {
	t.surf.Fill(color.RGBA{0, 0, 0, 255})
	for _, rect := range t.environmentRects {
		DrawRect(t.surf, rect, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 200, 0, 255})
	}
	DrawRect(t.surf, t.movable, color.RGBA{255, 255, 255, 255}, color.RGBA{0, 200, 255, 255})

	ebitenrenderutil.DrawAtScaled(t.surf, screen, 16, 16, 8, 8)
}

func (t *testApp) Layout(outsideHeight, outsideWidth int) (int, int) {
	return GAME_WIDTH * PIXEL_SCALE, GAME_HEIGHT * PIXEL_SCALE
}

// Sets up a simple game loop for testing this package.
func RunTestEnv() {
	gridTiles := make([][]bool, 0)
	gridTiles = append(gridTiles, []bool{true, true, true, true, true})
	gridTiles = append(gridTiles, []bool{true, false, true, false, true})
	gridTiles = append(gridTiles, []bool{true, false, true, true, true})
	gridTiles = append(gridTiles, []bool{true, true, false, false, true})
	gridTiles = append(gridTiles, []bool{true, true, true, true, true})

	a := &testApp{
		SlamboxEnvironment: NewSlamboxEnvironment(16, gridTiles, slamboxes),
		surf:               ebiten.NewImage(480, 270),
		movable:            maths.NewRect(16, 16, 16, 16),
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
