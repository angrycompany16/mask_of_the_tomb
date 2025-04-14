package root

import (
	"mask_of_the_tomb/internal/game/UI/node"
	"mask_of_the_tomb/internal/game/core/rendering"
)

type Root struct {
	*node.NodeData
}

func (r *Root) Update() {}

func (r *Root) Draw(offsetX, offsetY float64) {
	r.DrawChildren(offsetX+r.PosX, offsetY+r.PosY)
}

func New() *Root {
	return &Root{
		&node.NodeData{
			Width:    rendering.GameWidth,
			Height:   rendering.GameHeight,
			Children: make([]node.Node, 0),
		},
	}
}
