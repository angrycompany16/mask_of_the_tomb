package world

import (
	"fmt"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
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

type LevelSwapCtx struct {
	LevelIid      string
	DoorEntityIid string
}

type World struct {
	*assettypes.LDTKAsset
	currentBiome string
	ActiveLevel  *Level
	LevelSwapCtx LevelSwapCtx
}

func NewWorld() *World {
	return &World{
		LDTKAsset: errs.Must(assettypes.GetLDTKAsset("LDTKAsset")),
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
		newLevelLDTK, err = world.World.GetLevelByName(levelId)
		if err != nil {
			fmt.Println("Couldn't switch levels by name (id string), trying Iid...")
			var Ierr error
			newLevelLDTK, Ierr = world.World.GetLevelByIid(levelId)
			if Ierr != nil {
				fmt.Println("Error when switching levels by Iid (id string)")
				return "", Ierr
			}
		}
	case int:
		newLevelLDTK, err = world.World.GetLevelByUid(levelId)
		if err != nil {
			fmt.Println("Error when switching levels (id int)")
			return "", err
		}
	}

	newLevel, err := newLevel(&newLevelLDTK, &world.World.Defs, world.Tilesets)
	if err != nil {
		return "", err
	}

	resources.PreviousLevelName = newLevel.name

	newLevel.resetX, newLevel.resetY, newLevel.resetDir = newLevel.GetDefaultSpawnInfo()
	if doorIid != "" {
		newLevel.resetX, newLevel.resetY, newLevel.resetDir = newLevel.ResetInfoFromIid(doorIid)
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

func (w *World) ResetActiveLevel() (float64, float64, maths.Direction) {
	resetX, resetY, resetDir := w.ActiveLevel.GetResetInfo()

	// reset slambox positions
	w.ActiveLevel.reset()
	// set reset position O_o
	// w.ActiveLevel.resetX = resetX
	// w.ActiveLevel.resetY = resetY
	// w.ActiveLevel.resetDir = resetDir
	return resetX, resetY, resetDir
}
