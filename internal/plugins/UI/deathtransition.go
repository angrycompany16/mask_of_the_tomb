package ui

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	deathTransitionColor = []float64{0, 0, 0}
)

type DeathTransition struct {
	image            *ebiten.Image
	transitionShader *ebiten.Shader
}

func (d *DeathTransition) Draw(t float64, enter bool) {
	op := ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{"T": t}
	rendering.ScreenLayers.Overlay.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, d.transitionShader, &op)
}

func NewDeathTransition() OverlayContent {
	return &DeathTransition{
		transitionShader: errs.Must(assettypes.GetShaderAsset("deathTransitionShader")),
		image:            ebiten.NewImage(rendering.GAME_WIDTH, rendering.GAME_HEIGHT),
	}
}
