package musicplayer

import (
	"bytes"
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/gamestate"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// TODO: Fade music when launching game, switching tracks or (switching levels)
// The last one will have to be done later as it requires proper level transitions
// Ah and also dim when pausing the game

// TODO: Add multiple sound formats

type SongName int

const (
	MenuTheme SongName = iota
	BasementTheme
	LibraryTheme
)

type AmbienceName int

const (
	MenuAmbience AmbienceName = iota
	BasementAmbience
	LibraryAmbience
)

type MusicPlayer struct {
	songs      map[SongName]*audio.Player
	activeSong SongName
	ambience   map[AmbienceName]*audio.Player
}

func (m *MusicPlayer) Update(state gamestate.State, levelBiome string) {
	if state == gamestate.Loading {
		return
	}

	m.TryRestartSong()
	switch state {
	case gamestate.MainMenu:
		m.PlaySong(MenuTheme)
	case gamestate.Playing:
		if song, ok := m.songs[m.activeSong]; ok {
			song.SetVolume(1.0)
		}
		switch levelBiome {
		case "Basement":
			m.PlaySong(BasementTheme)
		case "Library":
			m.PlaySong(LibraryTheme)
		default:
			fmt.Println("Level has no biome, so no song will be played")
		}
	case gamestate.Paused:
		if song, ok := m.songs[m.activeSong]; ok {
			song.SetVolume(0.1)
		}
	}
}

func (m *MusicPlayer) TryRestartSong() {
	if song, ok := m.songs[m.activeSong]; ok {
		if !song.IsPlaying() {
			song.Rewind()
		}
	}
}

func (m *MusicPlayer) PlaySong(name SongName) {
	for _name, song := range m.songs {
		if _name != name && song.IsPlaying() {
			song.Pause()
			song.Rewind()
		}
		if _name == name && !song.IsPlaying() {
			m.activeSong = _name
			song.Play()
		}
	}
}

func NewMusicPlayer(ctx *audio.Context) *MusicPlayer {
	return &MusicPlayer{
		songs: map[SongName]*audio.Player{
			MenuTheme:     NewPlayer(assets.Menu_mp3, ctx),
			BasementTheme: NewPlayer(assets.Basement_mp3, ctx),
			LibraryTheme:  NewPlayer(assets.Library_mp3, ctx),
		},
	}
}

func NewPlayer(src []byte, ctx *audio.Context) *audio.Player {
	stream := errs.Must(mp3.DecodeF32(bytes.NewReader(src)))
	return errs.Must(ctx.NewPlayerF32(stream))
}
