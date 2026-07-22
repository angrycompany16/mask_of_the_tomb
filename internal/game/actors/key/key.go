package key

import (
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/utils"
)

// Here it would be nice to just create a bundle. Idea for how it could work:
// Key -- has logic, collision detection etc.
// -> Sprite
// Thats's it

type Key struct {
	*transform2D.Transform2D
}

func defaultKey(transform2D *transform2D.Transform2D) *Key {
	return &Key{
		Transform2D: transform2D,
	}
}

func NewKey(transform2D *transform2D.Transform2D, options ...utils.Option[Key]) *Key {
	key := defaultKey(transform2D)

	for _, option := range options {
		option(key)
	}

	return key
}
