package savemanager

import (
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/globaldata"
)

// Its only purpose is to call load and save lol
type SaveManager struct {
	*nodeactor.Node
}

func (s *SaveManager) Load(cmd *commands.Commands) {
	globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
	globaldata.Persist.Load(0, true)
	globaldata.Persist.Load(globaldata.Temp.CurrentSaveProfile, false)
}

func (s *SaveManager) Save(cmd *commands.Commands) {
	globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
	globaldata.Persist.Save(0, true)
	globaldata.Persist.Save(globaldata.Temp.CurrentSaveProfile, false)
}

func MakeSaveManager() *SaveManager {
	return &SaveManager{
		Node: nodeactor.NewNode(),
	}
}
