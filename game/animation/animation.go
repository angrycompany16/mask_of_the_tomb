package animation

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpritesheetOrientation int

const (
	StripOrientation SpritesheetOrientation = iota
	MultiRow
)

type Spritesheet struct {
	src           *ebiten.Image
	width, height float64
	frames        int
}

func NewSpritesheetAuto(img *ebiten.Image) *Spritesheet {
	// Do calculations
	return &Spritesheet{
		src:    img,
		width:  16,
		height: 16,
		frames: 0,
	}
}

type Animation struct {
	spritesheet *Spritesheet
	orientation SpritesheetOrientation
	xindex      int
	yindex      int
	frameDelay  float64
	t           float64
}

func (a *Animation) Update() {
	a.t += 0.01666666666667
	if a.t > a.frameDelay {
		a.t = 0
		a.SwitchFrame()
	}
}

func (a *Animation) SwitchFrame() {
	if a.orientation == StripOrientation {
		a.xindex++
		a.yindex %= a.spritesheet.frames
		return
	}
}

func (a *Animation) GetSprite() *ebiten.Image {
	if a.orientation == StripOrientation {
		return a.spritesheet.src.SubImage(
			image.Rect(
				a.xindex*int(a.spritesheet.width),
				0,
				a.yindex*int(a.spritesheet.width+1),
				0,
			),
		).(*ebiten.Image)
	}
	return a.spritesheet.src
}

func NewAnimation(spritesheet *Spritesheet, frameDelay float64, orientation SpritesheetOrientation) *Animation {
	return &Animation{
		spritesheet: spritesheet,
		orientation: orientation,
		xindex:      0,
		yindex:      0,
		frameDelay:  frameDelay,
		t:           0,
	}
}
