package assetloader

import (
	"errors"
	"fmt"
)

var (
	_assetLoader = assetLoader{
		assetPool: make(map[string]Asset),
	}
)

type Asset interface {
	Load() error
}

type assetLoader struct {
	assetPool map[string]Asset
}

func Exists(hash string) (Asset, bool) {
	asset, ok := _assetLoader.assetPool[hash]
	return asset, ok
}

func Load(hash string, asset Asset) {
	_assetLoader.assetPool[hash] = asset
}

func GetAsset(name string) (Asset, error) {
	asset, ok := _assetLoader.assetPool[name]
	if !ok {
		fmt.Printf("Could not find asset with name '%s'", name)
		PrintAssetRegistry()
		return nil, errors.New("AssetNotFound")
	}
	return asset, nil
}

func PrintAssetRegistry() {
	fmt.Printf("\n ---- ASSET REGISTRY ---- \n")
	for name := range _assetLoader.assetPool {
		fmt.Println(name)
	}
	fmt.Printf(" ------------------------ \n\n")
}

func LoadAll(doneChan chan<- int) {
	for name, asset := range _assetLoader.assetPool {
		err := asset.Load()
		if err != nil {
			fmt.Printf("Asset with name %s failed with error %s\n", name, err.Error())
		}
	}
	doneChan <- 1
}
