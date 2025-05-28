package node

type RelativePos int

const (
	Centered RelativePos = iota
	TopLeft
)

type Container struct {
	NodeData    `yaml:",inline"`
	RelativePos RelativePos `yaml:"RelativePos"`
}

func (c *Container) Update(confirmations map[string]ConfirmInfo) {
	c.UpdateChildren(confirmations)
}

func (c *Container) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(c.Width, c.Height, parentWidth, parentHeight)
	switch c.RelativePos {
	case Centered:
		c.DrawChildren(offsetX+c.PosX, offsetY+c.PosY, w, h)
		// c.DrawChildren(offsetX+c.PosX+parentWidth/2, offsetY+c.PosY+parentHeight/2, w, h)
	}
}

func (c *Container) Reset(overWriteInfo map[string]OverWriteInfo) {
	c.ResetChildren(overWriteInfo)
}
