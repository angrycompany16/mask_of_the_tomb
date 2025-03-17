package world

import (
	"fmt"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/entities/level"
	"mask_of_the_tomb/internal/libraries/errs"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

var (
	initLevelIid = "Level_0"
)

type world struct {
	*entities.Entity
	worldLDTK ebitenLDTK.World
}

func NewWorld() (*world, level.InitLevelInfo) {
	_world := world{
		worldLDTK: errs.Must(ebitenLDTK.LoadWorld(LDTKMapPath)),
	}
	_world.Entity = entities.RegisterEntity(&_world, "world")

	LDTKpath := filepath.Clean(filepath.Join(LDTKMapPath, ".."))

	for i := 0; i < len(_world.worldLDTK.Defs.Tilesets); i++ {
		tileset := &_world.worldLDTK.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		tileset.Image = errs.MustNewImageFromFile(tilesetPath)
	}

	firstLevel := GetLevel(&_world, initLevelIid)
	_level, initLevelInfo, _ := level.NewLevel(firstLevel, &_world.worldLDTK.Defs)
	_world.AddChild(_level.Entity)

	return &_world, initLevelInfo
}

func (w *world) Update() {
	// adv := advertisers.GetAdvertiser(pubgame.GameEntityName)
	// val := adv.Read().(pubgame.GameAdvertiser)

	// if val.State != pubgame.StatePlaying {
	// 	return
	// }
}

func (w *world) Draw() {
	// adv := advertisers.GetAdvertiser(pubgame.GameEntityName)
	// val := adv.Read().(pubgame.GameAdvertiser)
	// if val.State == pubgame.StateMainMenu {
	// 	return
	// }
}

func GetLevel[T string | int](world *world, id T) *ebitenLDTK.Level {
	var level ebitenLDTK.Level
	var err error

	switch levelId := any(id).(type) {
	case string:
		level, err = world.worldLDTK.GetLevelByName(levelId)
		if err != nil {
			fmt.Println("Couldn't switch levels by name (id string), trying Iid...")
			level = errs.Must(world.worldLDTK.GetLevelByIid(levelId))
		}
	case int:
		level = errs.Must(world.worldLDTK.GetLevelByUid(levelId))
	}

	return &level
}
