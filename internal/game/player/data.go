package player

import (
	"mask_of_the_tomb/internal/game/core/assetloader"
	"path/filepath"
)

var (
	playerSpritePath        = filepath.Join(assetloader.PlayerFolder, "player.png")
	idleSpritesheetPath     = filepath.Join(assetloader.PlayerFolder, "player-idle-Sheet.png")
	dashInitSpritesheetPath = filepath.Join(assetloader.PlayerFolder, "player-init-jump-Sheet.png")
	dashLoopSpritesheetPath = filepath.Join(assetloader.PlayerFolder, "player-loop-jump-Sheet.png")
	slamSpritesheetPath     = filepath.Join(assetloader.PlayerFolder, "player-slam-Sheet.png")
	jumpParticlesBroadPath  = filepath.Join("assets", "particlesystems", "player", "jump-broad.yaml")
	jumpParticlesTightPath  = filepath.Join("assets", "particlesystems", "player", "jump-tight.yaml")
)
