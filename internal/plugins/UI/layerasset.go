package ui

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
)

type layerAsset struct {
	path  string
	layer Layer
}

func (a *layerAsset) Load() {
	a.layer = *errs.Must(FromFile(a.path))
}

func NewLayerAsset(path string) *Layer {
	asset, exists := assetloader.Exists(path)
	if exists {
		return &asset.(*layerAsset).layer
	}

	layerAsset := layerAsset{
		path: path,
	}

	assetloader.Load(path, &layerAsset)
	return &layerAsset.layer
}
