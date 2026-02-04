package ui

import (
	"fmt"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/threads"
	"time"
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
	timeout       time.Duration
	OnFinishEnter *events.Event
	OnFinishExit  *events.Event
	OnIdleTimeout *events.Event
	IdleTimer     *time.Timer
}

func (oi *Overlay) Update() {
	switch oi.state {
	case enter:
		oi.fadeIn()
	case exit:
		oi.fadeOut()
	case off:
	case on:
		if _, timedout := threads.Poll(oi.IdleTimer.C); timedout {
			oi.OnIdleTimeout.Raise(events.EventInfo{})
		}
	}
}

func (oi *Overlay) Draw() {
	switch oi.state {
	case enter, on:
		oi.OverlayContent.Draw(oi.t, true)
	case exit, off:
		oi.OverlayContent.Draw(oi.t, false)
	}
}

func (oi *Overlay) StartFadeIn() {
	oi.state = enter
}

func (oi *Overlay) StartFadeOut() {
	oi.state = exit
}

func (oi *Overlay) fadeIn() {
	oi.t = maths.Lerp(oi.t, 3, 0.01)
	if 1-oi.t <= 0.01 {
		oi.t = 1
		oi.OnFinishEnter.Raise(events.EventInfo{})
		oi.IdleTimer = time.NewTimer(oi.timeout)
		oi.state = on
	}
}

func (oi *Overlay) fadeOut() {
	oi.t = maths.Lerp(oi.t, -2, 0.01)
	if oi.t <= 0.01 {
		fmt.Println("Faded out")
		oi.t = 0
		oi.OnFinishExit.Raise(events.EventInfo{})
		oi.state = off
	}
}

func NewOverlay(content OverlayContent, timeout time.Duration) *Overlay {
	return &Overlay{
		OverlayContent: content,
		state:          off,
		timeout:        timeout,
		OnIdleTimeout:  events.NewEvent(),
		OnFinishEnter:  events.NewEvent(),
		OnFinishExit:   events.NewEvent(),
	}
}

type OverlayContent interface {
	Draw(t float64, enter bool)
}
