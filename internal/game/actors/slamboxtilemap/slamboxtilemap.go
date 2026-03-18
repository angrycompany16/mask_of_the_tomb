package slamboxtilemap

import (
	"fmt"
	"image/color"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type SlamboxTilemap struct {
	*graphic.Graphic
	tiles       [][]int // TODO: Change to tilemap interface
	tileSize    int     `debug:"auto"`
	gizmosImage *ebiten.Image
}

func (st *SlamboxTilemap) Init(cmd *engine.Commands) {
	st.Graphic.Init(cmd)
	if cmd.SlamboxEnv().TileSize != float64(st.tileSize) {
		fmt.Println("Warning: Backend tilemap has a different tilesize than this actor!")
	}

	cmd.SlamboxEnv().SetTiles(st.tiles)

	for i := range st.tiles {
		for j := range st.tiles[i] {
			if st.tiles[i][j] == 0 {
				continue
			}
			x := j * st.tileSize
			y := i * st.tileSize
			vector64.FillRect(st.gizmosImage, float64(x), float64(y), float64(st.tileSize), float64(st.tileSize), color.RGBA{232, 128, 9, 173}, false)
		}
	}
}

func (st *SlamboxTilemap) DrawInspector(ctx *debugui.Context) {
	st.Graphic.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, st)
}

// TODO: Make it so that gizmos are only drawn when selected in the
// inspector. But how? Haven't got the slightest idea.
func (st *SlamboxTilemap) DrawGizmo(cmd *engine.Commands) {
	st.Graphic.DrawGizmo(cmd)
	gPosX, gPosY := st.GetPos(false)
	camX, camY := st.GetCamera().WorldToCam(gPosX, gPosY, false)
	cmd.Renderer().Request(opgen.Pos(st.gizmosImage, camX, camY, 0.5, 0.5), st.gizmosImage, "Overlay", 1)
}

func NewSlamboxTilemap(graphic *graphic.Graphic, tiles [][]int, tileSize int) *SlamboxTilemap {
	return &SlamboxTilemap{
		Graphic:     graphic,
		tiles:       tiles,
		tileSize:    tileSize,
		gizmosImage: ebiten.NewImage(len(tiles[0])*tileSize, len(tiles)*tileSize),
	}
}
