package autotilesprite

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/autotile"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"
	"slices"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type AutoTileSprite struct {
	*graphic.Graphic
	rect       *maths.Rect
	otherRects []*maths.Rect
	sprite     *ebiten.Image
	subsprites []*ebiten.Image
	tilemapSrc string
	tilemap    *assetloader.AssetRef[ebiten.Image]
	target     renderer.RenderTarget
}

func (a *AutoTileSprite) OnTreeAdd(node *engine.Node, cmd *commands.Commands) {
	a.Graphic.OnTreeAdd(node, cmd)
	a.tilemap = assetloader.StageAsset[ebiten.Image](
		cmd.AssetLoader,
		a.tilemapSrc,
		assettypes.NewImageAsset(a.tilemapSrc),
	)
}

func (a *AutoTileSprite) Init(cmd *commands.Commands) {
	a.Graphic.Init(cmd)
	a.createSprite(a.tilemap.Value())
}

func (a *AutoTileSprite) Update(cmd *commands.Commands) {
	a.Graphic.Update(cmd)

	gPosX, gPosY := a.GetPos(false)
	camX, camY := a.GetCamera().WorldToCam(gPosX, gPosY, true)
	cmd.Renderer.Request(opgen.Pos(a.sprite, camX, camY, 0, 0), a.sprite, a.target, 20)

	fmt.Println(len(a.subsprites))
	for i, subsprite := range a.subsprites {
		dx := a.otherRects[i].X - a.rect.X
		dy := a.otherRects[i].Y - a.rect.Y

		fmt.Println(dx, dy)
		cmd.Renderer.Request(opgen.Pos(subsprite, camX + dx, camY + dy, 0, 0), subsprite, a.target, 20)
	}
}

func (a *AutoTileSprite) DrawInspector(ctx *debugui.Context) {
	ctx.Header("AutoTileSprite", false, func() {
		utils.RenderFieldsAuto(ctx, a)
	})
	a.Graphic.DrawInspector(ctx)
}

func (a *AutoTileSprite) createSprite(slamboxTilemap *ebiten.Image) {
	autotile.CreateSprite(
		slamboxTilemap,
		a.sprite,
		autotile.GetDefaultTileRectData(0, 0, 8),
		autotile.GetDefaultTileRuleset(),
		8,
		a.rect,
		autotile.WALL,
		autotile.RectList{
			List: a.otherRects,
			Kind: autotile.WALL,
		},
	)


	for i, rect := range a.otherRects {
		neighbours := slices.Concat(a.otherRects[:i], a.otherRects[i+1:])
		neighbours = append(neighbours, a.rect)

		fmt.Println(a.otherRects, a.rect)

		autotile.CreateSprite(
			slamboxTilemap,
			a.subsprites[i],
			autotile.GetDefaultTileRectData(0, 0, 8),
			autotile.GetDefaultTileRuleset(),
			8,
			rect,
			autotile.WALL,
			autotile.RectList{
				List: neighbours,
				Kind: autotile.WALL,
			},
		)
	}
}

func defaultAutoTileSprite(graphic *graphic.Graphic) *AutoTileSprite {
	return &AutoTileSprite{
		Graphic:    graphic,
		rect:       maths.NewRect(0, 0, 8, 8),
		tilemapSrc: "sprites/environment/slambox_tilemap.png",
	}
}

func NewAutoTileSprite(graphic *graphic.Graphic, target renderer.RenderTarget, options ...utils.Option[AutoTileSprite]) *AutoTileSprite {
	newAutoTileSprite := defaultAutoTileSprite(graphic)
	newAutoTileSprite.target = target

	for _, option := range options {
		option(newAutoTileSprite)
	}

	return newAutoTileSprite
}

func WithRects(mainRect *maths.Rect, otherRects []*maths.Rect) utils.Option[AutoTileSprite] {
	return func(a *AutoTileSprite) {
		a.rect = mainRect
		a.otherRects = otherRects

		a.sprite = ebiten.NewImage(int(a.rect.Width), int(a.rect.Height))
		a.subsprites = make([]*ebiten.Image, len(a.otherRects))
		for i, rect := range otherRects {
			a.subsprites[i] = ebiten.NewImage(int(rect.Width), int(rect.Height))
		}
	}
}

func WithSize(width, height float64) utils.Option[AutoTileSprite] {
	return func(a *AutoTileSprite) {
		a.rect.SetSize(width, height)
		a.sprite = ebiten.NewImage(int(width), int(height))
	}
}

func WithTilemap(path string) utils.Option[AutoTileSprite] {
	return func(a *AutoTileSprite) {
		a.tilemapSrc = path
	}
}
