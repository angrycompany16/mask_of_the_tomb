package animator

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/animation"
	"mask_of_the_tomb/internal/backend/events"

	"github.com/hajimehoshi/ebiten/v2"
)

// Inherit from sprite. Control Image during update. Will need asset loading
type Animator struct {
	Clips          map[int]*animation.Animation `yaml:"Clips"`
	ActiveClip     int
	OnClipFinished *events.Event
}

func (a *Animator) Update() {
	activeClip := a.Clips[a.ActiveClip]

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
	activeClip, ok := a.Clips[a.ActiveClip]
	if !ok {
		fmt.Println("Tried to set animator to invalid clip", newClip)
		return
	}

	activeClip.Reset()
	activeClip.Play()
}

func (a *Animator) GetSprite() *ebiten.Image {
	clip, _ := a.Clips[a.ActiveClip]
	return clip.GetSprite()
}

func (a *Animator) AddAnimation(anim *animation.Animation, id int) {
	a.Clips[id] = anim
}

func MakeAnimator(clips map[int]*animation.Animation) *Animator {
	return &Animator{
		Clips:          clips,
		OnClipFinished: events.NewEvent(),
	}
}
