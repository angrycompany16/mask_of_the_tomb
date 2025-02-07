package world

import (
	"fmt"
	"mask_of_the_tomb/ebitenLDTK"
	"mask_of_the_tomb/files"

	"path"
)

const (
	playerSpaceLayerName    = "Playerspace"
	foreGroundLayerName     = "Foreground"
	spawnPosEntityName      = "SpawnPosition"
	doorEntityName          = "Door"
	doorOtherSideFieldName  = "OtherSide"
	collectibleLayerName    = "Collectibles"
	hazardLayerName         = "Hazards"
	hazardDamageFieldName   = "Damage"
	roomTransitionLayerName = "Doors"
	spawnPointLayerName     = "SpawnPoint"
)

type World struct {
	worldLDTK   ebitenLDTK.World
	ActiveLevel *Level
}

func (w *World) Init() {
	w.worldLDTK = *files.LazyLDTK(LDTKMapPath)

	// One folder back to access LDTK folder
	LDTKpath := path.Clean(path.Join(LDTKMapPath, ".."))
	for i := 0; i < len(w.worldLDTK.Defs.Tilesets); i++ {
		tileset := &w.worldLDTK.Defs.Tilesets[i]
		tilesetPath := path.Join(LDTKpath, tileset.RelPath)
		tileset.Image = files.LazyImage(tilesetPath)
	}

	ChangeActiveLevel(w, 0)
}

func (w *World) Update() {
	// Anything...
}

func ChangeActiveLevel[T string | int](world *World, id T) error {
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
				return Ierr
			}
		}
	case int:
		newLevelLDTK, err = world.worldLDTK.GetLevelByUid(levelId)
		if err != nil {
			fmt.Println("Error when switching levels (id int)")
			return err
		}
	}

	newLevel, err := newLevel(&newLevelLDTK, &world.worldLDTK.Defs)
	if err != nil {
		return err
	}
	world.ActiveLevel = &newLevel
	return nil
}
