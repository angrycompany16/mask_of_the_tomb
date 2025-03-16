package world

import (
	"mask_of_the_tomb/internal/game/assetloader"
	"path/filepath"
)

var (
	LDTKMapPath        = filepath.Join("assets", "LDTK", "world.ldtk")
	SlamboxTilemapPath = filepath.Join(assetloader.EnvironmentTilemapFolder, "playerspace_tilemap.png")
)
