package switcher

import (
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
)

type Switch struct {
	*nodeactor.Node
	sceneName string
}

func (s *Switch) Init(cmd *commands.Commands) {
	s.Node.Init(cmd)
	defaultScheme := input.NewInputScheme()
	cmd.InputHandler.InputSchemes["Default"] = defaultScheme
	defaultScheme.RegisterAction("Switch", input.KeyJustPressedAction(ebiten.KeyH))
}

func (s *Switch) Update(cmd *commands.Commands) {
	s.Node.Update(cmd)
	defaultScheme := cmd.InputHandler.InputSchemes["Default"]
	if defaultScheme.PollAction("Switch") {
		sceneManager, _ := commands.Get[engine.SceneManager](cmd)
		if s.sceneName == "TestScene1" {
			sceneManager.SpawnScene("TestScene2", cmd)
		} else if s.sceneName == "TestScene2" {
			sceneManager.SpawnScene("TestScene1", cmd)
		}
	}
}

func NewSwitch(node *nodeactor.Node, sceneName string) *Switch {
	return &Switch{
		Node:      node,
		sceneName: sceneName,
	}
}
