package ldtktilelayer

import (
	"image"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

type LDTKTilemapLayer struct {
	*transform2D.Transform2D
	LDTKlayer    *ebitenLDTK.Layer
	tilesetImage *ebiten.Image
	layerImage   *ebiten.Image
	layer        string
	drawOrder    int
	tileSize     float64
}

func (t *LDTKTilemapLayer) Init(servers *engine.Commands) {
	// Pre-render all tiles
	var tiles []ebitenLDTK.Tile
	if t.LDTKlayer.Type == ebitenLDTK.LayerTypeTiles {
		tiles = t.LDTKlayer.GridTiles
	} else if t.LDTKlayer.Type == ebitenLDTK.LayerTypeIntGrid {
		tiles = t.LDTKlayer.AutoLayerTiles
	}

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
		t.layerImage.DrawImage(tileImage, op)
	}
}

func (t *LDTKTilemapLayer) Update(servers *engine.Commands) {
	gPosX, gPosY := t.Transform2D.GetPos(false)
	gScaleX, gScaleY := t.Transform2D.GetScale(false)
	gAngle := t.Transform2D.GetAngle(false)

	t.Transform2D.Update(servers)
	servers.Renderer().Request(opgen.PosScaleRot(
		t.layerImage, gPosX, gPosY, gAngle, gScaleX, gScaleY,
	), t.layerImage, t.layer, t.drawOrder)
}

func NewLDTKTilemapLayer(
	transform2D *transform2D.Transform2D,
	LDTKLayer *ebitenLDTK.Layer,
	tilesetImg *ebiten.Image,
	layer string,
	drawOrder int,
	tileSize int,
	pxWidth, pxHeight int,
) *LDTKTilemapLayer {
	return &LDTKTilemapLayer{
		LDTKlayer:    LDTKLayer,
		Transform2D:  transform2D,
		tilesetImage: tilesetImg,
		layerImage:   ebiten.NewImage(pxWidth, pxHeight),
		layer:        layer,
		drawOrder:    drawOrder,
		tileSize:     float64(tileSize),
	}
}
