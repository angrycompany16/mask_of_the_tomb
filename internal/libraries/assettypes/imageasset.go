package assettypes

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"

	"github.com/hajimehoshi/ebiten/v2"
)

type imageAsset struct {
	path  string
	Image ebiten.Image
}

func (a *imageAsset) Load() {
	a.Image = *errs.MustNewImageFromFile(a.path)
}

func NewImageAsset(path string) *ebiten.Image {
	asset, exists := assetloader.Exists(path)
	if exists {
		return &asset.(*imageAsset).Image
	}

	imageAsset := imageAsset{
		path: path,
	}

	assetloader.Load(path, &imageAsset)
	return &imageAsset.Image
}
