package musicplayer

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/sound"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// TODO: Fade music when launching game, switching tracks or (switching levels)
// The last one will have to be done later as it requires proper level transitions

// TODO: Integrate music more tightly with LDTK

type ambienceName int

const (
	menuAmbience ambienceName = iota
	basementAmbience
	libraryAmbience
)

type MusicPlayer struct {
	songs      map[string]*audio.Player
	activeSong string
	ambience   map[ambienceName]*audio.Player
}

func (m *MusicPlayer) Init() {
	menuThemeStream := errs.Must(assettypes.GetMp3Stream("menuTheme"))
	basementThemeStream := errs.Must(assettypes.GetWavStream("basementTheme"))
	libraryThemeStream := errs.Must(assettypes.GetMp3Stream("libraryTheme"))

	m.songs["menuTheme"] = errs.Must(sound.FromStream(menuThemeStream))
	m.songs["basementTheme"] = errs.Must(sound.FromStream(basementThemeStream))
	m.songs["libraryTheme"] = errs.Must(sound.FromStream(libraryThemeStream))
}

func (m *MusicPlayer) PlayMenuMusic() {
	m.tryRestartSong()
	m.playSong("menuTheme")
}

func (m *MusicPlayer) PlayGameMusic(levelBiome string) {
	if song, ok := m.songs[m.activeSong]; ok {
		song.SetVolume(resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
	}
	switch levelBiome {
	case "Basement":
		m.playSong("basementTheme")
	case "Library":
		m.playSong("libraryTheme")
	default:
		m.stopMusic()
		fmt.Println("Level has no biome, so no song will be played")
	}
}

func (m *MusicPlayer) LowerMusic() {
	if song, ok := m.songs[m.activeSong]; ok {
		song.SetVolume(0.1 * resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
	}
}

func (m *MusicPlayer) ResetMusicVolume() {
	for _, song := range m.songs {
		song.SetVolume(resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
	}
}

func (m *MusicPlayer) tryRestartSong() {
	if song, ok := m.songs[m.activeSong]; ok {
		if !song.IsPlaying() {
			song.Rewind()
		}
	}
}

func (m *MusicPlayer) playSong(name string) {
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

func (m *MusicPlayer) stopMusic() {
	for _, song := range m.songs {
		if song.IsPlaying() {
			song.Pause()
			song.Rewind()
		}
	}
}

func NewMusicPlayer() *MusicPlayer {
	musicPlayer := MusicPlayer{
		songs: make(map[string]*audio.Player),
	}

	menuStream := errs.Must(assettypes.GetMp3Stream("menuTheme"))
	basementStream := errs.Must(assettypes.GetOggStream("basementTheme"))
	libraryStream := errs.Must(assettypes.GetMp3Stream("libraryTheme"))

	musicPlayer.songs["menuTheme"] = errs.Must(sound.FromStream(menuStream))
	musicPlayer.songs["basementTheme"] = errs.Must(sound.FromStream(basementStream))
	musicPlayer.songs["libraryTheme"] = errs.Must(sound.FromStream(libraryStream))

	return &musicPlayer
}
