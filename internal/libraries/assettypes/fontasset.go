package assettypes

import (
	"bytes"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type fontAsset struct {
	src  []byte
	Font text.GoTextFaceSource
}

func (a *fontAsset) Load() {
	a.Font = *errs.Must(text.NewGoTextFaceSource(bytes.NewReader(a.src)))
}

func NewFontAsset(src []byte) *text.GoTextFaceSource {
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
