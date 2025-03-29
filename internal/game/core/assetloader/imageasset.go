package assetloader

import (
	"mask_of_the_tomb/internal/errs"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAsset struct {
	path  string
	Image ebiten.Image
}

func (a *ImageAsset) load() {
	a.Image = *errs.MustNewImageFromFile(a.path)
}

func NewImageAsset(path string) *ImageAsset {
	return &ImageAsset{
		path: path,
	}
}
