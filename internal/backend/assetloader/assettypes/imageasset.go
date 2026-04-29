package assettypes

import (
	"image"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAsset struct {
	srcPath string
}

func (a *ImageAsset) Load(fs fs.FS) (any, error) {
	f, err := fs.Open(a.srcPath)
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

func NewImageAsset(srcPath string) *ImageAsset {
	return &ImageAsset{
		srcPath: srcPath,
	}
}
