package shaderactor

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Find a way to render the tile layers to texture instead
// of to screen. Then process them with shaders.
type Shader struct {
	*graphic.Graphic
	shaderRef *assetloader.AssetRef[ebiten.Shader]
	// Everything changes now. srcImage becomes a texture obtained from the
	// renderer struct. Stuff like player, ldtk tile layer, etc. will render
	// to this texture
	srcImage  *ebiten.Image // The source image for the shader
	Image     *ebiten.Image // The resulting image
	srcPath   string        `debug:"auto"`
	layer     string        `debug:"auto"`
	drawOrder int           `debug:"auto"`
	op        *ebiten.DrawRectShaderOptions
}

func (s *Shader) OnTreeAdd(node *engine.Node, cmd *commands.Commands) {
	s.Graphic.OnTreeAdd(node, cmd)
	s.shaderRef = assetloader.StageAsset[ebiten.Shader](
		cmd.AssetLoader,
		s.srcPath,
		assettypes.NewShaderAsset(s.srcPath),
	)
}

func (s *Shader) Init(cmd *commands.Commands) {
	s.Graphic.Init(cmd)
}

func (s *Shader) Update(cmd *commands.Commands) {
	s.Graphic.Update(cmd)
	if s.shaderRef.Status() != assetloader.LOADED {
		fmt.Println("Shader asset not loaded. Check error logs.")
		return
	}

	// shaderOp.Uniforms = map[string]any{
	// 	"Time":       resources.Time / 5,
	// 	"Amplitude":  1.0,
	// 	"Frequency":  0.025,
	// 	"Strength":   0.7,
	// 	"Threshold":  0.4,
	// 	"Color":      [4]float64{37.0 / 255, 49.0 / 255, 94.0 / 255, 1.0},
	// 	"Center":     [2]float64{0.5, 0.5},
	// 	"Resolution": [2]float64{rendering.GAME_WIDTH, rendering.GAME_HEIGHT},
	// 	"PlayerPos":  [2]float64{ctx.PlayerX, ctx.PlayerY},
	// }
	gw, gh := cmd.Renderer.GetGameSize()
	s.Image.DrawRectShader(int(gw), int(gh), s.shaderRef.Value(), s.op)

	gPosX, gPosY := s.Transform2D.GetPos(false)
	camX, camY := s.GetCamera().WorldToCam(gPosX, gPosY, true)

	cmd.Renderer.Request(opgen.Pos(
		s.Image,
		camX, camY,
	), s.Image, s.layer, s.drawOrder)
}

func (s *Shader) SetShaderOp(op *ebiten.DrawRectShaderOptions) {
	s.op = op
}

func (s *Shader) SetShaderOpUniform(key string, value any) {
	s.op.Uniforms[key] = value
}

func NewShader(graphic *graphic.Graphic, srcPath string, srcImage *ebiten.Image, layer string, drawOrder int) *Shader {

	return &Shader{
		Graphic:   graphic,
		Image:     ebiten.NewImage(srcImage.Bounds().Dx(), srcImage.Bounds().Dy()),
		srcPath:   srcPath,
		srcImage:  srcImage,
		layer:     layer,
		drawOrder: drawOrder,
		op:        &ebiten.DrawRectShaderOptions{Blend: ebiten.BlendSourceOver},
	}
}
