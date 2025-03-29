package assetloader

import (
	"path/filepath"
	"slices"
)

// Useful file paths
// TODO: Move to assets.go
// Or just rename into embeds.go or something
var (
	EnvironmentTilemapFolder = filepath.Join("assets", "sprites", "environment", "tilemaps", "export")
	PlayerFolder             = filepath.Join("assets", "sprites", "player", "export")
	_assetLoader             = assetLoader{}
)

type ShaderAsset struct {
}

type AnimationAsset struct {
}

type assetLoader struct {
	assetPool    []Asset
	loadedAssets []Asset
}

func AddAsset(assets ...Asset) {
	_assetLoader.assetPool = append(_assetLoader.assetPool, assets...)
}

func LoadAll(doneChan chan<- int) {
	for i, asset := range _assetLoader.assetPool {
		asset.Load()
		_assetLoader.loadedAssets = append(_assetLoader.loadedAssets, asset)
		_assetLoader.assetPool = slices.Delete(_assetLoader.assetPool, i, i)
	}
	doneChan <- 1
}
