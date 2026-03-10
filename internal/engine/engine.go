package engine

import (
	"mask_of_the_tomb/internal/engine/servers"
	"mask_of_the_tomb/internal/engine/servers/renderer"
	"mask_of_the_tomb/internal/node_v2"
	"mask_of_the_tomb/internal/node_v2/ebitenrender"
	"reflect"
	"unsafe"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
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
	Init()
	Update(*servers.Servers)
	// We need to implement LateUpdate
	OnTreeAdd(*Node, *servers.Servers)
	DrawInspector(ctx *debugui.Context)
}

type Node = node_v2.Node[Actor]
type NodeTree = node_v2.NodeTree[Actor]

type Scene struct {
	name     string
	nodeTree NodeTree
	servers  *servers.Servers
}

func (s *Scene) Update(servers *servers.Servers) {
	s.nodeTree.Traverse(func(n *Node) {
		(*n.GetValue()).Update(servers)
	})
}

func (s *Scene) Init() {
	s.nodeTree.Traverse(func(n *Node) {
		(*n.GetValue()).Init()
	})
}

func (s *Scene) MakeDrawFunc(w, h int) func(ctx *debugui.Context) error {
	return ebitenrender.MakeRenderFunc(s.name, w, h, &s.nodeTree, func(ctx *debugui.Context, nodeVal *Actor) {
		(*nodeVal).DrawInspector(ctx)
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

func (s *Scene) GetRoot() *Node {
	return s.nodeTree.GetRoot()
}

func (s *Scene) GetName() string {
	return s.name
}

func (s *Scene) Spawn(name string, actor Actor) *Node {
	node := s.nodeTree.GetRoot().AddChild(actor, name)
	actor.OnTreeAdd(node, s.servers)
	return node
}

// Spawn a child of type node with the specified parent
// This is sort of annoying: We should be able to add children by
// chaining methods, hence the Node is the type that should have
// an addchild method. However, that doesn't really work...
func (s *Scene) AddChild(actor Actor, name string, parent *Node) *Node {
	node := s.Spawn(name, actor)
	s.SetParent(node, parent)
	return node
}

func (s *Scene) SetParent(node *Node, parent *Node) {
	node.Reparent(parent)
}

func (s *Scene) Print() {
	s.nodeTree.Print()
}

// Returns the field T embedded in the actor passed in, i.e.
//
//	GetActor[Node](transform2D)
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

func NewScene(name string, root Actor, servers *servers.Servers) *Scene {
	nodeTree, rootNode := node_v2.NewNodeTree(root)
	root.OnTreeAdd(rootNode, servers)
	return &Scene{
		servers:  servers,
		name:     name,
		nodeTree: *nodeTree,
	}
}

type Game struct {
	servers     *servers.Servers
	editor      *Editor
	scenes      map[string]*Scene // Would it be more correct to call these scene templates?
	activeScene Scene
}

func NewGame(servers *servers.Servers, editor *Editor) *Game {
	return &Game{
		servers: servers,
		editor:  editor,
		scenes:  make(map[string]*Scene),
	}
}

func (g *Game) AddScene(scene *Scene) *Game {
	g.scenes[scene.name] = scene
	return g
}

func (g *Game) MakeScene(name string, root Actor) *Scene {
	g.scenes[name] = NewScene(name, root, g.servers)
	return g.scenes[name]
}

func (g *Game) SetActiveScene(name string) *Game {
	g.LoadStaged() // This will go into a different thread to avoid freezing
	g.activeScene = *g.scenes[name]
	g.activeScene.Init() // Reset values?
	return g
}

func (g *Game) ActiveScene() *Scene {
	return &g.activeScene
}

func (g *Game) Update() error {
	if _, err := g.editor.treeUI.Update(
		g.ActiveScene().MakeDrawFunc(int(300), int(500)),
	); err != nil {
		return err
	}

	g.activeScene.Update(g.servers) // consider returning this instead?
	g.editor.Update(g.servers)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.editor.Draw()
	g.servers.Renderer().Draw(screen)
}

// Loads all staged assets
func (g *Game) LoadStaged() {
	g.servers.AssetLoader().LoadAll()
}

type Editor struct {
	treeUI      debugui.DebugUI
	editorImage *ebiten.Image
}

func (e *Editor) Update(servers *servers.Servers) {
	servers.Renderer().Request(renderer.Pos(e.editorImage, 0, 0), e.editorImage, "EditorUI", 0)
}

func (e *Editor) Draw() {
	e.editorImage.Clear()
	e.treeUI.Draw(e.editorImage)
}

func NewEditor(w, h int) *Editor {
	return &Editor{
		editorImage: ebiten.NewImage(w, h),
	}
}
