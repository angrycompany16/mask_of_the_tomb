package assettypes

// don't think this is used
// type fontAsset struct {
// 	src []byte
// 	// Maybe turn this into a pointer would help?
// 	Font text.GoTextFaceSource
// }

// func (a *fontAsset) Load() error {
// 	font, err := text.NewGoTextFaceSource(bytes.NewReader(a.src))
// 	a.Font = *font
// 	return err
// }

// // Literally never used
// func NewFontAsset(src []byte) *text.GoTextFaceSource {
// 	// TODO: Do NOT do this
// 	asset, exists := assetloader.Exists(string(src))
// 	if exists {
// 		return &asset.(*fontAsset).Font
// 	}

// 	fontAsset := fontAsset{
// 		src: src,
// 	}

// 	assetloader.Add(string(src), &fontAsset)

// 	return &fontAsset.Font
// }
