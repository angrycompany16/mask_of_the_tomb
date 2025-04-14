package assetloader

import (
	"path/filepath"
)

// Useful file paths
// TODO: Move to assets.go
// Or just rename into embeds.go or something
var (
	EnvironmentTilemapFolder = filepath.Join("assets", "sprites", "environment", "tilemaps", "export")
	PlayerFolder             = filepath.Join("assets", "sprites", "player", "export")
	_assetLoader             = assetLoader{}
)

type assetLoader struct {
	assetPool []Asset
}

func AddAsset(assets ...Asset) {
	_assetLoader.assetPool = append(_assetLoader.assetPool, assets...)
}

func LoadAll(doneChan chan<- int) {
	for _, asset := range _assetLoader.assetPool {
		asset.Load()
	}
	doneChan <- 1
}
