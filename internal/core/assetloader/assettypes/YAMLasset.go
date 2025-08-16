package assettypes

import (
	"mask_of_the_tomb/internal/core/assetloader"

	"gopkg.in/yaml.v3"
)

// Note: `out` is expected to be a pointer to a struct
type YamlAsset struct {
	src []byte
	out any
}

func (a *YamlAsset) Load() error {
	err := yaml.Unmarshal(a.src, a.out)
	return err
}

func GetYamlAsset(name string) (any, error) {
	yamlAsset, err := assetloader.GetAsset(name)
	return yamlAsset.(*YamlAsset).out, err
}

func MakeYamlAsset(src []byte, out any) *YamlAsset {
	return &YamlAsset{
		src: src,
		out: out,
	}
}
