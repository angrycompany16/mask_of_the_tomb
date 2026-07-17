package sounddebug

import (
	sound_v2 "mask_of_the_tomb/internal/backend/sound"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SoundDebug struct {
	*nodeactor.Node
}

func (s *SoundDebug) Update(cmd *commands.Commands) {
	if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		sound_v2.DebugSoundServer()
	}
}

func CreateSoundDebug() *SoundDebug {
	return &SoundDebug{
		nodeactor.NewNode(),
	}
}
