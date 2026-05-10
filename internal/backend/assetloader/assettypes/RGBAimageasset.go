package assettypes

import (
	"image"
	"image/draw"
	_ "image/png"
	"io/fs"
)

type RGBAImageAsset struct {
	srcPath string
}

func (a *RGBAImageAsset) Load(fs fs.FS) (any, error) {
	f, err := fs.Open(a.srcPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	RGBAimg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))

	draw.Draw(RGBAimg, RGBAimg.Bounds(), img, img.Bounds().Min, draw.Src)

	return RGBAimg, nil
}

func NewRGBAImageAsset(srcPath string) *RGBAImageAsset {
	return &RGBAImageAsset{
		srcPath: srcPath,
	}
}
