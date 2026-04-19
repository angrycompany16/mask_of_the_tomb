package ldtktilelayer

import (
	"image"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type LDTKTilemapLayer struct {
	*graphic.Graphic
	LDTKlayer    *ebitenLDTK.Layer
	tilesetImage *ebiten.Image
	Image        *ebiten.Image
	renderTarget renderer.RenderTarget
	// layer        string  `debug:"auto"`
	drawOrder int     `debug:"auto"`
	tileSize  float64 `debug:"auto"`
	// drawToScreen bool    `debug:"auto"`
}

func (t *LDTKTilemapLayer) Init(cmd *commands.Commands) {
	t.Graphic.Init(cmd)
	// Pre-render all tiles
	var tiles []ebitenLDTK.Tile
	if t.LDTKlayer.Type == ebitenLDTK.LayerTypeTiles {
		tiles = t.LDTKlayer.GridTiles
	} else if t.LDTKlayer.Type == ebitenLDTK.LayerTypeIntGrid {
		tiles = t.LDTKlayer.AutoLayerTiles
	}

	// fmt.Println(t.LDTKlayer.Name, t.tileSize)
	for _, tile := range tiles {
		scaleX, scaleY := 1.0, 1.0
		switch tile.TileOrientation {
		case ebitenLDTK.OrientationFlipX:
			scaleX = -1
		case ebitenLDTK.OrientationFlipY:
			scaleY = -1
		case ebitenLDTK.OrientationFlipXY:
			scaleX, scaleY = -1, -1
		}

		tileImage := t.tilesetImage.SubImage(
			image.Rect(
				int(tile.Src[0]),
				int(tile.Src[1]),
				int(tile.Src[0]+t.tileSize),
				int(tile.Src[1]+t.tileSize),
			),
		).(*ebiten.Image)

		op := opgen.PosScale(tileImage, tile.Px[0], tile.Px[1], scaleX, scaleY, 0.5, 0.5)
		t.Image.DrawImage(tileImage, op)
	}
}

func (t *LDTKTilemapLayer) Update(cmd *commands.Commands) {
	t.Graphic.Update(cmd)
	// if !t.drawToScreen {
	// 	return
	// }
	gPosX, gPosY := t.Transform2D.GetPos(false)
	camX, camY := t.GetCamera().WorldToCam(gPosX, gPosY, true)
	gScaleX, gScaleY := t.Transform2D.GetScale(false)
	gAngle := t.Transform2D.GetAngle(false)

	cmd.Renderer.Request(opgen.PosRotScale(
		t.Image, camX, camY, gAngle, gScaleX, gScaleY, 0.0, 0.0,
	), t.Image, t.renderTarget, t.drawOrder)
}

func (t *LDTKTilemapLayer) DrawInspector(ctx *debugui.Context) {
	utils.RenderFieldsAuto(ctx, t)
	t.Graphic.DrawInspector(ctx)
}

func NewLDTKTilemapLayer(
	graphic *graphic.Graphic,
	LDTKLayer *ebitenLDTK.Layer,
	tilesetImg *ebiten.Image,
	renderTarget renderer.RenderTarget,
	drawOrder int,
	tileSize int,
	pxWidth, pxHeight int,
) *LDTKTilemapLayer {
	return &LDTKTilemapLayer{
		LDTKlayer:    LDTKLayer,
		Graphic:      graphic,
		tilesetImage: tilesetImg,
		Image:        ebiten.NewImage(pxWidth, pxHeight),
		renderTarget: renderTarget,
		drawOrder:    drawOrder,
		tileSize:     float64(tileSize),
	}
}
