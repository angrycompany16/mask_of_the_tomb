package sound

import (
	"bytes"
	"mask_of_the_tomb/internal/core/resources"

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
	e.Player.SetVolume(resources.Settings.MasterVolume * resources.Settings.SoundVolume / 10000.0)
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

func NewEffectPlayer(src []byte, format AudioFormat) (*EffectPlayer, error) {
	player, err := LoadAudio(src, format)
	return &EffectPlayer{player}, err
}

// TODO: Add volume parameter 0-1
// TODO: Add audio asset
// TODO: Propagate errors
func LoadAudio(src []byte, format AudioFormat) (*audio.Player, error) {
	var player *audio.Player
	var err error
	switch format {
	case Mp3:
		stream, _ := mp3.DecodeF32(bytes.NewReader(src))
		player, err = GetCurrentAudioContext().NewPlayerF32(stream)
	case Ogg:
		stream, _ := vorbis.DecodeF32(bytes.NewReader(src))
		player, err = GetCurrentAudioContext().NewPlayerF32(stream)
	case Wav:
		stream, _ := wav.DecodeF32(bytes.NewReader(src))
		player, err = GetCurrentAudioContext().NewPlayerF32(stream)
	}
	return player, err
}
