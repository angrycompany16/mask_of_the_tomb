package musicplayer

import (
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/commands"
)

type MusicPlayer struct {
	*sound.SoundPlayer
	tracks map[string]int
}

func (m *MusicPlayer) Update(cmd *commands.Commands) {
	
}
