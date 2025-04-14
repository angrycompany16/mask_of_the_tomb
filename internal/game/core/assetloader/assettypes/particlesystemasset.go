package assettypes

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"mask_of_the_tomb/internal/game/physics/particles"

	"github.com/hajimehoshi/ebiten/v2"
)

type ParticleSystemAsset struct {
	layer          *ebiten.Image
	path           string
	ParticleSystem particles.ParticleSystem
}

func (a *ParticleSystemAsset) Load() {
	a.ParticleSystem = *errs.Must(particles.FromFile(a.path, a.layer))
}

func NewParticleSystemAsset(path string, layer *ebiten.Image) *particles.ParticleSystem {
	particleSystemAsset := ParticleSystemAsset{layer: layer, path: path}

	assetloader.AddAsset(&particleSystemAsset)
	return &particleSystemAsset.ParticleSystem
}
