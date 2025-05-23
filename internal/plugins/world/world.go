package world

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/libraries/assettypes"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

var (
	slamboxTilemapPath = filepath.Join(assets.EnvironmentTilemapFolder, "slambox_tilemap.png")
	// slamboxTilemap     *ebiten.Image
)

const (
	firstLevelName = "Level_1"
)

type World struct {
	currentBiome     string
	worldLDTK        *ebitenLDTK.World
	ActiveLevel      *Level
	worldStateMemory map[string]LevelMemory
}

func NewWorld() *World {
	return &World{
		worldStateMemory: make(map[string]LevelMemory),
	}
}

func (w *World) Load(LDTKMapPath string) {
	w.worldLDTK = assettypes.NewLDTKAsset(LDTKMapPath)
	assetloader.Load("slamboxTilemap", assettypes.MakeImageAsset(assets.Slambox_tilemap))
	assetloader.Load("grassTilemap", assettypes.MakeImageAsset(assets.Grass_tiles))
}

func (w *World) Init(initLevelName string, gameData save.GameData) {
	if initLevelName == "" {
		if gameData.SpawnRoomName == "" {
			initLevelName = firstLevelName
		} else {
			initLevelName = gameData.SpawnRoomName
		}
	}

	// for id, levelmemory := range gameData.WorldStateMemory {
	// w.worldStateMemory[id] = levelmemory
	// }
	ChangeActiveLevel(w, initLevelName, "")
}

func (w *World) LoadMemory(memory map[string]LevelMemory) {
}

func (w *World) Update(playerX, playerY, playerVelX, playerVelY float64) {
	w.ActiveLevel.Update(playerX, playerY, playerVelX, playerVelY)
}

func (w *World) GetWorldStateMemory() map[string]LevelMemory {
	return w.worldStateMemory
}

func ChangeActiveLevel[T string | int](world *World, id T, doorIid string) (string, error) {
	var newLevelLDTK ebitenLDTK.Level
	var err error

	switch levelId := any(id).(type) {
	case string:
		newLevelLDTK, err = world.worldLDTK.GetLevelByName(levelId)
		if err != nil {
			fmt.Println("Couldn't switch levels by name (id string), trying Iid...")
			var Ierr error
			newLevelLDTK, Ierr = world.worldLDTK.GetLevelByIid(levelId)
			if Ierr != nil {
				fmt.Println("Error when switching levels by Iid (id string)")
				return "", Ierr
			}
		}
	case int:
		newLevelLDTK, err = world.worldLDTK.GetLevelByUid(levelId)
		if err != nil {
			fmt.Println("Error when switching levels (id int)")
			return "", err
		}
	}

	newLevel, err := newLevel(&newLevelLDTK, &world.worldLDTK.Defs)
	if err != nil {
		return "", err
	}

	if memory, ok := world.worldStateMemory[newLevelLDTK.Iid]; ok {
		newLevel.restoreFromMemory(&memory)
	}

	if world.ActiveLevel != nil {
		world.SaveLevel(world.ActiveLevel)
	}
	world.SaveLevel(newLevel)

	newLevel.resetX, newLevel.resetY = newLevel.GetDefaultSpawnPoint()
	if doorIid != "" {
		doorEntity := errs.Must(newLevel.levelLDTK.GetEntityByIid(doorIid))
		newLevel.resetX, newLevel.resetY = doorEntity.Px[0], doorEntity.Px[1]
		// play biome animation if updated
		if newLevel.GetBiome() != world.currentBiome {
			world.ActiveLevel = newLevel
			world.currentBiome = newLevel.GetBiome()
			return newLevel.GetBiome(), err
		}
	}
	world.ActiveLevel = newLevel
	world.currentBiome = newLevel.GetBiome()
	return "", nil
}

func (w *World) SaveLevel(level *Level) {
	// w.worldStateMemory[level.levelLDTK.Iid] = levelmemory.LevelMemory{level.GetSlamboxPositions()}
}

func (w *World) ResetActiveLevel() (float64, float64) {
	_resetX, _resetY := w.ActiveLevel.GetResetPoint()

	// reset slambox positions
	level := w.worldStateMemory[w.ActiveLevel.levelLDTK.Iid]
	w.ActiveLevel.restoreFromMemory(&level)
	// reset player position
	w.ActiveLevel.resetX = _resetX
	w.ActiveLevel.resetY = _resetY
	return _resetX, _resetY
}
