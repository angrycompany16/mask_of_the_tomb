package animatedsprite

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimatedSprite struct {
	*graphic.Graphic
	Clips          map[string]*Clip
	ActiveClipName string
	OnClipFinished *eventsv2.Event
	layer          string  `debug:"auto"`
	drawOrder      int     `debug:"auto"`
	pivotX         float64 `debug:"auto"`
	pivotY         float64 `debug:"auto"`
}

func (a *AnimatedSprite) OnTreeAdd(node *engine.Node, cmd *engine.Commands) {
	a.Graphic.OnTreeAdd(node, cmd)

	for _, clip := range a.Clips {
		clip.spritesheetImgAsset = assetloader.StageAsset[ebiten.Image](
			cmd.AssetLoader(),
			clip.spritesheetPath,
			assettypes.NewImageAsset(clip.spritesheetPath),
		)
	}
}

func (a *AnimatedSprite) Init(cmd *engine.Commands) {
	a.Graphic.Init(cmd)
	for _, clip := range a.Clips {
		clip.nFrames = clip.spritesheetImgAsset.Value().Bounds().Dx() / int(clip.frameWidth)
	}
}

func (a *AnimatedSprite) Update(cmd *engine.Commands) {
	a.Graphic.Update(cmd)
	activeClip := a.Clips[a.ActiveClipName]

	activeClip.Update()
	if activeClip.IsFinished() {
		if activeClip.GetNext() != "" {
			a.OnClipFinished.WithData("clip", a.ActiveClipName).Raise()
			a.ActiveClipName = activeClip.GetNext()
		}
	}

	gPosX, gPosY := a.GetPos(false)
	camX, camY := a.GetCamera().WorldToCam(gPosX, gPosY, true)
	angle := a.GetAngle(false)
	scaleX, scaleY := a.GetScale(false)

	// Create draw call
	cmd.Renderer().Request(
		opgen.PosRotScale(activeClip.GetSprite(), camX, camY, angle, scaleX, scaleY, a.pivotX, a.pivotY),
		activeClip.GetSprite(),
		a.layer,
		a.drawOrder,
	)
}

func (a *AnimatedSprite) SwitchClip(newClip string) {
	if newClip == a.ActiveClipName {
		return
	}

	a.ActiveClipName = newClip
	activeClip, ok := a.Clips[a.ActiveClipName]
	if !ok {
		fmt.Println("Tried to set animator to invalid clip", newClip)
		return
	}

	activeClip.Reset()
	activeClip.Play()
}

func (a *AnimatedSprite) AddAnimation(clip *Clip, name string) {
	a.Clips[name] = clip
}

func NewAnimatedSprite(
	graphic *graphic.Graphic,
	clips map[string]*Clip,
	layer string,
	drawOrder int,
	pivotX, pivotY float64,
	startClip string,
) *AnimatedSprite {
	return &AnimatedSprite{
		Graphic:        graphic,
		Clips:          clips,
		OnClipFinished: eventsv2.NewEvent(),
		layer:          layer,
		drawOrder:      drawOrder,
		pivotX:         pivotX,
		pivotY:         pivotY,
		ActiveClipName: startClip,
	}
}

type AnimationLoopMode int

const (
	Loop AnimationLoopMode = iota
	Once
)

type Clip struct {
	LoopMode                AnimationLoopMode
	FrameDelay              int
	NextAnimationName       string
	spritesheetImgAsset     *assetloader.AssetRef[ebiten.Image]
	spritesheetPath         string
	frameWidth, frameHeight float64
	nFrames                 int
	finished                bool
	paused                  bool
	ticker                  *time.Ticker
	xindex                  int
	yindex                  int
}

func (c *Clip) Update() {
	if c.paused {
		return
	}

	if _, tick := utils.PollThread(c.ticker.C); tick {
		c.switchFrame()
	}
}

func (c *Clip) IsFinished() bool {
	return c.finished
}

func (c *Clip) GetNext() string {
	return c.NextAnimationName
}

func (c *Clip) Reset() {
	c.ticker = time.NewTicker(time.Duration(c.FrameDelay * int(time.Millisecond)))
	c.xindex = 0
	c.yindex = 0
	c.finished = false
}

func (c *Clip) switchFrame() {
	c.xindex++
	if c.LoopMode == Loop {
		c.xindex %= c.nFrames
	} else if c.LoopMode == Once && c.xindex == c.nFrames {
		c.finished = true
		c.xindex--

		if c.NextAnimationName == "" {
			c.Pause()
		}
	}
}

func (c *Clip) Pause() {
	c.paused = true
}

func (c *Clip) Play() {
	c.paused = false
}

func (c *Clip) GetSprite() *ebiten.Image {
	return c.spritesheetImgAsset.Value().SubImage(
		image.Rect(
			c.xindex*int(c.frameWidth),
			0,
			(c.xindex+1)*int(c.frameWidth),
			int(c.frameHeight),
		),
	).(*ebiten.Image)
}

func NewClip(spritesheetPath string, frameWidth, frameHeight float64, loopMode AnimationLoopMode, frameDelay int, nextAnimationName string) *Clip {
	return &Clip{
		spritesheetPath:   spritesheetPath,
		frameWidth:        frameWidth,
		frameHeight:       frameHeight,
		LoopMode:          loopMode,
		FrameDelay:        frameDelay,
		NextAnimationName: nextAnimationName,
		finished:          false,
		paused:            false,
		ticker:            time.NewTicker(time.Duration(frameDelay * int(time.Millisecond))),
		xindex:            0,
		yindex:            0,
	}
}
