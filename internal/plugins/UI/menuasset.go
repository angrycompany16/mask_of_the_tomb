package ui

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
)

type menuAsset struct {
	path string
	Menu Display
}

func (a *menuAsset) Load() {
	a.Menu = *errs.Must(FromFile(a.path))
}

func NewMenuAsset(path string) *Display {
	asset, exists := assetloader.Exists(path)
	if exists {
		return &asset.(*menuAsset).Menu
	}

	menuAsset := menuAsset{
		path: path,
	}

	assetloader.Load(path, &menuAsset)
	return &menuAsset.Menu
}
