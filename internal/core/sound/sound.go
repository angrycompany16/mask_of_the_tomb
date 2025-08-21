package sound

import (
	"fmt"
	"io"
	"mask_of_the_tomb/internal/core/resources"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// TODO: Figure out how to add pitch control so that we can add some pitch randomization
type EffectPlayer struct {
	*audio.Player
	Volume float64
}

func (e *EffectPlayer) Play() {
	fmt.Println("Playing sound")
	e.Player.SetVolume(e.Volume * resources.Settings.MasterVolume * resources.Settings.SoundVolume / 20000.0)
	fmt.Println("We managed to set the colume level")

	// if e.IsPlaying() {
	// 	go playAudio(*e.Player)
	// 	return
	// }

	if !e.IsPlaying() {
		e.Rewind()
		e.Player.Play()
	}
}

func playAudio(player audio.Player) {
	player.Rewind()
	for {
		if !player.IsPlaying() {
			return
		}
	}
}

func FromStream[T io.Reader](stream T) (*audio.Player, error) {
	player, err := GetAudioContext().NewPlayerF32(stream)
	if err != nil {
		return nil, err
	}
	return player, nil
}
