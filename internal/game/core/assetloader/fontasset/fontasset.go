package fontasset

import (
	"bytes"
	"mask_of_the_tomb/internal/errs"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FontAsset struct {
	src  []byte
	Font *text.GoTextFaceSource
}

func (a *FontAsset) Load() {
	a.Font = errs.Must(text.NewGoTextFaceSource(bytes.NewReader(a.src)))
}

func New(src []byte) *FontAsset {
	return &FontAsset{src: src}
}
