package overlay

import (
	"mask_of_the_tomb/internal/game/core/events"
	"mask_of_the_tomb/internal/maths"
)

type overlayState int

const (
	enter overlayState = iota
	exit
	off
	on
)

type Overlay struct {
	OverlayContent
	t             float64
	state         overlayState
	OnFinishEnter *events.Event
	OnFinishExit  *events.Event
	IdleTime      float64
}

func (oi *Overlay) Update() {
	switch oi.state {
	case enter:
		oi.fadeIn()
	case exit:
		oi.fadeOut()
	case off:
	case on:
		oi.IdleTime += 1.0 / 60.0
	}
}

func (oi *Overlay) Draw() {
	oi.OverlayContent.Draw(oi.t)
}

func (oi *Overlay) StartFadeIn() {
	oi.IdleTime = 0
	oi.state = enter
}

func (oi *Overlay) StartFadeOut() {
	oi.IdleTime = 0
	oi.state = exit
}

func (oi *Overlay) fadeIn() {
	// oi.state = enter
	oi.t = maths.Lerp(oi.t, 3, 0.01)
	if 1-oi.t <= 0.01 {
		oi.t = 1
		oi.OnFinishEnter.Raise(events.EventInfo{})
		oi.state = on
	}
}

func (oi *Overlay) fadeOut() {
	// oi.state = exit
	oi.t = maths.Lerp(oi.t, -2, 0.01)
	if oi.t <= 0.01 {
		// fmt.Println("Faded out")
		oi.t = 0
		oi.OnFinishExit.Raise(events.EventInfo{})
		oi.state = off
	}
}

func NewOverlay(content OverlayContent) *Overlay {
	return &Overlay{
		OverlayContent: content,
		state:          off,
		OnFinishEnter:  events.NewEvent(),
		OnFinishExit:   events.NewEvent(),
	}
}

type OverlayContent interface {
	Draw(t float64)
}
