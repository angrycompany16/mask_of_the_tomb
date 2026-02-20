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

type ldtkNames struct {
	EntityLayer            string
	PlayerSpaceLayer       string
	SpawnPosEntity         string
	DoorEntity             string
	SpawnPointEntity       string
	SlamboxEntity          string
	SpikeIntGrid           string
	GameEntryPos           string
	GrassEntity            string
	HazardEntity           string
	TurretEntity           string
	CatcherEntity          string
	PlatformEntity         string
	LanternEntity          string
	ChainNodeEntity        string
	SlamboxChainEntity     string
	TestSpeechBubbleEntity string
	LevelTitleField        string
}

var LDTKNames = ldtkNames{
	EntityLayer:            "Entities",
	PlayerSpaceLayer:       "Playerspace",
	SpawnPosEntity:         "DefaultSpawnPos",
	DoorEntity:             "Door",
	SpawnPointEntity:       "SpawnPoint",
	SlamboxEntity:          "Slambox",
	SpikeIntGrid:           "Spikes",
	GameEntryPos:           "GameEntryPos",
	GrassEntity:            "Grass",
	HazardEntity:           "Hazard",
	TurretEntity:           "TurretEnemy",
	CatcherEntity:          "Catcher",
	PlatformEntity:         "OneWayPlatform",
	LanternEntity:          "Lantern",
	ChainNodeEntity:        "SlamboxChainNode",
	SlamboxChainEntity:     "SlamboxChain",
	TestSpeechBubbleEntity: "TestSpeechBubble",
	LevelTitleField:        "Title",
}
