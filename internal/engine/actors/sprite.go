package actors

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Transform2D
	SrcPath string
	image   *ebiten.Image
}

func (s *Sprite) Init() {

}

func (s *Sprite) Update() {

}
