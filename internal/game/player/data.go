package player

import "path/filepath"

var (
	playerFolder        = filepath.Join("assets", "sprites", "player", "export")
	PlayerSpritePath    = filepath.Join(playerFolder, "player-concept.png")
	IdleSpritesheetPath = filepath.Join(playerFolder, "player-animations-Sheet.png")
)
