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

// Note: Should only be used when the spritesheet is a strip of tiles
func NewSpritesheetAuto(img *ebiten.Image) *Spritesheet {
	tileSize := float64(img.Bounds().Size().Y)
	numTiles := float64(img.Bounds().Size().X) / tileSize
	return &Spritesheet{
		src:    img,
		width:  tileSize,
		height: tileSize,
		frames: int(numTiles),
	}
}

// IsFinished, GetNext, Reset, Play
type Animation struct {
	spritesheet *Spritesheet
	orientation SpritesheetOrientation
	loopMode    AnimationLoopMode
	xindex      int
	yindex      int
	frameDelay  float64
	t           float64
	paused      bool
	finished    bool
	next        int // id of the next animation we want to play. -1 if irrelevant
}

func (a *Animation) Update() {
	if a.paused {
		return
	}

	a.t += 0.008333333 // ? Why ? I don't knwo man. For some reason it is quite accurate
	if a.t > a.frameDelay {
		a.t = 0
		a.switchFrame()
	}
}

func (a *Animation) IsFinished() bool {
	return a.finished
}

func (a *Animation) GetNext() int {
	return a.next
}

func (a *Animation) Reset() {
	a.t = 0
	a.xindex = 0
	a.yindex = 0
	a.finished = false
}

func (a *Animation) switchFrame() {
	if a.orientation == Strip {
		a.xindex++
		if a.loopMode == Loop {
			a.xindex %= a.spritesheet.frames
		} else if a.loopMode == Once && a.xindex == a.spritesheet.frames {
			a.finished = true

			if a.next == -1 {
				a.Pause()
			}
		}
	}
}

func (a *Animation) Pause() {
	a.paused = true
}

func (a *Animation) Play() {
	a.paused = false
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

func NewAnimation(spritesheet *Spritesheet, frameDelay float64, orientation SpritesheetOrientation, loopMode AnimationLoopMode, next int) *Animation {
	return &Animation{
		spritesheet: spritesheet,
		orientation: orientation,
		xindex:      0,
		yindex:      0,
		frameDelay:  frameDelay,
		t:           0,
		loopMode:    loopMode,
		paused:      false,
		next:        next,
		finished:    false,
	}
}
