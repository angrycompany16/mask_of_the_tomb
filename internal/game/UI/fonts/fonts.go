package fonts

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"mask_of_the_tomb/internal/game/core/assetloader/fontasset"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	FontRegistry Fonts
)

type Fonts struct {
	M map[string]*text.GoTextFaceSource
}

func (f *Fonts) Load() {
	mainFontAsset := fontasset.New(assets.JSE_AmigaAMOS_ttf)
	assetloader.AddAsset(mainFontAsset)
	f.M["JSE_AmigaAMOS"] = &mainFontAsset.Font
}
