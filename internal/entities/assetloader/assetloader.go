package assetloader

import "path/filepath"

// Useful file paths
var (
	EnvironmentTilemapFolder = filepath.Join("assets", "sprites", "environment", "tilemaps", "export")
	PlayerFolder             = filepath.Join("assets", "sprites", "player", "export")
)

type Asset interface {
	load() // Very simple asset loader interface
	// Can now create loaders for images, ldtk etc..
}

type ImageAsset struct {
}

type FontAsset struct {
}

type ShaderAsset struct {
}

type ParticleSystemAsset struct {
}

type AnimationAsset struct {
}

type DelayAsset struct { // Debug asset which takes a long time to load
}

type AssetLoader struct {
	// Loading screen
	// Asset pool []*Asset
}

// How should this system work?
// Define load events where the asset loader will spawn a thread to load a bunch of
// assets
// When the asset loader is done it will signal this on some channel so that the
// game can switch states and start doing things again
// The asset loading periods will typically be level switching and first loading the
// game, however there may also be other cases such as when changing certain settings
// Loading assets at different points should not be hard

// Idea:
// We use an 'asset pool', which gets filled with objects. Then, when all the objects
// have been found, we call load_all() in a different thread, and display some kind of
// loading screen while the loader loads all of the assets into memory.
