package key

import (
	"mask_of_the_tomb/internal/engine/actors/missile"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/commands"
)

// Here it would be nice to just create a bundle. Idea for how it could work:
// Key -- has logic, collision detection etc.
// -> Sprite
// Thats's it

type CollectedKey struct {
	*missile.Missile
	Sprite *sprite.Sprite
}

func (k *CollectedKey) Init(cmd *commands.Commands) {
	k.Missile.Init(cmd)
	k.SetPos(k.Missile.TargetTransform.GetPos(false))
}

func (k *CollectedKey) Update(cmd *commands.Commands) {
	k.Missile.Update(cmd)

	// Detect when we try to unlock something. Find the door and move to the lock point.
	playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
	if playerControls.PollAction("Use") {
		k.Active = false
		k.Sprite.Hidden = true
	}
}

func NewCollectedKey(missile *missile.Missile) *CollectedKey {
	return &CollectedKey{
		Missile: missile,
	}
}
