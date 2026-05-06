package shaderactor

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
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
	Op        *ebiten.DrawRectShaderOptions
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
		// fmt.Println("Shader asset not loaded. Check error logs.")
		return
	}

	gw, gh := cmd.Renderer.GetGameSize()
	s.Image.Clear()
	s.Image.DrawRectShader(int(gw), int(gh), s.shaderRef.Value(), s.Op)

	gPosX, gPosY := s.Transform2D.GetPos(false)
	camX, camY := s.GetCamera().WorldToCam(gPosX, gPosY, true)

	cmd.Renderer.Request(opgen.Pos(
		s.Image,
		camX, camY,
	), s.Image, renderer.RenderTarget{
		renderer.SCREEN,
		s.layer,
	}, s.drawOrder)
}

func (s *Shader) DrawInspector(ctx *debugui.Context) {
	s.Graphic.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, s)
}

func (s *Shader) SetShaderOp(op *ebiten.DrawRectShaderOptions) {
	s.Op = op
}

func (s *Shader) GetSrcImage() *ebiten.Image {
	return s.srcImage
}

func (s *Shader) SetShaderOpUniform(key string, value any) {
	s.Op.Uniforms[key] = value
}

func NewShader(graphic *graphic.Graphic, srcPath string, srcImage *ebiten.Image, layer string, drawOrder int) *Shader {
	return &Shader{
		Graphic:   graphic,
		Image:     ebiten.NewImage(srcImage.Bounds().Dx(), srcImage.Bounds().Dy()),
		srcPath:   srcPath,
		srcImage:  srcImage,
		layer:     layer,
		drawOrder: drawOrder,
		Op:        &ebiten.DrawRectShaderOptions{Blend: ebiten.BlendSourceOver},
	}
}
