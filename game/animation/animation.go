package animation

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpritesheetOrientation int

const (
	Strip SpritesheetOrientation = iota
	MultiRow
)

type AnimationLoopMode int

const (
	Loop AnimationLoopMode = iota
	Once
)

type Spritesheet struct {
	src           *ebiten.Image
	width, height float64
	frames        int
}

// Note: Should only be used when there's a spritesheet with quadratic tiles
func NewSpritesheetAuto(img *ebiten.Image) *Spritesheet {
	// Do calculations
	tileSize := float64(img.Bounds().Size().Y)
	numTiles := float64(img.Bounds().Size().X) / tileSize
	return &Spritesheet{
		src:    img,
		width:  tileSize,
		height: tileSize,
		frames: int(numTiles),
	}
}

type Animation struct {
	spritesheet *Spritesheet
	orientation SpritesheetOrientation
	loopMode    AnimationLoopMode
	xindex      int
	yindex      int
	frameDelay  float64
	t           float64
	paused      bool
}

func (a *Animation) Update() {
	if a.paused {
		return
	}

	a.t += 0.01666666666667
	if a.t > a.frameDelay {
		a.t = 0
		a.switchFrame()
	}
}

func (a *Animation) switchFrame() {
	if a.orientation == Strip {
		a.xindex++
		if a.loopMode == Loop {
			a.xindex %= a.spritesheet.frames
		} else if a.loopMode == Once && a.xindex == a.spritesheet.frames {
			a.Pause()
		}
		return
	}
}

func (a *Animation) Pause() {
	a.paused = true
}

func (a *Animation) GetSprite() *ebiten.Image {
	if a.orientation == Strip {
		return a.spritesheet.src.SubImage(
			image.Rect(
				a.xindex*int(a.spritesheet.width),
				0,
				(a.xindex+1)*int(a.spritesheet.width),
				int(a.spritesheet.height),
			),
		).(*ebiten.Image)
	}
	return a.spritesheet.src
}

func NewAnimation(spritesheet *Spritesheet, frameDelay float64, orientation SpritesheetOrientation, loopMode AnimationLoopMode) *Animation {
	return &Animation{
		spritesheet: spritesheet,
		orientation: orientation,
		xindex:      0,
		yindex:      0,
		frameDelay:  frameDelay,
		t:           0,
		loopMode:    loopMode,
		paused:      false,
	}
}
