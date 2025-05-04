package entities

import (
	"mask_of_the_tomb/internal/game/core/assetloader"
	"path/filepath"
)

var (
	SlamboxTilemapPath = filepath.Join(assetloader.EnvironmentTilemapFolder, "playerspace_tilemap.png")
	doorSpritePath     = filepath.Join("assets", "sprites", "environment", "entities", "export", "teleporter.png")
)
