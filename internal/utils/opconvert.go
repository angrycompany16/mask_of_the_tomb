package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

func OpConvert(op *ebiten.DrawImageOptions) *colorm.DrawImageOptions {
	return &colorm.DrawImageOptions{GeoM: op.GeoM, Blend: op.Blend, Filter: op.Filter}
}
