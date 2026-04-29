package assettypes

import (
	"io"
	"io/fs"

	"gopkg.in/yaml.v3"
)

// Unsure how the future of this guy will be
// Well he may have some use tbh
// Note: `out` is expected to be a pointer to a struct
type YamlAsset struct {
	srcPath string
	out     any
}

func (a *YamlAsset) Load(fs fs.FS) (any, error) {
	f, err := fs.Open(a.srcPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	yamlBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlBytes, a.out)
	if err != nil {
		return nil, err
	}
	return a.out, err
}

// Make sure to pass in a pointer to the struct that the
// YAML asset should be marshalled into
func NewYamlAsset(srcPath string, out any) *YamlAsset {
	return &YamlAsset{
		srcPath: srcPath,
		out:     out,
	}
}
