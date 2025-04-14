package assettypes

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/assetloader"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAsset struct {
	path  string
	Image ebiten.Image
}

func (a *ImageAsset) Load() {
	a.Image = *errs.MustNewImageFromFile(a.path)
}

func NewImageAsset(path string) *ebiten.Image {
	imageAsset := ImageAsset{
		path: path,
	}

	assetloader.AddAsset(&imageAsset)
	return &imageAsset.Image
}
