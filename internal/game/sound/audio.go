package sound

import "github.com/hajimehoshi/ebiten/v2/audio"

// TODO: This is probably not needed
// At this point it's just a wrapper for audio.context, which is useless

type SoundManager struct {
	ctx *audio.Context
}

func (s *SoundManager) Update() {

}

func (s *SoundManager) AddSound() {

}

func NewSoundManager(ctx *audio.Context) *SoundManager {
	return &SoundManager{
		ctx: ctx,
	}
}

// How to structure this?
// One centralized sound manager
// Sound player component (should make an asset)
// Sound players are added to game entities and
