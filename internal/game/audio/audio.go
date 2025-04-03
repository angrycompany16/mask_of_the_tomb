package audio

import "github.com/hajimehoshi/ebiten/v2/audio"

const (
	sampleRate = 44100
)

// NOTE: audio needs to be a singleton because we only can have one audio context at a time
var (
	_audioManager = &audioManager{
		ctx: *audio.NewContext(sampleRate),
	}
)

type audioManager struct {
	ctx audio.Context
}

func Init() {

}

func Update() {

}

func NewSound() {
	sound := _audioManager.ctx.
}

// How to structure this?
// One centralized sound manager
// Sound player component (should make an asset)
// Sound players are added to game entities and
