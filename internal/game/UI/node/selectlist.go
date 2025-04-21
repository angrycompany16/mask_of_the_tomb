package node

import (
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SelectList struct {
	NodeData    `yaml:",inline"`
	SelectorPos int
}

func (s *SelectList) Update(confirmations map[string]bool) {
	s.UpdateChildren(confirmations)

	if len(s.Children) == 0 {
		return
	}

	inputDir := 0
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		inputDir = 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		inputDir = -1
	}

	s.SelectorPos += inputDir
	s.SelectorPos = maths.Mod(s.SelectorPos, len(s.Children))

	for i, child := range s.Children {
		selectable := child.Node.(selectable)
		if i == s.SelectorPos {
			selectable.SetSelected()
		} else {
			selectable.SetDeselected()
		}
		i++
	}
}

func (s *SelectList) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(s.Width, s.Height, parentWidth, parentHeight)
	s.DrawChildren(offsetX+s.PosX, offsetY+s.PosY, w, h)
}

func (s *SelectList) Reset() {
	s.SelectorPos = 0
	s.ResetChildren()
}
