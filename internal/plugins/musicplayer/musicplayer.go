package musicplayer

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/libraries/assettypes"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// TODO: Fade music when launching game, switching tracks or (switching levels)
// The last one will have to be done later as it requires proper level transitions

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

func (m *MusicPlayer) Load() {
	assetloader.Load("menuTheme", assettypes.MakeSoundAsset(assets.Menu_mp3, assettypes.Mp3))
	assetloader.Load("basementTheme", assettypes.MakeSoundAsset(assets.Basement_wav, assettypes.Wav))
	assetloader.Load("libraryTheme", assettypes.MakeSoundAsset(assets.Library_mp3, assettypes.Mp3))
}

func (m *MusicPlayer) Init() {
	m.songs[menuTheme] = errs.Must(assettypes.GetSoundAsset("menuTheme"))
	m.songs[basementTheme] = errs.Must(assettypes.GetSoundAsset("basementTheme"))
	m.songs[libraryTheme] = errs.Must(assettypes.GetSoundAsset("libraryTheme"))
}

func (m *MusicPlayer) PlayMenuMusic() {
	m.tryRestartSong()
	m.playSong(menuTheme)
}

func (m *MusicPlayer) PlayGameMusic(levelBiome string) {
	if song, ok := m.songs[m.activeSong]; ok {
		song.SetVolume(resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
	}
	switch levelBiome {
	case "Basement":
		m.playSong(basementTheme)
	case "Library":
		m.playSong(libraryTheme)
	default:
		fmt.Println("Level has no biome, so no song will be played")
	}
}

func (m *MusicPlayer) LowerMusic() {
	if song, ok := m.songs[m.activeSong]; ok {
		song.SetVolume(0.1 * resources.Settings.MasterVolume * resources.Settings.MusicVolume / 10000.0)
	}
}

func (m *MusicPlayer) ResetMusic() {
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

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		songs: make(map[songName]*audio.Player),
	}
}
