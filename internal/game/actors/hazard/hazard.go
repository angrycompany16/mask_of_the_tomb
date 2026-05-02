package hazard

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/game/actors/trigger"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

type Hazard struct {
	*trigger.Trigger
}

func NewHazard(graphic *graphic.Graphic, entity *ebitenLDTK.Entity) *Hazard {
	hitbox := maths.NewRect(
		0,
		0,
		entity.Width,
		entity.Height,
	)

	hazard := Hazard{
		Trigger: trigger.NewTrigger(graphic,
			trigger.WithName("Hazard"),
			trigger.WithRect(hitbox),
		),
	}

	hazard.SetPos(entity.Px[0], entity.Px[1])

	return &hazard
}
