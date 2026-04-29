package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ActionBinding func() bool

type InputHandler struct {
	InputSchemes map[string]*InputScheme
}

type InputScheme struct {
	Active   bool
	inputMap map[string][]ActionBinding
}

func (i *InputScheme) RegisterAction(name string, binding ActionBinding) {
	i.inputMap[name] = []ActionBinding{binding}
}

func (i *InputScheme) AddBinding(name string, binding ActionBinding) {
	i.inputMap[name] = append(i.inputMap[name], binding)
}

func (i *InputScheme) PollAction(name string) bool {
	if actionBindings, ok := i.inputMap[name]; ok {
		if !i.Active {
			return false
		}

		for _, actionBinding := range actionBindings {
			if actionBinding() {
				return true
			}
		}
	} else {
		fmt.Printf("Action [%s] not found\n", name)
		return false
	}
	return false
}

func KeyJustPressedAction(key ebiten.Key) ActionBinding {
	return func() bool {
		return inpututil.IsKeyJustPressed(key)
	}
}

func NewInputHandler() *InputHandler {
	return &InputHandler{
		InputSchemes: make(map[string]*InputScheme),
	}
}

func NewInputScheme() *InputScheme {
	return &InputScheme{
		Active:   true,
		inputMap: make(map[string][]ActionBinding),
	}
}
