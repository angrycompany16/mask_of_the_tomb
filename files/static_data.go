package files

import (
	"path/filepath"
)

// path registry for assets
var (
	// Sprites
	// PlayerSpritePath = filepath.Join("assets", "kenney_pixel-platformer", "Tiles", "Characters", "tile_0000.png")
	PlayerSpritePath = filepath.Join("assets", "aseprite", "player-test.png")
	TileSpritePath   = filepath.Join("assets", "kenney_pixel-platformer", "Tiles", "tile_0000.png")

	// LDTK
	LDTKMapPath = filepath.Join("assets", "LDTK", "test.ldtk")
)
