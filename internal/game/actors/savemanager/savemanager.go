package savemanager

import (
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/globaldata"
)

// Its only purpose is to call load on startup and save on exit lol
type SaveManager struct {
	*nodeactor.Node
}

func (s *SaveManager) Init(cmd *commands.Commands) {
	s.Node.Init(cmd)

	globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
	globaldata.Persist.Load(0, true)
}

func (s *SaveManager) Save()

func MakeSaveManager() *SaveManager {
	return &SaveManager{
		Node: nodeactor.NewNode(),
	}
}
