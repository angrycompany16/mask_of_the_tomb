package ldtkworld

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/actors/vectorgraphic"
	ldtktilelayer "mask_of_the_tomb/internal/game/actors/LDTKTileLayer"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

var layerMap = map[string]string{
	"Foreground":       "Foreground",
	"PlayerspaceAlt":   "Playerspace",
	"Playerspace":      "Playerspace",
	"Props":            "Midground",
	"MidgroundSprites": "Midground",
	"BackgroundTiles":  "Background",
}

type LDTKLevel struct {
	transform2D.Transform2D
	worldSrcPath string `debug:"auto"`
	levelName    string `debug:"auto"`
	LDTKData     *assetloader.AssetRef[assettypes.LDTKData]
}

func (l *LDTKLevel) OnTreeAdd(node *engine.Node, servers *engine.Servers) {
	l.Transform2D.OnTreeAdd(node, servers)
	l.LDTKData = assetloader.StageAsset[assettypes.LDTKData](
		servers.AssetLoader(),
		l.worldSrcPath,
		assettypes.NewLDTKAsset(l.worldSrcPath),
	)
}

func (l *LDTKLevel) Init(servers *engine.Servers) {
	// Add some sub-elements
	// servers.Scene().AddChild()
	world := l.LDTKData.Value().World
	tilesetMap := l.LDTKData.Value().Tilesets
	defs := world.Defs
	level := utils.Must(world.GetLevelByName(l.levelName))
	for i := range level.Layers {
		layer := level.Layers[i]

		if layer.Type == ebitenLDTK.LayerTypeEntities {
			continue
		}

		tileset := utils.Must(defs.GetTilesetByUid(layer.TilesetUid))
		tileSize := tileset.TileGridSize
		tilesetImg := tilesetMap[tileset.Name]

		servers.Scene().AddChild(ldtktilelayer.NewLDTKTilemapLayer(
			transform2D.NewTransform2D(
				*nodeactor.NewNode(),
			),
			&layer, tilesetImg, "Playerspace",
			-i,
			int(tileSize),
			int(level.PxWid),
			int(level.PxHei),
		), layer.Name, l.GetNode(), servers)
	}

	// fmt.Println(utils.Must(colors.HexToRGB(level.BgColorHex)))

	servers.Scene().AddChild(vectorgraphic.NewVectorGraphic(
		transform2D.NewTransform2D(
			*nodeactor.NewNode(),
		),
		func(img *ebiten.Image) {
			vector64.FillRect(
				img,
				0, 0, level.PxWid, level.PxHei,
				utils.Must(colors.HexToRGB(level.BgColorHex)), false,
			)
		},
		"Background",
		-len(level.Layers),
		int(level.PxWid),
		int(level.PxHei),
	), "BackgroundColor", l.GetNode(), servers)
}

func (l *LDTKLevel) Update(servers *engine.Servers) {

}

func (l *LDTKLevel) DrawInspector(ctx *debugui.Context) {
	l.Transform2D.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, l)
}

func NewLDTKLevel(transform2d transform2D.Transform2D, levelName, worldSrcPath string) *LDTKLevel {
	return &LDTKLevel{
		Transform2D:  transform2d,
		levelName:    levelName,
		worldSrcPath: worldSrcPath,
	}
}
