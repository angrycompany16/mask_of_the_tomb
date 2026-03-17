package inspector

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Inspector struct {
	*nodeactor.Node
	x, y          float64 `debug:"auto"`
	width, height int     `debug:"auto"`
	visible       bool    `debug:"auto"`
	UI            debugui.DebugUI
	editorImage   *ebiten.Image
}

func (i *Inspector) Init(cmd *engine.Commands) {
	i.Node.Init(cmd)
	cmd.InputHandler().RegisterAction("toggleInspector", input.KeyJustPressedAction(ebiten.KeyTab))
}

func (i *Inspector) Update(cmd *engine.Commands) {
	if _, err := i.UI.Update(
		cmd.Scene().MakeDrawFunc(i.width, i.height),
	); err != nil {
		fmt.Println("Error in Editor!")
	}

	if cmd.InputHandler().PollAction("toggleInspector") {
		i.visible = !i.visible
		if i.visible {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
		}
	}

	if i.visible {
		i.editorImage.Clear()
		i.UI.Draw(i.editorImage)
		cmd.Renderer().Request(opgen.Pos(i.editorImage, i.x, i.y), i.editorImage, "EditorUI", 0)
	}
}

func (i *Inspector) DrawInspector(ctx *debugui.Context) {
	i.Node.DrawInspector(ctx)
	utils.RenderFieldsAuto(ctx, i)
}

func defaultInspector(node *nodeactor.Node) *Inspector {
	return &Inspector{
		Node:        node,
		x:           0,
		y:           0,
		width:       300,
		height:      500,
		visible:     false,
		editorImage: ebiten.NewImage(300, 500),
	}
}

func NewInspector(node *nodeactor.Node, options ...utils.Option[Inspector]) *Inspector {
	inspector := defaultInspector(node)

	for _, option := range options {
		option(inspector)
	}

	return inspector
}

func WithPos(x, y float64) utils.Option[Inspector] {
	return func(i *Inspector) {
		i.x = x
		i.y = y
	}
}

func WithSize(width, height int) utils.Option[Inspector] {
	return func(i *Inspector) {
		i.width = width
		i.height = height
		i.editorImage = ebiten.NewImage(width, height)
	}
}
