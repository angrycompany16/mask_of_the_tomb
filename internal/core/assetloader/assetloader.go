package assetloader

import (
	"errors"
	"fmt"
)

// How to make this work:
// Asset loader has a channel for each asset
// Getting an asset correspons to finding the proper channel and listening
// Those channels will typically be handing out pointers to their respective assets'
// data
// But that doesn't help with thread-safety... two processes could still read
// a resource at the same time...

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
	for _, asset := range _assetLoader.assetPool {
		err := asset.Load()
		if err != nil {
			asset.status = fmt.Sprintf("FAILED: Failed with error %s\n", err.Error())
		} else {
			asset.status = "SUCCESS"
		}
	}
	doneChan <- 1
}
