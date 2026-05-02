package gameinfo

import "time"

// TODO: Probably will be merged into gamestate
type GameInfo struct {
	Exit      bool
	startTime time.Time
}

func (g *GameInfo) GetTime() float64 {
	return float64(time.Since(g.startTime)) / float64(time.Second)
}

func NewGameInfo() *GameInfo {
	return &GameInfo{
		startTime: time.Now(),
	}
}
