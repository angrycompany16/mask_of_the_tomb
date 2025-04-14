package node

type Node interface {
	Draw(offsetX, offsetY float64)
	Update()

	AddChild(node Node)
}

type NodeData struct {
	PosX, PosY    float64
	Width, Height float64 // Actually not needed maybe, although may be useful later
	Children      []Node
	Parent        Node
}

func (n *NodeData) DrawChildren(offsetX, offsetY float64) {
	for _, child := range n.Children {
		child.Draw(offsetX, offsetY)
	}
}

func (n *NodeData) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
