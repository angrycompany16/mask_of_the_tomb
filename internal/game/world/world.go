package world

import (
	"fmt"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/assetloader/assettypes"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

type World struct {
	worldLDTK   *ebitenLDTK.World
	ActiveLevel *Level
	levelMemory map[string]levelMemory
}

func NewWorld() *World {
	return &World{
		levelMemory: make(map[string]levelMemory),
	}
}

func (w *World) Load() {
	w.worldLDTK = assettypes.NewLDTKAsset(LDTKMapPath)
}

func (w *World) Init(initLevelName string) {
	if initLevelName == "" {
		initLevelName = w.worldLDTK.Levels[0].Name
	}
	ChangeActiveLevel(w, initLevelName, "")
}

func (w *World) Update() {
	w.ActiveLevel.Update()
}

func ChangeActiveLevel[T string | int](world *World, id T, doorIid string) error {
	if world.ActiveLevel != nil {
		world.levelMemory[world.ActiveLevel.levelLDTK.Iid] = levelMemory{world.ActiveLevel.slamboxes}
	}

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

	if memory, ok := world.levelMemory[newLevelLDTK.Iid]; ok {
		newLevel.restoreFromMemory(&memory)
	}

	newLevel.resetX, newLevel.resetY = newLevel.GetDefaultSpawnPoint()
	if doorIid != "" {
		doorEntity := errs.Must(newLevel.levelLDTK.GetEntityByIid(doorIid))
		newLevel.resetX, newLevel.resetY = doorEntity.Px[0], doorEntity.Px[1]
	}
	world.ActiveLevel = newLevel
	return nil
}

func (w *World) ResetActiveLevel() (float64, float64) {
	_resetX, _resetY := w.ActiveLevel.GetResetPoint()
	levelLDTK := *w.ActiveLevel.levelLDTK
	defs := *w.ActiveLevel.defs
	newLevel, err := newLevel(&levelLDTK, &defs)
	if err != nil {
		panic(fmt.Sprintln("Could not reset level", w.ActiveLevel.name))
	}

	w.ActiveLevel = newLevel
	w.ActiveLevel.resetX = _resetX
	w.ActiveLevel.resetY = _resetY
	return _resetX, _resetY
}
