package ldtkworld

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
)

type LDTKLevel struct {
	transform2D.Transform2D
	worldSrcPath string `debug:"auto"`
	levelName    string `debug:"auto"`
	LDTKData     *assetloader.AssetRef[assettypes.LDTKData]
}

func (l *LDTKLevel) OnTreeAdd(node *engine.Node, servers *engine.Servers) {
	l.Transform2D.OnTreeAdd(node, servers)
	l.LDTKData = assetloader.StageAsset[assettypes.LDTKData](
		servers.AssetLoader(),
		l.worldSrcPath,
		assettypes.NewLDTKAsset(l.worldSrcPath),
	)
}

func (l *LDTKLevel) Init() {
	// Needs to spawn children somehow
}

func (l *LDTKLevel) Update(servers *engine.Servers) {
	fmt.Println(l.LDTKData.Value())
}

func (l *LDTKLevel) DrawInspector(ctx *debugui.Context) {
	l.Transform2D.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, l)
}

func NewLDTKLevel(transform2d transform2D.Transform2D, levelName, worldSrcPath string) *LDTKLevel {
	return &LDTKLevel{
		Transform2D:  transform2d,
		levelName:    levelName,
		worldSrcPath: worldSrcPath,
	}
}
