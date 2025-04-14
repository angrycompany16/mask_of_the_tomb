package world

import (
	"fmt"
	"mask_of_the_tomb/internal/game/core/assetloader/assettypes"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

var (
	InitLevelIid = "Level_0"
)

type World struct {
	worldLDTK   *ebitenLDTK.World
	ActiveLevel *Level
}

func (w *World) Load() {
	w.worldLDTK = assettypes.NewLDTKAsset(LDTKMapPath)
	// assetloader.AddAsset(worldAsset)
	// w.worldLDTK = &worldAsset.World
}

func (w *World) Init() {
	ChangeActiveLevel(w, InitLevelIid)
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
	world.ActiveLevel = newLevel
	return nil
}

func (w *World) ResetActiveLevel() (float64, float64) {
	levelLDTK := *w.ActiveLevel.levelLDTK
	defs := *w.ActiveLevel.defs
	newLevel, err := newLevel(&levelLDTK, &defs)
	if err != nil {
		panic(fmt.Sprintln("Could not reset level", w.ActiveLevel.name))
	}

	w.ActiveLevel = newLevel
	return w.ActiveLevel.GetSpawnPoint()
}
