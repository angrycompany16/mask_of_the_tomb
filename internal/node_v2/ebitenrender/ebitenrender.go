package ebitenrender

import (
	"image"
	"mask_of_the_tomb/internal/node_v2"

	"github.com/ebitengine/debugui"
)

func MakeRenderFunc[T any](
	title string,
	w, h int,
	nodeTree *node_v2.NodeTree[T],
	nodeFunc func(ctx *debugui.Context, nodeVal T),
) func(ctx *debugui.Context) error {
	root := nodeTree.GetRoot()
	return func(ctx *debugui.Context) error {
		ctx.Window(title, image.Rect(5, 5, w, h), func(layout debugui.ContainerLayout) {
			ctx.TreeNode(root.GetName(), makeTreeNodeRecursion(root, ctx, nodeFunc))
		})
		return nil
	}
}

func makeTreeNodeRecursion[T any](
	node *node_v2.Node[T],
	ctx *debugui.Context,
	nodeFunc func(ctx *debugui.Context, nodeVal T),
) func() {
	return func() {
		children := node.GetChildren()
		ctx.Loop(len(children), func(i int) {
			ctx.TreeNode(children[i].GetName(),
				makeTreeNodeRecursion(children[i], ctx, nodeFunc),
			)
		})
		nodeFunc(ctx, node.GetValue())
	}
}
