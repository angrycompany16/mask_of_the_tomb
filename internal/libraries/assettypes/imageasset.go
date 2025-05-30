package assettypes

import (
	"bytes"
	"mask_of_the_tomb/internal/core/assetloader"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageAsset struct {
	src   []byte
	Image *ebiten.Image
}

func (a *ImageAsset) Load() error {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(a.src))
	a.Image = img
	return err
}

func GetImageAsset(name string) (*ebiten.Image, error) {
	imageAsset, err := assetloader.GetAsset(name)
	return imageAsset.(*ImageAsset).Image, err
}

func MakeImageAsset(src []byte) *ImageAsset {
	return &ImageAsset{
		src: src,
	}
}
