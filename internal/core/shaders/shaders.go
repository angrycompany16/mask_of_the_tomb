package shaders

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Light struct {
	X           float64
	Y           float64
	InnerRadius float64
	OuterRadius float64
	ZOffset     float64
	Intensity   float64
	R           float64
	G           float64
	B           float64
}

type lightParameters struct {
	PositionsX  [100]float64
	PositionsY  [100]float64
	InnerRadii  [100]float64
	OuterRadii  [100]float64
	ZOffsets    [100]float64
	Intensities [100]float64
	ColorsR     [100]float64
	ColorsG     [100]float64
	ColorsB     [100]float64
}

// camX, camY should be camera.GetStablePos(). Src should be the image to sample
// time should be resources.time / 5
func MakeShaderOp(
	lights []*Light,
	camX, camY, camShakeX, camShakeY float64,
	ambientR, ambientG, ambientB float64,
	time float64,
	game_width, game_height float64,
	src *ebiten.Image,
) ebiten.DrawRectShaderOptions {
	shaderOp := ebiten.DrawRectShaderOptions{}

	shaderOp.Images = [4]*ebiten.Image{
		// NEVER touch the first texture argument. EVER.
		nil,
		src.SubImage(image.Rect(int(camX), int(camY), int(camX+game_width), int(camY+game_height))).(*ebiten.Image),
		nil,
		nil,
	}

	lightParameters := lightParameters{}

	for i, light := range lights {
		lightParameters.PositionsX[i] = light.X - camX
		lightParameters.PositionsY[i] = light.Y - camY
		lightParameters.InnerRadii[i] = light.InnerRadius
		lightParameters.OuterRadii[i] = light.OuterRadius
		lightParameters.ZOffsets[i] = light.ZOffset
		lightParameters.Intensities[i] = light.Intensity
		lightParameters.ColorsR[i] = light.R
		lightParameters.ColorsG[i] = light.G
		lightParameters.ColorsB[i] = light.B

		if i == 99 {
			fmt.Println("WARNING: More than 100 lights, stopping early.")
			break
		}
	}

	shaderOp.Uniforms = map[string]any{
		"CamShake":     [2]float64{camShakeX, camShakeY},
		"Time":         time,
		"PositionsX":   lightParameters.PositionsX,
		"PositionsY":   lightParameters.PositionsY,
		"InnerRadii":   lightParameters.InnerRadii,
		"OuterRadii":   lightParameters.OuterRadii,
		"ZOffsets":     lightParameters.ZOffsets,
		"Intensities":  lightParameters.Intensities,
		"ColorsR":      lightParameters.ColorsR,
		"ColorsG":      lightParameters.ColorsG,
		"ColorsB":      lightParameters.ColorsB,
		"AmbientLight": [3]float64{ambientR, ambientG, ambientB},
	}

	return shaderOp
}

// camX, camY should be camera.GetStablePos(). Src should be the image to sample
// time should be resources.time / 5
func ChangeSrc(
	shaderOp ebiten.DrawRectShaderOptions,
	camX, camY float64,
	game_width, game_height float64,
	src *ebiten.Image,
) ebiten.DrawRectShaderOptions {
	shaderOp.Images = [4]*ebiten.Image{
		// NEVER touch the first texture argument. EVER.
		nil,
		src.SubImage(image.Rect(int(camX), int(camY), int(camX+game_width), int(camY+game_height))).(*ebiten.Image),
		nil,
		nil,
	}
	return shaderOp
}
