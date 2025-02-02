package world

import (
	"fmt"
	"mask_of_the_tomb/ebitenLDTK"
	"mask_of_the_tomb/files"
	. "mask_of_the_tomb/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpaceLayerName   = "Playerspace"
	spawnPosEntityName     = "SpawnPosition"
	doorEntityName         = "Door"
	doorOtherSideFieldName = "OtherSide"
	collectibleLayerName   = "Collectibles"
)

// A note on level saving - this may need layered saving
// When to actually save data? We have some options:
//   - Just load the entire world into RAM :^)
//   - Save data whenever we change levels
//   - Save data on exit, but keep local changes when switching scenes
//   - Save data within some even interval (i.e. resting on a bench or smth), and then
//     also local changes when switching scenes

// Idea:
// Use a save struct as the local version of the world
// Store some sort of local diff of the world
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

	changeActiveLevel(w, 0)
}

func (w *World) Update() {
	// Anything...
}

func (w *World) Draw(surf *ebiten.Image, camX, camY float64) {
	w.ActiveLevel.draw(surf, camX, camY)
}

func (w *World) ExitByDoor(doorEntity ebitenLDTK.EntityInstance) (float64, float64) {
	entity, err := w.worldLDTK.Defs.GetEntityByUid(doorEntity.Uid)
	HandleLazy(err)
	if entity.Name != doorEntityName {
		return 0, 0 // Should honestly be error handled
	}
	for _, fieldInstance := range doorEntity.FieldInstances {
		if fieldInstance.Name != doorOtherSideFieldName {
			continue
		}

		nextLevel, err := w.worldLDTK.GetLevelByIid(fieldInstance.EntityRefValue.LevelIid)
		HandleLazy(err)
		// Change level
		changeActiveLevel(w, nextLevel.Uid)

		otherSideEntityRef, err := w.ActiveLevel.levelLDTK.GetEntityInstanceByIid(fieldInstance.EntityRefValue.EntityIid)
		HandleLazy(err)
		return otherSideEntityRef.Px[0], otherSideEntityRef.Px[1]
	}
	return 0, 0
}

func changeActiveLevel[T string | int](world *World, id T) error {
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
