package musicplayer

import (
	"bytes"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/gamestate"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// TODO: Fade music when launching game, switching tracks or (switching levels)
// The last one will have to be done later as it requires proper level transitions
// Ah and also dim when pausing the game

type SongName int

const (
	MenuSong SongName = iota
	LevelSong
)

type MusicPlayer struct {
	songs map[SongName]*audio.Player
}

func (m *MusicPlayer) Update(state gamestate.State) {
	switch state {
	case gamestate.Loading:
	case gamestate.MainMenu:
		if !m.songs[MenuSong].IsPlaying() {
			m.songs[MenuSong].Play()
		}
	case gamestate.Playing:
		m.songs[MenuSong].Pause()
	case gamestate.Paused:
		// TODO: Dim the music
	}
}

func NewMusicPlayer(ctx *audio.Context) *MusicPlayer {
	stream := errs.Must(mp3.DecodeF32(bytes.NewReader(assets.Menu_mp3)))
	p := errs.Must(ctx.NewPlayerF32(stream))

	return &MusicPlayer{
		songs: map[SongName]*audio.Player{
			MenuSong: p,
		},
	}
}
