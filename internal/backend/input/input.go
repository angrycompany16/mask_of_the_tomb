package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputAction func() bool

type InputServer struct {
	inputMap map[string]InputAction
}

func (i *InputServer) RegisterAction(name string, action InputAction) {
	i.inputMap[name] = action
}

func (i *InputServer) PollAction(name string) bool {
	if isAction, ok := i.inputMap[name]; ok {
		return isAction()
	} else {
		fmt.Printf("Action [%s] not found\n", name)
		return false
	}
}

func KeyJustPressedAction(key ebiten.Key) InputAction {
	return func() bool {
		return inpututil.IsKeyJustPressed(key)
	}
}

func NewInputServer() *InputServer {
	return &InputServer{
		inputMap: make(map[string]InputAction),
	}
}
