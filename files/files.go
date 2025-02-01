package files

import (
	"bytes"
	"fmt"
	"log"
	"mask_of_the_tomb/ebitenLDTK"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func LazyImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func LazyLDTK(path string) *ebitenLDTK.World {
	ldtkWorld, err := ebitenLDTK.LoadWorld(path)
	if err != nil {
		log.Fatal(err)
	}
	return &ldtkWorld
}

func LazyFont(src []byte) *text.GoTextFaceSource {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(src))
	if err != nil {
		fmt.Println("Font death")
		log.Fatal(err)
	}
	return s
}
