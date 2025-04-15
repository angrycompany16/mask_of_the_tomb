package node

import (
	"errors"
	"slices"

	"gopkg.in/yaml.v3"
)

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
	}

	err := value.Decode(resultNode)
	if err != nil {
		return err
	}

	n.Node = resultNode
	return nil
}
