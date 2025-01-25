package files

import (
	"log"
	"mask_of_the_tomb/ebitenLDTK"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LazyImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func LazyLDTK(path string) *ebitenLDTK.LDTKWorld {
	ldtkWorld, err := ebitenLDTK.LoadLDTK(path)
	if err != nil {
		log.Fatal(err)
	}
	return &ldtkWorld
}
