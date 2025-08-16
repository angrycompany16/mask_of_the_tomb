package ui

import (
	"mask_of_the_tomb/internal/core/assetloader"
)

type layerAsset struct {
	path  string
	layer Layer
}

func (a *layerAsset) Load() error {
	layer, err := FromFile(a.path)
	a.layer = *layer
	return err
}

func NewLayerAsset(path string) *Layer {
	asset, exists := assetloader.Exists(path)
	if exists {
		return &asset.(*layerAsset).layer
	}

	layerAsset := layerAsset{
		path: path,
	}

	assetloader.Add(path, &layerAsset)
	return &layerAsset.layer
}
