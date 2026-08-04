package grok

import "mask_of_the_tomb/internal/engine/actors/animatedsprite"

// A note:
// Although we do need to (should) use LDtk for the first spawn position, from then the spawning of an entity is going to be
// so dynamic that we probably will just have to deal with it using persistent GlobalData storage.

// It's possible that its better to spawn NPCs as a bundle but we shall see
// Q: What components are needed for
type Grok struct {
	*animatedsprite.AnimatedSprite
}
