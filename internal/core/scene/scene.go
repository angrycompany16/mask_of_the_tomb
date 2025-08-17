package scene

// Note: This is slowly turning into the `node` system we are using for UI. At some point it will
// probably be possible to merge them

// Contians methods for appending and removing scenes
type SceneTree struct {
	root Scene
}

func (st *SceneTree) Update() {
	st.root.Update(st)
}

func (st *SceneTree) Draw() {
	st.root.Draw()
}

func (st *SceneTree) Switch(scene *Scene, transition SceneTransition) {
	switch transition.Kind {
	case Replace:
		// Add a child to parent
		// delete self
	case Sibling:
		// Add a child to parent
	case Child:
		// Add a child to self
	}
}

type Scene struct {
	SceneBehaviour
	parent     *Scene
	children   []*Scene
	name       string
	transition SceneTransition
}

func (s *Scene) Update(sceneTree *SceneTree) {
	s.SceneBehaviour.Update()
	for _, child := range s.children {
		child.Update(sceneTree)
	}

	if s.Exit() {
		sceneTree.Switch(s, s.transition)
	}
}

func (s *Scene) Draw() {
	s.SceneBehaviour.Draw()
	for _, child := range s.children {
		child.Draw()
	}
}

type SceneBehaviour interface {
	Init()
	Update()
	Draw()
	Exit() bool
}

type Kind int

const (
	Replace Kind = iota
	Sibling
	Child
)

type SceneTransition struct {
	Kind Kind
	Next *Scene
}
