package ebitenLDTK

import "github.com/hajimehoshi/ebiten/v2"

type Tileset struct {
	Image        *ebiten.Image
	Name         string  `json:"identifier"`
	Uid          int     `json:"uid"`
	RelPath      string  `json:"relPath"`
	PxWid        int     `json:"pxWid"`
	PxHei        int     `json:"pxHei"`
	TileGridSize float64 `json:"tileGridSize"`
	Spacing      int     `json:"spacing"`
	Padding      int     `json:"padding"`
}
