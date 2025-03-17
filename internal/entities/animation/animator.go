package animation

import (
	"fmt"
	"mask_of_the_tomb/internal/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
)

// Essentially a container for a map of animations
type Animator struct {
	clips             map[int]*Animation
	ActiveClip        int // The animation that is currently active
	FinishedClipEvent *events.Event
}

func (a *Animator) Update() {
	activeClip := a.clips[a.ActiveClip]

	activeClip.Update()
	if activeClip.finished {
		if activeClip.next != -1 {
			a.FinishedClipEvent.Raise(events.EventInfo{})
			a.ActiveClip = activeClip.next
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

	activeClip.reset()
	activeClip.play()
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
		clips:             clips,
		ActiveClip:        0,
		FinishedClipEvent: events.NewEvent(),
	}
}
