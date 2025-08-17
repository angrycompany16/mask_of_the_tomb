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
	src       []byte
	format    AudioFormat
	volume    float64
	mp3stream *mp3.Stream
	oggstream *vorbis.Stream
	wavstream *wav.Stream
}

func (a *AudioStreamAsset) Load() error {
	var err error
	switch a.format {
	case Mp3:
		a.mp3stream, err = mp3.DecodeF32(bytes.NewReader(a.src))
	case Ogg:
		a.oggstream, err = vorbis.DecodeF32(bytes.NewReader(a.src))
	case Wav:
		a.wavstream, err = wav.DecodeF32(bytes.NewReader(a.src))
	}
	return err
}

func GetMp3Stream(name string) (*mp3.Stream, error) {
	soundAsset, err := assetloader.GetAsset(name)
	return soundAsset.(*AudioStreamAsset).mp3stream, err
}

func GetOggStream(name string) (*vorbis.Stream, error) {
	soundAsset, err := assetloader.GetAsset(name)
	return soundAsset.(*AudioStreamAsset).oggstream, err
}

func GetWavStream(name string) (*wav.Stream, error) {
	soundAsset, err := assetloader.GetAsset(name)
	return soundAsset.(*AudioStreamAsset).wavstream, err
}

func MakeAudioStreamAsset(src []byte, format AudioFormat) *AudioStreamAsset {
	return &AudioStreamAsset{
		src:    src,
		format: format,
	}
}
