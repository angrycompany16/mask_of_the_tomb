package animation

import (
	"fmt"
	"mask_of_the_tomb/internal/core/events"

	"github.com/hajimehoshi/ebiten/v2"
)

// Essentially a container for a map of animations
type Animator struct {
	clips          map[int]*Animation
	ActiveClip     int // The animation that is currently active
	OnClipFinished *events.Event
}

func (a *Animator) Update() {
	activeClip := a.clips[a.ActiveClip]

	activeClip.Update()
	if activeClip.IsFinished() {
		if activeClip.GetNext() != -1 {
			a.OnClipFinished.Raise(events.EventInfo{})
			a.ActiveClip = activeClip.GetNext()
		}
	}
}

func (a *Animator) SwitchClip(newClip int) {
	if newClip == a.ActiveClip {
		return
	}

	a.ActiveClip = newClip
	activeClip, ok := a.clips[a.ActiveClip]
	if !ok {
		fmt.Println("Tried to set animator to invalid clip", newClip)
		return
	}

	activeClip.Reset()
	activeClip.Play()
}

func (a *Animator) GetSprite() *ebiten.Image {
	return a.clips[a.ActiveClip].GetSprite()
}

func (a *Animator) AddAnimation(anim *Animation, id int) {
	a.clips[id] = anim
}

func NewAnimator(clips map[int]*Animation) *Animator {
	// NOTE: An empty animator cannet exist which probably isn't great
	return &Animator{
		clips:          clips,
		ActiveClip:     0,
		OnClipFinished: events.NewEvent(),
	}
}
