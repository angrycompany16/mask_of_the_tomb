package audiocontext

import "github.com/hajimehoshi/ebiten/v2/audio"

var (
	_globalAudioContext *GlobalAudioContext
)

type GlobalAudioContext struct {
	*audio.Context
}

func Current() *GlobalAudioContext {
	if _globalAudioContext == nil {
		_globalAudioContext = &GlobalAudioContext{audio.NewContext(44100)}
	}
	return _globalAudioContext
}
