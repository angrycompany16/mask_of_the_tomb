package fonts

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/game/core/assetloader/fontasset"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TODO: Problem: Creating font assets is impossible because copying fonts by value is
// not allowed in ebiten. Thus the font has to somehow be copied as a reference, but
// then there's a problem as we get a nil reference which seems to just automatically
// get dropped?

var (
	FontRegistry = Fonts{M: make(map[string]*text.GoTextFaceSource)}
)

type Fonts struct {
	M map[string]*text.GoTextFaceSource
}

func (f *Fonts) LoadPreamble() {
	mainFontAsset := fontasset.New(assets.JSE_AmigaAMOS_ttf)
	mainFontAsset.Load()
	f.M["JSE_AmigaAMOS"] = mainFontAsset.Font
}

func (f *Fonts) Load() {

}
