package node

type Node interface {
	Draw(offsetX, offsetY float64, parentWidth, parentHeight float64)
	Update(confirmations map[string]bool)
	Reset()

	AddChild(node NodeContainer)
}

type NodeData struct {
	PosX     float64         `yaml:"PosX"`
	PosY     float64         `yaml:"PosY"`
	Width    float64         `yaml:"Width"`
	Height   float64         `yaml:"Height"`
	Parent   NodeContainer   `yaml:"Parent"`
	Children []NodeContainer `yaml:"Children"`
}

func (n *NodeData) UpdateChildren(confirmations map[string]bool) {
	for _, child := range n.Children {
		child.Update(confirmations)
	}
}

func (n *NodeData) DrawChildren(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	for _, child := range n.Children {
		child.Draw(offsetX, offsetY, parentWidth, parentHeight)
	}
}

func (n *NodeData) ResetChildren() {
	for _, child := range n.Children {
		child.Reset()
	}
}

func (n *NodeData) AddChild(node NodeContainer) {
	n.Children = append(n.Children, node)
}

func inheritSize(width, height, parentWidth, parentHeight float64) (outWidth float64, outHeight float64) {
	if width == 0 {
		outWidth = parentWidth
	} else {
		outWidth = width
	}
	if height == 0 {
		outHeight = parentHeight
	} else {
		outHeight = height
	}
	return
}
