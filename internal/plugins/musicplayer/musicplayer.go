package musicplayer

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/libraries/gamestate"
	"mask_of_the_tomb/internal/libraries/sound"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// TODO: Fade music when launching game, switching tracks or (switching levels)
// The last one will have to be done later as it requires proper level transitions
// Ah and also dim when pausing the game

// TODO: Integrate music more tightly with LDTK

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

// TODO: This will become a plugin
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
			// fmt.Println("Level has no biome, so no song will be played")
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
			menuTheme:     sound.LoadAudio(assets.Menu_mp3, sound.Mp3),
			basementTheme: sound.LoadAudio(assets.Basement_wav, sound.Wav),
			libraryTheme:  sound.LoadAudio(assets.Library_mp3, sound.Mp3),
		},
	}
}
