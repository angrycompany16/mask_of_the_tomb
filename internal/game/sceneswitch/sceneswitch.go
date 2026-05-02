package sceneswitch

import "mask_of_the_tomb/internal/backend/maths"

// TODO: Can probably merge this into gameState as well
type SceneSwitch struct {
	SpawnEntityIid string
	SpawnDirection maths.Direction
}
