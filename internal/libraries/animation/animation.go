package animation

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/threads"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpritesheetFormat int

const (
	Strip SpritesheetFormat = iota
	MultiRow
)

type AnimationLoopMode int

const (
	Loop AnimationLoopMode = iota
	Once
)

type AnimationInfo struct {
	SpriteSheetPath   string            `yaml:"SpriteSheetPath"`
	SpriteSheetFormat SpritesheetFormat `yaml:"SpriteSheetFormat"`
	LoopMode          AnimationLoopMode `yaml:"LoopMode"`
	FrameDelay        int               `yaml:"FrameDelay"`
	NextAnimationId   int               `yaml:"NextAnimationId"`
}

type Animation struct {
	info        AnimationInfo
	finished    bool
	paused      bool
	ticker      *time.Ticker
	xindex      int
	yindex      int
	spritesheet *Spritesheet
}

func (a *Animation) Update() {
	if a.paused {
		return
	}

	if _, tick := threads.Poll(a.ticker.C); tick {
		a.switchFrame()
		fmt.Println("switched frame")
		fmt.Println(a.xindex)
	} else {
		fmt.Println("Just a normal frame")
		fmt.Println(a.xindex)
	}
}

func (a *Animation) IsFinished() bool {
	return a.finished
}

func (a *Animation) GetNext() int {
	return a.info.NextAnimationId
}

func (a *Animation) Reset() {
	a.ticker = time.NewTicker(time.Duration(a.info.FrameDelay * int(time.Millisecond)))
	a.xindex = 0
	a.yindex = 0
	a.finished = false
}

func (a *Animation) switchFrame() {
	if a.info.SpriteSheetFormat == Strip {
		a.xindex++
		if a.info.LoopMode == Loop {
			a.xindex %= a.spritesheet.frames
		} else if a.info.LoopMode == Once && a.xindex == a.spritesheet.frames {
			a.finished = true

			if a.info.NextAnimationId == -1 {
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
	if a.info.SpriteSheetFormat == Strip {
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

func NewAnimation(info AnimationInfo) *Animation {
	return &Animation{
		info:        info,
		finished:    false,
		paused:      false,
		ticker:      time.NewTicker(time.Duration(info.FrameDelay * int(time.Millisecond))),
		xindex:      0,
		yindex:      0,
		spritesheet: NewSpritesheetAuto(errs.MustNewImageFromFile(info.SpriteSheetPath)),
	}
}

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
