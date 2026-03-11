package main

import (
	"mask_of_the_tomb/internal/node_v2"
	"mask_of_the_tomb/internal/node_v2/ebitenrender"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type object struct {
	name string
	x, y float64
}

type Game struct {
	debugui  debugui.DebugUI
	nodeTree *node_v2.NodeTree[object]
}

func (g *Game) Update() error {
	if _, err := g.debugui.Update(
		ebitenrender.MakeRenderFunc[object]("test UI", 320, 180, g.nodeTree, func(ctx *debugui.Context, nodeVal *object) {}),
	); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	nodeTree, _ := node_v2.NewNodeTree[object](
		object{
			x: 0,
			y: 0,
		},
	)
	root := nodeTree.GetRoot()
	root.AddChild(object{x: 1, y: 1}, "child")
	child2 := root.AddChild(object{x: 1, y: 1}, "child2")
	child2.AddChild(object{x: 2, y: 400}, "grandchild")

	game := &Game{
		nodeTree: nodeTree.DeepCopy(func(o object) object {
			return o
		}),
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
