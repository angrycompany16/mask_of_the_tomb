package resources

var (
	Time              float64
	GrassWindSeed     int64
	State             GameState
	PreviousLevelName string
	Settings          SettingsConfig
	DebugMode         bool
)

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
