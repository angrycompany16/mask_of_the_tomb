package animatedsprite

import (
	"fmt"
	"image"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"
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
	Name              string
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

	if _, tick := utils.PollThread(a.ticker.C); tick {
		a.switchFrame()
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
			a.xindex--

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

func NewAnimation(info AnimationInfo, width, height float64) *Animation {
	src := utils.MustNewImageFromFile(info.SpriteSheetPath)
	return &Animation{
		info:     info,
		finished: false,
		paused:   false,
		ticker:   time.NewTicker(time.Duration(info.FrameDelay * int(time.Millisecond))),
		xindex:   0,
		yindex:   0,
		// no no NO
		spritesheet: &Spritesheet{
			src:    src,
			width:  width,
			height: height,
			frames: src.Bounds().Dx() / int(width),
		},
	}
}

func NewAnimationAuto(info AnimationInfo) *Animation {
	return &Animation{
		info:     info,
		finished: false,
		paused:   false,
		ticker:   time.NewTicker(time.Duration(info.FrameDelay * int(time.Millisecond))),
		xindex:   0,
		yindex:   0,
		// no no NO
		spritesheet: NewSpritesheetAuto(utils.MustNewImageFromFile(info.SpriteSheetPath)),
	}
}

type Spritesheet struct {
	// TODO: Turn this into an AssetRef
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

// Inherit from sprite. Control Image during update. Will need asset loading
type AnimatedSprite struct {
	*graphic.Graphic
	Clips          map[int]*Animation
	ActiveClip     int
	OnClipFinished *eventsv2.Event
	layer          string `debug:"auto"`
	drawOrder      int    `debug:"auto"`
}

func (a *AnimatedSprite) Update(cmd *engine.Commands) {
	a.Graphic.Update(cmd)
	activeClip := a.Clips[a.ActiveClip]

	activeClip.Update()
	if activeClip.IsFinished() {
		if activeClip.GetNext() != -1 {
			a.OnClipFinished.WithData("clip", activeClip.info.Name).Raise()
			a.ActiveClip = activeClip.GetNext()
		}
	}

	gPosX, gPosY := a.GetPos(false)
	camX, camY := a.GetCamera().WorldToCam(gPosX, gPosY, true)
	angle := a.GetAngle(false)
	scaleX, scaleY := a.GetScale(false)

	// Create draw call
	cmd.Renderer().Request(
		opgen.PosRotScale(activeClip.GetSprite(), camX, camY, angle, scaleX, scaleY, 0.5, 0.5),
		activeClip.GetSprite(),
		a.layer,
		a.drawOrder,
	)
}

func (a *AnimatedSprite) SwitchClip(newClip int) {
	if newClip == a.ActiveClip {
		return
	}

	a.ActiveClip = newClip
	activeClip, ok := a.Clips[a.ActiveClip]
	if !ok {
		fmt.Println("Tried to set animator to invalid clip", newClip)
		return
	}

	activeClip.Reset()
	activeClip.Play()
}

func (a *AnimatedSprite) AddAnimation(anim *Animation, id int) {
	a.Clips[id] = anim
}

func NewAnimator(graphic *graphic.Graphic, clips map[int]*Animation, layer string, drawOrder int) *AnimatedSprite {
	return &AnimatedSprite{
		Graphic:        graphic,
		Clips:          clips,
		OnClipFinished: eventsv2.NewEvent(),
		layer:          layer,
		drawOrder:      drawOrder,
	}
}
