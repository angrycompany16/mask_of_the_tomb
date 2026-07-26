package gamestate

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"math/rand/v2"
)

// Figure out how to separate persistent / non-persistent data
// idea: Add a 'scope' to all types of data: can be none, level, area, world, game
// should essentially represent when some piece of data is forgotten (so for instance,
// if the scope is level, the data is forgotten upon exiting the level)
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

type Inventory struct {
	keys []*Key
}

type Key struct {
	opensIid string
}

type GameState struct {
	Config        Config
	LevelStates   map[string]LevelState
	GrassWindSeed int64
	Inventory     *Inventory
	UnlockedDoors map[string]bool
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

func (i *Inventory) AddKey(opensIid string) {
	i.keys = append(i.keys, &Key{opensIid: opensIid})
	fmt.Println("Added key, current inventory")
	i.Pritn()
}

func (i *Inventory) HasKey(doorIid string) bool {
	for _, key := range i.keys {
		if key.opensIid == doorIid {
			return true
		}
	}
	return false
}

func (i *Inventory) Pritn() {
	fmt.Println("Printing inventory...")
	for _, key := range i.keys {
		fmt.Println("A key, opening", key.opensIid)
	}
}

func NewGameState() *GameState {
	return &GameState{
		Config:        Config{},
		LevelStates:   make(map[string]LevelState),
		GrassWindSeed: rand.Int64(),
		Inventory: &Inventory{
			keys: make([]*Key, 0),
		},
	}
}

func NewLevelState(spawnX, spawnY float64, spawnDir maths.Direction) LevelState {
	return LevelState{
		PlayerSpawnPos:   maths.NewVec2(spawnX, spawnY),
		PlayerSpawnDir:   spawnDir,
		SlamboxPositions: make(map[string]maths.Vec2),
	}
}
