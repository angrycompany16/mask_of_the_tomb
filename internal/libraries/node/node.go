package node

import (
	"errors"
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/sound"
	"slices"

	"gopkg.in/yaml.v3"
)

type Node interface {
	Draw(offsetX, offsetY float64, parentWidth, parentHeight float64)
	Update(confirmations map[string]ConfirmInfo)
	Reset(overWriteInfo map[string]OverWriteInfo)

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
	case "container":
		resultNode = &Container{}
	case "textbox":
		resultNode = &Textbox{}
	case "button":
		// Looks like loading a stream into multiple players may not be very supported
		selectSoundStream := errs.Must(assettypes.GetOggStream("selectSound"))
		selectEffectPlayer := &sound.EffectPlayer{errs.Must(sound.FromStream(selectSoundStream)), 1.0}
		fmt.Println("Finished creating selectEffectPlayer")
		selectEffectPlayer.SetVolume(0.5)
		fmt.Println("Success for effect player", selectEffectPlayer)
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
		dialogueSoundStream := errs.Must(assettypes.GetOggStream("dialogueSound"))
		dialogueEffectPlayer := &sound.EffectPlayer{errs.Must(sound.FromStream(dialogueSoundStream)), 1.0}
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

func (n *NodeData) UpdateChildren(confirmations map[string]ConfirmInfo) {
	fmt.Println("node", n)
	fmt.Println("children:", n.Children)
	for _, child := range n.Children {
		fmt.Println("updating child", child.Node)
		child.Update(confirmations)
		fmt.Println("what")
		fmt.Println("Confirmations", confirmations)
		fmt.Println("Updated child", child)
	}
	fmt.Println("Exited loop")
}

func (n *NodeData) DrawChildren(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	for _, child := range n.Children {
		child.Draw(offsetX, offsetY, parentWidth, parentHeight)
	}
}

func (n *NodeData) ResetChildren(overWriteInfo map[string]OverWriteInfo) {
	for _, child := range n.Children {
		child.Reset(overWriteInfo)
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
