package screenfade

import (
	"image/color"
	ebitenrenderutil "mask_of_the_tomb/internal/ebitenrenderutil"
	"mask_of_the_tomb/internal/game/core/events"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Generalize into an overlay type?
// For suuure we can make an interface for this, or something

var (
	OverlayColor = []uint8{20, 16, 19}
)

type deathEffectState int

const (
	enter deathEffectState = iota
	exit
	idle
)

type DeathEffect struct {
	state         deathEffectState
	image         *ebiten.Image
	alpha         float64
	OnFinishEnter *events.Event
	OnFinishExit  *events.Event
}

func (d *DeathEffect) Update() {
	switch d.state {
	case enter:
		d.alpha = maths.Lerp(d.alpha, 3, 0.01)
		if 1-d.alpha <= 0.01 {
			d.alpha = 1
			d.OnFinishEnter.Raise(events.EventInfo{})
		}
	case exit:
		d.alpha = maths.Lerp(d.alpha, -2, 0.01)
		if d.alpha <= 0.01 {
			d.alpha = 0
			d.OnFinishExit.Raise(events.EventInfo{})
			d.state = idle
		}
	case idle:
	}
}

func (d *DeathEffect) Draw() {
	alpha := uint8(d.alpha * 255)
	d.image.Fill(color.RGBA{OverlayColor[0], OverlayColor[1], OverlayColor[2], alpha})
	ebitenrenderutil.DrawAt(d.image, rendering.RenderLayers.Overlay, 0, 0)
}

func (d *DeathEffect) StartEnter() {
	d.state = enter
}

func (d *DeathEffect) StartExit() {
	d.state = exit
}

func NewDeathEffect() *DeathEffect {
	return &DeathEffect{
		image:         ebiten.NewImage(rendering.GameWidth, rendering.GameHeight),
		alpha:         0.0,
		state:         idle,
		OnFinishEnter: events.NewEvent(),
		OnFinishExit:  events.NewEvent(),
	}
}
