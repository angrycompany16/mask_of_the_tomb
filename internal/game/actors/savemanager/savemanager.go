package savemanager

import (
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
)

// Its only purpose is to call load and save lol
type SaveManager struct {
	*nodeactor.Node
}
