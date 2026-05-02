package gamestate

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
)

// Figure out how to separate persistent / non-persistent data
type Config struct {
	SoundVolume float64
	MusicVolume float64
	SfxVolume   float64
}

type LevelState struct {
	SlamboxPositions map[string]maths.Vec2
	PlayerSpawnPos   maths.Vec2
	PlayerSpawnDir   maths.Direction
}

type GameState struct {
	Config      Config
	LevelStates map[string]LevelState
}

func (g *GameState) SaveLevelState(scene *engine.Scene) {
	levelState := g.LevelStates[scene.GetName()]

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

func NewGameState() *GameState {
	return &GameState{
		Config:      Config{},
		LevelStates: make(map[string]LevelState),
	}
}

func NewLevelState(spawnX, spawnY float64, spawnDir maths.Direction) LevelState {
	return LevelState{
		PlayerSpawnPos:   maths.NewVec2(spawnX, spawnY),
		PlayerSpawnDir:   spawnDir,
		SlamboxPositions: make(map[string]maths.Vec2),
	}
}
