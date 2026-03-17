package assetviewer

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"

	"github.com/ebitengine/debugui"
	om "github.com/wk8/go-ordered-map/v2"
)

type AssetViewer struct {
	*nodeactor.Node
	assetpool *om.OrderedMap[string, *assetloader.Asset]
}

func (a *AssetViewer) Init(cmd *engine.Commands) {
	a.Node.Init(cmd)
	a.assetpool = cmd.AssetLoader().GetAssetPool()
}

func (a *AssetViewer) DrawInspector(ctx *debugui.Context) {
	a.Node.DrawInspector(ctx)
	// Make a list of all assets
	ctx.SetGridLayout([]int{-3, -1}, []int{0, 0})

	pair := a.assetpool.Oldest()
	ctx.Loop(a.assetpool.Len(), func(i int) {
		ctx.Text(pair.Key)
		ctx.Text(pair.Value.GetStatusString())
		pair = pair.Next()
	})
}

func NewAssetViewer(node *nodeactor.Node) *AssetViewer {
	return &AssetViewer{
		Node: node,
	}
}
