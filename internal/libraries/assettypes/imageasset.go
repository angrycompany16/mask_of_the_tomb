package assettypes

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageAsset struct {
	src   []byte
	Image *ebiten.Image
}

func (a *ImageAsset) Load() {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(a.src))
	a.Image = img
	if err != nil {
		panic(err)
	}
}

func MakeImageAsset(src []byte) *ImageAsset {
	return &ImageAsset{
		src: src,
	}
}
