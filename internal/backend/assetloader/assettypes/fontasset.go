package assettypes

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FontAsset struct {
	srcPath string
}

func (f *FontAsset) Load(fs fs.FS) (any, error) {
	file, err := fs.Open(f.srcPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	font, err := text.NewGoTextFaceSource(file)
	if err != nil {
		return nil, err
	}

	return font, nil
}

func NewFontAsset(srcPath string) *FontAsset {
	return &FontAsset{
		srcPath: srcPath,
	}
}
