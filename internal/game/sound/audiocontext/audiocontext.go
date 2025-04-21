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
		_globalAudioContext = &GlobalAudioContext{audio.NewContext(48000)}
	}
	return _globalAudioContext
}
