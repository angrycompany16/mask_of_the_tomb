package assettypes

import (
	"image"
	"io/fs"
	"mask_of_the_tomb/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAsset struct {
	srcPath string
	Image   image.Image
}

func (a *ImageAsset) Load(fs fs.FS) (any, error) {
	f, err := assets.FS.Open(a.srcPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func MakeImageAsset(srcPath string) *ImageAsset {
	return &ImageAsset{
		srcPath: srcPath,
	}
}
