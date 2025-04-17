package sound

import (
	"bytes"
	"mask_of_the_tomb/internal/errs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

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

// TODO: Add volume parameter 0-1
func NewEffectPlayer(src []byte, ctx *audio.Context, format AudioFormat) *EffectPlayer {
	var player *audio.Player
	switch format {
	case Mp3:
		stream := errs.Must(mp3.DecodeF32(bytes.NewReader(src)))
		player = errs.Must(ctx.NewPlayerF32(stream))
	case Ogg:
		stream := errs.Must(vorbis.DecodeF32(bytes.NewReader(src)))
		player = errs.Must(ctx.NewPlayerF32(stream))
	case Wav:
		stream := errs.Must(wav.DecodeF32(bytes.NewReader(src)))
		player = errs.Must(ctx.NewPlayerF32(stream))
	}
	return &EffectPlayer{player}
}
