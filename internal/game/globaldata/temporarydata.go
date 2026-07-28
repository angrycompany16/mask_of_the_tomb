package globaldata

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"math/rand/v2"
)

type TemporaryData struct {
	LevelStates   map[string]LevelState
	GrassWindSeed int64
}

type LevelState struct {
	SlamboxPositions map[string]maths.Vec2
	PlayerSpawnPos   maths.Vec2
	PlayerSpawnDir   maths.Direction
}

func (t *TemporaryData) SaveLevelState(scene *engine.Scene) {
	levelState := t.LevelStates[scene.GetName()]

	slamboxes := scene.GetRoot().GetChildrenFunc(
		func(n *node.Node[engine.Actor]) bool {
			_, ok := engine.As[*slamboxactor.Slambox](n.GetValue())
			return ok
		},
	)

	for _, slambox := range slamboxes {
		slamboxactor, _ := engine.As[*slamboxactor.Slambox](slambox.GetValue())
		levelState.SlamboxPositions[slambox.GetID()] = maths.NewVec2(slamboxactor.GetPos())
	}
}

func NewLevelState(spawnX, spawnY float64, spawnDir maths.Direction) LevelState {
	return LevelState{
		PlayerSpawnPos:   maths.NewVec2(spawnX, spawnY),
		PlayerSpawnDir:   spawnDir,
		SlamboxPositions: make(map[string]maths.Vec2),
	}
}

func newTemporaryData() TemporaryData {
	return TemporaryData{
		LevelStates:   make(map[string]LevelState),
		GrassWindSeed: rand.Int64(),
	}
}
