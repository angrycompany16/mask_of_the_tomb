package node_v2

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type Node[T any] struct {
	id       string
	name     string // Non-unique
	value    T
	parent   *Node[T]
	children []*Node[T]
}

type nodeCtx struct {
	depth int
}

func (n *Node[T]) traverseRecursive(ctx nodeCtx, callback func(*Node[T])) {
	callback(n)
	for _, child := range n.children {
		child.traverseRecursive(nodeCtx{
			depth: ctx.depth + 1,
		}, callback)
	}
}

func (n *Node[T]) GetName() string {
	return n.name
}

func (n *Node[T]) GetID() string {
	return n.id
}

func (n *Node[T]) GetValue() *T {
	return &n.value
}

func (n *Node[T]) GetParent() *Node[T] {
	return n.parent
}

func (n *Node[T]) GetChildren() []*Node[T] {
	return n.children
}

func (n *Node[T]) GetChild(id string) (*Node[T], bool) {
	return n.GetChildFunc(func(child *Node[T]) bool {
		return child.id == id
	})
}

func (n *Node[T]) getChildRecursive(id string) (*Node[T], bool) {
	return n.getChildRecursiveFunc(func(child *Node[T]) bool {
		return child.id == id
	})
}

// Returns the first node for which the function evaluates to true.
func (n *Node[T]) GetChildFunc(f func(*Node[T]) bool) (*Node[T], bool) {
	i := slices.IndexFunc(n.children, func(child *Node[T]) bool {
		return f(child)
	})
	if i == -1 {
		// fmt.Println("The node search did not return anything.")
		// n.Print()
		return nil, false
	}
	return n.children[i], true
}

func (n *Node[T]) getChildRecursiveFunc(f func(*Node[T]) bool) (*Node[T], bool) {
	child, found := n.GetChildFunc(f)
	if !found {
		for _, child := range n.children {
			return child.getChildRecursiveFunc(f)
		}
	} else {
		return child, true
	}
	return nil, false
}

func (n *Node[T]) AddChild(value T, name string) string {
	child := NewNode(value, name)
	child.parent = n
	n.children = append(n.children, child)
	return child.id
}

func (n *Node[T]) DeleteChild(id string) {
	trimmedChildren := slices.DeleteFunc(n.children, func(node *Node[T]) bool {
		return node.id == id
	})
	n.children = trimmedChildren
}

func (n *Node[T]) Reparent(parent *Node[T]) {
	trimmedChildren := slices.DeleteFunc(n.parent.children, func(node *Node[T]) bool {
		return node.id == n.id
	})
	n.parent.children = trimmedChildren
	n.parent = parent
	n.parent.children = append(n.parent.children, n)
}

func (n *Node[T]) Print() {
	fmt.Println("Node info:")
	fmt.Println("ID:", n.id)
	fmt.Println("Name:", n.name)
	fmt.Println("children:")
	for i, child := range n.children {
		fmt.Printf("	ID %d: %s\n", i, child.id)
		fmt.Printf("	Name %d: %s\n", i, child.name)
	}
}

func NewNode[T any](value T, name string) *Node[T] {
	return &Node[T]{
		id:    uuid.NewString(),
		name:  name,
		value: value,
	}
}

type NodeTree[T any] struct {
	root *Node[T]
}

func (nt *NodeTree[T]) Traverse(callback func(*Node[T])) {
	nt.root.traverseRecursive(nodeCtx{0}, callback)
}

func (nt *NodeTree[T]) GetRoot() *Node[T] {
	return nt.root
}

func (n *NodeTree[T]) GetNode(id string) (*Node[T], bool) {
	return n.GetNodeFunc(func(child *Node[T]) bool {
		return child.id == id
	})
}

// Returns the first node for which the function evaluates to true.
func (n *NodeTree[T]) GetNodeFunc(f func(*Node[T]) bool) (*Node[T], bool) {
	return n.GetRoot().getChildRecursiveFunc(f)
}

func (nt *NodeTree[T]) Print() {
	nt.Traverse(func(n *Node[T]) {
		n.Print()
	})
}

// Returns the new node tree, and the ID of the root node
func NewNodeTree[T any](rootValue T) (*NodeTree[T], string) {
	rootNode := &Node[T]{
		id:    uuid.NewString(),
		name:  "root",
		value: rootValue,
	}
	newNodeTree := NodeTree[T]{
		root: rootNode,
	}
	return &newNodeTree, rootNode.id
}

// Probably make some sort of nice builder syntax
// Also a more advanced query system would be nice
