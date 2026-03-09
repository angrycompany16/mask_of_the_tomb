package engine

import (
	"fmt"
	"mask_of_the_tomb/internal/node_v2"
	"mask_of_the_tomb/internal/node_v2/ebitenrender"
	"reflect"

	"github.com/ebitengine/debugui"
)

// I'm pretty sure this is massively slow right now. Performance
// improvements are in order. For instance, we probably don't need
// to do as many tree searches, and maybe we can do some things upon
// changes to the nodeTree?

// Steal Godot's server architecture
// Notice that with the new audio system we've already started on a
// client-server system

// We may possibly add some arguments to these later
type Actor interface {
	Init()
	Update(*NodeTree)
	SetID(string)
	DrawInspector(ctx *debugui.Context)
	// We don't really need an explicit draw method i think?

	// No, we won't - Components will send commands to the rendering server during
	// Update(), which will then be executed, and then ebitengine will draw
	// whatever the server spits out.
}

type Node = node_v2.Node[Actor]
type NodeTree = node_v2.NodeTree[Actor]

// A scene then simply consists of a nodetree of INode's
// More or less contains wrapper functions for the node tree
type Scene struct {
	name     string
	nodeTree *NodeTree
}

func (s *Scene) Update() {
	s.nodeTree.Traverse(func(n *Node) {
		(*n.GetValue()).Update(s.nodeTree)
	})
}

func (s *Scene) Init() {
	s.nodeTree.Traverse(func(n *Node) {
		(*n.GetValue()).Init()
	})
}

func (s *Scene) MakeDrawFunc() func(ctx *debugui.Context) error {
	return ebitenrender.MakeRenderFunc(s.name, s.nodeTree, func(ctx *debugui.Context, nodeVal *Actor) {
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

// Spawn a child of type node with the specified parent
func (s *Scene) Spawn(name string, actor Actor) string {
	treeID := s.nodeTree.GetRoot().AddChild(actor, name)
	actor.SetID(treeID)
	return treeID
}

func (s *Scene) AddChild(name string, actor Actor, parent *Node) string {
	treeID := parent.AddChild(actor, name)
	actor.SetID(treeID)
	return treeID
}

func (s *Scene) SetParent(node *Node, parent *Node) {
	node.Reparent(parent)
}

func (s *Scene) Print() {
	s.nodeTree.Print()
}

// Returns the field T embedded in the actor passed in, i.e.
//
//	GetActor[nodes.Null](transform2D)
//
// returns the Null actor embedded in the transform2D passed in.
func GetActor[T Actor](actor Actor) (T, bool) {
	if val, ok := actor.(T); ok {
		return val, true
	}
	// Loop through embedded fields
	// oh boy, here we go again...
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

		fmt.Println(vf)
		recurseVal, ok := GetActor[T](vf.Addr().Interface().(Actor))
		if ok {
			return recurseVal, true
		}
	}
	return empty, false
}

// TODO: We'll need some builder methods for scene structs for QoL reasons
func NewScene(name string, root Actor) *Scene {
	nodeTree, rootID := node_v2.NewNodeTree(root)
	root.SetID(rootID)
	return &Scene{
		name:     name,
		nodeTree: nodeTree,
	}
}

// Idea for nodes/components/assets
// Make all components loadable from .yaml files, so that these can be used to
// store scenes and such things on disk instead of keeping all of it in memory.
// These become prefabs
// Of course nodes can also be instantiated via code.
// This is nice because it means that we get a lot of flexibility: For instance
// it is probably a good idea to implement some basic editor functionality in the
// debug ui, where we can make changes live and then record them as new scene files.
// That's actually hella nice
// We embed everything (probably just keep a FS), and then we create a function for
// adding a prefab to the loading stage, as well as a function for loading all
// prefabs into actual objects in memory (this runs in a separate thread)
// This is probably most easily achieved with a simple mutex?
// we could potentially try to make it more idiomatic by creating an asset server
// that responds to requests via channels, but this raises a large concern: What
// do we show when waiting for assets? For instance, this would require a
// 'missing image' file to show while waiting for an image to load, which is not
// really that great.
// Better to go the simple way I think.
// So to summarize:
//   We load prefabs (whic hare subtrees) from asset files stored in YAML format.
//   A scene is then really nothing but a root node plus some children.
// Now in this loading stage we also load any references to other files that are
// stored internally in the nodes. For instance, a sprite node will have a file path
// that tells us where the sprite is located. Here we should actuallty create a
// 'missing image' file so that we can tell when an asset path is wrong or an asset
// is broken (this should not crash the game!)
// Well, that's it. The core philosophy should be:
// - Any component must be readable from YAML
// - We must also be able to visualize it
// - I guess that's it?

// One more thing: How do we make components interact with each other?
// For user-written components it will suffice to use events for notifying the
// other components
// But for internal components, this could lead to some bugs
// For instance, a physics component must control the transform component
// A sprite component must read the transform component to find global position

// And this has to happen automatically. A physics component without transform cannot
// be allowed, an neither (?) can a sprite component

// An important thing to notice:
// - Godot's nodes work through inheritance, so each node is more than just a simple
//   behaviour
// - For instance, the physicsbody2D node will control its transform by simply
//   inheriting from it, and then controlling it in its update method.
// - But then a final question: How do we propagate the transforms?
//   I guess it'll be something like this? The transform component has a behaviour
//   which gets the parent node (unless it is root) and adds the parent transform
//   to itself, which is called during the update method.
//   More specifically: Get the parent node, type assert that it is a transform
//   (or inheriting from a transform, this may be tricky[it was indeed tricky]), get the transform values
//   apply them to our guy. If the parent has no transform (which will be very rare,
//   but should be handled), just set the global position equal to the local position.
//   We'll need a linear algebra library to represent the transform as a matrix and
//   speed things up a lil bit eventually...
//   or... just write it ourselves hehe
