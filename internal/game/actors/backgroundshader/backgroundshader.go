package backgroundshader

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/actors/shaderactor"
)

type BackgroundShader struct {
	*shaderactor.Shader
	player *player.Player
}

func (b *BackgroundShader) Init(cmd *commands.Commands) {
	b.Shader.Init(cmd)

	gw, gh := cmd.Renderer.GetGameSize()
	b.Shader.Op.Uniforms = map[string]any{
		"Amplitude":  1.0,
		"Frequency":  0.025,
		"Strength":   0.7,
		"Threshold":  0.4,
		"Color":      [4]float64{37.0 / 255, 49.0 / 255, 94.0 / 255, 1.0},
		"Center":     [2]float64{0.5, 0.5},
		"Resolution": [2]float64{gw, gh},
		"Time":       cmd.GameInfo.GetTime() / 5,
	}

	scene, ok := commands.Get[engine.Scene](cmd)
	if !ok {
		panic("Scene is missing from commands")
	}

	playerNode, ok := engine.GetNodeByType[*player.Player](scene)
	if !ok {
		panic("Player not found")
	}
	playerActor, ok := engine.As[*player.Player](playerNode.GetValue())
	b.player = playerActor
}

func (b *BackgroundShader) Update(cmd *commands.Commands) {
	b.Shader.Update(cmd)

	playerX, playerY := b.player.GetPos()

	b.Op.Uniforms["Time"] = cmd.GameInfo.GetTime() / 5
	b.Op.Uniforms["PlayerPos"] = [2]float64{playerX, playerY}
}

func NewBackgroundShader(shader *shaderactor.Shader) *BackgroundShader {
	return &BackgroundShader{
		Shader: shader,
	}
}
