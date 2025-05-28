package resources

var Time float64
var GrassWindSeed int64
var State GameState
var PreviousLevelName string
var Settings SettingsConfig

type GameState int

const (
	Loading GameState = iota
	MainMenu
	Intro
	Playing
	Paused
)

type SettingsConfig struct {
	MasterVolume float64
	SoundVolume  float64
	MusicVolume  float64
}
