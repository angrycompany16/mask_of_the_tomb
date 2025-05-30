package assettypes

import (
	"bytes"
	"mask_of_the_tomb/internal/core/assetloader"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type fontAsset struct {
	src  []byte
	Font text.GoTextFaceSource
}

func (a *fontAsset) Load() error {
	font, err := text.NewGoTextFaceSource(bytes.NewReader(a.src))
	a.Font = *font
	return err
}

func NewFontAsset(src []byte) *text.GoTextFaceSource {
	// TODO: Do NOT do this
	asset, exists := assetloader.Exists(string(src))
	if exists {
		return &asset.(*fontAsset).Font
	}

	fontAsset := fontAsset{
		src: src,
	}

	assetloader.Load(string(src), &fontAsset)

	return &fontAsset.Font
}
