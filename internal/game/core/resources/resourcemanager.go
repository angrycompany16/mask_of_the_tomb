package resources

// A major problem:
// This method enforces the dependence between different entites, which is not good
// For example, it requires there to always be a world publishing world data
// That's not too great

var (
	ResourceMan ResourceManager
)

type ResourceManager struct {
	PlayerData PlayerData
}

// Can also maybe contain events?
type PlayerData struct {
	// Contains all fields required for interactions with player entity
}

type WorldData struct {
	// Contains all fields required for interactions with world entity
}

type MenuData struct {
	// Contains all fields required for interactions with menu entity
}
