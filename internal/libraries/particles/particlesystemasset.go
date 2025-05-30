package particles

import (
	"mask_of_the_tomb/internal/core/assetloader"

	"github.com/hajimehoshi/ebiten/v2"
)

type particleSystemAsset struct {
	layer          *ebiten.Image
	path           string
	ParticleSystem *ParticleSystem
}

func (a *particleSystemAsset) Load() error {
	particleSys, err := FromFile(a.path, a.layer)
	a.ParticleSystem = particleSys
	return err
}

func GetParticleSystemAsset(name string) (*ParticleSystem, error) {
	_particleSystemAsset, err := assetloader.GetAsset(name)
	return _particleSystemAsset.(*particleSystemAsset).ParticleSystem, err
}

func NewParticleSystemAsset(path string, layer *ebiten.Image) *particleSystemAsset {
	return &particleSystemAsset{
		layer: layer,
		path:  path,
	}
}
