package assettypes

import (
	"mask_of_the_tomb/internal/errs"
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

func NewParticleSystemAsset(path string, layer *ebiten.Image) *ParticleSystemAsset {
	return &ParticleSystemAsset{layer: layer, path: path}
}
