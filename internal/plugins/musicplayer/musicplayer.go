package musicplayer

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/sound"

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
func (m *MusicPlayer) Update(levelBiome string) {
	if resources.State == resources.Loading {
		return
	}

	for _, song := range m.songs {
		song.SetVolume(resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
	}

	m.tryRestartSong()
	switch resources.State {
	case resources.MainMenu:
		m.playSong(menuTheme)
	case resources.Playing:
		if song, ok := m.songs[m.activeSong]; ok {
			song.SetVolume(resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
		}
		switch levelBiome {
		case "Basement":
			m.playSong(basementTheme)
		case "Library":
			m.playSong(libraryTheme)
		default:
			// fmt.Println("Level has no biome, so no song will be played")
		}
	case resources.Paused:
		if song, ok := m.songs[m.activeSong]; ok {
			song.SetVolume(0.1 * resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
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
			menuTheme:     errs.Must(sound.LoadAudio(assets.Menu_mp3, sound.Mp3)),
			basementTheme: errs.Must(sound.LoadAudio(assets.Basement_wav, sound.Wav)),
			libraryTheme:  errs.Must(sound.LoadAudio(assets.Library_mp3, sound.Mp3)),
		},
	}
}
