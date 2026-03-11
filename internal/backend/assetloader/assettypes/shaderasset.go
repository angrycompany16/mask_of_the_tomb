package assettypes

import (
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

type ShaderAsset struct {
	srcPath string
}

func (a *ShaderAsset) Load(fs fs.FS) (any, error) {
	f, err := fs.Open(a.srcPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	shaderBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	shader, err := ebiten.NewShader(shaderBytes)
	if err != nil {
		return nil, err
	}

	return shader, nil
}

func NewShaderAsset(srcPath string) *ShaderAsset {
	return &ShaderAsset{
		srcPath: srcPath,
	}
}
