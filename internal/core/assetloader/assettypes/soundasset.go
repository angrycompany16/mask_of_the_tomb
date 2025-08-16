package assettypes

import (
	"bytes"
	"mask_of_the_tomb/internal/core/assetloader"

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

type AudioStreamAsset struct {
	src    []byte
	format AudioFormat
	stream any
}

func (a *AudioStreamAsset) Load() error {
	var err error
	switch a.format {
	case Mp3:
		a.stream, err = mp3.DecodeF32(bytes.NewReader(a.src))
	case Ogg:
		a.stream, err = vorbis.DecodeF32(bytes.NewReader(a.src))
	case Wav:
		a.stream, err = wav.DecodeF32(bytes.NewReader(a.src))
	}
	return err
}

func GetAudioStreamAsset(name string) (any, error) {
	soundAsset, err := assetloader.GetAsset(name)
	return soundAsset.(*AudioStreamAsset).stream, err
}

// TODO: Add volume parameter 0-1
func MakeAudioStreamAsset(src []byte, format AudioFormat) *AudioStreamAsset {
	return &AudioStreamAsset{
		src:    src,
		format: format,
	}
}
