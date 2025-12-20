package scene

import (
	"fmt"
)

// Note: This is slowly turning into the `node` system we are using for UI. At some point it will
// probably be possible to merge them

type SceneStack struct {
	stack []Scene
}

func (s *SceneStack) Update() bool {
	for _, scene := range s.stack {
		// TODO: Validate transition: For instance, insure that if the operation is not Pop or PopN, a
		// scene was supplied.
		transition, ok := scene.Update(s)
		if ok {
			if transition.Kind == Quit {
				return true
			}
			s.Switch(transition)
		}
	}
	return false
}

func (s *SceneStack) Draw() {
	for _, scene := range s.stack {
		scene.Draw()
	}
}

func (s *SceneStack) Switch(transition *SceneTransition) {
	switch transition.Kind {
	case Replace:
		if _, n, ok := s.GetScene(transition.Name); ok {
			s.stack = s.stack[:n]
			s.Push(transition.OtherScene)
		}
	case Push:
		s.Push(transition.OtherScene)
	case Pop:
		s.stack = s.stack[:len(s.stack)-1]
	case PopN:
		s.stack = s.stack[:len(s.stack)-transition.N]
	case PopName:
		n := len(s.stack) - 1
		for i := len(s.stack) - 1; i >= 0; i-- {
			if s.stack[i].GetName() == transition.Name {
				n = i
				break
			}
		}
		s.stack[n] = s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
	}
}

func (s *SceneStack) GetScene(name string) (Scene, int, bool) {
	for i, scene := range s.stack {
		if scene.GetName() == name {
			return scene, i, true
		}
	}
	fmt.Println("Could not find scene with name", name)
	return nil, 0, false
}

func (s *SceneStack) Push(scene Scene) {
	s.stack = append(s.stack, scene)
	scene.Init()
}

func NewSceneStack() *SceneStack {
	return &SceneStack{
		stack: make([]Scene, 0),
	}
}

type Scene interface {
	Init()
	Update(sceneStack *SceneStack) (*SceneTransition, bool)
	Draw()
	GetName() string
}

type Kind int

const (
	Replace Kind = iota
	Push
	Pop
	PopN
	PopName
	Quit
)

type SceneTransition struct {
	Kind       Kind
	OtherScene Scene
	Name       string
	N          int // Only used for PopN
}
