package slambox

import "mask_of_the_tomb/internal/core/maths"

type ChainNodeConnection struct {
	node *ChainNode
	dir  maths.Direction
}

type ChainNode struct {
	connections []*ChainNodeConnection
}

// A chain that can connect multiple slamboxes / slambox groups.
type SlamboxChain struct {
	nodes         []*ChainNode
	slamboxes     []*Slambox
	slamboxGroups []*SlamboxGroup
}
