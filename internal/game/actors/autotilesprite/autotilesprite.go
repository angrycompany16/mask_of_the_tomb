package autotilesprite

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/autotile"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type AutoTileSprite struct {
	*graphic.Graphic
	rect       *maths.Rect
	sprite     *ebiten.Image
	tilemapSrc string
	tilemap    *assetloader.AssetRef[ebiten.Image]
}

func (a *AutoTileSprite) OnTreeAdd(node *engine.Node, cmd *engine.Commands) {
	a.Graphic.OnTreeAdd(node, cmd)
	a.tilemap = assetloader.StageAsset[ebiten.Image](
		cmd.AssetLoader(),
		a.tilemapSrc,
		assettypes.NewImageAsset(a.tilemapSrc),
	)
}

func (a *AutoTileSprite) Init(cmd *engine.Commands) {
	a.Graphic.Init(cmd)
	a.createSprite(a.tilemap.Value())
}

func (a *AutoTileSprite) Update(cmd *engine.Commands) {
	a.Graphic.Update(cmd)

	gPosX, gPosY := a.GetPos(false)
	camX, camY := a.GetCamera().WorldToCam(gPosX, gPosY, true)
	cmd.Renderer().Request(opgen.Pos(a.sprite, camX, camY, 0, 0), a.sprite, "Playerspace", 20)
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
			List: make([]*maths.Rect, 0),
			Kind: autotile.WALL,
		},
	)
}

func defaultAutoTileSprite(graphic *graphic.Graphic) *AutoTileSprite {
	return &AutoTileSprite{
		Graphic:    graphic,
		rect:       maths.NewRect(0, 0, 8, 8),
		tilemapSrc: "sprites/environment/slambox_tilemap.png",
	}
}

func NewAutoTileSprite(graphic *graphic.Graphic, options ...utils.Option[AutoTileSprite]) *AutoTileSprite {
	newAutoTileSprite := defaultAutoTileSprite(graphic)

	for _, option := range options {
		option(newAutoTileSprite)
	}

	return newAutoTileSprite
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
