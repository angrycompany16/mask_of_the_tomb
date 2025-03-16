package world

import (
	"fmt"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/entities"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

// TODO: Maybe find a better place for this

var (
	_world       = world{}
	initLevelIid = "Level_0"
)

// TODO: Implement some kind of readonly wrapper so that we can pass things in without worrying about whether
// they are being mutated?
type world struct {
	worldLDTK   ebitenLDTK.World
	activeLevel *Level
}

type worldInitInfo struct {
	LevelWidth, LevelHeight float64
	SpawnX, SpawnY          float64
}

// TODO: Fix game-breaking bug with slamboxes

func Init() worldInitInfo {
	entities.RegisterEntity(&_world, "World")
	_world.worldLDTK = errs.Must(ebitenLDTK.LoadWorld(LDTKMapPath))

	// One folder back to access LDTK folder
	LDTKpath := filepath.Clean(filepath.Join(LDTKMapPath, ".."))
	for i := 0; i < len(_world.worldLDTK.Defs.Tilesets); i++ {
		tileset := &_world.worldLDTK.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		// fmt.Println(tilesetPath)
		tileset.Image = errs.MustNewImageFromFile(tilesetPath)
	}
	level := ChangeActiveLevel(&_world, initLevelIid)
	_world.activeLevel = &level

	spawnX, spawnY := _world.activeLevel.GetSpawnPoint()
	return worldInitInfo{
		LevelWidth:  _world.activeLevel.levelLDTK.PxWid,
		LevelHeight: _world.activeLevel.levelLDTK.PxHei,
		SpawnX:      spawnX,
		SpawnY:      spawnY,
	}
}

func (w *world) Update() {
	w.activeLevel.Update()
}

func ChangeActiveLevel[T string | int](world *world, id T) Level {
	var newLevelLDTK ebitenLDTK.Level
	var err error

	switch levelId := any(id).(type) {
	case string:
		newLevelLDTK, err = world.worldLDTK.GetLevelByName(levelId)
		if err != nil {
			fmt.Println("Couldn't switch levels by name (id string), trying Iid...")
			newLevelLDTK = errs.Must(world.worldLDTK.GetLevelByIid(levelId))
		}
	case int:
		newLevelLDTK = errs.Must(world.worldLDTK.GetLevelByUid(levelId))
	}

	newLevel := errs.Must(newLevel(&newLevelLDTK, &world.worldLDTK.Defs))
	return newLevel
}
