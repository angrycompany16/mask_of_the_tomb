package assetloader

import (
	"errors"
	"fmt"
)

var (
	_assetLoader = assetLoader{
		assetPool: make(map[string]*AssetEntry),
	}
)

type Asset interface {
	Load() error
}

type AssetEntry struct {
	Asset
	status string
}

type assetLoader struct {
	assetPool map[string]*AssetEntry
}

func Exists(hash string) (Asset, bool) {
	asset, ok := _assetLoader.assetPool[hash]
	return asset.Asset, ok
}

func Add(hash string, asset Asset) {
	_assetLoader.assetPool[hash] = &AssetEntry{asset, "NOT LOADED"}
}

func GetAsset(name string) (Asset, error) {
	asset, ok := _assetLoader.assetPool[name]
	if !ok || asset.status == "NOT LOADED" {
		fmt.Printf("Could not find asset with name '%s'", name)
		PrintAssetRegistry()
		return nil, errors.New("AssetNotFound")
	}
	return asset.Asset, nil
}

func PrintAssetRegistry() {
	fmt.Printf("\n ---- ASSET REGISTRY ---- \n")
	for name, asset := range _assetLoader.assetPool {
		fmt.Println(name, asset.status)
	}
	fmt.Printf(" ------------------------ \n\n")
}

func LoadAll(doneChan chan<- int) {
	for name, asset := range _assetLoader.assetPool {
		err := asset.Load()
		if err != nil {
			fmt.Printf("Asset with name %s failed with error %s\n", name, err.Error())
			asset.status = "FAILED"
		} else {
			asset.status = "SUCCESS"
		}
	}
	doneChan <- 1
}
