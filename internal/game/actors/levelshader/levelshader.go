package levelshader

import (
	"mask_of_the_tomb/internal/backend/shaders"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/actors/shaderactor"

	"github.com/hajimehoshi/ebiten/v2"
)

type LevelShader struct {
	*shaderactor.Shader
	player *player.Player
	test   *ebiten.Image
}

func (l *LevelShader) Init(cmd *commands.Commands) {
	l.Shader.Init(cmd)

	scene, ok := commands.Get[engine.Scene](cmd)
	if !ok {
		panic("Scene is missing from commands")
	}

	playerNode, ok := engine.GetNodeByType[*player.Player](scene)
	if !ok {
		panic("Player not found")
	}
	playerActor, ok := engine.As[*player.Player](playerNode.GetValue())
	l.player = playerActor

	// l.test.Fill(color.RGBA{255, 255, 255, 255})

	// gw, gh := cmd.Renderer.GetGameSize()
	// camX, camY := l.GetCamera().GetPos(false)
	// shakeX, shakeY := l.GetCamera().GetShake()
	// l.Shader.Op = shaders.MakeShaderOp(
	// 	// slices.Concat(
	// 	// arrays.MapSlice(l.turrets, func(turret *entities.Turret) *shaders.Light { return turret.Light }),
	// 	// arrays.MapSlice(l.lanterns, func(lantern *entities.Lantern) *shaders.Light { return lantern.Light }),
	// 	// arrays.MapSlice(l.slamboxEntities, func(slambox *SlamboxEntity) *shaders.Light { return slambox.Light }),
	// 	[]*shaders.Light{l.player.Light},
	// 	// ),
	// 	camX,
	// 	camY,
	// 	shakeX,
	// 	shakeY,
	// 	0.45,
	// 	0.45,
	// 	0.45,
	// 	cmd.GameInfo.GetTime()/5,
	// 	gw,
	// 	gh,
	// 	l.GetSrcImage(),
	// )
}

func (l *LevelShader) Update(cmd *commands.Commands) {
	gw, gh := cmd.Renderer.GetGameSize()
	camX, camY := l.GetCamera().GetPos(false)
	// fmt.Println(camX, camY)
	shakeX, shakeY := l.GetCamera().GetShake()

	l.Op = shaders.MakeShaderOp(
		// slices.Concat(
		// arrays.MapSlice(l.turrets, func(turret *entities.Turret) *shaders.Light { return turret.Light }),
		// arrays.MapSlice(l.lanterns, func(lantern *entities.Lantern) *shaders.Light { return lantern.Light }),
		// arrays.MapSlice(l.slamboxEntities, func(slambox *SlamboxEntity) *shaders.Light { return slambox.Light }),
		[]*shaders.Light{l.player.Light},
		// []*shaders.Light{},
		// ),
		camX,
		camY,
		shakeX,
		shakeY,
		0.45,
		0.45,
		0.45,
		cmd.GameInfo.GetTime()/5,
		gw,
		gh,
		l.GetSrcImage(),
		// l.test,
	)
	l.Shader.Update(cmd)
}

func NewLevelShader(shader *shaderactor.Shader) *LevelShader {
	return &LevelShader{
		Shader: shader,
		test:   ebiten.NewImage(480, 270),
	}
}
