package assetloader

import (
	"fmt"
	"io/fs"
)

// HOLY SHIT.
// There's a reason why I'm known as the pointer God...

type AssetStatus int

const (
	STAGED AssetStatus = iota
	LOADED
	FAILED
)

type Loadable interface {
	Load(fs fs.FS) (any, error)
}

// This is what we'll be using in our objects. When loading finishes,
// the pointer becomes valid.
type AssetRef[T any] struct {
	value  any
	status *AssetStatus
}

func (a *AssetRef[T]) Value() *T {
	return (*a.value.(*any)).(*T)
}

func (a *AssetRef[T]) Status() AssetStatus {
	return *a.status
}

type Asset struct {
	value    any
	loadable Loadable
	status   AssetStatus
}

type AssetLoader struct {
	assetpool map[string]*Asset
	fs        fs.FS
}

func StageAsset[T any](a *AssetLoader, name string, loadable Loadable) *AssetRef[T] {
	assetRef := AssetRef[T]{}
	if asset, ok := a.assetpool[name]; ok {
		assetRef.value = &asset.value
		assetRef.status = &asset.status
		return &assetRef
	}

	asset := Asset{
		loadable: loadable,
		status:   STAGED,
	}

	assetRef.value = &asset.value
	assetRef.status = &asset.status
	a.assetpool[name] = &asset
	return &assetRef
}

func (a *AssetLoader) LoadAll() {
	for _, asset := range a.assetpool {
		if asset.status == LOADED || asset.status == FAILED {
			continue
		}
		val, err := asset.loadable.Load(a.fs)
		if err != nil {
			fmt.Printf("Asset failed with error %s\n", err.Error())
			asset.status = FAILED
			continue
		}
		asset.status = LOADED
		asset.value = val
	}
}

func NewAssetLoader(fs fs.FS) *AssetLoader {
	return &AssetLoader{
		assetpool: make(map[string]*Asset),
		fs:        fs,
	}
}
