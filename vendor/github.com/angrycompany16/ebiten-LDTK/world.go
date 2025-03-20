package ebitenLDTK

import (
	"encoding/json"
	"fmt"
	"os"
)

type World struct {
	GridWidth  int     `json:"worldGridWidth"`
	GridHeight int     `json:"worldGridHeight"`
	Defs       Defs    `json:"defs"`
	Levels     []Level `json:"levels"`
}

func LoadWorld(path string) (World, error) {
	world := World{}
	file, err := os.Open(path)
	if err != nil {
		return world, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&world)
	if err != nil {
		return world, err
	}
	return world, nil
}

func (w *World) GetLevelByUid(uid int) (Level, error) {
	for _, level := range w.Levels {
		if level.Uid == uid {
			return level, nil
		}
	}
	return Level{}, fmt.Errorf("level with uid [%d] was not found", uid)
}

func (w *World) GetLevelByName(name string) (Level, error) {
	for _, level := range w.Levels {
		if level.Name == name {
			return level, nil
		}
	}
	return Level{}, fmt.Errorf("level with name [%s] was not found", name)
}

func (w *World) GetLevelByIid(iid string) (Level, error) {
	for _, level := range w.Levels {
		if level.Iid == iid {
			return level, nil
		}
	}
	return Level{}, fmt.Errorf("level with iid [%s] was not found", iid)
}
