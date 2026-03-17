package engine

import (
	"errors"
	"fmt"
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/backend/node/ebitenrender"
	"reflect"
	"unsafe"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// I'm pretty sure this is massively slow right now. Performance
// improvements are in order. For instance, we probably don't need
// to do as many tree searches, and maybe we can do some things upon
// changes to the nodeTree?

// Interface representing an actor (essentially a component)
// Stage() is used for 'staging' assets, i.e. putting them into the
// assetloader to later load them. Called before initialization
// OnTreeAdd is called immediately after the node is added to the tree.
// A pointer to the Node instance corresponding to this actor is passed in
// The rest of the methods are obvious.
type Actor interface {
	Init(*Commands)
	Update(*Commands) // TODO: Change to handle errors in nodes independently
	OnTreeAdd(*Node, *Commands)
	DrawInspector(*debugui.Context)
	DrawGizmo(*Commands)
}

type Node = node.Node[Actor]
type NodeTree = node.NodeTree[Actor]

type SceneBuilder func(*Commands) *Scene

type Scene struct {
	name       string
	nodeTree   *NodeTree
	drawGizmos bool
}

func (s *Scene) Update(cmd *Commands) {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		s.drawGizmos = !s.drawGizmos
	}

	s.nodeTree.Traverse(func(n *Node) {
		n.GetValue().Update(cmd)
	})

	if s.drawGizmos {
		s.nodeTree.Traverse(func(n *Node) {
			n.GetValue().DrawGizmo(cmd)
		})
	}
}

func (s *Scene) Init(servers *Commands) {
	s.nodeTree.Traverse(func(n *Node) {
		n.GetValue().Init(servers)
	})
}

func (s *Scene) MakeDrawFunc(w, h int) func(ctx *debugui.Context) error {
	return ebitenrender.MakeRenderFunc(s.name, w, h, s.nodeTree, func(ctx *debugui.Context, nodeVal Actor) {
		nodeVal.DrawInspector(ctx)
	})
}

func (s *Scene) GetNodeByName(name string) (*Node, bool) {
	return s.nodeTree.GetNodeFunc(func(n *Node) bool {
		return n.GetName() == name
	})
}

func (s *Scene) GetNodeByID(id string) (*Node, bool) {
	return s.nodeTree.GetNode(id)
}

// Returns the first actor of type T, or false if none is found
func GetNodeByType[T Actor](s *Scene) (*Node, bool) {
	f := func(n *Node) bool {
		_, ok := n.GetValue().(T)
		return ok
	}
	return s.nodeTree.GetNodeFunc(f)
}

func (s *Scene) GetRoot() *Node {
	return s.nodeTree.GetRoot()
}

func (s *Scene) GetName() string {
	return s.name
}

func (s *Scene) SpawnActor(name string, actor Actor, cmd *Commands) *Node {
	node := s.nodeTree.GetRoot().AddChild(actor, name)
	actor.OnTreeAdd(node, cmd)
	return node
}

// Spawn a child of type node with the specified parent
// This is sort of annoying: We should be able to add children by
// chaining methods, hence the Node is the type that should have
// an addchild method. However, that doesn't really work...
func (s *Scene) AddChild(actor Actor, name string, parent *Node, cmd *Commands) *Node {
	node := s.SpawnActor(name, actor, cmd)
	s.SetParent(node, parent)
	return node
}

func (s *Scene) SetParent(node *Node, parent *Node) {
	node.Reparent(parent)
}

func (s *Scene) Print() {
	s.nodeTree.Print()
}

// TODO: Make it so that we don't have to use a pointer as type argument
// Returns the field T embedded in the actor passed in, i.e.
//
//	GetActor[*Node](transform2D)
//
// returns the Node actor embedded in the Transform2D passed in.
func GetActor[T Actor](actor Actor) (T, bool) {
	if val, ok := actor.(T); ok {
		return val, true
	}

	v := reflect.ValueOf(actor).Elem()
	t := v.Type()
	var empty T
	for i := range t.NumField() {
		tf := t.Field(i)
		if !tf.Anonymous {
			continue
		}

		vf := v.Field(i)
		val := extractFieldUnsafe(vf).Interface()
		if val, ok := val.(T); ok {
			return val, true
		}

		recurseVal, ok := GetActor[T](vf.Addr().Interface().(Actor))
		if ok {
			return recurseVal, true
		}
	}
	return empty, false
}

func extractFieldUnsafe(v reflect.Value) reflect.Value {
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

func NewScene(name string, root Actor, servers *Commands) *Scene {
	nodeTree, rootNode := node.NewNodeTree(root)
	root.OnTreeAdd(rootNode, servers)
	return &Scene{
		name:     name,
		nodeTree: nodeTree,
	}
}

var ErrTerminated = errors.New("Terminatednow")

type Game struct {
	cmd         *Commands
	scenes      map[string]SceneBuilder
	activeScene *Scene
}

func NewGame(cmd *Commands) *Game {
	return &Game{
		cmd:    cmd,
		scenes: make(map[string]SceneBuilder),
	}
}

func (g *Game) RegisterScene(name string, sceneBuilder SceneBuilder) *Game {
	g.scenes[name] = sceneBuilder
	return g
}

func (g *Game) SpawnScene(name string) *Game {
	// Create an instance of the scene's node tree
	sceneBuilder := g.scenes[name]
	sceneInst := sceneBuilder(g.cmd)

	// Load any staged assets
	g.LoadStaged() // This will go into a different thread to avoid freezing
	g.cmd.scene = sceneInst
	g.activeScene = sceneInst

	g.activeScene.Init(g.cmd)
	return g
}

func (g *Game) ActiveScene() *Scene {
	return g.activeScene
}

func (g *Game) Update() error {
	if g.activeScene == nil {
		fmt.Println("Warning: Running game without active scene")
		return nil
	}

	g.activeScene.Update(g.cmd) // consider returning this instead?
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cmd.Renderer().Draw(screen)
}

// Loads all staged assets
func (g *Game) LoadStaged() {
	g.cmd.AssetLoader().LoadAll()
}

func (g *Game) GetCmd() *Commands {
	return g.cmd
}
