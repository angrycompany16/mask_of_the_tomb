package align

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/container"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"
)

// Align takes the child items and arranges them, either vertically
// or horitzontally, according to the dimensions specified.
type Align struct {
	*container.Container
	// should contain one number for each item, indicating how much
	// of the relative space each cell should take. I.e. exactly
	// copy what ebiten.debugui does.
	// TODO: Add an option for automatically spacing all elements equally
	spacing []float64
	IsRow   bool
}

func (a *Align) Init(cmd *commands.Commands) {
	a.Container.Init(cmd)
	if len(a.spacing) == 1 && a.spacing[0] == -1 {
		children := a.GetNode().GetChildren()
		a.spacing = make([]float64, len(children))
		for i := range a.spacing {
			a.spacing[i] = 1
		}
	}
}

func (a *Align) Update(cmd *commands.Commands) {
	a.Container.Update(cmd)
	// Loop through children
	// Determine positions and rect sizes based on a.spacing
	var s, parentSize float64
	if a.IsRow {
		parentSize = a.Rect.Width
	} else {
		parentSize = a.Rect.Height
	}

	for _, r := range a.spacing {
		s += r
	}

	children := a.GetNode().GetChildren()
	pos := 0.0
	for i, child := range children {
		container, ok := engine.As[*container.Container](child.GetValue())
		if !ok {
			continue
		}

		if i >= len(a.spacing) {
			fmt.Println("Too many children, quitting.")
			break
		}

		childSize := parentSize * a.spacing[i] / s

		if a.IsRow {
			container.Resize(childSize, a.Rect.Height)
			container.SetPos(pos, container.Rect.Y)
		} else {
			container.Resize(a.Rect.Width, childSize)
			container.SetPos(container.Rect.X, pos)
		}
		pos += childSize
	}
}

func defaultAlign(container *container.Container) *Align {
	return &Align{
		Container: container,
		spacing:   []float64{-1},
		IsRow:     true,
	}
}

func NewAlign(container *container.Container, options ...utils.Option[Align]) *Align {
	align := defaultAlign(container)

	for _, option := range options {
		option(align)
	}

	return align
}

func WithSpacing(spacing []float64) utils.Option[Align] {
	return func(a *Align) {
		a.spacing = spacing
	}
}

func WithIsRow(isRow bool) utils.Option[Align] {
	return func(a *Align) {
		a.IsRow = isRow
	}
}
