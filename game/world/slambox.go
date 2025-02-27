package world

import (
	"mask_of_the_tomb/ebitenLDTK"
	"mask_of_the_tomb/game/physics"
	"mask_of_the_tomb/utils/rect"
)

type slambox struct {
	collider physics.RectCollider
}

func newSlambox(
	entity *ebitenLDTK.Entity,
) slambox {
	newSlambox := slambox{}
	newSlambox.collider = physics.NewRectCollider(*rect.FromEntity(entity))

	return newSlambox
}
