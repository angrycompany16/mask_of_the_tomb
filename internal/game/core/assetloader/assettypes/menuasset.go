package assettypes

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/UI/menu"
)

type MenuAsset struct {
	path string
	menu menu.Menu
}

func (a *MenuAsset) Load() {
	a.menu = *errs.Must(menu.FromFile(a.path))
}

func NewMenuAsset(path string) *ParticleSystemAsset {
	return &ParticleSystemAsset{
		path: path,
	}
}
