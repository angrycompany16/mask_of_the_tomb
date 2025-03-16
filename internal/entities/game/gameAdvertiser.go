package game

type gameAdvertiser struct {
	state GameState
}

func makeAdvertiser() *gameAdvertiser {
	return &gameAdvertiser{
		state: _game.state,
	}
}

// Q: What should the return type be? In the interface definition we only have one Read() method, but
// we might want to return multiple pieces of data, however we of course don't want to return them
// as a reference

// So essentially this function reads all the relevant values from the gamestate and then returns
// them as copies
func (g *gameAdvertiser) Read() any {
	g.state = _game.state
	return *g
}
