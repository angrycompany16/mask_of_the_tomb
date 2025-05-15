package particles

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"

	"github.com/hajimehoshi/ebiten/v2"
)

type particleSystemAsset struct {
	layer          *ebiten.Image
	path           string
	ParticleSystem ParticleSystem
}

func (a *particleSystemAsset) Load() {
	a.ParticleSystem = *errs.Must(FromFile(a.path, a.layer))
}

func NewParticleSystemAsset(path string, layer *ebiten.Image) *ParticleSystem {
	asset, exists := assetloader.Exists(path)
	if exists {
		return &asset.(*particleSystemAsset).ParticleSystem
	}

	particleSystemAsset := particleSystemAsset{layer: layer, path: path}

	assetloader.Load(path, &particleSystemAsset)
	return &particleSystemAsset.ParticleSystem
}
