package player

import "path/filepath"

var (
	playerFolder        = filepath.Join("assets", "sprites", "player", "export")
	PlayerSpritePath    = filepath.Join(playerFolder, "player.png")
	IdleSpritesheetPath = filepath.Join(playerFolder, "player-idle-animation-Sheet.png")
)
