package assettypes

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/UI/display"
	"mask_of_the_tomb/internal/game/core/assetloader"
)

type MenuAsset struct {
	path string
	Menu display.Display
}

func (a *MenuAsset) Load() {
	a.Menu = *errs.Must(display.FromFile(a.path))
}

func NewMenuAsset(path string) *display.Display {
	menuAsset := MenuAsset{
		path: path,
	}

	assetloader.AddAsset(&menuAsset)
	return &menuAsset.Menu
}
