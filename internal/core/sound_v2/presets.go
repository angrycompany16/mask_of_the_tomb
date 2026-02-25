package sound_v2

import (
	"errors"
	"strings"

	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

func Loop(path string) SoundData {
	return SoundData{
		Path:      path,
		Looping:   true,
		format:    GetFormat(path),
		QueueSize: 1,
	}
}

func Oneshot(path string, queueSize int) SoundData {
	return SoundData{
		Path:      path,
		Looping:   false,
		format:    GetFormat(path),
		QueueSize: queueSize,
	}
}

func GetFormat(path string) AudioFormat {
	if strings.HasSuffix(path, ".ogg") {
		return Ogg
	} else if strings.HasSuffix(path, ".mp3") {
		return Mp3
	} else if strings.HasSuffix(path, ".wav") {
		return Wav
	}
	return -1
}

func SetVolumeAction(volume float64) func(effect resound.IEffect) error {
	return func(effect resound.IEffect) error {
		effect, ok := effect.(*effects.Volume)
		if !ok {
			return errors.New("Could not convert effect interface to volume")
		}
		effect.(*effects.Volume).SetStrength(volume)
		return nil
	}
}
