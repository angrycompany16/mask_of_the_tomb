package uigraphic

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/container"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Add some kind of scaling method so that pixel art can be rendered
type UIGraphic struct {
	*container.Container
	imageRef  *assetloader.AssetRef[ebiten.Image]
	srcPath   string                `debug:"auto"`
	drawOrder int                   `debug:"auto"`
	target    renderer.RenderTarget `debug:"auto"`
}

func (u *UIGraphic) OnTreeAdd(node *engine.Node, cmd *commands.Commands) {
	u.Container.OnTreeAdd(node, cmd)
	u.imageRef = assetloader.StageAsset[ebiten.Image](
		cmd.AssetLoader,
		u.srcPath,
		assettypes.NewImageAsset(u.srcPath),
	)
}

func (u *UIGraphic) Update(cmd *commands.Commands) {
	u.Container.Update(cmd)
	s := u.imageRef.Value().Bounds().Size()
	srcW, srcH := s.X, s.Y

	scaleX := u.Rect.Width / float64(srcW)
	scaleY := u.Rect.Height / float64(srcH)
	op := opgen.PosScale(
		u.imageRef.Value(),
		u.Rect.X, u.Rect.Y,
		scaleX, scaleY,
		0, 0,
	)
	cmd.Renderer.Request(op, u.imageRef.Value(), u.target, u.drawOrder)
}

func (u *UIGraphic) DrawInspector(ctx *debugui.Context) {
	u.Container.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, u)
}

func NewUIGraphic(container *container.Container, srcPath string, drawOrder int, target renderer.RenderTarget) *UIGraphic {
	return &UIGraphic{
		Container: container,
		srcPath:   srcPath,
		drawOrder: drawOrder,
		target:    target,
	}
}
