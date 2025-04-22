package fonts

import (
	"bytes"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/errs"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gopkg.in/yaml.v3"
)

var (
	_fonts = make(map[string]*text.GoTextFaceSource)
)

type FontYAML struct {
	*text.GoTextFaceSource
}

func (f *FontYAML) UnmarshalYAML(value *yaml.Node) error {
	f.GoTextFaceSource = _fonts[value.Value]
	return nil
}

func LoadPreamble() {
	_fonts["JSE_AmigaAMOS"] = errs.Must(text.NewGoTextFaceSource(bytes.NewReader(assets.JSE_AmigaAMOS_ttf)))
}

func GetFont(name string) *text.GoTextFaceSource {
	return _fonts[name]
}

func Load() {

}
