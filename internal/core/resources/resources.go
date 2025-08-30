package resources

var (
	Time              float64
	GrassWindSeed     int64
	PreviousLevelName string
	Settings          SettingsConfig
	DebugMode         bool
)

type SettingsConfig struct {
	MasterVolume float64
	SoundVolume  float64
	MusicVolume  float64
}
