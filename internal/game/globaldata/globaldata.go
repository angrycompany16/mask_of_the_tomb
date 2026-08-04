package globaldata

// Figure out how to separate persistent / non-persistent data
// idea: Add a 'scope' to all types of data: can be none, level, area, world, game
// should essentially represent when some piece of data is forgotten (so for instance,
// if the scope is level, the data is forgotten upon exiting the level)

// Master class for all types of saved data beyond just actors. In reality we only need to distinguish if the data should be persistent
// or not.
type GlobalData struct {
	Config  Config
	Persist PersistentData
	Temp    TemporaryData
}

func NewGlobalData(bypassSave bool) *GlobalData {
	return &GlobalData{
		Config:  newConfig(),
		Persist: newPersistentData(),
		Temp:    newTemporaryData(bypassSave),
	}
}
