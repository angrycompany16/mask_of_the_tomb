package ui

import (
	"image"
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/assettypes"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	FONT_SIZE = 48.0
)

var (
	topLeftRect     = image.Rect(0, 0, 8, 8)
	topRightRect    = image.Rect(16, 0, 24, 8)
	bottomLeftRect  = image.Rect(0, 16, 8, 24)
	bottomRightRect = image.Rect(16, 16, 24, 24)
	topRect         = image.Rect(8, 0, 16, 8)
	bottomRect      = image.Rect(8, 16, 16, 24)
	leftRect        = image.Rect(0, 8, 8, 16)
	rightRect       = image.Rect(16, 8, 24, 16)
	centerRect      = image.Rect(8, 8, 16, 16)
)

type LevelCard struct {
	text string
	font *text.GoTextFaceSource
	img  *ebiten.Image
}

func (lc *LevelCard) Draw(t float64) {
	h := maths.Lerp((rendering.GAME_HEIGHT+32)*rendering.PIXEL_SCALE, (rendering.GAME_HEIGHT-32)*rendering.PIXEL_SCALE, maths.CubicInOut(t))
	ebitenrenderutil.DrawAt(lc.img, rendering.ScreenLayers.ScreenUI, rendering.GAME_WIDTH*rendering.PIXEL_SCALE/2, h, 0.5, 0.5)
}

// TODO: Rewrite with autotile module
func (lc *LevelCard) ChangeText(newText string) {
	tilemap := errs.Must(assettypes.GetImageAsset("titleCard"))
	// TILES
	topLeft := tilemap.SubImage(topLeftRect).(*ebiten.Image)
	topRight := tilemap.SubImage(topRightRect).(*ebiten.Image)
	bottomLeft := tilemap.SubImage(bottomLeftRect).(*ebiten.Image)
	bottomRight := tilemap.SubImage(bottomRightRect).(*ebiten.Image)
	top := tilemap.SubImage(topRect).(*ebiten.Image)
	bottom := tilemap.SubImage(bottomRect).(*ebiten.Image)
	left := tilemap.SubImage(leftRect).(*ebiten.Image)
	right := tilemap.SubImage(rightRect).(*ebiten.Image)
	center := tilemap.SubImage(centerRect).(*ebiten.Image)

	// ?
	if newText == "" {
		newText = "-UNTITLED-"
	}
	lc.text = newText
	width := (FONT_SIZE*float64(len(newText)) + 60) / rendering.PIXEL_SCALE
	height := (FONT_SIZE + 20) / rendering.PIXEL_SCALE
	widthTiles := math.Ceil(width / 8)
	heightTiles := math.Ceil(height / 8)
	lc.img = ebiten.NewImage(int(widthTiles*8*rendering.PIXEL_SCALE), int(heightTiles*8*rendering.PIXEL_SCALE))

	// DRAW BACKGROUND
	tileimg := ebiten.NewImage(int(widthTiles*8), int(heightTiles*8))
	ebitenrenderutil.DrawAt(topLeft, tileimg, 0, 0)
	ebitenrenderutil.DrawAt(topRight, tileimg, (widthTiles-1)*8, 0)
	ebitenrenderutil.DrawAt(bottomLeft, tileimg, 0, (heightTiles-1)*8)
	ebitenrenderutil.DrawAt(bottomRight, tileimg, (widthTiles-1)*8, (heightTiles-1)*8)
	for i := 1; i < int(widthTiles)-1; i++ {
		ebitenrenderutil.DrawAt(top, tileimg, float64(i*8), 0)
		ebitenrenderutil.DrawAt(bottom, tileimg, float64(i*8), (heightTiles-1)*8)
	}

	for i := 1; i < int(heightTiles)-1; i++ {
		ebitenrenderutil.DrawAt(left, tileimg, 0, float64(i*8))
		ebitenrenderutil.DrawAt(right, tileimg, (widthTiles-1)*8, float64(i*8))
	}

	for i := 1; i < int(widthTiles)-1; i++ {
		for j := 1; j < int(heightTiles)-1; j++ {
			ebitenrenderutil.DrawAt(center, tileimg, float64(i*8), float64(j*8))
		}
	}

	ebitenrenderutil.DrawAtScaled(tileimg, lc.img, 0, 0, rendering.PIXEL_SCALE, rendering.PIXEL_SCALE)

	// DRAW TEXT
	opText := &text.DrawOptions{}
	opText.LayoutOptions.PrimaryAlign = text.AlignCenter
	opText.LayoutOptions.SecondaryAlign = text.AlignCenter
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.SetR(1.0)
	opText.ColorScale.SetG(1.0)
	opText.ColorScale.SetB(1.0)
	opText.ColorScale.SetA(1.0)
	opText.GeoM.Translate(widthTiles*4*rendering.PIXEL_SCALE, heightTiles*4*rendering.PIXEL_SCALE)

	text.Draw(lc.img, lc.text,
		&text.GoTextFace{
			Source: lc.font,
			Size:   FONT_SIZE,
		}, opText)
}

func NewLevelCard() OverlayContent {
	return &LevelCard{
		text: "",
		img:  ebiten.NewImage(1, 1),
		font: assetloader.GetFont("C&C_Red_alert"),
	}
}
