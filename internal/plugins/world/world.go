package world

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/resources"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

// Just an idea: With this system, could we turn each level into a scene?
// Hmm...
// This opens up some new and interesting possibilities

var (
	slamboxTilemapPath = filepath.Join(assets.EnvironmentFolder, "slambox_tilemap.png")
	particleSystemPath = filepath.Join("assets", "particlesystems", "basement.yaml")
)

const (
	firstLevelName = "Level_1"
)

type World struct {
	currentBiome string
	worldLDTK    *ebitenLDTK.World
	ActiveLevel  *Level
}

func NewWorld() *World {
	return &World{
		worldLDTK: errs.Must(assettypes.GetLDTKAsset("LDTKAsset")),
	}
}

func (w *World) Init(initLevelName string, gameData save.SaveData) {
	if initLevelName == "" {
		if gameData.SpawnRoomName == "" {
			initLevelName = firstLevelName
		} else {
			initLevelName = gameData.SpawnRoomName
		}
	}
	errs.Must(ChangeActiveLevel(w, initLevelName, ""))
}

func (w *World) Update(playerX, playerY, playerVelX, playerVelY float64) {
	w.ActiveLevel.Update(playerX, playerY, playerVelX, playerVelY)
}

// Set up a small loading stage when switching levels
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

	resources.PreviousLevelName = newLevel.name

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

func (w *World) ResetActiveLevel() (float64, float64) {
	_resetX, _resetY := w.ActiveLevel.GetResetPoint()

	// reset slambox positions
	w.ActiveLevel.reset()
	// reset player position
	w.ActiveLevel.resetX = _resetX
	w.ActiveLevel.resetY = _resetY
	return _resetX, _resetY
}
