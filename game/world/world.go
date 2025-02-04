package world

import (
	"fmt"
	"mask_of_the_tomb/ebitenLDTK"
	"mask_of_the_tomb/files"

	// . "mask_of_the_tomb/utils"
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
	roomTransitionLayerName = "RoomTransitions"
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

// func (w *World) ExitByDoor(doorEntity ebitenLDTK.EntityInstance) (float64, float64) {
// 	entity, err := w.worldLDTK.Defs.GetEntityByUid(doorEntity.DefUid)
// 	HandleLazy(err)
// 	if entity.Name != doorEntityName {
// 		return 0, 0 // Should honestly be error handled
// 	}
// 	for _, fieldInstance := range doorEntity.FieldInstances {
// 		if fieldInstance.Name != doorOtherSideFieldName {
// 			continue
// 		}

// 		// entityRef, ok := fieldInstance.Value.
// 		// if !ok {
// 		// 	fmt.Println(fieldInstance.Value)
// 		// 	log.Fatal("could not get entityRef from fieldInstance in ExitByDoor")
// 		// }
// 		fmt.Println(fieldInstance.EntityRef)
// 		fmt.Println(fieldInstance.Float)
// 		nextLevel, err := w.worldLDTK.GetLevelByIid(fieldInstance.EntityRef.LevelIid)
// 		HandleLazy(err)
// 		// Change level
// 		ChangeActiveLevel(w, nextLevel.Uid)

// 		otherSideEntityRef, err := w.ActiveLevel.levelLDTK.GetEntityInstanceByIid(fieldInstance.EntityRef.EntityIid)
// 		HandleLazy(err)
// 		return otherSideEntityRef.Px[0], otherSideEntityRef.Px[1]
// 	}
// 	return 0, 0
// }

func ChangeActiveLevel[T string | int](world *World, id T) error {
	var newLevelLDTK ebitenLDTK.Level
	var err error

	switch v := any(id).(type) {
	case string:
		newLevelLDTK, err = world.worldLDTK.GetLevelByName(v)
		if err != nil {
			fmt.Println("Error when switching levels (id string)")
			return err
		}
	case int:
		newLevelLDTK, err = world.worldLDTK.GetLevelByUid(v)
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
