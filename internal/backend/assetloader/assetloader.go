package assetloader

import (
	"fmt"
	"io/fs"

	om "github.com/wk8/go-ordered-map/v2"
)

// HOLY SHIT.
// There's a reason why I'm known as the pointer God...

//go:generate stringer -type=AssetStatus
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

func (a *Asset) GetStatusString() string {
	return a.status.String()
}

type AssetLoader struct {
	assetpool *om.OrderedMap[string, *Asset]
	fs        fs.FS
}

func StageAsset[T any](a *AssetLoader, name string, loadable Loadable) *AssetRef[T] {
	assetRef := AssetRef[T]{}
	if asset, ok := a.assetpool.Get(name); ok {
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
	a.assetpool.Set(name, &asset)
	return &assetRef
}

func (a *AssetLoader) LoadAll() {
	for pair := a.assetpool.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value.status == LOADED || pair.Value.status == FAILED {
			continue
		}
		val, err := pair.Value.loadable.Load(a.fs)
		if err != nil {
			fmt.Printf("Asset failed with error %s\n", err.Error())
			pair.Value.status = FAILED
			continue
		}
		pair.Value.status = LOADED
		pair.Value.value = val
	}
}

func (a *AssetLoader) GetAssetPool() *om.OrderedMap[string, *Asset] {
	return a.assetpool
}

func NewAssetLoader(fs fs.FS) *AssetLoader {
	return &AssetLoader{
		assetpool: om.New[string, *Asset](),
		fs:        fs,
	}
}
