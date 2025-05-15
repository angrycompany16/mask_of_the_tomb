package sound

import (
	"bytes"
	"mask_of_the_tomb/internal/core/audiocontext"
	"mask_of_the_tomb/internal/core/errs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type AudioFormat int

const (
	Mp3 AudioFormat = iota
	Wav
	Ogg
)

// Called an effect player because it's mostly used for sound effects, where the
// same sound might play multiple times on top of itself.
// Has some very simple concurrency implemented
type EffectPlayer struct {
	*audio.Player
}

func (e *EffectPlayer) Play() {
	if e.IsPlaying() {
		go playAudio(*e.Player)
		return
	}
	e.Rewind()
	e.Player.Play()
}

func playAudio(player audio.Player) {
	player.Rewind()
	for {
		if !player.IsPlaying() {
			return
		}
	}
}

func NewEffectPlayer(src []byte, format AudioFormat) *EffectPlayer {
	return &EffectPlayer{LoadAudio(src, format)}
}

// TODO: Add volume parameter 0-1
// TODO: Add audio asset
// TODO: Propagate errors
func LoadAudio(src []byte, format AudioFormat) *audio.Player {
	var player *audio.Player
	switch format {
	case Mp3:
		stream := errs.Must(mp3.DecodeF32(bytes.NewReader(src)))
		player = errs.Must(audiocontext.Current().NewPlayerF32(stream))
	case Ogg:
		stream := errs.Must(vorbis.DecodeF32(bytes.NewReader(src)))
		player = errs.Must(audiocontext.Current().NewPlayerF32(stream))
	case Wav:
		stream := errs.Must(wav.DecodeF32(bytes.NewReader(src)))
		player = errs.Must(audiocontext.Current().NewPlayerF32(stream))
	}
	return player
}
