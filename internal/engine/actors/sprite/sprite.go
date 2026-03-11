package sprite

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/utils"
	"math"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Option func(*Sprite)

type Sprite struct {
	transform2D.Transform2D
	layer      string  `debug:"auto"`
	drawOrder  int     `debug:"auto"`
	scaling    float64 `debug:"auto"`
	srcPath    string  `debug:"auto"`
	imageAsset *assetloader.AssetRef[ebiten.Image]
}

func (s *Sprite) OnTreeAdd(node *engine.Node, servers *engine.Servers) {
	s.Transform2D.OnTreeAdd(node, servers)
	s.imageAsset = assetloader.StageAsset[ebiten.Image](
		servers.AssetLoader(),
		"sprite", // this is not a very smart name
		assettypes.NewImageAsset(s.srcPath),
	)
}

func (s *Sprite) Update(servers *engine.Servers) {
	s.Transform2D.Update(servers)
	if s.imageAsset.Status() != assetloader.LOADED {
		fmt.Println("Error: Sprite image asset not loaded")
		// This should in theory never happen. But humans make mistakes...
		return
	}
	gPosX, gPosY := s.Transform2D.GetPos(false)
	gAngle := s.Transform2D.GetAngle(false)
	gScaleX, gScaleY := s.Transform2D.GetScale(false)

	// Change this so that stuff is centered tbh
	servers.Renderer().Request(opgen.PosScaleRot(
		s.imageAsset.Value(),
		gPosX, gPosY,
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
	s.Transform2D.DrawInspector(ctx)
}

// We can remove layer and src as required args
func NewSprite(transform2D transform2D.Transform2D, layer string, srcPath string, options ...Option) *Sprite {
	s := defaultSprite(transform2D)
	s.layer = layer
	s.srcPath = srcPath

	for _, option := range options {
		option(s)
	}

	return s
}

func defaultSprite(transform2D transform2D.Transform2D) *Sprite {
	return &Sprite{
		layer:       "Playerspace",
		drawOrder:   0,
		scaling:     1.0,
		Transform2D: transform2D,
	}
}

func WithDrawOrder(drawOrder int) Option {
	return func(s *Sprite) {
		s.drawOrder = drawOrder
	}
}

func WithScaling(scaling float64) Option {
	return func(s *Sprite) {
		s.scaling = scaling
	}
}
