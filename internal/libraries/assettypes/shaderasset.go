package assettypes

import (
	"mask_of_the_tomb/internal/core/assetloader"

	"github.com/hajimehoshi/ebiten/v2"
)

type ShaderAsset struct {
	src    []byte
	Shader *ebiten.Shader
}

func (a *ShaderAsset) Load() error {
	shader, err := ebiten.NewShader(a.src)
	a.Shader = shader
	return err
}

func GetShaderAsset(name string) (*ebiten.Shader, error) {
	shaderAsset, err := assetloader.GetAsset(name)
	return shaderAsset.(*ShaderAsset).Shader, err
}

func MakeShaderAsset(src []byte) *ShaderAsset {
	return &ShaderAsset{
		src: src,
	}
}
