package sprite

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"
	"math"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	*graphic.Graphic
	layer      string  `debug:"auto"`
	drawOrder  int     `debug:"auto"`
	scaling    float64 `debug:"auto"`
	srcPath    string  `debug:"auto"`
	imageAsset *assetloader.AssetRef[ebiten.Image]
}

func (s *Sprite) OnTreeAdd(node *engine.Node, cmd *engine.Commands) {
	s.Graphic.OnTreeAdd(node, cmd)
	s.imageAsset = assetloader.StageAsset[ebiten.Image](
		cmd.AssetLoader(),
		s.srcPath,
		assettypes.NewImageAsset(s.srcPath),
	)
}

func (s *Sprite) Init(cmd *engine.Commands) {
	s.Graphic.Init(cmd)
}

func (s *Sprite) Update(cmd *engine.Commands) {
	s.Graphic.Update(cmd)
	if s.imageAsset.Status() != assetloader.LOADED {
		fmt.Println("Error: Sprite image asset not loaded")
		// This should in theory never happen. But humans make mistakes...
		return
	}
	gPosX, gPosY := s.Transform2D.GetPos(false)
	gAngle := s.Transform2D.GetAngle(false)
	gScaleX, gScaleY := s.Transform2D.GetScale(false)
	camX, camY := s.GetCamera().WorldToCam(gPosX, gPosY, true)

	// Change this so that stuff is centered tbh
	cmd.Renderer().Request(opgen.PosRotScale(
		s.imageAsset.Value(),
		camX, camY,
		gAngle,
		gScaleX*s.scaling, gScaleY*s.scaling,
		0.5, 0.5,
	), s.imageAsset.Value(), s.layer, s.drawOrder)
}

func (s *Sprite) DrawInspector(ctx *debugui.Context) {
	ctx.SetGridLayout([]int{0}, []int{0})
	ctx.Header("Sprite", false, func() {
		ctx.SetGridLayout([]int{-1, -2}, []int{-1, -5})
		ctx.Text("Image")
		ctx.GridCell(func(bounds image.Rectangle) {
			ctx.DrawOnlyWidget(func(screen *ebiten.Image) {
				scale := ctx.Scale()
				S := s.imageAsset.Value().Bounds().Size()
				trueWidth := float64(bounds.Dx() * scale)
				trueHeight := float64(bounds.Dy() * scale)
				scalingFactorX := trueWidth / float64(S.X)
				scalingFactorY := trueHeight / float64(S.Y)
				scalingFactor := math.Min(float64(scalingFactorX), float64(scalingFactorY))
				screen.DrawImage(s.imageAsset.Value(), opgen.PosScale(
					s.imageAsset.Value(),
					float64(bounds.Min.X*scale),
					float64(bounds.Min.Y*scale),
					scalingFactor,
					scalingFactor,
				))
			})
		})
		utils.RenderFieldsAuto(ctx, s)
	})
	s.Graphic.DrawInspector(ctx)
}

func (s *Sprite) GetLayer() string {
	return s.layer
}

// We can remove layer and src as required args
func NewSprite(graphic *graphic.Graphic, layer string, srcPath string, options ...utils.Option[Sprite]) *Sprite {
	s := defaultSprite(graphic)
	s.layer = layer
	s.srcPath = srcPath

	for _, option := range options {
		option(s)
	}

	return s
}

func defaultSprite(graphic *graphic.Graphic) *Sprite {
	return &Sprite{
		layer:     "Playerspace",
		drawOrder: 0,
		scaling:   1.0,
		Graphic:   graphic,
	}
}

func WithDrawOrder(drawOrder int) utils.Option[Sprite] {
	return func(s *Sprite) {
		s.drawOrder = drawOrder
	}
}

func WithScaling(scaling float64) utils.Option[Sprite] {
	return func(s *Sprite) {
		s.scaling = scaling
	}
}
