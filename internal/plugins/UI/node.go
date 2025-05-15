package ui

import (
	"errors"
	"slices"

	"gopkg.in/yaml.v3"
)

type Node interface {
	Draw(offsetX, offsetY float64, parentWidth, parentHeight float64)
	Update(confirmations map[string]bool)
	Reset()

	AddChild(node NodeContainer)
}

type NodeContainer struct {
	Node
}

func (n *NodeContainer) UnmarshalYAML(value *yaml.Node) error {
	i := slices.IndexFunc(value.Content, func(value *yaml.Node) bool { return value.Value == "Type" })
	if i == -1 {
		return errors.New("missing type, failed to unmarshal")
	}

	var resultNode Node
	nodeType := value.Content[i+1].Value
	switch nodeType {
	case "textbox":
		resultNode = &Textbox{}
	case "button":
		resultNode = &Button{}
	case "selectlist":
		resultNode = &SelectList{}
	case "inputfield":
		resultNode = &InputField{}
	case "filesearch":
		resultNode = &FileSearch{}
	case "slider":
		resultNode = &Slider{}
	case "dialogue":
		resultNode = &Dialogue{}
	}

	err := value.Decode(resultNode)
	if err != nil {
		return err
	}

	n.Node = resultNode
	return nil
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
