package ui

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	LevelTransitionColor = []float64{0, 0, 0}
)

type LevelTransition struct {
	image       *ebiten.Image
	enterShader *ebiten.Shader
	exitShader  *ebiten.Shader
}

func (l *LevelTransition) Draw(t float64, enter bool) {
	var h float64
	if enter {
		h = maths.Lerp(0, 1, maths.QuadIn(t))
	} else {
		h = maths.Lerp(0, 1, maths.QuadOut(t))
	}

	op := ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"T":          h,
		"Resolution": [2]float64{rendering.GAME_WIDTH, rendering.GAME_HEIGHT},
	}

	// Don't draw directly to rendering.ScreenLayers!
	if enter {
		rendering.ScreenLayers.Overlay.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.enterShader, &op)
	} else {
		rendering.ScreenLayers.Overlay.DrawRectShader(rendering.GAME_WIDTH, rendering.GAME_HEIGHT, l.exitShader, &op)
	}
}

func NewLevelTransition() OverlayContent {
	return &LevelTransition{
		enterShader: errs.Must(assettypes.GetShaderAsset("levelTransitionEnterShader")),
		exitShader:  errs.Must(assettypes.GetShaderAsset("levelTransitionExitShader")),
		image:       ebiten.NewImage(rendering.GAME_WIDTH, rendering.GAME_HEIGHT),
	}
}
