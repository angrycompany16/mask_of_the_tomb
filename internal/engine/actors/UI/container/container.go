package container

import (
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type AutoAlign int

const (
	None AutoAlign = iota
	Center
	Fill
)

type Container struct {
	*nodeactor.Node
	Rect     *maths.Rect
	absRect  *maths.Rect
	OnResize *events.Event
	// Controls whether the Container is aligned automatically in relation
	// to its parent. Note that this only has an effect on relatively
	// positioned containers
	autoAlign AutoAlign
	Relative  bool `debug:"auto"`
}

func (c *Container) Init(cmd *commands.Commands) {
	c.Node.Init(cmd)
	UIControls := cmd.InputHandler.InputSchemes["UIControls"]
	UIControls.AddBinding("UIConfirm", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeySpace)
	})
	UIControls.AddBinding("UIConfirm", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	})
	UIControls.AddBinding("UIClick", func() bool {
		return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	})
	UIControls.AddBinding("UIRight", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyRight)
	})
	UIControls.AddBinding("UILeft", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	})
	UIControls.AddBinding("UIUp", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyUp)
	})
	UIControls.AddBinding("UIDown", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyDown)
	})
}

func (c *Container) Update(cmd *commands.Commands) {
	c.Node.Update(cmd)

	parentNode := c.Node.GetNode().GetParent()
	if parentNode == nil {
		return
	}

	if parentContainer, ok := engine.As[*Container](parentNode.GetValue()); ok {
		if c.Relative {
			switch c.autoAlign {
			case None:
			case Center:
				wDiff := c.Rect.Width - parentContainer.Rect.Width
				hDiff := c.Rect.Height - parentContainer.Rect.Height
				offsetX := wDiff / 2
				offsetY := hDiff / 2
				c.Rect.SetPos(offsetX, offsetY)
			case Fill:
				c.Rect.SetPos(0, 0)
				c.Rect.SetSize(parentContainer.Rect.Size())
			}

			c.absRect.SetPos(
				parentContainer.absRect.X+c.Rect.X, parentContainer.absRect.Y+c.Rect.Y,
			)
		} else {
			c.absRect.SetPos(
				c.Rect.X,
				c.Rect.Y,
			)
		}
	} else {
		c.absRect.SetPos(
			c.Rect.X,
			c.Rect.Y,
		)
	}
}

func (c *Container) DrawInspector(ctx *debugui.Context) {
	c.Node.DrawInspector(ctx)
	ctx.SetGridLayout([]int{0}, []int{0})

	ctx.Header("Rect", false, func() {
		ctx.SetGridLayout([]int{-1, -1, -1}, []int{0, 0, 0})
		ctx.Text("Position")
		ctx.NumberFieldF(&c.Rect.X, 10, 0)
		ctx.NumberFieldF(&c.Rect.Y, 10, 0)

		ctx.Text("Size")
		ctx.NumberFieldF(&c.Rect.Width, 10, 2).On(func() { c.OnResize.WithData("Rect", *c.Rect).Raise() })
		ctx.NumberFieldF(&c.Rect.Height, 10, 2).On(func() { c.OnResize.WithData("Rect", *c.Rect).Raise() })
	})
	utils.RenderFieldsAuto(ctx, c)
}

func (c *Container) SetPos(x, y float64) {
	c.Rect.SetPos(x, y)
}

func (c *Container) Resize(w, h float64) {
	if w == c.Rect.Width && h == c.Rect.Height {
		return
	}
	c.Rect.SetSize(w, h)
	c.absRect.SetSize(w, h)
	c.OnResize.WithData("Rect", *c.Rect).Raise()
}

func (c *Container) GetAbsPos() (float64, float64) {
	return c.absRect.X, c.absRect.Y
}

func (c *Container) GetAbsRect() maths.Rect {
	return *c.absRect
}

func defaultContainer(node *nodeactor.Node) *Container {
	return &Container{
		Node:     node,
		Rect:     maths.NewRect(0, 0, 460, 100),
		absRect:  maths.NewRect(0, 0, 460, 100),
		OnResize: events.NewEvent(),
		Relative: true,
	}
}

func NewContainer(node *nodeactor.Node, options ...utils.Option[Container]) *Container {
	container := defaultContainer(node)

	for _, option := range options {
		option(container)
	}

	return container
}

func WithRect(rect *maths.Rect) utils.Option[Container] {
	return func(c *Container) {
		c.Rect = rect
		c.absRect = rect
	}
}

func WithRelative(relative bool) utils.Option[Container] {
	return func(c *Container) {
		c.Relative = relative
	}
}

func WithAutoAlign(autoAlign AutoAlign) utils.Option[Container] {
	return func(c *Container) {
		c.autoAlign = autoAlign
	}
}
