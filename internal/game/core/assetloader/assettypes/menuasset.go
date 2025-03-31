package assettypes

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/UI/menu"
)

type MenuAsset struct {
	path string
	Menu menu.Menu
}

func (a *MenuAsset) Load() {
	a.Menu = *errs.Must(menu.FromFile(a.path))
}

func NewMenuAsset(path string) *MenuAsset {
	return &MenuAsset{
		path: path,
	}
}
