package inspector

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Inspector struct {
	nodeactor.Node
	x, y          float64 `debug:"auto"`
	width, height int     `debug:"auto"`
	visible       bool    `debug:"auto"`
	UI            debugui.DebugUI
	editorImage   *ebiten.Image
}

func (e *Inspector) Update(cmd *engine.Commands) {
	if _, err := e.UI.Update(
		cmd.Scene().MakeDrawFunc(e.width, e.height),
	); err != nil {
		fmt.Println("Error in Editor!")
	}

	if cmd.InputHandler().PollAction("toggleInspector") {
		e.visible = !e.visible
		if e.visible {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
		}
	}

	if e.visible {
		cmd.Renderer().Request(opgen.Pos(e.editorImage, e.x, e.y), e.editorImage, "EditorUI", 0)
		e.editorImage.Clear()
		e.UI.Draw(e.editorImage)
	}
}

func (e *Inspector) DrawInspector(ctx *debugui.Context) {
	e.Node.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, e)
}

func NewInspector(x, y float64, w, h int) *Inspector {
	return &Inspector{
		x:           x,
		y:           y,
		width:       w,
		height:      h,
		editorImage: ebiten.NewImage(w, h),
	}
}
