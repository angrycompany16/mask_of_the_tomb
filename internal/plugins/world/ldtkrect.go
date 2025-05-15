package world

import (
	"mask_of_the_tomb/internal/core/maths"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

func RectFromEntity(entityInstance *ebitenLDTK.Entity) *maths.Rect {
	return maths.NewRect(
		entityInstance.Px[0],
		entityInstance.Px[1],
		entityInstance.Width,
		entityInstance.Height,
	)
}
