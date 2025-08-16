package assettypes

import (
	"bytes"
	"image"
	"mask_of_the_tomb/internal/core/assetloader"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAsset struct {
	src   []byte
	Image image.Image
}

func (a *ImageAsset) Load() error {
	img, _, err := image.Decode(bytes.NewReader(a.src))
	a.Image = img
	return err
}

func GetImageAsset(name string) (*ebiten.Image, error) {
	imageAsset, err := assetloader.GetAsset(name)
	return ebiten.NewImageFromImage((imageAsset.(*ImageAsset).Image)), err
}

func MakeImageAsset(src []byte) *ImageAsset {
	return &ImageAsset{
		src: src,
	}
}
