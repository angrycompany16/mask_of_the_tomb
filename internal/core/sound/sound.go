package sound

import (
	"mask_of_the_tomb/internal/core/resources"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// TODO: Add volume scaling for specific effect player
// for instance if we want to make the dash quieter
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

func NewEffectPlayer(player *audio.Player) *EffectPlayer {
	return &EffectPlayer{player}
}
