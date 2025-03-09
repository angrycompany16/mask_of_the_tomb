package world

import (
	"fmt"
	"mask_of_the_tomb/internal/errs"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

const (
	playerSpaceLayerName    = "Playerspace"
	spawnPosEntityName      = "SpawnPosition"
	doorEntityName          = "Door"
	doorOtherSideFieldName  = "OtherSide"
	hazardLayerName         = "Hazards"
	hazardDamageFieldName   = "Damage"
	roomTransitionLayerName = "Doors"
	spawnPointLayerName     = "SpawnPoint"
	slamboxLayerName        = "Slamboxes"
)

type World struct {
	worldLDTK   ebitenLDTK.World
	ActiveLevel *Level
}

func (w *World) Init() {
	w.worldLDTK = errs.Must(ebitenLDTK.LoadWorld(LDTKMapPath))

	// One folder back to access LDTK folder
	LDTKpath := filepath.Clean(filepath.Join(LDTKMapPath, ".."))
	for i := 0; i < len(w.worldLDTK.Defs.Tilesets); i++ {
		tileset := &w.worldLDTK.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		tileset.Image = errs.MustNewImageFromFile(tilesetPath)
	}

	ChangeActiveLevel(w, 0)
}

func (w *World) Update() {
	w.ActiveLevel.Update()
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
