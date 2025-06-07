package ui

import (
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/assettypes"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	OverlayColor = []float64{0, 0, 0}
)

type ScreenFade struct {
	image            *ebiten.Image
	transitionShader *ebiten.Shader
}

func (d *ScreenFade) Draw(t float64) {
	op := ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{"T": t}
	rendering.ScreenLayers.Overlay.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, d.transitionShader, &op)
}

func NewScreenFade() OverlayContent {
	return &ScreenFade{
		transitionShader: errs.Must(assettypes.GetShaderAsset("transitionShader")),
		image:            ebiten.NewImage(rendering.GAME_WIDTH, rendering.GAME_HEIGHT),
	}
}
