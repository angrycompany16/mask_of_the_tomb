package animation

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// Essentially a container for a map of animations
type Animator struct {
	clips      map[int]*Animation
	ActiveClip int // The animation that is currently active
}

func (a *Animator) Update() {
	activeClip, ok := a.clips[a.ActiveClip]

	if !ok {
		return
	}

	activeClip.Update()
	if activeClip.finished {
		if activeClip.next != -1 {
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
	activeClip, ok := a.clips[a.ActiveClip]
	if !ok {
		return ebiten.NewImage(1, 1)
	}
	return activeClip.GetSprite()

}

func (a *Animator) AddAnimation(anim *Animation, id int) {
	a.clips[id] = anim
}

func NewAnimator(clips map[int]*Animation) *Animator {
	return &Animator{
		clips:      clips,
		ActiveClip: 0,
	}
}
