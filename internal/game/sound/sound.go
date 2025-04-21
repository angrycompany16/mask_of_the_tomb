package sound

import (
	"bytes"
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/gamestate"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// TODO: Fade music when launching game, switching tracks or (switching levels)
// The last one will have to be done later as it requires proper level transitions
// Ah and also dim when pausing the game

// TODO: Integrate music more tightly with LDTK

type AudioFormat int

const (
	Mp3 AudioFormat = iota
	Wav
	Ogg
)

type songName int

const (
	menuTheme songName = iota
	basementTheme
	libraryTheme
)

type ambienceName int

const (
	menuAmbience ambienceName = iota
	basementAmbience
	libraryAmbience
)

type MusicPlayer struct {
	songs      map[songName]*audio.Player
	activeSong songName
	ambience   map[ambienceName]*audio.Player
}

func (m *MusicPlayer) Update(state gamestate.State, levelBiome string) {
	if state == gamestate.Loading {
		return
	}

	m.tryRestartSong()
	switch state {
	case gamestate.MainMenu:
		m.playSong(menuTheme)
	case gamestate.Playing:
		if song, ok := m.songs[m.activeSong]; ok {
			song.SetVolume(1.0)
		}
		switch levelBiome {
		case "Basement":
			m.playSong(basementTheme)
		case "Library":
			m.playSong(libraryTheme)
		default:
			fmt.Println("Level has no biome, so no song will be played")
		}
	case gamestate.Paused:
		if song, ok := m.songs[m.activeSong]; ok {
			song.SetVolume(0.1)
		}
	}
}

func (m *MusicPlayer) tryRestartSong() {
	if song, ok := m.songs[m.activeSong]; ok {
		if !song.IsPlaying() {
			song.Rewind()
		}
	}
}

func (m *MusicPlayer) playSong(name songName) {
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
		songs: map[songName]*audio.Player{
			menuTheme:     NewAudioPlayer(assets.Menu_mp3, ctx, Mp3),
			basementTheme: NewAudioPlayer(assets.Basement_mp3, ctx, Mp3),
			libraryTheme:  NewAudioPlayer(assets.Library_mp3, ctx, Mp3),
		},
	}
}

func NewAudioPlayer(src []byte, ctx *audio.Context, format AudioFormat) *audio.Player {
	switch format {
	case Mp3:
		stream := errs.Must(mp3.DecodeF32(bytes.NewReader(src)))
		return errs.Must(ctx.NewPlayerF32(stream))
	case Ogg:
		stream := errs.Must(vorbis.DecodeF32(bytes.NewReader(src)))
		return errs.Must(ctx.NewPlayerF32(stream))
	case Wav:
		stream := errs.Must(wav.DecodeF32(bytes.NewReader(src)))
		return errs.Must(ctx.NewPlayerF32(stream))
	}
	return nil
}
