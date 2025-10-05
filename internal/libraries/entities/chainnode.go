package entities

import (
	"image/color"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ChainNode struct {
	Rect            *maths.Rect
	NodeConnections []string
	ID              string
}

func (c *ChainNode) Update() {
}

func (c *ChainNode) Draw(ctx rendering.Ctx) {
	cx, cy := c.Rect.Center()
	vector.DrawFilledCircle(ctx.Dst, float32(cx), float32(cy), 3.0, color.RGBA{255, 0, 0, 255}, false)
}

func NewChainNode(
	entity *ebitenLDTK.Entity,
) *ChainNode {
	newChainNode := ChainNode{}
	newChainNode.ID = entity.Iid
	newChainNode.Rect = maths.NewRect(
		entity.Px[0], entity.Px[1], entity.Width, entity.Height,
	)

	newChainNode.NodeConnections = make([]string, 0)
	nextNodeField := errs.Must(entity.GetFieldByName("NodeConnector"))
	for _, entityRef := range nextNodeField.EntityRefArray {
		newChainNode.NodeConnections = append(newChainNode.NodeConnections, entityRef.EntityIid)
	}

	return &newChainNode
}
