package assettypes

import (
	"bytes"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/sound"

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

type SoundAsset struct {
	src    []byte
	format AudioFormat
	Player *audio.Player
}

func (a *SoundAsset) Load() error {
	var player *audio.Player
	var err error
	switch a.format {
	case Mp3:
		stream, _ := mp3.DecodeF32(bytes.NewReader(a.src))
		player, err = sound.GetCurrentAudioContext().NewPlayerF32(stream)
	case Ogg:
		stream, _ := vorbis.DecodeF32(bytes.NewReader(a.src))
		player, err = sound.GetCurrentAudioContext().NewPlayerF32(stream)
	case Wav:
		stream, _ := wav.DecodeF32(bytes.NewReader(a.src))
		player, err = sound.GetCurrentAudioContext().NewPlayerF32(stream)
	}
	a.Player = player
	return err
}

func GetSoundAsset(name string) (*audio.Player, error) {
	soundAsset, err := assetloader.GetAsset(name)
	return soundAsset.(*SoundAsset).Player, err
}

func GetEffectPlayerAsset(name string) (*sound.EffectPlayer, error) {
	soundAsset, err := assetloader.GetAsset(name)
	return sound.NewEffectPlayer(soundAsset.(*SoundAsset).Player), err
}

// TODO: Add volume parameter 0-1
func MakeSoundAsset(src []byte, format AudioFormat) *SoundAsset {
	return &SoundAsset{
		src:    src,
		format: format,
	}
}

// package assettypes

// import (
// 	"bytes"
// 	"mask_of_the_tomb/internal/core/assetloader"
// 	"mask_of_the_tomb/internal/core/sound"

// 	"github.com/hajimehoshi/ebiten/v2/audio"
// 	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
// 	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
// 	"github.com/hajimehoshi/ebiten/v2/audio/wav"
// )

// // TODO: Decouple from sound library itself
// // A goal: Always keep assets loading stuff that's exernal to the program
// type AudioFormat int

// const (
// 	Mp3 AudioFormat = iota
// 	Wav
// 	Ogg
// )

// type AudioStream interface {
// 	mp3.Stream | vorbis.Stream | wav.Stream
// }

// type SoundAsset[T AudioStream] struct {
// 	src    []byte
// 	format AudioFormat
// 	stream T
// }

// func (a *SoundAsset) Load() error {
// 	var player *audio.Player
// 	var err error
// 	switch a.format {
// 	case Mp3:
// 		stream, err := mp3.DecodeF32(bytes.NewReader(a.src))
// 		a.stream = stream
// 	case Ogg:
// 		stream, err := vorbis.DecodeF32(bytes.NewReader(a.src))
// 		a.stream = stream
// 	case Wav:
// 		stream, err := wav.DecodeF32(bytes.NewReader(a.src))
// 		a.stream = stream
// 	}
// 	a.Player = player
// 	return err
// }

// func GetSoundAsset(name string) (*audio.Player, error) {
// 	soundAsset, err := assetloader.GetAsset(name)
// 	return soundAsset.(*SoundAsset).Player, err
// }

// func GetEffectPlayerAsset(name string) (*sound.EffectPlayer, error) {
// 	soundAsset, err := assetloader.GetAsset(name)
// 	return sound.NewEffectPlayer(soundAsset.(*SoundAsset).Player), err
// }

// // TODO: Add volume parameter 0-1
// func MakeSoundAsset(src []byte, format AudioFormat) *SoundAsset {
// 	return &SoundAsset{
// 		src:    src,
// 		format: format,
// 	}
// }
