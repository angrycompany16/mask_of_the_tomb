package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ActionBinding func() bool

type InputServer struct {
	inputMap map[string][]ActionBinding
}

func (i *InputServer) RegisterAction(name string, binding ActionBinding) {
	i.inputMap[name] = []ActionBinding{binding}
}

func (i *InputServer) AddBinding(name string, binding ActionBinding) {
	i.inputMap[name] = append(i.inputMap[name], binding)
}

func (i *InputServer) PollAction(name string) bool {
	if actionBindings, ok := i.inputMap[name]; ok {
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

func NewInputServer() *InputServer {
	return &InputServer{
		inputMap: make(map[string][]ActionBinding),
	}
}
