package globals

import "time"

// I don't really think that this is very great but oh well
type Globals struct {
	gameStartTime   time.Time
	sceneSwitchTime time.Time
}

func (g *Globals) SecondsSinceInit() float64 {
	return time.Since(g.gameStartTime).Seconds()
}

func (g *Globals) SecondsSinceSceneSwitch() float64 {
	return time.Since(g.sceneSwitchTime).Seconds()
}

func NewGlobals() *Globals {
	return &Globals{
		gameStartTime: time.Now(),
	}
}
